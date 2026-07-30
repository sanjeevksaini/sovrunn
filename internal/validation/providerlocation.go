package validation

import (
	"context"
	"encoding/json"
	"regexp"
	"unicode/utf8"

	"github.com/sanjeevksaini/sovrunn/internal/apicond"
	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/apiref"
	"github.com/sanjeevksaini/sovrunn/internal/apivalid"
	"github.com/sanjeevksaini/sovrunn/internal/resources"
)

// Geographic descriptor syntax (design §6.2 / DD-04). Compiled once for the
// hot validation path. No ISO membership dataset or library is consulted.
var (
	geoCountryCodeRe     = regexp.MustCompile(`^[A-Z]{2}$`)
	geoSubdivisionCodeRe = regexp.MustCompile(`^[A-Z]{2}-[A-Z0-9]{1,3}$`)
)

const (
	providerLocationSchemaID = "api/schemas/provider-location.json"
	maxGeoCountryCodeChars   = 2
	maxGeoSubdivisionChars   = 6
	violationValidationFail  = apiproblem.ViolationCode("VALIDATION_FAILED")
)

// providerLocationScopeConstraint constrains metadata.scopeRef to Provider
// (F14-REQ-06). Offline only: no existence lookup.
var providerLocationScopeConstraint = apiref.Constraint{
	AllowedScopes: []apimeta.ScopeKind{apimeta.ScopeProvider},
}

// providerLocationSemanticCarrier adapts a decoded ProviderLocation to
// apivalid.SemanticCarrier so CommonSemantic can own common type-metadata,
// name, labels, annotations, scope presence, and limit checks without
// duplicating FEATURE-0012 grammar here.
type providerLocationSemanticCarrier struct {
	loc resources.ProviderLocation
}

var _ apivalid.SemanticCarrier = providerLocationSemanticCarrier{}

func (c providerLocationSemanticCarrier) APIVersion() string   { return c.loc.APIVersion }
func (c providerLocationSemanticCarrier) Kind() string         { return c.loc.Kind }
func (c providerLocationSemanticCarrier) ResourceName() string { return c.loc.Metadata.Name }
func (c providerLocationSemanticCarrier) GetScopeRef() *apimeta.ScopeRef {
	return c.loc.Metadata.ScopeRef
}
func (c providerLocationSemanticCarrier) GetOwnerRef() *apimeta.OwnerRef { return nil }
func (c providerLocationSemanticCarrier) Labels() map[string]string      { return c.loc.Metadata.Labels }
func (c providerLocationSemanticCarrier) Annotations() map[string]string {
	return c.loc.Metadata.Annotations
}
func (c providerLocationSemanticCarrier) Conditions() []apicond.Condition { return nil }
func (c providerLocationSemanticCarrier) Phase() string                   { return "" }
func (c providerLocationSemanticCarrier) Profile() (apimeta.Profile, bool) {
	return "", false
}
func (c providerLocationSemanticCarrier) Boundary() (apimeta.Boundary, bool) {
	return "", false
}
func (c providerLocationSemanticCarrier) Stability() (apimeta.Stability, bool) {
	return "", false
}
func (c providerLocationSemanticCarrier) DataClassification() (apimeta.DataClassification, bool) {
	return "", false
}

// ValidateProviderLocation performs offline structural and semantic validation
// of a supplied ProviderLocation JSON document (FEATURE-0014 Task 12).
//
// structural is an injected apivalid.StructuralValidator bound by the caller
// to api/schemas/provider-location.json. A nil or unavailable validator fails
// closed with INTERNAL_ERROR. Production code does not construct a registry,
// read schema files, or hold global validator state.
//
// Decode uses ModeCreateRequest + DefaultLimits (unknown/duplicate-field
// rejection, bounded size/nesting, status/system-owned rejection). Common
// semantic checks reuse apivalid.CommonSemantic. Local checks cover exact
// FEATURE-0014 identity, Provider-constrained scopeRef fields via
// apiref.Constraint, and optional spec.geo syntax/country-prefix agreement.
//
// Returns nil on success. Failures use inherited RFC 9457 Problem Details with
// stable codes and RFC 6901 paths. Messages carry no native identifiers or
// secrets (F14-REQ-31). The function performs no resource lookup and owns no
// state.
func ValidateProviderLocation(
	ctx context.Context,
	data []byte,
	structural apivalid.StructuralValidator,
) *apiproblem.Problem {
	lim := apivalid.DefaultLimits()
	pol := apivalid.PolicyFor(apivalid.ModeCreateRequest)

	var loc resources.ProviderLocation
	if prob := apivalid.DecodeJSON(data, lim, pol, &loc); prob != nil {
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
	structViolations, err := structural.Validate(raw, providerLocationSchemaID)
	if err != nil {
		return apiproblem.New(apiproblem.CodeInternalError).
			WithDetail("structural validator is unavailable")
	}
	if len(structViolations) > 0 {
		return validationFailedProblem(structViolations, lim)
	}

	violations, err := validateProviderLocationSemantic(ctx, loc, lim)
	if err != nil {
		return apiproblem.New(apiproblem.CodeInternalError).
			WithDetail("semantic validator is unavailable")
	}
	if len(violations) == 0 {
		return nil
	}
	return validationFailedProblem(violations, lim)
}

func validateProviderLocationSemantic(
	ctx context.Context,
	loc resources.ProviderLocation,
	lim apivalid.Limits,
) ([]apiproblem.Violation, error) {
	var out []apiproblem.Violation

	common := apivalid.NewCommonSemantic(lim, false)
	commonViolations, err := common.Validate(ctx, providerLocationSemanticCarrier{loc: loc})
	if err != nil {
		return nil, err
	}
	out = append(out, commonViolations...)

	// Exact FEATURE-0014 identity (CommonSemantic only checks group/version
	// shape and PascalCase kind).
	if loc.APIVersion != resources.FabricAPIVersion {
		out = append(out, apiproblem.Violation{
			Field:   "/apiVersion",
			Code:    apivalid.ViolationInvalidAPIVersion,
			Message: "apiVersion must be fabric.sovrunn.io/v1alpha1",
		})
	}
	if loc.Kind != resources.KindProviderLocation {
		out = append(out, apiproblem.Violation{
			Field:   "/kind",
			Code:    apivalid.ViolationInvalidKind,
			Message: "kind must be ProviderLocation",
		})
	}

	out = append(out, validateProviderLocationScopeRef(loc.Metadata.ScopeRef)...)
	out = append(out, validateProviderLocationGeo(loc.Spec.Geo)...)
	return out, nil
}

func validateProviderLocationScopeRef(scope *apimeta.ScopeRef) []apiproblem.Violation {
	if scope == nil {
		// CommonSemantic already reports SCOPE_REF_REQUIRED when
		// allowPlatformScope is false; avoid a duplicate violation here.
		return nil
	}
	return apivalid.RefIssuesToViolations(
		providerLocationScopeConstraint.ValidateRef(scope.TypedRef, "/metadata/scopeRef"),
	)
}

// validateProviderLocationGeo checks optional spec.geo syntax and country-
// prefix agreement (design §6.2). It performs no ISO assignment lookup:
// syntactically valid unassigned codes are accepted as declared metadata only.
func validateProviderLocationGeo(geo *resources.ProviderLocationGeo) []apiproblem.Violation {
	if geo == nil {
		return nil
	}

	var out []apiproblem.Violation
	if utf8.RuneCountInString(geo.CountryCode) != maxGeoCountryCodeChars || !geoCountryCodeRe.MatchString(geo.CountryCode) {
		out = append(out, apiproblem.Violation{
			Field:   "/spec/geo/countryCode",
			Code:    apiproblem.ViolationOutOfRange,
			Message: "countryCode must be exactly two uppercase ASCII letters (ISO 3166-1 alpha-2 syntax)",
		})
	}

	if geo.SubdivisionCode == "" {
		return out
	}
	if utf8.RuneCountInString(geo.SubdivisionCode) > maxGeoSubdivisionChars || !geoSubdivisionCodeRe.MatchString(geo.SubdivisionCode) {
		out = append(out, apiproblem.Violation{
			Field:   "/spec/geo/subdivisionCode",
			Code:    apiproblem.ViolationOutOfRange,
			Message: "subdivisionCode must match <countryCode>-<subdivision> (at most 6 characters)",
		})
		return out
	}
	// Country-prefix agreement is semantic and independent of membership.
	if len(geo.CountryCode) == maxGeoCountryCodeChars && geo.SubdivisionCode[:2] != geo.CountryCode {
		out = append(out, apiproblem.Violation{
			Field:   "/spec/geo/subdivisionCode",
			Code:    violationValidationFail,
			Message: "subdivisionCode country prefix must equal countryCode",
		})
	}
	return out
}

func validationFailedProblem(violations []apiproblem.Violation, lim apivalid.Limits) *apiproblem.Problem {
	if lim.MaxViolations > 0 && len(violations) > lim.MaxViolations {
		trimmed := make([]apiproblem.Violation, lim.MaxViolations)
		copy(trimmed, violations[:lim.MaxViolations])
		violations = trimmed
	}
	return apiproblem.New(apiproblem.CodeValidationFailed).
		WithDetail("one or more fields are invalid").
		WithViolations(violations)
}
