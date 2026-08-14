package api

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/sanjeevksaini/sovrunn/internal/apiconform"
	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/apivalid"
	"github.com/sanjeevksaini/sovrunn/internal/cloudmodel"
	"github.com/sanjeevksaini/sovrunn/internal/cloudmodel/model"
	"github.com/sanjeevksaini/sovrunn/internal/cloudmodel/validate"
)

// ParticipationActionHandler serves one CloudProviderParticipation item-action
// route under /actions/<action> (ADH-2026-053). Lifecycle mutations are never
// registered as /{uid}:<action>.
type ParticipationActionHandler struct {
	Action      cloudmodel.ParticipationAction
	Store       *cloudmodel.Store
	Grants      GrantResolver
	Publication *cloudmodel.PublicationCoordinator
	Idempotency *cloudmodel.IdempotencyCoordinator
	Guard       *cloudmodel.ParticipationActionGuard
}

// NewParticipationActionHandler constructs a single-action handler.
func NewParticipationActionHandler(
	action cloudmodel.ParticipationAction,
	store *cloudmodel.Store,
	grants GrantResolver,
	pub *cloudmodel.PublicationCoordinator,
	idemp *cloudmodel.IdempotencyCoordinator,
	guard *cloudmodel.ParticipationActionGuard,
) *ParticipationActionHandler {
	if guard == nil {
		guard = &cloudmodel.ParticipationActionGuard{}
	}
	return &ParticipationActionHandler{
		Action: action, Store: store, Grants: grants,
		Publication: pub, Idempotency: idemp, Guard: guard,
	}
}

// ParticipationActionPattern returns the registered ServeMux method/path
// pattern used for idempotency namespaces (design DD-08).
func ParticipationActionPattern(action cloudmodel.ParticipationAction) string {
	return "POST /apis/governance.sovrunn.io/v1alpha1/cloud-provider-participations/{uid}/actions/" + string(action)
}

// ServeHTTP handles POST only for the configured action.
func (h *ParticipationActionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeMethodNotAllowed(w, r, "POST")
		return
	}
	h.action(w, r, r.PathValue("uid"))
}

func (h *ParticipationActionHandler) action(w http.ResponseWriter, r *http.Request, uid string) {
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

	idemKey := r.Header.Get("Idempotency-Key")
	ifMatch := r.Header.Get("If-Match")
	if prob := validate.ValidateIdempotencyKey(idemKey); prob != nil {
		writeProblem(w, r, prob)
		return
	}
	if prob := validate.ValidateRequiredIfMatch(ifMatch); prob != nil {
		writeProblem(w, r, prob)
		return
	}

	body, prob := readParticipationActionBody(w, r)
	if prob != nil {
		writeProblem(w, r, prob)
		return
	}
	if len(body) > 0 {
		phaseOne, decodeProb := DecodePhaseOne(body, apivalid.DefaultLimits())
		if decodeProb != nil {
			writeProblem(w, r, decodeProb)
			return
		}
		if phaseOne.HasBootstrapGrant {
			writeAuditedDenial(w, r, h.Publication, cloudmodel.NewForgedGrantDenial(
				actorRef(principal), requestID(r), newAuditUID("forged-body"),
			))
			return
		}
		if emptyProb := validate.ValidateEmptyActionBody(body); emptyProb != nil {
			writeProblem(w, r, emptyProb)
			return
		}
	}

	if !h.coarseActionAuthorized() {
		writeProblem(w, r, apiproblem.New(apiproblem.CodeAuthorizationDenied).WithDetail(
			"participation action grant is required",
		))
		return
	}

	part, ok := h.Store.GetParticipation(uid)
	if uid == "" || !ok || !participationReadable(h.Grants, part.Metadata.ScopeRef, uid) {
		writeAuditedDenial(w, r, h.Publication, safeDenialGET(principal, requestID(r)))
		return
	}
	if _, authProb := h.authorizeAction(part); authProb != nil {
		writeProblem(w, r, authProb)
		return
	}

	pattern := ParticipationActionPattern(h.Action)
	digest, err := cloudmodel.CanonicalDigest(cloudmodel.CanonicalDigestInput{
		Body: body, Pattern: pattern, IfMatch: ifMatch,
	})
	if err != nil {
		writeProblem(w, r, apiproblem.New(apiproblem.CodeInternalError).WithDetail("failed to compute idempotency digest"))
		return
	}
	ns := cloudmodel.IdempotencyNamespace{
		PrincipalID: principal,
		Pattern:     pattern,
		TargetUID:   uid,
		Key:         idemKey,
	}
	reserve := h.Idempotency.Reserve(r.Context(), ns, digest, func(ctx context.Context) *apiproblem.Problem {
		if _, p := authenticatePrincipal(r, h.Grants); p != nil {
			return p
		}
		cur, exists := h.Store.GetParticipation(uid)
		if !exists || !participationReadable(h.Grants, cur.Metadata.ScopeRef, uid) {
			return apiproblem.New(apiproblem.CodeResourceNotFound).WithViolations([]apiproblem.Violation{{
				Code: violationAuthorizationSafeDenial, Message: "resource not found",
			}})
		}
		if _, p := h.authorizeAction(cur); p != nil {
			return p
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

	if prob := runFeature0015StageSet(
		r.Context(),
		validate.KindCloudProviderParticipation,
		validate.RequestClassItemAction,
		part,
	); prob != nil {
		h.Idempotency.Abort(ns)
		writeProblem(w, r, prob)
		return
	}

	h.Guard.Lock()
	defer h.Guard.Unlock()

	live, ok := h.Store.GetParticipation(uid)
	if !ok || !participationReadable(h.Grants, live.Metadata.ScopeRef, uid) {
		h.Idempotency.Abort(ns)
		writeAuditedDenial(w, r, h.Publication, safeDenialGET(principal, requestID(r)))
		return
	}
	party, authProb := h.authorizeAction(live)
	if authProb != nil {
		h.Idempotency.Abort(ns)
		writeProblem(w, r, authProb)
		return
	}
	// Stale/missing If-Match wins over lifecycle source-state validation
	// (design §5.1 steps 12–13; ADH-2026-046).
	if prob := validate.CheckIfMatchCurrent(ifMatch, live.Metadata.ResourceVersion); prob != nil {
		h.Idempotency.Abort(ns)
		writeProblem(w, r, prob)
		return
	}

	next, stateProb := cloudmodel.ApplyParticipationAction(live.Status, h.Action, party)
	if stateProb != nil {
		h.Idempotency.Abort(ns)
		writeProblem(w, r, stateProb)
		return
	}

	preRV := live.Metadata.ResourceVersion
	live.Status = next
	live.Metadata.UpdatedAt = time.Now().UTC().Format(time.RFC3339)

	h.Store.BeginPublication()
	defer h.Store.EndPublication()
	staged, stageProb := h.Store.StageUpdateParticipation(live)
	if stageProb != nil {
		h.Idempotency.Abort(ns)
		writeProblem(w, r, stageProb)
		return
	}

	resp := live
	resp.Metadata.ResourceVersion = nextOpaqueResourceVersion(preRV)
	respBody, err := json.Marshal(resp)
	if err != nil {
		h.Store.Abort(staged)
		h.Idempotency.Abort(ns)
		writeProblem(w, r, apiproblem.New(apiproblem.CodeInternalError).WithDetail("failed to encode action response"))
		return
	}

	event := cloudmodel.NewRedactedAuditEvent(cloudmodel.RedactedAuditInput{
		UID: newAuditUID("action-" + string(h.Action)), RequestID: requestID(r),
		Action:  cloudmodel.AuditActionParticipationAction,
		Outcome: apiconform.AuditOutcomeSucceeded,
		Reason:  "participation." + string(h.Action),
		Actor:   actorRef(principal),
		Subject: apimeta.TypedRef{
			APIVersion: model.APIVersionCloudProviderParticipation,
			Kind:       model.KindCloudProviderParticipation,
			Name:       live.Metadata.Name,
			UID:        uid,
		},
		SubjectResourceVersion: preRV,
	})
	if pubProb := h.Publication.PublishSuccess(r.Context(), cloudmodel.AuditedSuccess{
		Staged: staged,
		Event:  event,
		Idempotency: &cloudmodel.IdempotencyCompletion{
			Namespace: ns, Digest: digest,
			Result: cloudmodel.CompletedResult{StatusCode: http.StatusOK, Body: respBody},
		},
	}); pubProb != nil {
		writeProblem(w, r, pubProb)
		return
	}
	writeRawJSON(w, r, http.StatusOK, "application/json", respBody)
}

func readParticipationActionBody(w http.ResponseWriter, r *http.Request) ([]byte, *apiproblem.Problem) {
	lim := apivalid.DefaultLimits()
	body := r.Body
	if body == nil {
		body = http.NoBody
	}
	if lim.MaxObjectBytes > 0 {
		body = http.MaxBytesReader(w, body, int64(lim.MaxObjectBytes))
	}
	data, err := io.ReadAll(body)
	if err != nil {
		var maxErr *http.MaxBytesError
		if errors.As(err, &maxErr) {
			return nil, apiproblem.New(apiproblem.CodeRequestTooLarge).WithDetail("request body exceeds MaxObjectBytes")
		}
		return nil, apiproblem.New(apiproblem.CodeMalformedRequest).WithDetail("failed to read request body")
	}
	return data, nil
}

// nextOpaqueResourceVersion mirrors store.nextResourceVersion for response
// encoding before publication (caller's copy is not mutated by StageUpdate*).
func nextOpaqueResourceVersion(current string) string {
	if current == "" {
		return "1"
	}
	n, err := strconv.ParseUint(current, 10, 64)
	if err != nil {
		return "1"
	}
	return strconv.FormatUint(n+1, 10)
}

func (h *ParticipationActionHandler) coarseActionAuthorized() bool {
	if h == nil || h.Grants == nil {
		return false
	}
	switch h.Action {
	case cloudmodel.ParticipationActionSuspend:
		return len(h.Grants.CoarseLookup(cloudmodel.ActionSuspendPlatform))+
			len(h.Grants.CoarseLookup(cloudmodel.ActionSuspendProvider)) > 0
	case cloudmodel.ParticipationActionResume:
		return len(h.Grants.CoarseLookup(cloudmodel.ActionResumePlatform))+
			len(h.Grants.CoarseLookup(cloudmodel.ActionResumeProvider)) > 0
	default:
		action, ok := participationActionGrant(h.Action)
		if !ok {
			return false
		}
		return len(h.Grants.CoarseLookup(action)) > 0
	}
}

func (h *ParticipationActionHandler) authorizeAction(
	part model.CloudProviderParticipation,
) (cloudmodel.HoldParty, *apiproblem.Problem) {
	platformUID := part.Spec.CloudPlatformRef.UID
	providerUID := part.Spec.CloudProviderRef.UID
	switch h.Action {
	case cloudmodel.ParticipationActionSuspend, cloudmodel.ParticipationActionResume:
		return resolveSuspendResumeFromGrants(h.Grants, h.Action, platformUID, providerUID)
	case cloudmodel.ParticipationActionAccept, cloudmodel.ParticipationActionReject, cloudmodel.ParticipationActionRequestRelease:
		action, _ := participationActionGrant(h.Action)
		scope := apimeta.ScopeIdentity{Kind: apimeta.ScopeCloudPlatform, UID: platformUID}
		if !h.Grants.AuthorizeExact(action, scope, "", "") {
			return "", apiproblem.New(apiproblem.CodeAuthorizationDenied).WithDetail("participation platform action grant is required")
		}
		return "", nil
	case cloudmodel.ParticipationActionWithdraw, cloudmodel.ParticipationActionAcceptRelease, cloudmodel.ParticipationActionDeclineRelease:
		action, _ := participationActionGrant(h.Action)
		scope := apimeta.ScopeIdentity{Kind: apimeta.ScopeCloudProvider, UID: providerUID}
		if !h.Grants.AuthorizeExact(action, scope, "", "") {
			return "", apiproblem.New(apiproblem.CodeAuthorizationDenied).WithDetail("participation provider action grant is required")
		}
		return "", nil
	default:
		return "", apiproblem.New(apiproblem.CodeAuthorizationDenied).WithDetail("unknown participation action")
	}
}

func participationActionGrant(action cloudmodel.ParticipationAction) (string, bool) {
	switch action {
	case cloudmodel.ParticipationActionAccept:
		return cloudmodel.ActionParticipationAcceptPlatform, true
	case cloudmodel.ParticipationActionReject:
		return cloudmodel.ActionParticipationRejectPlatform, true
	case cloudmodel.ParticipationActionRequestRelease:
		return cloudmodel.ActionParticipationRequestReleasePlatform, true
	case cloudmodel.ParticipationActionWithdraw:
		return cloudmodel.ActionParticipationWithdrawProvider, true
	case cloudmodel.ParticipationActionAcceptRelease:
		return cloudmodel.ActionParticipationAcceptReleaseProvider, true
	case cloudmodel.ParticipationActionDeclineRelease:
		return cloudmodel.ActionParticipationDeclineReleaseProvider, true
	default:
		return "", false
	}
}

func resolveSuspendResumeFromGrants(
	grants GrantResolver,
	action cloudmodel.ParticipationAction,
	platformUID, providerUID string,
) (cloudmodel.HoldParty, *apiproblem.Problem) {
	var scoped []cloudmodel.ScopedActionGrant
	if grants != nil {
		for _, name := range []string{
			cloudmodel.ActionSuspendPlatform,
			cloudmodel.ActionSuspendProvider,
			cloudmodel.ActionResumePlatform,
			cloudmodel.ActionResumeProvider,
		} {
			for _, g := range grants.CoarseLookup(name) {
				scoped = append(scoped, g.ScopedActionGrant())
			}
		}
	}
	return cloudmodel.ResolveSuspendResumeParty(action, platformUID, providerUID, scoped)
}
