package apiconform

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"reflect"
	"strings"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/apivalid"
	"gopkg.in/yaml.v3"
)

// Deterministic seed for Property 2 reproducibility
// (F12-VALIDATION-001(2), F12-VALIDATION-002).
const property2Seed int64 = 20260723

const property2Iterations = 100

var property2Tiers = []string{"small", "large", "xlarge", "tiny"}

var property2Phases = []string{"Pending", "Ready", "Failed"}

// property2YAMLOnlySamples are YAML documents that use constructs outside the
// strict JSON-compatible subset. DecodeYAML must reject each before JSON
// normalization / structural validation (D-03a).
var property2YAMLOnlySamples = []struct {
	name           string
	raw            string
	detailContains string
}{
	{
		name: "anchor",
		raw: `apiVersion: v1
kind: Project
metadata:
  name: &x demo
spec:
  displayName: Demo
`,
		detailContains: "YAML anchors",
	},
	{
		name: "alias",
		raw: `apiVersion: v1
kind: Project
metadata:
  name: &x demo
spec:
  displayName: *x
`,
		detailContains: "YAML anchors",
	},
	{
		name: "merge-key",
		raw: `apiVersion: v1
kind: Project
metadata:
  name: demo
spec:
  <<:
    displayName: Demo
`,
		detailContains: "YAML merge keys",
	},
	{
		name: "custom-tag",
		raw: `apiVersion: v1
kind: Project
metadata:
  name: demo
spec:
  displayName: !foo Demo
`,
		// !foo is rejected by the TaggedStyle safety pass (same class as explicit tags).
		detailContains: "YAML",
	},
	{
		name: "explicit-tag",
		raw: `apiVersion: v1
kind: Project
metadata:
  name: demo
spec:
  displayName: !!str Demo
`,
		detailContains: "YAML explicit tags",
	},
	{
		name: "multiple-documents",
		raw: `apiVersion: v1
kind: Project
metadata:
  name: demo
spec:
  displayName: Demo
---
apiVersion: v1
kind: Project
metadata:
  name: other
spec:
  displayName: Other
`,
		detailContains: "multiple YAML documents",
	},
	{
		name: "non-string-int-key",
		raw: `apiVersion: v1
kind: Project
metadata:
  name: demo
spec:
  1: Demo
`,
		detailContains: "YAML mapping keys must be strings",
	},
	{
		name: "non-string-bool-key",
		raw: `apiVersion: v1
kind: Project
metadata:
  name: demo
spec:
  true: Demo
`,
		detailContains: "YAML mapping keys must be strings",
	},
	{
		name: "nan",
		raw: `apiVersion: v1
kind: Project
metadata:
  name: demo
spec:
  displayName: Demo
  count: .nan
`,
		detailContains: "non-finite YAML numbers",
	},
	{
		name: "inf",
		raw: `apiVersion: v1
kind: Project
metadata:
  name: demo
spec:
  displayName: Demo
  count: .inf
`,
		detailContains: "non-finite YAML numbers",
	},
	{
		name: "neg-inf",
		raw: `apiVersion: v1
kind: Project
metadata:
  name: demo
spec:
  displayName: Demo
  count: -.inf
`,
		detailContains: "non-finite YAML numbers",
	},
	{
		name: "timestamp",
		raw: `apiVersion: v1
kind: Project
metadata:
  name: demo
spec:
  displayName: 2020-01-01
`,
		detailContains: "YAML timestamp coercions",
	},
	{
		name: "binary",
		raw: `apiVersion: v1
kind: Project
metadata:
  name: demo
spec:
  displayName: !!binary aGVsbG8=
`,
		// !!binary may be reported via TaggedStyle ("explicit tags") or binary coercion.
		detailContains: "YAML",
	},
}

type property2CompatibleCase struct {
	JSON   []byte
	YAML   []byte
	Policy apivalid.FieldPolicy
	Mode   apivalid.DecodeMode
}

// Feature: api-resource-naming-status-and-validation-standard, Property 2: JSON/YAML validation equivalence
//
// For any object expressible in the strict JSON-compatible subset, decoding
// and validating its JSON form and its strict-YAML form produce equivalent
// results (same accept/reject outcome and, on success, the same normalized
// representation). For any YAML input using an alias, anchor, merge key,
// custom tag, multiple documents, a non-string mapping key, a non-finite
// number, or a YAML-only timestamp/binary coercion, decoding is rejected.
//
// Validates: Requirements 4.9 (F12-VALIDATION-001(2), F12-VALIDATION-002)
func TestProperty2_JSONYAMLValidationEquivalence(t *testing.T) {
	t.Parallel()

	v := yamlEquivValidator(t)
	lim := yamlEquivLimits()
	rng := rand.New(rand.NewSource(property2Seed))

	for i := 0; i < property2Iterations; i++ {
		compat := generateProperty2CompatibleCase(rng, i)
		if err := checkProperty2CompatibleCase(compat, v, lim, i); err != nil {
			t.Fatalf("property 2 compatible case failed at iteration %d (seed %d): %v", i, property2Seed, err)
		}

		only := property2YAMLOnlySamples[rng.Intn(len(property2YAMLOnlySamples))]
		if err := checkProperty2YAMLOnlyCase(only.raw, only.detailContains, lim, i, only.name); err != nil {
			t.Fatalf("property 2 YAML-only case %q failed at iteration %d (seed %d): %v",
				only.name, i, property2Seed, err)
		}
	}
}

func generateProperty2CompatibleCase(rng *rand.Rand, iteration int) property2CompatibleCase {
	modes := []apivalid.DecodeMode{
		apivalid.ModeCreateRequest,
		apivalid.ModeReplaceRequest,
		apivalid.ModeInternalObject,
		apivalid.ModeReadRepresentation,
	}
	mode := modes[rng.Intn(len(modes))]
	pol := apivalid.PolicyFor(mode)

	doc := map[string]any{
		"apiVersion": "v1",
		"kind":       "Project",
		"metadata": map[string]any{
			"name": fmt.Sprintf("obj-%d-%d", iteration, rng.Intn(10000)),
		},
		"spec": map[string]any{},
	}
	meta := doc["metadata"].(map[string]any)
	spec := doc["spec"].(map[string]any)

	// Shape buckets force coverage of accept, structural reject, FieldPolicy
	// reject, and unknown-field reject while remaining JSON-compatible.
	switch iteration % 8 {
	case 0:
		// Minimal valid project.
		spec["displayName"] = fmt.Sprintf("Demo-%d", rng.Intn(1000))
	case 1:
		// Valid with optional enum/integer.
		spec["displayName"] = fmt.Sprintf("Demo-%d", rng.Intn(1000))
		spec["tier"] = property2Tiers[rng.Intn(2)] // only schema-valid tiers
		spec["count"] = int64(rng.Intn(50))
	case 2:
		// Missing required displayName → structural violation (when decoded).
		if rng.Intn(2) == 0 {
			spec["tier"] = "small"
		}
	case 3:
		// Enum mismatch on tier.
		spec["displayName"] = fmt.Sprintf("Demo-%d", rng.Intn(1000))
		spec["tier"] = property2Tiers[2+rng.Intn(2)] // xlarge/tiny
	case 4:
		// Type mismatch on count (string instead of integer).
		spec["displayName"] = fmt.Sprintf("Demo-%d", rng.Intn(1000))
		spec["count"] = "nope"
	case 5:
		// Unknown top-level field → decode reject (typed path uses map dest).
		spec["displayName"] = fmt.Sprintf("Demo-%d", rng.Intn(1000))
		doc["mystery"] = true
	case 6:
		// Status + system-owned fields (FieldPolicy may accept or reject).
		spec["displayName"] = fmt.Sprintf("Demo-%d", rng.Intn(1000))
		meta["uid"] = fmt.Sprintf("%032x", rng.Uint64())
		meta["generation"] = int64(1 + rng.Intn(20))
		meta["resourceVersion"] = fmt.Sprintf("%d", 1+rng.Intn(1000))
		doc["status"] = map[string]any{
			"phase": property2Phases[rng.Intn(2)], // Pending/Ready only
		}
	default:
		// Nested unknown under spec + optional invalid kind.
		spec["displayName"] = fmt.Sprintf("Demo-%d", rng.Intn(1000))
		if rng.Intn(2) == 0 {
			spec["unexpected"] = true
		} else {
			doc["kind"] = "Tenant"
		}
	}

	jsonRaw, err := json.Marshal(doc)
	if err != nil {
		panic(fmt.Sprintf("property2 json marshal failed (seed %d iteration %d): %v", property2Seed, iteration, err))
	}
	yamlRaw, err := yaml.Marshal(doc)
	if err != nil {
		panic(fmt.Sprintf("property2 yaml marshal failed (seed %d iteration %d): %v", property2Seed, iteration, err))
	}

	return property2CompatibleCase{
		JSON:   jsonRaw,
		YAML:   yamlRaw,
		Policy: pol,
		Mode:   mode,
	}
}

func checkProperty2CompatibleCase(
	c property2CompatibleCase,
	v *StructuralValidator,
	lim apivalid.Limits,
	iteration int,
) error {
	var fromJSON, fromYAML map[string]any
	jsonProb, jsonViolations, jsonErr := decodeAndStructurallyValidate(
		c.JSON, false, lim, c.Policy, &fromJSON, v, yamlEquivSchemaID)
	yamlProb, yamlViolations, yamlErr := decodeAndStructurallyValidate(
		c.YAML, true, lim, c.Policy, &fromYAML, v, yamlEquivSchemaID)

	if err := compareProperty2DecodeProblems(jsonProb, yamlProb, iteration, c.Mode); err != nil {
		return err
	}

	// Decode failure: structural validation must not run; outcomes already match.
	if jsonProb != nil {
		if jsonErr != nil || yamlErr != nil {
			return fmt.Errorf("iteration %d mode %s: structural validate must not run after decode failure: jsonErr=%v yamlErr=%v",
				iteration, c.Mode, jsonErr, yamlErr)
		}
		if len(jsonViolations) != 0 || len(yamlViolations) != 0 {
			return fmt.Errorf("iteration %d mode %s: structural violations must be empty after decode failure: json=%#v yaml=%#v",
				iteration, c.Mode, jsonViolations, yamlViolations)
		}
		return nil
	}

	if (jsonErr == nil) != (yamlErr == nil) {
		return fmt.Errorf("iteration %d mode %s: structural err divergence: json=%v yaml=%v",
			iteration, c.Mode, jsonErr, yamlErr)
	}
	if jsonErr != nil {
		// Adapter unavailable is out of scope for Property 2; treat as failure.
		return fmt.Errorf("iteration %d mode %s: structural validate unavailable: %v", iteration, c.Mode, jsonErr)
	}

	if err := compareProperty2Violations(jsonViolations, yamlViolations, iteration, c.Mode); err != nil {
		return err
	}

	if len(jsonViolations) == 0 {
		if !reflect.DeepEqual(fromYAML, fromJSON) {
			return fmt.Errorf("iteration %d mode %s: accepted values diverge:\nYAML=%#v\nJSON=%#v\njsonBody=%s\nyamlBody=%s",
				iteration, c.Mode, fromYAML, fromJSON, c.JSON, c.YAML)
		}
	} else if !reflect.DeepEqual(fromYAML, fromJSON) {
		// Rejected instances may still decode to equivalent maps before
		// structural findings; require that equivalence too.
		return fmt.Errorf("iteration %d mode %s: rejected decoded values diverge:\nYAML=%#v\nJSON=%#v",
			iteration, c.Mode, fromYAML, fromJSON)
	}

	return nil
}

func checkProperty2YAMLOnlyCase(raw, detailContains string, lim apivalid.Limits, iteration int, name string) error {
	// Customer and permissive policies must both reject YAML-only constructs
	// before FieldPolicy / structural layers matter.
	policies := []apivalid.FieldPolicy{
		apivalid.PolicyFor(apivalid.ModeCreateRequest),
		apivalid.PolicyFor(apivalid.ModeReadRepresentation),
	}
	for _, pol := range policies {
		var dst map[string]any
		prob := apivalid.DecodeYAML([]byte(raw), lim, pol, &dst)
		if prob == nil {
			return fmt.Errorf("iteration %d YAML-only %q: DecodeYAML must reject, got nil (policy=%#v)",
				iteration, name, pol)
		}
		if prob.Code != apiproblem.CodeMalformedRequest {
			return fmt.Errorf("iteration %d YAML-only %q: Code=%q want %q detail=%q",
				iteration, name, prob.Code, apiproblem.CodeMalformedRequest, prob.Detail)
		}
		if detailContains != "" && !strings.Contains(prob.Detail, detailContains) {
			return fmt.Errorf("iteration %d YAML-only %q: Detail=%q want substring %q",
				iteration, name, prob.Detail, detailContains)
		}
		if len(prob.Violations) != 1 {
			return fmt.Errorf("iteration %d YAML-only %q: Violations=%#v want one",
				iteration, name, prob.Violations)
		}
		if !strings.HasPrefix(prob.Violations[0].Field, "/") {
			return fmt.Errorf("iteration %d YAML-only %q: Field %q is not an RFC 6901 JSON Pointer",
				iteration, name, prob.Violations[0].Field)
		}
	}
	return nil
}

func compareProperty2DecodeProblems(jsonProb, yamlProb *apiproblem.Problem, iteration int, mode apivalid.DecodeMode) error {
	if (jsonProb == nil) != (yamlProb == nil) {
		return fmt.Errorf("iteration %d mode %s: decode accept/reject diverge: json=%#v yaml=%#v",
			iteration, mode, jsonProb, yamlProb)
	}
	if jsonProb == nil {
		return nil
	}
	if jsonProb.Code != yamlProb.Code {
		return fmt.Errorf("iteration %d mode %s: decode Code diverge: json=%q yaml=%q",
			iteration, mode, jsonProb.Code, yamlProb.Code)
	}
	if len(jsonProb.Violations) != len(yamlProb.Violations) {
		return fmt.Errorf("iteration %d mode %s: decode Violations len diverge: json=%d yaml=%d\njson=%#v\nyaml=%#v",
			iteration, mode, len(jsonProb.Violations), len(yamlProb.Violations),
			jsonProb.Violations, yamlProb.Violations)
	}
	for i := range jsonProb.Violations {
		jv, yv := jsonProb.Violations[i], yamlProb.Violations[i]
		if jv.Field != yv.Field {
			return fmt.Errorf("iteration %d mode %s: decode Violations[%d].Field diverge: json=%q yaml=%q",
				iteration, mode, i, jv.Field, yv.Field)
		}
		if jv.Code != yv.Code {
			return fmt.Errorf("iteration %d mode %s: decode Violations[%d].Code diverge: json=%q yaml=%q",
				iteration, mode, i, jv.Code, yv.Code)
		}
		if !strings.HasPrefix(jv.Field, "/") {
			return fmt.Errorf("iteration %d mode %s: decode Field %q is not an RFC 6901 JSON Pointer",
				iteration, mode, jv.Field)
		}
	}
	return nil
}

func compareProperty2Violations(jsonVs, yamlVs []apiproblem.Violation, iteration int, mode apivalid.DecodeMode) error {
	if len(jsonVs) != len(yamlVs) {
		return fmt.Errorf("iteration %d mode %s: structural Violations len diverge: json=%d yaml=%d\njson=%#v\nyaml=%#v",
			iteration, mode, len(jsonVs), len(yamlVs), jsonVs, yamlVs)
	}
	for i := range jsonVs {
		if jsonVs[i].Field != yamlVs[i].Field {
			return fmt.Errorf("iteration %d mode %s: structural Violations[%d].Field diverge: json=%q yaml=%q",
				iteration, mode, i, jsonVs[i].Field, yamlVs[i].Field)
		}
		if jsonVs[i].Code != yamlVs[i].Code {
			return fmt.Errorf("iteration %d mode %s: structural Violations[%d].Code diverge: json=%q yaml=%q",
				iteration, mode, i, jsonVs[i].Code, yamlVs[i].Code)
		}
		if !strings.HasPrefix(jsonVs[i].Field, "/") {
			return fmt.Errorf("iteration %d mode %s: structural Field %q is not an RFC 6901 JSON Pointer",
				iteration, mode, jsonVs[i].Field)
		}
	}
	return nil
}
