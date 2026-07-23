package apivalid

import (
	"context"
	"errors"
	"net/http"
	"sync/atomic"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/apiref"
)

// Compile-time StructuralValidator conformance for pipeline stubs.
var (
	_ StructuralValidator           = stubStructuralValidator{}
	_ StructuralValidator           = trackingStructuralValidator{}
	_ ScopeAuthorizer               = trackingScopeAuthorizer{}
	_ AuthorizedTargetScopeResolver = trackingTargetScopeResolver{}
)

type stubStructuralValidator struct {
	violations []apiproblem.Violation
	err        error
}

func (s stubStructuralValidator) Validate(_ any, _ string) ([]apiproblem.Violation, error) {
	if s.err != nil {
		return nil, s.err
	}
	if len(s.violations) == 0 {
		return nil, nil
	}
	out := make([]apiproblem.Violation, len(s.violations))
	copy(out, s.violations)
	return out, nil
}

type trackingStructuralValidator struct {
	inner stubStructuralValidator
	calls *atomic.Int32
}

func (s trackingStructuralValidator) Validate(instance any, schemaID string) ([]apiproblem.Violation, error) {
	if s.calls != nil {
		s.calls.Add(1)
	}
	return s.inner.Validate(instance, schemaID)
}

type trackingScopeAuthorizer struct {
	decision Decision
	calls    *atomic.Int32
}

func (s trackingScopeAuthorizer) Authorize(
	_ context.Context,
	_ CallerContext,
	_ apiref.TypedRef,
	_ apimeta.ScopeIdentity,
) Decision {
	if s.calls != nil {
		s.calls.Add(1)
	}
	return s.decision
}

type trackingTargetScopeResolver struct {
	scope     apimeta.ScopeIdentity
	available bool
	calls     *atomic.Int32
}

func (s trackingTargetScopeResolver) ResolveAuthorizedTargetScope(
	_ context.Context,
	_ CallerContext,
	_ apiref.TypedRef,
) (apimeta.ScopeIdentity, bool) {
	if s.calls != nil {
		s.calls.Add(1)
	}
	return s.scope, s.available
}

func platformScope() apimeta.ScopeIdentity {
	return apimeta.ScopeIdentity{Kind: apimeta.ScopePlatform, UID: apimeta.PlatformScopeUID}
}

func tenantScope(uid string) apimeta.ScopeIdentity {
	return apimeta.ScopeIdentity{Kind: apimeta.ScopeTenant, UID: uid}
}

func sampleTargetRef() apiref.TypedRef {
	return apiref.TypedRef{
		APIVersion: "core.sovrunn.io/v1alpha1",
		Kind:       "Project",
		Name:       "demo",
		UID:        "target-uid-1",
	}
}

func assertInternalErrorAt(t *testing.T, res Result, want Layer) {
	t.Helper()
	if res.FailedAt != want {
		t.Fatalf("FailedAt = %v, want %v", res.FailedAt, want)
	}
	if res.Problem == nil {
		t.Fatal("Problem is nil, want 500 INTERNAL_ERROR")
	}
	if res.Problem.Status != http.StatusInternalServerError {
		t.Fatalf("Problem.Status = %d, want %d", res.Problem.Status, http.StatusInternalServerError)
	}
	if res.Problem.Code != apiproblem.CodeInternalError {
		t.Fatalf("Problem.Code = %q, want %q", res.Problem.Code, apiproblem.CodeInternalError)
	}
	if res.Err == nil {
		t.Fatal("Err is nil, want internal cause")
	}
	if len(res.Violations) != 0 {
		t.Fatalf("Violations = %#v, want empty", res.Violations)
	}
}

func TestValidateStructuralNilValidator(t *testing.T) {
	t.Parallel()

	authCalls := &atomic.Int32{}
	op := platformScope()
	target := sampleTargetRef()
	ts := platformScope()
	caller := &CallerContext{}

	res := Validate(context.Background(), Input{
		Validator:      nil,
		Dst:            map[string]any{"kind": "Project"},
		OperationScope: &op,
		TargetRef:      &target,
		TargetScope:    &ts,
		Authorizer:     trackingScopeAuthorizer{decision: Allow, calls: authCalls},
		Caller:         caller,
	}, DefaultLimits())

	assertInternalErrorAt(t, res, LayerStructural)
	if res.Err.Error() != "nil validator" {
		t.Fatalf("Err = %v, want nil validator", res.Err)
	}
	if authCalls.Load() != 0 {
		t.Fatalf("Authorizer called %d times; layers 5–8 must not run after structural fail-closed", authCalls.Load())
	}
}

func TestValidateStructuralValidatorError(t *testing.T) {
	t.Parallel()

	authCalls := &atomic.Int32{}
	cause := errors.New("schema registry unavailable")
	op := platformScope()
	target := sampleTargetRef()
	ts := platformScope()

	res := Validate(context.Background(), Input{
		Validator:      stubStructuralValidator{err: cause},
		SchemaID:       "project.json",
		Dst:            map[string]any{"kind": "Project"},
		OperationScope: &op,
		TargetRef:      &target,
		TargetScope:    &ts,
		Authorizer:     trackingScopeAuthorizer{decision: Allow, calls: authCalls},
		Caller:         &CallerContext{},
	}, DefaultLimits())

	assertInternalErrorAt(t, res, LayerStructural)
	if !errors.Is(res.Err, cause) {
		t.Fatalf("Err = %v, want %v", res.Err, cause)
	}
	if authCalls.Load() != 0 {
		t.Fatalf("Authorizer called %d times after structural error", authCalls.Load())
	}
}

func TestValidateStructuralOrdinaryViolations(t *testing.T) {
	t.Parallel()

	authCalls := &atomic.Int32{}
	v := apiproblem.Violation{
		Field:   "/metadata/name",
		Code:    apiproblem.ViolationOutOfRange,
		Message: "name is required",
	}
	op := platformScope()
	target := sampleTargetRef()
	ts := platformScope()

	res := Validate(context.Background(), Input{
		Validator:      stubStructuralValidator{violations: []apiproblem.Violation{v}},
		Dst:            map[string]any{},
		OperationScope: &op,
		TargetRef:      &target,
		TargetScope:    &ts,
		Authorizer:     trackingScopeAuthorizer{decision: Allow, calls: authCalls},
		Caller:         &CallerContext{},
	}, DefaultLimits())

	if res.FailedAt != LayerStructural {
		t.Fatalf("FailedAt = %v, want %v", res.FailedAt, LayerStructural)
	}
	if res.Problem != nil {
		t.Fatalf("Problem = %#v, want nil for ordinary violations", res.Problem)
	}
	if res.Err != nil {
		t.Fatalf("Err = %v, want nil for ordinary violations", res.Err)
	}
	if len(res.Violations) != 1 || res.Violations[0] != v {
		t.Fatalf("Violations = %#v, want %#v", res.Violations, []apiproblem.Violation{v})
	}
	if authCalls.Load() != 0 {
		t.Fatalf("Authorizer called %d times; layers 5–8 must not run on structural violations", authCalls.Load())
	}
}

func TestValidateStructuralPassGenericStopsAfterLayer7(t *testing.T) {
	t.Parallel()

	calls := &atomic.Int32{}
	res := Validate(context.Background(), Input{
		Validator: trackingStructuralValidator{
			inner: stubStructuralValidator{},
			calls: calls,
		},
		Dst: map[string]any{"ok": true},
		// OperationScope nil → generic non-Operation path.
	}, DefaultLimits())

	if calls.Load() != 1 {
		t.Fatalf("structural Validate calls = %d, want 1", calls.Load())
	}
	if res.FailedAt != 0 {
		t.Fatalf("FailedAt = %v, want 0 (success)", res.FailedAt)
	}
	if res.Problem != nil || res.Err != nil || len(res.Violations) != 0 {
		t.Fatalf("unexpected failure result: %#v", res)
	}
}

func TestValidateLayer8InvalidConfigs(t *testing.T) {
	t.Parallel()

	op := platformScope()
	target := sampleTargetRef()
	ts := platformScope()
	caller := &CallerContext{}
	resolverCalls := &atomic.Int32{}
	authCalls := &atomic.Int32{}

	passing := stubStructuralValidator{}

	tests := []struct {
		name string
		in   Input
	}{
		{
			name: "neither TargetScope nor TargetScopeResolver",
			in: Input{
				Validator:      passing,
				OperationScope: &op,
				TargetRef:      &target,
				Caller:         caller,
			},
		},
		{
			name: "both TargetScope and TargetScopeResolver",
			in: Input{
				Validator:           passing,
				OperationScope:      &op,
				TargetRef:           &target,
				TargetScope:         &ts,
				TargetScopeResolver: trackingTargetScopeResolver{scope: ts, available: true, calls: resolverCalls},
				Authorizer:          trackingScopeAuthorizer{decision: Allow, calls: authCalls},
				Caller:              caller,
			},
		},
		{
			name: "missing Caller",
			in: Input{
				Validator:      passing,
				OperationScope: &op,
				TargetRef:      &target,
				TargetScope:    &ts,
				Authorizer:     trackingScopeAuthorizer{decision: Allow, calls: authCalls},
			},
		},
		{
			name: "missing TargetRef",
			in: Input{
				Validator:      passing,
				OperationScope: &op,
				TargetScope:    &ts,
				Authorizer:     trackingScopeAuthorizer{decision: Allow, calls: authCalls},
				Caller:         caller,
			},
		},
		{
			name: "incomplete Path A Authorizer nil",
			in: Input{
				Validator:      passing,
				OperationScope: &op,
				TargetRef:      &target,
				TargetScope:    &ts,
				Caller:         caller,
			},
		},
		{
			name: "incomplete Path B Caller nil",
			in: Input{
				Validator:           passing,
				OperationScope:      &op,
				TargetRef:           &target,
				TargetScopeResolver: trackingTargetScopeResolver{scope: ts, available: true, calls: resolverCalls},
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			beforeAuth := authCalls.Load()
			beforeResolver := resolverCalls.Load()

			res := Validate(context.Background(), tc.in, DefaultLimits())
			assertInternalErrorAt(t, res, LayerAuthorization)

			if authCalls.Load() != beforeAuth {
				t.Fatalf("Authorizer was called on invalid config %q", tc.name)
			}
			if resolverCalls.Load() != beforeResolver {
				t.Fatalf("TargetScopeResolver was called on invalid config %q", tc.name)
			}
		})
	}
}

func TestValidateLayer8PathAAllowMatch(t *testing.T) {
	t.Parallel()

	op := tenantScope("tenant-1")
	target := sampleTargetRef()
	ts := tenantScope("tenant-1")
	authCalls := &atomic.Int32{}

	res := Validate(context.Background(), Input{
		Validator:      stubStructuralValidator{},
		OperationScope: &op,
		TargetRef:      &target,
		TargetScope:    &ts,
		Authorizer:     trackingScopeAuthorizer{decision: Allow, calls: authCalls},
		Caller:         &CallerContext{Scopes: []apimeta.ScopeIdentity{ts}},
	}, DefaultLimits())

	if authCalls.Load() != 1 {
		t.Fatalf("Authorizer calls = %d, want 1", authCalls.Load())
	}
	if res.FailedAt != 0 || res.Problem != nil || res.Err != nil || len(res.Violations) != 0 {
		t.Fatalf("unexpected result: %#v", res)
	}
}

func TestValidateLayer8PathAAllowMismatch(t *testing.T) {
	t.Parallel()

	op := tenantScope("tenant-1")
	target := sampleTargetRef()
	ts := tenantScope("tenant-2")

	res := Validate(context.Background(), Input{
		Validator:      stubStructuralValidator{},
		OperationScope: &op,
		TargetRef:      &target,
		TargetScope:    &ts,
		Authorizer:     stubScopeAuthorizer{decision: Allow},
		Caller:         &CallerContext{},
	}, DefaultLimits())

	if res.FailedAt != LayerAuthorization {
		t.Fatalf("FailedAt = %v, want %v", res.FailedAt, LayerAuthorization)
	}
	if res.Problem != nil || res.Err != nil {
		t.Fatalf("Problem/Err must be nil for ordinary mismatch: Problem=%v Err=%v", res.Problem, res.Err)
	}
	if len(res.Violations) != 1 {
		t.Fatalf("Violations len = %d, want 1", len(res.Violations))
	}
	if res.Violations[0].Code != apiproblem.ViolationOperationTargetScopeMismatch {
		t.Fatalf("code = %q, want %q", res.Violations[0].Code, apiproblem.ViolationOperationTargetScopeMismatch)
	}
	if res.Violations[0].Field != "/metadata/scopeRef" {
		t.Fatalf("field = %q, want /metadata/scopeRef", res.Violations[0].Field)
	}
}

func TestValidateLayer8PathADenyNotDisclosed(t *testing.T) {
	t.Parallel()

	op := tenantScope("tenant-1")
	target := sampleTargetRef()
	ts := tenantScope("tenant-1")

	res := Validate(context.Background(), Input{
		Validator:      stubStructuralValidator{},
		OperationScope: &op,
		TargetRef:      &target,
		TargetScope:    &ts,
		Authorizer:     stubScopeAuthorizer{decision: DenyNotDisclosed},
		Caller:         &CallerContext{},
	}, DefaultLimits())

	if res.FailedAt != LayerAuthorization {
		t.Fatalf("FailedAt = %v, want %v", res.FailedAt, LayerAuthorization)
	}
	if res.Problem == nil {
		t.Fatal("Problem is nil, want SafeDenial 404")
	}
	if res.Problem.Status != http.StatusNotFound || res.Problem.Code != apiproblem.CodeResourceNotFound {
		t.Fatalf("Problem = status %d code %q, want 404 RESOURCE_NOT_FOUND", res.Problem.Status, res.Problem.Code)
	}
	if len(res.Violations) != 0 {
		t.Fatalf("Violations must be empty on SafeDenial, got %#v", res.Violations)
	}
}

func TestValidateLayer8PathBAvailableMatch(t *testing.T) {
	t.Parallel()

	op := platformScope()
	target := sampleTargetRef()
	ts := platformScope()
	resolverCalls := &atomic.Int32{}

	res := Validate(context.Background(), Input{
		Validator:      stubStructuralValidator{},
		OperationScope: &op,
		TargetRef:      &target,
		TargetScopeResolver: trackingTargetScopeResolver{
			scope:     ts,
			available: true,
			calls:     resolverCalls,
		},
		Caller: &CallerContext{},
	}, DefaultLimits())

	if resolverCalls.Load() != 1 {
		t.Fatalf("resolver calls = %d, want 1", resolverCalls.Load())
	}
	if res.FailedAt != 0 || res.Problem != nil || res.Err != nil || len(res.Violations) != 0 {
		t.Fatalf("unexpected result: %#v", res)
	}
}

func TestValidateLayer8PathBUnavailableSafeDenial(t *testing.T) {
	t.Parallel()

	op := platformScope()
	target := sampleTargetRef()

	res := Validate(context.Background(), Input{
		Validator:      stubStructuralValidator{},
		OperationScope: &op,
		TargetRef:      &target,
		TargetScopeResolver: stubAuthorizedTargetScopeResolver{
			available: false,
		},
		Caller: &CallerContext{},
	}, DefaultLimits())

	if res.FailedAt != LayerAuthorization {
		t.Fatalf("FailedAt = %v, want %v", res.FailedAt, LayerAuthorization)
	}
	if res.Problem == nil || res.Problem.Status != http.StatusNotFound || res.Problem.Code != apiproblem.CodeResourceNotFound {
		t.Fatalf("Problem = %#v, want SafeDenial 404 RESOURCE_NOT_FOUND", res.Problem)
	}
	if len(res.Violations) != 0 {
		t.Fatalf("Violations must be empty; mismatch must not be disclosed, got %#v", res.Violations)
	}
}
