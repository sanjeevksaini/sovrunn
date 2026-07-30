package apiconform

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"

	"github.com/sanjeevksaini/sovrunn/internal/apicond"
	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/apischema"
	"github.com/sanjeevksaini/sovrunn/internal/resources"
)

// Feature0014TypeBindings registers the five FEATURE-0014 canonical schemas
// against their concrete Go contract types (design §3.1 I-6; DD-07; F14-REQ-21).
// Bindings stay in apiconform so internal/resources never imports this package.
// This registry is intentionally separate from TypeBindings (FEATURE-0012/0013).
var Feature0014TypeBindings = []apischema.TypeBinding{
	{SchemaPath: "api/schemas/provider.json", GoType: reflect.TypeOf(resources.Provider{})},
	{SchemaPath: "api/schemas/provider-location.json", GoType: reflect.TypeOf(resources.ProviderLocation{})},
	{SchemaPath: "api/schemas/provider-datacenter.json", GoType: reflect.TypeOf(resources.ProviderDatacenter{})},
	{SchemaPath: "api/schemas/datacenter-failure-domain.json", GoType: reflect.TypeOf(resources.DatacenterFailureDomain{})},
	{SchemaPath: "api/schemas/infrastructure-stack.json", GoType: reflect.TypeOf(resources.InfrastructureStack{})},
}

// Feature0014KindEntry is one kind and its HTTP collection (design §5).
type Feature0014KindEntry struct {
	Kind       string
	Collection string
}

// Feature0014KindInventory is the closed FEATURE-0014 kind/collection set
// (F14-REQ-01, F14-REQ-02, F14-REQ-03; F14-AD-001, F14-AD-006, F14-AD-021).
// Exactly five entries; no alias kinds.
var Feature0014KindInventory = []Feature0014KindEntry{
	{Kind: resources.KindProvider, Collection: "providers"},
	{Kind: resources.KindProviderLocation, Collection: "provider-locations"},
	{Kind: resources.KindProviderDatacenter, Collection: "provider-datacenters"},
	{Kind: resources.KindDatacenterFailureDomain, Collection: "datacenter-failure-domains"},
	{Kind: resources.KindInfrastructureStack, Collection: "infrastructure-stacks"},
}

// expectedFeature0014Annotations is the ManagedResource / operator-facing
// contract shared by all five FEATURE-0014 schemas (F14-REQ-21, F14-REQ-22).
var expectedFeature0014Annotations = map[string]struct {
	kind      string
	scopes    []apimeta.ScopeKind
	stability apimeta.Stability
}{
	"api/schemas/provider.json": {
		kind:      resources.KindProvider,
		scopes:    []apimeta.ScopeKind{apimeta.ScopePlatform, apimeta.ScopeOrganization},
		stability: apimeta.StabilityAlpha,
	},
	"api/schemas/provider-location.json": {
		kind:      resources.KindProviderLocation,
		scopes:    []apimeta.ScopeKind{apimeta.ScopeProvider},
		stability: apimeta.StabilityAlpha,
	},
	"api/schemas/provider-datacenter.json": {
		kind:      resources.KindProviderDatacenter,
		scopes:    []apimeta.ScopeKind{apimeta.ScopeProvider},
		stability: apimeta.StabilityAlpha,
	},
	"api/schemas/datacenter-failure-domain.json": {
		kind:      resources.KindDatacenterFailureDomain,
		scopes:    []apimeta.ScopeKind{apimeta.ScopeProvider},
		stability: apimeta.StabilityAlpha,
	},
	"api/schemas/infrastructure-stack.json": {
		kind:      resources.KindInfrastructureStack,
		scopes:    []apimeta.ScopeKind{apimeta.ScopeProvider},
		stability: apimeta.StabilityAlpha,
	},
}

// AssertFeature0014KindInventory reports whether got matches the closed
// five-kind inventory exactly (kinds, collections, and count). Unknown,
// alias, or alternate kinds fail.
func AssertFeature0014KindInventory(got []Feature0014KindEntry) error {
	if len(got) != len(Feature0014KindInventory) {
		return fmt.Errorf("FEATURE-0014 kind inventory count=%d want exactly %d",
			len(got), len(Feature0014KindInventory))
	}

	wantByKind := make(map[string]string, len(Feature0014KindInventory))
	for _, e := range Feature0014KindInventory {
		wantByKind[e.Kind] = e.Collection
	}

	seen := make(map[string]struct{}, len(got))
	for _, e := range got {
		if e.Kind == "" {
			return fmt.Errorf("FEATURE-0014 kind inventory entry has empty kind")
		}
		if e.Collection == "" {
			return fmt.Errorf("FEATURE-0014 kind %q has empty collection", e.Kind)
		}
		wantCollection, ok := wantByKind[e.Kind]
		if !ok {
			return fmt.Errorf("FEATURE-0014 kind inventory rejects unknown or alias kind %q", e.Kind)
		}
		if e.Collection != wantCollection {
			return fmt.Errorf("FEATURE-0014 kind %q collection=%q want %q",
				e.Kind, e.Collection, wantCollection)
		}
		if _, dup := seen[e.Kind]; dup {
			return fmt.Errorf("FEATURE-0014 kind inventory has duplicate kind %q", e.Kind)
		}
		seen[e.Kind] = struct{}{}
	}

	if len(seen) != len(Feature0014KindInventory) {
		return fmt.Errorf("FEATURE-0014 kind inventory resolved %d distinct kinds want %d",
			len(seen), len(Feature0014KindInventory))
	}
	return nil
}

// CheckFeature0014TypeBindings runs ValidateSchemaSupport and
// VerifyGoTypeAgainstSchema for every FEATURE-0014 TypeBinding.
func CheckFeature0014TypeBindings(moduleRoot string) []string {
	if strings.TrimSpace(moduleRoot) == "" {
		return []string{"module root is empty"}
	}
	var findings []string
	seen := make(map[string]struct{}, len(Feature0014TypeBindings))
	for _, binding := range Feature0014TypeBindings {
		if binding.SchemaPath == "" {
			findings = append(findings, "TypeBinding SchemaPath must not be empty")
			continue
		}
		if binding.GoType == nil {
			findings = append(findings, fmt.Sprintf("%s: GoType must not be nil", binding.SchemaPath))
			continue
		}
		if _, dup := seen[binding.SchemaPath]; dup {
			findings = append(findings, fmt.Sprintf("duplicate TypeBinding for %q", binding.SchemaPath))
			continue
		}
		seen[binding.SchemaPath] = struct{}{}

		schemaPath := filepath.Join(moduleRoot, filepath.FromSlash(binding.SchemaPath))
		schema, err := os.ReadFile(schemaPath) // #nosec G304 -- trusted repo-local schema path
		if err != nil {
			findings = append(findings, fmt.Sprintf("%s: read: %v", binding.SchemaPath, err))
			continue
		}
		if issues := apischema.ValidateSchemaSupport(schema); len(issues) > 0 {
			findings = append(findings, fmt.Sprintf("%s: schema support: %v", binding.SchemaPath, issues))
			continue
		}
		if issues := apischema.VerifyGoTypeAgainstSchema(schema, binding.GoType); len(issues) > 0 {
			for _, issue := range issues {
				findings = append(findings, fmt.Sprintf("%s: %s %s: %s",
					binding.SchemaPath, issue.Path, issue.Code, issue.Message))
			}
		}
	}
	if len(Feature0014TypeBindings) != len(Feature0014KindInventory) {
		findings = append(findings, fmt.Sprintf(
			"FEATURE-0014 TypeBindings count=%d want %d",
			len(Feature0014TypeBindings), len(Feature0014KindInventory)))
	}
	if len(seen) != len(Feature0014KindInventory) {
		findings = append(findings, fmt.Sprintf(
			"FEATURE-0014 TypeBindings distinct=%d want %d",
			len(seen), len(Feature0014KindInventory)))
	}
	return findings
}

// CheckFeature0014CanonicalShape validates FEATURE-0012 ManagedResource
// annotations and status/condition grammar on the five bindings (F14-REQ-21,
// F14-REQ-22).
func CheckFeature0014CanonicalShape(moduleRoot string) []string {
	if strings.TrimSpace(moduleRoot) == "" {
		return []string{"module root is empty"}
	}
	var findings []string
	for _, binding := range Feature0014TypeBindings {
		want, ok := expectedFeature0014Annotations[binding.SchemaPath]
		if !ok {
			findings = append(findings, fmt.Sprintf("%s: missing expected annotation contract", binding.SchemaPath))
			continue
		}
		schemaPath := filepath.Join(moduleRoot, filepath.FromSlash(binding.SchemaPath))
		schema, err := os.ReadFile(schemaPath) // #nosec G304 -- trusted repo-local schema path
		if err != nil {
			findings = append(findings, fmt.Sprintf("%s: read: %v", binding.SchemaPath, err))
			continue
		}
		meta, issues := apischema.ReadAnnotations(schema)
		if len(issues) > 0 {
			findings = append(findings, fmt.Sprintf("%s: annotations: %v", binding.SchemaPath, issues))
			continue
		}
		if meta.Profile != apimeta.ProfileManagedResource {
			findings = append(findings, fmt.Sprintf("%s: profile=%q want ManagedResource",
				binding.SchemaPath, meta.Profile))
		}
		if meta.Boundary != apimeta.BoundaryOperatorFacing {
			findings = append(findings, fmt.Sprintf("%s: boundary=%q want operator-facing",
				binding.SchemaPath, meta.Boundary))
		}
		if meta.Stability != want.stability {
			findings = append(findings, fmt.Sprintf("%s: stability=%q want %q",
				binding.SchemaPath, meta.Stability, want.stability))
		}
		if !scopeKindsEqual(meta.AllowedScopes, want.scopes) {
			findings = append(findings, fmt.Sprintf("%s: allowed-scopes=%v want %v",
				binding.SchemaPath, meta.AllowedScopes, want.scopes))
		}
		if binding.GoType.Name() != want.kind {
			findings = append(findings, fmt.Sprintf("%s: GoType=%q want kind %q",
				binding.SchemaPath, binding.GoType.Name(), want.kind))
		}
		if err := assertFeature0014StatusConditionGrammar(binding.GoType); err != nil {
			findings = append(findings, fmt.Sprintf("%s: %v", binding.SchemaPath, err))
		}
		if err := assertFeature0014SchemaStatusConditions(schema); err != nil {
			findings = append(findings, fmt.Sprintf("%s: %v", binding.SchemaPath, err))
		}
	}
	return findings
}

// CheckFeature0014ProviderNeutrality reuses FEATURE-0012 provider-neutrality
// scanners against the five FEATURE-0014 contracts (F14-REQ-18, F14-AD-019).
func CheckFeature0014ProviderNeutrality(moduleRoot string) []FitnessFinding {
	var findings []FitnessFinding
	if strings.TrimSpace(moduleRoot) == "" {
		return []FitnessFinding{{
			Check:   FitnessCheckNoProviderSDKInCoreCustomer,
			Schema:  "FEATURE-0014",
			Path:    "/",
			Code:    CodeFitnessSchemaLoadFailed,
			Message: "module root is empty",
		}}
	}
	for _, binding := range Feature0014TypeBindings {
		schemaPath := filepath.Join(moduleRoot, filepath.FromSlash(binding.SchemaPath))
		body, err := os.ReadFile(schemaPath) // #nosec G304 -- trusted repo-local schema path
		if err != nil {
			findings = append(findings, FitnessFinding{
				Check:   FitnessCheckNoProviderSDKInCoreCustomer,
				Schema:  binding.SchemaPath,
				Path:    "/",
				Code:    CodeFitnessSchemaLoadFailed,
				Message: err.Error(),
			})
			continue
		}
		for _, prop := range collectSchemaPropertyNamesFitness(body) {
			if token := bannedCoreNativeFieldToken(prop); token != "" {
				findings = append(findings, FitnessFinding{
					Check:   FitnessCheckNoProviderSDKInCoreCustomer,
					Schema:  binding.SchemaPath,
					Path:    "/" + prop,
					Code:    CodeFitnessProviderNativeField,
					Message: fmt.Sprintf("property %q embeds provider-native token %q", prop, token),
				})
			}
		}
		if err := checkGoTypeNoProviderSDK(binding.GoType); err != nil {
			findings = append(findings, FitnessFinding{
				Check:   FitnessCheckNoProviderSDKInCoreCustomer,
				Schema:  binding.SchemaPath,
				Path:    "/",
				Code:    CodeFitnessProviderSDKGoType,
				Message: err.Error(),
			})
		}
		if err := checkGoTypeNoBannedNativeFields(binding.GoType); err != nil {
			findings = append(findings, FitnessFinding{
				Check:   FitnessCheckNoProviderSDKInCoreCustomer,
				Schema:  binding.SchemaPath,
				Path:    "/",
				Code:    CodeFitnessProviderNativeField,
				Message: err.Error(),
			})
		}
	}
	return findings
}

func assertFeature0014StatusConditionGrammar(goType reflect.Type) error {
	if goType == nil {
		return fmt.Errorf("Go type is nil")
	}
	for goType.Kind() == reflect.Pointer {
		goType = goType.Elem()
	}
	statusField, ok := goType.FieldByName("Status")
	if !ok {
		return fmt.Errorf("Go type %s missing Status field", goType.Name())
	}
	statusType := statusField.Type
	for statusType.Kind() == reflect.Pointer {
		statusType = statusType.Elem()
	}
	if _, ok := statusType.FieldByName("ObservedGeneration"); !ok {
		return fmt.Errorf("status type %s missing ObservedGeneration", statusType.Name())
	}
	condsField, ok := statusType.FieldByName("Conditions")
	if !ok {
		return fmt.Errorf("status type %s missing Conditions", statusType.Name())
	}
	condsType := condsField.Type
	if condsType.Kind() != reflect.Slice {
		return fmt.Errorf("status.Conditions must be a slice, got %s", condsType.Kind())
	}
	elem := condsType.Elem()
	for elem.Kind() == reflect.Pointer {
		elem = elem.Elem()
	}
	if elem != reflect.TypeOf(apicond.Condition{}) {
		return fmt.Errorf("status.Conditions element type=%s want apicond.Condition", elem.String())
	}

	for _, c := range []apicond.Condition{
		{Type: "Valid", Status: apicond.ConditionTrue, Reason: "ValidationSucceeded"},
		{Type: "TopologyComplete", Status: apicond.ConditionTrue, Reason: "PathComplete"},
	} {
		if !c.Valid() {
			return fmt.Errorf("condition grammar rejected well-formed %s/%s", c.Type, c.Reason)
		}
	}
	return nil
}

func assertFeature0014SchemaStatusConditions(schema []byte) error {
	var root map[string]any
	if err := json.Unmarshal(schema, &root); err != nil {
		return fmt.Errorf("schema JSON: %w", err)
	}
	props, _ := root["properties"].(map[string]any)
	status, _ := props["status"].(map[string]any)
	if status == nil {
		return fmt.Errorf("schema missing status")
	}
	statusProps, _ := status["properties"].(map[string]any)
	if _, ok := statusProps["observedGeneration"]; !ok {
		return fmt.Errorf("schema status missing observedGeneration")
	}
	conds, _ := statusProps["conditions"].(map[string]any)
	if conds == nil {
		return fmt.Errorf("schema status missing conditions")
	}
	items, _ := conds["items"].(map[string]any)
	if items == nil {
		return fmt.Errorf("schema status.conditions missing items")
	}
	ref, _ := items["$ref"].(string)
	if ref != "_common/condition.json" {
		return fmt.Errorf("schema status.conditions items $ref=%q want _common/condition.json", ref)
	}
	return nil
}

func scopeKindsEqual(got, want []apimeta.ScopeKind) bool {
	if len(got) != len(want) {
		return false
	}
	g := make([]string, len(got))
	w := make([]string, len(want))
	for i := range got {
		g[i] = string(got[i])
		w[i] = string(want[i])
	}
	sort.Strings(g)
	sort.Strings(w)
	for i := range g {
		if g[i] != w[i] {
			return false
		}
	}
	return true
}
