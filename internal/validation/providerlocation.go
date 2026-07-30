package validation

import (
	"context"
	"fmt"
	"regexp"
	"unicode/utf8"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
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
	maxDisplayNameChars     = 253 // design §6.4
	maxGeoCountryCodeChars  = 2
	maxGeoSubdivisionChars  = 6
	violationValidationFail = apiproblem.ViolationCode("VALIDATION_FAILED")
)

// ValidateProviderLocation performs offline structural and semantic validation
// of a supplied ProviderLocation JSON document (FEATURE-0014 Task 12).
//
// Structural checks reuse apivalid.DecodeJSON with ModeCreateRequest and
// DefaultLimits: unknown/duplicate-field rejection, bounded object size and
// nesting, and rejection of user-authored status / system-owned metadata
// (F14-REQ-21, F14-REQ-22, F14-REQ-28).
//
// Semantic checks enforce metadata.scopeRef.kind == Provider and, when
// present, spec.geo syntax plus country-prefix agreement (F14-REQ-06,
// F14-REQ-27). Syntactically valid unassigned geographic codes are accepted
// without authoritative, residency, compliance, or placement meaning.
//
// Returns nil on success. Failures use inherited RFC 9457 Problem Details with
// stable codes and RFC 6901 paths. Messages carry no native identifiers or
// secrets (F14-REQ-31). The function performs no resource lookup and owns no
// state.
func ValidateProviderLocation(ctx context.Context, data []byte) *apiproblem.Problem {
	_ = ctx // reserved for request-scoped cancellation; offline path is non-blocking

	lim := apivalid.DefaultLimits()
	pol := apivalid.PolicyFor(apivalid.ModeCreateRequest)

	var loc resources.ProviderLocation
	if prob := apivalid.DecodeJSON(data, lim, pol, &loc); prob != nil {
		return prob
	}

	violations := validateProviderLocationDecoded(loc, lim)
	if len(violations) == 0 {
		return nil
	}
	if lim.MaxViolations > 0 && len(violations) > lim.MaxViolations {
		trimmed := make([]apiproblem.Violation, lim.MaxViolations)
		copy(trimmed, violations[:lim.MaxViolations])
		violations = trimmed
	}
	return apiproblem.New(apiproblem.CodeValidationFailed).
		WithDetail("one or more fields are invalid").
		WithViolations(violations)
}

func validateProviderLocationDecoded(loc resources.ProviderLocation, lim apivalid.Limits) []apiproblem.Violation {
	var out []apiproblem.Violation

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

	out = append(out, validateProviderLocationName(loc.Metadata.Name)...)
	out = append(out, validateProviderLocationDisplayName(loc.Metadata.DisplayName)...)
	out = append(out, validateProviderLocationLabels(loc.Metadata.Labels, lim)...)
	out = append(out, validateProviderLocationAnnotations(loc.Metadata.Annotations, lim)...)
	out = append(out, validateProviderLocationScopeRef(loc.Metadata.ScopeRef)...)
	out = append(out, validateProviderLocationGeo(loc.Spec.Geo)...)

	return out
}

func validateProviderLocationName(name string) []apiproblem.Violation {
	if name == "" {
		return []apiproblem.Violation{{
			Field:   "/metadata/name",
			Code:    apivalid.ViolationInvalidResourceName,
			Message: "resource name is required",
		}}
	}
	if utf8.RuneCountInString(name) > 63 || !dnsLabelRe.MatchString(name) {
		return []apiproblem.Violation{{
			Field:   "/metadata/name",
			Code:    apivalid.ViolationInvalidResourceName,
			Message: "resource name must be lowercase kebab-case DNS label, at most 63 characters",
		}}
	}
	return nil
}

func validateProviderLocationDisplayName(displayName string) []apiproblem.Violation {
	if displayName == "" {
		return nil
	}
	if utf8.RuneCountInString(displayName) > maxDisplayNameChars {
		return []apiproblem.Violation{{
			Field:   "/metadata/displayName",
			Code:    apiproblem.ViolationOutOfRange,
			Message: fmt.Sprintf("displayName must not exceed %d characters", maxDisplayNameChars),
		}}
	}
	return nil
}

func validateProviderLocationLabels(labels map[string]string, lim apivalid.Limits) []apiproblem.Violation {
	if labels == nil {
		return nil
	}
	var out []apiproblem.Violation
	if lim.MaxLabels > 0 && len(labels) > lim.MaxLabels {
		out = append(out, apiproblem.Violation{
			Field:   "/metadata/labels",
			Code:    apiproblem.ViolationOutOfRange,
			Message: fmt.Sprintf("labels exceed MaxLabels (%d)", lim.MaxLabels),
		})
	}
	for key, value := range labels {
		if lim.MaxLabelKeyChars > 0 && utf8.RuneCountInString(key) > lim.MaxLabelKeyChars {
			out = append(out, apiproblem.Violation{
				Field:   "/metadata/labels",
				Code:    apiproblem.ViolationOutOfRange,
				Message: fmt.Sprintf("label key exceeds MaxLabelKeyChars (%d)", lim.MaxLabelKeyChars),
			})
			break
		}
		if lim.MaxLabelValueChars > 0 && utf8.RuneCountInString(value) > lim.MaxLabelValueChars {
			out = append(out, apiproblem.Violation{
				Field:   "/metadata/labels",
				Code:    apiproblem.ViolationOutOfRange,
				Message: fmt.Sprintf("label value exceeds MaxLabelValueChars (%d)", lim.MaxLabelValueChars),
			})
			break
		}
	}
	return out
}

func validateProviderLocationAnnotations(annotations map[string]string, lim apivalid.Limits) []apiproblem.Violation {
	if annotations == nil || lim.MaxAnnotationsBytes <= 0 {
		return nil
	}
	total := 0
	for key, value := range annotations {
		total += len(key) + len(value)
		if total > lim.MaxAnnotationsBytes {
			return []apiproblem.Violation{{
				Field:   "/metadata/annotations",
				Code:    apiproblem.ViolationOutOfRange,
				Message: fmt.Sprintf("annotations exceed MaxAnnotationsBytes (%d)", lim.MaxAnnotationsBytes),
			}}
		}
	}
	return nil
}

func validateProviderLocationScopeRef(scope *apimeta.ScopeRef) []apiproblem.Violation {
	if scope == nil {
		return []apiproblem.Violation{{
			Field:   "/metadata/scopeRef",
			Code:    apivalid.ViolationScopeRefRequired,
			Message: "primary scopeRef is required and must identify a Provider",
		}}
	}
	if scope.Kind != string(apimeta.ScopeProvider) {
		return []apiproblem.Violation{{
			Field:   "/metadata/scopeRef/kind",
			Code:    apivalid.ViolationInvalidEnum,
			Message: "scopeRef.kind must be Provider",
		}}
	}
	return nil
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
