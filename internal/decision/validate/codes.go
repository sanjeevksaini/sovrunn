// Package validate implements FEATURE-0013 decision and AuditEvent validation
// against the closed public Problem Details violation-code registry.
package validate

import "github.com/sanjeevksaini/sovrunn/internal/apiproblem"

// Closed DECISION_* violations[].code constants (architecture §15.2; F13-ERR-002).
// These names are public compatibility contracts. No code may be added, renamed,
// aliased, or remapped without an approved architecture change (F13-ERR-004).
// Prefixes are for violations[].code only; they are never Problem.type URI
// families, top-level Problem.code values, or HTTP-status selectors (F13-ERR-003).
const (
	// Decision profile family
	CodeProfileUnknown            apiproblem.ViolationCode = "DECISION_PROFILE_UNKNOWN"
	CodeProfileVersionUnsupported apiproblem.ViolationCode = "DECISION_PROFILE_VERSION_UNSUPPORTED"
	CodeProfileInactive           apiproblem.ViolationCode = "DECISION_PROFILE_INACTIVE"
	CodeProfileSchemaInvalid      apiproblem.ViolationCode = "DECISION_PROFILE_SCHEMA_INVALID"
	CodeProfileLimitInvalid       apiproblem.ViolationCode = "DECISION_PROFILE_LIMIT_INVALID"

	// Evaluation family
	CodeEvaluationTypeUnsupported apiproblem.ViolationCode = "DECISION_EVALUATION_TYPE_UNSUPPORTED"
	CodeEvaluationResultInvalid   apiproblem.ViolationCode = "DECISION_EVALUATION_RESULT_INVALID"
	CodeEvaluationScopeMismatch   apiproblem.ViolationCode = "DECISION_EVALUATION_SCOPE_MISMATCH"
	CodeEvaluationLimitExceeded   apiproblem.ViolationCode = "DECISION_EVALUATION_LIMIT_EXCEEDED"

	// Composition family
	CodeCompositionStrategyUnsupported apiproblem.ViolationCode = "DECISION_COMPOSITION_STRATEGY_UNSUPPORTED"
	CodeCompositionGraphInvalid        apiproblem.ViolationCode = "DECISION_COMPOSITION_GRAPH_INVALID"
	CodeCompositionCycle               apiproblem.ViolationCode = "DECISION_COMPOSITION_CYCLE"
	CodeCompositionLimitExceeded       apiproblem.ViolationCode = "DECISION_COMPOSITION_LIMIT_EXCEEDED"
	CodeCompositionInputConflict       apiproblem.ViolationCode = "DECISION_COMPOSITION_INPUT_CONFLICT"

	// Obligation family
	CodeObligationUnknown              apiproblem.ViolationCode = "DECISION_OBLIGATION_UNKNOWN"
	CodeObligationInvalid              apiproblem.ViolationCode = "DECISION_OBLIGATION_INVALID"
	CodeObligationUnsupportedMandatory apiproblem.ViolationCode = "DECISION_OBLIGATION_UNSUPPORTED_MANDATORY"

	// Trust family
	CodeTrustRequired apiproblem.ViolationCode = "DECISION_TRUST_REQUIRED"
	CodeTrustUnknown  apiproblem.ViolationCode = "DECISION_TRUST_UNKNOWN"
	CodeTrustExpired  apiproblem.ViolationCode = "DECISION_TRUST_EXPIRED"
	CodeTrustRevoked  apiproblem.ViolationCode = "DECISION_TRUST_REVOKED"
	CodeTrustMismatch apiproblem.ViolationCode = "DECISION_TRUST_MISMATCH"

	// Scope family
	CodeScopeRequired apiproblem.ViolationCode = "DECISION_SCOPE_REQUIRED"
	CodeScopeInvalid  apiproblem.ViolationCode = "DECISION_SCOPE_INVALID"
	CodeScopeConflict apiproblem.ViolationCode = "DECISION_SCOPE_CONFLICT"
	CodeScopeMismatch apiproblem.ViolationCode = "DECISION_SCOPE_MISMATCH"
	CodeScopeWidening apiproblem.ViolationCode = "DECISION_SCOPE_WIDENING"

	// Relationship family
	CodeRelationshipKindInvalid        apiproblem.ViolationCode = "DECISION_RELATIONSHIP_KIND_INVALID"
	CodeRelationshipTargetInvalid      apiproblem.ViolationCode = "DECISION_RELATIONSHIP_TARGET_INVALID"
	CodeRelationshipCycle              apiproblem.ViolationCode = "DECISION_RELATIONSHIP_CYCLE"
	CodeRelationshipChainLimitExceeded apiproblem.ViolationCode = "DECISION_RELATIONSHIP_CHAIN_LIMIT_EXCEEDED"
	CodeRelationshipConflict           apiproblem.ViolationCode = "DECISION_RELATIONSHIP_CONFLICT"
)

// AllCodes returns the architecture-owned closed DECISION_* violation-code
// registry in stable family order (architecture §15.2). Exactly 32 codes.
func AllCodes() []apiproblem.ViolationCode {
	return []apiproblem.ViolationCode{
		CodeProfileUnknown,
		CodeProfileVersionUnsupported,
		CodeProfileInactive,
		CodeProfileSchemaInvalid,
		CodeProfileLimitInvalid,
		CodeEvaluationTypeUnsupported,
		CodeEvaluationResultInvalid,
		CodeEvaluationScopeMismatch,
		CodeEvaluationLimitExceeded,
		CodeCompositionStrategyUnsupported,
		CodeCompositionGraphInvalid,
		CodeCompositionCycle,
		CodeCompositionLimitExceeded,
		CodeCompositionInputConflict,
		CodeObligationUnknown,
		CodeObligationInvalid,
		CodeObligationUnsupportedMandatory,
		CodeTrustRequired,
		CodeTrustUnknown,
		CodeTrustExpired,
		CodeTrustRevoked,
		CodeTrustMismatch,
		CodeScopeRequired,
		CodeScopeInvalid,
		CodeScopeConflict,
		CodeScopeMismatch,
		CodeScopeWidening,
		CodeRelationshipKindInvalid,
		CodeRelationshipTargetInvalid,
		CodeRelationshipCycle,
		CodeRelationshipChainLimitExceeded,
		CodeRelationshipConflict,
	}
}

// Valid reports whether c is one of the closed FEATURE-0013 DECISION_*
// violations[].code values. Unknown codes fail closed.
func Valid(c apiproblem.ViolationCode) bool {
	switch c {
	case CodeProfileUnknown,
		CodeProfileVersionUnsupported,
		CodeProfileInactive,
		CodeProfileSchemaInvalid,
		CodeProfileLimitInvalid,
		CodeEvaluationTypeUnsupported,
		CodeEvaluationResultInvalid,
		CodeEvaluationScopeMismatch,
		CodeEvaluationLimitExceeded,
		CodeCompositionStrategyUnsupported,
		CodeCompositionGraphInvalid,
		CodeCompositionCycle,
		CodeCompositionLimitExceeded,
		CodeCompositionInputConflict,
		CodeObligationUnknown,
		CodeObligationInvalid,
		CodeObligationUnsupportedMandatory,
		CodeTrustRequired,
		CodeTrustUnknown,
		CodeTrustExpired,
		CodeTrustRevoked,
		CodeTrustMismatch,
		CodeScopeRequired,
		CodeScopeInvalid,
		CodeScopeConflict,
		CodeScopeMismatch,
		CodeScopeWidening,
		CodeRelationshipKindInvalid,
		CodeRelationshipTargetInvalid,
		CodeRelationshipCycle,
		CodeRelationshipChainLimitExceeded,
		CodeRelationshipConflict:
		return true
	default:
		return false
	}
}
