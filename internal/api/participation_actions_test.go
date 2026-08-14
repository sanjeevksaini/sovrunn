package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/cloudmodel"
	"github.com/sanjeevksaini/sovrunn/internal/cloudmodel/model"
	"github.com/sanjeevksaini/sovrunn/internal/cloudmodel/validate"
)

var participationActions = []cloudmodel.ParticipationAction{
	cloudmodel.ParticipationActionAccept,
	cloudmodel.ParticipationActionReject,
	cloudmodel.ParticipationActionWithdraw,
	cloudmodel.ParticipationActionSuspend,
	cloudmodel.ParticipationActionResume,
	cloudmodel.ParticipationActionRequestRelease,
	cloudmodel.ParticipationActionAcceptRelease,
	cloudmodel.ParticipationActionDeclineRelease,
}

func participationActionGrants(principal, platformUID, providerUID string) *cloudmodel.BootstrapGrantResolver {
	return cloudmodel.NewBootstrapGrantResolver(cloudmodel.BootstrapGrantConfig{
		PrincipalID:       principal,
		CloudPlatformUIDs: []string{platformUID},
		CloudProviderUIDs: []string{providerUID},
		ParticipationRelations: []cloudmodel.ParticipationRelation{{
			CloudPlatformUID: platformUID,
			CloudProviderUID: providerUID,
		}},
		SuspendResume: cloudmodel.SuspendResumeGrantSpec{
			PlatformUIDs: []string{platformUID},
		},
	})
}

func (rt *f15Runtime) participationAction(action cloudmodel.ParticipationAction) http.Handler {
	return NewParticipationActionHandler(action, rt.store, rt.grants, rt.pub, rt.idemp, &cloudmodel.ParticipationActionGuard{})
}

func muxParticipationWithActions(rt *f15Runtime) *http.ServeMux {
	mux := muxParticipation(rt)
	for _, action := range participationActions {
		pattern := "POST /apis/governance.sovrunn.io/v1alpha1/cloud-provider-participations/{uid}/actions/" + string(action)
		mux.Handle(pattern, rt.participationAction(action))
	}
	return mux
}

func actionPath(uid string, action cloudmodel.ParticipationAction) string {
	return "/apis/governance.sovrunn.io/v1alpha1/cloud-provider-participations/" + uid + "/actions/" + string(action)
}

func actionHeaders(rv, key string) map[string]string {
	return map[string]string{
		"If-Match":        rv,
		"Idempotency-Key": key,
	}
}

func seedPendingParticipation(t *testing.T, rt *f15Runtime, name string) model.CloudProviderParticipation {
	t.Helper()
	seedPlatformAndProviderForParticipation(t, rt, testPlatformUID, testProviderUID)
	mux := muxParticipation(rt)
	rec := doMux(t, mux, http.MethodPost, "/apis/governance.sovrunn.io/v1alpha1/cloud-provider-participations",
		participationCreateBody(name, testPlatformUID, testProviderUID), map[string]string{
			"Content-Type": "application/json", "Idempotency-Key": "seed-" + name,
		})
	if rec.Code != http.StatusCreated {
		t.Fatalf("seed create status=%d body=%s", rec.Code, rec.Body.String())
	}
	var created model.CloudProviderParticipation
	if err := json.Unmarshal(rec.Body.Bytes(), &created); err != nil {
		t.Fatalf("seed decode: %v", err)
	}
	return created
}

func TestParticipationActions_EightHappyPathsAndAudit(t *testing.T) {
	type step struct {
		action cloudmodel.ParticipationAction
		want   model.ParticipationPhase
	}
	// Independent runtimes so each chain starts from a fresh Pending participation.
	chains := [][]step{
		{{cloudmodel.ParticipationActionAccept, model.ParticipationPhaseActive}},
		{{cloudmodel.ParticipationActionReject, model.ParticipationPhaseRejected}},
		{{cloudmodel.ParticipationActionWithdraw, model.ParticipationPhaseWithdrawn}},
		{
			{cloudmodel.ParticipationActionAccept, model.ParticipationPhaseActive},
			{cloudmodel.ParticipationActionSuspend, model.ParticipationPhaseSuspended},
		},
		{
			{cloudmodel.ParticipationActionAccept, model.ParticipationPhaseActive},
			{cloudmodel.ParticipationActionSuspend, model.ParticipationPhaseSuspended},
			{cloudmodel.ParticipationActionResume, model.ParticipationPhaseActive},
		},
		{
			{cloudmodel.ParticipationActionAccept, model.ParticipationPhaseActive},
			{cloudmodel.ParticipationActionRequestRelease, model.ParticipationPhaseTerminating},
		},
		{
			{cloudmodel.ParticipationActionAccept, model.ParticipationPhaseActive},
			{cloudmodel.ParticipationActionRequestRelease, model.ParticipationPhaseTerminating},
			{cloudmodel.ParticipationActionAcceptRelease, model.ParticipationPhaseTerminated},
		},
		{
			{cloudmodel.ParticipationActionAccept, model.ParticipationPhaseActive},
			{cloudmodel.ParticipationActionRequestRelease, model.ParticipationPhaseTerminating},
			{cloudmodel.ParticipationActionDeclineRelease, model.ParticipationPhaseActive},
		},
	}

	for i, chain := range chains {
		rt := newF15Runtime(participationActionGrants(testPrincipal, testPlatformUID, testProviderUID))
		mux := muxParticipationWithActions(rt)
		part := seedPendingParticipation(t, rt, "happy-"+string(rune('a'+i)))
		before := rt.audit.Len()
		for j, st := range chain {
			headers := actionHeaders(part.Metadata.ResourceVersion, "happy-"+string(rune('a'+i))+"-"+string(st.action)+"-"+string(rune('0'+j)))
			rec := doMux(t, mux, http.MethodPost, actionPath(part.Metadata.UID, st.action), nil, headers)
			if rec.Code != http.StatusOK {
				t.Fatalf("chain %d step %s status=%d body=%s", i, st.action, rec.Code, rec.Body.String())
			}
			if err := json.Unmarshal(rec.Body.Bytes(), &part); err != nil {
				t.Fatalf("decode: %v", err)
			}
			if part.Status.Phase != st.want {
				t.Fatalf("chain %d step %s phase=%q want %q", i, st.action, part.Status.Phase, st.want)
			}
		}
		if rt.audit.Len() != before+len(chain) {
			t.Fatalf("chain %d audit=%d want %d", i, rt.audit.Len(), before+len(chain))
		}
	}
}

func TestParticipationActions_EmptyBodyContract(t *testing.T) {
	rt := newF15Runtime(participationActionGrants(testPrincipal, testPlatformUID, testProviderUID))
	mux := muxParticipationWithActions(rt)
	part := seedPendingParticipation(t, rt, "body-contract")

	// EOF/zero bytes accepted.
	headers := actionHeaders(part.Metadata.ResourceVersion, "body-empty")
	rec := doMux(t, mux, http.MethodPost, actionPath(part.Metadata.UID, cloudmodel.ParticipationActionAccept), nil, headers)
	if rec.Code != http.StatusOK {
		t.Fatalf("empty body status=%d body=%s", rec.Code, rec.Body.String())
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &part); err != nil {
		t.Fatalf("decode: %v", err)
	}

	// Body validation precedes lifecycle: non-empty bodies fail closed with no AuditEvent.
	before := rt.audit.Len()
	for i, body := range [][]byte{[]byte("   "), []byte("{}"), []byte(`{"x":1}`)} {
		h := actionHeaders(part.Metadata.ResourceVersion, "body-bad-"+string(rune('0'+i)))
		rec = doMux(t, mux, http.MethodPost, actionPath(part.Metadata.UID, cloudmodel.ParticipationActionSuspend), body, h)
		if rec.Code != http.StatusBadRequest || decodeProblem(t, rec).Code != apiproblem.CodeMalformedRequest {
			t.Fatalf("non-empty body %q: %d %#v", body, rec.Code, decodeProblem(t, rec))
		}
		if rt.audit.Len() != before {
			t.Fatalf("malformed body must not audit")
		}
	}
}

func TestParticipationActions_PreconditionsAndSourceState(t *testing.T) {
	rt := newF15Runtime(participationActionGrants(testPrincipal, testPlatformUID, testProviderUID))
	mux := muxParticipationWithActions(rt)
	part := seedPendingParticipation(t, rt, "precond")

	before := rt.audit.Len()
	rec := doMux(t, mux, http.MethodPost, actionPath(part.Metadata.UID, cloudmodel.ParticipationActionAccept), nil, map[string]string{
		"Idempotency-Key": "missing-if-match",
	})
	if rec.Code != http.StatusPreconditionFailed || decodeProblem(t, rec).Code != apiproblem.CodeStaleResourceVersion {
		t.Fatalf("missing If-Match: %d %#v", rec.Code, decodeProblem(t, rec))
	}
	if rt.audit.Len() != before {
		t.Fatal("missing If-Match must not audit")
	}

	rec = doMux(t, mux, http.MethodPost, actionPath(part.Metadata.UID, cloudmodel.ParticipationActionAccept), nil, map[string]string{
		"If-Match": "1",
	})
	if rec.Code != http.StatusBadRequest || decodeProblem(t, rec).Code != apiproblem.CodeMalformedRequest {
		t.Fatalf("missing Idempotency-Key: %d %#v", rec.Code, decodeProblem(t, rec))
	}

	rec = doMux(t, mux, http.MethodPost, actionPath(part.Metadata.UID, cloudmodel.ParticipationActionAccept), nil, actionHeaders("bad\x01", "malformed-if"))
	if rec.Code != http.StatusPreconditionFailed || decodeProblem(t, rec).Code != apiproblem.CodeStaleResourceVersion {
		t.Fatalf("malformed If-Match: %d %#v", rec.Code, decodeProblem(t, rec))
	}

	// Stale If-Match + invalid lifecycle: stale wins (Pending + request-release is invalid).
	before = rt.audit.Len()
	rec = doMux(t, mux, http.MethodPost, actionPath(part.Metadata.UID, cloudmodel.ParticipationActionRequestRelease), nil, actionHeaders("stale", "stale-wins"))
	if rec.Code != http.StatusPreconditionFailed || decodeProblem(t, rec).Code != apiproblem.CodeStaleResourceVersion {
		t.Fatalf("stale wins: %d %#v", rec.Code, decodeProblem(t, rec))
	}
	if rt.audit.Len() != before {
		t.Fatal("stale must not audit")
	}

	// Current If-Match + invalid source state → CONFLICT / VS0_PARTICIPATION_STATE_INVALID.
	before = rt.audit.Len()
	rec = doMux(t, mux, http.MethodPost, actionPath(part.Metadata.UID, cloudmodel.ParticipationActionRequestRelease), nil,
		actionHeaders(part.Metadata.ResourceVersion, "invalid-state"))
	if rec.Code != http.StatusConflict ||
		!hasViolation(decodeProblem(t, rec), cloudmodel.ViolationParticipationStateInvalid) ||
		rt.audit.Len() != before {
		t.Fatalf("invalid state: %d %#v audit=%d", rec.Code, decodeProblem(t, rec), rt.audit.Len())
	}
}

func TestParticipationActions_IdempotencyReplayAndMismatch(t *testing.T) {
	rt := newF15Runtime(participationActionGrants(testPrincipal, testPlatformUID, testProviderUID))
	mux := muxParticipationWithActions(rt)
	part := seedPendingParticipation(t, rt, "idemp")

	headers := actionHeaders(part.Metadata.ResourceVersion, "idemp-accept")
	rec := doMux(t, mux, http.MethodPost, actionPath(part.Metadata.UID, cloudmodel.ParticipationActionAccept), nil, headers)
	if rec.Code != http.StatusOK {
		t.Fatalf("accept status=%d", rec.Code)
	}
	before := rt.audit.Len()
	firstBody := append([]byte(nil), rec.Body.Bytes()...)

	// Same-key/same-digest replay — no second AuditEvent.
	rec = doMux(t, mux, http.MethodPost, actionPath(part.Metadata.UID, cloudmodel.ParticipationActionAccept), nil, headers)
	if rec.Code != http.StatusOK || rec.Body.String() != string(firstBody) {
		t.Fatalf("replay status=%d body mismatch", rec.Code)
	}
	if rt.audit.Len() != before {
		t.Fatalf("replay must not audit again, got %d", rt.audit.Len())
	}

	// Same-key/different-digest (changed If-Match) → CONFLICT.
	before = rt.audit.Len()
	headers["If-Match"] = "999"
	rec = doMux(t, mux, http.MethodPost, actionPath(part.Metadata.UID, cloudmodel.ParticipationActionAccept), nil, headers)
	if rec.Code != http.StatusConflict ||
		!hasViolation(decodeProblem(t, rec), cloudmodel.ViolationIdempotencyKeyReuseMismatch) ||
		rt.audit.Len() != before {
		t.Fatalf("digest mismatch: %d %#v audit=%d", rec.Code, decodeProblem(t, rec), rt.audit.Len())
	}

	// Target UID binding: same key on another participation must not replay.
	partB := seedPendingOnNewProvider(t, rt, "idemp-b", "dddddddddddddddddddddddddddddddd")
	headersB := actionHeaders(partB.Metadata.ResourceVersion, "idemp-accept") // same key, different target
	rec = doMux(t, mux, http.MethodPost, actionPath(partB.Metadata.UID, cloudmodel.ParticipationActionAccept), nil, headersB)
	if rec.Code != http.StatusOK {
		t.Fatalf("other target must be independent: %d %s", rec.Code, rec.Body.String())
	}
	var acceptedB model.CloudProviderParticipation
	if err := json.Unmarshal(rec.Body.Bytes(), &acceptedB); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if acceptedB.Metadata.UID == part.Metadata.UID {
		t.Fatal("must not disclose other target replay")
	}
}

func seedPendingOnNewProvider(t *testing.T, rt *f15Runtime, name, providerUID string) model.CloudProviderParticipation {
	t.Helper()
	// Rebuild runtime grants to include the new provider relation by swapping grants
	// is not supported; seed via store with grants that already authorize platform accept.
	if _, ok := rt.store.GetCloudProvider(providerUID); !ok {
		if _, prob := rt.store.CreateCloudProvider(model.CloudProvider{
			TypeMeta: apimeta.TypeMeta{APIVersion: model.APIVersionCloudProvider, Kind: model.KindCloudProvider},
			Metadata: apimeta.ObjectMeta{Name: "prov-" + providerUID[:4], UID: providerUID, ResourceVersion: "1"},
			Spec:     model.CloudProviderSpec{OperatingMarkets: []string{"US"}},
			Status:   model.CloudProviderStatus{Phase: model.CloudProviderPhaseActive},
		}); prob != nil {
			t.Fatalf("seed provider: %#v", prob)
		}
	}
	now := time.Now().UTC()
	uid, err := apimeta.GenerateUID()
	if err != nil {
		t.Fatalf("uid: %v", err)
	}
	p := model.CloudProviderParticipation{
		TypeMeta: apimeta.TypeMeta{APIVersion: model.APIVersionCloudProviderParticipation, Kind: model.KindCloudProviderParticipation},
		Metadata: apimeta.ObjectMeta{
			Name: name, UID: uid, ResourceVersion: "1", Generation: 1,
			ScopeRef: &apimeta.ScopeRef{TypedRef: apimeta.TypedRef{
				APIVersion: model.APIVersionCloudPlatform, Kind: string(apimeta.ScopeCloudPlatform),
				Name: "root-plat", UID: testPlatformUID,
			}},
			CreatedAt: now.Format(time.RFC3339), UpdatedAt: now.Format(time.RFC3339),
		},
		Spec: model.CloudProviderParticipationSpec{
			CloudPlatformRef: apimeta.TypedRef{
				APIVersion: model.APIVersionCloudPlatform, Kind: model.KindCloudPlatform,
				Name: "root-plat", UID: testPlatformUID,
			},
			CloudProviderRef: apimeta.TypedRef{
				APIVersion: model.APIVersionCloudProvider, Kind: model.KindCloudProvider,
				Name: "prov-b", UID: providerUID,
			},
			Environment: model.ParticipationEnvironmentDevelopment,
		},
		Status: cloudmodel.InitialParticipationStatus(now),
	}
	if _, prob := rt.store.CreateParticipation(p); prob != nil {
		t.Fatalf("seed participation: %#v", prob)
	}
	return p
}

func TestParticipationActions_ReplayRecheckRevokedGrant(t *testing.T) {
	grants := participationActionGrants(testPrincipal, testPlatformUID, testProviderUID)
	rt := newF15Runtime(grants)
	mux := muxParticipationWithActions(rt)
	part := seedPendingParticipation(t, rt, "revoke")

	headers := actionHeaders(part.Metadata.ResourceVersion, "revoke-accept")
	rec := doMux(t, mux, http.MethodPost, actionPath(part.Metadata.UID, cloudmodel.ParticipationActionAccept), nil, headers)
	if rec.Code != http.StatusOK {
		t.Fatalf("accept status=%d", rec.Code)
	}
	stored := append([]byte(nil), rec.Body.Bytes()...)

	// Revoke by replacing grants with read-only participation grants.
	rt.grants = &staticGrants{
		principal: testPrincipal,
		grants: []cloudmodel.Grant{
			{Action: cloudmodel.ActionParticipationRead, Scope: apimeta.ScopeIdentity{Kind: apimeta.ScopeCloudPlatform, UID: testPlatformUID}},
		},
	}
	mux = muxParticipationWithActions(rt)
	rec = doMux(t, mux, http.MethodPost, actionPath(part.Metadata.UID, cloudmodel.ParticipationActionAccept), nil, headers)
	if rec.Code == http.StatusOK || rec.Body.String() == string(stored) {
		t.Fatalf("revoked grant must not disclose stored result: %d %s", rec.Code, rec.Body.String())
	}
	if decodeProblem(t, rec).Code != apiproblem.CodeAuthorizationDenied {
		t.Fatalf("want AUTHORIZATION_DENIED, got %#v", decodeProblem(t, rec))
	}
}

func TestParticipationActions_ForgedGrantAndSafeDenialAudit(t *testing.T) {
	rt := newF15Runtime(participationActionGrants(testPrincipal, testPlatformUID, testProviderUID))
	mux := muxParticipationWithActions(rt)
	part := seedPendingParticipation(t, rt, "forge")

	before := rt.audit.Len()
	headers := actionHeaders(part.Metadata.ResourceVersion, "forge-hdr")
	headers[cloudmodel.HeaderBootstrapGrant] = "claim"
	rec := doMux(t, mux, http.MethodPost, actionPath(part.Metadata.UID, cloudmodel.ParticipationActionAccept), nil, headers)
	if rec.Code != http.StatusForbidden || rt.audit.Len() != before+1 {
		t.Fatalf("forged header: %d audit=%d", rec.Code, rt.audit.Len())
	}

	before = rt.audit.Len()
	headers = actionHeaders(part.Metadata.ResourceVersion, "forge-body")
	rec = doMux(t, mux, http.MethodPost, actionPath(part.Metadata.UID, cloudmodel.ParticipationActionAccept),
		[]byte(`{"bootstrapGrant":{"id":"g1"}}`), headers)
	if rec.Code != http.StatusForbidden || rt.audit.Len() != before+1 {
		t.Fatalf("forged body: %d audit=%d", rec.Code, rt.audit.Len())
	}

	before = rt.audit.Len()
	headers = actionHeaders("1", "safe-miss")
	rec = doMux(t, mux, http.MethodPost, actionPath("ffffffffffffffffffffffffffffffff", cloudmodel.ParticipationActionAccept), nil, headers)
	if rec.Code != http.StatusNotFound || !hasViolation(decodeProblem(t, rec), violationAuthorizationSafeDenial) || rt.audit.Len() != before+1 {
		t.Fatalf("safe denial: %d %#v audit=%d", rec.Code, decodeProblem(t, rec), rt.audit.Len())
	}

	rt.audit.SetFail(errors.New("audit down"))
	before = rt.audit.Len()
	headers = actionHeaders(part.Metadata.ResourceVersion, "audit-fail")
	headers[cloudmodel.HeaderBootstrapGrant] = "claim"
	rec = doMux(t, mux, http.MethodPost, actionPath(part.Metadata.UID, cloudmodel.ParticipationActionAccept), nil, headers)
	if rec.Code != http.StatusInternalServerError || rt.audit.Len() != before {
		t.Fatalf("audit append failure: %d audit=%d", rec.Code, rt.audit.Len())
	}
	got, _ := rt.store.GetParticipation(part.Metadata.UID)
	if got.Status.Phase != model.ParticipationPhasePending {
		t.Fatalf("audit failure must not mutate, phase=%q", got.Status.Phase)
	}
}

func TestParticipationActions_IndependentHoldsAndDeniedPhases(t *testing.T) {
	rt := newF15Runtime(participationActionGrants(testPrincipal, testPlatformUID, testProviderUID))
	mux := muxParticipationWithActions(rt)
	part := seedPendingParticipation(t, rt, "holds")

	// Suspend denied in Pending.
	before := rt.audit.Len()
	rec := doMux(t, mux, http.MethodPost, actionPath(part.Metadata.UID, cloudmodel.ParticipationActionSuspend), nil,
		actionHeaders(part.Metadata.ResourceVersion, "sus-pending"))
	if rec.Code != http.StatusConflict || !hasViolation(decodeProblem(t, rec), cloudmodel.ViolationParticipationStateInvalid) {
		t.Fatalf("suspend pending: %d %#v", rec.Code, decodeProblem(t, rec))
	}
	if rt.audit.Len() != before {
		t.Fatal("invalid suspend must not audit")
	}

	rec = doMux(t, mux, http.MethodPost, actionPath(part.Metadata.UID, cloudmodel.ParticipationActionAccept), nil,
		actionHeaders(part.Metadata.ResourceVersion, "holds-accept"))
	if err := json.Unmarshal(rec.Body.Bytes(), &part); err != nil {
		t.Fatalf("decode: %v", err)
	}

	// Platform suspend sets only platform hold.
	rec = doMux(t, mux, http.MethodPost, actionPath(part.Metadata.UID, cloudmodel.ParticipationActionSuspend), nil,
		actionHeaders(part.Metadata.ResourceVersion, "holds-sus"))
	if err := json.Unmarshal(rec.Body.Bytes(), &part); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if !part.Status.PlatformSuspended || part.Status.ProviderSuspended || part.Status.Phase != model.ParticipationPhaseSuspended {
		t.Fatalf("platform suspend holds: %#v", part.Status)
	}

	// Add provider hold via direct status mutation + stage (provider grant not in default SuspendResume).
	part.Status.ProviderSuspended = true
	part.Status.Phase = cloudmodel.DerivedEffectivePhase(true, true)
	if _, prob := rt.store.UpdateParticipation(part); prob != nil {
		t.Fatalf("set provider hold: %#v", prob)
	}
	part, _ = rt.store.GetParticipation(part.Metadata.UID)

	// Resume clears only platform hold; remains Suspended while provider hold true.
	rec = doMux(t, mux, http.MethodPost, actionPath(part.Metadata.UID, cloudmodel.ParticipationActionResume), nil,
		actionHeaders(part.Metadata.ResourceVersion, "holds-res"))
	if rec.Code != http.StatusOK {
		t.Fatalf("resume status=%d body=%s", rec.Code, rec.Body.String())
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &part); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if part.Status.PlatformSuspended || !part.Status.ProviderSuspended || part.Status.Phase != model.ParticipationPhaseSuspended {
		t.Fatalf("clear one hold must not reactivate: %#v", part.Status)
	}

	// Terminating: suspend denied.
	rec = doMux(t, mux, http.MethodPost, actionPath(part.Metadata.UID, cloudmodel.ParticipationActionRequestRelease), nil,
		actionHeaders(part.Metadata.ResourceVersion, "holds-rel"))
	if err := json.Unmarshal(rec.Body.Bytes(), &part); err != nil {
		t.Fatalf("decode: %v", err)
	}
	before = rt.audit.Len()
	rec = doMux(t, mux, http.MethodPost, actionPath(part.Metadata.UID, cloudmodel.ParticipationActionSuspend), nil,
		actionHeaders(part.Metadata.ResourceVersion, "sus-term"))
	if rec.Code != http.StatusConflict || !hasViolation(decodeProblem(t, rec), cloudmodel.ViolationParticipationStateInvalid) || rt.audit.Len() != before {
		t.Fatalf("suspend terminating: %d %#v", rec.Code, decodeProblem(t, rec))
	}

	// Decline-release restores hold-derived Suspended.
	rec = doMux(t, mux, http.MethodPost, actionPath(part.Metadata.UID, cloudmodel.ParticipationActionDeclineRelease), nil,
		actionHeaders(part.Metadata.ResourceVersion, "holds-decline"))
	if err := json.Unmarshal(rec.Body.Bytes(), &part); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if part.Status.Phase != model.ParticipationPhaseSuspended || !part.Status.ProviderSuspended {
		t.Fatalf("decline-release prior effective: %#v", part.Status)
	}
}

func TestParticipationActions_SuspendResumeExactOneGrant(t *testing.T) {
	// Zero suspend grants → fail closed.
	rtZero := newF15Runtime(participationGrants(testPrincipal, testPlatformUID, testProviderUID))
	muxZero := muxParticipationWithActions(rtZero)
	part := seedPendingParticipation(t, rtZero, "zero-grant")
	rec := doMux(t, muxZero, http.MethodPost, actionPath(part.Metadata.UID, cloudmodel.ParticipationActionAccept), nil,
		actionHeaders(part.Metadata.ResourceVersion, "zero-accept"))
	if err := json.Unmarshal(rec.Body.Bytes(), &part); err != nil {
		t.Fatalf("decode: %v", err)
	}
	rec = doMux(t, muxZero, http.MethodPost, actionPath(part.Metadata.UID, cloudmodel.ParticipationActionSuspend), nil,
		actionHeaders(part.Metadata.ResourceVersion, "zero-sus"))
	if rec.Code != http.StatusForbidden {
		t.Fatalf("zero grant status=%d %#v", rec.Code, decodeProblem(t, rec))
	}

	// Multiple matching grants → fail closed.
	rtBoth := newF15Runtime(cloudmodel.NewBootstrapGrantResolver(cloudmodel.BootstrapGrantConfig{
		PrincipalID:       testPrincipal,
		CloudPlatformUIDs: []string{testPlatformUID},
		CloudProviderUIDs: []string{testProviderUID},
		ParticipationRelations: []cloudmodel.ParticipationRelation{{
			CloudPlatformUID: testPlatformUID, CloudProviderUID: testProviderUID,
		}},
		SuspendResume: cloudmodel.SuspendResumeGrantSpec{
			PlatformUIDs: []string{testPlatformUID},
			ProviderUIDs: []string{testProviderUID},
		},
	}))
	muxBoth := muxParticipationWithActions(rtBoth)
	part = seedPendingParticipation(t, rtBoth, "both-grant")
	rec = doMux(t, muxBoth, http.MethodPost, actionPath(part.Metadata.UID, cloudmodel.ParticipationActionAccept), nil,
		actionHeaders(part.Metadata.ResourceVersion, "both-accept"))
	if err := json.Unmarshal(rec.Body.Bytes(), &part); err != nil {
		t.Fatalf("decode: %v", err)
	}
	rec = doMux(t, muxBoth, http.MethodPost, actionPath(part.Metadata.UID, cloudmodel.ParticipationActionSuspend), nil,
		actionHeaders(part.Metadata.ResourceVersion, "both-sus"))
	if rec.Code != http.StatusForbidden {
		t.Fatalf("multiple grant status=%d %#v", rec.Code, decodeProblem(t, rec))
	}
}

func TestParticipationActions_AuthRequiredAndRouteForm(t *testing.T) {
	rt := newF15Runtime(participationActionGrants(testPrincipal, testPlatformUID, testProviderUID))
	mux := muxParticipationWithActions(rt)
	part := seedPendingParticipation(t, rt, "auth")
	before := rt.audit.Len()

	rec := doMuxUnauth(t, mux, http.MethodPost, actionPath(part.Metadata.UID, cloudmodel.ParticipationActionAccept), nil,
		actionHeaders(part.Metadata.ResourceVersion, "auth-miss"))
	if rec.Code != http.StatusUnauthorized || decodeProblem(t, rec).Code != apiproblem.CodeAuthRequired {
		t.Fatalf("AUTH_REQUIRED: %d %#v", rec.Code, decodeProblem(t, rec))
	}
	if rt.audit.Len() != before {
		t.Fatalf("AUTH_REQUIRED must not audit, audit=%d", rt.audit.Len())
	}

	// Retired /{uid}:<action> form must not be registered as an action route.
	// Go 1.22 may match "{uid}:accept" as a single {uid} wildcard on the item
	// path, which returns 405 from the GET-only item handler — not an action.
	rec = doMux(t, mux, http.MethodPost,
		"/apis/governance.sovrunn.io/v1alpha1/cloud-provider-participations/"+part.Metadata.UID+":accept",
		nil, actionHeaders(part.Metadata.ResourceVersion, "colon-form"))
	if rec.Code == http.StatusOK {
		t.Fatalf("colon form must not execute an action, status=%d body=%s", rec.Code, rec.Body.String())
	}
	got, _ := rt.store.GetParticipation(part.Metadata.UID)
	if got.Status.Phase != model.ParticipationPhasePending {
		t.Fatalf("colon form must not transition, phase=%q", got.Status.Phase)
	}
}

func TestParticipationActions_UsesInheritedStageSet(t *testing.T) {
	stages, ok := validate.StageSetFor(validate.KindCloudProviderParticipation, validate.RequestClassItemAction)
	if !ok || stages.Defaulting == nil || stages.Semantic == nil || stages.Reference == nil {
		t.Fatal("FEATURE-0015 must plug item-action validators into inherited StageSet")
	}
	rt := newF15Runtime(participationActionGrants(testPrincipal, testPlatformUID, testProviderUID))
	mux := muxParticipationWithActions(rt)
	part := seedPendingParticipation(t, rt, "stageset")
	rec := doMux(t, mux, http.MethodPost, actionPath(part.Metadata.UID, cloudmodel.ParticipationActionAccept), nil,
		actionHeaders(part.Metadata.ResourceVersion, "stageset-accept"))
	if rec.Code != http.StatusOK {
		t.Fatalf("StageSet path status=%d body=%s", rec.Code, rec.Body.String())
	}
}

func TestParticipationActions_AuditAppendFailureOnSuccess(t *testing.T) {
	rt := newF15Runtime(participationActionGrants(testPrincipal, testPlatformUID, testProviderUID))
	mux := muxParticipationWithActions(rt)
	part := seedPendingParticipation(t, rt, "audit-ok-fail")
	rt.audit.SetFail(errors.New("down"))
	before := rt.audit.Len()
	rec := doMux(t, mux, http.MethodPost, actionPath(part.Metadata.UID, cloudmodel.ParticipationActionAccept), nil,
		actionHeaders(part.Metadata.ResourceVersion, "audit-ok-fail"))
	if rec.Code != http.StatusInternalServerError || rt.audit.Len() != before {
		t.Fatalf("success audit fail: %d audit=%d", rec.Code, rt.audit.Len())
	}
	got, _ := rt.store.GetParticipation(part.Metadata.UID)
	if got.Status.Phase != model.ParticipationPhasePending {
		t.Fatalf("must not publish mutation on audit failure, phase=%q", got.Status.Phase)
	}
}
