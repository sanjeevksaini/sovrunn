package conformance

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/apiconform"
	"github.com/sanjeevksaini/sovrunn/internal/apivalid"
)

// TestPositiveDecisionFixturesValidateCleanly exercises every positive/compat
// scenario ID against its on-disk T-029 fixtures via T-027 RunDecisionCheck
// (F13-CONF-001; F13-SCALE-004; AD-028).
func TestPositiveDecisionFixturesValidateCleanly(t *testing.T) {
	t.Parallel()

	root := moduleRoot(t)
	run := apiconform.DecisionCheckRun{
		ModuleRoot: root,
		DecodeMode: apivalid.ModeReadRepresentation,
	}

	for _, check := range apiconform.RegisteredDecisionChecks() {
		switch check.Kind {
		case apiconform.DecisionCheckPositive, apiconform.DecisionCheckCompat:
			// exercised below
		default:
			continue
		}
		check := check
		t.Run(string(check.ID), func(t *testing.T) {
			t.Parallel()
			res := apiconform.RunDecisionCheck(check, run)
			if !res.OK {
				t.Fatalf("positive fixture for %s failed: detail=%q problem=%#v",
					check.ID, res.Detail, res.Problem)
			}
		})
	}
}

// TestPositiveDecisionFixtureFilesExist asserts every positive DefaultFixturePath
// under decision/positive/ is present on disk (RID-07; T-029 acceptance).
func TestPositiveDecisionFixtureFilesExist(t *testing.T) {
	t.Parallel()

	root := moduleRoot(t)
	for _, check := range apiconform.RegisteredDecisionChecks() {
		if check.Kind != apiconform.DecisionCheckPositive && check.Kind != apiconform.DecisionCheckCompat {
			continue
		}
		for _, rel := range check.DefaultFixturePaths {
			path := filepath.Join(root, apiconform.ConformanceFixturesDir, filepath.FromSlash(rel))
			if _, err := os.ReadFile(path); err != nil {
				t.Fatalf("missing positive fixture %s for %s: %v", rel, check.ID, err)
			}
		}
	}
}

func moduleRoot(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	// tests/conformance/<thisfile> → repo root
	return filepath.Clean(filepath.Join(filepath.Dir(file), "../.."))
}
