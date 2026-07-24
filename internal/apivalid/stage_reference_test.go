package apivalid

import (
	"context"
	"errors"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/apiref"
)

// Compile-time interface conformance for the test resource stub.
var _ ReferenceCarrier = (*stubReferenceResource)(nil)

type stubReferenceResource struct {
	scope       *apimeta.ScopeRef
	singular    map[string]apiref.TypedRef
	absent      map[string]struct{}
	collections map[string]apiref.Refs
	absentCols  map[string]struct{}
}

func (r *stubReferenceResource) GetScopeRef() *apimeta.ScopeRef { return r.scope }

func (r *stubReferenceResource) RefAt(path string) (apiref.TypedRef, bool) {
	if _, missing := r.absent[path]; missing {
		return apiref.TypedRef{}, false
	}
	ref, ok := r.singular[path]
	return ref, ok
}

func (r *stubReferenceResource) RefsAt(path string) (apiref.Refs, bool) {
	if _, missing := r.absentCols[path]; missing {
		return nil, false
	}
	refs, ok := r.collections[path]
	return refs, ok
}

func validReferenceConfig() ReferenceConfig {
	return ReferenceConfig{
		AllowedScopes: []apimeta.ScopeKind{apimeta.ScopeTenant, apimeta.ScopeProject},
		Fields: []RefField{{
			Path: "/spec/resourcePoolRef",
			Constraint: apiref.Constraint{
				AllowedKinds: []string{"ResourcePool"},
				Direction:    apiref.DirectionOutbound,
			},
		}},
	}
}

func validReferenceStub() *stubReferenceResource {
	return &stubReferenceResource{
		scope: &apimeta.ScopeRef{TypedRef: apimeta.TypedRef{
			APIVersion: "tenancy.sovrunn.io/v1alpha1",
			Kind:       string(apimeta.ScopeTenant),
			Name:       "acme",
			UID:        "tenant-uid-1",
		}},
		singular: map[string]apiref.TypedRef{
			"/spec/resourcePoolRef": {
				APIVersion: "fabric.sovrunn.io/v1alpha1",
				Kind:       "ResourcePool",
				Name:       "sovereign-pool-a",
				UID:        "pool-uid-1",
			},
		},
	}
}

func TestCommonReferenceValidRefPasses(t *testing.T) {
	t.Parallel()

	stage := NewCommonReference(validReferenceConfig(), DefaultLimits())
	violations, err := stage.Validate(context.Background(), validReferenceStub())
	if err != nil {
		t.Fatalf("Validate: unexpected error: %v", err)
	}
	if len(violations) != 0 {
		t.Fatalf("valid refs: violations=%#v, want none", violations)
	}
}

func TestCommonReferenceDisallowedKindViolation(t *testing.T) {
	t.Parallel()

	stage := NewCommonReference(validReferenceConfig(), DefaultLimits())
	obj := validReferenceStub()
	ref := obj.singular["/spec/resourcePoolRef"]
	ref.Kind = "Project"
	obj.singular["/spec/resourcePoolRef"] = ref

	violations, err := stage.Validate(context.Background(), obj)
	if err != nil {
		t.Fatalf("Validate: unexpected error: %v", err)
	}
	wantCode := apiproblem.ViolationCode(apiref.CodeKindNotAllowed)
	if !hasViolationCode(violations, wantCode) {
		t.Fatalf("disallowed kind must yield %s, got %#v", wantCode, violations)
	}
	if !hasViolationField(violations, "/spec/resourcePoolRef/kind") {
		t.Fatalf("disallowed kind must point at /spec/resourcePoolRef/kind, got %#v", violations)
	}
}

func TestCommonReferenceScopeNotInAllowedSetViolation(t *testing.T) {
	t.Parallel()

	stage := NewCommonReference(validReferenceConfig(), DefaultLimits())
	obj := validReferenceStub()
	obj.scope = &apimeta.ScopeRef{TypedRef: apimeta.TypedRef{
		APIVersion: "tenancy.sovrunn.io/v1alpha1",
		Kind:       string(apimeta.ScopeOrganization),
		Name:       "org-1",
		UID:        "org-uid-1",
	}}

	violations, err := stage.Validate(context.Background(), obj)
	if err != nil {
		t.Fatalf("Validate: unexpected error: %v", err)
	}
	wantCode := apiproblem.ViolationCode(apiref.CodeScopeNotAllowed)
	if !hasViolationCode(violations, wantCode) {
		t.Fatalf("disallowed scope must yield %s, got %#v", wantCode, violations)
	}
	if !hasViolationField(violations, "/metadata/scopeRef/kind") {
		t.Fatalf("disallowed scope must point at /metadata/scopeRef/kind, got %#v", violations)
	}
}

func TestCommonReferenceNilScopeRejectedWhenPlatformNotAllowed(t *testing.T) {
	t.Parallel()

	// AllowedScopes excludes Platform: canonical nil scopeRef must fail.
	stage := NewCommonReference(validReferenceConfig(), DefaultLimits())
	obj := validReferenceStub()
	obj.scope = nil

	violations, err := stage.Validate(context.Background(), obj)
	if err != nil {
		t.Fatalf("Validate: unexpected error: %v", err)
	}
	wantCode := apiproblem.ViolationCode(apiref.CodeScopeNotAllowed)
	if !hasViolationCode(violations, wantCode) {
		t.Fatalf("nil scope without Platform allowed must yield %s, got %#v", wantCode, violations)
	}
	if !hasViolationField(violations, "/metadata/scopeRef") {
		t.Fatalf("nil disallowed platform scope must point at /metadata/scopeRef, got %#v", violations)
	}
}

func TestCommonReferenceNilScopeAcceptedWhenPlatformAllowed(t *testing.T) {
	t.Parallel()

	cfg := ReferenceConfig{
		AllowedScopes: []apimeta.ScopeKind{apimeta.ScopePlatform},
		Fields: []RefField{{
			Path: "/spec/targetRef",
			Constraint: apiref.Constraint{
				AllowedKinds: []string{"PluginDefinition"},
				Direction:    apiref.DirectionOutbound,
			},
		}},
	}
	stage := NewCommonReference(cfg, DefaultLimits())
	obj := &stubReferenceResource{
		scope: nil,
		singular: map[string]apiref.TypedRef{
			"/spec/targetRef": {
				APIVersion: "plugin.sovrunn.io/v1alpha1",
				Kind:       "PluginDefinition",
				Name:       "demo-plugin",
			},
		},
	}

	violations, err := stage.Validate(context.Background(), obj)
	if err != nil {
		t.Fatalf("Validate: unexpected error: %v", err)
	}
	if len(violations) != 0 {
		t.Fatalf("platform-allowed nil scope must pass, got %#v", violations)
	}

	got := apimeta.CanonicalScopeIdentity(obj.GetScopeRef())
	want := apimeta.ScopeIdentity{Kind: apimeta.ScopePlatform, UID: apimeta.PlatformScopeUID}
	if got != want {
		t.Fatalf("CanonicalScopeIdentity=%#v, want %#v", got, want)
	}
}

func TestCommonReferenceMissingConstraintConfigError(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		cfg  ReferenceConfig
	}{
		{name: "empty config", cfg: ReferenceConfig{}},
		{name: "missing allowed scopes", cfg: ReferenceConfig{
			Fields: []RefField{{
				Path: "/spec/resourcePoolRef",
				Constraint: apiref.Constraint{
					AllowedKinds: []string{"ResourcePool"},
				},
			}},
		}},
		{name: "missing fields", cfg: ReferenceConfig{
			AllowedScopes: []apimeta.ScopeKind{apimeta.ScopeTenant},
		}},
		{name: "invalid allowed scope", cfg: ReferenceConfig{
			AllowedScopes: []apimeta.ScopeKind{apimeta.ScopeKind("NotAScope")},
			Fields: []RefField{{
				Path: "/spec/resourcePoolRef",
				Constraint: apiref.Constraint{
					AllowedKinds: []string{"ResourcePool"},
				},
			}},
		}},
		{name: "empty field path", cfg: ReferenceConfig{
			AllowedScopes: []apimeta.ScopeKind{apimeta.ScopeTenant},
			Fields: []RefField{{
				Path: "",
				Constraint: apiref.Constraint{
					AllowedKinds: []string{"ResourcePool"},
				},
			}},
		}},
		{name: "invalid direction", cfg: ReferenceConfig{
			AllowedScopes: []apimeta.ScopeKind{apimeta.ScopeTenant},
			Fields: []RefField{{
				Path: "/spec/resourcePoolRef",
				Constraint: apiref.Constraint{
					AllowedKinds: []string{"ResourcePool"},
					Direction:    apiref.Direction("Sideways"),
				},
			}},
		}},
		{name: "empty field constraint", cfg: ReferenceConfig{
			AllowedScopes: []apimeta.ScopeKind{apimeta.ScopeTenant},
			Fields: []RefField{{
				Path:       "/spec/resourcePoolRef",
				Constraint: apiref.Constraint{},
			}},
		}},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			stage := NewCommonReference(tc.cfg, DefaultLimits())
			_, err := stage.Validate(context.Background(), validReferenceStub())
			if !errors.Is(err, ErrReferenceConfig) {
				t.Fatalf("err=%v, want ErrReferenceConfig", err)
			}
		})
	}
}

func TestCommonReferenceInternalFaultReturnsError(t *testing.T) {
	t.Parallel()

	stage := NewCommonReference(validReferenceConfig(), DefaultLimits())

	_, err := stage.Validate(context.Background(), nil)
	if !errors.Is(err, ErrReferenceInternal) {
		t.Fatalf("nil object: err=%v, want ErrReferenceInternal", err)
	}

	var typedNil *stubReferenceResource
	_, err = stage.Validate(context.Background(), typedNil)
	if !errors.Is(err, ErrReferenceInternal) {
		t.Fatalf("typed nil: err=%v, want ErrReferenceInternal", err)
	}

	var nilStage *CommonReference
	_, err = nilStage.Validate(context.Background(), validReferenceStub())
	if !errors.Is(err, ErrReferenceInternal) {
		t.Fatalf("nil stage: err=%v, want ErrReferenceInternal", err)
	}
}

func TestCommonReferenceNoOpForUnknownKind(t *testing.T) {
	t.Parallel()

	stage := NewCommonReference(validReferenceConfig(), DefaultLimits())
	unknown := &unknownKindObject{Name: "x"}
	violations, err := stage.Validate(context.Background(), unknown)
	if err != nil {
		t.Fatalf("unknown type: unexpected error: %v", err)
	}
	if len(violations) != 0 {
		t.Fatalf("unknown type must no-op, got %#v", violations)
	}
}

func TestCommonReferenceCollectionDisallowedKind(t *testing.T) {
	t.Parallel()

	cfg := ReferenceConfig{
		AllowedScopes: []apimeta.ScopeKind{apimeta.ScopeTenant},
		Fields: []RefField{{
			Path:       "/spec/resourcePoolRefs",
			Collection: true,
			Constraint: apiref.Constraint{
				AllowedKinds: []string{"ResourcePool"},
				Direction:    apiref.DirectionOutbound,
			},
		}},
	}
	stage := NewCommonReference(cfg, DefaultLimits())
	obj := &stubReferenceResource{
		scope: &apimeta.ScopeRef{TypedRef: apimeta.TypedRef{
			APIVersion: "tenancy.sovrunn.io/v1alpha1",
			Kind:       string(apimeta.ScopeTenant),
			Name:       "acme",
			UID:        "tenant-uid-1",
		}},
		collections: map[string]apiref.Refs{
			"/spec/resourcePoolRefs": {{
				APIVersion: "fabric.sovrunn.io/v1alpha1",
				Kind:       "Project",
				Name:       "payments",
			}},
		},
	}

	violations, err := stage.Validate(context.Background(), obj)
	if err != nil {
		t.Fatalf("Validate: unexpected error: %v", err)
	}
	wantCode := apiproblem.ViolationCode(apiref.CodeKindNotAllowed)
	if !hasViolationCode(violations, wantCode) {
		t.Fatalf("collection disallowed kind must yield %s, got %#v", wantCode, violations)
	}
	if !hasViolationField(violations, "/spec/resourcePoolRefs/0/kind") {
		t.Fatalf("collection kind path must be indexed, got %#v", violations)
	}
}

func TestCommonReferenceConfigIsImmutableCopy(t *testing.T) {
	t.Parallel()

	cfg := validReferenceConfig()
	stage := NewCommonReference(cfg, DefaultLimits())

	// Mutate caller-owned slices after construction; stage must be unaffected.
	cfg.AllowedScopes[0] = apimeta.ScopeOrganization
	cfg.Fields[0].Constraint.AllowedKinds[0] = "Project"
	cfg.Fields[0].Path = "/mutated"

	violations, err := stage.Validate(context.Background(), validReferenceStub())
	if err != nil {
		t.Fatalf("Validate: unexpected error: %v", err)
	}
	if len(violations) != 0 {
		t.Fatalf("immutable config copy must still accept valid object, got %#v", violations)
	}
}

func TestCommonReferenceAbsentOptionalFieldSkipped(t *testing.T) {
	t.Parallel()

	stage := NewCommonReference(validReferenceConfig(), DefaultLimits())
	obj := validReferenceStub()
	obj.absent = map[string]struct{}{"/spec/resourcePoolRef": {}}
	delete(obj.singular, "/spec/resourcePoolRef")

	violations, err := stage.Validate(context.Background(), obj)
	if err != nil {
		t.Fatalf("Validate: unexpected error: %v", err)
	}
	if len(violations) != 0 {
		t.Fatalf("absent optional ref must be skipped, got %#v", violations)
	}
}
