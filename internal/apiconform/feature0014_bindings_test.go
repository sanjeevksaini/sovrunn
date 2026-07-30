package apiconform

import (
	"path/filepath"
	"reflect"
	"sync"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/apischema"
	"github.com/sanjeevksaini/sovrunn/internal/resources"
)

func TestFeature0014TypeBindingsRegisterAndVerify(t *testing.T) {
	t.Parallel()

	root := moduleRoot(t)
	if findings := CheckFeature0014TypeBindings(root); len(findings) > 0 {
		t.Fatalf("CheckFeature0014TypeBindings: %v", findings)
	}
	if findings := CheckFeature0014CanonicalShape(root); len(findings) > 0 {
		t.Fatalf("CheckFeature0014CanonicalShape: %v", findings)
	}

	for _, binding := range Feature0014TypeBindings {
		binding := binding
		t.Run(binding.SchemaPath, func(t *testing.T) {
			t.Parallel()
			if binding.GoType == nil {
				t.Fatal("GoType must not be nil")
			}
			schemaPath := filepath.Join(root, filepath.FromSlash(binding.SchemaPath))
			schema, err := readRepoFile(schemaPath)
			if err != nil {
				t.Fatalf("read schema: %v", err)
			}
			if issues := apischema.ValidateSchemaSupport(schema); len(issues) > 0 {
				t.Fatalf("ValidateSchemaSupport: %v", issues)
			}
			if issues := apischema.VerifyGoTypeAgainstSchema(schema, binding.GoType); len(issues) > 0 {
				t.Fatalf("VerifyGoTypeAgainstSchema: %v", issues)
			}
		})
	}
}

func TestFeature0014KindInventoryExactlyFive(t *testing.T) {
	t.Parallel()

	if len(Feature0014KindInventory) != 5 {
		t.Fatalf("Feature0014KindInventory len=%d want 5", len(Feature0014KindInventory))
	}
	if len(Feature0014TypeBindings) != 5 {
		t.Fatalf("Feature0014TypeBindings len=%d want 5", len(Feature0014TypeBindings))
	}
	if err := AssertFeature0014KindInventory(Feature0014KindInventory); err != nil {
		t.Fatalf("canonical inventory must pass: %v", err)
	}

	wantKinds := []string{
		resources.KindProvider,
		resources.KindProviderLocation,
		resources.KindProviderDatacenter,
		resources.KindDatacenterFailureDomain,
		resources.KindInfrastructureStack,
	}
	wantCollections := []string{
		"providers",
		"provider-locations",
		"provider-datacenters",
		"datacenter-failure-domains",
		"infrastructure-stacks",
	}
	for i, e := range Feature0014KindInventory {
		if e.Kind != wantKinds[i] {
			t.Fatalf("inventory[%d].Kind=%q want %q", i, e.Kind, wantKinds[i])
		}
		if e.Collection != wantCollections[i] {
			t.Fatalf("inventory[%d].Collection=%q want %q", i, e.Collection, wantCollections[i])
		}
	}
}

func TestFeature0014KindInventoryRejectsAliasAlternateAndSuperseded(t *testing.T) {
	t.Parallel()

	base := append([]Feature0014KindEntry(nil), Feature0014KindInventory...)

	t.Run("sixth_kind", func(t *testing.T) {
		t.Parallel()
		got := append(append([]Feature0014KindEntry(nil), base...), Feature0014KindEntry{
			Kind:       "ProviderAlias",
			Collection: "provider-aliases",
		})
		if err := AssertFeature0014KindInventory(got); err == nil {
			t.Fatal("expected sixth/alias kind to fail inventory assertion")
		}
	})

	t.Run("alternate_location_term", func(t *testing.T) {
		t.Parallel()
		got := append([]Feature0014KindEntry(nil), base...)
		got[1] = Feature0014KindEntry{Kind: "ProviderRegion", Collection: "provider-regions"}
		if err := AssertFeature0014KindInventory(got); err == nil {
			t.Fatal("expected alternate location term ProviderRegion to fail inventory assertion")
		}
	})

	t.Run("superseded_stack_kind_name", func(t *testing.T) {
		t.Parallel()
		got := append([]Feature0014KindEntry(nil), base...)
		got[4] = Feature0014KindEntry{Kind: "IaaSStack", Collection: "iaas-stacks"}
		if err := AssertFeature0014KindInventory(got); err == nil {
			t.Fatal("expected superseded stack-kind name IaaSStack to fail inventory assertion")
		}
	})

	t.Run("wrong_collection", func(t *testing.T) {
		t.Parallel()
		got := append([]Feature0014KindEntry(nil), base...)
		got[0] = Feature0014KindEntry{Kind: resources.KindProvider, Collection: "provider"}
		if err := AssertFeature0014KindInventory(got); err == nil {
			t.Fatal("expected wrong collection to fail inventory assertion")
		}
	})
}

func TestFeature0014ProviderNeutrality(t *testing.T) {
	t.Parallel()

	findings := CheckFeature0014ProviderNeutrality(moduleRoot(t))
	if len(findings) > 0 {
		t.Fatalf("CheckFeature0014ProviderNeutrality: %v", findings)
	}
}

func TestFeature0014ProviderNeutralityRejectsEndpointFields(t *testing.T) {
	t.Parallel()

	for _, name := range []string{
		"endpoint",
		"endpoints",
		"apiEndpoint",
		"providerEndpoint",
		"managementEndpoints",
		"baseURL",
		"serviceUrl",
	} {
		name := name
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			if !isFeature0014EndpointField(name) {
				t.Fatalf("endpoint-bearing field %q was accepted", name)
			}
		})
	}

	for _, name := range []string{"technology", "countryCode", "providerLocationRef"} {
		if isFeature0014EndpointField(name) {
			t.Fatalf("ordinary FEATURE-0014 field %q was rejected as an endpoint", name)
		}
	}

	type nestedEndpoint struct {
		APIEndpoint string `json:"apiEndpoint"`
	}
	type contractWithEndpoint struct {
		Spec nestedEndpoint `json:"spec"`
	}
	if err := checkGoTypeNoFeature0014EndpointFields(reflect.TypeOf(contractWithEndpoint{})); err == nil {
		t.Fatal("recursive Go-type scan accepted apiEndpoint")
	}

	type nestedURL struct {
		BaseURL string `json:"baseURL"`
	}
	type contractWithURL struct {
		Items []nestedURL `json:"items"`
	}
	if err := checkGoTypeNoFeature0014EndpointFields(reflect.TypeOf(contractWithURL{})); err == nil {
		t.Fatal("recursive Go-type scan accepted baseURL")
	}
}

func TestFeature0014BindingsLoadViaSchemaRegistry(t *testing.T) {
	t.Parallel()

	root := moduleRoot(t)
	reg, err := NewRepositorySchemaRegistry(filepath.Join(root, CanonicalSchemasDir))
	if err != nil {
		t.Fatalf("NewRepositorySchemaRegistry: %v", err)
	}
	for _, binding := range Feature0014TypeBindings {
		body, err := reg.Load(binding.SchemaPath)
		if err != nil {
			t.Fatalf("Load(%s): %v", binding.SchemaPath, err)
		}
		if len(body) == 0 {
			t.Fatalf("Load(%s) returned empty schema", binding.SchemaPath)
		}
	}
}

func TestFeature0014BindingsIsolatedFromTypeBindings(t *testing.T) {
	t.Parallel()

	f14 := make(map[string]struct{}, len(Feature0014TypeBindings))
	for _, b := range Feature0014TypeBindings {
		f14[b.SchemaPath] = struct{}{}
		if b.GoType.PkgPath() != "github.com/sanjeevksaini/sovrunn/internal/resources" {
			t.Fatalf("%s GoType package=%q want internal/resources", b.SchemaPath, b.GoType.PkgPath())
		}
	}
	for _, b := range TypeBindings {
		if _, ok := f14[b.SchemaPath]; ok {
			t.Fatalf("FEATURE-0014 schema %q must not be registered in TypeBindings (isolation)", b.SchemaPath)
		}
	}

	// resources must not import apiconform (one-way binding direction).
	imports := packageImports(t, filepath.Join(moduleRoot(t), "internal", "resources"))
	for _, imp := range imports {
		if imp.path == modulePath+"/internal/apiconform" {
			t.Fatalf("internal/resources must not import apiconform (got %s)", imp.file)
		}
	}
}

func TestFeature0014BindingsConcurrentSafe(t *testing.T) {
	t.Parallel()

	root := moduleRoot(t)
	const goroutines = 8
	var wg sync.WaitGroup
	errCh := make(chan error, goroutines*3)
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			if err := AssertFeature0014KindInventory(Feature0014KindInventory); err != nil {
				errCh <- err
				return
			}
			if findings := CheckFeature0014TypeBindings(root); len(findings) > 0 {
				errCh <- errString("type bindings", findings)
				return
			}
			if findings := CheckFeature0014ProviderNeutrality(root); len(findings) > 0 {
				errCh <- errString("provider neutrality", findingsToStrings(findings))
			}
		}()
	}
	wg.Wait()
	close(errCh)
	for err := range errCh {
		if err != nil {
			t.Fatalf("concurrent FEATURE-0014 binding check failed: %v", err)
		}
	}

	// Touch bindings under concurrent readers to confirm immutable registry shape.
	var sum int
	var mu sync.Mutex
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			n := 0
			for _, b := range Feature0014TypeBindings {
				if b.GoType != nil && b.GoType.Kind() == reflect.Struct {
					n++
				}
			}
			mu.Lock()
			sum += n
			mu.Unlock()
		}()
	}
	wg.Wait()
	if sum != goroutines*len(Feature0014TypeBindings) {
		t.Fatalf("concurrent binding reads sum=%d want %d", sum, goroutines*len(Feature0014TypeBindings))
	}
}

func errString(label string, findings []string) error {
	return &feature0014CheckError{label: label, findings: findings}
}

type feature0014CheckError struct {
	label    string
	findings []string
}

func (e *feature0014CheckError) Error() string {
	return e.label + ": " + joinFindings(e.findings)
}

func findingsToStrings(findings []FitnessFinding) []string {
	out := make([]string, 0, len(findings))
	for _, f := range findings {
		out = append(out, f.String())
	}
	return out
}

func joinFindings(findings []string) string {
	if len(findings) == 0 {
		return ""
	}
	out := findings[0]
	for i := 1; i < len(findings); i++ {
		out += "; " + findings[i]
	}
	return out
}
