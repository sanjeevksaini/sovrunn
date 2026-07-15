package server

import (
	"context"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sanjeevksaini/sovrunn/internal/api"
	"github.com/sanjeevksaini/sovrunn/internal/config"
	"github.com/sanjeevksaini/sovrunn/internal/health"
	"github.com/sanjeevksaini/sovrunn/internal/resources"
)

// Server owns HTTP server lifecycle.
type Server struct {
	cfg        config.Config
	httpServer *http.Server
	readiness  *health.ReadinessState
	logger     *log.Logger
}

// New constructs the Server with routes and middleware registered.
func New(
	cfg config.Config,
	org *api.OrgHandler,
	bootstrap *api.BootstrapHandler,
	readiness *health.ReadinessState,
) *Server {
	mux := http.NewServeMux()
	logger := log.New(os.Stdout, "", log.LstdFlags)

	chain := func(h http.Handler) http.Handler {
		return requestIDMiddleware(loggingMiddleware(logger)(contentTypeMiddleware(h)))
	}
	bootstrapChain := func(h http.Handler) http.Handler {
		return requestIDMiddleware(loggingMiddleware(logger)(h))
	}

	mux.Handle("/v1/organizations", chain(http.HandlerFunc(org.HandleCollection)))
	mux.Handle("/v1/organizations/", chain(http.HandlerFunc(org.HandleItem)))

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		bootstrapChain(http.HandlerFunc(bootstrap.Healthz)).ServeHTTP(w, r)
	})
	mux.HandleFunc("/readyz", func(w http.ResponseWriter, r *http.Request) {
		bootstrapChain(http.HandlerFunc(bootstrap.Readyz)).ServeHTTP(w, r)
	})
	mux.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		bootstrapChain(http.HandlerFunc(bootstrap.Version)).ServeHTTP(w, r)
	})

	return &Server{
		cfg: cfg,
		httpServer: &http.Server{
			Addr:    cfg.Addr(),
			Handler: mux,
		},
		readiness: readiness,
		logger:    logger,
	}
}

// Start binds the listener, marks readiness true, and blocks until signal.
func (s *Server) Start() error {
	listener, err := net.Listen("tcp", s.httpServer.Addr)
	if err != nil {
		return err
	}

	s.readiness.SetReady(true)

	errCh := make(chan error, 1)
	go func() {
		s.logger.Printf("server listening on %s", s.cfg.Addr())
		if err := s.httpServer.Serve(listener); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	sigCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	select {
	case <-sigCtx.Done():
		timeout := time.Duration(s.cfg.Server.ShutdownTimeout) * time.Second
		if err := s.Shutdown(timeout); err != nil {
			return err
		}
		s.logger.Println("server shutdown complete")
		return nil
	case err := <-errCh:
		s.readiness.SetReady(false)
		return err
	}
}

// Shutdown stops accepting new connections and drains in-flight requests.
func (s *Server) Shutdown(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	s.readiness.SetReady(false)
	return s.httpServer.Shutdown(ctx)
}

func writeErrorBody(w http.ResponseWriter, code resources.ErrorCode, message, field, details string) error {
	envelope := resources.APIErrorEnvelope{
		Error: resources.APIError{
			Code:    code,
			Message: message,
			Field:   field,
			Details: details,
		},
	}
	return json.NewEncoder(w).Encode(envelope)
}
