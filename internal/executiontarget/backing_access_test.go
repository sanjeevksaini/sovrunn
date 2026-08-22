package executiontarget

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/cloudmodel"
	cmmodel "github.com/sanjeevksaini/sovrunn/internal/cloudmodel/model"
)

type staticBackingGrants struct {
	allow map[string]bool
}

func (g *staticBackingGrants) AuthorizeExact(action string, scope apimeta.ScopeIdentity, _, _ string) bool {
	if g == nil || g.allow == nil {
		return false
	}
	return g.allow[action+"|"+scope.UID]
}

func seedPairedBacking(t *testing.T, s *cloudmodel.Store) (partUID, stackUID, providerUID string) {
	t.Helper()
	platformUID := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	providerUID = "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
	partUID = "cccccccccccccccccccccccccccccccc"
	stackUID = "dddddddddddddddddddddddddddddddd"
	if _, prob := s.CreateCloudPlatform(cmmodel.CloudPlatform{
		Metadata: apimeta.ObjectMeta{Name: "plat", UID: platformUID},
		Spec: cmmodel.CloudPlatformSpec{OwnerRegistration: cmmodel.OwnerRegistration{
			LegalName: "Acme", RegistrationIdentifier: "R1", JurisdictionCode: "US",
		}},
		Status: cmmodel.CloudPlatformStatus{Phase: cmmodel.CloudPlatformPhaseActive},
	}); prob != nil {
		t.Fatalf("platform: %#v", prob)
	}
	if _, prob := s.CreateCloudProvider(cmmodel.CloudProvider{
		Metadata: apimeta.ObjectMeta{Name: "prov", UID: providerUID},
		Spec:     cmmodel.CloudProviderSpec{OperatingMarkets: []string{"US"}},
		Status:   cmmodel.CloudProviderStatus{Phase: cmmodel.CloudProviderPhaseActive},
	}); prob != nil {
		t.Fatalf("provider: %#v", prob)
	}
	if _, prob := s.CreateParticipation(cmmodel.CloudProviderParticipation{
		Metadata: apimeta.ObjectMeta{
			Name: "part", UID: partUID, ResourceVersion: "1", Generation: 7,
			ScopeRef: &apimeta.ScopeRef{TypedRef: apimeta.TypedRef{
				APIVersion: cmmodel.APIVersionCloudPlatform, Kind: string(apimeta.ScopeCloudPlatform),
				Name: "plat", UID: platformUID,
			}},
		},
		Spec: cmmodel.CloudProviderParticipationSpec{
			CloudPlatformRef: apimeta.TypedRef{UID: platformUID},
			CloudProviderRef: apimeta.TypedRef{UID: providerUID},
		},
		Status: cmmodel.CloudProviderParticipationStatus{Phase: cmmodel.ParticipationPhaseActive},
	}); prob != nil {
		t.Fatalf("participation: %#v", prob)
	}
	scope := &apimeta.ScopeRef{TypedRef: apimeta.TypedRef{
		APIVersion: cmmodel.APIVersionCloudProvider, Kind: string(apimeta.ScopeCloudProvider),
		Name: "prov", UID: providerUID,
	}}
	if _, prob := s.CreateInfrastructureStack(cmmodel.InfrastructureStack{
		Metadata: apimeta.ObjectMeta{Name: "stack", UID: stackUID, ScopeRef: scope, ResourceVersion: "1", Generation: 4},
		Status:   cmmodel.InfrastructureStackStatus{Phase: cmmodel.InfrastructureStackPhaseActive},
	}); prob != nil {
		t.Fatalf("stack: %#v", prob)
	}
	return partUID, stackUID, providerUID
}

func TestBackingAccessProvider_CallerRelativeDispositions(t *testing.T) {
	t.Parallel()
	store := cloudmodel.NewStore()
	partUID, stackUID, providerUID := seedPairedBacking(t, store)
	provider := NewBackingAccessProvider(store)

	t.Run("allowed", func(t *testing.T) {
		grants := &staticBackingGrants{allow: map[string]bool{
			"executiontarget.write|" + providerUID: true,
		}}
		res := provider.Access(context.Background(), "p1", "executiontarget.write", partUID, stackUID, grants)
		if res.Disposition != BackingAllowed {
			t.Fatalf("disposition=%v", res.Disposition)
		}
		if res.View.ParticipationUID != partUID || res.View.StackUID != stackUID {
			t.Fatalf("view=%#v", res.View)
		}
		if !res.View.ParticipationEffectiveActive || !res.View.StackActive() {
			t.Fatalf("viability=%#v", res.View)
		}
		if res.View.StackGeneration != 4 {
			t.Fatalf("generation=%d", res.View.StackGeneration)
		}
	})

	t.Run("safe-denied-missing", func(t *testing.T) {
		grants := &staticBackingGrants{allow: map[string]bool{
			"executiontarget.write|" + providerUID: true,
		}}
		res := provider.Access(context.Background(), "p1", "executiontarget.write", partUID, "ffffffffffffffffffffffffffffffff", grants)
		if res.Disposition != BackingSafeDenied {
			t.Fatalf("disposition=%v", res.Disposition)
		}
		if res.View.ParticipationUID != "" || res.View.StackUID != "" {
			t.Fatalf("safe denial must not expose view fields: %#v", res.View)
		}
	})

	t.Run("authorization-denied", func(t *testing.T) {
		grants := &staticBackingGrants{allow: map[string]bool{}}
		res := provider.Access(context.Background(), "p1", "executiontarget.write", partUID, stackUID, grants)
		if res.Disposition != BackingAuthorizationDenied {
			t.Fatalf("disposition=%v", res.Disposition)
		}
		if res.View.ParticipationUID != "" {
			t.Fatalf("authz denial must not expose view: %#v", res.View)
		}
	})
}

func TestBackingAccessProvider_FinalLeaseHoldsThroughCallback(t *testing.T) {
	t.Parallel()
	store := cloudmodel.NewStore()
	partUID, stackUID, providerUID := seedPairedBacking(t, store)
	provider := NewBackingAccessProvider(store)
	grants := &staticBackingGrants{allow: map[string]bool{
		"executiontarget.qualify|" + providerUID: true,
	}}

	entered := make(chan struct{})
	release := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		disposition := provider.WithFinalQualificationLease(
			context.Background(), "p1", "executiontarget.qualify", partUID, stackUID, grants,
			func(view BackingView) {
				if view.StackGeneration != 4 {
					t.Errorf("generation=%d", view.StackGeneration)
				}
				close(entered)
				<-release
			},
		)
		if disposition != BackingAllowed {
			t.Errorf("disposition=%v", disposition)
		}
	}()
	<-entered

	writeDone := make(chan struct{})
	go func() {
		cur, _ := store.GetParticipation(partUID)
		cur.Status.ProviderSuspended = true
		if _, prob := store.UpdateParticipation(cur); prob != nil {
			t.Errorf("update: %#v", prob)
		}
		close(writeDone)
	}()
	select {
	case <-writeDone:
		t.Fatal("backing write interleaved through final qualification lease")
	case <-time.After(50 * time.Millisecond):
	}
	close(release)
	<-writeDone
	wg.Wait()
}

func TestBackingAccessProvider_NoParticipationGenerationFence(t *testing.T) {
	t.Parallel()
	store := cloudmodel.NewStore()
	partUID, stackUID, providerUID := seedPairedBacking(t, store)
	provider := NewBackingAccessProvider(store)
	grants := &staticBackingGrants{allow: map[string]bool{
		"executiontarget.qualify|" + providerUID: true,
	}}
	before := provider.Access(context.Background(), "p1", "executiontarget.qualify", partUID, stackUID, grants)
	if before.Disposition != BackingAllowed {
		t.Fatalf("before=%v", before.Disposition)
	}
	fpBefore := before.View.Fingerprint()

	cur, _ := store.GetParticipation(partUID)
	cur.Metadata.Generation = 99
	if _, prob := store.UpdateParticipation(cur); prob != nil {
		t.Fatalf("update generation: %#v", prob)
	}
	after := provider.Access(context.Background(), "p1", "executiontarget.qualify", partUID, stackUID, grants)
	if after.Disposition != BackingAllowed {
		t.Fatalf("after=%v", after.Disposition)
	}
	if after.View.Fingerprint() != fpBefore {
		t.Fatal("participation generation must not affect viability fingerprint fence")
	}
	// BackingView has no participation-generation field by construction.
	_ = after.View
}

func TestBackingView_AsViability(t *testing.T) {
	t.Parallel()
	v := BackingView{
		ParticipationUID:              "p",
		ParticipationEffectiveActive:  true,
		ParticipationCloudProviderUID: "c",
		StackUID:                      "s",
		StackPhase:                    "Active",
		StackScopeUID:                 "c",
		StackGeneration:               2,
	}
	b := v.AsViability()
	if b.ParticipationUID != "p" || b.StackGeneration != 2 || !b.StackActive() {
		t.Fatalf("%#v", b)
	}
	if b.Fingerprint() != v.Fingerprint() {
		t.Fatal("fingerprint mismatch")
	}
}
