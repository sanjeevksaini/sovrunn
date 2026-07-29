package apiconform

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/apivalid"
	"github.com/sanjeevksaini/sovrunn/internal/decision"
	"github.com/sanjeevksaini/sovrunn/internal/decision/bundle"
	"github.com/sanjeevksaini/sovrunn/internal/decision/compose"
	"github.com/sanjeevksaini/sovrunn/internal/decision/validate"
)

// Expected FEATURE-0013 scenario ID count:
// 28 parents + 49 compound suffixes + 11 SCOPE + 8 EVAL + 8 SEC + 4 TRUST + 9 COMPAT = 117.
const expectedDecisionScenarioCount = 117

func TestRegisteredDecisionScenarioIDsResolve(t *testing.T) {
	t.Parallel()

	ids := RegisteredDecisionScenarioIDs()
	if len(ids) != expectedDecisionScenarioCount {
		t.Fatalf("RegisteredDecisionScenarioIDs len=%d, want %d", len(ids), expectedDecisionScenarioCount)
	}

	checks := RegisteredDecisionChecks()
	if len(checks) != len(ids) {
		t.Fatalf("RegisteredDecisionChecks len=%d, want %d", len(checks), len(ids))
	}

	seen := make(map[DecisionScenarioID]struct{}, len(ids))
	for i, id := range ids {
		if id == "" {
			t.Fatalf("empty scenario ID at index %d", i)
		}
		if _, dup := seen[id]; dup {
			t.Fatalf("duplicate scenario ID %q", id)
		}
		seen[id] = struct{}{}

		got, ok := ResolveDecisionCheck(id)
		if !ok {
			t.Fatalf("ResolveDecisionCheck(%q) = false", id)
		}
		if got.ID != id {
			t.Fatalf("ResolveDecisionCheck(%q).ID = %q", id, got.ID)
		}
		if got.Description == "" {
			t.Fatalf("ResolveDecisionCheck(%q) missing description", id)
		}
		if got.Kind == "" {
			t.Fatalf("ResolveDecisionCheck(%q) missing kind", id)
		}
		if checks[i].ID != id {
			t.Fatalf("registry order mismatch at %d: check=%q id=%q", i, checks[i].ID, id)
		}
	}

	required := []DecisionScenarioID{
		"F13-CF-01", "F13-CF-04.a", "F13-CF-04.i", "F13-CF-08.d", "F13-CF-11.e",
		"F13-CF-12.e", "F13-CF-16.f", "F13-CF-20.f", "F13-CF-21", "F13-CF-26.h",
		"F13-CF-28.f", "F13-SCOPE-01", "F13-SCOPE-11", "F13-EVAL-01", "F13-EVAL-08",
		"F13-SEC-01", "F13-SEC-08", "F13-TRUST-01", "F13-TRUST-04",
		"F13-COMPAT-01", "F13-COMPAT-09",
	}
	for _, id := range required {
		if _, ok := ResolveDecisionCheck(id); !ok {
			t.Fatalf("required scenario ID %q did not resolve", id)
		}
	}

	if _, ok := ResolveDecisionCheck("F13-CF-99"); ok {
		t.Fatal("unknown scenario ID must not resolve")
	}
}

func TestDecisionChecksExerciseBoundFixturesByKind(t *testing.T) {
	t.Parallel()

	root := moduleRoot(t)
	bundleJSON := mustDecisionCheckBundleJSON(t)
	positiveRecord := mustDecisionCheckRecordJSON(t, false)
	negativeRecord := mustDecisionCheckRecordJSON(t, true)
	auditJSON := mustReadConformanceFixture(t, root, "audit-event.json")

	strategy := compose.Strategy(func(in compose.StrategyInput) (compose.StrategyOutput, *apiproblem.Problem) {
		return compose.StrategyOutput{
			StrategyID:      in.Meta.ID,
			StrategyVersion: in.Meta.Version,
			InputRefs:       append([]string(nil), in.InputRefs...),
			TypedResult:     json.RawMessage(`{"selected":"pool-a"}`),
			Rationale:       decision.DecisionRationale{ReasonCodes: []string{"REPLAYED"}},
		}, nil
	})

	bound := map[DecisionScenarioID]DecisionCheckFixture{
		"F13-CF-01": {
			Kind:               DecisionCheckPositive,
			DecisionRecordJSON: positiveRecord,
			BundleJSON:         bundleJSON,
		},
		"F13-SCOPE-10": {
			Kind:               DecisionCheckNegative,
			DecisionRecordJSON: negativeRecord,
			BundleJSON:         bundleJSON,
			WantViolationCode:  validate.CodeScopeInvalid,
			WantViolationField: "/metadata/scopeRef/kind",
		},
		"F13-CF-17": {
			Kind: DecisionCheckSchema,
		},
		"F13-COMPAT-09": {
			Kind: DecisionCheckSchema,
		},
		"F13-CF-21": {
			Kind: DecisionCheckReplay,
			CapturedEvaluations: []decision.EvaluationResult{{
				Evaluator:            decision.EvaluatorIdentity{Type: "capacity", Version: "1.0.0"},
				ResultStatus:         decision.EvaluationResultStatusSuccess,
				Result:               json.RawMessage(`{"capacity":"ok"}`),
				EvaluatedAt:          "2026-07-29T10:00:01Z",
				StructuralTrustState: "TRUSTED",
			}},
			Strategy: strategy,
			StrategyInput: compose.StrategyInput{
				Meta:      compose.StrategyMetadata{ID: "all-of-v1", Version: "1.0.0"},
				InputRefs: []string{"input-1"},
			},
		},
		"F13-COMPAT-07": {
			Kind:           DecisionCheckCompat,
			AuditEventJSON: auditJSON,
		},
	}

	run := DecisionCheckRun{
		ModuleRoot:    root,
		BoundFixtures: bound,
		DecodeMode:    apivalid.ModeReadRepresentation,
	}

	for id := range bound {
		id := id
		check, ok := ResolveDecisionCheck(id)
		if !ok {
			t.Fatalf("ResolveDecisionCheck(%q) = false", id)
		}
		t.Run(string(id), func(t *testing.T) {
			t.Parallel()
			res := RunDecisionCheck(check, run)
			if !res.OK {
				t.Fatalf("RunDecisionCheck(%q) failed: detail=%q problem=%#v", id, res.Detail, res.Problem)
			}
		})
	}
}

func TestCheckStrictDecodeAndValidateHelpers(t *testing.T) {
	t.Parallel()

	bundleJSON := mustDecisionCheckBundleJSON(t)
	view, prob := LoadBundleViewStrictDecode(bundleJSON, apivalid.ModeReadRepresentation)
	if prob != nil {
		t.Fatalf("LoadBundleViewStrictDecode: %#v", prob)
	}

	recJSON := mustDecisionCheckRecordJSON(t, false)
	var decoded decision.DecisionRecord
	if prob := CheckStrictDecode(recJSON, apivalid.ModeReadRepresentation, &decoded); prob != nil {
		t.Fatalf("CheckStrictDecode: %#v", prob)
	}
	if decoded.Metadata.Name == "" {
		t.Fatal("expected decoded DecisionRecord name")
	}

	if _, prob := CheckDecisionRecordConformance(recJSON, apivalid.ModeReadRepresentation, view); prob != nil {
		t.Fatalf("CheckDecisionRecordConformance: %#v", prob)
	}

	negJSON := mustDecisionCheckRecordJSON(t, true)
	assertProb := CheckDecisionRecordNegative(
		negJSON,
		apivalid.ModeReadRepresentation,
		view,
		validate.CodeScopeInvalid,
		"/metadata/scopeRef/kind",
	)
	if assertProb != nil {
		t.Fatalf("CheckDecisionRecordNegative: %#v", assertProb)
	}
}

func TestCheckFeature0013TypeBindings(t *testing.T) {
	t.Parallel()

	findings := CheckFeature0013TypeBindings(moduleRoot(t))
	if len(findings) > 0 {
		t.Fatalf("CheckFeature0013TypeBindings findings: %v", findings)
	}
}

func TestCheckReplayKernelDeterministic(t *testing.T) {
	t.Parallel()

	captured := []decision.EvaluationResult{{
		Evaluator:            decision.EvaluatorIdentity{Type: "capacity", Version: "1.0.0"},
		ResultStatus:         decision.EvaluationResultStatusSuccess,
		Result:               json.RawMessage(`{"capacity":"ok"}`),
		EvaluatedAt:          "2026-07-29T10:00:01Z",
		StructuralTrustState: "TRUSTED",
	}}
	strategy := compose.Strategy(func(in compose.StrategyInput) (compose.StrategyOutput, *apiproblem.Problem) {
		return compose.StrategyOutput{
			StrategyID:      in.Meta.ID,
			StrategyVersion: in.Meta.Version,
			TypedResult:     append(json.RawMessage(nil), in.Evaluations[0].Result...),
		}, nil
	})
	in := compose.StrategyInput{
		Meta: compose.StrategyMetadata{ID: "replay-v1", Version: "1.0.0"},
	}

	out1, p1 := CheckReplayKernel(captured, strategy, in)
	out2, p2 := CheckReplayKernel(captured, strategy, in)
	if p1 != nil || p2 != nil {
		t.Fatalf("replay problems: %#v %#v", p1, p2)
	}
	if string(out1.TypedResult) != string(out2.TypedResult) {
		t.Fatalf("replay not deterministic: %s vs %s", out1.TypedResult, out2.TypedResult)
	}
}

func mustDecisionCheckBundleJSON(t *testing.T) []byte {
	t.Helper()
	doc := bundle.Bundle{
		APIVersion: bundle.APIVersionDecisionProfileBundle,
		Kind:       bundle.KindDecisionProfileBundle,
		Spec: bundle.BundleSpec{
			Profiles: []decision.DecisionProfile{decisionCheckProfile()},
			Obligations: []bundle.VersionedEntry{
				{ID: "audit-retain", Version: "1.0.0"},
			},
		},
	}
	raw, err := json.Marshal(doc)
	if err != nil {
		t.Fatalf("marshal bundle: %v", err)
	}
	return raw
}

func decisionCheckProfile() decision.DecisionProfile {
	return decision.DecisionProfile{
		TypeMeta: apimeta.TypeMeta{
			APIVersion: decision.APIVersionDecisionProfile,
			Kind:       decision.KindDecisionProfile,
		},
		Metadata: apimeta.ObjectMeta{
			Name: "placement-decision",
			UID:  "profile-uid-1",
			ScopeRef: &apimeta.ScopeRef{
				TypedRef: apimeta.TypedRef{
					APIVersion: "core.sovrunn.io/v1alpha1",
					Kind:       string(apimeta.ScopeProject),
					Name:       "payments-production",
					UID:        "project-123",
				},
			},
		},
		Spec: decision.ProfileSpec{
			Name:            "PlacementDecision",
			Family:          "placement",
			Version:         "1.0.0",
			Owner:           "platform-governance",
			LifecycleStatus: validate.LifecycleStatusActive,
			PrimaryForm:     decision.DecisionFormSelection,
			AllowedAuthorityLevels: []decision.DecisionAuthority{
				decision.DecisionAuthorityAuthoritative,
			},
			InputSchemaRef:       "schemas/placement-input@1",
			TypedResultSchemaRef: "schemas/placement-result@1",
			RationaleSchemaRef:   "schemas/placement-rationale@1",
			ObligationSchemaRef:  "schemas/placement-obligation@1",
			AcceptedEvaluationTypes: []decision.AcceptedEvaluationType{
				{Type: "capacity-evaluator", Version: "1.2.0"},
			},
			AcceptedCompositionStrategies: []string{"all-of-v1"},
			ScopeSemantics: decision.ProfileScopeSemantics{
				AllowedScopes: []apimeta.ScopeKind{apimeta.ScopeProject},
			},
			FailureBehavior: decision.FailureBehavior{
				FailPosture:          decision.FailPostureFailClosed,
				TimeoutBehavior:      "fail-closed",
				InsufficientEvidence: "fail-closed",
				DeclaredFailureModes: []string{"timeout", "evaluator-unavailable"},
			},
			SensitivityCeiling:  decision.SensitivityInternal,
			IdempotencyIdentity: "placement-v1",
			Limits: decision.ProfileLimits{
				MaxNodes:        16,
				MaxEdges:        32,
				MaxDepth:        8,
				MaxWidth:        8,
				MaxFanOut:       4,
				MaxPayloadBytes: 65536,
				MaxLatencyMs:    2000,
			},
		},
	}
}

func mustDecisionCheckRecordJSON(t *testing.T, invalidScope bool) []byte {
	t.Helper()
	scopeKind := string(apimeta.ScopeProject)
	if invalidScope {
		scopeKind = validate.KindServiceInstance
	}
	rec := decision.DecisionRecord{
		TypeMeta: apimeta.TypeMeta{
			APIVersion: decision.APIVersionDecisionRecord,
			Kind:       decision.KindDecisionRecord,
		},
		Metadata: apimeta.ObjectMeta{
			Name: "placement-decision-123",
			UID:  "decision-123",
			ScopeRef: &apimeta.ScopeRef{
				TypedRef: apimeta.TypedRef{
					APIVersion: "core.sovrunn.io/v1alpha1",
					Kind:       scopeKind,
					Name:       "payments-production",
					UID:        "project-123",
				},
			},
		},
		Record: decision.DecisionBody{
			ProfileRef: decision.ProfileRef{Name: "PlacementDecision", Version: "1.0.0"},
			Form:       decision.DecisionFormSelection,
			Authority:  decision.DecisionAuthorityAuthoritative,
			Purpose:    "workload-placement",
			SubjectRefs: []apimeta.TypedRef{{
				APIVersion: "platform.sovrunn.io/v1alpha1",
				Kind:       validate.KindServiceInstance,
				Name:       "postgres-primary",
				UID:        "service-456",
			}},
			Result: decision.DecisionResultBody{
				TypedResult: json.RawMessage(`{"selected":"pool-a"}`),
				Rationale: decision.DecisionRationale{
					ReasonCodes: []string{"CAPACITY_FIT"},
				},
				Obligations: []decision.Obligation{
					{ID: "audit-retain", Mandatory: true},
				},
			},
			SemanticIdentity: decision.SemanticDecisionIdentity{
				Basis:          "placement-v1",
				InputRefs:      []string{"input-1"},
				CanonicalScope: "Project:project-123",
			},
			Correlation: decision.Correlation{
				AuditEventRefs: []apimeta.TypedRef{{
					APIVersion: "governance.sovrunn.io/v1alpha1",
					Kind:       "AuditEvent",
					Name:       "audit-1",
					UID:        "audit-uid-1",
				}},
			},
			Finality: decision.FinalityFinal,
		},
	}
	raw, err := json.Marshal(rec)
	if err != nil {
		t.Fatalf("marshal record: %v", err)
	}
	return raw
}

func mustReadConformanceFixture(t *testing.T, moduleRoot, rel string) []byte {
	t.Helper()
	path := filepath.Join(moduleRoot, ConformanceFixturesDir, filepath.FromSlash(rel))
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read fixture %s: %v", path, err)
	}
	return raw
}
