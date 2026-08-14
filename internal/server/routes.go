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
