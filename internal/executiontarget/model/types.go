package model

import "github.com/sanjeevksaini/sovrunn/internal/apimeta"

// ExecutionTarget is the FEATURE-0016 CloudProvider-scoped resource
// (VS0-SCHEMA-015; design §4.1). Scope derives solely from
// spec.cloudProviderParticipationRef; no scopeRef is persisted or projected.
type ExecutionTarget struct {
	apimeta.TypeMeta                       // anonymous embed promotes apiVersion/kind
	Metadata         apimeta.ObjectMeta    `json:"metadata"`
	Spec             ExecutionTargetSpec   `json:"spec"`
	Status           ExecutionTargetStatus `json:"status,omitempty"`
	// CloudProviderScopeUID is the derived CloudProvider scope used for
	// process-local name uniqueness. It is never serialized as scopeRef and
	// is never projected (ADH-2026-058 clause 1).
	CloudProviderScopeUID string `json:"-"`
}

// ExecutionTargetSpec carries the closed immutable create-time fields.
type ExecutionTargetSpec struct {
	CloudProviderParticipationRef apimeta.TypedRef `json:"cloudProviderParticipationRef"`
	InfrastructureStackRef        apimeta.TypedRef `json:"infrastructureStackRef"`
	TargetClass                   TargetClass      `json:"targetClass"`
}

// ExecutionTargetStatus is api-server-owned observed state. There is no
// status.availability field; effectiveAvailability is response-only.
type ExecutionTargetStatus struct {
	Lifecycle              LifecycleState     `json:"lifecycle,omitempty"`
	Qualification          QualificationState `json:"qualification,omitempty"`
	MaintenanceEpoch       int64              `json:"maintenanceEpoch"`
	ObservedGeneration     int64              `json:"observedGeneration,omitempty"`
	FactSetRef             *apimeta.TypedRef  `json:"factSetRef,omitempty"`
	QualificationResultRef *apimeta.TypedRef  `json:"qualificationResultRef,omitempty"`
}

// NormalizedTargetFactSet is the internal-only observation record
// (VS0-SCHEMA-016; design §4.2). It is never projected.
type NormalizedTargetFactSet struct {
	UID                           string               `json:"uid"`
	TargetRef                     apimeta.TypedRef     `json:"targetRef"`
	InfrastructureStackGeneration int64                `json:"infrastructureStackGeneration"`
	MaintenanceEpoch              int64                `json:"maintenanceEpoch"`
	ViabilityFingerprint          ViabilityFingerprint `json:"viabilityFingerprint"`
	ObserverID                    string               `json:"observerID"`
	ObserverRevision              string               `json:"observerRevision"`
	FactVersion                   string               `json:"factVersion"`
	ObservedAt                    string               `json:"observedAt"`
	ExpiresAt                     string               `json:"expiresAt"`
	Facts                         FactSetTruths        `json:"facts"`
}

// FactSetTruths carries the four named fact truth values
// (compute.vm, storage.block, storage.object, network.private).
type FactSetTruths struct {
	Compute FactComputeTruths `json:"compute"`
	Storage FactStorageTruths `json:"storage"`
	Network FactNetworkTruths `json:"network"`
}

// FactComputeTruths holds compute.vm.
type FactComputeTruths struct {
	VM FactTruth `json:"vm"`
}

// FactStorageTruths holds storage.block and storage.object.
type FactStorageTruths struct {
	Block  FactTruth `json:"block"`
	Object FactTruth `json:"object"`
}

// FactNetworkTruths holds network.private.
type FactNetworkTruths struct {
	Private FactTruth `json:"private"`
}

// TargetQualificationResult is the internal-only qualification conclusion
// (VS0-SCHEMA-017; design §4.3). viabilityFingerprint is FactSet-only and is
// not duplicated here.
type TargetQualificationResult struct {
	UID                           string               `json:"uid"`
	TargetRef                     apimeta.TypedRef     `json:"targetRef"`
	FactSetRef                    apimeta.TypedRef     `json:"factSetRef"`
	InfrastructureStackGeneration int64                `json:"infrastructureStackGeneration"`
	MaintenanceEpoch              int64                `json:"maintenanceEpoch"`
	ProfileVersion                string               `json:"profileVersion"`
	Outcome                       QualificationOutcome `json:"outcome"`
	ReasonCodes                   []string             `json:"reasonCodes"`
	EvaluatedAt                   string               `json:"evaluatedAt"`
}

// CurrentMaintenanceMarker is the lifecycle-service-owned, target-bound,
// in-process Maintenance marker (ADH-2026-058 clause 7). It is never
// client-writable or projected and clears on process restart.
type CurrentMaintenanceMarker struct {
	TargetUID        string
	MaintenanceEpoch int64
	Active           bool
}

// ViabilityFingerprint is the opaque viability value composed only of
// participation UID/effective-active state and stack UID/phase
// (requirements glossary; design §4.2).
type ViabilityFingerprint string

// ComputeViabilityFingerprint builds the closed viability fingerprint value
// from participation and stack viability inputs only.
func ComputeViabilityFingerprint(participationUID string, participationEffectiveActive bool, stackUID, stackPhase string) ViabilityFingerprint {
	active := "inactive"
	if participationEffectiveActive {
		active = "active"
	}
	return ViabilityFingerprint(participationUID + "\x00" + active + "\x00" + stackUID + "\x00" + stackPhase)
}

// String returns the fingerprint opaque value.
func (f ViabilityFingerprint) String() string {
	return string(f)
}

// NewCreateProposal initializes a lifecycle-owned create proposal with
// derived CloudProvider scope, Active/Unqualified, maintenanceEpoch=0, and
// observedGeneration equal to the referenced InfrastructureStack generation
// (VS0-SCHEMA-015 serverAssigned; design §4.1).
//
// Callers must supply a non-empty uid and cloudProviderScopeUID. ScopeRef is
// intentionally left unset.
func NewCreateProposal(
	name, uid, cloudProviderScopeUID string,
	participationRef, stackRef apimeta.TypedRef,
	stackGeneration int64,
) ExecutionTarget {
	return ExecutionTarget{
		TypeMeta: apimeta.TypeMeta{
			APIVersion: APIVersionExecutionTarget,
			Kind:       KindExecutionTarget,
		},
		Metadata: apimeta.ObjectMeta{
			Name:            name,
			UID:             uid,
			Generation:      1,
			ResourceVersion: "1",
		},
		Spec: ExecutionTargetSpec{
			CloudProviderParticipationRef: participationRef,
			InfrastructureStackRef:        stackRef,
			TargetClass:                   TargetClassSyntheticIaaS,
		},
		Status: ExecutionTargetStatus{
			Lifecycle:          LifecycleActive,
			Qualification:      QualificationUnqualified,
			MaintenanceEpoch:   0,
			ObservedGeneration: stackGeneration,
		},
		CloudProviderScopeUID: cloudProviderScopeUID,
	}
}
