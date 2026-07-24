package apischema

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestClassifyChangeIdenticalSchemas(t *testing.T) {
	t.Parallel()

	schema := []byte(`{
		"type": "object",
		"properties": {
			"name": { "type": "string" }
		},
		"required": ["name"]
	}`)
	got := ClassifyChange(schema, schema)
	if len(got) != 0 {
		t.Fatalf("identical schemas: got %d changes %#v, want none", len(got), got)
	}
}

func TestClassifyChangeAddOptionalFieldCompatible(t *testing.T) {
	t.Parallel()

	oldSchema := []byte(`{
		"type": "object",
		"properties": {
			"name": { "type": "string" }
		},
		"required": ["name"]
	}`)
	newSchema := []byte(`{
		"type": "object",
		"properties": {
			"name": { "type": "string" },
			"description": { "type": "string" }
		},
		"required": ["name"]
	}`)

	got := ClassifyChange(oldSchema, newSchema)
	if !hasChange(got, ChangeCompatible, KindAddOptionalField, "/properties/description") {
		t.Fatalf("expected Compatible add_optional_field at /properties/description, got %#v", got)
	}
	if hasClass(got, ChangeBreaking) || hasClass(got, ChangeReviewRequired) {
		t.Fatalf("add optional field must not be Breaking or ReviewRequired, got %#v", got)
	}
}

func TestClassifyChangeRemoveFieldBreaking(t *testing.T) {
	t.Parallel()

	oldSchema := []byte(`{
		"type": "object",
		"properties": {
			"name": { "type": "string" },
			"description": { "type": "string" }
		}
	}`)
	newSchema := []byte(`{
		"type": "object",
		"properties": {
			"name": { "type": "string" }
		}
	}`)

	got := ClassifyChange(oldSchema, newSchema)
	if !hasChange(got, ChangeBreaking, KindRemoveField, "/properties/description") {
		t.Fatalf("expected Breaking remove_field at /properties/description, got %#v", got)
	}
}

func TestClassifyChangeAddEnumValueReviewRequired(t *testing.T) {
	t.Parallel()

	oldSchema := []byte(`{
		"type": "object",
		"properties": {
			"phase": {
				"type": "string",
				"enum": ["Pending", "Ready"]
			}
		}
	}`)
	newSchema := []byte(`{
		"type": "object",
		"properties": {
			"phase": {
				"type": "string",
				"enum": ["Pending", "Ready", "Failed"]
			}
		}
	}`)

	got := ClassifyChange(oldSchema, newSchema)
	if !hasChange(got, ChangeReviewRequired, KindAddEnumValue, "/properties/phase/enum") {
		t.Fatalf("expected ReviewRequired add_enum_value at /properties/phase/enum, got %#v", got)
	}
	if hasClass(got, ChangeBreaking) {
		t.Fatalf("pure enum addition must not be Breaking, got %#v", got)
	}
}

func TestClassifyChangeAddRequiredFieldBreaking(t *testing.T) {
	t.Parallel()

	oldSchema := []byte(`{
		"type": "object",
		"properties": {
			"name": { "type": "string" }
		}
	}`)
	newSchema := []byte(`{
		"type": "object",
		"properties": {
			"name": { "type": "string" },
			"owner": { "type": "string" }
		},
		"required": ["owner"]
	}`)

	got := ClassifyChange(oldSchema, newSchema)
	if !hasChange(got, ChangeBreaking, KindAddRequiredField, "/properties/owner") {
		t.Fatalf("expected Breaking add_required_field at /properties/owner, got %#v", got)
	}
}

func TestClassifyChangeNarrowEnumBreaking(t *testing.T) {
	t.Parallel()

	oldSchema := []byte(`{
		"type": "string",
		"enum": ["Pending", "Ready", "Failed"]
	}`)
	newSchema := []byte(`{
		"type": "string",
		"enum": ["Pending", "Ready"]
	}`)

	got := ClassifyChange(oldSchema, newSchema)
	if !hasChange(got, ChangeBreaking, KindNarrowEnum, "/enum") {
		t.Fatalf("expected Breaking narrow_enum at /enum, got %#v", got)
	}
}

func TestClassifyChangeNarrowValidationRangeBreaking(t *testing.T) {
	t.Parallel()

	oldSchema := []byte(`{
		"type": "string",
		"minLength": 1,
		"maxLength": 63
	}`)
	newSchema := []byte(`{
		"type": "string",
		"minLength": 3,
		"maxLength": 32
	}`)

	got := ClassifyChange(oldSchema, newSchema)
	if !hasChange(got, ChangeBreaking, KindNarrowValidationRange, "/minLength") {
		t.Fatalf("expected Breaking narrow minLength, got %#v", got)
	}
	if !hasChange(got, ChangeBreaking, KindNarrowValidationRange, "/maxLength") {
		t.Fatalf("expected Breaking narrow maxLength, got %#v", got)
	}
}

func TestClassifyChangeChangeRefTargetReviewRequired(t *testing.T) {
	t.Parallel()

	oldSchema := []byte(`{
		"$ref": "../_common/typed-ref.json"
	}`)
	newSchema := []byte(`{
		"$ref": "../_common/scope-ref.json"
	}`)

	got := ClassifyChange(oldSchema, newSchema)
	if !hasChange(got, ChangeReviewRequired, KindChangeReferenceTarget, "/$ref") {
		t.Fatalf("expected ReviewRequired change_reference_target at /$ref, got %#v", got)
	}
}

func TestClassifyChangeChangeAllowedScopesReviewRequired(t *testing.T) {
	t.Parallel()

	oldSchema := []byte(`{
		"type": "object",
		"x-sovrunn-allowed-scopes": ["Tenant"]
	}`)
	newSchema := []byte(`{
		"type": "object",
		"x-sovrunn-allowed-scopes": ["Tenant", "Project"]
	}`)

	got := ClassifyChange(oldSchema, newSchema)
	if !hasChange(got, ChangeReviewRequired, KindChangeAllowedScopes, "/x-sovrunn-allowed-scopes") {
		t.Fatalf("expected ReviewRequired change_allowed_scopes, got %#v", got)
	}
}

func TestClassifyChangeExposeInternalPubliclyReviewRequired(t *testing.T) {
	t.Parallel()

	oldSchema := []byte(`{
		"type": "string",
		"x-sovrunn-field-policy": {
			"classification": "Internal",
			"authorizedWriter": "system",
			"authorizedReaders": ["internal"],
			"mutability": "system-only",
			"retention": "standard",
			"redaction": "omit",
			"residency": "any",
			"auditRequired": true
		}
	}`)
	newSchema := []byte(`{
		"type": "string",
		"x-sovrunn-field-policy": {
			"classification": "Public",
			"authorizedWriter": "system",
			"authorizedReaders": ["customer"],
			"mutability": "system-only",
			"retention": "standard",
			"redaction": "none",
			"residency": "any",
			"auditRequired": true
		}
	}`)

	got := ClassifyChange(oldSchema, newSchema)
	if !hasChange(got, ChangeReviewRequired, KindExposeInternalPublicly, "/x-sovrunn-field-policy/classification") {
		t.Fatalf("expected ReviewRequired expose_internal_publicly, got %#v", got)
	}
}

func TestClassifyChangeMutabilityChangeBreaking(t *testing.T) {
	t.Parallel()

	oldSchema := []byte(`{
		"type": "string",
		"x-sovrunn-field-policy": {
			"classification": "Public",
			"authorizedWriter": "creator",
			"authorizedReaders": ["customer"],
			"mutability": "immutable",
			"retention": "standard",
			"redaction": "none",
			"residency": "any",
			"auditRequired": false
		}
	}`)
	newSchema := []byte(`{
		"type": "string",
		"x-sovrunn-field-policy": {
			"classification": "Public",
			"authorizedWriter": "creator",
			"authorizedReaders": ["customer"],
			"mutability": "mutable",
			"retention": "standard",
			"redaction": "none",
			"residency": "any",
			"auditRequired": false
		}
	}`)

	got := ClassifyChange(oldSchema, newSchema)
	if !hasChange(got, ChangeBreaking, KindChangeOwnerOrMutability, "/x-sovrunn-field-policy/mutability") {
		t.Fatalf("expected Breaking change_owner_or_mutability, got %#v", got)
	}
}

func TestClassifyChangeAddRegisteredExtensionCompatible(t *testing.T) {
	t.Parallel()

	oldSchema := []byte(`{
		"type": "object",
		"x-sovrunn-profile": "ManagedResource",
		"x-sovrunn-boundary": "customer-facing",
		"x-sovrunn-allowed-scopes": ["Tenant"],
		"x-sovrunn-stability": "alpha"
	}`)
	newSchema := []byte(`{
		"type": "object",
		"x-sovrunn-profile": "ManagedResource",
		"x-sovrunn-boundary": "customer-facing",
		"x-sovrunn-allowed-scopes": ["Tenant"],
		"x-sovrunn-stability": "alpha",
		"x-sovrunn-field-policy": {
			"classification": "Public",
			"authorizedWriter": "system",
			"authorizedReaders": ["customer"],
			"mutability": "system-only",
			"retention": "none",
			"redaction": "none",
			"residency": "any",
			"auditRequired": false
		}
	}`)

	got := ClassifyChange(oldSchema, newSchema)
	if !hasChange(got, ChangeCompatible, KindAddRegisteredExtension, "/x-sovrunn-field-policy") {
		t.Fatalf("expected Compatible add_registered_extension for field-policy, got %#v", got)
	}
}

func TestClassifyChangeMalformedSchemaReviewRequired(t *testing.T) {
	t.Parallel()

	got := ClassifyChange([]byte(`not-json`), []byte(`{"type":"object"}`))
	if !hasChange(got, ChangeReviewRequired, KindMalformedSchema, "/") {
		t.Fatalf("expected ReviewRequired malformed_schema, got %#v", got)
	}
}

func hasChange(changes []Change, class ChangeClass, kind ChangeKind, path string) bool {
	for _, c := range changes {
		if c.Class == class && c.Kind == kind && c.Path == path {
			return true
		}
	}
	return false
}

func hasClass(changes []Change, class ChangeClass) bool {
	for _, c := range changes {
		if c.Class == class {
			return true
		}
	}
	return false
}

func TestClassifyChangePromoteOptionalToRequiredBreaking(t *testing.T) {
	t.Parallel()

	oldSchema := []byte(`{
		"type": "object",
		"properties": {
			"name": { "type": "string" },
			"label": { "type": "string" }
		},
		"required": ["name"]
	}`)
	newSchema := []byte(`{
		"type": "object",
		"properties": {
			"name": { "type": "string" },
			"label": { "type": "string" }
		},
		"required": ["name", "label"]
	}`)

	got := ClassifyChange(oldSchema, newSchema)
	if !hasChange(got, ChangeBreaking, KindPromoteOptionalToRequired, "/properties/label") {
		t.Fatalf("expected Breaking promote_optional_to_required, got %#v", got)
	}
	// Sanity: message mentions promotion.
	found := false
	for _, c := range got {
		if c.Kind == KindPromoteOptionalToRequired && strings.Contains(c.Message, "label") {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected promote message mentioning label, got %#v", got)
	}
}

func TestVerifyBaselineIntegrityTamperedBaselineFails(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	schemaPath := filepath.Join(dir, "project.json")
	original := []byte(`{"type":"object","title":"project"}`)
	if err := os.WriteFile(schemaPath, original, 0o644); err != nil {
		t.Fatalf("write schema: %v", err)
	}
	digest := sha256Hex(original)
	writeManifest(t, dir, map[string]string{"project.json": digest})

	// Tamper baseline content without updating the manifest.
	if err := os.WriteFile(schemaPath, []byte(`{"type":"object","title":"tampered"}`), 0o644); err != nil {
		t.Fatalf("tamper schema: %v", err)
	}

	err := VerifyBaselineIntegrity(filepath.Join(dir, BaselineManifestFileName), dir)
	if err == nil {
		t.Fatal("expected integrity failure for tampered baseline, got nil")
	}
	if !strings.Contains(err.Error(), "digest mismatch") {
		t.Fatalf("expected digest mismatch error, got %v", err)
	}
}

func TestVerifyBaselineIntegrityMatchingPasses(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	schema := []byte(`{"type":"object","title":"project"}`)
	if err := os.WriteFile(filepath.Join(dir, "project.json"), schema, 0o644); err != nil {
		t.Fatalf("write schema: %v", err)
	}
	writeManifest(t, dir, map[string]string{"project.json": sha256Hex(schema)})

	if err := VerifyBaselineIntegrity(filepath.Join(dir, BaselineManifestFileName), dir); err != nil {
		t.Fatalf("expected integrity pass, got %v", err)
	}
}

func TestVerifyBaselineApprovalMissingEvidenceFails(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	oldContent := []byte(`{"type":"object","title":"old"}`)
	newContent := []byte(`{"type":"object","title":"new"}`)
	oldDigest := sha256Hex(oldContent)
	newDigest := sha256Hex(newContent)

	if err := os.WriteFile(filepath.Join(dir, "project.json"), newContent, 0o644); err != nil {
		t.Fatalf("write schema: %v", err)
	}
	// Co-edit: manifest updated to match the new baseline bytes.
	writeManifest(t, dir, map[string]string{"project.json": newDigest})
	// Recorded digests still reflect the prior approved baseline; no evidence.
	writeApprovals(t, dir, baselineApprovalsFile{
		RecordedDigests: map[string]string{"project.json": oldDigest},
		Approvals:       nil,
	})

	err := VerifyBaselineApproval(
		filepath.Join(dir, BaselineApprovalsFileName),
		filepath.Join(dir, BaselineManifestFileName),
		dir,
	)
	if err == nil {
		t.Fatal("expected approval failure when evidence is missing, got nil")
	}
	if !strings.Contains(err.Error(), "without recorded approval evidence") {
		t.Fatalf("expected missing-evidence error, got %v", err)
	}
}

func TestVerifyBaselineApprovalValidEvidencePasses(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	oldContent := []byte(`{"type":"object","title":"old"}`)
	newContent := []byte(`{"type":"object","title":"new"}`)
	oldDigest := sha256Hex(oldContent)
	newDigest := sha256Hex(newContent)

	if err := os.WriteFile(filepath.Join(dir, "project.json"), newContent, 0o644); err != nil {
		t.Fatalf("write schema: %v", err)
	}
	writeManifest(t, dir, map[string]string{"project.json": newDigest})
	writeApprovals(t, dir, baselineApprovalsFile{
		RecordedDigests: map[string]string{"project.json": oldDigest},
		Approvals: []baselineApproval{{
			Path:      "project.json",
			OldDigest: oldDigest,
			NewDigest: newDigest,
			ADH:       "ADH-2026-012",
			Reviewer:  "Sanjeev Kumar",
			Date:      "2026-07-23",
		}},
	})

	if err := VerifyBaselineApproval(
		filepath.Join(dir, BaselineApprovalsFileName),
		filepath.Join(dir, BaselineManifestFileName),
		dir,
	); err != nil {
		t.Fatalf("expected approval pass with matching evidence, got %v", err)
	}
}

func TestVerifyBaselineApprovalCoEditWithoutEvidenceFails(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	oldContent := []byte(`{"type":"object","properties":{"a":{"type":"string"}}}`)
	newContent := []byte(`{"type":"object","properties":{"a":{"type":"string"},"b":{"type":"string"}}}`)
	oldDigest := sha256Hex(oldContent)
	newDigest := sha256Hex(newContent)

	if err := os.WriteFile(filepath.Join(dir, "resource-pool.json"), newContent, 0o644); err != nil {
		t.Fatalf("write schema: %v", err)
	}
	writeManifest(t, dir, map[string]string{"resource-pool.json": newDigest})
	writeApprovals(t, dir, baselineApprovalsFile{
		RecordedDigests: map[string]string{"resource-pool.json": oldDigest},
		Approvals:       []baselineApproval{}, // co-edit only; no evidence
	})

	// Integrity alone would pass (manifest matches files).
	if err := VerifyBaselineIntegrity(filepath.Join(dir, BaselineManifestFileName), dir); err != nil {
		t.Fatalf("co-edit setup must pass integrity: %v", err)
	}
	err := VerifyBaselineApproval(
		filepath.Join(dir, BaselineApprovalsFileName),
		filepath.Join(dir, BaselineManifestFileName),
		dir,
	)
	if err == nil {
		t.Fatal("co-editing baseline+manifest without approval evidence must fail")
	}
}

func TestVerifyBaselineApprovalInitialEmptyPasses(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	schema := []byte(`{"type":"object"}`)
	if err := os.WriteFile(filepath.Join(dir, "project.json"), schema, 0o644); err != nil {
		t.Fatalf("write schema: %v", err)
	}
	writeManifest(t, dir, map[string]string{"project.json": sha256Hex(schema)})
	writeApprovals(t, dir, baselineApprovalsFile{}) // empty: initial baseline

	if err := VerifyBaselineApproval(
		filepath.Join(dir, BaselineApprovalsFileName),
		filepath.Join(dir, BaselineManifestFileName),
		dir,
	); err != nil {
		t.Fatalf("initial empty approvals must pass, got %v", err)
	}
}

func TestVerifyBaselineApprovalIncompleteEvidenceFails(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	oldContent := []byte(`{"type":"string"}`)
	newContent := []byte(`{"type":"integer"}`)
	oldDigest := sha256Hex(oldContent)
	newDigest := sha256Hex(newContent)

	if err := os.WriteFile(filepath.Join(dir, "plugin-definition.json"), newContent, 0o644); err != nil {
		t.Fatalf("write schema: %v", err)
	}
	writeManifest(t, dir, map[string]string{"plugin-definition.json": newDigest})
	writeApprovals(t, dir, baselineApprovalsFile{
		RecordedDigests: map[string]string{"plugin-definition.json": oldDigest},
		Approvals: []baselineApproval{{
			Path:      "plugin-definition.json",
			OldDigest: oldDigest,
			NewDigest: newDigest,
			ADH:       "ADH-2026-012",
			// missing reviewer and date
		}},
	})

	err := VerifyBaselineApproval(
		filepath.Join(dir, BaselineApprovalsFileName),
		filepath.Join(dir, BaselineManifestFileName),
		dir,
	)
	if err == nil {
		t.Fatal("expected failure for incomplete approval evidence")
	}
	if !strings.Contains(err.Error(), "reviewer") {
		t.Fatalf("expected reviewer error, got %v", err)
	}
}

func TestVerifyBaselineApprovalTokenAccepted(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	oldContent := []byte(`{"enum":["a"]}`)
	newContent := []byte(`{"enum":["a","b"]}`)
	oldDigest := sha256Hex(oldContent)
	newDigest := sha256Hex(newContent)

	if err := os.WriteFile(filepath.Join(dir, "operation.json"), newContent, 0o644); err != nil {
		t.Fatalf("write schema: %v", err)
	}
	writeManifest(t, dir, map[string]string{"operation.json": newDigest})
	writeApprovals(t, dir, baselineApprovalsFile{
		RecordedDigests: map[string]string{"operation.json": oldDigest},
		Approvals: []baselineApproval{{
			Path:          "operation.json",
			OldDigest:     oldDigest,
			NewDigest:     newDigest,
			ApprovalToken: "APPROVED_BASELINE_CHANGE",
			Reviewer:      "Reviewer",
			Date:          "2026-07-23",
		}},
	})

	if err := VerifyBaselineApproval(
		filepath.Join(dir, BaselineApprovalsFileName),
		filepath.Join(dir, BaselineManifestFileName),
		dir,
	); err != nil {
		t.Fatalf("approval token must be accepted in place of ADH, got %v", err)
	}
}

func writeManifest(t *testing.T, dir string, files map[string]string) {
	t.Helper()
	raw, err := json.Marshal(baselineManifest{Files: files})
	if err != nil {
		t.Fatalf("marshal manifest: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, BaselineManifestFileName), raw, 0o644); err != nil {
		t.Fatalf("write manifest: %v", err)
	}
}

func writeApprovals(t *testing.T, dir string, approvals baselineApprovalsFile) {
	t.Helper()
	raw, err := json.Marshal(approvals)
	if err != nil {
		t.Fatalf("marshal approvals: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, BaselineApprovalsFileName), raw, 0o644); err != nil {
		t.Fatalf("write approvals: %v", err)
	}
}

func sha256Hex(b []byte) string {
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:])
}
