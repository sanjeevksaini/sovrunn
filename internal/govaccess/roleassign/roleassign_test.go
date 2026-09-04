package roleassign_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/clock"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/evidence"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/idempotency"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/model"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/roleassign"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/state"
)

func sampleSpec(t *testing.T, validity model.AssignmentValidityMode) model.RoleAssignmentSpec {
	t.Helper()
	holder := model.PrincipalRef{
		Issuer: "https://idp.example", Subject: "user-1", PrincipalType: model.PrincipalTypeHuman,
	}
	spec := model.RoleAssignmentSpec{
		RoleHolderRef: model.RoleHolderRef{
			Kind:      model.RoleHolderPrincipal,
			Principal: &holder,
		},
		RoleDefinitionRef: apimeta.TypedRef{
			APIVersion: model.APIVersionIAM, Kind: "RoleDefinition", Name: "admin", UID: "rd-1",
		},
		RoleDefinitionVersion: "v1",
		Validity:              validity,
	}
	if validity == model.AssignmentValidityTimeBound {
		nb := time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC)
		exp := time.Date(2026, 12, 1, 0, 0, 0, 0, time.UTC)
		spec.NotBefore = &nb
		spec.ExpiresAt = &exp
	} else {
		rule := apimeta.TypedRef{
			APIVersion: model.APIVersionGov, Kind: "AccessReviewRule", Name: "standing", UID: "arr-1",
		}
		spec.AccessReviewRuleRef = &rule
	}
	return spec
}

func sampleScope() apimeta.ScopeRef {
	return apimeta.ScopeRef{TypedRef: apimeta.TypedRef{
		APIVersion: "core.sovrunn.io/v1alpha1", Kind: string(apimeta.ScopeProject), Name: "p1", UID: "proj-1",
	}}
}

func mustGrantIntent(t *testing.T) roleassign.PreparedRoleAssignmentIntent {
	t.Helper()
	spec := sampleSpec(t, model.AssignmentValidityTimeBound)
	trigger, fail := roleassign.NewGrantTriggerIntent(
		"p-grant", "ra-1", "ra-one", sampleScope(), spec, "0", "1", nil, state.OriginatingResult, nil,
	)
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	intent, fail := roleassign.NewGrantPort().Submit(trigger)
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	return intent
}

func issuePermit(
	t *testing.T,
	claim state.ParticipantClaim,
) (state.CallerMutationTransaction, state.MutationFinalizationPermit) {
	t.Helper()
	clk := clock.NewFixedClock(time.Date(2026, 9, 4, 12, 0, 0, 0, time.UTC))
	store := state.NewStore(clk)
	for _, p := range claim.Versions().Predicates() {
		store.SetResourceVersion(p.ResourceUID, p.ResourceVersion)
	}
	cd, fail := evidence.NewCandidateDigest([]byte("cand-ra"))
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	depSet, err := state.NewAuthorizationDependencyVersionSet(
		store.SnapshotVersion(), claim.Versions().Predicates(),
	)
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
	actor := model.PrincipalRef{Issuer: "https://idp.example", Subject: "actor", PrincipalType: model.PrincipalTypeHuman}
	key, fail := idempotency.NewIdempotencyLookupKey(actor, "roleassignment.grant", "ck-ra")
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	binding, fail := idempotency.NewRequestBinding(
		apimeta.TypedRef{APIVersion: model.APIVersionIAM, Kind: roleassign.KindRoleAssignment, Name: "ra-1", UID: "ra-1"},
		sampleScope(),
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
	desc, fail := evidence.NewAuthorizationCandidateDescriptor(cd, evidence.AuthorizationAllow, "roleassignment.grant", "ra-1")
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
	permit, fail := tx.FinalizationPermit(claim)
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	return tx, permit
}

func TestPreparedRoleAssignmentIntentConstruction(t *testing.T) {
	t.Parallel()
	intent := mustGrantIntent(t)
	if !intent.Sealed() {
		t.Fatal("intent must be sealed")
	}
	if intent.Operation() != roleassign.OpGrant {
		t.Fatalf("operation=%s", intent.Operation())
	}
	if intent.TriggerPort() != "grant" {
		t.Fatalf("triggerPort=%s", intent.TriggerPort())
	}
	if intent.HasPublicationDerivedFields() {
		t.Fatal("intent must be publication-time-free")
	}
	if !intent.ParticipantClaim().Sealed() || intent.ParticipantClaim().Owner() != state.RoleAssignmentOwner {
		t.Fatal("claim must be RoleAssignmentOwner")
	}
	if _, fail := roleassign.NewGrantPort().Submit(roleassign.GrantTriggerIntent{}); fail.Reason() == "" {
		t.Fatal("unsealed trigger must fail")
	}
}

func TestFinalizeAtUsesPermitPublicationContext(t *testing.T) {
	t.Parallel()
	intent := mustGrantIntent(t)
	tx, permit := issuePermit(t, intent.ParticipantClaim())
	defer tx.Abort()

	outcome := intent.FinalizeAt(permit)
	change, ok := outcome.AsFinalized()
	if !ok {
		if fail, fOK := outcome.AsMechanicalFailure(); fOK {
			t.Fatalf("finalize mechanical: %s", fail.Reason())
		}
		if domain, dOK := outcome.AsDomain(); dOK {
			t.Fatalf("unexpected domain: %s", domain.Code())
		}
		t.Fatal("expected finalized change")
	}
	if !change.Sealed() {
		t.Fatal("change must be sealed")
	}
	if !change.PermitBinding().Equal(permit.Binding()) {
		t.Fatal("change must embed exact permit binding")
	}
	if change.PublicationKind() != state.ResourceMutation {
		t.Fatal("publication kind")
	}
	descs := change.EvidenceDescriptors()
	if len(descs) != 1 || descs[0].EventType() != "roleassignment.created" {
		t.Fatalf("descriptors=%v", descs)
	}
	if _, ok := change.CanonicalResult(); !ok {
		t.Fatal("originating result required")
	}
	ctxBind := change.ContextBinding()
	want, fail := evidence.BindPublicationContext(permit.PublicationContext())
	if fail.Reason() != "" || !ctxBind.Equal(want) {
		t.Fatal("context binding must match permit publication context")
	}
}

func TestFinalizeAtExpiredDomainOutcome(t *testing.T) {
	t.Parallel()
	spec := sampleSpec(t, model.AssignmentValidityTimeBound)
	past := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	exp := time.Date(2020, 2, 1, 0, 0, 0, 0, time.UTC)
	spec.NotBefore = &past
	spec.ExpiresAt = &exp
	trigger, fail := roleassign.NewGrantTriggerIntent(
		"p-exp", "ra-exp", "ra-exp", sampleScope(), spec, "0", "1", nil, state.OriginatingResult, nil,
	)
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	intent, fail := roleassign.NewGrantPort().Submit(trigger)
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	tx, permit := issuePermit(t, intent.ParticipantClaim())
	defer tx.Abort()
	outcome := intent.FinalizeAt(permit)
	domain, ok := outcome.AsDomain()
	if !ok || domain.Code() != "assignment_expired" {
		t.Fatalf("expected assignment_expired domain, got %#v kind=%s", domain, outcome.KindName())
	}
}

func TestSealedFinalizedChangeRejectsForgedPermit(t *testing.T) {
	t.Parallel()
	intent := mustGrantIntent(t)
	tx, permit := issuePermit(t, intent.ParticipantClaim())
	defer tx.Abort()
	outcome := intent.FinalizeAt(permit)
	change, ok := outcome.AsFinalized()
	if !ok {
		t.Fatal("expected finalized")
	}
	if change.PermitBinding().Equal(state.PermitBinding{}) {
		t.Fatal("permit binding must be sealed non-zero")
	}
	unsealed := roleassign.PreparedRoleAssignmentIntent{}
	if fail, ok := unsealed.FinalizeAt(permit).AsMechanicalFailure(); !ok || fail.Reason() != "intent_unsealed" {
		t.Fatalf("unsealed intent: %#v", fail)
	}
}

func TestThreeWorkflowTriggerPorts(t *testing.T) {
	t.Parallel()
	spec := sampleSpec(t, model.AssignmentValidityTimeBound)
	scope := sampleScope()

	grantTrigger, fail := roleassign.NewGrantTriggerIntent(
		"p1", "ra-g", "ra-g", scope, spec, "0", "1", nil, state.OriginatingResult, nil,
	)
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	grantIntent, fail := roleassign.NewGrantPort().Submit(grantTrigger)
	if fail.Reason() != "" || grantIntent.TriggerPort() != "grant" || grantIntent.Operation() != roleassign.OpGrant {
		t.Fatalf("grant port: fail=%s port=%s op=%s", fail.Reason(), grantIntent.TriggerPort(), grantIntent.Operation())
	}

	actTrigger, fail := roleassign.NewActivationTriggerIntent(
		roleassign.ActivationActivate, "p2", "ra-a", "ra-a", scope, spec, true, "0", "1", nil, state.OriginatingResult,
	)
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	actIntent, fail := roleassign.NewActivationPort().Submit(actTrigger)
	if fail.Reason() != "" || actIntent.TriggerPort() != "activation" || actIntent.Operation() != roleassign.OpActivate {
		t.Fatalf("activation port: fail=%s", fail.Reason())
	}

	revTrigger, fail := roleassign.NewActivationTriggerIntent(
		roleassign.ActivationPrivilegedRevoke, "p3", "ra-pr", "ra-pr", scope, model.RoleAssignmentSpec{}, false, "1", "2", nil, state.OriginatingResult,
	)
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	revIntent, fail := roleassign.NewActivationPort().Submit(revTrigger)
	if fail.Reason() != "" || revIntent.Operation() != roleassign.OpPrivilegedRevoke {
		t.Fatalf("privileged revoke: fail=%s op=%s", fail.Reason(), revIntent.Operation())
	}

	due := time.Date(2026, 11, 1, 0, 0, 0, 0, time.UTC)
	reviewTrigger, fail := roleassign.NewReviewRemediationTriggerIntent(
		roleassign.RemediationAdvanceReviewDue, "p4", "ra-r", "ra-r", scope, spec, true,
		model.RoleAssignmentSpec{}, false, "1", "2", nil, state.OriginatingResult, &due,
	)
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	reviewIntent, fail := roleassign.NewReviewRemediationPort().Submit(reviewTrigger)
	if fail.Reason() != "" || reviewIntent.TriggerPort() != "review" || reviewIntent.Operation() != roleassign.OpAdvanceReviewDue {
		t.Fatalf("review port: fail=%s", fail.Reason())
	}

	replaceTrigger, fail := roleassign.NewReviewRemediationTriggerIntent(
		roleassign.RemediationReplace, "p5", "ra-rep", "ra-rep", scope, spec, true, spec, true, "1", "2", nil, state.OriginatingResult, nil,
	)
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	replaceIntent, fail := roleassign.NewReviewRemediationPort().Submit(replaceTrigger)
	if fail.Reason() != "" || replaceIntent.Operation() != roleassign.OpReviewReplace {
		t.Fatalf("replace: fail=%s", fail.Reason())
	}

	revokeTrigger, fail := roleassign.NewReviewRemediationTriggerIntent(
		roleassign.RemediationRevoke, "p6", "ra-rv", "ra-rv", scope, spec, true,
		model.RoleAssignmentSpec{}, false, "1", "2", nil, state.OriginatingResult, nil,
	)
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	revokeIntent, fail := roleassign.NewReviewRemediationPort().Submit(revokeTrigger)
	if fail.Reason() != "" || revokeIntent.Operation() != roleassign.OpReviewRevoke {
		t.Fatalf("review revoke: fail=%s", fail.Reason())
	}

	if grantIntent.TriggerPort() == actIntent.TriggerPort() ||
		actIntent.TriggerPort() == reviewIntent.TriggerPort() ||
		grantIntent.TriggerPort() == reviewIntent.TriggerPort() {
		t.Fatal("workflow trigger ports must remain distinct")
	}
}

func TestDirectRevokePortSemanticPreparation(t *testing.T) {
	t.Parallel()
	spec := sampleSpec(t, model.AssignmentValidityStanding)
	in, fail := roleassign.NewDirectRevokeIntent(
		"p-dr", "ra-dr", "ra-dr", sampleScope(), spec, true, "3", "4", nil, state.OriginatingResult,
	)
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	intent, fail := roleassign.NewDirectRevokePort().Prepare(in)
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	if intent.TriggerPort() != "direct_revoke" || intent.Operation() != roleassign.OpDirectRevoke {
		t.Fatalf("port=%s op=%s", intent.TriggerPort(), intent.Operation())
	}
	if intent.TriggerPort() == "grant" || intent.TriggerPort() == "activation" || intent.TriggerPort() == "review" {
		t.Fatal("direct revoke must not be classified as a workflow trigger")
	}
	tx, permit := issuePermit(t, intent.ParticipantClaim())
	defer tx.Abort()
	outcome := intent.FinalizeAt(permit)
	change, ok := outcome.AsFinalized()
	if !ok {
		t.Fatalf("finalize: %s", outcome.KindName())
	}
	if change.EvidenceDescriptors()[0].EventType() != "roleassignment.revoked" {
		t.Fatal("direct revoke must emit roleassignment.revoked")
	}
}

func TestRoleassignNeverConstructsApprovalRequirement(t *testing.T) {
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
	for _, pkg := range pkgs {
		for name, f := range pkg.Files {
			if strings.HasSuffix(name, "_test.go") {
				continue
			}
			for _, imp := range f.Imports {
				path := strings.Trim(imp.Path.Value, `"`)
				if strings.Contains(path, "/govaccess/approvalreq") {
					t.Fatalf("%s imports approvalreq", name)
				}
			}
			ast.Inspect(f, func(n ast.Node) bool {
				id, ok := n.(*ast.Ident)
				if !ok {
					return true
				}
				if id.Name == "ApprovalRequirement" {
					t.Fatalf("%s references ApprovalRequirement", name)
				}
				return true
			})
		}
	}
}

func TestPropertyOnlyRoleassignConstructsIntentAndChange(t *testing.T) {
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

	sealedIntentFiles := map[string]int{}
	sealedChangeFiles := map[string]int{}
	newIntentCallers := map[string]int{}
	for _, pkg := range pkgs {
		for name, f := range pkg.Files {
			if strings.HasSuffix(name, "_test.go") {
				continue
			}
			base := filepath.Base(name)
			ast.Inspect(f, func(n ast.Node) bool {
				switch node := n.(type) {
				case *ast.CallExpr:
					if ident, ok := node.Fun.(*ast.Ident); ok && ident.Name == "newPreparedRoleAssignmentIntent" {
						newIntentCallers[base]++
					}
				case *ast.CompositeLit:
					typeName := compositeTypeName(node.Type)
					if !compositeSetsSealedTrue(node) {
						return true
					}
					switch typeName {
					case "PreparedRoleAssignmentIntent":
						sealedIntentFiles[base]++
					case "FinalizedRoleAssignmentChange":
						sealedChangeFiles[base]++
					}
				}
				return true
			})
		}
	}
	if len(sealedIntentFiles) != 1 || sealedIntentFiles["intent.go"] == 0 {
		t.Fatalf("sealed PreparedRoleAssignmentIntent must be constructed only in intent.go, got %#v", sealedIntentFiles)
	}
	if len(sealedChangeFiles) != 1 || sealedChangeFiles["finalization.go"] == 0 {
		t.Fatalf("sealed FinalizedRoleAssignmentChange must be constructed only in finalization.go, got %#v", sealedChangeFiles)
	}
	allowedCallers := map[string]bool{
		"grantport.go": true, "activationport.go": true, "reviewport.go": true, "revokeport.go": true,
	}
	for file := range newIntentCallers {
		if !allowedCallers[file] {
			t.Fatalf("newPreparedRoleAssignmentIntent called from unauthorized file %s", file)
		}
	}
	for file := range allowedCallers {
		if newIntentCallers[file] == 0 {
			t.Fatalf("expected %s to call newPreparedRoleAssignmentIntent", file)
		}
	}
}

func compositeTypeName(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.SelectorExpr:
		return t.Sel.Name
	default:
		return ""
	}
}

func compositeSetsSealedTrue(cl *ast.CompositeLit) bool {
	for _, elt := range cl.Elts {
		kv, ok := elt.(*ast.KeyValueExpr)
		if !ok {
			continue
		}
		key, ok := kv.Key.(*ast.Ident)
		if !ok || key.Name != "sealed" {
			continue
		}
		if lit, ok := kv.Value.(*ast.Ident); ok && lit.Name == "true" {
			return true
		}
	}
	return false
}
