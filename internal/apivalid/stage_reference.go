package apivalid

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/apiref"
)

// ErrReferenceInternal is returned when layer-7 structural reference
// validation encounters an internal fault (nil stage, nil object, or
// typed-nil object). Callers map this to Result.Problem = 500
// INTERNAL_ERROR at LayerReference (D-04). The error text MUST NOT
// include secrets, credentials, tokens, or raw object payloads.
var ErrReferenceInternal = errors.New("apivalid: reference validation internal fault")

// ErrReferenceConfig is returned when trusted reference-constraint
// configuration is missing or malformed. Callers map this to Result.Problem
// = 500 INTERNAL_ERROR at LayerReference (D-04; F12-VALIDATION-004).
// The error text MUST NOT include secrets or arbitrary caller input.
var ErrReferenceConfig = errors.New("apivalid: reference constraint configuration missing or malformed")

// RefField binds a JSON Pointer path to an immutable trusted
// apiref.Constraint for a singular TypedRef or a Refs collection
// (F12-REF-001/004, D-04). Construction-time only; Validate MUST NOT
// accept arbitrary caller-supplied constraints.
type RefField struct {
	// Path is the RFC 6901 JSON Pointer of the reference field
	// (for example "/spec/resourcePoolRef" or "/spec/targetRef").
	Path string
	// Constraint is the trusted allowed-kinds/scopes/direction rule.
	Constraint apiref.Constraint
	// Collection is true when Path names a Refs collection field.
	Collection bool
}

// ReferenceConfig is the immutable trusted configuration for
// CommonReference (D-04, D-02). AllowedScopes captures the
// schema-declared x-sovrunn-allowed-scopes set at construction time;
// it is not accepted as arbitrary runtime Input and MUST NOT require
// importing apischema. Fields carries trusted per-field typed-reference
// constraints.
type ReferenceConfig struct {
	AllowedScopes []apimeta.ScopeKind
	Fields        []RefField
}

// ReferenceCarrier is the resource surface required for common layer-7
// structural reference/kind/scope validation (D-04, F12-REF-001/002,
// F12-SCOPE-002). Conformance and domain types expose scopeRef and
// path-addressable typed references so CommonReference can validate
// without importing concrete contract packages or apischema.
//
// Objects that do not implement ReferenceCarrier are treated as unknown
// kinds with no applicable reference rules and receive a deterministic
// no-op (pass).
type ReferenceCarrier interface {
	// GetScopeRef returns the (already defaulted) metadata.scopeRef.
	// Nil is the canonical Platform form (D-16).
	GetScopeRef() *apimeta.ScopeRef
	// RefAt returns the singular TypedRef at path. present=false means
	// the field is absent (nil optional); ValidateRef is then skipped
	// for that field.
	RefAt(path string) (ref apiref.TypedRef, present bool)
	// RefsAt returns the Refs collection at path. present=false means
	// the collection field is absent.
	RefsAt(path string) (refs apiref.Refs, present bool)
}

// CommonReference implements ValidationStage for the Reference slot in
// StageSet (D-04, D-02; F12-VALIDATION-004, F12-REF-001/002, F12-SCOPE-002).
//
// Trusted rule configuration is owned by the stage value and set at
// construction: schema-declared allowed scopes and per-field
// apiref.Constraint values. Arbitrary caller-supplied constraints are
// not accepted at Validate time. apivalid MUST NOT import apischema.
//
// Behavior:
//   - Validates the object's scope kind against AllowedScopes (nil
//     scopeRef is canonical Platform via CanonicalScopeIdentity).
//   - Applies apiref.Constraint.ValidateRef (and Refs.Validate for
//     collections) for each configured field.
//   - Translates RefIssue results via RefIssuesToViolations.
//   - Does NOT resolve targets from external state and does NOT perform
//     Operation target-scope equality (layer 8 / D-17).
//
// A nil *CommonReference, nil/typed-nil object, or missing/malformed
// configuration fails closed with an error (500 path).
type CommonReference struct {
	allowedScopes []apimeta.ScopeKind
	fields        []RefField
	limits        Limits
	configOK      bool
}

// NewCommonReference returns a deterministic common reference-validation
// stage. cfg must include a non-empty AllowedScopes set and at least one
// Field constraint; invalid or empty configuration is retained so
// Validate fails closed with ErrReferenceConfig rather than panicking.
// limits should normally be DefaultLimits(); MaxViolations and
// MaxReferencesPerField bound ordinary findings.
func NewCommonReference(cfg ReferenceConfig, limits Limits) *CommonReference {
	s := &CommonReference{limits: limits}
	if err := validateReferenceConfig(cfg); err != nil {
		// Retain a non-ready stage so Validate fails closed deterministically.
		return s
	}
	s.allowedScopes = copyScopeKinds(cfg.AllowedScopes)
	s.fields = copyRefFields(cfg.Fields)
	s.configOK = true
	return s
}

// Compile-time check that CommonReference satisfies ValidationStage.
var _ ValidationStage = (*CommonReference)(nil)

// Validate runs structural reference/kind/scope checks on the (defaulted)
// object.
//
// Return semantics (ValidationStage):
//   - err != nil: internal fault or bad config → 500 at LayerReference
//   - err == nil, len(violations) > 0: ordinary findings → 422 handling
//   - err == nil, len(violations) == 0: pass
//
// Non-ReferenceCarrier objects are an explicit deterministic no-op when
// configuration is valid.
func (s *CommonReference) Validate(_ context.Context, object any) ([]apiproblem.Violation, error) {
	if s == nil {
		return nil, ErrReferenceInternal
	}
	if !s.configOK {
		return nil, ErrReferenceConfig
	}
	if isNilReferenceObject(object) {
		return nil, ErrReferenceInternal
	}

	carrier, ok := object.(ReferenceCarrier)
	if !ok {
		// Unknown kind / no applicable reference rules: explicit no-op.
		return nil, nil
	}

	var issues []apiref.RefIssue
	issues = append(issues, s.validateAllowedScopes(carrier.GetScopeRef())...)
	issues = append(issues, s.validateFields(carrier)...)

	return capViolations(RefIssuesToViolations(issues), s.limits.MaxViolations), nil
}

// validateAllowedScopes enforces schema-declared allowed-scope membership
// against the canonical scope identity (D-16, F12-SCOPE-002). Nil scopeRef
// is Platform; Platform is accepted only when present in AllowedScopes.
func (s *CommonReference) validateAllowedScopes(scope *apimeta.ScopeRef) []apiref.RefIssue {
	identity := apimeta.CanonicalScopeIdentity(scope)
	if containsScopeKind(s.allowedScopes, identity.Kind) {
		// When a non-nil scopeRef is present, also enforce TypedRef
		// well-formedness and AllowedScopes via ValidateRef.
		if scope == nil {
			return nil
		}
		c := apiref.Constraint{AllowedScopes: s.allowedScopes}
		return c.ValidateRef(scope.TypedRef, "/metadata/scopeRef")
	}

	field := "/metadata/scopeRef"
	if scope != nil && scope.Kind != "" {
		field = "/metadata/scopeRef/kind"
	}
	return []apiref.RefIssue{{
		Path:    field,
		Code:    apiref.CodeScopeNotAllowed,
		Message: "scope kind is not in the schema-declared allowed-scopes set",
	}}
}

func (s *CommonReference) validateFields(carrier ReferenceCarrier) []apiref.RefIssue {
	var issues []apiref.RefIssue
	maxRefs := s.limits.MaxReferencesPerField
	if maxRefs <= 0 {
		maxRefs = apiref.DefaultMaxRefs
	}

	for _, field := range s.fields {
		// Scope membership for /metadata/scopeRef is owned by
		// validateAllowedScopes so nil Platform is handled once.
		if field.Path == "/metadata/scopeRef" {
			continue
		}

		if field.Collection {
			refs, present := carrier.RefsAt(field.Path)
			if !present {
				continue
			}
			issues = append(issues, refs.Validate(field.Constraint, field.Path, maxRefs)...)
			continue
		}

		ref, present := carrier.RefAt(field.Path)
		if !present {
			continue
		}
		issues = append(issues, field.Constraint.ValidateRef(ref, field.Path)...)
	}
	return issues
}

func validateReferenceConfig(cfg ReferenceConfig) error {
	if len(cfg.AllowedScopes) == 0 {
		return fmt.Errorf("%w: allowed scopes required", ErrReferenceConfig)
	}
	if len(cfg.Fields) == 0 {
		return fmt.Errorf("%w: reference field constraints required", ErrReferenceConfig)
	}
	for _, sk := range cfg.AllowedScopes {
		if !sk.Valid() {
			return fmt.Errorf("%w: invalid allowed scope kind", ErrReferenceConfig)
		}
	}
	seen := make(map[string]struct{}, len(cfg.Fields))
	for _, f := range cfg.Fields {
		if f.Path == "" {
			return fmt.Errorf("%w: field path required", ErrReferenceConfig)
		}
		if _, dup := seen[f.Path]; dup {
			return fmt.Errorf("%w: duplicate field path", ErrReferenceConfig)
		}
		seen[f.Path] = struct{}{}
		if f.Constraint.Direction != "" && !f.Constraint.Direction.Valid() {
			return fmt.Errorf("%w: invalid reference direction", ErrReferenceConfig)
		}
		if len(f.Constraint.AllowedKinds) == 0 && len(f.Constraint.AllowedScopes) == 0 {
			return fmt.Errorf("%w: field constraint must declare allowed kinds or scopes", ErrReferenceConfig)
		}
		for _, sk := range f.Constraint.AllowedScopes {
			if !sk.Valid() {
				return fmt.Errorf("%w: invalid field allowed scope kind", ErrReferenceConfig)
			}
		}
	}
	return nil
}

func copyScopeKinds(in []apimeta.ScopeKind) []apimeta.ScopeKind {
	if len(in) == 0 {
		return nil
	}
	out := make([]apimeta.ScopeKind, len(in))
	copy(out, in)
	return out
}

func copyRefFields(in []RefField) []RefField {
	if len(in) == 0 {
		return nil
	}
	out := make([]RefField, len(in))
	for i, f := range in {
		out[i] = RefField{
			Path:       f.Path,
			Collection: f.Collection,
			Constraint: apiref.Constraint{
				AllowedKinds:  copyStrings(f.Constraint.AllowedKinds),
				AllowedScopes: copyScopeKinds(f.Constraint.AllowedScopes),
				Direction:     f.Constraint.Direction,
			},
		}
	}
	return out
}

func copyStrings(in []string) []string {
	if len(in) == 0 {
		return nil
	}
	out := make([]string, len(in))
	copy(out, in)
	return out
}

func containsScopeKind(set []apimeta.ScopeKind, v apimeta.ScopeKind) bool {
	for _, s := range set {
		if s == v {
			return true
		}
	}
	return false
}

// isNilReferenceObject reports whether object is a nil interface value or a
// typed nil (pointer/interface/map/slice/func/chan). Typed nils are treated
// as an internal fault so Validate never method-calls through a nil receiver.
func isNilReferenceObject(object any) bool {
	if object == nil {
		return true
	}
	v := reflect.ValueOf(object)
	switch v.Kind() {
	case reflect.Pointer, reflect.Interface, reflect.Map, reflect.Slice, reflect.Func, reflect.Chan:
		return v.IsNil()
	default:
		return false
	}
}
