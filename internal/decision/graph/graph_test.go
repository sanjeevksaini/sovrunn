package graph

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
)

func testBounds() Bounds {
	return Bounds{
		MaxNodes:        16,
		MaxEdges:        32,
		MaxDepth:        8,
		MaxWidth:        8,
		MaxFanOut:       8,
		MaxPayloadBytes: 1024,
	}
}

func singleDecisionDef() Definition {
	return Definition{
		Ref:     "g-single",
		Version: "1.0.0",
		Nodes: []Node{
			{ID: "d1", Kind: NodeKindDecision, PayloadBytes: 10},
		},
		Edges:  []Edge{},
		Bounds: testBounds(),
	}
}

func linearDef() Definition {
	return Definition{
		Ref:     "g-linear",
		Version: "1.0.0",
		Nodes: []Node{
			{ID: "in1", Kind: NodeKindInput, PayloadBytes: 1},
			{ID: "ev1", Kind: NodeKindEvaluation, PayloadBytes: 2},
			{ID: "d1", Kind: NodeKindDecision, PayloadBytes: 4},
		},
		Edges: []Edge{
			{ID: "e-in-ev", From: "in1", To: "ev1"},
			{ID: "e-ev-d", From: "ev1", To: "d1"},
		},
		Bounds: testBounds(),
	}
}

func TestBuildSingleNodeEmptyEdges(t *testing.T) {
	t.Parallel()

	g, prob := Build(singleDecisionDef(), testBounds())
	if prob != nil {
		t.Fatalf("Build = %#v, want nil", prob)
	}
	m := g.Accounting()
	if m.Nodes != 1 || m.Edges != 0 || m.Depth != 0 || m.Width != 1 || m.FanOut != 0 || m.CriticalPath != 1 {
		t.Fatalf("Accounting = %+v, want single-node empty-edge metrics", m)
	}
	if m.PayloadBytes != 10 {
		t.Fatalf("PayloadBytes = %d, want 10", m.PayloadBytes)
	}
}

func TestBuildRejectsEmptyEdgesNonDecision(t *testing.T) {
	t.Parallel()

	def := Definition{
		Ref:     "g-bad",
		Version: "1.0.0",
		Nodes:   []Node{{ID: "in1", Kind: NodeKindInput}},
		Edges:   []Edge{},
		Bounds:  testBounds(),
	}
	_, prob := Build(def, testBounds())
	assertViolation(t, prob, violationGraphInvalid, "/nodes")
}

func TestBuildUniqueNodeAndEdgeIDs(t *testing.T) {
	t.Parallel()

	dupNode := linearDef()
	dupNode.Nodes = append(dupNode.Nodes, Node{ID: "in1", Kind: NodeKindInput})
	_, prob := Build(dupNode, testBounds())
	assertViolation(t, prob, violationGraphInvalid, "/nodes/3/id")

	dupEdge := linearDef()
	dupEdge.Edges = append(dupEdge.Edges, Edge{ID: "e-in-ev", From: "in1", To: "d1"})
	_, prob = Build(dupEdge, testBounds())
	assertViolation(t, prob, violationGraphInvalid, "/edges/2/id")
}

func TestBuildEdgeIDDistinctFromEndpoints(t *testing.T) {
	t.Parallel()

	def := linearDef()
	def.Edges[0].ID = "in1" // collides with From
	_, prob := Build(def, testBounds())
	assertViolation(t, prob, violationGraphInvalid, "/edges/0/id")

	def = linearDef()
	def.Edges[1].ID = "d1" // collides with To
	_, prob = Build(def, testBounds())
	assertViolation(t, prob, violationGraphInvalid, "/edges/1/id")
}

func TestBuildRejectsDuplicateFromToPair(t *testing.T) {
	t.Parallel()

	def := linearDef()
	def.Edges = append(def.Edges, Edge{ID: "e-dup-pair", From: "in1", To: "ev1"})
	_, prob := Build(def, testBounds())
	assertViolation(t, prob, violationGraphInvalid, "/edges/2")
}

func TestBuildParallelGroupIndependence(t *testing.T) {
	t.Parallel()

	independent := Definition{
		Ref:     "g-par",
		Version: "1.0.0",
		Nodes: []Node{
			{ID: "a", Kind: NodeKindInput},
			{ID: "b", Kind: NodeKindInput},
			{ID: "c", Kind: NodeKindComposition},
			{ID: "d", Kind: NodeKindDecision},
		},
		Edges: []Edge{
			{ID: "e-a-c", From: "a", To: "c"},
			{ID: "e-b-c", From: "b", To: "c"},
			{ID: "e-c-d", From: "c", To: "d"},
		},
		ParallelGroups: []ParallelGroup{
			{ID: "pg1", NodeIDs: []string{"a", "b"}},
		},
		Bounds: testBounds(),
	}
	g, prob := Build(independent, testBounds())
	if prob != nil {
		t.Fatalf("independent parallel group Build = %#v", prob)
	}
	if len(g.Definition().ParallelGroups) != 1 {
		t.Fatalf("ParallelGroups len = %d, want 1", len(g.Definition().ParallelGroups))
	}

	dependent := independent
	dependent.Edges = append(dependent.Edges, Edge{ID: "e-a-b", From: "a", To: "b"})
	_, prob = Build(dependent, testBounds())
	assertViolation(t, prob, violationGraphInvalid, "/parallelGroups/0/nodeIds/1")
}

func TestBuildDeterministicAccounting(t *testing.T) {
	t.Parallel()

	def := Definition{
		Ref:     "g-acct",
		Version: "1.0.0",
		Nodes: []Node{
			{ID: "i1", Kind: NodeKindInput, PayloadBytes: 3},
			{ID: "i2", Kind: NodeKindInput, PayloadBytes: 5},
			{ID: "e1", Kind: NodeKindEvaluation, PayloadBytes: 7},
			{ID: "c1", Kind: NodeKindComposition, PayloadBytes: 11},
			{ID: "d1", Kind: NodeKindDecision, PayloadBytes: 13},
		},
		Edges: []Edge{
			{ID: "e-i1-e1", From: "i1", To: "e1"},
			{ID: "e-i2-e1", From: "i2", To: "e1"},
			{ID: "e-e1-c1", From: "e1", To: "c1"},
			{ID: "e-c1-d1", From: "c1", To: "d1"},
		},
		ParallelGroups: []ParallelGroup{
			{ID: "pg", NodeIDs: []string{"i1", "i2"}},
		},
		Bounds: testBounds(),
	}

	g1, prob := Build(def, testBounds())
	if prob != nil {
		t.Fatalf("Build = %#v", prob)
	}
	g2, prob := Build(def, testBounds())
	if prob != nil {
		t.Fatalf("second Build = %#v", prob)
	}

	m1 := g1.Accounting()
	m2 := g2.Accounting()
	if !reflect.DeepEqual(m1, m2) {
		t.Fatalf("Accounting not deterministic: %+v vs %+v", m1, m2)
	}

	// Depth: i1→e1→c1→d1 = 3 edges; width at sources = 2; fan-out of e1 wait i1,i2 fan-out 1, e1 fan-out 1.
	if m1.Nodes != 5 || m1.Edges != 4 {
		t.Fatalf("Nodes/Edges = %d/%d, want 5/4", m1.Nodes, m1.Edges)
	}
	if m1.Depth != 3 {
		t.Fatalf("Depth = %d, want 3", m1.Depth)
	}
	if m1.Width != 2 {
		t.Fatalf("Width = %d, want 2", m1.Width)
	}
	if m1.FanOut != 1 {
		t.Fatalf("FanOut = %d, want 1", m1.FanOut)
	}
	if m1.CriticalPath != 4 {
		t.Fatalf("CriticalPath = %d, want 4", m1.CriticalPath)
	}
	if m1.PayloadBytes != 3+5+7+11+13 {
		t.Fatalf("PayloadBytes = %d, want 39", m1.PayloadBytes)
	}

	// Recount: Accounting must count each node/edge once across calls.
	m3 := g1.Accounting()
	if !reflect.DeepEqual(m1, m3) {
		t.Fatalf("recount changed metrics: %+v vs %+v", m1, m3)
	}
}

func TestBuildRejectsCycle(t *testing.T) {
	t.Parallel()

	// Cycle among non-decision nodes (decision remains a terminal sink).
	cyc := Definition{
		Ref:     "g-cyc",
		Version: "1.0.0",
		Nodes: []Node{
			{ID: "a", Kind: NodeKindEvaluation},
			{ID: "b", Kind: NodeKindComposition},
			{ID: "d", Kind: NodeKindDecision},
		},
		Edges: []Edge{
			{ID: "e-a-b", From: "a", To: "b"},
			{ID: "e-b-a", From: "b", To: "a"},
			{ID: "e-b-d", From: "b", To: "d"},
		},
		Bounds: testBounds(),
	}
	_, prob := Build(cyc, testBounds())
	assertViolation(t, prob, violationCycle, "/edges")
}

func TestBuildRejectsUnknownNodeKind(t *testing.T) {
	t.Parallel()

	def := singleDecisionDef()
	def.Nodes[0].Kind = "workflow"
	_, prob := Build(def, testBounds())
	assertViolation(t, prob, violationGraphInvalid, "/nodes/0/kind")
}

func TestBuildRejectsLimitExceeded(t *testing.T) {
	t.Parallel()

	def := linearDef()
	tight := Bounds{
		MaxNodes:        2,
		MaxEdges:        32,
		MaxDepth:        8,
		MaxWidth:        8,
		MaxFanOut:       8,
		MaxPayloadBytes: 1024,
	}
	_, prob := Build(def, tight)
	assertViolation(t, prob, violationLimitExceeded, "/bounds/maxNodes")
}

func TestBuildPerformsNoLookup(t *testing.T) {
	t.Parallel()

	// Build accepts any already-resolved Ref/Version pair without contacting a registry.
	def := singleDecisionDef()
	def.Ref = "never-registered-graph"
	def.Version = "9.9.9-unresolved-elsewhere"
	if _, prob := Build(def, testBounds()); prob != nil {
		t.Fatalf("Build must not perform version/registry lookup: %#v", prob)
	}
}

func TestDefinitionJSONRoundTrip(t *testing.T) {
	t.Parallel()

	in := linearDef()
	in.ParallelGroups = []ParallelGroup{{ID: "pg", NodeIDs: []string{"in1"}}}
	data, err := json.Marshal(in)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var out Definition
	if err := json.Unmarshal(data, &out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if !reflect.DeepEqual(in, out) {
		t.Fatalf("round-trip mismatch:\n in=%+v\nout=%+v", in, out)
	}
}

func TestValidNodeKind(t *testing.T) {
	t.Parallel()

	for _, k := range []string{NodeKindInput, NodeKindEvaluation, NodeKindComposition, NodeKindDecision} {
		if !ValidNodeKind(k) {
			t.Fatalf("ValidNodeKind(%q) = false", k)
		}
	}
	if ValidNodeKind("task") {
		t.Fatal("ValidNodeKind(task) = true, want false")
	}
}

func assertViolation(t *testing.T, prob *apiproblem.Problem, code apiproblem.ViolationCode, field string) {
	t.Helper()
	if prob == nil {
		t.Fatalf("problem = nil, want code=%s field=%s", code, field)
	}
	if prob.Code != apiproblem.CodeValidationFailed {
		t.Fatalf("problem.Code = %s, want VALIDATION_FAILED", prob.Code)
	}
	if prob.Status != 422 {
		t.Fatalf("problem.Status = %d, want 422", prob.Status)
	}
	if len(prob.Violations) != 1 {
		t.Fatalf("violations = %#v, want exactly one", prob.Violations)
	}
	v := prob.Violations[0]
	if v.Code != code {
		t.Fatalf("violation.code = %s, want %s", v.Code, code)
	}
	if v.Field != field {
		t.Fatalf("violation.field = %q, want %q", v.Field, field)
	}
}
