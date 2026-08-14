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

// CloudProviderCollectionHandler serves CloudProvider LIST and create.
type CloudProviderCollectionHandler struct {
	Store       *cloudmodel.Store
	Grants      GrantResolver
	Publication *cloudmodel.PublicationCoordinator
	Idempotency *cloudmodel.IdempotencyCoordinator
}

// NewCloudProviderCollectionHandler constructs the collection handler.
func NewCloudProviderCollectionHandler(
	store *cloudmodel.Store,
	grants GrantResolver,
	pub *cloudmodel.PublicationCoordinator,
	idemp *cloudmodel.IdempotencyCoordinator,
) *CloudProviderCollectionHandler {
	return &CloudProviderCollectionHandler{
		Store: store, Grants: grants, Publication: pub, Idempotency: idemp,
	}
}

// ServeHTTP dispatches GET (LIST) and POST (create).
func (h *CloudProviderCollectionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.list(w, r)
	case http.MethodPost:
		h.create(w, r)
	default:
		writeMethodNotAllowed(w, r, "GET, POST")
	}
}

func (h *CloudProviderCollectionHandler) list(w http.ResponseWriter, r *http.Request) {
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
	if len(h.Grants.CoarseLookup(cloudmodel.ActionCloudProviderRead)) == 0 {
		writeAuditedDenial(w, r, h.Publication, cloudmodel.AuditedDenial{
			Problem: apiproblem.New(apiproblem.CodeAuthorizationDenied).WithDetail("cloudprovider.read grant is required"),
			Event: cloudmodel.NewRedactedAuditEvent(cloudmodel.RedactedAuditInput{
				UID: newAuditUID("list-denied"), RequestID: requestID(r),
				Action:  cloudmodel.AuditActionListWithoutReadGrant,
				Outcome: apiconform.AuditOutcomeDenied, Reason: "list_without_read_grant",
				Actor: actorRef(principal),
			}),
		})
		return
	}
	items := h.Store.ListCloudProviders()
	visible := make([]model.CloudProvider, 0, len(items))
	for _, item := range items {
		if h.Grants.AuthorizeExact(cloudmodel.ActionCloudProviderRead, cloudmodel.PlatformRootScope, item.Metadata.UID, "") {
			visible = append(visible, item)
		}
	}
	writeJSONSuccess(w, r, http.StatusOK, apimeta.ListEnvelope[model.CloudProvider]{
		TypeMeta: apimeta.TypeMeta{APIVersion: model.APIVersionCloudProvider, Kind: model.KindCloudProvider + "List"},
		Items:    visible,
	})
}

func (h *CloudProviderCollectionHandler) create(w http.ResponseWriter, r *http.Request) {
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
	if len(h.Grants.CoarseLookup(cloudmodel.ActionCloudProviderWrite)) == 0 {
		writeProblem(w, r, apiproblem.New(apiproblem.CodeAuthorizationDenied).WithDetail("cloudprovider.write grant is required"))
		return
	}
	if !h.Store.HasCloudPlatform() {
		writeAuditedDenial(w, r, h.Publication, rootGateDenial(principal, requestID(r)))
		return
	}
	if !h.Grants.AuthorizeExact(cloudmodel.ActionCloudProviderWrite, cloudmodel.PlatformRootScope, "", "") {
		writeProblem(w, r, apiproblem.New(apiproblem.CodeAuthorizationDenied).WithDetail("cloudprovider.write is not authorized for Platform-root"))
		return
	}

	var decoded model.CloudProvider
	if prob := DecodePhaseTwoStrict(validate.KindCloudProvider, validate.RequestClassCollectionCreate, phaseOne.Body, &decoded); prob != nil {
		writeClassifiedCreateDenial(w, r, h.Publication, principal, prob)
		return
	}

	digest, err := cloudmodel.CanonicalDigest(cloudmodel.CanonicalDigestInput{
		Body: phaseOne.Body, Pattern: PatternCloudProviderCreate,
	})
	if err != nil {
		writeProblem(w, r, apiproblem.New(apiproblem.CodeInternalError).WithDetail("failed to compute idempotency digest"))
		return
	}
	ns := cloudmodel.IdempotencyNamespace{
		PrincipalID: principal, Pattern: PatternCloudProviderCreate, Key: r.Header.Get("Idempotency-Key"),
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
		if !h.Grants.AuthorizeExact(cloudmodel.ActionCloudProviderWrite, cloudmodel.PlatformRootScope, "", "") {
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
	decoded.TypeMeta = apimeta.TypeMeta{APIVersion: model.APIVersionCloudProvider, Kind: model.KindCloudProvider}
	decoded.Metadata.UID = uid
	decoded.Metadata.ScopeRef = nil
	decoded.Metadata.Generation = 1
	decoded.Metadata.CreatedAt = now
	decoded.Metadata.UpdatedAt = now
	decoded.Status = model.CloudProviderStatus{Phase: model.CloudProviderPhaseActive}

	if prob := runFeature0015StageSet(r.Context(), validate.KindCloudProvider, validate.RequestClassCollectionCreate, decoded); prob != nil {
		h.Idempotency.Abort(ns)
		writeProblem(w, r, prob)
		return
	}

	h.Store.BeginPublication()
	defer h.Store.EndPublication()
	staged, stageProb := h.Store.StageCreateCloudProvider(decoded)
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
		Reason: "cloudprovider.created", Actor: actorRef(principal),
		Subject: apimeta.TypedRef{
			APIVersion: model.APIVersionCloudProvider, Kind: model.KindCloudProvider,
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
