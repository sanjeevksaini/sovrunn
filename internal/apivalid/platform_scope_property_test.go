package apivalid

import (
	"context"
	"fmt"
	"math/rand"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/apiref"
)

// Deterministic seed for Property 4 (partial) reproducibility
// (F12-SCOPE-002, F12-REF-001; D-16).
const property4Seed int64 = 20260723

const property4Iterations = 100

// property4NonPlatformKinds are the five Matrix B kinds that are never
// canonicalized to nil by NormalizeScope.
var property4NonPlatformKinds = []apimeta.ScopeKind{
	apimeta.ScopeOrganization,
	apimeta.ScopeOrganizationUnit,
	apimeta.ScopeTenant,
	apimeta.ScopeProject,
	apimeta.ScopeProvider,
}

// property4Case pairs a ScopeRef form with a test-local allowed-scope
// contract ([]ScopeKind). Schema annotations are intentionally absent;
// complete Property 4 (schema x-sovrunn-allowed-scopes) is task 14.5.
type property4Case struct {
	Scope   *apimeta.ScopeRef
	Allowed []apimeta.ScopeKind
	Form    property4Form
}

type property4Form int

const (
	property4FormNil property4Form = iota
	property4FormExplicitPlatform
	property4FormNonPlatform
)

// Feature: api-resource-naming-status-and-validation-standard, Property 4 (partial): Canonical platform scope normalization and identity
//
// For any ScopeRef, NormalizeScope maps explicit Platform to nil;
// CanonicalScopeIdentity(nil) → {Platform, PlatformScopeUID};
// CanonicalScopeIdentity(non-platform) → {ref.Kind, ref.UID}; and
// NormalizeScope is idempotent. A test-local allowed-scope contract
// ([]ScopeKind) accepts nil/normalized Platform when Platform is allowed
// and rejects it when Platform is not allowed. Does not depend on
// canonical schemas or annotations (complete Property 4 is task 14.5).
//
// Validates: Requirements 4.4, 4.5 (F12-SCOPE-002, F12-REF-001)
func TestProperty4_CanonicalPlatformScopePartial(t *testing.T) {
	t.Parallel()

	rng := rand.New(rand.NewSource(property4Seed))
	for i := 0; i < property4Iterations; i++ {
		c := generateProperty4Case(rng, i)
		if err := checkProperty4Case(c, i); err != nil {
			t.Fatalf("property 4 (partial) failed at iteration %d (seed %d): %v", i, property4Seed, err)
		}
	}
}

func generateProperty4Case(rng *rand.Rand, iteration int) property4Case {
	// Force coverage buckets so each oracle class appears often.
	switch iteration % 6 {
	case 0:
		return property4Case{
			Scope:   nil,
			Allowed: property4AllowedWithPlatform(rng, true),
			Form:    property4FormNil,
		}
	case 1:
		return property4Case{
			Scope:   nil,
			Allowed: property4AllowedWithPlatform(rng, false),
			Form:    property4FormNil,
		}
	case 2:
		return property4Case{
			Scope:   property4ExplicitPlatform(rng),
			Allowed: property4AllowedWithPlatform(rng, true),
			Form:    property4FormExplicitPlatform,
		}
	case 3:
		return property4Case{
			Scope:   property4ExplicitPlatform(rng),
			Allowed: property4AllowedWithPlatform(rng, false),
			Form:    property4FormExplicitPlatform,
		}
	case 4:
		scope := property4NonPlatformScope(rng)
		kind := apimeta.ScopeKind(scope.Kind)
		return property4Case{
			Scope:   scope,
			Allowed: property4AllowedIncluding(rng, kind),
			Form:    property4FormNonPlatform,
		}
	default:
		scope := property4NonPlatformScope(rng)
		kind := apimeta.ScopeKind(scope.Kind)
		return property4Case{
			Scope:   scope,
			Allowed: property4AllowedExcluding(rng, kind),
			Form:    property4FormNonPlatform,
		}
	}
}

func property4ExplicitPlatform(rng *rand.Rand) *apimeta.ScopeRef {
	return &apimeta.ScopeRef{TypedRef: apimeta.TypedRef{
		APIVersion: "platform.sovrunn.io/v1alpha1",
		Kind:       string(apimeta.ScopePlatform),
		Name:       fmt.Sprintf("platform-%d", rng.Intn(10000)),
		UID:        apimeta.PlatformScopeUID,
	}}
}

func property4NonPlatformScope(rng *rand.Rand) *apimeta.ScopeRef {
	kind := property4NonPlatformKinds[rng.Intn(len(property4NonPlatformKinds))]
	return &apimeta.ScopeRef{TypedRef: apimeta.TypedRef{
		APIVersion: "tenancy.sovrunn.io/v1alpha1",
		Kind:       string(kind),
		Name:       fmt.Sprintf("%s-%d", kind, rng.Intn(10000)),
		UID:        fmt.Sprintf("%s-uid-%d", kind, rng.Intn(1_000_000)),
	}}
}

// property4AllowedWithPlatform returns a non-empty allowed-scope set that
// either includes or excludes Platform. When excluding Platform at least
// one other Matrix B kind is present so ReferenceConfig stays valid.
func property4AllowedWithPlatform(rng *rand.Rand, includePlatform bool) []apimeta.ScopeKind {
	out := make([]apimeta.ScopeKind, 0, 3)
	if includePlatform {
		out = append(out, apimeta.ScopePlatform)
	}
	// Optionally add 0..2 non-platform kinds.
	nExtra := rng.Intn(3)
	perm := rng.Perm(len(property4NonPlatformKinds))
	for i := 0; i < nExtra; i++ {
		out = append(out, property4NonPlatformKinds[perm[i]])
	}
	if len(out) == 0 {
		// Must remain non-empty for NewCommonReference.
		out = append(out, property4NonPlatformKinds[rng.Intn(len(property4NonPlatformKinds))])
	}
	return out
}

func property4AllowedIncluding(rng *rand.Rand, required apimeta.ScopeKind) []apimeta.ScopeKind {
	out := []apimeta.ScopeKind{required}
	if rng.Intn(2) == 0 {
		out = append(out, apimeta.ScopePlatform)
	}
	for _, k := range property4NonPlatformKinds {
		if k == required {
			continue
		}
		if rng.Intn(3) == 0 {
			out = append(out, k)
		}
	}
	return out
}

func property4AllowedExcluding(rng *rand.Rand, excluded apimeta.ScopeKind) []apimeta.ScopeKind {
	candidates := make([]apimeta.ScopeKind, 0, len(property4NonPlatformKinds)+1)
	candidates = append(candidates, apimeta.ScopePlatform)
	for _, k := range property4NonPlatformKinds {
		if k != excluded {
			candidates = append(candidates, k)
		}
	}
	// At least one allowed kind, never the excluded one.
	n := 1 + rng.Intn(len(candidates))
	perm := rng.Perm(len(candidates))
	out := make([]apimeta.ScopeKind, 0, n)
	for i := 0; i < n; i++ {
		out = append(out, candidates[perm[i]])
	}
	return out
}

func checkProperty4Case(c property4Case, iteration int) error {
	if len(c.Allowed) == 0 {
		return fmt.Errorf("iteration %d: generator produced empty allowed-scope contract", iteration)
	}
	for _, sk := range c.Allowed {
		if !sk.Valid() {
			return fmt.Errorf("iteration %d: invalid allowed scope kind %q", iteration, sk)
		}
	}

	normalized := apimeta.NormalizeScope(c.Scope)

	// Normalization rules by form.
	switch c.Form {
	case property4FormNil:
		if c.Scope != nil {
			return fmt.Errorf("iteration %d: nil-form case has non-nil Scope", iteration)
		}
		if normalized != nil {
			return fmt.Errorf("iteration %d: NormalizeScope(nil)=%#v, want nil", iteration, normalized)
		}
	case property4FormExplicitPlatform:
		if c.Scope == nil || apimeta.ScopeKind(c.Scope.Kind) != apimeta.ScopePlatform {
			return fmt.Errorf("iteration %d: explicit-platform form missing Platform kind", iteration)
		}
		if normalized != nil {
			return fmt.Errorf("iteration %d: NormalizeScope(Platform)=%#v, want nil", iteration, normalized)
		}
	case property4FormNonPlatform:
		if c.Scope == nil {
			return fmt.Errorf("iteration %d: non-platform form has nil Scope", iteration)
		}
		kind := apimeta.ScopeKind(c.Scope.Kind)
		if !kind.Valid() || kind == apimeta.ScopePlatform {
			return fmt.Errorf("iteration %d: non-platform form has invalid kind %q", iteration, c.Scope.Kind)
		}
		if normalized != c.Scope {
			return fmt.Errorf("iteration %d: NormalizeScope must leave non-platform unchanged (same pointer)", iteration)
		}
	default:
		return fmt.Errorf("iteration %d: unknown form %d", iteration, c.Form)
	}

	// Idempotence: NormalizeScope ∘ NormalizeScope ≡ NormalizeScope.
	again := apimeta.NormalizeScope(normalized)
	if again != normalized {
		// For nil both are nil (equal). For non-platform, same pointer.
		if again != nil || normalized != nil {
			return fmt.Errorf("iteration %d: NormalizeScope not idempotent: first=%#v second=%#v",
				iteration, normalized, again)
		}
	}

	// Identity oracle.
	idNorm := apimeta.CanonicalScopeIdentity(normalized)
	idRaw := apimeta.CanonicalScopeIdentity(c.Scope)
	platformID := apimeta.ScopeIdentity{Kind: apimeta.ScopePlatform, UID: apimeta.PlatformScopeUID}

	switch c.Form {
	case property4FormNil, property4FormExplicitPlatform:
		if idNorm != platformID {
			return fmt.Errorf("iteration %d: CanonicalScopeIdentity(normalized)=%#v, want %#v",
				iteration, idNorm, platformID)
		}
		if idRaw != platformID {
			return fmt.Errorf("iteration %d: CanonicalScopeIdentity(raw Platform form)=%#v, want %#v",
				iteration, idRaw, platformID)
		}
		// Explicit Platform and nil share identity after (and before) normalization.
		if idNorm != idRaw {
			return fmt.Errorf("iteration %d: platform forms must share identity: norm=%#v raw=%#v",
				iteration, idNorm, idRaw)
		}
	case property4FormNonPlatform:
		want := apimeta.ScopeIdentity{Kind: apimeta.ScopeKind(c.Scope.Kind), UID: c.Scope.UID}
		if idNorm != want {
			return fmt.Errorf("iteration %d: CanonicalScopeIdentity(non-platform)=%#v, want %#v",
				iteration, idNorm, want)
		}
		if idRaw != want {
			return fmt.Errorf("iteration %d: CanonicalScopeIdentity(raw non-platform)=%#v, want %#v",
				iteration, idRaw, want)
		}
	}

	// Test-local allowed-scope contract via CommonReference (no schema annotations).
	return checkProperty4AllowedScopeContract(c, normalized, idNorm, iteration)
}

// checkProperty4AllowedScopeContract validates Platform allowed → nil
// accepted and Platform not allowed → violation, using the generated
// []ScopeKind set as the trusted allowed-scope contract.
func checkProperty4AllowedScopeContract(
	c property4Case,
	normalized *apimeta.ScopeRef,
	identity apimeta.ScopeIdentity,
	iteration int,
) error {
	stage := NewCommonReference(ReferenceConfig{
		AllowedScopes: append([]apimeta.ScopeKind(nil), c.Allowed...),
		Fields: []RefField{{
			Path: "/spec/targetRef",
			Constraint: apiref.Constraint{
				AllowedKinds: []string{"PluginDefinition"},
				Direction:    apiref.DirectionOutbound,
			},
		}},
	}, DefaultLimits())

	obj := &stubReferenceResource{
		scope: normalized,
		singular: map[string]apiref.TypedRef{
			"/spec/targetRef": {
				APIVersion: "plugin.sovrunn.io/v1alpha1",
				Kind:       "PluginDefinition",
				Name:       fmt.Sprintf("demo-%d", iteration),
			},
		},
	}

	violations, err := stage.Validate(context.Background(), obj)
	if err != nil {
		return fmt.Errorf("iteration %d: CommonReference unexpected error: %v", iteration, err)
	}

	allowed := containsScopeKind(c.Allowed, identity.Kind)
	wantCode := apiproblem.ViolationCode(apiref.CodeScopeNotAllowed)

	if allowed {
		if hasViolationCode(violations, wantCode) {
			return fmt.Errorf("iteration %d: scope kind %q is allowed but got %s (allowed=%v violations=%#v)",
				iteration, identity.Kind, wantCode, c.Allowed, violations)
		}
		// Platform forms: nil accepted when Platform is in the contract.
		if identity.Kind == apimeta.ScopePlatform && normalized != nil {
			return fmt.Errorf("iteration %d: accepted platform scope must be canonical nil, got %#v",
				iteration, normalized)
		}
		return nil
	}

	if !hasViolationCode(violations, wantCode) {
		return fmt.Errorf("iteration %d: scope kind %q not in allowed=%v; want %s, got %#v",
			iteration, identity.Kind, c.Allowed, wantCode, violations)
	}
	// Nil/normalized platform points at /metadata/scopeRef; non-nil kinds at /kind.
	wantField := "/metadata/scopeRef"
	if normalized != nil && normalized.Kind != "" {
		wantField = "/metadata/scopeRef/kind"
	}
	if !hasViolationField(violations, wantField) {
		return fmt.Errorf("iteration %d: want field %q for disallowed scope, got %#v",
			iteration, wantField, violations)
	}

	// Explicit Platform (pre-normalization) and nil must produce the same
	// accept/reject class for a fixed allowed-scope contract.
	if c.Form == property4FormExplicitPlatform || c.Form == property4FormNil {
		if containsScopeKind(c.Allowed, apimeta.ScopePlatform) {
			return fmt.Errorf("iteration %d: platform form rejected but Platform is allowed=%v",
				iteration, c.Allowed)
		}
	}
	return nil
}
