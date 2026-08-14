package cloudmodel

import (
	"sort"
	"strconv"
	"sync"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/cloudmodel/model"
	"github.com/sanjeevksaini/sovrunn/internal/cloudmodel/validate"
)

// Store is the FEATURE-0015 in-memory registry for the seven owned kinds
// (design DD-05). It enforces name uniqueness within server-derived scope,
// non-terminal participation pair uniqueness, resourceVersion increments,
// and publication lock/staging primitives. Task 8 owns AuditEvent
// append-before-publication orchestration on top of these primitives.
type Store struct {
	mu sync.Mutex

	platforms      map[string]model.CloudPlatform              // uid -> resource
	providers      map[string]model.CloudProvider              // uid -> resource
	participations map[string]model.CloudProviderParticipation // uid -> resource
	locations      map[string]model.HostingLocation            // uid -> resource
	datacenters    map[string]model.Datacenter                 // uid -> resource
	faultDomains   map[string]model.FaultDomain                // uid -> resource
	stacks         map[string]model.InfrastructureStack        // uid -> resource

	// nameIndex keys are scopeUID + "\x00" + kind + "\x00" + name -> uid
	nameIndex map[string]string
	// pairIndex keys are platformUID + "\x00" + providerUID -> participation uid
	// for non-terminal participations only.
	pairIndex map[string]string
}

// NewStore returns an empty in-memory Store.
func NewStore() *Store {
	return &Store{
		platforms:      make(map[string]model.CloudPlatform),
		providers:      make(map[string]model.CloudProvider),
		participations: make(map[string]model.CloudProviderParticipation),
		locations:      make(map[string]model.HostingLocation),
		datacenters:    make(map[string]model.Datacenter),
		faultDomains:   make(map[string]model.FaultDomain),
		stacks:         make(map[string]model.InfrastructureStack),
		nameIndex:      make(map[string]string),
		pairIndex:      make(map[string]string),
	}
}

// StagedOutcome is a mutation prepared under the publication lock but not yet
// published. Task 8 appends AuditEvent before Publish; on append failure it
// must Abort so neither the resource nor an idempotency completion is visible.
type StagedOutcome struct {
	store  *Store
	apply  func()
	active bool
}

// BeginPublication acquires the publication lock. Callers must pair it with
// EndPublication (typically via defer).
func (s *Store) BeginPublication() {
	s.mu.Lock()
}

// EndPublication releases the publication lock.
func (s *Store) EndPublication() {
	s.mu.Unlock()
}

// Publish applies a staged mutation. The publication lock must be held.
func (s *Store) Publish(staged *StagedOutcome) {
	if staged == nil || !staged.active || staged.store != s {
		return
	}
	staged.apply()
	staged.active = false
	staged.apply = nil
}

// Abort discards a staged mutation without publishing. The publication lock
// must be held.
func (s *Store) Abort(staged *StagedOutcome) {
	if staged == nil {
		return
	}
	staged.active = false
	staged.apply = nil
}

func (s *Store) stage(apply func()) *StagedOutcome {
	return &StagedOutcome{store: s, apply: apply, active: true}
}

// ScopeUID returns the uniqueness scope key for a ScopeRef. Platform-root
// uses the reserved PlatformScopeUID sentinel.
func ScopeUID(scope *apimeta.ScopeRef) string {
	n := apimeta.NormalizeScope(scope)
	if n == nil {
		return apimeta.PlatformScopeUID
	}
	return n.UID
}

func nameKey(scopeUID, kind, name string) string {
	return scopeUID + "\x00" + kind + "\x00" + name
}

func pairKey(platformUID, providerUID string) string {
	return platformUID + "\x00" + providerUID
}

func nextResourceVersion(current string) string {
	if current == "" {
		return "1"
	}
	n, err := strconv.ParseUint(current, 10, 64)
	if err != nil {
		return "1"
	}
	return strconv.FormatUint(n+1, 10)
}

func alreadyExists(field string) *apiproblem.Problem {
	return apiproblem.New(apiproblem.CodeAlreadyExists).WithViolations([]apiproblem.Violation{{
		Field:   field,
		Code:    apiproblem.ViolationOutOfRange,
		Message: "resource name already exists in scope",
	}})
}

func participationDuplicate() *apiproblem.Problem {
	return apiproblem.New(apiproblem.CodeAlreadyExists).WithViolations([]apiproblem.Violation{{
		Field:   "/spec",
		Code:    validate.ViolationParticipationDuplicate,
		Message: "non-terminal participation already exists for this pair",
	}})
}

func isTerminalParticipation(phase model.ParticipationPhase) bool {
	switch phase {
	case model.ParticipationPhaseRejected,
		model.ParticipationPhaseWithdrawn,
		model.ParticipationPhaseExpired,
		model.ParticipationPhaseTerminated:
		return true
	default:
		return false
	}
}

// StageCreateCloudPlatform rechecks Platform-root name uniqueness and stages
// a create. The publication lock must be held.
func (s *Store) StageCreateCloudPlatform(cp model.CloudPlatform) (*StagedOutcome, *apiproblem.Problem) {
	scope := ScopeUID(cp.Metadata.ScopeRef)
	key := nameKey(scope, model.KindCloudPlatform, cp.Metadata.Name)
	if _, exists := s.nameIndex[key]; exists {
		return nil, alreadyExists("/metadata/name")
	}
	if cp.Metadata.UID == "" {
		return nil, apiproblem.New(apiproblem.CodeInternalError).WithDetail("uid is required before staging")
	}
	if _, exists := s.platforms[cp.Metadata.UID]; exists {
		return nil, alreadyExists("/metadata/uid")
	}
	cp.Metadata.ResourceVersion = "1"
	staged := cp
	return s.stage(func() {
		s.platforms[staged.Metadata.UID] = staged
		s.nameIndex[key] = staged.Metadata.UID
	}), nil
}

// StageCreateCloudProvider rechecks Platform-root name uniqueness and stages
// a create. The publication lock must be held.
func (s *Store) StageCreateCloudProvider(cp model.CloudProvider) (*StagedOutcome, *apiproblem.Problem) {
	scope := ScopeUID(cp.Metadata.ScopeRef)
	key := nameKey(scope, model.KindCloudProvider, cp.Metadata.Name)
	if _, exists := s.nameIndex[key]; exists {
		return nil, alreadyExists("/metadata/name")
	}
	if cp.Metadata.UID == "" {
		return nil, apiproblem.New(apiproblem.CodeInternalError).WithDetail("uid is required before staging")
	}
	cp.Metadata.ResourceVersion = "1"
	staged := cp
	return s.stage(func() {
		s.providers[staged.Metadata.UID] = staged
		s.nameIndex[key] = staged.Metadata.UID
	}), nil
}

// StageCreateParticipation rechecks CloudPlatform-scope name uniqueness and
// non-terminal pair uniqueness, then stages a create.
func (s *Store) StageCreateParticipation(p model.CloudProviderParticipation) (*StagedOutcome, *apiproblem.Problem) {
	scope := ScopeUID(p.Metadata.ScopeRef)
	nkey := nameKey(scope, model.KindCloudProviderParticipation, p.Metadata.Name)
	if _, exists := s.nameIndex[nkey]; exists {
		return nil, alreadyExists("/metadata/name")
	}
	pkey := pairKey(p.Spec.CloudPlatformRef.UID, p.Spec.CloudProviderRef.UID)
	if _, exists := s.pairIndex[pkey]; exists {
		return nil, participationDuplicate()
	}
	if p.Metadata.UID == "" {
		return nil, apiproblem.New(apiproblem.CodeInternalError).WithDetail("uid is required before staging")
	}
	p.Metadata.ResourceVersion = "1"
	staged := p
	return s.stage(func() {
		s.participations[staged.Metadata.UID] = staged
		s.nameIndex[nkey] = staged.Metadata.UID
		if !isTerminalParticipation(staged.Status.Phase) {
			s.pairIndex[pkey] = staged.Metadata.UID
		}
	}), nil
}

// StageCreateHostingLocation rechecks CloudProvider-scope name uniqueness.
func (s *Store) StageCreateHostingLocation(hl model.HostingLocation) (*StagedOutcome, *apiproblem.Problem) {
	return s.stageTopologyCreate(model.KindHostingLocation, ScopeUID(hl.Metadata.ScopeRef), hl.Metadata.Name, hl.Metadata.UID, func(rv string) {
		hl.Metadata.ResourceVersion = rv
		s.locations[hl.Metadata.UID] = hl
	})
}

// StageCreateDatacenter rechecks CloudProvider-scope name uniqueness.
func (s *Store) StageCreateDatacenter(dc model.Datacenter) (*StagedOutcome, *apiproblem.Problem) {
	return s.stageTopologyCreate(model.KindDatacenter, ScopeUID(dc.Metadata.ScopeRef), dc.Metadata.Name, dc.Metadata.UID, func(rv string) {
		dc.Metadata.ResourceVersion = rv
		s.datacenters[dc.Metadata.UID] = dc
	})
}

// StageCreateFaultDomain rechecks CloudProvider-scope name uniqueness.
func (s *Store) StageCreateFaultDomain(fd model.FaultDomain) (*StagedOutcome, *apiproblem.Problem) {
	return s.stageTopologyCreate(model.KindFaultDomain, ScopeUID(fd.Metadata.ScopeRef), fd.Metadata.Name, fd.Metadata.UID, func(rv string) {
		fd.Metadata.ResourceVersion = rv
		s.faultDomains[fd.Metadata.UID] = fd
	})
}

// StageCreateInfrastructureStack rechecks CloudProvider-scope name uniqueness.
func (s *Store) StageCreateInfrastructureStack(st model.InfrastructureStack) (*StagedOutcome, *apiproblem.Problem) {
	return s.stageTopologyCreate(model.KindInfrastructureStack, ScopeUID(st.Metadata.ScopeRef), st.Metadata.Name, st.Metadata.UID, func(rv string) {
		st.Metadata.ResourceVersion = rv
		s.stacks[st.Metadata.UID] = st
	})
}

func (s *Store) stageTopologyCreate(kind, scopeUID, name, uid string, put func(rv string)) (*StagedOutcome, *apiproblem.Problem) {
	key := nameKey(scopeUID, kind, name)
	if _, exists := s.nameIndex[key]; exists {
		return nil, alreadyExists("/metadata/name")
	}
	if uid == "" {
		return nil, apiproblem.New(apiproblem.CodeInternalError).WithDetail("uid is required before staging")
	}
	return s.stage(func() {
		put("1")
		s.nameIndex[key] = uid
	}), nil
}

// StageUpdateCloudPlatform stages a mutation that increments resourceVersion.
func (s *Store) StageUpdateCloudPlatform(cp model.CloudPlatform) (*StagedOutcome, *apiproblem.Problem) {
	cur, ok := s.platforms[cp.Metadata.UID]
	if !ok {
		return nil, apiproblem.New(apiproblem.CodeResourceNotFound)
	}
	cp.Metadata.ResourceVersion = nextResourceVersion(cur.Metadata.ResourceVersion)
	staged := cp
	return s.stage(func() {
		s.platforms[staged.Metadata.UID] = staged
	}), nil
}

// StageUpdateCloudProvider stages a mutation that increments resourceVersion.
func (s *Store) StageUpdateCloudProvider(cp model.CloudProvider) (*StagedOutcome, *apiproblem.Problem) {
	cur, ok := s.providers[cp.Metadata.UID]
	if !ok {
		return nil, apiproblem.New(apiproblem.CodeResourceNotFound)
	}
	cp.Metadata.ResourceVersion = nextResourceVersion(cur.Metadata.ResourceVersion)
	staged := cp
	return s.stage(func() {
		s.providers[staged.Metadata.UID] = staged
	}), nil
}

// StageUpdateParticipation stages a participation mutation, maintaining the
// non-terminal pair index.
func (s *Store) StageUpdateParticipation(p model.CloudProviderParticipation) (*StagedOutcome, *apiproblem.Problem) {
	cur, ok := s.participations[p.Metadata.UID]
	if !ok {
		return nil, apiproblem.New(apiproblem.CodeResourceNotFound)
	}
	p.Metadata.ResourceVersion = nextResourceVersion(cur.Metadata.ResourceVersion)
	staged := p
	pkey := pairKey(staged.Spec.CloudPlatformRef.UID, staged.Spec.CloudProviderRef.UID)
	return s.stage(func() {
		s.participations[staged.Metadata.UID] = staged
		if isTerminalParticipation(staged.Status.Phase) {
			if s.pairIndex[pkey] == staged.Metadata.UID {
				delete(s.pairIndex, pkey)
			}
		} else {
			s.pairIndex[pkey] = staged.Metadata.UID
		}
	}), nil
}

// StageUpdateHostingLocation stages a topology mutation.
func (s *Store) StageUpdateHostingLocation(hl model.HostingLocation) (*StagedOutcome, *apiproblem.Problem) {
	cur, ok := s.locations[hl.Metadata.UID]
	if !ok {
		return nil, apiproblem.New(apiproblem.CodeResourceNotFound)
	}
	hl.Metadata.ResourceVersion = nextResourceVersion(cur.Metadata.ResourceVersion)
	staged := hl
	return s.stage(func() { s.locations[staged.Metadata.UID] = staged }), nil
}

// StageUpdateDatacenter stages a topology mutation.
func (s *Store) StageUpdateDatacenter(dc model.Datacenter) (*StagedOutcome, *apiproblem.Problem) {
	cur, ok := s.datacenters[dc.Metadata.UID]
	if !ok {
		return nil, apiproblem.New(apiproblem.CodeResourceNotFound)
	}
	dc.Metadata.ResourceVersion = nextResourceVersion(cur.Metadata.ResourceVersion)
	staged := dc
	return s.stage(func() { s.datacenters[staged.Metadata.UID] = staged }), nil
}

// StageUpdateFaultDomain stages a topology mutation.
func (s *Store) StageUpdateFaultDomain(fd model.FaultDomain) (*StagedOutcome, *apiproblem.Problem) {
	cur, ok := s.faultDomains[fd.Metadata.UID]
	if !ok {
		return nil, apiproblem.New(apiproblem.CodeResourceNotFound)
	}
	fd.Metadata.ResourceVersion = nextResourceVersion(cur.Metadata.ResourceVersion)
	staged := fd
	return s.stage(func() { s.faultDomains[staged.Metadata.UID] = staged }), nil
}

// StageUpdateInfrastructureStack stages a topology mutation.
func (s *Store) StageUpdateInfrastructureStack(st model.InfrastructureStack) (*StagedOutcome, *apiproblem.Problem) {
	cur, ok := s.stacks[st.Metadata.UID]
	if !ok {
		return nil, apiproblem.New(apiproblem.CodeResourceNotFound)
	}
	st.Metadata.ResourceVersion = nextResourceVersion(cur.Metadata.ResourceVersion)
	staged := st
	return s.stage(func() { s.stacks[staged.Metadata.UID] = staged }), nil
}

// HasCloudPlatform reports whether any CloudPlatform is published.
// Handlers use this for the CloudProvider/topology/participation root gate;
// the Store does not make authorization decisions.
func (s *Store) HasCloudPlatform() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.platforms) > 0
}

// ListCloudPlatforms returns an independent snapshot ordered by ascending
// metadata.uid. Empty stores return a non-nil empty slice.
func (s *Store) ListCloudPlatforms() []model.CloudPlatform {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]model.CloudPlatform, 0, len(s.platforms))
	for _, cp := range s.platforms {
		out = append(out, cp)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].Metadata.UID < out[j].Metadata.UID
	})
	return out
}

// ListCloudProviders returns an independent snapshot ordered by ascending
// metadata.uid. Empty stores return a non-nil empty slice.
func (s *Store) ListCloudProviders() []model.CloudProvider {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]model.CloudProvider, 0, len(s.providers))
	for _, cp := range s.providers {
		out = append(out, cp)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].Metadata.UID < out[j].Metadata.UID
	})
	return out
}

// ListParticipations returns an independent snapshot ordered by ascending
// metadata.uid. Empty stores return a non-nil empty slice.
func (s *Store) ListParticipations() []model.CloudProviderParticipation {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]model.CloudProviderParticipation, 0, len(s.participations))
	for _, p := range s.participations {
		out = append(out, p)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].Metadata.UID < out[j].Metadata.UID
	})
	return out
}

// ListHostingLocations returns an independent snapshot ordered by ascending
// metadata.uid. Empty stores return a non-nil empty slice.
func (s *Store) ListHostingLocations() []model.HostingLocation {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]model.HostingLocation, 0, len(s.locations))
	for _, hl := range s.locations {
		out = append(out, hl)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].Metadata.UID < out[j].Metadata.UID
	})
	return out
}

// ListDatacenters returns an independent snapshot ordered by ascending
// metadata.uid. Empty stores return a non-nil empty slice.
func (s *Store) ListDatacenters() []model.Datacenter {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]model.Datacenter, 0, len(s.datacenters))
	for _, dc := range s.datacenters {
		out = append(out, dc)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].Metadata.UID < out[j].Metadata.UID
	})
	return out
}

// ListFaultDomains returns an independent snapshot ordered by ascending
// metadata.uid. Empty stores return a non-nil empty slice.
func (s *Store) ListFaultDomains() []model.FaultDomain {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]model.FaultDomain, 0, len(s.faultDomains))
	for _, fd := range s.faultDomains {
		out = append(out, fd)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].Metadata.UID < out[j].Metadata.UID
	})
	return out
}

// ListInfrastructureStacks returns an independent snapshot ordered by ascending
// metadata.uid. Empty stores return a non-nil empty slice.
func (s *Store) ListInfrastructureStacks() []model.InfrastructureStack {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]model.InfrastructureStack, 0, len(s.stacks))
	for _, st := range s.stacks {
		out = append(out, st)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].Metadata.UID < out[j].Metadata.UID
	})
	return out
}

// GetCloudPlatform returns a copy by UID.
func (s *Store) GetCloudPlatform(uid string) (model.CloudPlatform, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	cp, ok := s.platforms[uid]
	return cp, ok
}

// LookupCloudPlatform returns a copy by UID. The publication lock must be held.
func (s *Store) LookupCloudPlatform(uid string) (model.CloudPlatform, bool) {
	cp, ok := s.platforms[uid]
	return cp, ok
}

// LookupCloudProvider returns a copy by UID. The publication lock must be held.
func (s *Store) LookupCloudProvider(uid string) (model.CloudProvider, bool) {
	cp, ok := s.providers[uid]
	return cp, ok
}

// GetCloudProvider returns a copy by UID.
func (s *Store) GetCloudProvider(uid string) (model.CloudProvider, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	cp, ok := s.providers[uid]
	return cp, ok
}

// GetParticipation returns a copy by UID.
func (s *Store) GetParticipation(uid string) (model.CloudProviderParticipation, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	p, ok := s.participations[uid]
	return p, ok
}

// GetHostingLocation returns a copy by UID.
func (s *Store) GetHostingLocation(uid string) (model.HostingLocation, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	hl, ok := s.locations[uid]
	return hl, ok
}

// LookupHostingLocation returns a copy by UID. The publication lock must be held.
func (s *Store) LookupHostingLocation(uid string) (model.HostingLocation, bool) {
	hl, ok := s.locations[uid]
	return hl, ok
}

// GetDatacenter returns a copy by UID.
func (s *Store) GetDatacenter(uid string) (model.Datacenter, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	dc, ok := s.datacenters[uid]
	return dc, ok
}

// LookupDatacenter returns a copy by UID. The publication lock must be held.
func (s *Store) LookupDatacenter(uid string) (model.Datacenter, bool) {
	dc, ok := s.datacenters[uid]
	return dc, ok
}

// GetFaultDomain returns a copy by UID.
func (s *Store) GetFaultDomain(uid string) (model.FaultDomain, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	fd, ok := s.faultDomains[uid]
	return fd, ok
}

// LookupFaultDomain returns a copy by UID. The publication lock must be held.
func (s *Store) LookupFaultDomain(uid string) (model.FaultDomain, bool) {
	fd, ok := s.faultDomains[uid]
	return fd, ok
}

// GetInfrastructureStack returns a copy by UID.
func (s *Store) GetInfrastructureStack(uid string) (model.InfrastructureStack, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	st, ok := s.stacks[uid]
	return st, ok
}

// LookupInfrastructureStack returns a copy by UID. The publication lock must be held.
func (s *Store) LookupInfrastructureStack(uid string) (model.InfrastructureStack, bool) {
	st, ok := s.stacks[uid]
	return st, ok
}

// CreateCloudPlatform is a convenience helper that stages and publishes under
// the publication lock (used by store tests and later coordinators).
func (s *Store) CreateCloudPlatform(cp model.CloudPlatform) (model.CloudPlatform, *apiproblem.Problem) {
	s.BeginPublication()
	defer s.EndPublication()
	staged, prob := s.StageCreateCloudPlatform(cp)
	if prob != nil {
		return model.CloudPlatform{}, prob
	}
	s.Publish(staged)
	return s.platforms[cp.Metadata.UID], nil
}

// CreateCloudProvider stages and publishes under the publication lock.
func (s *Store) CreateCloudProvider(cp model.CloudProvider) (model.CloudProvider, *apiproblem.Problem) {
	s.BeginPublication()
	defer s.EndPublication()
	staged, prob := s.StageCreateCloudProvider(cp)
	if prob != nil {
		return model.CloudProvider{}, prob
	}
	s.Publish(staged)
	return s.providers[cp.Metadata.UID], nil
}

// CreateParticipation stages and publishes under the publication lock.
func (s *Store) CreateParticipation(p model.CloudProviderParticipation) (model.CloudProviderParticipation, *apiproblem.Problem) {
	s.BeginPublication()
	defer s.EndPublication()
	staged, prob := s.StageCreateParticipation(p)
	if prob != nil {
		return model.CloudProviderParticipation{}, prob
	}
	s.Publish(staged)
	return s.participations[p.Metadata.UID], nil
}

// CreateHostingLocation stages and publishes under the publication lock.
func (s *Store) CreateHostingLocation(hl model.HostingLocation) (model.HostingLocation, *apiproblem.Problem) {
	s.BeginPublication()
	defer s.EndPublication()
	staged, prob := s.StageCreateHostingLocation(hl)
	if prob != nil {
		return model.HostingLocation{}, prob
	}
	s.Publish(staged)
	return s.locations[hl.Metadata.UID], nil
}

// CreateDatacenter stages and publishes under the publication lock.
func (s *Store) CreateDatacenter(dc model.Datacenter) (model.Datacenter, *apiproblem.Problem) {
	s.BeginPublication()
	defer s.EndPublication()
	staged, prob := s.StageCreateDatacenter(dc)
	if prob != nil {
		return model.Datacenter{}, prob
	}
	s.Publish(staged)
	return s.datacenters[dc.Metadata.UID], nil
}

// CreateFaultDomain stages and publishes under the publication lock.
func (s *Store) CreateFaultDomain(fd model.FaultDomain) (model.FaultDomain, *apiproblem.Problem) {
	s.BeginPublication()
	defer s.EndPublication()
	staged, prob := s.StageCreateFaultDomain(fd)
	if prob != nil {
		return model.FaultDomain{}, prob
	}
	s.Publish(staged)
	return s.faultDomains[fd.Metadata.UID], nil
}

// CreateInfrastructureStack stages and publishes under the publication lock.
func (s *Store) CreateInfrastructureStack(st model.InfrastructureStack) (model.InfrastructureStack, *apiproblem.Problem) {
	s.BeginPublication()
	defer s.EndPublication()
	staged, prob := s.StageCreateInfrastructureStack(st)
	if prob != nil {
		return model.InfrastructureStack{}, prob
	}
	s.Publish(staged)
	return s.stacks[st.Metadata.UID], nil
}

// UpdateCloudPlatform stages and publishes a mutation under the lock.
func (s *Store) UpdateCloudPlatform(cp model.CloudPlatform) (model.CloudPlatform, *apiproblem.Problem) {
	s.BeginPublication()
	defer s.EndPublication()
	staged, prob := s.StageUpdateCloudPlatform(cp)
	if prob != nil {
		return model.CloudPlatform{}, prob
	}
	s.Publish(staged)
	return s.platforms[cp.Metadata.UID], nil
}

// UpdateParticipation stages and publishes a mutation under the lock.
func (s *Store) UpdateParticipation(p model.CloudProviderParticipation) (model.CloudProviderParticipation, *apiproblem.Problem) {
	s.BeginPublication()
	defer s.EndPublication()
	staged, prob := s.StageUpdateParticipation(p)
	if prob != nil {
		return model.CloudProviderParticipation{}, prob
	}
	s.Publish(staged)
	return s.participations[p.Metadata.UID], nil
}
