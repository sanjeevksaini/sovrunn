package decision

import (
	"encoding/json"
	"strings"
	"testing"
)

func sampleTrustCarrier() TrustCarrier {
	return TrustCarrier{
		State:                  "verified",
		CanonicalizationMethod: "canonicalization:profile-declared",
		DigestAlgorithm:        "digest:profile-declared",
		CoveredFields:          "covered:decision-body",
		SignatureAlgorithm:     "signature:profile-declared",
	}
}

func TestTrustCarrierJSONRoundTrip(t *testing.T) {
	t.Parallel()

	in := sampleTrustCarrier()
	data, err := json.Marshal(in)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var out TrustCarrier
	if err := json.Unmarshal(data, &out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if out != in {
		t.Fatalf("round-trip mismatch: got %+v want %+v", out, in)
	}
}

func TestTrustCarrierOmitemptyOptionalAlgorithmFields(t *testing.T) {
	t.Parallel()

	in := TrustCarrier{State: "opaque-token"}
	data, err := json.Marshal(in)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	raw := string(data)

	if !strings.Contains(raw, `"state"`) {
		t.Fatalf("required state must be present; json=%s", raw)
	}
	for _, key := range []string{
		`"canonicalizationMethod"`,
		`"digestAlgorithm"`,
		`"coveredFields"`,
		`"signatureAlgorithm"`,
	} {
		if strings.Contains(raw, key) {
			t.Fatalf("optional field %s must be omitted when unset; json=%s", key, raw)
		}
	}
}

func TestTrustCarrierPresentButEmptyIsStructurallyEmpty(t *testing.T) {
	t.Parallel()

	empty := TrustCarrier{}
	if !empty.IsStructurallyEmpty() {
		t.Fatal("zero-value TrustCarrier must be structurally empty (RID-08)")
	}

	populated := sampleTrustCarrier()
	if populated.IsStructurallyEmpty() {
		t.Fatal("populated TrustCarrier must not be structurally empty")
	}

	stateOnly := TrustCarrier{State: "unknown"}
	if stateOnly.IsStructurallyEmpty() {
		t.Fatal("carrier with state set must not be structurally empty")
	}
}

func TestTrustCarrierNoCryptographicBehavior(t *testing.T) {
	t.Parallel()

	// Structural round-trip only: identifiers are opaque tokens and are not
	// interpreted, selected, computed, signed, or verified.
	in := TrustCarrier{
		State:                  "unverified",
		CanonicalizationMethod: "illustrative-only",
		DigestAlgorithm:        "illustrative-only",
		CoveredFields:          "/record",
		SignatureAlgorithm:     "illustrative-only",
	}
	data, err := json.Marshal(in)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var out TrustCarrier
	if err := json.Unmarshal(data, &out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if out.DigestAlgorithm != "illustrative-only" || out.SignatureAlgorithm != "illustrative-only" {
		t.Fatalf("algorithm identifiers must remain opaque structural tokens; got %+v", out)
	}
}

func TestDecisionBodyTrustOmitempty(t *testing.T) {
	t.Parallel()

	in := sampleDecisionRecord()
	in.Record.Trust = nil

	data, err := json.Marshal(in)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if strings.Contains(string(data), `"trust"`) {
		t.Fatalf("optional trust must be omitted when unset; json=%s", data)
	}

	carrier := sampleTrustCarrier()
	in.Record.Trust = &carrier
	data, err = json.Marshal(in)
	if err != nil {
		t.Fatalf("marshal with trust: %v", err)
	}
	var out DecisionRecord
	if err := json.Unmarshal(data, &out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if out.Record.Trust == nil || *out.Record.Trust != carrier {
		t.Fatalf("trust = %#v, want %#v", out.Record.Trust, carrier)
	}
}
