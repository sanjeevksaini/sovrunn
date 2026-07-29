package decision

import (
	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
)

// Canonical DecisionProfile TypeMeta values (F13-OBJ-003; AD-029).
// DecisionProfile uses the FEATURE-0012 VersionedDefinition shape; definition
// content lives in Spec.
const (
	APIVersionDecisionProfile = "governance.sovrunn.io/v1alpha1"
	KindDecisionProfile       = "DecisionProfile"
)

// FailPosture is the profile failure-posture vocabulary required by
// F13-PROF-002 / architecture §6.8. Security-class profiles must use
// FAIL_CLOSED or REQUIRES_APPROVAL; FAIL_OPEN is valid only when the
// governing profile carries a valid SecurityExceptionRef.
type FailPosture string

const (
	FailPostureFailClosed       FailPosture = "FAIL_CLOSED"
	FailPostureRequiresApproval FailPosture = "REQUIRES_APPROVAL"
	FailPostureFailOpen         FailPosture = "FAIL_OPEN"
)

// AllFailPostures returns the FailPosture set in stable order.
func AllFailPostures() []FailPosture {
	return []FailPosture{
		FailPostureFailClosed,
		FailPostureRequiresApproval,
		FailPostureFailOpen,
	}
}

// Valid reports whether p is one of the three architecture-mandated postures.
func (p FailPosture) Valid() bool {
	switch p {
	case FailPostureFailClosed, FailPostureRequiresApproval, FailPostureFailOpen:
		return true
	default:
		return false
	}
}

// DecisionProfile is the versioned decision-family contract
// (VersionedDefinition; design §7.3; F13-OBJ-003, F13-PROF-001).
// metadata.scopeRef remains the sole scope authority for the profile object;
// ExceptionScopeRef inside SecurityExceptionRef is exception-evidence only.
type DecisionProfile struct {
	apimeta.TypeMeta                    // apiVersion=governance.sovrunn.io/v1alpha1, kind=DecisionProfile
	Metadata         apimeta.ObjectMeta `json:"metadata"`
	Spec             ProfileSpec        `json:"spec"`
}

// ProfileSpec is the DecisionProfile definition payload (design §7.3).
// Fields reproduce F13-PROF-001 declaration requirements; no approval-workflow
// or exception-issuance field is introduced.
type ProfileSpec struct {
	Name            string `json:"name"`            // stable profile id; required non-empty
	Family          string `json:"family"`          // decision family; required non-empty
	Version         string `json:"version"`         // semver; required; (name,version) composite key
	Owner           string `json:"owner"`           // profile owner identity; required non-empty
	LifecycleStatus string `json:"lifecycleStatus"` // profile lifecycle; inactive rejected by validators

	PrimaryForm            DecisionForm        `json:"primaryForm"`
	SecondaryFacets        []DecisionForm      `json:"secondaryFacets,omitempty"`
	AllowedAuthorityLevels []DecisionAuthority `json:"allowedAuthorityLevels"`

	InputSchemaRef       string `json:"inputSchemaRef"`
	TypedResultSchemaRef string `json:"typedResultSchemaRef"`
	RationaleSchemaRef   string `json:"rationaleSchemaRef"`
	ObligationSchemaRef  string `json:"obligationSchemaRef"`

	AcceptedEvaluationTypes       []AcceptedEvaluationType `json:"acceptedEvaluationTypes,omitempty"`
	AcceptedCompositionStrategies []string                 `json:"acceptedCompositionStrategies,omitempty"`

	ScopeSemantics    ProfileScopeSemantics    `json:"scopeSemantics"`
	ActorSemantics    ProfileActorSemantics    `json:"actorSemantics"`
	SubjectSemantics  ProfileSubjectSemantics  `json:"subjectSemantics"`
	PolicySemantics   ProfilePolicySemantics   `json:"policySemantics"`
	EvidenceSemantics ProfileEvidenceSemantics `json:"evidenceSemantics"`
	AuditSemantics    ProfileAuditSemantics    `json:"auditSemantics"`

	DeterminismBehavior DeterminismBehavior `json:"determinismBehavior"`
	FailureBehavior     FailureBehavior     `json:"failureBehavior"`

	ValidityRules ValidityRules `json:"validityRules"`

	SensitivityCeiling Sensitivity       `json:"sensitivityCeiling"`
	ContentCategories  []string          `json:"contentCategories,omitempty"`
	ResidencyProfile   *ResidencyProfile `json:"residencyProfile,omitempty"`
	RetentionProfile   *RetentionProfile `json:"retentionProfile,omitempty"`
	LegalHoldProfile   *LegalHoldProfile `json:"legalHoldProfile,omitempty"`
	ProjectionRules    []ProjectionRule  `json:"projectionRules,omitempty"`

	IdempotencyIdentity  string              `json:"idempotencyIdentity"`
	Limits               ProfileLimits       `json:"limits"`
	CompatibilityPolicy  CompatibilityPolicy `json:"compatibilityPolicy"`
	FixtureRefs          []string            `json:"fixtureRefs,omitempty"`
	ReassessmentTriggers []string            `json:"reassessmentTriggers,omitempty"`
}

// AcceptedEvaluationType declares one accepted evaluator type and optional
// evaluator-specific ceilings that may only narrow ProfileLimits (F13-PROF-005).
type AcceptedEvaluationType struct {
	Type    string         `json:"type"`              // registered evaluator type id
	Version string         `json:"version,omitempty"` // optional version pin
	Limits  *ProfileLimits `json:"limits,omitempty"`  // may only narrow profile ceilings
}

// ProfileLimits are profile-owned ceilings within FEATURE-0012 outer bounds
// (F13-PROF-005; design §7.3; AD-040). Missing/unknown/above-ceiling mandatory
// limits fail closed under profile validation (T-019).
type ProfileLimits struct {
	MaxNodes         int    `json:"maxNodes,omitempty"`
	MaxEdges         int    `json:"maxEdges,omitempty"`
	MaxDepth         int    `json:"maxDepth,omitempty"`
	MaxWidth         int    `json:"maxWidth,omitempty"`
	MaxFanOut        int    `json:"maxFanOut,omitempty"`
	MaxPayloadBytes  int64  `json:"maxPayloadBytes,omitempty"`
	MaxFieldCount    int    `json:"maxFieldCount,omitempty"`
	MaxCalls         int    `json:"maxCalls,omitempty"`
	MaxConcurrency   int    `json:"maxConcurrency,omitempty"`
	MaxRetries       int    `json:"maxRetries,omitempty"`
	MaxLatencyMs     int64  `json:"maxLatencyMs,omitempty"`
	MaxJurisdictions int    `json:"maxJurisdictions,omitempty"`
	MaxCostBudget    string `json:"maxCostBudget,omitempty"` // opaque structural budget tag
}

// ProfileScopeSemantics declares allowed governance scopes for decisions
// governed by this profile. It does not create a second scope authority on
// DecisionRecord/AuditEvent; metadata.scopeRef remains sole.
type ProfileScopeSemantics struct {
	AllowedScopes   []apimeta.ScopeKind `json:"allowedScopes"`
	PlatformAllowed bool                `json:"platformAllowed"`
}

// ProfileActorSemantics declares actor requirements for the profile.
type ProfileActorSemantics struct {
	RequireActor bool `json:"requireActor"`
}

// ProfileSubjectSemantics declares subject requirements for the profile.
type ProfileSubjectSemantics struct {
	RequireSubject bool     `json:"requireSubject"`
	AllowedKinds   []string `json:"allowedKinds,omitempty"`
	MaxSubjects    int      `json:"maxSubjects,omitempty"`
}

// ProfilePolicySemantics declares policy-context requirements.
type ProfilePolicySemantics struct {
	RequireEffectivePolicyContext bool `json:"requireEffectivePolicyContext"`
}

// ProfileEvidenceSemantics declares evidence requirements.
type ProfileEvidenceSemantics struct {
	RequireEvidence       bool     `json:"requireEvidence"`
	AcceptedEvidenceKinds []string `json:"acceptedEvidenceKinds,omitempty"`
}

// ProfileAuditSemantics declares audit-linkage requirements for the profile.
type ProfileAuditSemantics struct {
	RequireDurableRecord bool `json:"requireDurableRecord"`
}

// DeterminismBehavior declares determinism and AI-evaluator posture.
type DeterminismBehavior struct {
	RequireDeterministic    bool `json:"requireDeterministic"`
	AllowNondeterministicAI bool `json:"allowNondeterministicAI"`
}

// FailureBehavior declares failure, timeout, and insufficient-evidence posture
// (design §7.3/§7.7). Fail-open requires a valid SecurityExceptionRef.
type FailureBehavior struct {
	FailPosture          FailPosture           `json:"failPosture"`
	TimeoutBehavior      string                `json:"timeoutBehavior"`
	InsufficientEvidence string                `json:"insufficientEvidence"`
	DeclaredFailureModes []string              `json:"declaredFailureModes,omitempty"` // profile-declared modes for exception coverage
	SecurityExceptionRef *SecurityExceptionRef `json:"securityExceptionRef,omitempty"`
}

// ValidityRules declares correction/supersession/revocation and validity
// carriers (append-only; no in-place rewrite).
type ValidityRules struct {
	DefaultValidityMs int64 `json:"defaultValidityMs,omitempty"`
	AllowCorrection   bool  `json:"allowCorrection"`
	AllowSupersession bool  `json:"allowSupersession"`
	AllowRevocation   bool  `json:"allowRevocation"`
	MaxChainDepth     int   `json:"maxChainDepth,omitempty"`
}

// ResidencyProfile is a structural residency declaration (jurisdiction tags only).
type ResidencyProfile struct {
	AllowedJurisdictions []string `json:"allowedJurisdictions,omitempty"`
}

// RetentionProfile is a structural retention-class declaration.
type RetentionProfile struct {
	RetentionClass string `json:"retentionClass,omitempty"`
}

// LegalHoldProfile is a structural legal-hold declaration.
type LegalHoldProfile struct {
	AllowLegalHold bool `json:"allowLegalHold"`
}

// CompatibilityPolicy declares version-compatibility and ownership reassessment
// metadata for the profile (F13-PROF-001).
type CompatibilityPolicy struct {
	Owner                string `json:"owner,omitempty"`
	MinCompatibleVersion string `json:"minCompatibleVersion,omitempty"`
	BackwardCompatible   bool   `json:"backwardCompatible"`
}

// ProjectionRule is a structural audience projection declaration
// (design §7.9; F13-SEC-007). Conformance is structural only; no redaction
// execution or projected-output generation is performed.
type ProjectionRule struct {
	Audience string   `json:"audience"`           // profile-declared audience id; required
	Includes []string `json:"includes,omitempty"` // RFC 6901 pointers into the canonical view
	Excludes []string `json:"excludes,omitempty"` // RFC 6901 pointers; overlap with Includes is invalid
}

// SecurityExceptionRef is a bounded structural approval-evidence reference
// carried inside FailureBehavior (design §7.7; AD-043; F13-PROF-002).
// It is evidence only: it never approves, issues, revokes, executes, stores,
// authorizes, or workflows an exception.
//
// ExceptionScopeRef is exception-evidence scope only and never overrides,
// duplicates, replaces, or widens the canonical metadata.scopeRef.
// CompensatingControls is a bounded non-empty []string of control identifiers;
// no ControlRef or control object type is introduced.
type SecurityExceptionRef struct {
	ApprovalRef           apimeta.TypedRef  `json:"approvalRef"`                 // apiVersion/kind/name/uid; uid mandatory at validation
	ExceptionScopeRef     *apimeta.ScopeRef `json:"exceptionScopeRef,omitempty"` // exception-evidence scope only; absent = Platform where allowed
	OwnerRef              apimeta.TypedRef  `json:"ownerRef"`                    // uid required at validation
	ApprovingAuthorityRef apimeta.TypedRef  `json:"approvingAuthorityRef"`       // uid required at validation
	CompensatingControls  []string          `json:"compensatingControls"`        // required non-empty; bounded control identifiers only
	EffectiveFrom         string            `json:"effectiveFrom"`               // UTC RFC3339
	EffectiveUntil        string            `json:"effectiveUntil"`              // UTC RFC3339; after EffectiveFrom
	AuditTreatment        string            `json:"auditTreatment"`              // bounded registry identifier; structural only
	ReassessmentTrigger   string            `json:"reassessmentTrigger"`         // bounded registry identifier
	CoveredFailureModes   []string          `json:"coveredFailureModes"`         // non-empty; profile-declared failure modes only
	Purpose               string            `json:"purpose"`                     // non-empty bounded string
}
