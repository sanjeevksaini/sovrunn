package validation

import (
	"sync"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/resources"
)

func TestEvaluateScopeAndHierarchy_PositiveSameProviderCorrectOrder(t *testing.T) {
	t.Parallel()

	providerUID := "provider-uid-a"

	cases := []struct {
		name   string
		parent TopologyValue
		child  TopologyValue
	}{
		{
			name: "Provider/ProviderLocation",
			parent: TopologyValue{
				Kind:             resources.KindProvider,
				UID:              providerUID,
				ProviderScopeUID: providerUID,
			},
			child: TopologyValue{
				Kind:             resources.KindProviderLocation,
				UID:              "loc-uid-1",
				ProviderScopeUID: providerUID,
			},
		},
		{
			name: "ProviderLocation/ProviderDatacenter",
			parent: TopologyValue{
				Kind:             resources.KindProviderLocation,
				UID:              "loc-uid-1",
				ProviderScopeUID: providerUID,
			},
			child: TopologyValue{
				Kind:             resources.KindProviderDatacenter,
				UID:              "dc-uid-1",
				ProviderScopeUID: providerUID,
			},
		},
		{
			name: "ProviderDatacenter/DatacenterFailureDomain",
			parent: TopologyValue{
				Kind:             resources.KindProviderDatacenter,
				UID:              "dc-uid-1",
				ProviderScopeUID: providerUID,
			},
			child: TopologyValue{
				Kind:             resources.KindDatacenterFailureDomain,
				UID:              "fd-uid-1",
				ProviderScopeUID: providerUID,
			},
		},
		{
			name: "DatacenterFailureDomain/InfrastructureStack",
			parent: TopologyValue{
				Kind:             resources.KindDatacenterFailureDomain,
				UID:              "fd-uid-1",
				ProviderScopeUID: providerUID,
			},
			child: TopologyValue{
				Kind:             resources.KindInfrastructureStack,
				UID:              "stack-uid-1",
				ProviderScopeUID: providerUID,
			},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if !EvaluateScopeAndHierarchy(tc.parent, tc.child) {
				t.Fatalf("expected agreement for %s", tc.name)
			}
		})
	}
}

func TestEvaluateScopeAndHierarchy_NegativeDifferentProviderScopeUIDs(t *testing.T) {
	t.Parallel()

	parent := TopologyValue{
		Kind:             resources.KindProviderLocation,
		UID:              "loc-a",
		ProviderScopeUID: "provider-uid-a",
	}
	child := TopologyValue{
		Kind:             resources.KindProviderDatacenter,
		UID:              "dc-b",
		ProviderScopeUID: "provider-uid-b",
	}
	if EvaluateScopeAndHierarchy(parent, child) {
		t.Fatal("expected rejection for different Provider scope UIDs")
	}
}

func TestEvaluateScopeAndHierarchy_NegativeWrongHierarchySkip(t *testing.T) {
	t.Parallel()

	providerUID := "provider-uid-a"
	// Skip: Provider directly parenting ProviderDatacenter (missing ProviderLocation).
	parent := TopologyValue{
		Kind:             resources.KindProvider,
		UID:              providerUID,
		ProviderScopeUID: providerUID,
	}
	child := TopologyValue{
		Kind:             resources.KindProviderDatacenter,
		UID:              "dc-uid-1",
		ProviderScopeUID: providerUID,
	}
	if EvaluateScopeAndHierarchy(parent, child) {
		t.Fatal("expected rejection for skipped hierarchy level")
	}
}

func TestEvaluateScopeAndHierarchy_NegativeWrongHierarchyReorder(t *testing.T) {
	t.Parallel()

	providerUID := "provider-uid-a"
	// Reorder: ProviderDatacenter as parent of ProviderLocation.
	parent := TopologyValue{
		Kind:             resources.KindProviderDatacenter,
		UID:              "dc-uid-1",
		ProviderScopeUID: providerUID,
	}
	child := TopologyValue{
		Kind:             resources.KindProviderLocation,
		UID:              "loc-uid-1",
		ProviderScopeUID: providerUID,
	}
	if EvaluateScopeAndHierarchy(parent, child) {
		t.Fatal("expected rejection for reordered hierarchy level")
	}
}

func TestEvaluateScopeAndHierarchy_IsolationTwoProvidersSameOwnerOrganization(t *testing.T) {
	t.Parallel()

	// Two Providers under the same owner Organization remain distinct scope UIDs
	// and must not agree across providers. Ownership Organization is external to
	// this helper; isolation is expressed only by ProviderScopeUID (F14-REQ-30).

	providerA := TopologyValue{
		Kind:             resources.KindProvider,
		UID:              "provider-uid-a",
		ProviderScopeUID: "provider-uid-a",
	}
	providerB := TopologyValue{
		Kind:             resources.KindProvider,
		UID:              "provider-uid-b",
		ProviderScopeUID: "provider-uid-b",
	}
	locUnderA := TopologyValue{
		Kind:             resources.KindProviderLocation,
		UID:              "loc-under-a",
		ProviderScopeUID: "provider-uid-a",
	}
	locUnderB := TopologyValue{
		Kind:             resources.KindProviderLocation,
		UID:              "loc-under-b",
		ProviderScopeUID: "provider-uid-b",
	}

	if !EvaluateScopeAndHierarchy(providerA, locUnderA) {
		t.Fatal("expected agreement within provider A")
	}
	if !EvaluateScopeAndHierarchy(providerB, locUnderB) {
		t.Fatal("expected agreement within provider B")
	}
	if EvaluateScopeAndHierarchy(providerA, locUnderB) {
		t.Fatal("expected rejection across providers under the same owner Organization")
	}
	if EvaluateScopeAndHierarchy(providerB, locUnderA) {
		t.Fatal("expected rejection across providers under the same owner Organization")
	}
}

func TestEvaluateScopeAndHierarchy_IsolationSameOperatorDifferentOwners(t *testing.T) {
	t.Parallel()

	// Same external operator name under two owner Organizations yields
	// independent Provider scope UIDs. The helper keys only on scope UID,
	// never on operator name or native ID (F14-REQ-30).
	providerOwner1 := TopologyValue{
		Kind:             resources.KindProvider,
		UID:              "provider-uid-owner1",
		ProviderScopeUID: "provider-uid-owner1",
	}
	providerOwner2 := TopologyValue{
		Kind:             resources.KindProvider,
		UID:              "provider-uid-owner2",
		ProviderScopeUID: "provider-uid-owner2",
	}

	locOwner1 := TopologyValue{
		Kind:             resources.KindProviderLocation,
		UID:              "loc-owner1",
		ProviderScopeUID: "provider-uid-owner1",
	}
	locOwner2 := TopologyValue{
		Kind:             resources.KindProviderLocation,
		UID:              "loc-owner2",
		ProviderScopeUID: "provider-uid-owner2",
	}

	if !EvaluateScopeAndHierarchy(providerOwner1, locOwner1) {
		t.Fatal("expected agreement under owner 1")
	}
	if !EvaluateScopeAndHierarchy(providerOwner2, locOwner2) {
		t.Fatal("expected agreement under owner 2")
	}
	if EvaluateScopeAndHierarchy(providerOwner1, locOwner2) {
		t.Fatal("expected rejection across owners for same external operator name")
	}
	if EvaluateScopeAndHierarchy(providerOwner2, locOwner1) {
		t.Fatal("expected rejection across owners for same external operator name")
	}
}

func TestEvaluateScopeAndHierarchy_EmptyScopeUIDsReject(t *testing.T) {
	t.Parallel()

	parent := TopologyValue{
		Kind:             resources.KindProvider,
		UID:              "provider-uid-a",
		ProviderScopeUID: "provider-uid-a",
	}
	child := TopologyValue{
		Kind:             resources.KindProviderLocation,
		UID:              "loc-uid-1",
		ProviderScopeUID: "",
	}
	if EvaluateScopeAndHierarchy(parent, child) {
		t.Fatal("expected rejection when child ProviderScopeUID is empty")
	}
}

func TestRequiredImmediateParentKind(t *testing.T) {
	t.Parallel()

	cases := []struct {
		child  string
		parent string
		ok     bool
	}{
		{resources.KindProviderLocation, resources.KindProvider, true},
		{resources.KindProviderDatacenter, resources.KindProviderLocation, true},
		{resources.KindDatacenterFailureDomain, resources.KindProviderDatacenter, true},
		{resources.KindInfrastructureStack, resources.KindDatacenterFailureDomain, true},
		{resources.KindProvider, "", false},
		{"UnknownKind", "", false},
	}
	for _, tc := range cases {
		got, ok := RequiredImmediateParentKind(tc.child)
		if ok != tc.ok || got != tc.parent {
			t.Fatalf("child=%s: got (%q, %v), want (%q, %v)", tc.child, got, ok, tc.parent, tc.ok)
		}
	}
}

func TestTopologyValueFromResources_Projection(t *testing.T) {
	t.Parallel()

	providerUID := "provider-uid-a"
	provider := resources.Provider{
		Metadata: apimeta.ObjectMeta{
			Name: "provider-operator-a",
			UID:  providerUID,
			ScopeRef: &apimeta.ScopeRef{TypedRef: apimeta.TypedRef{
				APIVersion: "core.sovrunn.io/v1alpha1",
				Kind:       string(apimeta.ScopeOrganization),
				Name:       "owner-org-a",
				UID:        "org-uid-a",
			}},
		},
	}
	provider.Kind = resources.KindProvider

	loc := resources.ProviderLocation{
		Metadata: apimeta.ObjectMeta{
			Name: "loc-a",
			UID:  "loc-uid-a",
			ScopeRef: &apimeta.ScopeRef{TypedRef: apimeta.TypedRef{
				APIVersion: resources.FabricAPIVersion,
				Kind:       string(apimeta.ScopeProvider),
				Name:       "provider-operator-a",
				UID:        providerUID,
			}},
		},
	}
	loc.Kind = resources.KindProviderLocation

	pv := TopologyValueFromProvider(provider)
	lv := TopologyValueFromProviderLocation(loc)
	if pv.ProviderScopeUID != providerUID {
		t.Fatalf("provider scope UID: got %q want %q", pv.ProviderScopeUID, providerUID)
	}
	if lv.ProviderScopeUID != providerUID {
		t.Fatalf("location scope UID: got %q want %q", lv.ProviderScopeUID, providerUID)
	}
	if !EvaluateScopeAndHierarchy(pv, lv) {
		t.Fatal("expected agreement from projected Provider/ProviderLocation values")
	}

	// Non-Provider scopeRef on a descendant yields empty ProviderScopeUID.
	badLoc := loc
	badLoc.Metadata.ScopeRef = &apimeta.ScopeRef{TypedRef: apimeta.TypedRef{
		Kind: string(apimeta.ScopeOrganization),
		UID:  "org-uid-a",
	}}
	if TopologyValueFromProviderLocation(badLoc).ProviderScopeUID != "" {
		t.Fatal("expected empty ProviderScopeUID for non-Provider scopeRef")
	}
}

func TestEvaluateScopeAndHierarchy_ConcurrentRace(t *testing.T) {
	t.Parallel()

	providerUID := "provider-uid-race"
	agreeParent := TopologyValue{
		Kind:             resources.KindProvider,
		UID:              providerUID,
		ProviderScopeUID: providerUID,
	}
	agreeChild := TopologyValue{
		Kind:             resources.KindProviderLocation,
		UID:              "loc-race",
		ProviderScopeUID: providerUID,
	}
	rejectParent := TopologyValue{
		Kind:             resources.KindProvider,
		UID:              "provider-uid-other",
		ProviderScopeUID: "provider-uid-other",
	}

	const goroutines = 32
	const iterations = 200
	var wg sync.WaitGroup
	wg.Add(goroutines)
	errCh := make(chan string, goroutines)

	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				if !EvaluateScopeAndHierarchy(agreeParent, agreeChild) {
					errCh <- "expected concurrent agreement"
					return
				}
				if EvaluateScopeAndHierarchy(rejectParent, agreeChild) {
					errCh <- "expected concurrent rejection across providers"
					return
				}
				if EvaluateScopeAndHierarchy(agreeParent, TopologyValue{
					Kind:             resources.KindProviderDatacenter,
					UID:              "dc-skip",
					ProviderScopeUID: providerUID,
				}) {
					errCh <- "expected concurrent rejection for hierarchy skip"
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
