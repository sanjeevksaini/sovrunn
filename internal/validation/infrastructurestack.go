package validation

import (
	"context"
	"encoding/json"
	"regexp"
	"strings"

	"github.com/sanjeevksaini/sovrunn/internal/apicond"
	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/apiref"
	"github.com/sanjeevksaini/sovrunn/internal/apivalid"
	"github.com/sanjeevksaini/sovrunn/internal/resources"
)

const (
	infrastructureStackSchemaID = "api/schemas/infrastructure-stack.json"
	maxTechnologyChars          = 100
)

// technologyRe is the post-normalization accepted form (design §6.3 / DD-05).
// Compiled once for the hot validation path. Technology is never matched
// against a closed enum and is never used as identity or a lookup key.
var technologyRe = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9 ._+-]{0,99}$`)

// infrastructureStackScopeConstraint constrains metadata.scopeRef to Provider
// (F14-REQ-06). Offline only: no existence lookup.
var infrastructureStackScopeConstraint = apiref.Constraint{
	AllowedScopes: []apimeta.ScopeKind{apimeta.ScopeCloudProvider},
}

// infrastructureStackParentConstraint constrains the single immediate parent
// reference to kind DatacenterFailureDomain (F14-REQ-08, F14-REQ-09,
// F14-REQ-16, DD-08, F14-AD-010). Offline only: live parent existence and
// same-Provider resolution are CONTRACT_ONLY / NO_TASK (design §7.1 stage 3).
// Spanning (zero/multiple parents) is rejected by the singular required field.
var infrastructureStackParentConstraint = apiref.Constraint{
	AllowedKinds: []string{resources.KindDatacenterFailureDomain},
}

// infrastructureStackSemanticCarrier adapts a decoded InfrastructureStack to
// apivalid.SemanticCarrier so CommonSemantic can own common type-metadata,
// name, labels, annotations, scope presence, and limit checks without
// duplicating FEATURE-0012 grammar here.
type infrastructureStackSemanticCarrier struct {
	stack resources.InfrastructureStack
}

var _ apivalid.SemanticCarrier = infrastructureStackSemanticCarrier{}

func (c infrastructureStackSemanticCarrier) APIVersion() string { return c.stack.APIVersion }
func (c infrastructureStackSemanticCarrier) Kind() string       { return c.stack.Kind }
func (c infrastructureStackSemanticCarrier) ResourceName() string {
	return c.stack.Metadata.Name
}
func (c infrastructureStackSemanticCarrier) GetScopeRef() *apimeta.ScopeRef {
	return c.stack.Metadata.ScopeRef
}
func (c infrastructureStackSemanticCarrier) GetOwnerRef() *apimeta.OwnerRef { return nil }
func (c infrastructureStackSemanticCarrier) Labels() map[string]string {
	return c.stack.Metadata.Labels
}
func (c infrastructureStackSemanticCarrier) Annotations() map[string]string {
	return c.stack.Metadata.Annotations
}
func (c infrastructureStackSemanticCarrier) Conditions() []apicond.Condition { return nil }
func (c infrastructureStackSemanticCarrier) Phase() string                   { return "" }
func (c infrastructureStackSemanticCarrier) Profile() (apimeta.Profile, bool) {
	return "", false
}
func (c infrastructureStackSemanticCarrier) Boundary() (apimeta.Boundary, bool) {
	return "", false
}
func (c infrastructureStackSemanticCarrier) Stability() (apimeta.Stability, bool) {
	return "", false
}
func (c infrastructureStackSemanticCarrier) DataClassification() (apimeta.DataClassification, bool) {
	return "", false
}

// ValidateInfrastructureStack performs offline structural and semantic
// validation of a supplied InfrastructureStack JSON document
// (FEATURE-0014 Task 16).
//
// structural is an injected apivalid.StructuralValidator bound by the caller
// to api/schemas/infrastructure-stack.json. A nil or unavailable validator
// fails closed with INTERNAL_ERROR. Production code does not construct a
// registry, read schema files, or hold global validator state.
//
// Decode uses ModeCreateRequest + DefaultLimits (unknown/duplicate-field
// rejection, bounded size/nesting, status/system-owned rejection). Common
// semantic checks reuse apivalid.CommonSemantic. Local checks cover exact
// FEATURE-0014 identity, Provider-constrained scopeRef via apiref.Constraint,
// exactly one datacenterFailureDomainRef of kind DatacenterFailureDomain,
// and optional bounded descriptive technology with design §6.3 deterministic
// normalization. Technology is validated as the normalized verbatim string;
// it is not an enum, identity, or lookup key (F14-REQ-15, F14-REQ-17).
//
// Returns nil on success. Failures use inherited RFC 9457 Problem Details
// with stable codes and RFC 6901 paths. Messages carry no native identifiers
// or secrets (F14-REQ-31). The function performs no resource lookup and owns
// no state. The superseded stack-kind name is not introduced.
func ValidateInfrastructureStack(
	ctx context.Context,
	data []byte,
	structural apivalid.StructuralValidator,
) *apiproblem.Problem {
	_, prob := ValidateAndNormalizeInfrastructureStack(ctx, data, structural)
	return prob
}

// ValidateAndNormalizeInfrastructureStack performs the same pure offline
// validation as ValidateInfrastructureStack and returns the canonical decoded
// resource on success. The returned value carries technology after design §6.3
// normalization, making the canonical representation available without adding
// storage, lookup, registry, or runtime ownership to FEATURE-0014.
//
// On failure the returned resource is the zero value and prob describes the
// failure. Callers must use the returned resource, rather than the input bytes,
// when a later separately authorized component needs the canonical value.
func ValidateAndNormalizeInfrastructureStack(
	ctx context.Context,
	data []byte,
	structural apivalid.StructuralValidator,
) (resources.InfrastructureStack, *apiproblem.Problem) {
	lim := apivalid.DefaultLimits()
	pol := apivalid.PolicyFor(apivalid.ModeCreateRequest)

	var stack resources.InfrastructureStack
	if prob := apivalid.DecodeJSON(data, lim, pol, &stack); prob != nil {
		return resources.InfrastructureStack{}, prob
	}

	if structural == nil {
		return resources.InfrastructureStack{}, apiproblem.New(apiproblem.CodeInternalError).
			WithDetail("structural validator is unavailable")
	}

	// Validate the raw JSON document so required-field absence (e.g. missing
	// datacenterFailureDomainRef) is visible. Typed round-trip would invent
	// zero-value objects and hide REQUIRED_FIELD findings.
	var raw any
	if err := json.Unmarshal(data, &raw); err != nil {
		return resources.InfrastructureStack{}, apiproblem.New(apiproblem.CodeMalformedRequest).
			WithDetail("malformed JSON")
	}

	// Normalize technology before canonical-value validation, as required by
	// design §6.3. Update both the typed result and the raw structural input so
	// leading/trailing or repeated ASCII spaces are evaluated in canonical form.
	normalizedTechnology, technologyViolations := normalizeAndValidateTechnology(stack.Spec.Technology)
	if len(technologyViolations) > 0 {
		return resources.InfrastructureStack{}, validationFailedProblem(technologyViolations, lim)
	}
	stack.Spec.Technology = normalizedTechnology
	setNormalizedTechnology(raw, normalizedTechnology)

	structViolations, err := structural.Validate(raw, infrastructureStackSchemaID)
	if err != nil {
		return resources.InfrastructureStack{}, apiproblem.New(apiproblem.CodeInternalError).
			WithDetail("structural validator is unavailable")
	}
	if len(structViolations) > 0 {
		return resources.InfrastructureStack{}, validationFailedProblem(structViolations, lim)
	}

	violations, err := validateInfrastructureStackSemantic(ctx, stack, lim)
	if err != nil {
		return resources.InfrastructureStack{}, apiproblem.New(apiproblem.CodeInternalError).
			WithDetail("semantic validator is unavailable")
	}
	if len(violations) == 0 {
		return stack, nil
	}
	return resources.InfrastructureStack{}, validationFailedProblem(violations, lim)
}

func validateInfrastructureStackSemantic(
	ctx context.Context,
	stack resources.InfrastructureStack,
	lim apivalid.Limits,
) ([]apiproblem.Violation, error) {
	var out []apiproblem.Violation

	common := apivalid.NewCommonSemantic(lim, false)
	commonViolations, err := common.Validate(ctx, infrastructureStackSemanticCarrier{stack: stack})
	if err != nil {
		return nil, err
	}
	out = append(out, commonViolations...)

	// Exact FEATURE-0014 identity (CommonSemantic only checks group/version
	// shape and PascalCase kind).
	if stack.APIVersion != resources.FabricAPIVersion {
		out = append(out, apiproblem.Violation{
			Field:   "/apiVersion",
			Code:    apivalid.ViolationInvalidAPIVersion,
			Message: "apiVersion must be fabric.sovrunn.io/v1alpha1",
		})
	}
	if stack.Kind != resources.KindInfrastructureStack {
		out = append(out, apiproblem.Violation{
			Field:   "/kind",
			Code:    apivalid.ViolationInvalidKind,
			Message: "kind must be InfrastructureStack",
		})
	}

	out = append(out, validateInfrastructureStackScopeRef(stack.Metadata.ScopeRef)...)
	out = append(out, validateInfrastructureStackParentRef(stack.Spec.DatacenterFailureDomainRef)...)
	out = append(out, validateInfrastructureStackTechnology(stack.Spec.Technology)...)
	return out, nil
}

func validateInfrastructureStackScopeRef(scope *apimeta.ScopeRef) []apiproblem.Violation {
	if scope == nil {
		// CommonSemantic already reports SCOPE_REF_REQUIRED when
		// allowPlatformScope is false; avoid a duplicate violation here.
		return nil
	}
	return apivalid.RefIssuesToViolations(
		infrastructureStackScopeConstraint.ValidateRef(scope.TypedRef, "/metadata/scopeRef"),
	)
}

func validateInfrastructureStackParentRef(ref resources.DatacenterFailureDomainRef) []apiproblem.Violation {
	return apivalid.RefIssuesToViolations(
		infrastructureStackParentConstraint.ValidateRef(ref.TypedRef, "/spec/datacenterFailureDomainRef"),
	)
}

// validateInfrastructureStackTechnology applies design §6.3 deterministic
// normalization in fixed order, then validates the accepted pattern.
// Absent/empty technology is allowed (optional field). The normalized value
// is the canonical in-memory form; it is not an enum or identity.
func validateInfrastructureStackTechnology(tech string) []apiproblem.Violation {
	_, violations := normalizeAndValidateTechnology(tech)
	return violations
}

func normalizeAndValidateTechnology(tech string) (string, []apiproblem.Violation) {
	if tech == "" {
		return "", nil
	}

	normalized, ok := normalizeTechnology(tech)
	if !ok {
		return "", []apiproblem.Violation{{
			Field:   "/spec/technology",
			Code:    apiproblem.ViolationOutOfRange,
			Message: "technology must use only ASCII letters, digits, space, and ._+-",
		}}
	}
	if normalized == "" ||
		normalized[0] == ' ' ||
		normalized[len(normalized)-1] == ' ' ||
		len(normalized) > maxTechnologyChars ||
		!technologyRe.MatchString(normalized) {
		return "", []apiproblem.Violation{{
			Field:   "/spec/technology",
			Code:    apiproblem.ViolationOutOfRange,
			Message: "technology must be 1–100 characters matching the bounded descriptive pattern after normalization",
		}}
	}
	return normalized, nil
}

// setNormalizedTechnology updates the already-decoded raw document only when
// spec.technology was present in the input. Presence remains a schema concern:
// an absent optional value is not synthesized.
func setNormalizedTechnology(raw any, normalized string) {
	doc, ok := raw.(map[string]any)
	if !ok {
		return
	}
	spec, ok := doc["spec"].(map[string]any)
	if !ok {
		return
	}
	if _, present := spec["technology"]; present {
		spec["technology"] = normalized
	}
}

// normalizeTechnology applies design §6.3 normalization in fixed order:
// (1) reject any character outside the allowed ASCII set;
// (2) trim leading/trailing ASCII spaces;
// (3) collapse each internal run of spaces to a single space.
// Case is preserved. ok is false when step (1) fails.
func normalizeTechnology(s string) (string, bool) {
	for i := 0; i < len(s); i++ {
		if !isAllowedTechnologyByte(s[i]) {
			return "", false
		}
	}

	trimmed := strings.Trim(s, " ")
	if trimmed == "" {
		return "", true
	}

	var b strings.Builder
	b.Grow(len(trimmed))
	prevSpace := false
	for i := 0; i < len(trimmed); i++ {
		c := trimmed[i]
		if c == ' ' {
			if prevSpace {
				continue
			}
			prevSpace = true
			b.WriteByte(c)
			continue
		}
		prevSpace = false
		b.WriteByte(c)
	}
	return b.String(), true
}

func isAllowedTechnologyByte(c byte) bool {
	switch {
	case c >= 'A' && c <= 'Z':
		return true
	case c >= 'a' && c <= 'z':
		return true
	case c >= '0' && c <= '9':
		return true
	case c == ' ' || c == '.' || c == '_' || c == '+' || c == '-':
		return true
	default:
		return false
	}
}
