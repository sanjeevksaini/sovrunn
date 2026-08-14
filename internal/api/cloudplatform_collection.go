package api

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/sanjeevksaini/sovrunn/internal/apiconform"
	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/apivalid"
	"github.com/sanjeevksaini/sovrunn/internal/cloudmodel"
	"github.com/sanjeevksaini/sovrunn/internal/cloudmodel/model"
	"github.com/sanjeevksaini/sovrunn/internal/cloudmodel/validate"
	"github.com/sanjeevksaini/sovrunn/internal/requestctx"
)

// FEATURE-0015 registered method/path patterns used for idempotency namespaces.
const (
	PatternCloudPlatformCreate = "POST /apis/core.sovrunn.io/v1alpha1/cloud-platforms"
	PatternCloudProviderCreate = "POST /apis/core.sovrunn.io/v1alpha1/cloud-providers"

	mediaMergePatchJSON = "application/merge-patch+json"
)

// Slice 0 violation codes emitted only in violations[].code.
const (
	violationAuthRequired              apiproblem.ViolationCode = "VS0_AUTH_REQUIRED"
	violationAuthorizationSafeDenial   apiproblem.ViolationCode = "VS0_AUTHORIZATION_SAFE_DENIAL"
	violationCloudPlatformRootRequired apiproblem.ViolationCode = "VS0_CLOUDPLATFORM_ROOT_REQUIRED"
)

// GrantResolver is the FEATURE-0015 server-resolved grant surface used by
// handlers. *cloudmodel.BootstrapGrantResolver satisfies it; tests may supply
// narrower fakes for wrong-administrator / missing-read cases.
type GrantResolver interface {
	PrincipalID() string
	CoarseLookup(action string) []cloudmodel.Grant
	AuthorizeExact(action string, scope apimeta.ScopeIdentity, resourceUID, relatedPlatformUID string) bool
}

// CloudPlatformCollectionHandler serves CloudPlatform LIST and create.
type CloudPlatformCollectionHandler struct {
	Store       *cloudmodel.Store
	Grants      GrantResolver
	Publication *cloudmodel.PublicationCoordinator
	Idempotency *cloudmodel.IdempotencyCoordinator
}

// NewCloudPlatformCollectionHandler constructs the collection handler.
func NewCloudPlatformCollectionHandler(
	store *cloudmodel.Store,
	grants GrantResolver,
	pub *cloudmodel.PublicationCoordinator,
	idemp *cloudmodel.IdempotencyCoordinator,
) *CloudPlatformCollectionHandler {
	return &CloudPlatformCollectionHandler{
		Store: store, Grants: grants, Publication: pub, Idempotency: idemp,
	}
}

// ServeHTTP dispatches GET (LIST) and POST (create). PUT/DELETE/HEAD → 405.
func (h *CloudPlatformCollectionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.list(w, r)
	case http.MethodPost:
		h.create(w, r)
	default:
		writeMethodNotAllowed(w, r, "GET, POST")
	}
}

func (h *CloudPlatformCollectionHandler) list(w http.ResponseWriter, r *http.Request) {
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
	if len(h.Grants.CoarseLookup(cloudmodel.ActionCloudPlatformRead)) == 0 {
		writeAuditedDenial(w, r, h.Publication, cloudmodel.AuditedDenial{
			Problem: apiproblem.New(apiproblem.CodeAuthorizationDenied).WithDetail("cloudplatform.read grant is required"),
			Event: cloudmodel.NewRedactedAuditEvent(cloudmodel.RedactedAuditInput{
				UID: newAuditUID("list-denied"), RequestID: requestID(r),
				Action:  cloudmodel.AuditActionListWithoutReadGrant,
				Outcome: apiconform.AuditOutcomeDenied, Reason: "list_without_read_grant",
				Actor: actorRef(principal),
			}),
		})
		return
	}
	items := h.Store.ListCloudPlatforms()
	visible := make([]model.CloudPlatform, 0, len(items))
	for _, item := range items {
		if h.Grants.AuthorizeExact(cloudmodel.ActionCloudPlatformRead, cloudmodel.PlatformRootScope, item.Metadata.UID, "") {
			visible = append(visible, item)
		}
	}
	writeJSONSuccess(w, r, http.StatusOK, apimeta.ListEnvelope[model.CloudPlatform]{
		TypeMeta: apimeta.TypeMeta{APIVersion: model.APIVersionCloudPlatform, Kind: model.KindCloudPlatform + "List"},
		Items:    visible,
	})
}

func (h *CloudPlatformCollectionHandler) create(w http.ResponseWriter, r *http.Request) {
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
	if len(h.Grants.CoarseLookup(cloudmodel.ActionCloudPlatformWrite)) == 0 {
		writeProblem(w, r, apiproblem.New(apiproblem.CodeAuthorizationDenied).WithDetail("cloudplatform.write grant is required"))
		return
	}
	if !h.Grants.AuthorizeExact(cloudmodel.ActionCloudPlatformWrite, cloudmodel.PlatformRootScope, "", "") {
		writeProblem(w, r, apiproblem.New(apiproblem.CodeAuthorizationDenied).WithDetail("cloudplatform.write is not authorized for Platform-root"))
		return
	}

	var decoded model.CloudPlatform
	if prob := DecodePhaseTwoStrict(validate.KindCloudPlatform, validate.RequestClassCollectionCreate, phaseOne.Body, &decoded); prob != nil {
		writeClassifiedCreateDenial(w, r, h.Publication, principal, prob)
		return
	}

	digest, err := cloudmodel.CanonicalDigest(cloudmodel.CanonicalDigestInput{
		Body: phaseOne.Body, Pattern: PatternCloudPlatformCreate,
	})
	if err != nil {
		writeProblem(w, r, apiproblem.New(apiproblem.CodeInternalError).WithDetail("failed to compute idempotency digest"))
		return
	}
	ns := cloudmodel.IdempotencyNamespace{
		PrincipalID: principal, Pattern: PatternCloudPlatformCreate, Key: r.Header.Get("Idempotency-Key"),
	}
	reserve := h.Idempotency.Reserve(r.Context(), ns, digest, func(ctx context.Context) *apiproblem.Problem {
		if _, p := authenticatePrincipal(r, h.Grants); p != nil {
			return p
		}
		if !h.Grants.AuthorizeExact(cloudmodel.ActionCloudPlatformWrite, cloudmodel.PlatformRootScope, "", "") {
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
	decoded.TypeMeta = apimeta.TypeMeta{APIVersion: model.APIVersionCloudPlatform, Kind: model.KindCloudPlatform}
	decoded.Metadata.UID = uid
	decoded.Metadata.ScopeRef = nil // canonical Platform-root
	decoded.Metadata.Generation = 1
	decoded.Metadata.CreatedAt = now
	decoded.Metadata.UpdatedAt = now
	decoded.Status = model.CloudPlatformStatus{Phase: model.CloudPlatformPhaseActive}

	if prob := runFeature0015StageSet(r.Context(), validate.KindCloudPlatform, validate.RequestClassCollectionCreate, decoded); prob != nil {
		h.Idempotency.Abort(ns)
		writeProblem(w, r, prob)
		return
	}

	h.Store.BeginPublication()
	defer h.Store.EndPublication()
	staged, stageProb := h.Store.StageCreateCloudPlatform(decoded)
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
		Reason: "cloudplatform.created", Actor: actorRef(principal),
		Subject: apimeta.TypedRef{
			APIVersion: model.APIVersionCloudPlatform, Kind: model.KindCloudPlatform,
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

func writeClassifiedCreateDenial(
	w http.ResponseWriter,
	r *http.Request,
	pub *cloudmodel.PublicationCoordinator,
	principal string,
	prob *apiproblem.Problem,
) {
	if prob == nil {
		return
	}
	if prob.Code == apiproblem.CodeAuthorizationDenied && hasViolation(prob, validate.ViolationStatusFieldWrite) {
		writeAuditedDenial(w, r, pub, cloudmodel.AuditedDenial{
			Problem: prob,
			Event: cloudmodel.NewRedactedAuditEvent(cloudmodel.RedactedAuditInput{
				UID: newAuditUID("status-write"), RequestID: requestID(r),
				Action: cloudmodel.AuditActionStatusWrite, Outcome: apiconform.AuditOutcomeDenied,
				Reason: "status_field_write", Actor: actorRef(principal),
			}),
		})
		return
	}
	if prob.Code == apiproblem.CodeAuthorizationDenied && hasViolation(prob, validate.ViolationSystemOwnedFieldWrite) {
		writeAuditedDenial(w, r, pub, cloudmodel.AuditedDenial{
			Problem: prob,
			Event: cloudmodel.NewRedactedAuditEvent(cloudmodel.RedactedAuditInput{
				UID: newAuditUID("system-owned"), RequestID: requestID(r),
				Action: cloudmodel.AuditActionSystemOwnedFieldWrite, Outcome: apiconform.AuditOutcomeDenied,
				Reason: "system_owned_field_write", Actor: actorRef(principal),
			}),
		})
		return
	}
	writeProblem(w, r, prob)
}

func runFeature0015StageSet(ctx context.Context, kind validate.Kind, class validate.RequestClass, object any) *apiproblem.Problem {
	stages, ok := validate.StageSetFor(kind, class)
	if !ok {
		return apiproblem.New(apiproblem.CodeInternalError).WithDetail("StageSet is not configured")
	}
	obj, err := stages.Defaulting.Apply(ctx, object)
	if err != nil {
		return apiproblem.New(apiproblem.CodeInternalError).WithDetail("defaulting stage failed")
	}
	if v, err := stages.Semantic.Validate(ctx, obj); err != nil {
		return apiproblem.New(apiproblem.CodeInternalError).WithDetail("semantic stage failed")
	} else if len(v) > 0 {
		return apiproblem.New(apiproblem.CodeValidationFailed).WithViolations(v)
	}
	if v, err := stages.Reference.Validate(ctx, obj); err != nil {
		return apiproblem.New(apiproblem.CodeInternalError).WithDetail("reference stage failed")
	} else if len(v) > 0 {
		return apiproblem.New(apiproblem.CodeValidationFailed).WithViolations(v)
	}
	return nil
}

func authenticatePrincipal(r *http.Request, grants GrantResolver) (string, *apiproblem.Problem) {
	if grants == nil {
		return "", authRequiredProblem()
	}
	raw := strings.TrimSpace(r.Header.Get("Authorization"))
	if raw == "" {
		return "", authRequiredProblem()
	}
	const prefix = "Bearer "
	if !strings.HasPrefix(raw, prefix) {
		return "", authRequiredProblem()
	}
	token := strings.TrimSpace(strings.TrimPrefix(raw, prefix))
	if token == "" {
		return "", authRequiredProblem()
	}
	want := grants.PrincipalID()
	if want == "" || token != want {
		return "", authRequiredProblem()
	}
	return token, nil
}

func authRequiredProblem() *apiproblem.Problem {
	return apiproblem.New(apiproblem.CodeAuthRequired).WithViolations([]apiproblem.Violation{{
		Code:    violationAuthRequired,
		Message: "authentication is required",
	}})
}

func actorRef(principalID string) apimeta.TypedRef {
	return apimeta.TypedRef{
		APIVersion: "identity.sovrunn.io/v1alpha1",
		Kind:       "Principal",
		Name:       principalID,
		UID:        principalID,
	}
}

func requestID(r *http.Request) string {
	if id := requestctx.RequestIDFromContext(r.Context()); id != "" {
		return id
	}
	if id := strings.TrimSpace(r.Header.Get("X-Sovrunn-Request-ID")); id != "" {
		return id
	}
	return "req-unknown"
}

func newAuditUID(suffix string) string {
	uid, err := apimeta.GenerateUID()
	if err != nil {
		return "audit-" + suffix
	}
	return "audit-" + suffix + "-" + uid
}

func hasViolation(prob *apiproblem.Problem, code apiproblem.ViolationCode) bool {
	if prob == nil {
		return false
	}
	for _, v := range prob.Violations {
		if v.Code == code {
			return true
		}
	}
	return false
}

func writeProblem(w http.ResponseWriter, r *http.Request, prob *apiproblem.Problem) {
	if prob == nil {
		prob = apiproblem.New(apiproblem.CodeInternalError)
	}
	prob = prob.WithRequestID(requestID(r))
	writeRawJSON(w, r, prob.Status, apiproblem.MediaTypeProblemJSON, mustJSON(prob))
}

func writeJSONSuccess(w http.ResponseWriter, r *http.Request, status int, v any) {
	writeRawJSON(w, r, status, "application/json", mustJSON(v))
}

func writeRawJSON(w http.ResponseWriter, r *http.Request, status int, contentType string, body []byte) {
	w.Header().Set("Content-Type", contentType)
	if id := requestID(r); id != "" {
		w.Header().Set("X-Sovrunn-Request-ID", id)
	}
	w.WriteHeader(status)
	_, _ = w.Write(body)
	if len(body) > 0 && body[len(body)-1] != '\n' {
		_, _ = w.Write([]byte("\n"))
	}
}

func writeReplay(w http.ResponseWriter, r *http.Request, replay *cloudmodel.ReplayResponse) {
	if replay == nil {
		writeProblem(w, r, apiproblem.New(apiproblem.CodeInternalError).WithDetail("missing replay payload"))
		return
	}
	ct := replay.ContentType
	if ct == "" {
		ct = "application/json"
	}
	writeRawJSON(w, r, replay.StatusCode, ct, replay.Body)
}

func writeAuditedDenial(w http.ResponseWriter, r *http.Request, pub *cloudmodel.PublicationCoordinator, denial cloudmodel.AuditedDenial) {
	if pub == nil {
		writeProblem(w, r, apiproblem.New(apiproblem.CodeInternalError).WithDetail("publication coordinator is not configured"))
		return
	}
	prob := pub.PublishDenial(r.Context(), denial)
	writeProblem(w, r, prob)
}

func writeMethodNotAllowed(w http.ResponseWriter, r *http.Request, allow string) {
	w.Header().Set("Allow", allow)
	if id := requestID(r); id != "" {
		w.Header().Set("X-Sovrunn-Request-ID", id)
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func mustJSON(v any) []byte {
	b, err := json.Marshal(v)
	if err != nil {
		p := apiproblem.New(apiproblem.CodeInternalError).WithDetail("json encode failed")
		b, _ = json.Marshal(p)
	}
	return b
}

func rootGateDenial(principal, requestID string) cloudmodel.AuditedDenial {
	return cloudmodel.AuditedDenial{
		Problem: apiproblem.New(apiproblem.CodeConflict).WithViolations([]apiproblem.Violation{{
			Code:    violationCloudPlatformRootRequired,
			Message: "CloudPlatform must exist before creating dependent resources",
		}}),
		Event: cloudmodel.NewRedactedAuditEvent(cloudmodel.RedactedAuditInput{
			UID: newAuditUID("root"), RequestID: requestID,
			Action: cloudmodel.AuditActionRootDenial, Outcome: apiconform.AuditOutcomeDenied,
			Reason: "cloudplatform_root_required", Actor: actorRef(principal),
		}),
	}
}

func safeDenialGET(principal, requestID string) cloudmodel.AuditedDenial {
	return cloudmodel.AuditedDenial{
		Problem: apiproblem.New(apiproblem.CodeResourceNotFound).WithViolations([]apiproblem.Violation{{
			Code:    violationAuthorizationSafeDenial,
			Message: "resource not found",
		}}),
		Event: cloudmodel.NewRedactedAuditEvent(cloudmodel.RedactedAuditInput{
			UID: newAuditUID("safe-denial"), RequestID: requestID,
			Action: cloudmodel.AuditActionSafeDenial, Outcome: apiconform.AuditOutcomeDenied,
			Reason: "authorization_safe_denial", Actor: actorRef(principal),
		}),
	}
}

func mergePatchMediaTypeOK(r *http.Request) bool {
	ct := strings.TrimSpace(r.Header.Get("Content-Type"))
	if ct == mediaMergePatchJSON {
		return true
	}
	return strings.HasPrefix(ct, mediaMergePatchJSON+";")
}
