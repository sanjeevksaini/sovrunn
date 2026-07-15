package server

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRequestIDMiddleware_GeneratesWhenAbsent(t *testing.T) {
	var captured string
	h := requestIDMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		captured = r.Header.Get(headerRequestID)
	}))
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Header().Get(headerRequestID) == "" {
		t.Fatal("expected generated request ID in response header")
	}
	_ = captured
}

func TestRequestIDMiddleware_PropagatesWhenPresent(t *testing.T) {
	h := requestIDMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(headerRequestID, "custom-id-123")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if got := rec.Header().Get(headerRequestID); got != "custom-id-123" {
		t.Errorf("got %q, want custom-id-123", got)
	}
}

func TestContentTypeMiddleware_RejectsPOSTWithoutJSON(t *testing.T) {
	h := requestIDMiddleware(contentTypeMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})))
	req := httptest.NewRequest(http.MethodPost, "/v1/organizations", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "text/plain")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnsupportedMediaType {
		t.Fatalf("status = %d, want 415", rec.Code)
	}
}

func TestContentTypeMiddleware_PassThroughGET(t *testing.T) {
	called := false
	h := requestIDMiddleware(contentTypeMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	})))
	req := httptest.NewRequest(http.MethodGet, "/v1/organizations", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if !called {
		t.Fatal("handler not called for GET")
	}
}

func TestLoggingMiddleware_WritesStructuredLog(t *testing.T) {
	var buf bytes.Buffer
	logger := log.New(&buf, "", 0)
	h := requestIDMiddleware(loggingMiddleware(logger)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})))
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	out := buf.String()
	for _, field := range []string{"request_id=", "method=GET", "path=/healthz", "status_code=200", "latency_ms="} {
		if !strings.Contains(out, field) {
			t.Errorf("log missing %q: %s", field, out)
		}
	}
}
