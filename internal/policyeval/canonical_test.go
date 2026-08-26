package policyeval

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
)

// Frozen ASCII conformance preimage and digest from CDG-F17-02 / AC-F17-10.
const (
	frozenCanonicalPreimage = `{"action":"service.read","candidateRefs":[],"profileRefs":[{"apiVersion":"iam.sovrunn.io/v1alpha1","kind":"RoleDefinition","name":"reader","uid":"role-001"}],"schema":"sovrunn.policy-evaluation-request/v1","subjectRef":{"apiVersion":"services.sovrunn.io/v1alpha1","kind":"ServiceInstance","name":"reporting-api","uid":"service-001"}}`
	frozenInputDigest       = "daca15fd0310c46b45d5aff9bbe4f1a5dedd4b788c44c3780838a8be40a56103"
)

func frozenVectorRequest() PolicyEvaluationRequest {
	return PolicyEvaluationRequest{
		SubjectRef: apimeta.TypedRef{
			APIVersion: "services.sovrunn.io/v1alpha1",
			Kind:       "ServiceInstance",
			Name:       "reporting-api",
			UID:        "service-001",
		},
		Action: "service.read",
		ProfileRefs: []apimeta.TypedRef{{
			APIVersion: "iam.sovrunn.io/v1alpha1",
			Kind:       "RoleDefinition",
			Name:       "reader",
			UID:        "role-001",
		}},
		RequestID: "corr-001",
	}
}

func TestCanonicalFrozenASCIIVector(t *testing.T) {
	t.Parallel()

	req := frozenVectorRequest()
	gotBytes, err := canonicalPolicyEvaluationRequestBytes(req)
	if err != nil {
		t.Fatalf("canonicalPolicyEvaluationRequestBytes: %v", err)
	}
	if !bytes.Equal(gotBytes, []byte(frozenCanonicalPreimage)) {
		t.Fatalf("canonical bytes mismatch\ngot:  %s\nwant: %s", gotBytes, frozenCanonicalPreimage)
	}

	gotDigest, err := digestPolicyEvaluationRequest(req)
	if err != nil {
		t.Fatalf("digestPolicyEvaluationRequest: %v", err)
	}
	if gotDigest != frozenInputDigest {
		t.Fatalf("digest=%q, want %q", gotDigest, frozenInputDigest)
	}
	if len(gotDigest) != 64 {
		t.Fatalf("digest length=%d, want 64", len(gotDigest))
	}
	if gotDigest != strings.ToLower(gotDigest) {
		t.Fatalf("digest must be lowercase hex: %q", gotDigest)
	}

	sum := sha256.Sum256([]byte(frozenCanonicalPreimage))
	if hex.EncodeToString(sum[:]) != frozenInputDigest {
		t.Fatal("test fixture preimage/digest mismatch")
	}
}

func TestCanonicalNoWhitespaceBOMOrTrailingNewline(t *testing.T) {
	t.Parallel()

	got, err := canonicalPolicyEvaluationRequestBytes(frozenVectorRequest())
	if err != nil {
		t.Fatalf("canonicalPolicyEvaluationRequestBytes: %v", err)
	}
	if bytes.ContainsAny(got, " \t\n\r") {
		t.Fatalf("canonical bytes contain whitespace: %q", got)
	}
	if bytes.HasPrefix(got, []byte{0xEF, 0xBB, 0xBF}) {
		t.Fatal("canonical bytes must not start with BOM")
	}
	if len(got) == 0 || got[len(got)-1] == '\n' {
		t.Fatal("canonical bytes must not end with trailing newline")
	}
	if !bytes.HasPrefix(got, []byte{'{'}) || !bytes.HasSuffix(got, []byte{'}'}) {
		t.Fatalf("canonical bytes must be a single JSON object, got %q", got)
	}
	if bytes.Contains(got, []byte("requestId")) {
		t.Fatal("requestId must be excluded from normalized object")
	}
}

func TestCanonicalOmittedVsEmptyCandidateRefs(t *testing.T) {
	t.Parallel()

	omitted := frozenVectorRequest()
	omitted.CandidateRefs = nil

	empty := frozenVectorRequest()
	empty.CandidateRefs = []apimeta.TypedRef{}

	omittedBytes, err := canonicalPolicyEvaluationRequestBytes(omitted)
	if err != nil {
		t.Fatalf("omitted candidates: %v", err)
	}
	emptyBytes, err := canonicalPolicyEvaluationRequestBytes(empty)
	if err != nil {
		t.Fatalf("empty candidates: %v", err)
	}
	if !bytes.Equal(omittedBytes, emptyBytes) {
		t.Fatalf("omitted vs [] candidateRefs diverge\nomitted: %s\nempty:   %s", omittedBytes, emptyBytes)
	}
	if !bytes.Contains(omittedBytes, []byte(`"candidateRefs":[]`)) {
		t.Fatalf("expected candidateRefs:[], got %s", omittedBytes)
	}

	omittedDigest, err := digestPolicyEvaluationRequest(omitted)
	if err != nil {
		t.Fatalf("omitted digest: %v", err)
	}
	emptyDigest, err := digestPolicyEvaluationRequest(empty)
	if err != nil {
		t.Fatalf("empty digest: %v", err)
	}
	if omittedDigest != emptyDigest {
		t.Fatalf("digests differ: omitted=%q empty=%q", omittedDigest, emptyDigest)
	}
	if omittedDigest != frozenInputDigest {
		t.Fatalf("empty-set digest=%q, want frozen %q", omittedDigest, frozenInputDigest)
	}
}

func TestCanonicalContextRefAbsentAndPresent(t *testing.T) {
	t.Parallel()

	absent := frozenVectorRequest()
	absentBytes, err := canonicalPolicyEvaluationRequestBytes(absent)
	if err != nil {
		t.Fatalf("absent contextRef: %v", err)
	}
	if bytes.Contains(absentBytes, []byte(`"contextRef"`)) {
		t.Fatalf("omitted contextRef must be absent from normalized object: %s", absentBytes)
	}

	present := frozenVectorRequest()
	present.ContextRef = &apimeta.TypedRef{
		APIVersion: "governance.sovrunn.io/v1alpha1",
		Kind:       "EffectiveGovernanceContext",
		Name:       "default",
		UID:        "egc-001",
	}
	presentBytes, err := canonicalPolicyEvaluationRequestBytes(present)
	if err != nil {
		t.Fatalf("present contextRef: %v", err)
	}
	wantFragment := `"contextRef":{"apiVersion":"governance.sovrunn.io/v1alpha1","kind":"EffectiveGovernanceContext","name":"default","uid":"egc-001"}`
	if !bytes.Contains(presentBytes, []byte(wantFragment)) {
		t.Fatalf("present contextRef missing/wrong\ngot: %s\nwant fragment: %s", presentBytes, wantFragment)
	}
	// JCS root order: action, candidateRefs, contextRef, profileRefs, schema, subjectRef
	actionIdx := bytes.Index(presentBytes, []byte(`"action"`))
	candIdx := bytes.Index(presentBytes, []byte(`"candidateRefs"`))
	ctxIdx := bytes.Index(presentBytes, []byte(`"contextRef"`))
	profIdx := bytes.Index(presentBytes, []byte(`"profileRefs"`))
	schemaIdx := bytes.Index(presentBytes, []byte(`"schema"`))
	subjIdx := bytes.Index(presentBytes, []byte(`"subjectRef"`))
	if actionIdx >= candIdx || candIdx >= ctxIdx || ctxIdx >= profIdx || profIdx >= schemaIdx || schemaIdx >= subjIdx {
		t.Fatalf("root key JCS order incorrect: %s", presentBytes)
	}

	absentDigest, err := digestPolicyEvaluationRequest(absent)
	if err != nil {
		t.Fatalf("absent digest: %v", err)
	}
	presentDigest, err := digestPolicyEvaluationRequest(present)
	if err != nil {
		t.Fatalf("present digest: %v", err)
	}
	if absentDigest == presentDigest {
		t.Fatal("present contextRef must alter digest")
	}
}

func TestCanonicalListOrderIndependentDigest(t *testing.T) {
	t.Parallel()

	base := frozenVectorRequest()
	base.ProfileRefs = []apimeta.TypedRef{
		{APIVersion: "iam.sovrunn.io/v1alpha1", Kind: "RoleDefinition", Name: "reader", UID: "role-001"},
		{APIVersion: "iam.sovrunn.io/v1alpha1", Kind: "RoleDefinition", Name: "writer", UID: "role-002"},
	}
	base.CandidateRefs = []apimeta.TypedRef{
		{APIVersion: "services.sovrunn.io/v1alpha1", Kind: "ServiceInstance", Name: "alpha", UID: "svc-a"},
		{APIVersion: "services.sovrunn.io/v1alpha1", Kind: "ServiceInstance", Name: "beta", UID: "svc-b"},
	}

	reordered := frozenVectorRequest()
	reordered.ProfileRefs = []apimeta.TypedRef{
		{APIVersion: "iam.sovrunn.io/v1alpha1", Kind: "RoleDefinition", Name: "writer", UID: "role-002"},
		{APIVersion: "iam.sovrunn.io/v1alpha1", Kind: "RoleDefinition", Name: "reader", UID: "role-001"},
	}
	reordered.CandidateRefs = []apimeta.TypedRef{
		{APIVersion: "services.sovrunn.io/v1alpha1", Kind: "ServiceInstance", Name: "beta", UID: "svc-b"},
		{APIVersion: "services.sovrunn.io/v1alpha1", Kind: "ServiceInstance", Name: "alpha", UID: "svc-a"},
	}
	reordered.RequestID = "different-correlation"

	baseDigest, err := digestPolicyEvaluationRequest(base)
	if err != nil {
		t.Fatalf("base digest: %v", err)
	}
	reorderedDigest, err := digestPolicyEvaluationRequest(reordered)
	if err != nil {
		t.Fatalf("reordered digest: %v", err)
	}
	if baseDigest != reorderedDigest {
		t.Fatalf("reordered lists must preserve digest: base=%q reordered=%q", baseDigest, reorderedDigest)
	}

	baseBytes, err := canonicalPolicyEvaluationRequestBytes(base)
	if err != nil {
		t.Fatalf("base bytes: %v", err)
	}
	// Sorted order: alpha before beta; reader before writer.
	readerIdx := bytes.Index(baseBytes, []byte(`"name":"reader"`))
	writerIdx := bytes.Index(baseBytes, []byte(`"name":"writer"`))
	alphaIdx := bytes.Index(baseBytes, []byte(`"name":"alpha"`))
	betaIdx := bytes.Index(baseBytes, []byte(`"name":"beta"`))
	if readerIdx < 0 || writerIdx < 0 || alphaIdx < 0 || betaIdx < 0 {
		t.Fatalf("missing sorted names in %s", baseBytes)
	}
	if readerIdx >= writerIdx || alphaIdx >= betaIdx {
		t.Fatalf("refs not sorted: %s", baseBytes)
	}
}

func TestCanonicalSemanticChangeAltersDigest(t *testing.T) {
	t.Parallel()

	baseDigest, err := digestPolicyEvaluationRequest(frozenVectorRequest())
	if err != nil {
		t.Fatalf("base: %v", err)
	}

	changedAction := frozenVectorRequest()
	changedAction.Action = "service.write"
	actionDigest, err := digestPolicyEvaluationRequest(changedAction)
	if err != nil {
		t.Fatalf("action: %v", err)
	}
	if actionDigest == baseDigest {
		t.Fatal("action change must alter digest")
	}

	changedSubject := frozenVectorRequest()
	changedSubject.SubjectRef.Name = "other-api"
	subjectDigest, err := digestPolicyEvaluationRequest(changedSubject)
	if err != nil {
		t.Fatalf("subject: %v", err)
	}
	if subjectDigest == baseDigest {
		t.Fatal("subjectRef change must alter digest")
	}

	changedProfile := frozenVectorRequest()
	changedProfile.ProfileRefs[0].UID = "role-999"
	profileDigest, err := digestPolicyEvaluationRequest(changedProfile)
	if err != nil {
		t.Fatalf("profile: %v", err)
	}
	if profileDigest == baseDigest {
		t.Fatal("profileRefs change must alter digest")
	}

	withCandidate := frozenVectorRequest()
	withCandidate.CandidateRefs = []apimeta.TypedRef{{
		APIVersion: "services.sovrunn.io/v1alpha1",
		Kind:       "ServiceInstance",
		Name:       "peer",
		UID:        "svc-peer",
	}}
	candDigest, err := digestPolicyEvaluationRequest(withCandidate)
	if err != nil {
		t.Fatalf("candidate: %v", err)
	}
	if candDigest == baseDigest {
		t.Fatal("candidateRefs change must alter digest")
	}
}

func TestCanonicalRequestIDOnlyChangePreservesDigest(t *testing.T) {
	t.Parallel()

	a := frozenVectorRequest()
	a.RequestID = "id-a"
	b := frozenVectorRequest()
	b.RequestID = "id-b-completely-different"

	da, err := digestPolicyEvaluationRequest(a)
	if err != nil {
		t.Fatalf("a: %v", err)
	}
	db, err := digestPolicyEvaluationRequest(b)
	if err != nil {
		t.Fatalf("b: %v", err)
	}
	if da != db {
		t.Fatalf("requestId-only change must preserve digest: %q vs %q", da, db)
	}
	if da != frozenInputDigest {
		t.Fatalf("digest=%q, want frozen %q", da, frozenInputDigest)
	}
}

func TestCanonicalSnapshotOwnership(t *testing.T) {
	t.Parallel()

	req := frozenVectorRequest()
	req.ProfileRefs = []apimeta.TypedRef{
		{APIVersion: "iam.sovrunn.io/v1alpha1", Kind: "RoleDefinition", Name: "z", UID: "u-z"},
		{APIVersion: "iam.sovrunn.io/v1alpha1", Kind: "RoleDefinition", Name: "a", UID: "u-a"},
	}
	req.ContextRef = &apimeta.TypedRef{
		APIVersion: "governance.sovrunn.io/v1alpha1",
		Kind:       "EffectiveGovernanceContext",
		Name:       "ctx",
		UID:        "ctx-1",
	}
	originalProfileOrder := []string{req.ProfileRefs[0].Name, req.ProfileRefs[1].Name}

	before, err := digestPolicyEvaluationRequest(req)
	if err != nil {
		t.Fatalf("before: %v", err)
	}
	if req.ProfileRefs[0].Name != originalProfileOrder[0] || req.ProfileRefs[1].Name != originalProfileOrder[1] {
		t.Fatalf("canonicalization must not reorder caller-owned profileRefs: got %q,%q",
			req.ProfileRefs[0].Name, req.ProfileRefs[1].Name)
	}

	// Mutate caller-owned storage after digest; previously computed digest is fixed.
	req.Action = "mutated"
	req.ProfileRefs[0].Name = "mutated"
	req.ContextRef.Name = "mutated"
	req.RequestID = "mutated"

	sameShape := frozenVectorRequest()
	sameShape.ProfileRefs = []apimeta.TypedRef{
		{APIVersion: "iam.sovrunn.io/v1alpha1", Kind: "RoleDefinition", Name: "z", UID: "u-z"},
		{APIVersion: "iam.sovrunn.io/v1alpha1", Kind: "RoleDefinition", Name: "a", UID: "u-a"},
	}
	sameShape.ContextRef = &apimeta.TypedRef{
		APIVersion: "governance.sovrunn.io/v1alpha1",
		Kind:       "EffectiveGovernanceContext",
		Name:       "ctx",
		UID:        "ctx-1",
	}
	sameDigest, err := digestPolicyEvaluationRequest(sameShape)
	if err != nil {
		t.Fatalf("sameShape: %v", err)
	}
	if before != sameDigest {
		t.Fatalf("digest must depend on snapshot values, not later caller mutation: before=%q same=%q", before, sameDigest)
	}

	reordered := frozenVectorRequest()
	reordered.ProfileRefs = []apimeta.TypedRef{
		{APIVersion: "iam.sovrunn.io/v1alpha1", Kind: "RoleDefinition", Name: "a", UID: "u-a"},
		{APIVersion: "iam.sovrunn.io/v1alpha1", Kind: "RoleDefinition", Name: "z", UID: "u-z"},
	}
	reordered.ContextRef = &apimeta.TypedRef{
		APIVersion: "governance.sovrunn.io/v1alpha1",
		Kind:       "EffectiveGovernanceContext",
		Name:       "ctx",
		UID:        "ctx-1",
	}
	reorderedDigest, err := digestPolicyEvaluationRequest(reordered)
	if err != nil {
		t.Fatalf("reordered: %v", err)
	}
	if reorderedDigest != before {
		t.Fatalf("order-independent digest mismatch: %q vs %q", reorderedDigest, before)
	}
}

func TestCanonicalOmitEmptyUID(t *testing.T) {
	t.Parallel()

	req := frozenVectorRequest()
	req.SubjectRef.UID = ""
	req.ProfileRefs = []apimeta.TypedRef{{
		APIVersion: "iam.sovrunn.io/v1alpha1",
		Kind:       "RoleDefinition",
		Name:       "reader",
		UID:        "role-001", // still present on profile (required semantically later)
	}}
	req.CandidateRefs = []apimeta.TypedRef{{
		APIVersion: "services.sovrunn.io/v1alpha1",
		Kind:       "ServiceInstance",
		Name:       "peer",
		UID:        "",
	}}

	got, err := canonicalPolicyEvaluationRequestBytes(req)
	if err != nil {
		t.Fatalf("canonical: %v", err)
	}
	wantSubject := `"subjectRef":{"apiVersion":"services.sovrunn.io/v1alpha1","kind":"ServiceInstance","name":"reporting-api"}`
	if !bytes.Contains(got, []byte(wantSubject)) {
		t.Fatalf("subjectRef without uid incorrect: %s", got)
	}
	wantCandidate := `"candidateRefs":[{"apiVersion":"services.sovrunn.io/v1alpha1","kind":"ServiceInstance","name":"peer"}]`
	if !bytes.Contains(got, []byte(wantCandidate)) {
		t.Fatalf("candidateRef without uid incorrect: %s", got)
	}
	if !bytes.Contains(got, []byte(`"uid":"role-001"`)) {
		t.Fatalf("non-empty profile uid must remain: %s", got)
	}
}

func TestCanonicalJCSStringEmitter(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name       string
		action     string
		wantSubstr string
	}{
		{
			name:       "quote_and_reverse_solidus",
			action:     `a"b\c`,
			wantSubstr: `"action":"a\"b\\c"`,
		},
		{
			name:       "short_escapes",
			action:     "a\b\t\n\f\rb",
			wantSubstr: `"action":"a\b\t\n\f\rb"`,
		},
		{
			name:       "other_controls_lowercase_u",
			action:     "a\u0000\u0001\u001fb",
			wantSubstr: `"action":"a\u0000\u0001\u001fb"`,
		},
		{
			name:       "html_sensitive_unchanged",
			action:     `<>&/`,
			wantSubstr: `"action":"<>&/"`,
		},
		{
			name:       "bmp_non_ascii",
			action:     "café",
			wantSubstr: `"action":"café"`,
		},
		{
			name:       "supplementary_emoji",
			action:     "grin😀",
			wantSubstr: `"action":"grin😀"`,
		},
		{
			name:       "line_separator_u2028",
			action:     "a\u2028b",
			wantSubstr: "\"action\":\"a\u2028b\"",
		},
		{
			name:       "paragraph_separator_u2029",
			action:     "a\u2029b",
			wantSubstr: "\"action\":\"a\u2029b\"",
		},
		{
			name:       "solidus_unescaped",
			action:     "a/b",
			wantSubstr: `"action":"a/b"`,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			req := frozenVectorRequest()
			req.Action = tc.action
			got, err := canonicalPolicyEvaluationRequestBytes(req)
			if err != nil {
				t.Fatalf("canonical: %v", err)
			}
			if !bytes.Contains(got, []byte(tc.wantSubstr)) {
				t.Fatalf("missing expected emission\ngot:  %s\nwant: %s", got, tc.wantSubstr)
			}
			if bytes.Contains(got, []byte("\\u2028")) || bytes.Contains(got, []byte("\\u2029")) {
				t.Fatalf("U+2028/U+2029 must be emitted as UTF-8, not \\u escapes: %s", got)
			}
		})
	}
}

func TestCanonicalRejectInvalidUTF8(t *testing.T) {
	t.Parallel()

	req := frozenVectorRequest()
	req.Action = string([]byte{'a', 0xff, 'b'})
	if utf8.ValidString(req.Action) {
		t.Fatal("test setup: action unexpectedly valid UTF-8")
	}
	if _, err := canonicalPolicyEvaluationRequestBytes(req); err == nil {
		t.Fatal("expected invalid UTF-8 rejection")
	} else if err != errRequestInvalid {
		t.Fatalf("err=%v, want errRequestInvalid", err)
	}
	if _, err := digestPolicyEvaluationRequest(req); err == nil {
		t.Fatal("digest must fail for invalid UTF-8")
	}
}

func TestCanonicalRejectLoneSurrogateUTF8(t *testing.T) {
	t.Parallel()

	// UTF-8 encoding of U+D800 (lone high surrogate) — invalid / non-scalar.
	req := frozenVectorRequest()
	req.Action = string([]byte{0xed, 0xa0, 0x80})
	if utf8.ValidString(req.Action) {
		t.Fatal("test setup: lone surrogate unexpectedly valid")
	}
	if _, err := canonicalPolicyEvaluationRequestBytes(req); err == nil {
		t.Fatal("expected lone surrogate rejection")
	} else if err != errRequestInvalid {
		t.Fatalf("err=%v, want errRequestInvalid", err)
	}
}

func TestCanonicalNoUnicodeNormalization(t *testing.T) {
	t.Parallel()

	composed := frozenVectorRequest()
	composed.Action = "é" // U+00E9

	decomposed := frozenVectorRequest()
	decomposed.Action = "e\u0301" // e + combining acute

	if composed.Action == decomposed.Action {
		t.Fatal("test setup: composed and decomposed must differ as Go strings")
	}

	cBytes, err := canonicalPolicyEvaluationRequestBytes(composed)
	if err != nil {
		t.Fatalf("composed: %v", err)
	}
	dBytes, err := canonicalPolicyEvaluationRequestBytes(decomposed)
	if err != nil {
		t.Fatalf("decomposed: %v", err)
	}
	if bytes.Equal(cBytes, dBytes) {
		t.Fatal("composed vs decomposed must not normalize to the same JCS bytes")
	}

	cDigest, err := digestPolicyEvaluationRequest(composed)
	if err != nil {
		t.Fatalf("composed digest: %v", err)
	}
	dDigest, err := digestPolicyEvaluationRequest(decomposed)
	if err != nil {
		t.Fatalf("decomposed digest: %v", err)
	}
	if cDigest == dDigest {
		t.Fatal("composed vs decomposed must produce different digests")
	}
}

func TestCanonicalAppendJCSStringDirect(t *testing.T) {
	t.Parallel()

	got, err := appendJCSString(nil, "x\"\\y")
	if err != nil {
		t.Fatalf("appendJCSString: %v", err)
	}
	if string(got) != `"x\"\\y"` {
		t.Fatalf("got %s", got)
	}

	if _, err := appendJCSString(nil, string([]byte{0xff})); err != errRequestInvalid {
		t.Fatalf("invalid UTF-8 err=%v, want errRequestInvalid", err)
	}
}
