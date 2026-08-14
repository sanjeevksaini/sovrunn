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

// DatacenterCollectionHandler serves Datacenter LIST and create.
type DatacenterCollectionHandler struct {
	Store       *cloudmodel.Store
	Grants      GrantResolver
	Publication *cloudmodel.PublicationCoordinator
	Idempotency *cloudmodel.IdempotencyCoordinator
}

// NewDatacenterCollectionHandler constructs the collection handler.
func NewDatacenterCollectionHandler(
	store *cloudmodel.Store,
	grants GrantResolver,
	pub *cloudmodel.PublicationCoordinator,
	idemp *cloudmodel.IdempotencyCoordinator,
) *DatacenterCollectionHandler {
	return &DatacenterCollectionHandler{
		Store: store, Grants: grants, Publication: pub, Idempotency: idemp,
	}
}

// ServeHTTP dispatches GET (LIST) and POST (create).
func (h *DatacenterCollectionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.list(w, r)
	case http.MethodPost:
		h.create(w, r)
	default:
		writeMethodNotAllowed(w, r, "GET, POST")
	}
}

func (h *DatacenterCollectionHandler) list(w http.ResponseWriter, r *http.Request) {
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
	items := h.Store.ListDatacenters()
	visible := make([]model.Datacenter, 0, len(items))
	for _, item := range items {
		if topologyReadable(h.Grants, item.Metadata.ScopeRef, item.Metadata.UID) {
			visible = append(visible, item)
		}
	}
	writeJSONSuccess(w, r, http.StatusOK, apimeta.ListEnvelope[model.Datacenter]{
		TypeMeta: apimeta.TypeMeta{APIVersion: model.APIVersionDatacenter, Kind: model.KindDatacenter + "List"},
		Items:    visible,
	})
}

func (h *DatacenterCollectionHandler) create(w http.ResponseWriter, r *http.Request) {
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
	if len(h.Grants.CoarseLookup(cloudmodel.ActionTopologyWrite)) == 0 {
		writeProblem(w, r, apiproblem.New(apiproblem.CodeAuthorizationDenied).WithDetail("topology.write grant is required"))
		return
	}
	if !h.Store.HasCloudPlatform() {
		writeAuditedDenial(w, r, h.Publication, rootGateDenial(principal, requestID(r)))
		return
	}

	parentRef, hasRef := phaseOne.Refs["spec.hostingLocationRef"]
	var (
		parent   model.HostingLocation
		scopeRef *apimeta.ScopeRef
	)
	if hasRef && parentRef.UID != "" {
		var parentOK bool
		parent, parentOK = h.Store.GetHostingLocation(parentRef.UID)
		var visible bool
		var mismatch *apiproblem.Problem
		scopeRef, visible, mismatch = resolveTopologyParentAccess(h.Grants, parent.Metadata.ScopeRef, parent.Metadata.UID, parentOK)
		if !visible {
			writeAuditedDenial(w, r, h.Publication, safeDenialGET(principal, requestID(r)))
			return
		}
		if mismatch != nil {
			writeProblem(w, r, mismatch)
			return
		}
	}

	var decoded model.Datacenter
	if prob := DecodePhaseTwoStrict(validate.KindDatacenter, validate.RequestClassCollectionCreate, phaseOne.Body, &decoded); prob != nil {
		writeClassifiedCreateDenial(w, r, h.Publication, principal, prob)
		return
	}
	if scopeRef == nil {
		// Parent ref absent/incomplete after phase one: closed-contract decode
		// already accepted the body shape only when the ref object was present;
		// reject before idempotency digest.
		writeProblem(w, r, apiproblem.New(apiproblem.CodeValidationFailed).WithViolations([]apiproblem.Violation{{
			Field:   "/spec/hostingLocationRef/uid",
			Code:    apiproblem.ViolationOutOfRange,
			Message: "hostingLocationRef uid is required",
		}}))
		return
	}

	digest, err := cloudmodel.CanonicalDigest(cloudmodel.CanonicalDigestInput{
		Body: phaseOne.Body, Pattern: PatternDatacenterCreate,
	})
	if err != nil {
		writeProblem(w, r, apiproblem.New(apiproblem.CodeInternalError).WithDetail("failed to compute idempotency digest"))
		return
	}
	ns := cloudmodel.IdempotencyNamespace{
		PrincipalID: principal, Pattern: PatternDatacenterCreate, Key: r.Header.Get("Idempotency-Key"),
	}
	parentScopeID := apimeta.CanonicalScopeIdentity(scopeRef)
	reserve := h.Idempotency.Reserve(r.Context(), ns, digest, func(ctx context.Context) *apiproblem.Problem {
		if _, p := authenticatePrincipal(r, h.Grants); p != nil {
			return p
		}
		if !h.Store.HasCloudPlatform() {
			return apiproblem.New(apiproblem.CodeConflict).WithViolations([]apiproblem.Violation{{
				Code: violationCloudPlatformRootRequired, Message: "CloudPlatform must exist",
			}})
		}
		if !h.Grants.AuthorizeExact(cloudmodel.ActionTopologyWrite, parentScopeID, "", "") {
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
	decoded.TypeMeta = apimeta.TypeMeta{APIVersion: model.APIVersionDatacenter, Kind: model.KindDatacenter}
	decoded.Metadata.UID = uid
	decoded.Metadata.ScopeRef = scopeRef
	decoded.Metadata.Generation = 1
	decoded.Metadata.CreatedAt = now
	decoded.Metadata.UpdatedAt = now
	decoded.Status = model.DatacenterStatus{Phase: model.DatacenterPhaseActive}

	if prob := runFeature0015StageSet(r.Context(), validate.KindDatacenter, validate.RequestClassCollectionCreate, decoded); prob != nil {
		h.Idempotency.Abort(ns)
		writeProblem(w, r, prob)
		return
	}
	if prob := validate.ValidateDatacenterParentChain(decoded, parent); prob != nil {
		h.Idempotency.Abort(ns)
		writeProblem(w, r, prob)
		return
	}

	h.Store.BeginPublication()
	defer h.Store.EndPublication()
	staged, stageProb := h.Store.StageCreateDatacenter(decoded)
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
		Reason: "datacenter.created", Actor: actorRef(principal),
		Subject: apimeta.TypedRef{
			APIVersion: model.APIVersionDatacenter, Kind: model.KindDatacenter,
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
