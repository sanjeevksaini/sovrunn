package decision

import (
	"encoding/json"
	"strings"
	"testing"
)

func sampleSemanticDecisionIdentity() SemanticDecisionIdentity {
	return SemanticDecisionIdentity{
		Basis:          "placement-v1",
		InputRefs:      []string{"input-snapshot-1", "policy-context-1"},
		CanonicalScope: "Project:project-123",
		Digest:         "opaque-identity-digest-token",
	}
}

func TestSemanticDecisionIdentityJSONRoundTrip(t *testing.T) {
	t.Parallel()

	in := sampleSemanticDecisionIdentity()
	data, err := json.Marshal(in)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var out SemanticDecisionIdentity
	if err := json.Unmarshal(data, &out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if out.Basis != in.Basis ||
		out.CanonicalScope != in.CanonicalScope ||
		out.Digest != in.Digest ||
		len(out.InputRefs) != len(in.InputRefs) ||
		out.InputRefs[0] != in.InputRefs[0] ||
		out.InputRefs[1] != in.InputRefs[1] {
		t.Fatalf("round-trip mismatch: got %+v want %+v", out, in)
	}
}

func TestSemanticDecisionIdentityDigestOmitempty(t *testing.T) {
	t.Parallel()

	in := SemanticDecisionIdentity{
		Basis:          "placement-v1",
		InputRefs:      []string{"input-1"},
		CanonicalScope: "Platform",
	}
	data, err := json.Marshal(in)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	raw := string(data)
	if strings.Contains(raw, `"digest"`) {
		t.Fatalf("optional digest must be omitted when unset; json=%s", raw)
	}
	for _, key := range []string{`"basis"`, `"inputRefs"`, `"canonicalScope"`} {
		if !strings.Contains(raw, key) {
			t.Fatalf("required field %s missing; json=%s", key, raw)
		}
	}
}

func TestRetryKeyDistinctFromSemanticIdentity(t *testing.T) {
	t.Parallel()

	// AD-036 / F13-SCALE-001/002: retry idempotency and semantic identity are
	// separate carriers. Equal retry keys must not imply equal semantic identity.
	retryA := RetryIdempotencyKey("retry-shared")
	retryB := RetryIdempotencyKey("retry-shared")
	if retryA != retryB {
		t.Fatal("identical retry keys must compare equal as retry identity")
	}

	identityA := SemanticDecisionIdentity{
		Basis:          "placement-v1",
		InputRefs:      []string{"input-a"},
		CanonicalScope: "Project:project-123",
	}
	identityB := SemanticDecisionIdentity{
		Basis:          "placement-v1",
		InputRefs:      []string{"input-b"},
		CanonicalScope: "Project:project-123",
	}
	if identityA.Basis == identityB.Basis &&
		len(identityA.InputRefs) == len(identityB.InputRefs) &&
		identityA.InputRefs[0] == identityB.InputRefs[0] {
		t.Fatal("test setup error: identities must differ")
	}

	recA := sampleDecisionRecord()
	recA.Record.RetryKey = string(retryA)
	recA.Record.SemanticIdentity = identityA

	recB := sampleDecisionRecord()
	recB.Record.RetryKey = string(retryB)
	recB.Record.SemanticIdentity = identityB

	if recA.Record.RetryKey != recB.Record.RetryKey {
		t.Fatal("retry keys should match for this case")
	}
	if recA.Record.SemanticIdentity.InputRefs[0] == recB.Record.SemanticIdentity.InputRefs[0] {
		t.Fatal("semantic identities must remain distinct despite equal retry keys")
	}

	data, err := json.Marshal(recA.Record)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	raw := string(data)
	if !strings.Contains(raw, `"retryKey"`) {
		t.Fatalf("retryKey field missing; json=%s", raw)
	}
	if !strings.Contains(raw, `"semanticIdentity"`) {
		t.Fatalf("semanticIdentity field missing; json=%s", raw)
	}

	var body map[string]json.RawMessage
	if err := json.Unmarshal(data, &body); err != nil {
		t.Fatalf("unmarshal body: %v", err)
	}
	if _, ok := body["retryKey"]; !ok {
		t.Fatal("retryKey must be a distinct JSON field")
	}
	if _, ok := body["semanticIdentity"]; !ok {
		t.Fatal("semanticIdentity must be a distinct JSON field")
	}
}

func TestCanonicalScopeIsDerivedDescriptorNotAuthority(t *testing.T) {
	t.Parallel()

	// CanonicalScope is a derived descriptor string. It must not introduce a
	// second scope source on the identity carrier.
	in := sampleSemanticDecisionIdentity()
	data, err := json.Marshal(in)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var asMap map[string]json.RawMessage
	if err := json.Unmarshal(data, &asMap); err != nil {
		t.Fatalf("unmarshal map: %v", err)
	}
	for _, forbidden := range []string{"scopeRef", "scope", "ownerScope", "subjectScope", "auditScope"} {
		if _, ok := asMap[forbidden]; ok {
			t.Fatalf("SemanticDecisionIdentity must not carry independent scope field %q", forbidden)
		}
	}
	if _, ok := asMap["canonicalScope"]; !ok {
		t.Fatal("canonicalScope derived descriptor must be present")
	}
}
