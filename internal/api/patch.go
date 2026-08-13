package api

import (
	"encoding/json"

	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/cloudmodel/model"
	"github.com/sanjeevksaini/sovrunn/internal/cloudmodel/validate"
)

// ApplyMergePatch applies RFC 7396 merge-patch to a staged clone of the
// current resource, then runs post-merge schema/immutable/system-owned/
// deferred/unknown-field validation before publication (design §6.1.1).
//
// Before merge, every patch member is classified against the per-kind
// allowlist. null removes only optional mutable spec.description or
// spec.displayName; null for required spec.operatingMarkets is
// VALIDATION_FAILED.
func ApplyMergePatch(kind validate.Kind, current any, patch []byte) (any, *apiproblem.Problem) {
	if prob := validate.ClassifyPatchContract(kind, patch); prob != nil {
		return nil, prob
	}
	var patchObj map[string]any
	if err := json.Unmarshal(patch, &patchObj); err != nil {
		return nil, apiproblem.New(apiproblem.CodeMalformedRequest).WithDetail("malformed JSON")
	}
	if prob := rejectDisallowedNulls(kind, patchObj); prob != nil {
		return nil, prob
	}

	currentBytes, err := json.Marshal(current)
	if err != nil {
		return nil, apiproblem.New(apiproblem.CodeInternalError).WithDetail("failed to clone resource")
	}
	var staged map[string]any
	if err := json.Unmarshal(currentBytes, &staged); err != nil {
		return nil, apiproblem.New(apiproblem.CodeInternalError).WithDetail("failed to clone resource")
	}

	merged := mergePatch(staged, patchObj)
	mergedBytes, err := json.Marshal(merged)
	if err != nil {
		return nil, apiproblem.New(apiproblem.CodeInternalError).WithDetail("failed to encode merged resource")
	}

	after, prob := decodeKind(kind, mergedBytes)
	if prob != nil {
		return nil, prob
	}
	if prob := validate.ValidateImmutablePatch(kind, current, after); prob != nil {
		return nil, prob
	}
	if prob := validate.ValidateSchema(kind, after); prob != nil {
		return nil, prob
	}
	return after, nil
}

func rejectDisallowedNulls(kind validate.Kind, patch map[string]any) *apiproblem.Problem {
	spec, _ := patch["spec"].(map[string]any)
	if spec == nil {
		return nil
	}
	if v, ok := spec["operatingMarkets"]; ok && v == nil {
		return apiproblem.New(apiproblem.CodeValidationFailed).WithViolations([]apiproblem.Violation{{
			Field:   "/spec/operatingMarkets",
			Code:    apiproblem.ViolationOutOfRange,
			Message: "required operatingMarkets cannot be removed",
		}})
	}
	// null is permitted only for optional mutable description/displayName.
	for key, v := range spec {
		if v != nil {
			continue
		}
		switch key {
		case "description":
			if kind == validate.KindCloudProvider {
				return apiproblem.New(apiproblem.CodeUnknownField).WithViolations([]apiproblem.Violation{{
					Field:   "/spec/description",
					Code:    apiproblem.ViolationUnknownField,
					Message: "spec.description is not PATCHable for CloudProvider",
				}})
			}
		case "displayName":
			if kind != validate.KindCloudProvider {
				return apiproblem.New(apiproblem.CodeUnknownField).WithViolations([]apiproblem.Violation{{
					Field:   "/spec/displayName",
					Code:    apiproblem.ViolationUnknownField,
					Message: "spec.displayName is not PATCHable for this kind",
				}})
			}
		default:
			return apiproblem.New(apiproblem.CodeValidationFailed).WithViolations([]apiproblem.Violation{{
				Field:   "/spec/" + key,
				Code:    validate.ViolationPatchImmutableField,
				Message: "null removal is not permitted for this field",
			}})
		}
	}
	return nil
}

// mergePatch implements RFC 7396 against generic JSON values.
func mergePatch(target, patch any) any {
	patchObj, patchIsObj := patch.(map[string]any)
	if !patchIsObj {
		return patch
	}
	targetObj, targetIsObj := target.(map[string]any)
	if !targetIsObj {
		targetObj = map[string]any{}
	}
	// Copy target keys so the staged clone is not mutated in place unexpectedly.
	out := make(map[string]any, len(targetObj)+len(patchObj))
	for k, v := range targetObj {
		out[k] = v
	}
	for k, v := range patchObj {
		if v == nil {
			delete(out, k)
			continue
		}
		if _, vObj := v.(map[string]any); vObj {
			out[k] = mergePatch(out[k], v)
			continue
		}
		out[k] = v
	}
	return out
}

func decodeKind(kind validate.Kind, data []byte) (any, *apiproblem.Problem) {
	switch kind {
	case validate.KindCloudPlatform:
		var v model.CloudPlatform
		if err := json.Unmarshal(data, &v); err != nil {
			return nil, apiproblem.New(apiproblem.CodeMalformedRequest).WithDetail("malformed merged resource")
		}
		return v, nil
	case validate.KindCloudProvider:
		var v model.CloudProvider
		if err := json.Unmarshal(data, &v); err != nil {
			return nil, apiproblem.New(apiproblem.CodeMalformedRequest).WithDetail("malformed merged resource")
		}
		return v, nil
	case validate.KindHostingLocation:
		var v model.HostingLocation
		if err := json.Unmarshal(data, &v); err != nil {
			return nil, apiproblem.New(apiproblem.CodeMalformedRequest).WithDetail("malformed merged resource")
		}
		return v, nil
	case validate.KindDatacenter:
		var v model.Datacenter
		if err := json.Unmarshal(data, &v); err != nil {
			return nil, apiproblem.New(apiproblem.CodeMalformedRequest).WithDetail("malformed merged resource")
		}
		return v, nil
	case validate.KindFaultDomain:
		var v model.FaultDomain
		if err := json.Unmarshal(data, &v); err != nil {
			return nil, apiproblem.New(apiproblem.CodeMalformedRequest).WithDetail("malformed merged resource")
		}
		return v, nil
	case validate.KindInfrastructureStack:
		var v model.InfrastructureStack
		if err := json.Unmarshal(data, &v); err != nil {
			return nil, apiproblem.New(apiproblem.CodeMalformedRequest).WithDetail("malformed merged resource")
		}
		return v, nil
	default:
		return nil, apiproblem.New(apiproblem.CodeValidationFailed).WithDetail("unsupported PATCH kind")
	}
}
