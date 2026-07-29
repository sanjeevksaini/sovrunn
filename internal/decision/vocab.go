package decision

// DecisionForm is the closed decision-form vocabulary (F13-FORM-001; AD-030).
// Values are reproduced exactly from design §7.5; no value may be added.
type DecisionForm string

const (
	DecisionFormAdjudication   DecisionForm = "ADJUDICATION"
	DecisionFormSelection      DecisionForm = "SELECTION"
	DecisionFormRanking        DecisionForm = "RANKING"
	DecisionFormClassification DecisionForm = "CLASSIFICATION"
	DecisionFormResolution     DecisionForm = "RESOLUTION"
	DecisionFormAllocation     DecisionForm = "ALLOCATION"
	DecisionFormPlan           DecisionForm = "PLAN"
	DecisionFormAssessment     DecisionForm = "ASSESSMENT"
	DecisionFormRecommendation DecisionForm = "RECOMMENDATION"
)

// AllDecisionForms returns the closed DecisionForm set in stable order.
func AllDecisionForms() []DecisionForm {
	return []DecisionForm{
		DecisionFormAdjudication,
		DecisionFormSelection,
		DecisionFormRanking,
		DecisionFormClassification,
		DecisionFormResolution,
		DecisionFormAllocation,
		DecisionFormPlan,
		DecisionFormAssessment,
		DecisionFormRecommendation,
	}
}

// Valid reports whether f is one of the nine closed DecisionForm values.
func (f DecisionForm) Valid() bool {
	switch f {
	case DecisionFormAdjudication,
		DecisionFormSelection,
		DecisionFormRanking,
		DecisionFormClassification,
		DecisionFormResolution,
		DecisionFormAllocation,
		DecisionFormPlan,
		DecisionFormAssessment,
		DecisionFormRecommendation:
		return true
	default:
		return false
	}
}

// DecisionAuthority is the closed authority vocabulary (F13-AUTH-001; AD-031).
// Values are reproduced exactly from design §7.5; no value may be added.
type DecisionAuthority string

const (
	DecisionAuthorityAuthoritative  DecisionAuthority = "AUTHORITATIVE"
	DecisionAuthorityAdvisory       DecisionAuthority = "ADVISORY"
	DecisionAuthorityRecommendation DecisionAuthority = "RECOMMENDATION"
	DecisionAuthoritySimulation     DecisionAuthority = "SIMULATION"
)

// AllDecisionAuthorities returns the closed DecisionAuthority set in stable order.
func AllDecisionAuthorities() []DecisionAuthority {
	return []DecisionAuthority{
		DecisionAuthorityAuthoritative,
		DecisionAuthorityAdvisory,
		DecisionAuthorityRecommendation,
		DecisionAuthoritySimulation,
	}
}

// Valid reports whether a is one of the four closed DecisionAuthority values.
func (a DecisionAuthority) Valid() bool {
	switch a {
	case DecisionAuthorityAuthoritative,
		DecisionAuthorityAdvisory,
		DecisionAuthorityRecommendation,
		DecisionAuthoritySimulation:
		return true
	default:
		return false
	}
}

// AdjudicationOutcome is the closed adjudication-outcome vocabulary
// (F13-ADJ-001; AD-003). Values are reproduced exactly from design §7.5.
type AdjudicationOutcome string

const (
	AdjudicationOutcomeAllowed          AdjudicationOutcome = "ALLOWED"
	AdjudicationOutcomeDenied           AdjudicationOutcome = "DENIED"
	AdjudicationOutcomeRequiresApproval AdjudicationOutcome = "REQUIRES_APPROVAL"
)

// AllAdjudicationOutcomes returns the closed AdjudicationOutcome set in stable order.
func AllAdjudicationOutcomes() []AdjudicationOutcome {
	return []AdjudicationOutcome{
		AdjudicationOutcomeAllowed,
		AdjudicationOutcomeDenied,
		AdjudicationOutcomeRequiresApproval,
	}
}

// Valid reports whether o is one of the three closed AdjudicationOutcome values.
func (o AdjudicationOutcome) Valid() bool {
	switch o {
	case AdjudicationOutcomeAllowed,
		AdjudicationOutcomeDenied,
		AdjudicationOutcomeRequiresApproval:
		return true
	default:
		return false
	}
}

// AuditOutcome is the closed audit-outcome vocabulary (F13-AUDIT-002).
// Values match the FEATURE-0012 AuditEvent outcome set exactly
// (Succeeded, Denied, Failed); no value may be added or renamed.
type AuditOutcome string

const (
	AuditOutcomeSucceeded AuditOutcome = "Succeeded"
	AuditOutcomeDenied    AuditOutcome = "Denied"
	AuditOutcomeFailed    AuditOutcome = "Failed"
)

// AllAuditOutcomes returns the closed AuditOutcome set in stable order.
func AllAuditOutcomes() []AuditOutcome {
	return []AuditOutcome{
		AuditOutcomeSucceeded,
		AuditOutcomeDenied,
		AuditOutcomeFailed,
	}
}

// Valid reports whether o is one of the three closed AuditOutcome values.
func (o AuditOutcome) Valid() bool {
	switch o {
	case AuditOutcomeSucceeded, AuditOutcomeDenied, AuditOutcomeFailed:
		return true
	default:
		return false
	}
}

// RelationshipKind is the closed immutable-relationship vocabulary
// (F13-IMMUT-002; AD-004). Values are reproduced exactly from design §7.5.
type RelationshipKind string

const (
	RelationshipKindCorrects   RelationshipKind = "CORRECTS"
	RelationshipKindSupersedes RelationshipKind = "SUPERSEDES"
	RelationshipKindRevokes    RelationshipKind = "REVOKES"
)

// AllRelationshipKinds returns the closed RelationshipKind set in stable order.
func AllRelationshipKinds() []RelationshipKind {
	return []RelationshipKind{
		RelationshipKindCorrects,
		RelationshipKindSupersedes,
		RelationshipKindRevokes,
	}
}

// Valid reports whether k is one of the three closed RelationshipKind values.
func (k RelationshipKind) Valid() bool {
	switch k {
	case RelationshipKindCorrects, RelationshipKindSupersedes, RelationshipKindRevokes:
		return true
	default:
		return false
	}
}

// Sensitivity is the closed ordered sensitivity vocabulary (F13-SEC-001; AD-019).
// Order is PUBLIC < INTERNAL < CONFIDENTIAL < RESTRICTED. Values are reproduced
// exactly from design §7.5; the classification floor map is added separately.
type Sensitivity string

const (
	SensitivityPublic       Sensitivity = "PUBLIC"
	SensitivityInternal     Sensitivity = "INTERNAL"
	SensitivityConfidential Sensitivity = "CONFIDENTIAL"
	SensitivityRestricted   Sensitivity = "RESTRICTED"
)

// AllSensitivities returns the closed Sensitivity set in stable ascending order.
func AllSensitivities() []Sensitivity {
	return []Sensitivity{
		SensitivityPublic,
		SensitivityInternal,
		SensitivityConfidential,
		SensitivityRestricted,
	}
}

// Valid reports whether s is one of the four closed Sensitivity values.
func (s Sensitivity) Valid() bool {
	switch s {
	case SensitivityPublic,
		SensitivityInternal,
		SensitivityConfidential,
		SensitivityRestricted:
		return true
	default:
		return false
	}
}

// Finality is the closed finality vocabulary (F13-IMMUT-001). A persisted
// DecisionRecord always carries FINAL; no pending/non-final value exists.
type Finality string

const (
	FinalityFinal Finality = "FINAL"
)

// AllFinalities returns the closed Finality set in stable order.
func AllFinalities() []Finality {
	return []Finality{
		FinalityFinal,
	}
}

// Valid reports whether f is the sole closed Finality value.
func (f Finality) Valid() bool {
	return f == FinalityFinal
}

// EffectiveState is a projection-only effective-state vocabulary
// (F13-IMMUT-003; AD-034). It is never written back to the immutable record.
type EffectiveState string

const (
	EffectiveStateEffective  EffectiveState = "EFFECTIVE"
	EffectiveStateSuperseded EffectiveState = "SUPERSEDED"
	EffectiveStateRevoked    EffectiveState = "REVOKED"
)

// AllEffectiveStates returns the closed EffectiveState set in stable order.
func AllEffectiveStates() []EffectiveState {
	return []EffectiveState{
		EffectiveStateEffective,
		EffectiveStateSuperseded,
		EffectiveStateRevoked,
	}
}

// Valid reports whether s is one of the three closed EffectiveState values.
func (s EffectiveState) Valid() bool {
	switch s {
	case EffectiveStateEffective, EffectiveStateSuperseded, EffectiveStateRevoked:
		return true
	default:
		return false
	}
}
