package apiconform

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/apischema"
)

// Deterministic seed for Property 9 reproducibility
// (F12-NAMING-005, F12-VERIFY-001(13); D-01b).
const property9Seed int64 = 20260723

const property9Iterations = 100

// property9Scenario classifies the generated TypeBinding consistency case.
type property9Scenario string

const (
	property9MatchRegistered     property9Scenario = "match_registered"
	property9MismatchNil         property9Scenario = "mismatch_nil"
	property9MismatchPrimitive   property9Scenario = "mismatch_primitive"
	property9MismatchEmptyStruct property9Scenario = "mismatch_empty_struct"
	property9MismatchMap         property9Scenario = "mismatch_map"
	property9MismatchExtraField  property9Scenario = "mismatch_extra_field"
	property9MismatchWrongType   property9Scenario = "mismatch_wrong_type"
	property9MismatchCrossBind   property9Scenario = "mismatch_cross_bind"
)

type property9Case struct {
	Scenario   property9Scenario
	SchemaPath string
	Schema     []byte
	GoType     reflect.Type
	WantAccept bool
}

// Feature: api-resource-naming-status-and-validation-standard, Property 9: Derivative Go-type / schema consistency
//
// For any registered TypeBinding, VerifyGoTypeAgainstSchema accepts the
// derivative Go type if and only if it matches its canonical schema across
// the supported subset (property names, JSON tags, required vs optional,
// primitives, arrays/maps, embedded fields, enum-backed types, and
// additionalProperties behavior). Deliberate mismatches are rejected.
// Fixture round-tripping is supporting evidence only, not proof.
//
// Validates: Requirements 4.1, 4.9, 4.16 (F12-NAMING-005, F12-VERIFY-001(13))
func TestProperty9_DerivativeGoTypeSchemaConsistency(t *testing.T) {
	t.Parallel()

	schemas := loadProperty9Schemas(t)
	if len(TypeBindings) == 0 {
		t.Fatal("TypeBindings registry must not be empty")
	}
	if len(schemas) != len(TypeBindings) {
		t.Fatalf("loaded %d schemas for %d TypeBindings", len(schemas), len(TypeBindings))
	}

	rng := rand.New(rand.NewSource(property9Seed))
	for i := 0; i < property9Iterations; i++ {
		c := generateProperty9Case(rng, i, schemas)
		if err := checkProperty9Case(c, i); err != nil {
			t.Fatalf("property 9 failed at iteration %d (seed %d scenario %s schema %s): %v",
				i, property9Seed, c.Scenario, c.SchemaPath, err)
		}
	}
}

func loadProperty9Schemas(t *testing.T) map[string][]byte {
	t.Helper()

	root := moduleRoot(t)
	out := make(map[string][]byte, len(TypeBindings))
	for _, binding := range TypeBindings {
		path := filepath.Join(root, filepath.FromSlash(binding.SchemaPath))
		raw, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("read schema %s: %v", path, err)
		}
		support := apischema.ValidateSchemaSupport(raw)
		if len(support) > 0 {
			t.Fatalf("ValidateSchemaSupport failed for %s: %v", binding.SchemaPath, support)
		}
		out[binding.SchemaPath] = raw
	}
	return out
}

func generateProperty9Case(rng *rand.Rand, iteration int, schemas map[string][]byte) property9Case {
	// Force coverage of match and mismatch classes; occasional rng shuffle
	// keeps the 100 iterations seed-reproducible without a pure round-robin.
	bucket := iteration % 8
	if rng.Intn(20) == 0 {
		bucket = rng.Intn(8)
	}

	switch bucket {
	case 0, 1, 2:
		// Match: registered TypeBinding with its own derivative Go type.
		// Cycle bindings so every registry entry is exercised often.
		b := TypeBindings[iteration%len(TypeBindings)]
		if rng.Intn(3) == 0 {
			b = TypeBindings[rng.Intn(len(TypeBindings))]
		}
		return property9Case{
			Scenario:   property9MatchRegistered,
			SchemaPath: b.SchemaPath,
			Schema:     schemas[b.SchemaPath],
			GoType:     b.GoType,
			WantAccept: true,
		}
	case 3:
		b := TypeBindings[rng.Intn(len(TypeBindings))]
		return property9Case{
			Scenario:   property9MismatchNil,
			SchemaPath: b.SchemaPath,
			Schema:     schemas[b.SchemaPath],
			GoType:     nil,
			WantAccept: false,
		}
	case 4:
		b := TypeBindings[rng.Intn(len(TypeBindings))]
		prims := []reflect.Type{
			reflect.TypeOf(""),
			reflect.TypeOf(0),
			reflect.TypeOf(false),
			reflect.TypeOf([]string{}),
		}
		return property9Case{
			Scenario:   property9MismatchPrimitive,
			SchemaPath: b.SchemaPath,
			Schema:     schemas[b.SchemaPath],
			GoType:     prims[rng.Intn(len(prims))],
			WantAccept: false,
		}
	case 5:
		b := TypeBindings[rng.Intn(len(TypeBindings))]
		return property9Case{
			Scenario:   property9MismatchEmptyStruct,
			SchemaPath: b.SchemaPath,
			Schema:     schemas[b.SchemaPath],
			GoType:     reflect.TypeOf(struct{}{}),
			WantAccept: false,
		}
	case 6:
		b := TypeBindings[rng.Intn(len(TypeBindings))]
		maps := []reflect.Type{
			reflect.TypeOf(map[string]string{}),
			reflect.TypeOf(map[string]any{}),
			reflect.TypeOf(map[string]int{}),
		}
		return property9Case{
			Scenario:   property9MismatchMap,
			SchemaPath: b.SchemaPath,
			Schema:     schemas[b.SchemaPath],
			GoType:     maps[rng.Intn(len(maps))],
			WantAccept: false,
		}
	case 7:
		// Alternate between synthetic field mismatches and incompatible
		// cross-binding (never the TypedRef/ScopeRef/OwnerRef family pair).
		if rng.Intn(2) == 0 {
			return property9SyntheticFieldMismatch(rng, schemas)
		}
		return property9CrossBindMismatch(rng, schemas)
	default:
		b := TypeBindings[0]
		return property9Case{
			Scenario:   property9MatchRegistered,
			SchemaPath: b.SchemaPath,
			Schema:     schemas[b.SchemaPath],
			GoType:     b.GoType,
			WantAccept: true,
		}
	}
}

func property9SyntheticFieldMismatch(rng *rand.Rand, schemas map[string][]byte) property9Case {
	b := TypeBindings[rng.Intn(len(TypeBindings))]
	if rng.Intn(2) == 0 {
		// Extra JSON field never declared in any closed schema.
		goType := reflect.StructOf([]reflect.StructField{{
			Name: "DeliberateExtra",
			Type: reflect.TypeOf(""),
			Tag:  reflect.StructTag(`json:"__deliberate_extra_field__"`),
		}})
		return property9Case{
			Scenario:   property9MismatchExtraField,
			SchemaPath: b.SchemaPath,
			Schema:     schemas[b.SchemaPath],
			GoType:     goType,
			WantAccept: false,
		}
	}

	// Wrong primitive for a string-shaped field name that exists on Page and
	// is absent (extra) on other schemas — either way the type is rejected.
	goType := reflect.StructOf([]reflect.StructField{{
		Name: "NextPageToken",
		Type: reflect.TypeOf(0),
		Tag:  reflect.StructTag(`json:"nextPageToken,omitempty"`),
	}})
	return property9Case{
		Scenario:   property9MismatchWrongType,
		SchemaPath: b.SchemaPath,
		Schema:     schemas[b.SchemaPath],
		GoType:     goType,
		WantAccept: false,
	}
}

func property9CrossBindMismatch(rng *rand.Rand, schemas map[string][]byte) property9Case {
	// Prefer Page as the alien Go type: its single optional field cannot
	// satisfy any other registered object schema, and other types cannot
	// satisfy page.json either. Avoid TypedRef/ScopeRef/OwnerRef pairs —
	// they share the same JSON surface and may accept under the subset check.
	pageIdx := property9BindingIndex("api/schemas/_common/page.json")
	schemaIdx := rng.Intn(len(TypeBindings))
	if schemaIdx == pageIdx {
		schemaIdx = (schemaIdx + 1 + rng.Intn(len(TypeBindings)-1)) % len(TypeBindings)
	}
	schemaBinding := TypeBindings[schemaIdx]
	pageBinding := TypeBindings[pageIdx]

	usePageAsGoType := rng.Intn(2) == 0
	if usePageAsGoType {
		return property9Case{
			Scenario:   property9MismatchCrossBind,
			SchemaPath: schemaBinding.SchemaPath,
			Schema:     schemas[schemaBinding.SchemaPath],
			GoType:     pageBinding.GoType,
			WantAccept: false,
		}
	}
	return property9Case{
		Scenario:   property9MismatchCrossBind,
		SchemaPath: pageBinding.SchemaPath,
		Schema:     schemas[pageBinding.SchemaPath],
		GoType:     schemaBinding.GoType,
		WantAccept: false,
	}
}

func property9BindingIndex(schemaPath string) int {
	for i, b := range TypeBindings {
		if b.SchemaPath == schemaPath {
			return i
		}
	}
	return 0
}

func checkProperty9Case(c property9Case, iteration int) error {
	if len(c.Schema) == 0 {
		return fmt.Errorf("iteration %d: empty schema for %s", iteration, c.SchemaPath)
	}

	issues := apischema.VerifyGoTypeAgainstSchema(c.Schema, c.GoType)
	accepted := len(issues) == 0

	if c.WantAccept && !accepted {
		return fmt.Errorf("iteration %d: registered match rejected for %s → %s:%s",
			iteration, c.SchemaPath, property9TypePkg(c.GoType), property9FormatIssues(issues))
	}
	if !c.WantAccept && accepted {
		return fmt.Errorf("iteration %d: deliberate mismatch accepted for %s → %s (scenario %s)",
			iteration, c.SchemaPath, property9TypePkg(c.GoType), c.Scenario)
	}

	// Determinism: same inputs always yield the same accept/reject class.
	again := apischema.VerifyGoTypeAgainstSchema(c.Schema, c.GoType)
	if (len(again) == 0) != accepted {
		return fmt.Errorf("iteration %d: VerifyGoTypeAgainstSchema non-deterministic for %s",
			iteration, c.SchemaPath)
	}

	if !c.WantAccept {
		if err := property9AssertMismatchCodes(issues); err != nil {
			return fmt.Errorf("iteration %d: %w", iteration, err)
		}
	}
	return nil
}

func property9AssertMismatchCodes(issues []apischema.SchemaIssue) error {
	allowed := map[string]struct{}{
		apischema.CodeGoTypeMismatch:                 {},
		apischema.CodeGoFieldMissing:                 {},
		apischema.CodeGoFieldExtra:                   {},
		apischema.CodeGoRequiredMismatch:             {},
		apischema.CodeGoAdditionalPropertiesMismatch: {},
		apischema.CodeGoEnumTypeMismatch:             {},
		apischema.CodeMalformedSchema:                {},
	}
	for _, issue := range issues {
		if _, ok := allowed[issue.Code]; !ok {
			return fmt.Errorf("unexpected mismatch code %q at %s: %s",
				issue.Code, issue.Path, issue.Message)
		}
	}
	return nil
}

func property9TypePkg(t reflect.Type) string {
	if t == nil {
		return "<nil>"
	}
	if t.Name() != "" {
		if t.PkgPath() != "" {
			return t.PkgPath() + "." + t.Name()
		}
		return t.Name()
	}
	return t.String()
}

func property9FormatIssues(issues []apischema.SchemaIssue) string {
	var b strings.Builder
	for _, issue := range issues {
		b.WriteString("\n  ")
		b.WriteString(issue.Path)
		b.WriteString(" ")
		b.WriteString(issue.Code)
		b.WriteString(": ")
		b.WriteString(issue.Message)
	}
	return b.String()
}
