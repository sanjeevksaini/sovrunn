package apiconform

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"gopkg.in/yaml.v3"
)

func TestAllFitnessChecksRegisteredAndExecuted(t *testing.T) {
	t.Parallel()

	root := moduleRoot(t)
	regs := RegisteredFitnessChecks()
	if len(regs) == 0 {
		t.Fatal("RegisteredFitnessChecks returned empty registry")
	}

	reqIDs := RequiredFitnessCheckIDs()
	if len(reqIDs) != 15 {
		t.Fatalf("RequiredFitnessCheckIDs len=%d want 15", len(reqIDs))
	}

	byID := make(map[string]FitnessCheckRegistration, len(regs))
	for _, reg := range regs {
		if reg.ID == "" {
			t.Fatal("registered fitness check has empty ID")
		}
		if reg.Requirement == "" {
			t.Fatalf("check %q missing requirement text", reg.ID)
		}
		if reg.Run == nil {
			t.Fatalf("check %q has nil Run", reg.ID)
		}
		if _, dup := byID[reg.ID]; dup {
			t.Fatalf("duplicate registered fitness check ID %q", reg.ID)
		}
		byID[reg.ID] = reg
	}

	reqMap := FitnessCheckRequirementMap()
	for _, id := range reqIDs {
		reg, ok := byID[id]
		if !ok {
			t.Fatalf("F12-VERIFY-001 check %q is not registered", id)
		}
		wantReq, ok := reqMap[id]
		if !ok || wantReq == "" {
			t.Fatalf("check %q missing from FitnessCheckRequirementMap", id)
		}
		if reg.Requirement != wantReq {
			t.Fatalf("check %q requirement mismatch:\n got %q\nwant %q", id, reg.Requirement, wantReq)
		}
	}

	// Check 1a (field-policy) is registered in addition to checks 1–15.
	if _, ok := byID[FitnessCheckFieldPolicyCoverage]; !ok {
		t.Fatalf("check %q (field-policy coverage) must be registered with aggregation", FitnessCheckFieldPolicyCoverage)
	}

	findings := RunAllFitnessChecks(root)
	if len(findings) != 0 {
		t.Fatalf("aggregated fitness checks 1–15 (+1a) failed: %#v", findings)
	}

	// Prove every registered check was executed by re-running each Run and
	// confirming the aggregate equals the concatenation of individual results.
	var expected []FitnessFinding
	for _, reg := range regs {
		expected = append(expected, reg.Run(root)...)
	}
	aggregated := RunAllFitnessChecks(root)
	if len(aggregated) != len(expected) {
		t.Fatalf("RunAllFitnessChecks len=%d want %d (per-check concatenation)", len(aggregated), len(expected))
	}
	for i := range expected {
		if aggregated[i] != expected[i] {
			t.Fatalf("aggregated finding[%d]=%#v want %#v", i, aggregated[i], expected[i])
		}
	}
}

func TestFitnessCheckBoundaryLedgerPasses(t *testing.T) {
	t.Parallel()

	findings := CheckBoundaryLedger(moduleRoot(t))
	if len(findings) != 0 {
		t.Fatalf("boundary ledger fitness check failed: %#v", findings)
	}
}

func TestFitnessCheckBoundaryLedgerMissingCategoryFails(t *testing.T) {
	t.Parallel()

	root := moduleRoot(t)
	src, err := os.ReadFile(filepath.Join(root, BoundaryLedgerPath))
	if err != nil {
		t.Fatalf("read ledger: %v", err)
	}
	doc, err := ParseBoundaryLedgerYAML(src)
	if err != nil {
		t.Fatalf("parse ledger: %v", err)
	}
	if len(doc.Boundaries) == 0 {
		t.Fatal("ledger has no boundaries")
	}
	doc.Boundaries[0].Purpose = ""

	dir := t.TempDir()
	ledgerPath := filepath.Join(dir, "boundary-ledger.yaml")
	if err := os.WriteFile(ledgerPath, mustRenderMinimalLedgerYAML(t, doc), 0o644); err != nil {
		t.Fatalf("write ledger: %v", err)
	}

	// Point schemas at the real canonical schemas so schema→ledger coverage
	// still resolves; only the ledger under test is mutated.
	findings := checkBoundaryLedgerAt(ledgerPath, filepath.Join(root, CanonicalSchemasDir))
	if !hasFitnessFinding(findings, FitnessCheckBoundaryLedger, BoundaryLedgerPath,
		"/boundaries/0/purpose", CodeFitnessLedgerCategoryMissing) {
		t.Fatalf("expected missing purpose category finding, got %#v", findings)
	}
}

func TestFitnessCheckBoundaryLedgerSchemaBoundaryWithoutEntryFails(t *testing.T) {
	t.Parallel()

	root := moduleRoot(t)
	src, err := os.ReadFile(filepath.Join(root, BoundaryLedgerPath))
	if err != nil {
		t.Fatalf("read ledger: %v", err)
	}
	doc, err := ParseBoundaryLedgerYAML(src)
	if err != nil {
		t.Fatalf("parse ledger: %v", err)
	}

	// Drop customer-facing (used by project.json) from the ledger.
	filtered := make([]BoundaryLedgerEntry, 0, len(doc.Boundaries))
	for _, entry := range doc.Boundaries {
		if entry.ID == string(apimeta.BoundaryCustomerFacing) {
			continue
		}
		filtered = append(filtered, entry)
	}
	doc.Boundaries = filtered

	dir := t.TempDir()
	ledgerPath := filepath.Join(dir, "boundary-ledger.yaml")
	if err := os.WriteFile(ledgerPath, mustRenderMinimalLedgerYAML(t, doc), 0o644); err != nil {
		t.Fatalf("write ledger: %v", err)
	}

	findings := checkBoundaryLedgerAt(ledgerPath, filepath.Join(root, CanonicalSchemasDir))
	if !hasFitnessFinding(findings, FitnessCheckBoundaryLedger, BoundaryLedgerPath,
		"/boundaries", CodeFitnessLedgerSchemaBoundaryMissing) {
		t.Fatalf("expected schema-boundary-without-ledger finding, got %#v", findings)
	}
	found := false
	for _, f := range findings {
		if f.Code == CodeFitnessLedgerSchemaBoundaryMissing &&
			strings.Contains(f.Message, string(apimeta.BoundaryCustomerFacing)) {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected message mentioning %q, got %#v", apimeta.BoundaryCustomerFacing, findings)
	}
}

func TestLedgerEntryCategoryGaps(t *testing.T) {
	t.Parallel()

	complete := BoundaryLedgerEntry{
		ID:                  string(apimeta.BoundaryCustomerFacing),
		Purpose:             "p",
		Owner:               "o",
		Producers:           []string{"prod"},
		Consumers:           []string{"cons"},
		AllowedData:         []string{"a"},
		ProhibitedData:      []string{"b"},
		Authorization:       "authz",
		Audit:               "audit",
		Observability:       "obs",
		FailureBehavior:     "fail",
		Versioning:          "ver",
		ReplacementPath:     "rep",
		MigrationPath:       "mig",
		ReassessmentTrigger: "trig",
	}
	if gaps := LedgerEntryCategoryGaps(complete); len(gaps) != 0 {
		t.Fatalf("complete entry gaps=%v want none", gaps)
	}

	incomplete := complete
	incomplete.Audit = ""
	incomplete.Producers = nil
	gaps := LedgerEntryCategoryGaps(incomplete)
	want := map[string]bool{"audit": true, "producers": true}
	if len(gaps) != 2 {
		t.Fatalf("gaps=%v want audit and producers", gaps)
	}
	for _, g := range gaps {
		if !want[g] {
			t.Fatalf("unexpected gap %q in %v", g, gaps)
		}
	}
}

func TestFitnessAggregationInventory(t *testing.T) {
	t.Parallel()

	if FitnessCheckBoundaryLedger != "ledger" {
		t.Fatalf("FitnessCheckBoundaryLedger=%q", FitnessCheckBoundaryLedger)
	}
	if len(RequiredFitnessCheckIDs()) != 15 {
		t.Fatalf("want 15 required check IDs")
	}
	regs := RegisteredFitnessChecks()
	// 15 numbered checks + 1a
	if len(regs) != 16 {
		t.Fatalf("RegisteredFitnessChecks len=%d want 16 (1–15 + 1a)", len(regs))
	}
}

func mustRenderMinimalLedgerYAML(t *testing.T, doc BoundaryLedger) []byte {
	t.Helper()
	raw, err := yaml.Marshal(doc)
	if err != nil {
		t.Fatalf("marshal ledger: %v", err)
	}
	return raw
}
