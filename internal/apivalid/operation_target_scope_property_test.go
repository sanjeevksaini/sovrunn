package apivalid

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"sync/atomic"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/apiref"
)

// Deterministic seed for Property 11 reproducibility
// (F12-SCOPE-002, F12-REF-001; D-17).
const property11Seed int64 = 20260723

const property11Iterations = 100

// property11ScopeKinds covers all six Matrix B / D-17 governance scopes.
var property11ScopeKinds = []apimeta.ScopeKind{
	apimeta.ScopePlatform,
	apimeta.ScopeOrganization,
	apimeta.ScopeOrganizationUnit,
	apimeta.ScopeTenant,
	apimeta.ScopeProject,
	apimeta.ScopeProvider,
}

// property11CaseClass selects the oracle bucket for one generated iteration.
type property11CaseClass int

const (
	property11PathAMatch property11CaseClass = iota
	property11PathAPlatformNilMatch
	property11PathANonPlatformUIDMatch
	property11PathAKindMismatch
	property11PathAUIDMismatch
	property11PathADenyNotDisclosed
	property11PathBMatch
	property11PathBKindMismatch
	property11PathBUIDMismatch
	property11PathBUnavailable
	property11PathBUnauthorized
	property11ConfigNeither
	property11ConfigBoth
	property11ConfigMissingCaller
	property11ConfigMissingTargetRef
	property11ConfigIncompletePathA
	property11ConfigIncompletePathB
	property11GenericNonOperation
)

const property11CaseClassCount = int(property11GenericNonOperation) + 1

// property11Case models one Operation target-scope equality scenario under
// Path A, Path B, an invalid layer-8 configuration, or the generic
// non-Operation path (nil OperationScope).
type property11Case struct {
	Class       property11CaseClass
	OpScope     apimeta.ScopeIdentity
	TargetScope apimeta.ScopeIdentity
	TargetRef   apiref.TypedRef
	Caller      CallerContext
}

// Feature: api-resource-naming-status-and-validation-standard, Property 11: Operation target-scope equality
//
// For any Operation, scopeRef MUST equal the resolved canonical governance
// scope of targetRef. Matching scopes (all six kinds, platform nil /
// PlatformScopeUID form, and matching non-platform UIDs) succeed on complete
// Path A (Allow) and Path B (available=true). Kind or UID mismatch yields
// OPERATION_TARGET_SCOPE_MISMATCH at /metadata/scopeRef. Unauthorized or
// unavailable targets map through SafeDenial 404 with no mismatch disclosed.
// Invalid layer-8 configurations fail closed with 500 INTERNAL_ERROR and no
// target lookup. Generic non-Operation validation (nil OperationScope) stops
// successfully after layer 7.
//
// Validates: Requirements 4.4, 4.5 (F12-SCOPE-002, F12-REF-001; D-17)
func TestProperty11_OperationTargetScopeEquality(t *testing.T) {
	t.Parallel()

	rng := rand.New(rand.NewSource(property11Seed))
	for i := 0; i < property11Iterations; i++ {
		c := generateProperty11Case(rng, i)
		if err := checkProperty11Case(c, i); err != nil {
			t.Fatalf("property 11 failed at iteration %d (seed %d): %v", i, property11Seed, err)
		}
	}
}

func generateProperty11Case(rng *rand.Rand, iteration int) property11Case {
	class := property11CaseClass(iteration % property11CaseClassCount)
	targetRef := property11RandomTargetRef(rng)
	caller := CallerContext{Scopes: []apimeta.ScopeIdentity{property11RandomScope(rng)}}

	switch class {
	case property11PathAMatch, property11PathBMatch:
		// Cycle all six scopes across iterations so each matching kind
		// appears often in the positive Path A / Path B buckets.
		op := property11ScopeByIndex(iteration / property11CaseClassCount)
		return property11Case{
			Class:       class,
			OpScope:     op,
			TargetScope: op,
			TargetRef:   targetRef,
			Caller:      caller,
		}

	case property11PathAPlatformNilMatch:
		op := apimeta.CanonicalScopeIdentity(nil) // {Platform, PlatformScopeUID}
		return property11Case{
			Class:       class,
			OpScope:     op,
			TargetScope: op,
			TargetRef:   targetRef,
			Caller:      caller,
		}

	case property11PathANonPlatformUIDMatch:
		op := property11RandomNonPlatformScope(rng)
		return property11Case{
			Class:       class,
			OpScope:     op,
			TargetScope: op,
			TargetRef:   targetRef,
			Caller:      caller,
		}

	case property11PathAKindMismatch, property11PathBKindMismatch:
		op := property11RandomScope(rng)
		target := property11DistinctKindScope(rng, op)
		return property11Case{
			Class:       class,
			OpScope:     op,
			TargetScope: target,
			TargetRef:   targetRef,
			Caller:      caller,
		}

	case property11PathAUIDMismatch, property11PathBUIDMismatch:
		op := property11RandomNonPlatformScope(rng)
		target := apimeta.ScopeIdentity{Kind: op.Kind, UID: property11DistinctUID(rng, op.UID)}
		return property11Case{
			Class:       class,
			OpScope:     op,
			TargetScope: target,
			TargetRef:   targetRef,
			Caller:      caller,
		}

	case property11PathADenyNotDisclosed:
		op := property11RandomScope(rng)
		// Latent mismatch must never leak through SafeDenial.
		target := property11DistinctScope(rng, op)
		return property11Case{
			Class:       class,
			OpScope:     op,
			TargetScope: target,
			TargetRef:   targetRef,
			Caller:      caller,
		}

	case property11PathBUnavailable, property11PathBUnauthorized:
		op := property11RandomScope(rng)
		target := property11DistinctScope(rng, op)
		return property11Case{
			Class:       class,
			OpScope:     op,
			TargetScope: target,
			TargetRef:   targetRef,
			Caller:      caller,
		}

	case property11ConfigNeither, property11ConfigBoth,
		property11ConfigMissingCaller, property11ConfigMissingTargetRef,
		property11ConfigIncompletePathA, property11ConfigIncompletePathB:
		op := property11RandomScope(rng)
		return property11Case{
			Class:       class,
			OpScope:     op,
			TargetScope: op,
			TargetRef:   targetRef,
			Caller:      caller,
		}

	case property11GenericNonOperation:
		return property11Case{
			Class:     class,
			TargetRef: targetRef,
			Caller:    caller,
		}

	default:
		return property11Case{Class: class, TargetRef: targetRef, Caller: caller}
	}
}

func property11ScopeByIndex(n int) apimeta.ScopeIdentity {
	kind := property11ScopeKinds[n%len(property11ScopeKinds)]
	if kind == apimeta.ScopePlatform {
		return apimeta.ScopeIdentity{Kind: apimeta.ScopePlatform, UID: apimeta.PlatformScopeUID}
	}
	return apimeta.ScopeIdentity{
		Kind: kind,
		UID:  fmt.Sprintf("%s-uid-%d", kind, n),
	}
}

func property11RandomScope(rng *rand.Rand) apimeta.ScopeIdentity {
	kind := property11ScopeKinds[rng.Intn(len(property11ScopeKinds))]
	if kind == apimeta.ScopePlatform {
		return apimeta.ScopeIdentity{Kind: apimeta.ScopePlatform, UID: apimeta.PlatformScopeUID}
	}
	return apimeta.ScopeIdentity{
		Kind: kind,
		UID:  fmt.Sprintf("%s-uid-%d", kind, rng.Intn(1_000_000)),
	}
}

func property11RandomNonPlatformScope(rng *rand.Rand) apimeta.ScopeIdentity {
	nonPlatform := property11ScopeKinds[1:] // exclude Platform
	kind := nonPlatform[rng.Intn(len(nonPlatform))]
	return apimeta.ScopeIdentity{
		Kind: kind,
		UID:  fmt.Sprintf("%s-uid-%d", kind, rng.Intn(1_000_000)),
	}
}

func property11DistinctScope(rng *rand.Rand, other apimeta.ScopeIdentity) apimeta.ScopeIdentity {
	for {
		got := property11RandomScope(rng)
		if got != other {
			return got
		}
	}
}

func property11DistinctKindScope(rng *rand.Rand, other apimeta.ScopeIdentity) apimeta.ScopeIdentity {
	for {
		got := property11RandomScope(rng)
		if got.Kind != other.Kind {
			return got
		}
	}
}

func property11DistinctUID(rng *rand.Rand, other string) string {
	for {
		uid := fmt.Sprintf("uid-%d", rng.Intn(1_000_000))
		if uid != other {
			return uid
		}
	}
}

func property11RandomTargetRef(rng *rand.Rand) apiref.TypedRef {
	kinds := []string{"Project", "ResourcePool", "PluginDefinition", "ServiceInstance"}
	return apiref.TypedRef{
		APIVersion: "core.sovrunn.io/v1alpha1",
		Kind:       kinds[rng.Intn(len(kinds))],
		Name:       fmt.Sprintf("target-%d", rng.Intn(100000)),
		UID:        fmt.Sprintf("uid-%d", rng.Intn(1_000_000)),
	}
}

func checkProperty11Case(c property11Case, iteration int) error {
	switch c.Class {
	case property11PathAMatch, property11PathAPlatformNilMatch, property11PathANonPlatformUIDMatch:
		return checkProperty11PathAMatch(c, iteration)
	case property11PathAKindMismatch, property11PathAUIDMismatch:
		return checkProperty11PathAMismatch(c, iteration)
	case property11PathADenyNotDisclosed:
		return checkProperty11PathADenyNotDisclosed(c, iteration)
	case property11PathBMatch:
		return checkProperty11PathBMatch(c, iteration)
	case property11PathBKindMismatch, property11PathBUIDMismatch:
		return checkProperty11PathBMismatch(c, iteration)
	case property11PathBUnavailable, property11PathBUnauthorized:
		return checkProperty11PathBUnavailable(c, iteration)
	case property11ConfigNeither, property11ConfigBoth,
		property11ConfigMissingCaller, property11ConfigMissingTargetRef,
		property11ConfigIncompletePathA, property11ConfigIncompletePathB:
		return checkProperty11ConfigFailure(c, iteration)
	case property11GenericNonOperation:
		return checkProperty11GenericNonOperation(c, iteration)
	default:
		return fmt.Errorf("iteration %d: unknown case class %d", iteration, c.Class)
	}
}

func checkProperty11PathAMatch(c property11Case, iteration int) error {
	op := c.OpScope
	ts := c.TargetScope
	target := c.TargetRef
	caller := c.Caller
	authCalls := &atomic.Int32{}

	res := Validate(context.Background(), Input{
		Validator:      stubStructuralValidator{},
		Stages:         passThroughStages(),
		OperationScope: &op,
		TargetRef:      &target,
		TargetScope:    &ts,
		Authorizer:     trackingScopeAuthorizer{decision: Allow, calls: authCalls},
		Caller:         &caller,
	}, DefaultLimits())

	if authCalls.Load() != 1 {
		return fmt.Errorf("iteration %d path A match: Authorizer calls = %d, want 1", iteration, authCalls.Load())
	}
	if v := CheckOperationTargetScopeMatch(op, ts); v != nil {
		return fmt.Errorf("iteration %d path A match: pure comparator unexpectedly returned %#v", iteration, v)
	}
	return property11AssertSuccess(res, iteration, "path A match")
}

func checkProperty11PathAMismatch(c property11Case, iteration int) error {
	op := c.OpScope
	ts := c.TargetScope
	target := c.TargetRef
	caller := c.Caller
	authCalls := &atomic.Int32{}

	res := Validate(context.Background(), Input{
		Validator:      stubStructuralValidator{},
		Stages:         passThroughStages(),
		OperationScope: &op,
		TargetRef:      &target,
		TargetScope:    &ts,
		Authorizer:     trackingScopeAuthorizer{decision: Allow, calls: authCalls},
		Caller:         &caller,
	}, DefaultLimits())

	if authCalls.Load() != 1 {
		return fmt.Errorf("iteration %d path A mismatch: Authorizer calls = %d, want 1", iteration, authCalls.Load())
	}
	if v := CheckOperationTargetScopeMatch(op, ts); v == nil {
		return fmt.Errorf("iteration %d path A mismatch: pure comparator returned nil for op=%#v target=%#v",
			iteration, op, ts)
	}
	return property11AssertMismatch(res, iteration, "path A mismatch")
}

func checkProperty11PathADenyNotDisclosed(c property11Case, iteration int) error {
	op := c.OpScope
	ts := c.TargetScope
	target := c.TargetRef
	caller := c.Caller
	authCalls := &atomic.Int32{}
	want := SafeDenial(DenyNotDisclosed)

	res := Validate(context.Background(), Input{
		Validator:      stubStructuralValidator{},
		Stages:         passThroughStages(),
		OperationScope: &op,
		TargetRef:      &target,
		TargetScope:    &ts,
		Authorizer:     trackingScopeAuthorizer{decision: DenyNotDisclosed, calls: authCalls},
		Caller:         &caller,
	}, DefaultLimits())

	if authCalls.Load() != 1 {
		return fmt.Errorf("iteration %d path A DenyNotDisclosed: Authorizer calls = %d, want 1", iteration, authCalls.Load())
	}
	return property11AssertSafeDenial(res, want, iteration, "path A DenyNotDisclosed")
}

func checkProperty11PathBMatch(c property11Case, iteration int) error {
	op := c.OpScope
	ts := c.TargetScope
	target := c.TargetRef
	caller := c.Caller
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
		Caller: &caller,
	}, DefaultLimits())

	if resolverCalls.Load() != 1 {
		return fmt.Errorf("iteration %d path B match: resolver calls = %d, want 1", iteration, resolverCalls.Load())
	}
	if v := CheckOperationTargetScopeMatch(op, ts); v != nil {
		return fmt.Errorf("iteration %d path B match: pure comparator unexpectedly returned %#v", iteration, v)
	}
	return property11AssertSuccess(res, iteration, "path B match")
}

func checkProperty11PathBMismatch(c property11Case, iteration int) error {
	op := c.OpScope
	ts := c.TargetScope
	target := c.TargetRef
	caller := c.Caller
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
		Caller: &caller,
	}, DefaultLimits())

	if resolverCalls.Load() != 1 {
		return fmt.Errorf("iteration %d path B mismatch: resolver calls = %d, want 1", iteration, resolverCalls.Load())
	}
	return property11AssertMismatch(res, iteration, "path B mismatch")
}

func checkProperty11PathBUnavailable(c property11Case, iteration int) error {
	// available=false is the uniform outcome for absent and unauthorized.
	// Latent TargetScope mismatch must not surface OPERATION_TARGET_SCOPE_MISMATCH.
	op := c.OpScope
	target := c.TargetRef
	caller := c.Caller
	resolverCalls := &atomic.Int32{}
	want := SafeDenial(DenyNotDisclosed)

	res := Validate(context.Background(), Input{
		Validator:      stubStructuralValidator{},
		Stages:         passThroughStages(),
		OperationScope: &op,
		TargetRef:      &target,
		TargetScopeResolver: trackingTargetScopeResolver{
			scope:     c.TargetScope, // would mismatch if disclosed
			available: false,
			calls:     resolverCalls,
		},
		Caller: &caller,
	}, DefaultLimits())

	if resolverCalls.Load() != 1 {
		return fmt.Errorf("iteration %d path B unavailable: resolver calls = %d, want 1", iteration, resolverCalls.Load())
	}
	label := "path B unavailable"
	if c.Class == property11PathBUnauthorized {
		label = "path B unauthorized"
	}
	return property11AssertSafeDenial(res, want, iteration, label)
}

func checkProperty11ConfigFailure(c property11Case, iteration int) error {
	op := c.OpScope
	ts := c.TargetScope
	target := c.TargetRef
	caller := c.Caller
	authCalls := &atomic.Int32{}
	resolverCalls := &atomic.Int32{}

	var in Input
	switch c.Class {
	case property11ConfigNeither:
		in = Input{
			Validator:      stubStructuralValidator{},
			Stages:         passThroughStages(),
			OperationScope: &op,
			TargetRef:      &target,
			Caller:         &caller,
		}
	case property11ConfigBoth:
		in = Input{
			Validator:      stubStructuralValidator{},
			Stages:         passThroughStages(),
			OperationScope: &op,
			TargetRef:      &target,
			TargetScope:    &ts,
			TargetScopeResolver: trackingTargetScopeResolver{
				scope:     ts,
				available: true,
				calls:     resolverCalls,
			},
			Authorizer: trackingScopeAuthorizer{decision: Allow, calls: authCalls},
			Caller:     &caller,
		}
	case property11ConfigMissingCaller:
		in = Input{
			Validator:      stubStructuralValidator{},
			Stages:         passThroughStages(),
			OperationScope: &op,
			TargetRef:      &target,
			TargetScope:    &ts,
			Authorizer:     trackingScopeAuthorizer{decision: Allow, calls: authCalls},
		}
	case property11ConfigMissingTargetRef:
		in = Input{
			Validator:      stubStructuralValidator{},
			Stages:         passThroughStages(),
			OperationScope: &op,
			TargetScope:    &ts,
			Authorizer:     trackingScopeAuthorizer{decision: Allow, calls: authCalls},
			Caller:         &caller,
		}
	case property11ConfigIncompletePathA:
		in = Input{
			Validator:      stubStructuralValidator{},
			Stages:         passThroughStages(),
			OperationScope: &op,
			TargetRef:      &target,
			TargetScope:    &ts,
			Caller:         &caller,
			// Authorizer nil → incomplete Path A
		}
	case property11ConfigIncompletePathB:
		in = Input{
			Validator:      stubStructuralValidator{},
			Stages:         passThroughStages(),
			OperationScope: &op,
			TargetRef:      &target,
			TargetScopeResolver: trackingTargetScopeResolver{
				scope:     ts,
				available: true,
				calls:     resolverCalls,
			},
			// Caller nil → incomplete Path B (also covered by missing-Caller
			// matrix rule; Path B requires Caller).
		}
	default:
		return fmt.Errorf("iteration %d: not a config-failure class: %d", iteration, c.Class)
	}

	res := Validate(context.Background(), in, DefaultLimits())
	if err := property11AssertInternalError(res, iteration, fmt.Sprintf("config class %d", c.Class)); err != nil {
		return err
	}
	if authCalls.Load() != 0 {
		return fmt.Errorf("iteration %d config class %d: Authorizer called %d times; want 0 (no target lookup)",
			iteration, c.Class, authCalls.Load())
	}
	if resolverCalls.Load() != 0 {
		return fmt.Errorf("iteration %d config class %d: TargetScopeResolver called %d times; want 0 (no target lookup)",
			iteration, c.Class, resolverCalls.Load())
	}
	return nil
}

func checkProperty11GenericNonOperation(c property11Case, iteration int) error {
	structCalls := &atomic.Int32{}
	defCalls := &atomic.Int32{}
	semCalls := &atomic.Int32{}
	refCalls := &atomic.Int32{}
	authCalls := &atomic.Int32{}
	resolverCalls := &atomic.Int32{}

	res := Validate(context.Background(), Input{
		Validator: trackingStructuralValidator{
			inner: stubStructuralValidator{},
			calls: structCalls,
		},
		Dst:    map[string]any{"ok": true},
		Stages: trackingStages(defCalls, semCalls, refCalls),
		// OperationScope nil → generic non-Operation; layer 8 must not run.
		Authorizer: trackingScopeAuthorizer{decision: Allow, calls: authCalls},
		TargetScopeResolver: trackingTargetScopeResolver{
			available: true,
			calls:     resolverCalls,
		},
	}, DefaultLimits())

	if structCalls.Load() != 1 {
		return fmt.Errorf("iteration %d generic: structural calls = %d, want 1", iteration, structCalls.Load())
	}
	if defCalls.Load() != 1 || semCalls.Load() != 1 || refCalls.Load() != 1 {
		return fmt.Errorf("iteration %d generic: layers 5–7 must run: defaulting=%d semantic=%d reference=%d",
			iteration, defCalls.Load(), semCalls.Load(), refCalls.Load())
	}
	if authCalls.Load() != 0 || resolverCalls.Load() != 0 {
		return fmt.Errorf("iteration %d generic: layer 8 must not run: auth=%d resolver=%d",
			iteration, authCalls.Load(), resolverCalls.Load())
	}
	_ = c // TargetRef/Caller unused on generic path by design
	return property11AssertSuccess(res, iteration, "generic non-Operation")
}

func property11AssertSuccess(res Result, iteration int, label string) error {
	if res.FailedAt != 0 {
		return fmt.Errorf("iteration %d %s: FailedAt = %v, want 0", iteration, label, res.FailedAt)
	}
	if res.Problem != nil || res.Err != nil || len(res.Violations) != 0 {
		return fmt.Errorf("iteration %d %s: unexpected failure %#v", iteration, label, res)
	}
	return nil
}

func property11AssertMismatch(res Result, iteration int, label string) error {
	if res.FailedAt != LayerAuthorization {
		return fmt.Errorf("iteration %d %s: FailedAt = %v, want %v", iteration, label, res.FailedAt, LayerAuthorization)
	}
	if res.Problem != nil || res.Err != nil {
		return fmt.Errorf("iteration %d %s: Problem/Err must be nil for ordinary mismatch: Problem=%v Err=%v",
			iteration, label, res.Problem, res.Err)
	}
	if len(res.Violations) != 1 {
		return fmt.Errorf("iteration %d %s: Violations len = %d, want 1; got %#v",
			iteration, label, len(res.Violations), res.Violations)
	}
	v := res.Violations[0]
	if v.Code != apiproblem.ViolationOperationTargetScopeMismatch {
		return fmt.Errorf("iteration %d %s: code = %q, want %q",
			iteration, label, v.Code, apiproblem.ViolationOperationTargetScopeMismatch)
	}
	if v.Field != "/metadata/scopeRef" {
		return fmt.Errorf("iteration %d %s: field = %q, want /metadata/scopeRef", iteration, label, v.Field)
	}
	return nil
}

func property11AssertSafeDenial(res Result, want *apiproblem.Problem, iteration int, label string) error {
	if res.FailedAt != LayerAuthorization {
		return fmt.Errorf("iteration %d %s: FailedAt = %v, want %v", iteration, label, res.FailedAt, LayerAuthorization)
	}
	if res.Err != nil {
		return fmt.Errorf("iteration %d %s: Err = %v, want nil", iteration, label, res.Err)
	}
	if len(res.Violations) != 0 {
		return fmt.Errorf("iteration %d %s: Violations must be empty (no mismatch disclosed); got %#v",
			iteration, label, res.Violations)
	}
	if res.Problem == nil || want == nil {
		return fmt.Errorf("iteration %d %s: Problem nil (got nil=%v want nil=%v)",
			iteration, label, res.Problem == nil, want == nil)
	}
	if res.Problem.Status != http.StatusNotFound || res.Problem.Code != apiproblem.CodeResourceNotFound {
		return fmt.Errorf("iteration %d %s: Problem = status %d code %q, want 404 RESOURCE_NOT_FOUND",
			iteration, label, res.Problem.Status, res.Problem.Code)
	}
	gotJSON, err := json.Marshal(res.Problem)
	if err != nil {
		return fmt.Errorf("iteration %d %s: marshal got: %v", iteration, label, err)
	}
	wantJSON, err := json.Marshal(want)
	if err != nil {
		return fmt.Errorf("iteration %d %s: marshal want: %v", iteration, label, err)
	}
	if !bytes.Equal(gotJSON, wantJSON) {
		return fmt.Errorf("iteration %d %s: SafeDenial not byte-identical:\ngot=%s\nwant=%s",
			iteration, label, gotJSON, wantJSON)
	}
	return nil
}

func property11AssertInternalError(res Result, iteration int, label string) error {
	if res.FailedAt != LayerAuthorization {
		return fmt.Errorf("iteration %d %s: FailedAt = %v, want %v", iteration, label, res.FailedAt, LayerAuthorization)
	}
	if res.Problem == nil {
		return fmt.Errorf("iteration %d %s: Problem is nil, want 500 INTERNAL_ERROR", iteration, label)
	}
	if res.Problem.Status != http.StatusInternalServerError {
		return fmt.Errorf("iteration %d %s: Problem.Status = %d, want %d",
			iteration, label, res.Problem.Status, http.StatusInternalServerError)
	}
	if res.Problem.Code != apiproblem.CodeInternalError {
		return fmt.Errorf("iteration %d %s: Problem.Code = %q, want %q",
			iteration, label, res.Problem.Code, apiproblem.CodeInternalError)
	}
	if res.Err == nil {
		return fmt.Errorf("iteration %d %s: Err is nil, want internal cause", iteration, label)
	}
	if len(res.Violations) != 0 {
		return fmt.Errorf("iteration %d %s: Violations = %#v, want empty", iteration, label, res.Violations)
	}
	return nil
}
