package conformance

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/apiconform"
	"github.com/sanjeevksaini/sovrunn/internal/apivalid"
)

// TestNegativeDecisionFixturesRejectWithExpectedViolation exercises every
// negative scenario ID against its on-disk T-030 fixtures via T-027
// RunDecisionCheck (F13-CONF-001; F13-CONF-004; AD-028).
func TestNegativeDecisionFixturesRejectWithExpectedViolation(t *testing.T) {
	t.Parallel()

	root := moduleRoot(t)
	run := apiconform.DecisionCheckRun{
		ModuleRoot: root,
		DecodeMode: apivalid.ModeReadRepresentation,
	}

	for _, check := range apiconform.RegisteredDecisionChecks() {
		if check.Kind != apiconform.DecisionCheckNegative {
			continue
		}
		check := check
		t.Run(string(check.ID), func(t *testing.T) {
			t.Parallel()
			res := apiconform.RunDecisionCheck(check, run)
			if !res.OK {
				t.Fatalf("negative fixture for %s failed: detail=%q problem=%#v",
					check.ID, res.Detail, res.Problem)
			}
		})
	}
}

// TestNegativeDecisionFixtureFilesExist asserts every negative DefaultFixturePath
// under negative/decision/ is present on disk (RID-07; T-030 acceptance).
func TestNegativeDecisionFixtureFilesExist(t *testing.T) {
	t.Parallel()

	root := moduleRoot(t)
	for _, check := range apiconform.RegisteredDecisionChecks() {
		if check.Kind != apiconform.DecisionCheckNegative {
			continue
		}
		for _, rel := range check.DefaultFixturePaths {
			path := filepath.Join(root, apiconform.ConformanceFixturesDir, filepath.FromSlash(rel))
			if _, err := os.ReadFile(path); err != nil {
				t.Fatalf("missing negative fixture %s for %s: %v", rel, check.ID, err)
			}
		}
	}
}
