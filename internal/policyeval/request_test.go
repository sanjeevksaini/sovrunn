package policyeval

import (
	"bytes"
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
)

func TestRequestValidAtRegisteredBounds(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		raw  string
	}{
		{
			name: "action_min_1",
			raw:  validRequestJSON(withAction(strings.Repeat("a", 1))),
		},
		{
			name: "action_max_63",
			raw:  validRequestJSON(withAction(strings.Repeat("a", 63))),
		},
		{
			name: "profileRefs_min_1",
			raw:  validRequestJSON(withProfileCount(1)),
		},
		{
			name: "profileRefs_max_32",
			raw:  validRequestJSON(withProfileCount(32)),
		},
		{
			name: "candidateRefs_min_0_omitted",
			raw:  validRequestJSON(omitCandidates()),
		},
		{
			name: "candidateRefs_min_0_empty_array",
			raw:  validRequestJSON(withCandidateCount(0)),
		},
		{
			name: "candidateRefs_max_64",
			raw:  validRequestJSON(withCandidateCount(64)),
		},
		{
			name: "requestId_min_1",
			raw:  validRequestJSON(withRequestID(strings.Repeat("r", 1))),
		},
		{
			name: "requestId_max_128",
			raw:  validRequestJSON(withRequestID(strings.Repeat("r", 128))),
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			req, err := DecodePolicyEvaluationRequestJSON([]byte(tc.raw))
			if err != nil {
				t.Fatalf("DecodePolicyEvaluationRequestJSON: %v\nraw=%s", err, tc.raw)
			}
			if err := req.Validate(); err != nil {
				t.Fatalf("Validate: %v", err)
			}
		})
	}
}

func TestRequestRejectInvalidUTF8(t *testing.T) {
	t.Parallel()
	// Inject invalid UTF-8 into the action string value.
	invalid := append([]byte(`{"subjectRef":{"apiVersion":"v1","kind":"ServiceInstance","name":"svc"},"action":"`),
		[]byte{0xff}...)
	invalid = append(invalid, []byte(`","profileRefs":[{"apiVersion":"v1","kind":"RoleDefinition","name":"r","uid":"u1"}],"requestId":"id"}`)...)
	if utf8.Valid(invalid) {
		t.Fatal("test setup: expected invalid UTF-8")
	}
	_, err := DecodePolicyEvaluationRequestJSON(invalid)
	if !errors.Is(err, errRequestInvalid) {
		t.Fatalf("err=%v, want errRequestInvalid", err)
	}
}

func TestRequestRejectLoneHighSurrogateEscape(t *testing.T) {
	t.Parallel()
	raw := validRequestJSON(withAction(`\uD800`))
	_, err := DecodePolicyEvaluationRequestJSON([]byte(raw))
	if !errors.Is(err, errRequestInvalid) {
		t.Fatalf("err=%v, want errRequestInvalid", err)
	}
}

func TestRequestRejectLoneLowSurrogateEscape(t *testing.T) {
	t.Parallel()
	raw := validRequestJSON(withAction(`\uDC00`))
	_, err := DecodePolicyEvaluationRequestJSON([]byte(raw))
	if !errors.Is(err, errRequestInvalid) {
		t.Fatalf("err=%v, want errRequestInvalid", err)
	}
}

func TestRequestRejectMalformedSurrogatePair(t *testing.T) {
	t.Parallel()
	// High surrogate followed by non-surrogate escape.
	raw := validRequestJSON(withAction(`\uD800\u0041`))
	_, err := DecodePolicyEvaluationRequestJSON([]byte(raw))
	if !errors.Is(err, errRequestInvalid) {
		t.Fatalf("err=%v, want errRequestInvalid", err)
	}
}

func TestRequestRejectUnpairedSurrogateInStringValue(t *testing.T) {
	t.Parallel()
	raw := validRequestJSON(withRequestID(`prefix-\uD83D-suffix`))
	_, err := DecodePolicyEvaluationRequestJSON([]byte(raw))
	if !errors.Is(err, errRequestInvalid) {
		t.Fatalf("err=%v, want errRequestInvalid", err)
	}
}

func TestRequestAcceptValidSurrogatePair(t *testing.T) {
	t.Parallel()
	raw := validRequestJSON(withAction(`\uD83D\uDE00`))
	req, err := DecodePolicyEvaluationRequestJSON([]byte(raw))
	if err != nil {
		t.Fatalf("DecodePolicyEvaluationRequestJSON: %v", err)
	}
	if req.Action != "😀" {
		t.Fatalf("Action=%q, want grinning face emoji", req.Action)
	}
}

func TestRequestRejectUnknownRootField(t *testing.T) {
	t.Parallel()
	raw := validRequestJSON(withExtraRoot(`"extraField":true`))
	_, err := DecodePolicyEvaluationRequestJSON([]byte(raw))
	if !errors.Is(err, errRequestInvalid) {
		t.Fatalf("err=%v, want errRequestInvalid", err)
	}
}

func TestRequestRejectUnknownNestedReferenceField(t *testing.T) {
	t.Parallel()
	raw := `{
		"subjectRef":{"apiVersion":"v1","kind":"ServiceInstance","name":"svc","extra":"x"},
		"action":"read",
		"profileRefs":[{"apiVersion":"v1","kind":"RoleDefinition","name":"r","uid":"u1"}],
		"requestId":"id-1"
	}`
	_, err := DecodePolicyEvaluationRequestJSON([]byte(raw))
	if !errors.Is(err, errRequestInvalid) {
		t.Fatalf("err=%v, want errRequestInvalid", err)
	}
}

func TestRequestRejectDuplicateRootMember(t *testing.T) {
	t.Parallel()
	raw := `{
		"subjectRef":{"apiVersion":"v1","kind":"ServiceInstance","name":"svc"},
		"action":"read",
		"action":"write",
		"profileRefs":[{"apiVersion":"v1","kind":"RoleDefinition","name":"r","uid":"u1"}],
		"requestId":"id-1"
	}`
	_, err := DecodePolicyEvaluationRequestJSON([]byte(raw))
	if !errors.Is(err, errRequestInvalid) {
		t.Fatalf("err=%v, want errRequestInvalid", err)
	}
}

func TestRequestRejectDuplicateNestedReferenceMember(t *testing.T) {
	t.Parallel()
	raw := `{
		"subjectRef":{"apiVersion":"v1","kind":"ServiceInstance","name":"svc","name":"other"},
		"action":"read",
		"profileRefs":[{"apiVersion":"v1","kind":"RoleDefinition","name":"r","uid":"u1"}],
		"requestId":"id-1"
	}`
	_, err := DecodePolicyEvaluationRequestJSON([]byte(raw))
	if !errors.Is(err, errRequestInvalid) {
		t.Fatalf("err=%v, want errRequestInvalid", err)
	}
}

func TestRequestRejectTrailingJSONContent(t *testing.T) {
	t.Parallel()
	raw := validRequestJSON() + `{"trailing":true}`
	_, err := DecodePolicyEvaluationRequestJSON([]byte(raw))
	if !errors.Is(err, errRequestInvalid) {
		t.Fatalf("err=%v, want errRequestInvalid", err)
	}
}

func TestRequestRejectMalformedReference(t *testing.T) {
	t.Parallel()
	raw := `{
		"subjectRef":{"apiVersion":"v1","kind":"ServiceInstance"},
		"action":"read",
		"profileRefs":[{"apiVersion":"v1","kind":"RoleDefinition","name":"r","uid":"u1"}],
		"requestId":"id-1"
	}`
	_, err := DecodePolicyEvaluationRequestJSON([]byte(raw))
	if !errors.Is(err, errRequestInvalid) {
		t.Fatalf("err=%v, want errRequestInvalid", err)
	}
}

func TestRequestContextRefAbsentAndNullRules(t *testing.T) {
	t.Parallel()

	t.Run("absent_accepted", func(t *testing.T) {
		t.Parallel()
		req, err := DecodePolicyEvaluationRequestJSON([]byte(validRequestJSON(omitContext())))
		if err != nil {
			t.Fatalf("Decode: %v", err)
		}
		if req.ContextRef != nil {
			t.Fatalf("ContextRef=%v, want nil", req.ContextRef)
		}
	})

	t.Run("valid_present_accepted", func(t *testing.T) {
		t.Parallel()
		req, err := DecodePolicyEvaluationRequestJSON([]byte(validRequestJSON(withContextRef(
			`{"apiVersion":"governance.sovrunn.io/v1alpha1","kind":"EffectiveGovernanceContext","name":"egc","uid":"egc-1"}`,
		))))
		if err != nil {
			t.Fatalf("Decode: %v", err)
		}
		if req.ContextRef == nil || req.ContextRef.Name != "egc" {
			t.Fatalf("ContextRef=%v", req.ContextRef)
		}
	})

	t.Run("explicit_null_rejected", func(t *testing.T) {
		t.Parallel()
		_, err := DecodePolicyEvaluationRequestJSON([]byte(validRequestJSON(withContextRef(`null`))))
		if !errors.Is(err, errRequestInvalid) {
			t.Fatalf("err=%v, want errRequestInvalid", err)
		}
	})

	t.Run("empty_object_rejected", func(t *testing.T) {
		t.Parallel()
		_, err := DecodePolicyEvaluationRequestJSON([]byte(validRequestJSON(withContextRef(`{}`))))
		if !errors.Is(err, errRequestInvalid) {
			t.Fatalf("err=%v, want errRequestInvalid", err)
		}
	})

	t.Run("partial_rejected", func(t *testing.T) {
		t.Parallel()
		_, err := DecodePolicyEvaluationRequestJSON([]byte(validRequestJSON(withContextRef(
			`{"apiVersion":"governance.sovrunn.io/v1alpha1","kind":"EffectiveGovernanceContext"}`,
		))))
		if !errors.Is(err, errRequestInvalid) {
			t.Fatalf("err=%v, want errRequestInvalid", err)
		}
	})
}

func TestRequestCandidateRefsNullEmptyOmitted(t *testing.T) {
	t.Parallel()

	t.Run("omitted_constructs_empty_slice", func(t *testing.T) {
		t.Parallel()
		req, err := DecodePolicyEvaluationRequestJSON([]byte(validRequestJSON(omitCandidates())))
		if err != nil {
			t.Fatalf("Decode: %v", err)
		}
		if req.CandidateRefs == nil {
			t.Fatal("CandidateRefs is nil; want non-nil empty slice")
		}
		if len(req.CandidateRefs) != 0 {
			t.Fatalf("len(CandidateRefs)=%d, want 0", len(req.CandidateRefs))
		}
	})

	t.Run("empty_array_accepted", func(t *testing.T) {
		t.Parallel()
		req, err := DecodePolicyEvaluationRequestJSON([]byte(validRequestJSON(withCandidateCount(0))))
		if err != nil {
			t.Fatalf("Decode: %v", err)
		}
		if len(req.CandidateRefs) != 0 {
			t.Fatalf("len(CandidateRefs)=%d, want 0", len(req.CandidateRefs))
		}
	})

	t.Run("null_rejected", func(t *testing.T) {
		t.Parallel()
		raw := validRequestJSON(withRawCandidates(`null`))
		_, err := DecodePolicyEvaluationRequestJSON([]byte(raw))
		if !errors.Is(err, errRequestInvalid) {
			t.Fatalf("err=%v, want errRequestInvalid", err)
		}
	})
}

func TestRequestRejectEmptyRequiredAndEmptyProfiles(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		raw  string
	}{
		{"empty_action", validRequestJSON(withAction(""))},
		{"empty_requestId", validRequestJSON(withRequestID(""))},
		{"null_action", validRequestJSON(withRawAction(`null`))},
		{"null_requestId", validRequestJSON(withRawRequestID(`null`))},
		{"null_subjectRef", validRequestJSON(withRawSubject(`null`))},
		{"null_profileRefs", validRequestJSON(withRawProfiles(`null`))},
		{"empty_profileRefs", validRequestJSON(withProfileCount(0))},
		{"omitted_profileRefs", validRequestJSON(omitProfiles())},
		{"action_too_long", validRequestJSON(withAction(strings.Repeat("a", 64)))},
		{"requestId_too_long", validRequestJSON(withRequestID(strings.Repeat("r", 129)))},
		{"profileRefs_too_many", validRequestJSON(withProfileCount(33))},
		{"candidateRefs_too_many", validRequestJSON(withCandidateCount(65))},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			_, err := DecodePolicyEvaluationRequestJSON([]byte(tc.raw))
			if !errors.Is(err, errRequestInvalid) {
				t.Fatalf("err=%v, want errRequestInvalid\nraw=%s", err, tc.raw)
			}
		})
	}
}

func TestRequestFEATURE0012StructuralValidationOnEveryReference(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		raw  string
	}{
		{
			name: "subject_missing_kind",
			raw: `{
				"subjectRef":{"apiVersion":"v1","name":"svc"},
				"action":"read",
				"profileRefs":[{"apiVersion":"v1","kind":"RoleDefinition","name":"r","uid":"u1"}],
				"requestId":"id-1"
			}`,
		},
		{
			name: "profile_missing_apiVersion",
			raw: `{
				"subjectRef":{"apiVersion":"v1","kind":"ServiceInstance","name":"svc"},
				"action":"read",
				"profileRefs":[{"kind":"RoleDefinition","name":"r","uid":"u1"}],
				"requestId":"id-1"
			}`,
		},
		{
			name: "candidate_missing_name",
			raw: `{
				"subjectRef":{"apiVersion":"v1","kind":"ServiceInstance","name":"svc"},
				"action":"read",
				"profileRefs":[{"apiVersion":"v1","kind":"RoleDefinition","name":"r","uid":"u1"}],
				"candidateRefs":[{"apiVersion":"v1","kind":"ServicePlan"}],
				"requestId":"id-1"
			}`,
		},
		{
			name: "provider_native_subject_rejected",
			raw: `{
				"subjectRef":{"apiVersion":"v1","kind":"AWS::RDS::DBInstance","name":"svc"},
				"action":"read",
				"profileRefs":[{"apiVersion":"v1","kind":"RoleDefinition","name":"r","uid":"u1"}],
				"requestId":"id-1"
			}`,
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			_, err := DecodePolicyEvaluationRequestJSON([]byte(tc.raw))
			if !errors.Is(err, errRequestInvalid) {
				t.Fatalf("err=%v, want errRequestInvalid", err)
			}
		})
	}
}

func TestRequestDirectGoConstructionSemantics(t *testing.T) {
	t.Parallel()

	base := PolicyEvaluationRequest{
		SubjectRef: apimeta.TypedRef{
			APIVersion: "services.sovrunn.io/v1alpha1",
			Kind:       "ServiceInstance",
			Name:       "reporting-api",
		},
		Action: "service.read",
		ProfileRefs: []apimeta.TypedRef{{
			APIVersion: "iam.sovrunn.io/v1alpha1",
			Kind:       "RoleDefinition",
			Name:       "reader",
			UID:        "role-001",
		}},
		RequestID: "corr-1",
	}

	t.Run("nil_context_is_absence", func(t *testing.T) {
		t.Parallel()
		req := base
		req.ContextRef = nil
		req.CandidateRefs = nil
		if err := req.Validate(); err != nil {
			t.Fatalf("Validate: %v", err)
		}
	})

	t.Run("nil_candidates_are_empty_set", func(t *testing.T) {
		t.Parallel()
		req := base
		req.CandidateRefs = nil
		if err := req.Validate(); err != nil {
			t.Fatalf("Validate: %v", err)
		}
	})

	t.Run("empty_candidates_are_empty_set", func(t *testing.T) {
		t.Parallel()
		req := base
		req.CandidateRefs = []apimeta.TypedRef{}
		if err := req.Validate(); err != nil {
			t.Fatalf("Validate: %v", err)
		}
	})
}

func TestRequestUnmarshalJSONDelegatesToStrictDecoder(t *testing.T) {
	t.Parallel()

	t.Run("happy_path", func(t *testing.T) {
		t.Parallel()
		var req PolicyEvaluationRequest
		if err := json.Unmarshal([]byte(validRequestJSON()), &req); err != nil {
			t.Fatalf("json.Unmarshal: %v", err)
		}
		if req.Action != "service.read" {
			t.Fatalf("Action=%q", req.Action)
		}
	})

	t.Run("unknown_field_rejected", func(t *testing.T) {
		t.Parallel()
		var req PolicyEvaluationRequest
		err := json.Unmarshal([]byte(validRequestJSON(withExtraRoot(`"extra":1`))), &req)
		if !errors.Is(err, errRequestInvalid) {
			t.Fatalf("err=%v, want errRequestInvalid", err)
		}
	})

	t.Run("error_text_does_not_echo_request_values", func(t *testing.T) {
		t.Parallel()
		secret := "super-secret-action-value"
		_, err := DecodePolicyEvaluationRequestJSON([]byte(validRequestJSON(withAction(secret), withExtraRoot(`"nope":true`))))
		if err == nil {
			t.Fatal("expected error")
		}
		if strings.Contains(err.Error(), secret) {
			t.Fatalf("error echoed request value: %v", err)
		}
		if err.Error() != errRequestInvalid.Error() {
			t.Fatalf("error=%q, want %q", err.Error(), errRequestInvalid.Error())
		}
	})
}

// --- test helpers ---

type requestOpt func(map[string]json.RawMessage)

func validRequestJSON(opts ...requestOpt) string {
	root := map[string]json.RawMessage{
		"subjectRef":  json.RawMessage(`{"apiVersion":"services.sovrunn.io/v1alpha1","kind":"ServiceInstance","name":"reporting-api"}`),
		"action":      json.RawMessage(`"service.read"`),
		"profileRefs": json.RawMessage(`[{"apiVersion":"iam.sovrunn.io/v1alpha1","kind":"RoleDefinition","name":"reader","uid":"role-001"}]`),
		"requestId":   json.RawMessage(`"corr-001"`),
	}
	for _, opt := range opts {
		opt(root)
	}
	var b bytes.Buffer
	b.WriteByte('{')
	first := true
	// Stable key order for readability; encoding/json map marshal is fine too.
	keys := []string{"subjectRef", "action", "contextRef", "profileRefs", "candidateRefs", "requestId"}
	extras := make([]string, 0)
	for k := range root {
		switch k {
		case "subjectRef", "action", "contextRef", "profileRefs", "candidateRefs", "requestId":
		default:
			extras = append(extras, k)
		}
	}
	keys = append(keys, extras...)
	for _, k := range keys {
		v, ok := root[k]
		if !ok {
			continue
		}
		if !first {
			b.WriteByte(',')
		}
		first = false
		b.WriteString(`"` + k + `":`)
		b.Write(v)
	}
	b.WriteByte('}')
	return b.String()
}

func withAction(action string) requestOpt {
	return func(m map[string]json.RawMessage) {
		// Allow callers to pass already-escaped JSON string bodies (e.g. \uD800).
		if strings.Contains(action, `\u`) || action == "" {
			m["action"] = json.RawMessage(`"` + action + `"`)
			return
		}
		enc, err := json.Marshal(action)
		if err != nil {
			panic(err)
		}
		m["action"] = json.RawMessage(enc)
	}
}

func withRawAction(raw string) requestOpt {
	return func(m map[string]json.RawMessage) {
		m["action"] = json.RawMessage(raw)
	}
}

func withRequestID(id string) requestOpt {
	return func(m map[string]json.RawMessage) {
		if strings.Contains(id, `\u`) || id == "" {
			m["requestId"] = json.RawMessage(`"` + id + `"`)
			return
		}
		enc, err := json.Marshal(id)
		if err != nil {
			panic(err)
		}
		m["requestId"] = json.RawMessage(enc)
	}
}

func withRawRequestID(raw string) requestOpt {
	return func(m map[string]json.RawMessage) {
		m["requestId"] = json.RawMessage(raw)
	}
}

func withRawSubject(raw string) requestOpt {
	return func(m map[string]json.RawMessage) {
		m["subjectRef"] = json.RawMessage(raw)
	}
}

func withContextRef(raw string) requestOpt {
	return func(m map[string]json.RawMessage) {
		m["contextRef"] = json.RawMessage(raw)
	}
}

func omitContext() requestOpt {
	return func(m map[string]json.RawMessage) {
		delete(m, "contextRef")
	}
}

func withProfileCount(n int) requestOpt {
	return func(m map[string]json.RawMessage) {
		var b strings.Builder
		b.WriteByte('[')
		for i := 0; i < n; i++ {
			if i > 0 {
				b.WriteByte(',')
			}
			b.WriteString(`{"apiVersion":"iam.sovrunn.io/v1alpha1","kind":"RoleDefinition","name":"role-`)
			b.WriteString(itoa(i))
			b.WriteString(`","uid":"uid-`)
			b.WriteString(itoa(i))
			b.WriteString(`"}`)
		}
		b.WriteByte(']')
		m["profileRefs"] = json.RawMessage(b.String())
	}
}

func withRawProfiles(raw string) requestOpt {
	return func(m map[string]json.RawMessage) {
		m["profileRefs"] = json.RawMessage(raw)
	}
}

func omitProfiles() requestOpt {
	return func(m map[string]json.RawMessage) {
		delete(m, "profileRefs")
	}
}

func withCandidateCount(n int) requestOpt {
	return func(m map[string]json.RawMessage) {
		var b strings.Builder
		b.WriteByte('[')
		for i := 0; i < n; i++ {
			if i > 0 {
				b.WriteByte(',')
			}
			b.WriteString(`{"apiVersion":"services.sovrunn.io/v1alpha1","kind":"ServicePlan","name":"plan-`)
			b.WriteString(itoa(i))
			b.WriteString(`"}`)
		}
		b.WriteByte(']')
		m["candidateRefs"] = json.RawMessage(b.String())
	}
}

func withRawCandidates(raw string) requestOpt {
	return func(m map[string]json.RawMessage) {
		m["candidateRefs"] = json.RawMessage(raw)
	}
}

func omitCandidates() requestOpt {
	return func(m map[string]json.RawMessage) {
		delete(m, "candidateRefs")
	}
}

func withExtraRoot(fragment string) requestOpt {
	return func(m map[string]json.RawMessage) {
		// fragment like `"extraField":true`
		parts := strings.SplitN(fragment, ":", 2)
		if len(parts) != 2 {
			panic("bad extra fragment")
		}
		key := strings.Trim(parts[0], ` "`)
		m[key] = json.RawMessage(strings.TrimSpace(parts[1]))
	}
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var digits [16]byte
	i := len(digits)
	for n > 0 {
		i--
		digits[i] = byte('0' + n%10)
		n /= 10
	}
	return string(digits[i:])
}
