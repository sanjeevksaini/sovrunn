package apiref

import (
	"strconv"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
)

// TypedRef is re-exported from apimeta as the common typed-reference base
// (F12-REF-001, D-16). Singular fields end in Ref; collections end in Refs.
// Provider-native IDs MUST NOT act as core refs (F12-REF-003).
type TypedRef = apimeta.TypedRef

// Refs is a typed-reference collection (F12-REF-001). Collection fields end
// in Refs; singular fields end in Ref.
type Refs []TypedRef

// DefaultMaxRefs is the reviewed default upper bound for references per field
// (design Limits.MaxReferencesPerField). Callers may pass a tighter max to
// Validate; non-positive max selects this default.
const DefaultMaxRefs = 64

// Stable codes for package-local RefIssue values (F12-REF-002/003).
// Translation to apiproblem.Violation is owned by apivalid (task 2.7).
const (
	CodeKindNotAllowed    = "REF_KIND_NOT_ALLOWED"
	CodeScopeNotAllowed   = "REF_SCOPE_NOT_ALLOWED"
	CodeDirectionInvalid  = "REF_DIRECTION_INVALID"
	CodeNameUIDMismatch   = "REF_NAME_UID_MISMATCH"
	CodeProviderNativeID  = "REF_PROVIDER_NATIVE_ID"
	CodeMissingAPIVersion = "REF_MISSING_API_VERSION"
	CodeMissingKind       = "REF_MISSING_KIND"
	CodeMissingName       = "REF_MISSING_NAME"
	CodeRefsExceedLimit   = "REFS_EXCEED_LIMIT"
)

// CheckNameUIDAgreement enforces F12-REF-002 against a resolved object
// identity. UID MAY be omitted on the ref (human-authored input); when both
// name and uid are present on the ref they MUST identify the same resolved
// object. A disagreement is reported as REF_NAME_UID_MISMATCH.
//
// resolvedName and resolvedUID are the authoritative identity of the object
// the caller already resolved (offline fixtures or adopter lookup). Empty
// resolved values skip the check so structural ValidateRef can remain
// store-free.
func CheckNameUIDAgreement(ref TypedRef, path, resolvedName, resolvedUID string) []RefIssue {
	if ref.Name == "" || ref.UID == "" {
		return nil
	}
	if resolvedName == "" || resolvedUID == "" {
		return nil
	}
	if ref.Name == resolvedName && ref.UID == resolvedUID {
		return nil
	}
	return []RefIssue{{
		Path:    path,
		Code:    CodeNameUIDMismatch,
		Message: "reference name and uid must identify the same object",
	}}
}

// Validate applies Constraint to every element and enforces the per-field
// reference count bound (F12-VALIDATION-007). path is the JSON Pointer of the
// collection field (for example "/spec/resourcePoolRefs").
func (r Refs) Validate(c Constraint, path string, max int) []RefIssue {
	if max <= 0 {
		max = DefaultMaxRefs
	}
	var issues []RefIssue
	if len(r) > max {
		issues = append(issues, RefIssue{
			Path:    path,
			Code:    CodeRefsExceedLimit,
			Message: "reference collection exceeds the finite per-field limit",
		})
	}
	for i, ref := range r {
		elemPath := path + "/" + strconv.Itoa(i)
		issues = append(issues, c.ValidateRef(ref, elemPath)...)
	}
	return issues
}
