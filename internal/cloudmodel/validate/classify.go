package validate

import (
	"encoding/json"
	"sort"
	"strings"

	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/cloudmodel/model"
)

// Slice 0 violation codes carried only in violations[].code (design §5.2).
const (
	ViolationScopeReferenceMismatch   apiproblem.ViolationCode = "VS0_SCOPE_REFERENCE_MISMATCH"
	ViolationTopologyProviderMismatch apiproblem.ViolationCode = "VS0_TOPOLOGY_PROVIDER_MISMATCH"
	ViolationPatchImmutableField      apiproblem.ViolationCode = "VS0_PATCH_IMMUTABLE_FIELD"
	ViolationScopeKindInvalid         apiproblem.ViolationCode = "VS0_SCOPE_KIND_INVALID"
	ViolationParticipationDuplicate   apiproblem.ViolationCode = "VS0_PARTICIPATION_DUPLICATE"
	ViolationStatusFieldWrite         apiproblem.ViolationCode = "VS0_STATUS_FIELD_WRITE"
	ViolationSystemOwnedFieldWrite    apiproblem.ViolationCode = "VS0_SYSTEM_OWNED_FIELD_WRITE"
)

// Shared length limits from VS-000 schemaDefaults / VS0-SCHEMA-008..014.
const (
	MaxStringChars       = 253
	MaxIdempotencyKeyLen = 128
	MaxOperatingMarkets  = 32
	MaxAdminAreaCodeLen  = 16
	MaxConditions        = 32
)

// Kind identifies an owned FEATURE-0015 resource kind for validation.
type Kind string

const (
	KindCloudPlatform              Kind = model.KindCloudPlatform
	KindCloudProvider              Kind = model.KindCloudProvider
	KindCloudProviderParticipation Kind = model.KindCloudProviderParticipation
	KindHostingLocation            Kind = model.KindHostingLocation
	KindDatacenter                 Kind = model.KindDatacenter
	KindFaultDomain                Kind = model.KindFaultDomain
	KindInfrastructureStack        Kind = model.KindInfrastructureStack
)

// Classification is the closed create-contract field disposition for one kind.
type Classification struct {
	ClientRequired []string
	ClientOptional []string
	// Deferred are FEATURE-0021-owned fields rejected as UNKNOWN_FIELD.
	Deferred []string
}

// CreateClassification returns the closed create-contract allowlists for kind
// (ADH-2026-047 decision 1; VS0-SCHEMA-008..014 requestContract).
func CreateClassification(kind Kind) (Classification, bool) {
	switch kind {
	case KindCloudPlatform:
		return Classification{
			ClientRequired: []string{
				"metadata.name",
				"spec.ownerRegistration.legalName",
				"spec.ownerRegistration.registrationIdentifier",
				"spec.ownerRegistration.jurisdictionCode",
			},
			ClientOptional: []string{"metadata.displayName", "spec.description"},
		}, true
	case KindCloudProvider:
		return Classification{
			ClientRequired: []string{"metadata.name", "spec.operatingMarkets"},
			ClientOptional: []string{"metadata.displayName", "spec.displayName"},
		}, true
	case KindCloudProviderParticipation:
		return Classification{
			ClientRequired: []string{
				"metadata.name",
				"spec.cloudPlatformRef",
				"spec.cloudProviderRef",
				"spec.environment",
			},
			ClientOptional: nil,
			Deferred: []string{
				"spec.providerSelectionModes",
				"spec.permittedHostingLocationRefs",
			},
		}, true
	case KindHostingLocation:
		return Classification{
			ClientRequired: []string{"metadata.name", "spec.countryCode", "spec.locality"},
			ClientOptional: []string{"spec.administrativeAreaCode", "spec.description"},
		}, true
	case KindDatacenter:
		return Classification{
			ClientRequired: []string{"metadata.name", "spec.hostingLocationRef"},
			ClientOptional: []string{"spec.description"},
		}, true
	case KindFaultDomain:
		return Classification{
			ClientRequired: []string{"metadata.name", "spec.datacenterRef"},
			ClientOptional: []string{"spec.description"},
		}, true
	case KindInfrastructureStack:
		return Classification{
			ClientRequired: []string{"metadata.name", "spec.faultDomainRef"},
			ClientOptional: []string{"spec.description"},
		}, true
	default:
		return Classification{}, false
	}
}

// PatchAllowlist returns the exact PATCHable field paths for kind (design §4.2).
func PatchAllowlist(kind Kind) ([]string, bool) {
	switch kind {
	case KindCloudPlatform:
		return []string{"spec.description"}, true
	case KindCloudProvider:
		return []string{"spec.displayName", "spec.operatingMarkets"}, true
	case KindHostingLocation, KindDatacenter, KindFaultDomain, KindInfrastructureStack:
		return []string{"spec.description"}, true
	case KindCloudProviderParticipation:
		return nil, true // no PATCH surface
	default:
		return nil, false
	}
}

var systemOwnedMetadataFields = []string{
	"uid", "generation", "resourceVersion", "createdAt", "updatedAt",
}

// ClassifyCreateContract classifies a raw JSON object against the closed
// create contract for kind. It rejects server-owned metadata, status,
// client-supplied scopeRef, deferred FEATURE-0021 fields, and unknown fields.
// Missing client-required fields are reported as VALIDATION_FAILED.
//
// The function is pure: it does not mutate data and does not depend on map
// iteration order for emitted violation ordering (paths are sorted).
func ClassifyCreateContract(kind Kind, data []byte) *apiproblem.Problem {
	class, ok := CreateClassification(kind)
	if !ok {
		return apiproblem.New(apiproblem.CodeValidationFailed).WithDetail("unknown resource kind")
	}
	root, prob := decodeJSONObject(data)
	if prob != nil {
		return prob
	}

	clientFields := make(map[string]struct{}, len(class.ClientRequired)+len(class.ClientOptional))
	for _, p := range class.ClientRequired {
		clientFields[p] = struct{}{}
	}
	for _, p := range class.ClientOptional {
		clientFields[p] = struct{}{}
	}

	present := collectLeafPaths(root, "")
	sort.Strings(present)

	var unknown []string
	var deferred []string
	var systemOwned []string
	var statusPaths []string
	var scopePaths []string
	presentSet := make(map[string]struct{}, len(present))
	for _, p := range present {
		presentSet[p] = struct{}{}
		switch {
		case p == "status" || strings.HasPrefix(p, "status."):
			statusPaths = append(statusPaths, p)
		case p == "metadata.scopeRef" || strings.HasPrefix(p, "metadata.scopeRef."):
			scopePaths = append(scopePaths, p)
		case isSystemOwnedMetadataPath(p):
			systemOwned = append(systemOwned, p)
		case isDeferredPath(p, class.Deferred):
			deferred = append(deferred, p)
		case !isAllowedContractPath(p, clientFields, true):
			unknown = append(unknown, p)
		}
	}

	// Audited-denial categories surface as AUTHORIZATION_DENIED with VS0 codes.
	if len(statusPaths) > 0 {
		return apiproblem.New(apiproblem.CodeAuthorizationDenied).WithViolations([]apiproblem.Violation{{
			Field:   "/status",
			Code:    ViolationStatusFieldWrite,
			Message: "status is api-server-owned and cannot be client-authored",
		}})
	}
	if len(systemOwned) > 0 {
		return apiproblem.New(apiproblem.CodeAuthorizationDenied).WithViolations([]apiproblem.Violation{{
			Field:   "/" + strings.ReplaceAll(systemOwned[0], ".", "/"),
			Code:    ViolationSystemOwnedFieldWrite,
			Message: "api-server-owned metadata field cannot be client-authored",
		}})
	}
	if len(deferred) > 0 {
		return apiproblem.New(apiproblem.CodeUnknownField).WithViolations([]apiproblem.Violation{{
			Field:   "/" + strings.ReplaceAll(deferred[0], ".", "/"),
			Code:    apiproblem.ViolationUnknownField,
			Message: "field is owned by another feature and is not accepted",
		}})
	}
	if len(unknown) > 0 {
		return apiproblem.New(apiproblem.CodeUnknownField).WithViolations([]apiproblem.Violation{{
			Field:   "/" + strings.ReplaceAll(unknown[0], ".", "/"),
			Code:    apiproblem.ViolationUnknownField,
			Message: "unknown field",
		}})
	}
	if len(scopePaths) > 0 {
		return apiproblem.New(apiproblem.CodeValidationFailed).WithViolations([]apiproblem.Violation{{
			Field:   "/metadata/scopeRef",
			Code:    ViolationScopeKindInvalid,
			Message: "client cannot supply or select metadata.scopeRef",
		}})
	}

	for _, req := range class.ClientRequired {
		if !pathPresent(presentSet, req) {
			return apiproblem.New(apiproblem.CodeValidationFailed).WithViolations([]apiproblem.Violation{{
				Field:   "/" + strings.ReplaceAll(req, ".", "/"),
				Code:    apiproblem.ViolationOutOfRange,
				Message: "required field is missing",
			}})
		}
	}
	return nil
}

// ClassifyPatchContract classifies a merge-patch JSON object against the
// per-kind PATCH allowlist before RFC 7396 merge (design §6.1.1).
func ClassifyPatchContract(kind Kind, data []byte) *apiproblem.Problem {
	allow, ok := PatchAllowlist(kind)
	if !ok {
		return apiproblem.New(apiproblem.CodeValidationFailed).WithDetail("unknown resource kind")
	}
	if kind == KindCloudProviderParticipation {
		return apiproblem.New(apiproblem.CodeValidationFailed).WithDetail("CloudProviderParticipation has no PATCH surface")
	}
	root, prob := decodeJSONObject(data)
	if prob != nil {
		return prob
	}
	clientFields := make(map[string]struct{}, len(allow))
	for _, p := range allow {
		clientFields[p] = struct{}{}
	}

	present := collectLeafPaths(root, "")
	sort.Strings(present)
	for _, p := range present {
		switch {
		case p == "status" || strings.HasPrefix(p, "status."):
			return apiproblem.New(apiproblem.CodeAuthorizationDenied).WithViolations([]apiproblem.Violation{{
				Field:   "/status",
				Code:    ViolationStatusFieldWrite,
				Message: "status is api-server-owned and cannot be client-authored",
			}})
		case isSystemOwnedMetadataPath(p):
			return apiproblem.New(apiproblem.CodeAuthorizationDenied).WithViolations([]apiproblem.Violation{{
				Field:   "/" + strings.ReplaceAll(p, ".", "/"),
				Code:    ViolationSystemOwnedFieldWrite,
				Message: "api-server-owned metadata field cannot be client-authored",
			}})
		case p == "metadata.scopeRef" || strings.HasPrefix(p, "metadata.scopeRef."):
			return apiproblem.New(apiproblem.CodeValidationFailed).WithViolations([]apiproblem.Violation{{
				Field:   "/metadata/scopeRef",
				Code:    ViolationPatchImmutableField,
				Message: "metadata.scopeRef is immutable",
			}})
		case p == "metadata.name" || strings.HasPrefix(p, "metadata.name."):
			return apiproblem.New(apiproblem.CodeValidationFailed).WithViolations([]apiproblem.Violation{{
				Field:   "/metadata/name",
				Code:    ViolationPatchImmutableField,
				Message: "metadata.name is immutable identity",
			}})
		case strings.HasPrefix(p, "spec.ownerRegistration"):
			return apiproblem.New(apiproblem.CodeValidationFailed).WithViolations([]apiproblem.Violation{{
				Field:   "/spec/ownerRegistration",
				Code:    ViolationPatchImmutableField,
				Message: "ownerRegistration is immutable",
			}})
		case isImmutableReferencePath(kind, p):
			return apiproblem.New(apiproblem.CodeValidationFailed).WithViolations([]apiproblem.Violation{{
				Field:   "/" + strings.ReplaceAll(p, ".", "/"),
				Code:    ViolationPatchImmutableField,
				Message: "relationship reference is immutable",
			}})
		case p == "apiVersion" || p == "kind" || strings.HasPrefix(p, "metadata."):
			return apiproblem.New(apiproblem.CodeValidationFailed).WithViolations([]apiproblem.Violation{{
				Field:   "/" + strings.ReplaceAll(p, ".", "/"),
				Code:    ViolationPatchImmutableField,
				Message: "field is not PATCHable",
			}})
		case !isAllowedContractPath(p, clientFields, false):
			return apiproblem.New(apiproblem.CodeUnknownField).WithViolations([]apiproblem.Violation{{
				Field:   "/" + strings.ReplaceAll(p, ".", "/"),
				Code:    apiproblem.ViolationUnknownField,
				Message: "unknown or non-PATCHable field",
			}})
		}
	}
	return nil
}

func isImmutableReferencePath(kind Kind, path string) bool {
	switch kind {
	case KindDatacenter:
		return path == "spec.hostingLocationRef" || strings.HasPrefix(path, "spec.hostingLocationRef.")
	case KindFaultDomain:
		return path == "spec.datacenterRef" || strings.HasPrefix(path, "spec.datacenterRef.")
	case KindInfrastructureStack:
		return path == "spec.faultDomainRef" || strings.HasPrefix(path, "spec.faultDomainRef.")
	case KindCloudProviderParticipation:
		return path == "spec.cloudPlatformRef" || strings.HasPrefix(path, "spec.cloudPlatformRef.") ||
			path == "spec.cloudProviderRef" || strings.HasPrefix(path, "spec.cloudProviderRef.") ||
			path == "spec.environment"
	default:
		return false
	}
}

func isSystemOwnedMetadataPath(path string) bool {
	for _, f := range systemOwnedMetadataFields {
		if path == "metadata."+f || strings.HasPrefix(path, "metadata."+f+".") {
			return true
		}
	}
	return false
}

func isDeferredPath(path string, deferred []string) bool {
	for _, d := range deferred {
		if path == d || strings.HasPrefix(path, d+".") {
			return true
		}
	}
	return false
}

// isAllowedContractPath reports whether a collected leaf path is within the
// closed client field set. Bare structural containers (metadata/spec) are
// accepted only as empty objects; they never authorize arbitrary children.
func isAllowedContractPath(path string, clientFields map[string]struct{}, allowTypeMeta bool) bool {
	if allowTypeMeta && (path == "apiVersion" || path == "kind") {
		return true
	}
	switch path {
	case "metadata", "spec", "spec.ownerRegistration":
		return true
	}
	if _, ok := clientFields[path]; ok {
		return true
	}
	for a := range clientFields {
		if strings.HasPrefix(path, a+".") {
			return true
		}
	}
	return false
}

func pathPresent(present map[string]struct{}, req string) bool {
	if _, ok := present[req]; ok {
		return true
	}
	prefix := req + "."
	for p := range present {
		if strings.HasPrefix(p, prefix) {
			return true
		}
	}
	// operatingMarkets may be present as an array leaf.
	if _, ok := present[req]; ok {
		return true
	}
	return false
}

func decodeJSONObject(data []byte) (map[string]any, *apiproblem.Problem) {
	if len(bytesTrimSpace(data)) == 0 {
		return nil, apiproblem.New(apiproblem.CodeMalformedRequest).WithDetail("request body is required")
	}
	var root any
	if err := json.Unmarshal(data, &root); err != nil {
		return nil, apiproblem.New(apiproblem.CodeMalformedRequest).WithDetail("malformed JSON")
	}
	obj, ok := root.(map[string]any)
	if !ok {
		return nil, apiproblem.New(apiproblem.CodeMalformedRequest).WithDetail("request body must be a JSON object")
	}
	return obj, nil
}

func collectLeafPaths(v any, prefix string) []string {
	switch t := v.(type) {
	case map[string]any:
		if len(t) == 0 && prefix != "" {
			return []string{prefix}
		}
		keys := make([]string, 0, len(t))
		for k := range t {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		var out []string
		for _, k := range keys {
			child := k
			if prefix != "" {
				child = prefix + "." + k
			}
			out = append(out, collectLeafPaths(t[k], child)...)
		}
		return out
	case []any:
		if prefix != "" {
			return []string{prefix}
		}
		return nil
	default:
		if prefix == "" {
			return nil
		}
		return []string{prefix}
	}
}

func bytesTrimSpace(b []byte) []byte {
	return []byte(strings.TrimSpace(string(b)))
}
