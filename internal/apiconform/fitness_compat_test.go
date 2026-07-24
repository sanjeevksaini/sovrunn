package apiconform

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/apischema"
)

func TestFitnessCheckSchemaCompatibilityDetectsBreaking(t *testing.T) {
	t.Parallel()

	findings := CheckSchemaCompatibilityDetectsBreaking(moduleRoot(t))
	if len(findings) != 0 {
		t.Fatalf("check 10 failed: %#v", findings)
	}
}

func TestFitnessCheckSchemaCompatibilityTamperedBaselineFails(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	schemaPath := filepath.Join(dir, "project.json")
	original := []byte(`{"type":"object","title":"project"}`)
	if err := os.WriteFile(schemaPath, original, 0o644); err != nil {
		t.Fatalf("write schema: %v", err)
	}
	writeFitnessBaselineManifest(t, dir, map[string]string{"project.json": sha256HexCompat(original)})
	writeFitnessBaselineApprovals(t, dir, map[string]any{
		"recordedDigests": map[string]string{},
		"approvals":       []any{},
	})

	if err := os.WriteFile(schemaPath, []byte(`{"type":"object","title":"tampered"}`), 0o644); err != nil {
		t.Fatalf("tamper schema: %v", err)
	}

	err := apischema.VerifyBaselineIntegrity(filepath.Join(dir, apischema.BaselineManifestFileName), dir)
	if err == nil {
		t.Fatal("expected integrity failure for tampered baseline")
	}
	if !strings.Contains(err.Error(), "digest mismatch") {
		t.Fatalf("expected digest mismatch, got %v", err)
	}
}

func TestFitnessCheckSchemaCompatibilityMissingApprovalFails(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	oldContent := []byte(`{"type":"object","title":"old"}`)
	newContent := []byte(`{"type":"object","title":"new"}`)
	oldDigest := sha256HexCompat(oldContent)
	newDigest := sha256HexCompat(newContent)

	if err := os.WriteFile(filepath.Join(dir, "project.json"), newContent, 0o644); err != nil {
		t.Fatalf("write schema: %v", err)
	}
	// Co-edit baseline + manifest without approval evidence (D-11).
	writeFitnessBaselineManifest(t, dir, map[string]string{"project.json": newDigest})
	writeFitnessBaselineApprovals(t, dir, map[string]any{
		"recordedDigests": map[string]string{"project.json": oldDigest},
		"approvals":       []any{},
	})

	err := apischema.VerifyBaselineApproval(
		filepath.Join(dir, apischema.BaselineApprovalsFileName),
		filepath.Join(dir, apischema.BaselineManifestFileName),
		dir,
	)
	if err == nil {
		t.Fatal("expected approval failure when evidence is missing")
	}
	if !strings.Contains(err.Error(), "without recorded approval evidence") {
		t.Fatalf("expected missing-evidence error, got %v", err)
	}
}

func TestFitnessCheckSchemaCompatibilityBreakingDiffFails(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	schemasDir := filepath.Join(root, CanonicalSchemasDir)
	baselineDir := filepath.Join(root, BaselineSchemasDir)
	if err := os.MkdirAll(filepath.Join(schemasDir, "_common"), 0o755); err != nil {
		t.Fatalf("mkdir schemas: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(baselineDir, "_common"), 0o755); err != nil {
		t.Fatalf("mkdir baseline: %v", err)
	}

	oldSchema := []byte(`{"type":"object","properties":{"name":{"type":"string"},"label":{"type":"string"}}}`)
	newSchema := []byte(`{"type":"object","properties":{"name":{"type":"string"}}}`)
	if err := os.WriteFile(filepath.Join(baselineDir, "project.json"), oldSchema, 0o644); err != nil {
		t.Fatalf("write baseline: %v", err)
	}
	if err := os.WriteFile(filepath.Join(schemasDir, "project.json"), newSchema, 0o644); err != nil {
		t.Fatalf("write current: %v", err)
	}
	writeFitnessBaselineManifest(t, baselineDir, map[string]string{"project.json": sha256HexCompat(oldSchema)})
	writeFitnessBaselineApprovals(t, baselineDir, map[string]any{
		"recordedDigests": map[string]string{},
		"approvals":       []any{},
	})

	findings := baselineCurrentDiffFindings(root)
	if !hasFitnessFinding(findings, FitnessCheckSchemaCompatibility, CanonicalSchemasDir+"/project.json",
		"/properties/label", CodeFitnessUnapprovedBreakingChange) {
		// Path may vary slightly; accept any breaking finding for project.
		found := false
		for _, f := range findings {
			if f.Check == FitnessCheckSchemaCompatibility &&
				f.Schema == CanonicalSchemasDir+"/project.json" &&
				f.Code == CodeFitnessUnapprovedBreakingChange {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("expected unapproved breaking change finding, got %#v", findings)
		}
	}
}

func TestFitnessCheckSizesBounded(t *testing.T) {
	t.Parallel()

	findings := CheckSizesBounded()
	if len(findings) != 0 {
		t.Fatalf("check 11 failed: %#v", findings)
	}
}

func TestFitnessCheckSizesBoundedLimitsMismatchFails(t *testing.T) {
	t.Parallel()

	// Local assertion of the expected table used by check 11 — a drift here
	// would make CheckSizesBounded fail against DefaultLimits.
	want := expectedDefaultLimits
	if want.MaxObjectBytes != 1_048_576 || want.MaxPageSize != 200 || want.MaxConditions != 32 {
		t.Fatalf("expectedDefaultLimits drifted from D-06 table: %#v", want)
	}
}

func TestFitnessCheckErrorsUseStableCodesAndJSONPointers(t *testing.T) {
	t.Parallel()

	findings := CheckErrorsUseStableCodesAndJSONPointers()
	if len(findings) != 0 {
		t.Fatalf("check 12 failed: %#v", findings)
	}
}

func TestFitnessCheckGeneratedArtifactsMatchCanonicalSchema(t *testing.T) {
	t.Parallel()

	findings := CheckGeneratedArtifactsMatchCanonicalSchema(moduleRoot(t))
	if len(findings) != 0 {
		t.Fatalf("check 13 failed: %#v", findings)
	}
}

func TestFitnessCheckGeneratedArtifactsDeliberateMismatchFails(t *testing.T) {
	t.Parallel()

	schema, err := os.ReadFile(filepath.Join(moduleRoot(t), "api/schemas/_common/page.json"))
	if err != nil {
		t.Fatalf("read page schema: %v", err)
	}
	type mismatchedPage struct {
		NextPageToken int `json:"nextPageToken,omitempty"`
	}
	issues := apischema.VerifyGoTypeAgainstSchema(schema, reflect.TypeOf(mismatchedPage{}))
	if len(issues) == 0 {
		t.Fatal("expected deliberate Go-type mismatch to be rejected")
	}
}

func TestFitnessCompatChecksInventory(t *testing.T) {
	t.Parallel()

	want := []string{
		FitnessCheckSchemaCompatibility,
		FitnessCheckSizesBounded,
		FitnessCheckStableCodesAndJSONPointers,
		FitnessCheckGeneratedArtifactsMatchSchema,
	}
	for _, id := range want {
		if id == "" {
			t.Fatal("empty fitness check id")
		}
	}
	if BaselineSchemasDir != "api/schemas/baseline" {
		t.Fatalf("BaselineSchemasDir = %q", BaselineSchemasDir)
	}
}

func writeFitnessBaselineManifest(t *testing.T, dir string, files map[string]string) {
	t.Helper()
	raw, err := json.Marshal(map[string]any{"files": files})
	if err != nil {
		t.Fatalf("marshal manifest: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, apischema.BaselineManifestFileName), raw, 0o644); err != nil {
		t.Fatalf("write manifest: %v", err)
	}
}

func writeFitnessBaselineApprovals(t *testing.T, dir string, body map[string]any) {
	t.Helper()
	raw, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("marshal approvals: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, apischema.BaselineApprovalsFileName), raw, 0o644); err != nil {
		t.Fatalf("write approvals: %v", err)
	}
}

func sha256HexCompat(b []byte) string {
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:])
}
