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
	"github.com/sanjeevksaini/sovrunn/internal/cloudmodel"
	"github.com/sanjeevksaini/sovrunn/internal/config"
	"github.com/sanjeevksaini/sovrunn/internal/health"
	"github.com/sanjeevksaini/sovrunn/internal/resources"
)

// CloudModelRuntime is the FEATURE-0015 composition attached to the API server:
// store, idempotency coordinator, audit-publication coordinator, and the
// audit-aware participation expiry scheduler (Task 8).
type CloudModelRuntime struct {
	Store       *cloudmodel.Store
	Idempotency *cloudmodel.IdempotencyCoordinator
	Publication *cloudmodel.PublicationCoordinator
	Scheduler   *cloudmodel.ParticipationExpiryScheduler
	Guard       *cloudmodel.ParticipationActionGuard
}

// Server owns HTTP server lifecycle and optional FEATURE-0015 cloud-model
// coordination (scheduler + idempotency abort on shutdown).
type Server struct {
	cfg        config.Config
	httpServer *http.Server
	readiness  *health.ReadinessState
	logger     *log.Logger

	cloud          *CloudModelRuntime
	schedCtxCancel context.CancelFunc
}

// New constructs the Server with routes and middleware registered.
// ServiceClass, ServicePlan, Plugin, Capability, ServiceInstance, and
// ServiceBinding routes are registered only when the corresponding handlers
// are non-nil.
func New(
	cfg config.Config,
	org *api.OrgHandler,
	ou *api.OUHandler,
	tenant *api.TenantHandler,
	project *api.ProjectHandler,
	operation *api.OperationHandler,
	serviceClass *api.ServiceClassHandler,
	servicePlan *api.ServicePlanHandler,
	plugin *api.PluginHandler,
	capability *api.CapabilityHandler,
	serviceInstance *api.ServiceInstanceHandler,
	serviceBinding *api.ServiceBindingHandler,
	bootstrap *api.BootstrapHandler,
	readiness *health.ReadinessState,
) *Server {
	mux := http.NewServeMux()
	logger := log.New(os.Stdout, "", log.LstdFlags)

	chain := func(h http.Handler) http.Handler {
		return requestIDMiddleware(loggingMiddleware(logger)(contentTypeMiddleware(h)))
	}
	bootstrapChain := func(h http.Handler) http.Handler {
		return requestIDMiddleware(loggingMiddleware(logger)(methodGET(h)))
	}

	mux.Handle("/v1/organizations", chain(http.HandlerFunc(org.HandleCollection)))
	mux.Handle("/v1/organizations/", chain(http.HandlerFunc(org.HandleItem)))

	mux.Handle("/v1/organization-units", chain(http.HandlerFunc(ou.HandleCollection)))
	mux.Handle("/v1/organization-units/", chain(http.HandlerFunc(ou.HandleItem)))

	mux.Handle("/v1/tenants", chain(http.HandlerFunc(tenant.HandleCollection)))
	mux.Handle("/v1/tenants/", chain(http.HandlerFunc(tenant.HandleItem)))

	mux.Handle("/v1/projects", chain(http.HandlerFunc(project.HandleCollection)))
	mux.Handle("/v1/projects/", chain(http.HandlerFunc(project.HandleItem)))

	mux.Handle("/v1/operations", chain(http.HandlerFunc(operation.HandleCollection)))
	mux.Handle("/v1/operations/", chain(http.HandlerFunc(operation.HandleItem)))

	if serviceClass != nil {
		mux.Handle("/v1/service-classes", chain(http.HandlerFunc(serviceClass.HandleCollection)))
		mux.Handle("/v1/service-classes/", chain(http.HandlerFunc(serviceClass.HandleItem)))
	}
	if servicePlan != nil {
		mux.Handle("/v1/service-plans", chain(http.HandlerFunc(servicePlan.HandleCollection)))
		mux.Handle("/v1/service-plans/", chain(http.HandlerFunc(servicePlan.HandleItem)))
	}
	if plugin != nil {
		mux.Handle("/v1/plugins", chain(http.HandlerFunc(plugin.HandleCollection)))
		mux.Handle("/v1/plugins/", chain(http.HandlerFunc(plugin.HandleItem)))
	}
	if capability != nil {
		mux.Handle("/v1/capabilities", chain(http.HandlerFunc(capability.HandleCollection)))
		mux.Handle("/v1/capabilities/", chain(http.HandlerFunc(capability.HandleItem)))
	}
	if serviceInstance != nil {
		mux.Handle("/v1/service-instances", chain(http.HandlerFunc(serviceInstance.HandleCollection)))
		mux.Handle("/v1/service-instances/", chain(http.HandlerFunc(serviceInstance.HandleItem)))
	}
	if serviceBinding != nil {
		mux.Handle("/v1/service-bindings", chain(http.HandlerFunc(serviceBinding.HandleCollection)))
		mux.Handle("/v1/service-bindings/", chain(http.HandlerFunc(serviceBinding.HandleItem)))
	}

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
			Addr:              cfg.Addr(),
			Handler:           mux,
			ReadHeaderTimeout: 5 * time.Second,
			ReadTimeout:       15 * time.Second,
			WriteTimeout:      30 * time.Second,
			IdleTimeout:       60 * time.Second,
		},
		readiness: readiness,
		logger:    logger,
	}
}

// AttachCloudModel attaches FEATURE-0015 runtime coordination. The scheduler
// must already be constructed with PublicationCoordinator.ExpireParticipation
// as its injected expiry operation.
func (s *Server) AttachCloudModel(rt *CloudModelRuntime) {
	s.cloud = rt
}

// CloudModel returns the attached FEATURE-0015 runtime, if any.
func (s *Server) CloudModel() *CloudModelRuntime {
	return s.cloud
}

// NewCloudModelRuntime builds the composed store, idempotency coordinator,
// audit-publication coordinator, and audit-aware expiry scheduler.
func NewCloudModelRuntime(audit cloudmodel.AuditAppender) *CloudModelRuntime {
	store := cloudmodel.NewStore()
	guard := &cloudmodel.ParticipationActionGuard{}
	idemp := cloudmodel.NewIdempotencyCoordinator(cloudmodel.IdempotencyCoordinatorConfig{})
	if audit == nil {
		audit = &cloudmodel.MemoryAuditAppender{}
	}
	pub := cloudmodel.NewPublicationCoordinator(cloudmodel.PublicationCoordinatorConfig{
		Store:       store,
		Audit:       audit,
		Idempotency: idemp,
	})
	sched := cloudmodel.NewParticipationExpiryScheduler(cloudmodel.ParticipationExpirySchedulerConfig{
		Store:  store,
		Guard:  guard,
		Expire: pub.ExpireParticipation,
	})
	return &CloudModelRuntime{
		Store:       store,
		Idempotency: idemp,
		Publication: pub,
		Scheduler:   sched,
		Guard:       guard,
	}
}

// Start binds the listener, starts the FEATURE-0015 scheduler when attached,
// marks readiness true, and blocks until signal.
func (s *Server) Start() error {
	listener, err := net.Listen("tcp", s.httpServer.Addr)
	if err != nil {
		return err
	}

	s.startCloudModel()
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
		_ = s.stopCloudModel()
		return err
	}
}

// Shutdown stops the FEATURE-0015 scheduler, aborts in-flight idempotency
// reservations (waking waiters), then drains HTTP connections and completes.
func (s *Server) Shutdown(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	s.readiness.SetShuttingDown()
	_ = s.stopCloudModel()
	return s.httpServer.Shutdown(ctx)
}

func (s *Server) startCloudModel() {
	if s == nil || s.cloud == nil || s.cloud.Scheduler == nil {
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	s.schedCtxCancel = cancel
	s.cloud.Scheduler.Start(ctx)
	s.logger.Println("cloudmodel participation expiry scheduler started")
}

func (s *Server) stopCloudModel() error {
	if s == nil || s.cloud == nil {
		return nil
	}
	// Ordered shutdown: scheduler stop, then idempotency abort/waiter wake-up.
	if s.schedCtxCancel != nil {
		s.schedCtxCancel()
		s.schedCtxCancel = nil
	}
	if s.cloud.Scheduler != nil {
		s.cloud.Scheduler.Stop()
		s.logger.Println("cloudmodel participation expiry scheduler stopped")
	}
	if s.cloud.Idempotency != nil {
		s.cloud.Idempotency.AbortAll()
		s.logger.Println("cloudmodel in-flight idempotency reservations aborted")
	}
	return nil
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
