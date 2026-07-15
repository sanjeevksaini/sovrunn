package api

import (
	"encoding/json"
	"net/http"

	"github.com/sanjeevksaini/sovrunn/internal/requestctx"
	"github.com/sanjeevksaini/sovrunn/internal/resources"
)

// writeError writes a JSON-encoded APIErrorEnvelope with the given HTTP
// status code. It always sets Content-Type: application/json and the
// X-Sovrunn-Request-ID header.
func writeError(
	w http.ResponseWriter,
	r *http.Request,
	status int,
	code resources.ErrorCode,
	message, field, details string,
) {
	envelope := resources.APIErrorEnvelope{
		Error: resources.APIError{
			Code:    code,
			Message: message,
			Field:   field,
			Details: details,
		},
	}
	writeJSON(w, r, status, envelope)
}

// writeJSON writes v as JSON with the given status code.
func writeJSON(w http.ResponseWriter, r *http.Request, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	if reqID := requestctx.RequestIDFromContext(r.Context()); reqID != "" {
		w.Header().Set("X-Sovrunn-Request-ID", reqID)
	}
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
