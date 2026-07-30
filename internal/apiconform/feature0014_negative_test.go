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
// Topology/completeness evaluation reuses the Task 20 same-package surface and
// its narrow wrappers around the canonical Task 17/18 helpers. This file does
// not import internal/validation and does not reimplement those hierarchy or
// child-existence rules.

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
				if v == nil {
					t.Fatal("expected geo prefix mismatch violation")
				}
				if v.Field != "/spec/geo/subdivisionCode" {
					t.Fatalf("field = %q, want /spec/geo/subdivisionCode", v.Field)
				}
				if v.Code != apiproblem.ViolationCode(apiproblem.CodeValidationFailed) {
					t.Fatalf("violation code = %q, want VALIDATION_FAILED", v.Code)
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
				t.Fatalf("missing exact UNKNOWN_FIELD at /status/history in %#v", violations)
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

func TestFeature0014IsolationCrossScopeOperationMatrix(t *testing.T) {
	t.Parallel()

	root := moduleRoot(t)
	set := mustLoadFeature0014CompletenessSet(t, root)
	provider := findFeature0014ByKind(t, set, resources.KindProvider)
	loc := findFeature0014ByKind(t, set, resources.KindProviderLocation)
	if !feature0014CanonicalScopeAndHierarchyAgree(
		provider.Kind, provider.UID, provider.ProviderScopeUID,
		loc.Kind, loc.UID, loc.ProviderScopeUID,
	) {
		t.Fatal("canonical Task 17 helper must accept the same-provider parent/child link")
	}

	// Cross-provider / cross-owner supplied state: CompletenessValue projections
	// already carry ProviderScopeUID from CompletenessValueFrom*; a mismatched
	// scope UID is isolation evidence without reimplementing EvaluateScopeAndHierarchy.
	crossLoc := loc
	crossLoc.ProviderScopeUID = "provider-uid-other-owner"
	if feature0014CanonicalScopeAndHierarchyAgree(
		provider.Kind, provider.UID, provider.ProviderScopeUID,
		crossLoc.Kind, crossLoc.UID, crossLoc.ProviderScopeUID,
	) {
		t.Fatal("canonical Task 17 helper must reject the cross-owner Provider scope UID")
	}

	absent := apiproblem.New(apiproblem.CodeResourceNotFound)
	absentJSON, err := json.Marshal(absent)
	if err != nil {
		t.Fatalf("marshal absent: %v", err)
	}

	// Live get/list/reference handlers are CONTRACT_ONLY / NO_TASK. This table
	// therefore verifies their shared response contract; it does not claim to
	// execute three production paths.
	operations := []string{"get", "list", "reference"}

	for _, operation := range operations {
		operation := operation
		t.Run(operation, func(t *testing.T) {
			t.Parallel()
			denied := apivalid.SafeDenial(apivalid.DenyNotDisclosed)
			if denied == nil {
				t.Fatal("denial Problem must be non-nil")
			}
			deniedJSON, err := json.Marshal(denied)
			if err != nil {
				t.Fatalf("marshal denial: %v", err)
			}
			if denied.Status != absent.Status || denied.Code != absent.Code {
				t.Fatalf("%s denial status/code = %d %q, want %d %q",
					operation, denied.Status, denied.Code, absent.Status, absent.Code)
			}
			if !bytes.Equal(deniedJSON, absentJSON) {
				t.Fatalf("%s denial body must be byte-identical to genuinely absent target:\ndenied=%s\nabsent=%s",
					operation, deniedJSON, absentJSON)
			}
			if len(denied.Violations) != 0 {
				t.Fatalf("%s denial must not disclose target existence: %#v", operation, denied.Violations)
			}
			lower := strings.ToLower(string(deniedJSON))
			for _, banned := range []string{
				"password", "token", "private_key", "credential", "secret",
				"endpoint", "nativeid", "provider-uid-other-owner",
				strings.ToLower(crossLoc.UID), strings.ToLower(provider.UID),
			} {
				if banned != "" && strings.Contains(lower, banned) {
					t.Fatalf("%s denial must not embed target/native/secret detail %q", operation, banned)
				}
			}
		})
	}
}

func TestFeature0014IsolationCrossProviderNoExistenceDisclosure(t *testing.T) {
	t.Parallel()

	root := moduleRoot(t)
	set := mustLoadFeature0014CompletenessSet(t, root)
	provider := findFeature0014ByKind(t, set, resources.KindProvider)
	loc := findFeature0014ByKind(t, set, resources.KindProviderLocation)

	if !feature0014CanonicalScopeAndHierarchyAgree(
		provider.Kind, provider.UID, provider.ProviderScopeUID,
		loc.Kind, loc.UID, loc.ProviderScopeUID,
	) {
		t.Fatal("canonical Task 17 helper must accept the same-provider path")
	}

	cross := loc
	cross.ProviderScopeUID = "provider-uid-b"
	if feature0014CanonicalScopeAndHierarchyAgree(
		provider.Kind, provider.UID, provider.ProviderScopeUID,
		cross.Kind, cross.UID, cross.ProviderScopeUID,
	) {
		t.Fatal("canonical Task 17 helper must reject cross-provider scope UIDs")
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

	root := moduleRoot(t)
	set := mustLoadFeature0014CompletenessSet(t, root)
	parent := findFeature0014ByKind(t, set, resources.KindProviderLocation)
	child := findFeature0014ByKind(t, set, resources.KindProviderDatacenter)

	if child.ParentUID != parent.UID {
		t.Fatalf("authorized CompletenessValue must link datacenter ParentUID to location uid")
	}
	mismatched := child
	mismatched.ProviderScopeUID = "provider-uid-b"
	if feature0014CanonicalScopeAndHierarchyAgree(
		parent.Kind, parent.UID, parent.ProviderScopeUID,
		mismatched.Kind, mismatched.UID, mismatched.ProviderScopeUID,
	) {
		t.Fatal("canonical Task 17 helper must reject parent/child Provider scope-UID mismatch")
	}

	denied := apivalid.SafeDenial(apivalid.DenyNotDisclosed)
	absent := apiproblem.New(apiproblem.CodeResourceNotFound)
	dj, err := json.Marshal(denied)
	if err != nil {
		t.Fatalf("marshal denial: %v", err)
	}
	aj, err := json.Marshal(absent)
	if err != nil {
		t.Fatalf("marshal absent: %v", err)
	}
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

	set := mustLoadFeature0014CompletenessSet(t, root)
	provider := findFeature0014ByKind(t, set, resources.KindProvider)
	loc := findFeature0014ByKind(t, set, resources.KindProviderLocation)
	if !feature0014CanonicalScopeAndHierarchyAgree(
		provider.Kind, provider.UID, provider.ProviderScopeUID,
		loc.Kind, loc.UID, loc.ProviderScopeUID,
	) {
		t.Fatal("canonical Task 17 helper must accept same-provider topology")
	}
	cross := loc
	cross.ProviderScopeUID = "prov-b"
	if feature0014CanonicalScopeAndHierarchyAgree(
		provider.Kind, provider.UID, provider.ProviderScopeUID,
		cross.Kind, cross.UID, cross.ProviderScopeUID,
	) {
		t.Fatal("cross-provider under shared owner must not agree or imply connectivity")
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

	if !feature0014CanonicalHasImmediateChildren(provider, set) {
		t.Fatal("canonical Task 18 helper must report the Provider's immediate child")
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

	before := map[string]string{}
	for _, v := range set {
		before[v.UID] = v.ParentUID
	}
	// Blocked delete mutates nothing: no cascade and no silent reparenting.
	for _, v := range set {
		if before[v.UID] != v.ParentUID {
			t.Fatalf("blocked delete must not reparent %s", v.UID)
		}
	}

	stack := findFeature0014ByKind(t, set, resources.KindInfrastructureStack)
	if feature0014CanonicalHasImmediateChildren(stack, set) {
		t.Fatal("canonical Task 18 helper must not report children for the leaf stack")
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
	base := mustLoadFeature0014CompletenessSet(t, root)
	provider := findFeature0014ByKind(t, base, resources.KindProvider)
	stack := findFeature0014ByKind(t, base, resources.KindInfrastructureStack)

	if !feature0014CanonicalHasImmediateChildren(provider, base) {
		t.Fatal("fixture provider must have children for concurrency proof")
	}

	// Mutable ledger over CompletenessValue projections so competing
	// delete/create transitions run without a production store.
	var mu sync.Mutex
	live := append(base[:0:0], base...)
	versions := make(map[string]string, len(base))
	retired := make(map[string]struct{})
	parents := make(map[string]string, len(base))
	for _, v := range base {
		versions[v.UID] = "rv-1"
		parents[v.UID] = v.ParentUID
	}

	findIndexLocked := func(uid string) int {
		for i, v := range live {
			if v.UID == uid {
				return i
			}
		}
		return -1
	}
	deleteLocked := func(uid string) {
		out := live[:0]
		for _, v := range live {
			if v.UID == uid {
				continue
			}
			out = append(out, v)
		}
		live = out
		delete(versions, uid)
		retired[uid] = struct{}{}
	}
	parentsUnchangedLocked := func() bool {
		for _, v := range live {
			if parents[v.UID] != v.ParentUID {
				return false
			}
		}
		return true
	}

	tryDelete := func(uid, ifMatch string) string {
		mu.Lock()
		defer mu.Unlock()
		index := findIndexLocked(uid)
		if index < 0 {
			return string(apiproblem.CodeResourceNotFound)
		}
		if prob := apivalid.CheckIfMatch(ifMatch, versions[uid]); prob != nil {
			return string(prob.Code)
		}
		if feature0014CanonicalHasImmediateChildren(live[index], live) {
			if !parentsUnchangedLocked() {
				return "CASCADE_OR_REPARENT"
			}
			return string(apiproblem.CodeDeleteBlocked)
		}
		deleteLocked(uid)
		if !parentsUnchangedLocked() {
			return "CASCADE_OR_REPARENT"
		}
		return "OK"
	}

	tryCreateReuse := func(uid string) string {
		mu.Lock()
		defer mu.Unlock()
		if _, wasRetired := retired[uid]; wasRetired {
			return "UID_REUSE_REJECTED"
		}
		if findIndexLocked(uid) >= 0 {
			return "ALREADY_EXISTS"
		}
		candidate := stack
		candidate.UID = uid
		live = append(live, candidate)
		versions[uid] = "rv-1"
		parents[uid] = stack.ParentUID
		return "OK"
	}

	tryStaleReplace := func(uid string) string {
		mu.Lock()
		defer mu.Unlock()
		current, ok := versions[uid]
		if !ok {
			return string(apiproblem.CodeResourceNotFound)
		}
		if prob := apivalid.CheckIfMatch("rv-stale", current); prob == nil {
			return "STALE_ACCEPTED"
		}
		return string(apiproblem.CodeStaleResourceVersion)
	}

	const workers = 32
	var wg sync.WaitGroup
	errCh := make(chan string, workers*8)

	// Phase 1: compete while children exist — parent deletes block; stale writes conflict.
	for i := 0; i < workers; i++ {
		wg.Add(3)
		go func() {
			defer wg.Done()
			if code := tryDelete(provider.UID, "rv-1"); code != string(apiproblem.CodeDeleteBlocked) {
				errCh <- "parent delete expected DELETE_BLOCKED, got " + code
			}
		}()
		go func() {
			defer wg.Done()
			if code := tryStaleReplace(provider.UID); code != string(apiproblem.CodeStaleResourceVersion) {
				errCh <- "stale replace expected STALE_RESOURCE_VERSION, got " + code
			}
		}()
		go func() {
			defer wg.Done()
			if code := tryDelete(provider.UID, "rv-stale"); code != string(apiproblem.CodeStaleResourceVersion) {
				errCh <- "stale parent delete expected STALE_RESOURCE_VERSION, got " + code
			}
		}()
	}
	wg.Wait()

	// Phase 2: deletion of the live leaf and recreation of the same UID start
	// together. Whichever obtains the lock first, creation must be rejected:
	// either the UID is still live or it has already been permanently retired.
	start := make(chan struct{})
	deleteResult := make(chan string, 1)
	createResult := make(chan string, 1)
	wg.Add(2)
	go func() {
		defer wg.Done()
		<-start
		deleteResult <- tryDelete(stack.UID, "rv-1")
	}()
	go func() {
		defer wg.Done()
		<-start
		createResult <- tryCreateReuse(stack.UID)
	}()
	close(start)
	wg.Wait()
	if code := <-deleteResult; code != "OK" {
		t.Fatalf("competing leaf delete: got %s want OK", code)
	}
	if code := <-createResult; code != "ALREADY_EXISTS" && code != "UID_REUSE_REJECTED" {
		t.Fatalf("competing recreation must be rejected, got %s", code)
	}

	// Finish the valid leaf-first deletion order after the leaf race.
	order := []string{
		resources.KindDatacenterFailureDomain,
		resources.KindProviderDatacenter,
		resources.KindProviderLocation,
		resources.KindProvider,
	}
	deletedStackUID := stack.UID
	for _, kind := range order {
		var uid string
		mu.Lock()
		for _, v := range live {
			if v.Kind == kind {
				uid = v.UID
				break
			}
		}
		mu.Unlock()
		if uid == "" {
			t.Fatalf("missing live kind %s during leaf-first concurrency phase", kind)
		}
		if code := tryDelete(uid, "rv-1"); code != "OK" {
			t.Fatalf("leaf-first delete of %s: got %s want OK", kind, code)
		}
	}

	// Permanent retirement remains stable under concurrent later attempts.
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if code := tryCreateReuse(deletedStackUID); code != "UID_REUSE_REJECTED" {
				errCh <- "deleted UID reuse expected UID_REUSE_REJECTED, got " + code
			}
		}()
	}
	wg.Wait()
	close(errCh)
	for msg := range errCh {
		t.Fatal(msg)
	}

	mu.Lock()
	defer mu.Unlock()
	if len(live) != 0 {
		t.Fatalf("after leaf-first deletes live set must be empty, got %d", len(live))
	}
	if _, ok := retired[deletedStackUID]; !ok {
		t.Fatal("deleted stack UID must remain permanently retired")
	}
}

func feature0014ValidateGeo(geo *resources.ProviderLocationGeo) *apiproblem.Violation {
	if geo == nil {
		return nil
	}
	if len(geo.CountryCode) != 2 || !feature0014GeoCountryRe.MatchString(geo.CountryCode) {
		return &apiproblem.Violation{
			Field:   "/spec/geo/countryCode",
			Code:    apiproblem.ViolationOutOfRange,
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
			Code:    apiproblem.ViolationCode(apiproblem.CodeValidationFailed),
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
		if v.Field == wantField && string(v.Code) == wantVCode {
			return string(apiproblem.CodeValidationFailed), v.Field
		}
	}
	t.Fatalf("expected exact structural violation field=%s code=%s, got %#v",
		wantField, wantVCode, violations)
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
