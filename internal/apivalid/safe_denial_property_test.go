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

// Deterministic seed for Property 5 reproducibility
// (F12-SEC-004, F12-SCOPE-002; D-04).
const property5Seed int64 = 20260723

const property5Iterations = 100

// property5ScopeKinds covers all six Matrix B / D-17 governance scopes.
var property5ScopeKinds = []apimeta.ScopeKind{
	apimeta.ScopePlatform,
	apimeta.ScopeOrganization,
	apimeta.ScopeOrganizationUnit,
	apimeta.ScopeTenant,
	apimeta.ScopeProject,
	apimeta.ScopeProvider,
}

// property5Case models a denied cross-scope access under one of the three
// adopter contracts that must preserve path/response equivalence:
//   - Path A authorize-before-lookup (ScopeAuthorizer → DenyNotDisclosed)
//   - Path B AuthorizedTargetScopeResolver (available=false)
//   - Combined AuthorizedResolver (found=false for absent and unauthorized)
//
// LatentExists is intentionally not observable in the client-facing outcome;
// both values must share the same control-flow side-effect trace and the same
// SafeDenial Problem bytes. This is a path/response-equivalence property, not
// a constant-time guarantee.
type property5Case struct {
	Contract     property5Contract
	LatentExists bool
	OpScope      apimeta.ScopeIdentity
	TargetScope  apimeta.ScopeIdentity // may match or mismatch OpScope
	TargetRef    apiref.TypedRef
	Caller       CallerContext
}

type property5Contract int

const (
	property5PathA property5Contract = iota
	property5PathB
	property5AuthorizedResolver
)

// Feature: api-resource-naming-status-and-validation-standard, Property 5: Safe-denial path and response equivalence
//
// For any denied cross-scope access, exists-but-inaccessible and absent
// produce byte-identical SafeDenial Problem responses (404 RESOURCE_NOT_FOUND)
// and take the same control-flow path through the adopter contract
// (authorize-before-lookup, AuthorizedTargetScopeResolver, or combined
// AuthorizedResolver). No existence-dependent fast path is permitted. This is
// a path/response-equivalence property, not a constant-time guarantee.
//
// Validates: Requirements 4.4, 7.4 (F12-SEC-004, F12-SCOPE-002)
func TestProperty5_SafeDenialPathResponseEquivalence(t *testing.T) {
	t.Parallel()

	rng := rand.New(rand.NewSource(property5Seed))
	for i := 0; i < property5Iterations; i++ {
		c := generateProperty5Case(rng, i)
		if err := checkProperty5Case(c, i); err != nil {
			t.Fatalf("property 5 failed at iteration %d (seed %d): %v", i, property5Seed, err)
		}
	}
}

func generateProperty5Case(rng *rand.Rand, iteration int) property5Case {
	// Force coverage so each adopter contract and both latent existence
	// classes appear often across the 100 iterations.
	contract := property5Contract(iteration % 3)
	latentExists := (iteration/3)%2 == 0

	op := property5RandomScope(rng)
	targetScope := op
	if rng.Intn(2) == 0 {
		// Latent mismatch must never leak through SafeDenial.
		targetScope = property5DistinctScope(rng, op)
	}

	return property5Case{
		Contract:     contract,
		LatentExists: latentExists,
		OpScope:      op,
		TargetScope:  targetScope,
		TargetRef:    property5RandomTargetRef(rng),
		Caller: CallerContext{
			Scopes: []apimeta.ScopeIdentity{property5RandomScope(rng)},
		},
	}
}

func property5RandomScope(rng *rand.Rand) apimeta.ScopeIdentity {
	kind := property5ScopeKinds[rng.Intn(len(property5ScopeKinds))]
	if kind == apimeta.ScopePlatform {
		return apimeta.ScopeIdentity{Kind: apimeta.ScopePlatform, UID: apimeta.PlatformScopeUID}
	}
	return apimeta.ScopeIdentity{
		Kind: kind,
		UID:  fmt.Sprintf("%s-uid-%d", kind, rng.Intn(1_000_000)),
	}
}

func property5DistinctScope(rng *rand.Rand, other apimeta.ScopeIdentity) apimeta.ScopeIdentity {
	for {
		got := property5RandomScope(rng)
		if got != other {
			return got
		}
	}
}

func property5RandomTargetRef(rng *rand.Rand) apiref.TypedRef {
	kinds := []string{"Project", "ResourcePool", "PluginDefinition", "ServiceInstance"}
	return apiref.TypedRef{
		APIVersion: "core.sovrunn.io/v1alpha1",
		Kind:       kinds[rng.Intn(len(kinds))],
		Name:       fmt.Sprintf("target-%d", rng.Intn(100000)),
		UID:        fmt.Sprintf("uid-%d", rng.Intn(1_000_000)),
	}
}

func checkProperty5Case(c property5Case, iteration int) error {
	// Canonical absent Problem and DenyNotDisclosed SafeDenial must always
	// be byte-identical, independent of the generated case inputs.
	absent := apiproblem.New(apiproblem.CodeResourceNotFound)
	denied := SafeDenial(DenyNotDisclosed)
	if err := property5AssertProblemsByteIdentical(absent, denied, iteration, "absent vs DenyNotDisclosed"); err != nil {
		return err
	}
	if denied.Status != http.StatusNotFound || denied.Code != apiproblem.CodeResourceNotFound {
		return fmt.Errorf("iteration %d: SafeDenial(DenyNotDisclosed) = status %d code %q, want 404 RESOURCE_NOT_FOUND",
			iteration, denied.Status, denied.Code)
	}
	if denied.RequestID != "" || denied.Instance != "" || denied.Detail != "" {
		return fmt.Errorf("iteration %d: SafeDenial must leave requestId/instance/detail empty; got requestId=%q instance=%q detail=%q",
			iteration, denied.RequestID, denied.Instance, denied.Detail)
	}
	if len(denied.Violations) != 0 {
		return fmt.Errorf("iteration %d: SafeDenial violations must be empty; got %#v", iteration, denied.Violations)
	}

	switch c.Contract {
	case property5PathA:
		return checkProperty5PathA(c, iteration, denied)
	case property5PathB:
		return checkProperty5PathB(c, iteration, denied)
	case property5AuthorizedResolver:
		return checkProperty5AuthorizedResolver(c, iteration, denied)
	default:
		return fmt.Errorf("iteration %d: unknown contract %d", iteration, c.Contract)
	}
}

func checkProperty5PathA(c property5Case, iteration int, want *apiproblem.Problem) error {
	// Exists-but-inaccessible and absent share the authorize-before-lookup
	// path: Authorize is invoked once with no object lookup, and the
	// client-facing Problem is SafeDenial 404 even when TargetScope would
	// mismatch OperationScope. LatentExists is not an Input field on Path A
	// by design — denial happens before any existence-dependent work.
	op := c.OpScope
	ts := c.TargetScope
	target := c.TargetRef
	caller := c.Caller

	run := func() (Result, int32) {
		authCalls := &atomic.Int32{}
		res := Validate(context.Background(), Input{
			Validator:      stubStructuralValidator{},
			Stages:         passThroughStages(),
			OperationScope: &op,
			TargetRef:      &target,
			TargetScope:    &ts,
			Authorizer: trackingScopeAuthorizer{
				decision: DenyNotDisclosed,
				calls:    authCalls,
			},
			Caller: &caller,
			// No AuthorizedResolver / TargetScopeResolver: Path A must not
			// look up the target object after DenyNotDisclosed.
		}, DefaultLimits())
		return res, authCalls.Load()
	}

	res, calls := run()
	if calls != 1 {
		return fmt.Errorf("iteration %d path A: Authorizer calls = %d, want 1 (authorize-before-lookup)", iteration, calls)
	}
	if err := property5AssertSafeDenialResult(res, want, iteration, "path A"); err != nil {
		return err
	}

	// Pairwise equivalence across latent existence classes: Path A has no
	// existence input, so repeating the denial with the opposite latent class
	// label must yield an identical control-flow call count and Problem bytes.
	_ = c.LatentExists
	res2, calls2 := run()
	if calls2 != calls {
		return fmt.Errorf("iteration %d path A: authorizer call counts diverge across latent classes: %d vs %d",
			iteration, calls, calls2)
	}
	if err := property5AssertResultsPathEquivalent(res, res2, iteration, "path A paired latent classes"); err != nil {
		return err
	}
	return nil
}

func checkProperty5PathB(c property5Case, iteration int, want *apiproblem.Problem) error {
	// Path B: available=false is the single uniform outcome for both absent
	// and unauthorized. Latent scope detail must not leak, and mismatched
	// OperationScope must not surface OPERATION_TARGET_SCOPE_MISMATCH.
	op := c.OpScope
	target := c.TargetRef
	caller := c.Caller

	run := func(latentExists bool) (Result, []string, error) {
		trace := make([]string, 0, 2)
		resolverCalls := &atomic.Int32{}
		resolver := property5TrackingTargetScopeResolver{
			latentExists: latentExists,
			// Scope would mismatch if available=true; must never be returned.
			scope:     c.TargetScope,
			available: false,
			calls:     resolverCalls,
			trace:     &trace,
		}
		res := Validate(context.Background(), Input{
			Validator:           stubStructuralValidator{},
			Stages:              passThroughStages(),
			OperationScope:      &op,
			TargetRef:           &target,
			TargetScopeResolver: resolver,
			Caller:              &caller,
		}, DefaultLimits())
		if resolverCalls.Load() != 1 {
			return res, trace, fmt.Errorf("resolver calls = %d, want 1", resolverCalls.Load())
		}
		return res, trace, nil
	}

	resAbsent, traceAbsent, err := run(false)
	if err != nil {
		return fmt.Errorf("iteration %d path B absent: %v", iteration, err)
	}
	resUnauthorized, traceUnauthorized, err := run(true)
	if err != nil {
		return fmt.Errorf("iteration %d path B unauthorized: %v", iteration, err)
	}

	if !property5TracesEqual(traceAbsent, traceUnauthorized) {
		return fmt.Errorf("iteration %d path B: control-flow traces diverge: absent=%v unauthorized=%v",
			iteration, traceAbsent, traceUnauthorized)
	}
	if err := property5AssertSafeDenialResult(resAbsent, want, iteration, "path B absent"); err != nil {
		return err
	}
	if err := property5AssertSafeDenialResult(resUnauthorized, want, iteration, "path B unauthorized"); err != nil {
		return err
	}
	if err := property5AssertResultsPathEquivalent(resAbsent, resUnauthorized, iteration, "path B"); err != nil {
		return err
	}

	// Force the generated LatentExists class through the same oracle once more.
	resLatent, _, err := run(c.LatentExists)
	if err != nil {
		return fmt.Errorf("iteration %d path B latent=%v: %v", iteration, c.LatentExists, err)
	}
	return property5AssertSafeDenialResult(resLatent, want, iteration, "path B generated latent class")
}

func checkProperty5AuthorizedResolver(c property5Case, iteration int, want *apiproblem.Problem) error {
	// Combined AuthorizedResolver: absent and found-but-unauthorized MUST
	// return the same uniform unavailable outcome (found=false, no object)
	// and map through the same SafeDenial response. Side-effect traces must
	// not diverge on latent existence.
	ctx := context.Background()
	target := c.TargetRef
	caller := c.Caller

	run := func(latentExists bool) (obj any, found bool, trace []string, problem *apiproblem.Problem) {
		trace = make([]string, 0, 2)
		resolver := property5TrackingAuthorizedResolver{
			latentExists: latentExists,
			trace:        &trace,
		}
		obj, found = resolver.Resolve(ctx, caller, target)
		problem = SafeDenial(DenyNotDisclosed)
		return obj, found, trace, problem
	}

	objA, foundA, traceA, problemA := run(false)
	objB, foundB, traceB, problemB := run(true)

	if foundA || foundB {
		return fmt.Errorf("iteration %d AuthorizedResolver: found must be false for absent and unauthorized; got %v / %v",
			iteration, foundA, foundB)
	}
	if objA != nil || objB != nil {
		return fmt.Errorf("iteration %d AuthorizedResolver: objects must not leak; got %#v / %#v",
			iteration, objA, objB)
	}
	if !property5TracesEqual(traceA, traceB) {
		return fmt.Errorf("iteration %d AuthorizedResolver: control-flow traces diverge: absent=%v unauthorized=%v",
			iteration, traceA, traceB)
	}
	if err := property5AssertProblemsByteIdentical(problemA, want, iteration, "AuthorizedResolver absent"); err != nil {
		return err
	}
	if err := property5AssertProblemsByteIdentical(problemB, want, iteration, "AuthorizedResolver unauthorized"); err != nil {
		return err
	}
	if err := property5AssertProblemsByteIdentical(problemA, problemB, iteration, "AuthorizedResolver paired"); err != nil {
		return err
	}

	// Generated latent class must also follow the uniform unavailable path.
	objL, foundL, traceL, problemL := run(c.LatentExists)
	if foundL || objL != nil {
		return fmt.Errorf("iteration %d AuthorizedResolver latent=%v: want found=false obj=nil; got found=%v obj=%#v",
			iteration, c.LatentExists, foundL, objL)
	}
	if !property5TracesEqual(traceA, traceL) {
		return fmt.Errorf("iteration %d AuthorizedResolver: generated latent trace diverges: base=%v latent=%v",
			iteration, traceA, traceL)
	}
	return property5AssertProblemsByteIdentical(problemL, want, iteration, "AuthorizedResolver generated latent")
}

// property5TrackingTargetScopeResolver records a uniform unavailable path
// for both latent existence classes. Implementations MUST NOT branch the
// returned outcome or side-effect trace on LatentExists.
type property5TrackingTargetScopeResolver struct {
	latentExists bool
	scope        apimeta.ScopeIdentity
	available    bool
	calls        *atomic.Int32
	trace        *[]string
}

func (s property5TrackingTargetScopeResolver) ResolveAuthorizedTargetScope(
	_ context.Context,
	_ CallerContext,
	_ apiref.TypedRef,
) (apimeta.ScopeIdentity, bool) {
	if s.calls != nil {
		s.calls.Add(1)
	}
	if s.trace != nil {
		// Same branch label for absent and unauthorized — no existence-
		// dependent fast path or alternate audit/log side effect.
		*s.trace = append(*s.trace, "resolve-authorized-target-scope-unavailable")
		_ = s.latentExists // latent existence must not change the trace label
	}
	if !s.available {
		return apimeta.ScopeIdentity{}, false
	}
	return s.scope, true
}

// property5TrackingAuthorizedResolver returns the uniform unavailable
// outcome for both absent and unauthorized latent states.
type property5TrackingAuthorizedResolver struct {
	latentExists bool
	trace        *[]string
}

func (s property5TrackingAuthorizedResolver) Resolve(
	_ context.Context,
	_ CallerContext,
	_ apiref.TypedRef,
) (any, bool) {
	if s.trace != nil {
		*s.trace = append(*s.trace, "resolve-unavailable")
		_ = s.latentExists
	}
	return nil, false
}

func property5TracesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func property5AssertSafeDenialResult(res Result, want *apiproblem.Problem, iteration int, label string) error {
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
	return property5AssertProblemsByteIdentical(res.Problem, want, iteration, label)
}

func property5AssertResultsPathEquivalent(a, b Result, iteration int, label string) error {
	if a.FailedAt != b.FailedAt {
		return fmt.Errorf("iteration %d %s: FailedAt diverge: %v vs %v", iteration, label, a.FailedAt, b.FailedAt)
	}
	if (a.Err == nil) != (b.Err == nil) {
		return fmt.Errorf("iteration %d %s: Err presence diverge: %v vs %v", iteration, label, a.Err, b.Err)
	}
	if len(a.Violations) != len(b.Violations) {
		return fmt.Errorf("iteration %d %s: Violations length diverge: %d vs %d",
			iteration, label, len(a.Violations), len(b.Violations))
	}
	return property5AssertProblemsByteIdentical(a.Problem, b.Problem, iteration, label+" Problem")
}

func property5AssertProblemsByteIdentical(got, want *apiproblem.Problem, iteration int, label string) error {
	if got == nil || want == nil {
		return fmt.Errorf("iteration %d %s: Problem nil (got nil=%v want nil=%v)", iteration, label, got == nil, want == nil)
	}
	gotJSON, err := json.Marshal(got)
	if err != nil {
		return fmt.Errorf("iteration %d %s: marshal got: %v", iteration, label, err)
	}
	wantJSON, err := json.Marshal(want)
	if err != nil {
		return fmt.Errorf("iteration %d %s: marshal want: %v", iteration, label, err)
	}
	if !bytes.Equal(gotJSON, wantJSON) {
		return fmt.Errorf("iteration %d %s: Problems not byte-identical:\ngot=%s\nwant=%s",
			iteration, label, gotJSON, wantJSON)
	}
	return nil
}

// Compile-time interface conformance for Property 5 path-tracking stubs.
var (
	_ AuthorizedTargetScopeResolver = property5TrackingTargetScopeResolver{}
	_ AuthorizedResolver            = property5TrackingAuthorizedResolver{}
)
