package apiproblem

// ErrorCode is a stable machine-contract problem code (F12-ERROR-003).
// Human detail/message text MUST NOT be treated as a contract.
type ErrorCode string

// Baseline problem-level ErrorCode values from the F12-ERROR-002 table.
const (
	CodeMalformedRequest      ErrorCode = "MALFORMED_REQUEST"
	CodeUnknownField          ErrorCode = "UNKNOWN_FIELD"
	CodeDuplicateField        ErrorCode = "DUPLICATE_FIELD"
	CodeRequestTooLarge       ErrorCode = "REQUEST_TOO_LARGE"
	CodeAuthRequired          ErrorCode = "AUTH_REQUIRED"        // reserved for adopters
	CodeAuthorizationDenied   ErrorCode = "AUTHORIZATION_DENIED" // reserved for adopters
	CodeResourceNotFound      ErrorCode = "RESOURCE_NOT_FOUND"
	CodeConflict              ErrorCode = "CONFLICT"
	CodeAlreadyExists         ErrorCode = "ALREADY_EXISTS"
	CodeDeleteBlocked         ErrorCode = "DELETE_BLOCKED"
	CodeStaleResourceVersion  ErrorCode = "STALE_RESOURCE_VERSION"
	CodeUnsupportedMediaType  ErrorCode = "UNSUPPORTED_MEDIA_TYPE"
	CodeValidationFailed      ErrorCode = "VALIDATION_FAILED"
	CodeInternalError         ErrorCode = "INTERNAL_ERROR"
	CodeDependencyUnavailable ErrorCode = "DEPENDENCY_UNAVAILABLE"
)

// AllErrorCodes returns the closed baseline ErrorCode set in stable order.
func AllErrorCodes() []ErrorCode {
	return []ErrorCode{
		CodeMalformedRequest,
		CodeUnknownField,
		CodeDuplicateField,
		CodeRequestTooLarge,
		CodeAuthRequired,
		CodeAuthorizationDenied,
		CodeResourceNotFound,
		CodeConflict,
		CodeAlreadyExists,
		CodeDeleteBlocked,
		CodeStaleResourceVersion,
		CodeUnsupportedMediaType,
		CodeValidationFailed,
		CodeInternalError,
		CodeDependencyUnavailable,
	}
}

// Valid reports whether c is a registered baseline ErrorCode.
func (c ErrorCode) Valid() bool {
	switch c {
	case CodeMalformedRequest,
		CodeUnknownField,
		CodeDuplicateField,
		CodeRequestTooLarge,
		CodeAuthRequired,
		CodeAuthorizationDenied,
		CodeResourceNotFound,
		CodeConflict,
		CodeAlreadyExists,
		CodeDeleteBlocked,
		CodeStaleResourceVersion,
		CodeUnsupportedMediaType,
		CodeValidationFailed,
		CodeInternalError,
		CodeDependencyUnavailable:
		return true
	default:
		return false
	}
}

// TitleFor returns the stable human title for a baseline ErrorCode.
// Titles MAY evolve; codes MUST NOT (F12-ERROR-003).
func TitleFor(code ErrorCode) string {
	switch code {
	case CodeMalformedRequest:
		return "Malformed request"
	case CodeUnknownField:
		return "Unknown field"
	case CodeDuplicateField:
		return "Duplicate field"
	case CodeRequestTooLarge:
		return "Request too large"
	case CodeAuthRequired:
		return "Authentication required"
	case CodeAuthorizationDenied:
		return "Authorization denied"
	case CodeResourceNotFound:
		return "Resource not found"
	case CodeConflict:
		return "Conflict"
	case CodeAlreadyExists:
		return "Already exists"
	case CodeDeleteBlocked:
		return "Delete blocked"
	case CodeStaleResourceVersion:
		return "Stale resource version"
	case CodeUnsupportedMediaType:
		return "Unsupported media type"
	case CodeValidationFailed:
		return "Validation failed"
	case CodeInternalError:
		return "Internal error"
	case CodeDependencyUnavailable:
		return "Dependency unavailable"
	default:
		return "Error"
	}
}

// ViolationCode is a stable field-level machine contract carried on
// Violation.Code (F12-ERROR-001/003).
type ViolationCode string

// Registered violation codes used by the shared grammar (D-05).
// OPERATION_TARGET_SCOPE_MISMATCH is required by D-17 / F12-SCOPE-002.
const (
	ViolationUnknownField                 ViolationCode = "UNKNOWN_FIELD"
	ViolationDuplicateField               ViolationCode = "DUPLICATE_FIELD"
	ViolationOutOfRange                   ViolationCode = "OUT_OF_RANGE"
	ViolationOperationTargetScopeMismatch ViolationCode = "OPERATION_TARGET_SCOPE_MISMATCH"
)

// AllViolationCodes returns the registered violation-code set in stable order.
func AllViolationCodes() []ViolationCode {
	return []ViolationCode{
		ViolationUnknownField,
		ViolationDuplicateField,
		ViolationOutOfRange,
		ViolationOperationTargetScopeMismatch,
	}
}

// Valid reports whether c is a registered ViolationCode.
func (c ViolationCode) Valid() bool {
	switch c {
	case ViolationUnknownField,
		ViolationDuplicateField,
		ViolationOutOfRange,
		ViolationOperationTargetScopeMismatch:
		return true
	default:
		return false
	}
}
