package executiontarget

import (
	"sort"
	"strconv"
	"sync"

	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/executiontarget/model"
)

// Store is the FEATURE-0016 process-local registry for ExecutionTarget,
// NormalizedTargetFactSet, and TargetQualificationResult (design §3.1).
//
// Mutation helpers are package-private staging/index primitives. They do not
// publish independently; TASK-F16-04 invokes them only while the lifecycle
// service mutex is held (via beginPublication/endPublication).
type Store struct {
	mu sync.Mutex

	targets  map[string]model.ExecutionTarget           // uid -> resource
	factSets map[string]model.NormalizedTargetFactSet   // uid -> record
	results  map[string]model.TargetQualificationResult // uid -> record

	// maintenanceMarkers keys are target UID.
	maintenanceMarkers map[string]model.CurrentMaintenanceMarker

	// nameIndex keys are cloudProviderScopeUID + "\x00" + name -> uid.
	// Reservations are retained after retirement for the process lifetime.
	nameIndex map[string]string
	// tupleIndex keys are participationUID + "\x00" + stackUID + "\x00" +
	// targetClass -> uid for live non-Retired targets only.
	tupleIndex map[string]string
}

// NewStore returns an empty in-memory Store.
func NewStore() *Store {
	return &Store{
		targets:            make(map[string]model.ExecutionTarget),
		factSets:           make(map[string]model.NormalizedTargetFactSet),
		results:            make(map[string]model.TargetQualificationResult),
		maintenanceMarkers: make(map[string]model.CurrentMaintenanceMarker),
		nameIndex:          make(map[string]string),
		tupleIndex:         make(map[string]string),
	}
}

// stagedOutcome is a mutation prepared under the publication lock but not yet
// published. Lifecycle commit paths append AuditEvent before publish; on append
// failure they must abort so neither the resource nor indexes become visible.
type stagedOutcome struct {
	store  *Store
	apply  func()
	active bool
}

// beginPublication acquires the publication lock. Callers must pair it with
// endPublication (typically via defer).
func (s *Store) beginPublication() {
	s.mu.Lock()
}

// endPublication releases the publication lock.
func (s *Store) endPublication() {
	s.mu.Unlock()
}

// publish applies a staged mutation. The publication lock must be held.
func (s *Store) publish(staged *stagedOutcome) {
	if staged == nil || !staged.active || staged.store != s {
		return
	}
	staged.apply()
	staged.active = false
	staged.apply = nil
}

// abort discards a staged mutation without publishing. The publication lock
// must be held.
func (s *Store) abort(staged *stagedOutcome) {
	if staged == nil {
		return
	}
	staged.active = false
	staged.apply = nil
}

func (s *Store) stage(apply func()) *stagedOutcome {
	return &stagedOutcome{store: s, apply: apply, active: true}
}

func nameKey(scopeUID, name string) string {
	return scopeUID + "\x00" + name
}

func tupleKey(participationUID, stackUID string, targetClass model.TargetClass) string {
	return participationUID + "\x00" + stackUID + "\x00" + string(targetClass)
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
		Message: "resource name or live tuple already exists in scope",
	}})
}

// stageCreateExecutionTarget rechecks scope-unique name reservation and
// live-tuple uniqueness, then stages a create. The publication lock must be
// held. Staging alone does not publish.
func (s *Store) stageCreateExecutionTarget(et model.ExecutionTarget) (*stagedOutcome, *apiproblem.Problem) {
	if et.Metadata.UID == "" {
		return nil, apiproblem.New(apiproblem.CodeInternalError).WithDetail("uid is required before staging")
	}
	if et.CloudProviderScopeUID == "" {
		return nil, apiproblem.New(apiproblem.CodeInternalError).WithDetail("derived cloudProvider scope is required before staging")
	}
	if _, exists := s.targets[et.Metadata.UID]; exists {
		return nil, alreadyExists("/metadata/uid")
	}
	nkey := nameKey(et.CloudProviderScopeUID, et.Metadata.Name)
	if _, exists := s.nameIndex[nkey]; exists {
		return nil, alreadyExists("/metadata/name")
	}
	tkey := tupleKey(
		et.Spec.CloudProviderParticipationRef.UID,
		et.Spec.InfrastructureStackRef.UID,
		et.Spec.TargetClass,
	)
	if _, exists := s.tupleIndex[tkey]; exists {
		return nil, alreadyExists("/spec")
	}
	if et.Metadata.ResourceVersion == "" {
		et.Metadata.ResourceVersion = "1"
	}
	staged := et
	return s.stage(func() {
		s.targets[staged.Metadata.UID] = staged
		s.nameIndex[nkey] = staged.Metadata.UID
		if staged.Status.Lifecycle != model.LifecycleRetired {
			s.tupleIndex[tkey] = staged.Metadata.UID
		}
	}), nil
}

// stageUpdateExecutionTarget stages a target mutation and maintains the
// live-tuple index across Active/Retired transitions. Name reservations are
// never released. The publication lock must be held.
func (s *Store) stageUpdateExecutionTarget(et model.ExecutionTarget) (*stagedOutcome, *apiproblem.Problem) {
	cur, ok := s.targets[et.Metadata.UID]
	if !ok {
		return nil, apiproblem.New(apiproblem.CodeResourceNotFound)
	}
	et.Metadata.ResourceVersion = nextResourceVersion(cur.Metadata.ResourceVersion)
	// Preserve derived scope and immutable name/spec identity for indexes.
	et.CloudProviderScopeUID = cur.CloudProviderScopeUID
	et.Metadata.Name = cur.Metadata.Name
	et.Spec = cur.Spec

	oldKey := tupleKey(
		cur.Spec.CloudProviderParticipationRef.UID,
		cur.Spec.InfrastructureStackRef.UID,
		cur.Spec.TargetClass,
	)
	newKey := tupleKey(
		et.Spec.CloudProviderParticipationRef.UID,
		et.Spec.InfrastructureStackRef.UID,
		et.Spec.TargetClass,
	)
	staged := et
	return s.stage(func() {
		s.targets[staged.Metadata.UID] = staged
		if cur.Status.Lifecycle != model.LifecycleRetired {
			if s.tupleIndex[oldKey] == staged.Metadata.UID {
				delete(s.tupleIndex, oldKey)
			}
		}
		if staged.Status.Lifecycle != model.LifecycleRetired {
			s.tupleIndex[newKey] = staged.Metadata.UID
		}
	}), nil
}

// stagePutFactSet stages insertion or replacement of an internal FactSet
// record. The publication lock must be held.
func (s *Store) stagePutFactSet(fs model.NormalizedTargetFactSet) (*stagedOutcome, *apiproblem.Problem) {
	if fs.UID == "" {
		return nil, apiproblem.New(apiproblem.CodeInternalError).WithDetail("fact set uid is required before staging")
	}
	staged := fs
	return s.stage(func() {
		s.factSets[staged.UID] = staged
	}), nil
}

// stagePutQualificationResult stages insertion or replacement of an internal
// TargetQualificationResult record. The publication lock must be held.
func (s *Store) stagePutQualificationResult(r model.TargetQualificationResult) (*stagedOutcome, *apiproblem.Problem) {
	if r.UID == "" {
		return nil, apiproblem.New(apiproblem.CodeInternalError).WithDetail("qualification result uid is required before staging")
	}
	staged := r
	return s.stage(func() {
		s.results[staged.UID] = staged
	}), nil
}

// stageDeleteFactSet stages removal of a FactSet by UID. The publication lock
// must be held.
func (s *Store) stageDeleteFactSet(uid string) *stagedOutcome {
	return s.stage(func() {
		delete(s.factSets, uid)
	})
}

// stageDeleteQualificationResult stages removal of a Result by UID. The
// publication lock must be held.
func (s *Store) stageDeleteQualificationResult(uid string) *stagedOutcome {
	return s.stage(func() {
		delete(s.results, uid)
	})
}

// stageSetMaintenanceMarker stages the current-Maintenance marker for a
// target. The publication lock must be held.
func (s *Store) stageSetMaintenanceMarker(m model.CurrentMaintenanceMarker) (*stagedOutcome, *apiproblem.Problem) {
	if m.TargetUID == "" {
		return nil, apiproblem.New(apiproblem.CodeInternalError).WithDetail("maintenance marker target uid is required")
	}
	staged := m
	return s.stage(func() {
		s.maintenanceMarkers[staged.TargetUID] = staged
	}), nil
}

// stageClearMaintenanceMarker stages removal of the current-Maintenance
// marker for a target. The publication lock must be held.
func (s *Store) stageClearMaintenanceMarker(targetUID string) *stagedOutcome {
	return s.stage(func() {
		delete(s.maintenanceMarkers, targetUID)
	})
}

// lookupExecutionTarget returns a copy by UID. The publication lock must be held.
func (s *Store) lookupExecutionTarget(uid string) (model.ExecutionTarget, bool) {
	et, ok := s.targets[uid]
	return et, ok
}

// lookupFactSet returns a copy by UID. The publication lock must be held.
func (s *Store) lookupFactSet(uid string) (model.NormalizedTargetFactSet, bool) {
	fs, ok := s.factSets[uid]
	return fs, ok
}

// lookupQualificationResult returns a copy by UID. The publication lock must be held.
func (s *Store) lookupQualificationResult(uid string) (model.TargetQualificationResult, bool) {
	r, ok := s.results[uid]
	return r, ok
}

// lookupMaintenanceMarker returns a copy by target UID. The publication lock
// must be held.
func (s *Store) lookupMaintenanceMarker(targetUID string) (model.CurrentMaintenanceMarker, bool) {
	m, ok := s.maintenanceMarkers[targetUID]
	return m, ok
}

// nameReserved reports whether the scope/name reservation exists. The
// publication lock must be held.
func (s *Store) nameReserved(scopeUID, name string) bool {
	_, ok := s.nameIndex[nameKey(scopeUID, name)]
	return ok
}

// liveTupleOccupied reports whether a live non-Retired tuple is indexed. The
// publication lock must be held.
func (s *Store) liveTupleOccupied(participationUID, stackUID string, targetClass model.TargetClass) bool {
	_, ok := s.tupleIndex[tupleKey(participationUID, stackUID, targetClass)]
	return ok
}

// GetExecutionTarget returns a copy by UID.
func (s *Store) GetExecutionTarget(uid string) (model.ExecutionTarget, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	et, ok := s.targets[uid]
	return et, ok
}

// GetFactSet returns a copy by UID.
func (s *Store) GetFactSet(uid string) (model.NormalizedTargetFactSet, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	fs, ok := s.factSets[uid]
	return fs, ok
}

// GetQualificationResult returns a copy by UID.
func (s *Store) GetQualificationResult(uid string) (model.TargetQualificationResult, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	r, ok := s.results[uid]
	return r, ok
}

// GetMaintenanceMarker returns a copy by target UID.
func (s *Store) GetMaintenanceMarker(targetUID string) (model.CurrentMaintenanceMarker, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	m, ok := s.maintenanceMarkers[targetUID]
	return m, ok
}

// ListExecutionTargets returns an independent snapshot ordered by ascending
// metadata.uid. Empty stores return a non-nil empty slice.
func (s *Store) ListExecutionTargets() []model.ExecutionTarget {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]model.ExecutionTarget, 0, len(s.targets))
	for _, et := range s.targets {
		out = append(out, et)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].Metadata.UID < out[j].Metadata.UID
	})
	return out
}
