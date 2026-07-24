package apivalid

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
)

const sampleJSONBody = `{"apiVersion":"platform.sovrunn.io/v1","kind":"Project","metadata":{"name":"demo"},"spec":{"displayName":"Demo"}}`

const sampleYAMLBody = `apiVersion: platform.sovrunn.io/v1
kind: Project
metadata:
  name: demo
spec:
  displayName: Demo
`

func TestStrictDecodeJSONHappyPath(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/apis/test", strings.NewReader(sampleJSONBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	var dst decodeSample
	prob := StrictDecode(rec, req, testLimits, ModeCreateRequest, &dst)
	if prob != nil {
		t.Fatalf("StrictDecode returned problem: %#v", prob)
	}
	if dst.Kind != "Project" || dst.Metadata.Name != "demo" {
		t.Fatalf("unexpected decode result: %#v", dst)
	}
}

func TestStrictDecodeYAMLHappyPath(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/apis/test", strings.NewReader(sampleYAMLBody))
	req.Header.Set("Content-Type", "application/yaml")
	rec := httptest.NewRecorder()

	var dst decodeSample
	prob := StrictDecode(rec, req, testLimits, ModeCreateRequest, &dst)
	if prob != nil {
		t.Fatalf("StrictDecode returned problem: %#v", prob)
	}
	if dst.Kind != "Project" || dst.Metadata.Name != "demo" {
		t.Fatalf("unexpected decode result: %#v", dst)
	}
}

func TestStrictDecodeSelectsYAMLAliases(t *testing.T) {
	cases := []string{"application/x-yaml", "text/yaml", "Application/X-YAML; charset=utf-8"}
	for _, ct := range cases {
		t.Run(ct, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/apis/test", strings.NewReader(sampleYAMLBody))
			req.Header.Set("Content-Type", ct)
			rec := httptest.NewRecorder()

			var dst decodeSample
			prob := StrictDecode(rec, req, testLimits, ModeCreateRequest, &dst)
			if prob != nil {
				t.Fatalf("StrictDecode returned problem for %q: %#v", ct, prob)
			}
			if dst.Kind != "Project" || dst.Metadata.Name != "demo" {
				t.Fatalf("unexpected decode result for %q: %#v", ct, dst)
			}
		})
	}
}

func TestStrictDecodeJSONWithCharset(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/apis/test", strings.NewReader(sampleJSONBody))
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	rec := httptest.NewRecorder()

	var dst decodeSample
	prob := StrictDecode(rec, req, testLimits, ModeCreateRequest, &dst)
	if prob != nil {
		t.Fatalf("StrictDecode returned problem: %#v", prob)
	}
	if dst.Metadata.Name != "demo" {
		t.Fatalf("unexpected decode result: %#v", dst)
	}
}

func TestStrictDecodeUnsupportedMediaType415(t *testing.T) {
	cases := []string{
		"",
		"text/plain",
		"application/xml",
		"multipart/form-data",
		"application/jsonn",
		"not-a-media-type",
	}
	for _, ct := range cases {
		t.Run(ct, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/apis/test", strings.NewReader(sampleJSONBody))
			if ct != "" {
				req.Header.Set("Content-Type", ct)
			}
			rec := httptest.NewRecorder()

			var dst decodeSample
			prob := StrictDecode(rec, req, testLimits, ModeCreateRequest, &dst)
			if prob == nil {
				t.Fatal("expected UNSUPPORTED_MEDIA_TYPE Problem, got nil")
			}
			if prob.Code != apiproblem.CodeUnsupportedMediaType {
				t.Fatalf("Code = %q, want %q", prob.Code, apiproblem.CodeUnsupportedMediaType)
			}
			if prob.Status != http.StatusUnsupportedMediaType {
				t.Fatalf("Status = %d, want %d", prob.Status, http.StatusUnsupportedMediaType)
			}
			if dst.Kind != "" {
				t.Fatalf("dst must remain empty on media-type rejection, got %#v", dst)
			}
		})
	}
}

func TestStrictDecodeOversizedBody400(t *testing.T) {
	lim := DefaultLimits()
	lim.MaxObjectBytes = 64

	body := `{"apiVersion":"v1","kind":"Project","metadata":{"name":"` + strings.Repeat("x", 200) + `"},"spec":{}}`
	req := httptest.NewRequest(http.MethodPost, "/apis/test", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	var dst decodeSample
	prob := StrictDecode(rec, req, lim, ModeCreateRequest, &dst)
	if prob == nil {
		t.Fatal("expected REQUEST_TOO_LARGE Problem, got nil")
	}
	if prob.Code != apiproblem.CodeRequestTooLarge {
		t.Fatalf("Code = %q, want %q", prob.Code, apiproblem.CodeRequestTooLarge)
	}
	if prob.Status != http.StatusBadRequest {
		t.Fatalf("Status = %d, want 400 (not 413)", prob.Status)
	}
	if prob.Status == http.StatusRequestEntityTooLarge {
		t.Fatal("REQUEST_TOO_LARGE must not use 413")
	}
}

func TestStrictDecodeSelectsJSONDecoder(t *testing.T) {
	// YAML-only merge key must not be evaluated when Content-Type is JSON;
	// the JSON decoder rejects it as malformed JSON.
	req := httptest.NewRequest(http.MethodPost, "/apis/test", strings.NewReader("<<: *anchor\n"))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	var dst decodeSample
	prob := StrictDecode(rec, req, testLimits, ModeCreateRequest, &dst)
	if prob == nil {
		t.Fatal("expected JSON decode failure, got nil")
	}
	if prob.Code != apiproblem.CodeMalformedRequest {
		t.Fatalf("Code = %q, want %q (JSON path selected)", prob.Code, apiproblem.CodeMalformedRequest)
	}
}

func TestStrictDecodeSelectsYAMLDecoder(t *testing.T) {
	// Valid JSON bytes are accepted by the YAML path after normalization,
	// proving the YAML decoder was selected for application/yaml.
	req := httptest.NewRequest(http.MethodPost, "/apis/test", strings.NewReader(sampleJSONBody))
	req.Header.Set("Content-Type", "application/yaml")
	rec := httptest.NewRecorder()

	var dst decodeSample
	prob := StrictDecode(rec, req, testLimits, ModeCreateRequest, &dst)
	if prob != nil {
		t.Fatalf("YAML decoder should accept JSON-compatible input: %#v", prob)
	}
	if dst.Kind != "Project" || dst.Metadata.Name != "demo" {
		t.Fatalf("unexpected decode result: %#v", dst)
	}
}

func TestStrictDecodeAppliesFieldPolicy(t *testing.T) {
	const withStatus = `{"apiVersion":"v1","kind":"Project","metadata":{"name":"demo"},"spec":{},"status":{"phase":"Ready"}}`
	req := httptest.NewRequest(http.MethodPost, "/apis/test", strings.NewReader(withStatus))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	var dst decodeSample
	prob := StrictDecode(rec, req, testLimits, ModeCreateRequest, &dst)
	if prob == nil {
		t.Fatal("expected FieldPolicy rejection for status under ModeCreateRequest")
	}
	if prob.Code != apiproblem.CodeValidationFailed {
		t.Fatalf("Code = %q, want %q", prob.Code, apiproblem.CodeValidationFailed)
	}
	if len(prob.Violations) == 0 || prob.Violations[0].Field != "/status" {
		t.Fatalf("Violations = %#v, want /status", prob.Violations)
	}
}

func TestStrictDecodeDoesNotWriteResponse(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/apis/test", strings.NewReader(sampleJSONBody))
	req.Header.Set("Content-Type", "text/plain")
	rec := httptest.NewRecorder()

	var dst decodeSample
	_ = StrictDecode(rec, req, testLimits, ModeCreateRequest, &dst)
	if rec.Code != http.StatusOK || rec.Body.Len() != 0 {
		t.Fatalf("StrictDecode must not write the response; code=%d body=%q", rec.Code, rec.Body.String())
	}
}

func TestStrictDecodeNilBody(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/apis/test", nil)
	req.Body = nil
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	var dst decodeSample
	prob := StrictDecode(rec, req, testLimits, ModeCreateRequest, &dst)
	if prob == nil {
		t.Fatal("expected MALFORMED_REQUEST for empty body, got nil")
	}
	if prob.Code != apiproblem.CodeMalformedRequest {
		t.Fatalf("Code = %q, want %q", prob.Code, apiproblem.CodeMalformedRequest)
	}
}

func TestStrictDecodeMaxBytesReaderRejectsOversizedRawBytes(t *testing.T) {
	lim := DefaultLimits()
	lim.MaxObjectBytes = 32

	oversized := bytes.Repeat([]byte("a"), 1024)
	req := httptest.NewRequest(http.MethodPost, "/apis/test", bytes.NewReader(oversized))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	var dst decodeSample
	prob := StrictDecode(rec, req, lim, ModeCreateRequest, &dst)
	if prob == nil || prob.Code != apiproblem.CodeRequestTooLarge {
		t.Fatalf("expected REQUEST_TOO_LARGE, got %#v", prob)
	}
	if prob.Status != http.StatusBadRequest {
		t.Fatalf("Status = %d, want 400", prob.Status)
	}
}
