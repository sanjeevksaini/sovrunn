package validate

import (
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/decision"
)

func assertRelationshipViolation(t *testing.T, p *apiproblem.Problem, wantCode apiproblem.ViolationCode, wantField string) {
	t.Helper()
	if p == nil {
		t.Fatalf("expected Problem with %s at %s, got nil", wantCode, wantField)
	}
	if p.Code != apiproblem.CodeValidationFailed {
		t.Fatalf("top-level code = %q, want %q", p.Code, apiproblem.CodeValidationFailed)
	}
	if p.Status != 422 {
		t.Fatalf("status = %d, want 422", p.Status)
	}
	if p.Type != "urn:sovrunn:problem:validation-failed" {
		t.Fatalf("type = %q, want urn:sovrunn:problem:validation-failed", p.Type)
	}
	if len(p.Violations) != 1 {
		t.Fatalf("violations len=%d, want 1: %#v", len(p.Violations), p.Violations)
	}
	v := p.Violations[0]
	if v.Code != wantCode {
		t.Fatalf("violations[0].code = %q, want %q", v.Code, wantCode)
	}
	if v.Field != wantField {
		t.Fatalf("violations[0].field = %q, want %q", v.Field, wantField)
	}
	if v.Message == "" {
		t.Fatal("violations[0].message must be non-empty (redactable)")
	}
}

func decisionRef(name, uid string) apimeta.TypedRef {
	return apimeta.TypedRef{
		APIVersion: decision.APIVersionDecisionRecord,
		Kind:       decision.KindDecisionRecord,
		Name:       name,
		UID:        uid,
	}
}

func rootNode(name, uid string) RelationshipNode {
	return RelationshipNode{Ref: decisionRef(name, uid)}
}

func linkNode(name, uid string, kind decision.RelationshipKind, target apimeta.TypedRef) RelationshipNode {
	return RelationshipNode{
		Ref: decisionRef(name, uid),
		Relationship: &decision.DecisionRelationship{
			Kind:      kind,
			TargetRef: target,
		},
	}
}

func TestValidateRelationship_NoopWhenAbsent(t *testing.T) {
	t.Parallel()

	p := ValidateRelationship(RelationshipInput{})
	if p != nil {
		t.Fatalf("absent relationship must pass, got %#v", p)
	}
}

func TestValidateRelationship_KindInvalid_Precondition(t *testing.T) {
	t.Parallel()

	root := rootNode("root", "uid-root")
	self := decisionRef("new", "uid-new")
	p := ValidateRelationship(RelationshipInput{
		Relationship: &decision.DecisionRelationship{
			Kind:      decision.RelationshipKind("AMENDS"),
			TargetRef: root.Ref,
		},
		Self:        self,
		Established: []RelationshipNode{root},
		MaxChain:    8,
	})
	assertRelationshipViolation(t, p, CodeRelationshipKindInvalid, ptrRelationshipKind)
}

func TestValidateRelationship_TargetMissingUID_Step1(t *testing.T) {
	t.Parallel()

	root := rootNode("root", "uid-root")
	p := ValidateRelationship(RelationshipInput{
		Relationship: &decision.DecisionRelationship{
			Kind: decision.RelationshipKindSupersedes,
			TargetRef: apimeta.TypedRef{
				APIVersion: decision.APIVersionDecisionRecord,
				Kind:       decision.KindDecisionRecord,
				Name:       "root",
			},
		},
		Self:        decisionRef("new", "uid-new"),
		Established: []RelationshipNode{root},
		MaxChain:    8,
	})
	assertRelationshipViolation(t, p, CodeRelationshipTargetInvalid, ptrRelationshipTarget)
}

func TestValidateRelationship_TargetUnresolvable_Step1(t *testing.T) {
	t.Parallel()

	root := rootNode("root", "uid-root")
	p := ValidateRelationship(RelationshipInput{
		Relationship: &decision.DecisionRelationship{
			Kind:      decision.RelationshipKindSupersedes,
			TargetRef: decisionRef("missing", "uid-missing"),
		},
		Self:        decisionRef("new", "uid-new"),
		Established: []RelationshipNode{root},
		MaxChain:    8,
	})
	assertRelationshipViolation(t, p, CodeRelationshipTargetInvalid, ptrRelationshipTarget)
}

func TestValidateRelationship_TargetSelfReferential_Step1(t *testing.T) {
	t.Parallel()

	self := decisionRef("self", "uid-self")
	p := ValidateRelationship(RelationshipInput{
		Relationship: &decision.DecisionRelationship{
			Kind:      decision.RelationshipKindCorrects,
			TargetRef: self,
		},
		Self:        self,
		Established: []RelationshipNode{rootNode("self", "uid-self")},
		MaxChain:    8,
	})
	assertRelationshipViolation(t, p, CodeRelationshipTargetInvalid, ptrRelationshipTarget)
}

func TestValidateRelationship_TargetWrongKind_Step1(t *testing.T) {
	t.Parallel()

	target := apimeta.TypedRef{
		APIVersion: "core.sovrunn.io/v1alpha1",
		Kind:       "Project",
		Name:       "payments",
		UID:        "uid-project",
	}
	p := ValidateRelationship(RelationshipInput{
		Relationship: &decision.DecisionRelationship{
			Kind:      decision.RelationshipKindSupersedes,
			TargetRef: target,
		},
		Self:        decisionRef("new", "uid-new"),
		Established: []RelationshipNode{{Ref: target}},
		MaxChain:    8,
	})
	assertRelationshipViolation(t, p, CodeRelationshipTargetInvalid, ptrRelationshipTarget)
}

func TestValidateRelationship_Cycle_Step2(t *testing.T) {
	t.Parallel()

	// a → b → (candidate self → a) forms a cycle.
	a := rootNode("a", "uid-a")
	b := linkNode("b", "uid-b", decision.RelationshipKindSupersedes, a.Ref)
	self := decisionRef("self", "uid-self")
	// Make `a` point at self so walking from target `a` reaches self.
	// Cycle via: self → a → self.
	aLinked := linkNode("a", "uid-a", decision.RelationshipKindSupersedes, self)
	p := ValidateRelationship(RelationshipInput{
		Relationship: &decision.DecisionRelationship{
			Kind:      decision.RelationshipKindCorrects,
			TargetRef: aLinked.Ref,
		},
		Self:        self,
		Established: []RelationshipNode{aLinked, b},
		MaxChain:    8,
	})
	assertRelationshipViolation(t, p, CodeRelationshipCycle, ptrRelationshipTarget)
}

func TestValidateRelationship_ChainLimitExceeded_Step3(t *testing.T) {
	t.Parallel()

	r0 := rootNode("r0", "uid-0")
	r1 := linkNode("r1", "uid-1", decision.RelationshipKindSupersedes, r0.Ref)
	r2 := linkNode("r2", "uid-2", decision.RelationshipKindSupersedes, r1.Ref)
	// Candidate → r2 has depth 3; MaxChain=2 → exceed.
	p := ValidateRelationship(RelationshipInput{
		Relationship: &decision.DecisionRelationship{
			Kind:      decision.RelationshipKindSupersedes,
			TargetRef: r2.Ref,
		},
		Self:        decisionRef("r3", "uid-3"),
		Established: []RelationshipNode{r0, r1, r2},
		MaxChain:    2,
	})
	assertRelationshipViolation(t, p, CodeRelationshipChainLimitExceeded, ptrRelationship)
}

func TestValidateRelationship_SameTargetDifferentKinds_Step4(t *testing.T) {
	t.Parallel()

	root := rootNode("root", "uid-root")
	corrects := linkNode("c1", "uid-c1", decision.RelationshipKindCorrects, root.Ref)
	p := ValidateRelationship(RelationshipInput{
		Relationship: &decision.DecisionRelationship{
			Kind:      decision.RelationshipKindSupersedes,
			TargetRef: root.Ref,
		},
		Self:        decisionRef("s1", "uid-s1"),
		Established: []RelationshipNode{root, corrects},
		MaxChain:    8,
	})
	assertRelationshipViolation(t, p, CodeRelationshipConflict, ptrRelationship)
}

func TestValidateRelationship_SameTargetSameKind_CreatesSecondCurrent_Step6(t *testing.T) {
	t.Parallel()

	root := rootNode("root", "uid-root")
	first := linkNode("s1", "uid-s1", decision.RelationshipKindSupersedes, root.Ref)
	p := ValidateRelationship(RelationshipInput{
		Relationship: &decision.DecisionRelationship{
			Kind:      decision.RelationshipKindSupersedes,
			TargetRef: root.Ref,
		},
		Self:        decisionRef("s2", "uid-s2"),
		Established: []RelationshipNode{root, first},
		MaxChain:    8,
	})
	assertRelationshipViolation(t, p, CodeRelationshipConflict, ptrRelationship)
}

func TestValidateRelationship_SameTargetSameKind_HistoricalFork_Valid(t *testing.T) {
	t.Parallel()

	// Structural fork: s1 SUPERSEDES root, then s1b SUPERSEDES s1 (s1 no longer
	// current). Candidate s2 SUPERSEDES root → only one current successor of root.
	root := rootNode("root", "uid-root")
	s1 := linkNode("s1", "uid-s1", decision.RelationshipKindSupersedes, root.Ref)
	s1b := linkNode("s1b", "uid-s1b", decision.RelationshipKindSupersedes, s1.Ref)
	p := ValidateRelationship(RelationshipInput{
		Relationship: &decision.DecisionRelationship{
			Kind:      decision.RelationshipKindSupersedes,
			TargetRef: root.Ref,
		},
		Self:        decisionRef("s2", "uid-s2"),
		Established: []RelationshipNode{root, s1, s1b},
		MaxChain:    8,
	})
	if p != nil {
		t.Fatalf("historical fork with single current head must pass, got %#v", p)
	}
}

func TestValidateRelationship_ChainConvergence_Valid(t *testing.T) {
	t.Parallel()

	// Shared ancestor, acyclic: r0 ← r1 ← r2 (candidate extends r2).
	r0 := rootNode("r0", "uid-0")
	r1 := linkNode("r1", "uid-1", decision.RelationshipKindSupersedes, r0.Ref)
	r2 := linkNode("r2", "uid-2", decision.RelationshipKindCorrects, r1.Ref)
	p := ValidateRelationship(RelationshipInput{
		Relationship: &decision.DecisionRelationship{
			Kind:      decision.RelationshipKindSupersedes,
			TargetRef: r2.Ref,
		},
		Self:        decisionRef("r3", "uid-3"),
		Established: []RelationshipNode{r0, r1, r2},
		MaxChain:    8,
	})
	if p != nil {
		t.Fatalf("acyclic shared-ancestor chain must pass, got %#v", p)
	}
}

func TestValidateRelationship_MultipleCurrentSuccessors_Step6(t *testing.T) {
	t.Parallel()

	root := rootNode("root", "uid-root")
	c1 := linkNode("c1", "uid-c1", decision.RelationshipKindCorrects, root.Ref)
	p := ValidateRelationship(RelationshipInput{
		Relationship: &decision.DecisionRelationship{
			Kind:      decision.RelationshipKindCorrects,
			TargetRef: root.Ref,
		},
		Self:        decisionRef("c2", "uid-c2"),
		Established: []RelationshipNode{root, c1},
		MaxChain:    8,
	})
	assertRelationshipViolation(t, p, CodeRelationshipConflict, ptrRelationship)
}

func TestValidateRelationship_RevokesLiveTarget_Valid(t *testing.T) {
	t.Parallel()

	root := rootNode("root", "uid-root")
	live := linkNode("live", "uid-live", decision.RelationshipKindSupersedes, root.Ref)
	p := ValidateRelationship(RelationshipInput{
		Relationship: &decision.DecisionRelationship{
			Kind:      decision.RelationshipKindRevokes,
			TargetRef: live.Ref,
		},
		Self:        decisionRef("rev", "uid-rev"),
		Established: []RelationshipNode{root, live},
		MaxChain:    8,
	})
	if p != nil {
		t.Fatalf("REVOKES of a live target must pass, got %#v", p)
	}

	graph := []RelationshipNode{root, live, linkNode("rev", "uid-rev", decision.RelationshipKindRevokes, live.Ref)}
	if got := ProjectEffectiveState(live.Ref, graph); got != decision.EffectiveStateRevoked {
		t.Fatalf("projected state of revoked target = %q, want REVOKED", got)
	}
}

func TestValidateRelationship_CorrectsAfterRevocation_Step5(t *testing.T) {
	t.Parallel()

	root := rootNode("root", "uid-root")
	rev := linkNode("rev", "uid-rev", decision.RelationshipKindRevokes, root.Ref)
	p := ValidateRelationship(RelationshipInput{
		Relationship: &decision.DecisionRelationship{
			Kind:      decision.RelationshipKindCorrects,
			TargetRef: root.Ref,
		},
		Self:        decisionRef("c1", "uid-c1"),
		Established: []RelationshipNode{root, rev},
		MaxChain:    8,
	})
	assertRelationshipViolation(t, p, CodeRelationshipConflict, ptrRelationshipTarget)
}

func TestValidateRelationship_SupersedesAfterRevocation_Step5(t *testing.T) {
	t.Parallel()

	root := rootNode("root", "uid-root")
	rev := linkNode("rev", "uid-rev", decision.RelationshipKindRevokes, root.Ref)
	p := ValidateRelationship(RelationshipInput{
		Relationship: &decision.DecisionRelationship{
			Kind:      decision.RelationshipKindSupersedes,
			TargetRef: root.Ref,
		},
		Self:        decisionRef("s1", "uid-s1"),
		Established: []RelationshipNode{root, rev},
		MaxChain:    8,
	})
	assertRelationshipViolation(t, p, CodeRelationshipConflict, ptrRelationshipTarget)
}

func TestValidateRelationship_ReRevocation_Step5(t *testing.T) {
	t.Parallel()

	root := rootNode("root", "uid-root")
	rev := linkNode("rev", "uid-rev", decision.RelationshipKindRevokes, root.Ref)
	p := ValidateRelationship(RelationshipInput{
		Relationship: &decision.DecisionRelationship{
			Kind:      decision.RelationshipKindRevokes,
			TargetRef: root.Ref,
		},
		Self:        decisionRef("rev2", "uid-rev2"),
		Established: []RelationshipNode{root, rev},
		MaxChain:    8,
	})
	assertRelationshipViolation(t, p, CodeRelationshipConflict, ptrRelationshipTarget)
}

func TestValidateRelationship_CorrectsRevokedChainSuccessor_Step5(t *testing.T) {
	t.Parallel()

	// root revoked; child SUPERSEDES root was established before revocation is
	// not re-creatable, but targeting the child after root revocation conflicts
	// as revoked-chain.
	root := rootNode("root", "uid-root")
	child := linkNode("child", "uid-child", decision.RelationshipKindSupersedes, root.Ref)
	rev := linkNode("rev", "uid-rev", decision.RelationshipKindRevokes, root.Ref)
	p := ValidateRelationship(RelationshipInput{
		Relationship: &decision.DecisionRelationship{
			Kind:      decision.RelationshipKindCorrects,
			TargetRef: child.Ref,
		},
		Self:        decisionRef("c1", "uid-c1"),
		Established: []RelationshipNode{root, child, rev},
		MaxChain:    8,
	})
	assertRelationshipViolation(t, p, CodeRelationshipConflict, ptrRelationshipTarget)
}

func TestValidateRelationship_HappyPathSupersedesRoot(t *testing.T) {
	t.Parallel()

	root := rootNode("root", "uid-root")
	self := decisionRef("s1", "uid-s1")
	rel := &decision.DecisionRelationship{
		Kind:      decision.RelationshipKindSupersedes,
		TargetRef: root.Ref,
	}
	p := ValidateRelationship(RelationshipInput{
		Relationship: rel,
		Self:         self,
		Established:  []RelationshipNode{root},
		MaxChain:     8,
	})
	if p != nil {
		t.Fatalf("supersession of live root must pass, got %#v", p)
	}

	graph := []RelationshipNode{root, {Ref: self, Relationship: rel}}
	if got := ProjectEffectiveState(root.Ref, graph); got != decision.EffectiveStateSuperseded {
		t.Fatalf("projected predecessor state = %q, want SUPERSEDED", got)
	}
	if got := ProjectEffectiveState(self, graph); got != decision.EffectiveStateEffective {
		t.Fatalf("projected successor state = %q, want EFFECTIVE", got)
	}
}

func TestValidateRelationship_CurrentVsHistoricalProjection_Step7(t *testing.T) {
	t.Parallel()

	root := rootNode("root", "uid-root")
	s1 := linkNode("s1", "uid-s1", decision.RelationshipKindSupersedes, root.Ref)
	s2 := linkNode("s2", "uid-s2", decision.RelationshipKindSupersedes, s1.Ref)
	graph := []RelationshipNode{root, s1, s2}

	if got := ProjectEffectiveState(root.Ref, graph); got != decision.EffectiveStateSuperseded {
		t.Fatalf("root = %q, want SUPERSEDED", got)
	}
	if got := ProjectEffectiveState(s1.Ref, graph); got != decision.EffectiveStateSuperseded {
		t.Fatalf("s1 = %q, want SUPERSEDED", got)
	}
	if got := ProjectEffectiveState(s2.Ref, graph); got != decision.EffectiveStateEffective {
		t.Fatalf("s2 = %q, want EFFECTIVE", got)
	}
}

func TestValidateRelationship_MaxChainZeroMeansUnlimited(t *testing.T) {
	t.Parallel()

	r0 := rootNode("r0", "uid-0")
	r1 := linkNode("r1", "uid-1", decision.RelationshipKindSupersedes, r0.Ref)
	r2 := linkNode("r2", "uid-2", decision.RelationshipKindSupersedes, r1.Ref)
	p := ValidateRelationship(RelationshipInput{
		Relationship: &decision.DecisionRelationship{
			Kind:      decision.RelationshipKindSupersedes,
			TargetRef: r2.Ref,
		},
		Self:        decisionRef("r3", "uid-3"),
		Established: []RelationshipNode{r0, r1, r2},
		MaxChain:    0,
	})
	if p != nil {
		t.Fatalf("MaxChain<=0 must not enforce a ceiling, got %#v", p)
	}
}

func TestValidateRelationship_ChainDepthAtLimit_Valid(t *testing.T) {
	t.Parallel()

	r0 := rootNode("r0", "uid-0")
	r1 := linkNode("r1", "uid-1", decision.RelationshipKindSupersedes, r0.Ref)
	// Candidate → r1 has depth 2; MaxChain=2 → accept.
	p := ValidateRelationship(RelationshipInput{
		Relationship: &decision.DecisionRelationship{
			Kind:      decision.RelationshipKindSupersedes,
			TargetRef: r1.Ref,
		},
		Self:        decisionRef("r2", "uid-2"),
		Established: []RelationshipNode{r0, r1},
		MaxChain:    2,
	})
	if p != nil {
		t.Fatalf("chain depth at declared limit must pass, got %#v", p)
	}
}
