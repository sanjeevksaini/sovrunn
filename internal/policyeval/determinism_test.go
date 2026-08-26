package policyeval

import (
	"context"
	"encoding/json"
	"reflect"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/decision"
)

// TASK-F17-10 / CDG-F17-06 / design §5.3 — consolidated determinism,
// concurrency, bounds, reason-code sort, and copy-ownership proofs.

func TestDeterminism_FixedTimeSourceReplay(t *testing.T) {
	t.Parallel()

	fixed := time.Date(2026, 8, 26, 6, 30, 0, 123456789, time.UTC)
	wantAt := fixed.Format(time.RFC3339Nano)

	fake := mustFake(t, []FakeFixture{{
		InputDigest: frozenInputDigest,
		Conclusion: &Conclusion{
			Outcome:     OutcomeRequiresApproval,
			ReasonCodes: []string{"NEED_APPROVAL", "POLICY_CHECK"},
		},
	}})
	b := mustBoundary(t, fake, NewFixedTimeSource(fixed))
	req := frozenVectorRequest()

	r1, t1, nr1, err1 := b.Evaluate(context.Background(), req)
	if err1 != nil || nr1 != "" {
		t.Fatalf("first Evaluate: nr=%q err=%v", nr1, err1)
	}
	r2, t2, nr2, err2 := b.Evaluate(context.Background(), req)
	if err2 != nil || nr2 != "" {
		t.Fatalf("second Evaluate: nr=%q err=%v", nr2, err2)
	}

	if !reflect.DeepEqual(r1, r2) {
		t.Fatalf("results differ across replay\nfirst=%+v\nsecond=%+v", r1, r2)
	}
	if !reflect.DeepEqual(t1, t2) {
		t.Fatalf("timings differ across replay\nfirst=%+v\nsecond=%+v", t1, t2)
	}
	if r1.EvaluatedAt != wantAt || r1.EvaluatedAt != t1.CompletedAt {
		t.Fatalf("evaluatedAt=%q completedAt=%q want %q", r1.EvaluatedAt, t1.CompletedAt, wantAt)
	}
	if t1.StartedAt != wantAt || t1.CompletedAt != wantAt {
		t.Fatalf("timing=%+v, want started/completed %q", t1, wantAt)
	}
	if r1.Outcome != OutcomeRequiresApproval {
		t.Fatalf("outcome=%q, want RequiresApproval", r1.Outcome)
	}
	if r1.InputDigest != frozenInputDigest {
		t.Fatalf("inputDigest=%q, want frozen %q", r1.InputDigest, frozenInputDigest)
	}
	wantCodes := []string{"NEED_APPROVAL", "POLICY_CHECK"}
	if !reflect.DeepEqual(r1.ReasonCodes, wantCodes) {
		t.Fatalf("reasonCodes=%v, want %v", r1.ReasonCodes, wantCodes)
	}
	if r1.Obligations != nil {
		t.Fatalf("Obligations must remain nil, got %v", r1.Obligations)
	}
}

func TestDeterminism_ConcurrentEquivalentCalls(t *testing.T) {
	t.Parallel()

	fixed := time.Date(2026, 8, 26, 7, 0, 0, 0, time.UTC)
	fake := mustFake(t, []FakeFixture{{
		InputDigest: frozenInputDigest,
		Conclusion: &Conclusion{
			Outcome:     OutcomeDeny,
			ReasonCodes: []string{"ZULU", "ALPHA", "MIKE"},
		},
	}})
	b := mustBoundary(t, fake, NewFixedTimeSource(fixed))
	req := frozenVectorRequest()

	baseResult, baseTiming, nr, err := b.Evaluate(context.Background(), req)
	if err != nil || nr != "" {
		t.Fatalf("baseline Evaluate: nr=%q err=%v", nr, err)
	}
	wantCodes := []string{"ALPHA", "MIKE", "ZULU"}
	if !reflect.DeepEqual(baseResult.ReasonCodes, wantCodes) {
		t.Fatalf("baseline reasonCodes=%v, want %v", baseResult.ReasonCodes, wantCodes)
	}

	const goroutines = 64
	results := make([]PolicyEvaluationResult, goroutines)
	timings := make([]decision.TimingEnvelope, goroutines)
	errs := make([]error, goroutines)
	nrs := make([]NonResult, goroutines)

	var wg sync.WaitGroup
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			// Equivalent call: same request value, shared boundary/fake/time.
			r, tm, gotNR, gotErr := b.Evaluate(context.Background(), req)
			results[i] = r
			timings[i] = tm
			nrs[i] = gotNR
			errs[i] = gotErr
		}(i)
	}
	wg.Wait()

	for i := 0; i < goroutines; i++ {
		if errs[i] != nil || nrs[i] != "" {
			t.Fatalf("goroutine %d: nr=%q err=%v", i, nrs[i], errs[i])
		}
		if !reflect.DeepEqual(results[i], baseResult) {
			t.Fatalf("goroutine %d result differs\ngot=%+v\nwant=%+v", i, results[i], baseResult)
		}
		if !reflect.DeepEqual(timings[i], baseTiming) {
			t.Fatalf("goroutine %d timing differs\ngot=%+v\nwant=%+v", i, timings[i], baseTiming)
		}
	}
}

func TestDeterminism_CopyOwnershipCallerFixtureAdapterOutput(t *testing.T) {
	t.Parallel()

	fixed := time.Date(2026, 8, 26, 8, 0, 0, 0, time.UTC)
	wantAt := fixed.Format(time.RFC3339Nano)

	t.Run("D10_caller_and_adapter_return", func(t *testing.T) {
		t.Parallel()

		adapterCodes := []string{"ZULU", "ALPHA", "MIKE"}
		spy := &evaluateSpyAdapter{
			conclusion: Conclusion{
				Outcome:     OutcomeDeny,
				ReasonCodes: adapterCodes,
			},
		}
		b := mustBoundary(t, spy, NewFixedTimeSource(fixed))

		req := frozenVectorRequest()
		req.ProfileRefs = append([]apimeta.TypedRef(nil), req.ProfileRefs...)
		ctxRef := apimeta.TypedRef{
			APIVersion: "governance.sovrunn.io/v1alpha1",
			Kind:       "EffectiveGovernanceContext",
			Name:       "egc-1",
			UID:        "egc-uid-1",
		}
		req.ContextRef = &ctxRef
		req.CandidateRefs = []apimeta.TypedRef{{
			APIVersion: "topology.sovrunn.io/v1alpha1",
			Kind:       "CloudProvider",
			Name:       "provider-a",
			UID:        "cp-a",
		}}

		result, timing, nr, err := b.Evaluate(context.Background(), req)
		if err != nil || nr != "" {
			t.Fatalf("Evaluate: nr=%q err=%v", nr, err)
		}
		wantCodes := []string{"ALPHA", "MIKE", "ZULU"}
		if !reflect.DeepEqual(result.ReasonCodes, wantCodes) {
			t.Fatalf("reasonCodes=%v, want %v", result.ReasonCodes, wantCodes)
		}
		if result.EvaluatedAt != wantAt || timing.CompletedAt != wantAt {
			t.Fatalf("evaluatedAt/completedAt=%q/%q, want %q", result.EvaluatedAt, timing.CompletedAt, wantAt)
		}
		ownedDigest := result.InputDigest
		ownedCodes := append([]string(nil), result.ReasonCodes...)

		// Mutate caller-owned request after handoff (D-10).
		req.Action = "mutated.action"
		req.RequestID = "mutated-request-id"
		req.ProfileRefs[0].Name = "mutated-profile"
		req.CandidateRefs[0].Name = "mutated-candidate"
		ctxRef.Name = "mutated-context"
		req.ContextRef.UID = "mutated-uid"

		// Mutate adapter-owned conclusion after return (D-10).
		adapterCodes[0] = "MUTATED_ADAPTER"
		spy.conclusion.Outcome = OutcomeAllow
		spy.conclusion.ReasonCodes[1] = "MUTATED_SPY"
		spy.conclusion.ReasonCodes = append(spy.conclusion.ReasonCodes, "APPENDED")

		if result.Outcome != OutcomeDeny {
			t.Fatalf("result outcome mutated via adapter: %q", result.Outcome)
		}
		if !reflect.DeepEqual(result.ReasonCodes, ownedCodes) {
			t.Fatalf("result reasonCodes mutated via adapter: %v", result.ReasonCodes)
		}
		if result.InputDigest != ownedDigest {
			t.Fatalf("result digest mutated: %q", result.InputDigest)
		}
		if adapterCodes[0] == "ZULU" {
			t.Fatal("test setup: adapterCodes mutation did not take effect")
		}
		if spy.conclusion.ReasonCodes[0] != "MUTATED_ADAPTER" {
			t.Fatalf("adapter conclusion not mutated as expected: %v", spy.conclusion.ReasonCodes)
		}

		// Mutate boundary output after return; a later call must be independent.
		result.ReasonCodes[0] = "MUTATED_OUTPUT"
		result.Outcome = OutcomeAllow
		result.InputDigest = "deadbeef"
		result.EvaluatedAt = "mutated-time"
		timing.StartedAt = "mutated-started"
		timing.CompletedAt = "mutated-completed"

		// Restore a valid adapter conclusion so the second call proves output
		// mutation did not poison boundary state (adapter config is separate).
		spy.mu.Lock()
		spy.conclusion = Conclusion{
			Outcome:     OutcomeDeny,
			ReasonCodes: []string{"ZULU", "ALPHA", "MIKE"},
		}
		spy.mu.Unlock()

		r2, t2, nr2, err2 := b.Evaluate(context.Background(), frozenVectorRequest())
		if err2 != nil || nr2 != "" {
			t.Fatalf("second Evaluate: nr=%q err=%v", nr2, err2)
		}
		if !reflect.DeepEqual(r2.ReasonCodes, wantCodes) {
			t.Fatalf("second result reasonCodes=%v, want %v", r2.ReasonCodes, wantCodes)
		}
		if r2.EvaluatedAt != wantAt || t2.CompletedAt != wantAt {
			t.Fatalf("second evaluatedAt/completedAt=%q/%q, want %q", r2.EvaluatedAt, t2.CompletedAt, wantAt)
		}
		if r2.InputDigest != frozenInputDigest {
			t.Fatalf("second digest=%q, want %q", r2.InputDigest, frozenInputDigest)
		}
	})

	t.Run("D12_fixture_and_fake_return", func(t *testing.T) {
		t.Parallel()

		fixtureCodes := []string{"POLICY_ALLOW", "POLICY_EXTRA"}
		conclusion := &Conclusion{Outcome: OutcomeAllow, ReasonCodes: fixtureCodes}
		fixtures := []FakeFixture{{
			InputDigest: frozenInputDigest,
			Conclusion:  conclusion,
		}}
		fake := mustFake(t, fixtures)
		b := mustBoundary(t, fake, NewFixedTimeSource(fixed))

		result, _, nr, err := b.Evaluate(context.Background(), frozenVectorRequest())
		if err != nil || nr != "" {
			t.Fatalf("Evaluate: nr=%q err=%v", nr, err)
		}
		wantCodes := []string{"POLICY_ALLOW", "POLICY_EXTRA"}
		if !reflect.DeepEqual(result.ReasonCodes, wantCodes) {
			t.Fatalf("reasonCodes=%v, want %v", result.ReasonCodes, wantCodes)
		}

		// Mutate caller-owned fixture material after construction (D-12).
		fixtureCodes[0] = "MUTATED_FIXTURE"
		conclusion.Outcome = OutcomeDeny
		conclusion.ReasonCodes = []string{"MUTATED_AFTER"}
		fixtures[0].InputDigest = fakeDigestDeny
		fixtures[0].AdapterFailure = true
		fixtures[0].Conclusion = nil

		// Mutate returned result / fake Evaluate conclusion after handoff.
		result.ReasonCodes[0] = "MUTATED_RESULT"
		direct, derr := fake.Evaluate(context.Background(), NormalizedInput{}, frozenInputDigest)
		if derr != nil {
			t.Fatalf("fake.Evaluate: %v", derr)
		}
		direct.Outcome = OutcomeIndeterminate
		direct.ReasonCodes[0] = "MUTATED_DIRECT"
		direct.ReasonCodes = append(direct.ReasonCodes, "APPENDED")

		r2, _, nr2, err2 := b.Evaluate(context.Background(), frozenVectorRequest())
		if err2 != nil || nr2 != "" {
			t.Fatalf("second Evaluate after fixture mutation: nr=%q err=%v", nr2, err2)
		}
		if r2.Outcome != OutcomeAllow {
			t.Fatalf("outcome=%q, want Allow (fixture mutation leaked)", r2.Outcome)
		}
		if !reflect.DeepEqual(r2.ReasonCodes, wantCodes) {
			t.Fatalf("reasonCodes=%v, want %v (fixture/output mutation leaked)", r2.ReasonCodes, wantCodes)
		}
	})

	t.Run("D14_mapper_output", func(t *testing.T) {
		t.Parallel()

		fake := mustFake(t, []FakeFixture{{
			InputDigest: frozenInputDigest,
			Conclusion:  fakeConclusion(OutcomeAllow, "POLICY_ALLOW", "POLICY_EXTRA"),
		}})
		b := mustBoundary(t, fake, NewFixedTimeSource(fixed))

		result, timing, nr, err := b.Evaluate(context.Background(), frozenVectorRequest())
		if err != nil || nr != "" {
			t.Fatalf("Evaluate: nr=%q err=%v", nr, err)
		}

		ident := validMapperSnapshotIdent()
		integrity := &decision.OpaqueIntegrityCarrier{State: "bound"}
		ident.InputIntegrity = integrity

		mapped, err := MapToEvaluationResult(result, timing, validMapperEvaluator(), ident)
		if err != nil {
			t.Fatalf("MapToEvaluationResult: %v", err)
		}
		ownedRaw := append(json.RawMessage(nil), mapped.Result...)

		// Mutate caller-owned Evaluate/mapper inputs after successful mapping.
		result.ReasonCodes[0] = "MUTATED"
		result.Outcome = OutcomeDeny
		timing.StartedAt = "mutated"
		timing.CompletedAt = "mutated"
		ident.InputSnapshotRef.Name = "mutated-name"
		integrity.State = "mutated"

		if mapped.InputSnapshotRef == nil || mapped.InputSnapshotRef.Name != "placement-input-1" {
			t.Fatalf("snapshot aliased or mutated: %+v", mapped.InputSnapshotRef)
		}
		if mapped.InputIntegrity == nil || mapped.InputIntegrity.State != "bound" {
			t.Fatalf("integrity aliased or mutated: %+v", mapped.InputIntegrity)
		}
		if mapped.EvaluatedAt != wantAt || mapped.Timing.CompletedAt != wantAt {
			t.Fatalf("mapped times mutated: evaluatedAt=%q completedAt=%q", mapped.EvaluatedAt, mapped.Timing.CompletedAt)
		}

		var embedded PolicyEvaluationResult
		if err := json.Unmarshal(ownedRaw, &embedded); err != nil {
			t.Fatalf("unmarshal owned raw: %v", err)
		}
		if embedded.ReasonCodes[0] != "POLICY_ALLOW" || embedded.Outcome != OutcomeAllow {
			t.Fatalf("embedded result mutated via caller: %+v", embedded)
		}

		// Mutating mapper output must not rewrite the retained raw copy.
		mapped.Result[0] = 'X'
		if string(mapped.Result) == string(ownedRaw) {
			t.Fatal("mapped.Result must be distinct from retained raw copy after mutation")
		}
	})
}

func TestDeterminism_ReasonCodeSortStableAndTotal(t *testing.T) {
	t.Parallel()

	fixed := NewFixedTimeSource(time.Date(2026, 8, 26, 9, 0, 0, 0, time.UTC))
	perms := [][]string{
		{"ZULU", "ALPHA", "MIKE", "BRAVO"},
		{"MIKE", "BRAVO", "ZULU", "ALPHA"},
		{"ALPHA", "BRAVO", "MIKE", "ZULU"},
		{"BRAVO", "ZULU", "ALPHA", "MIKE"},
	}
	want := append([]string(nil), perms[0]...)
	sort.Strings(want)

	var baseline PolicyEvaluationResult
	for i, codes := range perms {
		codes := append([]string(nil), codes...)
		spy := &evaluateSpyAdapter{
			conclusion: Conclusion{Outcome: OutcomeIndeterminate, ReasonCodes: codes},
		}
		b := mustBoundary(t, spy, fixed)
		result, _, nr, err := b.Evaluate(context.Background(), frozenVectorRequest())
		if err != nil || nr != "" {
			t.Fatalf("perm %d: nr=%q err=%v", i, nr, err)
		}
		if !reflect.DeepEqual(result.ReasonCodes, want) {
			t.Fatalf("perm %d: reasonCodes=%v, want %v", i, result.ReasonCodes, want)
		}
		// Adapter-owned storage must remain unsorted (no in-place sort).
		if !reflect.DeepEqual(spy.conclusion.ReasonCodes, codes) {
			t.Fatalf("perm %d: adapter reasonCodes mutated: got %v want %v", i, spy.conclusion.ReasonCodes, codes)
		}
		if i == 0 {
			baseline = result
			continue
		}
		if !reflect.DeepEqual(result.ReasonCodes, baseline.ReasonCodes) {
			t.Fatalf("sort not stable across permutations: %v vs %v", result.ReasonCodes, baseline.ReasonCodes)
		}
	}

	// Totality: every pairwise order matches bytewise string ordering.
	for i := 0; i < len(want)-1; i++ {
		if want[i] >= want[i+1] {
			t.Fatalf("sorted sequence not strictly increasing at %d: %v", i, want)
		}
		if strings.Compare(want[i], want[i+1]) >= 0 {
			t.Fatalf("string compare not total/strict at %d: %v", i, want)
		}
	}
}

func TestDeterminism_BoundsAcceptedThroughEvaluate(t *testing.T) {
	t.Parallel()

	fixed := NewFixedTimeSource(time.Date(2026, 8, 26, 10, 0, 0, 0, time.UTC))
	spy := &evaluateSpyAdapter{conclusion: validAllowConclusion()}
	b := mustBoundary(t, spy, fixed)

	cases := []struct {
		name string
		raw  string
	}{
		{name: "action_min_1", raw: validRequestJSON(withAction(strings.Repeat("a", 1)))},
		{name: "action_max_63", raw: validRequestJSON(withAction(strings.Repeat("a", 63)))},
		{name: "profileRefs_min_1", raw: validRequestJSON(withProfileCount(1))},
		{name: "profileRefs_max_32", raw: validRequestJSON(withProfileCount(32))},
		{name: "candidateRefs_min_0_omitted", raw: validRequestJSON(omitCandidates())},
		{name: "candidateRefs_min_0_empty_array", raw: validRequestJSON(withCandidateCount(0))},
		{name: "candidateRefs_max_64", raw: validRequestJSON(withCandidateCount(64))},
		{name: "requestId_min_1", raw: validRequestJSON(withRequestID(strings.Repeat("r", 1)))},
		{name: "requestId_max_128", raw: validRequestJSON(withRequestID(strings.Repeat("r", 128)))},
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

			result, timing, nr, evalErr := b.Evaluate(context.Background(), req)
			if evalErr != nil || nr != "" {
				t.Fatalf("Evaluate: nr=%q err=%v", nr, evalErr)
			}
			if result.Outcome != OutcomeAllow {
				t.Fatalf("outcome=%q, want Allow", result.Outcome)
			}
			if result.EvaluatedAt == "" || result.EvaluatedAt != timing.CompletedAt {
				t.Fatalf("evaluatedAt=%q completedAt=%q", result.EvaluatedAt, timing.CompletedAt)
			}
			if result.InputDigest == "" || len(result.InputDigest) != 64 {
				t.Fatalf("inputDigest=%q, want 64-char digest", result.InputDigest)
			}
			if result.Obligations != nil {
				t.Fatalf("Obligations must remain nil, got %v", result.Obligations)
			}
		})
	}
}

func TestDeterminism_DigestPropertiesUnderEvaluate(t *testing.T) {
	t.Parallel()

	fixed := NewFixedTimeSource(time.Date(2026, 8, 26, 11, 0, 0, 0, time.UTC))
	spy := &evaluateSpyAdapter{conclusion: validAllowConclusion()}
	b := mustBoundary(t, spy, fixed)

	base := frozenVectorRequest()
	base.ProfileRefs = []apimeta.TypedRef{
		{APIVersion: "iam.sovrunn.io/v1alpha1", Kind: "RoleDefinition", Name: "writer", UID: "role-002"},
		{APIVersion: "iam.sovrunn.io/v1alpha1", Kind: "RoleDefinition", Name: "reader", UID: "role-001"},
	}
	base.CandidateRefs = []apimeta.TypedRef{
		{APIVersion: "topology.sovrunn.io/v1alpha1", Kind: "CloudProvider", Name: "b", UID: "cp-b"},
		{APIVersion: "topology.sovrunn.io/v1alpha1", Kind: "CloudProvider", Name: "a", UID: "cp-a"},
	}

	reordered := frozenVectorRequest()
	reordered.ProfileRefs = []apimeta.TypedRef{
		{APIVersion: "iam.sovrunn.io/v1alpha1", Kind: "RoleDefinition", Name: "reader", UID: "role-001"},
		{APIVersion: "iam.sovrunn.io/v1alpha1", Kind: "RoleDefinition", Name: "writer", UID: "role-002"},
	}
	reordered.CandidateRefs = []apimeta.TypedRef{
		{APIVersion: "topology.sovrunn.io/v1alpha1", Kind: "CloudProvider", Name: "a", UID: "cp-a"},
		{APIVersion: "topology.sovrunn.io/v1alpha1", Kind: "CloudProvider", Name: "b", UID: "cp-b"},
	}

	r1, _, nr1, err1 := b.Evaluate(context.Background(), base)
	if err1 != nil || nr1 != "" {
		t.Fatalf("base Evaluate: nr=%q err=%v", nr1, err1)
	}
	r2, _, nr2, err2 := b.Evaluate(context.Background(), reordered)
	if err2 != nil || nr2 != "" {
		t.Fatalf("reordered Evaluate: nr=%q err=%v", nr2, err2)
	}
	if r1.InputDigest != r2.InputDigest {
		t.Fatalf("reordered refs altered digest: %q vs %q", r1.InputDigest, r2.InputDigest)
	}

	idOnly := base
	idOnly.RequestID = "completely-different-correlation-id"
	r3, _, nr3, err3 := b.Evaluate(context.Background(), idOnly)
	if err3 != nil || nr3 != "" {
		t.Fatalf("requestId-only Evaluate: nr=%q err=%v", nr3, err3)
	}
	if r3.InputDigest != r1.InputDigest {
		t.Fatalf("requestId-only change altered digest: %q vs %q", r3.InputDigest, r1.InputDigest)
	}

	semantic := base
	semantic.Action = "service.write"
	r4, _, nr4, err4 := b.Evaluate(context.Background(), semantic)
	if err4 != nil || nr4 != "" {
		t.Fatalf("semantic Evaluate: nr=%q err=%v", nr4, err4)
	}
	if r4.InputDigest == r1.InputDigest {
		t.Fatal("semantic action change must alter digest")
	}
}
