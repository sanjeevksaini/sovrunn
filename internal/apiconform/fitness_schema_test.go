package apiconform

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/apischema"
)

func TestFitnessCheckExternalSchemaAnnotations(t *testing.T) {
	t.Parallel()

	root := filepath.Join(moduleRoot(t), CanonicalSchemasDir)
	findings := CheckExternalSchemaAnnotations(root)
	if len(findings) != 0 {
		t.Fatalf("check 1 failed: %#v", findings)
	}
}

func TestFitnessCheckExternalSchemaAnnotationsMissingFails(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	for _, name := range externalCanonicalSchemaFiles {
		src, err := os.ReadFile(filepath.Join(moduleRoot(t), CanonicalSchemasDir, name))
		if err != nil {
			t.Fatalf("read %s: %v", name, err)
		}
		if name == "project.json" {
			var root map[string]any
			if err := json.Unmarshal(src, &root); err != nil {
				t.Fatalf("unmarshal: %v", err)
			}
			delete(root, apischema.ExtProfile)
			src, err = json.Marshal(root)
			if err != nil {
				t.Fatalf("marshal: %v", err)
			}
		}
		if err := os.WriteFile(filepath.Join(dir, name), src, 0o644); err != nil {
			t.Fatalf("write %s: %v", name, err)
		}
	}
	findings := CheckExternalSchemaAnnotations(dir)
	if !hasFitnessFinding(findings, FitnessCheckExternalSchemaAnnotations, CanonicalSchemasDir+"/project.json", "/"+apischema.ExtProfile, CodeFitnessAnnotationMissing) {
		t.Fatalf("expected missing profile finding, got %#v", findings)
	}
}

func TestFitnessCheckFieldPolicyCoverage(t *testing.T) {
	t.Parallel()

	root := filepath.Join(moduleRoot(t), CanonicalSchemasDir)
	findings := CheckFieldPolicyCoverage(root)
	if len(findings) != 0 {
		t.Fatalf("check 1a failed: %#v", findings)
	}
}

func TestFitnessCheckFieldPolicyCoverageMissingPolicyFails(t *testing.T) {
	t.Parallel()

	var root map[string]any
	if err := json.Unmarshal(loadCanonicalSchema(t, "project.json"), &root); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	props := root["properties"].(map[string]any)
	kind := props["kind"].(map[string]any)
	delete(kind, apischema.ExtFieldPolicy)
	raw, err := json.Marshal(root)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	issues := CheckFieldPolicyCompleteness(raw)
	if !hasSchemaIssue(issues, CodeFitnessFieldPolicyMissing, "/properties/kind") {
		t.Fatalf("expected missing policy field, got %#v", issues)
	}
}

func TestFitnessCheckFieldPolicyCoverageUnknownPolicyFieldFails(t *testing.T) {
	t.Parallel()

	var root map[string]any
	if err := json.Unmarshal(loadCanonicalSchema(t, "project.json"), &root); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	props := root["properties"].(map[string]any)
	kind := props["kind"].(map[string]any)
	fp := kind[apischema.ExtFieldPolicy].(map[string]any)
	fp["extra"] = "nope"
	raw, err := json.Marshal(root)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	issues := CheckFieldPolicyCompleteness(raw)
	wantPath := "/properties/kind/" + apischema.ExtFieldPolicy + "/extra"
	if !hasSchemaIssue(issues, apischema.CodeFieldPolicyUnknownField, wantPath) {
		t.Fatalf("expected unknown policy field at %s, got %#v", wantPath, issues)
	}
}

func TestFitnessCheckFieldPolicyCoverageUnknownExtensionFails(t *testing.T) {
	t.Parallel()

	var root map[string]any
	if err := json.Unmarshal(loadCanonicalSchema(t, "project.json"), &root); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	root["x-sovrunn-foo"] = "bar"
	raw, err := json.Marshal(root)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	_, issues := apischema.ReadAnnotations(raw)
	if !hasSchemaIssue(issues, apischema.CodeUnknownExtension, "/x-sovrunn-foo") {
		t.Fatalf("expected unknown x-sovrunn-* rejection, got %#v", issues)
	}
	extIssues := scanUnknownSovrunnExtensions(raw)
	if !hasSchemaIssue(extIssues, apischema.CodeUnknownExtension, "/x-sovrunn-foo") {
		t.Fatalf("expected scanUnknownSovrunnExtensions rejection, got %#v", extIssues)
	}
}

func TestFitnessCheckMutableFieldOwnership(t *testing.T) {
	t.Parallel()

	root := filepath.Join(moduleRoot(t), CanonicalSchemasDir)
	findings := CheckMutableFieldOwnership(root)
	if len(findings) != 0 {
		t.Fatalf("check 3 failed: %#v", findings)
	}
}

func TestFitnessCheckMutableFieldOwnershipMissingWriterFails(t *testing.T) {
	t.Parallel()

	var root map[string]any
	if err := json.Unmarshal(loadCanonicalSchema(t, "project.json"), &root); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	props := root["properties"].(map[string]any)
	spec := props["spec"].(map[string]any)
	fp := spec[apischema.ExtFieldPolicy].(map[string]any)
	delete(fp, "authorizedWriter")
	raw, err := json.Marshal(root)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	findings := ownershipFindingsForSchema(CanonicalSchemasDir+"/project.json", raw, false)
	if !hasFitnessFinding(findings, FitnessCheckMutableFieldOwnership, CanonicalSchemasDir+"/project.json",
		"/properties/spec/"+apischema.ExtFieldPolicy+"/authorizedWriter", CodeFitnessOwnershipMissing) {
		t.Fatalf("expected missing owner finding, got %#v", findings)
	}
}

func TestFitnessCheckMutableFieldOwnershipConditionWriterFails(t *testing.T) {
	t.Parallel()

	var root map[string]any
	if err := json.Unmarshal(loadCanonicalSchema(t, "project.json"), &root); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	status := root["properties"].(map[string]any)["status"].(map[string]any)
	conds := status["properties"].(map[string]any)["conditions"].(map[string]any)
	fp := conds[apischema.ExtFieldPolicy].(map[string]any)
	fp["authorizedWriter"] = apischema.WriterCreator
	raw, err := json.Marshal(root)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	findings := ownershipFindingsForSchema(CanonicalSchemasDir+"/project.json", raw, false)
	if !hasFitnessFinding(findings, FitnessCheckMutableFieldOwnership, CanonicalSchemasDir+"/project.json",
		"/properties/status/properties/conditions/"+apischema.ExtFieldPolicy+"/authorizedWriter",
		CodeFitnessConditionOwnerInvalid) {
		t.Fatalf("expected condition owner finding, got %#v", findings)
	}
}

func TestFitnessCheckUnknownAndDuplicateFieldsFail(t *testing.T) {
	t.Parallel()

	findings := CheckUnknownAndDuplicateFieldsFail()
	if len(findings) != 0 {
		t.Fatalf("check 4 failed: %#v", findings)
	}
}

func TestFitnessCheckPublishedDefinitionsImmutable(t *testing.T) {
	t.Parallel()

	root := filepath.Join(moduleRoot(t), CanonicalSchemasDir)
	findings := CheckPublishedDefinitionsImmutable(root)
	if len(findings) != 0 {
		t.Fatalf("check 9 failed: %#v", findings)
	}
}

func TestFitnessCheckPublishedDefinitionsImmutableMutableContractFails(t *testing.T) {
	t.Parallel()

	var root map[string]any
	if err := json.Unmarshal(loadCanonicalSchema(t, "plugin-definition.json"), &root); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	spec := root["properties"].(map[string]any)["spec"].(map[string]any)
	version := spec["properties"].(map[string]any)["version"].(map[string]any)
	fp := version[apischema.ExtFieldPolicy].(map[string]any)
	fp["mutability"] = apischema.MutabilityMutable
	raw, err := json.Marshal(root)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	findings := versionedDefinitionImmutabilityFindings(CanonicalSchemasDir+"/plugin-definition.json", raw)
	if !hasFitnessFinding(findings, FitnessCheckPublishedDefinitionsImmutable,
		CanonicalSchemasDir+"/plugin-definition.json",
		"/properties/spec/properties/version/"+apischema.ExtFieldPolicy+"/mutability",
		CodeFitnessPublishedFieldMutable) {
		t.Fatalf("expected published-field mutability finding, got %#v", findings)
	}
}

func TestFitnessCheckImmutableRecordMutablePayloadFails(t *testing.T) {
	t.Parallel()

	var root map[string]any
	if err := json.Unmarshal(loadCanonicalSchema(t, "audit-event.json"), &root); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	record := root["properties"].(map[string]any)["record"].(map[string]any)
	fp := record[apischema.ExtFieldPolicy].(map[string]any)
	fp["mutability"] = apischema.MutabilityMutable
	raw, err := json.Marshal(root)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	findings := immutableRecordMutabilityFindings(CanonicalSchemasDir+"/audit-event.json", raw)
	if !hasFitnessFinding(findings, FitnessCheckPublishedDefinitionsImmutable,
		CanonicalSchemasDir+"/audit-event.json",
		"/properties/record/"+apischema.ExtFieldPolicy+"/mutability",
		CodeFitnessImmutableRecordMutable) {
		t.Fatalf("expected ImmutableRecord mutability finding, got %#v", findings)
	}
}

func TestFitnessSchemaChecksInventory(t *testing.T) {
	t.Parallel()

	if got := ExternalCanonicalSchemaFiles(); len(got) != 8 {
		t.Fatalf("ExternalCanonicalSchemaFiles len=%d want 8", len(got))
	}
	if got := CommonSubSchemaFiles(); len(got) != 9 {
		t.Fatalf("CommonSubSchemaFiles len=%d want 9", len(got))
	}
}

func hasFitnessFinding(findings []FitnessFinding, check, schema, path, code string) bool {
	for _, f := range findings {
		if f.Check == check && f.Schema == schema && f.Path == path && f.Code == code {
			return true
		}
	}
	return false
}
