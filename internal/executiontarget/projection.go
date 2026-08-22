package executiontarget

import (
	"time"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/executiontarget/model"
)

// ProjectionSnapshot is a coherent immutable view of every input required to
// compute the closed safe ExecutionTarget projection and response-only
// effectiveAvailability (design §4.1, §4.6). It is captured under the DD-02
// lifecycle mutex; Project never re-reads store maps, markers, facts, or
// viability.
type ProjectionSnapshot struct {
	// Target is a copy of the committed ExecutionTarget. Internal links and
	// non-projected status fields may be present on Target; Project omits them.
	Target model.ExecutionTarget

	// MaintenanceActive is true when the current-Maintenance marker for this
	// target is active at capture time.
	MaintenanceActive bool

	// FactsFresh is true when the target has a current FactSet link, the
	// FactSet record exists, and facts have not expired (expiresAt > now and
	// not marked expired by the injected-clock expiry worker).
	FactsFresh bool

	// BackingViable is true when the already-resolved FEATURE-0015 backing
	// snapshot is effective-Active participation with an Active stack.
	BackingViable bool
}

// SafeExecutionTarget is the closed safe JSON projection for create/item
// GET/qualify/retire/LIST (design §4.1). It contains exactly the approved
// fields; no public conditions; facts, results, observer identity, handles,
// and internal record links are omitted.
type SafeExecutionTarget struct {
	Metadata              SafeMetadata                `json:"metadata"`
	Spec                  SafeSpec                    `json:"spec"`
	Status                SafeStatus                  `json:"status"`
	EffectiveAvailability model.EffectiveAvailability `json:"effectiveAvailability"`
}

// SafeMetadata is the closed projected metadata subset.
type SafeMetadata struct {
	UID             string `json:"uid"`
	Name            string `json:"name"`
	Generation      int64  `json:"generation"`
	ResourceVersion string `json:"resourceVersion"`
}

// SafeSpec is the closed projected immutable spec subset.
type SafeSpec struct {
	CloudProviderParticipationRef apimeta.TypedRef  `json:"cloudProviderParticipationRef"`
	InfrastructureStackRef        apimeta.TypedRef  `json:"infrastructureStackRef"`
	TargetClass                   model.TargetClass `json:"targetClass"`
}

// SafeStatus is the closed projected status subset (lifecycle + qualification
// only). There is no status.availability or conditions member.
type SafeStatus struct {
	Lifecycle     model.LifecycleState     `json:"lifecycle"`
	Qualification model.QualificationState `json:"qualification"`
}

// CaptureProjectionSnapshot returns a coherent immutable snapshot for
// targetUID under the DD-02 mutex. backing is the already-resolved
// FEATURE-0015 viability input (never read from F0016 store maps). ok is
// false when the target does not exist.
func (s *ExecutionTargetLifecycleService) CaptureProjectionSnapshot(
	targetUID string,
	backing BackingViability,
) (ProjectionSnapshot, bool) {
	s.lock()
	defer s.unlock()

	et, ok := s.store.lookupExecutionTarget(targetUID)
	if !ok {
		return ProjectionSnapshot{}, false
	}

	maintenanceActive := false
	if marker, has := s.store.lookupMaintenanceMarker(targetUID); has && marker.Active {
		maintenanceActive = true
	}

	return ProjectionSnapshot{
		Target:            et,
		MaintenanceActive: maintenanceActive,
		FactsFresh:        s.factsFreshLocked(et),
		BackingViable:     backingViable(backing),
	}, true
}

// Project builds the closed safe ExecutionTarget JSON projection from a
// coherent snapshot. It never reads mutable store maps, markers, facts, or
// viability and never includes facts/results/observer identity/handles/
// internal record links in its output.
func Project(snap ProjectionSnapshot) SafeExecutionTarget {
	return SafeExecutionTarget{
		Metadata: SafeMetadata{
			UID:             snap.Target.Metadata.UID,
			Name:            snap.Target.Metadata.Name,
			Generation:      snap.Target.Metadata.Generation,
			ResourceVersion: snap.Target.Metadata.ResourceVersion,
		},
		Spec: SafeSpec{
			CloudProviderParticipationRef: snap.Target.Spec.CloudProviderParticipationRef,
			InfrastructureStackRef:        snap.Target.Spec.InfrastructureStackRef,
			TargetClass:                   snap.Target.Spec.TargetClass,
		},
		Status: SafeStatus{
			Lifecycle:     snap.Target.Status.Lifecycle,
			Qualification: snap.Target.Status.Qualification,
		},
		EffectiveAvailability: ComputeEffectiveAvailability(snap),
	}
}

// ComputeEffectiveAvailability applies the ordered response-only truth table
// (design §4.6):
//  1. Retired → Unavailable
//  2. else Active with current Maintenance marker active → Maintenance
//  3. else only Active + Qualified + fresh current facts + viable backing → Available
//  4. every other Active combination → Unavailable
func ComputeEffectiveAvailability(snap ProjectionSnapshot) model.EffectiveAvailability {
	if snap.Target.Status.Lifecycle == model.LifecycleRetired {
		return model.EffectiveUnavailable
	}
	// Remaining cases are Active (Qualifying is never persisted or projected).
	if snap.MaintenanceActive {
		return model.EffectiveMaintenance
	}
	if snap.Target.Status.Lifecycle == model.LifecycleActive &&
		snap.Target.Status.Qualification == model.QualificationQualified &&
		snap.FactsFresh &&
		snap.BackingViable {
		return model.EffectiveAvailable
	}
	return model.EffectiveUnavailable
}

// factsFreshLocked reports whether the target has fresh current facts.
// Caller must hold the DD-02 mutex (and store publication lock).
func (s *ExecutionTargetLifecycleService) factsFreshLocked(et model.ExecutionTarget) bool {
	if et.Status.FactSetRef == nil || et.Status.FactSetRef.UID == "" {
		return false
	}
	fsUID := et.Status.FactSetRef.UID
	fs, ok := s.store.lookupFactSet(fsUID)
	if !ok {
		return false
	}
	if st := s.factExpiry[fsUID]; st != nil && st.expired {
		return false
	}
	exp, err := time.Parse(time.RFC3339, fs.ExpiresAt)
	if err != nil {
		return false
	}
	return exp.After(s.now())
}

// backingViable reports whether FEATURE-0015 backing is viable for Available
// projection (participation effective-Active and stack phase Active).
func backingViable(b BackingViability) bool {
	return b.ParticipationEffectiveActive && b.StackActive()
}
