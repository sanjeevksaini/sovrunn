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

const datacenterFailureDomainSchemaID = "api/schemas/datacenter-failure-domain.json"

// datacenterFailureDomainScopeConstraint constrains metadata.scopeRef to
// Provider (F14-REQ-06). Offline only: no existence lookup.
var datacenterFailureDomainScopeConstraint = apiref.Constraint{
	AllowedScopes: []apimeta.ScopeKind{apimeta.ScopeCloudProvider},
}

// datacenterFailureDomainParentConstraint constrains the single immediate
// parent reference to kind ProviderDatacenter (F14-REQ-08, F14-REQ-09,
// DD-08). Offline only: live parent existence and same-Provider resolution
// are CONTRACT_ONLY / NO_TASK (design §7.1 stage 3).
var datacenterFailureDomainParentConstraint = apiref.Constraint{
	AllowedKinds: []string{resources.KindProviderDatacenter},
}

// datacenterFailureDomainSemanticCarrier adapts a decoded
// DatacenterFailureDomain to apivalid.SemanticCarrier so CommonSemantic can
// own common type-metadata, name, labels, annotations, scope presence, and
// limit checks without duplicating FEATURE-0012 grammar here.
type datacenterFailureDomainSemanticCarrier struct {
	fd resources.DatacenterFailureDomain
}

var _ apivalid.SemanticCarrier = datacenterFailureDomainSemanticCarrier{}

func (c datacenterFailureDomainSemanticCarrier) APIVersion() string   { return c.fd.APIVersion }
func (c datacenterFailureDomainSemanticCarrier) Kind() string         { return c.fd.Kind }
func (c datacenterFailureDomainSemanticCarrier) ResourceName() string { return c.fd.Metadata.Name }
func (c datacenterFailureDomainSemanticCarrier) GetScopeRef() *apimeta.ScopeRef {
	return c.fd.Metadata.ScopeRef
}
func (c datacenterFailureDomainSemanticCarrier) GetOwnerRef() *apimeta.OwnerRef { return nil }
func (c datacenterFailureDomainSemanticCarrier) Labels() map[string]string {
	return c.fd.Metadata.Labels
}
func (c datacenterFailureDomainSemanticCarrier) Annotations() map[string]string {
	return c.fd.Metadata.Annotations
}
func (c datacenterFailureDomainSemanticCarrier) Conditions() []apicond.Condition { return nil }
func (c datacenterFailureDomainSemanticCarrier) Phase() string                   { return "" }
func (c datacenterFailureDomainSemanticCarrier) Profile() (apimeta.Profile, bool) {
	return "", false
}
func (c datacenterFailureDomainSemanticCarrier) Boundary() (apimeta.Boundary, bool) {
	return "", false
}
func (c datacenterFailureDomainSemanticCarrier) Stability() (apimeta.Stability, bool) {
	return "", false
}
func (c datacenterFailureDomainSemanticCarrier) DataClassification() (apimeta.DataClassification, bool) {
	return "", false
}

// ValidateDatacenterFailureDomain performs offline structural and semantic
// validation of a supplied DatacenterFailureDomain JSON document
// (FEATURE-0014 Task 15).
//
// structural is an injected apivalid.StructuralValidator bound by the caller
// to api/schemas/datacenter-failure-domain.json. A nil or unavailable
// validator fails closed with INTERNAL_ERROR. Production code does not
// construct a registry, read schema files, or hold global validator state.
//
// Decode uses ModeCreateRequest + DefaultLimits (unknown/duplicate-field
// rejection, bounded size/nesting, status/system-owned rejection). Common
// semantic checks reuse apivalid.CommonSemantic. Local checks cover exact
// FEATURE-0014 identity, Provider-constrained scopeRef via apiref.Constraint,
// and exactly one providerDatacenterRef of kind ProviderDatacenter.
//
// Returns nil on success. Failures use inherited RFC 9457 Problem Details
// with stable codes and RFC 6901 paths. Messages carry no native identifiers
// or secrets (F14-REQ-31). The function performs no resource lookup and owns
// no state. No connectivity field or inference is introduced (F14-REQ-19,
// F14-REQ-20).
func ValidateDatacenterFailureDomain(
	ctx context.Context,
	data []byte,
	structural apivalid.StructuralValidator,
) *apiproblem.Problem {
	lim := apivalid.DefaultLimits()
	pol := apivalid.PolicyFor(apivalid.ModeCreateRequest)

	var fd resources.DatacenterFailureDomain
	if prob := apivalid.DecodeJSON(data, lim, pol, &fd); prob != nil {
		return prob
	}

	if structural == nil {
		return apiproblem.New(apiproblem.CodeInternalError).
			WithDetail("structural validator is unavailable")
	}

	// Validate the raw JSON document so required-field absence (e.g. missing
	// providerDatacenterRef) is visible. Typed round-trip would invent
	// zero-value objects and hide REQUIRED_FIELD findings.
	var raw any
	if err := json.Unmarshal(data, &raw); err != nil {
		return apiproblem.New(apiproblem.CodeMalformedRequest).
			WithDetail("malformed JSON")
	}
	structViolations, err := structural.Validate(raw, datacenterFailureDomainSchemaID)
	if err != nil {
		return apiproblem.New(apiproblem.CodeInternalError).
			WithDetail("structural validator is unavailable")
	}
	if len(structViolations) > 0 {
		return validationFailedProblem(structViolations, lim)
	}

	violations, err := validateDatacenterFailureDomainSemantic(ctx, fd, lim)
	if err != nil {
		return apiproblem.New(apiproblem.CodeInternalError).
			WithDetail("semantic validator is unavailable")
	}
	if len(violations) == 0 {
		return nil
	}
	return validationFailedProblem(violations, lim)
}

func validateDatacenterFailureDomainSemantic(
	ctx context.Context,
	fd resources.DatacenterFailureDomain,
	lim apivalid.Limits,
) ([]apiproblem.Violation, error) {
	var out []apiproblem.Violation

	common := apivalid.NewCommonSemantic(lim, false)
	commonViolations, err := common.Validate(ctx, datacenterFailureDomainSemanticCarrier{fd: fd})
	if err != nil {
		return nil, err
	}
	out = append(out, commonViolations...)

	// Exact FEATURE-0014 identity (CommonSemantic only checks group/version
	// shape and PascalCase kind).
	if fd.APIVersion != resources.FabricAPIVersion {
		out = append(out, apiproblem.Violation{
			Field:   "/apiVersion",
			Code:    apivalid.ViolationInvalidAPIVersion,
			Message: "apiVersion must be fabric.sovrunn.io/v1alpha1",
		})
	}
	if fd.Kind != resources.KindDatacenterFailureDomain {
		out = append(out, apiproblem.Violation{
			Field:   "/kind",
			Code:    apivalid.ViolationInvalidKind,
			Message: "kind must be DatacenterFailureDomain",
		})
	}

	out = append(out, validateDatacenterFailureDomainScopeRef(fd.Metadata.ScopeRef)...)
	out = append(out, validateDatacenterFailureDomainParentRef(fd.Spec.ProviderDatacenterRef)...)
	return out, nil
}

func validateDatacenterFailureDomainScopeRef(scope *apimeta.ScopeRef) []apiproblem.Violation {
	if scope == nil {
		// CommonSemantic already reports SCOPE_REF_REQUIRED when
		// allowPlatformScope is false; avoid a duplicate violation here.
		return nil
	}
	return apivalid.RefIssuesToViolations(
		datacenterFailureDomainScopeConstraint.ValidateRef(scope.TypedRef, "/metadata/scopeRef"),
	)
}

func validateDatacenterFailureDomainParentRef(ref resources.ProviderDatacenterRef) []apiproblem.Violation {
	return apivalid.RefIssuesToViolations(
		datacenterFailureDomainParentConstraint.ValidateRef(ref.TypedRef, "/spec/providerDatacenterRef"),
	)
}
