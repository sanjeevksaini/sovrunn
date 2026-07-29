// Package bundle provides the FEATURE-0013 declarative decision-profile
// bundle and (id, version)-keyed registry views (design §6.3, §8, §10;
// F13-PROF-004/007/008, F13-COMPAT-005; AD-025, AD-038; closures 8/13).
//
// The JSON/YAML bundle document is the semantic source of truth. Go maps in
// BundleView are derivative views. Offline verification only: no network
// fetch, digest computation, or signature verification. Unknown, expired,
// revoked, or mismatched trust states fail closed with no fallback to an
// older version.
package bundle

import (
	"strconv"

	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/decision"
	"github.com/sanjeevksaini/sovrunn/internal/decision/graph"
)

// Canonical DecisionProfileBundle TypeMeta values (design §6.1/§10).
const (
	APIVersionDecisionProfileBundle = "governance.sovrunn.io/v1alpha1"
	KindDecisionProfileBundle       = "DecisionProfileBundle"
)

// Closed composition/trust violation codes used by Load and ResolveGraph.
// Canonical constants live in internal/decision/validate (T-017); these
// reproduce the same closed strings so the bundle package can emit Problem
// Details now.
const (
	violationProfileSchemaInvalid apiproblem.ViolationCode = "DECISION_PROFILE_SCHEMA_INVALID"
	violationGraphInvalid         apiproblem.ViolationCode = "DECISION_COMPOSITION_GRAPH_INVALID"
	violationTrustUnknown         apiproblem.ViolationCode = "DECISION_TRUST_UNKNOWN"
	violationTrustExpired         apiproblem.ViolationCode = "DECISION_TRUST_EXPIRED"
	violationTrustRevoked         apiproblem.ViolationCode = "DECISION_TRUST_REVOKED"
	violationTrustMismatch        apiproblem.ViolationCode = "DECISION_TRUST_MISMATCH"
)

// Opaque structural trust-state tokens that fail closed at Load time
// (design §10; F13-PROF-008; F13-CF-23). These are not a new closed
// vocabulary type; they are the opaque fail-closed transitions architecture
// requires before ADR-F13-002 cryptographic selection.
const (
	TrustStateUnknown = "unknown"
	TrustStateExpired = "expired"
	TrustStateRevoked = "revoked"
)

// VersionKey is the composite (id, version) registry key (closure 8).
type VersionKey struct {
	ID      string
	Version string
}

// GraphRef is a version-pinned graph identity resolved by ResolveGraph
// before graph.Build (closure 4; design §8).
type GraphRef struct {
	Ref     string `json:"ref"`
	Version string `json:"version"`
}

// VersionedEntry is one versioned semantic-registry identity pin
// (F13-COMPAT-005). Entries are structural identifiers only.
type VersionedEntry struct {
	ID      string `json:"id"`
	Version string `json:"version"`
}

// Bundle is the declarative decision-profile / semantic-registry document.
// It is the semantic source of truth; BundleView maps are derivative.
type Bundle struct {
	APIVersion string     `json:"apiVersion"`
	Kind       string     `json:"kind"`
	Spec       BundleSpec `json:"spec"`

	view BundleView // populated by Load; not serialized
}

// BundleSpec is the bundle definition payload (design §10).
// Profiles and graphs are keyed by composite (id, version) after Load.
// Semantic registries follow the same key rule (closure 8).
type BundleSpec struct {
	Profiles []decision.DecisionProfile `json:"profiles"`
	Graphs   []graph.Definition         `json:"graphs,omitempty"`

	Forms          []VersionedEntry `json:"forms,omitempty"`
	Authorities    []VersionedEntry `json:"authorities,omitempty"`
	EvaluatorTypes []VersionedEntry `json:"evaluatorTypes,omitempty"`
	Strategies     []VersionedEntry `json:"strategies,omitempty"`
	ReasonCodes    []VersionedEntry `json:"reasonCodes,omitempty"`
	Obligations    []VersionedEntry `json:"obligations,omitempty"`
	EventTypes     []VersionedEntry `json:"eventTypes,omitempty"`
	EvidenceTypes  []VersionedEntry `json:"evidenceTypes,omitempty"`
	Projections    []VersionedEntry `json:"projections,omitempty"`

	// Trust is the algorithm-agile structural trust carrier for the bundle
	// (AD-022, AD-038). Present-but-empty is invalid (RID-08).
	Trust *decision.TrustCarrier `json:"trust,omitempty"`
	// ExpectedTrustState is the opaque expected trust-state token. When set,
	// it must equal Trust.State; mismatch fails closed with no digest check.
	ExpectedTrustState string `json:"expectedTrustState,omitempty"`
}

// BundleView is a derivative (id, version)-keyed view of a loaded Bundle.
// It is never the semantic authority (F13-PROF-004).
type BundleView struct {
	profiles map[VersionKey]decision.DecisionProfile
	graphs   map[VersionKey]graph.Definition

	forms          map[VersionKey]VersionedEntry
	authorities    map[VersionKey]VersionedEntry
	evaluatorTypes map[VersionKey]VersionedEntry
	strategies     map[VersionKey]VersionedEntry
	reasonCodes    map[VersionKey]VersionedEntry
	obligations    map[VersionKey]VersionedEntry
	eventTypes     map[VersionKey]VersionedEntry
	evidenceTypes  map[VersionKey]VersionedEntry
	projections    map[VersionKey]VersionedEntry
}

// View returns the derivative registry view built at Load time.
func (b Bundle) View() BundleView {
	return b.view
}

// Profile looks up a DecisionProfile by composite (id, version).
// Lookup is exact; there is no fallback to another version (F13-PROF-007).
func (v BundleView) Profile(id, version string) (decision.DecisionProfile, bool) {
	if v.profiles == nil {
		return decision.DecisionProfile{}, false
	}
	p, ok := v.profiles[VersionKey{ID: id, Version: version}]
	return p, ok
}

// Graph looks up a resolved graph Definition by composite (ref, version).
func (v BundleView) Graph(ref, version string) (graph.Definition, bool) {
	if v.graphs == nil {
		return graph.Definition{}, false
	}
	g, ok := v.graphs[VersionKey{ID: ref, Version: version}]
	return g, ok
}

// Forms returns the forms registry entry for (id, version), if present.
func (v BundleView) Forms(id, version string) (VersionedEntry, bool) {
	return lookupEntry(v.forms, id, version)
}

// Authorities returns the authorities registry entry for (id, version).
func (v BundleView) Authorities(id, version string) (VersionedEntry, bool) {
	return lookupEntry(v.authorities, id, version)
}

// EvaluatorTypes returns the evaluator-types registry entry for (id, version).
func (v BundleView) EvaluatorTypes(id, version string) (VersionedEntry, bool) {
	return lookupEntry(v.evaluatorTypes, id, version)
}

// Strategies returns the strategies registry entry for (id, version).
func (v BundleView) Strategies(id, version string) (VersionedEntry, bool) {
	return lookupEntry(v.strategies, id, version)
}

// ReasonCodes returns the reason-codes registry entry for (id, version).
func (v BundleView) ReasonCodes(id, version string) (VersionedEntry, bool) {
	return lookupEntry(v.reasonCodes, id, version)
}

// Obligations returns the obligations registry entry for (id, version).
func (v BundleView) Obligations(id, version string) (VersionedEntry, bool) {
	return lookupEntry(v.obligations, id, version)
}

// EventTypes returns the event-types registry entry for (id, version).
func (v BundleView) EventTypes(id, version string) (VersionedEntry, bool) {
	return lookupEntry(v.eventTypes, id, version)
}

// EvidenceTypes returns the evidence-types registry entry for (id, version).
func (v BundleView) EvidenceTypes(id, version string) (VersionedEntry, bool) {
	return lookupEntry(v.evidenceTypes, id, version)
}

// Projections returns the projections registry entry for (id, version).
func (v BundleView) Projections(id, version string) (VersionedEntry, bool) {
	return lookupEntry(v.projections, id, version)
}

// ResolveGraph performs graph reference/version resolution against the
// declarative bundle and returns the already-resolved Definition for
// graph.Build. It performs no Build and no fallback to another version
// (closure 4; design §8/§10).
func (b Bundle) ResolveGraph(ref GraphRef) (graph.Definition, *apiproblem.Problem) {
	if ref.Ref == "" {
		return graph.Definition{}, bundleProblem(violationGraphInvalid, "/ref", "graph ref is required")
	}
	if ref.Version == "" {
		return graph.Definition{}, bundleProblem(violationGraphInvalid, "/version", "graph version is required")
	}
	def, ok := b.view.Graph(ref.Ref, ref.Version)
	if !ok {
		return graph.Definition{}, bundleProblem(
			violationGraphInvalid,
			"/",
			"graph reference/version not found in bundle; no version fallback",
		)
	}
	// Defensive pin: returned definition must match the requested key.
	if def.Ref != ref.Ref || def.Version != ref.Version {
		return graph.Definition{}, bundleProblem(
			violationGraphInvalid,
			"/",
			"graph reference/version mismatch; no version fallback",
		)
	}
	return cloneGraphDefinition(def), nil
}

func lookupEntry(m map[VersionKey]VersionedEntry, id, version string) (VersionedEntry, bool) {
	if m == nil {
		return VersionedEntry{}, false
	}
	e, ok := m[VersionKey{ID: id, Version: version}]
	return e, ok
}

// buildView constructs the derivative BundleView and rejects duplicate
// (id, version) keys fail-closed (closure 8).
func buildView(spec BundleSpec) (BundleView, *apiproblem.Problem) {
	v := BundleView{
		profiles:       make(map[VersionKey]decision.DecisionProfile),
		graphs:         make(map[VersionKey]graph.Definition),
		forms:          make(map[VersionKey]VersionedEntry),
		authorities:    make(map[VersionKey]VersionedEntry),
		evaluatorTypes: make(map[VersionKey]VersionedEntry),
		strategies:     make(map[VersionKey]VersionedEntry),
		reasonCodes:    make(map[VersionKey]VersionedEntry),
		obligations:    make(map[VersionKey]VersionedEntry),
		eventTypes:     make(map[VersionKey]VersionedEntry),
		evidenceTypes:  make(map[VersionKey]VersionedEntry),
		projections:    make(map[VersionKey]VersionedEntry),
	}

	for i, p := range spec.Profiles {
		base := "/spec/profiles/" + strconv.Itoa(i)
		id := p.Spec.Name
		version := p.Spec.Version
		if id == "" {
			return BundleView{}, bundleProblem(violationProfileSchemaInvalid, base+"/spec/name", "profile name is required")
		}
		if version == "" {
			return BundleView{}, bundleProblem(violationProfileSchemaInvalid, base+"/spec/version", "profile version is required")
		}
		key := VersionKey{ID: id, Version: version}
		if _, exists := v.profiles[key]; exists {
			return BundleView{}, bundleProblem(violationProfileSchemaInvalid, base, "duplicate profile (id,version)")
		}
		v.profiles[key] = p
	}

	for i, g := range spec.Graphs {
		base := "/spec/graphs/" + strconv.Itoa(i)
		if g.Ref == "" {
			return BundleView{}, bundleProblem(violationGraphInvalid, base+"/ref", "graph ref is required")
		}
		if g.Version == "" {
			return BundleView{}, bundleProblem(violationGraphInvalid, base+"/version", "graph version is required")
		}
		key := VersionKey{ID: g.Ref, Version: g.Version}
		if _, exists := v.graphs[key]; exists {
			return BundleView{}, bundleProblem(violationGraphInvalid, base, "duplicate graph (ref,version)")
		}
		v.graphs[key] = g
	}

	if prob := indexEntries(spec.Forms, "/spec/forms", v.forms); prob != nil {
		return BundleView{}, prob
	}
	if prob := indexEntries(spec.Authorities, "/spec/authorities", v.authorities); prob != nil {
		return BundleView{}, prob
	}
	if prob := indexEntries(spec.EvaluatorTypes, "/spec/evaluatorTypes", v.evaluatorTypes); prob != nil {
		return BundleView{}, prob
	}
	if prob := indexEntries(spec.Strategies, "/spec/strategies", v.strategies); prob != nil {
		return BundleView{}, prob
	}
	if prob := indexEntries(spec.ReasonCodes, "/spec/reasonCodes", v.reasonCodes); prob != nil {
		return BundleView{}, prob
	}
	if prob := indexEntries(spec.Obligations, "/spec/obligations", v.obligations); prob != nil {
		return BundleView{}, prob
	}
	if prob := indexEntries(spec.EventTypes, "/spec/eventTypes", v.eventTypes); prob != nil {
		return BundleView{}, prob
	}
	if prob := indexEntries(spec.EvidenceTypes, "/spec/evidenceTypes", v.evidenceTypes); prob != nil {
		return BundleView{}, prob
	}
	if prob := indexEntries(spec.Projections, "/spec/projections", v.projections); prob != nil {
		return BundleView{}, prob
	}

	return v, nil
}

func indexEntries(entries []VersionedEntry, basePath string, dst map[VersionKey]VersionedEntry) *apiproblem.Problem {
	for i, e := range entries {
		base := basePath + "/" + strconv.Itoa(i)
		if e.ID == "" {
			return bundleProblem(violationProfileSchemaInvalid, base+"/id", "registry entry id is required")
		}
		if e.Version == "" {
			return bundleProblem(violationProfileSchemaInvalid, base+"/version", "registry entry version is required")
		}
		key := VersionKey{ID: e.ID, Version: e.Version}
		if _, exists := dst[key]; exists {
			return bundleProblem(violationProfileSchemaInvalid, base, "duplicate registry (id,version)")
		}
		dst[key] = e
	}
	return nil
}

// validateBundleTrust applies offline structural trust fail-closed rules
// (design §10; RID-08; F13-PROF-008). No digest, signature, or network check.
func validateBundleTrust(spec BundleSpec) *apiproblem.Problem {
	if spec.Trust == nil {
		if spec.ExpectedTrustState != "" {
			return bundleProblem(violationTrustMismatch, "/spec/expectedTrustState", "expected trust state set but trust carrier absent")
		}
		return nil
	}
	if spec.Trust.IsStructurallyEmpty() {
		return bundleProblem(violationTrustUnknown, "/spec/trust", "present-but-empty trust carrier is invalid")
	}
	switch spec.Trust.State {
	case TrustStateUnknown:
		return bundleProblem(violationTrustUnknown, "/spec/trust/state", "unknown trust state fails closed")
	case TrustStateExpired:
		return bundleProblem(violationTrustExpired, "/spec/trust/state", "expired trust state fails closed")
	case TrustStateRevoked:
		return bundleProblem(violationTrustRevoked, "/spec/trust/state", "revoked trust state fails closed")
	}
	if spec.ExpectedTrustState != "" && spec.Trust.State != spec.ExpectedTrustState {
		return bundleProblem(violationTrustMismatch, "/spec/trust/state", "expected versus observed trust state mismatch; no fallback")
	}
	return nil
}

func bundleProblem(code apiproblem.ViolationCode, field, message string) *apiproblem.Problem {
	return apiproblem.New(apiproblem.CodeValidationFailed).
		WithDetail(message).
		WithViolations([]apiproblem.Violation{{
			Field:   field,
			Code:    code,
			Message: message,
		}})
}

func cloneGraphDefinition(def graph.Definition) graph.Definition {
	out := graph.Definition{
		Ref:     def.Ref,
		Version: def.Version,
		Bounds:  def.Bounds,
	}
	if len(def.Nodes) > 0 {
		out.Nodes = make([]graph.Node, len(def.Nodes))
		copy(out.Nodes, def.Nodes)
	}
	if len(def.Edges) > 0 {
		out.Edges = make([]graph.Edge, len(def.Edges))
		copy(out.Edges, def.Edges)
	}
	if len(def.ParallelGroups) > 0 {
		out.ParallelGroups = make([]graph.ParallelGroup, len(def.ParallelGroups))
		for i, g := range def.ParallelGroups {
			out.ParallelGroups[i] = graph.ParallelGroup{ID: g.ID}
			if len(g.NodeIDs) > 0 {
				out.ParallelGroups[i].NodeIDs = make([]string, len(g.NodeIDs))
				copy(out.ParallelGroups[i].NodeIDs, g.NodeIDs)
			}
		}
	}
	return out
}
