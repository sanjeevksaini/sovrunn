package model

import (
	"bytes"
	"time"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
)

// ApprovalTerminalResult is the F18-RD-19 approval-decision aggregate result.
type ApprovalTerminalResult string

const (
	ApprovalTerminalApproved ApprovalTerminalResult = "Approved"
	ApprovalTerminalDenied   ApprovalTerminalResult = "Denied"
)

// Valid reports whether r is a closed ApprovalTerminalResult.
func (r ApprovalTerminalResult) Valid() bool {
	switch r {
	case ApprovalTerminalApproved, ApprovalTerminalDenied:
		return true
	default:
		return false
	}
}

// ExceptionTerminalResult is the F18-RD-19 exception-decision terminal result.
type ExceptionTerminalResult string

const (
	ExceptionTerminalGrant ExceptionTerminalResult = "Grant"
	ExceptionTerminalDeny  ExceptionTerminalResult = "Deny"
)

// Valid reports whether r is a closed ExceptionTerminalResult.
func (r ExceptionTerminalResult) Valid() bool {
	switch r {
	case ExceptionTerminalGrant, ExceptionTerminalDeny:
		return true
	default:
		return false
	}
}

// ApprovalDecisionFacts is an immutable deep-copied evidence-carrier input for
// approval-decision/v1 (design §3.2; F18-RD-19). It is not a resource,
// DecisionRecord, writer, route, or public contract.
type ApprovalDecisionFacts struct {
	sealed             bool
	policyVersion      string
	subjectDigest      []byte
	proposalDigest     []byte
	activeStage        string
	eligiblePrincipals []PrincipalRef
	countedDecisions   int
	result             ApprovalTerminalResult
	expiresAt          time.Time
}

// NewApprovalDecisionFacts constructs an immutable deep-copied carrier input.
func NewApprovalDecisionFacts(
	policyVersion string,
	subjectDigest []byte,
	proposalDigest []byte,
	activeStage string,
	eligiblePrincipals []PrincipalRef,
	countedDecisions int,
	result ApprovalTerminalResult,
	expiresAt time.Time,
) (ApprovalDecisionFacts, bool) {
	if policyVersion == "" || activeStage == "" || !result.Valid() || expiresAt.IsZero() {
		return ApprovalDecisionFacts{}, false
	}
	if len(subjectDigest) == 0 || len(proposalDigest) == 0 {
		return ApprovalDecisionFacts{}, false
	}
	if countedDecisions < 0 {
		return ApprovalDecisionFacts{}, false
	}
	return ApprovalDecisionFacts{
		sealed:             true,
		policyVersion:      policyVersion,
		subjectDigest:      copyBytes(subjectDigest),
		proposalDigest:     copyBytes(proposalDigest),
		activeStage:        activeStage,
		eligiblePrincipals: copyPrincipals(eligiblePrincipals),
		countedDecisions:   countedDecisions,
		result:             result,
		expiresAt:          expiresAt.UTC(),
	}, true
}

// Sealed reports whether f was constructed through NewApprovalDecisionFacts.
func (f ApprovalDecisionFacts) Sealed() bool { return f.sealed }

// PolicyVersion returns the sealed policy version.
func (f ApprovalDecisionFacts) PolicyVersion() string { return f.policyVersion }

// SubjectDigest returns a defensive copy of the subject digest.
func (f ApprovalDecisionFacts) SubjectDigest() []byte { return copyBytes(f.subjectDigest) }

// ProposalDigest returns a defensive copy of the proposal digest.
func (f ApprovalDecisionFacts) ProposalDigest() []byte { return copyBytes(f.proposalDigest) }

// ActiveStage returns the sealed active stage.
func (f ApprovalDecisionFacts) ActiveStage() string { return f.activeStage }

// EligiblePrincipals returns a defensive copy of eligible principals.
func (f ApprovalDecisionFacts) EligiblePrincipals() []PrincipalRef {
	return copyPrincipals(f.eligiblePrincipals)
}

// CountedDecisions returns the sealed counted-decisions value.
func (f ApprovalDecisionFacts) CountedDecisions() int { return f.countedDecisions }

// Result returns the sealed terminal approval result.
func (f ApprovalDecisionFacts) Result() ApprovalTerminalResult { return f.result }

// ExpiresAt returns the sealed ApprovalRequest expiresAt bound.
func (f ApprovalDecisionFacts) ExpiresAt() time.Time { return f.expiresAt }

// Equal reports byte-stable equality of sealed facts.
func (f ApprovalDecisionFacts) Equal(other ApprovalDecisionFacts) bool {
	if !f.sealed || !other.sealed {
		return false
	}
	if f.policyVersion != other.policyVersion ||
		f.activeStage != other.activeStage ||
		f.countedDecisions != other.countedDecisions ||
		f.result != other.result ||
		!f.expiresAt.Equal(other.expiresAt) ||
		!bytes.Equal(f.subjectDigest, other.subjectDigest) ||
		!bytes.Equal(f.proposalDigest, other.proposalDigest) {
		return false
	}
	if len(f.eligiblePrincipals) != len(other.eligiblePrincipals) {
		return false
	}
	for i := range f.eligiblePrincipals {
		if !f.eligiblePrincipals[i].Equal(other.eligiblePrincipals[i]) {
			return false
		}
	}
	return true
}

// ExceptionDecisionFacts is an immutable deep-copied evidence-carrier input for
// exception-decision/v1 (design §3.2; F18-RD-19). It is not a resource,
// DecisionRecord, writer, route, or public contract.
type ExceptionDecisionFacts struct {
	sealed               bool
	controlVersion       string
	subjectUID           string
	scopeRef             apimeta.ScopeRef
	notBefore            time.Time
	expiresAt            time.Time
	typedOverride        string
	compensatingControls []string
	approvalEvidence     []byte
	result               ExceptionTerminalResult
	terminalReason       string
	proposalUID          string
	grantUID             string // set only for Grant; empty for Deny
}

// NewExceptionDecisionFacts constructs an immutable deep-copied carrier input.
func NewExceptionDecisionFacts(
	controlVersion string,
	subjectUID string,
	scopeRef apimeta.ScopeRef,
	notBefore time.Time,
	expiresAt time.Time,
	typedOverride string,
	compensatingControls []string,
	approvalEvidence []byte,
	result ExceptionTerminalResult,
	terminalReason string,
	proposalUID string,
	grantUID string,
) (ExceptionDecisionFacts, bool) {
	if controlVersion == "" || subjectUID == "" || typedOverride == "" ||
		terminalReason == "" || proposalUID == "" || !result.Valid() {
		return ExceptionDecisionFacts{}, false
	}
	if scopeRef.Kind == "" || !apimeta.ScopeKind(scopeRef.Kind).Valid() {
		return ExceptionDecisionFacts{}, false
	}
	if notBefore.IsZero() || expiresAt.IsZero() || !expiresAt.After(notBefore) {
		return ExceptionDecisionFacts{}, false
	}
	if result == ExceptionTerminalGrant && grantUID == "" {
		return ExceptionDecisionFacts{}, false
	}
	if result == ExceptionTerminalDeny && grantUID != "" {
		return ExceptionDecisionFacts{}, false
	}
	return ExceptionDecisionFacts{
		sealed:               true,
		controlVersion:       controlVersion,
		subjectUID:           subjectUID,
		scopeRef:             scopeRef,
		notBefore:            notBefore.UTC(),
		expiresAt:            expiresAt.UTC(),
		typedOverride:        typedOverride,
		compensatingControls: copyStrings(compensatingControls),
		approvalEvidence:     copyBytes(approvalEvidence),
		result:               result,
		terminalReason:       terminalReason,
		proposalUID:          proposalUID,
		grantUID:             grantUID,
	}, true
}

// Sealed reports whether f was constructed through NewExceptionDecisionFacts.
func (f ExceptionDecisionFacts) Sealed() bool { return f.sealed }

// ControlVersion returns the sealed control version.
func (f ExceptionDecisionFacts) ControlVersion() string { return f.controlVersion }

// SubjectUID returns the sealed subject UID.
func (f ExceptionDecisionFacts) SubjectUID() string { return f.subjectUID }

// ScopeRef returns the sealed scope reference.
func (f ExceptionDecisionFacts) ScopeRef() apimeta.ScopeRef { return f.scopeRef }

// NotBefore returns the sealed interval start.
func (f ExceptionDecisionFacts) NotBefore() time.Time { return f.notBefore }

// ExpiresAt returns the sealed interval end.
func (f ExceptionDecisionFacts) ExpiresAt() time.Time { return f.expiresAt }

// TypedOverride returns the sealed typed override identifier.
func (f ExceptionDecisionFacts) TypedOverride() string { return f.typedOverride }

// CompensatingControls returns a defensive copy of compensating controls.
func (f ExceptionDecisionFacts) CompensatingControls() []string {
	return copyStrings(f.compensatingControls)
}

// ApprovalEvidence returns a defensive copy of approval evidence bytes.
func (f ExceptionDecisionFacts) ApprovalEvidence() []byte {
	return copyBytes(f.approvalEvidence)
}

// Result returns the sealed terminal exception result.
func (f ExceptionDecisionFacts) Result() ExceptionTerminalResult { return f.result }

// TerminalReason returns the sealed terminal reason.
func (f ExceptionDecisionFacts) TerminalReason() string { return f.terminalReason }

// ProposalUID returns the sealed ExceptionProposal UID.
func (f ExceptionDecisionFacts) ProposalUID() string { return f.proposalUID }

// GrantUID returns the sealed ExceptionGrant UID when result is Grant.
func (f ExceptionDecisionFacts) GrantUID() string { return f.grantUID }

// Equal reports byte-stable equality of sealed facts.
func (f ExceptionDecisionFacts) Equal(other ExceptionDecisionFacts) bool {
	if !f.sealed || !other.sealed {
		return false
	}
	return f.controlVersion == other.controlVersion &&
		f.subjectUID == other.subjectUID &&
		f.scopeRef == other.scopeRef &&
		f.notBefore.Equal(other.notBefore) &&
		f.expiresAt.Equal(other.expiresAt) &&
		f.typedOverride == other.typedOverride &&
		stringSliceEqual(f.compensatingControls, other.compensatingControls) &&
		bytes.Equal(f.approvalEvidence, other.approvalEvidence) &&
		f.result == other.result &&
		f.terminalReason == other.terminalReason &&
		f.proposalUID == other.proposalUID &&
		f.grantUID == other.grantUID
}

// MutationEventFacts is an immutable deep-copied evidence-carrier input for
// mutation AuditEvent descriptors (design §3.2). It contains only registered
// event taxonomy, resource/version linkage, action, and after-state digest.
type MutationEventFacts struct {
	sealed           bool
	eventType        string
	resourceRef      apimeta.TypedRef
	resourceVersion  string
	action           string
	afterStateDigest []byte
}

// NewMutationEventFacts constructs an immutable deep-copied carrier input.
func NewMutationEventFacts(
	eventType string,
	resourceRef apimeta.TypedRef,
	resourceVersion string,
	action string,
	afterStateDigest []byte,
) (MutationEventFacts, bool) {
	if eventType == "" || action == "" || resourceVersion == "" {
		return MutationEventFacts{}, false
	}
	if resourceRef.APIVersion == "" || resourceRef.Kind == "" || resourceRef.Name == "" {
		return MutationEventFacts{}, false
	}
	if len(afterStateDigest) == 0 {
		return MutationEventFacts{}, false
	}
	return MutationEventFacts{
		sealed:           true,
		eventType:        eventType,
		resourceRef:      resourceRef,
		resourceVersion:  resourceVersion,
		action:           action,
		afterStateDigest: copyBytes(afterStateDigest),
	}, true
}

// Sealed reports whether f was constructed through NewMutationEventFacts.
func (f MutationEventFacts) Sealed() bool { return f.sealed }

// EventType returns the sealed taxonomy event type.
func (f MutationEventFacts) EventType() string { return f.eventType }

// ResourceRef returns the sealed resource reference.
func (f MutationEventFacts) ResourceRef() apimeta.TypedRef { return f.resourceRef }

// ResourceVersion returns the sealed resource version.
func (f MutationEventFacts) ResourceVersion() string { return f.resourceVersion }

// Action returns the sealed action.
func (f MutationEventFacts) Action() string { return f.action }

// AfterStateDigest returns a defensive copy of the after-state digest.
func (f MutationEventFacts) AfterStateDigest() []byte {
	return copyBytes(f.afterStateDigest)
}

// Equal reports byte-stable equality of sealed facts.
func (f MutationEventFacts) Equal(other MutationEventFacts) bool {
	if !f.sealed || !other.sealed {
		return false
	}
	return f.eventType == other.eventType &&
		f.resourceRef == other.resourceRef &&
		f.resourceVersion == other.resourceVersion &&
		f.action == other.action &&
		bytes.Equal(f.afterStateDigest, other.afterStateDigest)
}

func copyBytes(in []byte) []byte {
	if len(in) == 0 {
		return nil
	}
	out := make([]byte, len(in))
	copy(out, in)
	return out
}

func copyStrings(in []string) []string {
	if len(in) == 0 {
		return nil
	}
	out := make([]string, len(in))
	copy(out, in)
	return out
}

func copyPrincipals(in []PrincipalRef) []PrincipalRef {
	if len(in) == 0 {
		return nil
	}
	out := make([]PrincipalRef, len(in))
	copy(out, in)
	return out
}

func stringSliceEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
