package apiconform

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/apiref"
	"github.com/sanjeevksaini/sovrunn/internal/apischema"
	"github.com/sanjeevksaini/sovrunn/internal/apivalid"
)

// Fitness check IDs for F12-VERIFY-001 checks implemented in this file
// (task 16.2: checks 2, 5, 6, 7, 8). Aggregation across 1–15 is task 16.5.
const (
	FitnessCheckNoProviderSDKInCoreCustomer       = "2"
	FitnessCheckReferencesConstrainKindsAndScopes = "5"
	FitnessCheckCrossTenantNoExistenceDisclosure  = "6"
	FitnessCheckNoRawSecretLikeValues             = "7"
	FitnessCheckObservationProvenanceAndFreshness = "8"
)

// Stable fitness finding codes for reference/scope/boundary/security checks.
const (
	CodeFitnessProviderSDKImport        = "FITNESS_PROVIDER_SDK_IMPORT"
	CodeFitnessProviderNativeField      = "FITNESS_PROVIDER_NATIVE_FIELD"
	CodeFitnessProviderSDKGoType        = "FITNESS_PROVIDER_SDK_GO_TYPE"
	CodeFitnessScopeKindUnconstrained   = "FITNESS_SCOPE_KIND_UNCONSTRAINED"
	CodeFitnessTypedRefIncomplete       = "FITNESS_TYPED_REF_INCOMPLETE"
	CodeFitnessRefNotConstrained        = "FITNESS_REF_NOT_CONSTRAINED"
	CodeFitnessAllowedScopesMissing     = "FITNESS_ALLOWED_SCOPES_MISSING"
	CodeFitnessRefConstraintIneffective = "FITNESS_REF_CONSTRAINT_INEFFECTIVE"
	CodeFitnessSafeDenialMismatch       = "FITNESS_SAFE_DENIAL_MISMATCH"
	CodeFitnessSafeDenialDiscloses      = "FITNESS_SAFE_DENIAL_DISCLOSES"
	CodeFitnessSecretLikeValue          = "FITNESS_SECRET_LIKE_VALUE"
	CodeFitnessProvenanceMissing        = "FITNESS_PROVENANCE_MISSING"
	CodeFitnessFreshnessMissing         = "FITNESS_FRESHNESS_MISSING"
)

// Core grammar packages that must remain provider-neutral (F12-VERIFY-001(2),
// F12-SEC-006, Property 7). apischema/apiconform are scanned for SDK imports
// only; schema property bans apply to core/customer boundaries.
var fitnessCoreGrammarPackages = []string{
	"apimeta",
	"apiref",
	"apicond",
	"apiproblem",
	"apivalid",
	"apischema",
	"apiconform",
}

// providerSDKImportPrefixes are banned import path prefixes in grammar packages.
var providerSDKImportPrefixes = []string{
	"k8s.io/",
	"sigs.k8s.io/",
	"github.com/kubernetes/",
	"github.com/aws/",
	"github.com/Azure/",
	"github.com/microsoft/azure-",
	"cloud.google.com/",
	"google.golang.org/api/",
	"google.golang.org/genproto",
	"github.com/googleapis/",
	"github.com/oracle/oci-go-sdk",
	"github.com/digitalocean/",
	"github.com/hetznercloud/",
	"github.com/linode/",
	"github.com/vultr/",
	"github.com/IBM/",
	"github.com/hashicorp/terraform-provider",
	"github.com/pulumi/",
}

// bannedCoreNativeFieldTokens are matched against lowercased alphanumeric-only
// property/field identifiers in core/customer schemas and Go types. Descriptions
// are not scanned (avoid false positives on "provider-neutral").
var bannedCoreNativeFieldTokens = []string{
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

// secretLikeTokensFitness are scanned case-insensitively against metadata
// labels/annotations, status maps, error/problem fixtures, and property names
// under metadata/status/errors (F12-SEC-003). Composite phrases only —
// plain "key" is intentionally excluded.
var secretLikeTokensFitness = []string{
	"password", "secret", "token", "credential",
	"apikey", "accesskey", "secretkey", "privatekey", "private_key",
	"connectionstring", "secretvalue",
}

// approvedTypedRefTargets are the only $ref values permitted for reference
// fields in FEATURE-0012 schemas (relative forms under api/schemas).
var approvedTypedRefTargets = map[string]struct{}{
	"_common/typed-ref.json": {},
	"_common/scope-ref.json": {},
	"_common/owner-ref.json": {},
	"typed-ref.json":         {},
	"scope-ref.json":         {},
	"owner-ref.json":         {},
}

// coreCustomerGoBindings maps core/customer (non-adapter/plugin) schemas to
// their conformance Go types for SDK/native-field scanning.
var coreCustomerGoBindings = []struct {
	SchemaFile string
	GoType     reflect.Type
}{
	{"project.json", reflect.TypeOf(Project{})},
	{"resource-pool.json", reflect.TypeOf(ResourcePool{})},
	{"placement-evaluation-request.json", reflect.TypeOf(PlacementEvaluationRequest{})},
	{"audit-event.json", reflect.TypeOf(AuditEvent{})},
}

// CheckNoProviderSDKInCoreCustomer implements F12-VERIFY-001 check 2:
// no core/customer schema imports or embeds provider SDK/native types
// (F12-SEC-006, F12-BOUNDARY-001, F12-R04).
//
// moduleRoot is the repository root (directory containing go.mod). Schemas are
// loaded from moduleRoot/api/schemas.
func CheckNoProviderSDKInCoreCustomer(moduleRoot string) []FitnessFinding {
	var findings []FitnessFinding
	schemasRoot := filepath.Join(moduleRoot, CanonicalSchemasDir)

	for _, pkg := range fitnessCoreGrammarPackages {
		dir := filepath.Join(moduleRoot, "internal", pkg)
		imports, err := listGoImportPaths(dir)
		if err != nil {
			findings = append(findings, FitnessFinding{
				Check:   FitnessCheckNoProviderSDKInCoreCustomer,
				Schema:  "internal/" + pkg,
				Path:    "/",
				Code:    CodeFitnessSchemaLoadFailed,
				Message: err.Error(),
			})
			continue
		}
		for _, imp := range imports {
			if isFitnessProviderSDK(imp) {
				findings = append(findings, FitnessFinding{
					Check:   FitnessCheckNoProviderSDKInCoreCustomer,
					Schema:  "internal/" + pkg,
					Path:    "/",
					Code:    CodeFitnessProviderSDKImport,
					Message: fmt.Sprintf("imports provider SDK %q", imp),
				})
			}
		}
	}

	for _, name := range externalCanonicalSchemaFiles {
		schemaID := CanonicalSchemasDir + "/" + name
		body, err := os.ReadFile(filepath.Join(schemasRoot, name))
		if err != nil {
			findings = append(findings, FitnessFinding{
				Check:   FitnessCheckNoProviderSDKInCoreCustomer,
				Schema:  schemaID,
				Path:    "/",
				Code:    CodeFitnessSchemaLoadFailed,
				Message: err.Error(),
			})
			continue
		}
		meta, issues := apischema.ReadAnnotations(body)
		if len(issues) != 0 {
			continue // annotation failures owned by check 1
		}
		if !isCoreCustomerBoundary(meta.Boundary) {
			continue
		}
		for _, prop := range collectSchemaPropertyNamesFitness(body) {
			if token := bannedCoreNativeFieldToken(prop); token != "" {
				findings = append(findings, FitnessFinding{
					Check:   FitnessCheckNoProviderSDKInCoreCustomer,
					Schema:  schemaID,
					Path:    "/properties/" + escapeFitnessPointerToken(prop),
					Code:    CodeFitnessProviderNativeField,
					Message: fmt.Sprintf("core/customer schema property embeds provider-native token %q", token),
				})
			}
		}
	}

	commonRoot := filepath.Join(schemasRoot, "_common")
	for _, name := range commonSubSchemaFiles {
		schemaID := CanonicalSchemasDir + "/_common/" + name
		body, err := os.ReadFile(filepath.Join(commonRoot, name))
		if err != nil {
			findings = append(findings, FitnessFinding{
				Check:   FitnessCheckNoProviderSDKInCoreCustomer,
				Schema:  schemaID,
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
					Schema:  schemaID,
					Path:    "/properties/" + escapeFitnessPointerToken(prop),
					Code:    CodeFitnessProviderNativeField,
					Message: fmt.Sprintf("_common schema property embeds provider-native token %q", token),
				})
			}
		}
	}

	for _, b := range coreCustomerGoBindings {
		if err := checkGoTypeNoProviderSDK(b.GoType); err != nil {
			findings = append(findings, FitnessFinding{
				Check:   FitnessCheckNoProviderSDKInCoreCustomer,
				Schema:  CanonicalSchemasDir + "/" + b.SchemaFile,
				Path:    "/",
				Code:    CodeFitnessProviderSDKGoType,
				Message: err.Error(),
			})
		}
		if err := checkGoTypeNoBannedNativeFields(b.GoType); err != nil {
			findings = append(findings, FitnessFinding{
				Check:   FitnessCheckNoProviderSDKInCoreCustomer,
				Schema:  CanonicalSchemasDir + "/" + b.SchemaFile,
				Path:    "/",
				Code:    CodeFitnessProviderNativeField,
				Message: err.Error(),
			})
		}
	}

	return sortFindings(findings)
}

// CheckReferencesConstrainKindsAndScopes implements F12-VERIFY-001 check 5:
// references constrain kinds and scopes (F12-REF-001/002, F12-SCOPE-002).
func CheckReferencesConstrainKindsAndScopes(schemasRoot string) []FitnessFinding {
	var findings []FitnessFinding

	findings = append(findings, checkScopeRefKindEnum(schemasRoot)...)
	findings = append(findings, checkTypedRefRequiredFields(schemasRoot)...)
	findings = append(findings, checkObjectMetaScopeRefTarget(schemasRoot)...)
	findings = append(findings, checkExternalSchemaAllowedScopes(schemasRoot)...)
	findings = append(findings, checkSchemaReferenceFieldTargets(schemasRoot)...)
	findings = append(findings, checkApirefConstraintBehavior()...)

	return sortFindings(findings)
}

// CheckCrossTenantAccessNoExistenceDisclosure implements F12-VERIFY-001 check 6:
// cross-tenant access fails without existence disclosure via SafeDenial
// (owned by apivalid/authz.go; F12-SEC-004, F12-SCOPE-002).
//
// This check proves path/response equivalence of DenyNotDisclosed vs absent
// (identical 404 RESOURCE_NOT_FOUND) and DenyKnown → 403. It does not claim
// perfect constant-time execution. Request/operation correlation fields
// (requestId/instance) are intentionally left empty on SafeDenial so callers
// may attach them without changing the stable denial shape; secrets and
// inaccessible resource details MUST NOT be added.
func CheckCrossTenantAccessNoExistenceDisclosure() []FitnessFinding {
	var findings []FitnessFinding

	denied := apivalid.SafeDenial(apivalid.DenyNotDisclosed)
	absent := apiproblem.New(apiproblem.CodeResourceNotFound)
	if denied == nil || absent == nil {
		findings = append(findings, FitnessFinding{
			Check:   FitnessCheckCrossTenantNoExistenceDisclosure,
			Schema:  "apivalid.SafeDenial",
			Path:    "/",
			Code:    CodeFitnessSafeDenialMismatch,
			Message: "SafeDenial/absent Problem must be non-nil",
		})
		return sortFindings(findings)
	}

	deniedJSON, err1 := json.Marshal(denied)
	absentJSON, err2 := json.Marshal(absent)
	if err1 != nil || err2 != nil {
		findings = append(findings, FitnessFinding{
			Check:   FitnessCheckCrossTenantNoExistenceDisclosure,
			Schema:  "apivalid.SafeDenial",
			Path:    "/",
			Code:    CodeFitnessSafeDenialMismatch,
			Message: "failed to marshal SafeDenial/absent Problem",
		})
		return sortFindings(findings)
	}
	if !bytes.Equal(deniedJSON, absentJSON) {
		findings = append(findings, FitnessFinding{
			Check:   FitnessCheckCrossTenantNoExistenceDisclosure,
			Schema:  "apivalid.SafeDenial",
			Path:    "/",
			Code:    CodeFitnessSafeDenialMismatch,
			Message: "DenyNotDisclosed must be byte-identical to absent RESOURCE_NOT_FOUND",
		})
	}
	if denied.Status != 404 || denied.Code != apiproblem.CodeResourceNotFound {
		findings = append(findings, FitnessFinding{
			Check:   FitnessCheckCrossTenantNoExistenceDisclosure,
			Schema:  "apivalid.SafeDenial",
			Path:    "/",
			Code:    CodeFitnessSafeDenialMismatch,
			Message: fmt.Sprintf("DenyNotDisclosed status/code = %d %q, want 404 RESOURCE_NOT_FOUND", denied.Status, denied.Code),
		})
	}
	if len(denied.Violations) != 0 {
		findings = append(findings, FitnessFinding{
			Check:   FitnessCheckCrossTenantNoExistenceDisclosure,
			Schema:  "apivalid.SafeDenial",
			Path:    "/",
			Code:    CodeFitnessSafeDenialDiscloses,
			Message: "SafeDenial must not disclose mismatch/existence violations",
		})
	}
	if hit := secretLikeHitFitness(denied.Detail) + secretLikeHitFitness(denied.Title); hit != "" {
		findings = append(findings, FitnessFinding{
			Check:   FitnessCheckCrossTenantNoExistenceDisclosure,
			Schema:  "apivalid.SafeDenial",
			Path:    "/",
			Code:    CodeFitnessSecretLikeValue,
			Message: "SafeDenial must not embed secret-like text",
		})
	}

	known := apivalid.SafeDenial(apivalid.DenyKnown)
	if known == nil || known.Status != 403 || known.Code != apiproblem.CodeAuthorizationDenied {
		status := 0
		code := apiproblem.ErrorCode("")
		if known != nil {
			status = known.Status
			code = known.Code
		}
		findings = append(findings, FitnessFinding{
			Check:   FitnessCheckCrossTenantNoExistenceDisclosure,
			Schema:  "apivalid.SafeDenial",
			Path:    "/",
			Code:    CodeFitnessSafeDenialMismatch,
			Message: fmt.Sprintf("DenyKnown status/code = %d %q, want 403 AUTHORIZATION_DENIED", status, code),
		})
	}

	if allow := apivalid.SafeDenial(apivalid.Allow); allow != nil {
		findings = append(findings, FitnessFinding{
			Check:   FitnessCheckCrossTenantNoExistenceDisclosure,
			Schema:  "apivalid.SafeDenial",
			Path:    "/",
			Code:    CodeFitnessSafeDenialMismatch,
			Message: "Allow must map to nil Problem",
		})
	}

	return sortFindings(findings)
}

// CheckNoRawSecretLikeValues implements F12-VERIFY-001 check 7:
// raw secret-like values are prohibited from metadata/status/errors
// (F12-SEC-003, F12-R10). Scans positive fixtures and _common problem/
// object-meta schemas. Finding messages never echo raw secret material.
func CheckNoRawSecretLikeValues(moduleRoot string) []FitnessFinding {
	var findings []FitnessFinding
	fixturesDir := filepath.Join(moduleRoot, ConformanceFixturesDir)

	entries, err := os.ReadDir(fixturesDir)
	if err != nil {
		return []FitnessFinding{{
			Check:   FitnessCheckNoRawSecretLikeValues,
			Schema:  ConformanceFixturesDir,
			Path:    "/",
			Code:    CodeFitnessSchemaLoadFailed,
			Message: err.Error(),
		}}
	}
	for _, ent := range entries {
		if ent.IsDir() {
			continue
		}
		name := ent.Name()
		if !strings.HasSuffix(name, ".json") {
			continue
		}
		raw, err := os.ReadFile(filepath.Join(fixturesDir, name))
		if err != nil {
			findings = append(findings, FitnessFinding{
				Check:   FitnessCheckNoRawSecretLikeValues,
				Schema:  ConformanceFixturesDir + "/" + name,
				Path:    "/",
				Code:    CodeFitnessSchemaLoadFailed,
				Message: "fixture read failed",
			})
			continue
		}
		if hit, path := findSecretLikeInObjectJSON(raw); hit != "" {
			findings = append(findings, FitnessFinding{
				Check:   FitnessCheckNoRawSecretLikeValues,
				Schema:  ConformanceFixturesDir + "/" + name,
				Path:    path,
				Code:    CodeFitnessSecretLikeValue,
				Message: fmt.Sprintf("secret-like token %q prohibited in metadata/status/errors", hit),
			})
		}
	}

	schemasRoot := filepath.Join(moduleRoot, CanonicalSchemasDir)
	for _, name := range []string{"_common/object-meta.json", "_common/problem.json", "_common/condition.json"} {
		schemaID := CanonicalSchemasDir + "/" + name
		body, err := os.ReadFile(filepath.Join(schemasRoot, filepath.FromSlash(name)))
		if err != nil {
			findings = append(findings, FitnessFinding{
				Check:   FitnessCheckNoRawSecretLikeValues,
				Schema:  schemaID,
				Path:    "/",
				Code:    CodeFitnessSchemaLoadFailed,
				Message: err.Error(),
			})
			continue
		}
		for _, prop := range collectSchemaPropertyNamesFitness(body) {
			if hit := secretLikeHitFitness(prop); hit != "" {
				findings = append(findings, FitnessFinding{
					Check:   FitnessCheckNoRawSecretLikeValues,
					Schema:  schemaID,
					Path:    "/properties/" + escapeFitnessPointerToken(prop),
					Code:    CodeFitnessSecretLikeValue,
					Message: fmt.Sprintf("schema property name embeds secret-like token %q", hit),
				})
			}
		}
	}

	return sortFindings(findings)
}

// CheckObservationProvenanceAndFreshness implements F12-VERIFY-001 check 8:
// externally sourced observations include provenance and freshness
// (F12-SEC-005, F12-STATUS-003).
func CheckObservationProvenanceAndFreshness(schemasRoot string) []FitnessFinding {
	var findings []FitnessFinding
	observedCount := 0

	for _, name := range externalCanonicalSchemaFiles {
		schemaID := CanonicalSchemasDir + "/" + name
		body, err := os.ReadFile(filepath.Join(schemasRoot, name))
		if err != nil {
			findings = append(findings, FitnessFinding{
				Check:   FitnessCheckObservationProvenanceAndFreshness,
				Schema:  schemaID,
				Path:    "/",
				Code:    CodeFitnessSchemaLoadFailed,
				Message: err.Error(),
			})
			continue
		}
		meta, issues := apischema.ReadAnnotations(body)
		if len(issues) != 0 {
			continue
		}
		if meta.Profile != apimeta.ProfileObservedExternalResource {
			continue
		}
		observedCount++
		findings = append(findings, observationProvenanceFreshnessFindings(schemaID, body)...)
	}

	if observedCount == 0 {
		findings = append(findings, FitnessFinding{
			Check:   FitnessCheckObservationProvenanceAndFreshness,
			Schema:  CanonicalSchemasDir,
			Path:    "/",
			Code:    CodeFitnessProvenanceMissing,
			Message: "no ObservedExternalResource schema found; Matrix D requires at least one",
		})
	}

	return sortFindings(findings)
}

func observationProvenanceFreshnessFindings(schemaID string, body []byte) []FitnessFinding {
	var root map[string]any
	if err := json.Unmarshal(body, &root); err != nil {
		return []FitnessFinding{{
			Check:   FitnessCheckObservationProvenanceAndFreshness,
			Schema:  schemaID,
			Path:    "/",
			Code:    apischema.CodeMalformedSchema,
			Message: "schema document is not valid JSON",
		}}
	}
	var findings []FitnessFinding
	props, _ := root["properties"].(map[string]any)
	required := stringSetFromAny(root["required"])

	prov, hasProv := props["provenance"].(map[string]any)
	if !hasProv {
		findings = append(findings, FitnessFinding{
			Check:   FitnessCheckObservationProvenanceAndFreshness,
			Schema:  schemaID,
			Path:    "/properties/provenance",
			Code:    CodeFitnessProvenanceMissing,
			Message: "ObservedExternalResource must declare provenance",
		})
	} else {
		if !required["provenance"] {
			findings = append(findings, FitnessFinding{
				Check:   FitnessCheckObservationProvenanceAndFreshness,
				Schema:  schemaID,
				Path:    "/required",
				Code:    CodeFitnessProvenanceMissing,
				Message: "provenance must be required on ObservedExternalResource",
			})
		}
		provReq := stringSetFromAny(prov["required"])
		for _, field := range []string{"sourceRef", "observedAt"} {
			if !provReq[field] {
				findings = append(findings, FitnessFinding{
					Check:   FitnessCheckObservationProvenanceAndFreshness,
					Schema:  schemaID,
					Path:    "/properties/provenance/required",
					Code:    CodeFitnessProvenanceMissing,
					Message: fmt.Sprintf("provenance.%s must be required", field),
				})
			}
		}
		provProps, _ := prov["properties"].(map[string]any)
		if _, ok := provProps["sourceRef"]; !ok {
			findings = append(findings, FitnessFinding{
				Check:   FitnessCheckObservationProvenanceAndFreshness,
				Schema:  schemaID,
				Path:    "/properties/provenance/properties/sourceRef",
				Code:    CodeFitnessProvenanceMissing,
				Message: "provenance.sourceRef property missing",
			})
		}
		if _, ok := provProps["observedAt"]; !ok {
			findings = append(findings, FitnessFinding{
				Check:   FitnessCheckObservationProvenanceAndFreshness,
				Schema:  schemaID,
				Path:    "/properties/provenance/properties/observedAt",
				Code:    CodeFitnessProvenanceMissing,
				Message: "provenance.observedAt property missing",
			})
		}
	}

	fresh, hasFresh := props["freshness"].(map[string]any)
	if !hasFresh {
		findings = append(findings, FitnessFinding{
			Check:   FitnessCheckObservationProvenanceAndFreshness,
			Schema:  schemaID,
			Path:    "/properties/freshness",
			Code:    CodeFitnessFreshnessMissing,
			Message: "ObservedExternalResource must declare freshness",
		})
	} else {
		if !required["freshness"] {
			findings = append(findings, FitnessFinding{
				Check:   FitnessCheckObservationProvenanceAndFreshness,
				Schema:  schemaID,
				Path:    "/required",
				Code:    CodeFitnessFreshnessMissing,
				Message: "freshness must be required on ObservedExternalResource",
			})
		}
		freshReq := stringSetFromAny(fresh["required"])
		if !freshReq["freshnessState"] {
			findings = append(findings, FitnessFinding{
				Check:   FitnessCheckObservationProvenanceAndFreshness,
				Schema:  schemaID,
				Path:    "/properties/freshness/required",
				Code:    CodeFitnessFreshnessMissing,
				Message: "freshness.freshnessState must be required",
			})
		}
		freshProps, _ := fresh["properties"].(map[string]any)
		if _, ok := freshProps["freshnessState"]; !ok {
			findings = append(findings, FitnessFinding{
				Check:   FitnessCheckObservationProvenanceAndFreshness,
				Schema:  schemaID,
				Path:    "/properties/freshness/properties/freshnessState",
				Code:    CodeFitnessFreshnessMissing,
				Message: "freshness.freshnessState property missing",
			})
		}
	}

	return findings
}

func checkScopeRefKindEnum(schemasRoot string) []FitnessFinding {
	schemaID := CanonicalSchemasDir + "/_common/scope-ref.json"
	body, err := os.ReadFile(filepath.Join(schemasRoot, "_common", "scope-ref.json"))
	if err != nil {
		return []FitnessFinding{{
			Check:   FitnessCheckReferencesConstrainKindsAndScopes,
			Schema:  schemaID,
			Path:    "/",
			Code:    CodeFitnessSchemaLoadFailed,
			Message: err.Error(),
		}}
	}
	var root map[string]any
	if err := json.Unmarshal(body, &root); err != nil {
		return []FitnessFinding{{
			Check:   FitnessCheckReferencesConstrainKindsAndScopes,
			Schema:  schemaID,
			Path:    "/",
			Code:    apischema.CodeMalformedSchema,
			Message: "schema document is not valid JSON",
		}}
	}
	kind, _ := root["properties"].(map[string]any)["kind"].(map[string]any)
	rawEnum, _ := kind["enum"].([]any)
	want := apimeta.AllScopeKinds()
	if len(rawEnum) != len(want) {
		return []FitnessFinding{{
			Check:   FitnessCheckReferencesConstrainKindsAndScopes,
			Schema:  schemaID,
			Path:    "/properties/kind/enum",
			Code:    CodeFitnessScopeKindUnconstrained,
			Message: fmt.Sprintf("scope-ref kind enum len=%d want %d Matrix B values", len(rawEnum), len(want)),
		}}
	}
	var findings []FitnessFinding
	for i, w := range want {
		got, _ := rawEnum[i].(string)
		if got != string(w) {
			findings = append(findings, FitnessFinding{
				Check:   FitnessCheckReferencesConstrainKindsAndScopes,
				Schema:  schemaID,
				Path:    fmt.Sprintf("/properties/kind/enum/%d", i),
				Code:    CodeFitnessScopeKindUnconstrained,
				Message: fmt.Sprintf("scope-ref kind enum[%d]=%q want %q", i, got, w),
			})
		}
	}
	return findings
}

func checkTypedRefRequiredFields(schemasRoot string) []FitnessFinding {
	schemaID := CanonicalSchemasDir + "/_common/typed-ref.json"
	body, err := os.ReadFile(filepath.Join(schemasRoot, "_common", "typed-ref.json"))
	if err != nil {
		return []FitnessFinding{{
			Check:   FitnessCheckReferencesConstrainKindsAndScopes,
			Schema:  schemaID,
			Path:    "/",
			Code:    CodeFitnessSchemaLoadFailed,
			Message: err.Error(),
		}}
	}
	var root map[string]any
	if err := json.Unmarshal(body, &root); err != nil {
		return []FitnessFinding{{
			Check:   FitnessCheckReferencesConstrainKindsAndScopes,
			Schema:  schemaID,
			Path:    "/",
			Code:    apischema.CodeMalformedSchema,
			Message: "schema document is not valid JSON",
		}}
	}
	req := stringSetFromAny(root["required"])
	var findings []FitnessFinding
	for _, field := range []string{"apiVersion", "kind", "name"} {
		if !req[field] {
			findings = append(findings, FitnessFinding{
				Check:   FitnessCheckReferencesConstrainKindsAndScopes,
				Schema:  schemaID,
				Path:    "/required",
				Code:    CodeFitnessTypedRefIncomplete,
				Message: fmt.Sprintf("typed-ref must require %s", field),
			})
		}
	}
	props, _ := root["properties"].(map[string]any)
	kind, _ := props["kind"].(map[string]any)
	if _, ok := kind["pattern"]; !ok {
		if _, hasEnum := kind["enum"]; !hasEnum {
			findings = append(findings, FitnessFinding{
				Check:   FitnessCheckReferencesConstrainKindsAndScopes,
				Schema:  schemaID,
				Path:    "/properties/kind",
				Code:    CodeFitnessTypedRefIncomplete,
				Message: "typed-ref kind must constrain values via pattern or enum",
			})
		}
	}
	return findings
}

func checkObjectMetaScopeRefTarget(schemasRoot string) []FitnessFinding {
	schemaID := CanonicalSchemasDir + "/_common/object-meta.json"
	body, err := os.ReadFile(filepath.Join(schemasRoot, "_common", "object-meta.json"))
	if err != nil {
		return []FitnessFinding{{
			Check:   FitnessCheckReferencesConstrainKindsAndScopes,
			Schema:  schemaID,
			Path:    "/",
			Code:    CodeFitnessSchemaLoadFailed,
			Message: err.Error(),
		}}
	}
	var root map[string]any
	if err := json.Unmarshal(body, &root); err != nil {
		return []FitnessFinding{{
			Check:   FitnessCheckReferencesConstrainKindsAndScopes,
			Schema:  schemaID,
			Path:    "/",
			Code:    apischema.CodeMalformedSchema,
			Message: "schema document is not valid JSON",
		}}
	}
	scopeRef, _ := root["properties"].(map[string]any)["scopeRef"].(map[string]any)
	ref, _ := scopeRef["$ref"].(string)
	if _, ok := approvedTypedRefTargets[ref]; !ok || !strings.Contains(ref, "scope-ref") {
		return []FitnessFinding{{
			Check:   FitnessCheckReferencesConstrainKindsAndScopes,
			Schema:  schemaID,
			Path:    "/properties/scopeRef/$ref",
			Code:    CodeFitnessRefNotConstrained,
			Message: fmt.Sprintf("object-meta.scopeRef must $ref scope-ref.json, got %q", ref),
		}}
	}
	return nil
}

func checkExternalSchemaAllowedScopes(schemasRoot string) []FitnessFinding {
	var findings []FitnessFinding
	for _, name := range externalCanonicalSchemaFiles {
		schemaID := CanonicalSchemasDir + "/" + name
		body, err := os.ReadFile(filepath.Join(schemasRoot, name))
		if err != nil {
			findings = append(findings, FitnessFinding{
				Check:   FitnessCheckReferencesConstrainKindsAndScopes,
				Schema:  schemaID,
				Path:    "/",
				Code:    CodeFitnessSchemaLoadFailed,
				Message: err.Error(),
			})
			continue
		}
		meta, issues := apischema.ReadAnnotations(body)
		if len(issues) != 0 {
			continue
		}
		if len(meta.AllowedScopes) == 0 {
			findings = append(findings, FitnessFinding{
				Check:   FitnessCheckReferencesConstrainKindsAndScopes,
				Schema:  schemaID,
				Path:    "/" + apischema.ExtAllowedScopes,
				Code:    CodeFitnessAllowedScopesMissing,
				Message: "external schema must declare x-sovrunn-allowed-scopes for reference/scope constraints",
			})
			continue
		}
		for i, sk := range meta.AllowedScopes {
			if !sk.Valid() {
				findings = append(findings, FitnessFinding{
					Check:   FitnessCheckReferencesConstrainKindsAndScopes,
					Schema:  schemaID,
					Path:    fmt.Sprintf("/%s/%d", apischema.ExtAllowedScopes, i),
					Code:    CodeFitnessScopeKindUnconstrained,
					Message: fmt.Sprintf("allowed-scopes contains non-Matrix-B kind %q", sk),
				})
			}
		}
	}
	return findings
}

func checkSchemaReferenceFieldTargets(schemasRoot string) []FitnessFinding {
	var findings []FitnessFinding
	scan := func(schemaID string, body []byte) {
		var root any
		if err := json.Unmarshal(body, &root); err != nil {
			findings = append(findings, FitnessFinding{
				Check:   FitnessCheckReferencesConstrainKindsAndScopes,
				Schema:  schemaID,
				Path:    "/",
				Code:    apischema.CodeMalformedSchema,
				Message: "schema document is not valid JSON",
			})
			return
		}
		walkReferenceFieldTargets(root, "", schemaID, &findings)
	}

	for _, name := range externalCanonicalSchemaFiles {
		schemaID := CanonicalSchemasDir + "/" + name
		body, err := os.ReadFile(filepath.Join(schemasRoot, name))
		if err != nil {
			findings = append(findings, FitnessFinding{
				Check:   FitnessCheckReferencesConstrainKindsAndScopes,
				Schema:  schemaID,
				Path:    "/",
				Code:    CodeFitnessSchemaLoadFailed,
				Message: err.Error(),
			})
			continue
		}
		scan(schemaID, body)
	}
	commonRoot := filepath.Join(schemasRoot, "_common")
	for _, name := range []string{"object-meta.json", "owner-ref.json", "scope-ref.json"} {
		schemaID := CanonicalSchemasDir + "/_common/" + name
		body, err := os.ReadFile(filepath.Join(commonRoot, name))
		if err != nil {
			findings = append(findings, FitnessFinding{
				Check:   FitnessCheckReferencesConstrainKindsAndScopes,
				Schema:  schemaID,
				Path:    "/",
				Code:    CodeFitnessSchemaLoadFailed,
				Message: err.Error(),
			})
			continue
		}
		scan(schemaID, body)
	}
	return findings
}

func walkReferenceFieldTargets(node any, path, schemaID string, findings *[]FitnessFinding) {
	obj, ok := node.(map[string]any)
	if !ok {
		return
	}
	if rawProps, hasProps := obj["properties"]; hasProps {
		props, ok := rawProps.(map[string]any)
		if !ok {
			return
		}
		names := make([]string, 0, len(props))
		for name := range props {
			names = append(names, name)
		}
		sort.Strings(names)
		for _, name := range names {
			propPath := joinFitnessPointer(joinFitnessPointer(path, "properties"), name)
			propSchema, ok := props[name].(map[string]any)
			if !ok {
				continue
			}
			if looksLikeReferencePropertyName(name) {
				if !referencePropertyIsConstrained(propSchema) {
					*findings = append(*findings, FitnessFinding{
						Check:   FitnessCheckReferencesConstrainKindsAndScopes,
						Schema:  schemaID,
						Path:    propPath,
						Code:    CodeFitnessRefNotConstrained,
						Message: fmt.Sprintf("reference field %q must $ref typed-ref/scope-ref/owner-ref (not a bare string)", name),
					})
				}
			}
			walkReferenceFieldTargets(propSchema, propPath, schemaID, findings)
		}
	}
	if rawItems, hasItems := obj["items"]; hasItems {
		walkReferenceFieldTargets(rawItems, joinFitnessPointer(path, "items"), schemaID, findings)
	}
}

func looksLikeReferencePropertyName(name string) bool {
	switch name {
	case "scopeRef", "ownerRef", "targetRef", "sourceRef", "subjectRef",
		"actorRef", "operationRef", "resourcePoolRef", "projectRef",
		"serviceClassRef", "secretRef":
		return true
	}
	return strings.HasSuffix(name, "Ref")
}

func referencePropertyIsConstrained(propSchema map[string]any) bool {
	if ref, ok := propSchema["$ref"].(string); ok {
		if _, approved := approvedTypedRefTargets[ref]; approved {
			return true
		}
		// Relative refs that clearly target the approved common files.
		base := filepath.Base(ref)
		if _, approved := approvedTypedRefTargets[base]; approved {
			return true
		}
	}
	// Inline object with required apiVersion/kind/name is also constrained.
	if t, _ := propSchema["type"].(string); t == "object" {
		req := stringSetFromAny(propSchema["required"])
		if req["apiVersion"] && req["kind"] && req["name"] {
			return true
		}
	}
	return false
}

func checkApirefConstraintBehavior() []FitnessFinding {
	var findings []FitnessFinding
	// AllowedKinds-only constraint: AllowedScopes is omitted so a non-scope
	// resource kind is not also evaluated as a ScopeKind.
	c := apiref.Constraint{
		AllowedKinds: []string{"ResourcePool"},
		Direction:    apiref.DirectionOutbound,
	}

	okRef := apiref.TypedRef{
		APIVersion: "core.sovrunn.io/v1alpha1",
		Kind:       "ResourcePool",
		Name:       "pool-a",
	}
	if issues := c.ValidateRef(okRef, "/spec/resourcePoolRef"); len(issues) != 0 {
		findings = append(findings, FitnessFinding{
			Check:   FitnessCheckReferencesConstrainKindsAndScopes,
			Schema:  "apiref.Constraint",
			Path:    "/spec/resourcePoolRef",
			Code:    CodeFitnessRefConstraintIneffective,
			Message: "allowed kind must pass ValidateRef",
		})
	}

	badKind := okRef
	badKind.Kind = "Project"
	if issues := c.ValidateRef(badKind, "/spec/resourcePoolRef"); !hasRefIssueCode(issues, apiref.CodeKindNotAllowed) {
		findings = append(findings, FitnessFinding{
			Check:   FitnessCheckReferencesConstrainKindsAndScopes,
			Schema:  "apiref.Constraint",
			Path:    "/spec/resourcePoolRef/kind",
			Code:    CodeFitnessRefConstraintIneffective,
			Message: "disallowed kind must fail ValidateRef with KIND_NOT_ALLOWED",
		})
	}

	scopeC := apiref.Constraint{
		AllowedScopes: []apimeta.ScopeKind{apimeta.ScopeTenant},
		Direction:     apiref.DirectionOutbound,
	}
	badScope := apiref.TypedRef{
		APIVersion: "core.sovrunn.io/v1alpha1",
		Kind:       string(apimeta.ScopeProvider),
		Name:       "prov-a",
	}
	if issues := scopeC.ValidateRef(badScope, "/metadata/scopeRef"); !hasRefIssueCode(issues, apiref.CodeScopeNotAllowed) {
		findings = append(findings, FitnessFinding{
			Check:   FitnessCheckReferencesConstrainKindsAndScopes,
			Schema:  "apiref.Constraint",
			Path:    "/metadata/scopeRef/kind",
			Code:    CodeFitnessRefConstraintIneffective,
			Message: "disallowed scope kind must fail ValidateRef with SCOPE_NOT_ALLOWED",
		})
	}

	native := apiref.TypedRef{
		APIVersion: "core.sovrunn.io/v1alpha1",
		Kind:       "AWS::RDS::DBInstance",
		Name:       "db-a",
	}
	nativeC := apiref.Constraint{AllowedKinds: []string{"AWS::RDS::DBInstance", "ResourcePool"}}
	if issues := nativeC.ValidateRef(native, "/spec/resourcePoolRef"); !hasRefIssueCode(issues, apiref.CodeProviderNativeID) {
		findings = append(findings, FitnessFinding{
			Check:   FitnessCheckReferencesConstrainKindsAndScopes,
			Schema:  "apiref.Constraint",
			Path:    "/spec/resourcePoolRef",
			Code:    CodeFitnessRefConstraintIneffective,
			Message: "provider-native kind must fail ValidateRef with PROVIDER_NATIVE_ID",
		})
	}

	return findings
}

func hasRefIssueCode(issues []apiref.RefIssue, code string) bool {
	for _, issue := range issues {
		if issue.Code == code {
			return true
		}
	}
	return false
}

func isCoreCustomerBoundary(b apimeta.Boundary) bool {
	switch b {
	case apimeta.BoundaryCustomerFacing,
		apimeta.BoundaryOperatorFacing,
		apimeta.BoundaryInternalEngineFacing,
		apimeta.BoundaryGovernanceOnly:
		return true
	default:
		return false
	}
}

func isFitnessProviderSDK(importPath string) bool {
	for _, p := range providerSDKImportPrefixes {
		if strings.HasPrefix(importPath, p) {
			return true
		}
	}
	return false
}

func listGoImportPaths(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	fset := token.NewFileSet()
	var out []string
	for _, ent := range entries {
		if ent.IsDir() || !strings.HasSuffix(ent.Name(), ".go") {
			continue
		}
		path := filepath.Join(dir, ent.Name())
		file, err := parser.ParseFile(fset, path, nil, parser.ImportsOnly)
		if err != nil {
			return nil, fmt.Errorf("parse %s: %w", path, err)
		}
		for _, is := range file.Imports {
			if is.Path == nil {
				continue
			}
			out = append(out, strings.Trim(is.Path.Value, `"`))
		}
	}
	return out, nil
}

func bannedCoreNativeFieldToken(name string) string {
	norm := normalizeFitnessIdent(name)
	for _, token := range bannedCoreNativeFieldTokens {
		if norm == token || strings.Contains(norm, token) {
			return token
		}
	}
	return ""
}

func normalizeFitnessIdent(s string) string {
	var b strings.Builder
	for _, r := range strings.ToLower(s) {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			b.WriteRune(r)
		}
	}
	return b.String()
}

func collectSchemaPropertyNamesFitness(schema []byte) []string {
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

func checkGoTypeNoProviderSDK(t reflect.Type) error {
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
		if pkg := tt.PkgPath(); pkg != "" && isFitnessProviderSDK(pkg) {
			return fmt.Errorf("Go type %s embeds provider SDK package %q", tt.String(), pkg)
		}
		switch tt.Kind() {
		case reflect.Struct:
			for i := 0; i < tt.NumField(); i++ {
				if err := walk(tt.Field(i).Type); err != nil {
					return err
				}
			}
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
	return walk(t)
}

func checkGoTypeNoBannedNativeFields(t reflect.Type) error {
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
			if token := bannedCoreNativeFieldToken(name); token != "" {
				return fmt.Errorf("core Go field %q embeds provider-native token %q", name, token)
			}
			if err := walk(f.Type); err != nil {
				return err
			}
		}
		return nil
	}
	return walk(t)
}

func stringSetFromAny(v any) map[string]bool {
	out := map[string]bool{}
	arr, ok := v.([]any)
	if !ok {
		return out
	}
	for _, item := range arr {
		if s, ok := item.(string); ok {
			out[s] = true
		}
	}
	return out
}

// findSecretLikeInObjectJSON scans metadata.labels/annotations, status maps,
// and top-level error/problem objects for secret-like tokens. Returns the
// token hit and a JSON Pointer path. Values are not echoed in findings.
func findSecretLikeInObjectJSON(raw []byte) (hit, path string) {
	var top map[string]json.RawMessage
	if err := json.Unmarshal(raw, &top); err != nil {
		return "", ""
	}
	if metaRaw, ok := top["metadata"]; ok {
		if h, p := scanSecretLikeMapField(metaRaw, "/metadata"); h != "" {
			return h, p
		}
	}
	if statusRaw, ok := top["status"]; ok {
		if h, p := scanSecretLikeNested(statusRaw, "/status"); h != "" {
			return h, p
		}
	}
	for _, key := range []string{"error", "problem", "detail", "message"} {
		if fieldRaw, ok := top[key]; ok {
			if h := secretLikeHitFitness(string(fieldRaw)); h != "" {
				return h, "/" + key
			}
		}
	}
	return "", ""
}

func scanSecretLikeMapField(metaRaw json.RawMessage, base string) (string, string) {
	var meta map[string]json.RawMessage
	if err := json.Unmarshal(metaRaw, &meta); err != nil {
		return "", ""
	}
	for _, field := range []string{"labels", "annotations"} {
		fieldRaw, ok := meta[field]
		if !ok {
			continue
		}
		var kv map[string]string
		if err := json.Unmarshal(fieldRaw, &kv); err != nil {
			continue
		}
		keys := make([]string, 0, len(kv))
		for k := range kv {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			if hit := secretLikeHitFitness(k); hit != "" {
				return hit, base + "/" + field + "/" + escapeFitnessPointerToken(k)
			}
			if hit := secretLikeHitFitness(kv[k]); hit != "" {
				return hit, base + "/" + field + "/" + escapeFitnessPointerToken(k)
			}
		}
	}
	return "", ""
}

func scanSecretLikeNested(raw json.RawMessage, base string) (string, string) {
	var obj map[string]any
	if err := json.Unmarshal(raw, &obj); err != nil {
		return "", ""
	}
	keys := make([]string, 0, len(obj))
	for k := range obj {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		path := base + "/" + escapeFitnessPointerToken(k)
		if hit := secretLikeHitFitness(k); hit != "" {
			return hit, path
		}
		switch v := obj[k].(type) {
		case string:
			if hit := secretLikeHitFitness(v); hit != "" {
				return hit, path
			}
		case map[string]any:
			sub, err := json.Marshal(v)
			if err != nil {
				continue
			}
			if hit, p := scanSecretLikeNested(sub, path); hit != "" {
				return hit, p
			}
		}
	}
	return "", ""
}

func secretLikeHitFitness(s string) string {
	lower := strings.ToLower(s)
	for _, tok := range secretLikeTokensFitness {
		if strings.Contains(lower, tok) {
			return tok
		}
	}
	return ""
}
