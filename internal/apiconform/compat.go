package apiconform

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
)

// Phase1CompatibilityReportPath is the repository-relative Phase 1
// compatibility report required by D-13 / F12-COMPAT-001.
const Phase1CompatibilityReportPath = "docs/api/PHASE1_COMPATIBILITY_REPORT.md"

// requiredPhase1Contracts is the exact F12-COMPAT-001 inventory. Order matches
// the requirement and the Phase 1 compatibility report coverage table.
var requiredPhase1Contracts = []string{
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

// RequiredPhase1Contracts returns a copy of the F12-COMPAT-001 contract IDs
// that docs/api/PHASE1_COMPATIBILITY_REPORT.md must cover.
func RequiredPhase1Contracts() []string {
	out := make([]string, len(requiredPhase1Contracts))
	copy(out, requiredPhase1Contracts)
	return out
}

// AssertPhase1CompatibilityCoverage verifies that reportMarkdown covers every
// required Phase 1 resource/endpoint from F12-COMPAT-001.
//
// Coverage evidence is a dedicated Markdown H3 heading whose text is exactly
// the contract identifier (for example "### Organization", "### health/readiness").
// A missing required heading fails closed with a deterministic error listing
// the uncovered contract IDs.
func AssertPhase1CompatibilityCoverage(reportMarkdown []byte) error {
	if len(bytes.TrimSpace(reportMarkdown)) == 0 {
		return fmt.Errorf("phase 1 compatibility report is empty")
	}

	covered := phase1ContractHeadings(reportMarkdown)
	missing := make([]string, 0)
	for _, id := range requiredPhase1Contracts {
		if !covered[id] {
			missing = append(missing, id)
		}
	}
	if len(missing) > 0 {
		return fmt.Errorf("phase 1 compatibility report missing coverage for: %s", strings.Join(missing, ", "))
	}
	return nil
}

// phase1ContractHeadings returns the set of exact H3 heading texts found in
// the report. Only headings that match a required contract ID are retained so
// incidental H3 headings do not satisfy coverage.
func phase1ContractHeadings(reportMarkdown []byte) map[string]bool {
	required := make(map[string]struct{}, len(requiredPhase1Contracts))
	for _, id := range requiredPhase1Contracts {
		required[id] = struct{}{}
	}

	covered := make(map[string]bool, len(requiredPhase1Contracts))
	scanner := bufio.NewScanner(bytes.NewReader(reportMarkdown))
	// Allow long Markdown lines without truncating contract IDs.
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if !strings.HasPrefix(line, "### ") {
			continue
		}
		heading := strings.TrimSpace(strings.TrimPrefix(line, "### "))
		if _, ok := required[heading]; ok {
			covered[heading] = true
		}
	}
	return covered
}
