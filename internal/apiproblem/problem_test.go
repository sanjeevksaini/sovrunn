package apiproblem

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestErrorCodesMatchBaselineTable(t *testing.T) {
	t.Parallel()

	want := map[ErrorCode]struct {
		class  FailureClass
		status int
	}{
		CodeMalformedRequest:      {FailureMalformedRequest, http.StatusBadRequest},
		CodeUnknownField:          {FailureMalformedRequest, http.StatusBadRequest},
		CodeDuplicateField:        {FailureMalformedRequest, http.StatusBadRequest},
		CodeRequestTooLarge:       {FailureMalformedRequest, http.StatusBadRequest},
		CodeAuthRequired:          {FailureAuthenticationRequired, http.StatusUnauthorized},
		CodeAuthorizationDenied:   {FailureAuthorizationDenied, http.StatusForbidden},
		CodeResourceNotFound:      {FailureNotFound, http.StatusNotFound},
		CodeConflict:              {FailureConflict, http.StatusConflict},
		CodeAlreadyExists:         {FailureConflict, http.StatusConflict},
		CodeDeleteBlocked:         {FailureConflict, http.StatusConflict},
		CodeStaleResourceVersion:  {FailureStaleResourceVersion, http.StatusPreconditionFailed},
		CodeUnsupportedMediaType:  {FailureUnsupportedMediaType, http.StatusUnsupportedMediaType},
		CodeValidationFailed:      {FailureValidationFailed, http.StatusUnprocessableEntity},
		CodeInternalError:         {FailureInternal, http.StatusInternalServerError},
		CodeDependencyUnavailable: {FailureDependencyUnavailable, http.StatusServiceUnavailable},
	}

	codes := AllErrorCodes()
	if len(codes) != len(want) {
		t.Fatalf("AllErrorCodes len = %d, want %d", len(codes), len(want))
	}
	for _, code := range codes {
		entry, ok := want[code]
		if !ok {
			t.Fatalf("unexpected ErrorCode %q in AllErrorCodes", code)
		}
		if !code.Valid() {
			t.Fatalf("%q must be Valid", code)
		}
		if got := ClassForCode(code); got != entry.class {
			t.Fatalf("ClassForCode(%q) = %q, want %q", code, got, entry.class)
		}
		if got := StatusForCode(code); got != entry.status {
			t.Fatalf("StatusForCode(%q) = %d, want %d", code, got, entry.status)
		}
		if TitleFor(code) == "" || TitleFor(code) == "Error" {
			t.Fatalf("TitleFor(%q) must be a non-default title", code)
		}
	}
}

func TestFailureClassHTTPStatusBaseline(t *testing.T) {
	t.Parallel()

	want := map[FailureClass]int{
		FailureMalformedRequest:       http.StatusBadRequest,
		FailureAuthenticationRequired: http.StatusUnauthorized,
		FailureAuthorizationDenied:    http.StatusForbidden,
		FailureNotFound:               http.StatusNotFound,
		FailureConflict:               http.StatusConflict,
		FailureStaleResourceVersion:   http.StatusPreconditionFailed,
		FailureUnsupportedMediaType:   http.StatusUnsupportedMediaType,
		FailureValidationFailed:       http.StatusUnprocessableEntity,
		FailureInternal:               http.StatusInternalServerError,
		FailureDependencyUnavailable:  http.StatusServiceUnavailable,
	}

	classes := AllFailureClasses()
	if len(classes) != len(want) {
		t.Fatalf("AllFailureClasses len = %d, want %d", len(classes), len(want))
	}
	for _, class := range classes {
		status, ok := want[class]
		if !ok {
			t.Fatalf("unexpected FailureClass %q", class)
		}
		if !class.Valid() {
			t.Fatalf("%q must be Valid", class)
		}
		if got := StatusFor(class); got != status {
			t.Fatalf("StatusFor(%q) = %d, want %d", class, got, status)
		}
	}

	// Oversized body is 400, not 413 (F12-ERROR-002 founder-approval note).
	if StatusForCode(CodeRequestTooLarge) != http.StatusBadRequest {
		t.Fatalf("REQUEST_TOO_LARGE must map to 400, got %d", StatusForCode(CodeRequestTooLarge))
	}
	if StatusForCode(CodeRequestTooLarge) == http.StatusRequestEntityTooLarge {
		t.Fatal("REQUEST_TOO_LARGE must not use 413")
	}
}

func TestViolationCodeRegistry(t *testing.T) {
	t.Parallel()

	if !ViolationOperationTargetScopeMismatch.Valid() {
		t.Fatal("OPERATION_TARGET_SCOPE_MISMATCH must be registered")
	}
	if ViolationOperationTargetScopeMismatch != "OPERATION_TARGET_SCOPE_MISMATCH" {
		t.Fatalf("code value = %q", ViolationOperationTargetScopeMismatch)
	}

	found := false
	for _, c := range AllViolationCodes() {
		if !c.Valid() {
			t.Fatalf("%q from AllViolationCodes must be Valid", c)
		}
		if c == ViolationOperationTargetScopeMismatch {
			found = true
		}
	}
	if !found {
		t.Fatal("AllViolationCodes must include OPERATION_TARGET_SCOPE_MISMATCH")
	}

	if ViolationCode("NOT_A_REAL_CODE").Valid() {
		t.Fatal("unknown violation code must be invalid")
	}
	if ErrorCode("NOT_A_REAL_CODE").Valid() {
		t.Fatal("unknown error code must be invalid")
	}
	if FailureClass("not_a_class").Valid() {
		t.Fatal("unknown failure class must be invalid")
	}
}

func TestProblemSerializesToRFC9457(t *testing.T) {
	t.Parallel()

	p := New(CodeValidationFailed).
		WithDetail("One or more fields are invalid.").
		WithInstance("/apis/core.sovrunn.io/v1alpha1/projects").
		WithRequestID("opaque-request-id").
		WithViolations([]Violation{{
			Field:   "/spec/storage/sizeGiB",
			Code:    ViolationOutOfRange,
			Message: "Value must be greater than zero.",
		}})

	raw, err := json.Marshal(p)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}

	var got map[string]any
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}

	assertString := func(key, want string) {
		t.Helper()
		v, ok := got[key].(string)
		if !ok || v != want {
			t.Fatalf("%s = %#v, want %q", key, got[key], want)
		}
	}
	assertString("type", "urn:sovrunn:problem:validation-failed")
	assertString("title", "Validation failed")
	assertString("detail", "One or more fields are invalid.")
	assertString("instance", "/apis/core.sovrunn.io/v1alpha1/projects")
	assertString("code", "VALIDATION_FAILED")
	assertString("requestId", "opaque-request-id")

	status, ok := got["status"].(float64)
	if !ok || int(status) != http.StatusUnprocessableEntity {
		t.Fatalf("status = %#v, want %d", got["status"], http.StatusUnprocessableEntity)
	}

	violations, ok := got["violations"].([]any)
	if !ok || len(violations) != 1 {
		t.Fatalf("violations = %#v", got["violations"])
	}
	v0, ok := violations[0].(map[string]any)
	if !ok {
		t.Fatalf("violation[0] = %#v", violations[0])
	}
	if v0["field"] != "/spec/storage/sizeGiB" {
		t.Fatalf("field = %#v", v0["field"])
	}
	if v0["code"] != "OUT_OF_RANGE" {
		t.Fatalf("violation code = %#v", v0["code"])
	}
	if v0["message"] != "Value must be greater than zero." {
		t.Fatalf("message = %#v", v0["message"])
	}
}

func TestProblemOmitsEmptyOptionalFields(t *testing.T) {
	t.Parallel()

	raw, err := json.Marshal(New(CodeResourceNotFound))
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	var got map[string]any
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	for _, key := range []string{"detail", "instance", "requestId", "violations"} {
		if _, present := got[key]; present {
			t.Fatalf("optional %q must be omitted when empty, got %#v", key, got)
		}
	}
	if got["type"] != "urn:sovrunn:problem:resource-not-found" {
		t.Fatalf("type = %#v", got["type"])
	}
	if got["code"] != "RESOURCE_NOT_FOUND" {
		t.Fatalf("code = %#v", got["code"])
	}
}

func TestWithViolationsCopiesSlice(t *testing.T) {
	t.Parallel()

	src := []Violation{{Field: "/metadata/scopeRef", Code: ViolationOperationTargetScopeMismatch, Message: "scope mismatch"}}
	p := New(CodeValidationFailed).WithViolations(src)
	src[0].Message = "mutated"
	if p.Violations[0].Message != "scope mismatch" {
		t.Fatal("WithViolations must copy the violations slice")
	}
}

func TestTypeURN(t *testing.T) {
	t.Parallel()

	cases := map[ErrorCode]string{
		CodeValidationFailed:     "urn:sovrunn:problem:validation-failed",
		CodeRequestTooLarge:      "urn:sovrunn:problem:request-too-large",
		CodeStaleResourceVersion: "urn:sovrunn:problem:stale-resource-version",
	}
	for code, want := range cases {
		if got := TypeURN(code); got != want {
			t.Fatalf("TypeURN(%q) = %q, want %q", code, got, want)
		}
	}
}

func TestUnknownCodeMapsToInternal(t *testing.T) {
	t.Parallel()

	if ClassForCode(ErrorCode("WEIRD")) != FailureInternal {
		t.Fatal("unknown ErrorCode must map to FailureInternal")
	}
	if StatusFor(FailureClass("weird")) != http.StatusInternalServerError {
		t.Fatal("unknown FailureClass must map to 500")
	}
}
