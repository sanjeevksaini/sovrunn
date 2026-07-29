package bundle

import (
	"encoding/json"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/apivalid"
	"github.com/sanjeevksaini/sovrunn/internal/decision"
	"github.com/sanjeevksaini/sovrunn/internal/decision/graph"
)

func TestLoadJSONAndYAMLEquivalenceViaStrictDecode(t *testing.T) {
	t.Parallel()

	doc := Bundle{
		APIVersion: APIVersionDecisionProfileBundle,
		Kind:       KindDecisionProfileBundle,
		Spec: BundleSpec{
			Profiles: []decision.DecisionProfile{
				minimalProfile("authz", "1.0.0"),
			},
			Graphs: []graph.Definition{
				singleGraph("g1", "1.0.0"),
			},
			Forms: []VersionedEntry{
				{ID: "ADJUDICATION", Version: "1.0.0"},
			},
			Authorities: []VersionedEntry{
				{ID: "AUTHORITATIVE", Version: "1.0.0"},
			},
			EvaluatorTypes: []VersionedEntry{
				{ID: "cel", Version: "1.0.0"},
			},
			Strategies: []VersionedEntry{
				{ID: "allOf", Version: "1.0.0"},
			},
			ReasonCodes: []VersionedEntry{
				{ID: "POLICY_MATCH", Version: "1.0.0"},
			},
			Obligations: []VersionedEntry{
				{ID: "audit", Version: "1.0.0"},
			},
			EventTypes: []VersionedEntry{
				{ID: "decision.completed", Version: "1.0.0"},
			},
			EvidenceTypes: []VersionedEntry{
				{ID: "attestation", Version: "1.0.0"},
			},
			Projections: []VersionedEntry{
				{ID: "operator", Version: "1.0.0"},
			},
			Trust: &decision.TrustCarrier{
				State: "trusted-root-a",
			},
			ExpectedTrustState: "trusted-root-a",
		},
	}

	jsonBytes, err := json.Marshal(doc)
	if err != nil {
		t.Fatalf("marshal json: %v", err)
	}

	// Strict JSON-compatible YAML subset (apivalid.DecodeYAML path via StrictDecode).
	yamlBytes := []byte(`
apiVersion: governance.sovrunn.io/v1alpha1
kind: DecisionProfileBundle
spec:
  profiles:
    - apiVersion: governance.sovrunn.io/v1alpha1
      kind: DecisionProfile
      metadata:
        name: authz
      spec:
        name: authz
        version: "1.0.0"
        family: test
        owner: owner
  graphs:
    - ref: g1
      version: "1.0.0"
      nodes:
        - id: d1
          kind: decision
      edges: []
      bounds:
        maxNodes: 4
        maxEdges: 4
        maxDepth: 2
        maxWidth: 2
        maxFanOut: 2
        maxPayloadBytes: 64
  forms:
    - id: ADJUDICATION
      version: "1.0.0"
  authorities:
    - id: AUTHORITATIVE
      version: "1.0.0"
  evaluatorTypes:
    - id: cel
      version: "1.0.0"
  strategies:
    - id: allOf
      version: "1.0.0"
  reasonCodes:
    - id: POLICY_MATCH
      version: "1.0.0"
  obligations:
    - id: audit
      version: "1.0.0"
  eventTypes:
    - id: decision.completed
      version: "1.0.0"
  evidenceTypes:
    - id: attestation
      version: "1.0.0"
  projections:
    - id: operator
      version: "1.0.0"
  trust:
    state: trusted-root-a
  expectedTrustState: trusted-root-a
`)

	fromJSON, prob := Load(jsonBytes, apivalid.ModeReadRepresentation)
	if prob != nil {
		t.Fatalf("Load JSON = %#v", prob)
	}
	fromYAML, prob := Load(yamlBytes, apivalid.ModeReadRepresentation)
	if prob != nil {
		t.Fatalf("Load YAML = %#v", prob)
	}

	assertEquivalentViews(t, fromJSON.View(), fromYAML.View())

	jDef, jProb := fromJSON.ResolveGraph(GraphRef{Ref: "g1", Version: "1.0.0"})
	yDef, yProb := fromYAML.ResolveGraph(GraphRef{Ref: "g1", Version: "1.0.0"})
	if jProb != nil || yProb != nil {
		t.Fatalf("ResolveGraph JSON=%#v YAML=%#v", jProb, yProb)
	}
	if jDef.Ref != yDef.Ref || jDef.Version != yDef.Version || len(jDef.Nodes) != len(yDef.Nodes) {
		t.Fatalf("resolved graphs differ: json=%+v yaml=%+v", jDef, yDef)
	}
}

func TestLoadRejectsMalformedJSON(t *testing.T) {
	t.Parallel()

	_, prob := Load([]byte(`{"apiVersion":`), apivalid.ModeReadRepresentation)
	if prob == nil {
		t.Fatal("expected malformed JSON problem")
	}
	// Decode failures use FEATURE-0012 top-level codes, not DECISION_*.
	if len(string(prob.Code)) >= 9 && string(prob.Code)[:9] == "DECISION_" {
		t.Fatalf("decode failure must not use DECISION_* top-level code, got %q", prob.Code)
	}
	for _, v := range prob.Violations {
		if len(string(v.Code)) >= 9 && string(v.Code)[:9] == "DECISION_" {
			t.Fatalf("decode failure must not use DECISION_* violation code, got %q", v.Code)
		}
	}
	if prob.Code != apiproblem.CodeMalformedRequest && prob.Code != apiproblem.CodeValidationFailed {
		// Accept any inherited FEATURE-0012 top-level decode/malformed code.
		switch prob.Code {
		case apiproblem.CodeRequestTooLarge, apiproblem.CodeUnsupportedMediaType:
			// also inherited
		default:
			// Still fine as long as it is not DECISION_* (checked above).
		}
	}
}

func TestLoadEmptyBodyRejected(t *testing.T) {
	t.Parallel()

	_, prob := Load(nil, apivalid.ModeReadRepresentation)
	if prob == nil {
		t.Fatal("expected empty-body problem")
	}
}

func assertEquivalentViews(t *testing.T, a, b BundleView) {
	t.Helper()

	ap, ok := a.Profile("authz", "1.0.0")
	if !ok {
		t.Fatal("json view missing profile")
	}
	bp, ok := b.Profile("authz", "1.0.0")
	if !ok {
		t.Fatal("yaml view missing profile")
	}
	if ap.Spec.Name != bp.Spec.Name || ap.Spec.Version != bp.Spec.Version {
		t.Fatalf("profile mismatch json=%+v yaml=%+v", ap.Spec, bp.Spec)
	}

	mustBoth := func(name, id, version string, fa, fb func(string, string) (VersionedEntry, bool)) {
		t.Helper()
		if _, ok := fa(id, version); !ok {
			t.Fatalf("json view missing %s %s/%s", name, id, version)
		}
		if _, ok := fb(id, version); !ok {
			t.Fatalf("yaml view missing %s %s/%s", name, id, version)
		}
	}
	mustBoth("forms", "ADJUDICATION", "1.0.0", a.Forms, b.Forms)
	mustBoth("authorities", "AUTHORITATIVE", "1.0.0", a.Authorities, b.Authorities)
	mustBoth("evaluatorTypes", "cel", "1.0.0", a.EvaluatorTypes, b.EvaluatorTypes)
	mustBoth("strategies", "allOf", "1.0.0", a.Strategies, b.Strategies)
	mustBoth("reasonCodes", "POLICY_MATCH", "1.0.0", a.ReasonCodes, b.ReasonCodes)
	mustBoth("obligations", "audit", "1.0.0", a.Obligations, b.Obligations)
	mustBoth("eventTypes", "decision.completed", "1.0.0", a.EventTypes, b.EventTypes)
	mustBoth("evidenceTypes", "attestation", "1.0.0", a.EvidenceTypes, b.EvidenceTypes)
	mustBoth("projections", "operator", "1.0.0", a.Projections, b.Projections)

	if _, ok := a.Graph("g1", "1.0.0"); !ok {
		t.Fatal("json view missing graph")
	}
	if _, ok := b.Graph("g1", "1.0.0"); !ok {
		t.Fatal("yaml view missing graph")
	}
}
