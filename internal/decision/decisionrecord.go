package decision

import (
	"encoding/json"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
)

// Canonical DecisionRecord TypeMeta values (AD-001; F13-OBJ-001).
const (
	APIVersionDecisionRecord = "governance.sovrunn.io/v1alpha1"
	KindDecisionRecord       = "DecisionRecord"
)

// DecisionRecord is the immutable append-only governed decision envelope
// (ImmutableRecord; design §7.1). It has no mutable spec or status.
// metadata.scopeRef is the sole scope authority.
type DecisionRecord struct {
	apimeta.TypeMeta                    // apiVersion=governance.sovrunn.io/v1alpha1, kind=DecisionRecord
	Metadata         apimeta.ObjectMeta `json:"metadata"`
	Record           DecisionBody       `json:"record"`
}

// DecisionBody is the append-only DecisionRecord payload (design §7.1).
type DecisionBody struct {
	ProfileRef        ProfileRef               `json:"profileRef"`
	Form              DecisionForm             `json:"form"`
	Authority         DecisionAuthority        `json:"authority"`
	Purpose           string                   `json:"purpose"`
	Adjudication      *AdjudicationOutcome     `json:"adjudication,omitempty"` // adjudication facets only
	SubjectRefs       []apimeta.TypedRef       `json:"subjectRefs,omitempty"`
	EvaluationResults []EvaluationResult       `json:"evaluationResults,omitempty"` // EmbeddedValue (design §7.1/§7.2)
	Composition       *CompositionRef          `json:"composition,omitempty"`
	Result            DecisionResultBody       `json:"result"`
	Relationship      *DecisionRelationship    `json:"relationship,omitempty"`
	RetryKey          string                   `json:"retryKey,omitempty"` // RetryIdempotencyKey carrier (F13-SCALE-001; AD-036)
	SemanticIdentity  SemanticDecisionIdentity `json:"semanticIdentity"`   // F13-SCALE-002; AD-036
	Trust             *TrustCarrier            `json:"trust,omitempty"`    // algorithm-agile structural trust (F13-TRUST-001; RID-08)
	Sovereignty       *SovereigntyCarrier      `json:"sovereignty,omitempty"`
	Correlation       Correlation              `json:"correlation"`
	Finality          Finality                 `json:"finality"` // always FINAL when persisted
}

// ProfileRef is the versioned DecisionProfile identity (design §7.9).
// Lookup uses the composite (name, version) key.
type ProfileRef struct {
	Name    string `json:"name"`    // profile id; required non-empty
	Version string `json:"version"` // semver; required
}

// DecisionResultBody carries the profile-declared typed result, structured
// rationale, and obligations (architecture §6.1; design §7.1).
type DecisionResultBody struct {
	TypedResult json.RawMessage   `json:"typedResult"` // profile-declared opaque object
	Rationale   DecisionRationale `json:"rationale"`
	Obligations []Obligation      `json:"obligations"` // required-present; empty allowed
}

// DecisionRationale is structured explainability embedded in the decision
// (F13-OBJ-001; architecture DecisionRationale).
type DecisionRationale struct {
	ReasonCodes           []string `json:"reasonCodes,omitempty"`
	Reasons               []string `json:"reasons,omitempty"`
	Alternatives          []string `json:"alternatives,omitempty"`
	Warnings              []string `json:"warnings,omitempty"`
	CorrectiveSuggestions []string `json:"correctiveSuggestions,omitempty"`
}

// Obligation is an enforceable condition tied to the decision
// (F13-OBLIG-001; design §7.1).
type Obligation struct {
	ID        string          `json:"id"`                // vocabulary identity; required non-empty
	Mandatory bool            `json:"mandatory"`         // required explicit bool
	Payload   json.RawMessage `json:"payload,omitempty"` // profile-declared opaque payload
}

// Correlation links the decision to preallocated AuditEvent identities and
// optional operation/trace correlation (design §7.1/§7.9; F13-AUDIT-001).
// AuditEventRefs is required non-empty (minItems: 1).
type Correlation struct {
	AuditEventRefs []apimeta.TypedRef `json:"auditEventRefs"`         // required non-empty
	OperationRef   *apimeta.TypedRef  `json:"operationRef,omitempty"` // optional singular ref
	TraceRef       string             `json:"traceRef,omitempty"`     // opaque operational correlation
}

// CompositionRef records the chosen composition strategy and optional graph
// identity for deterministic replay (architecture §7.2; design §7.1).
type CompositionRef struct {
	StrategyName    string   `json:"strategyName"`
	StrategyVersion string   `json:"strategyVersion"`
	GraphRef        string   `json:"graphRef,omitempty"`
	GraphVersion    string   `json:"graphVersion,omitempty"`
	InputRefs       []string `json:"inputRefs,omitempty"`
}

// SovereigntyCarrier is a structural carrier for the seven sovereignty
// dimensions (AD-016; F13-SOV-001). It carries tags/references only and must
// not embed secrets, credentials, or provider-native resource types.
type SovereigntyCarrier struct {
	Data                string `json:"data,omitempty"`
	Operational         string `json:"operational,omitempty"`
	Technical           string `json:"technical,omitempty"`
	Legal               string `json:"legal,omitempty"`
	Cryptographic       string `json:"cryptographic,omitempty"`
	SoftwareSupplyChain string `json:"softwareSupplyChain,omitempty"`
	AIModel             string `json:"aiModel,omitempty"`
}

// DecisionRelationship links a new immutable record to a predecessor
// (design §7.9; F13-IMMUT-002). Orientation is new record -> predecessor.
type DecisionRelationship struct {
	Kind      RelationshipKind `json:"kind"`      // CORRECTS | SUPERSEDES | REVOKES
	TargetRef apimeta.TypedRef `json:"targetRef"` // predecessor; uid required at validation
}
