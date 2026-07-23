package apivalid

import (
	"context"
	"errors"
	"reflect"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
)

// ErrDefaultingInternal is returned when layer-5 defaulting encounters an
// internal fault (nil stage, nil object, or typed-nil object). Callers map
// this to Result.Problem = 500 INTERNAL_ERROR at LayerDefaulting (D-04).
// The error text MUST NOT include secrets, credentials, tokens, or raw
// object payloads.
var ErrDefaultingInternal = errors.New("apivalid: defaulting internal fault")

// ScopeDefaultable is the resource surface required for common layer-5
// defaulting (D-04, D-16). Conformance and domain types that carry
// metadata.scopeRef implement this so CommonDefaulting can normalize
// Platform scope without importing concrete contract packages.
//
// Objects that do not implement ScopeDefaultable are treated as unknown
// kinds with no applicable defaulting rules and receive a deterministic
// no-op.
type ScopeDefaultable interface {
	// DefaultingKind returns the singular PascalCase resource kind
	// (TypeMeta.kind). An empty kind is treated as unknown and no-ops.
	DefaultingKind() string
	// GetScopeRef returns the current metadata.scopeRef (may be nil).
	GetScopeRef() *apimeta.ScopeRef
	// SetScopeRef replaces metadata.scopeRef with the canonical form.
	SetScopeRef(scope *apimeta.ScopeRef)
}

// CommonDefaulting implements DefaultingStage with documented, versioned,
// deterministic common defaults only (F12-VALIDATION-003/004, D-04, D-16).
//
// Common rule: canonicalize Platform scope via apimeta.NormalizeScope so an
// explicit Kind=="Platform" scopeRef input alternate becomes the canonical
// absent (nil) form before identity, authorization, concurrency,
// persistence, and output processing (F12-SCOPE-002).
//
// The stage never depends on external state. Context is accepted for the
// DefaultingStage contract but is not consulted. Trusted rule configuration
// is owned by the stage value (immutable; currently the fixed common rule
// set). Arbitrary caller-supplied defaulting rules are not accepted at
// Apply time.
//
// A zero-value CommonDefaulting is usable. A nil *CommonDefaulting fails
// closed with ErrDefaultingInternal.
type CommonDefaulting struct{}

// NewCommonDefaulting returns a deterministic common defaulting stage.
func NewCommonDefaulting() *CommonDefaulting {
	return &CommonDefaulting{}
}

// Compile-time check that CommonDefaulting satisfies DefaultingStage.
var _ DefaultingStage = (*CommonDefaulting)(nil)

// Apply returns the defaulted object used by all later layers.
//
// Behavior:
//   - nil stage or nil/typed-nil object → ErrDefaultingInternal (fail closed)
//   - object does not implement ScopeDefaultable, or DefaultingKind is empty
//     → deterministic no-op (object returned unchanged)
//   - ScopeDefaultable with a non-empty kind → NormalizeScope on scopeRef;
//     Platform → nil; non-platform left unchanged
func (d *CommonDefaulting) Apply(_ context.Context, object any) (any, error) {
	if d == nil {
		return nil, ErrDefaultingInternal
	}
	if isNilDefaultingObject(object) {
		return nil, ErrDefaultingInternal
	}

	carrier, ok := object.(ScopeDefaultable)
	if !ok || carrier.DefaultingKind() == "" {
		// Unknown kind / no applicable defaulting rules: explicit no-op.
		return object, nil
	}

	carrier.SetScopeRef(apimeta.NormalizeScope(carrier.GetScopeRef()))
	return object, nil
}

// isNilDefaultingObject reports whether object is a nil interface value or a
// typed nil (pointer/interface/map/slice/func/chan). Typed nils are treated
// as an internal fault so Apply never method-calls through a nil receiver.
func isNilDefaultingObject(object any) bool {
	if object == nil {
		return true
	}
	v := reflect.ValueOf(object)
	switch v.Kind() {
	case reflect.Ptr, reflect.Interface, reflect.Map, reflect.Slice, reflect.Func, reflect.Chan:
		return v.IsNil()
	default:
		return false
	}
}
