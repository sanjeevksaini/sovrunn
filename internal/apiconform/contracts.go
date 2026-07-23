package apiconform

import (
	"github.com/sanjeevksaini/sovrunn/internal/apicond"
	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
)

// Conformance-only Go contract types for the eight canonical schemas
// (D-01b, D-17, Matrix D). These prove schema fit via TypeBinding checks;
// they do NOT implement domain behavior, runtime services, or later-phase
// execution (F12-FIXTURE-001, F12-FIXTURE-002, F12-NAMING-005).
//
// JSON tags and required/optional encoding follow VerifyGoTypeAgainstSchema
// rules: required fields are non-pointer without omitempty; optional fields
// use omitempty or pointers. TypeMeta is anonymously embedded so
// apiVersion/kind promote to the top-level JSON object.

// ---------------------------------------------------------------------------
// Project — ManagedResource / customer-facing / Tenant
// ---------------------------------------------------------------------------

// ProjectPhase is the closed Project status.phase vocabulary.
type ProjectPhase string

const (
	ProjectPhasePending  ProjectPhase = "Pending"
	ProjectPhaseActive   ProjectPhase = "Active"
	ProjectPhaseInactive ProjectPhase = "Inactive"
	ProjectPhaseDeleting ProjectPhase = "Deleting"
	ProjectPhaseFailed   ProjectPhase = "Failed"
)

// Project is the conformance-only Go type for api/schemas/project.json.
type Project struct {
	apimeta.TypeMeta                    // anonymous embed promotes apiVersion/kind (F12-NAMING-002)
	Metadata         apimeta.ObjectMeta `json:"metadata"`
	Spec             ProjectSpec        `json:"spec"`
	Status           ProjectStatus      `json:"status,omitempty"`
}

// ProjectSpec is the customer-authored desired state for Project.
type ProjectSpec struct {
	Description string `json:"description,omitempty"`
}

// ProjectStatus is the system-owned observed state for Project.
type ProjectStatus struct {
	Phase              ProjectPhase        `json:"phase,omitempty"`
	ObservedGeneration int64               `json:"observedGeneration,omitempty"`
	Conditions         []apicond.Condition `json:"conditions,omitempty"`
}

// ---------------------------------------------------------------------------
// ResourcePool — ManagedResource / operator-facing / Provider
// ---------------------------------------------------------------------------

// ResourcePoolPhase is the closed ResourcePool status.phase vocabulary.
type ResourcePoolPhase string

const (
	ResourcePoolPhasePending     ResourcePoolPhase = "Pending"
	ResourcePoolPhaseReady       ResourcePoolPhase = "Ready"
	ResourcePoolPhaseDegraded    ResourcePoolPhase = "Degraded"
	ResourcePoolPhaseUnavailable ResourcePoolPhase = "Unavailable"
	ResourcePoolPhaseDeleting    ResourcePoolPhase = "Deleting"
	ResourcePoolPhaseFailed      ResourcePoolPhase = "Failed"
)

// ResourcePool is the conformance-only Go type for api/schemas/resource-pool.json.
type ResourcePool struct {
	apimeta.TypeMeta
	Metadata apimeta.ObjectMeta `json:"metadata"`
	Spec     ResourcePoolSpec   `json:"spec"`
	Status   ResourcePoolStatus `json:"status,omitempty"`
}

// ResourcePoolSpec is the provider-neutral pool declaration.
type ResourcePoolSpec struct {
	CapabilityClass  string `json:"capabilityClass"`
	JurisdictionCode string `json:"jurisdictionCode,omitempty"`
}

// ResourcePoolStatus is the system-owned observed pool state.
type ResourcePoolStatus struct {
	Phase              ResourcePoolPhase   `json:"phase,omitempty"`
	ObservedGeneration int64               `json:"observedGeneration,omitempty"`
	Conditions         []apicond.Condition `json:"conditions,omitempty"`
}

// ---------------------------------------------------------------------------
// DiscoveredDatabase — ObservedExternalResource / adapter-facing / Provider
// ---------------------------------------------------------------------------

// ObservationState is the closed observation freshness/existence vocabulary.
type ObservationState string

const (
	ObservationStateCurrent ObservationState = "Current"
	ObservationStateStale   ObservationState = "Stale"
	ObservationStateUnknown ObservationState = "Unknown"
	ObservationStateAbsent  ObservationState = "Absent"
)

// FreshnessState is the closed freshness classification vocabulary.
type FreshnessState string

const (
	FreshnessStateFresh   FreshnessState = "Fresh"
	FreshnessStateStale   FreshnessState = "Stale"
	FreshnessStateUnknown FreshnessState = "Unknown"
)

// DiscoveredDatabase is the conformance-only Go type for
// api/schemas/discovered-database.json.
type DiscoveredDatabase struct {
	apimeta.TypeMeta
	Metadata   apimeta.ObjectMeta           `json:"metadata"`
	Status     DiscoveredDatabaseStatus     `json:"status"`
	Provenance DiscoveredDatabaseProvenance `json:"provenance"`
	Freshness  DiscoveredDatabaseFreshness  `json:"freshness"`
}

// DiscoveredDatabaseStatus is the normalized observation state.
type DiscoveredDatabaseStatus struct {
	ObservationState ObservationState    `json:"observationState"`
	ExternalName     string              `json:"externalName,omitempty"`
	Conditions       []apicond.Condition `json:"conditions,omitempty"`
}

// DiscoveredDatabaseProvenance is required source provenance for observations.
type DiscoveredDatabaseProvenance struct {
	SourceRef  apimeta.TypedRef `json:"sourceRef"`
	ObservedAt string           `json:"observedAt"`
}

// DiscoveredDatabaseFreshness is the required freshness contract.
type DiscoveredDatabaseFreshness struct {
	TTLSeconds     int64          `json:"ttlSeconds,omitempty"`
	FreshnessState FreshnessState `json:"freshnessState"`
}

// ---------------------------------------------------------------------------
// PluginDefinition — VersionedDefinition / plugin-facing / Platform
// ---------------------------------------------------------------------------

// PublicationState is the closed plugin publication lifecycle vocabulary.
type PublicationState string

const (
	PublicationStateDraft      PublicationState = "Draft"
	PublicationStatePublished  PublicationState = "Published"
	PublicationStateSuperseded PublicationState = "Superseded"
)

// PluginDefinitionPhase is the closed PluginDefinition status.phase vocabulary.
type PluginDefinitionPhase string

const (
	PluginDefinitionPhasePending  PluginDefinitionPhase = "Pending"
	PluginDefinitionPhaseAccepted PluginDefinitionPhase = "Accepted"
	PluginDefinitionPhaseRejected PluginDefinitionPhase = "Rejected"
)

// PluginDefinition is the conformance-only Go type for
// api/schemas/plugin-definition.json.
type PluginDefinition struct {
	apimeta.TypeMeta
	Metadata apimeta.ObjectMeta     `json:"metadata"`
	Spec     PluginDefinitionSpec   `json:"spec"`
	Status   PluginDefinitionStatus `json:"status,omitempty"`
}

// PluginDefinitionSpec is the published plugin contract payload.
type PluginDefinitionSpec struct {
	Version            string           `json:"version"`
	CompatibilityRange string           `json:"compatibilityRange"`
	PublicationState   PublicationState `json:"publicationState"`
}

// PluginDefinitionStatus is optional system-owned definition acceptance state.
type PluginDefinitionStatus struct {
	Phase      PluginDefinitionPhase `json:"phase,omitempty"`
	Conditions []apicond.Condition   `json:"conditions,omitempty"`
}

// ---------------------------------------------------------------------------
// AdapterConfiguration — ManagedResource / adapter-facing / Provider
// ---------------------------------------------------------------------------

// AdapterConfigurationPhase is the closed AdapterConfiguration status.phase
// vocabulary.
type AdapterConfigurationPhase string

const (
	AdapterConfigurationPhasePending  AdapterConfigurationPhase = "Pending"
	AdapterConfigurationPhaseReady    AdapterConfigurationPhase = "Ready"
	AdapterConfigurationPhaseDegraded AdapterConfigurationPhase = "Degraded"
	AdapterConfigurationPhaseFailed   AdapterConfigurationPhase = "Failed"
	AdapterConfigurationPhaseDeleting AdapterConfigurationPhase = "Deleting"
)

// AdapterConfiguration is the conformance-only Go type for
// api/schemas/adapter-configuration.json.
type AdapterConfiguration struct {
	apimeta.TypeMeta
	Metadata apimeta.ObjectMeta         `json:"metadata"`
	Spec     AdapterConfigurationSpec   `json:"spec"`
	Status   AdapterConfigurationStatus `json:"status,omitempty"`
}

// AdapterConfigurationSpec is desired adapter configuration with
// secret-reference isolation (raw credentials MUST NOT appear here).
type AdapterConfigurationSpec struct {
	AdapterClass         string            `json:"adapterClass"`
	CredentialsSecretRef apimeta.TypedRef  `json:"credentialsSecretRef"`
	NativeConfigRef      *apimeta.TypedRef `json:"nativeConfigRef,omitempty"`
}

// AdapterConfigurationStatus is system-owned adapter configuration observation.
type AdapterConfigurationStatus struct {
	Phase              AdapterConfigurationPhase `json:"phase,omitempty"`
	ObservedGeneration int64                     `json:"observedGeneration,omitempty"`
	Conditions         []apicond.Condition       `json:"conditions,omitempty"`
}

// ---------------------------------------------------------------------------
// PlacementEvaluationRequest — TransientRequestResult / internal-engine / Project
// ---------------------------------------------------------------------------

// PlacementOutcome is the closed placement evaluation outcome vocabulary.
type PlacementOutcome string

const (
	PlacementOutcomeAllow          PlacementOutcome = "Allow"
	PlacementOutcomeDeny           PlacementOutcome = "Deny"
	PlacementOutcomeReviewRequired PlacementOutcome = "ReviewRequired"
)

// PlacementEvaluationRequest is the conformance-only Go type for
// api/schemas/placement-evaluation-request.json.
type PlacementEvaluationRequest struct {
	apimeta.TypeMeta
	Metadata apimeta.ObjectMeta             `json:"metadata,omitempty"`
	Request  PlacementEvaluationRequestBody `json:"request"`
	Result   *PlacementEvaluationResult     `json:"result,omitempty"`
}

// PlacementEvaluationRequestBody is the typed placement evaluation input.
type PlacementEvaluationRequestBody struct {
	SubjectRef              apimeta.TypedRef `json:"subjectRef"`
	RequiredCapabilityClass string           `json:"requiredCapabilityClass"`
	JurisdictionCode        string           `json:"jurisdictionCode,omitempty"`
}

// PlacementEvaluationResult is the typed evaluation result.
type PlacementEvaluationResult struct {
	Outcome         PlacementOutcome  `json:"outcome"`
	ReasonCode      string            `json:"reasonCode"`
	SelectedPoolRef *apimeta.TypedRef `json:"selectedPoolRef,omitempty"`
}

// ---------------------------------------------------------------------------
// Operation — LongRunningOperation / plugin-facing / six Matrix B scopes (D-17)
// ---------------------------------------------------------------------------

// OperationPhase is the closed Operation status.phase vocabulary.
type OperationPhase string

const (
	OperationPhasePending   OperationPhase = "Pending"
	OperationPhaseRunning   OperationPhase = "Running"
	OperationPhaseSucceeded OperationPhase = "Succeeded"
	OperationPhaseFailed    OperationPhase = "Failed"
	OperationPhaseCancelled OperationPhase = "Cancelled"
)

// TriState is the closed True/False/Unknown vocabulary used by Operation
// status.retryable.
type TriState string

const (
	TriStateTrue    TriState = "True"
	TriStateFalse   TriState = "False"
	TriStateUnknown TriState = "Unknown"
)

// BoolString is the closed True/False vocabulary used by Operation
// status.cancelRequested.
type BoolString string

const (
	BoolStringTrue  BoolString = "True"
	BoolStringFalse BoolString = "False"
)

// Operation is the conformance-only Go type for api/schemas/operation.json.
//
// D-17 fields:
//   - Spec.TargetRef — operation target; governance scope must match Metadata.ScopeRef
//   - Metadata.ScopeRef — canonical nil for Platform targets; otherwise target scope UID
//   - OwnerRef — optional lifecycle containment only; MUST NOT replace scopeRef
type Operation struct {
	apimeta.TypeMeta
	Metadata apimeta.ObjectMeta `json:"metadata"`
	OwnerRef *apimeta.OwnerRef  `json:"ownerRef,omitempty"`
	Spec     OperationSpec      `json:"spec"`
	Status   OperationStatus    `json:"status,omitempty"`
}

// OperationSpec is the immutable request payload after acceptance.
type OperationSpec struct {
	TargetRef      apimeta.TypedRef  `json:"targetRef"`
	Action         string            `json:"action"`
	RequesterRef   *apimeta.TypedRef `json:"requesterRef,omitempty"`
	IdempotencyKey string            `json:"idempotencyKey,omitempty"`
	RequestID      string            `json:"requestId,omitempty"`
}

// OperationStatus is executor-owned progress, retryability, and terminal result.
type OperationStatus struct {
	Phase              OperationPhase      `json:"phase,omitempty"`
	ProgressPercent    int64               `json:"progressPercent,omitempty"`
	Retryable          TriState            `json:"retryable,omitempty"`
	CancelRequested    BoolString          `json:"cancelRequested,omitempty"`
	TerminalCode       string              `json:"terminalCode,omitempty"`
	TerminalMessage    string              `json:"terminalMessage,omitempty"`
	ObservedGeneration int64               `json:"observedGeneration,omitempty"`
	Conditions         []apicond.Condition `json:"conditions,omitempty"`
}

// ---------------------------------------------------------------------------
// AuditEvent — ImmutableRecord / governance-only / Organization
// ---------------------------------------------------------------------------

// AuditOutcome is the closed audit outcome vocabulary.
type AuditOutcome string

const (
	AuditOutcomeSucceeded AuditOutcome = "Succeeded"
	AuditOutcomeDenied    AuditOutcome = "Denied"
	AuditOutcomeFailed    AuditOutcome = "Failed"
)

// AuditEvent is the conformance-only Go type for api/schemas/audit-event.json.
type AuditEvent struct {
	apimeta.TypeMeta
	Metadata apimeta.ObjectMeta `json:"metadata"`
	Record   AuditEventRecord   `json:"record"`
}

// AuditEventRecord is the append-only audit payload.
type AuditEventRecord struct {
	ActorRef               apimeta.TypedRef  `json:"actorRef"`
	RequestID              string            `json:"requestId"`
	OperationRef           *apimeta.TypedRef `json:"operationRef,omitempty"`
	SubjectRef             apimeta.TypedRef  `json:"subjectRef"`
	SubjectResourceVersion string            `json:"subjectResourceVersion,omitempty"`
	Action                 string            `json:"action"`
	Outcome                AuditOutcome      `json:"outcome"`
	ReasonCode             string            `json:"reasonCode"`
	CorrectionOfRef        *apimeta.TypedRef `json:"correctionOfRef,omitempty"`
}
