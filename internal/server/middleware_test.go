package server

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/resources"
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

func apiChain(logger *log.Logger, handler http.Handler) http.Handler {
	return requestIDMiddleware(loggingMiddleware(logger)(contentTypeMiddleware(handler)))
}

func TestContentTypeMiddleware_RejectsPOSTWithoutJSON(t *testing.T) {
	h := apiChain(log.Default(), http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	req := httptest.NewRequest(http.MethodPost, "/v1/organizations", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "text/plain")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnsupportedMediaType {
		t.Fatalf("status = %d, want 415", rec.Code)
	}
}

func TestContentTypeMiddleware_RejectsPOSTWithoutJSON_LogsRequest(t *testing.T) {
	var buf bytes.Buffer
	logger := log.New(&buf, "", 0)
	h := apiChain(logger, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("handler should not be called for rejected content type")
	}))
	req := httptest.NewRequest(http.MethodPost, "/v1/organizations", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "text/plain")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnsupportedMediaType {
		t.Fatalf("status = %d, want 415", rec.Code)
	}
	out := buf.String()
	for _, field := range []string{"request_id=", "method=POST", "path=/v1/organizations", "status_code=415", "latency_ms=", "error_code=VALIDATION_FAILED"} {
		if !strings.Contains(out, field) {
			t.Errorf("log missing %q: %s", field, out)
		}
	}
}

func TestContentTypeMiddleware_RejectsPOSTWithoutJSON_IncludesRequestID(t *testing.T) {
	h := apiChain(log.Default(), http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("handler should not be called for rejected content type")
	}))
	req := httptest.NewRequest(http.MethodPost, "/v1/organizations", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "text/plain")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnsupportedMediaType {
		t.Fatalf("status = %d, want 415", rec.Code)
	}
	if rec.Header().Get(headerRequestID) == "" {
		t.Fatal("expected X-Sovrunn-Request-ID on 415 response")
	}
}

func TestContentTypeMiddleware_PassThroughGET(t *testing.T) {
	called := false
	h := apiChain(log.Default(), http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	}))
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

func TestMethodGET_AllowsGET(t *testing.T) {
	called := false
	h := methodGET(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	}))
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if !called {
		t.Fatal("handler not called for GET")
	}
	if rec.Code == http.StatusMethodNotAllowed {
		t.Fatal("GET must not return 405")
	}
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
}

func TestMethodGET_RejectsPOST(t *testing.T) {
	assertMethodGETRejectsWithBody(t, http.MethodPost)
}

func TestMethodGET_RejectsPUT(t *testing.T) {
	assertMethodGETRejectsWithBody(t, http.MethodPut)
}

func TestMethodGET_RejectsDELETE(t *testing.T) {
	assertMethodGETRejectsWithBody(t, http.MethodDelete)
}

func TestMethodGET_RejectsPATCH(t *testing.T) {
	assertMethodGETRejectsWithBody(t, http.MethodPatch)
}

func TestMethodGET_RejectsHEAD(t *testing.T) {
	h := methodGET(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("handler should not be called for HEAD")
	}))
	req := httptest.NewRequest(http.MethodHead, "/healthz", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status = %d, want 405", rec.Code)
	}
	if got := rec.Header().Get("Allow"); got != "GET" {
		t.Errorf("Allow = %q, want GET", got)
	}
	if got := rec.Header().Get("Content-Type"); got != "application/json" {
		t.Errorf("Content-Type = %q, want application/json", got)
	}
	if rec.Body.Len() != 0 {
		t.Fatalf("HEAD body length = %d, want 0", rec.Body.Len())
	}
}

func assertMethodGETRejectsWithBody(t *testing.T, method string) {
	t.Helper()
	h := methodGET(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatalf("handler should not be called for %s", method)
	}))
	req := httptest.NewRequest(method, "/healthz", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status = %d, want 405", rec.Code)
	}
	if got := rec.Header().Get("Allow"); got != "GET" {
		t.Errorf("Allow = %q, want GET", got)
	}
	if got := rec.Header().Get("Content-Type"); got != "application/json" {
		t.Errorf("Content-Type = %q, want application/json", got)
	}
	var envelope resources.APIErrorEnvelope
	if err := json.NewDecoder(rec.Body).Decode(&envelope); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if envelope.Error.Code != resources.ErrCodeMethodNotAllowed {
		t.Errorf("code = %q, want METHOD_NOT_ALLOWED", envelope.Error.Code)
	}
	if envelope.Error.Message != "only GET is supported" {
		t.Errorf("message = %q, want %q", envelope.Error.Message, "only GET is supported")
	}
}
