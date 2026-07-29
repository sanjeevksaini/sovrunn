//go:build ignore

// gen_decision_negative_fixtures writes T-030 negative DecisionNegativeFixture
// envelopes under tests/conformance/fixtures/negative/decision/. Each fixture is
// verified against validate.DecodeAndValidateDecisionRecordWith before write.
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/apivalid"
	"github.com/sanjeevksaini/sovrunn/internal/decision"
	"github.com/sanjeevksaini/sovrunn/internal/decision/bundle"
	"github.com/sanjeevksaini/sovrunn/internal/decision/validate"
)

const kindDecisionNegativeFixture = "DecisionNegativeFixture"

type fixtureEnvelope struct {
	Kind               string          `json:"kind"`
	WantViolationCode  string          `json:"wantViolationCode"`
	WantViolationField string          `json:"wantViolationField"`
	Options            json.RawMessage `json:"options,omitempty"`
	Bundle             json.RawMessage `json:"bundle,omitempty"`
	DecisionRecord     json.RawMessage `json:"decisionRecord"`
}

type scenario struct {
	id        string
	wantCode  apiproblem.ViolationCode
	wantField string
	build     func(ctx *genCtx) (record, bundleMap map[string]any, opts map[string]any)
}

type genCtx struct {
	shared map[string]any
}

func main() {
	root := moduleRoot()
	outDir := filepath.Join(root, "tests/conformance/fixtures/negative/decision")
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		fatal(err)
	}

	shared, err := loadJSON(filepath.Join(root, "tests/conformance/fixtures/decision/positive/_shared-bundle.json"))
	if err != nil {
		fatal(err)
	}
	ctx := &genCtx{shared: shared}

	for _, sc := range allScenarios() {
		record, bundleMap, opts := sc.build(ctx)
		if err := verify(record, bundleMap, opts, sc.wantCode, sc.wantField); err != nil {
			fatal(fmt.Errorf("%s: %w", sc.id, err))
		}
		env, err := marshalEnvelope(record, bundleMap, opts, sc.wantCode, sc.wantField)
		if err != nil {
			fatal(fmt.Errorf("%s marshal: %w", sc.id, err))
		}
		path := filepath.Join(outDir, sc.id+".json")
		if err := os.WriteFile(path, env, 0o644); err != nil {
			fatal(fmt.Errorf("%s write: %w", sc.id, err))
		}
		fmt.Println("wrote", path)
	}
}

func allScenarios() []scenario {
	// Parent IDs alias the representative suffix (architecture §17).
	cf08a := buildCF08a
	cf11d := buildCF11d
	cf12a := buildCF12a
	cf16a := buildCF16a
	cf26c := buildCF26c

	return []scenario{
		{id: "F13-CF-08", wantCode: validate.CodeEvaluationLimitExceeded, wantField: "/record/evaluationResults/0/timing/durationMs", build: cf08a},
		{id: "F13-CF-08.a", wantCode: validate.CodeEvaluationLimitExceeded, wantField: "/record/evaluationResults/0/timing/durationMs", build: cf08a},
		{id: "F13-CF-08.b", wantCode: validate.CodeTrustRequired, wantField: "/record/trust", build: buildCF08b},
		{id: "F13-CF-08.c", wantCode: validate.CodeEvaluationTypeUnsupported, wantField: "/record/evaluationResults/0/evaluator/type", build: buildCF08c},
		{id: "F13-CF-08.d", wantCode: validate.CodeCompositionInputConflict, wantField: "/nodes/1", build: buildCF08d},

		{id: "F13-CF-11", wantCode: validate.CodeRelationshipCycle, wantField: "/record/relationship/targetRef", build: cf11d},
		{id: "F13-CF-11.a", wantCode: validate.CodeRelationshipConflict, wantField: "/record/relationship/targetRef", build: buildCF11a},
		{id: "F13-CF-11.b", wantCode: validate.CodeRelationshipConflict, wantField: "/record/relationship/targetRef", build: buildCF11b},
		{id: "F13-CF-11.c", wantCode: validate.CodeRelationshipConflict, wantField: "/record/relationship/targetRef", build: buildCF11c},
		{id: "F13-CF-11.d", wantCode: validate.CodeRelationshipCycle, wantField: "/record/relationship/targetRef", build: cf11d},
		{id: "F13-CF-11.e", wantCode: validate.CodeRelationshipChainLimitExceeded, wantField: "/record/relationship", build: buildCF11e},

		{id: "F13-CF-12", wantCode: validate.CodeProfileUnknown, wantField: "/record/profileRef/name", build: cf12a},
		{id: "F13-CF-12.a", wantCode: validate.CodeProfileUnknown, wantField: "/record/profileRef/name", build: cf12a},
		{id: "F13-CF-12.b", wantCode: validate.CodeProfileSchemaInvalid, wantField: "/spec/failureBehavior/securityExceptionRef", build: buildCF12b},
		{id: "F13-CF-12.c", wantCode: validate.CodeProfileSchemaInvalid, wantField: "/spec/projectionRules/0/includes/0", build: buildCF12c},
		{id: "F13-CF-12.d", wantCode: validate.CodeObligationUnsupportedMandatory, wantField: "/record/result/obligations/0/id", build: buildCF12d},
		{id: "F13-CF-12.e", wantCode: validate.CodeScopeMismatch, wantField: "/metadata/scopeRef", build: buildCF12e},

		{id: "F13-CF-16", wantCode: validate.CodeCompositionLimitExceeded, wantField: "/bounds/maxNodes", build: cf16a},
		{id: "F13-CF-16.a", wantCode: validate.CodeCompositionLimitExceeded, wantField: "/bounds/maxNodes", build: cf16a},
		{id: "F13-CF-16.b", wantCode: validate.CodeCompositionLimitExceeded, wantField: "/bounds/maxDepth", build: buildCF16b},
		{id: "F13-CF-16.c", wantCode: validate.CodeCompositionLimitExceeded, wantField: "/bounds/maxWidth", build: buildCF16c},
		{id: "F13-CF-16.d", wantCode: validate.CodeCompositionLimitExceeded, wantField: "/bounds/maxFanOut", build: buildCF16d},
		{id: "F13-CF-16.e", wantCode: validate.CodeEvaluationLimitExceeded, wantField: "/record/evaluationResults/0/timing/durationMs", build: buildCF16e},
		{id: "F13-CF-16.f", wantCode: validate.CodeEvaluationLimitExceeded, wantField: "/record/evaluationResults/0/result", build: buildCF16f},

		{id: "F13-CF-23", wantCode: validate.CodeProfileUnknown, wantField: "/record/profileRef/name", build: buildCF23},
		{id: "F13-CF-24", wantCode: validate.CodeObligationUnsupportedMandatory, wantField: "/record/result/obligations/0/id", build: buildCF24},
		{id: "F13-CF-25", wantCode: validate.CodeTrustExpired, wantField: "/record/trust/state", build: buildCF25},

		{id: "F13-CF-26", wantCode: validate.CodeProfileLimitInvalid, wantField: "/spec/limits/maxPayloadBytes", build: cf26c},
		{id: "F13-CF-26.a", wantCode: validate.CodeEvaluationResultInvalid, wantField: "/record/evaluationResults/0/evaluator/type", build: buildCF26a},
		{id: "F13-CF-26.b", wantCode: validate.CodeProfileLimitInvalid, wantField: "/spec/limits/maxCalls", build: buildCF26b},
		{id: "F13-CF-26.c", wantCode: validate.CodeProfileLimitInvalid, wantField: "/spec/limits/maxPayloadBytes", build: cf26c},
		{id: "F13-CF-26.d", wantCode: validate.CodeProfileLimitInvalid, wantField: "/spec/limits/maxConcurrency", build: buildCF26d},
		{id: "F13-CF-26.e", wantCode: validate.CodeProfileLimitInvalid, wantField: "/spec/limits/maxRetries", build: buildCF26e},
		{id: "F13-CF-26.f", wantCode: validate.CodeProfileLimitInvalid, wantField: "/spec/limits/maxJurisdictions", build: buildCF26f},
		{id: "F13-CF-26.g", wantCode: validate.CodeEvaluationLimitExceeded, wantField: "/record/evaluationResults/0/timing/durationMs", build: buildCF26g},
		{id: "F13-CF-26.h", wantCode: validate.CodeProfileLimitInvalid, wantField: "/spec/acceptedEvaluationTypes/1/limits/maxCostBudget", build: buildCF26h},

		{id: "F13-SCOPE-07", wantCode: validate.CodeScopeRequired, wantField: "/metadata/scopeRef", build: buildScope07},
		{id: "F13-SCOPE-09", wantCode: validate.CodeScopeConflict, wantField: "/scopeRef", build: buildScope09},
		{id: "F13-SCOPE-10", wantCode: validate.CodeScopeInvalid, wantField: "/metadata/scopeRef/kind", build: buildScope10},

		{id: "F13-EVAL-03", wantCode: validate.CodeProfileLimitInvalid, wantField: "/spec/acceptedEvaluationTypes/1/limits/maxPayloadBytes", build: buildEval03},
		{id: "F13-EVAL-04", wantCode: validate.CodeProfileLimitInvalid, wantField: "/spec/limits/maxPayloadBytes", build: buildEval04},
		{id: "F13-EVAL-05", wantCode: validate.CodeEvaluationResultInvalid, wantField: "/record/evaluationResults/0/resultStatus", build: buildEval05},
		{id: "F13-EVAL-06", wantCode: validate.CodeEvaluationScopeMismatch, wantField: "/record/evaluationResults/0/scopeEvidence", build: buildEval06},
		{id: "F13-EVAL-07", wantCode: validate.CodeTrustMismatch, wantField: "/record/trust/state", build: buildEval07},

		{id: "F13-SEC-01", wantCode: validate.CodeEvaluationLimitExceeded, wantField: "/declaredSensitivity", build: buildSEC01},
		{id: "F13-SEC-02", wantCode: validate.CodeProfileLimitInvalid, wantField: "/spec/limits/maxFieldCount", build: buildSEC02},
		{id: "F13-SEC-03", wantCode: validate.CodeProfileLimitInvalid, wantField: "/spec/limits/maxPayloadBytes", build: buildSEC03},
		{id: "F13-SEC-04", wantCode: validate.CodeProfileLimitInvalid, wantField: "/spec/limits/maxDepth", build: buildSEC04},
		{id: "F13-SEC-05", wantCode: validate.CodeEvaluationResultInvalid, wantField: "/contentCategories/0", build: buildSEC05},
		{id: "F13-SEC-06", wantCode: validate.CodeEvaluationResultInvalid, wantField: "/contentCategories/0", build: buildSEC06},
		{id: "F13-SEC-07", wantCode: validate.CodeProfileSchemaInvalid, wantField: "/declaredSensitivity", build: buildSEC07},
		{id: "F13-SEC-08", wantCode: validate.CodeProfileSchemaInvalid, wantField: "/classifications/0", build: buildSEC08},

		{id: "F13-TRUST-03", wantCode: validate.CodeTrustRevoked, wantField: "/record/trust/state", build: buildTrust03},
	}
}

// --- shared helpers ---

func moduleRoot() string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		panic("runtime.Caller failed")
	}
	return filepath.Clean(filepath.Join(filepath.Dir(file), "../.."))
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

func loadJSON(path string) (map[string]any, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var out map[string]any
	if err := json.Unmarshal(raw, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func cloneMap(in map[string]any) map[string]any {
	raw, _ := json.Marshal(in)
	var out map[string]any
	_ = json.Unmarshal(raw, &out)
	return out
}

func profileSpec(bundleMap map[string]any) map[string]any {
	spec := bundleMap["spec"].(map[string]any)
	profiles := spec["profiles"].([]any)
	prof := profiles[0].(map[string]any)
	return prof["spec"].(map[string]any)
}

func profileLimits(bundleMap map[string]any) map[string]any {
	return profileSpec(bundleMap)["limits"].(map[string]any)
}

func baseRecord(id string) map[string]any {
	slug := strings.ToLower(strings.ReplaceAll(id, ".", "-"))
	return map[string]any{
		"apiVersion": "governance.sovrunn.io/v1alpha1",
		"kind":       "DecisionRecord",
		"metadata": map[string]any{
			"name": slug,
			"scopeRef": map[string]any{
				"apiVersion": "core.sovrunn.io/v1alpha1",
				"kind":       "Project",
				"name":       "payments-production",
				"uid":        "project-123",
			},
			"uid": "decision-uid-" + slug,
		},
		"record": map[string]any{
			"authority": "AUTHORITATIVE",
			"correlation": map[string]any{
				"auditEventRefs": []any{
					map[string]any{
						"apiVersion": "governance.sovrunn.io/v1alpha1",
						"kind":       "AuditEvent",
						"name":       "audit-" + slug,
						"uid":        "audit-uid-" + slug,
					},
				},
			},
			"finality": "FINAL",
			"form":     "SELECTION",
			"profileRef": map[string]any{
				"name":    "FixtureLocalDecision",
				"version": "1.0.0",
			},
			"purpose": "conformance-negative-" + slug,
			"result": map[string]any{
				"obligations": []any{},
				"rationale": map[string]any{
					"reasonCodes": []any{"CAPACITY_FIT"},
				},
				"typedResult": map[string]any{
					"selected": "pool-a",
				},
			},
			"semanticIdentity": map[string]any{
				"basis":          "fixture-local-v1",
				"canonicalScope": "Project:project-123",
				"inputRefs":      []any{"input-1"},
			},
		},
	}
}

func baseEvalResult() map[string]any {
	return map[string]any{
		"evaluatedAt": "2026-07-29T10:00:01Z",
		"evaluator": map[string]any{
			"type":    "capacity-evaluator",
			"version": "1.2.0",
		},
		"inputSnapshotRef": map[string]any{
			"apiVersion": "evidence.sovrunn.io/v1alpha1",
			"kind":       "InputSnapshot",
			"name":       "snapshot-1",
			"uid":        "snapshot-uid-1",
		},
		"result": map[string]any{
			"capacity": "ok",
		},
		"resultStatus": "SUCCESS",
		"timing": map[string]any{
			"completedAt": "2026-07-29T10:00:01Z",
			"durationMs":  1000,
			"startedAt":   "2026-07-29T10:00:00Z",
		},
	}
}

func withEval(record map[string]any, eval map[string]any) {
	rec := record["record"].(map[string]any)
	rec["evaluationResults"] = []any{eval}
}

func withComposition(record map[string]any) {
	rec := record["record"].(map[string]any)
	rec["composition"] = map[string]any{
		"strategyName":    "all-of-v1",
		"strategyVersion": "1.0.0",
		"graphRef":        "fixture-graph",
		"graphVersion":    "1.0.0",
		"inputRefs":       []any{"input-1"},
	}
}

func tightCompositionBundle(ctx *genCtx) map[string]any {
	b := cloneMap(ctx.shared)
	lim := profileLimits(b)
	lim["maxNodes"] = 16
	lim["maxEdges"] = 32
	lim["maxDepth"] = 8
	lim["maxWidth"] = 8
	lim["maxFanOut"] = 4
	lim["maxPayloadBytes"] = 65536
	lim["maxLatencyMs"] = 2000
	return b
}

func decisionRef(name, uid string) map[string]any {
	return map[string]any{
		"apiVersion": "governance.sovrunn.io/v1alpha1",
		"kind":       "DecisionRecord",
		"name":       name,
		"uid":        uid,
	}
}

func establishedNode(name, uid string, kind string, target map[string]any) map[string]any {
	node := map[string]any{"ref": decisionRef(name, uid)}
	if kind != "" {
		node["relationship"] = map[string]any{
			"kind":      kind,
			"targetRef": target,
		}
	}
	return node
}

func canonicalView() []any {
	return []any{
		"/record/result/typedResult",
		"/record/result/rationale",
		"/record/trust",
		"/record/evaluationResults",
	}
}

// --- scenario builders ---

func buildCF08a(ctx *genCtx) (map[string]any, map[string]any, map[string]any) {
	record := baseRecord("F13-CF-08.a")
	eval := baseEvalResult()
	eval["timing"].(map[string]any)["durationMs"] = 99999
	withEval(record, eval)
	return record, cloneMap(ctx.shared), nil
}

func buildCF08b(ctx *genCtx) (map[string]any, map[string]any, map[string]any) {
	record := baseRecord("F13-CF-08.b")
	withEval(record, baseEvalResult())
	return record, cloneMap(ctx.shared), map[string]any{"trustRequired": true}
}

func buildCF08c(ctx *genCtx) (map[string]any, map[string]any, map[string]any) {
	record := baseRecord("F13-CF-08.c")
	eval := baseEvalResult()
	eval["evaluator"] = map[string]any{"type": "unknown-evaluator", "version": "9.9.9"}
	withEval(record, eval)
	return record, cloneMap(ctx.shared), nil
}

func buildCF08d(ctx *genCtx) (map[string]any, map[string]any, map[string]any) {
	record := baseRecord("F13-CF-08.d")
	withComposition(record)
	rec := record["record"].(map[string]any)
	comp := rec["composition"].(map[string]any)
	comp["inputRefs"] = []any{"eval-1"}
	withEval(record, baseEvalResult())
	return record, cloneMap(ctx.shared), nil
}

func buildCF11a(ctx *genCtx) (map[string]any, map[string]any, map[string]any) {
	record := baseRecord("F13-CF-11.a")
	root := decisionRef("root", "uid-root")
	rec := record["record"].(map[string]any)
	rec["relationship"] = map[string]any{
		"kind":      "SUPERSEDES",
		"targetRef": root,
	}
	opts := map[string]any{
		"established": []any{
			establishedNode("root", "uid-root", "", nil),
			establishedNode("rev", "uid-rev", "REVOKES", root),
		},
	}
	return record, cloneMap(ctx.shared), opts
}

func buildCF11b(ctx *genCtx) (map[string]any, map[string]any, map[string]any) {
	record := baseRecord("F13-CF-11.b")
	root := decisionRef("root", "uid-root")
	rec := record["record"].(map[string]any)
	rec["relationship"] = map[string]any{
		"kind":      "CORRECTS",
		"targetRef": root,
	}
	opts := map[string]any{
		"established": []any{
			establishedNode("root", "uid-root", "", nil),
			establishedNode("rev", "uid-rev", "REVOKES", root),
		},
	}
	return record, cloneMap(ctx.shared), opts
}

func buildCF11c(ctx *genCtx) (map[string]any, map[string]any, map[string]any) {
	record := baseRecord("F13-CF-11.c")
	root := decisionRef("root", "uid-root")
	rec := record["record"].(map[string]any)
	rec["relationship"] = map[string]any{
		"kind":      "REVOKES",
		"targetRef": root,
	}
	opts := map[string]any{
		"established": []any{
			establishedNode("root", "uid-root", "", nil),
			establishedNode("rev", "uid-rev", "REVOKES", root),
		},
	}
	return record, cloneMap(ctx.shared), opts
}

func buildCF11d(ctx *genCtx) (map[string]any, map[string]any, map[string]any) {
	record := baseRecord("F13-CF-11.d")
	meta := record["metadata"].(map[string]any)
	self := decisionRef(meta["name"].(string), meta["uid"].(string))
	a := decisionRef("a", "uid-a")
	rec := record["record"].(map[string]any)
	rec["relationship"] = map[string]any{
		"kind":      "CORRECTS",
		"targetRef": a,
	}
	opts := map[string]any{
		"established": []any{
			establishedNode("a", "uid-a", "SUPERSEDES", self),
			establishedNode("b", "uid-b", "SUPERSEDES", a),
		},
	}
	return record, cloneMap(ctx.shared), opts
}

func buildCF11e(ctx *genCtx) (map[string]any, map[string]any, map[string]any) {
	b := cloneMap(ctx.shared)
	ps := profileSpec(b)
	vr := map[string]any{"maxChainDepth": 1}
	ps["validityRules"] = vr

	record := baseRecord("F13-CF-11.e")
	r0 := decisionRef("r0", "uid-0")
	r1 := decisionRef("r1", "uid-1")
	r2 := decisionRef("r2", "uid-2")
	rec := record["record"].(map[string]any)
	rec["relationship"] = map[string]any{
		"kind":      "SUPERSEDES",
		"targetRef": r2,
	}
	opts := map[string]any{
		"established": []any{
			establishedNode("r0", "uid-0", "", nil),
			establishedNode("r1", "uid-1", "SUPERSEDES", r0),
			establishedNode("r2", "uid-2", "SUPERSEDES", r1),
		},
	}
	return record, b, opts
}

func buildCF12a(ctx *genCtx) (map[string]any, map[string]any, map[string]any) {
	record := baseRecord("F13-CF-12.a")
	rec := record["record"].(map[string]any)
	rec["profileRef"] = map[string]any{"name": "UnknownProfile", "version": "9.9.9"}
	return record, cloneMap(ctx.shared), nil
}

func buildCF12b(ctx *genCtx) (map[string]any, map[string]any, map[string]any) {
	b := cloneMap(ctx.shared)
	ps := profileSpec(b)
	fb := ps["failureBehavior"].(map[string]any)
	fb["failPosture"] = "FAIL_OPEN"
	delete(fb, "securityExceptionRef")
	record := baseRecord("F13-CF-12.b")
	return record, b, nil
}

func buildCF12c(ctx *genCtx) (map[string]any, map[string]any, map[string]any) {
	b := cloneMap(ctx.shared)
	ps := profileSpec(b)
	ps["projectionRules"] = []any{
		map[string]any{
			"audience": "operator",
			"includes": []any{"/record/result/typedResult"},
			"excludes": []any{"/record/result/typedResult"},
		},
	}
	record := baseRecord("F13-CF-12.c")
	return record, b, map[string]any{"canonicalView": canonicalView()}
}

func buildCF12d(ctx *genCtx) (map[string]any, map[string]any, map[string]any) {
	record := baseRecord("F13-CF-12.d")
	rec := record["record"].(map[string]any)
	res := rec["result"].(map[string]any)
	res["obligations"] = []any{map[string]any{"id": "audit-retain", "mandatory": true}}
	opts := map[string]any{"supportedObligations": map[string]any{"audit-retain": false}}
	return record, cloneMap(ctx.shared), opts
}

func buildCF12e(ctx *genCtx) (map[string]any, map[string]any, map[string]any) {
	record := baseRecord("F13-CF-12.e")
	meta := record["metadata"].(map[string]any)
	meta["scopeRef"] = map[string]any{
		"apiVersion": "core.sovrunn.io/v1alpha1",
		"kind":       "Organization",
		"name":       "acme-corp",
		"uid":        "org-456",
	}
	rec := record["record"].(map[string]any)
	rec["subjectRefs"] = []any{
		map[string]any{
			"apiVersion": "platform.sovrunn.io/v1alpha1",
			"kind":       "ServiceInstance",
			"name":       "postgres-primary",
			"uid":        "service-456",
		},
	}
	sem := rec["semanticIdentity"].(map[string]any)
	sem["canonicalScope"] = "Organization:org-456"
	return record, cloneMap(ctx.shared), nil
}

func cf16Bundle(ctx *genCtx, tweak func(map[string]any)) map[string]any {
	b := tightCompositionBundle(ctx)
	lim := profileLimits(b)
	tweak(lim)
	// Keep accepted evaluator registrations narrowed vs profile (pass 3).
	ps := profileSpec(b)
	types := ps["acceptedEvaluationTypes"].([]any)
	for _, raw := range types {
		entry := raw.(map[string]any)
		if entry["limits"] == nil {
			continue
		}
		el := entry["limits"].(map[string]any)
		for k, v := range lim {
			pv, ok := numericLimit(v)
			if !ok || pv <= 0 {
				continue
			}
			if ev, ok := el[k]; ok {
				if evn, ok := numericLimit(ev); ok && evn <= pv {
					continue
				}
			}
			el[k] = pv
		}
	}
	return b
}

func numericLimit(v any) (float64, bool) {
	switch n := v.(type) {
	case float64:
		return n, true
	case int:
		return float64(n), true
	case int64:
		return float64(n), true
	default:
		return 0, false
	}
}

func buildCF16a(ctx *genCtx) (map[string]any, map[string]any, map[string]any) {
	record := baseRecord("F13-CF-16.a")
	withComposition(record)
	b := cf16Bundle(ctx, func(lim map[string]any) { lim["maxNodes"] = 2 })
	return record, b, nil
}

func buildCF16b(ctx *genCtx) (map[string]any, map[string]any, map[string]any) {
	record := baseRecord("F13-CF-16.b")
	withComposition(record)
	b := cf16Bundle(ctx, func(lim map[string]any) { lim["maxDepth"] = 1 })
	return record, b, nil
}

func buildCF16c(ctx *genCtx) (map[string]any, map[string]any, map[string]any) {
	record := baseRecord("F13-CF-16.c")
	withComposition(record)
	b := cf16Bundle(ctx, func(lim map[string]any) { lim["maxWidth"] = 1 })
	return record, b, nil
}

func buildCF16d(ctx *genCtx) (map[string]any, map[string]any, map[string]any) {
	record := baseRecord("F13-CF-16.d")
	withComposition(record)
	b := cf16Bundle(ctx, func(lim map[string]any) { lim["maxFanOut"] = 1 })
	return record, b, nil
}

func buildCF16e(ctx *genCtx) (map[string]any, map[string]any, map[string]any) {
	record := baseRecord("F13-CF-16.e")
	withComposition(record)
	eval := baseEvalResult()
	eval["timing"].(map[string]any)["durationMs"] = 99999
	withEval(record, eval)
	return record, tightCompositionBundle(ctx), nil
}

func buildCF16f(ctx *genCtx) (map[string]any, map[string]any, map[string]any) {
	record := baseRecord("F13-CF-16.f")
	withComposition(record)
	eval := baseEvalResult()
	eval["result"] = map[string]any{"payload": strings.Repeat("x", 256)}
	withEval(record, eval)
	b := cf16Bundle(ctx, func(lim map[string]any) {
		lim["maxPayloadBytes"] = 16
		lim["maxLatencyMs"] = 2000
	})
	return record, b, nil
}

func buildCF23(ctx *genCtx) (map[string]any, map[string]any, map[string]any) {
	record := baseRecord("F13-CF-23")
	rec := record["record"].(map[string]any)
	rec["profileRef"] = map[string]any{"name": "ExpiredProfile", "version": "0.0.1"}
	return record, cloneMap(ctx.shared), nil
}

func buildCF24(ctx *genCtx) (map[string]any, map[string]any, map[string]any) {
	return buildCF12d(ctx)
}

func buildCF25(ctx *genCtx) (map[string]any, map[string]any, map[string]any) {
	record := baseRecord("F13-CF-25")
	rec := record["record"].(map[string]any)
	rec["trust"] = map[string]any{
		"canonicalizationMethod": "fixture-canon-v1",
		"digestAlgorithm":        "fixture-digest-v1",
		"signatureAlgorithm":     "fixture-sig-v1",
		"state":                  "expired",
	}
	return record, cloneMap(ctx.shared), map[string]any{"trustRequired": true}
}

func buildCF26a(ctx *genCtx) (map[string]any, map[string]any, map[string]any) {
	record := baseRecord("F13-CF-26.a")
	eval := baseEvalResult()
	eval["evaluator"] = map[string]any{"type": "", "version": "1.2.0"}
	withEval(record, eval)
	return record, cloneMap(ctx.shared), nil
}

func buildCF26b(ctx *genCtx) (map[string]any, map[string]any, map[string]any) {
	b := cloneMap(ctx.shared)
	profileLimits(b)["maxCalls"] = -1
	record := baseRecord("F13-CF-26.b")
	return record, b, nil
}

func buildCF26c(ctx *genCtx) (map[string]any, map[string]any, map[string]any) {
	b := cloneMap(ctx.shared)
	profileLimits(b)["maxPayloadBytes"] = 1048577
	record := baseRecord("F13-CF-26.c")
	return record, b, nil
}

func buildCF26d(ctx *genCtx) (map[string]any, map[string]any, map[string]any) {
	b := cloneMap(ctx.shared)
	profileLimits(b)["maxConcurrency"] = -1
	record := baseRecord("F13-CF-26.d")
	return record, b, nil
}

func buildCF26e(ctx *genCtx) (map[string]any, map[string]any, map[string]any) {
	b := cloneMap(ctx.shared)
	profileLimits(b)["maxRetries"] = -1
	record := baseRecord("F13-CF-26.e")
	return record, b, nil
}

func buildCF26f(ctx *genCtx) (map[string]any, map[string]any, map[string]any) {
	b := cloneMap(ctx.shared)
	profileLimits(b)["maxJurisdictions"] = -1
	record := baseRecord("F13-CF-26.f")
	return record, b, nil
}

func buildCF26g(ctx *genCtx) (map[string]any, map[string]any, map[string]any) {
	record := baseRecord("F13-CF-26.g")
	eval := baseEvalResult()
	eval["timing"].(map[string]any)["durationMs"] = 99999
	withEval(record, eval)
	return record, cloneMap(ctx.shared), nil
}

func buildCF26h(ctx *genCtx) (map[string]any, map[string]any, map[string]any) {
	b := cloneMap(ctx.shared)
	ps := profileSpec(b)
	ps["limits"].(map[string]any)["maxCostBudget"] = "profile-budget"
	types := ps["acceptedEvaluationTypes"].([]any)
	narrow := types[1].(map[string]any)
	limits := narrow["limits"].(map[string]any)
	limits["maxCostBudget"] = "eval-budget"
	record := baseRecord("F13-CF-26.h")
	withEval(record, baseEvalResult())
	return record, b, nil
}

func buildScope07(ctx *genCtx) (map[string]any, map[string]any, map[string]any) {
	b := cloneMap(ctx.shared)
	ss := profileSpec(b)["scopeSemantics"].(map[string]any)
	ss["platformAllowed"] = false
	record := baseRecord("F13-SCOPE-07")
	delete(record["metadata"].(map[string]any), "scopeRef")
	rec := record["record"].(map[string]any)
	rec["semanticIdentity"].(map[string]any)["canonicalScope"] = "Platform"
	return record, b, nil
}

func buildScope09(ctx *genCtx) (map[string]any, map[string]any, map[string]any) {
	record := baseRecord("F13-SCOPE-09")
	return record, cloneMap(ctx.shared), map[string]any{"parallelScopeFields": []any{"/scopeRef"}}
}

func buildScope10(ctx *genCtx) (map[string]any, map[string]any, map[string]any) {
	record := baseRecord("F13-SCOPE-10")
	meta := record["metadata"].(map[string]any)
	meta["scopeRef"] = map[string]any{
		"apiVersion": "platform.sovrunn.io/v1alpha1",
		"kind":       "ServiceInstance",
		"name":       "svc-1",
		"uid":        "svc-uid-1",
	}
	rec := record["record"].(map[string]any)
	rec["semanticIdentity"].(map[string]any)["canonicalScope"] = "ServiceInstance:svc-uid-1"
	return record, cloneMap(ctx.shared), nil
}

func buildEval03(ctx *genCtx) (map[string]any, map[string]any, map[string]any) {
	b := cloneMap(ctx.shared)
	ps := profileSpec(b)
	types := ps["acceptedEvaluationTypes"].([]any)
	narrow := types[1].(map[string]any)
	limits := narrow["limits"].(map[string]any)
	limits["maxPayloadBytes"] = 131072 // wider than profile 65536
	record := baseRecord("F13-EVAL-03")
	withEval(record, baseEvalResult())
	return record, b, nil
}

func buildEval04(ctx *genCtx) (map[string]any, map[string]any, map[string]any) {
	b := cloneMap(ctx.shared)
	profileLimits(b)["maxPayloadBytes"] = 1048577
	record := baseRecord("F13-EVAL-04")
	return record, b, nil
}

func buildEval05(ctx *genCtx) (map[string]any, map[string]any, map[string]any) {
	record := baseRecord("F13-EVAL-05")
	eval := baseEvalResult()
	eval["resultStatus"] = "NOT_A_STATUS"
	withEval(record, eval)
	return record, cloneMap(ctx.shared), nil
}

func buildEval06(ctx *genCtx) (map[string]any, map[string]any, map[string]any) {
	record := baseRecord("F13-EVAL-06")
	eval := baseEvalResult()
	eval["scopeEvidence"] = map[string]any{
		"apiVersion": "core.sovrunn.io/v1alpha1",
		"kind":       "Organization",
		"name":       "acme-corp",
		"uid":        "org-456",
	}
	withEval(record, eval)
	return record, cloneMap(ctx.shared), nil
}

func buildEval07(ctx *genCtx) (map[string]any, map[string]any, map[string]any) {
	record := baseRecord("F13-EVAL-07")
	rec := record["record"].(map[string]any)
	rec["trust"] = map[string]any{
		"canonicalizationMethod": "fixture-canon-v1",
		"digestAlgorithm":        "fixture-digest-v1",
		"signatureAlgorithm":     "fixture-sig-v1",
		"state":                  "TRUSTED",
	}
	return record, cloneMap(ctx.shared), map[string]any{"expectedTrustState": "OTHER"}
}

func buildSEC01(ctx *genCtx) (map[string]any, map[string]any, map[string]any) {
	record := baseRecord("F13-SEC-01")
	opts := map[string]any{
		"declaredSensitivity": "CONFIDENTIAL",
		"classifications":     []any{"Public"},
		"canonicalView":       canonicalView(),
	}
	return record, cloneMap(ctx.shared), opts
}

func buildSEC02(ctx *genCtx) (map[string]any, map[string]any, map[string]any) {
	b := cloneMap(ctx.shared)
	profileLimits(b)["maxFieldCount"] = 65
	record := baseRecord("F13-SEC-02")
	return record, b, nil
}

func buildSEC03(ctx *genCtx) (map[string]any, map[string]any, map[string]any) {
	b := cloneMap(ctx.shared)
	profileLimits(b)["maxPayloadBytes"] = 1048577
	record := baseRecord("F13-SEC-03")
	return record, b, nil
}

func buildSEC04(ctx *genCtx) (map[string]any, map[string]any, map[string]any) {
	b := cloneMap(ctx.shared)
	profileLimits(b)["maxDepth"] = 33
	record := baseRecord("F13-SEC-04")
	return record, b, nil
}

func buildSEC05(ctx *genCtx) (map[string]any, map[string]any, map[string]any) {
	b := cloneMap(ctx.shared)
	profileSpec(b)["contentCategories"] = []any{"capacity"}
	record := baseRecord("F13-SEC-05")
	opts := map[string]any{
		"contentCategories": []any{"unregistered"},
		"canonicalView":     canonicalView(),
	}
	return record, b, opts
}

func buildSEC06(ctx *genCtx) (map[string]any, map[string]any, map[string]any) {
	record := baseRecord("F13-SEC-06")
	opts := map[string]any{
		"contentCategories":    []any{"raw-secret"},
		"prohibitedCategories": []any{"raw-secret"},
		"canonicalView":        canonicalView(),
	}
	return record, cloneMap(ctx.shared), opts
}

func buildSEC07(ctx *genCtx) (map[string]any, map[string]any, map[string]any) {
	record := baseRecord("F13-SEC-07")
	opts := map[string]any{
		"classifications":     []any{"Sensitive"},
		"declaredSensitivity": "INTERNAL",
		"canonicalView":       canonicalView(),
	}
	return record, cloneMap(ctx.shared), opts
}

func buildSEC08(ctx *genCtx) (map[string]any, map[string]any, map[string]any) {
	record := baseRecord("F13-SEC-08")
	opts := map[string]any{
		"classifications":     []any{"Top-secret"},
		"declaredSensitivity": "RESTRICTED",
		"canonicalView":       canonicalView(),
	}
	return record, cloneMap(ctx.shared), opts
}

func buildTrust03(ctx *genCtx) (map[string]any, map[string]any, map[string]any) {
	record := baseRecord("F13-TRUST-03")
	rec := record["record"].(map[string]any)
	rec["trust"] = map[string]any{
		"canonicalizationMethod": "fixture-canon-v1",
		"digestAlgorithm":        "fixture-digest-v1",
		"signatureAlgorithm":     "fixture-sig-v1",
		"state":                  "revoked",
	}
	return record, cloneMap(ctx.shared), map[string]any{"trustRequired": true}
}

// --- verification ---

func verify(record, bundleMap map[string]any, opts map[string]any, wantCode apiproblem.ViolationCode, wantField string) error {
	recordRaw, err := json.Marshal(record)
	if err != nil {
		return err
	}
	bundleRaw, err := json.Marshal(bundleMap)
	if err != nil {
		return err
	}
	loaded, prob := bundle.Load(bundleRaw, apivalid.ModeReadRepresentation)
	if prob != nil {
		return fmt.Errorf("bundle load: %s", prob.Detail)
	}
	view := loaded.View()

	var recOpts validate.DecisionRecordOptions
	if opts != nil {
		optsRaw, _ := json.Marshal(opts)
		recOpts, err = parseOpts(optsRaw)
		if err != nil {
			return err
		}
	}

	_, prob = validate.DecodeAndValidateDecisionRecordWith(recordRaw, apivalid.ModeReadRepresentation, view, recOpts)
	if prob == nil {
		return fmt.Errorf("expected violation %s at %s, got success", wantCode, wantField)
	}
	if prob.Code != apiproblem.CodeValidationFailed {
		return fmt.Errorf("expected VALIDATION_FAILED, got %s: %s", prob.Code, prob.Detail)
	}
	if len(prob.Violations) != 1 {
		return fmt.Errorf("expected exactly one violation, got %d: %#v", len(prob.Violations), prob.Violations)
	}
	v := prob.Violations[0]
	if v.Code != wantCode {
		return fmt.Errorf("expected code %s, got %s (field %s)", wantCode, v.Code, v.Field)
	}
	if wantField != "" && v.Field != wantField {
		return fmt.Errorf("expected field %s, got %s (code %s)", wantField, v.Field, v.Code)
	}
	return nil
}

func parseOpts(raw []byte) (validate.DecisionRecordOptions, error) {
	var dto struct {
		ParallelScopeFields  []string          `json:"parallelScopeFields,omitempty"`
		TrustRequired        *bool             `json:"trustRequired,omitempty"`
		ExpectedTrustState   string            `json:"expectedTrustState,omitempty"`
		Classifications      []string          `json:"classifications,omitempty"`
		DeclaredSensitivity  string            `json:"declaredSensitivity,omitempty"`
		ContentCategories    []string          `json:"contentCategories,omitempty"`
		ProhibitedCategories []string          `json:"prohibitedCategories,omitempty"`
		CanonicalView        []string          `json:"canonicalView,omitempty"`
		SupportedObligations map[string]bool   `json:"supportedObligations,omitempty"`
		Established          []json.RawMessage `json:"established,omitempty"`
	}
	if err := json.Unmarshal(raw, &dto); err != nil {
		return validate.DecisionRecordOptions{}, err
	}
	out := validate.DecisionRecordOptions{
		ParallelScopeFields:  append([]string(nil), dto.ParallelScopeFields...),
		TrustRequired:        dto.TrustRequired,
		ExpectedTrustState:   dto.ExpectedTrustState,
		DeclaredSensitivity:  decision.Sensitivity(dto.DeclaredSensitivity),
		ContentCategories:    append([]string(nil), dto.ContentCategories...),
		ProhibitedCategories: append([]string(nil), dto.ProhibitedCategories...),
		CanonicalView:        append([]string(nil), dto.CanonicalView...),
	}
	if len(dto.SupportedObligations) > 0 {
		out.SupportedObligations = make(map[string]bool, len(dto.SupportedObligations))
		for k, v := range dto.SupportedObligations {
			out.SupportedObligations[k] = v
		}
	}
	if len(dto.Classifications) > 0 {
		out.Classifications = make([]apimeta.DataClassification, 0, len(dto.Classifications))
		for _, c := range dto.Classifications {
			out.Classifications = append(out.Classifications, apimeta.DataClassification(c))
		}
	}
	if len(dto.Established) > 0 {
		out.Established = make([]validate.RelationshipNode, 0, len(dto.Established))
		for _, rawNode := range dto.Established {
			var node struct {
				Ref          apimeta.TypedRef               `json:"ref"`
				Relationship *decision.DecisionRelationship `json:"relationship,omitempty"`
			}
			if err := json.Unmarshal(rawNode, &node); err != nil {
				return validate.DecisionRecordOptions{}, err
			}
			out.Established = append(out.Established, validate.RelationshipNode{
				Ref:          node.Ref,
				Relationship: node.Relationship,
			})
		}
	}
	return out, nil
}

func marshalEnvelope(record, bundleMap map[string]any, opts map[string]any, wantCode apiproblem.ViolationCode, wantField string) ([]byte, error) {
	recordRaw, err := json.Marshal(record)
	if err != nil {
		return nil, err
	}
	env := fixtureEnvelope{
		Kind:               kindDecisionNegativeFixture,
		WantViolationCode:  string(wantCode),
		WantViolationField: wantField,
		DecisionRecord:     recordRaw,
	}
	if bundleMap != nil {
		env.Bundle, err = json.Marshal(bundleMap)
		if err != nil {
			return nil, err
		}
	}
	if opts != nil {
		env.Options, err = json.Marshal(opts)
		if err != nil {
			return nil, err
		}
	}
	return json.MarshalIndent(env, "", "  ")
}
