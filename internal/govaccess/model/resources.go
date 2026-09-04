package model

import (
	"time"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
)

// AssuranceEvidence is an operation-local assurance value (REQ-F18-03).
// Raw tokens, assertions, and provider-native method payloads are never retained.
type AssuranceEvidence struct {
	PrincipalRef       PrincipalRef   `json:"principalRef"`
	AuthenticatedAt    time.Time      `json:"authenticatedAt"`
	Level              AssuranceLevel `json:"level"`
	PhishingResistant  bool           `json:"phishingResistant"`
	SourceAuthority    string         `json:"sourceAuthority"`
	IntegrityReference string         `json:"integrityReference"`
	ExpiresAt          *time.Time     `json:"expiresAt,omitempty"`
}

// AccessGroup is the FEATURE-0018 AccessGroup resource (VS0-SCHEMA-062; REQ-F18-04).
type AccessGroup struct {
	apimeta.TypeMeta
	Metadata apimeta.ObjectMeta `json:"metadata"`
	Spec     AccessGroupSpec    `json:"spec"`
	Status   AccessGroupStatus  `json:"status,omitempty"`
}

// AccessGroupSpec carries the responsible Human owner and display intent.
type AccessGroupSpec struct {
	OwnerRef    PrincipalRef `json:"ownerRef"`
	DisplayName string       `json:"displayName,omitempty"`
}

// AccessGroupStatus is system-owned lifecycle state.
type AccessGroupStatus struct {
	Phase AccessGroupPhase `json:"phase,omitempty"`
}

// Membership is the FEATURE-0018 Membership resource (VS0-SCHEMA-021; REQ-F18-05).
type Membership struct {
	apimeta.TypeMeta
	Metadata  apimeta.ObjectMeta  `json:"metadata"`
	Spec      MembershipSpec      `json:"spec"`
	Protected MembershipProtected `json:"protected"`
	Status    MembershipStatus    `json:"status,omitempty"`
}

// MembershipSpec is the immutable relationship specification.
type MembershipSpec struct {
	PrincipalRef        PrincipalRef          `json:"principalRef"`
	ContextKind         MembershipContextKind `json:"contextKind"`
	MembershipType      MembershipType        `json:"membershipType"`
	AccessGroupRef      *AccessGroupRef       `json:"accessGroupRef,omitempty"`
	ExpiresAt           *time.Time            `json:"expiresAt,omitempty"`
	SponsorRef          *PrincipalRef         `json:"sponsorRef,omitempty"`
	ResponsiblePartyRef *PrincipalRef         `json:"responsiblePartyRef,omitempty"`
}

// MembershipProtected holds system-owned provenance/freshness fields.
type MembershipProtected struct {
	SourceAuthority    string     `json:"sourceAuthority"`
	ProvisionedAt      time.Time  `json:"provisionedAt"`
	SourceObjectRef    string     `json:"sourceObjectRef,omitempty"`
	LastSynchronizedAt *time.Time `json:"lastSynchronizedAt,omitempty"`
	FreshUntil         *time.Time `json:"freshUntil,omitempty"`
}

// MembershipStatus is system-owned lifecycle state.
type MembershipStatus struct {
	Phase MembershipPhase `json:"phase,omitempty"`
}

// RoleDefinition is the FEATURE-0018 RoleDefinition resource (VS0-SCHEMA-022; REQ-F18-07).
type RoleDefinition struct {
	apimeta.TypeMeta
	Metadata apimeta.ObjectMeta `json:"metadata"`
	Spec     RoleDefinitionSpec `json:"spec"`
}

// RoleDefinitionSpec carries versioned actions and publication state.
type RoleDefinitionSpec struct {
	Version            string             `json:"version"`
	Actions            []string           `json:"actions"`
	PublicationState   PublicationState   `json:"publicationState"`
	BaseClassification RoleClassification `json:"baseClassification,omitempty"`
	SupersedesRef      *apimeta.TypedRef  `json:"supersedesRef,omitempty"`
}

// RoleAssignment is the FEATURE-0018 RoleAssignment resource (VS0-SCHEMA-023; REQ-F18-08).
type RoleAssignment struct {
	apimeta.TypeMeta
	Metadata apimeta.ObjectMeta   `json:"metadata"`
	Spec     RoleAssignmentSpec   `json:"spec"`
	Status   RoleAssignmentStatus `json:"status,omitempty"`
}

// RoleAssignmentSpec binds holder, role version, validity, and optional pins.
type RoleAssignmentSpec struct {
	RoleHolderRef         RoleHolderRef          `json:"roleHolderRef"`
	RoleDefinitionRef     apimeta.TypedRef       `json:"roleDefinitionRef"`
	RoleDefinitionVersion string                 `json:"roleDefinitionVersion"`
	Validity              AssignmentValidityMode `json:"validity"`
	NotBefore             *time.Time             `json:"notBefore,omitempty"`
	ExpiresAt             *time.Time             `json:"expiresAt,omitempty"`
	ResourceRef           *apimeta.TypedRef      `json:"resourceRef,omitempty"`
	MembershipRef         *apimeta.TypedRef      `json:"membershipRef,omitempty"`
	ResponsiblePartyRef   *PrincipalRef          `json:"responsiblePartyRef,omitempty"`
	AccessReviewRuleRef   *apimeta.TypedRef      `json:"accessReviewRuleRef,omitempty"`
}

// RoleAssignmentStatus is system-owned lifecycle/effect state.
type RoleAssignmentStatus struct {
	Phase           string     `json:"phase,omitempty"`
	NextReviewDueAt *time.Time `json:"nextReviewDueAt,omitempty"`
}

// PrivilegedAccessRequest is the FEATURE-0018 PrivilegedAccessRequest resource
// (VS0-SCHEMA-063; REQ-F18-14/15).
type PrivilegedAccessRequest struct {
	apimeta.TypeMeta
	Metadata apimeta.ObjectMeta            `json:"metadata"`
	Spec     PrivilegedAccessRequestSpec   `json:"spec"`
	Status   PrivilegedAccessRequestStatus `json:"status,omitempty"`
}

// PrivilegedAccessRequestSpec carries the immutable JIT/BreakGlass request intent.
type PrivilegedAccessRequestSpec struct {
	Mode               PrivilegedAccessMode `json:"mode"`
	RequesterRef       PrincipalRef         `json:"requesterRef"`
	RoleDefinitionRef  apimeta.TypedRef     `json:"roleDefinitionRef"`
	ActivationDeadline time.Time            `json:"activationDeadline"`
	RequestedDuration  string               `json:"requestedDuration"`
}

// PrivilegedAccessRequestStatus is system-owned LRO state.
type PrivilegedAccessRequestStatus struct {
	Phase string `json:"phase,omitempty"`
}

// AccessReview is the FEATURE-0018 AccessReview resource (VS0-SCHEMA-064; REQ-F18-16).
type AccessReview struct {
	apimeta.TypeMeta
	Metadata apimeta.ObjectMeta `json:"metadata"`
	Spec     AccessReviewSpec   `json:"spec"`
	Status   AccessReviewStatus `json:"status,omitempty"`
}

// AccessReviewSpec carries campaign configuration.
type AccessReviewSpec struct {
	AccessReviewRuleRef apimeta.TypedRef `json:"accessReviewRuleRef"`
	ReviewerEligibility []EligibilityRef `json:"reviewerEligibilityRefs"`
}

// AccessReviewStatus is system-owned LRO state.
type AccessReviewStatus struct {
	Phase string `json:"phase,omitempty"`
}

// ApprovalPolicy is the FEATURE-0018 ApprovalPolicy resource (VS0-SCHEMA-065; REQ-F18-12).
type ApprovalPolicy struct {
	apimeta.TypeMeta
	Metadata apimeta.ObjectMeta `json:"metadata"`
	Spec     ApprovalPolicySpec `json:"spec"`
}

// ApprovalPolicySpec carries versioned approval stages and eligibility.
type ApprovalPolicySpec struct {
	Version             string           `json:"version"`
	PublicationState    PublicationState `json:"publicationState"`
	ApproverEligibility []EligibilityRef `json:"approverEligibilityRefs"`
}

// ApprovalRequest is the FEATURE-0018 ApprovalRequest resource (VS0-SCHEMA-066; REQ-F18-13).
type ApprovalRequest struct {
	apimeta.TypeMeta
	Metadata apimeta.ObjectMeta    `json:"metadata"`
	Spec     ApprovalRequestSpec   `json:"spec"`
	Status   ApprovalRequestStatus `json:"status,omitempty"`
}

// ApprovalRequestSpec pins the immutable bounded subject and policy.
type ApprovalRequestSpec struct {
	ApprovalPolicyRef apimeta.TypedRef `json:"approvalPolicyRef"`
	SubjectKind       string           `json:"subjectKind"`
	SubjectUID        string           `json:"subjectUid"`
}

// ApprovalRequestStatus is system-owned decision/LRO state.
type ApprovalRequestStatus struct {
	Phase    string           `json:"phase,omitempty"`
	Decision ApprovalDecision `json:"decision,omitempty"`
}

// ExceptionGrant is the FEATURE-0018 ExceptionGrant resource (VS0-SCHEMA-067; REQ-F18-17).
type ExceptionGrant struct {
	apimeta.TypeMeta
	Metadata apimeta.ObjectMeta   `json:"metadata"`
	Record   ExceptionGrantRecord `json:"record"`
}

// ExceptionGrantRecord is the immutable exception evidence payload.
type ExceptionGrantRecord struct {
	Effect     ExceptionGrantEffect `json:"effect"`
	ControlRef apimeta.TypedRef     `json:"controlRef"`
	SubjectUID string               `json:"subjectUid"`
	NotBefore  time.Time            `json:"notBefore"`
	ExpiresAt  time.Time            `json:"expiresAt"`
}

// GovernanceProfile is the FEATURE-0018 GovernanceProfile resource
// (VS0-SCHEMA-024; REQ-F18-18).
type GovernanceProfile struct {
	apimeta.TypeMeta
	Metadata apimeta.ObjectMeta    `json:"metadata"`
	Spec     GovernanceProfileSpec `json:"spec"`
}

// GovernanceProfileSpec composes FEATURE-0018-owned governance references.
type GovernanceProfileSpec struct {
	Version            string             `json:"version"`
	PublicationState   PublicationState   `json:"publicationState"`
	ApprovalPolicyRefs []apimeta.TypedRef `json:"approvalPolicyRefs,omitempty"`
}
