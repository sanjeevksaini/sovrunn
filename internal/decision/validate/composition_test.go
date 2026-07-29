package validate

import (
	"encoding/json"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/apivalid"
	"github.com/sanjeevksaini/sovrunn/internal/decision"
	"github.com/sanjeevksaini/sovrunn/internal/decision/bundle"
	"github.com/sanjeevksaini/sovrunn/internal/decision/graph"
)

func compositionProfile() decision.DecisionProfile {
	p := validProfile()
	p.Spec.AcceptedCompositionStrategies = []string{"all-of-v1", "any-of@2.0.0"}
	p.Spec.Limits = decision.ProfileLimits{
		MaxNodes:        8,
		MaxEdges:        16,
		MaxDepth:        4,
		MaxWidth:        4,
		MaxFanOut:       4,
		MaxPayloadBytes: 1024,
	}
	return p
}

func compositionBounds() graph.Bounds {
	return graph.Bounds{
		MaxNodes:        8,
		MaxEdges:        16,
		MaxDepth:        4,
		MaxWidth:        4,
		MaxFanOut:       4,
		MaxPayloadBytes: 1024,
	}
}

func linearCompositionGraph() graph.Definition {
	return graph.Definition{
		Ref:     "g-linear",
		Version: "1.0.0",
		Nodes: []graph.Node{
			{ID: "in1", Kind: graph.NodeKindInput, PayloadBytes: 1},
			{ID: "ev1", Kind: graph.NodeKindEvaluation, PayloadBytes: 2},
			{ID: "d1", Kind: graph.NodeKindDecision, PayloadBytes: 4},
		},
		Edges: []graph.Edge{
			{ID: "e-in-ev", From: "in1", To: "ev1"},
			{ID: "e-ev-d", From: "ev1", To: "d1"},
		},
		Bounds: compositionBounds(),
	}
}

func singleDecisionGraph() graph.Definition {
	return graph.Definition{
		Ref:     "g-single",
		Version: "1.0.0",
		Nodes: []graph.Node{
			{ID: "d1", Kind: graph.NodeKindDecision, PayloadBytes: 1},
		},
		Edges:  []graph.Edge{},
		Bounds: compositionBounds(),
	}
}

func assertCompositionViolation(t *testing.T, p *apiproblem.Problem, wantCode apiproblem.ViolationCode, wantField string) {
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

func mustLoadCompositionBundle(t *testing.T, graphs ...graph.Definition) bundle.Bundle {
	t.Helper()
	prof := compositionProfile()
	doc := bundle.Bundle{
		APIVersion: bundle.APIVersionDecisionProfileBundle,
		Kind:       bundle.KindDecisionProfileBundle,
		Spec: bundle.BundleSpec{
			Profiles: []decision.DecisionProfile{prof},
			Graphs:   graphs,
			Strategies: []bundle.VersionedEntry{
				{ID: "all-of-v1", Version: "1.0.0"},
				{ID: "any-of", Version: "2.0.0"},
			},
		},
	}
	raw, err := json.Marshal(doc)
	if err != nil {
		t.Fatalf("marshal bundle: %v", err)
	}
	loaded, prob := bundle.Load(raw, apivalid.ModeReadRepresentation)
	if prob != nil {
		t.Fatalf("Load = %#v", prob)
	}
	return loaded
}

func TestValidateComposition_HappyPathStrategyOnly(t *testing.T) {
	t.Parallel()

	p := ValidateComposition(CompositionInput{
		Composition: decision.CompositionRef{
			StrategyName:    "all-of-v1",
			StrategyVersion: "1.0.0",
		},
		Profile: compositionProfile(),
	})
	if p != nil {
		t.Fatalf("strategy-only composition must pass, got %#v", p)
	}
}

func TestValidateComposition_HappyPathWithResolvedGraph(t *testing.T) {
	t.Parallel()

	def := linearCompositionGraph()
	p := ValidateComposition(CompositionInput{
		Composition: decision.CompositionRef{
			StrategyName:    "all-of-v1",
			StrategyVersion: "1.0.0",
			GraphRef:        def.Ref,
			GraphVersion:    def.Version,
			InputRefs:       []string{"in1"},
		},
		Profile:       compositionProfile(),
		ResolvedGraph: &def,
	})
	if p != nil {
		t.Fatalf("valid composition+graph must pass, got %#v", p)
	}
}

func TestValidateComposition_HappyPathBundleResolve(t *testing.T) {
	t.Parallel()

	def := linearCompositionGraph()
	loaded := mustLoadCompositionBundle(t, def)
	p := ValidateComposition(CompositionInput{
		Composition: decision.CompositionRef{
			StrategyName:    "all-of-v1",
			StrategyVersion: "1.0.0",
			GraphRef:        def.Ref,
			GraphVersion:    def.Version,
		},
		Profile: compositionProfile(),
		Bundle:  &loaded,
	})
	if p != nil {
		t.Fatalf("bundle-resolved composition must pass, got %#v", p)
	}
}

func TestValidateComposition_NoopWhenAbsent(t *testing.T) {
	t.Parallel()

	p := ValidateComposition(CompositionInput{Profile: compositionProfile()})
	if p != nil {
		t.Fatalf("absent composition must be a no-op, got %#v", p)
	}
}

func TestValidateComposition_StrategyUnsupported(t *testing.T) {
	t.Parallel()

	p := ValidateComposition(CompositionInput{
		Composition: decision.CompositionRef{
			StrategyName:    "unknown-strategy",
			StrategyVersion: "1.0.0",
		},
		Profile: compositionProfile(),
	})
	assertCompositionViolation(t, p, CodeCompositionStrategyUnsupported, ptrCompositionStrategyName)
}

func TestValidateComposition_StrategyVersionRequired(t *testing.T) {
	t.Parallel()

	p := ValidateComposition(CompositionInput{
		Composition: decision.CompositionRef{
			StrategyName: "all-of-v1",
		},
		Profile: compositionProfile(),
	})
	assertCompositionViolation(t, p, CodeCompositionStrategyUnsupported, ptrCompositionStrategyVersion)
}

func TestValidateComposition_StrategyCompositeAcceptedForm(t *testing.T) {
	t.Parallel()

	p := ValidateComposition(CompositionInput{
		Composition: decision.CompositionRef{
			StrategyName:    "any-of",
			StrategyVersion: "2.0.0",
		},
		Profile: compositionProfile(),
	})
	if p != nil {
		t.Fatalf("name@version accepted form must pass, got %#v", p)
	}
}

func TestValidateComposition_GraphRefNotFound(t *testing.T) {
	t.Parallel()

	loaded := mustLoadCompositionBundle(t, singleDecisionGraph())
	p := ValidateComposition(CompositionInput{
		Composition: decision.CompositionRef{
			StrategyName:    "all-of-v1",
			StrategyVersion: "1.0.0",
			GraphRef:        "missing",
			GraphVersion:    "1.0.0",
		},
		Profile: compositionProfile(),
		Bundle:  &loaded,
	})
	assertCompositionViolation(t, p, CodeCompositionGraphInvalid, ptrCompositionGraphRef)
}

func TestValidateGraph_MissingNodeID_GraphLocalPointer(t *testing.T) {
	t.Parallel()

	def := linearCompositionGraph()
	def.Nodes[0].ID = ""
	p := ValidateGraph(def, compositionBounds(), "")
	assertCompositionViolation(t, p, CodeCompositionGraphInvalid, "/nodes/0/id")
}

func TestValidateGraph_MissingNodeID_BundleEmbeddedPointer(t *testing.T) {
	t.Parallel()

	def := linearCompositionGraph()
	def.Nodes[0].ID = ""
	p := ValidateGraph(def, compositionBounds(), "/spec/graphs/0")
	assertCompositionViolation(t, p, CodeCompositionGraphInvalid, "/spec/graphs/0/nodes/0/id")
}

func TestValidateGraph_DuplicateNodeID_BundleEmbeddedPointer(t *testing.T) {
	t.Parallel()

	def := linearCompositionGraph()
	def.Nodes = append(def.Nodes, graph.Node{ID: "in1", Kind: graph.NodeKindInput})
	p := ValidateGraph(def, compositionBounds(), "/spec/graphs/2")
	assertCompositionViolation(t, p, CodeCompositionGraphInvalid, "/spec/graphs/2/nodes/3/id")
}

func TestValidateGraph_InvalidNodeKind_GraphLocalPointer(t *testing.T) {
	t.Parallel()

	def := linearCompositionGraph()
	def.Nodes[1].Kind = "workflow"
	p := ValidateGraph(def, compositionBounds(), "")
	assertCompositionViolation(t, p, CodeCompositionGraphInvalid, "/nodes/1/kind")
}

func TestValidateGraph_Cycle(t *testing.T) {
	t.Parallel()

	def := graph.Definition{
		Ref:     "g-cycle",
		Version: "1.0.0",
		Nodes: []graph.Node{
			{ID: "a", Kind: graph.NodeKindEvaluation},
			{ID: "b", Kind: graph.NodeKindEvaluation},
			{ID: "d1", Kind: graph.NodeKindDecision},
		},
		Edges: []graph.Edge{
			{ID: "e-ab", From: "a", To: "b"},
			{ID: "e-ba", From: "b", To: "a"},
			{ID: "e-bd", From: "b", To: "d1"},
		},
		Bounds: compositionBounds(),
	}
	p := ValidateGraph(def, compositionBounds(), "")
	assertCompositionViolation(t, p, CodeCompositionCycle, "/edges")

	p = ValidateGraph(def, compositionBounds(), "/spec/graphs/1")
	assertCompositionViolation(t, p, CodeCompositionCycle, "/spec/graphs/1/edges")
}

func TestValidateGraph_BoundExceeded_MaxDepth(t *testing.T) {
	t.Parallel()

	def := linearCompositionGraph()
	tight := compositionBounds()
	tight.MaxDepth = 1 // linear path depth is 2
	p := ValidateGraph(def, tight, "")
	assertCompositionViolation(t, p, CodeCompositionLimitExceeded, "/bounds/maxDepth")

	p = ValidateGraph(def, tight, "/spec/graphs/0")
	assertCompositionViolation(t, p, CodeCompositionLimitExceeded, "/spec/graphs/0/bounds/maxDepth")
}

func TestValidateGraph_BoundExceeded_MaxNodes(t *testing.T) {
	t.Parallel()

	def := linearCompositionGraph()
	tight := compositionBounds()
	tight.MaxNodes = 2
	p := ValidateGraph(def, tight, "/spec/graphs/3")
	assertCompositionViolation(t, p, CodeCompositionLimitExceeded, "/spec/graphs/3/bounds/maxNodes")
}

func TestValidateGraph_MalformedBounds(t *testing.T) {
	t.Parallel()

	def := singleDecisionGraph()
	bad := compositionBounds()
	bad.MaxNodes = 0
	p := ValidateGraph(def, bad, "")
	assertCompositionViolation(t, p, CodeCompositionGraphInvalid, "/bounds/maxNodes")
}

func TestValidateComposition_InputConflictNonInputNode(t *testing.T) {
	t.Parallel()

	def := linearCompositionGraph()
	// nodes[1] is evaluation "ev1"; referencing it as an input is a conflict.
	p := ValidateComposition(CompositionInput{
		Composition: decision.CompositionRef{
			StrategyName:    "all-of-v1",
			StrategyVersion: "1.0.0",
			InputRefs:       []string{"ev1"},
		},
		Profile:       compositionProfile(),
		ResolvedGraph: &def,
	})
	assertCompositionViolation(t, p, CodeCompositionInputConflict, "/nodes/1")
}

func TestValidateComposition_InputConflictNonInputNode_BundleEmbeddedPointer(t *testing.T) {
	t.Parallel()

	def := linearCompositionGraph()
	p := ValidateComposition(CompositionInput{
		Composition: decision.CompositionRef{
			StrategyName:    "all-of-v1",
			StrategyVersion: "1.0.0",
			InputRefs:       []string{"ev1"},
		},
		Profile:       compositionProfile(),
		ResolvedGraph: &def,
		PointerBase:   "/spec/graphs/0",
	})
	assertCompositionViolation(t, p, CodeCompositionInputConflict, "/spec/graphs/0/nodes/1")
}

func TestValidateComposition_InputConflictDuplicate(t *testing.T) {
	t.Parallel()

	def := linearCompositionGraph()
	p := ValidateComposition(CompositionInput{
		Composition: decision.CompositionRef{
			StrategyName:    "all-of-v1",
			StrategyVersion: "1.0.0",
			InputRefs:       []string{"in1", "in1"},
		},
		Profile:       compositionProfile(),
		ResolvedGraph: &def,
	})
	assertCompositionViolation(t, p, CodeCompositionInputConflict, "/nodes/0")
}

func TestValidateComposition_EffectiveBoundsUseProfileCeiling(t *testing.T) {
	t.Parallel()

	def := linearCompositionGraph()
	prof := compositionProfile()
	prof.Spec.Limits.MaxDepth = 1 // tighter than definition bounds
	p := ValidateComposition(CompositionInput{
		Composition: decision.CompositionRef{
			StrategyName:    "all-of-v1",
			StrategyVersion: "1.0.0",
		},
		Profile:       prof,
		ResolvedGraph: &def,
	})
	assertCompositionViolation(t, p, CodeCompositionLimitExceeded, "/bounds/maxDepth")
}

func TestEffectiveGraphBounds_MostRestrictiveWins(t *testing.T) {
	t.Parallel()

	profile := decision.ProfileLimits{
		MaxNodes:        8,
		MaxEdges:        16,
		MaxDepth:        4,
		MaxWidth:        4,
		MaxFanOut:       2,
		MaxPayloadBytes: 512,
	}
	def := graph.Bounds{
		MaxNodes:        16,
		MaxEdges:        8,
		MaxDepth:        8,
		MaxWidth:        2,
		MaxFanOut:       4,
		MaxPayloadBytes: 1024,
	}
	got := EffectiveGraphBounds(profile, def)
	if got.MaxNodes != 8 {
		t.Fatalf("MaxNodes = %d, want 8", got.MaxNodes)
	}
	if got.MaxEdges != 8 {
		t.Fatalf("MaxEdges = %d, want 8", got.MaxEdges)
	}
	if got.MaxWidth != 2 {
		t.Fatalf("MaxWidth = %d, want 2", got.MaxWidth)
	}
	if got.MaxFanOut != 2 {
		t.Fatalf("MaxFanOut = %d, want 2", got.MaxFanOut)
	}
	if got.MaxPayloadBytes != 512 {
		t.Fatalf("MaxPayloadBytes = %d, want 512", got.MaxPayloadBytes)
	}
}

func TestJoinPointerBase(t *testing.T) {
	t.Parallel()

	if got := joinPointerBase("", "/nodes/0/id"); got != "/nodes/0/id" {
		t.Fatalf("standalone = %q", got)
	}
	if got := joinPointerBase("/spec/graphs/0", "/nodes/0/id"); got != "/spec/graphs/0/nodes/0/id" {
		t.Fatalf("embedded = %q", got)
	}
	if got := joinPointerBase("/spec/graphs/0/", "/edges"); got != "/spec/graphs/0/edges" {
		t.Fatalf("trim base slash = %q", got)
	}
}

func TestValidateComposition_RequiresBundleWhenGraphRefSet(t *testing.T) {
	t.Parallel()

	p := ValidateComposition(CompositionInput{
		Composition: decision.CompositionRef{
			StrategyName:    "all-of-v1",
			StrategyVersion: "1.0.0",
			GraphRef:        "g-linear",
			GraphVersion:    "1.0.0",
		},
		Profile: compositionProfile(),
	})
	assertCompositionViolation(t, p, CodeCompositionGraphInvalid, ptrCompositionGraphRef)
}
