package cloudmodel_test

import (
	"sync"
	"testing"
	"time"

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

func TestStore_HasCloudPlatformAndTypedLists(t *testing.T) {
	t.Parallel()
	s := cloudmodel.NewStore()
	if s.HasCloudPlatform() {
		t.Fatal("empty store must report no CloudPlatform")
	}
	assertEmptyLists(t, s)

	uida := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	uidb := "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
	uidc := "cccccccccccccccccccccccccccccccc"
	if _, prob := s.CreateCloudPlatform(platformScoped("plat-b", uidb)); prob != nil {
		t.Fatalf("create platform b: %#v", prob)
	}
	if _, prob := s.CreateCloudPlatform(platformScoped("plat-a", uida)); prob != nil {
		t.Fatalf("create platform a: %#v", prob)
	}
	if !s.HasCloudPlatform() {
		t.Fatal("published CloudPlatform must be visible to HasCloudPlatform")
	}

	platforms := s.ListCloudPlatforms()
	if len(platforms) != 2 || platforms[0].Metadata.UID != uida || platforms[1].Metadata.UID != uidb {
		t.Fatalf("ListCloudPlatforms ascending uid = %#v", platforms)
	}
	platforms[0].Metadata.Name = "mutated"
	if got, _ := s.GetCloudPlatform(uida); got.Metadata.Name != "plat-a" {
		t.Fatal("ListCloudPlatforms must return an independent snapshot")
	}

	if _, prob := s.CreateCloudProvider(providerScoped("prov-c", uidc)); prob != nil {
		t.Fatalf("create provider: %#v", prob)
	}
	if _, prob := s.CreateCloudProvider(providerScoped("prov-a", uida)); prob != nil {
		t.Fatalf("create provider a: %#v", prob)
	}
	providers := s.ListCloudProviders()
	if len(providers) != 2 || providers[0].Metadata.UID != uida || providers[1].Metadata.UID != uidc {
		t.Fatalf("ListCloudProviders ascending uid = %#v", providers)
	}

	scope := &apimeta.ScopeRef{TypedRef: apimeta.TypedRef{
		APIVersion: model.APIVersionCloudPlatform, Kind: model.KindCloudPlatform,
		Name: "plat-a", UID: uida,
	}}
	partB := model.CloudProviderParticipation{
		Metadata: apimeta.ObjectMeta{Name: "part-b", UID: uidb, ScopeRef: scope},
		Spec: model.CloudProviderParticipationSpec{
			CloudPlatformRef: apimeta.TypedRef{APIVersion: model.APIVersionCloudPlatform, Kind: model.KindCloudPlatform, Name: "plat-a", UID: uida},
			CloudProviderRef: apimeta.TypedRef{APIVersion: model.APIVersionCloudProvider, Kind: model.KindCloudProvider, Name: "prov-c", UID: uidc},
			Environment:      model.ParticipationEnvironmentDevelopment,
		},
		Status: model.CloudProviderParticipationStatus{Phase: model.ParticipationPhasePending},
	}
	partA := partB
	partA.Metadata.Name = "part-a"
	partA.Metadata.UID = uida
	partA.Spec.CloudProviderRef.UID = uida
	if _, prob := s.CreateParticipation(partB); prob != nil {
		t.Fatalf("create part b: %#v", prob)
	}
	if _, prob := s.CreateParticipation(partA); prob != nil {
		t.Fatalf("create part a: %#v", prob)
	}
	parts := s.ListParticipations()
	if len(parts) != 2 || parts[0].Metadata.UID != uida || parts[1].Metadata.UID != uidb {
		t.Fatalf("ListParticipations ascending uid = %#v", parts)
	}

	providerScope := &apimeta.ScopeRef{TypedRef: apimeta.TypedRef{
		APIVersion: model.APIVersionCloudProvider, Kind: string(apimeta.ScopeCloudProvider),
		Name: "prov-a", UID: uida,
	}}
	hlB := model.HostingLocation{
		Metadata: apimeta.ObjectMeta{Name: "loc-b", UID: uidb, ScopeRef: providerScope},
		Spec:     model.HostingLocationSpec{CountryCode: "US", Locality: "Austin"},
	}
	hlA := hlB
	hlA.Metadata.Name = "loc-a"
	hlA.Metadata.UID = uida
	if _, prob := s.CreateHostingLocation(hlB); prob != nil {
		t.Fatalf("create location b: %#v", prob)
	}
	if _, prob := s.CreateHostingLocation(hlA); prob != nil {
		t.Fatalf("create location a: %#v", prob)
	}
	locations := s.ListHostingLocations()
	if len(locations) != 2 || locations[0].Metadata.UID != uida || locations[1].Metadata.UID != uidb {
		t.Fatalf("ListHostingLocations ascending uid = %#v", locations)
	}

	dcB := model.Datacenter{
		Metadata: apimeta.ObjectMeta{Name: "dc-b", UID: uidb, ScopeRef: providerScope},
		Spec: model.DatacenterSpec{HostingLocationRef: apimeta.TypedRef{
			APIVersion: model.APIVersionHostingLocation, Kind: model.KindHostingLocation, Name: "loc-a", UID: uida,
		}},
	}
	dcA := dcB
	dcA.Metadata.Name = "dc-a"
	dcA.Metadata.UID = uida
	if _, prob := s.CreateDatacenter(dcB); prob != nil {
		t.Fatalf("create datacenter b: %#v", prob)
	}
	if _, prob := s.CreateDatacenter(dcA); prob != nil {
		t.Fatalf("create datacenter a: %#v", prob)
	}
	datacenters := s.ListDatacenters()
	if len(datacenters) != 2 || datacenters[0].Metadata.UID != uida || datacenters[1].Metadata.UID != uidb {
		t.Fatalf("ListDatacenters ascending uid = %#v", datacenters)
	}

	fdB := model.FaultDomain{
		Metadata: apimeta.ObjectMeta{Name: "fd-b", UID: uidb, ScopeRef: providerScope},
		Spec: model.FaultDomainSpec{DatacenterRef: apimeta.TypedRef{
			APIVersion: model.APIVersionDatacenter, Kind: model.KindDatacenter, Name: "dc-a", UID: uida,
		}},
	}
	fdA := fdB
	fdA.Metadata.Name = "fd-a"
	fdA.Metadata.UID = uida
	if _, prob := s.CreateFaultDomain(fdB); prob != nil {
		t.Fatalf("create fault domain b: %#v", prob)
	}
	if _, prob := s.CreateFaultDomain(fdA); prob != nil {
		t.Fatalf("create fault domain a: %#v", prob)
	}
	faultDomains := s.ListFaultDomains()
	if len(faultDomains) != 2 || faultDomains[0].Metadata.UID != uida || faultDomains[1].Metadata.UID != uidb {
		t.Fatalf("ListFaultDomains ascending uid = %#v", faultDomains)
	}

	stB := model.InfrastructureStack{
		Metadata: apimeta.ObjectMeta{Name: "stack-b", UID: uidb, ScopeRef: providerScope},
		Spec: model.InfrastructureStackSpec{FaultDomainRef: apimeta.TypedRef{
			APIVersion: model.APIVersionFaultDomain, Kind: model.KindFaultDomain, Name: "fd-a", UID: uida,
		}},
	}
	stA := stB
	stA.Metadata.Name = "stack-a"
	stA.Metadata.UID = uida
	if _, prob := s.CreateInfrastructureStack(stB); prob != nil {
		t.Fatalf("create stack b: %#v", prob)
	}
	if _, prob := s.CreateInfrastructureStack(stA); prob != nil {
		t.Fatalf("create stack a: %#v", prob)
	}
	stacks := s.ListInfrastructureStacks()
	if len(stacks) != 2 || stacks[0].Metadata.UID != uida || stacks[1].Metadata.UID != uidb {
		t.Fatalf("ListInfrastructureStacks ascending uid = %#v", stacks)
	}
}

func assertEmptyLists(t *testing.T, s *cloudmodel.Store) {
	t.Helper()
	if got := s.ListCloudPlatforms(); got == nil || len(got) != 0 {
		t.Fatalf("ListCloudPlatforms empty-safe = %#v", got)
	}
	if got := s.ListCloudProviders(); got == nil || len(got) != 0 {
		t.Fatalf("ListCloudProviders empty-safe = %#v", got)
	}
	if got := s.ListParticipations(); got == nil || len(got) != 0 {
		t.Fatalf("ListParticipations empty-safe = %#v", got)
	}
	if got := s.ListHostingLocations(); got == nil || len(got) != 0 {
		t.Fatalf("ListHostingLocations empty-safe = %#v", got)
	}
	if got := s.ListDatacenters(); got == nil || len(got) != 0 {
		t.Fatalf("ListDatacenters empty-safe = %#v", got)
	}
	if got := s.ListFaultDomains(); got == nil || len(got) != 0 {
		t.Fatalf("ListFaultDomains empty-safe = %#v", got)
	}
	if got := s.ListInfrastructureStacks(); got == nil || len(got) != 0 {
		t.Fatalf("ListInfrastructureStacks empty-safe = %#v", got)
	}
}

func TestStore_TopologyLookupUnderPublicationLock(t *testing.T) {
	t.Parallel()
	s := cloudmodel.NewStore()
	providerUID := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	scope := &apimeta.ScopeRef{TypedRef: apimeta.TypedRef{
		APIVersion: model.APIVersionCloudProvider, Kind: string(apimeta.ScopeCloudProvider),
		Name: "prov", UID: providerUID,
	}}
	hlUID := "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
	dcUID := "cccccccccccccccccccccccccccccccc"
	fdUID := "dddddddddddddddddddddddddddddddd"
	stUID := "eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee"

	if _, prob := s.CreateHostingLocation(model.HostingLocation{
		Metadata: apimeta.ObjectMeta{Name: "loc", UID: hlUID, ScopeRef: scope, ResourceVersion: "1"},
		Spec:     model.HostingLocationSpec{CountryCode: "US", Locality: "Austin"},
	}); prob != nil {
		t.Fatalf("create HL: %#v", prob)
	}
	if _, prob := s.CreateDatacenter(model.Datacenter{
		Metadata: apimeta.ObjectMeta{Name: "dc", UID: dcUID, ScopeRef: scope, ResourceVersion: "1"},
		Spec: model.DatacenterSpec{HostingLocationRef: apimeta.TypedRef{
			APIVersion: model.APIVersionHostingLocation, Kind: model.KindHostingLocation, Name: "loc", UID: hlUID,
		}},
	}); prob != nil {
		t.Fatalf("create DC: %#v", prob)
	}
	if _, prob := s.CreateFaultDomain(model.FaultDomain{
		Metadata: apimeta.ObjectMeta{Name: "fd", UID: fdUID, ScopeRef: scope, ResourceVersion: "1"},
		Spec: model.FaultDomainSpec{DatacenterRef: apimeta.TypedRef{
			APIVersion: model.APIVersionDatacenter, Kind: model.KindDatacenter, Name: "dc", UID: dcUID,
		}},
	}); prob != nil {
		t.Fatalf("create FD: %#v", prob)
	}
	if _, prob := s.CreateInfrastructureStack(model.InfrastructureStack{
		Metadata: apimeta.ObjectMeta{Name: "stack", UID: stUID, ScopeRef: scope, ResourceVersion: "1"},
		Spec: model.InfrastructureStackSpec{FaultDomainRef: apimeta.TypedRef{
			APIVersion: model.APIVersionFaultDomain, Kind: model.KindFaultDomain, Name: "fd", UID: fdUID,
		}},
	}); prob != nil {
		t.Fatalf("create stack: %#v", prob)
	}

	s.BeginPublication()
	hl, ok := s.LookupHostingLocation(hlUID)
	if !ok || hl.Metadata.UID != hlUID || hl.Metadata.Name != "loc" {
		s.EndPublication()
		t.Fatalf("LookupHostingLocation=%#v ok=%v", hl, ok)
	}
	hl.Metadata.Name = "mutated"
	if got, _ := s.LookupHostingLocation(hlUID); got.Metadata.Name != "loc" {
		s.EndPublication()
		t.Fatal("LookupHostingLocation must return a copy")
	}

	dc, ok := s.LookupDatacenter(dcUID)
	if !ok || dc.Metadata.UID != dcUID {
		s.EndPublication()
		t.Fatalf("LookupDatacenter=%#v ok=%v", dc, ok)
	}
	fd, ok := s.LookupFaultDomain(fdUID)
	if !ok || fd.Metadata.UID != fdUID {
		s.EndPublication()
		t.Fatalf("LookupFaultDomain=%#v ok=%v", fd, ok)
	}
	st, ok := s.LookupInfrastructureStack(stUID)
	if !ok || st.Metadata.UID != stUID {
		s.EndPublication()
		t.Fatalf("LookupInfrastructureStack=%#v ok=%v", st, ok)
	}
	s.EndPublication()

	// Public Get… methods remain for unlocked reads and acquire the Store lock.
	if got, ok := s.GetHostingLocation(hlUID); !ok || got.Metadata.UID != hlUID {
		t.Fatalf("GetHostingLocation unlocked read failed: %#v ok=%v", got, ok)
	}
}

func TestStore_PairedBackingRead_NonInterleaving(t *testing.T) {
	t.Parallel()
	s := cloudmodel.NewStore()
	platformUID := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	providerUID := "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
	partUID := "cccccccccccccccccccccccccccccccc"
	stackUID := "dddddddddddddddddddddddddddddddd"

	if _, prob := s.CreateCloudPlatform(platformScoped("plat", platformUID)); prob != nil {
		t.Fatalf("platform: %#v", prob)
	}
	if _, prob := s.CreateCloudProvider(providerScoped("prov", providerUID)); prob != nil {
		t.Fatalf("provider: %#v", prob)
	}
	part := model.CloudProviderParticipation{
		Metadata: apimeta.ObjectMeta{
			Name: "part", UID: partUID, ResourceVersion: "1", Generation: 1,
			ScopeRef: &apimeta.ScopeRef{TypedRef: apimeta.TypedRef{
				APIVersion: model.APIVersionCloudPlatform, Kind: string(apimeta.ScopeCloudPlatform),
				Name: "plat", UID: platformUID,
			}},
		},
		Spec: model.CloudProviderParticipationSpec{
			CloudPlatformRef: apimeta.TypedRef{UID: platformUID},
			CloudProviderRef: apimeta.TypedRef{UID: providerUID},
		},
		Status: model.CloudProviderParticipationStatus{Phase: model.ParticipationPhaseActive},
	}
	if _, prob := s.CreateParticipation(part); prob != nil {
		t.Fatalf("participation: %#v", prob)
	}
	scope := &apimeta.ScopeRef{TypedRef: apimeta.TypedRef{
		APIVersion: model.APIVersionCloudProvider, Kind: string(apimeta.ScopeCloudProvider),
		Name: "prov", UID: providerUID,
	}}
	stack := model.InfrastructureStack{
		Metadata: apimeta.ObjectMeta{Name: "stack", UID: stackUID, ScopeRef: scope, ResourceVersion: "1", Generation: 3},
		Status:   model.InfrastructureStackStatus{Phase: model.InfrastructureStackPhaseActive},
	}
	if _, prob := s.CreateInfrastructureStack(stack); prob != nil {
		t.Fatalf("stack: %#v", prob)
	}

	entered := make(chan struct{})
	release := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		s.WithPairedBackingRead(partUID, stackUID, func(snap cloudmodel.PairedBackingSnapshot) {
			if !snap.ParticipationOK || !snap.StackOK {
				t.Errorf("paired snapshot missing resources")
			}
			if snap.Participation.Status.Phase != model.ParticipationPhaseActive {
				t.Errorf("unexpected participation phase=%s", snap.Participation.Status.Phase)
			}
			if snap.Stack.Metadata.Generation != 3 {
				t.Errorf("unexpected stack generation=%d", snap.Stack.Metadata.Generation)
			}
			close(entered)
			<-release
			if snap.Participation.Status.Phase != model.ParticipationPhaseActive {
				t.Errorf("lease view mutated while held")
			}
			if snap.Stack.Metadata.Generation != 3 {
				t.Errorf("lease stack generation mutated while held")
			}
		})
	}()
	<-entered

	writeDone := make(chan struct{})
	go func() {
		cur, ok := s.GetParticipation(partUID)
		if !ok {
			t.Errorf("participation vanished")
			close(writeDone)
			return
		}
		cur.Status.Phase = model.ParticipationPhaseSuspended
		if _, prob := s.UpdateParticipation(cur); prob != nil {
			t.Errorf("update while lease held: %#v", prob)
		}
		close(writeDone)
	}()

	select {
	case <-writeDone:
		t.Fatal("FEATURE-0015 write interleaved through paired backing-read lease")
	case <-time.After(50 * time.Millisecond):
	}
	close(release)
	<-writeDone
	wg.Wait()

	after, ok := s.GetParticipation(partUID)
	if !ok || after.Status.Phase != model.ParticipationPhaseSuspended {
		t.Fatalf("write must apply after lease release: %#v ok=%v", after, ok)
	}
}

func TestStore_PairedBackingRead_MissingCopies(t *testing.T) {
	t.Parallel()
	s := cloudmodel.NewStore()
	var got cloudmodel.PairedBackingSnapshot
	s.WithPairedBackingRead("missing-part", "missing-stack", func(snap cloudmodel.PairedBackingSnapshot) {
		got = snap
	})
	if got.ParticipationOK || got.StackOK {
		t.Fatalf("missing resources must report !OK: %#v", got)
	}
}
