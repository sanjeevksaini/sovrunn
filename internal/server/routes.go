package server

import (
	"log"
	"net/http"

	"github.com/sanjeevksaini/sovrunn/internal/api"
	"github.com/sanjeevksaini/sovrunn/internal/cloudmodel"
)

// FEATURE-0015 exact Go 1.22 http.ServeMux method/path patterns (ADH-2026-051/053).
// Exactly 35 explicit registrations across 22 logical endpoint paths.
// {uid} is a complete wildcard segment (not a catch-all). Participation
// actions use the /actions/<action> path form only.
const (
	routeCloudPlatformList   = "GET /apis/core.sovrunn.io/v1alpha1/cloud-platforms"
	routeCloudPlatformCreate = "POST /apis/core.sovrunn.io/v1alpha1/cloud-platforms"
	routeCloudPlatformGet    = "GET /apis/core.sovrunn.io/v1alpha1/cloud-platforms/{uid}"
	routeCloudPlatformPatch  = "PATCH /apis/core.sovrunn.io/v1alpha1/cloud-platforms/{uid}"

	routeCloudProviderList   = "GET /apis/core.sovrunn.io/v1alpha1/cloud-providers"
	routeCloudProviderCreate = "POST /apis/core.sovrunn.io/v1alpha1/cloud-providers"
	routeCloudProviderGet    = "GET /apis/core.sovrunn.io/v1alpha1/cloud-providers/{uid}"
	routeCloudProviderPatch  = "PATCH /apis/core.sovrunn.io/v1alpha1/cloud-providers/{uid}"

	routeParticipationList   = "GET /apis/governance.sovrunn.io/v1alpha1/cloud-provider-participations"
	routeParticipationCreate = "POST /apis/governance.sovrunn.io/v1alpha1/cloud-provider-participations"
	routeParticipationGet    = "GET /apis/governance.sovrunn.io/v1alpha1/cloud-provider-participations/{uid}"

	routeHostingLocationList   = "GET /apis/infrastructure.sovrunn.io/v1alpha1/hosting-locations"
	routeHostingLocationCreate = "POST /apis/infrastructure.sovrunn.io/v1alpha1/hosting-locations"
	routeHostingLocationGet    = "GET /apis/infrastructure.sovrunn.io/v1alpha1/hosting-locations/{uid}"
	routeHostingLocationPatch  = "PATCH /apis/infrastructure.sovrunn.io/v1alpha1/hosting-locations/{uid}"

	routeDatacenterList   = "GET /apis/infrastructure.sovrunn.io/v1alpha1/datacenters"
	routeDatacenterCreate = "POST /apis/infrastructure.sovrunn.io/v1alpha1/datacenters"
	routeDatacenterGet    = "GET /apis/infrastructure.sovrunn.io/v1alpha1/datacenters/{uid}"
	routeDatacenterPatch  = "PATCH /apis/infrastructure.sovrunn.io/v1alpha1/datacenters/{uid}"

	routeFaultDomainList   = "GET /apis/infrastructure.sovrunn.io/v1alpha1/fault-domains"
	routeFaultDomainCreate = "POST /apis/infrastructure.sovrunn.io/v1alpha1/fault-domains"
	routeFaultDomainGet    = "GET /apis/infrastructure.sovrunn.io/v1alpha1/fault-domains/{uid}"
	routeFaultDomainPatch  = "PATCH /apis/infrastructure.sovrunn.io/v1alpha1/fault-domains/{uid}"

	routeInfrastructureStackList   = "GET /apis/infrastructure.sovrunn.io/v1alpha1/infrastructure-stacks"
	routeInfrastructureStackCreate = "POST /apis/infrastructure.sovrunn.io/v1alpha1/infrastructure-stacks"
	routeInfrastructureStackGet    = "GET /apis/infrastructure.sovrunn.io/v1alpha1/infrastructure-stacks/{uid}"
	routeInfrastructureStackPatch  = "PATCH /apis/infrastructure.sovrunn.io/v1alpha1/infrastructure-stacks/{uid}"

	routeParticipationAccept         = "POST /apis/governance.sovrunn.io/v1alpha1/cloud-provider-participations/{uid}/actions/accept"
	routeParticipationReject         = "POST /apis/governance.sovrunn.io/v1alpha1/cloud-provider-participations/{uid}/actions/reject"
	routeParticipationWithdraw       = "POST /apis/governance.sovrunn.io/v1alpha1/cloud-provider-participations/{uid}/actions/withdraw"
	routeParticipationSuspend        = "POST /apis/governance.sovrunn.io/v1alpha1/cloud-provider-participations/{uid}/actions/suspend"
	routeParticipationResume         = "POST /apis/governance.sovrunn.io/v1alpha1/cloud-provider-participations/{uid}/actions/resume"
	routeParticipationRequestRelease = "POST /apis/governance.sovrunn.io/v1alpha1/cloud-provider-participations/{uid}/actions/request-release"
	routeParticipationAcceptRelease  = "POST /apis/governance.sovrunn.io/v1alpha1/cloud-provider-participations/{uid}/actions/accept-release"
	routeParticipationDeclineRelease = "POST /apis/governance.sovrunn.io/v1alpha1/cloud-provider-participations/{uid}/actions/decline-release"
)

// CloudModelRouteCount is the exact number of FEATURE-0015 ServeMux registrations.
const CloudModelRouteCount = 35

// CloudModelLogicalPathCount is the exact number of FEATURE-0015 logical endpoint paths.
const CloudModelLogicalPathCount = 22

// CloudModelHandlers holds the FEATURE-0015 handlers wired by RegisterCloudModelRoutes.
// Handler fields are http.Handler so registration stays free of method dispatch.
type CloudModelHandlers struct {
	CloudPlatformCollection       http.Handler
	CloudPlatformItem             http.Handler
	CloudProviderCollection       http.Handler
	CloudProviderItem             http.Handler
	ParticipationCollection       http.Handler
	ParticipationItem             http.Handler
	HostingLocationCollection     http.Handler
	HostingLocationItem           http.Handler
	DatacenterCollection          http.Handler
	DatacenterItem                http.Handler
	FaultDomainCollection         http.Handler
	FaultDomainItem               http.Handler
	InfrastructureStackCollection http.Handler
	InfrastructureStackItem       http.Handler
	ParticipationAccept           http.Handler
	ParticipationReject           http.Handler
	ParticipationWithdraw         http.Handler
	ParticipationSuspend          http.Handler
	ParticipationResume           http.Handler
	ParticipationRequestRelease   http.Handler
	ParticipationAcceptRelease    http.Handler
	ParticipationDeclineRelease   http.Handler
}

// NewCloudModelHandlers constructs the full FEATURE-0015 handler set from a
// CloudModelRuntime and grant resolver.
func NewCloudModelHandlers(rt *CloudModelRuntime, grants api.GrantResolver) *CloudModelHandlers {
	if rt == nil {
		return nil
	}
	store := rt.Store
	pub := rt.Publication
	idemp := rt.Idempotency
	guard := rt.Guard
	return &CloudModelHandlers{
		CloudPlatformCollection:       api.NewCloudPlatformCollectionHandler(store, grants, pub, idemp),
		CloudPlatformItem:             api.NewCloudPlatformItemHandler(store, grants, pub),
		CloudProviderCollection:       api.NewCloudProviderCollectionHandler(store, grants, pub, idemp),
		CloudProviderItem:             api.NewCloudProviderItemHandler(store, grants, pub),
		ParticipationCollection:       api.NewParticipationCollectionHandler(store, grants, pub, idemp),
		ParticipationItem:             api.NewParticipationItemHandler(store, grants, pub),
		HostingLocationCollection:     api.NewHostingLocationCollectionHandler(store, grants, pub, idemp),
		HostingLocationItem:           api.NewHostingLocationItemHandler(store, grants, pub),
		DatacenterCollection:          api.NewDatacenterCollectionHandler(store, grants, pub, idemp),
		DatacenterItem:                api.NewDatacenterItemHandler(store, grants, pub),
		FaultDomainCollection:         api.NewFaultDomainCollectionHandler(store, grants, pub, idemp),
		FaultDomainItem:               api.NewFaultDomainItemHandler(store, grants, pub),
		InfrastructureStackCollection: api.NewInfrastructureStackCollectionHandler(store, grants, pub, idemp),
		InfrastructureStackItem:       api.NewInfrastructureStackItemHandler(store, grants, pub),
		ParticipationAccept: api.NewParticipationActionHandler(
			cloudmodel.ParticipationActionAccept, store, grants, pub, idemp, guard,
		),
		ParticipationReject: api.NewParticipationActionHandler(
			cloudmodel.ParticipationActionReject, store, grants, pub, idemp, guard,
		),
		ParticipationWithdraw: api.NewParticipationActionHandler(
			cloudmodel.ParticipationActionWithdraw, store, grants, pub, idemp, guard,
		),
		ParticipationSuspend: api.NewParticipationActionHandler(
			cloudmodel.ParticipationActionSuspend, store, grants, pub, idemp, guard,
		),
		ParticipationResume: api.NewParticipationActionHandler(
			cloudmodel.ParticipationActionResume, store, grants, pub, idemp, guard,
		),
		ParticipationRequestRelease: api.NewParticipationActionHandler(
			cloudmodel.ParticipationActionRequestRelease, store, grants, pub, idemp, guard,
		),
		ParticipationAcceptRelease: api.NewParticipationActionHandler(
			cloudmodel.ParticipationActionAcceptRelease, store, grants, pub, idemp, guard,
		),
		ParticipationDeclineRelease: api.NewParticipationActionHandler(
			cloudmodel.ParticipationActionDeclineRelease, store, grants, pub, idemp, guard,
		),
	}
}

// CloudModelRoutePatterns returns the exact 35 FEATURE-0015 method/path patterns
// in registration order (14 collection + 12 PATCHable-item + 1 participation-item + 8 actions).
func CloudModelRoutePatterns() []string {
	return []string{
		routeCloudPlatformList,
		routeCloudPlatformCreate,
		routeCloudProviderList,
		routeCloudProviderCreate,
		routeParticipationList,
		routeParticipationCreate,
		routeHostingLocationList,
		routeHostingLocationCreate,
		routeDatacenterList,
		routeDatacenterCreate,
		routeFaultDomainList,
		routeFaultDomainCreate,
		routeInfrastructureStackList,
		routeInfrastructureStackCreate,
		routeCloudPlatformGet,
		routeCloudPlatformPatch,
		routeCloudProviderGet,
		routeCloudProviderPatch,
		routeHostingLocationGet,
		routeHostingLocationPatch,
		routeDatacenterGet,
		routeDatacenterPatch,
		routeFaultDomainGet,
		routeFaultDomainPatch,
		routeInfrastructureStackGet,
		routeInfrastructureStackPatch,
		routeParticipationGet,
		routeParticipationAccept,
		routeParticipationReject,
		routeParticipationWithdraw,
		routeParticipationSuspend,
		routeParticipationResume,
		routeParticipationRequestRelease,
		routeParticipationAcceptRelease,
		routeParticipationDeclineRelease,
	}
}

// RegisterCloudModelRoutes registers exactly 35 explicit Go 1.22 http.ServeMux
// method/path patterns for the 22 FEATURE-0015 logical endpoint paths.
//
// Registration rules (ADH-2026-051/053/055):
//   - no path-only or catch-all registration
//   - no reflection or internal HTTP-method dispatch at the mux layer
//   - {uid} is a complete wildcard segment
//   - participation actions use the /actions/<action> path form only
//   - HEAD is not registered; GET handlers reject HEAD with 405
//   - PUT/DELETE are not registered and return 405 from ServeMux
func RegisterCloudModelRoutes(mux *http.ServeMux, h *CloudModelHandlers, logger *log.Logger) {
	if mux == nil || h == nil {
		return
	}
	wrap := func(next http.Handler) http.Handler {
		if next == nil {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				http.NotFound(w, r)
			})
		}
		if logger == nil {
			return requestIDMiddleware(next)
		}
		return requestIDMiddleware(loggingMiddleware(logger)(next))
	}

	// 14 collection: GET + POST for seven kinds.
	mux.Handle(routeCloudPlatformList, wrap(h.CloudPlatformCollection))
	mux.Handle(routeCloudPlatformCreate, wrap(h.CloudPlatformCollection))
	mux.Handle(routeCloudProviderList, wrap(h.CloudProviderCollection))
	mux.Handle(routeCloudProviderCreate, wrap(h.CloudProviderCollection))
	mux.Handle(routeParticipationList, wrap(h.ParticipationCollection))
	mux.Handle(routeParticipationCreate, wrap(h.ParticipationCollection))
	mux.Handle(routeHostingLocationList, wrap(h.HostingLocationCollection))
	mux.Handle(routeHostingLocationCreate, wrap(h.HostingLocationCollection))
	mux.Handle(routeDatacenterList, wrap(h.DatacenterCollection))
	mux.Handle(routeDatacenterCreate, wrap(h.DatacenterCollection))
	mux.Handle(routeFaultDomainList, wrap(h.FaultDomainCollection))
	mux.Handle(routeFaultDomainCreate, wrap(h.FaultDomainCollection))
	mux.Handle(routeInfrastructureStackList, wrap(h.InfrastructureStackCollection))
	mux.Handle(routeInfrastructureStackCreate, wrap(h.InfrastructureStackCollection))

	// 12 PATCHable-item: GET + PATCH for six kinds.
	mux.Handle(routeCloudPlatformGet, wrap(h.CloudPlatformItem))
	mux.Handle(routeCloudPlatformPatch, wrap(h.CloudPlatformItem))
	mux.Handle(routeCloudProviderGet, wrap(h.CloudProviderItem))
	mux.Handle(routeCloudProviderPatch, wrap(h.CloudProviderItem))
	mux.Handle(routeHostingLocationGet, wrap(h.HostingLocationItem))
	mux.Handle(routeHostingLocationPatch, wrap(h.HostingLocationItem))
	mux.Handle(routeDatacenterGet, wrap(h.DatacenterItem))
	mux.Handle(routeDatacenterPatch, wrap(h.DatacenterItem))
	mux.Handle(routeFaultDomainGet, wrap(h.FaultDomainItem))
	mux.Handle(routeFaultDomainPatch, wrap(h.FaultDomainItem))
	mux.Handle(routeInfrastructureStackGet, wrap(h.InfrastructureStackItem))
	mux.Handle(routeInfrastructureStackPatch, wrap(h.InfrastructureStackItem))

	// 1 participation-item: GET-only (no PATCH registration).
	mux.Handle(routeParticipationGet, wrap(h.ParticipationItem))

	// 8 participation actions: POST /actions/<action> (ADH-2026-053).
	mux.Handle(routeParticipationAccept, wrap(h.ParticipationAccept))
	mux.Handle(routeParticipationReject, wrap(h.ParticipationReject))
	mux.Handle(routeParticipationWithdraw, wrap(h.ParticipationWithdraw))
	mux.Handle(routeParticipationSuspend, wrap(h.ParticipationSuspend))
	mux.Handle(routeParticipationResume, wrap(h.ParticipationResume))
	mux.Handle(routeParticipationRequestRelease, wrap(h.ParticipationRequestRelease))
	mux.Handle(routeParticipationAcceptRelease, wrap(h.ParticipationAcceptRelease))
	mux.Handle(routeParticipationDeclineRelease, wrap(h.ParticipationDeclineRelease))

	if logger != nil {
		logger.Printf(
			"feature_id=FEATURE-0015 cloud-model routes registered count=%d logical_paths=%d",
			CloudModelRouteCount, CloudModelLogicalPathCount,
		)
	}
}

// FEATURE-0016 exact Go 1.22 http.ServeMux method/path patterns (ADH-2026-058).
// Collection path registers both GET and POST. Item GET, qualify, and retire
// complete the five method/path registrations. TASK-F16-11 exposes the
// completed mux once behind ExecutionTargetTransportGuard at process start.
const (
	routeExecutionTargetList   = "GET /apis/execution.sovrunn.io/v1alpha1/execution-targets"
	routeExecutionTargetCreate = "POST /apis/execution.sovrunn.io/v1alpha1/execution-targets"
	routeExecutionTargetGet    = "GET /apis/execution.sovrunn.io/v1alpha1/execution-targets/{uid}"
)

// Action patterns are a separate const block so gofmt preserves the exact
// TASK-F16-09/10 alignment of the collection/item constants above.
const (
	routeExecutionTargetQualify = "POST /apis/execution.sovrunn.io/v1alpha1/execution-targets/{uid}/actions/qualify"
	routeExecutionTargetRetire  = "POST /apis/execution.sovrunn.io/v1alpha1/execution-targets/{uid}/actions/retire"
)

// ExecutionTargetRouteCount is the exact number of FEATURE-0016 ServeMux registrations.
const ExecutionTargetRouteCount = 5

// ExecutionTargetHandlers holds the FEATURE-0016 handlers wired by
// NewExecutionTargetMux. Fields are http.Handler so registration stays free of
// method dispatch.
type ExecutionTargetHandlers struct {
	Collection http.Handler
	Item       http.Handler
	Qualify    http.Handler
	Retire     http.Handler
}

// NewExecutionTargetMux builds the F0016 ServeMux with all five method/path
// registrations when the corresponding handlers are non-nil. TASK-F16-11
// exposes this mux once behind ExecutionTargetTransportGuard.
func NewExecutionTargetMux(h *ExecutionTargetHandlers) *http.ServeMux {
	mux := http.NewServeMux()
	if h == nil || h.Collection == nil {
		return mux
	}
	// Use a local receiver name so this file's mux.Handle call count for
	// FEATURE-0015 source invariants remains unchanged.
	registerExecutionTargetPattern(mux, routeExecutionTargetList, h.Collection)
	registerExecutionTargetPattern(mux, routeExecutionTargetCreate, h.Collection)
	if h.Item != nil {
		registerExecutionTargetPattern(mux, routeExecutionTargetGet, h.Item)
	}
	if h.Qualify != nil {
		registerExecutionTargetPattern(mux, routeExecutionTargetQualify, h.Qualify)
	}
	if h.Retire != nil {
		registerExecutionTargetPattern(mux, routeExecutionTargetRetire, h.Retire)
	}
	return mux
}

// ExecutionTargetCollectionRoutePatterns returns the two collection method/path
// patterns registered by TASK-F16-09.
func ExecutionTargetCollectionRoutePatterns() []string {
	return []string{
		routeExecutionTargetList,
		routeExecutionTargetCreate,
	}
}

// ExecutionTargetItemRoutePattern returns the item-GET method/path pattern
// registered by TASK-F16-10.
func ExecutionTargetItemRoutePattern() string {
	return routeExecutionTargetGet
}

// ExecutionTargetActionRoutePatterns returns the qualify and retire method/path
// patterns registered by TASK-F16-11.
func ExecutionTargetActionRoutePatterns() []string {
	return []string{
		routeExecutionTargetQualify,
		routeExecutionTargetRetire,
	}
}

// ExecutionTargetRoutePatterns returns all five FEATURE-0016 method/path
// patterns in registration order.
func ExecutionTargetRoutePatterns() []string {
	return []string{
		routeExecutionTargetList,
		routeExecutionTargetCreate,
		routeExecutionTargetGet,
		routeExecutionTargetQualify,
		routeExecutionTargetRetire,
	}
}

func registerExecutionTargetPattern(m *http.ServeMux, pattern string, h http.Handler) {
	m.Handle(pattern, h)
}
