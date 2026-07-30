package validation

import (
	"sync"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/apicond"
	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/resources"
)

func TestEvaluateTopologyCompleteness_PositiveFullFiveLevelPath(t *testing.T) {
	t.Parallel()

	set := fullCompletePathSet("provider-uid-a")
	provider := findByKind(t, set, resources.KindProvider)

	got := EvaluateTopologyCompleteness(provider, set)
	if !got.Complete {
		t.Fatalf("expected complete five-level path; reason=%s", got.Reason)
	}
	if got.Reason != ReasonPathComplete {
		t.Fatalf("reason: got %q want %q", got.Reason, ReasonPathComplete)
	}
	if !got.HasChildren {
		t.Fatal("expected HasChildren true for provider with locations")
	}

	// Every ancestor on the path is complete.
	for _, kind := range []string{
		resources.KindProvider,
		resources.KindProviderLocation,
		resources.KindProviderDatacenter,
		resources.KindDatacenterFailureDomain,
	} {
		v := findByKind(t, set, kind)
		ev := EvaluateTopologyCompleteness(v, set)
		if !ev.Complete || ev.Reason != ReasonPathComplete {
			t.Fatalf("%s: complete=%v reason=%q", kind, ev.Complete, ev.Reason)
		}
	}
}

func TestEvaluateTopologyCompleteness_PositiveLeafStackValidIsComplete(t *testing.T) {
	t.Parallel()

	set := fullCompletePathSet("provider-uid-leaf")
	stack := findByKind(t, set, resources.KindInfrastructureStack)

	got := EvaluateTopologyCompleteness(stack, set)
	if !got.Complete {
		t.Fatal("expected Valid leaf InfrastructureStack to be complete")
	}
	if got.Reason != ReasonLeafValid {
		t.Fatalf("reason: got %q want %q", got.Reason, ReasonLeafValid)
	}
	if got.HasChildren {
		t.Fatal("leaf stack must not report topology children")
	}

	// Invalid leaf is incomplete (registered≠complete distinction for leaf).
	invalid := stack
	invalid.Valid = false
	got = EvaluateTopologyCompleteness(invalid, set)
	if got.Complete || got.Reason != ReasonPathIncomplete {
		t.Fatalf("invalid leaf: complete=%v reason=%q", got.Complete, got.Reason)
	}
}

func TestEvaluateTopologyCompleteness_NegativeMissingLevelIncomplete(t *testing.T) {
	t.Parallel()

	providerUID := "provider-uid-missing"
	full := fullCompletePathSet(providerUID)

	// Drop DatacenterFailureDomain and InfrastructureStack → path broken.
	truncated := make([]CompletenessValue, 0, 3)
	for _, v := range full {
		switch v.Kind {
		case resources.KindDatacenterFailureDomain, resources.KindInfrastructureStack:
			continue
		default:
			truncated = append(truncated, v)
		}
	}
	provider := findByKind(t, truncated, resources.KindProvider)
	got := EvaluateTopologyCompleteness(provider, truncated)
	if got.Complete {
		t.Fatal("expected incomplete when failure-domain/stack levels are missing")
	}
	if got.Reason != ReasonPathIncomplete {
		t.Fatalf("reason: got %q want %q", got.Reason, ReasonPathIncomplete)
	}
}

func TestEvaluateTopologyCompleteness_NegativeEmptyRegisteredParentIncomplete(t *testing.T) {
	t.Parallel()

	providerUID := "provider-uid-empty"
	// Registered Valid parent with zero children: incomplete, not complete
	// (F14-REQ-12, F14-REQ-13).
	provider := CompletenessValue{
		Kind:             resources.KindProvider,
		UID:              providerUID,
		ProviderScopeUID: providerUID,
		Valid:            true,
	}
	got := EvaluateTopologyCompleteness(provider, []CompletenessValue{provider})
	if got.Complete {
		t.Fatal("empty registered parent must be incomplete")
	}
	if got.Reason != ReasonPathIncomplete {
		t.Fatalf("reason: got %q want %q", got.Reason, ReasonPathIncomplete)
	}
	if got.HasChildren {
		t.Fatal("empty parent must not report children")
	}
	if !provider.Valid {
		t.Fatal("registered Valid fact must remain distinct from completeness")
	}
}

func TestEvaluateTopologyCompleteness_NegativeSkipNeverComplete(t *testing.T) {
	t.Parallel()

	providerUID := "provider-uid-skip"
	provider := CompletenessValue{
		Kind:             resources.KindProvider,
		UID:              providerUID,
		ProviderScopeUID: providerUID,
		Valid:            true,
	}
	// Skip: datacenter claims Provider as parent (missing ProviderLocation).
	skippedDC := CompletenessValue{
		Kind:             resources.KindProviderDatacenter,
		UID:              "dc-skip",
		ProviderScopeUID: providerUID,
		ParentUID:        providerUID,
		Valid:            true,
	}
	stackUnderSkip := CompletenessValue{
		Kind:             resources.KindInfrastructureStack,
		UID:              "stack-skip",
		ProviderScopeUID: providerUID,
		ParentUID:        "dc-skip",
		Valid:            true,
	}
	set := []CompletenessValue{provider, skippedDC, stackUnderSkip}
	got := EvaluateTopologyCompleteness(provider, set)
	if got.Complete {
		t.Fatal("hierarchy skip must never evaluate as complete")
	}
	if got.HasChildren {
		t.Fatal("skipped wrong-kind link must not count as immediate child")
	}
}

func TestHasImmediateTopologyChildren_DeleteBlockedOutcome(t *testing.T) {
	t.Parallel()

	set := fullCompletePathSet("provider-uid-delete")
	provider := findByKind(t, set, resources.KindProvider)
	loc := findByKind(t, set, resources.KindProviderLocation)

	if !HasImmediateTopologyChildren(provider, set) {
		t.Fatal("expected child-existence (delete-blocked) for provider with location")
	}
	if !HasImmediateTopologyChildren(loc, set) {
		t.Fatal("expected child-existence for location with datacenter")
	}

	stack := findByKind(t, set, resources.KindInfrastructureStack)
	if HasImmediateTopologyChildren(stack, set) {
		t.Fatal("leaf must not be delete-blocked by topology children")
	}

	empty := CompletenessValue{
		Kind:             resources.KindProvider,
		UID:              "provider-uid-lonely",
		ProviderScopeUID: "provider-uid-lonely",
		Valid:            true,
	}
	if HasImmediateTopologyChildren(empty, []CompletenessValue{empty}) {
		t.Fatal("zero-child parent must not be delete-blocked")
	}
}

func TestEvaluateTopologyCompleteness_BoundaryOrderIndependent(t *testing.T) {
	t.Parallel()

	providerUID := "provider-uid-order"
	base := fullCompletePathSet(providerUID)
	provider := findByKind(t, base, resources.KindProvider)

	orders := [][]CompletenessValue{
		base,
		reverseCompleteness(base),
		rotateCompleteness(base, 2),
		rotateCompleteness(base, 4),
	}

	var first TopologyCompletenessEvaluation
	for i, set := range orders {
		got := EvaluateTopologyCompleteness(provider, set)
		if i == 0 {
			first = got
			continue
		}
		if got.Complete != first.Complete || got.Reason != first.Reason || got.HasChildren != first.HasChildren {
			t.Fatalf("order %d: got %+v want %+v", i, got, first)
		}
	}
}

func TestEvaluateTopologyCompleteness_IsolationDistinctProviders(t *testing.T) {
	t.Parallel()

	setA := fullCompletePathSet("provider-uid-a")
	setB := fullCompletePathSet("provider-uid-b")
	combined := append(append([]CompletenessValue{}, setA...), setB...)

	providerA := findByKind(t, setA, resources.KindProvider)
	providerB := findByKind(t, setB, resources.KindProvider)

	gotA := EvaluateTopologyCompleteness(providerA, combined)
	gotB := EvaluateTopologyCompleteness(providerB, combined)
	if !gotA.Complete || !gotB.Complete {
		t.Fatalf("each provider path must remain complete in combined set; A=%+v B=%+v", gotA, gotB)
	}

	// Provider A must not treat Provider B's children as its own.
	onlyB := setB
	gotCross := EvaluateTopologyCompleteness(providerA, onlyB)
	if gotCross.Complete || gotCross.HasChildren {
		t.Fatalf("cross-provider set must not complete provider A; got %+v", gotCross)
	}
}

func TestEvaluateTopologyCompleteness_ConcurrentRaceStable(t *testing.T) {
	t.Parallel()

	set := fullCompletePathSet("provider-uid-race")
	provider := findByKind(t, set, resources.KindProvider)
	stack := findByKind(t, set, resources.KindInfrastructureStack)
	empty := CompletenessValue{
		Kind:             resources.KindProvider,
		UID:              "provider-uid-empty-race",
		ProviderScopeUID: "provider-uid-empty-race",
		Valid:            true,
	}

	wantProvider := EvaluateTopologyCompleteness(provider, set)
	wantStack := EvaluateTopologyCompleteness(stack, set)
	wantEmpty := EvaluateTopologyCompleteness(empty, []CompletenessValue{empty})

	const goroutines = 32
	const iterations = 200
	var wg sync.WaitGroup
	wg.Add(goroutines)
	errCh := make(chan string, goroutines)

	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				gotP := EvaluateTopologyCompleteness(provider, set)
				if gotP != wantProvider {
					errCh <- "provider completeness drifted under concurrency"
					return
				}
				gotS := EvaluateTopologyCompleteness(stack, set)
				if gotS != wantStack {
					errCh <- "stack completeness drifted under concurrency"
					return
				}
				gotE := EvaluateTopologyCompleteness(empty, []CompletenessValue{empty})
				if gotE != wantEmpty {
					errCh <- "empty-parent completeness drifted under concurrency"
					return
				}
				if HasImmediateTopologyChildren(provider, set) != wantProvider.HasChildren {
					errCh <- "child-existence drifted under concurrency"
					return
				}
			}
		}()
	}
	wg.Wait()
	close(errCh)
	for msg := range errCh {
		t.Fatal(msg)
	}
}

func TestCompletenessValueFromResources_Projection(t *testing.T) {
	t.Parallel()

	providerUID := "provider-uid-proj"
	provider := resources.Provider{
		Metadata: apimeta.ObjectMeta{Name: "prov", UID: providerUID},
		Status: resources.ProviderStatus{
			Conditions: []apicond.Condition{{
				Type:   ConditionTypeValid,
				Status: apicond.ConditionTrue,
				Reason: "ValidationSucceeded",
			}},
		},
	}
	provider.Kind = resources.KindProvider

	loc := resources.ProviderLocation{
		Metadata: apimeta.ObjectMeta{
			Name: "loc",
			UID:  "loc-uid",
			ScopeRef: &apimeta.ScopeRef{TypedRef: apimeta.TypedRef{
				APIVersion: resources.FabricAPIVersion,
				Kind:       string(apimeta.ScopeProvider),
				Name:       "prov",
				UID:        providerUID,
			}},
		},
		Status: resources.ProviderLocationStatus{
			Conditions: []apicond.Condition{{
				Type:   ConditionTypeValid,
				Status: apicond.ConditionTrue,
				Reason: "ValidationSucceeded",
			}},
		},
	}
	loc.Kind = resources.KindProviderLocation

	dc := resources.ProviderDatacenter{
		Metadata: apimeta.ObjectMeta{
			Name: "dc",
			UID:  "dc-uid",
			ScopeRef: &apimeta.ScopeRef{TypedRef: apimeta.TypedRef{
				Kind: string(apimeta.ScopeProvider),
				UID:  providerUID,
			}},
		},
		Spec: resources.ProviderDatacenterSpec{
			ProviderLocationRef: resources.ProviderLocationRef{
				TypedRef: apimeta.TypedRef{UID: "loc-uid", Kind: resources.KindProviderLocation},
			},
		},
		Status: resources.ProviderDatacenterStatus{
			Conditions: []apicond.Condition{{
				Type:   ConditionTypeValid,
				Status: apicond.ConditionTrue,
				Reason: "ValidationSucceeded",
			}},
		},
	}
	dc.Kind = resources.KindProviderDatacenter

	fd := resources.DatacenterFailureDomain{
		Metadata: apimeta.ObjectMeta{
			Name: "fd",
			UID:  "fd-uid",
			ScopeRef: &apimeta.ScopeRef{TypedRef: apimeta.TypedRef{
				Kind: string(apimeta.ScopeProvider),
				UID:  providerUID,
			}},
		},
		Spec: resources.DatacenterFailureDomainSpec{
			ProviderDatacenterRef: resources.ProviderDatacenterRef{
				TypedRef: apimeta.TypedRef{UID: "dc-uid", Kind: resources.KindProviderDatacenter},
			},
		},
		Status: resources.DatacenterFailureDomainStatus{
			Conditions: []apicond.Condition{{
				Type:   ConditionTypeValid,
				Status: apicond.ConditionTrue,
				Reason: "ValidationSucceeded",
			}},
		},
	}
	fd.Kind = resources.KindDatacenterFailureDomain

	stack := resources.InfrastructureStack{
		Metadata: apimeta.ObjectMeta{
			Name: "stack",
			UID:  "stack-uid",
			ScopeRef: &apimeta.ScopeRef{TypedRef: apimeta.TypedRef{
				Kind: string(apimeta.ScopeProvider),
				UID:  providerUID,
			}},
		},
		Spec: resources.InfrastructureStackSpec{
			DatacenterFailureDomainRef: resources.DatacenterFailureDomainRef{
				TypedRef: apimeta.TypedRef{UID: "fd-uid", Kind: resources.KindDatacenterFailureDomain},
			},
			Technology: "Apache CloudStack",
		},
		Status: resources.InfrastructureStackStatus{
			Conditions: []apicond.Condition{{
				Type:   ConditionTypeValid,
				Status: apicond.ConditionTrue,
				Reason: "ValidationSucceeded",
			}},
		},
	}
	stack.Kind = resources.KindInfrastructureStack

	set := []CompletenessValue{
		CompletenessValueFromProvider(provider),
		CompletenessValueFromProviderLocation(loc),
		CompletenessValueFromProviderDatacenter(dc),
		CompletenessValueFromDatacenterFailureDomain(fd),
		CompletenessValueFromInfrastructureStack(stack),
	}

	pv := CompletenessValueFromProvider(provider)
	if pv.ParentUID != "" || !pv.Valid || pv.ProviderScopeUID != providerUID {
		t.Fatalf("provider projection: %+v", pv)
	}
	lv := CompletenessValueFromProviderLocation(loc)
	if lv.ParentUID != providerUID {
		t.Fatalf("location ParentUID: got %q want %q", lv.ParentUID, providerUID)
	}

	got := EvaluateTopologyCompleteness(pv, set)
	if !got.Complete || got.Reason != ReasonPathComplete {
		t.Fatalf("projected set incomplete: %+v", got)
	}
	gotStack := EvaluateTopologyCompleteness(CompletenessValueFromInfrastructureStack(stack), set)
	if !gotStack.Complete || gotStack.Reason != ReasonLeafValid {
		t.Fatalf("projected leaf incomplete: %+v", gotStack)
	}
}

func TestCompletenessValueFromInfrastructureStack_ValidMustBeCurrentGeneration(t *testing.T) {
	t.Parallel()

	makeStack := func(metadataGeneration, statusObservedGeneration, conditionObservedGeneration int64) resources.InfrastructureStack {
		stack := resources.InfrastructureStack{
			Metadata: apimeta.ObjectMeta{
				UID:        "stack-generation-uid",
				Generation: metadataGeneration,
			},
			Status: resources.InfrastructureStackStatus{
				ObservedGeneration: statusObservedGeneration,
				Conditions: []apicond.Condition{{
					Type:               ConditionTypeValid,
					Status:             apicond.ConditionTrue,
					Reason:             "ValidationSucceeded",
					ObservedGeneration: conditionObservedGeneration,
				}},
			},
		}
		stack.Kind = resources.KindInfrastructureStack
		return stack
	}

	cases := []struct {
		name                        string
		metadataGeneration          int64
		statusObservedGeneration    int64
		conditionObservedGeneration int64
		wantValid                   bool
	}{
		{
			name:                        "all-current",
			metadataGeneration:          2,
			statusObservedGeneration:    2,
			conditionObservedGeneration: 2,
			wantValid:                   true,
		},
		{
			name:                        "stale-status-and-condition",
			metadataGeneration:          2,
			statusObservedGeneration:    1,
			conditionObservedGeneration: 1,
		},
		{
			name:                        "stale-status",
			metadataGeneration:          2,
			statusObservedGeneration:    1,
			conditionObservedGeneration: 2,
		},
		{
			name:                        "stale-condition",
			metadataGeneration:          2,
			statusObservedGeneration:    2,
			conditionObservedGeneration: 1,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			value := CompletenessValueFromInfrastructureStack(makeStack(
				tc.metadataGeneration,
				tc.statusObservedGeneration,
				tc.conditionObservedGeneration,
			))
			if value.Valid != tc.wantValid {
				t.Fatalf("Valid = %v, want %v", value.Valid, tc.wantValid)
			}
			evaluation := EvaluateTopologyCompleteness(value, []CompletenessValue{value})
			if evaluation.Complete != tc.wantValid {
				t.Fatalf("Complete = %v, want %v", evaluation.Complete, tc.wantValid)
			}
		})
	}
}

func fullCompletePathSet(providerUID string) []CompletenessValue {
	locUID := providerUID + "-loc"
	dcUID := providerUID + "-dc"
	fdUID := providerUID + "-fd"
	stackUID := providerUID + "-stack"
	return []CompletenessValue{
		{
			Kind:             resources.KindProvider,
			UID:              providerUID,
			ProviderScopeUID: providerUID,
			Valid:            true,
		},
		{
			Kind:             resources.KindProviderLocation,
			UID:              locUID,
			ProviderScopeUID: providerUID,
			ParentUID:        providerUID,
			Valid:            true,
		},
		{
			Kind:             resources.KindProviderDatacenter,
			UID:              dcUID,
			ProviderScopeUID: providerUID,
			ParentUID:        locUID,
			Valid:            true,
		},
		{
			Kind:             resources.KindDatacenterFailureDomain,
			UID:              fdUID,
			ProviderScopeUID: providerUID,
			ParentUID:        dcUID,
			Valid:            true,
		},
		{
			Kind:             resources.KindInfrastructureStack,
			UID:              stackUID,
			ProviderScopeUID: providerUID,
			ParentUID:        fdUID,
			Valid:            true,
		},
	}
}

func findByKind(t *testing.T, set []CompletenessValue, kind string) CompletenessValue {
	t.Helper()
	for _, v := range set {
		if v.Kind == kind {
			return v
		}
	}
	t.Fatalf("kind %s not found in set", kind)
	return CompletenessValue{}
}

func reverseCompleteness(in []CompletenessValue) []CompletenessValue {
	out := make([]CompletenessValue, len(in))
	for i := range in {
		out[len(in)-1-i] = in[i]
	}
	return out
}

func rotateCompleteness(in []CompletenessValue, n int) []CompletenessValue {
	if len(in) == 0 {
		return nil
	}
	n = n % len(in)
	out := make([]CompletenessValue, len(in))
	copy(out, in[n:])
	copy(out[len(in)-n:], in[:n])
	return out
}
