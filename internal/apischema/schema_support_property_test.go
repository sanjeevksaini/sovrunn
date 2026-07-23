package apischema

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"sort"
	"testing"
)

// Deterministic seed for Property 1 reproducibility
// (F12-NAMING-005, F12-VALIDATION-001(4)).
const property1Seed int64 = 20260723

const property1Iterations = 100

// Unsupported keywords deliberately outside SupportedKeywords. Used both as
// injectables at schema positions and as trap identifiers under properties /
// extension objects (those traps must NOT be treated as keywords).
var property1UnsupportedKeywords = []string{
	"oneOf",
	"anyOf",
	"allOf",
	"not",
	"if",
	"then",
	"else",
	"$defs",
	"definitions",
	"format",
	"unevaluatedProperties",
	"contentEncoding",
	"dependentRequired",
	"prefixItems",
	"contains",
	"minItems",
	"maxItems",
	"uniqueItems",
	"patternProperties",
	"dependentSchemas",
	"x-sovrunn-foo",
	"x-sovrunn-unknown",
}

var property1SafePropNames = []string{
	"name",
	"count",
	"phase",
	"tags",
	"items", // may collide with keyword name; under properties it is an identifier
	"metadata",
	"spec",
	"status",
	"a/b~c", // exercises JSON Pointer escaping when nested findings appear
}

var property1Profiles = []string{
	"ManagedResource",
	"ObservedExternalResource",
	"VersionedDefinition",
	"TransientRequestResult",
	"LongRunningOperation",
	"ImmutableRecord",
}

var property1Boundaries = []string{
	"customer-facing",
	"operator-facing",
	"adapter-facing",
	"plugin-facing",
	"internal-engine-facing",
	"governance-only",
}

var property1Scopes = []string{
	"Platform",
	"Organization",
	"OrganizationUnit",
	"Tenant",
	"Project",
	"Provider",
}

// property1Case is a generated schema with an oracle of unsupported
// schema-position keyword paths. Paths are RFC 6901 JSON Pointers.
type property1Case struct {
	Raw              []byte
	UnsupportedPaths []string // schema-position keywords outside SupportedKeywords
}

// Feature: api-resource-naming-status-and-validation-standard, Property 1: Fail-closed schema support
//
// For any schema document, if it contains a keyword outside the supported
// subset at a schema position, ValidateSchemaSupport rejects it with
// SCHEMA_UNSUPPORTED_KEYWORD; a schema passes support validation iff every
// schema-position keyword is in the declared SupportedKeywords set. Property
// identifiers under "properties" and fields inside registered x-sovrunn-*
// extension objects are not keywords. Fail-closed: every unsupported
// schema-position keyword is reported (none silently ignored).
//
// Validates: Requirements 4.1, 4.9 (F12-NAMING-005, F12-VALIDATION-001(4))
func TestProperty1_FailClosedSchemaSupport(t *testing.T) {
	t.Parallel()

	rng := rand.New(rand.NewSource(property1Seed))
	for i := 0; i < property1Iterations; i++ {
		c := generateProperty1Case(rng, i)
		if err := checkProperty1Case(c, i); err != nil {
			t.Fatalf("property 1 failed at iteration %d (seed %d): %v", i, property1Seed, err)
		}
	}
}

func generateProperty1Case(rng *rand.Rand, iteration int) property1Case {
	var unsupported []string

	// Occasionally emit a boolean schema (always supported).
	if rng.Intn(12) == 0 {
		raw := []byte("true")
		if rng.Intn(2) == 0 {
			raw = []byte("false")
		}
		return property1Case{Raw: raw, UnsupportedPaths: nil}
	}

	root := generateSupportedObjectSchema(rng, 0)

	// Mix of modes: supported-only, inject unsupported, traps only, or both.
	mode := rng.Intn(4)
	switch mode {
	case 0:
		// supported-only (possibly with keyword-looking property names / extension traps)
		if rng.Intn(2) == 0 {
			addPropertyNameTraps(rng, root)
		}
		if rng.Intn(2) == 0 {
			addExtensionFieldTraps(rng, root)
		}
	case 1:
		injectUnsupportedKeywords(rng, root, "", &unsupported, 1+rng.Intn(3))
	case 2:
		addPropertyNameTraps(rng, root)
		addExtensionFieldTraps(rng, root)
		injectUnsupportedKeywords(rng, root, "", &unsupported, 1+rng.Intn(2))
	default:
		injectUnsupportedKeywords(rng, root, "", &unsupported, 1+rng.Intn(4))
		if rng.Intn(2) == 0 {
			addPropertyNameTraps(rng, root)
		}
	}

	raw, err := json.Marshal(root)
	if err != nil {
		panic(fmt.Sprintf("property1 marshal failed (seed %d iteration %d): %v", property1Seed, iteration, err))
	}

	// Normalize path list for stable assertions.
	sort.Strings(unsupported)
	unsupported = uniqueSorted(unsupported)

	return property1Case{
		Raw:              raw,
		UnsupportedPaths: unsupported,
	}
}

func generateSupportedObjectSchema(rng *rand.Rand, depth int) map[string]any {
	obj := map[string]any{
		"type": "object",
	}

	if depth == 0 {
		obj["$schema"] = "https://json-schema.org/draft/2020-12/schema"
		obj["$id"] = fmt.Sprintf("https://sovrunn.io/schemas/prop1-%d.json", rng.Intn(100000))
		obj["title"] = fmt.Sprintf("Prop1-%d", rng.Intn(1000))
		obj["description"] = "property-1 generated schema"
		obj["x-sovrunn-profile"] = property1Profiles[rng.Intn(len(property1Profiles))]
		obj["x-sovrunn-boundary"] = property1Boundaries[rng.Intn(len(property1Boundaries))]
		obj["x-sovrunn-stability"] = "alpha"
		scope := property1Scopes[rng.Intn(len(property1Scopes))]
		obj["x-sovrunn-allowed-scopes"] = []any{scope}
	}

	nProps := 1 + rng.Intn(4)
	props := make(map[string]any, nProps)
	required := make([]any, 0, nProps)
	usedNames := make(map[string]struct{}, nProps)

	for i := 0; i < nProps; i++ {
		name := property1SafePropNames[rng.Intn(len(property1SafePropNames))]
		if _, ok := usedNames[name]; ok {
			name = fmt.Sprintf("field-%d-%d", depth, i)
		}
		usedNames[name] = struct{}{}

		var propSchema map[string]any
		switch rng.Intn(5) {
		case 0:
			propSchema = map[string]any{
				"type":      "string",
				"minLength": 1,
				"maxLength": 63,
				"pattern":   "^[a-z0-9-]+$",
				"default":   "demo",
				"examples":  []any{"demo", "payments"},
			}
		case 1:
			propSchema = map[string]any{
				"type":    "integer",
				"minimum": 0,
				"maximum": 100,
			}
		case 2:
			propSchema = map[string]any{
				"type": "string",
				"enum": []any{"Pending", "Ready", "Failed"},
			}
		case 3:
			item := map[string]any{"type": "string"}
			if depth < 2 && rng.Intn(3) == 0 {
				item = generateSupportedObjectSchema(rng, depth+1)
			}
			propSchema = map[string]any{
				"type":  "array",
				"items": item,
			}
		default:
			if depth < 2 && rng.Intn(2) == 0 {
				propSchema = generateSupportedObjectSchema(rng, depth+1)
			} else {
				propSchema = map[string]any{
					"type":                 "object",
					"additionalProperties": false,
				}
			}
		}

		if rng.Intn(3) == 0 {
			propSchema["x-sovrunn-field-policy"] = map[string]any{
				"classification":    "Public",
				"authorizedWriter":  "creator",
				"authorizedReaders": []any{"customer"},
				"mutability":        "immutable",
				"retention":         "standard",
				"redaction":         "none",
				"residency":         "any",
				"auditRequired":     true,
			}
		}

		if rng.Intn(8) == 0 {
			propSchema["$ref"] = "../_common/typed-ref.json"
			// $ref schemas are still support-valid; instance validation is out of scope.
		}

		props[name] = propSchema
		if rng.Intn(2) == 0 {
			required = append(required, name)
		}
	}

	obj["properties"] = props
	if len(required) > 0 {
		obj["required"] = required
	}
	if rng.Intn(2) == 0 {
		obj["additionalProperties"] = false
	} else if depth < 2 && rng.Intn(3) == 0 {
		obj["additionalProperties"] = map[string]any{"type": "string"}
	}

	return obj
}

// addPropertyNameTraps adds property identifiers that collide with unsupported
// keyword names. Under "properties" these must never be treated as keywords.
func addPropertyNameTraps(rng *rand.Rand, root map[string]any) {
	props, ok := root["properties"].(map[string]any)
	if !ok {
		props = map[string]any{}
		root["properties"] = props
	}
	n := 1 + rng.Intn(3)
	for i := 0; i < n; i++ {
		name := property1UnsupportedKeywords[rng.Intn(len(property1UnsupportedKeywords))]
		props[name] = map[string]any{"type": "string"}
	}
}

// addExtensionFieldTraps places keyword-looking fields inside a registered
// extension object. Those fields are not schema keywords.
func addExtensionFieldTraps(rng *rand.Rand, root map[string]any) {
	props, _ := root["properties"].(map[string]any)
	target := root
	if props != nil && len(props) > 0 && rng.Intn(2) == 0 {
		// Attach to a random property schema when available.
		keys := make([]string, 0, len(props))
		for k := range props {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		if ps, ok := props[keys[rng.Intn(len(keys))]].(map[string]any); ok {
			target = ps
		}
	}

	ext, ok := target["x-sovrunn-field-policy"].(map[string]any)
	if !ok {
		ext = map[string]any{
			"classification":    "Public",
			"authorizedWriter":  "creator",
			"authorizedReaders": []any{"customer"},
			"mutability":        "immutable",
			"retention":         "standard",
			"redaction":         "none",
			"residency":         "any",
			"auditRequired":     true,
		}
		target["x-sovrunn-field-policy"] = ext
	}
	n := 1 + rng.Intn(3)
	for i := 0; i < n; i++ {
		name := property1UnsupportedKeywords[rng.Intn(len(property1UnsupportedKeywords))]
		ext[name] = "not-a-schema-keyword-here"
	}
}

// injectUnsupportedKeywords places unsupported keywords at actual schema
// positions (root and/or nested property schemas). Records JSON Pointer paths.
func injectUnsupportedKeywords(rng *rand.Rand, node map[string]any, path string, unsupported *[]string, count int) {
	targets := collectInjectableSchemaObjects(node, path)
	if len(targets) == 0 {
		return
	}
	for i := 0; i < count; i++ {
		t := targets[rng.Intn(len(targets))]
		key := property1UnsupportedKeywords[rng.Intn(len(property1UnsupportedKeywords))]
		// Avoid overwriting an already-injected key on the same object so the
		// oracle path list stays aligned with distinct schema-position findings.
		if _, exists := t.obj[key]; exists {
			key = property1UnsupportedKeywords[(rng.Intn(len(property1UnsupportedKeywords))+i+1)%len(property1UnsupportedKeywords)]
			if _, exists := t.obj[key]; exists {
				continue
			}
		}
		t.obj[key] = map[string]any{"type": "string"}
		*unsupported = append(*unsupported, joinPointer(t.path, key))
	}
}

type schemaTarget struct {
	obj  map[string]any
	path string
}

func collectInjectableSchemaObjects(node map[string]any, path string) []schemaTarget {
	out := []schemaTarget{{obj: node, path: path}}
	if props, ok := node["properties"].(map[string]any); ok {
		keys := make([]string, 0, len(props))
		for k := range props {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			if child, ok := props[k].(map[string]any); ok {
				childPath := joinPointer(joinPointer(path, "properties"), k)
				out = append(out, collectInjectableSchemaObjects(child, childPath)...)
			}
		}
	}
	if items, ok := node["items"].(map[string]any); ok {
		out = append(out, collectInjectableSchemaObjects(items, joinPointer(path, "items"))...)
	}
	if ap, ok := node["additionalProperties"].(map[string]any); ok {
		out = append(out, collectInjectableSchemaObjects(ap, joinPointer(path, "additionalProperties"))...)
	}
	return out
}

func checkProperty1Case(c property1Case, iteration int) error {
	issues := ValidateSchemaSupport(c.Raw)

	wantReject := len(c.UnsupportedPaths) > 0
	gotUnsupported := unsupportedKeywordPaths(issues)

	if !wantReject {
		if len(gotUnsupported) != 0 {
			return fmt.Errorf("iteration %d: supported-only schema must pass, got unsupported paths %v issues=%#v body=%s",
				iteration, gotUnsupported, issues, c.Raw)
		}
		// May only contain non-unsupported findings for malformed docs; our
		// generator never emits those, so issues must be empty.
		if len(issues) != 0 {
			return fmt.Errorf("iteration %d: supported-only schema must have zero issues, got %#v body=%s",
				iteration, issues, c.Raw)
		}
		return nil
	}

	if len(gotUnsupported) == 0 {
		return fmt.Errorf("iteration %d: schema with unsupported keywords must be rejected, want paths %v body=%s",
			iteration, c.UnsupportedPaths, c.Raw)
	}

	// Fail-closed: every intentionally injected unsupported keyword is reported.
	gotSet := make(map[string]struct{}, len(gotUnsupported))
	for _, p := range gotUnsupported {
		gotSet[p] = struct{}{}
	}
	for _, want := range c.UnsupportedPaths {
		if _, ok := gotSet[want]; !ok {
			return fmt.Errorf("iteration %d: missing fail-closed report for %q; got %v issues=%#v body=%s",
				iteration, want, gotUnsupported, issues, c.Raw)
		}
	}

	// No false positives: every SCHEMA_UNSUPPORTED_KEYWORD must be a real
	// unsupported schema-position keyword (oracle: not in SupportedKeywords,
	// and path corresponds to a schema-position key in the document).
	var root any
	if err := json.Unmarshal(c.Raw, &root); err != nil {
		return fmt.Errorf("iteration %d: generated schema must be valid JSON: %v", iteration, err)
	}
	oracle := map[string]struct{}{}
	collectUnsupportedSchemaKeywordPaths(root, "", oracle)
	for _, p := range gotUnsupported {
		if _, ok := oracle[p]; !ok {
			return fmt.Errorf("iteration %d: unexpected unsupported path %q (oracle=%v) body=%s",
				iteration, p, keysOf(oracle), c.Raw)
		}
	}
	// Completeness vs independent oracle (not only vs injected set).
	if len(gotSet) != len(oracle) {
		return fmt.Errorf("iteration %d: unsupported path set mismatch: got %v oracle %v body=%s",
			iteration, gotUnsupported, keysOf(oracle), c.Raw)
	}

	return nil
}

func unsupportedKeywordPaths(issues []SchemaIssue) []string {
	var paths []string
	for _, iss := range issues {
		if iss.Code == CodeUnsupportedKeyword {
			paths = append(paths, iss.Path)
		}
	}
	sort.Strings(paths)
	return paths
}

// collectUnsupportedSchemaKeywordPaths is a test-local oracle that mirrors the
// context-aware walk rules: property identifiers and registered extension
// object fields are not keywords; only schema-position keys outside
// SupportedKeywords are recorded.
func collectUnsupportedSchemaKeywordPaths(node any, path string, out map[string]struct{}) {
	switch node.(type) {
	case bool, nil:
		return
	}
	obj, ok := node.(map[string]any)
	if !ok {
		return
	}

	keys := make([]string, 0, len(obj))
	for k := range obj {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		val := obj[key]
		childPath := joinPointer(path, key)
		if !IsSupportedKeyword(key) {
			out[childPath] = struct{}{}
			continue
		}
		switch {
		case key == "properties":
			if props, ok := val.(map[string]any); ok {
				propKeys := make([]string, 0, len(props))
				for pk := range props {
					propKeys = append(propKeys, pk)
				}
				sort.Strings(propKeys)
				for _, pk := range propKeys {
					collectUnsupportedSchemaKeywordPaths(props[pk], joinPointer(childPath, pk), out)
				}
			}
		case key == "items":
			switch v := val.(type) {
			case []any:
				for i, item := range v {
					collectUnsupportedSchemaKeywordPaths(item, joinPointer(childPath, fmt.Sprintf("%d", i)), out)
				}
			default:
				collectUnsupportedSchemaKeywordPaths(v, childPath, out)
			}
		case key == "additionalProperties":
			if _, isBool := val.(bool); !isBool {
				collectUnsupportedSchemaKeywordPaths(val, childPath, out)
			}
		case IsRegisteredExtension(key):
			// Extension object fields are not schema keywords.
		}
	}
}

func uniqueSorted(in []string) []string {
	if len(in) == 0 {
		return nil
	}
	out := make([]string, 0, len(in))
	var prev string
	first := true
	for _, s := range in {
		if first || s != prev {
			out = append(out, s)
			prev = s
			first = false
		}
	}
	return out
}

func keysOf(m map[string]struct{}) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}
