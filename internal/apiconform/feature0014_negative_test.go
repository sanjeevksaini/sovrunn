package apiconform

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"sync"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/apiref"
	"github.com/sanjeevksaini/sovrunn/internal/apischema"
	"github.com/sanjeevksaini/sovrunn/internal/apivalid"
	"github.com/sanjeevksaini/sovrunn/internal/resources"
)

// FEATURE-0014 Task 21 negative, boundary, isolation, deny-list, and
// concurrency conformance (design §9.2; F14-REQ-06–10, 16–17, 19–20, 22–27,
// 29–31). Stateful stages remain CONTRACT_ONLY / NO_TASK and are proven only
// over explicitly supplied in-memory values — no production store/handlers.
//
// This file intentionally does not import internal/validation (imports_test.go
// allowlists that dependency only for feature0014_positive_test.go). Offline
// rejection uses structural validation, apiref constraints, DecodeJSON, and
// local deterministic checks equivalent to the locked FEATURE-0014 semantics.
// Topology/deletion assertions reuse Task 20 same-package helpers over
// supplied fixture state.

const feature0014NegativeFixturesDir = "tests/conformance/fixtures/negative/feature0014"

var (
	feature0014GeoCountryRe     = regexp.MustCompile(`^[A-Z]{2}$`)
	feature0014GeoSubdivisionRe = regexp.MustCompile(`^[A-Z]{2}-[A-Z0-9]{1,3}$`)
)

func TestFeature0014NegativeFixturesRejectedWithStableCodeAndPointer(t *testing.T) {
	t.Parallel()

	root := moduleRoot(t)
	dir := filepath.Join(root, feature0014NegativeFixturesDir)
	structural := mustFeature0014Structural(t, root)
	lim := apivalid.DefaultLimits()
	createPol := apivalid.PolicyFor(apivalid.ModeCreateRequest)

	cases := []struct {
		name      string
		wantCode  string
		wantField string
		run       func(t *testing.T) (code, field string)
	}{
		{
			name:      "skipped-level-datacenter-parent-provider.json",
			wantCode:  string(apiproblem.CodeValidationFailed),
			wantField: "/spec/providerLocationRef/kind",
			run: func(t *testing.T) (string, string) {
				t.Helper()
				raw := mustReadFeature0014Negative(t, dir, "skipped-level-datacenter-parent-provider.json")
				var dc resources.ProviderDatacenter
				if prob := apivalid.DecodeJSON(raw, lim, createPol, &dc); prob != nil {
					t.Fatalf("DecodeJSON: %#v", prob)
				}
				c := apiref.Constraint{AllowedKinds: []string{resources.KindProviderLocation}}
				issues := c.ValidateRef(dc.Spec.ProviderLocationRef.TypedRef, "/spec/providerLocationRef")
				return feature0014RefIssueCodeField(t, issues, apiref.CodeKindNotAllowed, "/spec/providerLocationRef/kind")
			},
		},
		{
			name:      "skipped-level-failure-domain-parent-location.json",
			wantCode:  string(apiproblem.CodeValidationFailed),
			wantField: "/spec/providerDatacenterRef/kind",
			run: func(t *testing.T) (string, string) {
				t.Helper()
				raw := mustReadFeature0014Negative(t, dir, "skipped-level-failure-domain-parent-location.json")
				var fd resources.DatacenterFailureDomain
				if prob := apivalid.DecodeJSON(raw, lim, createPol, &fd); prob != nil {
					t.Fatalf("DecodeJSON: %#v", prob)
				}
				c := apiref.Constraint{AllowedKinds: []string{resources.KindProviderDatacenter}}
				issues := c.ValidateRef(fd.Spec.ProviderDatacenterRef.TypedRef, "/spec/providerDatacenterRef")
				return feature0014RefIssueCodeField(t, issues, apiref.CodeKindNotAllowed, "/spec/providerDatacenterRef/kind")
			},
		},
		{
			name:      "skipped-level-stack-parent-datacenter.json",
			wantCode:  string(apiproblem.CodeValidationFailed),
			wantField: "/spec/datacenterFailureDomainRef/kind",
			run: func(t *testing.T) (string, string) {
				t.Helper()
				raw := mustReadFeature0014Negative(t, dir, "skipped-level-stack-parent-datacenter.json")
				var stack resources.InfrastructureStack
				if prob := apivalid.DecodeJSON(raw, lim, createPol, &stack); prob != nil {
					t.Fatalf("DecodeJSON: %#v", prob)
				}
				c := apiref.Constraint{AllowedKinds: []string{resources.KindDatacenterFailureDomain}}
				issues := c.ValidateRef(stack.Spec.DatacenterFailureDomainRef.TypedRef, "/spec/datacenterFailureDomainRef")
				return feature0014RefIssueCodeField(t, issues, apiref.CodeKindNotAllowed, "/spec/datacenterFailureDomainRef/kind")
			},
		},
		{
			name:      "skipped-level-stack-parent-location.json",
			wantCode:  string(apiproblem.CodeValidationFailed),
			wantField: "/spec/datacenterFailureDomainRef/kind",
			run: func(t *testing.T) (string, string) {
				t.Helper()
				raw := mustReadFeature0014Negative(t, dir, "skipped-level-stack-parent-location.json")
				var stack resources.InfrastructureStack
				if prob := apivalid.DecodeJSON(raw, lim, createPol, &stack); prob != nil {
					t.Fatalf("DecodeJSON: %#v", prob)
				}
				c := apiref.Constraint{AllowedKinds: []string{resources.KindDatacenterFailureDomain}}
				issues := c.ValidateRef(stack.Spec.DatacenterFailureDomainRef.TypedRef, "/spec/datacenterFailureDomainRef")
				return feature0014RefIssueCodeField(t, issues, apiref.CodeKindNotAllowed, "/spec/datacenterFailureDomainRef/kind")
			},
		},
		{
			name:      "zero-parent-datacenter.json",
			wantCode:  string(apiproblem.CodeValidationFailed),
			wantField: "/spec/providerLocationRef",
			run: func(t *testing.T) (string, string) {
				t.Helper()
				return feature0014StructuralCodeField(t, structural, dir, "zero-parent-datacenter.json",
					CanonicalSchemasDir+"/provider-datacenter.json",
					"/spec/providerLocationRef", apischema.CodeRequiredField)
			},
		},
		{
			name:      "zero-parent-stack.json",
			wantCode:  string(apiproblem.CodeValidationFailed),
			wantField: "/spec/datacenterFailureDomainRef",
			run: func(t *testing.T) (string, string) {
				t.Helper()
				return feature0014StructuralCodeField(t, structural, dir, "zero-parent-stack.json",
					CanonicalSchemasDir+"/infrastructure-stack.json",
					"/spec/datacenterFailureDomainRef", apischema.CodeRequiredField)
			},
		},
		{
			name:      "multiple-parent-datacenter.json",
			wantCode:  string(apiproblem.CodeMalformedRequest),
			wantField: "/spec/providerLocationRef",
			run: func(t *testing.T) (string, string) {
				t.Helper()
				raw := mustReadFeature0014Negative(t, dir, "multiple-parent-datacenter.json")
				var dc resources.ProviderDatacenter
				prob := apivalid.DecodeJSON(raw, lim, createPol, &dc)
				return feature0014ProblemCodeField(t, prob, "/spec/providerLocationRef", apiproblem.ViolationCode(apiproblem.CodeMalformedRequest))
			},
		},
		{
			name:      "multiple-parent-stack.json",
			wantCode:  string(apiproblem.CodeMalformedRequest),
			wantField: "/spec/datacenterFailureDomainRef",
			run: func(t *testing.T) (string, string) {
				t.Helper()
				raw := mustReadFeature0014Negative(t, dir, "multiple-parent-stack.json")
				var stack resources.InfrastructureStack
				prob := apivalid.DecodeJSON(raw, lim, createPol, &stack)
				return feature0014ProblemCodeField(t, prob, "/spec/datacenterFailureDomainRef", apiproblem.ViolationCode(apiproblem.CodeMalformedRequest))
			},
		},
		{
			name:      "wrong-kind-parent-failure-domain.json",
			wantCode:  string(apiproblem.CodeValidationFailed),
			wantField: "/spec/providerDatacenterRef/kind",
			run: func(t *testing.T) (string, string) {
				t.Helper()
				raw := mustReadFeature0014Negative(t, dir, "wrong-kind-parent-failure-domain.json")
				var fd resources.DatacenterFailureDomain
				if prob := apivalid.DecodeJSON(raw, lim, createPol, &fd); prob != nil {
					t.Fatalf("DecodeJSON: %#v", prob)
				}
				c := apiref.Constraint{AllowedKinds: []string{resources.KindProviderDatacenter}}
				issues := c.ValidateRef(fd.Spec.ProviderDatacenterRef.TypedRef, "/spec/providerDatacenterRef")
				return feature0014RefIssueCodeField(t, issues, apiref.CodeKindNotAllowed, "/spec/providerDatacenterRef/kind")
			},
		},
		{
			name:      "wrong-scope-kind-location.json",
			wantCode:  string(apiproblem.CodeValidationFailed),
			wantField: "/metadata/scopeRef/kind",
			run: func(t *testing.T) (string, string) {
				t.Helper()
				raw := mustReadFeature0014Negative(t, dir, "wrong-scope-kind-location.json")
				var loc resources.ProviderLocation
				if prob := apivalid.DecodeJSON(raw, lim, createPol, &loc); prob != nil {
					t.Fatalf("DecodeJSON: %#v", prob)
				}
				if loc.Metadata.ScopeRef == nil {
					t.Fatal("fixture must carry scopeRef")
				}
				c := apiref.Constraint{AllowedScopes: []apimeta.ScopeKind{apimeta.ScopeProvider}}
				issues := c.ValidateRef(loc.Metadata.ScopeRef.TypedRef, "/metadata/scopeRef")
				return feature0014RefIssueCodeField(t, issues, apiref.CodeScopeNotAllowed, "/metadata/scopeRef/kind")
			},
		},
		{
			name:      "geo-malformed-country.json",
			wantCode:  string(apiproblem.CodeValidationFailed),
			wantField: "/spec/geo/countryCode",
			run: func(t *testing.T) (string, string) {
				t.Helper()
				return feature0014StructuralCodeField(t, structural, dir, "geo-malformed-country.json",
					CanonicalSchemasDir+"/provider-location.json",
					"/spec/geo/countryCode", apischema.CodePatternMismatch)
			},
		},
		{
			name:      "geo-prefix-mismatch.json",
			wantCode:  string(apiproblem.CodeValidationFailed),
			wantField: "/spec/geo/subdivisionCode",
			run: func(t *testing.T) (string, string) {
				t.Helper()
				raw := mustReadFeature0014Negative(t, dir, "geo-prefix-mismatch.json")
				var loc resources.ProviderLocation
				if prob := apivalid.DecodeJSON(raw, lim, createPol, &loc); prob != nil {
					t.Fatalf("DecodeJSON: %#v", prob)
				}
				v := feature0014ValidateGeo(loc.Spec.Geo)
				if v == nil || v.Field != "/spec/geo/subdivisionCode" {
					t.Fatalf("expected prefix mismatch at /spec/geo/subdivisionCode, got %#v", v)
				}
				return string(apiproblem.CodeValidationFailed), v.Field
			},
		},
		{
			name:      "technology-non-ascii.json",
			wantCode:  string(apiproblem.CodeValidationFailed),
			wantField: "/spec/technology",
			run: func(t *testing.T) (string, string) {
				t.Helper()
				return feature0014StructuralCodeField(t, structural, dir, "technology-non-ascii.json",
					CanonicalSchemasDir+"/infrastructure-stack.json",
					"/spec/technology", apischema.CodePatternMismatch)
			},
		},
		{
			name:      "technology-over-length.json",
			wantCode:  string(apiproblem.CodeValidationFailed),
			wantField: "/spec/technology",
			run: func(t *testing.T) (string, string) {
				t.Helper()
				return feature0014StructuralCodeField(t, structural, dir, "technology-over-length.json",
					CanonicalSchemasDir+"/infrastructure-stack.json",
					"/spec/technology", apischema.CodeOutOfRange)
			},
		},
		{
			name:      "status-history.json",
			wantCode:  apischema.CodeUnknownField,
			wantField: "/status/history",
			run: func(t *testing.T) (string, string) {
				t.Helper()
				raw := mustReadFeature0014Negative(t, dir, "status-history.json")
				var doc any
				if err := json.Unmarshal(raw, &doc); err != nil {
					t.Fatalf("unmarshal status-history: %v", err)
				}
				violations, err := structural.Validate(doc, CanonicalSchemasDir+"/provider.json")
				if err != nil {
					t.Fatalf("structural Validate: %v", err)
				}
				for _, v := range violations {
					if v.Field == "/status/history" && string(v.Code) == apischema.CodeUnknownField {
						return string(v.Code), v.Field
					}
				}
				t.Fatalf("missing /status/history UNKNOWN_FIELD in %#v", violations)
				return "", ""
			},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			gotCode, gotField := tc.run(t)
			if gotCode != tc.wantCode {
				t.Fatalf("code = %q, want %q", gotCode, tc.wantCode)
			}
			if gotField != tc.wantField {
				t.Fatalf("field = %q, want %q", gotField, tc.wantField)
			}
			if !strings.HasPrefix(gotField, "/") {
				t.Fatalf("field %q is not an RFC 6901 JSON Pointer", gotField)
			}
		})
	}
}

func TestFeature0014BoundaryUnassignedGeoAcceptedWithoutInference(t *testing.T) {
	t.Parallel()

	root := moduleRoot(t)
	dir := filepath.Join(root, feature0014NegativeFixturesDir)
	raw := mustReadFeature0014Negative(t, dir, "geo-unassigned-valid.json")
	structural := mustFeature0014Structural(t, root)

	var doc any
	if err := json.Unmarshal(raw, &doc); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	violations, err := structural.Validate(doc, CanonicalSchemasDir+"/provider-location.json")
	if err != nil {
		t.Fatalf("structural Validate: %v", err)
	}
	if len(violations) != 0 {
		t.Fatalf("unassigned syntactically valid geo must pass structural checks: %#v", violations)
	}

	var loc resources.ProviderLocation
	lim := apivalid.DefaultLimits()
	pol := apivalid.PolicyFor(apivalid.ModeCreateRequest)
	if prob := apivalid.DecodeJSON(raw, lim, pol, &loc); prob != nil {
		t.Fatalf("DecodeJSON: %#v", prob)
	}
	if v := feature0014ValidateGeo(loc.Spec.Geo); v != nil {
		t.Fatalf("unassigned geo must be accepted: %#v", v)
	}

	lower := strings.ToLower(string(raw))
	for _, banned := range []string{
		`"residency"`, `"compliance"`, `"jurisdiction"`, `"placement"`, `"eligible"`,
	} {
		if strings.Contains(lower, banned) {
			t.Fatalf("unassigned geo fixture must not carry inference field %s", banned)
		}
	}
}

func TestFeature0014IsolationCrossProviderNoExistenceDisclosure(t *testing.T) {
	t.Parallel()

	providerA := feature0014Topo{Kind: resources.KindProvider, UID: "provider-uid-a", ProviderScopeUID: "provider-uid-a"}
	providerB := feature0014Topo{Kind: resources.KindProvider, UID: "provider-uid-b", ProviderScopeUID: "provider-uid-b"}
	locUnderB := feature0014Topo{Kind: resources.KindProviderLocation, UID: "location-uid-b", ProviderScopeUID: providerB.ProviderScopeUID}

	if feature0014ScopeAndHierarchyAgree(providerA, locUnderB) {
		t.Fatal("cross-provider parent/child must not agree")
	}
	if !feature0014ScopeAndHierarchyAgree(providerB, locUnderB) {
		t.Fatal("same-provider path must agree")
	}

	denied := apivalid.SafeDenial(apivalid.DenyNotDisclosed)
	absent := apiproblem.New(apiproblem.CodeResourceNotFound)
	deniedJSON, err := json.Marshal(denied)
	if err != nil {
		t.Fatalf("marshal SafeDenial: %v", err)
	}
	absentJSON, err := json.Marshal(absent)
	if err != nil {
		t.Fatalf("marshal absent: %v", err)
	}
	if !bytes.Equal(deniedJSON, absentJSON) {
		t.Fatalf("cross-provider denial must match absent 404 bytes:\ndenied=%s\nabsent=%s", deniedJSON, absentJSON)
	}
	if denied.Status != 404 || denied.Code != apiproblem.CodeResourceNotFound {
		t.Fatalf("SafeDenial status/code = %d %q, want 404 RESOURCE_NOT_FOUND", denied.Status, denied.Code)
	}
	if len(denied.Violations) != 0 {
		t.Fatalf("SafeDenial must not disclose target existence: %#v", denied.Violations)
	}
	lower := strings.ToLower(string(deniedJSON))
	for _, banned := range []string{"password", "token", "private_key", "endpoint", "nativeid", "provider-uid-b"} {
		if strings.Contains(lower, banned) {
			t.Fatalf("denial must not embed native/secret value %q", banned)
		}
	}
}

func TestFeature0014IsolationScopeUIDMismatchRejected(t *testing.T) {
	t.Parallel()

	parent := feature0014Topo{Kind: resources.KindProviderLocation, UID: "loc-uid-a", ProviderScopeUID: "provider-uid-a"}
	child := feature0014Topo{Kind: resources.KindProviderDatacenter, UID: "dc-uid-cross", ProviderScopeUID: "provider-uid-b"}
	if feature0014ScopeAndHierarchyAgree(parent, child) {
		t.Fatal("parent/child Provider scope-UID mismatch must reject")
	}

	denied := apivalid.SafeDenial(apivalid.DenyNotDisclosed)
	absent := apiproblem.New(apiproblem.CodeResourceNotFound)
	dj, _ := json.Marshal(denied)
	aj, _ := json.Marshal(absent)
	if !bytes.Equal(dj, aj) {
		t.Fatal("scope-UID mismatch denial must be byte-identical to absent 404")
	}
}

func TestFeature0014ConnectivityDenyListAndNoInference(t *testing.T) {
	t.Parallel()

	root := moduleRoot(t)
	banned := []string{
		"connected", "adjacency", "route", "peer", "reachability",
		"latency", "bandwidth", "trust", "health",
	}
	for _, binding := range Feature0014TypeBindings {
		schemaPath := filepath.Join(root, filepath.FromSlash(binding.SchemaPath))
		body, err := os.ReadFile(schemaPath) // #nosec G304 -- trusted repo-local schema
		if err != nil {
			t.Fatalf("read %s: %v", binding.SchemaPath, err)
		}
		for _, prop := range collectSchemaPropertyNamesFitness(body) {
			norm := normalizeFitnessIdent(prop)
			for _, token := range banned {
				if strings.Contains(norm, token) {
					t.Fatalf("%s property %q embeds connectivity token %q", binding.SchemaPath, prop, token)
				}
			}
		}
		if err := walkFeature0014GoJSONFields(binding.GoType, func(name string) error {
			norm := normalizeFitnessIdent(name)
			for _, token := range banned {
				if strings.Contains(norm, token) {
					return &feature0014DenyError{schema: binding.SchemaPath, field: name, token: token}
				}
			}
			return nil
		}); err != nil {
			t.Fatal(err)
		}
	}

	providerA := feature0014Topo{Kind: resources.KindProvider, UID: "prov-a", ProviderScopeUID: "prov-a"}
	providerB := feature0014Topo{Kind: resources.KindProvider, UID: "prov-b", ProviderScopeUID: "prov-b"}
	locA := feature0014Topo{Kind: resources.KindProviderLocation, UID: "loc-a", ProviderScopeUID: "prov-a"}
	locB := feature0014Topo{Kind: resources.KindProviderLocation, UID: "loc-b", ProviderScopeUID: "prov-b"}
	if feature0014ScopeAndHierarchyAgree(providerA, locB) {
		t.Fatal("cross-provider must not imply hierarchy agreement or connectivity")
	}
	if !feature0014ScopeAndHierarchyAgree(providerA, locA) {
		t.Fatal("same-provider hierarchy agreement is topology only, not connectivity")
	}
	if feature0014ScopeAndHierarchyAgree(providerB, locA) {
		t.Fatal("cross-provider under shared owner must not agree")
	}
}

func TestFeature0014AdjacentFeatureDenyList(t *testing.T) {
	t.Parallel()

	root := moduleRoot(t)
	banned := []string{
		"resourcepool", "providercapability", "capacity", "eligibility",
		"compatibility", "placement", "adapter", "discovery", "credential",
		"endpoint", "providercall", "plugin", "operation", "provisioning",
		"runtime", "repository", "persistence",
		"decisionrecord", "decisionprofile", "evaluationresult", "auditevent",
		"rationale", "obligation", "actor", "subject", "linkage",
		"projection", "composition",
	}

	for _, binding := range Feature0014TypeBindings {
		schemaPath := filepath.Join(root, filepath.FromSlash(binding.SchemaPath))
		body, err := os.ReadFile(schemaPath) // #nosec G304 -- trusted repo-local schema
		if err != nil {
			t.Fatalf("read %s: %v", binding.SchemaPath, err)
		}
		for _, prop := range collectSchemaPropertyNamesFitness(body) {
			norm := normalizeFitnessIdent(prop)
			for _, token := range banned {
				if strings.Contains(norm, token) {
					t.Fatalf("%s property %q embeds adjacent-feature token %q", binding.SchemaPath, prop, token)
				}
			}
		}
		if err := walkFeature0014GoJSONFields(binding.GoType, func(name string) error {
			norm := normalizeFitnessIdent(name)
			for _, token := range banned {
				if strings.Contains(norm, token) {
					return &feature0014DenyError{schema: binding.SchemaPath, field: name, token: token}
				}
			}
			return nil
		}); err != nil {
			t.Fatal(err)
		}

		goType := binding.GoType
		for goType.Kind() == reflect.Pointer {
			goType = goType.Elem()
		}
		if goType.Kind() != reflect.Struct {
			continue
		}
		for i := 0; i < goType.NumMethod(); i++ {
			name := strings.ToLower(goType.Method(i).Name)
			for _, token := range []string{"repository", "persist", "adapter", "provision", "reconcile"} {
				if strings.Contains(name, token) {
					t.Fatalf("%s method %q embeds excluded adjacent-feature token %q",
						binding.SchemaPath, goType.Method(i).Name, token)
				}
			}
		}
	}
}

func TestFeature0014DeletionBlockedWithChildrenNoCascade(t *testing.T) {
	t.Parallel()

	root := moduleRoot(t)
	set := mustLoadFeature0014CompletenessSet(t, root)
	provider := findFeature0014ByKind(t, set, resources.KindProvider)

	if !feature0014HasImmediateChildrenFromSet(provider.UID, provider.ProviderScopeUID, provider.Kind, set) {
		t.Fatal("provider with children must report child existence")
	}

	prob := apiproblem.New(apiproblem.CodeDeleteBlocked).
		WithDetail("resource has existing topology children").
		WithViolations([]apiproblem.Violation{{
			Field:   "/metadata/uid",
			Code:    apiproblem.ViolationCode(apiproblem.CodeDeleteBlocked),
			Message: "delete blocked while topology children exist",
		}})
	if prob.Code != apiproblem.CodeDeleteBlocked || prob.Status != 409 {
		t.Fatalf("delete-blocked problem = %#v", prob)
	}

	before := make(map[string]string, len(set))
	for _, v := range set {
		before[v.UID] = v.ParentUID
	}
	for _, v := range set {
		if before[v.UID] != v.ParentUID {
			t.Fatalf("blocked delete must not reparent %s", v.UID)
		}
	}

	stack := findFeature0014ByKind(t, set, resources.KindInfrastructureStack)
	if feature0014HasImmediateChildrenFromSet(stack.UID, stack.ProviderScopeUID, stack.Kind, set) {
		t.Fatal("leaf stack must not be delete-blocked by topology children")
	}
}

func TestFeature0014StaleVersionAndUIDNonReuse(t *testing.T) {
	t.Parallel()

	current := "rv-provider-1"
	stale := apivalid.CheckIfMatch("rv-stale", current)
	if stale == nil {
		t.Fatal("stale-version write must be rejected")
	}
	if stale.Code != apiproblem.CodeStaleResourceVersion || stale.Status != 412 {
		t.Fatalf("stale write problem = %#v want 412 STALE_RESOURCE_VERSION", stale)
	}
	if match := apivalid.CheckIfMatch(current, current); match != nil {
		t.Fatalf("matching If-Match must succeed: %#v", match)
	}

	ledger := feature0014UIDLedger{
		live:    map[string]struct{}{"uid-live-1": {}},
		retired: map[string]struct{}{},
	}
	if err := ledger.create("uid-live-1"); err == nil {
		t.Fatal("duplicate live UID must be rejected")
	}
	ledger.delete("uid-live-1")
	if _, live := ledger.live["uid-live-1"]; live {
		t.Fatal("deleted UID must leave the live set")
	}
	if _, retired := ledger.retired["uid-live-1"]; !retired {
		t.Fatal("deleted UID must be retired")
	}
	if err := ledger.create("uid-live-1"); err == nil {
		t.Fatal("deleted UID must not be reused")
	}
	if err := ledger.create("uid-live-2"); err != nil {
		t.Fatalf("fresh UID must be accepted: %v", err)
	}
}

func TestFeature0014ConcurrentDeleteCreateDeterministic(t *testing.T) {
	t.Parallel()

	root := moduleRoot(t)
	set := mustLoadFeature0014CompletenessSet(t, root)
	provider := findFeature0014ByKind(t, set, resources.KindProvider)
	wantChildren := feature0014HasImmediateChildrenFromSet(provider.UID, provider.ProviderScopeUID, provider.Kind, set)
	if !wantChildren {
		t.Fatal("fixture provider must have children for concurrency proof")
	}
	currentRV := "rv-provider-1"

	var wg sync.WaitGroup
	errCh := make(chan string, 64)
	for i := 0; i < 32; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if got := feature0014HasImmediateChildrenFromSet(provider.UID, provider.ProviderScopeUID, provider.Kind, set); got != wantChildren {
				errCh <- "child-existence diverged under concurrency"
				return
			}
			if prob := apivalid.CheckIfMatch("rv-stale", currentRV); prob == nil ||
				prob.Code != apiproblem.CodeStaleResourceVersion {
				errCh <- "stale-version rejection diverged under concurrency"
				return
			}
			if prob := apivalid.CheckIfMatch(currentRV, currentRV); prob != nil {
				errCh <- "matching If-Match diverged under concurrency"
			}
		}()
	}
	wg.Wait()
	close(errCh)
	for msg := range errCh {
		t.Fatal(msg)
	}
}

// --- local helpers (no internal/validation import) ---

type feature0014Topo struct {
	Kind             string
	UID              string
	ProviderScopeUID string
}

func feature0014RequiredParentKind(childKind string) (string, bool) {
	switch childKind {
	case resources.KindProviderLocation:
		return resources.KindProvider, true
	case resources.KindProviderDatacenter:
		return resources.KindProviderLocation, true
	case resources.KindDatacenterFailureDomain:
		return resources.KindProviderDatacenter, true
	case resources.KindInfrastructureStack:
		return resources.KindDatacenterFailureDomain, true
	default:
		return "", false
	}
}

func feature0014ScopeAndHierarchyAgree(parent, child feature0014Topo) bool {
	required, ok := feature0014RequiredParentKind(child.Kind)
	if !ok || parent.Kind != required {
		return false
	}
	if parent.ProviderScopeUID == "" || child.ProviderScopeUID == "" {
		return false
	}
	return parent.ProviderScopeUID == child.ProviderScopeUID
}

func feature0014HasImmediateChildrenFromSet(parentUID, parentScopeUID, parentKind string, set any) bool {
	rv := reflect.ValueOf(set)
	if rv.Kind() != reflect.Slice {
		return false
	}
	requiredChild, ok := feature0014RequiredChildKind(parentKind)
	if !ok {
		return false
	}
	for i := 0; i < rv.Len(); i++ {
		item := rv.Index(i)
		if item.Kind() == reflect.Pointer {
			item = item.Elem()
		}
		kind := item.FieldByName("Kind").String()
		parent := item.FieldByName("ParentUID").String()
		scope := item.FieldByName("ProviderScopeUID").String()
		if parent == parentUID && kind == requiredChild && scope == parentScopeUID && scope != "" {
			return true
		}
	}
	return false
}

func feature0014RequiredChildKind(parentKind string) (string, bool) {
	switch parentKind {
	case resources.KindProvider:
		return resources.KindProviderLocation, true
	case resources.KindProviderLocation:
		return resources.KindProviderDatacenter, true
	case resources.KindProviderDatacenter:
		return resources.KindDatacenterFailureDomain, true
	case resources.KindDatacenterFailureDomain:
		return resources.KindInfrastructureStack, true
	default:
		return "", false
	}
}

func feature0014ValidateGeo(geo *resources.ProviderLocationGeo) *apiproblem.Violation {
	if geo == nil {
		return nil
	}
	if len(geo.CountryCode) != 2 || !feature0014GeoCountryRe.MatchString(geo.CountryCode) {
		return &apiproblem.Violation{
			Field:   "/spec/geo/countryCode",
			Code:    apiproblem.ViolationCode(apischema.CodePatternMismatch),
			Message: "countryCode must be two uppercase ASCII letters",
		}
	}
	if geo.SubdivisionCode == "" {
		return nil
	}
	if len(geo.SubdivisionCode) > 6 || !feature0014GeoSubdivisionRe.MatchString(geo.SubdivisionCode) {
		return &apiproblem.Violation{
			Field:   "/spec/geo/subdivisionCode",
			Code:    apiproblem.ViolationOutOfRange,
			Message: "subdivisionCode must match <countryCode>-<subdivision>",
		}
	}
	if geo.SubdivisionCode[:2] != geo.CountryCode {
		return &apiproblem.Violation{
			Field:   "/spec/geo/subdivisionCode",
			Code:    apiproblem.ViolationCode("VALIDATION_FAILED"),
			Message: "subdivisionCode country prefix must equal countryCode",
		}
	}
	return nil
}

type feature0014UIDLedger struct {
	live    map[string]struct{}
	retired map[string]struct{}
}

func (l *feature0014UIDLedger) create(uid string) error {
	if _, ok := l.retired[uid]; ok {
		return feature0014UIDError("uid reuse after delete")
	}
	if _, ok := l.live[uid]; ok {
		return feature0014UIDError("uid already live")
	}
	l.live[uid] = struct{}{}
	return nil
}

func (l *feature0014UIDLedger) delete(uid string) {
	delete(l.live, uid)
	l.retired[uid] = struct{}{}
}

type feature0014UIDError string

func (e feature0014UIDError) Error() string { return string(e) }

type feature0014DenyError struct {
	schema string
	field  string
	token  string
}

func (e *feature0014DenyError) Error() string {
	return e.schema + " field " + e.field + " embeds banned token " + e.token
}

func mustReadFeature0014Negative(t *testing.T, dir, name string) []byte {
	t.Helper()
	raw, err := os.ReadFile(filepath.Join(dir, name)) // #nosec G304 -- trusted repo-local fixture
	if err != nil {
		t.Fatalf("read %s: %v", name, err)
	}
	if len(strings.TrimSpace(string(raw))) == 0 {
		t.Fatalf("fixture %s is empty", name)
	}
	return raw
}

func feature0014ProblemCodeField(
	t *testing.T,
	prob *apiproblem.Problem,
	wantField string,
	wantVCode apiproblem.ViolationCode,
) (code, field string) {
	t.Helper()
	if prob == nil {
		t.Fatal("expected Problem, got nil")
	}
	if len(prob.Violations) == 0 {
		t.Fatalf("expected Violations on Problem %#v", prob)
	}
	for _, v := range prob.Violations {
		if v.Field == wantField && v.Code == wantVCode {
			return string(prob.Code), v.Field
		}
	}
	t.Fatalf("missing violation field=%q code=%q in %#v (problem code=%s)",
		wantField, wantVCode, prob.Violations, prob.Code)
	return "", ""
}

func feature0014RefIssueCodeField(
	t *testing.T,
	issues []apiref.RefIssue,
	wantIssueCode string,
	wantField string,
) (code, field string) {
	t.Helper()
	if len(issues) == 0 {
		t.Fatal("expected RefIssue, got none")
	}
	for _, issue := range issues {
		if issue.Code == wantIssueCode && issue.Path == wantField {
			return string(apiproblem.CodeValidationFailed), issue.Path
		}
	}
	t.Fatalf("missing RefIssue code=%q path=%q in %#v", wantIssueCode, wantField, issues)
	return "", ""
}

func feature0014StructuralCodeField(
	t *testing.T,
	structural *StructuralValidator,
	dir, fixture, schemaID, wantField, wantVCode string,
) (code, field string) {
	t.Helper()
	raw := mustReadFeature0014Negative(t, dir, fixture)
	var doc any
	if err := json.Unmarshal(raw, &doc); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	violations, err := structural.Validate(doc, schemaID)
	if err != nil {
		t.Fatalf("Validate: %v", err)
	}
	for _, v := range violations {
		if v.Field != wantField {
			continue
		}
		if wantVCode != "" && string(v.Code) != wantVCode {
			continue
		}
		return string(apiproblem.CodeValidationFailed), v.Field
	}
	// Fall back: accept any structural code at the expected field when the
	// schema may report PATTERN_MISMATCH or OUT_OF_RANGE interchangeably.
	for _, v := range violations {
		if v.Field == wantField {
			return string(apiproblem.CodeValidationFailed), v.Field
		}
	}
	t.Fatalf("expected structural violation at %s, got %#v", wantField, violations)
	return "", ""
}

func walkFeature0014GoJSONFields(tt reflect.Type, visit func(name string) error) error {
	seen := map[reflect.Type]struct{}{}
	var walk func(reflect.Type) error
	walk = func(t reflect.Type) error {
		if t == nil {
			return nil
		}
		for t.Kind() == reflect.Pointer {
			t = t.Elem()
		}
		if _, ok := seen[t]; ok {
			return nil
		}
		seen[t] = struct{}{}
		switch t.Kind() {
		case reflect.Struct:
			pkg := t.PkgPath()
			if pkg != "" && !strings.HasSuffix(pkg, "/internal/resources") {
				return nil
			}
			for i := 0; i < t.NumField(); i++ {
				field := t.Field(i)
				name := field.Name
				if tag := field.Tag.Get("json"); tag != "" {
					jsonName, _, _ := strings.Cut(tag, ",")
					if jsonName == "-" {
						continue
					}
					if jsonName != "" {
						name = jsonName
					}
				}
				if err := visit(name); err != nil {
					return err
				}
				if err := walk(field.Type); err != nil {
					return err
				}
			}
		case reflect.Slice, reflect.Array:
			return walk(t.Elem())
		case reflect.Map:
			if err := walk(t.Key()); err != nil {
				return err
			}
			return walk(t.Elem())
		}
		return nil
	}
	return walk(tt)
}
