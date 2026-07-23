package apiconform

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/apischema"
)

// Deterministic seed for Property 7 reproducibility
// (F12-SCOPE-001, F12-SEC-006, F12-BOUNDARY-001, F12-REF-003).
const property7Seed int64 = 20260723

const property7Iterations = 100

// Core grammar primitives named by Property 7 / design D-02 neutrality scope.
var property7CoreGrammarPackages = []string{
	"apimeta",
	"apiref",
	"apicond",
	"apiproblem",
	"apivalid",
}

// Canonical schemas whose boundary is core/customer (not adapter/plugin).
// Provider-native identifiers/fields are prohibited here.
var property7CoreSchemaBindings = []struct {
	SchemaPath string
	GoType     reflect.Type
	Boundary   apimeta.Boundary
}{
	{SchemaPath: "api/schemas/project.json", GoType: reflect.TypeOf(Project{}), Boundary: apimeta.BoundaryCustomerFacing},
	{SchemaPath: "api/schemas/resource-pool.json", GoType: reflect.TypeOf(ResourcePool{}), Boundary: apimeta.BoundaryOperatorFacing},
	{SchemaPath: "api/schemas/placement-evaluation-request.json", GoType: reflect.TypeOf(PlacementEvaluationRequest{}), Boundary: apimeta.BoundaryInternalEngineFacing},
	{SchemaPath: "api/schemas/audit-event.json", GoType: reflect.TypeOf(AuditEvent{}), Boundary: apimeta.BoundaryGovernanceOnly},
	{SchemaPath: "api/schemas/_common/type-meta.json", GoType: reflect.TypeOf(apimeta.TypeMeta{}), Boundary: ""},
	{SchemaPath: "api/schemas/_common/object-meta.json", GoType: reflect.TypeOf(apimeta.ObjectMeta{}), Boundary: ""},
	{SchemaPath: "api/schemas/_common/typed-ref.json", GoType: reflect.TypeOf(apimeta.TypedRef{}), Boundary: ""},
	{SchemaPath: "api/schemas/_common/scope-ref.json", GoType: reflect.TypeOf(apimeta.ScopeRef{}), Boundary: ""},
	{SchemaPath: "api/schemas/_common/owner-ref.json", GoType: reflect.TypeOf(apimeta.OwnerRef{}), Boundary: ""},
}

// Adapter/plugin schemas may isolate opaque native references; they must not
// embed provider SDK types in derivative Go contracts.
var property7AdapterPluginBindings = []struct {
	SchemaPath string
	GoType     reflect.Type
	Boundary   apimeta.Boundary
}{
	{SchemaPath: "api/schemas/discovered-database.json", GoType: reflect.TypeOf(DiscoveredDatabase{}), Boundary: apimeta.BoundaryAdapterFacing},
	{SchemaPath: "api/schemas/adapter-configuration.json", GoType: reflect.TypeOf(AdapterConfiguration{}), Boundary: apimeta.BoundaryAdapterFacing},
	{SchemaPath: "api/schemas/plugin-definition.json", GoType: reflect.TypeOf(PluginDefinition{}), Boundary: apimeta.BoundaryPluginFacing},
	{SchemaPath: "api/schemas/operation.json", GoType: reflect.TypeOf(Operation{}), Boundary: apimeta.BoundaryPluginFacing},
}

// Provider-native property/field identifier tokens banned in core schemas.
// Matched against lowercased property names with non-alphanumerics stripped
// so camelCase forms like awsAccountId / azureSubscriptionId are caught.
// Descriptions are not scanned (avoid false positives on "provider-neutral").
var property7BannedCoreFieldTokens = []string{
	"aws",
	"amazon",
	"azure",
	"gcp",
	"gke",
	"eks",
	"aks",
	"kubernetes",
	"k8s",
	"arn",
	"subscriptionid",
	"resourcegroup",
	"vpcid",
	"instanceid",
	"cloudformation",
	"armtemplate",
	"providerid",
	"providerarn",
	"nativeconfig",
	"nativeconfigref",
	"adapterclass",
}

// Opaque native-isolation field names allowed only behind adapter/plugin.
var property7AllowedAdapterNativeFields = map[string]struct{}{
	"nativeconfigref": {},
	"adapterclass":    {},
	"externalname":    {},
}

type property7Scenario string

const (
	property7CoreGrammar            property7Scenario = "core_grammar_imports"
	property7CoreSchema             property7Scenario = "core_schema_neutral"
	property7AdapterPluginSchema    property7Scenario = "adapter_plugin_isolation"
	property7SyntheticCoreViolation property7Scenario = "synthetic_core_violation"
	property7SyntheticSDKImport     property7Scenario = "synthetic_sdk_import"
	property7SyntheticAdapterOK     property7Scenario = "synthetic_adapter_ok"
)

type property7Case struct {
	Scenario   property7Scenario
	WantAccept bool

	// Grammar import checks.
	PackageName string
	ImportPaths []string

	// Schema / Go-type checks.
	SchemaPath string
	Boundary   apimeta.Boundary
	Schema     []byte
	GoType     reflect.Type
	PropNames  []string
}

// Feature: api-resource-naming-status-and-validation-standard, Property 7: Provider neutrality of core
//
// For any core or customer-facing schema and any core primitive
// (apimeta/apiref/apicond/apiproblem/apivalid), no provider-native
// identifier, provider SDK type, or provider-specific field is present;
// provider-native data is expressible only behind adapter-facing/
// plugin-facing schemas. The provider-neutrality check fails on any
// violation.
//
// Validates: Requirements 4.4, 4.5, 4.6, 7.6 (F12-SCOPE-001, F12-SEC-006, F12-BOUNDARY-001)
func TestProperty7_ProviderNeutralityOfCore(t *testing.T) {
	t.Parallel()

	root := moduleRoot(t)
	coreSchemas := loadProperty7Schemas(t, root, property7CoreSchemaBindings)
	adapterSchemas := loadProperty7Schemas(t, root, property7AdapterPluginBindings)
	grammarImports := loadProperty7GrammarImports(t, root)

	rng := rand.New(rand.NewSource(property7Seed))
	for i := 0; i < property7Iterations; i++ {
		c := generateProperty7Case(rng, i, grammarImports, coreSchemas, adapterSchemas)
		if err := checkProperty7Case(c, i); err != nil {
			t.Fatalf("property 7 failed at iteration %d (seed %d scenario %s): %v",
				i, property7Seed, c.Scenario, err)
		}
	}
}

func loadProperty7GrammarImports(t *testing.T, root string) map[string][]string {
	t.Helper()

	out := make(map[string][]string, len(property7CoreGrammarPackages))
	for _, pkg := range property7CoreGrammarPackages {
		imports := packageImports(t, filepath.Join(root, "internal", pkg))
		paths := make([]string, 0, len(imports))
		for _, imp := range imports {
			paths = append(paths, imp.path)
		}
		out[pkg] = paths
	}
	return out
}

func loadProperty7Schemas(t *testing.T, root string, bindings []struct {
	SchemaPath string
	GoType     reflect.Type
	Boundary   apimeta.Boundary
}) map[string][]byte {
	t.Helper()

	out := make(map[string][]byte, len(bindings))
	for _, b := range bindings {
		path := filepath.Join(root, filepath.FromSlash(b.SchemaPath))
		raw, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("read schema %s: %v", path, err)
		}
		if issues := apischema.ValidateSchemaSupport(raw); len(issues) > 0 {
			t.Fatalf("ValidateSchemaSupport failed for %s: %v", b.SchemaPath, issues)
		}
		if b.Boundary != "" {
			meta, issues := apischema.ReadAnnotations(raw)
			if len(issues) > 0 {
				t.Fatalf("ReadAnnotations failed for %s: %v", b.SchemaPath, issues)
			}
			if meta.Boundary != b.Boundary {
				t.Fatalf("%s: boundary=%q, want %q", b.SchemaPath, meta.Boundary, b.Boundary)
			}
		}
		out[b.SchemaPath] = raw
	}
	return out
}

func generateProperty7Case(
	rng *rand.Rand,
	iteration int,
	grammarImports map[string][]string,
	coreSchemas map[string][]byte,
	adapterSchemas map[string][]byte,
) property7Case {
	bucket := iteration % 6
	if rng.Intn(20) == 0 {
		bucket = rng.Intn(6)
	}

	switch bucket {
	case 0:
		pkg := property7CoreGrammarPackages[iteration%len(property7CoreGrammarPackages)]
		if rng.Intn(3) == 0 {
			pkg = property7CoreGrammarPackages[rng.Intn(len(property7CoreGrammarPackages))]
		}
		return property7Case{
			Scenario:    property7CoreGrammar,
			WantAccept:  true,
			PackageName: pkg,
			ImportPaths: append([]string(nil), grammarImports[pkg]...),
		}
	case 1:
		b := property7CoreSchemaBindings[iteration%len(property7CoreSchemaBindings)]
		if rng.Intn(3) == 0 {
			b = property7CoreSchemaBindings[rng.Intn(len(property7CoreSchemaBindings))]
		}
		return property7Case{
			Scenario:   property7CoreSchema,
			WantAccept: true,
			SchemaPath: b.SchemaPath,
			Boundary:   b.Boundary,
			Schema:     coreSchemas[b.SchemaPath],
			GoType:     b.GoType,
			PropNames:  collectSchemaPropertyNames(coreSchemas[b.SchemaPath]),
		}
	case 2:
		b := property7AdapterPluginBindings[iteration%len(property7AdapterPluginBindings)]
		if rng.Intn(3) == 0 {
			b = property7AdapterPluginBindings[rng.Intn(len(property7AdapterPluginBindings))]
		}
		return property7Case{
			Scenario:   property7AdapterPluginSchema,
			WantAccept: true,
			SchemaPath: b.SchemaPath,
			Boundary:   b.Boundary,
			Schema:     adapterSchemas[b.SchemaPath],
			GoType:     b.GoType,
			PropNames:  collectSchemaPropertyNames(adapterSchemas[b.SchemaPath]),
		}
	case 3:
		token := property7BannedCoreFieldTokens[rng.Intn(len(property7BannedCoreFieldTokens))]
		field := property7SyntheticFieldName(token, rng)
		return property7Case{
			Scenario:   property7SyntheticCoreViolation,
			WantAccept: false,
			SchemaPath: "synthetic/core-violation.json",
			Boundary:   apimeta.BoundaryCustomerFacing,
			PropNames:  []string{field},
			GoType:     nil,
		}
	case 4:
		sdk := providerSDKPrefixes[rng.Intn(len(providerSDKPrefixes))]
		return property7Case{
			Scenario:    property7SyntheticSDKImport,
			WantAccept:  false,
			PackageName: "apimeta",
			ImportPaths: []string{"fmt", sdk + "fake-sdk"},
		}
	default:
		field := "nativeConfigRef"
		if rng.Intn(2) == 0 {
			field = "adapterClass"
		}
		return property7Case{
			Scenario:   property7SyntheticAdapterOK,
			WantAccept: true,
			SchemaPath: "synthetic/adapter-ok.json",
			Boundary:   apimeta.BoundaryAdapterFacing,
			PropNames:  []string{field},
			GoType:     nil,
		}
	}
}

func checkProperty7Case(c property7Case, iteration int) error {
	var err error
	switch c.Scenario {
	case property7CoreGrammar, property7SyntheticSDKImport:
		err = checkProperty7Imports(c.PackageName, c.ImportPaths)
	case property7CoreSchema, property7SyntheticCoreViolation:
		err = checkProperty7CoreSchema(c.Boundary, c.PropNames, c.GoType)
	case property7AdapterPluginSchema, property7SyntheticAdapterOK:
		err = checkProperty7AdapterPluginSchema(c.Boundary, c.PropNames, c.GoType)
	default:
		return fmt.Errorf("unknown scenario %q", c.Scenario)
	}

	accepted := err == nil
	if accepted != c.WantAccept {
		if c.WantAccept {
			return fmt.Errorf("iteration %d: expected accept, got error: %v", iteration, err)
		}
		return fmt.Errorf("iteration %d: expected provider-neutrality violation, got accept", iteration)
	}
	return nil
}

func checkProperty7Imports(pkg string, imports []string) error {
	for _, path := range imports {
		if isProviderSDK(path) {
			return fmt.Errorf("%s imports provider SDK %q (provider neutrality violated)", pkg, path)
		}
	}
	return nil
}

func checkProperty7CoreSchema(boundary apimeta.Boundary, props []string, goType reflect.Type) error {
	if boundary == apimeta.BoundaryAdapterFacing || boundary == apimeta.BoundaryPluginFacing {
		return fmt.Errorf("core check must not classify boundary %q as core", boundary)
	}
	for _, name := range props {
		if token := property7BannedCoreFieldToken(name); token != "" {
			return fmt.Errorf("core schema property %q embeds provider-native token %q", name, token)
		}
	}
	if goType != nil {
		if err := checkProperty7GoTypeNoSDK(goType); err != nil {
			return err
		}
		if err := checkProperty7GoTypeNoBannedFields(goType, true); err != nil {
			return err
		}
	}
	return nil
}

func checkProperty7AdapterPluginSchema(boundary apimeta.Boundary, props []string, goType reflect.Type) error {
	if boundary != apimeta.BoundaryAdapterFacing && boundary != apimeta.BoundaryPluginFacing {
		return fmt.Errorf("adapter/plugin check requires adapter-facing or plugin-facing, got %q", boundary)
	}
	// Adapter/plugin may carry opaque native-isolation fields, but must not
	// introduce raw vendor SDK property shapes beyond the approved opaque set.
	for _, name := range props {
		norm := property7NormalizeIdent(name)
		if _, ok := property7AllowedAdapterNativeFields[norm]; ok {
			continue
		}
		if token := property7BannedCoreFieldToken(name); token != "" {
			// Still reject explicit cloud-vendor identifiers even behind adapter
			// boundaries when they are not the approved opaque isolation fields.
			switch token {
			case "aws", "amazon", "azure", "gcp", "gke", "eks", "aks",
				"kubernetes", "k8s", "arn", "subscriptionid", "resourcegroup",
				"vpcid", "instanceid", "cloudformation", "armtemplate",
				"providerid", "providerarn":
				return fmt.Errorf("adapter/plugin schema property %q embeds provider-native token %q", name, token)
			}
		}
	}
	if goType != nil {
		if err := checkProperty7GoTypeNoSDK(goType); err != nil {
			return err
		}
		// Field names may include approved opaque refs; SDK embeds still banned.
		if err := checkProperty7GoTypeNoBannedFields(goType, false); err != nil {
			return err
		}
	}
	return nil
}

func checkProperty7GoTypeNoSDK(t reflect.Type) error {
	seen := map[reflect.Type]struct{}{}
	var walk func(reflect.Type) error
	walk = func(tt reflect.Type) error {
		if tt == nil {
			return nil
		}
		for tt.Kind() == reflect.Pointer {
			tt = tt.Elem()
		}
		if _, ok := seen[tt]; ok {
			return nil
		}
		seen[tt] = struct{}{}

		if pkg := tt.PkgPath(); pkg != "" && isProviderSDK(pkg) {
			return fmt.Errorf("go type %s embeds provider SDK package %q", tt.String(), pkg)
		}

		switch tt.Kind() {
		case reflect.Struct:
			for i := 0; i < tt.NumField(); i++ {
				f := tt.Field(i)
				if err := walk(f.Type); err != nil {
					return err
				}
			}
		case reflect.Slice, reflect.Array, reflect.Chan:
			return walk(tt.Elem())
		case reflect.Map:
			if err := walk(tt.Key()); err != nil {
				return err
			}
			return walk(tt.Elem())
		}
		return nil
	}
	return walk(t)
}

func checkProperty7GoTypeNoBannedFields(t reflect.Type, coreStrict bool) error {
	seen := map[reflect.Type]struct{}{}
	var walk func(reflect.Type) error
	walk = func(tt reflect.Type) error {
		if tt == nil {
			return nil
		}
		for tt.Kind() == reflect.Pointer {
			tt = tt.Elem()
		}
		if _, ok := seen[tt]; ok {
			return nil
		}
		seen[tt] = struct{}{}
		if tt.Kind() != reflect.Struct {
			switch tt.Kind() {
			case reflect.Slice, reflect.Array:
				return walk(tt.Elem())
			case reflect.Map:
				if err := walk(tt.Key()); err != nil {
					return err
				}
				return walk(tt.Elem())
			}
			return nil
		}
		for i := 0; i < tt.NumField(); i++ {
			f := tt.Field(i)
			name := f.Name
			if tag := f.Tag.Get("json"); tag != "" {
				if n, _, _ := strings.Cut(tag, ","); n != "" && n != "-" {
					name = n
				}
			}
			norm := property7NormalizeIdent(name)
			if coreStrict {
				if token := property7BannedCoreFieldToken(name); token != "" {
					return fmt.Errorf("core Go field %q embeds provider-native token %q", name, token)
				}
			} else {
				if _, ok := property7AllowedAdapterNativeFields[norm]; !ok {
					if token := property7BannedCoreFieldToken(name); token != "" {
						switch token {
						case "aws", "amazon", "azure", "gcp", "gke", "eks", "aks",
							"kubernetes", "k8s", "arn", "subscriptionid", "resourcegroup",
							"vpcid", "instanceid", "cloudformation", "armtemplate",
							"providerid", "providerarn":
							return fmt.Errorf("adapter/plugin Go field %q embeds provider-native token %q", name, token)
						}
					}
				}
			}
			if err := walk(f.Type); err != nil {
				return err
			}
		}
		return nil
	}
	return walk(t)
}

func property7BannedCoreFieldToken(name string) string {
	norm := property7NormalizeIdent(name)
	for _, token := range property7BannedCoreFieldTokens {
		if norm == token || strings.Contains(norm, token) {
			return token
		}
	}
	return ""
}

func property7NormalizeIdent(s string) string {
	var b strings.Builder
	for _, r := range strings.ToLower(s) {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			b.WriteRune(r)
		}
	}
	return b.String()
}

func property7SyntheticFieldName(token string, rng *rand.Rand) string {
	switch token {
	case "aws":
		return "awsAccountId"
	case "amazon":
		return "amazonResourceName"
	case "azure":
		return "azureSubscriptionId"
	case "gcp":
		return "gcpProjectNumber"
	case "gke", "eks", "aks":
		return token + "ClusterName"
	case "kubernetes":
		return "kubernetesClusterID"
	case "k8s":
		return "k8sNamespace"
	case "arn":
		return "arn"
	case "subscriptionid":
		return "subscriptionId"
	case "resourcegroup":
		return "resourceGroupName"
	case "vpcid":
		return "vpcId"
	case "instanceid":
		return "instanceId"
	case "cloudformation":
		return "cloudFormationStack"
	case "armtemplate":
		return "armTemplateId"
	case "providerid":
		return "providerId"
	case "providerarn":
		return "providerArn"
	case "nativeconfig":
		return "nativeConfig"
	case "nativeconfigref":
		return "nativeConfigRef"
	case "adapterclass":
		return "adapterClass"
	default:
		return token + fmt.Sprintf("Field%d", rng.Intn(1000))
	}
}

func collectSchemaPropertyNames(schema []byte) []string {
	var doc any
	if err := json.Unmarshal(schema, &doc); err != nil {
		return nil
	}
	var out []string
	var walk func(any)
	walk = func(node any) {
		obj, ok := node.(map[string]any)
		if !ok {
			return
		}
		if props, ok := obj["properties"].(map[string]any); ok {
			for name, sub := range props {
				out = append(out, name)
				walk(sub)
			}
		}
		if items, ok := obj["items"]; ok {
			walk(items)
		}
		if addl, ok := obj["additionalProperties"]; ok {
			walk(addl)
		}
	}
	walk(doc)
	return out
}
