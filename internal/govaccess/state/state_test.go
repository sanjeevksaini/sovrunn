package state_test

import (
	"testing"
	"time"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/clock"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/evidence"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/idempotency"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/model"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/operation"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/state"
)

type fakeChange struct {
	claim       state.ParticipantClaim
	binding     state.PermitBinding
	ctxBind     evidence.PublicationContextBinding
	kind        state.PublicationKind
	descriptors []evidence.MutationDescriptor
	conclusion  evidence.DomainConclusionDescriptor
	hasConc     bool
	materials   []evidence.DomainDecisionMaterial
	result      operation.ResultMaterial
	hasResult   bool
	putUID      string
	putVer      string
}

func (f fakeChange) ParticipantClaim() state.ParticipantClaim { return f.claim }
func (f fakeChange) ContextBinding() evidence.PublicationContextBinding {
	return f.ctxBind
}
func (f fakeChange) PermitBinding() state.PermitBinding     { return f.binding }
func (f fakeChange) PublicationKind() state.PublicationKind { return f.kind }
func (f fakeChange) EvidenceDescriptors() []evidence.MutationDescriptor {
	return f.descriptors
}
func (f fakeChange) ConclusionDescriptor() (evidence.DomainConclusionDescriptor, bool) {
	return f.conclusion, f.hasConc
}
func (f fakeChange) DomainDecisionMaterials() []evidence.DomainDecisionMaterial {
	return f.materials
}
func (f fakeChange) CanonicalResult() (operation.ResultMaterial, bool) {
	return f.result, f.hasResult
}
func (f fakeChange) ApplyTo(ed state.StateEditor) error {
	if f.kind == state.ResourceMutation {
		return ed.PutResource(f.putUID, f.putVer, []byte("payload"))
	}
	return nil
}

func seedCaller(t *testing.T) (*state.Store, state.CallerMutationTransaction, state.ParticipantClaim, evidence.CurrentAllowProof) {
	t.Helper()
	clk := clock.NewFixedClock(time.Date(2026, 9, 4, 12, 0, 0, 0, time.UTC))
	store := state.NewStore(clk)
	store.SetResourceVersion("ra-1", "1")

	versions, err := state.NewExpectedVersionSet([]state.ExpectedVersionPredicate{{ResourceUID: "ra-1", ResourceVersion: "1"}})
	if err != nil {
		t.Fatal(err)
	}
	claim, err := state.NewParticipantClaim(state.RoleAssignmentOwner, "p1", []byte("intent"), versions, state.OriginatingResult)
	if err != nil {
		t.Fatal(err)
	}
	cd, fail := evidence.NewCandidateDigest([]byte("cand"))
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	depSet, err := state.NewAuthorizationDependencyVersionSet(store.SnapshotVersion(), []state.ExpectedVersionPredicate{{ResourceUID: "ra-1", ResourceVersion: "1"}})
	if err != nil {
		t.Fatal(err)
	}
	curr, err := state.NewAuthorizationCurrentnessClaim(cd, depSet)
	if err != nil {
		t.Fatal(err)
	}
	adm, err := state.NewCallerMutationAdmission([]state.ParticipantClaim{claim}, curr)
	if err != nil {
		t.Fatal(err)
	}

	actor := model.PrincipalRef{Issuer: "https://idp.example", Subject: "u1", PrincipalType: model.PrincipalTypeHuman}
	key, fail := idempotency.NewIdempotencyLookupKey(actor, "roleassignment.create", "ck-1")
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	binding, fail := idempotency.NewRequestBinding(
		apimeta.TypedRef{APIVersion: "gov.sovrunn.io/v1alpha1", Kind: "RoleAssignment", Name: "ra-1", UID: "ra-1"},
		apimeta.ScopeRef{TypedRef: apimeta.TypedRef{APIVersion: "core.sovrunn.io/v1alpha1", Kind: "Project", Name: "p1"}},
		[]byte("req-digest"),
	)
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	out := store.InspectOrReserveCallerMutation(key, binding)
	lease, ok := out.Owner()
	if !ok {
		t.Fatal("expected owner lease")
	}
	tx, err := lease.Begin(adm)
	if err != nil {
		t.Fatal(err)
	}

	desc, fail := evidence.NewAuthorizationCandidateDescriptor(cd, evidence.AuthorizationAllow, "roleassignment.create", "ra-1")
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	fin, fail := evidence.NewFinalizedAuthorizationCandidate(desc, tx.AuthorizationAdmissionProof(), tx.PublicationContext())
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	req := evidence.RequireCurrentAllow(fin, tx.AuthorizationAdmissionProof(), tx.PublicationContext())
	allow, ok := req.AsAllow()
	if !ok {
		t.Fatal("expected current allow")
	}
	if fail := tx.AcceptCurrentAllow(allow); fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	return store, tx, claim, allow
}

func TestSealedTransactionLifecycleAndPermit(t *testing.T) {
	t.Parallel()
	_, tx, claim, _ := seedCaller(t)
	defer tx.Abort()

	permit, fail := tx.FinalizationPermit(claim)
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	if !permit.Binding().Sealed() || !permit.Sealed() {
		t.Fatal("permit must be sealed")
	}

	ctxBind, fail := evidence.BindPublicationContext(tx.PublicationContext())
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	mf, ok := model.NewMutationEventFacts(
		"roleassignment.created",
		apimeta.TypedRef{APIVersion: "gov.sovrunn.io/v1alpha1", Kind: "RoleAssignment", Name: "ra-1", UID: "ra-1"},
		"2", "create", []byte("after"),
	)
	if !ok {
		t.Fatal("mutation facts")
	}
	mdesc, fail := evidence.NewMutationDescriptor(mf, tx.PublicationContext())
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	rm, fail := operation.NewResultMaterial([]byte(`{"ok":true}`))
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	change := fakeChange{
		claim: claim, binding: permit.Binding(), ctxBind: ctxBind,
		kind: state.ResourceMutation, descriptors: []evidence.MutationDescriptor{mdesc},
		result: rm, hasResult: true, putUID: "ra-1", putVer: "2",
	}
	receipt, fail := tx.ApplyFinalized(change)
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	if receipt.PublicationKind() != state.ResourceMutation {
		t.Fatal("receipt kind")
	}
	if _, ok := receipt.EvidenceProof(); !ok {
		t.Fatal("expected mutation proof")
	}

	forged := fakeChange{
		claim: claim, binding: state.PermitBinding{}, ctxBind: ctxBind,
		kind: state.ResourceMutation, descriptors: []evidence.MutationDescriptor{mdesc},
		putUID: "ra-1", putVer: "3",
	}
	if _, fail := tx.ApplyFinalized(forged); fail.Reason() == "" {
		t.Fatal("forged permit binding must fail")
	}
}

func TestPermitUnforgeableAndAuthzGate(t *testing.T) {
	t.Parallel()
	clk := clock.NewFixedClock(time.Date(2026, 9, 4, 12, 0, 0, 0, time.UTC))
	store := state.NewStore(clk)
	store.SetResourceVersion("ra-1", "1")
	versions, err := state.NewExpectedVersionSet([]state.ExpectedVersionPredicate{{ResourceUID: "ra-1", ResourceVersion: "1"}})
	if err != nil {
		t.Fatal(err)
	}
	claim, err := state.NewParticipantClaim(state.RoleAssignmentOwner, "p1", []byte("intent"), versions, state.OriginatingResult)
	if err != nil {
		t.Fatal(err)
	}
	cd, fail := evidence.NewCandidateDigest([]byte("cand"))
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	depSet, err := state.NewAuthorizationDependencyVersionSet(store.SnapshotVersion(), []state.ExpectedVersionPredicate{{ResourceUID: "ra-1", ResourceVersion: "1"}})
	if err != nil {
		t.Fatal(err)
	}
	curr, err := state.NewAuthorizationCurrentnessClaim(cd, depSet)
	if err != nil {
		t.Fatal(err)
	}
	adm, err := state.NewCallerMutationAdmission([]state.ParticipantClaim{claim}, curr)
	if err != nil {
		t.Fatal(err)
	}
	actor := model.PrincipalRef{Issuer: "https://idp.example", Subject: "u2", PrincipalType: model.PrincipalTypeHuman}
	key, fail := idempotency.NewIdempotencyLookupKey(actor, "op", "k2")
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	binding, fail := idempotency.NewRequestBinding(
		apimeta.TypedRef{APIVersion: "gov.sovrunn.io/v1alpha1", Kind: "RoleAssignment", Name: "ra-1", UID: "ra-1"},
		apimeta.ScopeRef{TypedRef: apimeta.TypedRef{APIVersion: "core.sovrunn.io/v1alpha1", Kind: "Project", Name: "p1"}},
		[]byte("digest-2"),
	)
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	lease, ok := store.InspectOrReserveCallerMutation(key, binding).Owner()
	if !ok {
		t.Fatal("owner")
	}
	tx, err := lease.Begin(adm)
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Abort()
	if _, fail := tx.FinalizationPermit(claim); fail.Reason() == "" {
		t.Fatal("authorization-bearing tx must not issue permit before AcceptCurrentAllow")
	}
}

func TestDependencyVersionRevalidation(t *testing.T) {
	t.Parallel()
	clk := clock.NewFixedClock(time.Date(2026, 9, 4, 12, 0, 0, 0, time.UTC))
	store := state.NewStore(clk)
	store.SetResourceVersion("ra-1", "1")
	versions, err := state.NewExpectedVersionSet([]state.ExpectedVersionPredicate{{ResourceUID: "ra-1", ResourceVersion: "9"}})
	if err != nil {
		t.Fatal(err)
	}
	claim, err := state.NewParticipantClaim(state.RoleAssignmentOwner, "p1", []byte("intent"), versions, state.OriginatingResult)
	if err != nil {
		t.Fatal(err)
	}
	cd, fail := evidence.NewCandidateDigest([]byte("cand"))
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	depSet, err := state.NewAuthorizationDependencyVersionSet(store.SnapshotVersion(), []state.ExpectedVersionPredicate{{ResourceUID: "ra-1", ResourceVersion: "9"}})
	if err != nil {
		t.Fatal(err)
	}
	curr, err := state.NewAuthorizationCurrentnessClaim(cd, depSet)
	if err != nil {
		t.Fatal(err)
	}
	adm, err := state.NewCallerMutationAdmission([]state.ParticipantClaim{claim}, curr)
	if err != nil {
		t.Fatal(err)
	}
	actor := model.PrincipalRef{Issuer: "https://idp.example", Subject: "u3", PrincipalType: model.PrincipalTypeHuman}
	key, fail := idempotency.NewIdempotencyLookupKey(actor, "op", "k3")
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	binding, fail := idempotency.NewRequestBinding(
		apimeta.TypedRef{APIVersion: "gov.sovrunn.io/v1alpha1", Kind: "RoleAssignment", Name: "ra-1", UID: "ra-1"},
		apimeta.ScopeRef{TypedRef: apimeta.TypedRef{APIVersion: "core.sovrunn.io/v1alpha1", Kind: "Project", Name: "p1"}},
		[]byte("digest-3"),
	)
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	lease, ok := store.InspectOrReserveCallerMutation(key, binding).Owner()
	if !ok {
		t.Fatal("owner")
	}
	if _, err := lease.Begin(adm); err != state.ErrConflict {
		t.Fatalf("expected conflict on version mismatch, got %v", err)
	}
}

func TestTransactionSealCommitOnce(t *testing.T) {
	t.Parallel()
	_, tx, claim, _ := seedCaller(t)
	defer tx.Abort()

	permit, fail := tx.FinalizationPermit(claim)
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	ctxBind, fail := evidence.BindPublicationContext(tx.PublicationContext())
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	mf, ok := model.NewMutationEventFacts(
		"roleassignment.created",
		apimeta.TypedRef{APIVersion: "gov.sovrunn.io/v1alpha1", Kind: "RoleAssignment", Name: "ra-1", UID: "ra-1"},
		"2", "create", []byte("after"),
	)
	if !ok {
		t.Fatal("facts")
	}
	mdesc, fail := evidence.NewMutationDescriptor(mf, tx.PublicationContext())
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	rm, fail := operation.NewResultMaterial([]byte(`{"ok":true}`))
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	receipt, fail := tx.ApplyFinalized(fakeChange{
		claim: claim, binding: permit.Binding(), ctxBind: ctxBind,
		kind: state.ResourceMutation, descriptors: []evidence.MutationDescriptor{mdesc},
		result: rm, hasResult: true, putUID: "ra-1", putVer: "2",
	})
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}

	fin, fail := evidence.NewFinalizedAuthorizationCandidate(
		mustDesc(t, evidence.AuthorizationAllow),
		tx.AuthorizationAdmissionProof(),
		tx.PublicationContext(),
	)
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	eval, fail := evidence.CompleteAuthorizationCandidate(fin, tx.AuthorizationAdmissionProof(), tx.PublicationContext())
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	proof, ok := receipt.EvidenceProof()
	if !ok {
		t.Fatal("proof")
	}
	plan, fail := evidence.NewAuthorizedMutationPlan(tx.PublicationContext(), eval, []evidence.AppliedMutationProof{proof})
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	prep, fail := evidence.PrepareMutationCarrierSet(&plan, nil)
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	if fail := tx.AcceptPreparedEvidence(prep); fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}

	digest, fail := evidence.NewMutationDigestFromProofs([]evidence.AppliedMutationProof{proof})
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	cb, fail := evidence.CompletionBindingForPlan(tx.PublicationContext(), digest)
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	actor := model.PrincipalRef{Issuer: "https://idp.example", Subject: "u1", PrincipalType: model.PrincipalTypeHuman}
	key, fail := idempotency.NewIdempotencyLookupKey(actor, "roleassignment.create", "ck-1")
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	reqBinding, fail := idempotency.NewRequestBinding(
		apimeta.TypedRef{APIVersion: "gov.sovrunn.io/v1alpha1", Kind: "RoleAssignment", Name: "ra-1", UID: "ra-1"},
		apimeta.ScopeRef{TypedRef: apimeta.TypedRef{APIVersion: "core.sovrunn.io/v1alpha1", Kind: "Project", Name: "p1"}},
		[]byte("req-digest"),
	)
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	completed, fail := idempotency.NewPreparedCompletedResult(key, reqBinding, rm, cb)
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	if fail := tx.StageCompletedResult(completed); fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	if fail := tx.Seal(); fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	if fail := tx.Seal(); fail.Reason() == "" {
		t.Fatal("second seal must fail")
	}
	cr, fail := tx.Commit()
	if fail.Reason() != "" || !cr.Sealed() {
		t.Fatal(fail.Reason())
	}
	if _, fail := tx.Commit(); fail.Reason() == "" {
		t.Fatal("second commit must fail")
	}
}

func mustDesc(t *testing.T, outcome evidence.AuthorizationOutcome) evidence.AuthorizationCandidateDescriptor {
	t.Helper()
	cd, fail := evidence.NewCandidateDigest([]byte("cand"))
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	d, fail := evidence.NewAuthorizationCandidateDescriptor(cd, outcome, "roleassignment.create", "ra-1")
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	return d
}

func TestLockOrderingReservationThenState(t *testing.T) {
	t.Parallel()
	clk := clock.NewFixedClock(time.Date(2026, 9, 4, 12, 0, 0, 0, time.UTC))
	store := state.NewStore(clk)
	done := make(chan struct{})
	store.LockOrderObservation(func() {
		close(done)
	})
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("lock order observation timed out")
	}
}
