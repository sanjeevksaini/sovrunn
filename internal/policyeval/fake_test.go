package policyeval

import (
	"context"
	"errors"
	"strings"
	"testing"
)

const (
	fakeDigestAllow            = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	fakeDigestDeny             = "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
	fakeDigestIndeterminate    = "cccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccc"
	fakeDigestRequiresApproval = "dddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd"
	fakeDigestFailure          = "eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee"
	fakeDigestUnused           = "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"
)

func fakeConclusion(outcome Outcome, codes ...string) *Conclusion {
	return &Conclusion{Outcome: outcome, ReasonCodes: append([]string(nil), codes...)}
}

func mustFake(t *testing.T, fixtures []FakeFixture) *DeterministicFake {
	t.Helper()
	f, err := NewDeterministicFake(fixtures)
	if err != nil {
		t.Fatalf("NewDeterministicFake: %v", err)
	}
	return f
}

func assertFakeAdapterFailure(t *testing.T, got Conclusion, err error) {
	t.Helper()
	if err == nil {
		t.Fatal("expected adapter failure error, got nil")
	}
	if !errors.Is(err, errFakeAdapterFailure) {
		t.Fatalf("err=%v, want %v", err, errFakeAdapterFailure)
	}
	if got.Outcome != "" || got.ReasonCodes != nil {
		t.Fatalf("conclusion not zero on failure: %+v", got)
	}
}

func TestFake_AllFourOutcomes(t *testing.T) {
	t.Parallel()

	fixtures := []FakeFixture{
		{InputDigest: fakeDigestAllow, Conclusion: fakeConclusion(OutcomeAllow, "POLICY_ALLOW")},
		{InputDigest: fakeDigestDeny, Conclusion: fakeConclusion(OutcomeDeny, "POLICY_DENY")},
		{InputDigest: fakeDigestIndeterminate, Conclusion: fakeConclusion(OutcomeIndeterminate, "POLICY_INDETERMINATE")},
		{InputDigest: fakeDigestRequiresApproval, Conclusion: fakeConclusion(OutcomeRequiresApproval, "POLICY_REQUIRES_APPROVAL")},
	}
	f := mustFake(t, fixtures)

	cases := []struct {
		digest  string
		outcome Outcome
		code    string
	}{
		{fakeDigestAllow, OutcomeAllow, "POLICY_ALLOW"},
		{fakeDigestDeny, OutcomeDeny, "POLICY_DENY"},
		{fakeDigestIndeterminate, OutcomeIndeterminate, "POLICY_INDETERMINATE"},
		{fakeDigestRequiresApproval, OutcomeRequiresApproval, "POLICY_REQUIRES_APPROVAL"},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(string(tc.outcome), func(t *testing.T) {
			t.Parallel()
			got, err := f.Evaluate(context.Background(), NormalizedInput{}, tc.digest)
			if err != nil {
				t.Fatalf("Evaluate: %v", err)
			}
			if got.Outcome != tc.outcome {
				t.Fatalf("outcome=%q, want %q", got.Outcome, tc.outcome)
			}
			if len(got.ReasonCodes) != 1 || got.ReasonCodes[0] != tc.code {
				t.Fatalf("reasonCodes=%v, want [%q]", got.ReasonCodes, tc.code)
			}
		})
	}
}

func TestFake_MissingDigestAdapterFailure(t *testing.T) {
	t.Parallel()

	f := mustFake(t, []FakeFixture{
		{InputDigest: fakeDigestAllow, Conclusion: fakeConclusion(OutcomeAllow, "POLICY_ALLOW")},
	})

	got, err := f.Evaluate(context.Background(), NormalizedInput{}, fakeDigestUnused)
	assertFakeAdapterFailure(t, got, err)
}

func TestFake_EmptyFixturesEveryDigestUnconfigured(t *testing.T) {
	t.Parallel()

	f := mustFake(t, nil)
	got, err := f.Evaluate(context.Background(), NormalizedInput{}, fakeDigestAllow)
	assertFakeAdapterFailure(t, got, err)

	f2 := mustFake(t, []FakeFixture{})
	got2, err2 := f2.Evaluate(context.Background(), NormalizedInput{}, fakeDigestAllow)
	assertFakeAdapterFailure(t, got2, err2)
}

func TestFake_ConfiguredAdapterFailure(t *testing.T) {
	t.Parallel()

	f := mustFake(t, []FakeFixture{
		{InputDigest: fakeDigestFailure, AdapterFailure: true},
	})

	got, err := f.Evaluate(context.Background(), NormalizedInput{}, fakeDigestFailure)
	assertFakeAdapterFailure(t, got, err)
}

func TestFake_DuplicateDigestFailsConstruction(t *testing.T) {
	t.Parallel()

	_, err := NewDeterministicFake([]FakeFixture{
		{InputDigest: fakeDigestAllow, Conclusion: fakeConclusion(OutcomeAllow, "POLICY_ALLOW")},
		{InputDigest: fakeDigestAllow, AdapterFailure: true},
	})
	if !errors.Is(err, errFakeConfigurationInvalid) {
		t.Fatalf("err=%v, want %v", err, errFakeConfigurationInvalid)
	}
	if strings.Contains(err.Error(), fakeDigestAllow) {
		t.Fatalf("construction error echoed digest: %v", err)
	}
}

func TestFake_MalformedDigestFailsConstruction(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name   string
		digest string
	}{
		{"empty", ""},
		{"too_short", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},
		{"too_long", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaab"},
		{"uppercase", "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"},
		{"non_hex", "gggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggg"},
		{"mixed_case", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaAA"},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			_, err := NewDeterministicFake([]FakeFixture{
				{InputDigest: tc.digest, Conclusion: fakeConclusion(OutcomeAllow, "POLICY_ALLOW")},
			})
			if !errors.Is(err, errFakeConfigurationInvalid) {
				t.Fatalf("err=%v, want %v", err, errFakeConfigurationInvalid)
			}
			if tc.digest != "" && strings.Contains(err.Error(), tc.digest) {
				t.Fatalf("construction error echoed digest: %v", err)
			}
		})
	}
}

func TestFake_InvalidOneOfFailsConstruction(t *testing.T) {
	t.Parallel()

	t.Run("neither", func(t *testing.T) {
		t.Parallel()
		_, err := NewDeterministicFake([]FakeFixture{
			{InputDigest: fakeDigestAllow},
		})
		if !errors.Is(err, errFakeConfigurationInvalid) {
			t.Fatalf("err=%v, want %v", err, errFakeConfigurationInvalid)
		}
	})

	t.Run("both", func(t *testing.T) {
		t.Parallel()
		_, err := NewDeterministicFake([]FakeFixture{
			{
				InputDigest:    fakeDigestAllow,
				Conclusion:     fakeConclusion(OutcomeAllow, "POLICY_ALLOW"),
				AdapterFailure: true,
			},
		})
		if !errors.Is(err, errFakeConfigurationInvalid) {
			t.Fatalf("err=%v, want %v", err, errFakeConfigurationInvalid)
		}
	})
}

func TestFake_InvalidConclusionFailsConstruction(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name       string
		conclusion *Conclusion
	}{
		{"bad_outcome", &Conclusion{Outcome: Outcome("Unknown"), ReasonCodes: []string{"POLICY_ALLOW"}}},
		{"empty_outcome", &Conclusion{Outcome: "", ReasonCodes: []string{"POLICY_ALLOW"}}},
		{"no_reason_codes", &Conclusion{Outcome: OutcomeAllow, ReasonCodes: nil}},
		{"empty_reason_codes", &Conclusion{Outcome: OutcomeAllow, ReasonCodes: []string{}}},
		{"malformed_reason", &Conclusion{Outcome: OutcomeAllow, ReasonCodes: []string{"not-valid"}}},
		{"duplicate_reason", &Conclusion{Outcome: OutcomeAllow, ReasonCodes: []string{"POLICY_ALLOW", "POLICY_ALLOW"}}},
		{"too_many_reasons", &Conclusion{
			Outcome: OutcomeAllow,
			ReasonCodes: func() []string {
				codes := make([]string, MaxReasonCodes+1)
				for i := range codes {
					codes[i] = "CODE_" + string(rune('A'+i/26)) + string(rune('A'+i%26))
				}
				return codes
			}(),
		}},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			_, err := NewDeterministicFake([]FakeFixture{
				{InputDigest: fakeDigestAllow, Conclusion: tc.conclusion},
			})
			if !errors.Is(err, errFakeConfigurationInvalid) {
				t.Fatalf("err=%v, want %v", err, errFakeConfigurationInvalid)
			}
			for _, code := range tc.conclusion.ReasonCodes {
				if code != "" && strings.Contains(err.Error(), code) {
					t.Fatalf("construction error echoed reason code %q: %v", code, err)
				}
			}
		})
	}
}

func TestFake_ConstructorDeepCopyOwnership(t *testing.T) {
	t.Parallel()

	codes := []string{"POLICY_ALLOW", "POLICY_EXTRA"}
	conclusion := &Conclusion{Outcome: OutcomeAllow, ReasonCodes: codes}
	fixtures := []FakeFixture{
		{InputDigest: fakeDigestAllow, Conclusion: conclusion},
	}

	f := mustFake(t, fixtures)

	// Mutate caller-owned fixture after construction.
	codes[0] = "MUTATED_CODE"
	conclusion.Outcome = OutcomeDeny
	conclusion.ReasonCodes = []string{"MUTATED_AFTER"}
	fixtures[0].InputDigest = fakeDigestDeny
	fixtures[0].AdapterFailure = true
	fixtures[0].Conclusion = nil

	got, err := f.Evaluate(context.Background(), NormalizedInput{}, fakeDigestAllow)
	if err != nil {
		t.Fatalf("Evaluate after fixture mutation: %v", err)
	}
	if got.Outcome != OutcomeAllow {
		t.Fatalf("outcome=%q, want Allow (fixture mutation leaked)", got.Outcome)
	}
	if len(got.ReasonCodes) != 2 || got.ReasonCodes[0] != "POLICY_ALLOW" || got.ReasonCodes[1] != "POLICY_EXTRA" {
		t.Fatalf("reasonCodes=%v, want [POLICY_ALLOW POLICY_EXTRA]", got.ReasonCodes)
	}

	// Original digest must still succeed; mutated digest must be unconfigured.
	failGot, failErr := f.Evaluate(context.Background(), NormalizedInput{}, fakeDigestDeny)
	assertFakeAdapterFailure(t, failGot, failErr)
}

func TestFake_ReturnedConclusionDeepCopyOwnership(t *testing.T) {
	t.Parallel()

	f := mustFake(t, []FakeFixture{
		{InputDigest: fakeDigestAllow, Conclusion: fakeConclusion(OutcomeAllow, "POLICY_ALLOW", "POLICY_EXTRA")},
	})

	first, err := f.Evaluate(context.Background(), NormalizedInput{}, fakeDigestAllow)
	if err != nil {
		t.Fatalf("Evaluate: %v", err)
	}
	first.Outcome = OutcomeDeny
	first.ReasonCodes[0] = "MUTATED"
	first.ReasonCodes = append(first.ReasonCodes, "APPENDED")

	second, err := f.Evaluate(context.Background(), NormalizedInput{}, fakeDigestAllow)
	if err != nil {
		t.Fatalf("second Evaluate: %v", err)
	}
	if second.Outcome != OutcomeAllow {
		t.Fatalf("outcome=%q, want Allow (returned mutation leaked)", second.Outcome)
	}
	if len(second.ReasonCodes) != 2 || second.ReasonCodes[0] != "POLICY_ALLOW" || second.ReasonCodes[1] != "POLICY_EXTRA" {
		t.Fatalf("reasonCodes=%v, want original configured set", second.ReasonCodes)
	}
}

func TestFake_IgnoresNormalizedInput(t *testing.T) {
	t.Parallel()

	f := mustFake(t, []FakeFixture{
		{InputDigest: fakeDigestAllow, Conclusion: fakeConclusion(OutcomeAllow, "POLICY_ALLOW")},
	})

	// Different semantic inputs with the same digest must not change lookup.
	got, err := f.Evaluate(context.Background(), NormalizedInput{Action: "create"}, fakeDigestAllow)
	if err != nil {
		t.Fatalf("Evaluate: %v", err)
	}
	if got.Outcome != OutcomeAllow {
		t.Fatalf("outcome=%q, want Allow", got.Outcome)
	}
}

func TestFake_PreservesConfiguredReasonCodeOrder(t *testing.T) {
	t.Parallel()

	// Fake returns configured order; boundary sorting is elsewhere.
	f := mustFake(t, []FakeFixture{
		{InputDigest: fakeDigestAllow, Conclusion: fakeConclusion(OutcomeAllow, "ZULU", "ALPHA")},
	})
	got, err := f.Evaluate(context.Background(), NormalizedInput{}, fakeDigestAllow)
	if err != nil {
		t.Fatalf("Evaluate: %v", err)
	}
	if len(got.ReasonCodes) != 2 || got.ReasonCodes[0] != "ZULU" || got.ReasonCodes[1] != "ALPHA" {
		t.Fatalf("reasonCodes=%v, want [ZULU ALPHA]", got.ReasonCodes)
	}
}
