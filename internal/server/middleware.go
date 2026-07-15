package server

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/sanjeevksaini/sovrunn/internal/requestctx"
	"github.com/sanjeevksaini/sovrunn/internal/resources"
)

const headerRequestID = "X-Sovrunn-Request-ID"

// requestIDMiddleware reads X-Sovrunn-Request-ID from the request.
// If absent or empty, it generates a new ID. The resolved ID is stored
// in the request context and written to the response header.
func requestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := strings.TrimSpace(r.Header.Get(headerRequestID))
		if reqID == "" {
			reqID = generateRequestID()
		}
		w.Header().Set(headerRequestID, reqID)
		ctx := requestctx.WithRequestID(r.Context(), reqID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func generateRequestID() string {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		return hex.EncodeToString([]byte("fallback-request-id"))
	}
	return hex.EncodeToString(b[:])
}

// contentTypeMiddleware rejects requests whose Content-Type header does
// not equal "application/json" for methods that carry a request body.
func contentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost, http.MethodPut, http.MethodPatch:
			ct := r.Header.Get("Content-Type")
			if ct != "application/json" && !strings.HasPrefix(ct, "application/json;") {
				w.Header().Set("Content-Type", "application/json")
				if reqID := requestctx.RequestIDFromContext(r.Context()); reqID != "" {
					w.Header().Set(headerRequestID, reqID)
				}
				w.WriteHeader(http.StatusUnsupportedMediaType)
				_ = writeErrorBody(w, resources.ErrCodeValidationFailed, "content type must be application/json", "", "")
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

// loggingMiddleware logs request_id, method, path, status_code, latency_ms.
func loggingMiddleware(logger *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
			next.ServeHTTP(rec, r)
			latency := time.Since(start).Milliseconds()
			reqID := requestctx.RequestIDFromContext(r.Context())
			if rec.status >= 400 {
				logger.Printf("request_id=%s method=%s path=%s status_code=%d latency_ms=%d error_code=%s",
					reqID, r.Method, r.URL.Path, rec.status, latency, errorCodeForStatus(rec.status))
			} else {
				logger.Printf("request_id=%s method=%s path=%s status_code=%d latency_ms=%d",
					reqID, r.Method, r.URL.Path, rec.status, latency)
			}
		})
	}
}

func errorCodeForStatus(status int) string {
	switch status {
	case http.StatusBadRequest, http.StatusUnsupportedMediaType, http.StatusRequestEntityTooLarge:
		return string(resources.ErrCodeValidationFailed)
	case http.StatusNotFound:
		return string(resources.ErrCodeResourceNotFound)
	case http.StatusConflict:
		return string(resources.ErrCodeResourceAlreadyExists)
	case http.StatusInternalServerError:
		return string(resources.ErrCodeInternalError)
	default:
		return string(resources.ErrCodeValidationFailed)
	}
}
