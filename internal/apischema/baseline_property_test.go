package apischema

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Deterministic seed for Property 6 reproducibility
// (F12-EVOLVE-002, F12-VERIFY-001(10)).
const property6Seed int64 = 20260723

const property6Iterations = 100

var property6SchemaNames = []string{
	"project.json",
	"resource-pool.json",
	"plugin-definition.json",
	"operation.json",
	"audit-event.json",
}

// property6Scenario classifies the generated baseline-gate case.
type property6Scenario string

const (
	property6Unchanged           property6Scenario = "unchanged"
	property6Bootstrap           property6Scenario = "bootstrap"
	property6ValidADH            property6Scenario = "valid_adh"
	property6ValidToken          property6Scenario = "valid_token"
	property6CoEditNoEvidence    property6Scenario = "coedit_no_evidence"
	property6WrongOldDigest      property6Scenario = "wrong_old_digest"
	property6WrongNewDigest      property6Scenario = "wrong_new_digest"
	property6IncompleteEvidence  property6Scenario = "incomplete_evidence"
	property6IntegrityTamper     property6Scenario = "integrity_tamper"
	property6NewFileWithEvidence property6Scenario = "new_file_with_evidence"
	property6NewFileNoEvidence   property6Scenario = "new_file_no_evidence"
	property6EvidenceWrongPath   property6Scenario = "evidence_wrong_path"
)

type property6Case struct {
	Scenario      property6Scenario
	WantPass      bool
	WantIntegrity bool // false only for integrity_tamper
	FailHint      string
}

// Feature: api-resource-naming-status-and-validation-standard, Property 6: Controlled baseline updates
//
// For any baseline change, VerifyBaselineApproval fails unless accompanied by
// recorded approval evidence with matching old/new digests and an approving
// ADH or approval token plus reviewer and date. Co-editing a baseline file and
// its BASELINE_MANIFEST.json digest in the same commit, without matching
// approval evidence, is never sufficient to pass the gate. Integrity failures
// (tampered baseline vs manifest) also fail the gate.
//
// Validates: Requirements 4.14, 4.16 (F12-EVOLVE-002, F12-VERIFY-001(10))
func TestProperty6_ControlledBaselineUpdates(t *testing.T) {
	t.Parallel()

	rng := rand.New(rand.NewSource(property6Seed))
	for i := 0; i < property6Iterations; i++ {
		c := generateProperty6Case(rng, i)
		if err := checkProperty6Case(t, c, i); err != nil {
			t.Fatalf("property 6 failed at iteration %d (seed %d scenario %s): %v",
				i, property6Seed, c.Scenario, err)
		}
	}
}

func generateProperty6Case(rng *rand.Rand, iteration int) property6Case {
	// Force coverage of every oracle class; rng occasionally shuffles bucket
	// selection so the 100 iterations remain seed-reproducible but not purely
	// round-robin when a later failure seed is reported.
	bucket := iteration % 12
	if rng.Intn(20) == 0 {
		bucket = rng.Intn(12)
	}
	switch bucket {
	case 0:
		return property6Case{
			Scenario:      property6Unchanged,
			WantPass:      true,
			WantIntegrity: true,
		}
	case 1:
		return property6Case{
			Scenario:      property6Bootstrap,
			WantPass:      true,
			WantIntegrity: true,
		}
	case 2:
		return property6Case{
			Scenario:      property6ValidADH,
			WantPass:      true,
			WantIntegrity: true,
		}
	case 3:
		return property6Case{
			Scenario:      property6ValidToken,
			WantPass:      true,
			WantIntegrity: true,
		}
	case 4:
		return property6Case{
			Scenario:      property6CoEditNoEvidence,
			WantPass:      false,
			WantIntegrity: true,
			FailHint:      "without recorded approval evidence",
		}
	case 5:
		return property6Case{
			Scenario:      property6WrongOldDigest,
			WantPass:      false,
			WantIntegrity: true,
			FailHint:      "without recorded approval evidence",
		}
	case 6:
		return property6Case{
			Scenario:      property6WrongNewDigest,
			WantPass:      false,
			WantIntegrity: true,
			FailHint:      "without recorded approval evidence",
		}
	case 7:
		return property6Case{
			Scenario:      property6IncompleteEvidence,
			WantPass:      false,
			WantIntegrity: true,
			FailHint:      "reviewer",
		}
	case 8:
		return property6Case{
			Scenario:      property6IntegrityTamper,
			WantPass:      false,
			WantIntegrity: false,
			FailHint:      "digest mismatch",
		}
	case 9:
		return property6Case{
			Scenario:      property6NewFileWithEvidence,
			WantPass:      true,
			WantIntegrity: true,
		}
	case 10:
		return property6Case{
			Scenario:      property6NewFileNoEvidence,
			WantPass:      false,
			WantIntegrity: true,
			FailHint:      "without recorded approval evidence",
		}
	default:
		return property6Case{
			Scenario:      property6EvidenceWrongPath,
			WantPass:      false,
			WantIntegrity: true,
			FailHint:      "without recorded approval evidence",
		}
	}
}

func checkProperty6Case(t *testing.T, c property6Case, iteration int) error {
	t.Helper()

	dir := t.TempDir()
	name := property6SchemaNames[iteration%len(property6SchemaNames)]
	oldContent := property6SchemaBytes(iteration, "old")
	newContent := property6SchemaBytes(iteration, "new")
	oldDigest := sha256Hex(oldContent)
	newDigest := sha256Hex(newContent)
	if oldDigest == newDigest {
		return fmt.Errorf("iteration %d: generator produced identical digests", iteration)
	}

	switch c.Scenario {
	case property6Unchanged:
		if err := os.WriteFile(filepath.Join(dir, name), oldContent, 0o644); err != nil {
			return err
		}
		writeManifest(t, dir, map[string]string{name: oldDigest})
		writeApprovals(t, dir, baselineApprovalsFile{
			RecordedDigests: map[string]string{name: oldDigest},
			Approvals:       []baselineApproval{},
		})

	case property6Bootstrap:
		if err := os.WriteFile(filepath.Join(dir, name), newContent, 0o644); err != nil {
			return err
		}
		writeManifest(t, dir, map[string]string{name: newDigest})
		writeApprovals(t, dir, baselineApprovalsFile{}) // empty recordedDigests

	case property6ValidADH:
		if err := os.WriteFile(filepath.Join(dir, name), newContent, 0o644); err != nil {
			return err
		}
		writeManifest(t, dir, map[string]string{name: newDigest})
		writeApprovals(t, dir, baselineApprovalsFile{
			RecordedDigests: map[string]string{name: oldDigest},
			Approvals: []baselineApproval{{
				Path:      name,
				OldDigest: oldDigest,
				NewDigest: newDigest,
				ADH:       fmt.Sprintf("ADH-2026-%03d", iteration%1000),
				Reviewer:  "Sanjeev Kumar",
				Date:      "2026-07-23",
			}},
		})

	case property6ValidToken:
		if err := os.WriteFile(filepath.Join(dir, name), newContent, 0o644); err != nil {
			return err
		}
		writeManifest(t, dir, map[string]string{name: newDigest})
		writeApprovals(t, dir, baselineApprovalsFile{
			RecordedDigests: map[string]string{name: oldDigest},
			Approvals: []baselineApproval{{
				Path:          name,
				OldDigest:     oldDigest,
				NewDigest:     newDigest,
				ApprovalToken: "APPROVED_BASELINE_CHANGE",
				Reviewer:      "Baseline Reviewer",
				Date:          "2026-07-23",
			}},
		})

	case property6CoEditNoEvidence:
		// Co-edit: baseline + manifest updated together; no approval evidence.
		if err := os.WriteFile(filepath.Join(dir, name), newContent, 0o644); err != nil {
			return err
		}
		writeManifest(t, dir, map[string]string{name: newDigest})
		writeApprovals(t, dir, baselineApprovalsFile{
			RecordedDigests: map[string]string{name: oldDigest},
			Approvals:       []baselineApproval{},
		})

	case property6WrongOldDigest:
		if err := os.WriteFile(filepath.Join(dir, name), newContent, 0o644); err != nil {
			return err
		}
		writeManifest(t, dir, map[string]string{name: newDigest})
		writeApprovals(t, dir, baselineApprovalsFile{
			RecordedDigests: map[string]string{name: oldDigest},
			Approvals: []baselineApproval{{
				Path:      name,
				OldDigest: sha256Hex([]byte("not-the-old-content")),
				NewDigest: newDigest,
				ADH:       "ADH-2026-012",
				Reviewer:  "Sanjeev Kumar",
				Date:      "2026-07-23",
			}},
		})

	case property6WrongNewDigest:
		if err := os.WriteFile(filepath.Join(dir, name), newContent, 0o644); err != nil {
			return err
		}
		writeManifest(t, dir, map[string]string{name: newDigest})
		writeApprovals(t, dir, baselineApprovalsFile{
			RecordedDigests: map[string]string{name: oldDigest},
			Approvals: []baselineApproval{{
				Path:      name,
				OldDigest: oldDigest,
				NewDigest: sha256Hex([]byte("not-the-new-content")),
				ADH:       "ADH-2026-012",
				Reviewer:  "Sanjeev Kumar",
				Date:      "2026-07-23",
			}},
		})

	case property6IncompleteEvidence:
		if err := os.WriteFile(filepath.Join(dir, name), newContent, 0o644); err != nil {
			return err
		}
		writeManifest(t, dir, map[string]string{name: newDigest})
		writeApprovals(t, dir, baselineApprovalsFile{
			RecordedDigests: map[string]string{name: oldDigest},
			Approvals: []baselineApproval{{
				Path:      name,
				OldDigest: oldDigest,
				NewDigest: newDigest,
				ADH:       "ADH-2026-012",
				// missing reviewer and date
			}},
		})

	case property6IntegrityTamper:
		if err := os.WriteFile(filepath.Join(dir, name), oldContent, 0o644); err != nil {
			return err
		}
		writeManifest(t, dir, map[string]string{name: oldDigest})
		writeApprovals(t, dir, baselineApprovalsFile{
			RecordedDigests: map[string]string{name: oldDigest},
		})
		// Tamper on-disk baseline without updating the manifest.
		if err := os.WriteFile(filepath.Join(dir, name), newContent, 0o644); err != nil {
			return err
		}

	case property6NewFileWithEvidence:
		other := property6OtherName(name)
		otherContent := property6SchemaBytes(iteration, "other")
		otherDigest := sha256Hex(otherContent)
		if err := os.WriteFile(filepath.Join(dir, other), otherContent, 0o644); err != nil {
			return err
		}
		if err := os.WriteFile(filepath.Join(dir, name), newContent, 0o644); err != nil {
			return err
		}
		writeManifest(t, dir, map[string]string{
			other: otherDigest,
			name:  newDigest,
		})
		writeApprovals(t, dir, baselineApprovalsFile{
			RecordedDigests: map[string]string{other: otherDigest},
			Approvals: []baselineApproval{{
				Path:      name,
				OldDigest: "",
				NewDigest: newDigest,
				ADH:       "ADH-2026-013",
				Reviewer:  "Sanjeev Kumar",
				Date:      "2026-07-23",
			}},
		})

	case property6NewFileNoEvidence:
		other := property6OtherName(name)
		otherContent := property6SchemaBytes(iteration, "other")
		otherDigest := sha256Hex(otherContent)
		if err := os.WriteFile(filepath.Join(dir, other), otherContent, 0o644); err != nil {
			return err
		}
		if err := os.WriteFile(filepath.Join(dir, name), newContent, 0o644); err != nil {
			return err
		}
		writeManifest(t, dir, map[string]string{
			other: otherDigest,
			name:  newDigest,
		})
		writeApprovals(t, dir, baselineApprovalsFile{
			RecordedDigests: map[string]string{other: otherDigest},
			Approvals:       []baselineApproval{},
		})

	case property6EvidenceWrongPath:
		if err := os.WriteFile(filepath.Join(dir, name), newContent, 0o644); err != nil {
			return err
		}
		writeManifest(t, dir, map[string]string{name: newDigest})
		writeApprovals(t, dir, baselineApprovalsFile{
			RecordedDigests: map[string]string{name: oldDigest},
			Approvals: []baselineApproval{{
				Path:      property6OtherName(name),
				OldDigest: oldDigest,
				NewDigest: newDigest,
				ADH:       "ADH-2026-012",
				Reviewer:  "Sanjeev Kumar",
				Date:      "2026-07-23",
			}},
		})

	default:
		return fmt.Errorf("iteration %d: unknown scenario %q", iteration, c.Scenario)
	}

	integrityErr := VerifyBaselineIntegrity(filepath.Join(dir, BaselineManifestFileName), dir)
	if c.WantIntegrity {
		if integrityErr != nil {
			return fmt.Errorf("iteration %d scenario %s: integrity must pass before approval oracle, got %v",
				iteration, c.Scenario, integrityErr)
		}
	} else if integrityErr == nil {
		return fmt.Errorf("iteration %d scenario %s: integrity must fail for tampered baseline",
			iteration, c.Scenario)
	}

	err := VerifyBaselineApproval(
		filepath.Join(dir, BaselineApprovalsFileName),
		filepath.Join(dir, BaselineManifestFileName),
		dir,
	)

	if c.WantPass {
		if err != nil {
			return fmt.Errorf("iteration %d scenario %s: want pass, got %v", iteration, c.Scenario, err)
		}
		return nil
	}

	if err == nil {
		return fmt.Errorf("iteration %d scenario %s: gate must fail (co-edit without evidence / bad evidence / integrity never passes)",
			iteration, c.Scenario)
	}
	if c.FailHint != "" && !strings.Contains(err.Error(), c.FailHint) {
		return fmt.Errorf("iteration %d scenario %s: error %q missing hint %q",
			iteration, c.Scenario, err.Error(), c.FailHint)
	}

	// Co-edit without evidence must never be accepted: integrity alone is insufficient.
	if c.Scenario == property6CoEditNoEvidence {
		if integrityErr != nil {
			return fmt.Errorf("iteration %d: co-edit setup must pass integrity: %v", iteration, integrityErr)
		}
		if !strings.Contains(err.Error(), "without recorded approval evidence") {
			return fmt.Errorf("iteration %d: co-edit without evidence must cite missing evidence, got %v",
				iteration, err)
		}
	}

	return nil
}

func property6SchemaBytes(iteration int, tag string) []byte {
	return []byte(fmt.Sprintf(
		`{"$schema":"https://json-schema.org/draft/2020-12/schema","type":"object","title":"%s-%d","description":"property6 %s"}`,
		tag, iteration, tag,
	))
}

func property6OtherName(name string) string {
	for _, n := range property6SchemaNames {
		if n != name {
			return n
		}
	}
	return "other.json"
}
