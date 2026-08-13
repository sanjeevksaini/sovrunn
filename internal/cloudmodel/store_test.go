package cloudmodel_test

import (
	"sync"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/cloudmodel"
	"github.com/sanjeevksaini/sovrunn/internal/cloudmodel/model"
)

func platformScoped(name, uid string) model.CloudPlatform {
	return model.CloudPlatform{
		Metadata: apimeta.ObjectMeta{Name: name, UID: uid},
		Spec: model.CloudPlatformSpec{
			OwnerRegistration: model.OwnerRegistration{
				LegalName: "Acme", RegistrationIdentifier: "R1", JurisdictionCode: "US",
			},
		},
		Status: model.CloudPlatformStatus{Phase: model.CloudPlatformPhaseActive},
	}
}

func providerScoped(name, uid string) model.CloudProvider {
	return model.CloudProvider{
		Metadata: apimeta.ObjectMeta{Name: name, UID: uid},
		Spec:     model.CloudProviderSpec{OperatingMarkets: []string{"US"}},
		Status:   model.CloudProviderStatus{Phase: model.CloudProviderPhaseActive},
	}
}

func TestStore_NameUniquenessPlatformRoot(t *testing.T) {
	t.Parallel()
	s := cloudmodel.NewStore()
	if _, prob := s.CreateCloudPlatform(platformScoped("plat", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")); prob != nil {
		t.Fatalf("create: %#v", prob)
	}
	_, prob := s.CreateCloudPlatform(platformScoped("plat", "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"))
	if prob == nil || prob.Code != apiproblem.CodeAlreadyExists {
		t.Fatalf("duplicate platform name: %#v", prob)
	}
	if _, prob := s.CreateCloudProvider(providerScoped("prov", "cccccccccccccccccccccccccccccccc")); prob != nil {
		t.Fatalf("provider create: %#v", prob)
	}
	_, prob = s.CreateCloudProvider(providerScoped("prov", "dddddddddddddddddddddddddddddddd"))
	if prob == nil || prob.Code != apiproblem.CodeAlreadyExists {
		t.Fatalf("duplicate provider name: %#v", prob)
	}
}

func TestStore_TopologyNameUniquenessWithinCloudProviderScope(t *testing.T) {
	t.Parallel()
	s := cloudmodel.NewStore()
	scopeA := &apimeta.ScopeRef{TypedRef: apimeta.TypedRef{
		APIVersion: model.APIVersionCloudProvider, Kind: string(apimeta.ScopeCloudProvider),
		Name: "a", UID: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
	}}
	scopeB := &apimeta.ScopeRef{TypedRef: apimeta.TypedRef{
		APIVersion: model.APIVersionCloudProvider, Kind: string(apimeta.ScopeCloudProvider),
		Name: "b", UID: "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb",
	}}
	hl := model.HostingLocation{
		Metadata: apimeta.ObjectMeta{Name: "loc", UID: "cccccccccccccccccccccccccccccccc", ScopeRef: scopeA},
		Spec:     model.HostingLocationSpec{CountryCode: "US", Locality: "Austin"},
	}
	if _, prob := s.CreateHostingLocation(hl); prob != nil {
		t.Fatalf("create: %#v", prob)
	}
	dup := hl
	dup.Metadata.UID = "dddddddddddddddddddddddddddddddd"
	if _, prob := s.CreateHostingLocation(dup); prob == nil || prob.Code != apiproblem.CodeAlreadyExists {
		t.Fatalf("same scope duplicate: %#v", prob)
	}
	other := hl
	other.Metadata.UID = "eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee"
	other.Metadata.ScopeRef = scopeB
	if _, prob := s.CreateHostingLocation(other); prob != nil {
		t.Fatalf("different provider scope should allow same name: %#v", prob)
	}
}

func TestStore_ParticipationNameAndPairUniqueness(t *testing.T) {
	t.Parallel()
	s := cloudmodel.NewStore()
	platformUID := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	providerUID := "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
	scope := &apimeta.ScopeRef{TypedRef: apimeta.TypedRef{
		APIVersion: model.APIVersionCloudPlatform, Kind: model.KindCloudPlatform,
		Name: "plat", UID: platformUID,
	}}
	p := model.CloudProviderParticipation{
		Metadata: apimeta.ObjectMeta{Name: "part", UID: "cccccccccccccccccccccccccccccccc", ScopeRef: scope},
		Spec: model.CloudProviderParticipationSpec{
			CloudPlatformRef: apimeta.TypedRef{APIVersion: model.APIVersionCloudPlatform, Kind: model.KindCloudPlatform, Name: "plat", UID: platformUID},
			CloudProviderRef: apimeta.TypedRef{APIVersion: model.APIVersionCloudProvider, Kind: model.KindCloudProvider, Name: "prov", UID: providerUID},
			Environment:      model.ParticipationEnvironmentDevelopment,
		},
		Status: model.CloudProviderParticipationStatus{Phase: model.ParticipationPhasePending},
	}
	if _, prob := s.CreateParticipation(p); prob != nil {
		t.Fatalf("create: %#v", prob)
	}
	nameDup := p
	nameDup.Metadata.UID = "dddddddddddddddddddddddddddddddd"
	nameDup.Spec.CloudProviderRef.UID = "eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee"
	if _, prob := s.CreateParticipation(nameDup); prob == nil || prob.Code != apiproblem.CodeAlreadyExists {
		t.Fatalf("name uniqueness: %#v", prob)
	}
	pairDup := p
	pairDup.Metadata.Name = "part-2"
	pairDup.Metadata.UID = "ffffffffffffffffffffffffffffffff"
	if _, prob := s.CreateParticipation(pairDup); prob == nil || prob.Code != apiproblem.CodeAlreadyExists {
		t.Fatalf("pair uniqueness: %#v", prob)
	}

	// Terminal frees the pair slot.
	cur, _ := s.GetParticipation(p.Metadata.UID)
	cur.Status.Phase = model.ParticipationPhaseRejected
	if _, prob := s.UpdateParticipation(cur); prob != nil {
		t.Fatalf("terminate: %#v", prob)
	}
	again := p
	again.Metadata.Name = "part-3"
	again.Metadata.UID = "11111111111111111111111111111111"
	if _, prob := s.CreateParticipation(again); prob != nil {
		t.Fatalf("pair after terminal: %#v", prob)
	}
}

func TestStore_ResourceVersionIncrements(t *testing.T) {
	t.Parallel()
	s := cloudmodel.NewStore()
	created, prob := s.CreateCloudPlatform(platformScoped("plat", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"))
	if prob != nil {
		t.Fatalf("create: %#v", prob)
	}
	if created.Metadata.ResourceVersion != "1" {
		t.Fatalf("rv=%q", created.Metadata.ResourceVersion)
	}
	created.Spec.Description = "x"
	updated, prob := s.UpdateCloudPlatform(created)
	if prob != nil {
		t.Fatalf("update: %#v", prob)
	}
	if updated.Metadata.ResourceVersion != "2" {
		t.Fatalf("rv after update=%q", updated.Metadata.ResourceVersion)
	}
}

func TestStore_ConcurrentCreateNameRace(t *testing.T) {
	s := cloudmodel.NewStore()
	const n = 32
	var wg sync.WaitGroup
	wg.Add(n)
	errs := make([]*apiproblem.Problem, n)
	for i := 0; i < n; i++ {
		i := i
		go func() {
			defer wg.Done()
			uid, _ := apimeta.GenerateUID()
			_, errs[i] = s.CreateCloudPlatform(platformScoped("same-name", uid))
		}()
	}
	wg.Wait()
	success := 0
	for _, e := range errs {
		if e == nil {
			success++
			continue
		}
		if e.Code != apiproblem.CodeAlreadyExists {
			t.Fatalf("unexpected %#v", e)
		}
	}
	if success != 1 {
		t.Fatalf("want exactly one winner, got %d", success)
	}
}

func TestStore_ConcurrentParticipationPairRace(t *testing.T) {
	s := cloudmodel.NewStore()
	platformUID := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	providerUID := "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
	scope := &apimeta.ScopeRef{TypedRef: apimeta.TypedRef{
		APIVersion: model.APIVersionCloudPlatform, Kind: model.KindCloudPlatform,
		Name: "plat", UID: platformUID,
	}}
	const n = 32
	var wg sync.WaitGroup
	wg.Add(n)
	errs := make([]*apiproblem.Problem, n)
	for i := 0; i < n; i++ {
		i := i
		go func() {
			defer wg.Done()
			uid, _ := apimeta.GenerateUID()
			p := model.CloudProviderParticipation{
				Metadata: apimeta.ObjectMeta{Name: "part-" + uid[:8], UID: uid, ScopeRef: scope},
				Spec: model.CloudProviderParticipationSpec{
					CloudPlatformRef: apimeta.TypedRef{APIVersion: model.APIVersionCloudPlatform, Kind: model.KindCloudPlatform, Name: "plat", UID: platformUID},
					CloudProviderRef: apimeta.TypedRef{APIVersion: model.APIVersionCloudProvider, Kind: model.KindCloudProvider, Name: "prov", UID: providerUID},
					Environment:      model.ParticipationEnvironmentDevelopment,
				},
				Status: model.CloudProviderParticipationStatus{Phase: model.ParticipationPhasePending},
			}
			_, errs[i] = s.CreateParticipation(p)
		}()
	}
	wg.Wait()
	success := 0
	for _, e := range errs {
		if e == nil {
			success++
			continue
		}
		if e.Code != apiproblem.CodeAlreadyExists {
			t.Fatalf("unexpected %#v", e)
		}
	}
	if success != 1 {
		t.Fatalf("want exactly one pair winner, got %d", success)
	}
}

func TestStore_StagingAbortDoesNotPublish(t *testing.T) {
	t.Parallel()
	s := cloudmodel.NewStore()
	s.BeginPublication()
	staged, prob := s.StageCreateCloudPlatform(platformScoped("plat", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"))
	if prob != nil {
		s.EndPublication()
		t.Fatalf("stage: %#v", prob)
	}
	s.Abort(staged)
	s.EndPublication()
	if _, ok := s.GetCloudPlatform("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"); ok {
		t.Fatal("aborted stage must not publish")
	}
}
