package apivalid

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/sanjeevksaini/sovrunn/internal/apicond"
	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
)

// ErrSemanticInternal is returned when layer-6 semantic validation encounters
// an internal fault (nil stage, nil object, or typed-nil object). Callers map
// this to Result.Problem = 500 INTERNAL_ERROR at LayerSemantic (D-04).
// The error text MUST NOT include secrets, credentials, tokens, or raw
// object payloads.
var ErrSemanticInternal = errors.New("apivalid: semantic validation internal fault")

// Stable field-level violation codes for common semantic validation
// (F12-VALIDATION-006). Messages are informational and MUST NOT carry secrets.
const (
	ViolationInvalidResourceName      apiproblem.ViolationCode = "INVALID_RESOURCE_NAME"
	ViolationInvalidAPIVersion        apiproblem.ViolationCode = "INVALID_API_VERSION"
	ViolationInvalidKind              apiproblem.ViolationCode = "INVALID_KIND"
	ViolationInvalidEnum              apiproblem.ViolationCode = "INVALID_ENUM"
	ViolationInvalidCondition         apiproblem.ViolationCode = "INVALID_CONDITION"
	ViolationPhaseConditionIncoherent apiproblem.ViolationCode = "PHASE_CONDITION_INCOHERENT"
	ViolationScopeRefRequired         apiproblem.ViolationCode = "SCOPE_REF_REQUIRED"
)

// resourceNameRe is the F12-NAMING-002 resource-name grammar: lowercase
// kebab-case / DNS label (URL-safe), compiled once for the hot path.
var resourceNameRe = regexp.MustCompile(`^[a-z0-9]([-a-z0-9]*[a-z0-9])?$`)

const maxResourceNameChars = 63

// SemanticCarrier is the resource surface required for common layer-6
// semantic validation (D-04, D-06, F12-NAMING-001/002, F12-OWNER-001,
// F12-STATUS-003/005). Conformance and domain types that carry type
// metadata, ObjectMeta, optional ownerRef, and optional status phase/
// conditions implement this so CommonSemantic can validate without
// importing concrete contract packages.
//
// Objects that do not implement SemanticCarrier are treated as unknown
// kinds with no applicable semantic rules and receive a deterministic
// no-op (pass).
//
// Optional controlled-vocabulary fields return ok=false when absent on
// the object; when present (ok=true) they are checked against closed enums.
type SemanticCarrier interface {
	APIVersion() string
	Kind() string
	ResourceName() string
	GetScopeRef() *apimeta.ScopeRef
	GetOwnerRef() *apimeta.OwnerRef
	Labels() map[string]string
	Annotations() map[string]string
	Conditions() []apicond.Condition
	// Phase returns the coarse lifecycle summary when present; empty means
	// the object does not define a phase for coherence checks.
	Phase() string
	// Optional enums: ok=false means the field is not present on this object.
	Profile() (apimeta.Profile, bool)
	Boundary() (apimeta.Boundary, bool)
	Stability() (apimeta.Stability, bool)
	DataClassification() (apimeta.DataClassification, bool)
}

// CommonSemantic implements ValidationStage for the Semantic slot in
// StageSet (D-04, D-06; F12-VALIDATION-004/005, F12-NAMING-001, F12-OWNER-001).
//
// Trusted rule configuration is owned by the stage value and set at
// construction: Limits and whether nil scopeRef (canonical Platform) is
// allowed. Arbitrary caller-supplied semantic rules are not accepted at
// Validate time.
//
// Governance scope is always taken from scopeRef (via CanonicalScopeIdentity
// when adopters authorize). ownerRef is lifecycle containment only and
// MUST NOT substitute for a required scopeRef (F12-OWNER-001, D-17).
//
// A nil *CommonSemantic fails closed with ErrSemanticInternal.
type CommonSemantic struct {
	limits             Limits
	allowPlatformScope bool
}

// NewCommonSemantic returns a deterministic common semantic-validation stage.
// limits should normally be DefaultLimits(); zero fields disable the
// corresponding bound (see Limits). allowPlatformScope=true accepts a nil
// scopeRef as canonical Platform; allowPlatformScope=false requires a
// primary scopeRef and rejects ownerRef substitution for that requirement.
func NewCommonSemantic(limits Limits, allowPlatformScope bool) *CommonSemantic {
	return &CommonSemantic{
		limits:             limits,
		allowPlatformScope: allowPlatformScope,
	}
}

// Compile-time check that CommonSemantic satisfies ValidationStage.
var _ ValidationStage = (*CommonSemantic)(nil)

// Validate runs common semantic checks on the (defaulted) object.
//
// Return semantics (ValidationStage):
//   - err != nil: internal fault → 500 at LayerSemantic
//   - err == nil, len(violations) > 0: ordinary findings → 422 handling
//   - err == nil, len(violations) == 0: pass
//
// Non-SemanticCarrier objects are an explicit deterministic no-op.
func (s *CommonSemantic) Validate(_ context.Context, object any) ([]apiproblem.Violation, error) {
	if s == nil {
		return nil, ErrSemanticInternal
	}
	if isNilSemanticObject(object) {
		return nil, ErrSemanticInternal
	}

	carrier, ok := object.(SemanticCarrier)
	if !ok {
		// Unknown kind / no applicable semantic rules: explicit no-op.
		return nil, nil
	}

	var violations []apiproblem.Violation
	violations = append(violations, s.validateTypeMeta(carrier)...)
	violations = append(violations, s.validateResourceName(carrier.ResourceName())...)
	violations = append(violations, s.validateEnums(carrier)...)
	violations = append(violations, s.validateScopeAndOwner(carrier)...)
	violations = append(violations, s.validateLabels(carrier.Labels())...)
	violations = append(violations, s.validateAnnotations(carrier.Annotations())...)
	violations = append(violations, s.validatePhaseAndConditions(carrier.Phase(), carrier.Conditions())...)

	return capViolations(violations, s.limits.MaxViolations), nil
}

func (s *CommonSemantic) validateTypeMeta(c SemanticCarrier) []apiproblem.Violation {
	var out []apiproblem.Violation

	apiVersion := c.APIVersion()
	group, version, ok := apimeta.ParseAPIVersion(apiVersion)
	if !ok || !apimeta.IsKnownVersion(version) || !isSovrunnAPIGroup(group) {
		out = append(out, apiproblem.Violation{
			Field:   "/apiVersion",
			Code:    ViolationInvalidAPIVersion,
			Message: "apiVersion must be <domain>.sovrunn.io/{v1alpha1|v1beta1|v1}",
		})
	}

	kind := c.Kind()
	if !apicond.IsPascalCase(kind) {
		out = append(out, apiproblem.Violation{
			Field:   "/kind",
			Code:    ViolationInvalidKind,
			Message: "kind must be a singular PascalCase identifier",
		})
	}
	return out
}

func (s *CommonSemantic) validateResourceName(name string) []apiproblem.Violation {
	if name == "" {
		return []apiproblem.Violation{{
			Field:   "/metadata/name",
			Code:    ViolationInvalidResourceName,
			Message: "resource name is required",
		}}
	}
	if utf8.RuneCountInString(name) > maxResourceNameChars || !resourceNameRe.MatchString(name) {
		return []apiproblem.Violation{{
			Field:   "/metadata/name",
			Code:    ViolationInvalidResourceName,
			Message: "resource name must be lowercase kebab-case DNS label, at most 63 characters",
		}}
	}
	return nil
}

func (s *CommonSemantic) validateEnums(c SemanticCarrier) []apiproblem.Violation {
	var out []apiproblem.Violation

	if scope := c.GetScopeRef(); scope != nil {
		sk := apimeta.ScopeKind(scope.Kind)
		if !sk.Valid() {
			out = append(out, apiproblem.Violation{
				Field:   "/metadata/scopeRef/kind",
				Code:    ViolationInvalidEnum,
				Message: "scopeRef.kind must be a Matrix B scope kind",
			})
		}
	}

	if profile, present := c.Profile(); present && !profile.Valid() {
		out = append(out, apiproblem.Violation{
			Field:   "/metadata/profile",
			Code:    ViolationInvalidEnum,
			Message: "profile must be an approved Matrix A value",
		})
	}
	if boundary, present := c.Boundary(); present && !boundary.Valid() {
		out = append(out, apiproblem.Violation{
			Field:   "/metadata/boundary",
			Code:    ViolationInvalidEnum,
			Message: "boundary must be an approved Matrix C1 value",
		})
	}
	if stability, present := c.Stability(); present && !stability.Valid() {
		out = append(out, apiproblem.Violation{
			Field:   "/metadata/stability",
			Code:    ViolationInvalidEnum,
			Message: "stability must be alpha, beta, or stable",
		})
	}
	if class, present := c.DataClassification(); present && !class.Valid() {
		out = append(out, apiproblem.Violation{
			Field:   "/metadata/dataClassification",
			Code:    ViolationInvalidEnum,
			Message: "dataClassification must be an approved F12-SEC-002 value",
		})
	}
	return out
}

// validateScopeAndOwner enforces F12-OWNER-001 / D-17: ownerRef MUST NOT
// replace a required scopeRef or act as a governance/security scope.
// Authorization identity is derived only from scopeRef
// (apimeta.CanonicalScopeIdentity); ownerRef is ignored for that purpose.
// Identical ownerRef and scopeRef targets are allowed.
func (s *CommonSemantic) validateScopeAndOwner(c SemanticCarrier) []apiproblem.Violation {
	scope := c.GetScopeRef()
	if scope != nil {
		// scopeRef present: ownerRef may identify the same or a different
		// target; identity overlap alone is not a violation.
		return nil
	}
	if s.allowPlatformScope {
		// Canonical Platform (nil scopeRef). ownerRef may express lifecycle
		// containment but is never consulted for governance scope identity.
		return nil
	}
	// Required primary scopeRef is missing. ownerRef is ignored here so it
	// cannot substitute for governance/security scope (F12-OWNER-001).
	return []apiproblem.Violation{{
		Field:   "/metadata/scopeRef",
		Code:    ViolationScopeRefRequired,
		Message: "primary scopeRef is required; ownerRef must not substitute for governance scope",
	}}
}

func (s *CommonSemantic) validateLabels(labels map[string]string) []apiproblem.Violation {
	if labels == nil {
		return nil
	}
	var out []apiproblem.Violation
	if s.limits.MaxLabels > 0 && len(labels) > s.limits.MaxLabels {
		out = append(out, apiproblem.Violation{
			Field:   "/metadata/labels",
			Code:    apiproblem.ViolationOutOfRange,
			Message: fmt.Sprintf("labels exceed MaxLabels (%d)", s.limits.MaxLabels),
		})
	}
	// Deterministic key order is unnecessary for validity; iterate the map.
	for key, value := range labels {
		if s.limits.MaxLabelKeyChars > 0 && utf8.RuneCountInString(key) > s.limits.MaxLabelKeyChars {
			out = append(out, apiproblem.Violation{
				Field:   "/metadata/labels",
				Code:    apiproblem.ViolationOutOfRange,
				Message: fmt.Sprintf("label key exceeds MaxLabelKeyChars (%d)", s.limits.MaxLabelKeyChars),
			})
			break
		}
		if s.limits.MaxLabelValueChars > 0 && utf8.RuneCountInString(value) > s.limits.MaxLabelValueChars {
			out = append(out, apiproblem.Violation{
				Field:   "/metadata/labels",
				Code:    apiproblem.ViolationOutOfRange,
				Message: fmt.Sprintf("label value exceeds MaxLabelValueChars (%d)", s.limits.MaxLabelValueChars),
			})
			break
		}
	}
	return out
}

func (s *CommonSemantic) validateAnnotations(annotations map[string]string) []apiproblem.Violation {
	if annotations == nil || s.limits.MaxAnnotationsBytes <= 0 {
		return nil
	}
	total := 0
	for key, value := range annotations {
		total += len(key) + len(value)
		if total > s.limits.MaxAnnotationsBytes {
			return []apiproblem.Violation{{
				Field:   "/metadata/annotations",
				Code:    apiproblem.ViolationOutOfRange,
				Message: fmt.Sprintf("annotations exceed MaxAnnotationsBytes (%d)", s.limits.MaxAnnotationsBytes),
			}}
		}
	}
	return nil
}

func (s *CommonSemantic) validatePhaseAndConditions(phase string, conds []apicond.Condition) []apiproblem.Violation {
	var out []apiproblem.Violation

	if s.limits.MaxConditions > 0 && len(conds) > s.limits.MaxConditions {
		out = append(out, apiproblem.Violation{
			Field:   "/status/conditions",
			Code:    apiproblem.ViolationOutOfRange,
			Message: fmt.Sprintf("conditions exceed MaxConditions (%d)", s.limits.MaxConditions),
		})
	}

	seenTypes := make(map[string]struct{}, len(conds))
	for i, c := range conds {
		ptr := fmt.Sprintf("/status/conditions/%d", i)
		if !c.Status.Valid() {
			out = append(out, apiproblem.Violation{
				Field:   ptr + "/status",
				Code:    ViolationInvalidEnum,
				Message: "condition status must be True, False, or Unknown",
			})
		}
		if !apicond.IsPascalCase(c.Type) {
			out = append(out, apiproblem.Violation{
				Field:   ptr + "/type",
				Code:    ViolationInvalidCondition,
				Message: "condition type must be a stable PascalCase identifier",
			})
		}
		if !apicond.IsPascalCase(c.Reason) {
			out = append(out, apiproblem.Violation{
				Field:   ptr + "/reason",
				Code:    ViolationInvalidCondition,
				Message: "condition reason must be a stable PascalCase identifier",
			})
		}
		if c.Type != "" {
			if _, dup := seenTypes[c.Type]; dup {
				out = append(out, apiproblem.Violation{
					Field:   ptr + "/type",
					Code:    ViolationPhaseConditionIncoherent,
					Message: "condition types must be unique; conditions are current facts, not history",
				})
			} else {
				seenTypes[c.Type] = struct{}{}
			}
		}
	}

	if phase != "" {
		if !apicond.IsPascalCase(phase) {
			out = append(out, apiproblem.Violation{
				Field:   "/status/phase",
				Code:    ViolationPhaseConditionIncoherent,
				Message: "phase must be a PascalCase lifecycle summary when present",
			})
		}
		// Phase is a coarse lifecycle summary; it must not reuse condition
		// status vocabulary (True/False/Unknown) when both surfaces exist.
		if apicond.ConditionStatus(phase).Valid() {
			out = append(out, apiproblem.Violation{
				Field:   "/status/phase",
				Code:    ViolationPhaseConditionIncoherent,
				Message: "phase must not use condition status values True, False, or Unknown",
			})
		}
	}
	return out
}

func isSovrunnAPIGroup(group string) bool {
	return strings.HasSuffix(group, ".sovrunn.io") || group == "sovrunn.io"
}

func capViolations(v []apiproblem.Violation, max int) []apiproblem.Violation {
	if max <= 0 || len(v) <= max {
		return v
	}
	out := make([]apiproblem.Violation, max)
	copy(out, v[:max])
	return out
}

// isNilSemanticObject reports whether object is a nil interface value or a
// typed nil (pointer/interface/map/slice/func/chan). Typed nils are treated
// as an internal fault so Validate never method-calls through a nil receiver.
func isNilSemanticObject(object any) bool {
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
