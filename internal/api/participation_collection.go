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

// FEATURE-0015 CloudProviderParticipation collection-create idempotency pattern
// (design §4.2).
const PatternParticipationCreate = "POST /apis/governance.sovrunn.io/v1alpha1/cloud-provider-participations"

// ParticipationCollectionHandler serves CloudProviderParticipation LIST and create.
type ParticipationCollectionHandler struct {
	Store       *cloudmodel.Store
	Grants      GrantResolver
	Publication *cloudmodel.PublicationCoordinator
	Idempotency *cloudmodel.IdempotencyCoordinator
}

// NewParticipationCollectionHandler constructs the collection handler.
func NewParticipationCollectionHandler(
	store *cloudmodel.Store,
	grants GrantResolver,
	pub *cloudmodel.PublicationCoordinator,
	idemp *cloudmodel.IdempotencyCoordinator,
) *ParticipationCollectionHandler {
	return &ParticipationCollectionHandler{
		Store: store, Grants: grants, Publication: pub, Idempotency: idemp,
	}
}

// ServeHTTP dispatches GET (LIST) and POST (create).
func (h *ParticipationCollectionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.list(w, r)
	case http.MethodPost:
		h.create(w, r)
	default:
		writeMethodNotAllowed(w, r, "GET, POST")
	}
}

func (h *ParticipationCollectionHandler) list(w http.ResponseWriter, r *http.Request) {
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
	if len(h.Grants.CoarseLookup(cloudmodel.ActionParticipationRead)) == 0 {
		writeAuditedDenial(w, r, h.Publication, cloudmodel.AuditedDenial{
			Problem: apiproblem.New(apiproblem.CodeAuthorizationDenied).WithDetail("participation.read grant is required"),
			Event: cloudmodel.NewRedactedAuditEvent(cloudmodel.RedactedAuditInput{
				UID: newAuditUID("list-denied"), RequestID: requestID(r),
				Action:  cloudmodel.AuditActionListWithoutReadGrant,
				Outcome: apiconform.AuditOutcomeDenied, Reason: "list_without_read_grant",
				Actor: actorRef(principal),
			}),
		})
		return
	}
	items := h.Store.ListParticipations()
	visible := make([]model.CloudProviderParticipation, 0, len(items))
	for _, item := range items {
		if participationReadable(h.Grants, item.Metadata.ScopeRef, item.Metadata.UID) {
			visible = append(visible, item)
		}
	}
	writeJSONSuccess(w, r, http.StatusOK, apimeta.ListEnvelope[model.CloudProviderParticipation]{
		TypeMeta: apimeta.TypeMeta{
			APIVersion: model.APIVersionCloudProviderParticipation,
			Kind:       model.KindCloudProviderParticipation + "List",
		},
		Items: visible,
	})
}

func (h *ParticipationCollectionHandler) create(w http.ResponseWriter, r *http.Request) {
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
	if prob := validate.ValidateParticipationCreateHeaders(
		r.Header.Get("Idempotency-Key"), r.Header.Get("If-Match"),
	); prob != nil {
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
	if len(h.Grants.CoarseLookup(cloudmodel.ActionParticipationRequestProvider)) == 0 {
		writeProblem(w, r, apiproblem.New(apiproblem.CodeAuthorizationDenied).WithDetail("participation.request.provider grant is required"))
		return
	}
	if !h.Store.HasCloudPlatform() {
		writeAuditedDenial(w, r, h.Publication, rootGateDenial(principal, requestID(r)))
		return
	}

	platformRef, hasPlatform := phaseOne.Refs["spec.cloudPlatformRef"]
	providerRef, hasProvider := phaseOne.Refs["spec.cloudProviderRef"]
	var (
		scopeRef      *apimeta.ScopeRef
		providerScope apimeta.ScopeIdentity
		platformUID   string
		providerUID   string
	)
	if hasPlatform && platformRef.UID != "" && hasProvider && providerRef.UID != "" {
		platform, platformOK := h.Store.GetCloudPlatform(platformRef.UID)
		if !platformOK || !h.Grants.AuthorizeExact(
			cloudmodel.ActionCloudPlatformRead, cloudmodel.PlatformRootScope, platform.Metadata.UID, "",
		) {
			writeAuditedDenial(w, r, h.Publication, safeDenialGET(principal, requestID(r)))
			return
		}
		provider, providerOK := h.Store.GetCloudProvider(providerRef.UID)
		if !providerOK || !h.Grants.AuthorizeExact(
			cloudmodel.ActionCloudProviderRead, cloudmodel.PlatformRootScope, provider.Metadata.UID, "",
		) {
			writeAuditedDenial(w, r, h.Publication, safeDenialGET(principal, requestID(r)))
			return
		}
		providerScope = apimeta.ScopeIdentity{Kind: apimeta.ScopeCloudProvider, UID: provider.Metadata.UID}
		platformUID = platform.Metadata.UID
		providerUID = provider.Metadata.UID
		if !h.Grants.AuthorizeExact(
			cloudmodel.ActionParticipationRequestProvider, providerScope, "", platformUID,
		) {
			writeAuditedDenial(w, r, h.Publication, safeDenialGET(principal, requestID(r)))
			return
		}
		scopeRef = cloudPlatformScopeRef(platform)
	}

	var decoded model.CloudProviderParticipation
	if prob := DecodePhaseTwoStrict(validate.KindCloudProviderParticipation, validate.RequestClassCollectionCreate, phaseOne.Body, &decoded); prob != nil {
		writeClassifiedCreateDenial(w, r, h.Publication, principal, prob)
		return
	}
	if scopeRef == nil {
		writeProblem(w, r, apiproblem.New(apiproblem.CodeValidationFailed).WithViolations([]apiproblem.Violation{{
			Field:   "/spec/cloudPlatformRef/uid",
			Code:    apiproblem.ViolationOutOfRange,
			Message: "cloudPlatformRef and cloudProviderRef uids are required",
		}}))
		return
	}

	digest, err := cloudmodel.CanonicalDigest(cloudmodel.CanonicalDigestInput{
		Body: phaseOne.Body, Pattern: PatternParticipationCreate,
	})
	if err != nil {
		writeProblem(w, r, apiproblem.New(apiproblem.CodeInternalError).WithDetail("failed to compute idempotency digest"))
		return
	}
	ns := cloudmodel.IdempotencyNamespace{
		PrincipalID: principal, Pattern: PatternParticipationCreate, Key: r.Header.Get("Idempotency-Key"),
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
		if !h.Grants.AuthorizeExact(
			cloudmodel.ActionParticipationRequestProvider, providerScope, "", platformUID,
		) {
			return apiproblem.New(apiproblem.CodeAuthorizationDenied)
		}
		if _, ok := h.Store.GetCloudPlatform(platformUID); !ok {
			return apiproblem.New(apiproblem.CodeResourceNotFound).WithViolations([]apiproblem.Violation{{
				Code: violationAuthorizationSafeDenial, Message: "resource not found",
			}})
		}
		if _, ok := h.Store.GetCloudProvider(providerUID); !ok {
			return apiproblem.New(apiproblem.CodeResourceNotFound).WithViolations([]apiproblem.Violation{{
				Code: violationAuthorizationSafeDenial, Message: "resource not found",
			}})
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

	now := time.Now().UTC()
	nowStr := now.Format(time.RFC3339)
	uid, err := apimeta.GenerateUID()
	if err != nil {
		h.Idempotency.Abort(ns)
		writeProblem(w, r, apiproblem.New(apiproblem.CodeInternalError).WithDetail("failed to allocate uid"))
		return
	}
	decoded.TypeMeta = apimeta.TypeMeta{
		APIVersion: model.APIVersionCloudProviderParticipation,
		Kind:       model.KindCloudProviderParticipation,
	}
	decoded.Metadata.UID = uid
	decoded.Metadata.ScopeRef = scopeRef
	decoded.Metadata.Generation = 1
	decoded.Metadata.CreatedAt = nowStr
	decoded.Metadata.UpdatedAt = nowStr
	decoded.Status = cloudmodel.InitialParticipationStatus(now)

	if prob := runFeature0015StageSet(r.Context(), validate.KindCloudProviderParticipation, validate.RequestClassCollectionCreate, decoded); prob != nil {
		h.Idempotency.Abort(ns)
		writeProblem(w, r, prob)
		return
	}

	h.Store.BeginPublication()
	defer h.Store.EndPublication()
	staged, stageProb := h.Store.StageCreateParticipation(decoded)
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
		Action: cloudmodel.AuditActionParticipationCreate, Outcome: apiconform.AuditOutcomeSucceeded,
		Reason: "participation.created", Actor: actorRef(principal),
		Subject: apimeta.TypedRef{
			APIVersion: model.APIVersionCloudProviderParticipation,
			Kind:       model.KindCloudProviderParticipation,
			Name:       decoded.Metadata.Name,
			UID:        uid,
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

func participationReadable(grants GrantResolver, scope *apimeta.ScopeRef, resourceUID string) bool {
	id := apimeta.CanonicalScopeIdentity(scope)
	if id.Kind != apimeta.ScopeCloudPlatform || id.UID == "" {
		return false
	}
	return grants.AuthorizeExact(cloudmodel.ActionParticipationRead, id, resourceUID, "")
}

func cloudPlatformScopeRef(platform model.CloudPlatform) *apimeta.ScopeRef {
	return &apimeta.ScopeRef{TypedRef: apimeta.TypedRef{
		APIVersion: model.APIVersionCloudPlatform,
		Kind:       string(apimeta.ScopeCloudPlatform),
		Name:       platform.Metadata.Name,
		UID:        platform.Metadata.UID,
	}}
}
