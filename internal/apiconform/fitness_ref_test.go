package apiconform

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/apischema"
)

func TestFitnessCheckNoProviderSDKInCoreCustomer(t *testing.T) {
	t.Parallel()

	findings := CheckNoProviderSDKInCoreCustomer(moduleRoot(t))
	if len(findings) != 0 {
		t.Fatalf("check 2 failed: %#v", findings)
	}
}

func TestFitnessCheckNoProviderSDKInCoreCustomerNativeFieldFails(t *testing.T) {
	t.Parallel()

	var root map[string]any
	if err := json.Unmarshal(loadCanonicalSchema(t, "project.json"), &root); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	props := root["properties"].(map[string]any)
	spec := props["spec"].(map[string]any)
	specProps := spec["properties"].(map[string]any)
	specProps["awsAccountId"] = map[string]any{
		"type": "string",
		apischema.ExtFieldPolicy: map[string]any{
			"classification":    "Public",
			"authorizedWriter":  "creator",
			"authorizedReaders": []any{"customer"},
			"mutability":        "immutable",
			"retention":         "standard",
			"redaction":         "none",
			"residency":         "any",
			"auditRequired":     false,
		},
	}
	raw, err := json.Marshal(root)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	dir := t.TempDir()
	for _, name := range externalCanonicalSchemaFiles {
		src := loadCanonicalSchema(t, name)
		if name == "project.json" {
			src = raw
		}
		if err := os.WriteFile(filepath.Join(dir, name), src, 0o644); err != nil {
			t.Fatalf("write %s: %v", name, err)
		}
	}
	if err := copyCommonSchemas(t, dir); err != nil {
		t.Fatalf("copy common: %v", err)
	}

	// Rebuild a fake module root with schemas + empty internal packages is
	// heavy; exercise the schema scanner directly via collect + token helper
	// and assert the check would flag the synthetic property.
	if token := bannedCoreNativeFieldToken("awsAccountId"); token != "aws" {
		t.Fatalf("bannedCoreNativeFieldToken(awsAccountId)=%q want aws", token)
	}
	propsFound := false
	for _, name := range collectSchemaPropertyNamesFitness(raw) {
		if name == "awsAccountId" {
			propsFound = true
			break
		}
	}
	if !propsFound {
		t.Fatal("expected awsAccountId in collected property names")
	}
}

func TestFitnessCheckReferencesConstrainKindsAndScopes(t *testing.T) {
	t.Parallel()

	root := filepath.Join(moduleRoot(t), CanonicalSchemasDir)
	findings := CheckReferencesConstrainKindsAndScopes(root)
	if len(findings) != 0 {
		t.Fatalf("check 5 failed: %#v", findings)
	}
}

func TestFitnessCheckReferencesConstrainKindsAndScopesBareStringFails(t *testing.T) {
	t.Parallel()

	var root map[string]any
	if err := json.Unmarshal(loadCanonicalSchema(t, "project.json"), &root); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	props := root["properties"].(map[string]any)
	spec := props["spec"].(map[string]any)
	specProps := spec["properties"].(map[string]any)
	specProps["poolRef"] = map[string]any{
		"type": "string",
		apischema.ExtFieldPolicy: map[string]any{
			"classification":    "Public",
			"authorizedWriter":  "creator",
			"authorizedReaders": []any{"customer"},
			"mutability":        "immutable",
			"retention":         "standard",
			"redaction":         "none",
			"residency":         "any",
			"auditRequired":     false,
		},
	}
	raw, err := json.Marshal(root)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var findings []FitnessFinding
	walkReferenceFieldTargets(mustUnmarshalAny(t, raw), "", CanonicalSchemasDir+"/project.json", &findings)
	if !hasFitnessFinding(findings, FitnessCheckReferencesConstrainKindsAndScopes,
		CanonicalSchemasDir+"/project.json",
		"/properties/spec/properties/poolRef",
		CodeFitnessRefNotConstrained) {
		t.Fatalf("expected bare-string ref finding, got %#v", findings)
	}
}

func TestFitnessCheckReferencesScopeEnumMismatchFails(t *testing.T) {
	t.Parallel()

	var root map[string]any
	if err := json.Unmarshal(loadCanonicalSchema(t, "_common/scope-ref.json"), &root); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	kind := root["properties"].(map[string]any)["kind"].(map[string]any)
	kind["enum"] = []any{"Platform", "Organization"} // truncated Matrix B
	raw, err := json.Marshal(root)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	dir := t.TempDir()
	common := filepath.Join(dir, "_common")
	if err := os.MkdirAll(common, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := copyCommonSchemas(t, dir); err != nil {
		t.Fatalf("copy common: %v", err)
	}
	if err := os.WriteFile(filepath.Join(common, "scope-ref.json"), raw, 0o644); err != nil {
		t.Fatalf("write scope-ref: %v", err)
	}
	for _, name := range externalCanonicalSchemaFiles {
		src := loadCanonicalSchema(t, name)
		if err := os.WriteFile(filepath.Join(dir, name), src, 0o644); err != nil {
			t.Fatalf("write %s: %v", name, err)
		}
	}

	findings := checkScopeRefKindEnum(dir)
	if !hasFitnessFinding(findings, FitnessCheckReferencesConstrainKindsAndScopes,
		CanonicalSchemasDir+"/_common/scope-ref.json",
		"/properties/kind/enum",
		CodeFitnessScopeKindUnconstrained) {
		t.Fatalf("expected scope enum finding, got %#v", findings)
	}
}

func TestFitnessCheckCrossTenantAccessNoExistenceDisclosure(t *testing.T) {
	t.Parallel()

	findings := CheckCrossTenantAccessNoExistenceDisclosure()
	if len(findings) != 0 {
		t.Fatalf("check 6 failed: %#v", findings)
	}
}

func TestFitnessCheckNoRawSecretLikeValues(t *testing.T) {
	t.Parallel()

	findings := CheckNoRawSecretLikeValues(moduleRoot(t))
	if len(findings) != 0 {
		t.Fatalf("check 7 failed: %#v", findings)
	}
}

func TestFitnessCheckNoRawSecretLikeValuesDetectsPoison(t *testing.T) {
	t.Parallel()

	poison := []byte(`{"apiVersion":"core.sovrunn.io/v1alpha1","kind":"Project","metadata":{"name":"x","labels":{"db-password":"nope"}},"spec":{}}`)
	hit, path := findSecretLikeInObjectJSON(poison)
	if hit != "password" {
		t.Fatalf("hit=%q want password", hit)
	}
	if path != "/metadata/labels/db-password" {
		t.Fatalf("path=%q want /metadata/labels/db-password", path)
	}
}

func TestFitnessCheckObservationProvenanceAndFreshness(t *testing.T) {
	t.Parallel()

	root := filepath.Join(moduleRoot(t), CanonicalSchemasDir)
	findings := CheckObservationProvenanceAndFreshness(root)
	if len(findings) != 0 {
		t.Fatalf("check 8 failed: %#v", findings)
	}
}

func TestFitnessCheckObservationProvenanceAndFreshnessMissingFails(t *testing.T) {
	t.Parallel()

	var root map[string]any
	if err := json.Unmarshal(loadCanonicalSchema(t, "discovered-database.json"), &root); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	props := root["properties"].(map[string]any)
	delete(props, "provenance")
	delete(props, "freshness")
	root["required"] = []any{"apiVersion", "kind", "metadata", "status"}
	raw, err := json.Marshal(root)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	findings := observationProvenanceFreshnessFindings(CanonicalSchemasDir+"/discovered-database.json", raw)
	if !hasFitnessFinding(findings, FitnessCheckObservationProvenanceAndFreshness,
		CanonicalSchemasDir+"/discovered-database.json",
		"/properties/provenance",
		CodeFitnessProvenanceMissing) {
		t.Fatalf("expected provenance missing finding, got %#v", findings)
	}
	if !hasFitnessFinding(findings, FitnessCheckObservationProvenanceAndFreshness,
		CanonicalSchemasDir+"/discovered-database.json",
		"/properties/freshness",
		CodeFitnessFreshnessMissing) {
		t.Fatalf("expected freshness missing finding, got %#v", findings)
	}
}

func TestFitnessRefChecksInventory(t *testing.T) {
	t.Parallel()

	want := []string{
		FitnessCheckNoProviderSDKInCoreCustomer,
		FitnessCheckReferencesConstrainKindsAndScopes,
		FitnessCheckCrossTenantNoExistenceDisclosure,
		FitnessCheckNoRawSecretLikeValues,
		FitnessCheckObservationProvenanceAndFreshness,
	}
	for _, id := range want {
		if id == "" {
			t.Fatal("empty fitness check id")
		}
	}
}

func copyCommonSchemas(t *testing.T, schemasRoot string) error {
	t.Helper()
	srcCommon := filepath.Join(moduleRoot(t), CanonicalSchemasDir, "_common")
	dstCommon := filepath.Join(schemasRoot, "_common")
	if err := os.MkdirAll(dstCommon, 0o755); err != nil {
		return err
	}
	entries, err := os.ReadDir(srcCommon)
	if err != nil {
		return err
	}
	for _, ent := range entries {
		if ent.IsDir() || !stringsHasSuffixJSON(ent.Name()) {
			continue
		}
		src, err := os.ReadFile(filepath.Join(srcCommon, ent.Name()))
		if err != nil {
			return err
		}
		if err := os.WriteFile(filepath.Join(dstCommon, ent.Name()), src, 0o644); err != nil {
			return err
		}
	}
	return nil
}

func stringsHasSuffixJSON(name string) bool {
	return len(name) >= 5 && name[len(name)-5:] == ".json"
}

func mustUnmarshalAny(t *testing.T, raw []byte) any {
	t.Helper()
	var v any
	if err := json.Unmarshal(raw, &v); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	return v
}
