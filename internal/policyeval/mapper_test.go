package policyeval

import (
	"bytes"
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/decision"
)

const mapperFrozenDigest = "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"

func mapperCanonicalAt() string {
	return time.Date(2026, 8, 25, 12, 0, 0, 0, time.UTC).Format(time.RFC3339Nano)
}

func validMapperResult(outcome Outcome) PolicyEvaluationResult {
	return PolicyEvaluationResult{
		Outcome:     outcome,
		ReasonCodes: []string{"POLICY_ALLOW"},
		InputDigest: mapperFrozenDigest,
		EvaluatedAt: mapperCanonicalAt(),
	}
}

func validMapperTiming() decision.TimingEnvelope {
	at := mapperCanonicalAt()
	return decision.TimingEnvelope{
		StartedAt:   at,
		CompletedAt: at,
	}
}

func validMapperEvaluator() decision.EvaluatorIdentity {
	return decision.EvaluatorIdentity{
		Type:    "policy-evaluator",
		Version: "1.0.0",
	}
}

func validMapperSnapshotIdent() InputIdentity {
	return InputIdentity{
		InputSnapshotRef: &apimeta.TypedRef{
			APIVersion: "governance.sovrunn.io/v1alpha1",
			Kind:       "InputSnapshot",
			Name:       "placement-input-1",
			UID:        "input-uid-1",
		},
	}
}

func validMapperIntegrityIdent() InputIdentity {
	return InputIdentity{
		InputIntegrity: &decision.OpaqueIntegrityCarrier{State: "bound"},
	}
}

func assertMappingError(t *testing.T, err error) {
	t.Helper()
	if !errors.Is(err, errMappingInputInvalid) {
		t.Fatalf("err=%v, want %v", err, errMappingInputInvalid)
	}
	if err != nil && strings.Contains(err.Error(), mapperFrozenDigest) {
		t.Fatalf("mapping error must not echo digest: %v", err)
	}
}

func TestMapper_AllOutcomesMapStatus(t *testing.T) {
	t.Parallel()

	cases := []struct {
		outcome Outcome
		want    decision.EvaluationResultStatus
	}{
		{OutcomeAllow, decision.EvaluationResultStatusSuccess},
		{OutcomeDeny, decision.EvaluationResultStatusSuccess},
		{OutcomeRequiresApproval, decision.EvaluationResultStatusSuccess},
		{OutcomeIndeterminate, decision.EvaluationResultStatusIndeterminate},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(string(tc.outcome), func(t *testing.T) {
			t.Parallel()
			got, err := MapToEvaluationResult(
				validMapperResult(tc.outcome),
				validMapperTiming(),
				validMapperEvaluator(),
				validMapperSnapshotIdent(),
			)
			if err != nil {
				t.Fatalf("MapToEvaluationResult: %v", err)
			}
			if got.ResultStatus != tc.want {
				t.Fatalf("resultStatus=%q, want %q", got.ResultStatus, tc.want)
			}
			if got.EvaluatedAt != validMapperTiming().CompletedAt {
				t.Fatalf("evaluatedAt=%q != completedAt=%q", got.EvaluatedAt, validMapperTiming().CompletedAt)
			}
		})
	}
}

func TestMapper_ExactPopulatedAndOmittedFields(t *testing.T) {
	t.Parallel()

	result := validMapperResult(OutcomeAllow)
	timing := validMapperTiming()
	evaluator := validMapperEvaluator()
	ident := validMapperSnapshotIdent()

	got, err := MapToEvaluationResult(result, timing, evaluator, ident)
	if err != nil {
		t.Fatalf("MapToEvaluationResult: %v", err)
	}

	if got.Evaluator != evaluator {
		t.Fatalf("Evaluator=%+v, want %+v", got.Evaluator, evaluator)
	}
	if got.InputSnapshotRef == nil || *got.InputSnapshotRef != *ident.InputSnapshotRef {
		t.Fatalf("InputSnapshotRef=%v, want %v", got.InputSnapshotRef, ident.InputSnapshotRef)
	}
	if got.InputIntegrity != nil {
		t.Fatalf("InputIntegrity must be absent, got %+v", got.InputIntegrity)
	}
	if got.ResultStatus != decision.EvaluationResultStatusSuccess {
		t.Fatalf("ResultStatus=%q, want SUCCESS", got.ResultStatus)
	}
	if len(got.Result) == 0 {
		t.Fatal("Result payload must be present")
	}
	var embedded PolicyEvaluationResult
	if err := json.Unmarshal(got.Result, &embedded); err != nil {
		t.Fatalf("unmarshal Result: %v", err)
	}
	if embedded.Outcome != OutcomeAllow || embedded.InputDigest != mapperFrozenDigest {
		t.Fatalf("embedded result=%+v", embedded)
	}
	if got.Timing.StartedAt != timing.StartedAt || got.Timing.CompletedAt != timing.CompletedAt {
		t.Fatalf("Timing=%+v, want %+v", got.Timing, timing)
	}
	if got.EvaluatedAt != timing.CompletedAt || got.EvaluatedAt != result.EvaluatedAt {
		t.Fatalf("EvaluatedAt=%q mismatch", got.EvaluatedAt)
	}

	// Exact omitted fields.
	if got.Timing.DurationMs != 0 {
		t.Fatalf("DurationMs=%d, want 0 (omitted)", got.Timing.DurationMs)
	}
	if got.ExecutionConfig != nil {
		t.Fatalf("ExecutionConfig must be omitted, got %s", got.ExecutionConfig)
	}
	if got.TrustBoundary != "" {
		t.Fatalf("TrustBoundary must be omitted, got %q", got.TrustBoundary)
	}
	if got.SafetyPolicyFilters != nil {
		t.Fatalf("SafetyPolicyFilters must be omitted, got %v", got.SafetyPolicyFilters)
	}
	if got.StructuralTrustState != "" {
		t.Fatalf("StructuralTrustState must be omitted, got %q", got.StructuralTrustState)
	}
	if got.ScopeEvidence != nil {
		t.Fatalf("ScopeEvidence must be omitted, got %+v", got.ScopeEvidence)
	}
}

func TestMapper_BothInputIdentityCarriers(t *testing.T) {
	t.Parallel()

	ident := InputIdentity{
		InputSnapshotRef: validMapperSnapshotIdent().InputSnapshotRef,
		InputIntegrity:   validMapperIntegrityIdent().InputIntegrity,
	}
	got, err := MapToEvaluationResult(
		validMapperResult(OutcomeDeny),
		validMapperTiming(),
		validMapperEvaluator(),
		ident,
	)
	if err != nil {
		t.Fatalf("MapToEvaluationResult: %v", err)
	}
	if got.InputSnapshotRef == nil || got.InputIntegrity == nil {
		t.Fatal("both input identity carriers must be populated")
	}
}

func TestMapper_IntegrityOnly(t *testing.T) {
	t.Parallel()

	got, err := MapToEvaluationResult(
		validMapperResult(OutcomeRequiresApproval),
		validMapperTiming(),
		validMapperEvaluator(),
		validMapperIntegrityIdent(),
	)
	if err != nil {
		t.Fatalf("MapToEvaluationResult: %v", err)
	}
	if got.InputSnapshotRef != nil || got.InputIntegrity == nil {
		t.Fatalf("want integrity-only, snapshot=%v integrity=%v", got.InputSnapshotRef, got.InputIntegrity)
	}
}

func TestMapper_MalformedEvaluator(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name string
		ev   decision.EvaluatorIdentity
	}{
		{"empty_type", decision.EvaluatorIdentity{Type: "", Version: "1.0.0"}},
		{"blank_type", decision.EvaluatorIdentity{Type: "  ", Version: "1.0.0"}},
		{"empty_version", decision.EvaluatorIdentity{Type: "policy-evaluator", Version: ""}},
		{"blank_version", decision.EvaluatorIdentity{Type: "policy-evaluator", Version: "\t"}},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			_, err := MapToEvaluationResult(
				validMapperResult(OutcomeAllow),
				validMapperTiming(),
				tc.ev,
				validMapperSnapshotIdent(),
			)
			assertMappingError(t, err)
		})
	}
}

func TestMapper_InvalidTiming(t *testing.T) {
	t.Parallel()

	at := mapperCanonicalAt()
	cases := []struct {
		name   string
		timing decision.TimingEnvelope
		result PolicyEvaluationResult
	}{
		{
			name:   "zero_started",
			timing: decision.TimingEnvelope{StartedAt: "", CompletedAt: at},
			result: validMapperResult(OutcomeAllow),
		},
		{
			name:   "zero_completed",
			timing: decision.TimingEnvelope{StartedAt: at, CompletedAt: ""},
			result: validMapperResult(OutcomeAllow),
		},
		{
			name: "evaluated_completed_mismatch",
			timing: decision.TimingEnvelope{
				StartedAt:   at,
				CompletedAt: time.Date(2026, 8, 25, 13, 0, 0, 0, time.UTC).Format(time.RFC3339Nano),
			},
			result: validMapperResult(OutcomeAllow),
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			_, err := MapToEvaluationResult(tc.result, tc.timing, validMapperEvaluator(), validMapperSnapshotIdent())
			assertMappingError(t, err)
		})
	}
}

func TestMapper_CanonicalTimestampValidation(t *testing.T) {
	t.Parallel()

	good := mapperCanonicalAt()
	reject := []string{
		"",
		"not-a-timestamp",
		"2026-08-25 12:00:00Z",
		"2026-08-25T12:00:00+00:00",
		"2026-08-25T12:00:00.000Z",
		"2026-08-25T17:30:00+05:30",
		"2026-08-25T12:00:00.050Z", // non-canonical trailing zero
	}
	for _, bad := range reject {
		bad := bad
		t.Run("reject_"+bad, func(t *testing.T) {
			t.Parallel()
			timing := decision.TimingEnvelope{StartedAt: bad, CompletedAt: good}
			result := validMapperResult(OutcomeAllow)
			_, err := MapToEvaluationResult(result, timing, validMapperEvaluator(), validMapperSnapshotIdent())
			assertMappingError(t, err)

			timing = decision.TimingEnvelope{StartedAt: good, CompletedAt: bad}
			result.EvaluatedAt = bad
			_, err = MapToEvaluationResult(result, timing, validMapperEvaluator(), validMapperSnapshotIdent())
			assertMappingError(t, err)
		})
	}

	// Accept canonical Z value produced by the boundary formatter.
	frac := time.Date(2026, 8, 25, 12, 0, 0, 50_000_000, time.UTC).Format(time.RFC3339Nano)
	timing := decision.TimingEnvelope{StartedAt: frac, CompletedAt: frac}
	result := validMapperResult(OutcomeAllow)
	result.EvaluatedAt = frac
	if _, err := MapToEvaluationResult(result, timing, validMapperEvaluator(), validMapperSnapshotIdent()); err != nil {
		t.Fatalf("canonical fractional Z must pass: %v", err)
	}
}

func TestMapper_InvalidSnapshotRef(t *testing.T) {
	t.Parallel()

	ident := InputIdentity{InputSnapshotRef: &apimeta.TypedRef{
		APIVersion: "",
		Kind:       "InputSnapshot",
		Name:       "x",
	}}
	_, err := MapToEvaluationResult(
		validMapperResult(OutcomeAllow),
		validMapperTiming(),
		validMapperEvaluator(),
		ident,
	)
	assertMappingError(t, err)
}

func TestMapper_InvalidIntegrityCarrier(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name  string
		ident InputIdentity
	}{
		{"nil_state_empty", InputIdentity{InputIntegrity: &decision.OpaqueIntegrityCarrier{}}},
		{"empty_state", InputIdentity{InputIntegrity: &decision.OpaqueIntegrityCarrier{State: ""}}},
		{"blank_state", InputIdentity{InputIntegrity: &decision.OpaqueIntegrityCarrier{State: "  "}}},
		{"partial_carrier", InputIdentity{InputIntegrity: &decision.OpaqueIntegrityCarrier{
			State:                  "",
			CanonicalizationMethod: "jcs",
			DigestAlgorithm:        "sha256",
		}}},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			_, err := MapToEvaluationResult(
				validMapperResult(OutcomeAllow),
				validMapperTiming(),
				validMapperEvaluator(),
				tc.ident,
			)
			assertMappingError(t, err)
		})
	}
}

func TestMapper_BothInputIdentityNil(t *testing.T) {
	t.Parallel()

	_, err := MapToEvaluationResult(
		validMapperResult(OutcomeAllow),
		validMapperTiming(),
		validMapperEvaluator(),
		InputIdentity{},
	)
	assertMappingError(t, err)
}

func TestMapper_InvalidOutcome(t *testing.T) {
	t.Parallel()

	result := validMapperResult(OutcomeAllow)
	result.Outcome = Outcome("Unknown")
	_, err := MapToEvaluationResult(result, validMapperTiming(), validMapperEvaluator(), validMapperSnapshotIdent())
	assertMappingError(t, err)
}

func TestMapper_InvalidReasonCodes(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name  string
		codes []string
	}{
		{"empty", nil},
		{"empty_slice", []string{}},
		{"grammar", []string{"bad-code"}},
		{"duplicate", []string{"POLICY_ALLOW", "POLICY_ALLOW"}},
		{"too_many", reasonCodesN(33)},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			result := validMapperResult(OutcomeAllow)
			result.ReasonCodes = tc.codes
			_, err := MapToEvaluationResult(result, validMapperTiming(), validMapperEvaluator(), validMapperSnapshotIdent())
			assertMappingError(t, err)
		})
	}
}

func TestMapper_InvalidDigest(t *testing.T) {
	t.Parallel()

	for _, dig := range []string{"", "abc", "ZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZ", mapperFrozenDigest[:63]} {
		dig := dig
		t.Run("digest_"+dig, func(t *testing.T) {
			t.Parallel()
			result := validMapperResult(OutcomeAllow)
			result.InputDigest = dig
			_, err := MapToEvaluationResult(result, validMapperTiming(), validMapperEvaluator(), validMapperSnapshotIdent())
			assertMappingError(t, err)
		})
	}
}

func TestMapper_RejectNonEmptyObligations(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		obs  []json.RawMessage
	}{
		{"scalar", []json.RawMessage{json.RawMessage(`"scalar"`)}},
		{"object", []json.RawMessage{json.RawMessage(`{"k":1}`)}},
		{"array", []json.RawMessage{json.RawMessage(`[1,2]`)}},
		{"malformed", []json.RawMessage{json.RawMessage(`not-json`)}},
		{"mixed", []json.RawMessage{json.RawMessage(``), json.RawMessage(`1`)}},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			result := validMapperResult(OutcomeAllow)
			result.Obligations = tc.obs
			_, err := MapToEvaluationResult(result, validMapperTiming(), validMapperEvaluator(), validMapperSnapshotIdent())
			assertMappingError(t, err)
		})
	}

	// nil and empty represent absence and must pass.
	for _, obs := range [][]json.RawMessage{nil, {}} {
		result := validMapperResult(OutcomeAllow)
		result.Obligations = obs
		if _, err := MapToEvaluationResult(result, validMapperTiming(), validMapperEvaluator(), validMapperSnapshotIdent()); err != nil {
			t.Fatalf("nil/empty obligations must pass, err=%v", err)
		}
	}
}

func TestMapper_FEATURE0013StructuralFailureViaBoolWrapper(t *testing.T) {
	t.Parallel()

	// Whitespace-only TypedRef fields pass apiref.Constraint.ValidateRef
	// (non-empty string checks) but fail FEATURE-0013 TrimSpace structural
	// rules inside IsEvaluationResultStructurallyValid.
	ident := InputIdentity{InputSnapshotRef: &apimeta.TypedRef{
		APIVersion: "governance.sovrunn.io/v1alpha1",
		Kind:       "InputSnapshot",
		Name:       "   ",
		UID:        "input-uid-1",
	}}
	_, err := MapToEvaluationResult(
		validMapperResult(OutcomeAllow),
		validMapperTiming(),
		validMapperEvaluator(),
		ident,
	)
	assertMappingError(t, err)
}

func TestMapper_FailureDoesNotMutateInput(t *testing.T) {
	t.Parallel()

	result := validMapperResult(OutcomeAllow)
	result.ReasonCodes = []string{"POLICY_ALLOW", "POLICY_EXTRA"}
	origCodes := append([]string(nil), result.ReasonCodes...)
	origDigest := result.InputDigest
	origAt := result.EvaluatedAt

	timing := validMapperTiming()
	timing.CompletedAt = time.Date(2026, 8, 25, 13, 0, 0, 0, time.UTC).Format(time.RFC3339Nano)

	before, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("marshal before: %v", err)
	}

	_, mapErr := MapToEvaluationResult(result, timing, validMapperEvaluator(), validMapperSnapshotIdent())
	assertMappingError(t, mapErr)

	after, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("marshal after: %v", err)
	}
	if !bytes.Equal(before, after) {
		t.Fatalf("policy result mutated on mapping failure\nbefore=%s\nafter=%s", before, after)
	}
	if result.InputDigest != origDigest || result.EvaluatedAt != origAt {
		t.Fatal("scalar fields mutated")
	}
	for i := range origCodes {
		if result.ReasonCodes[i] != origCodes[i] {
			t.Fatalf("reasonCodes mutated: %v", result.ReasonCodes)
		}
	}
}

func TestMapper_CopyOwnership(t *testing.T) {
	t.Parallel()

	result := validMapperResult(OutcomeAllow)
	result.ReasonCodes = []string{"POLICY_ALLOW", "POLICY_EXTRA"}
	ident := validMapperSnapshotIdent()
	integrity := &decision.OpaqueIntegrityCarrier{State: "bound"}
	ident.InputIntegrity = integrity

	got, err := MapToEvaluationResult(result, validMapperTiming(), validMapperEvaluator(), ident)
	if err != nil {
		t.Fatalf("MapToEvaluationResult: %v", err)
	}
	if got.InputSnapshotRef == ident.InputSnapshotRef {
		t.Fatal("output InputSnapshotRef must not alias caller pointer")
	}
	if got.InputIntegrity == integrity {
		t.Fatal("output InputIntegrity must not alias caller pointer")
	}

	ownedResult := append(json.RawMessage(nil), got.Result...)

	// Mutate caller-owned inputs after success.
	result.ReasonCodes[0] = "MUTATED"
	ident.InputSnapshotRef.Name = "mutated-name"
	integrity.State = "mutated"

	if got.InputSnapshotRef == nil || got.InputSnapshotRef.Name != "placement-input-1" {
		t.Fatalf("snapshot aliased or mutated: %+v", got.InputSnapshotRef)
	}
	if got.InputIntegrity == nil || got.InputIntegrity.State != "bound" {
		t.Fatalf("integrity aliased or mutated: %+v", got.InputIntegrity)
	}

	var embedded PolicyEvaluationResult
	if err := json.Unmarshal(ownedResult, &embedded); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if embedded.ReasonCodes[0] != "POLICY_ALLOW" {
		t.Fatalf("marshaled result aliased caller reasonCodes: %v", embedded.ReasonCodes)
	}

	// Mutating the returned raw JSON must not affect a re-marshal of owned inputs.
	got.Result[0] = 'X'
	if bytes.Equal(got.Result, ownedResult) {
		t.Fatal("expected distinct raw JSON buffer after mutation")
	}
}

func TestMapper_EvaluatedAtEqualsCompletedAt(t *testing.T) {
	t.Parallel()

	got, err := MapToEvaluationResult(
		validMapperResult(OutcomeIndeterminate),
		validMapperTiming(),
		validMapperEvaluator(),
		validMapperSnapshotIdent(),
	)
	if err != nil {
		t.Fatalf("MapToEvaluationResult: %v", err)
	}
	if got.EvaluatedAt != got.Timing.CompletedAt {
		t.Fatalf("evaluatedAt=%q != timing.completedAt=%q", got.EvaluatedAt, got.Timing.CompletedAt)
	}
}
