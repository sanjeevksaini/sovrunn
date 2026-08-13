package validate_test

import (
	"strings"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/cloudmodel/validate"
)

func TestValidateIdempotencyKey(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name  string
		value string
		code  apiproblem.ErrorCode
	}{
		{name: "missing", value: "", code: apiproblem.CodeMalformedRequest},
		{name: "empty", value: "   ", code: apiproblem.CodeMalformedRequest},
		{name: "malformed space", value: "bad key", code: apiproblem.CodeMalformedRequest},
		{name: "over-length", value: strings.Repeat("a", validate.MaxIdempotencyKeyLen+1), code: apiproblem.CodeMalformedRequest},
		{name: "ok", value: "key-1", code: ""},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			prob := validate.ValidateIdempotencyKey(tc.value)
			if tc.code == "" {
				if prob != nil {
					t.Fatalf("want nil, got %#v", prob)
				}
				return
			}
			if prob == nil || prob.Code != tc.code {
				t.Fatalf("got %#v, want %s", prob, tc.code)
			}
		})
	}
}

func TestValidateParticipationCreateHeaders_IfMatchRejected(t *testing.T) {
	t.Parallel()
	prob := validate.ValidateParticipationCreateHeaders("key-1", "1")
	if prob == nil || prob.Code != apiproblem.CodeMalformedRequest {
		t.Fatalf("If-Match on create: %#v", prob)
	}
	if prob := validate.ValidateParticipationCreateHeaders("key-1", ""); prob != nil {
		t.Fatalf("valid create headers: %#v", prob)
	}
}

func TestValidateEmptyActionBody(t *testing.T) {
	t.Parallel()
	if prob := validate.ValidateEmptyActionBody(nil); prob != nil {
		t.Fatalf("nil body: %#v", prob)
	}
	if prob := validate.ValidateEmptyActionBody([]byte{}); prob != nil {
		t.Fatalf("empty: %#v", prob)
	}
	if prob := validate.ValidateEmptyActionBody([]byte(`{}`)); prob == nil || prob.Code != apiproblem.CodeMalformedRequest {
		t.Fatalf("{}: %#v", prob)
	}
	if prob := validate.ValidateEmptyActionBody([]byte(`{"x":1}`)); prob == nil || prob.Code != apiproblem.CodeMalformedRequest {
		t.Fatalf("non-empty: %#v", prob)
	}
	prob := validate.ValidateEmptyActionBody([]byte(`{"bootstrapGrant":{"id":"g1"}}`))
	if prob == nil || prob.Code != apiproblem.CodeAuthorizationDenied {
		t.Fatalf("bootstrapGrant: %#v", prob)
	}
}
