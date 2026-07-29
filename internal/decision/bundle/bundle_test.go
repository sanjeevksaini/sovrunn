package bundle

import (
	"encoding/json"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/apivalid"
	"github.com/sanjeevksaini/sovrunn/internal/decision"
	"github.com/sanjeevksaini/sovrunn/internal/decision/graph"
)

func minimalProfile(name, version string) decision.DecisionProfile {
	return decision.DecisionProfile{
		TypeMeta: apimeta.TypeMeta{
			APIVersion: decision.APIVersionDecisionProfile,
			Kind:       decision.KindDecisionProfile,
		},
		Metadata: apimeta.ObjectMeta{Name: name},
		Spec: decision.ProfileSpec{
			Name:    name,
			Version: version,
			Family:  "test",
			Owner:   "owner",
		},
	}
}

func singleGraph(ref, version string) graph.Definition {
	return graph.Definition{
		Ref:     ref,
		Version: version,
		Nodes: []graph.Node{
			{ID: "d1", Kind: graph.NodeKindDecision},
		},
		Edges: []graph.Edge{},
		Bounds: graph.Bounds{
			MaxNodes: 4, MaxEdges: 4, MaxDepth: 2, MaxWidth: 2, MaxFanOut: 2, MaxPayloadBytes: 64,
		},
	}
}

func mustLoadJSON(t *testing.T, b Bundle) Bundle {
	t.Helper()
	raw, err := json.Marshal(b)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	out, prob := Load(raw, apivalid.ModeReadRepresentation)
	if prob != nil {
		t.Fatalf("Load = %#v, want nil", prob)
	}
	return out
}

func TestDuplicateProfileIDVersionRejected(t *testing.T) {
	t.Parallel()

	doc := Bundle{
		APIVersion: APIVersionDecisionProfileBundle,
		Kind:       KindDecisionProfileBundle,
		Spec: BundleSpec{
			Profiles: []decision.DecisionProfile{
				minimalProfile("authz", "1.0.0"),
				minimalProfile("authz", "1.0.0"),
			},
		},
	}
	raw, err := json.Marshal(doc)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	_, prob := Load(raw, apivalid.ModeReadRepresentation)
	assertViolation(t, prob, violationProfileSchemaInvalid, "/spec/profiles/1")
}

func TestDuplicateGraphRefVersionRejected(t *testing.T) {
	t.Parallel()

	doc := Bundle{
		APIVersion: APIVersionDecisionProfileBundle,
		Kind:       KindDecisionProfileBundle,
		Spec: BundleSpec{
			Profiles: []decision.DecisionProfile{minimalProfile("authz", "1.0.0")},
			Graphs: []graph.Definition{
				singleGraph("g1", "1.0.0"),
				singleGraph("g1", "1.0.0"),
			},
		},
	}
	raw, err := json.Marshal(doc)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	_, prob := Load(raw, apivalid.ModeReadRepresentation)
	assertViolation(t, prob, violationGraphInvalid, "/spec/graphs/1")
}

func TestDuplicateRegistryEntryRejected(t *testing.T) {
	t.Parallel()

	doc := Bundle{
		APIVersion: APIVersionDecisionProfileBundle,
		Kind:       KindDecisionProfileBundle,
		Spec: BundleSpec{
			Profiles: []decision.DecisionProfile{minimalProfile("authz", "1.0.0")},
			Forms: []VersionedEntry{
				{ID: "ADJUDICATION", Version: "1.0.0"},
				{ID: "ADJUDICATION", Version: "1.0.0"},
			},
		},
	}
	raw, err := json.Marshal(doc)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	_, prob := Load(raw, apivalid.ModeReadRepresentation)
	assertViolation(t, prob, violationProfileSchemaInvalid, "/spec/forms/1")
}

func TestProfileLookupByCompositeKey(t *testing.T) {
	t.Parallel()

	loaded := mustLoadJSON(t, Bundle{
		APIVersion: APIVersionDecisionProfileBundle,
		Kind:       KindDecisionProfileBundle,
		Spec: BundleSpec{
			Profiles: []decision.DecisionProfile{
				minimalProfile("authz", "1.0.0"),
				minimalProfile("authz", "2.0.0"),
				minimalProfile("placement", "1.0.0"),
			},
			Forms: []VersionedEntry{
				{ID: "ADJUDICATION", Version: "1.0.0"},
			},
			Strategies: []VersionedEntry{
				{ID: "allOf", Version: "1.0.0"},
			},
		},
	})

	view := loaded.View()
	p, ok := view.Profile("authz", "2.0.0")
	if !ok {
		t.Fatal("Profile(authz, 2.0.0) missing")
	}
	if p.Spec.Name != "authz" || p.Spec.Version != "2.0.0" {
		t.Fatalf("Profile = %+v", p.Spec)
	}

	// Exact key only — no fallback to older version.
	if _, ok := view.Profile("authz", "3.0.0"); ok {
		t.Fatal("unexpected fallback to another profile version")
	}
	if _, ok := view.Profile("missing", "1.0.0"); ok {
		t.Fatal("unexpected hit for unknown profile id")
	}

	if _, ok := view.Forms("ADJUDICATION", "1.0.0"); !ok {
		t.Fatal("Forms lookup missed")
	}
	if _, ok := view.Strategies("allOf", "1.0.0"); !ok {
		t.Fatal("Strategies lookup missed")
	}
	if _, ok := view.Forms("ADJUDICATION", "9.9.9"); ok {
		t.Fatal("unexpected forms version fallback")
	}
}

func TestResolveGraphBeforeBuild(t *testing.T) {
	t.Parallel()

	loaded := mustLoadJSON(t, Bundle{
		APIVersion: APIVersionDecisionProfileBundle,
		Kind:       KindDecisionProfileBundle,
		Spec: BundleSpec{
			Profiles: []decision.DecisionProfile{minimalProfile("authz", "1.0.0")},
			Graphs: []graph.Definition{
				singleGraph("g-linear", "1.0.0"),
				singleGraph("g-linear", "2.0.0"),
			},
		},
	})

	def, prob := loaded.ResolveGraph(GraphRef{Ref: "g-linear", Version: "2.0.0"})
	if prob != nil {
		t.Fatalf("ResolveGraph = %#v, want nil", prob)
	}
	if def.Ref != "g-linear" || def.Version != "2.0.0" {
		t.Fatalf("resolved def = %+v", def)
	}

	// Resolution succeeds; Build is a separate step and must not be required here.
	g, buildProb := graph.Build(def, def.Bounds)
	if buildProb != nil {
		t.Fatalf("Build after ResolveGraph = %#v", buildProb)
	}
	if g.Accounting().Nodes != 1 {
		t.Fatalf("Accounting.Nodes = %d, want 1", g.Accounting().Nodes)
	}

	// Unknown version fails closed — no fallback to 1.0.0.
	_, prob = loaded.ResolveGraph(GraphRef{Ref: "g-linear", Version: "9.0.0"})
	assertViolation(t, prob, violationGraphInvalid, "/")

	_, prob = loaded.ResolveGraph(GraphRef{Ref: "missing", Version: "1.0.0"})
	assertViolation(t, prob, violationGraphInvalid, "/")

	_, prob = loaded.ResolveGraph(GraphRef{Ref: "", Version: "1.0.0"})
	assertViolation(t, prob, violationGraphInvalid, "/ref")

	_, prob = loaded.ResolveGraph(GraphRef{Ref: "g-linear", Version: ""})
	assertViolation(t, prob, violationGraphInvalid, "/version")
}

func TestTrustFailClosedStates(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		trust    *decision.TrustCarrier
		expected string
		code     apiproblem.ViolationCode
		field    string
	}{
		{
			name:  "unknown",
			trust: &decision.TrustCarrier{State: TrustStateUnknown},
			code:  violationTrustUnknown,
			field: "/spec/trust/state",
		},
		{
			name:  "expired",
			trust: &decision.TrustCarrier{State: TrustStateExpired},
			code:  violationTrustExpired,
			field: "/spec/trust/state",
		},
		{
			name:  "revoked",
			trust: &decision.TrustCarrier{State: TrustStateRevoked},
			code:  violationTrustRevoked,
			field: "/spec/trust/state",
		},
		{
			name:     "mismatch",
			trust:    &decision.TrustCarrier{State: "trusted-root-a"},
			expected: "trusted-root-b",
			code:     violationTrustMismatch,
			field:    "/spec/trust/state",
		},
		{
			name:  "present-empty",
			trust: &decision.TrustCarrier{},
			code:  violationTrustUnknown,
			field: "/spec/trust",
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			doc := Bundle{
				APIVersion: APIVersionDecisionProfileBundle,
				Kind:       KindDecisionProfileBundle,
				Spec: BundleSpec{
					Profiles:           []decision.DecisionProfile{minimalProfile("authz", "1.0.0")},
					Trust:              tc.trust,
					ExpectedTrustState: tc.expected,
				},
			}
			raw, err := json.Marshal(doc)
			if err != nil {
				t.Fatalf("marshal: %v", err)
			}
			_, prob := Load(raw, apivalid.ModeReadRepresentation)
			assertViolation(t, prob, tc.code, tc.field)
		})
	}
}

func TestTrustedBundleAccepted(t *testing.T) {
	t.Parallel()

	loaded := mustLoadJSON(t, Bundle{
		APIVersion: APIVersionDecisionProfileBundle,
		Kind:       KindDecisionProfileBundle,
		Spec: BundleSpec{
			Profiles: []decision.DecisionProfile{minimalProfile("authz", "1.0.0")},
			Trust: &decision.TrustCarrier{
				State:           "trusted-root-a",
				DigestAlgorithm: "opaque-unselected",
			},
			ExpectedTrustState: "trusted-root-a",
		},
	})
	if _, ok := loaded.View().Profile("authz", "1.0.0"); !ok {
		t.Fatal("expected profile present after trusted load")
	}
}

func assertViolation(t *testing.T, prob *apiproblem.Problem, code apiproblem.ViolationCode, field string) {
	t.Helper()
	if prob == nil {
		t.Fatalf("expected problem with code %s field %s", code, field)
	}
	if prob.Code != apiproblem.CodeValidationFailed {
		t.Fatalf("top-level code = %q, want VALIDATION_FAILED", prob.Code)
	}
	if len(prob.Violations) == 0 {
		t.Fatalf("expected violations, got %#v", prob)
	}
	v := prob.Violations[0]
	if v.Code != code {
		t.Fatalf("violation code = %q, want %q", v.Code, code)
	}
	if v.Field != field {
		t.Fatalf("violation field = %q, want %q", v.Field, field)
	}
}
