package api

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/sanjeevksaini/sovrunn/internal/apiconform"
	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/apivalid"
	"github.com/sanjeevksaini/sovrunn/internal/cloudmodel"
	"github.com/sanjeevksaini/sovrunn/internal/cloudmodel/model"
	"github.com/sanjeevksaini/sovrunn/internal/cloudmodel/validate"
)

// FEATURE-0015 topology collection-create idempotency patterns (design §4.2).
const (
	PatternHostingLocationCreate     = "POST /apis/infrastructure.sovrunn.io/v1alpha1/hosting-locations"
	PatternDatacenterCreate          = "POST /apis/infrastructure.sovrunn.io/v1alpha1/datacenters"
	PatternFaultDomainCreate         = "POST /apis/infrastructure.sovrunn.io/v1alpha1/fault-domains"
	PatternInfrastructureStackCreate = "POST /apis/infrastructure.sovrunn.io/v1alpha1/infrastructure-stacks"
)

// HostingLocationCollectionHandler serves HostingLocation LIST and create.
type HostingLocationCollectionHandler struct {
	Store       *cloudmodel.Store
	Grants      GrantResolver
	Publication *cloudmodel.PublicationCoordinator
	Idempotency *cloudmodel.IdempotencyCoordinator
}

// NewHostingLocationCollectionHandler constructs the collection handler.
func NewHostingLocationCollectionHandler(
	store *cloudmodel.Store,
	grants GrantResolver,
	pub *cloudmodel.PublicationCoordinator,
	idemp *cloudmodel.IdempotencyCoordinator,
) *HostingLocationCollectionHandler {
	return &HostingLocationCollectionHandler{
		Store: store, Grants: grants, Publication: pub, Idempotency: idemp,
	}
}

// ServeHTTP dispatches GET (LIST) and POST (create).
func (h *HostingLocationCollectionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.list(w, r)
	case http.MethodPost:
		h.create(w, r)
	default:
		writeMethodNotAllowed(w, r, "GET, POST")
	}
}

func (h *HostingLocationCollectionHandler) list(w http.ResponseWriter, r *http.Request) {
	principal, prob := authenticatePrincipal(r, h.Grants)
	if prob != nil {
		writeProblem(w, r, prob)
		return
	}
	if cloudmodel.HasForgedGrantHeader(r.Header) {
		writeAuditedDenial(w, r, h.Publication, cloudmodel.NewForgedGrantDenial(
			actorRef(principal), requestID(r), newAuditUID("forged"),
		))
		return
	}
	if len(h.Grants.CoarseLookup(cloudmodel.ActionTopologyRead)) == 0 {
		writeAuditedDenial(w, r, h.Publication, cloudmodel.AuditedDenial{
			Problem: apiproblem.New(apiproblem.CodeAuthorizationDenied).WithDetail("topology.read grant is required"),
			Event: cloudmodel.NewRedactedAuditEvent(cloudmodel.RedactedAuditInput{
				UID: newAuditUID("list-denied"), RequestID: requestID(r),
				Action:  cloudmodel.AuditActionListWithoutReadGrant,
				Outcome: apiconform.AuditOutcomeDenied, Reason: "list_without_read_grant",
				Actor: actorRef(principal),
			}),
		})
		return
	}
	items := h.Store.ListHostingLocations()
	visible := make([]model.HostingLocation, 0, len(items))
	for _, item := range items {
		if topologyReadable(h.Grants, item.Metadata.ScopeRef, item.Metadata.UID) {
			visible = append(visible, item)
		}
	}
	writeJSONSuccess(w, r, http.StatusOK, apimeta.ListEnvelope[model.HostingLocation]{
		TypeMeta: apimeta.TypeMeta{APIVersion: model.APIVersionHostingLocation, Kind: model.KindHostingLocation + "List"},
		Items:    visible,
	})
}

func (h *HostingLocationCollectionHandler) create(w http.ResponseWriter, r *http.Request) {
	principal, prob := authenticatePrincipal(r, h.Grants)
	if prob != nil {
		writeProblem(w, r, prob)
		return
	}
	if cloudmodel.HasForgedGrantHeader(r.Header) {
		writeAuditedDenial(w, r, h.Publication, cloudmodel.NewForgedGrantDenial(
			actorRef(principal), requestID(r), newAuditUID("forged"),
		))
		return
	}
	if prob := validate.ValidateIdempotencyKey(r.Header.Get("Idempotency-Key")); prob != nil {
		writeProblem(w, r, prob)
		return
	}

	phaseOne, prob := DecodePhaseOneHTTP(w, r, apivalid.DefaultLimits())
	if prob != nil {
		writeProblem(w, r, prob)
		return
	}
	if phaseOne.HasBootstrapGrant {
		writeAuditedDenial(w, r, h.Publication, cloudmodel.NewForgedGrantDenial(
			actorRef(principal), requestID(r), newAuditUID("forged-body"),
		))
		return
	}
	writes := h.Grants.CoarseLookup(cloudmodel.ActionTopologyWrite)
	if len(writes) == 0 {
		writeProblem(w, r, apiproblem.New(apiproblem.CodeAuthorizationDenied).WithDetail("topology.write grant is required"))
		return
	}
	if !h.Store.HasCloudPlatform() {
		writeAuditedDenial(w, r, h.Publication, rootGateDenial(principal, requestID(r)))
		return
	}

	// HostingLocation scope is derived solely from the topology.write grant
	// CloudProvider UID (ADH-2026-047 decision 2).
	writeGrant := writes[0]
	if !h.Grants.AuthorizeExact(cloudmodel.ActionTopologyWrite, writeGrant.Scope, "", "") {
		writeProblem(w, r, apiproblem.New(apiproblem.CodeAuthorizationDenied).WithDetail("topology.write is not authorized for CloudProvider scope"))
		return
	}
	scopeRef := cloudProviderScopeRef(h.Store, writeGrant.Scope.UID)

	var decoded model.HostingLocation
	if prob := DecodePhaseTwoStrict(validate.KindHostingLocation, validate.RequestClassCollectionCreate, phaseOne.Body, &decoded); prob != nil {
		writeClassifiedCreateDenial(w, r, h.Publication, principal, prob)
		return
	}

	digest, err := cloudmodel.CanonicalDigest(cloudmodel.CanonicalDigestInput{
		Body: phaseOne.Body, Pattern: PatternHostingLocationCreate,
	})
	if err != nil {
		writeProblem(w, r, apiproblem.New(apiproblem.CodeInternalError).WithDetail("failed to compute idempotency digest"))
		return
	}
	ns := cloudmodel.IdempotencyNamespace{
		PrincipalID: principal, Pattern: PatternHostingLocationCreate, Key: r.Header.Get("Idempotency-Key"),
	}
	reserve := h.Idempotency.Reserve(r.Context(), ns, digest, func(ctx context.Context) *apiproblem.Problem {
		if _, p := authenticatePrincipal(r, h.Grants); p != nil {
			return p
		}
		if !h.Store.HasCloudPlatform() {
			return apiproblem.New(apiproblem.CodeConflict).WithViolations([]apiproblem.Violation{{
				Code: violationCloudPlatformRootRequired, Message: "CloudPlatform must exist",
			}})
		}
		if !h.Grants.AuthorizeExact(cloudmodel.ActionTopologyWrite, writeGrant.Scope, "", "") {
			return apiproblem.New(apiproblem.CodeAuthorizationDenied)
		}
		return nil
	})
	switch reserve.Outcome {
	case cloudmodel.ReserveOutcomeReplay:
		writeReplay(w, r, reserve.Replay)
		return
	case cloudmodel.ReserveOutcomeConflict, cloudmodel.ReserveOutcomeDenied:
		writeProblem(w, r, reserve.Problem)
		return
	}

	now := time.Now().UTC().Format(time.RFC3339)
	uid, err := apimeta.GenerateUID()
	if err != nil {
		h.Idempotency.Abort(ns)
		writeProblem(w, r, apiproblem.New(apiproblem.CodeInternalError).WithDetail("failed to allocate uid"))
		return
	}
	decoded.TypeMeta = apimeta.TypeMeta{APIVersion: model.APIVersionHostingLocation, Kind: model.KindHostingLocation}
	decoded.Metadata.UID = uid
	decoded.Metadata.ScopeRef = scopeRef
	decoded.Metadata.Generation = 1
	decoded.Metadata.CreatedAt = now
	decoded.Metadata.UpdatedAt = now
	decoded.Status = model.HostingLocationStatus{Phase: model.HostingLocationPhaseActive}

	if prob := runFeature0015StageSet(r.Context(), validate.KindHostingLocation, validate.RequestClassCollectionCreate, decoded); prob != nil {
		h.Idempotency.Abort(ns)
		writeProblem(w, r, prob)
		return
	}

	h.Store.BeginPublication()
	defer h.Store.EndPublication()
	staged, stageProb := h.Store.StageCreateHostingLocation(decoded)
	if stageProb != nil {
		h.Idempotency.Abort(ns)
		writeProblem(w, r, stageProb)
		return
	}
	decoded.Metadata.ResourceVersion = "1"
	respBody, err := json.Marshal(decoded)
	if err != nil {
		h.Store.Abort(staged)
		h.Idempotency.Abort(ns)
		writeProblem(w, r, apiproblem.New(apiproblem.CodeInternalError).WithDetail("failed to encode create response"))
		return
	}
	event := cloudmodel.NewRedactedAuditEvent(cloudmodel.RedactedAuditInput{
		UID: newAuditUID("create"), RequestID: requestID(r),
		Action: cloudmodel.AuditActionCollectionCreate, Outcome: apiconform.AuditOutcomeSucceeded,
		Reason: "hostinglocation.created", Actor: actorRef(principal),
		Subject: apimeta.TypedRef{
			APIVersion: model.APIVersionHostingLocation, Kind: model.KindHostingLocation,
			Name: decoded.Metadata.Name, UID: uid,
		},
		SubjectResourceVersion: "1",
	})
	if pubProb := h.Publication.PublishSuccess(r.Context(), cloudmodel.AuditedSuccess{
		Staged: staged,
		Event:  event,
		Idempotency: &cloudmodel.IdempotencyCompletion{
			Namespace: ns, Digest: digest,
			Result: cloudmodel.CompletedResult{StatusCode: http.StatusCreated, Body: respBody},
		},
	}); pubProb != nil {
		writeProblem(w, r, pubProb)
		return
	}
	writeRawJSON(w, r, http.StatusCreated, "application/json", respBody)
}

func topologyReadable(grants GrantResolver, scope *apimeta.ScopeRef, resourceUID string) bool {
	id := apimeta.CanonicalScopeIdentity(scope)
	if id.Kind != apimeta.ScopeCloudProvider || id.UID == "" {
		return false
	}
	return grants.AuthorizeExact(cloudmodel.ActionTopologyRead, id, resourceUID, "")
}

func topologyWritable(grants GrantResolver, scope *apimeta.ScopeRef, resourceUID string) bool {
	id := apimeta.CanonicalScopeIdentity(scope)
	if id.Kind != apimeta.ScopeCloudProvider || id.UID == "" {
		return false
	}
	return grants.AuthorizeExact(cloudmodel.ActionTopologyWrite, id, resourceUID, "")
}

func cloudProviderScopeRef(store *cloudmodel.Store, providerUID string) *apimeta.ScopeRef {
	name := "cloud-provider"
	if store != nil {
		if cp, ok := store.GetCloudProvider(providerUID); ok {
			name = cp.Metadata.Name
		}
	}
	return &apimeta.ScopeRef{TypedRef: apimeta.TypedRef{
		APIVersion: model.APIVersionCloudProvider,
		Kind:       string(apimeta.ScopeCloudProvider),
		Name:       name,
		UID:        providerUID,
	}}
}

func copyScopeRef(src *apimeta.ScopeRef) *apimeta.ScopeRef {
	if src == nil {
		return nil
	}
	cp := *src
	return &cp
}

func topologyProviderMismatch(field string) *apiproblem.Problem {
	return apiproblem.New(apiproblem.CodeValidationFailed).WithViolations([]apiproblem.Violation{{
		Field:   field,
		Code:    validate.ViolationTopologyProviderMismatch,
		Message: "visible topology parent is not writable under the caller's CloudProvider scope",
	}})
}

// resolveTopologyParentAccess reports whether a resolved parent is visible and
// writable under topology.read / topology.write at the parent's CloudProvider
// scope. Callers treat !visible as audited safe 404 and mismatch as unaudited
// VALIDATION_FAILED / VS0_TOPOLOGY_PROVIDER_MISMATCH.
func resolveTopologyParentAccess(
	grants GrantResolver,
	parentScope *apimeta.ScopeRef,
	parentUID string,
	parentFound bool,
) (scope *apimeta.ScopeRef, visible bool, mismatch *apiproblem.Problem) {
	if !parentFound || !topologyReadable(grants, parentScope, parentUID) {
		return nil, false, nil
	}
	if !topologyWritable(grants, parentScope, "") {
		return nil, true, topologyProviderMismatch("/metadata/scopeRef")
	}
	return copyScopeRef(parentScope), true, nil
}
