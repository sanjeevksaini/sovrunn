package state

import (
	"sync"

	"github.com/sanjeevksaini/sovrunn/internal/govaccess/evidence"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/idempotency"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/operation"
)

type transactionState uint8

const (
	txOpen transactionState = iota
	txSealed
	txCommitted
	txAborted
)

type baseTransaction struct {
	store *Store

	mu sync.Mutex

	state             transactionState
	admission         MutationAdmission
	context           evidence.PublicationContext
	admissionProof    evidence.AuthorizationAdmissionProof
	hasAdmissionProof bool

	currentAllow evidence.CurrentAllowProof
	hasAllow     bool

	permits  map[string]*permitRegistryEntry // binding hex -> entry
	receipts []AppliedChangeReceipt

	prepared    evidence.PreparedEvidenceChange
	hasPrepared bool

	completed    idempotency.PreparedCompletedResult
	hasCompleted bool

	editor *memoryEditor
}

func (t *baseTransaction) PublicationContext() evidence.PublicationContext {
	return t.context
}

func (t *baseTransaction) AuthorizationAdmissionProof() (evidence.AuthorizationAdmissionProof, bool) {
	return t.admissionProof, t.hasAdmissionProof
}

func (t *baseTransaction) AcceptCurrentAllow(proof evidence.CurrentAllowProof) operation.MechanicalFailure {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.state != txOpen {
		return operation.NewMechanicalFailure("tx_not_open")
	}
	if t.admission.authority != AuthorizationBearing {
		return operation.NewMechanicalFailure("accept_allow_forbidden_in_automatic_mode")
	}
	if !proof.Sealed() {
		return operation.NewMechanicalFailure("current_allow_unsealed")
	}
	if !proof.PublicationContext().Equal(t.context) {
		return operation.NewMechanicalFailure("current_allow_context_mismatch")
	}
	if !t.hasAdmissionProof || !proof.AdmissionDigest().Equal(t.admissionProof.AdmissionDigest()) {
		return operation.NewMechanicalFailure("current_allow_admission_mismatch")
	}
	t.currentAllow = proof
	t.hasAllow = true
	return operation.MechanicalFailure{}
}

func (t *baseTransaction) FinalizationPermit(claim ParticipantClaim) (MutationFinalizationPermit, operation.MechanicalFailure) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.state != txOpen {
		return MutationFinalizationPermit{}, operation.NewMechanicalFailure("tx_not_open")
	}
	if t.admission.authority == AuthorizationBearing && !t.hasAllow {
		return MutationFinalizationPermit{}, operation.NewMechanicalFailure("finalization_requires_current_allow")
	}
	if !claim.sealed {
		return MutationFinalizationPermit{}, operation.NewMechanicalFailure("finalization_claim_unsealed")
	}
	found := false
	for _, c := range t.admission.participants {
		if c.Equal(claim) {
			found = true
			break
		}
	}
	if !found {
		return MutationFinalizationPermit{}, operation.NewMechanicalFailure("finalization_claim_not_admitted")
	}
	token, fail := newPermitToken()
	if fail.Reason() != "" {
		return MutationFinalizationPermit{}, fail
	}
	binding := mintPermitBinding(token)
	key := bindingHex(binding)
	if t.permits == nil {
		t.permits = map[string]*permitRegistryEntry{}
	}
	t.permits[key] = &permitRegistryEntry{claim: claim, binding: binding, consumed: false}
	return MutationFinalizationPermit{
		sealed:    true,
		claim:     claim,
		context:   t.context,
		authority: t.admission.authority,
		binding:   binding,
		token:     token,
	}, operation.MechanicalFailure{}
}

func (t *baseTransaction) ApplyFinalized(change ContextBoundChange) (AppliedChangeReceipt, operation.MechanicalFailure) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.state != txOpen {
		return AppliedChangeReceipt{}, operation.NewMechanicalFailure("tx_not_open")
	}
	if t.admission.authority == AuthorizationBearing && !t.hasAllow {
		return AppliedChangeReceipt{}, operation.NewMechanicalFailure("apply_requires_current_allow")
	}
	binding := change.PermitBinding()
	entry, ok := t.permits[bindingHex(binding)]
	if !ok || entry.consumed {
		return AppliedChangeReceipt{}, operation.NewMechanicalFailure("permit_binding_unforgeable_or_consumed")
	}
	if !entry.claim.Equal(change.ParticipantClaim()) {
		return AppliedChangeReceipt{}, operation.NewMechanicalFailure("permit_claim_mismatch")
	}
	ctxBind := change.ContextBinding()
	expected, fail := evidence.BindPublicationContext(t.context)
	if fail.Reason() != "" {
		return AppliedChangeReceipt{}, fail
	}
	if !ctxBind.Equal(expected) {
		return AppliedChangeReceipt{}, operation.NewMechanicalFailure("apply_context_binding_mismatch")
	}
	kind := change.PublicationKind()
	if !kind.Valid() {
		return AppliedChangeReceipt{}, operation.NewMechanicalFailure("publication_kind_invalid")
	}
	if err := change.ApplyTo(t.editor); err != nil {
		return AppliedChangeReceipt{}, operation.NewMechanicalFailure("apply_to_failed")
	}
	entry.consumed = true

	receipt := AppliedChangeReceipt{
		sealed:  true,
		claim:   change.ParticipantClaim(),
		context: t.context,
		binding: binding,
		kind:    kind,
	}

	switch kind {
	case ResourceMutation:
		if t.editor.deltaEmpty() {
			return AppliedChangeReceipt{}, operation.NewMechanicalFailure("resource_mutation_delta_empty")
		}
		descriptors := change.EvidenceDescriptors()
		if len(descriptors) == 0 {
			return AppliedChangeReceipt{}, operation.NewMechanicalFailure("resource_mutation_descriptors_empty")
		}
		participant, fail := operation.NewParticipantBinding(change.ParticipantClaim().IntentDigest())
		if fail.Reason() != "" {
			return AppliedChangeReceipt{}, fail
		}
		publication, fail := operation.NewPublicationBinding(expected.Digest())
		if fail.Reason() != "" {
			return AppliedChangeReceipt{}, fail
		}
		shadow, fail := operation.NewShadowDeltaDigest(t.editor.ShadowDigest())
		if fail.Reason() != "" {
			return AppliedChangeReceipt{}, fail
		}
		index, fail := operation.NewIndexDeltaDigest(t.editor.IndexDigest())
		if fail.Reason() != "" {
			return AppliedChangeReceipt{}, fail
		}
		result, hasResult := change.CanonicalResult()
		link, fail := operation.NewAppliedChangeLink(participant, publication, shadow, index, result, hasResult)
		if fail.Reason() != "" {
			return AppliedChangeReceipt{}, fail
		}
		proof, fail := evidence.BindAppliedChange(link, descriptors, change.DomainDecisionMaterials())
		if fail.Reason() != "" {
			return AppliedChangeReceipt{}, fail
		}
		receipt.mutationProof = proof
		receipt.hasMutationProof = true
	case DomainConclusion:
		if !t.editor.deltaEmpty() {
			return AppliedChangeReceipt{}, operation.NewMechanicalFailure("domain_conclusion_delta_nonempty")
		}
		conclusion, ok := change.ConclusionDescriptor()
		if !ok || !conclusion.Sealed() {
			return AppliedChangeReceipt{}, operation.NewMechanicalFailure("domain_conclusion_descriptor_required")
		}
		participant, fail := operation.NewParticipantBinding(change.ParticipantClaim().IntentDigest())
		if fail.Reason() != "" {
			return AppliedChangeReceipt{}, fail
		}
		publication, fail := operation.NewPublicationBinding(expected.Digest())
		if fail.Reason() != "" {
			return AppliedChangeReceipt{}, fail
		}
		link, fail := operation.NewAppliedConclusionLink(participant, publication)
		if fail.Reason() != "" {
			return AppliedChangeReceipt{}, fail
		}
		proof, fail := evidence.BindAppliedConclusion(link, conclusion, change.DomainDecisionMaterials())
		if fail.Reason() != "" {
			return AppliedChangeReceipt{}, fail
		}
		receipt.conclusionProof = proof
		receipt.hasConclusion = true
	}
	t.receipts = append(t.receipts, receipt)
	return receipt, operation.MechanicalFailure{}
}

func (t *baseTransaction) StageCompletedResult(result idempotency.PreparedCompletedResult) operation.MechanicalFailure {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.state != txOpen {
		return operation.NewMechanicalFailure("tx_not_open")
	}
	if t.admission.mode != CallerMutationMode {
		return operation.NewMechanicalFailure("stage_completed_caller_only")
	}
	if !result.Sealed() {
		return operation.NewMechanicalFailure("completed_result_unsealed")
	}
	t.completed = result
	t.hasCompleted = true
	return operation.MechanicalFailure{}
}

func (t *baseTransaction) AcceptPreparedEvidence(change evidence.PreparedEvidenceChange) operation.MechanicalFailure {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.state != txOpen {
		return operation.NewMechanicalFailure("tx_not_open")
	}
	if !change.Sealed() {
		return operation.NewMechanicalFailure("prepared_evidence_unsealed")
	}
	if !change.PublicationContext().Equal(t.context) {
		return operation.NewMechanicalFailure("prepared_evidence_context_mismatch")
	}
	t.prepared = change
	t.hasPrepared = true
	return operation.MechanicalFailure{}
}

func (t *baseTransaction) sealCommon(requireCompleted bool) operation.MechanicalFailure {
	if t.state != txOpen {
		return operation.NewMechanicalFailure("tx_not_open")
	}
	if t.admission.authority == AuthorizationBearing {
		if !t.hasAllow || !t.hasAdmissionProof {
			return operation.NewMechanicalFailure("seal_requires_allow_chain")
		}
	}
	if len(t.receipts) != len(t.admission.participants) {
		return operation.NewMechanicalFailure("seal_receipt_cardinality")
	}
	if !t.hasPrepared {
		return operation.NewMechanicalFailure("seal_requires_prepared_evidence")
	}
	if requireCompleted && !t.hasCompleted {
		return operation.NewMechanicalFailure("seal_requires_completed_result")
	}
	if !requireCompleted && t.hasCompleted {
		return operation.NewMechanicalFailure("seal_forbids_completed_result")
	}
	t.state = txSealed
	return operation.MechanicalFailure{}
}

func (t *baseTransaction) commitLocked() (CommitReceipt, operation.MechanicalFailure) {
	if t.state != txSealed {
		return CommitReceipt{}, operation.NewMechanicalFailure("commit_requires_seal")
	}
	t.store.installRoot(t.editor.snapshot())
	if t.hasCompleted {
		t.store.storeCompleted(t.completed)
	}
	t.state = txCommitted
	return CommitReceipt{
		sealed:   true,
		context:  t.context,
		evidence: t.prepared,
		pins:     t.prepared.CarrierSet().DecisionProfiles(),
	}, operation.MechanicalFailure{}
}

func (t *baseTransaction) abortLocked() {
	if t.state == txCommitted || t.state == txAborted {
		return
	}
	t.state = txAborted
	t.store.releaseStateMu()
}

// CallerMutationTransaction is the sealed caller mutation transaction.
type CallerMutationTransaction struct {
	base *baseTransaction
}

func (t CallerMutationTransaction) PublicationContext() evidence.PublicationContext {
	return t.base.PublicationContext()
}
func (t CallerMutationTransaction) AuthorizationAdmissionProof() evidence.AuthorizationAdmissionProof {
	p, _ := t.base.AuthorizationAdmissionProof()
	return p
}
func (t CallerMutationTransaction) AcceptCurrentAllow(p evidence.CurrentAllowProof) operation.MechanicalFailure {
	return t.base.AcceptCurrentAllow(p)
}
func (t CallerMutationTransaction) FinalizationPermit(c ParticipantClaim) (MutationFinalizationPermit, operation.MechanicalFailure) {
	return t.base.FinalizationPermit(c)
}
func (t CallerMutationTransaction) ApplyFinalized(c ContextBoundChange) (AppliedChangeReceipt, operation.MechanicalFailure) {
	return t.base.ApplyFinalized(c)
}
func (t CallerMutationTransaction) StageCompletedResult(r idempotency.PreparedCompletedResult) operation.MechanicalFailure {
	return t.base.StageCompletedResult(r)
}
func (t CallerMutationTransaction) AcceptPreparedEvidence(c evidence.PreparedEvidenceChange) operation.MechanicalFailure {
	return t.base.AcceptPreparedEvidence(c)
}
func (t CallerMutationTransaction) Seal() operation.MechanicalFailure {
	t.base.mu.Lock()
	defer t.base.mu.Unlock()
	return t.base.sealCommon(true)
}
func (t CallerMutationTransaction) Commit() (CommitReceipt, operation.MechanicalFailure) {
	t.base.mu.Lock()
	defer t.base.mu.Unlock()
	receipt, fail := t.base.commitLocked()
	if fail.Reason() == "" {
		t.base.store.releaseStateMu()
	}
	return receipt, fail
}
func (t CallerMutationTransaction) Abort() {
	t.base.mu.Lock()
	defer t.base.mu.Unlock()
	t.base.abortLocked()
}

// ControllerMutationTransaction is the sealed controller mutation transaction.
type ControllerMutationTransaction struct {
	base *baseTransaction
}

func (t ControllerMutationTransaction) PublicationContext() evidence.PublicationContext {
	return t.base.PublicationContext()
}
func (t ControllerMutationTransaction) AuthorizationAdmissionProof() (evidence.AuthorizationAdmissionProof, bool) {
	return t.base.AuthorizationAdmissionProof()
}
func (t ControllerMutationTransaction) AcceptCurrentAllow(p evidence.CurrentAllowProof) operation.MechanicalFailure {
	return t.base.AcceptCurrentAllow(p)
}
func (t ControllerMutationTransaction) FinalizationPermit(c ParticipantClaim) (MutationFinalizationPermit, operation.MechanicalFailure) {
	return t.base.FinalizationPermit(c)
}
func (t ControllerMutationTransaction) ApplyFinalized(c ContextBoundChange) (AppliedChangeReceipt, operation.MechanicalFailure) {
	return t.base.ApplyFinalized(c)
}
func (t ControllerMutationTransaction) AcceptPreparedEvidence(c evidence.PreparedEvidenceChange) operation.MechanicalFailure {
	return t.base.AcceptPreparedEvidence(c)
}
func (t ControllerMutationTransaction) Seal() operation.MechanicalFailure {
	t.base.mu.Lock()
	defer t.base.mu.Unlock()
	return t.base.sealCommon(false)
}
func (t ControllerMutationTransaction) Commit() (CommitReceipt, operation.MechanicalFailure) {
	t.base.mu.Lock()
	defer t.base.mu.Unlock()
	receipt, fail := t.base.commitLocked()
	if fail.Reason() == "" {
		t.base.store.releaseStateMu()
	}
	return receipt, fail
}
func (t ControllerMutationTransaction) Abort() {
	t.base.mu.Lock()
	defer t.base.mu.Unlock()
	t.base.abortLocked()
}

// EvidenceTransaction is the sealed evidence-only transaction.
type EvidenceTransaction struct {
	base *baseTransaction
}

func (t EvidenceTransaction) PublicationContext() evidence.PublicationContext {
	return t.base.PublicationContext()
}
func (t EvidenceTransaction) AuthorizationAdmissionProof() evidence.AuthorizationAdmissionProof {
	p, _ := t.base.AuthorizationAdmissionProof()
	return p
}
func (t EvidenceTransaction) AcceptPreparedEvidence(c evidence.PreparedEvidenceChange) operation.MechanicalFailure {
	return t.base.AcceptPreparedEvidence(c)
}
func (t EvidenceTransaction) Seal() operation.MechanicalFailure {
	t.base.mu.Lock()
	defer t.base.mu.Unlock()
	if t.base.state != txOpen {
		return operation.NewMechanicalFailure("tx_not_open")
	}
	if !t.base.hasAdmissionProof || !t.base.hasPrepared {
		return operation.NewMechanicalFailure("evidence_seal_incomplete")
	}
	t.base.state = txSealed
	return operation.MechanicalFailure{}
}
func (t EvidenceTransaction) Commit() (CommitReceipt, operation.MechanicalFailure) {
	t.base.mu.Lock()
	defer t.base.mu.Unlock()
	receipt, fail := t.base.commitLocked()
	if fail.Reason() == "" {
		t.base.store.releaseStateMu()
	}
	return receipt, fail
}
func (t EvidenceTransaction) Abort() {
	t.base.mu.Lock()
	defer t.base.mu.Unlock()
	t.base.abortLocked()
}
