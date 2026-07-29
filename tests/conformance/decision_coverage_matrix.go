// Package conformance holds FEATURE-0013 scenario-ID → fixture coverage
// mapping (T-028; design §6.4/§13/§18.4; F13-CONF-002/004; AD-028).
//
// Coverage is counted by scenario ID: one physical fixture may cover several
// IDs, but every parent, mandatory suffix, and additional case has its own
// matrix row. Physical positive/negative bodies are materialized by T-029 and
// T-030; this matrix is the authoritative ID→fixture binding source.
package conformance

import (
	"fmt"
	"sort"
	"strings"

	"github.com/sanjeevksaini/sovrunn/internal/apiconform"
)

// DecisionCoverageScenarioCount is the closed FEATURE-0013 conformance ID count:
// 28 parents + 49 compound suffixes + 11 SCOPE + 8 EVAL + 8 SEC + 4 TRUST + 9 COMPAT.
const DecisionCoverageScenarioCount = 117

// Logical fixture references for contract proofs that do not use a dedicated
// JSON body under tests/conformance/fixtures (schema agreement and replay kernel).
const (
	LogicalFixtureTypeBindings = "logical:apiconform-type-bindings"
	LogicalFixtureReplayKernel = "logical:compose-replay-kernel"
)

// DecisionCoverageRow maps one architecture-owned scenario ID to ≥1 fixture
// reference (RID-07; F13-CONF-004). Fixture paths are relative to
// apiconform.ConformanceFixturesDir unless they use a logical: prefix.
type DecisionCoverageRow struct {
	ID          apiconform.DecisionScenarioID
	Kind        apiconform.DecisionCheckKind
	Description string
	Fixtures    []string
}

// DecisionCoverageMatrix returns every FEATURE-0013 coverage row in stable
// registration order (parents, mandatory suffixes, then additional cases).
// Rows are derived from the T-027 executable check registry so IDs cannot
// drift; fixture bindings follow RID-07 path conventions and may share one
// physical fixture across multiple IDs.
func DecisionCoverageMatrix() []DecisionCoverageRow {
	checks := apiconform.RegisteredDecisionChecks()
	out := make([]DecisionCoverageRow, 0, len(checks))
	for _, c := range checks {
		out = append(out, DecisionCoverageRow{
			ID:          c.ID,
			Kind:        c.Kind,
			Description: c.Description,
			Fixtures:    fixturesForDecisionCoverage(c),
		})
	}
	return out
}

// ResolveDecisionCoverage looks up one scenario ID in the coverage matrix.
// Missing IDs fail closed (ok=false).
func ResolveDecisionCoverage(id apiconform.DecisionScenarioID) (DecisionCoverageRow, bool) {
	if id == "" {
		return DecisionCoverageRow{}, false
	}
	for _, row := range DecisionCoverageMatrix() {
		if row.ID == id {
			return row, true
		}
	}
	return DecisionCoverageRow{}, false
}

// DecisionCoverageByID returns a map of scenario ID → coverage row.
func DecisionCoverageByID() map[apiconform.DecisionScenarioID]DecisionCoverageRow {
	rows := DecisionCoverageMatrix()
	out := make(map[apiconform.DecisionScenarioID]DecisionCoverageRow, len(rows))
	for _, row := range rows {
		out[row.ID] = row
	}
	return out
}

// RequiredDecisionScenarioIDs returns the architecture-owned closed scenario ID
// set in stable registration order (architecture §17/§18.4).
func RequiredDecisionScenarioIDs() []apiconform.DecisionScenarioID {
	return apiconform.RegisteredDecisionScenarioIDs()
}

// AssertDecisionCoverageComplete verifies full ID coverage: every required
// scenario ID has exactly one matrix row with ≥1 fixture, no invented IDs,
// and no empty fixture lists (F13-CONF-002/004).
func AssertDecisionCoverageComplete() error {
	rows := DecisionCoverageMatrix()
	if len(rows) != DecisionCoverageScenarioCount {
		return fmt.Errorf("decision coverage matrix length = %d, want %d", len(rows), DecisionCoverageScenarioCount)
	}

	required := RequiredDecisionScenarioIDs()
	if len(required) != DecisionCoverageScenarioCount {
		return fmt.Errorf("required decision scenario IDs length = %d, want %d", len(required), DecisionCoverageScenarioCount)
	}

	byID := make(map[apiconform.DecisionScenarioID]DecisionCoverageRow, len(rows))
	for i, row := range rows {
		if row.ID == "" {
			return fmt.Errorf("coverage row %d has empty scenario ID", i)
		}
		if _, dup := byID[row.ID]; dup {
			return fmt.Errorf("duplicate coverage row for scenario ID %q", row.ID)
		}
		if row.Kind == "" {
			return fmt.Errorf("coverage row %q has empty kind", row.ID)
		}
		if len(row.Fixtures) == 0 {
			return fmt.Errorf("coverage row %q resolves to zero fixtures", row.ID)
		}
		for j, fx := range row.Fixtures {
			if strings.TrimSpace(fx) == "" {
				return fmt.Errorf("coverage row %q fixture[%d] is empty", row.ID, j)
			}
		}
		byID[row.ID] = row
		if required[i] != row.ID {
			return fmt.Errorf("coverage matrix order drift at %d: matrix=%q required=%q", i, row.ID, required[i])
		}
	}

	var missing []string
	for _, id := range required {
		if _, ok := byID[id]; !ok {
			missing = append(missing, string(id))
		}
	}
	if len(missing) > 0 {
		sort.Strings(missing)
		return fmt.Errorf("coverage matrix missing required IDs: %s", strings.Join(missing, ", "))
	}

	// Fail closed on invented IDs (matrix larger than required set).
	requiredSet := make(map[apiconform.DecisionScenarioID]struct{}, len(required))
	for _, id := range required {
		requiredSet[id] = struct{}{}
	}
	var invented []string
	for id := range byID {
		if _, ok := requiredSet[id]; !ok {
			invented = append(invented, string(id))
		}
	}
	if len(invented) > 0 {
		sort.Strings(invented)
		return fmt.Errorf("coverage matrix has invented IDs: %s", strings.Join(invented, ", "))
	}

	if err := assertMandatorySuffixesPresent(byID); err != nil {
		return err
	}
	return nil
}

func fixturesForDecisionCoverage(c apiconform.DecisionCheck) []string {
	// Prefer registry default paths (RID-07 conventions used by T-027/T-029/T-030).
	if len(c.DefaultFixturePaths) > 0 {
		out := append([]string(nil), c.DefaultFixturePaths...)
		// Shared FEATURE-0012 AuditEvent body also covers COMPAT-08 six-scope
		// expansion alongside the planned positive companion fixture.
		if c.ID == "F13-COMPAT-08" {
			out = ensureFixture(out, "audit-event.json")
		}
		return out
	}

	switch c.Kind {
	case apiconform.DecisionCheckSchema:
		return []string{LogicalFixtureTypeBindings}
	case apiconform.DecisionCheckReplay:
		return []string{LogicalFixtureReplayKernel}
	case apiconform.DecisionCheckNegative:
		return []string{"negative/decision/" + string(c.ID) + ".json"}
	case apiconform.DecisionCheckCompat:
		if c.ID == "F13-COMPAT-07" || c.ID == "F13-COMPAT-08" {
			return []string{"audit-event.json", "decision/positive/" + string(c.ID) + ".json"}
		}
		return []string{"decision/positive/" + string(c.ID) + ".json"}
	default:
		return []string{"decision/positive/" + string(c.ID) + ".json"}
	}
}

func ensureFixture(paths []string, want string) []string {
	for _, p := range paths {
		if p == want {
			return paths
		}
	}
	return append(paths, want)
}

func assertMandatorySuffixesPresent(byID map[apiconform.DecisionScenarioID]DecisionCoverageRow) error {
	// Architecture §17 / design §18.4 — suffixes must not be merged, omitted,
	// renumbered, or invented.
	mandatory := []string{
		"F13-CF-04.a", "F13-CF-04.b", "F13-CF-04.c", "F13-CF-04.d", "F13-CF-04.e",
		"F13-CF-04.f", "F13-CF-04.g", "F13-CF-04.h", "F13-CF-04.i",
		"F13-CF-08.a", "F13-CF-08.b", "F13-CF-08.c", "F13-CF-08.d",
		"F13-CF-11.a", "F13-CF-11.b", "F13-CF-11.c", "F13-CF-11.d", "F13-CF-11.e",
		"F13-CF-12.a", "F13-CF-12.b", "F13-CF-12.c", "F13-CF-12.d", "F13-CF-12.e",
		"F13-CF-16.a", "F13-CF-16.b", "F13-CF-16.c", "F13-CF-16.d", "F13-CF-16.e", "F13-CF-16.f",
		"F13-CF-20.a", "F13-CF-20.b", "F13-CF-20.c", "F13-CF-20.d", "F13-CF-20.e", "F13-CF-20.f",
		"F13-CF-26.a", "F13-CF-26.b", "F13-CF-26.c", "F13-CF-26.d",
		"F13-CF-26.e", "F13-CF-26.f", "F13-CF-26.g", "F13-CF-26.h",
		"F13-CF-28.a", "F13-CF-28.b", "F13-CF-28.c", "F13-CF-28.d", "F13-CF-28.e", "F13-CF-28.f",
	}
	for _, id := range mandatory {
		row, ok := byID[apiconform.DecisionScenarioID(id)]
		if !ok {
			return fmt.Errorf("mandatory suffix %q missing from coverage matrix", id)
		}
		if len(row.Fixtures) == 0 {
			return fmt.Errorf("mandatory suffix %q has no fixture binding", id)
		}
	}
	return nil
}
