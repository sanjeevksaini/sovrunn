package apivalid

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/apiref"
)

// Compile-time interface conformance for stub adopter types used in tests.
var (
	_ ScopeAuthorizer               = stubScopeAuthorizer{}
	_ AuthorizedResolver            = stubAuthorizedResolver{}
	_ AuthorizedTargetScopeResolver = stubAuthorizedTargetScopeResolver{}
)

type stubScopeAuthorizer struct {
	decision Decision
}

func (s stubScopeAuthorizer) Authorize(
	_ context.Context,
	_ CallerContext,
	_ apiref.TypedRef,
	_ apimeta.ScopeIdentity,
) Decision {
	return s.decision
}

type stubAuthorizedResolver struct {
	obj   any
	found bool
}

func (s stubAuthorizedResolver) Resolve(
	_ context.Context,
	_ CallerContext,
	_ apiref.TypedRef,
) (any, bool) {
	return s.obj, s.found
}

type stubAuthorizedTargetScopeResolver struct {
	scope     apimeta.ScopeIdentity
	available bool
}

func (s stubAuthorizedTargetScopeResolver) ResolveAuthorizedTargetScope(
	_ context.Context,
	_ CallerContext,
	_ apiref.TypedRef,
) (apimeta.ScopeIdentity, bool) {
	return s.scope, s.available
}

func TestSafeDenialDenyNotDisclosedMatchesAbsent(t *testing.T) {
	t.Parallel()

	// Canonical "absent" outcome uses the same Problem constructor as
	// DenyNotDisclosed so client-facing responses are byte-identical.
	absent := apiproblem.New(apiproblem.CodeResourceNotFound)
	denied := SafeDenial(DenyNotDisclosed)

	if denied == nil {
		t.Fatal("SafeDenial(DenyNotDisclosed) must not return nil")
	}
	if denied.Status != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", denied.Status, http.StatusNotFound)
	}
	if denied.Code != apiproblem.CodeResourceNotFound {
		t.Fatalf("code = %q, want %q", denied.Code, apiproblem.CodeResourceNotFound)
	}
	if denied.Title != apiproblem.TitleFor(apiproblem.CodeResourceNotFound) {
		t.Fatalf("title = %q, want %q", denied.Title, apiproblem.TitleFor(apiproblem.CodeResourceNotFound))
	}

	absentJSON, err := json.Marshal(absent)
	if err != nil {
		t.Fatalf("marshal absent: %v", err)
	}
	deniedJSON, err := json.Marshal(denied)
	if err != nil {
		t.Fatalf("marshal denied: %v", err)
	}
	if !bytes.Equal(absentJSON, deniedJSON) {
		t.Fatalf("DenyNotDisclosed and absent Problems differ:\nabsent=%s\ndenied=%s", absentJSON, deniedJSON)
	}

	// Repeated DenyNotDisclosed mappings remain byte-identical to each other.
	againJSON, err := json.Marshal(SafeDenial(DenyNotDisclosed))
	if err != nil {
		t.Fatalf("marshal second denial: %v", err)
	}
	if !bytes.Equal(deniedJSON, againJSON) {
		t.Fatalf("repeated DenyNotDisclosed responses differ:\nfirst=%s\nsecond=%s", deniedJSON, againJSON)
	}
}

func TestSafeDenialDenyKnownIs403(t *testing.T) {
	t.Parallel()

	got := SafeDenial(DenyKnown)
	if got == nil {
		t.Fatal("SafeDenial(DenyKnown) must not return nil")
	}
	if got.Status != http.StatusForbidden {
		t.Fatalf("status = %d, want %d", got.Status, http.StatusForbidden)
	}
	if got.Code != apiproblem.CodeAuthorizationDenied {
		t.Fatalf("code = %q, want %q", got.Code, apiproblem.CodeAuthorizationDenied)
	}
	if got.Title != apiproblem.TitleFor(apiproblem.CodeAuthorizationDenied) {
		t.Fatalf("title = %q, want %q", got.Title, apiproblem.TitleFor(apiproblem.CodeAuthorizationDenied))
	}
	if got.Detail != "" {
		t.Fatalf("detail must be empty to avoid leaking policy inputs; got %q", got.Detail)
	}
	if got.RequestID != "" {
		t.Fatalf("requestId must be unset by SafeDenial; got %q", got.RequestID)
	}
	if len(got.Violations) != 0 {
		t.Fatalf("violations must be empty; got %#v", got.Violations)
	}
}

func TestSafeDenialAllowReturnsNil(t *testing.T) {
	t.Parallel()

	if got := SafeDenial(Allow); got != nil {
		t.Fatalf("SafeDenial(Allow) = %#v, want nil", got)
	}
}

func TestSafeDenialUnknownDecisionFailsClosedWithoutDisclosure(t *testing.T) {
	t.Parallel()

	got := SafeDenial(Decision(99))
	want := SafeDenial(DenyNotDisclosed)
	gotJSON, err := json.Marshal(got)
	if err != nil {
		t.Fatalf("marshal unknown: %v", err)
	}
	wantJSON, err := json.Marshal(want)
	if err != nil {
		t.Fatalf("marshal want: %v", err)
	}
	if !bytes.Equal(gotJSON, wantJSON) {
		t.Fatalf("unknown Decision must map like DenyNotDisclosed:\ngot=%s\nwant=%s", gotJSON, wantJSON)
	}
}

func TestAuthorizedResolverUniformUnavailableOutcome(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	caller := CallerContext{
		Scopes: []apimeta.ScopeIdentity{{
			Kind: apimeta.ScopeTenant,
			UID:  "tenant-a",
		}},
	}
	target := apiref.TypedRef{
		APIVersion: "platform.sovrunn.io/v1alpha1",
		Kind:       "Project",
		Name:       "demo",
		UID:        "proj-1",
	}

	absent := stubAuthorizedResolver{obj: nil, found: false}
	unauthorized := stubAuthorizedResolver{obj: nil, found: false}

	objA, foundA := absent.Resolve(ctx, caller, target)
	objB, foundB := unauthorized.Resolve(ctx, caller, target)

	if foundA || foundB {
		t.Fatalf("both absent and unauthorized must report found=false; got %v / %v", foundA, foundB)
	}
	if objA != nil || objB != nil {
		t.Fatalf("uniform unavailable outcome must not leak objects; got %#v / %#v", objA, objB)
	}

	// Client-facing SafeDenial for the unavailable path is identical to absent.
	unavailableProblem := SafeDenial(DenyNotDisclosed)
	absentProblem := apiproblem.New(apiproblem.CodeResourceNotFound)
	uJSON, err := json.Marshal(unavailableProblem)
	if err != nil {
		t.Fatalf("marshal unavailable: %v", err)
	}
	aJSON, err := json.Marshal(absentProblem)
	if err != nil {
		t.Fatalf("marshal absent: %v", err)
	}
	if !bytes.Equal(uJSON, aJSON) {
		t.Fatalf("unavailable SafeDenial must match absent:\nunavail=%s\nabsent=%s", uJSON, aJSON)
	}
}

func TestAuthorizedTargetScopeResolverHidesAbsentVersusUnauthorized(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	caller := CallerContext{}
	target := apiref.TypedRef{Kind: "Project", Name: "x"}

	absent := stubAuthorizedTargetScopeResolver{available: false}
	unauthorized := stubAuthorizedTargetScopeResolver{available: false}

	scopeA, okA := absent.ResolveAuthorizedTargetScope(ctx, caller, target)
	scopeB, okB := unauthorized.ResolveAuthorizedTargetScope(ctx, caller, target)

	if okA || okB {
		t.Fatalf("available must be false for absent and unauthorized; got %v / %v", okA, okB)
	}
	if scopeA != (apimeta.ScopeIdentity{}) || scopeB != (apimeta.ScopeIdentity{}) {
		t.Fatalf("unavailable resolution must not return scope detail; got %#v / %#v", scopeA, scopeB)
	}

	problem := SafeDenial(DenyNotDisclosed)
	if problem.Status != http.StatusNotFound || problem.Code != apiproblem.CodeResourceNotFound {
		t.Fatalf("available=false must map through SafeDenial 404 RESOURCE_NOT_FOUND; got status=%d code=%q",
			problem.Status, problem.Code)
	}
}

func TestScopeAuthorizerStubDecisions(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	caller := CallerContext{}
	target := apiref.TypedRef{Kind: "Project", Name: "demo"}
	scope := apimeta.ScopeIdentity{Kind: apimeta.ScopeTenant, UID: "t1"}

	if got := (stubScopeAuthorizer{decision: Allow}).Authorize(ctx, caller, target, scope); got != Allow {
		t.Fatalf("Allow stub = %v, want Allow", got)
	}
	if got := (stubScopeAuthorizer{decision: DenyNotDisclosed}).Authorize(ctx, caller, target, scope); got != DenyNotDisclosed {
		t.Fatalf("DenyNotDisclosed stub = %v, want DenyNotDisclosed", got)
	}
	if got := (stubScopeAuthorizer{decision: DenyKnown}).Authorize(ctx, caller, target, scope); got != DenyKnown {
		t.Fatalf("DenyKnown stub = %v, want DenyKnown", got)
	}
}
