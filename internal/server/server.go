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
	"github.com/sanjeevksaini/sovrunn/internal/executiontarget"
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

// ExecutionTargetRuntime is the FEATURE-0016 composition: one shared lifecycle
// service and scheduler injected into every F0016 handler (TASK-F16-11).
type ExecutionTargetRuntime struct {
	Lifecycle *executiontarget.ExecutionTargetLifecycleService
	Scheduler *executiontarget.Scheduler
	Audit     executiontarget.AuditAppender
	Fixtures  *executiontarget.MapFixtureSource
}

// Server owns HTTP server lifecycle and optional FEATURE-0015 / FEATURE-0016
// coordination (scheduler + idempotency/lifecycle abort on shutdown).
type Server struct {
	cfg        config.Config
	httpServer *http.Server
	readiness  *health.ReadinessState
	logger     *log.Logger

	cloud          *CloudModelRuntime
	execution      *ExecutionTargetRuntime
	schedCtxCancel context.CancelFunc
	etSchedCancel  context.CancelFunc
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

// NewExecutionTargetRuntime constructs one shared lifecycle service and
// scheduler for FEATURE-0016 (TASK-F16-11).
func NewExecutionTargetRuntime(audit executiontarget.AuditAppender) *ExecutionTargetRuntime {
	if audit == nil {
		audit = &cloudmodel.MemoryAuditAppender{}
	}
	fixtures := executiontarget.NewMapFixtureSource()
	lifecycle := executiontarget.NewExecutionTargetLifecycleService(executiontarget.LifecycleConfig{
		Audit:    audit,
		Observer: executiontarget.NewSyntheticObserver(fixtures, nil),
	})
	sched := executiontarget.NewScheduler(executiontarget.SchedulerConfig{
		Lifecycle: lifecycle,
	})
	return &ExecutionTargetRuntime{
		Lifecycle: lifecycle,
		Scheduler: sched,
		Audit:     audit,
		Fixtures:  fixtures,
	}
}

// AttachExecutionTarget installs the completed five-registration F0016 mux
// once behind ExecutionTargetTransportGuard, injecting the shared lifecycle
// service into every F0016 handler (TASK-F16-11). cloud must be the
// FEATURE-0015 store used for backing resolution.
func (s *Server) AttachExecutionTarget(rt *ExecutionTargetRuntime, cloud *cloudmodel.Store, grants api.GrantResolver) {
	if s == nil || rt == nil || rt.Lifecycle == nil {
		return
	}
	s.execution = rt

	collection := api.NewExecutionTargetCollectionHandler(rt.Lifecycle, cloud, grants, rt.Audit)
	item := api.NewExecutionTargetItemHandler(rt.Lifecycle, cloud, grants, rt.Audit)
	qualify := api.NewExecutionTargetQualifyHandler(rt.Lifecycle, cloud, grants, rt.Audit)
	qualify.Fixtures = rt.Fixtures
	retire := api.NewExecutionTargetRetireHandler(rt.Lifecycle, cloud, grants, rt.Audit)

	f16mux := NewExecutionTargetMux(&ExecutionTargetHandlers{
		Collection: collection,
		Item:       item,
		Qualify:    qualify,
		Retire:     retire,
	})

	// Correlate F0016 requests like FEATURE-0015 route wrappers.
	f16handler := requestIDMiddleware(loggingMiddleware(s.logger)(f16mux))

	prev := s.httpServer.Handler
	combined := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, matched, _ := classifyExecutionTargetPath(r.URL.Path)
		if matched {
			f16handler.ServeHTTP(w, r)
			return
		}
		prev.ServeHTTP(w, r)
	})
	// Install the TASK-F16-08 guard exactly once around the completed F0016
	// surface (and non-F0016 pass-through via combined).
	s.httpServer.Handler = ExecutionTargetTransportGuard(combined)
}

// ExecutionTarget returns the attached FEATURE-0016 runtime, if any.
func (s *Server) ExecutionTarget() *ExecutionTargetRuntime {
	return s.execution
}

// Start binds the listener, starts FEATURE-0015 and FEATURE-0016 schedulers
// when attached, marks readiness true, and blocks until signal.
func (s *Server) Start() error {
	listener, err := net.Listen("tcp", s.httpServer.Addr)
	if err != nil {
		return err
	}

	s.startCloudModel()
	s.startExecutionTarget()
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
		_ = s.stopExecutionTarget()
		_ = s.stopCloudModel()
		return err
	}
}

// Shutdown stops FEATURE-0016 acceptance and lifecycle admission first, then
// FEATURE-0015 coordination, then drains HTTP connections (design §6.5;
// ADH-2026-058 clause 8).
func (s *Server) Shutdown(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	s.readiness.SetShuttingDown()
	_ = s.stopExecutionTarget()
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

func (s *Server) startExecutionTarget() {
	if s == nil || s.execution == nil || s.execution.Scheduler == nil {
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	s.etSchedCancel = cancel
	s.execution.Scheduler.Start(ctx)
	s.logger.Println("executiontarget expiry and maintenance-trigger intake started")
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

func (s *Server) stopExecutionTarget() error {
	if s == nil || s.execution == nil {
		return nil
	}
	// Ordered shutdown (TASK-F16-06 / TASK-F16-11): stop expiry and
	// Maintenance-trigger acceptance → abort lifecycle/idempotency admission
	// under the sole-committer mutex → wake waiters after unlock → HTTP shutdown.
	if s.etSchedCancel != nil {
		s.etSchedCancel()
		s.etSchedCancel = nil
	}
	if s.execution.Scheduler != nil {
		s.execution.Scheduler.Stop()
		s.logger.Println("executiontarget scheduler acceptance stopped")
	}
	if s.execution.Lifecycle != nil {
		s.execution.Lifecycle.Shutdown()
		s.logger.Println("executiontarget lifecycle admission closed and reservations aborted")
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
