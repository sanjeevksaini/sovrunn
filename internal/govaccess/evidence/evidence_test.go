package evidence_test

import (
	"encoding/json"
	"go/parser"
	"go/token"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/apivalid"
	"github.com/sanjeevksaini/sovrunn/internal/decision/bundle"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/evidence"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/model"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/operation"
)

func mustCtx(t *testing.T) evidence.PublicationContext {
	t.Helper()
	ctx, fail := evidence.NewPublicationContext(time.Date(2026, 9, 4, 12, 0, 0, 0, time.UTC), 1)
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	return ctx
}

func mustAllowCandidate(t *testing.T, ctx evidence.PublicationContext, outcome evidence.AuthorizationOutcome) (evidence.FinalizedAuthorizationCandidate, evidence.AuthorizationAdmissionProof) {
	t.Helper()
	cd, fail := evidence.NewCandidateDigest([]byte("cand-1"))
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	dd, fail := evidence.NewDependencyVersionDigest([]byte("dep-1"))
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	desc, fail := evidence.NewAuthorizationCandidateDescriptor(cd, outcome, "roleassignment.create", "ra-1")
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	adm, fail := evidence.BindAuthorizationAdmission(cd, dd, ctx)
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	fin, fail := evidence.NewFinalizedAuthorizationCandidate(desc, adm, ctx)
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	return fin, adm
}

func TestRequireCurrentAllowAllowAndDeny(t *testing.T) {
	t.Parallel()
	ctx := mustCtx(t)
	fin, adm := mustAllowCandidate(t, ctx, evidence.AuthorizationAllow)
	out := evidence.RequireCurrentAllow(fin, adm, ctx)
	if _, ok := out.AsAllow(); !ok {
		t.Fatal("expected Allow proof")
	}
	finDeny, admDeny := mustAllowCandidate(t, ctx, evidence.AuthorizationDeny)
	out = evidence.RequireCurrentAllow(finDeny, admDeny, ctx)
	if _, ok := out.AsDenyRoute(); !ok {
		t.Fatal("expected RouteDenyToEvidenceOnly")
	}
	if _, ok := out.AsAllow(); ok {
		t.Fatal("Deny must not yield CurrentAllowProof")
	}
}

func TestDecisionMaterialsConsumeCarrierInputs(t *testing.T) {
	t.Parallel()
	ctx := mustCtx(t)
	af, ok := model.NewApprovalDecisionFacts(
		"v1", []byte("s"), []byte("p"), "stage", nil, 1,
		model.ApprovalTerminalApproved, time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC),
	)
	if !ok {
		t.Fatal("approval facts")
	}
	mat, fail := evidence.NewApprovalDecisionMaterial(af, ctx)
	if fail.Reason() != "" || mat.ProfileID() != evidence.ProfileApprovalDecision {
		t.Fatalf("approval material: %s", fail.Reason())
	}
	got, ok := mat.ApprovalFacts()
	if !ok || !got.Equal(af) {
		t.Fatal("material must consume ApprovalDecisionFacts")
	}

	ef, ok := model.NewExceptionDecisionFacts(
		"cv", "subj",
		apimeta.ScopeRef{TypedRef: apimeta.TypedRef{APIVersion: "core.sovrunn.io/v1alpha1", Kind: "Project", Name: "p"}},
		time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC),
		"ov", []string{"c"}, []byte("ae"),
		model.ExceptionTerminalDeny, "reason", "prop-1", "",
	)
	if !ok {
		t.Fatal("exception facts")
	}
	emat, fail := evidence.NewExceptionDecisionMaterial(ef, ctx)
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	egot, ok := emat.ExceptionFacts()
	if !ok || !egot.Equal(ef) {
		t.Fatal("material must consume ExceptionDecisionFacts")
	}
	desc, fail := evidence.NewDomainConclusionDescriptor(ef, ctx)
	if fail.Reason() != "" || desc.EventType() != evidence.EventTypeExceptionProposalDecided {
		t.Fatal(fail.Reason())
	}
	fail = evidence.ValidateExceptionTerminalDescriptorSet(ef, nil, &desc, []evidence.DomainDecisionMaterial{emat}, ctx)
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
}

func TestMutationDescriptorOrderingAndTerminalGrant(t *testing.T) {
	t.Parallel()
	ctx := mustCtx(t)
	mf1, ok := model.NewMutationEventFacts(
		evidence.EventTypeExceptionProposalDecided,
		apimeta.TypedRef{APIVersion: "gov.sovrunn.io/v1alpha1", Kind: "ExceptionProposal", Name: "ep", UID: "ep-1"},
		"1", "decide", []byte("d1"),
	)
	if !ok {
		t.Fatal("facts1")
	}
	mf2, ok := model.NewMutationEventFacts(
		evidence.EventTypeExceptionGrantIssued,
		apimeta.TypedRef{APIVersion: "gov.sovrunn.io/v1alpha1", Kind: "ExceptionGrant", Name: "eg", UID: "eg-1"},
		"1", "issue", []byte("d2"),
	)
	if !ok {
		t.Fatal("facts2")
	}
	d1, fail := evidence.NewMutationDescriptor(mf1, ctx)
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	d2, fail := evidence.NewMutationDescriptor(mf2, ctx)
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	ef, ok := model.NewExceptionDecisionFacts(
		"cv", "subj",
		apimeta.ScopeRef{TypedRef: apimeta.TypedRef{APIVersion: "core.sovrunn.io/v1alpha1", Kind: "Project", Name: "p"}},
		time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC),
		"ov", nil, nil,
		model.ExceptionTerminalGrant, "reason", "prop-1", "grant-1",
	)
	if !ok {
		t.Fatal("exception grant facts")
	}
	mat, fail := evidence.NewExceptionDecisionMaterial(ef, ctx)
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	// reverse order input; validator sorts
	fail = evidence.ValidateExceptionTerminalDescriptorSet(ef, []evidence.MutationDescriptor{d2, d1}, nil, []evidence.DomainDecisionMaterial{mat}, ctx)
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
}

func TestCarrierPlans(t *testing.T) {
	t.Parallel()
	ctx := mustCtx(t)
	fin, adm := mustAllowCandidate(t, ctx, evidence.AuthorizationAllow)
	eval, fail := evidence.CompleteAuthorizationCandidate(fin, adm, ctx)
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	mf, ok := model.NewMutationEventFacts(
		"roleassignment.created",
		apimeta.TypedRef{APIVersion: "gov.sovrunn.io/v1alpha1", Kind: "RoleAssignment", Name: "ra", UID: "ra-1"},
		"1", "create", []byte("after"),
	)
	if !ok {
		t.Fatal("mutation facts")
	}
	desc, fail := evidence.NewMutationDescriptor(mf, ctx)
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	pb, fail := operation.NewParticipantBinding([]byte("part"))
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	pub, fail := operation.NewPublicationBinding([]byte("pub"))
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	sh, fail := operation.NewShadowDeltaDigest([]byte("sh"))
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	ix, fail := operation.NewIndexDeltaDigest([]byte("ix"))
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	link, fail := operation.NewAppliedChangeLink(pb, pub, sh, ix, operation.ResultMaterial{}, false)
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	proof, fail := evidence.BindAppliedChange(link, []evidence.MutationDescriptor{desc}, nil)
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	aplan, fail := evidence.NewAuthorizedMutationPlan(ctx, eval, []evidence.AppliedMutationProof{proof})
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	prep, fail := evidence.PrepareMutationCarrierSet(&aplan, nil)
	if fail.Reason() != "" || !prep.Sealed() {
		t.Fatal(fail.Reason())
	}
	auto, fail := evidence.NewAutomaticMutationPlan(ctx, []evidence.AppliedMutationProof{proof})
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	prep2, fail := evidence.PrepareMutationCarrierSet(nil, &auto)
	if fail.Reason() != "" || !prep2.Sealed() {
		t.Fatal(fail.Reason())
	}
}

func TestProfilesStrictLocalThreeOnly(t *testing.T) {
	t.Parallel()
	raw := evidence.LocalDecisionProfileBundleBytes()
	if len(raw) == 0 {
		t.Fatal("expected profile bundle bytes")
	}
	view, prob := evidence.ValidateLocalDecisionProfileBundle(raw)
	if prob != nil {
		t.Fatalf("validate: %#v", prob)
	}
	ids := evidence.ExactAdoptingProfileIDs()
	if len(ids) != 3 {
		t.Fatalf("expected exactly three profiles, got %d", len(ids))
	}
	for _, id := range ids {
		sem, ok := evidence.LookupProfileSemantics(id, evidence.ProfileVersionV1)
		if !ok {
			t.Fatalf("missing semantics for %s", id)
		}
		if sem.Authority == "" || sem.InputSummary == "" || sem.ResultSummary == "" ||
			sem.ValiditySummary == "" || sem.ObligationSummary == "" ||
			sem.ProjectionSummary == "" || sem.AuditLinkageSummary == "" {
			t.Fatalf("incomplete F18-RD-19 semantics for %s", id)
		}
		if _, ok := view.Profile(id, evidence.ProfileVersionV1); !ok {
			t.Fatalf("view missing %s/%s", id, evidence.ProfileVersionV1)
		}
	}
	if _, ok := evidence.LookupProfileSemantics("fourth-profile", evidence.ProfileVersionV1); ok {
		t.Fatal("fourth profile must be absent")
	}
	if _, ok := view.Profile("authorization-decision", "v2"); ok {
		t.Fatal("unknown version must be absent")
	}
}

func TestProfilesDuplicateAndMalformedFail(t *testing.T) {
	t.Parallel()
	raw := evidence.LocalDecisionProfileBundleBytes()
	var doc bundle.Bundle
	if err := json.Unmarshal(raw, &doc); err != nil {
		t.Fatal(err)
	}
	doc.Spec.Profiles = append(doc.Spec.Profiles, doc.Spec.Profiles[0])
	dup, err := json.Marshal(doc)
	if err != nil {
		t.Fatal(err)
	}
	_, prob := bundle.Load(dup, apivalid.ModeReadRepresentation)
	if prob == nil {
		t.Fatal("duplicate (id,version) must fail FEATURE-0013 validation")
	}
	_, prob = bundle.Load([]byte(`{"apiVersion":`), apivalid.ModeReadRepresentation)
	if prob == nil {
		t.Fatal("malformed profile bundle must fail")
	}
}

func TestEvidenceConstructsBytesWithoutRegistration(t *testing.T) {
	t.Parallel()
	// Constructing bytes twice yields usable independent Load results; no global registry.
	a := evidence.LocalDecisionProfileBundleBytes()
	b := evidence.LocalDecisionProfileBundleBytes()
	va, pa := bundle.Load(a, apivalid.ModeReadRepresentation)
	vb, pb := bundle.Load(b, apivalid.ModeReadRepresentation)
	if pa != nil || pb != nil {
		t.Fatal("local load must succeed independently")
	}
	if _, ok := va.View().Profile(evidence.ProfileAuthorizationDecision, evidence.ProfileVersionV1); !ok {
		t.Fatal("missing authz profile")
	}
	if _, ok := vb.View().Profile(evidence.ProfileApprovalDecision, evidence.ProfileVersionV1); !ok {
		t.Fatal("missing approval profile")
	}
}

func TestEvidenceNeverImportsStateUowAuthzeval(t *testing.T) {
	t.Parallel()
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("caller")
	}
	dir := filepath.Dir(thisFile)
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dir, nil, 0)
	if err != nil {
		t.Fatal(err)
	}
	forbidden := []string{"/govaccess/state", "/govaccess/uow", "/govaccess/authzeval"}
	for _, pkg := range pkgs {
		for name, f := range pkg.Files {
			if strings.HasSuffix(name, "_test.go") {
				continue
			}
			for _, imp := range f.Imports {
				path := strings.Trim(imp.Path.Value, `"`)
				for _, bad := range forbidden {
					if strings.Contains(path, bad) {
						t.Fatalf("%s imports forbidden package %s", name, path)
					}
				}
			}
		}
	}
}
