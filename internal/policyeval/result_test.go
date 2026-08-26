package policyeval

import (
	"bytes"
	"encoding/json"
	"regexp"
	"testing"
)

func TestResultOutcomeExactlyFourValues(t *testing.T) {
	t.Parallel()

	all := AllOutcomes()
	if len(all) != 4 {
		t.Fatalf("AllOutcomes len = %d, want 4", len(all))
	}

	want := map[Outcome]struct{}{
		OutcomeAllow:            {},
		OutcomeDeny:             {},
		OutcomeIndeterminate:    {},
		OutcomeRequiresApproval: {},
	}
	seen := make(map[Outcome]struct{}, len(all))
	for _, o := range all {
		if !o.Valid() {
			t.Fatalf("outcome %q reported invalid", o)
		}
		if _, ok := want[o]; !ok {
			t.Fatalf("unexpected outcome %q", o)
		}
		if _, dup := seen[o]; dup {
			t.Fatalf("duplicate outcome %q", o)
		}
		seen[o] = struct{}{}
	}
	if len(seen) != 4 {
		t.Fatalf("unique outcomes = %d, want 4", len(seen))
	}

	if Outcome("").Valid() {
		t.Fatal("empty Outcome must be invalid")
	}
	if Outcome("Unknown").Valid() {
		t.Fatal(`Outcome("Unknown") must be invalid`)
	}
	if Outcome("AdapterFailure").Valid() {
		t.Fatal("NonResult AdapterFailure must not be a valid Outcome")
	}
}

func TestResultNonResultExactlyFiveCategories(t *testing.T) {
	t.Parallel()

	all := AllNonResults()
	if len(all) != 5 {
		t.Fatalf("AllNonResults len = %d, want 5", len(all))
	}

	want := map[NonResult]struct{}{
		NonResultRequestInvalid:       {},
		NonResultCanceled:             {},
		NonResultDeadlineExceeded:     {},
		NonResultAdapterFailure:       {},
		NonResultInvalidAdapterResult: {},
	}
	seen := make(map[NonResult]struct{}, len(all))
	for _, n := range all {
		if !n.Valid() {
			t.Fatalf("non-result %q reported invalid", n)
		}
		if _, ok := want[n]; !ok {
			t.Fatalf("unexpected non-result %q", n)
		}
		if _, dup := seen[n]; dup {
			t.Fatalf("duplicate non-result %q", n)
		}
		seen[n] = struct{}{}
	}
	if len(seen) != 5 {
		t.Fatalf("unique non-results = %d, want 5", len(seen))
	}

	if NonResult("").Valid() {
		t.Fatal("empty NonResult must be invalid")
	}
	if NonResult("Allow").Valid() {
		t.Fatal("Outcome Allow must not be a valid NonResult")
	}
}

func TestResultReasonCodeGrammarRegex(t *testing.T) {
	t.Parallel()

	re, err := regexp.Compile(ReasonCodePattern)
	if err != nil {
		t.Fatalf("ReasonCodePattern compile: %v", err)
	}
	if ReasonCodePattern != `^[A-Z][A-Z0-9_]{0,62}$` {
		t.Fatalf("ReasonCodePattern = %q, want ^[A-Z][A-Z0-9_]{0,62}$", ReasonCodePattern)
	}
	if ReasonCodeRegexp.String() != ReasonCodePattern {
		t.Fatalf("ReasonCodeRegexp pattern = %q, want %q", ReasonCodeRegexp.String(), ReasonCodePattern)
	}

	valid := []string{
		"A",
		"ALLOW",
		"POLICY_DENY",
		"R" + string(bytes.Repeat([]byte("0"), 62)), // 63 chars total
	}
	for _, code := range valid {
		if !re.MatchString(code) {
			t.Fatalf("expected valid reason code %q", code)
		}
		if !ReasonCodeRegexp.MatchString(code) {
			t.Fatalf("ReasonCodeRegexp rejected valid code %q", code)
		}
		if len(code) > MaxReasonCodeLen {
			t.Fatalf("fixture longer than MaxReasonCodeLen: %d", len(code))
		}
	}

	invalid := []string{
		"",
		"allow",
		"9START",
		"_LEADING",
		"HAS-DASH",
		"HAS.DOT",
		"HAS SPACE",
		"Å",
		"AÅ",
		"A" + string(bytes.Repeat([]byte("0"), 63)), // 64 chars
	}
	for _, code := range invalid {
		if re.MatchString(code) {
			t.Fatalf("expected invalid reason code %q", code)
		}
		if ReasonCodeRegexp.MatchString(code) {
			t.Fatalf("ReasonCodeRegexp accepted invalid code %q", code)
		}
	}

	if MinReasonCodes != 1 || MaxReasonCodes != 32 {
		t.Fatalf("reason-code count bounds = %d..%d, want 1..32", MinReasonCodes, MaxReasonCodes)
	}
}

func TestResultObligationsOmittedWhenNilOrEmpty(t *testing.T) {
	t.Parallel()

	base := PolicyEvaluationResult{
		Outcome:     OutcomeAllow,
		ReasonCodes: []string{"ALLOW"},
		InputDigest: "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
		EvaluatedAt: "2026-01-01T00:00:00Z",
	}

	cases := []struct {
		name        string
		obligations []json.RawMessage
	}{
		{name: "nil", obligations: nil},
		{name: "empty", obligations: []json.RawMessage{}},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			r := base
			r.Obligations = tc.obligations
			raw, err := json.Marshal(r)
			if err != nil {
				t.Fatalf("json.Marshal: %v", err)
			}
			if bytes.Contains(raw, []byte(`"obligations"`)) {
				t.Fatalf("obligations key present in JSON for %s: %s", tc.name, raw)
			}
			var obj map[string]json.RawMessage
			if err := json.Unmarshal(raw, &obj); err != nil {
				t.Fatalf("json.Unmarshal: %v", err)
			}
			if _, ok := obj["obligations"]; ok {
				t.Fatalf("obligations member present for %s: %s", tc.name, raw)
			}
		})
	}
}
