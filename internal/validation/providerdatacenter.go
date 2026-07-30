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

const providerDatacenterSchemaID = "api/schemas/provider-datacenter.json"

// providerDatacenterScopeConstraint constrains metadata.scopeRef to Provider
// (F14-REQ-06). Offline only: no existence lookup.
var providerDatacenterScopeConstraint = apiref.Constraint{
	AllowedScopes: []apimeta.ScopeKind{apimeta.ScopeProvider},
}

// providerDatacenterParentConstraint constrains the single immediate parent
// reference to kind ProviderLocation (F14-REQ-08, F14-REQ-09, DD-08).
// Offline only: live parent existence and same-Provider resolution are
// CONTRACT_ONLY / NO_TASK (design §7.1 stage 3).
var providerDatacenterParentConstraint = apiref.Constraint{
	AllowedKinds: []string{resources.KindProviderLocation},
}

// providerDatacenterSemanticCarrier adapts a decoded ProviderDatacenter to
// apivalid.SemanticCarrier so CommonSemantic can own common type-metadata,
// name, labels, annotations, scope presence, and limit checks without
// duplicating FEATURE-0012 grammar here.
type providerDatacenterSemanticCarrier struct {
	dc resources.ProviderDatacenter
}

var _ apivalid.SemanticCarrier = providerDatacenterSemanticCarrier{}

func (c providerDatacenterSemanticCarrier) APIVersion() string   { return c.dc.APIVersion }
func (c providerDatacenterSemanticCarrier) Kind() string         { return c.dc.Kind }
func (c providerDatacenterSemanticCarrier) ResourceName() string { return c.dc.Metadata.Name }
func (c providerDatacenterSemanticCarrier) GetScopeRef() *apimeta.ScopeRef {
	return c.dc.Metadata.ScopeRef
}
func (c providerDatacenterSemanticCarrier) GetOwnerRef() *apimeta.OwnerRef { return nil }
func (c providerDatacenterSemanticCarrier) Labels() map[string]string      { return c.dc.Metadata.Labels }
func (c providerDatacenterSemanticCarrier) Annotations() map[string]string {
	return c.dc.Metadata.Annotations
}
func (c providerDatacenterSemanticCarrier) Conditions() []apicond.Condition { return nil }
func (c providerDatacenterSemanticCarrier) Phase() string                   { return "" }
func (c providerDatacenterSemanticCarrier) Profile() (apimeta.Profile, bool) {
	return "", false
}
func (c providerDatacenterSemanticCarrier) Boundary() (apimeta.Boundary, bool) {
	return "", false
}
func (c providerDatacenterSemanticCarrier) Stability() (apimeta.Stability, bool) {
	return "", false
}
func (c providerDatacenterSemanticCarrier) DataClassification() (apimeta.DataClassification, bool) {
	return "", false
}

// ValidateProviderDatacenter performs offline structural and semantic
// validation of a supplied ProviderDatacenter JSON document
// (FEATURE-0014 Task 14).
//
// structural is an injected apivalid.StructuralValidator bound by the caller
// to api/schemas/provider-datacenter.json. A nil or unavailable validator
// fails closed with INTERNAL_ERROR. Production code does not construct a
// registry, read schema files, or hold global validator state.
//
// Decode uses ModeCreateRequest + DefaultLimits (unknown/duplicate-field
// rejection, bounded size/nesting, status/system-owned rejection). Common
// semantic checks reuse apivalid.CommonSemantic. Local checks cover exact
// FEATURE-0014 identity, Provider-constrained scopeRef via apiref.Constraint,
// and exactly one providerLocationRef of kind ProviderLocation.
//
// Returns nil on success. Failures use inherited RFC 9457 Problem Details
// with stable codes and RFC 6901 paths. Messages carry no native identifiers
// or secrets (F14-REQ-31). The function performs no resource lookup and owns
// no state.
func ValidateProviderDatacenter(
	ctx context.Context,
	data []byte,
	structural apivalid.StructuralValidator,
) *apiproblem.Problem {
	lim := apivalid.DefaultLimits()
	pol := apivalid.PolicyFor(apivalid.ModeCreateRequest)

	var dc resources.ProviderDatacenter
	if prob := apivalid.DecodeJSON(data, lim, pol, &dc); prob != nil {
		return prob
	}

	if structural == nil {
		return apiproblem.New(apiproblem.CodeInternalError).
			WithDetail("structural validator is unavailable")
	}

	// Validate the raw JSON document so required-field absence (e.g. missing
	// providerLocationRef) is visible. Typed round-trip would invent
	// zero-value objects and hide REQUIRED_FIELD findings.
	var raw any
	if err := json.Unmarshal(data, &raw); err != nil {
		return apiproblem.New(apiproblem.CodeMalformedRequest).
			WithDetail("malformed JSON")
	}
	structViolations, err := structural.Validate(raw, providerDatacenterSchemaID)
	if err != nil {
		return apiproblem.New(apiproblem.CodeInternalError).
			WithDetail("structural validator is unavailable")
	}
	if len(structViolations) > 0 {
		return validationFailedProblem(structViolations, lim)
	}

	violations, err := validateProviderDatacenterSemantic(ctx, dc, lim)
	if err != nil {
		return apiproblem.New(apiproblem.CodeInternalError).
			WithDetail("semantic validator is unavailable")
	}
	if len(violations) == 0 {
		return nil
	}
	return validationFailedProblem(violations, lim)
}

func validateProviderDatacenterSemantic(
	ctx context.Context,
	dc resources.ProviderDatacenter,
	lim apivalid.Limits,
) ([]apiproblem.Violation, error) {
	var out []apiproblem.Violation

	common := apivalid.NewCommonSemantic(lim, false)
	commonViolations, err := common.Validate(ctx, providerDatacenterSemanticCarrier{dc: dc})
	if err != nil {
		return nil, err
	}
	out = append(out, commonViolations...)

	// Exact FEATURE-0014 identity (CommonSemantic only checks group/version
	// shape and PascalCase kind).
	if dc.APIVersion != resources.FabricAPIVersion {
		out = append(out, apiproblem.Violation{
			Field:   "/apiVersion",
			Code:    apivalid.ViolationInvalidAPIVersion,
			Message: "apiVersion must be fabric.sovrunn.io/v1alpha1",
		})
	}
	if dc.Kind != resources.KindProviderDatacenter {
		out = append(out, apiproblem.Violation{
			Field:   "/kind",
			Code:    apivalid.ViolationInvalidKind,
			Message: "kind must be ProviderDatacenter",
		})
	}

	out = append(out, validateProviderDatacenterScopeRef(dc.Metadata.ScopeRef)...)
	out = append(out, validateProviderDatacenterParentRef(dc.Spec.ProviderLocationRef)...)
	return out, nil
}

func validateProviderDatacenterScopeRef(scope *apimeta.ScopeRef) []apiproblem.Violation {
	if scope == nil {
		// CommonSemantic already reports SCOPE_REF_REQUIRED when
		// allowPlatformScope is false; avoid a duplicate violation here.
		return nil
	}
	return apivalid.RefIssuesToViolations(
		providerDatacenterScopeConstraint.ValidateRef(scope.TypedRef, "/metadata/scopeRef"),
	)
}

func validateProviderDatacenterParentRef(ref resources.ProviderLocationRef) []apiproblem.Violation {
	return apivalid.RefIssuesToViolations(
		providerDatacenterParentConstraint.ValidateRef(ref.TypedRef, "/spec/providerLocationRef"),
	)
}
