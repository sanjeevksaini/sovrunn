package resources

// ErrorCode is a stable, versioned string code used in APIError
// responses and registry sentinel errors.
type ErrorCode string

const (
	ErrCodeValidationFailed      ErrorCode = "VALIDATION_FAILED"
	ErrCodeResourceNotFound      ErrorCode = "RESOURCE_NOT_FOUND"
	ErrCodeResourceAlreadyExists ErrorCode = "RESOURCE_ALREADY_EXISTS"
	ErrCodeDeleteBlocked         ErrorCode = "DELETE_BLOCKED"
	ErrCodeMethodNotAllowed      ErrorCode = "METHOD_NOT_ALLOWED"
	ErrCodeInternalError         ErrorCode = "INTERNAL_ERROR"
)

// APIError is the wire shape for all error responses.
// The HTTP response body is always: {"error": <APIError>}.
type APIError struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
	Field   string    `json:"field,omitempty"`
	Details string    `json:"details,omitempty"`
}

// APIErrorEnvelope wraps APIError so the JSON shape is
// {"error": {"code": ..., "message": ...}}.
type APIErrorEnvelope struct {
	Error APIError `json:"error"`
}

// FieldError is returned by ValidateOrganization for each invalid field.
// Field is the dot-separated JSON path (e.g. "metadata.name").
type FieldError struct {
	Field   string
	Message string
}
