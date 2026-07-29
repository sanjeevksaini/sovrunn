package conformance

import (
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/apiconform"
)

func TestDecisionCoverageMatrixEveryIDResolvesToFixture(t *testing.T) {
	t.Parallel()

	if err := AssertDecisionCoverageComplete(); err != nil {
		t.Fatalf("AssertDecisionCoverageComplete: %v", err)
	}

	rows := DecisionCoverageMatrix()
	if len(rows) != DecisionCoverageScenarioCount {
		t.Fatalf("DecisionCoverageMatrix len=%d, want %d", len(rows), DecisionCoverageScenarioCount)
	}

	required := RequiredDecisionScenarioIDs()
	for _, id := range required {
		row, ok := ResolveDecisionCoverage(id)
		if !ok {
			t.Fatalf("ResolveDecisionCoverage(%q) = false", id)
		}
		if row.ID != id {
			t.Fatalf("ResolveDecisionCoverage(%q).ID = %q", id, row.ID)
		}
		if len(row.Fixtures) < 1 {
			t.Fatalf("scenario %q resolves to %d fixtures, want ≥1", id, len(row.Fixtures))
		}
		for _, fx := range row.Fixtures {
			if strings.TrimSpace(fx) == "" {
				t.Fatalf("scenario %q has empty fixture path", id)
			}
			if strings.HasPrefix(fx, "logical:") {
				switch fx {
				case LogicalFixtureTypeBindings, LogicalFixtureReplayKernel:
					// OK — schema/replay contract proofs.
				default:
					t.Fatalf("scenario %q uses unknown logical fixture %q", id, fx)
				}
				continue
			}
			// Disk paths are relative to tests/conformance/fixtures (RID-07).
			if filepath.IsAbs(fx) {
				t.Fatalf("scenario %q fixture %q must be relative to %s", id, fx, apiconform.ConformanceFixturesDir)
			}
			if strings.Contains(fx, `\`) {
				t.Fatalf("scenario %q fixture %q must use slash separators", id, fx)
			}
		}
	}

	if _, ok := ResolveDecisionCoverage("F13-CF-99"); ok {
		t.Fatal("unknown scenario ID must not resolve")
	}
}

func TestDecisionCoverageMatrixMatchesApiconformRegistry(t *testing.T) {
	t.Parallel()

	matrix := DecisionCoverageByID()
	checks := apiconform.RegisteredDecisionChecks()
	if len(matrix) != len(checks) {
		t.Fatalf("coverage map size=%d, registry size=%d", len(matrix), len(checks))
	}
	for _, c := range checks {
		row, ok := matrix[c.ID]
		if !ok {
			t.Fatalf("registry scenario %q missing from coverage matrix", c.ID)
		}
		if row.Kind != c.Kind {
			t.Fatalf("scenario %q kind matrix=%q registry=%q", c.ID, row.Kind, c.Kind)
		}
	}
}

func TestDecisionCoverageMatrixCompoundSuffixesNotMerged(t *testing.T) {
	t.Parallel()

	// Each suffix must own a distinct row; parents must also remain present.
	parents := []apiconform.DecisionScenarioID{
		"F13-CF-04", "F13-CF-08", "F13-CF-11", "F13-CF-12",
		"F13-CF-16", "F13-CF-20", "F13-CF-26", "F13-CF-28",
	}
	for _, id := range parents {
		if _, ok := ResolveDecisionCoverage(id); !ok {
			t.Fatalf("compound parent %q missing", id)
		}
	}

	suffixes := []apiconform.DecisionScenarioID{
		"F13-CF-04.a", "F13-CF-04.i",
		"F13-CF-08.a", "F13-CF-08.d",
		"F13-CF-11.a", "F13-CF-11.e",
		"F13-CF-12.a", "F13-CF-12.e",
		"F13-CF-16.a", "F13-CF-16.f",
		"F13-CF-20.a", "F13-CF-20.f",
		"F13-CF-26.a", "F13-CF-26.h",
		"F13-CF-28.a", "F13-CF-28.f",
	}
	seenFixtures := make(map[apiconform.DecisionScenarioID][]string, len(suffixes))
	for _, id := range suffixes {
		row, ok := ResolveDecisionCoverage(id)
		if !ok {
			t.Fatalf("mandatory suffix %q missing", id)
		}
		seenFixtures[id] = append([]string(nil), row.Fixtures...)
	}
	// Suffixes must not be collapsed into a single invented ID.
	if len(seenFixtures) != len(suffixes) {
		t.Fatalf("suffix rows collapsed: got %d want %d", len(seenFixtures), len(suffixes))
	}
}

func TestDecisionCoverageSharedFixtureMayCoverMultipleIDs(t *testing.T) {
	t.Parallel()

	// F13-CONF-004: one physical fixture may cover several IDs; each ID still
	// has its own row. COMPAT-07 and COMPAT-08 both bind audit-event.json.
	compat07, ok := ResolveDecisionCoverage("F13-COMPAT-07")
	if !ok {
		t.Fatal("F13-COMPAT-07 missing")
	}
	compat08, ok := ResolveDecisionCoverage("F13-COMPAT-08")
	if !ok {
		t.Fatal("F13-COMPAT-08 missing")
	}
	if !containsFixture(compat07.Fixtures, "audit-event.json") {
		t.Fatalf("F13-COMPAT-07 fixtures = %v, want audit-event.json", compat07.Fixtures)
	}
	if !containsFixture(compat08.Fixtures, "audit-event.json") {
		t.Fatalf("F13-COMPAT-08 fixtures = %v, want audit-event.json", compat08.Fixtures)
	}
}

func TestDecisionCoverageSchemaAndReplayLogicalFixtures(t *testing.T) {
	t.Parallel()

	schema, ok := ResolveDecisionCoverage("F13-CF-17")
	if !ok {
		t.Fatal("F13-CF-17 missing")
	}
	if !containsFixture(schema.Fixtures, LogicalFixtureTypeBindings) {
		t.Fatalf("F13-CF-17 fixtures = %v, want %s", schema.Fixtures, LogicalFixtureTypeBindings)
	}

	compat09, ok := ResolveDecisionCoverage("F13-COMPAT-09")
	if !ok {
		t.Fatal("F13-COMPAT-09 missing")
	}
	if !containsFixture(compat09.Fixtures, LogicalFixtureTypeBindings) {
		t.Fatalf("F13-COMPAT-09 fixtures = %v, want %s", compat09.Fixtures, LogicalFixtureTypeBindings)
	}

	replay, ok := ResolveDecisionCoverage("F13-CF-21")
	if !ok {
		t.Fatal("F13-CF-21 missing")
	}
	if !containsFixture(replay.Fixtures, LogicalFixtureReplayKernel) {
		t.Fatalf("F13-CF-21 fixtures = %v, want %s", replay.Fixtures, LogicalFixtureReplayKernel)
	}
}

func TestDecisionCoverageAdditionalFamiliesPresent(t *testing.T) {
	t.Parallel()

	want := []apiconform.DecisionScenarioID{
		"F13-SCOPE-01", "F13-SCOPE-11",
		"F13-EVAL-01", "F13-EVAL-08",
		"F13-SEC-01", "F13-SEC-08",
		"F13-TRUST-01", "F13-TRUST-04",
		"F13-COMPAT-01", "F13-COMPAT-09",
	}
	for _, id := range want {
		row, ok := ResolveDecisionCoverage(id)
		if !ok {
			t.Fatalf("additional case %q missing", id)
		}
		if len(row.Fixtures) < 1 {
			t.Fatalf("additional case %q has no fixtures", id)
		}
	}
}

func containsFixture(fixtures []string, want string) bool {
	for _, fx := range fixtures {
		if fx == want {
			return true
		}
	}
	return false
}

// Ensure the package lives at tests/conformance (compile-time path sanity for
// humans reading failures); not used as a runtime fixture existence gate.
func TestDecisionCoveragePackagePath(t *testing.T) {
	t.Parallel()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	if !strings.Contains(filepath.ToSlash(file), "tests/conformance/") {
		t.Fatalf("coverage matrix test path = %q, want tests/conformance/", file)
	}
}
