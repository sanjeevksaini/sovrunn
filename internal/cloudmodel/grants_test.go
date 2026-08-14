package cloudmodel_test

import (
	"context"
	"encoding/json"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/apiconform"
	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/cloudmodel"
)

const (
	testPlatformUID = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	testProviderUID = "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
	testOtherUID    = "cccccccccccccccccccccccccccccccc"
)

func testBootstrapResolver(t *testing.T) *cloudmodel.BootstrapGrantResolver {
	t.Helper()
	return cloudmodel.NewBootstrapGrantResolver(cloudmodel.BootstrapGrantConfig{
		PrincipalID:       "bootstrap-principal",
		CloudPlatformUIDs: []string{testPlatformUID},
		CloudProviderUIDs: []string{testProviderUID},
		ParticipationRelations: []cloudmodel.ParticipationRelation{{
			CloudPlatformUID: testPlatformUID,
			CloudProviderUID: testProviderUID,
		}},
		SuspendResume: cloudmodel.SuspendResumeGrantSpec{
			PlatformUIDs: []string{testPlatformUID},
		},
	})
}

func grantActions(grants []cloudmodel.Grant) map[string][]cloudmodel.Grant {
	out := make(map[string][]cloudmodel.Grant)
	for _, g := range grants {
		out[g.Action] = append(out[g.Action], g)
	}
	return out
}

func TestBootstrapGrants_ExactActionsAndScopes(t *testing.T) {
	t.Parallel()
	r := testBootstrapResolver(t)
	byAction := grantActions(r.Grants())

	assertPlatformRoot := func(action string) {
		t.Helper()
		gs := byAction[action]
		if len(gs) != 1 {
			t.Fatalf("%s: want 1 grant, got %d", action, len(gs))
		}
		if gs[0].Scope != cloudmodel.PlatformRootScope {
			t.Fatalf("%s scope=%#v, want Platform-root", action, gs[0].Scope)
		}
		if gs[0].ResourceUID != "" {
			t.Fatalf("%s must be unrestricted at create time", action)
		}
	}
	assertPlatformRoot(cloudmodel.ActionCloudPlatformWrite)
	assertPlatformRoot(cloudmodel.ActionCloudProviderWrite)
	assertPlatformRoot(cloudmodel.ActionCloudPlatformRead)
	assertPlatformRoot(cloudmodel.ActionCloudProviderRead)

	topoWrite := byAction[cloudmodel.ActionTopologyWrite]
	if len(topoWrite) != 1 || topoWrite[0].Scope != (apimeta.ScopeIdentity{Kind: apimeta.ScopeCloudProvider, UID: testProviderUID}) {
		t.Fatalf("topology.write=%#v", topoWrite)
	}
	topoRead := byAction[cloudmodel.ActionTopologyRead]
	if len(topoRead) != 1 || topoRead[0].Scope.UID != testProviderUID {
		t.Fatalf("topology.read=%#v", topoRead)
	}

	req := byAction[cloudmodel.ActionParticipationRequestProvider]
	if len(req) != 1 {
		t.Fatalf("request.provider count=%d", len(req))
	}
	if req[0].Scope.Kind != apimeta.ScopeCloudProvider || req[0].Scope.UID != testProviderUID {
		t.Fatalf("request.provider scope=%#v", req[0].Scope)
	}
	if req[0].RelatedCloudPlatformUID != testPlatformUID {
		t.Fatalf("request.provider related platform=%q", req[0].RelatedCloudPlatformUID)
	}

	for _, action := range []string{
		cloudmodel.ActionParticipationAcceptPlatform,
		cloudmodel.ActionParticipationRejectPlatform,
		cloudmodel.ActionParticipationRequestReleasePlatform,
		cloudmodel.ActionParticipationRead,
	} {
		gs := byAction[action]
		if len(gs) != 1 || gs[0].Scope != (apimeta.ScopeIdentity{Kind: apimeta.ScopeCloudPlatform, UID: testPlatformUID}) {
			t.Fatalf("%s=%#v", action, gs)
		}
	}

	for _, action := range []string{
		cloudmodel.ActionParticipationWithdrawProvider,
		cloudmodel.ActionParticipationAcceptReleaseProvider,
		cloudmodel.ActionParticipationDeclineReleaseProvider,
	} {
		gs := byAction[action]
		if len(gs) != 1 || gs[0].Scope != (apimeta.ScopeIdentity{Kind: apimeta.ScopeCloudProvider, UID: testProviderUID}) {
			t.Fatalf("%s=%#v", action, gs)
		}
	}

	suspendPlat := byAction[cloudmodel.ActionSuspendPlatform]
	resumePlat := byAction[cloudmodel.ActionResumePlatform]
	if len(suspendPlat) != 1 || suspendPlat[0].Scope.UID != testPlatformUID {
		t.Fatalf("suspend.platform=%#v", suspendPlat)
	}
	if len(resumePlat) != 1 || resumePlat[0].Scope.UID != testPlatformUID {
		t.Fatalf("resume.platform=%#v", resumePlat)
	}
	if len(byAction[cloudmodel.ActionSuspendProvider]) != 0 || len(byAction[cloudmodel.ActionResumeProvider]) != 0 {
		t.Fatal("provider suspend/resume must not be present when only platform UIDs are configured")
	}
}

func TestBootstrapGrants_ActionReturnedOnlyForAllowedResourceScope(t *testing.T) {
	t.Parallel()
	r := testBootstrapResolver(t)

	platformScope := apimeta.ScopeIdentity{Kind: apimeta.ScopeCloudPlatform, UID: testPlatformUID}
	providerScope := apimeta.ScopeIdentity{Kind: apimeta.ScopeCloudProvider, UID: testProviderUID}
	otherProvider := apimeta.ScopeIdentity{Kind: apimeta.ScopeCloudProvider, UID: testOtherUID}
	otherPlatform := apimeta.ScopeIdentity{Kind: apimeta.ScopeCloudPlatform, UID: testOtherUID}

	cases := []struct {
		name               string
		action             string
		scope              apimeta.ScopeIdentity
		resourceUID        string
		relatedPlatformUID string
		want               bool
	}{
		{"cloudplatform.write platform-root", cloudmodel.ActionCloudPlatformWrite, cloudmodel.PlatformRootScope, "", "", true},
		{"cloudplatform.write wrong scope", cloudmodel.ActionCloudPlatformWrite, platformScope, "", "", false},
		{"cloudprovider.write platform-root", cloudmodel.ActionCloudProviderWrite, cloudmodel.PlatformRootScope, "", "", true},
		{"topology.write provider", cloudmodel.ActionTopologyWrite, providerScope, "", "", true},
		{"topology.write other provider", cloudmodel.ActionTopologyWrite, otherProvider, "", "", false},
		{"topology.write platform scope denied", cloudmodel.ActionTopologyWrite, platformScope, "", "", false},
		{"accept.platform", cloudmodel.ActionParticipationAcceptPlatform, platformScope, "", "", true},
		{"accept.platform wrong uid", cloudmodel.ActionParticipationAcceptPlatform, otherPlatform, "", "", false},
		{"withdraw.provider", cloudmodel.ActionParticipationWithdrawProvider, providerScope, "", "", true},
		{"withdraw.provider wrong uid", cloudmodel.ActionParticipationWithdrawProvider, otherProvider, "", "", false},
		{"request.provider with relation", cloudmodel.ActionParticipationRequestProvider, providerScope, "", testPlatformUID, true},
		{"request.provider wrong relation", cloudmodel.ActionParticipationRequestProvider, providerScope, "", testOtherUID, false},
		{"participation.read platform", cloudmodel.ActionParticipationRead, platformScope, "", "", true},
		{"participation.read provider denied", cloudmodel.ActionParticipationRead, providerScope, "", "", false},
		{"cloudplatform.read", cloudmodel.ActionCloudPlatformRead, cloudmodel.PlatformRootScope, "", "", true},
		{"cloudprovider.read", cloudmodel.ActionCloudProviderRead, cloudmodel.PlatformRootScope, "", "", true},
		{"topology.read", cloudmodel.ActionTopologyRead, providerScope, "", "", true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := r.AuthorizeExact(tc.action, tc.scope, tc.resourceUID, tc.relatedPlatformUID)
			if got != tc.want {
				t.Fatalf("AuthorizeExact=%v, want %v", got, tc.want)
			}
		})
	}
}

func TestBootstrapGrants_CoarseLookupIgnoresScope(t *testing.T) {
	t.Parallel()
	r := testBootstrapResolver(t)
	got := r.CoarseLookup(cloudmodel.ActionTopologyWrite)
	if len(got) != 1 {
		t.Fatalf("coarse topology.write=%#v", got)
	}
	if len(r.CoarseLookup("participation.invented")) != 0 {
		t.Fatal("unknown action must return no grants")
	}
	if len(r.CoarseLookup("")) != 0 {
		t.Fatal("empty action must return no grants")
	}
}

func TestBootstrapGrants_ResourceUIDNarrowing(t *testing.T) {
	t.Parallel()
	r := testBootstrapResolver(t)
	target := "dddddddddddddddddddddddddddddddd"
	narrowed, ok := r.NarrowGrant(cloudmodel.ActionCloudPlatformWrite, cloudmodel.PlatformRootScope, target)
	if !ok || narrowed.ResourceUID != target {
		t.Fatalf("narrow=%#v ok=%v", narrowed, ok)
	}
	// Unrestricted stored grant still authorizes the target (empty ResourceUID).
	if !r.AuthorizeExact(cloudmodel.ActionCloudPlatformWrite, cloudmodel.PlatformRootScope, target, "") {
		t.Fatal("unrestricted write must authorize PATCH target")
	}
}

func TestBootstrapGrants_ForgedHeaderIgnoredAndDenied(t *testing.T) {
	t.Parallel()
	r := testBootstrapResolver(t)

	h := http.Header{}
	h.Set(cloudmodel.HeaderBootstrapGrant, `{"action":"cloudplatform.write","scope":"Platform"}`)
	if !cloudmodel.HasForgedGrantHeader(h) {
		t.Fatal("expected forged header detection")
	}
	prob := cloudmodel.EvaluateForgedGrantCarriers(true, false)
	if prob == nil || prob.Code != apiproblem.CodeAuthorizationDenied {
		t.Fatalf("header denial=%#v", prob)
	}
	// Claim must not expand or substitute for server-resolved grants.
	if !r.AuthorizeExact(cloudmodel.ActionCloudPlatformWrite, cloudmodel.PlatformRootScope, "", "") {
		t.Fatal("server-resolved grant must remain authoritative")
	}
	if r.AuthorizeExact("forged.action.from.header", cloudmodel.PlatformRootScope, "", "") {
		t.Fatal("header claim must never authorize an action")
	}

	_, audit, _, pub := newPub(t)
	denial := cloudmodel.NewForgedGrantDenial(actorRef("bootstrap-principal"), "req-forged-header", "audit-forged-header")
	got := pub.PublishDenial(context.Background(), denial)
	if got == nil || got.Code != apiproblem.CodeAuthorizationDenied {
		t.Fatalf("PublishDenial=%#v", got)
	}
	events := audit.Events()
	if len(events) != 1 || events[0].Record.Action != cloudmodel.AuditActionForgedGrant {
		t.Fatalf("audit=%#v", events)
	}
	if events[0].Record.Outcome != apiconform.AuditOutcomeDenied {
		t.Fatalf("outcome=%q", events[0].Record.Outcome)
	}
	raw, _ := json.Marshal(events[0])
	claim := h.Get(cloudmodel.HeaderBootstrapGrant)
	if claim != "" && strings.Contains(string(raw), claim) {
		t.Fatal("forged header claim must not be disclosed in AuditEvent")
	}
}

func TestBootstrapGrants_ForgedBodyMemberIgnoredAndDenied(t *testing.T) {
	t.Parallel()
	r := testBootstrapResolver(t)

	root := map[string]json.RawMessage{
		cloudmodel.BodyMemberBootstrapGrant: json.RawMessage(`{"action":"topology.write"}`),
	}
	if !cloudmodel.HasForgedGrantBodyMember(root) {
		t.Fatal("expected forged body member detection")
	}
	prob := cloudmodel.EvaluateForgedGrantCarriers(false, true)
	if prob == nil || prob.Code != apiproblem.CodeAuthorizationDenied {
		t.Fatalf("body denial=%#v", prob)
	}
	if r.AuthorizeExact("topology.write.from.body", apimeta.ScopeIdentity{Kind: apimeta.ScopeCloudProvider, UID: testProviderUID}, "", "") {
		t.Fatal("body claim must never authorize an action")
	}

	_, audit, _, pub := newPub(t)
	denial := cloudmodel.NewForgedGrantDenial(actorRef("bootstrap-principal"), "req-forged-body", "audit-forged-body")
	got := pub.PublishDenial(context.Background(), denial)
	if got == nil || got.Code != apiproblem.CodeAuthorizationDenied {
		t.Fatalf("PublishDenial=%#v", got)
	}
	events := audit.Events()
	if len(events) != 1 || events[0].Record.Action != cloudmodel.AuditActionForgedGrant {
		t.Fatalf("audit=%#v", events)
	}
	raw, _ := json.Marshal(events[0])
	claim := string(root[cloudmodel.BodyMemberBootstrapGrant])
	if claim != "" && strings.Contains(string(raw), claim) {
		t.Fatal("forged body claim must not be disclosed in AuditEvent")
	}
}

func TestBootstrapGrants_NoForgedCarrierAllowsServerResolvedAuthz(t *testing.T) {
	t.Parallel()
	if prob := cloudmodel.EvaluateForgedGrantCarriers(false, false); prob != nil {
		t.Fatalf("unexpected denial=%#v", prob)
	}
	h := http.Header{}
	h.Set("Authorization", "Bearer token")
	if cloudmodel.HasForgedGrantHeader(h) {
		t.Fatal("unrelated headers must not be treated as forged grants")
	}
	root := map[string]json.RawMessage{"metadata": json.RawMessage(`{"name":"x"}`)}
	if cloudmodel.HasForgedGrantBodyMember(root) {
		t.Fatal("ordinary body members must not be forged-grant carriers")
	}
}

func TestBootstrapGrants_SuspendResumeExactOneAndFailClosed(t *testing.T) {
	t.Parallel()

	platformOnly := cloudmodel.NewBootstrapGrantResolver(cloudmodel.BootstrapGrantConfig{
		PrincipalID: "bootstrap",
		SuspendResume: cloudmodel.SuspendResumeGrantSpec{
			PlatformUIDs: []string{testPlatformUID},
		},
	})
	party, prob := platformOnly.ResolveSuspendResume(cloudmodel.ParticipationActionSuspend, testPlatformUID, testProviderUID)
	if prob != nil || party != cloudmodel.HoldPartyPlatform {
		t.Fatalf("platform-only: party=%q prob=%#v", party, prob)
	}

	providerOnly := cloudmodel.NewBootstrapGrantResolver(cloudmodel.BootstrapGrantConfig{
		PrincipalID: "bootstrap",
		SuspendResume: cloudmodel.SuspendResumeGrantSpec{
			ProviderUIDs: []string{testProviderUID},
		},
	})
	party, prob = providerOnly.ResolveSuspendResume(cloudmodel.ParticipationActionResume, testPlatformUID, testProviderUID)
	if prob != nil || party != cloudmodel.HoldPartyProvider {
		t.Fatalf("provider-only: party=%q prob=%#v", party, prob)
	}

	zero := cloudmodel.NewBootstrapGrantResolver(cloudmodel.BootstrapGrantConfig{PrincipalID: "bootstrap"})
	_, prob = zero.ResolveSuspendResume(cloudmodel.ParticipationActionSuspend, testPlatformUID, testProviderUID)
	if prob == nil || prob.Code != apiproblem.CodeAuthorizationDenied {
		t.Fatalf("zero matches: %#v", prob)
	}

	both := cloudmodel.NewBootstrapGrantResolver(cloudmodel.BootstrapGrantConfig{
		PrincipalID: "bootstrap",
		SuspendResume: cloudmodel.SuspendResumeGrantSpec{
			PlatformUIDs: []string{testPlatformUID},
			ProviderUIDs: []string{testProviderUID},
		},
	})
	_, prob = both.ResolveSuspendResume(cloudmodel.ParticipationActionSuspend, testPlatformUID, testProviderUID)
	if prob == nil || prob.Code != apiproblem.CodeAuthorizationDenied {
		t.Fatalf("multiple matches: %#v", prob)
	}
}

func TestBootstrapGrants_NoPersistedRolesMembershipsOrAssignments(t *testing.T) {
	t.Parallel()
	r := testBootstrapResolver(t)
	if r.PersistsRolesMembershipsOrAssignments() {
		t.Fatal("FEATURE-0015 must not persist roles/memberships/assignments")
	}

	typ := reflect.TypeOf(*r)
	for i := 0; i < typ.NumField(); i++ {
		name := typ.Field(i).Name
		switch name {
		case "principalID", "grants":
			// allowed process-local fields
		default:
			t.Fatalf("unexpected resolver field %q (roles/memberships/assignments forbidden)", name)
		}
	}

	// Grants() returns a defensive copy.
	gs := r.Grants()
	if len(gs) == 0 {
		t.Fatal("expected grants")
	}
	gs[0].Action = "mutated"
	if r.Grants()[0].Action == "mutated" {
		t.Fatal("Grants() must return a copy")
	}
}
