package validate

import (
	"regexp"
	"unicode/utf8"

	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/cloudmodel/model"
)

// resourceNameRe is the URL-safe lowercase kebab-case identity grammar used
// with VS0-SCHEMA string(1..253) length bounds.
var resourceNameRe = regexp.MustCompile(`^[a-z0-9]([-a-z0-9]*[a-z0-9])?$`)

// ValidateSchema runs complete VS0-SCHEMA-008..014 field, format, range,
// reference, mutability, and request-contract constraints for the typed
// resource. Scope subset and ISO checks are included. Topology parent-chain
// existence is not checked here (see references.go).
func ValidateSchema(kind Kind, object any) *apiproblem.Problem {
	switch kind {
	case KindCloudPlatform:
		cp, ok := object.(model.CloudPlatform)
		if !ok {
			return typeMismatchProblem()
		}
		return validateCloudPlatform(cp)
	case KindCloudProvider:
		cp, ok := object.(model.CloudProvider)
		if !ok {
			return typeMismatchProblem()
		}
		return validateCloudProvider(cp)
	case KindCloudProviderParticipation:
		p, ok := object.(model.CloudProviderParticipation)
		if !ok {
			return typeMismatchProblem()
		}
		return validateParticipation(p)
	case KindHostingLocation:
		hl, ok := object.(model.HostingLocation)
		if !ok {
			return typeMismatchProblem()
		}
		return validateHostingLocation(hl)
	case KindDatacenter:
		dc, ok := object.(model.Datacenter)
		if !ok {
			return typeMismatchProblem()
		}
		return validateDatacenter(dc)
	case KindFaultDomain:
		fd, ok := object.(model.FaultDomain)
		if !ok {
			return typeMismatchProblem()
		}
		return validateFaultDomain(fd)
	case KindInfrastructureStack:
		st, ok := object.(model.InfrastructureStack)
		if !ok {
			return typeMismatchProblem()
		}
		return validateInfrastructureStack(st)
	default:
		return apiproblem.New(apiproblem.CodeValidationFailed).WithDetail("unknown resource kind")
	}
}

func validateCloudPlatform(cp model.CloudPlatform) *apiproblem.Problem {
	if prob := validateResourceName(cp.Metadata.Name); prob != nil {
		return prob
	}
	if prob := validateOptionalString("/metadata/displayName", cp.Metadata.DisplayName); prob != nil {
		return prob
	}
	if prob := ValidateScopeKindSubset(KindCloudPlatform, cp.Metadata.ScopeRef); prob != nil {
		return prob
	}
	if prob := validateRequiredString("/spec/ownerRegistration/legalName", cp.Spec.OwnerRegistration.LegalName); prob != nil {
		return prob
	}
	if prob := validateRequiredString("/spec/ownerRegistration/registrationIdentifier", cp.Spec.OwnerRegistration.RegistrationIdentifier); prob != nil {
		return prob
	}
	if prob := ValidateAssignedAlpha2("/spec/ownerRegistration/jurisdictionCode", cp.Spec.OwnerRegistration.JurisdictionCode); prob != nil {
		return prob
	}
	if prob := validateOptionalString("/spec/description", cp.Spec.Description); prob != nil {
		return prob
	}
	if len(cp.Status.Conditions) > MaxConditions {
		return outOfRange("/status/conditions", "conditions exceed maximum length")
	}
	return nil
}

func validateCloudProvider(cp model.CloudProvider) *apiproblem.Problem {
	if prob := validateResourceName(cp.Metadata.Name); prob != nil {
		return prob
	}
	if prob := validateOptionalString("/metadata/displayName", cp.Metadata.DisplayName); prob != nil {
		return prob
	}
	if prob := ValidateScopeKindSubset(KindCloudProvider, cp.Metadata.ScopeRef); prob != nil {
		return prob
	}
	if prob := validateOptionalString("/spec/displayName", cp.Spec.DisplayName); prob != nil {
		return prob
	}
	if prob := ValidateOperatingMarkets(cp.Spec.OperatingMarkets); prob != nil {
		return prob
	}
	if len(cp.Status.Conditions) > MaxConditions {
		return outOfRange("/status/conditions", "conditions exceed maximum length")
	}
	return nil
}

func validateParticipation(p model.CloudProviderParticipation) *apiproblem.Problem {
	if prob := validateResourceName(p.Metadata.Name); prob != nil {
		return prob
	}
	if prob := ValidateScopeKindSubset(KindCloudProviderParticipation, p.Metadata.ScopeRef); prob != nil {
		return prob
	}
	if prob := ValidateParticipationRefs(p); prob != nil {
		return prob
	}
	if p.Spec.Environment != model.ParticipationEnvironmentDevelopment {
		return outOfRange("/spec/environment", "environment must be development")
	}
	if p.Metadata.ScopeRef != nil {
		if prob := ValidateParticipationScopeUIDInvariant(p); prob != nil {
			return prob
		}
	}
	if len(p.Status.Conditions) > MaxConditions {
		return outOfRange("/status/conditions", "conditions exceed maximum length")
	}
	return nil
}

func validateHostingLocation(hl model.HostingLocation) *apiproblem.Problem {
	if prob := validateResourceName(hl.Metadata.Name); prob != nil {
		return prob
	}
	if prob := ValidateScopeKindSubset(KindHostingLocation, hl.Metadata.ScopeRef); prob != nil {
		return prob
	}
	if prob := ValidateAssignedAlpha2("/spec/countryCode", hl.Spec.CountryCode); prob != nil {
		return prob
	}
	if prob := validateRequiredString("/spec/locality", hl.Spec.Locality); prob != nil {
		return prob
	}
	if prob := ValidateAdministrativeAreaCode(hl.Spec.CountryCode, hl.Spec.AdministrativeAreaCode); prob != nil {
		return prob
	}
	if prob := validateOptionalString("/spec/description", hl.Spec.Description); prob != nil {
		return prob
	}
	if len(hl.Status.Conditions) > MaxConditions {
		return outOfRange("/status/conditions", "conditions exceed maximum length")
	}
	return nil
}

func validateDatacenter(dc model.Datacenter) *apiproblem.Problem {
	if prob := validateResourceName(dc.Metadata.Name); prob != nil {
		return prob
	}
	if prob := ValidateScopeKindSubset(KindDatacenter, dc.Metadata.ScopeRef); prob != nil {
		return prob
	}
	if prob := ValidateTypedRefStructural(dc.Spec.HostingLocationRef, "/spec/hostingLocationRef", model.KindHostingLocation, model.APIVersionHostingLocation); prob != nil {
		return prob
	}
	if prob := validateOptionalString("/spec/description", dc.Spec.Description); prob != nil {
		return prob
	}
	if len(dc.Status.Conditions) > MaxConditions {
		return outOfRange("/status/conditions", "conditions exceed maximum length")
	}
	return nil
}

func validateFaultDomain(fd model.FaultDomain) *apiproblem.Problem {
	if prob := validateResourceName(fd.Metadata.Name); prob != nil {
		return prob
	}
	if prob := ValidateScopeKindSubset(KindFaultDomain, fd.Metadata.ScopeRef); prob != nil {
		return prob
	}
	if prob := ValidateTypedRefStructural(fd.Spec.DatacenterRef, "/spec/datacenterRef", model.KindDatacenter, model.APIVersionDatacenter); prob != nil {
		return prob
	}
	if prob := validateOptionalString("/spec/description", fd.Spec.Description); prob != nil {
		return prob
	}
	if len(fd.Status.Conditions) > MaxConditions {
		return outOfRange("/status/conditions", "conditions exceed maximum length")
	}
	return nil
}

func validateInfrastructureStack(st model.InfrastructureStack) *apiproblem.Problem {
	if prob := validateResourceName(st.Metadata.Name); prob != nil {
		return prob
	}
	if prob := ValidateScopeKindSubset(KindInfrastructureStack, st.Metadata.ScopeRef); prob != nil {
		return prob
	}
	if prob := ValidateTypedRefStructural(st.Spec.FaultDomainRef, "/spec/faultDomainRef", model.KindFaultDomain, model.APIVersionFaultDomain); prob != nil {
		return prob
	}
	if prob := validateOptionalString("/spec/description", st.Spec.Description); prob != nil {
		return prob
	}
	if len(st.Status.Conditions) > MaxConditions {
		return outOfRange("/status/conditions", "conditions exceed maximum length")
	}
	return nil
}

func validateResourceName(name string) *apiproblem.Problem {
	if name == "" {
		return outOfRange("/metadata/name", "name is required")
	}
	if utf8.RuneCountInString(name) > MaxStringChars {
		return outOfRange("/metadata/name", "name exceeds maximum length")
	}
	if !resourceNameRe.MatchString(name) {
		return outOfRange("/metadata/name", "name must be lowercase kebab-case")
	}
	return nil
}

func validateRequiredString(field, value string) *apiproblem.Problem {
	if value == "" {
		return outOfRange(field, "field is required")
	}
	if utf8.RuneCountInString(value) > MaxStringChars {
		return outOfRange(field, "field exceeds maximum length")
	}
	return nil
}

func validateOptionalString(field, value string) *apiproblem.Problem {
	if value == "" {
		return nil
	}
	if utf8.RuneCountInString(value) > MaxStringChars {
		return outOfRange(field, "field exceeds maximum length")
	}
	return nil
}

func outOfRange(field, message string) *apiproblem.Problem {
	return apiproblem.New(apiproblem.CodeValidationFailed).WithViolations([]apiproblem.Violation{{
		Field:   field,
		Code:    apiproblem.ViolationOutOfRange,
		Message: message,
	}})
}

func typeMismatchProblem() *apiproblem.Problem {
	return apiproblem.New(apiproblem.CodeInternalError).WithDetail("schema validation type mismatch")
}
