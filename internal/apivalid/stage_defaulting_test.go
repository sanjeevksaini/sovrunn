package apivalid

import (
	"context"
	"errors"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
)

// Compile-time interface conformance for the test resource stub.
var _ ScopeDefaultable = (*stubScopedResource)(nil)

type stubScopedResource struct {
	kind  string
	scope *apimeta.ScopeRef
}

func (r *stubScopedResource) DefaultingKind() string { return r.kind }

func (r *stubScopedResource) GetScopeRef() *apimeta.ScopeRef { return r.scope }

func (r *stubScopedResource) SetScopeRef(scope *apimeta.ScopeRef) { r.scope = scope }

type unknownKindObject struct {
	Name string
}

func TestCommonDefaultingPlatformScopeNormalized(t *testing.T) {
	t.Parallel()

	stage := NewCommonDefaulting()
	platform := &apimeta.ScopeRef{TypedRef: apimeta.TypedRef{
		APIVersion: "tenancy.sovrunn.io/v1alpha1",
		Kind:       string(apimeta.ScopePlatform),
		Name:       "platform",
		UID:        apimeta.PlatformScopeUID,
	}}
	obj := &stubScopedResource{kind: "PluginDefinition", scope: platform}

	got, err := stage.Apply(context.Background(), obj)
	if err != nil {
		t.Fatalf("Apply: unexpected error: %v", err)
	}
	out, ok := got.(*stubScopedResource)
	if !ok {
		t.Fatalf("Apply returned %T, want *stubScopedResource", got)
	}
	if out != obj {
		t.Fatal("Apply must return the same object instance for later layers")
	}
	if out.scope != nil {
		t.Fatalf("Platform scopeRef = %#v, want nil after NormalizeScope", out.scope)
	}
}

func TestCommonDefaultingNonPlatformScopeUnchanged(t *testing.T) {
	t.Parallel()

	stage := NewCommonDefaulting()
	tenant := &apimeta.ScopeRef{TypedRef: apimeta.TypedRef{
		APIVersion: "tenancy.sovrunn.io/v1alpha1",
		Kind:       string(apimeta.ScopeTenant),
		Name:       "acme",
		UID:        "tenant-uid-1",
	}}
	obj := &stubScopedResource{kind: "Project", scope: tenant}

	got, err := stage.Apply(context.Background(), obj)
	if err != nil {
		t.Fatalf("Apply: unexpected error: %v", err)
	}
	out := got.(*stubScopedResource)
	if out.scope != tenant {
		t.Fatalf("non-platform scopeRef changed: got %#v, want same pointer %#v", out.scope, tenant)
	}
	if out.scope.Kind != string(apimeta.ScopeTenant) || out.scope.UID != "tenant-uid-1" {
		t.Fatalf("non-platform scope fields mutated: %#v", out.scope)
	}
}

func TestCommonDefaultingNilScopeUnchanged(t *testing.T) {
	t.Parallel()

	stage := NewCommonDefaulting()
	obj := &stubScopedResource{kind: "PluginDefinition", scope: nil}

	got, err := stage.Apply(context.Background(), obj)
	if err != nil {
		t.Fatalf("Apply: unexpected error: %v", err)
	}
	out := got.(*stubScopedResource)
	if out.scope != nil {
		t.Fatalf("nil scopeRef = %#v, want nil", out.scope)
	}
}

func TestCommonDefaultingNoOpForUnknownKind(t *testing.T) {
	t.Parallel()

	stage := NewCommonDefaulting()

	// Type that does not implement ScopeDefaultable.
	unknown := &unknownKindObject{Name: "x"}
	got, err := stage.Apply(context.Background(), unknown)
	if err != nil {
		t.Fatalf("unknown type: unexpected error: %v", err)
	}
	if got != unknown {
		t.Fatal("unknown type: must return the same object unchanged")
	}

	// ScopeDefaultable with empty kind is also treated as unknown.
	emptyKind := &stubScopedResource{
		kind: "",
		scope: &apimeta.ScopeRef{TypedRef: apimeta.TypedRef{
			Kind: string(apimeta.ScopePlatform),
			Name: "platform",
		}},
	}
	got, err = stage.Apply(context.Background(), emptyKind)
	if err != nil {
		t.Fatalf("empty kind: unexpected error: %v", err)
	}
	out := got.(*stubScopedResource)
	if out.scope == nil || out.scope.Kind != string(apimeta.ScopePlatform) {
		t.Fatalf("empty kind must no-op without normalizing; scope=%#v", out.scope)
	}
}

func TestCommonDefaultingInternalFaultReturnsError(t *testing.T) {
	t.Parallel()

	stage := NewCommonDefaulting()

	_, err := stage.Apply(context.Background(), nil)
	if !errors.Is(err, ErrDefaultingInternal) {
		t.Fatalf("nil object: err=%v, want ErrDefaultingInternal", err)
	}

	var typedNil *stubScopedResource
	_, err = stage.Apply(context.Background(), typedNil)
	if !errors.Is(err, ErrDefaultingInternal) {
		t.Fatalf("typed nil: err=%v, want ErrDefaultingInternal", err)
	}

	var nilStage *CommonDefaulting
	_, err = nilStage.Apply(context.Background(), &stubScopedResource{kind: "Project"})
	if !errors.Is(err, ErrDefaultingInternal) {
		t.Fatalf("nil stage: err=%v, want ErrDefaultingInternal", err)
	}
}

func TestCommonDefaultingIdempotentPlatformNormalization(t *testing.T) {
	t.Parallel()

	stage := NewCommonDefaulting()
	obj := &stubScopedResource{
		kind: "Operation",
		scope: &apimeta.ScopeRef{TypedRef: apimeta.TypedRef{
			Kind: string(apimeta.ScopePlatform),
			Name: "platform",
		}},
	}

	if _, err := stage.Apply(context.Background(), obj); err != nil {
		t.Fatalf("first Apply: %v", err)
	}
	if obj.scope != nil {
		t.Fatalf("after first Apply: scope=%#v, want nil", obj.scope)
	}
	if _, err := stage.Apply(context.Background(), obj); err != nil {
		t.Fatalf("second Apply: %v", err)
	}
	if obj.scope != nil {
		t.Fatalf("after second Apply: scope=%#v, want nil", obj.scope)
	}
}
