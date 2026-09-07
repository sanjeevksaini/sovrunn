package model

// Closed FEATURE-0018 supporting enums (requirements §3.3; F18-RD-02; DD-04).

// MembershipContextKind distinguishes Organization, CloudProvider, and AccessGroup contexts.
type MembershipContextKind string

const (
	MembershipContextOrganization  MembershipContextKind = "Organization"
	MembershipContextCloudProvider MembershipContextKind = "CloudProvider"
	MembershipContextAccessGroup   MembershipContextKind = "AccessGroup"
)

// AllMembershipContextKinds returns the closed set in stable order.
func AllMembershipContextKinds() []MembershipContextKind {
	return []MembershipContextKind{
		MembershipContextOrganization,
		MembershipContextCloudProvider,
		MembershipContextAccessGroup,
	}
}

// Valid reports whether k is a closed MembershipContextKind.
func (k MembershipContextKind) Valid() bool {
	switch k {
	case MembershipContextOrganization, MembershipContextCloudProvider, MembershipContextAccessGroup:
		return true
	default:
		return false
	}
}

// MembershipType is Standard | Guest.
type MembershipType string

const (
	MembershipTypeStandard MembershipType = "Standard"
	MembershipTypeGuest    MembershipType = "Guest"
)

// AllMembershipTypes returns the closed set in stable order.
func AllMembershipTypes() []MembershipType {
	return []MembershipType{MembershipTypeStandard, MembershipTypeGuest}
}

// Valid reports whether t is a closed MembershipType.
func (t MembershipType) Valid() bool {
	switch t {
	case MembershipTypeStandard, MembershipTypeGuest:
		return true
	default:
		return false
	}
}

// PrivilegedAccessMode is JIT | BreakGlass.
type PrivilegedAccessMode string

const (
	PrivilegedAccessModeJIT        PrivilegedAccessMode = "JIT"
	PrivilegedAccessModeBreakGlass PrivilegedAccessMode = "BreakGlass"
)

// AllPrivilegedAccessModes returns the closed set in stable order.
func AllPrivilegedAccessModes() []PrivilegedAccessMode {
	return []PrivilegedAccessMode{PrivilegedAccessModeJIT, PrivilegedAccessModeBreakGlass}
}

// Valid reports whether m is a closed PrivilegedAccessMode.
func (m PrivilegedAccessMode) Valid() bool {
	switch m {
	case PrivilegedAccessModeJIT, PrivilegedAccessModeBreakGlass:
		return true
	default:
		return false
	}
}

// ApprovalDecision is Approve | Deny.
type ApprovalDecision string

const (
	ApprovalDecisionApprove ApprovalDecision = "Approve"
	ApprovalDecisionDeny    ApprovalDecision = "Deny"
)

// AllApprovalDecisions returns the closed set in stable order.
func AllApprovalDecisions() []ApprovalDecision {
	return []ApprovalDecision{ApprovalDecisionApprove, ApprovalDecisionDeny}
}

// Valid reports whether d is a closed ApprovalDecision.
func (d ApprovalDecision) Valid() bool {
	switch d {
	case ApprovalDecisionApprove, ApprovalDecisionDeny:
		return true
	default:
		return false
	}
}

// ExceptionGrantEffect is Grant | Revoke.
type ExceptionGrantEffect string

const (
	ExceptionGrantEffectGrant  ExceptionGrantEffect = "Grant"
	ExceptionGrantEffectRevoke ExceptionGrantEffect = "Revoke"
)

// AllExceptionGrantEffects returns the closed set in stable order.
func AllExceptionGrantEffects() []ExceptionGrantEffect {
	return []ExceptionGrantEffect{ExceptionGrantEffectGrant, ExceptionGrantEffectRevoke}
}

// Valid reports whether e is a closed ExceptionGrantEffect.
func (e ExceptionGrantEffect) Valid() bool {
	switch e {
	case ExceptionGrantEffectGrant, ExceptionGrantEffectRevoke:
		return true
	default:
		return false
	}
}

// AssignmentValidityMode is Standing | TimeBound.
type AssignmentValidityMode string

const (
	AssignmentValidityStanding  AssignmentValidityMode = "Standing"
	AssignmentValidityTimeBound AssignmentValidityMode = "TimeBound"
)

// AllAssignmentValidityModes returns the closed set in stable order.
func AllAssignmentValidityModes() []AssignmentValidityMode {
	return []AssignmentValidityMode{AssignmentValidityStanding, AssignmentValidityTimeBound}
}

// Valid reports whether m is a closed AssignmentValidityMode.
func (m AssignmentValidityMode) Valid() bool {
	switch m {
	case AssignmentValidityStanding, AssignmentValidityTimeBound:
		return true
	default:
		return false
	}
}

// ActionTargetMode is ExactResource | CreateParent | ScopeOnly (REQ-F18-06).
type ActionTargetMode string

const (
	ActionTargetExactResource ActionTargetMode = "ExactResource"
	ActionTargetCreateParent  ActionTargetMode = "CreateParent"
	ActionTargetScopeOnly     ActionTargetMode = "ScopeOnly"
)

// AllActionTargetModes returns the closed set in stable order.
func AllActionTargetModes() []ActionTargetMode {
	return []ActionTargetMode{
		ActionTargetExactResource,
		ActionTargetCreateParent,
		ActionTargetScopeOnly,
	}
}

// Valid reports whether m is a closed ActionTargetMode.
func (m ActionTargetMode) Valid() bool {
	switch m {
	case ActionTargetExactResource, ActionTargetCreateParent, ActionTargetScopeOnly:
		return true
	default:
		return false
	}
}

// AssuranceLevel is the closed AAL vocabulary for AssuranceEvidence (REQ-F18-03).
type AssuranceLevel string

const (
	AssuranceLevelAAL1 AssuranceLevel = "AAL1"
	AssuranceLevelAAL2 AssuranceLevel = "AAL2"
	AssuranceLevelAAL3 AssuranceLevel = "AAL3"
)

// AllAssuranceLevels returns the closed set in stable order.
func AllAssuranceLevels() []AssuranceLevel {
	return []AssuranceLevel{AssuranceLevelAAL1, AssuranceLevelAAL2, AssuranceLevelAAL3}
}

// Valid reports whether l is a closed AssuranceLevel.
func (l AssuranceLevel) Valid() bool {
	switch l {
	case AssuranceLevelAAL1, AssuranceLevelAAL2, AssuranceLevelAAL3:
		return true
	default:
		return false
	}
}

// RoleClassification is Ordinary | Privileged (REQ-F18-07).
type RoleClassification string

const (
	RoleClassificationOrdinary   RoleClassification = "Ordinary"
	RoleClassificationPrivileged RoleClassification = "Privileged"
)

// AllRoleClassifications returns the closed set in stable order.
func AllRoleClassifications() []RoleClassification {
	return []RoleClassification{RoleClassificationOrdinary, RoleClassificationPrivileged}
}

// Valid reports whether c is a closed RoleClassification.
func (c RoleClassification) Valid() bool {
	switch c {
	case RoleClassificationOrdinary, RoleClassificationPrivileged:
		return true
	default:
		return false
	}
}

// PublicationState covers VersionedDefinition publication states used by
// RoleDefinition, ApprovalPolicy, and GovernanceProfile.
type PublicationState string

const (
	PublicationStateDraft     PublicationState = "Draft"
	PublicationStatePublished PublicationState = "Published"
	PublicationStateSuspended PublicationState = "Suspended"
	PublicationStateRetired   PublicationState = "Retired"
)

// AllPublicationStates returns the closed set in stable order.
func AllPublicationStates() []PublicationState {
	return []PublicationState{
		PublicationStateDraft,
		PublicationStatePublished,
		PublicationStateSuspended,
		PublicationStateRetired,
	}
}

// Valid reports whether s is a closed PublicationState.
func (s PublicationState) Valid() bool {
	switch s {
	case PublicationStateDraft, PublicationStatePublished, PublicationStateSuspended, PublicationStateRetired:
		return true
	default:
		return false
	}
}

// AccessGroupPhase is Active | Suspended | Retired (REQ-F18-04).
type AccessGroupPhase string

const (
	AccessGroupPhaseActive    AccessGroupPhase = "Active"
	AccessGroupPhaseSuspended AccessGroupPhase = "Suspended"
	AccessGroupPhaseRetired   AccessGroupPhase = "Retired"
)

// AllAccessGroupPhases returns the closed set in stable order.
func AllAccessGroupPhases() []AccessGroupPhase {
	return []AccessGroupPhase{AccessGroupPhaseActive, AccessGroupPhaseSuspended, AccessGroupPhaseRetired}
}

// Valid reports whether p is a closed AccessGroupPhase.
func (p AccessGroupPhase) Valid() bool {
	switch p {
	case AccessGroupPhaseActive, AccessGroupPhaseSuspended, AccessGroupPhaseRetired:
		return true
	default:
		return false
	}
}

// MembershipPhase is Active | Suspended | Revoked (VS0-SCHEMA-021).
type MembershipPhase string

const (
	MembershipPhaseActive    MembershipPhase = "Active"
	MembershipPhaseSuspended MembershipPhase = "Suspended"
	MembershipPhaseRevoked   MembershipPhase = "Revoked"
)

// AllMembershipPhases returns the closed set in stable order.
func AllMembershipPhases() []MembershipPhase {
	return []MembershipPhase{MembershipPhaseActive, MembershipPhaseSuspended, MembershipPhaseRevoked}
}

// Valid reports whether p is a closed MembershipPhase.
func (p MembershipPhase) Valid() bool {
	switch p {
	case MembershipPhaseActive, MembershipPhaseSuspended, MembershipPhaseRevoked:
		return true
	default:
		return false
	}
}
