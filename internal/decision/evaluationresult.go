package decision

import (
	"encoding/json"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
)

// Canonical EvaluationResult identity constants for TransientRequestResult
// exchange (F13-OBJ-004). EvaluationResult itself carries no TypeMeta or
// ObjectMeta: as an EmbeddedValue it must not receive artificial identity
// (FEATURE-0012 EmbeddedValue invariant).
const (
	APIVersionEvaluationResult = "governance.sovrunn.io/v1alpha1"
	KindEvaluationResult       = "EvaluationResult"
)

// EvaluationResultStatus is the closed evaluation result-status vocabulary
// (F13-EVAL-001). Values are reproduced from architecture §7.1; no value may
// be added.
type EvaluationResultStatus string

const (
	EvaluationResultStatusSuccess       EvaluationResultStatus = "SUCCESS"
	EvaluationResultStatusNonMatch      EvaluationResultStatus = "NON_MATCH"
	EvaluationResultStatusIndeterminate EvaluationResultStatus = "INDETERMINATE"
	EvaluationResultStatusTimeout       EvaluationResultStatus = "TIMEOUT"
	EvaluationResultStatusError         EvaluationResultStatus = "ERROR"
)

// AllEvaluationResultStatuses returns the closed EvaluationResultStatus set
// in stable order.
func AllEvaluationResultStatuses() []EvaluationResultStatus {
	return []EvaluationResultStatus{
		EvaluationResultStatusSuccess,
		EvaluationResultStatusNonMatch,
		EvaluationResultStatusIndeterminate,
		EvaluationResultStatusTimeout,
		EvaluationResultStatusError,
	}
}

// Valid reports whether s is one of the five closed EvaluationResultStatus values.
func (s EvaluationResultStatus) Valid() bool {
	switch s {
	case EvaluationResultStatusSuccess,
		EvaluationResultStatusNonMatch,
		EvaluationResultStatusIndeterminate,
		EvaluationResultStatusTimeout,
		EvaluationResultStatusError:
		return true
	default:
		return false
	}
}

// EvaluatorIdentity is the registered evaluator type and version
// (F13-EVAL-001/002; design §7.2).
type EvaluatorIdentity struct {
	Type    string `json:"type"`    // registered evaluator type id; required non-empty
	Version string `json:"version"` // evaluator version; required non-empty
}

// TimingEnvelope captures evaluation wall-clock boundaries (F13-EVAL-001).
type TimingEnvelope struct {
	StartedAt   string `json:"startedAt"`            // UTC RFC3339; required
	CompletedAt string `json:"completedAt"`          // UTC RFC3339; required
	DurationMs  int64  `json:"durationMs,omitempty"` // optional observed duration
}

// OpaqueIntegrityCarrier is a structural integrity carrier for a governed
// evaluation input snapshot (F13-EVAL-001/002). Fields are opaque structural
// identifiers only; no digest, signing, or verification is performed.
// A present-but-empty optional carrier is invalid (RID-08).
type OpaqueIntegrityCarrier struct {
	State                  string `json:"state"`                            // structural trust/integrity state; required when carrier present
	CanonicalizationMethod string `json:"canonicalizationMethod,omitempty"` // structural id; no execution
	DigestAlgorithm        string `json:"digestAlgorithm,omitempty"`        // structural id; no computation
	CoveredFields          string `json:"coveredFields,omitempty"`          // covered-field descriptor
	SignatureAlgorithm     string `json:"signatureAlgorithm,omitempty"`     // structural id; no signing/verification
}

// EvaluationResult is one registered evaluator's normalized finding
// (EmbeddedValue when captured in a DecisionRecord; TransientRequestResult
// when exchanged independently). It has no independent durable lifecycle and
// no independent scope field (F13-OBJ-004, F13-EVAL-003). Scope is derived
// from the containing DecisionRecord's metadata.scopeRef.
type EvaluationResult struct {
	Evaluator EvaluatorIdentity `json:"evaluator"`

	// Input identity: bounded input-snapshot reference and/or opaque integrity
	// carrier (F13-EVAL-001). Requiredness of the pair is enforced by
	// validation (T-020); both may be present (snapshot + optional carrier).
	InputSnapshotRef *apimeta.TypedRef       `json:"inputSnapshotRef,omitempty"`
	InputIntegrity   *OpaqueIntegrityCarrier `json:"inputIntegrity,omitempty"`

	ResultStatus EvaluationResultStatus `json:"resultStatus"`
	Result       json.RawMessage        `json:"result,omitempty"` // normalized evaluator output

	Timing      TimingEnvelope `json:"timing"`
	EvaluatedAt string         `json:"evaluatedAt"` // evaluation time; UTC RFC3339

	// Relevant execution configuration (opaque profile-declared object).
	ExecutionConfig json.RawMessage `json:"executionConfig,omitempty"`

	TrustBoundary       string   `json:"trustBoundary,omitempty"`
	SafetyPolicyFilters []string `json:"safetyPolicyFilters,omitempty"`

	// Structural trust state captured with the result (F13-EVAL-002).
	// Opaque; no cryptographic-validity claim.
	StructuralTrustState string `json:"structuralTrustState,omitempty"`

	// ScopeEvidence is evaluator-supplied scope metadata for structural
	// comparison only (F13-EVAL-003/004). It is evidence, never a second
	// scope authority. Canonical scope remains the containing record's
	// metadata.scopeRef.
	ScopeEvidence *apimeta.ScopeRef `json:"scopeEvidence,omitempty"`
}
