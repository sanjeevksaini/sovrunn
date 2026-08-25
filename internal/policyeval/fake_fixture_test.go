package policyeval

import (
	"bytes"
	"errors"
	"io"
	"testing"

	"gopkg.in/yaml.v3"
)

// Fixed diagnostics for the test-only D-13 YAML rules. These are not
// production API codes; semantic fixture validation remains in
// NewDeterministicFake.
var (
	errFakeFixtureYAMLRequired          = errors.New("YAML document is required")
	errFakeFixtureYAMLMultipleDocuments = errors.New("multiple YAML documents are not permitted")
	errFakeFixtureYAMLTrailingContent   = errors.New("trailing YAML content is not permitted")
)

// decodeFakeFixturesYAMLForTest is the D-13 test-only strict YAML helper.
// It uses gopkg.in/yaml.v3 with KnownFields(true), rejects duplicate mapping
// keys (yaml.v3 typed decode), requires exactly one document, and rejects
// trailing content. It returns a duplicate-preserving []FakeFixture; all
// semantic validation remains in NewDeterministicFake. This is not a
// production YAML input surface.
func decodeFakeFixturesYAMLForTest(data []byte) ([]FakeFixture, error) {
	dec := yaml.NewDecoder(bytes.NewReader(data))
	dec.KnownFields(true)

	var fixtures []FakeFixture
	if err := dec.Decode(&fixtures); err != nil {
		if errors.Is(err, io.EOF) {
			return nil, errFakeFixtureYAMLRequired
		}
		return nil, err
	}

	var trailing any
	err := dec.Decode(&trailing)
	switch {
	case err == nil:
		return nil, errFakeFixtureYAMLMultipleDocuments
	case errors.Is(err, io.EOF):
		return fixtures, nil
	default:
		return nil, errFakeFixtureYAMLTrailingContent
	}
}

func TestFakeFixtureYAML_Valid(t *testing.T) {
	t.Parallel()

	raw := []byte(`
- inputDigest: aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa
  conclusion:
    outcome: Allow
    reasonCodes:
      - POLICY_ALLOW
- inputDigest: eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee
  adapterFailure: true
`)

	got, err := decodeFakeFixturesYAMLForTest(raw)
	if err != nil {
		t.Fatalf("decodeFakeFixturesYAMLForTest: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("len=%d, want 2", len(got))
	}
	if got[0].InputDigest != fakeDigestAllow {
		t.Fatalf("fixtures[0].InputDigest=%q, want %q", got[0].InputDigest, fakeDigestAllow)
	}
	if got[0].Conclusion == nil {
		t.Fatal("fixtures[0].Conclusion is nil")
	}
	if got[0].Conclusion.Outcome != OutcomeAllow {
		t.Fatalf("outcome=%q, want %q", got[0].Conclusion.Outcome, OutcomeAllow)
	}
	if len(got[0].Conclusion.ReasonCodes) != 1 || got[0].Conclusion.ReasonCodes[0] != "POLICY_ALLOW" {
		t.Fatalf("reasonCodes=%v, want [POLICY_ALLOW]", got[0].Conclusion.ReasonCodes)
	}
	if got[0].AdapterFailure {
		t.Fatal("fixtures[0].AdapterFailure unexpectedly true")
	}
	if got[1].InputDigest != fakeDigestFailure {
		t.Fatalf("fixtures[1].InputDigest=%q, want %q", got[1].InputDigest, fakeDigestFailure)
	}
	if got[1].Conclusion != nil {
		t.Fatalf("fixtures[1].Conclusion=%v, want nil", got[1].Conclusion)
	}
	if !got[1].AdapterFailure {
		t.Fatal("fixtures[1].AdapterFailure unexpectedly false")
	}

	if _, err := NewDeterministicFake(got); err != nil {
		t.Fatalf("NewDeterministicFake after valid YAML: %v", err)
	}
}

func TestFakeFixtureYAML_UnknownField(t *testing.T) {
	t.Parallel()

	raw := []byte(`
- inputDigest: aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa
  unknownField: true
  adapterFailure: true
`)

	_, err := decodeFakeFixturesYAMLForTest(raw)
	if err == nil {
		t.Fatal("expected unknown-field rejection, got nil")
	}
}

func TestFakeFixtureYAML_DuplicateKey(t *testing.T) {
	t.Parallel()

	raw := []byte(`
- inputDigest: aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa
  inputDigest: bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb
  adapterFailure: true
`)

	_, err := decodeFakeFixturesYAMLForTest(raw)
	if err == nil {
		t.Fatal("expected duplicate-key rejection, got nil")
	}
}

func TestFakeFixtureYAML_MultipleDocuments(t *testing.T) {
	t.Parallel()

	raw := []byte(`
- inputDigest: aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa
  adapterFailure: true
---
- inputDigest: bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb
  adapterFailure: true
`)

	_, err := decodeFakeFixturesYAMLForTest(raw)
	if !errors.Is(err, errFakeFixtureYAMLMultipleDocuments) {
		t.Fatalf("err=%v, want %v", err, errFakeFixtureYAMLMultipleDocuments)
	}
}

func TestFakeFixtureYAML_TrailingContent(t *testing.T) {
	t.Parallel()

	// Complete one document, then non-document trailing content that yields a
	// non-EOF second Decode (D-13 trailing-content rejection).
	raw := []byte(`
[{inputDigest: aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa, adapterFailure: true}]
true
`)

	_, err := decodeFakeFixturesYAMLForTest(raw)
	if !errors.Is(err, errFakeFixtureYAMLTrailingContent) {
		t.Fatalf("err=%v, want %v", err, errFakeFixtureYAMLTrailingContent)
	}
}

func TestFakeFixtureYAML_DuplicateFixtureSlicePreserved(t *testing.T) {
	t.Parallel()

	raw := []byte(`
- inputDigest: aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa
  conclusion:
    outcome: Allow
    reasonCodes:
      - POLICY_ALLOW
- inputDigest: aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa
  adapterFailure: true
`)

	got, err := decodeFakeFixturesYAMLForTest(raw)
	if err != nil {
		t.Fatalf("decodeFakeFixturesYAMLForTest: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("len=%d, want 2 (duplicate-preserving)", len(got))
	}
	if got[0].InputDigest != fakeDigestAllow || got[1].InputDigest != fakeDigestAllow {
		t.Fatalf("digests=%q,%q, want both %q", got[0].InputDigest, got[1].InputDigest, fakeDigestAllow)
	}
	if got[0].Conclusion == nil || got[1].AdapterFailure != true {
		t.Fatalf("unexpected fixture bodies: %+v", got)
	}

	_, err = NewDeterministicFake(got)
	if !errors.Is(err, errFakeConfigurationInvalid) {
		t.Fatalf("constructor err=%v, want %v", err, errFakeConfigurationInvalid)
	}
}
