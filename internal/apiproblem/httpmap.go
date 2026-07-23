package apiproblem

import "net/http"

// FailureClass is the F12-ERROR-002 baseline failure taxonomy.
// Mapping is generic only: Decision/SafeDenial live in apivalid, not here.
type FailureClass string

const (
	FailureMalformedRequest       FailureClass = "malformed_request"
	FailureAuthenticationRequired FailureClass = "authentication_required"
	FailureAuthorizationDenied    FailureClass = "authorization_denied"
	FailureNotFound               FailureClass = "not_found"
	FailureConflict               FailureClass = "lifecycle_uniqueness_conflict"
	FailureStaleResourceVersion   FailureClass = "stale_resource_version"
	FailureUnsupportedMediaType   FailureClass = "unsupported_media_type"
	FailureValidationFailed       FailureClass = "structurally_semantically_invalid"
	FailureInternal               FailureClass = "internal_failure"
	FailureDependencyUnavailable  FailureClass = "temporary_dependency_failure"
)

// AllFailureClasses returns the closed F12-ERROR-002 failure-class set in
// stable order.
func AllFailureClasses() []FailureClass {
	return []FailureClass{
		FailureMalformedRequest,
		FailureAuthenticationRequired,
		FailureAuthorizationDenied,
		FailureNotFound,
		FailureConflict,
		FailureStaleResourceVersion,
		FailureUnsupportedMediaType,
		FailureValidationFailed,
		FailureInternal,
		FailureDependencyUnavailable,
	}
}

// Valid reports whether c is a registered FailureClass.
func (c FailureClass) Valid() bool {
	switch c {
	case FailureMalformedRequest,
		FailureAuthenticationRequired,
		FailureAuthorizationDenied,
		FailureNotFound,
		FailureConflict,
		FailureStaleResourceVersion,
		FailureUnsupportedMediaType,
		FailureValidationFailed,
		FailureInternal,
		FailureDependencyUnavailable:
		return true
	default:
		return false
	}
}

// StatusFor returns the baseline HTTP status for a FailureClass
// (F12-ERROR-002). Unknown classes map to 500.
func StatusFor(class FailureClass) int {
	switch class {
	case FailureMalformedRequest:
		return http.StatusBadRequest // 400
	case FailureAuthenticationRequired:
		return http.StatusUnauthorized // 401
	case FailureAuthorizationDenied:
		return http.StatusForbidden // 403
	case FailureNotFound:
		return http.StatusNotFound // 404
	case FailureConflict:
		return http.StatusConflict // 409
	case FailureStaleResourceVersion:
		return http.StatusPreconditionFailed // 412
	case FailureUnsupportedMediaType:
		return http.StatusUnsupportedMediaType // 415
	case FailureValidationFailed:
		return http.StatusUnprocessableEntity // 422
	case FailureInternal:
		return http.StatusInternalServerError // 500
	case FailureDependencyUnavailable:
		return http.StatusServiceUnavailable // 503
	default:
		return http.StatusInternalServerError
	}
}

// ClassForCode returns the FailureClass for a baseline ErrorCode.
// Unknown codes map to FailureInternal.
func ClassForCode(code ErrorCode) FailureClass {
	switch code {
	case CodeMalformedRequest, CodeUnknownField, CodeDuplicateField, CodeRequestTooLarge:
		return FailureMalformedRequest
	case CodeAuthRequired:
		return FailureAuthenticationRequired
	case CodeAuthorizationDenied:
		return FailureAuthorizationDenied
	case CodeResourceNotFound:
		return FailureNotFound
	case CodeConflict, CodeAlreadyExists, CodeDeleteBlocked:
		return FailureConflict
	case CodeStaleResourceVersion:
		return FailureStaleResourceVersion
	case CodeUnsupportedMediaType:
		return FailureUnsupportedMediaType
	case CodeValidationFailed:
		return FailureValidationFailed
	case CodeInternalError:
		return FailureInternal
	case CodeDependencyUnavailable:
		return FailureDependencyUnavailable
	default:
		return FailureInternal
	}
}

// StatusForCode returns the baseline HTTP status for a Problem ErrorCode
// (F12-ERROR-002). Oversized bodies use 400 REQUEST_TOO_LARGE (not 413).
func StatusForCode(code ErrorCode) int {
	return StatusFor(ClassForCode(code))
}
