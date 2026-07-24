package apiconform

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRequiredPhase1Contracts_F12Compat001Inventory(t *testing.T) {
	t.Parallel()

	want := []string{
		"Organization",
		"OrganizationUnit",
		"Tenant",
		"Project",
		"Operation",
		"ServiceClass",
		"ServicePlan",
		"Plugin",
		"Capability",
		"ServiceInstance",
		"ServiceBinding",
		"health/readiness",
		"demo-flow",
	}
	got := RequiredPhase1Contracts()
	if len(got) != len(want) {
		t.Fatalf("RequiredPhase1Contracts length: got %d want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("RequiredPhase1Contracts[%d]: got %q want %q", i, got[i], want[i])
		}
	}

	// Defensive copy: mutating the returned slice must not change the source.
	got[0] = "mutated"
	again := RequiredPhase1Contracts()
	if again[0] != want[0] {
		t.Fatalf("RequiredPhase1Contracts must return a defensive copy")
	}
}

func TestPhase1CompatibilityCoverage_ReportCoversAll(t *testing.T) {
	t.Parallel()

	raw, err := os.ReadFile(filepath.Join(moduleRoot(t), Phase1CompatibilityReportPath))
	if err != nil {
		t.Fatalf("read %s: %v", Phase1CompatibilityReportPath, err)
	}
	if err := AssertPhase1CompatibilityCoverage(raw); err != nil {
		t.Fatalf("AssertPhase1CompatibilityCoverage(%s): %v", Phase1CompatibilityReportPath, err)
	}
}

func TestPhase1CompatibilityCoverage_MissingContractFails(t *testing.T) {
	t.Parallel()

	raw, err := os.ReadFile(filepath.Join(moduleRoot(t), Phase1CompatibilityReportPath))
	if err != nil {
		t.Fatalf("read %s: %v", Phase1CompatibilityReportPath, err)
	}

	// Remove the exact Tenant coverage heading while leaving other mentions.
	stripped := strings.Replace(string(raw), "### Tenant\n", "### TenantRemoved\n", 1)
	if stripped == string(raw) {
		t.Fatal("test setup failed: ### Tenant heading not found in report")
	}
	err = AssertPhase1CompatibilityCoverage([]byte(stripped))
	if err == nil {
		t.Fatal("expected missing Tenant coverage to fail")
	}
	if !strings.Contains(err.Error(), "Tenant") {
		t.Fatalf("error must name missing contract Tenant, got: %v", err)
	}
	if !strings.Contains(err.Error(), "missing coverage") {
		t.Fatalf("error must indicate missing coverage, got: %v", err)
	}
}

func TestPhase1CompatibilityCoverage_EmptyFails(t *testing.T) {
	t.Parallel()

	if err := AssertPhase1CompatibilityCoverage(nil); err == nil {
		t.Fatal("nil report must fail")
	}
	if err := AssertPhase1CompatibilityCoverage([]byte("   \n\t")); err == nil {
		t.Fatal("whitespace-only report must fail")
	}
	if err := AssertPhase1CompatibilityCoverage([]byte("# No contracts\n")); err == nil {
		t.Fatal("report without required headings must fail")
	}
}

func TestPhase1CompatibilityCoverage_IncidentalHeadingDoesNotSatisfy(t *testing.T) {
	t.Parallel()

	// A heading that only contains a required ID as a substring must not count.
	report := []byte("### OrganizationExtra\n### health/readiness-extra\n")
	err := AssertPhase1CompatibilityCoverage(report)
	if err == nil {
		t.Fatal("substring/incidental headings must not satisfy coverage")
	}
	if !strings.Contains(err.Error(), "Organization") {
		t.Fatalf("expected Organization missing, got: %v", err)
	}
}
