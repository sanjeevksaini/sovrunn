package apiconform

import (
	"path/filepath"
	"testing"
)

// TestMatrixDScenariosRepresentable verifies Task 14.4: the Matrix D scenario
// table maps all seventeen scenarios to fixtures plus required-proof
// assertions, and the suite fails if any scenario is unrepresentable
// (F12-FIXTURE-001, F12-FIXTURE-002, F12-VERIFY-001, D-17).
func TestMatrixDScenariosRepresentable(t *testing.T) {
	t.Parallel()

	root := moduleRoot(t)
	loader, err := NewFixtureLoader(root)
	if err != nil {
		t.Fatalf("NewFixtureLoader: %v", err)
	}
	if got := loader.Dir(); got != filepath.Join(root, ConformanceFixturesDir) {
		t.Fatalf("Dir = %q, want %q", got, filepath.Join(root, ConformanceFixturesDir))
	}

	scenarios := MatrixDScenarios()
	if len(scenarios) != MatrixDScenarioCount {
		t.Fatalf("MatrixDScenarios length = %d, want %d", len(scenarios), MatrixDScenarioCount)
	}

	for _, sc := range scenarios {
		sc := sc
		t.Run(sc.Name, func(t *testing.T) {
			t.Parallel()
			if sc.RequiredProof == "" {
				t.Fatal("RequiredProof is empty")
			}
			if sc.Prove == nil {
				t.Fatal("Prove assertion is nil")
			}
			if len(sc.Fixtures) == 0 {
				t.Fatal("Fixtures list is empty (scenario unrepresentable)")
			}
			if err := loader.RequireFixtures(sc.Fixtures...); err != nil {
				t.Fatalf("unrepresentable: %v", err)
			}
			if err := sc.Prove(loader); err != nil {
				t.Fatalf("required proof failed: %v", err)
			}
		})
	}
}

func TestAssertAllMatrixDScenarios(t *testing.T) {
	t.Parallel()

	loader, err := NewFixtureLoader(moduleRoot(t))
	if err != nil {
		t.Fatalf("NewFixtureLoader: %v", err)
	}
	if err := AssertAllMatrixDScenarios(loader); err != nil {
		t.Fatalf("AssertAllMatrixDScenarios: %v", err)
	}
}

func TestFixtureLoaderJSONAndYAML(t *testing.T) {
	t.Parallel()

	loader, err := NewFixtureLoader(moduleRoot(t))
	if err != nil {
		t.Fatalf("NewFixtureLoader: %v", err)
	}
	var fromJSON Project
	if err := loader.DecodeJSON("project.json", &fromJSON); err != nil {
		t.Fatalf("DecodeJSON: %v", err)
	}
	var fromYAML Project
	if err := loader.DecodeYAML("project.yaml", &fromYAML); err != nil {
		t.Fatalf("DecodeYAML: %v", err)
	}
	if fromJSON.Metadata.Name == "" || fromYAML.Metadata.Name == "" {
		t.Fatal("expected non-empty project name from JSON and YAML")
	}
	if fromJSON.Metadata.Name != fromYAML.Metadata.Name {
		t.Fatalf("name JSON=%q YAML=%q", fromJSON.Metadata.Name, fromYAML.Metadata.Name)
	}
}

func TestFutureProvisioningUsesSixOperationScopes(t *testing.T) {
	t.Parallel()

	var found *MatrixDScenario
	for i := range MatrixDScenarios() {
		sc := MatrixDScenarios()[i]
		if sc.Name == "Future provisioning executes" {
			found = &sc
			break
		}
	}
	if found == nil {
		t.Fatal("Future provisioning executes scenario missing from Matrix D table")
	}

	want := []string{
		"operation-platform.json",
		"operation-organization.json",
		"operation-organizationunit.json",
		"operation-tenant.json",
		"operation-project.json",
		"operation-provider.json",
	}
	for _, name := range want {
		foundFile := false
		for _, f := range found.Fixtures {
			if f == name {
				foundFile = true
				break
			}
		}
		if !foundFile {
			t.Fatalf("Future provisioning fixtures missing Operation scope variant %s", name)
		}
	}
}
