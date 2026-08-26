package policyeval

import (
	"context"
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"gopkg.in/yaml.v3"
)

// Exact closed production import allowlist (D-01 / D-03 / D-14 / AC-F17-24).
// Compared against parsed import paths only — never against source keywords.
var productionImportAllowlist = map[string]struct{}{
	"github.com/sanjeevksaini/sovrunn/internal/apimeta":           {},
	"github.com/sanjeevksaini/sovrunn/internal/apiref":            {},
	"github.com/sanjeevksaini/sovrunn/internal/decision":          {},
	"github.com/sanjeevksaini/sovrunn/internal/decision/validate": {},
	"bytes":         {},
	"context":       {},
	"crypto/sha256": {},
	"encoding/hex":  {},
	"encoding/json": {},
	"errors":        {},
	"io":            {},
	"reflect":       {},
	"sort":          {},
	"strings":       {},
	"time":          {},
	"unicode/utf8":  {},
}

// Forbidden production type/symbol identifiers (AC-F17-22). Built by
// concatenation so this test file's own string literals are not relied on as
// a vacuous source-grep match against production files.
var forbiddenProductionSymbols = []string{
	"Decision" + "Record",
	"Audit" + "Event",
	"Route",
	"Router",
	"RouteHandler",
	"Controller",
	"Store",
	"Storage",
}

// Curated deny-list tokens for fixture reason-code review (AC-F17-17).
var fixtureReasonCodeDenyTokens = []string{
	"aws",
	"azure",
	"gcp",
	"customer",
	"secret",
	"credential",
	"password",
	"token",
}

func TestConformance_ClosedProductionImportSurface(t *testing.T) {
	t.Parallel()

	files, err := listProductionGoFiles(t)
	if err != nil {
		t.Fatalf("list production files: %v", err)
	}
	if len(files) == 0 {
		t.Fatal("no production .go files found; import-surface proof is vacuous")
	}

	fset := token.NewFileSet()
	seenImport := false
	for _, path := range files {
		file, err := parser.ParseFile(fset, path, nil, parser.ImportsOnly)
		if err != nil {
			t.Fatalf("parse %s: %v", path, err)
		}
		for _, imp := range file.Imports {
			pathLit, err := strconv.Unquote(imp.Path.Value)
			if err != nil {
				t.Fatalf("%s: unquote import %s: %v", path, imp.Path.Value, err)
			}
			seenImport = true
			if _, ok := productionImportAllowlist[pathLit]; !ok {
				t.Errorf("%s imports %q which is outside the closed D-01/D-03/D-14 allowlist", path, pathLit)
			}
		}
	}
	if !seenImport {
		t.Fatal("parsed zero production imports; import-surface proof is vacuous")
	}
}

func TestConformance_NoForbiddenProductionSymbols(t *testing.T) {
	t.Parallel()

	files, err := listProductionGoFiles(t)
	if err != nil {
		t.Fatalf("list production files: %v", err)
	}
	if len(files) == 0 {
		t.Fatal("no production .go files found; symbol-absence proof is vacuous")
	}

	forbidden := make(map[string]struct{}, len(forbiddenProductionSymbols))
	for _, name := range forbiddenProductionSymbols {
		forbidden[name] = struct{}{}
	}

	fset := token.NewFileSet()
	for _, path := range files {
		file, err := parser.ParseFile(fset, path, nil, 0)
		if err != nil {
			t.Fatalf("parse %s: %v", path, err)
		}
		ast.Inspect(file, func(n ast.Node) bool {
			ident, ok := n.(*ast.Ident)
			if !ok || ident == nil {
				return true
			}
			if _, hit := forbidden[ident.Name]; hit {
				pos := fset.Position(ident.Pos())
				t.Errorf("%s:%d:%d references forbidden symbol %q (AC-F17-22)", path, pos.Line, pos.Column, ident.Name)
			}
			return true
		})
	}
}

func TestConformance_FixtureCorpusReasonCodes(t *testing.T) {
	t.Parallel()

	dir := packageDir(t)
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("readdir: %v", err)
	}

	fset := token.NewFileSet()
	var (
		goFixtureSources   int
		yamlFixtureSources int
		reasonCodes        []string
	)

	for _, ent := range entries {
		name := ent.Name()
		if ent.IsDir() || !strings.HasSuffix(name, "_test.go") {
			continue
		}
		path := filepath.Join(dir, name)
		file, err := parser.ParseFile(fset, path, nil, 0)
		if err != nil {
			t.Fatalf("parse %s: %v", path, err)
		}

		stringSliceAssigns := collectStringSliceAssignments(file)

		var stack []ast.Node
		ast.Inspect(file, func(n ast.Node) bool {
			if n == nil {
				if len(stack) > 0 {
					stack = stack[:len(stack)-1]
				}
				return true
			}
			parents := stack
			stack = append(stack, n)

			switch node := n.(type) {
			case *ast.CompositeLit:
				typeName := compositeTypeName(node.Type)
				if typeName == "" && elidedFakeFixtureElement(node, parents) {
					typeName = "FakeFixture"
				}
				switch typeName {
				case "FakeFixture":
					goFixtureSources++
					codes, ok := reasonCodesFromFakeFixtureLit(node, stringSliceAssigns)
					if !ok {
						pos := fset.Position(node.Pos())
						t.Errorf("%s:%d: FakeFixture source could not be classified for reason-code extraction", path, pos.Line)
						return true
					}
					reasonCodes = append(reasonCodes, codes...)
				case "Conclusion":
					goFixtureSources++
					codes, ok := reasonCodesFromConclusionLit(node, stringSliceAssigns)
					if !ok {
						pos := fset.Position(node.Pos())
						t.Errorf("%s:%d: Conclusion source could not be classified for reason-code extraction", path, pos.Line)
						return true
					}
					reasonCodes = append(reasonCodes, codes...)
				}
			case *ast.CallExpr:
				fn := callName(node.Fun)
				switch fn {
				case "fakeConclusion":
					for i, arg := range node.Args {
						if i == 0 {
							continue
						}
						if s, ok := stringBasicLit(arg); ok {
							reasonCodes = append(reasonCodes, s)
						}
					}
				case "decodeFakeFixturesYAMLForTest":
					if len(node.Args) != 1 {
						pos := fset.Position(node.Pos())
						t.Errorf("%s:%d: decodeFakeFixturesYAMLForTest arity=%d; fixture source unclassified", path, pos.Line, len(node.Args))
						return true
					}
					raw, ok := resolveByteOrStringLiteral(node.Args[0], file)
					if !ok {
						pos := fset.Position(node.Pos())
						t.Errorf("%s:%d: YAML fixture argument could not be classified as a literal", path, pos.Line)
						return true
					}
					yamlFixtureSources++
					codes, err := reasonCodesFromFixtureYAML(raw)
					if err != nil {
						pos := fset.Position(node.Pos())
						t.Errorf("%s:%d: YAML fixture corpus parse failed: %v", path, pos.Line, err)
						return true
					}
					reasonCodes = append(reasonCodes, codes...)
				}
			}
			return true
		})
	}

	if goFixtureSources == 0 {
		t.Fatal("fixture corpus discovered zero FakeFixture/Conclusion composite literals; proof is vacuous")
	}
	if yamlFixtureSources == 0 {
		t.Fatal("fixture corpus discovered zero decodeFakeFixturesYAMLForTest YAML literals; proof is vacuous")
	}
	if len(reasonCodes) == 0 {
		t.Fatal("fixture corpus discovered zero reason codes; proof is vacuous")
	}

	for _, code := range reasonCodes {
		lower := strings.ToLower(code)
		for _, deny := range fixtureReasonCodeDenyTokens {
			if strings.Contains(lower, deny) {
				t.Errorf("fixture reason code %q contains denied token %q (AC-F17-17)", code, deny)
			}
		}
	}
}

func elidedFakeFixtureElement(lit *ast.CompositeLit, parents []ast.Node) bool {
	if lit.Type != nil || len(parents) == 0 {
		return false
	}
	parent, ok := parents[len(parents)-1].(*ast.CompositeLit)
	if !ok || parent == nil {
		return false
	}
	arr, ok := parent.Type.(*ast.ArrayType)
	if !ok {
		return false
	}
	return compositeTypeName(arr.Elt) == "FakeFixture"
}

func TestLeakage_SentinelInputsSanitized(t *testing.T) {
	t.Parallel()

	const (
		sentinelAction        = "SENTINEL_LEAK_ACTION_F17_11_A7K9"
		sentinelRequestID     = "SENTINEL_LEAK_REQID_F17_11_B3M2"
		sentinelSubjectName   = "SENTINEL_LEAK_SUBJECT_F17_11_C8P4"
		sentinelAdapterErr    = "SENTINEL_LEAK_ADAPTER_ERR_F17_11_D1Q5"
		sentinelReasonCode    = "sentinel-leak-reason-f17-11-e2r6" // grammar-invalid
		sentinelDigest        = "SENTINEL_LEAK_DIGEST_F17_11_F3S7_NOT_HEX"
		sentinelUnknownRoot   = "SENTINEL_LEAK_UNK_ROOT_F17_11_G4T8"
		sentinelUnknownNested = "SENTINEL_LEAK_UNK_NEST_F17_11_H5U9"
		sentinelDupRoot       = "SENTINEL_LEAK_DUP_ROOT_F17_11_I6V0"
		sentinelDupNested     = "SENTINEL_LEAK_DUP_NEST_F17_11_J7W1"
	)

	allSentinels := []string{
		sentinelAction,
		sentinelRequestID,
		sentinelSubjectName,
		sentinelAdapterErr,
		sentinelReasonCode,
		sentinelDigest,
		sentinelUnknownRoot,
		sentinelUnknownNested,
		sentinelDupRoot,
		sentinelDupNested,
	}
	assertNoSentinel := func(t *testing.T, label string, err error, nr NonResult) {
		t.Helper()
		var haystacks []string
		if err != nil {
			haystacks = append(haystacks, err.Error())
		}
		if nr != "" {
			haystacks = append(haystacks, string(nr))
		}
		for _, hay := range haystacks {
			for _, s := range allSentinels {
				if strings.Contains(hay, s) {
					t.Fatalf("%s leaked sentinel %q via %q", label, s, hay)
				}
			}
		}
	}

	fixed := NewFixedTimeSource(time.Date(2026, 8, 26, 7, 30, 0, 0, time.UTC))

	t.Run("request_field_values", func(t *testing.T) {
		t.Parallel()
		spy := &evaluateSpyAdapter{conclusion: validAllowConclusion()}
		b := mustBoundary(t, spy, fixed)

		req := frozenVectorRequest()
		req.Action = sentinelAction
		req.RequestID = sentinelRequestID
		req.SubjectRef.Name = sentinelSubjectName
		req.ProfileRefs = nil // force RequestInvalid after sentinel injection

		result, timing, nr, err := b.Evaluate(context.Background(), req)
		assertZeroNonResult(t, result, timing, nr, NonResultRequestInvalid, err, errRequestInvalid, spy.callCount())
		assertNoSentinel(t, "request_field_values", err, nr)
	})

	t.Run("adapter_error_message", func(t *testing.T) {
		t.Parallel()
		spy := &evaluateSpyAdapter{err: errors.New(sentinelAdapterErr)}
		b := mustBoundary(t, spy, fixed)

		result, timing, nr, err := b.Evaluate(context.Background(), frozenVectorRequest())
		if spy.callCount() != 1 {
			t.Fatalf("adapter calls=%d, want 1", spy.callCount())
		}
		if nr != NonResultAdapterFailure || !errors.Is(err, errAdapterFailed) {
			t.Fatalf("nr=%q err=%v, want AdapterFailure/errAdapterFailed", nr, err)
		}
		assertZeroResultTiming(t, result, timing)
		assertNoSentinel(t, "adapter_error_message", err, nr)
	})

	t.Run("reason_code_text", func(t *testing.T) {
		t.Parallel()
		spy := &evaluateSpyAdapter{
			conclusion: Conclusion{
				Outcome:     OutcomeAllow,
				ReasonCodes: []string{sentinelReasonCode}, // fails grammar → InvalidAdapterResult
			},
		}
		b := mustBoundary(t, spy, fixed)

		result, timing, nr, err := b.Evaluate(context.Background(), frozenVectorRequest())
		if spy.callCount() != 1 {
			t.Fatalf("adapter calls=%d, want 1", spy.callCount())
		}
		if nr != NonResultInvalidAdapterResult || !errors.Is(err, errInvalidAdapterResult) {
			t.Fatalf("nr=%q err=%v, want InvalidAdapterResult", nr, err)
		}
		assertZeroResultTiming(t, result, timing)
		assertNoSentinel(t, "reason_code_text", err, nr)
	})

	t.Run("fixture_digest", func(t *testing.T) {
		t.Parallel()
		_, err := NewDeterministicFake([]FakeFixture{{
			InputDigest:    sentinelDigest,
			AdapterFailure: true,
		}})
		if !errors.Is(err, errFakeConfigurationInvalid) {
			t.Fatalf("err=%v, want errFakeConfigurationInvalid", err)
		}
		assertNoSentinel(t, "fixture_digest", err, "")
	})

	t.Run("unknown_root_member", func(t *testing.T) {
		t.Parallel()
		raw := fmt.Sprintf(`{
			"subjectRef":{"apiVersion":"services.sovrunn.io/v1alpha1","kind":"ServiceInstance","name":"reporting-api"},
			"action":"service.read",
			"profileRefs":[{"apiVersion":"iam.sovrunn.io/v1alpha1","kind":"RoleDefinition","name":"reader","uid":"role-001"}],
			"requestId":"corr-001",
			%q:true
		}`, sentinelUnknownRoot)
		_, err := DecodePolicyEvaluationRequestJSON([]byte(raw))
		if !errors.Is(err, errRequestInvalid) {
			t.Fatalf("err=%v, want errRequestInvalid", err)
		}
		assertNoSentinel(t, "unknown_root_member", err, "")
	})

	t.Run("unknown_nested_member", func(t *testing.T) {
		t.Parallel()
		raw := fmt.Sprintf(`{
			"subjectRef":{"apiVersion":"services.sovrunn.io/v1alpha1","kind":"ServiceInstance","name":"reporting-api",%q:"x"},
			"action":"service.read",
			"profileRefs":[{"apiVersion":"iam.sovrunn.io/v1alpha1","kind":"RoleDefinition","name":"reader","uid":"role-001"}],
			"requestId":"corr-001"
		}`, sentinelUnknownNested)
		_, err := DecodePolicyEvaluationRequestJSON([]byte(raw))
		if !errors.Is(err, errRequestInvalid) {
			t.Fatalf("err=%v, want errRequestInvalid", err)
		}
		assertNoSentinel(t, "unknown_nested_member", err, "")
	})

	t.Run("duplicate_root_member", func(t *testing.T) {
		t.Parallel()
		raw := fmt.Sprintf(`{
			"subjectRef":{"apiVersion":"services.sovrunn.io/v1alpha1","kind":"ServiceInstance","name":"reporting-api"},
			"action":"service.read",
			"profileRefs":[{"apiVersion":"iam.sovrunn.io/v1alpha1","kind":"RoleDefinition","name":"reader","uid":"role-001"}],
			"requestId":"corr-001",
			%q:"one",
			%q:"two"
		}`, sentinelDupRoot, sentinelDupRoot)
		_, err := DecodePolicyEvaluationRequestJSON([]byte(raw))
		if !errors.Is(err, errRequestInvalid) {
			t.Fatalf("err=%v, want errRequestInvalid", err)
		}
		assertNoSentinel(t, "duplicate_root_member", err, "")
	})

	t.Run("duplicate_nested_member", func(t *testing.T) {
		t.Parallel()
		raw := fmt.Sprintf(`{
			"subjectRef":{"apiVersion":"services.sovrunn.io/v1alpha1","kind":"ServiceInstance","name":"reporting-api",%q:"one",%q:"two"},
			"action":"service.read",
			"profileRefs":[{"apiVersion":"iam.sovrunn.io/v1alpha1","kind":"RoleDefinition","name":"reader","uid":"role-001"}],
			"requestId":"corr-001"
		}`, sentinelDupNested, sentinelDupNested)
		_, err := DecodePolicyEvaluationRequestJSON([]byte(raw))
		if !errors.Is(err, errRequestInvalid) {
			t.Fatalf("err=%v, want errRequestInvalid", err)
		}
		assertNoSentinel(t, "duplicate_nested_member", err, "")
	})
}

func listProductionGoFiles(t *testing.T) ([]string, error) {
	t.Helper()
	dir := packageDir(t)
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var out []string
	for _, ent := range entries {
		name := ent.Name()
		if ent.IsDir() || !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") {
			continue
		}
		out = append(out, filepath.Join(dir, name))
	}
	return out, nil
}

func packageDir(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("getcwd: %v", err)
	}
	return dir
}

func compositeTypeName(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.SelectorExpr:
		return t.Sel.Name
	default:
		return ""
	}
}

func callName(fun ast.Expr) string {
	switch f := fun.(type) {
	case *ast.Ident:
		return f.Name
	case *ast.SelectorExpr:
		return f.Sel.Name
	default:
		return ""
	}
}

func stringBasicLit(expr ast.Expr) (string, bool) {
	lit, ok := expr.(*ast.BasicLit)
	if !ok || lit.Kind != token.STRING {
		return "", false
	}
	s, err := strconv.Unquote(lit.Value)
	if err != nil {
		return "", false
	}
	return s, true
}

func collectStringSliceAssignments(file *ast.File) map[string][]string {
	out := make(map[string][]string)
	ast.Inspect(file, func(n ast.Node) bool {
		assign, ok := n.(*ast.AssignStmt)
		if !ok || len(assign.Lhs) != 1 || len(assign.Rhs) != 1 {
			return true
		}
		ident, ok := assign.Lhs[0].(*ast.Ident)
		if !ok {
			return true
		}
		codes, ok := stringSliceComposite(assign.Rhs[0])
		if !ok {
			return true
		}
		out[ident.Name] = codes
		return true
	})
	return out
}

func stringSliceComposite(expr ast.Expr) ([]string, bool) {
	comp, ok := expr.(*ast.CompositeLit)
	if !ok {
		return nil, false
	}
	switch typ := comp.Type.(type) {
	case *ast.ArrayType:
		if ident, ok := typ.Elt.(*ast.Ident); !ok || ident.Name != "string" {
			return nil, false
		}
	default:
		return nil, false
	}
	var codes []string
	for _, elt := range comp.Elts {
		switch e := elt.(type) {
		case *ast.BasicLit:
			s, ok := stringBasicLit(e)
			if !ok {
				return nil, false
			}
			codes = append(codes, s)
		case *ast.Ident, *ast.SelectorExpr:
			// Named constants / table fields (e.g. tc.reasonCode): classified
			// without a static literal at this node.
			continue
		default:
			return nil, false
		}
	}
	return codes, true
}

func reasonCodesFromFakeFixtureLit(lit *ast.CompositeLit, assigns map[string][]string) ([]string, bool) {
	for _, elt := range lit.Elts {
		kv, ok := elt.(*ast.KeyValueExpr)
		if !ok {
			continue
		}
		key, ok := kv.Key.(*ast.Ident)
		if !ok || key.Name != "Conclusion" {
			continue
		}
		switch v := kv.Value.(type) {
		case *ast.Ident:
			if v.Name == "nil" {
				return nil, true
			}
			// Conclusion: tc.conclusion / conclusion variable — classified;
			// nested Conclusion composite lits are scanned independently.
			return nil, true
		case *ast.UnaryExpr:
			if v.Op != token.AND {
				return nil, false
			}
			inner, ok := v.X.(*ast.CompositeLit)
			if !ok || compositeTypeName(inner.Type) != "Conclusion" {
				return nil, false
			}
			return reasonCodesFromConclusionLit(inner, assigns)
		case *ast.CompositeLit:
			if compositeTypeName(v.Type) != "Conclusion" {
				return nil, false
			}
			return reasonCodesFromConclusionLit(v, assigns)
		case *ast.CallExpr:
			if callName(v.Fun) == "fakeConclusion" {
				var codes []string
				for i, arg := range v.Args {
					if i == 0 {
						continue
					}
					s, ok := stringBasicLit(arg)
					if !ok {
						return nil, false
					}
					codes = append(codes, s)
				}
				return codes, true
			}
			return nil, false
		case *ast.SelectorExpr:
			// e.g. Conclusion: tc.conclusion — nested Conclusion lits / table
			// entries are classified independently.
			return nil, true
		default:
			return nil, false
		}
	}
	// AdapterFailure-only fixtures have no conclusion/reason codes.
	return nil, true
}

func reasonCodesFromConclusionLit(lit *ast.CompositeLit, assigns map[string][]string) ([]string, bool) {
	for _, elt := range lit.Elts {
		kv, ok := elt.(*ast.KeyValueExpr)
		if !ok {
			continue
		}
		key, ok := kv.Key.(*ast.Ident)
		if !ok || key.Name != "ReasonCodes" {
			continue
		}
		switch v := kv.Value.(type) {
		case *ast.Ident:
			if v.Name == "nil" {
				return nil, true
			}
			if codes, ok := assigns[v.Name]; ok {
				return codes, true
			}
			return nil, true
		case *ast.SelectorExpr:
			return nil, true
		case *ast.CompositeLit:
			return stringSliceComposite(v)
		case *ast.CallExpr:
			return nil, true
		default:
			return nil, false
		}
	}
	return nil, true
}

func resolveByteOrStringLiteral(expr ast.Expr, file *ast.File) ([]byte, bool) {
	switch v := expr.(type) {
	case *ast.BasicLit:
		if v.Kind != token.STRING {
			return nil, false
		}
		s, err := strconv.Unquote(v.Value)
		if err != nil {
			return nil, false
		}
		return []byte(s), true
	case *ast.CallExpr:
		if len(v.Args) != 1 {
			return nil, false
		}
		arr, ok := v.Fun.(*ast.ArrayType)
		if !ok {
			return nil, false
		}
		elt, ok := arr.Elt.(*ast.Ident)
		if !ok || elt.Name != "byte" {
			return nil, false
		}
		return resolveByteOrStringLiteral(v.Args[0], file)
	case *ast.Ident:
		var found []byte
		ok := false
		ast.Inspect(file, func(n ast.Node) bool {
			assign, isAssign := n.(*ast.AssignStmt)
			if !isAssign || len(assign.Lhs) != 1 || len(assign.Rhs) != 1 {
				return true
			}
			lhs, isIdent := assign.Lhs[0].(*ast.Ident)
			if !isIdent || lhs.Name != v.Name {
				return true
			}
			if raw, resolved := resolveByteOrStringLiteral(assign.Rhs[0], file); resolved {
				found = raw
				ok = true
			}
			return true
		})
		return found, ok
	default:
		return nil, false
	}
}

func reasonCodesFromFixtureYAML(raw []byte) ([]string, error) {
	var fixtures []FakeFixture
	if err := yaml.Unmarshal(raw, &fixtures); err != nil {
		return nil, err
	}
	var codes []string
	for _, fx := range fixtures {
		if fx.Conclusion == nil {
			continue
		}
		codes = append(codes, fx.Conclusion.ReasonCodes...)
	}
	return codes, nil
}
