package validation

import (
	"context"
	"encoding/json"

	"github.com/sanjeevksaini/sovrunn/internal/apicond"
	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/apiref"
	"github.com/sanjeevksaini/sovrunn/internal/apivalid"
	"github.com/sanjeevksaini/sovrunn/internal/resources"
)

const providerSchemaID = "api/schemas/provider.json"

// providerScopeConstraint constrains an explicit metadata.scopeRef to
// Platform or Organization (F14-REQ-05). Offline only: no existence lookup.
// Nil/absent scopeRef is canonical Platform and is handled separately.
var providerScopeConstraint = apiref.Constraint{
	AllowedScopes: []apimeta.ScopeKind{apimeta.ScopePlatform, apimeta.ScopeOrganization},
}

// providerSemanticCarrier adapts a decoded Provider to apivalid.SemanticCarrier
// so CommonSemantic can own common type-metadata, name, labels, annotations,
// scope presence, and limit checks without duplicating FEATURE-0012 grammar.
// GetOwnerRef always returns nil: ObjectMeta has no ownerRef, and ownerRef
// must never act as governance scope for Provider (F14-REQ-05, F14-AD-003).
type providerSemanticCarrier struct {
	prov resources.Provider
}

var _ apivalid.SemanticCarrier = providerSemanticCarrier{}

func (c providerSemanticCarrier) APIVersion() string   { return c.prov.APIVersion }
func (c providerSemanticCarrier) Kind() string         { return c.prov.Kind }
func (c providerSemanticCarrier) ResourceName() string { return c.prov.Metadata.Name }
func (c providerSemanticCarrier) GetScopeRef() *apimeta.ScopeRef {
	return c.prov.Metadata.ScopeRef
}
func (c providerSemanticCarrier) GetOwnerRef() *apimeta.OwnerRef { return nil }
func (c providerSemanticCarrier) Labels() map[string]string      { return c.prov.Metadata.Labels }
func (c providerSemanticCarrier) Annotations() map[string]string {
	return c.prov.Metadata.Annotations
}
func (c providerSemanticCarrier) Conditions() []apicond.Condition { return nil }
func (c providerSemanticCarrier) Phase() string                   { return "" }
func (c providerSemanticCarrier) Profile() (apimeta.Profile, bool) {
	return "", false
}
func (c providerSemanticCarrier) Boundary() (apimeta.Boundary, bool) {
	return "", false
}
func (c providerSemanticCarrier) Stability() (apimeta.Stability, bool) {
	return "", false
}
func (c providerSemanticCarrier) DataClassification() (apimeta.DataClassification, bool) {
	return "", false
}

// ValidateProvider performs offline structural and semantic validation of a
// supplied Provider JSON document (FEATURE-0014 Task 13).
//
// structural is an injected apivalid.StructuralValidator bound by the caller
// to api/schemas/provider.json. A nil or unavailable validator fails closed
// with INTERNAL_ERROR. Production code does not construct a registry, read
// schema files, or hold global validator state.
//
// Decode uses ModeCreateRequest + DefaultLimits (unknown/duplicate-field
// rejection, bounded size/nesting, status/system-owned rejection). Common
// semantic checks reuse apivalid.CommonSemantic with allowPlatformScope so
// absent/nil scopeRef is canonical Platform. Explicit scopeRef is constrained
// to Platform or Organization via apiref.Constraint. Spec must carry no domain
// fields (empty ProviderSpec + schema additionalProperties:false).
//
// Returns nil on success. Failures use inherited RFC 9457 Problem Details with
// stable codes and RFC 6901 paths. Messages carry no native identifiers or
// secrets (F14-REQ-31). The function performs no resource lookup and owns no
// state.
func ValidateProvider(
	ctx context.Context,
	data []byte,
	structural apivalid.StructuralValidator,
) *apiproblem.Problem {
	lim := apivalid.DefaultLimits()
	pol := apivalid.PolicyFor(apivalid.ModeCreateRequest)

	var prov resources.Provider
	if prob := apivalid.DecodeJSON(data, lim, pol, &prov); prob != nil {
		return prob
	}

	if structural == nil {
		return apiproblem.New(apiproblem.CodeInternalError).
			WithDetail("structural validator is unavailable")
	}

	// Validate the raw JSON document so required-field absence (e.g. missing
	// top-level spec) is visible. Typed round-trip would invent zero-value
	// objects and hide REQUIRED_FIELD findings.
	var raw any
	if err := json.Unmarshal(data, &raw); err != nil {
		return apiproblem.New(apiproblem.CodeMalformedRequest).
			WithDetail("malformed JSON")
	}
	structViolations, err := structural.Validate(raw, providerSchemaID)
	if err != nil {
		return apiproblem.New(apiproblem.CodeInternalError).
			WithDetail("structural validator is unavailable")
	}
	if len(structViolations) > 0 {
		return validationFailedProblem(structViolations, lim)
	}

	violations, err := validateProviderSemantic(ctx, prov, lim)
	if err != nil {
		return apiproblem.New(apiproblem.CodeInternalError).
			WithDetail("semantic validator is unavailable")
	}
	if len(violations) == 0 {
		return nil
	}
	return validationFailedProblem(violations, lim)
}

func validateProviderSemantic(
	ctx context.Context,
	prov resources.Provider,
	lim apivalid.Limits,
) ([]apiproblem.Violation, error) {
	var out []apiproblem.Violation

	// allowPlatformScope=true: absent/nil scopeRef is canonical Platform
	// (F14-REQ-05, FEATURE-0012 D-16). ownerRef never substitutes for scope.
	common := apivalid.NewCommonSemantic(lim, true)
	commonViolations, err := common.Validate(ctx, providerSemanticCarrier{prov: prov})
	if err != nil {
		return nil, err
	}
	out = append(out, commonViolations...)

	// Exact FEATURE-0014 identity (CommonSemantic only checks group/version
	// shape and PascalCase kind).
	if prov.APIVersion != resources.FabricAPIVersion {
		out = append(out, apiproblem.Violation{
			Field:   "/apiVersion",
			Code:    apivalid.ViolationInvalidAPIVersion,
			Message: "apiVersion must be fabric.sovrunn.io/v1alpha1",
		})
	}
	if prov.Kind != resources.KindProvider {
		out = append(out, apiproblem.Violation{
			Field:   "/kind",
			Code:    apivalid.ViolationInvalidKind,
			Message: "kind must be Provider",
		})
	}

	out = append(out, validateProviderScopeRef(prov.Metadata.ScopeRef)...)
	return out, nil
}

func validateProviderScopeRef(scope *apimeta.ScopeRef) []apiproblem.Violation {
	if scope == nil {
		// Canonical Platform: CommonSemantic with allowPlatformScope already
		// accepted absent scopeRef. No explicit reference fields to check.
		return nil
	}
	// Explicit Platform or Organization requires a complete typed reference
	// constrained to those two scope kinds (F14-REQ-05).
	return apivalid.RefIssuesToViolations(
		providerScopeConstraint.ValidateRef(scope.TypedRef, "/metadata/scopeRef"),
	)
}
