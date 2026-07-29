// Package graph provides the FEATURE-0013 bounded in-memory decision-requirements
// graph (design §7.8; F13-COMP-003/004; AD-006, AD-040; closures 3/4).
//
// Build receives an already-resolved Definition and performs no registry or
// version lookup. The graph describes dependency and composition only; it does
// not execute provisioning, human tasks, compensations, arbitrary code, or
// unbounded loops (F13-COMP-004).
package graph

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
)

// Closed node-kind vocabulary (design §7.8). No other role is accepted.
const (
	NodeKindInput       = "input"
	NodeKindEvaluation  = "evaluation"
	NodeKindComposition = "composition"
	NodeKindDecision    = "decision"
)

// Closed composition violation codes used by Build (architecture §15.2).
// Canonical constants live in internal/decision/validate (T-017); these
// reproduce the same closed strings so Build can emit Problem Details now.
const (
	violationGraphInvalid  apiproblem.ViolationCode = "DECISION_COMPOSITION_GRAPH_INVALID"
	violationCycle         apiproblem.ViolationCode = "DECISION_COMPOSITION_CYCLE"
	violationLimitExceeded apiproblem.ViolationCode = "DECISION_COMPOSITION_LIMIT_EXCEEDED"
)

// Definition is the declarative graph definition (GraphDefinition).
// Version/reference compatibility is resolved by bundle.ResolveGraph before
// Build; Build never performs lookup (closure 4).
type Definition struct {
	Ref            string          `json:"ref"`
	Version        string          `json:"version"`
	Nodes          []Node          `json:"nodes"`
	Edges          []Edge          `json:"edges"`
	ParallelGroups []ParallelGroup `json:"parallelGroups,omitempty"`
	Bounds         Bounds          `json:"bounds"`
}

// Node is a typed graph node (GraphNode).
type Node struct {
	ID           string `json:"id"`
	Kind         string `json:"kind"`
	PayloadBytes int64  `json:"payloadBytes,omitempty"`
}

// Edge is a GraphEdge with a stable edge ID distinct from its From/To endpoints.
// Orientation is dependency → dependent.
type Edge struct {
	ID   string `json:"id"`
	From string `json:"from"`
	To   string `json:"to"`
}

// ParallelGroup declares mutually independent nodes that may evaluate in parallel.
// Independence means no direct or transitive dependency path exists between members.
type ParallelGroup struct {
	ID      string   `json:"id"`
	NodeIDs []string `json:"nodeIds"`
}

// Bounds are effective graph ceilings (GraphBounds) applied by Build.
// Callers supply the effective Bounds resolved under the inherited limit
// hierarchy (FEATURE-0012 → DecisionProfile → evaluator); Build does not
// resolve that hierarchy.
type Bounds struct {
	MaxNodes        int   `json:"maxNodes"`
	MaxEdges        int   `json:"maxEdges"`
	MaxDepth        int   `json:"maxDepth"`
	MaxWidth        int   `json:"maxWidth"`
	MaxFanOut       int   `json:"maxFanOut"`
	MaxPayloadBytes int64 `json:"maxPayloadBytes"`
}

// Metrics are deterministic graph accounting results (architecture §27.9(3)).
// Each node and each edge is counted once.
type Metrics struct {
	Nodes        int
	Edges        int
	Depth        int   // longest dependency path length in edges
	Width        int   // maximum node count at any depth level
	FanOut       int   // maximum outbound edge count from any single node
	CriticalPath int   // longest path length in nodes (each node once per path)
	PayloadBytes int64 // sum of node PayloadBytes (each node once)
}

// Graph is a bounded in-memory adjacency representation of a resolved Definition.
type Graph struct {
	def       Definition
	bounds    Bounds
	nodesByID map[string]Node
	edgesByID map[string]Edge
	nodeOrder []string // stable Definition.Nodes order
	edgeOrder []string // stable Definition.Edges order
	outEdges  map[string][]Edge
	inEdges   map[string][]Edge
}

// ValidNodeKind reports whether kind is one of the four approved node roles.
func ValidNodeKind(kind string) bool {
	switch kind {
	case NodeKindInput, NodeKindEvaluation, NodeKindComposition, NodeKindDecision:
		return true
	default:
		return false
	}
}

// Build constructs an adjacency Graph from an already-resolved Definition and
// the effective Bounds. It performs structural, cycle, parallel-group, and
// bound checks and returns inherited FEATURE-0012 Problem Details on failure.
// Build performs no registry or version lookup (closure 4).
//
// Duplicate (from,to) pairs are rejected (fail closed). Profile-permitted
// labelled multi-edges are a later composition/profile concern; this builder
// never invents a multi-edge permission.
func Build(def Definition, bounds Bounds) (Graph, *apiproblem.Problem) {
	if prob := validateBounds(bounds); prob != nil {
		return Graph{}, prob
	}
	if def.Ref == "" {
		return Graph{}, graphProblem(violationGraphInvalid, "/ref", "graph ref is required")
	}
	if def.Version == "" {
		return Graph{}, graphProblem(violationGraphInvalid, "/version", "graph version is required")
	}
	if len(def.Nodes) == 0 {
		return Graph{}, graphProblem(violationGraphInvalid, "/nodes", "graph nodes must be non-empty")
	}

	nodesByID := make(map[string]Node, len(def.Nodes))
	nodeOrder := make([]string, 0, len(def.Nodes))
	decisionIdx := -1
	var payloadSum int64

	for i, n := range def.Nodes {
		base := "/nodes/" + strconv.Itoa(i)
		if n.ID == "" {
			return Graph{}, graphProblem(violationGraphInvalid, base+"/id", "node id is required")
		}
		if _, exists := nodesByID[n.ID]; exists {
			return Graph{}, graphProblem(violationGraphInvalid, base+"/id", "duplicate node id")
		}
		if !ValidNodeKind(n.Kind) {
			return Graph{}, graphProblem(violationGraphInvalid, base+"/kind", "invalid node kind")
		}
		if n.PayloadBytes < 0 {
			return Graph{}, graphProblem(violationGraphInvalid, base+"/payloadBytes", "payloadBytes must be non-negative")
		}
		if bounds.MaxPayloadBytes > 0 && n.PayloadBytes > bounds.MaxPayloadBytes {
			return Graph{}, graphProblem(violationLimitExceeded, "/bounds/maxPayloadBytes", "node payload exceeds maxPayloadBytes")
		}
		if n.Kind == NodeKindDecision {
			if decisionIdx >= 0 {
				return Graph{}, graphProblem(violationGraphInvalid, "/nodes", "graph requires exactly one decision node")
			}
			decisionIdx = i
		}
		nodesByID[n.ID] = n
		nodeOrder = append(nodeOrder, n.ID)
		payloadSum += n.PayloadBytes
	}
	if decisionIdx < 0 {
		return Graph{}, graphProblem(violationGraphInvalid, "/nodes", "graph requires exactly one decision node")
	}

	if len(def.Nodes) == 1 {
		if len(def.Edges) != 0 {
			return Graph{}, graphProblem(violationGraphInvalid, "/edges", "single-node graph must have empty edges")
		}
		if def.Nodes[0].Kind != NodeKindDecision {
			return Graph{}, graphProblem(violationGraphInvalid, "/nodes", "empty-edge graph requires a sole decision node")
		}
	} else if len(def.Edges) == 0 {
		return Graph{}, graphProblem(violationGraphInvalid, "/nodes", "empty edges are legal only for a single-node decision graph")
	}

	edgesByID := make(map[string]Edge, len(def.Edges))
	edgeOrder := make([]string, 0, len(def.Edges))
	pairSeen := make(map[string]int, len(def.Edges)) // "from\0to" → first edge index
	outEdges := make(map[string][]Edge, len(def.Nodes))
	inEdges := make(map[string][]Edge, len(def.Nodes))
	for _, id := range nodeOrder {
		outEdges[id] = nil
		inEdges[id] = nil
	}

	for i, e := range def.Edges {
		base := "/edges/" + strconv.Itoa(i)
		if e.ID == "" {
			return Graph{}, graphProblem(violationGraphInvalid, base+"/id", "edge id is required")
		}
		if _, exists := edgesByID[e.ID]; exists {
			return Graph{}, graphProblem(violationGraphInvalid, base+"/id", "duplicate edge id")
		}
		if e.From == "" {
			return Graph{}, graphProblem(violationGraphInvalid, base+"/from", "edge from is required")
		}
		if e.To == "" {
			return Graph{}, graphProblem(violationGraphInvalid, base+"/to", "edge to is required")
		}
		if e.ID == e.From || e.ID == e.To {
			return Graph{}, graphProblem(violationGraphInvalid, base+"/id", "edge id must be distinct from from/to endpoints")
		}
		if _, ok := nodesByID[e.From]; !ok {
			return Graph{}, graphProblem(violationGraphInvalid, base+"/from", "edge from references unknown node")
		}
		if _, ok := nodesByID[e.To]; !ok {
			return Graph{}, graphProblem(violationGraphInvalid, base+"/to", "edge to references unknown node")
		}
		if e.From == e.To {
			return Graph{}, graphProblem(violationGraphInvalid, base, "self-loop edges are invalid")
		}
		pairKey := e.From + "\x00" + e.To
		if prev, ok := pairSeen[pairKey]; ok {
			_ = prev
			return Graph{}, graphProblem(violationGraphInvalid, base, "duplicate (from,to) pair is invalid")
		}
		pairSeen[pairKey] = i

		edgesByID[e.ID] = e
		edgeOrder = append(edgeOrder, e.ID)
		outEdges[e.From] = append(outEdges[e.From], e)
		inEdges[e.To] = append(inEdges[e.To], e)
	}

	// Deterministic adjacency order by edge ID.
	for _, id := range nodeOrder {
		sort.SliceStable(outEdges[id], func(i, j int) bool {
			return outEdges[id][i].ID < outEdges[id][j].ID
		})
		sort.SliceStable(inEdges[id], func(i, j int) bool {
			return inEdges[id][i].ID < inEdges[id][j].ID
		})
	}

	decisionID := def.Nodes[decisionIdx].ID
	if len(outEdges[decisionID]) != 0 {
		return Graph{}, graphProblem(violationGraphInvalid, "/nodes", "decision node must be a terminal sink")
	}
	if len(def.Edges) > 0 {
		if prob := validateReachability(nodeOrder, outEdges, decisionID); prob != nil {
			return Graph{}, prob
		}
	}

	groupIDs := make(map[string]struct{}, len(def.ParallelGroups))
	for i, g := range def.ParallelGroups {
		base := "/parallelGroups/" + strconv.Itoa(i)
		if g.ID == "" {
			return Graph{}, graphProblem(violationGraphInvalid, base+"/id", "parallel group id is required")
		}
		if _, exists := groupIDs[g.ID]; exists {
			return Graph{}, graphProblem(violationGraphInvalid, base+"/id", "duplicate parallel group id")
		}
		groupIDs[g.ID] = struct{}{}
		if len(g.NodeIDs) == 0 {
			return Graph{}, graphProblem(violationGraphInvalid, base+"/nodeIds", "parallel group nodeIds must be non-empty")
		}
		seenMember := make(map[string]struct{}, len(g.NodeIDs))
		for j, member := range g.NodeIDs {
			memberPath := base + "/nodeIds/" + strconv.Itoa(j)
			if member == "" {
				return Graph{}, graphProblem(violationGraphInvalid, memberPath, "parallel group member id is required")
			}
			if _, ok := nodesByID[member]; !ok {
				return Graph{}, graphProblem(violationGraphInvalid, memberPath, "parallel group member references unknown node")
			}
			if _, dup := seenMember[member]; dup {
				return Graph{}, graphProblem(violationGraphInvalid, memberPath, "duplicate parallel group member")
			}
			seenMember[member] = struct{}{}
		}
		if prob := validateParallelIndependence(g, i, outEdges); prob != nil {
			return Graph{}, prob
		}
	}

	if hasCycle(nodeOrder, outEdges) {
		return Graph{}, graphProblem(violationCycle, "/edges", "composition graph contains a cycle")
	}

	g := Graph{
		def:       cloneDefinition(def),
		bounds:    bounds,
		nodesByID: nodesByID,
		edgesByID: edgesByID,
		nodeOrder: nodeOrder,
		edgeOrder: edgeOrder,
		outEdges:  outEdges,
		inEdges:   inEdges,
	}

	m := g.Accounting()
	if m.Nodes > bounds.MaxNodes {
		return Graph{}, graphProblem(violationLimitExceeded, "/bounds/maxNodes", "node count exceeds maxNodes")
	}
	if m.Edges > bounds.MaxEdges {
		return Graph{}, graphProblem(violationLimitExceeded, "/bounds/maxEdges", "edge count exceeds maxEdges")
	}
	if m.Depth > bounds.MaxDepth {
		return Graph{}, graphProblem(violationLimitExceeded, "/bounds/maxDepth", "depth exceeds maxDepth")
	}
	if m.Width > bounds.MaxWidth {
		return Graph{}, graphProblem(violationLimitExceeded, "/bounds/maxWidth", "width exceeds maxWidth")
	}
	if m.FanOut > bounds.MaxFanOut {
		return Graph{}, graphProblem(violationLimitExceeded, "/bounds/maxFanOut", "fan-out exceeds maxFanOut")
	}
	if payloadSum > bounds.MaxPayloadBytes {
		return Graph{}, graphProblem(violationLimitExceeded, "/bounds/maxPayloadBytes", "total payload exceeds maxPayloadBytes")
	}

	return g, nil
}

// Accounting returns deterministic depth/width/fan-out/critical-path metrics.
// Each node and each edge is counted once (architecture §27.9(3)).
func (g Graph) Accounting() Metrics {
	n := len(g.nodeOrder)
	e := len(g.edgeOrder)
	if n == 0 {
		return Metrics{}
	}

	depthByNode := make(map[string]int, n)
	var payloadSum int64
	fanOut := 0

	// Kahn-style level assignment from sources; depth = longest inbound path (edges).
	inDegree := make(map[string]int, n)
	for _, id := range g.nodeOrder {
		inDegree[id] = len(g.inEdges[id])
		payloadSum += g.nodesByID[id].PayloadBytes
		if fo := len(g.outEdges[id]); fo > fanOut {
			fanOut = fo
		}
	}

	queue := make([]string, 0, n)
	for _, id := range g.nodeOrder {
		if inDegree[id] == 0 {
			depthByNode[id] = 0
			queue = append(queue, id)
		}
	}

	seen := 0
	for len(queue) > 0 {
		// Stable: process current frontier in nodeOrder relative order.
		sort.SliceStable(queue, func(i, j int) bool {
			return indexOf(g.nodeOrder, queue[i]) < indexOf(g.nodeOrder, queue[j])
		})
		id := queue[0]
		queue = queue[1:]
		seen++
		for _, edge := range g.outEdges[id] {
			cand := depthByNode[id] + 1
			if cur, ok := depthByNode[edge.To]; !ok || cand > cur {
				depthByNode[edge.To] = cand
			}
			inDegree[edge.To]--
			if inDegree[edge.To] == 0 {
				queue = append(queue, edge.To)
			}
		}
	}
	if seen != n {
		// Should not occur on a successfully Built graph; keep metrics zero-safe.
		return Metrics{Nodes: n, Edges: e, FanOut: fanOut, PayloadBytes: payloadSum}
	}

	depth := 0
	widthByLevel := map[int]int{}
	for _, id := range g.nodeOrder {
		d := depthByNode[id]
		if d > depth {
			depth = d
		}
		widthByLevel[d]++
	}
	width := 0
	for _, c := range widthByLevel {
		if c > width {
			width = c
		}
	}

	criticalPath := 0
	for _, id := range g.nodeOrder {
		// Longest path node count = depth(id) + 1.
		cp := depthByNode[id] + 1
		if cp > criticalPath {
			criticalPath = cp
		}
	}

	return Metrics{
		Nodes:        n,
		Edges:        e,
		Depth:        depth,
		Width:        width,
		FanOut:       fanOut,
		CriticalPath: criticalPath,
		PayloadBytes: payloadSum,
	}
}

// Definition returns a deep copy of the graph definition.
func (g Graph) Definition() Definition {
	return cloneDefinition(g.def)
}

// Bounds returns the effective bounds applied at Build time.
func (g Graph) Bounds() Bounds {
	return g.bounds
}

// Node returns the node with the given id.
func (g Graph) Node(id string) (Node, bool) {
	n, ok := g.nodesByID[id]
	return n, ok
}

// OutEdges returns outbound edges from node id in stable edge-ID order.
func (g Graph) OutEdges(id string) []Edge {
	src := g.outEdges[id]
	if len(src) == 0 {
		return nil
	}
	out := make([]Edge, len(src))
	copy(out, src)
	return out
}

// InEdges returns inbound edges to node id in stable edge-ID order.
func (g Graph) InEdges(id string) []Edge {
	src := g.inEdges[id]
	if len(src) == 0 {
		return nil
	}
	out := make([]Edge, len(src))
	copy(out, src)
	return out
}

// NodeIDs returns node IDs in Definition order.
func (g Graph) NodeIDs() []string {
	out := make([]string, len(g.nodeOrder))
	copy(out, g.nodeOrder)
	return out
}

func validateBounds(b Bounds) *apiproblem.Problem {
	switch {
	case b.MaxNodes <= 0:
		return graphProblem(violationGraphInvalid, "/bounds/maxNodes", "maxNodes must be positive")
	case b.MaxEdges < 0:
		return graphProblem(violationGraphInvalid, "/bounds/maxEdges", "maxEdges must be non-negative")
	case b.MaxDepth < 0:
		return graphProblem(violationGraphInvalid, "/bounds/maxDepth", "maxDepth must be non-negative")
	case b.MaxWidth <= 0:
		return graphProblem(violationGraphInvalid, "/bounds/maxWidth", "maxWidth must be positive")
	case b.MaxFanOut < 0:
		return graphProblem(violationGraphInvalid, "/bounds/maxFanOut", "maxFanOut must be non-negative")
	case b.MaxPayloadBytes < 0:
		return graphProblem(violationGraphInvalid, "/bounds/maxPayloadBytes", "maxPayloadBytes must be non-negative")
	default:
		return nil
	}
}

func validateReachability(nodeOrder []string, outEdges map[string][]Edge, decisionID string) *apiproblem.Problem {
	reachable := map[string]struct{}{decisionID: {}}
	// Walk reverse: build reverse adjacency from outEdges.
	preds := make(map[string][]string, len(nodeOrder))
	for from, edges := range outEdges {
		for _, e := range edges {
			preds[e.To] = append(preds[e.To], from)
		}
	}
	stack := []string{decisionID}
	for len(stack) > 0 {
		cur := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		for _, p := range preds[cur] {
			if _, ok := reachable[p]; ok {
				continue
			}
			reachable[p] = struct{}{}
			stack = append(stack, p)
		}
	}
	for _, id := range nodeOrder {
		if _, ok := reachable[id]; !ok {
			return graphProblem(violationGraphInvalid, "/nodes", "all nodes must reach the decision terminal")
		}
	}
	return nil
}

func validateParallelIndependence(g ParallelGroup, groupIndex int, outEdges map[string][]Edge) *apiproblem.Problem {
	base := "/parallelGroups/" + strconv.Itoa(groupIndex)
	for i := 0; i < len(g.NodeIDs); i++ {
		for j := i + 1; j < len(g.NodeIDs); j++ {
			a, b := g.NodeIDs[i], g.NodeIDs[j]
			if pathExists(a, b, outEdges) || pathExists(b, a, outEdges) {
				// Point at the later member that violates independence.
				return graphProblem(
					violationGraphInvalid,
					base+"/nodeIds/"+strconv.Itoa(j),
					fmt.Sprintf("parallel group members %q and %q are not independent", a, b),
				)
			}
		}
	}
	return nil
}

func pathExists(from, to string, outEdges map[string][]Edge) bool {
	if from == to {
		return true
	}
	seen := map[string]struct{}{from: {}}
	stack := []string{from}
	for len(stack) > 0 {
		cur := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		for _, e := range outEdges[cur] {
			if e.To == to {
				return true
			}
			if _, ok := seen[e.To]; ok {
				continue
			}
			seen[e.To] = struct{}{}
			stack = append(stack, e.To)
		}
	}
	return false
}

func hasCycle(nodeOrder []string, outEdges map[string][]Edge) bool {
	const (
		white = 0
		gray  = 1
		black = 2
	)
	color := make(map[string]int, len(nodeOrder))
	var visit func(string) bool
	visit = func(id string) bool {
		color[id] = gray
		for _, e := range outEdges[id] {
			switch color[e.To] {
			case gray:
				return true
			case white:
				if visit(e.To) {
					return true
				}
			}
		}
		color[id] = black
		return false
	}
	for _, id := range nodeOrder {
		if color[id] == white {
			if visit(id) {
				return true
			}
		}
	}
	return false
}

func graphProblem(code apiproblem.ViolationCode, field, message string) *apiproblem.Problem {
	return apiproblem.New(apiproblem.CodeValidationFailed).
		WithDetail(message).
		WithViolations([]apiproblem.Violation{{
			Field:   field,
			Code:    code,
			Message: message,
		}})
}

func cloneDefinition(def Definition) Definition {
	out := Definition{
		Ref:     def.Ref,
		Version: def.Version,
		Bounds:  def.Bounds,
	}
	if len(def.Nodes) > 0 {
		out.Nodes = make([]Node, len(def.Nodes))
		copy(out.Nodes, def.Nodes)
	}
	if len(def.Edges) > 0 {
		out.Edges = make([]Edge, len(def.Edges))
		copy(out.Edges, def.Edges)
	}
	if len(def.ParallelGroups) > 0 {
		out.ParallelGroups = make([]ParallelGroup, len(def.ParallelGroups))
		for i, g := range def.ParallelGroups {
			out.ParallelGroups[i] = ParallelGroup{ID: g.ID}
			if len(g.NodeIDs) > 0 {
				out.ParallelGroups[i].NodeIDs = make([]string, len(g.NodeIDs))
				copy(out.ParallelGroups[i].NodeIDs, g.NodeIDs)
			}
		}
	}
	return out
}

func indexOf(order []string, id string) int {
	for i, v := range order {
		if v == id {
			return i
		}
	}
	return len(order)
}
