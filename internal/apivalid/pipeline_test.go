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

// Compile-time StructuralValidator / stage conformance for pipeline stubs.
var (
	_ StructuralValidator           = stubStructuralValidator{}
	_ StructuralValidator           = trackingStructuralValidator{}
	_ ScopeAuthorizer               = trackingScopeAuthorizer{}
	_ AuthorizedTargetScopeResolver = trackingTargetScopeResolver{}
	_ DefaultingStage               = trackingDefaultingStage{}
	_ ValidationStage               = trackingValidationStage{}
	_ ValidationStage               = capturingValidationStage{}
)

// trackingDefaultingStage is a deterministic stage stub for pipeline tests.
// A zero value is an explicit no-op that returns the input object unchanged.
type trackingDefaultingStage struct {
	calls *atomic.Int32
	err   error
	out   any // when non-nil, returned instead of the input object
}

func (s trackingDefaultingStage) Apply(_ context.Context, object any) (any, error) {
	if s.calls != nil {
		s.calls.Add(1)
	}
	if s.err != nil {
		return nil, s.err
	}
	if s.out != nil {
		return s.out, nil
	}
	return object, nil
}

// trackingValidationStage is a deterministic semantic/reference stage stub.
// A zero value is an explicit no-op that returns no violations.
type trackingValidationStage struct {
	calls      *atomic.Int32
	err        error
	violations []apiproblem.Violation
}

func (s trackingValidationStage) Validate(_ context.Context, _ any) ([]apiproblem.Violation, error) {
	if s.calls != nil {
		s.calls.Add(1)
	}
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

// passThroughStages returns an explicit StageSet of deterministic no-op
// stages so layer-8 and structural-pass tests can reach later layers
// without silently omitting stages 5–7.
func passThroughStages() StageSet {
	return StageSet{
		Defaulting: trackingDefaultingStage{},
		Semantic:   trackingValidationStage{},
		Reference:  trackingValidationStage{},
	}
}

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

// trackingStages returns a StageSet with call counters for layers 5–7 so
// structural fail-closed tests can prove those layers were not executed.
func trackingStages(def, sem, ref *atomic.Int32) StageSet {
	return StageSet{
		Defaulting: trackingDefaultingStage{calls: def},
		Semantic:   trackingValidationStage{calls: sem},
		Reference:  trackingValidationStage{calls: ref},
	}
}

func assertLayers57NotExecuted(t *testing.T, def, sem, ref *atomic.Int32) {
	t.Helper()
	if def.Load() != 0 || sem.Load() != 0 || ref.Load() != 0 {
		t.Fatalf("layers 5–7 must not execute after structural fail-closed: defaulting=%d semantic=%d reference=%d",
			def.Load(), sem.Load(), ref.Load())
	}
}

func TestValidateStructuralNilValidator(t *testing.T) {
	t.Parallel()

	defCalls := &atomic.Int32{}
	semCalls := &atomic.Int32{}
	refCalls := &atomic.Int32{}
	authCalls := &atomic.Int32{}
	op := platformScope()
	target := sampleTargetRef()
	ts := platformScope()
	caller := &CallerContext{}

	res := Validate(context.Background(), Input{
		Validator:      nil,
		Dst:            map[string]any{"kind": "Project"},
		Stages:         trackingStages(defCalls, semCalls, refCalls),
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
	assertLayers57NotExecuted(t, defCalls, semCalls, refCalls)
	if authCalls.Load() != 0 {
		t.Fatalf("Authorizer called %d times; layer 8 must not run after structural fail-closed", authCalls.Load())
	}
}

func TestValidateStructuralValidatorError(t *testing.T) {
	t.Parallel()

	defCalls := &atomic.Int32{}
	semCalls := &atomic.Int32{}
	refCalls := &atomic.Int32{}
	authCalls := &atomic.Int32{}
	cause := errors.New("schema registry unavailable")
	op := platformScope()
	target := sampleTargetRef()
	ts := platformScope()

	res := Validate(context.Background(), Input{
		Validator:      stubStructuralValidator{err: cause},
		SchemaID:       "project.json",
		Dst:            map[string]any{"kind": "Project"},
		Stages:         trackingStages(defCalls, semCalls, refCalls),
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
	assertLayers57NotExecuted(t, defCalls, semCalls, refCalls)
	if authCalls.Load() != 0 {
		t.Fatalf("Authorizer called %d times after structural error", authCalls.Load())
	}
}

func TestValidateStructuralOrdinaryViolations(t *testing.T) {
	t.Parallel()

	defCalls := &atomic.Int32{}
	semCalls := &atomic.Int32{}
	refCalls := &atomic.Int32{}
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
		Stages:         trackingStages(defCalls, semCalls, refCalls),
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
	assertLayers57NotExecuted(t, defCalls, semCalls, refCalls)
	if authCalls.Load() != 0 {
		t.Fatalf("Authorizer called %d times; layers 5–8 must not run on structural violations", authCalls.Load())
	}
}

func TestValidateStructuralPassGenericStopsAfterLayer7(t *testing.T) {
	t.Parallel()

	structCalls := &atomic.Int32{}
	defCalls := &atomic.Int32{}
	semCalls := &atomic.Int32{}
	refCalls := &atomic.Int32{}

	res := Validate(context.Background(), Input{
		Validator: trackingStructuralValidator{
			inner: stubStructuralValidator{},
			calls: structCalls,
		},
		Dst:    map[string]any{"ok": true},
		Stages: trackingStages(defCalls, semCalls, refCalls),
		// OperationScope nil → generic non-Operation path.
	}, DefaultLimits())

	if structCalls.Load() != 1 {
		t.Fatalf("structural Validate calls = %d, want 1", structCalls.Load())
	}
	if defCalls.Load() != 1 || semCalls.Load() != 1 || refCalls.Load() != 1 {
		t.Fatalf("clean structural pass must invoke layers 5–7: defaulting=%d semantic=%d reference=%d",
			defCalls.Load(), semCalls.Load(), refCalls.Load())
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
				Stages:         passThroughStages(),
				OperationScope: &op,
				TargetRef:      &target,
				Caller:         caller,
			},
		},
		{
			name: "both TargetScope and TargetScopeResolver",
			in: Input{
				Validator:           passing,
				Stages:              passThroughStages(),
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
				Stages:         passThroughStages(),
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
				Stages:         passThroughStages(),
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
				Stages:         passThroughStages(),
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
				Stages:              passThroughStages(),
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
		Stages:         passThroughStages(),
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
		Stages:         passThroughStages(),
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
		Stages:         passThroughStages(),
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
		Stages:         passThroughStages(),
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
		Stages:         passThroughStages(),
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

// --- Task 6.5e: layers 5–7 stage invocation ---

func TestValidateNilRequiredDefaultingStage(t *testing.T) {
	t.Parallel()

	semCalls := &atomic.Int32{}
	refCalls := &atomic.Int32{}

	res := Validate(context.Background(), Input{
		Validator: stubStructuralValidator{},
		Dst:       map[string]any{"ok": true},
		Stages: StageSet{
			Defaulting: nil,
			Semantic:   trackingValidationStage{calls: semCalls},
			Reference:  trackingValidationStage{calls: refCalls},
		},
	}, DefaultLimits())

	assertInternalErrorAt(t, res, LayerDefaulting)
	if res.Err.Error() != "nil defaulting stage" {
		t.Fatalf("Err = %v, want nil defaulting stage", res.Err)
	}
	if semCalls.Load() != 0 || refCalls.Load() != 0 {
		t.Fatalf("later stages ran: semantic=%d reference=%d", semCalls.Load(), refCalls.Load())
	}
}

func TestValidateNilRequiredSemanticStage(t *testing.T) {
	t.Parallel()

	defCalls := &atomic.Int32{}
	refCalls := &atomic.Int32{}

	res := Validate(context.Background(), Input{
		Validator: stubStructuralValidator{},
		Dst:       map[string]any{"ok": true},
		Stages: StageSet{
			Defaulting: trackingDefaultingStage{calls: defCalls},
			Semantic:   nil,
			Reference:  trackingValidationStage{calls: refCalls},
		},
	}, DefaultLimits())

	assertInternalErrorAt(t, res, LayerSemantic)
	if res.Err.Error() != "nil semantic stage" {
		t.Fatalf("Err = %v, want nil semantic stage", res.Err)
	}
	if defCalls.Load() != 1 {
		t.Fatalf("defaulting calls = %d, want 1", defCalls.Load())
	}
	if refCalls.Load() != 0 {
		t.Fatalf("reference called %d times after nil semantic", refCalls.Load())
	}
}

func TestValidateNilRequiredReferenceStage(t *testing.T) {
	t.Parallel()

	defCalls := &atomic.Int32{}
	semCalls := &atomic.Int32{}

	res := Validate(context.Background(), Input{
		Validator: stubStructuralValidator{},
		Dst:       map[string]any{"ok": true},
		Stages: StageSet{
			Defaulting: trackingDefaultingStage{calls: defCalls},
			Semantic:   trackingValidationStage{calls: semCalls},
			Reference:  nil,
		},
	}, DefaultLimits())

	assertInternalErrorAt(t, res, LayerReference)
	if res.Err.Error() != "nil reference stage" {
		t.Fatalf("Err = %v, want nil reference stage", res.Err)
	}
	if defCalls.Load() != 1 || semCalls.Load() != 1 {
		t.Fatalf("earlier stage calls: defaulting=%d semantic=%d", defCalls.Load(), semCalls.Load())
	}
}

func TestValidateDefaultingStageErrorStopsLaterLayers(t *testing.T) {
	t.Parallel()

	cause := errors.New("defaulting fault")
	semCalls := &atomic.Int32{}
	refCalls := &atomic.Int32{}

	res := Validate(context.Background(), Input{
		Validator: stubStructuralValidator{},
		Dst:       map[string]any{"ok": true},
		Stages: StageSet{
			Defaulting: trackingDefaultingStage{err: cause},
			Semantic:   trackingValidationStage{calls: semCalls},
			Reference:  trackingValidationStage{calls: refCalls},
		},
	}, DefaultLimits())

	assertInternalErrorAt(t, res, LayerDefaulting)
	if !errors.Is(res.Err, cause) {
		t.Fatalf("Err = %v, want %v", res.Err, cause)
	}
	if semCalls.Load() != 0 || refCalls.Load() != 0 {
		t.Fatalf("later stages ran after defaulting error: semantic=%d reference=%d", semCalls.Load(), refCalls.Load())
	}
}

func TestValidateSemanticViolationsStopBeforeReference(t *testing.T) {
	t.Parallel()

	v := apiproblem.Violation{
		Field:   "/metadata/name",
		Code:    ViolationInvalidResourceName,
		Message: "invalid name",
	}
	defCalls := &atomic.Int32{}
	semCalls := &atomic.Int32{}
	refCalls := &atomic.Int32{}
	authCalls := &atomic.Int32{}
	op := platformScope()
	target := sampleTargetRef()
	ts := platformScope()

	res := Validate(context.Background(), Input{
		Validator: stubStructuralValidator{},
		Dst:       map[string]any{"ok": true},
		Stages: StageSet{
			Defaulting: trackingDefaultingStage{calls: defCalls},
			Semantic: trackingValidationStage{
				calls:      semCalls,
				violations: []apiproblem.Violation{v},
			},
			Reference: trackingValidationStage{calls: refCalls},
		},
		OperationScope: &op,
		TargetRef:      &target,
		TargetScope:    &ts,
		Authorizer:     trackingScopeAuthorizer{decision: Allow, calls: authCalls},
		Caller:         &CallerContext{},
	}, DefaultLimits())

	if res.FailedAt != LayerSemantic {
		t.Fatalf("FailedAt = %v, want %v", res.FailedAt, LayerSemantic)
	}
	if res.Problem != nil || res.Err != nil {
		t.Fatalf("Problem/Err must be nil for ordinary violations: Problem=%v Err=%v", res.Problem, res.Err)
	}
	if len(res.Violations) != 1 || res.Violations[0] != v {
		t.Fatalf("Violations = %#v, want %#v", res.Violations, []apiproblem.Violation{v})
	}
	if defCalls.Load() != 1 || semCalls.Load() != 1 {
		t.Fatalf("calls: defaulting=%d semantic=%d", defCalls.Load(), semCalls.Load())
	}
	if refCalls.Load() != 0 {
		t.Fatalf("reference called %d times; must stop at LayerSemantic", refCalls.Load())
	}
	if authCalls.Load() != 0 {
		t.Fatalf("layer-8 Authorizer called %d times after semantic violations", authCalls.Load())
	}
}

func TestValidateReferenceViolationsStopAtLayerReference(t *testing.T) {
	t.Parallel()

	v := apiproblem.Violation{
		Field:   "/spec/targetRef",
		Code:    apiproblem.ViolationCode(apiref.CodeKindNotAllowed),
		Message: "kind not allowed",
	}
	defCalls := &atomic.Int32{}
	semCalls := &atomic.Int32{}
	refCalls := &atomic.Int32{}
	authCalls := &atomic.Int32{}
	op := platformScope()
	target := sampleTargetRef()
	ts := platformScope()

	res := Validate(context.Background(), Input{
		Validator: stubStructuralValidator{},
		Dst:       map[string]any{"ok": true},
		Stages: StageSet{
			Defaulting: trackingDefaultingStage{calls: defCalls},
			Semantic:   trackingValidationStage{calls: semCalls},
			Reference: trackingValidationStage{
				calls:      refCalls,
				violations: []apiproblem.Violation{v},
			},
		},
		OperationScope: &op,
		TargetRef:      &target,
		TargetScope:    &ts,
		Authorizer:     trackingScopeAuthorizer{decision: Allow, calls: authCalls},
		Caller:         &CallerContext{},
	}, DefaultLimits())

	if res.FailedAt != LayerReference {
		t.Fatalf("FailedAt = %v, want %v", res.FailedAt, LayerReference)
	}
	if res.Problem != nil || res.Err != nil {
		t.Fatalf("Problem/Err must be nil for ordinary violations: Problem=%v Err=%v", res.Problem, res.Err)
	}
	if len(res.Violations) != 1 || res.Violations[0] != v {
		t.Fatalf("Violations = %#v, want %#v", res.Violations, []apiproblem.Violation{v})
	}
	if defCalls.Load() != 1 || semCalls.Load() != 1 || refCalls.Load() != 1 {
		t.Fatalf("calls: defaulting=%d semantic=%d reference=%d", defCalls.Load(), semCalls.Load(), refCalls.Load())
	}
	if authCalls.Load() != 0 {
		t.Fatalf("layer-8 Authorizer called %d times after reference violations", authCalls.Load())
	}
}

func TestValidateNoOpStagesPassThroughAndAreInvoked(t *testing.T) {
	t.Parallel()

	defCalls := &atomic.Int32{}
	semCalls := &atomic.Int32{}
	refCalls := &atomic.Int32{}
	defaultedMarker := &struct{ label string }{label: "defaulted"}
	var gotSemantic any
	var gotReference any

	sem := capturingValidationStage{calls: semCalls, capture: &gotSemantic}
	ref := capturingValidationStage{calls: refCalls, capture: &gotReference}

	res := Validate(context.Background(), Input{
		Validator: stubStructuralValidator{},
		Dst:       map[string]any{"raw": true},
		Stages: StageSet{
			Defaulting: trackingDefaultingStage{calls: defCalls, out: defaultedMarker},
			Semantic:   sem,
			Reference:  ref,
		},
	}, DefaultLimits())

	if res.FailedAt != 0 || res.Problem != nil || res.Err != nil || len(res.Violations) != 0 {
		t.Fatalf("unexpected failure: %#v", res)
	}
	if defCalls.Load() != 1 || semCalls.Load() != 1 || refCalls.Load() != 1 {
		t.Fatalf("no-op stages must still be invoked: defaulting=%d semantic=%d reference=%d",
			defCalls.Load(), semCalls.Load(), refCalls.Load())
	}
	if gotSemantic != defaultedMarker {
		t.Fatalf("semantic received %#v, want defaulted object", gotSemantic)
	}
	if gotReference != defaultedMarker {
		t.Fatalf("reference received %#v, want defaulted object", gotReference)
	}
}

func TestValidateSemanticStageErrorStopsBeforeReference(t *testing.T) {
	t.Parallel()

	cause := errors.New("semantic fault")
	refCalls := &atomic.Int32{}

	res := Validate(context.Background(), Input{
		Validator: stubStructuralValidator{},
		Dst:       map[string]any{"ok": true},
		Stages: StageSet{
			Defaulting: trackingDefaultingStage{},
			Semantic:   trackingValidationStage{err: cause},
			Reference:  trackingValidationStage{calls: refCalls},
		},
	}, DefaultLimits())

	assertInternalErrorAt(t, res, LayerSemantic)
	if !errors.Is(res.Err, cause) {
		t.Fatalf("Err = %v, want %v", res.Err, cause)
	}
	if refCalls.Load() != 0 {
		t.Fatalf("reference called %d times after semantic error", refCalls.Load())
	}
}

func TestValidateReferenceStageError(t *testing.T) {
	t.Parallel()

	cause := errors.New("reference fault")
	authCalls := &atomic.Int32{}
	op := platformScope()
	target := sampleTargetRef()
	ts := platformScope()

	res := Validate(context.Background(), Input{
		Validator: stubStructuralValidator{},
		Dst:       map[string]any{"ok": true},
		Stages: StageSet{
			Defaulting: trackingDefaultingStage{},
			Semantic:   trackingValidationStage{},
			Reference:  trackingValidationStage{err: cause},
		},
		OperationScope: &op,
		TargetRef:      &target,
		TargetScope:    &ts,
		Authorizer:     trackingScopeAuthorizer{decision: Allow, calls: authCalls},
		Caller:         &CallerContext{},
	}, DefaultLimits())

	assertInternalErrorAt(t, res, LayerReference)
	if !errors.Is(res.Err, cause) {
		t.Fatalf("Err = %v, want %v", res.Err, cause)
	}
	if authCalls.Load() != 0 {
		t.Fatalf("layer-8 ran after reference error")
	}
}

// capturingValidationStage records the object passed to Validate so tests
// can assert Defaulting output is forwarded to Semantic and Reference.
type capturingValidationStage struct {
	calls   *atomic.Int32
	capture *any
	err     error
}

func (s capturingValidationStage) Validate(_ context.Context, object any) ([]apiproblem.Violation, error) {
	if s.calls != nil {
		s.calls.Add(1)
	}
	if s.capture != nil {
		*s.capture = object
	}
	if s.err != nil {
		return nil, s.err
	}
	return nil, nil
}
