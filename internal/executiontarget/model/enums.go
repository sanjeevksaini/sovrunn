package model

// LifecycleState is the persisted ExecutionTarget status.lifecycle vocabulary
// (VS0-STATE-004; Active/Retired only). Qualifying is never a persisted state.
type LifecycleState string

const (
	LifecycleActive  LifecycleState = "Active"
	LifecycleRetired LifecycleState = "Retired"
)

// QualificationState is the persisted ExecutionTarget status.qualification
// vocabulary (VS0-STATE-004). Qualifying is a lifecycle-service-only
// in-flight reservation and is never persisted or projected.
type QualificationState string

const (
	QualificationUnqualified   QualificationState = "Unqualified"
	QualificationQualified     QualificationState = "Qualified"
	QualificationRejected      QualificationState = "Rejected"
	QualificationIndeterminate QualificationState = "Indeterminate"
)

// EffectiveAvailability is the response-only projection vocabulary
// (Available/Unavailable/Maintenance). It is never persisted on status.
type EffectiveAvailability string

const (
	EffectiveAvailable   EffectiveAvailability = "Available"
	EffectiveUnavailable EffectiveAvailability = "Unavailable"
	EffectiveMaintenance EffectiveAvailability = "Maintenance"
)

// FactTruth is the closed truth vocabulary for NormalizedTargetFactSet facts
// (Supported/Unsupported/Unknown).
type FactTruth string

const (
	FactSupported   FactTruth = "Supported"
	FactUnsupported FactTruth = "Unsupported"
	FactUnknown     FactTruth = "Unknown"
)

// QualificationOutcome is the TargetQualificationResult outcome vocabulary
// (Qualified/Rejected/Indeterminate).
type QualificationOutcome string

const (
	OutcomeQualified     QualificationOutcome = "Qualified"
	OutcomeRejected      QualificationOutcome = "Rejected"
	OutcomeIndeterminate QualificationOutcome = "Indeterminate"
)

// TargetClass is the immutable ExecutionTarget spec.targetClass vocabulary.
type TargetClass string

const (
	TargetClassSyntheticIaaS TargetClass = "synthetic-iaas"
)

// Fixed provenance and observer identity constants (design §4.4).
const (
	ObserverID                    = "sovrunn.synthetic-iaas-observer/v1"
	ProfileVersion                = "synthetic-iaas/v1"
	FactVersionV1                 = "v1"
	FactSchemaVersionV1           = "v1" // observer provenance; not a second FactSet field
	APIVersionExecutionTarget     = "execution.sovrunn.io/v1alpha1"
	KindExecutionTarget           = "ExecutionTarget"
	APIVersionFactSet             = "adapter.sovrunn.io/v1alpha1"
	KindNormalizedTargetFactSet   = "NormalizedTargetFactSet"
	APIVersionQualificationResult = "adapter.sovrunn.io/v1alpha1"
	KindTargetQualificationResult = "TargetQualificationResult"
)
