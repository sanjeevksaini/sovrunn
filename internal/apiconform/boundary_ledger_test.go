package apiconform

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
)

// BoundaryLedgerPath is the repository-relative machine-readable ledger
// (D-12, F12-LEDGER-001). BOUNDARY_LEDGER.md is a later derivative view.
const BoundaryLedgerPath = "docs/api/boundary-ledger.yaml"

// F12-LEDGER-001 categories. allowed_data and prohibited_data together
// represent the "allowed/prohibited data" category.
var requiredLedgerCategories = []string{
	"purpose",
	"owner",
	"producers",
	"consumers",
	"allowed_data",
	"prohibited_data",
	"authorization",
	"audit",
	"observability",
	"failure_behavior",
	"versioning",
	"replacement_path",
	"migration_path",
	"reassessment_trigger",
}

// ledgerDoc is the strict typed shape for docs/api/boundary-ledger.yaml.
// Unknown top-level or entry fields fail KnownFields decoding.
type ledgerDoc struct {
	APIVersion string           `yaml:"apiVersion"`
	Kind       string           `yaml:"kind"`
	Metadata   ledgerMetadata   `yaml:"metadata"`
	Boundaries []ledgerBoundary `yaml:"boundaries"`
}

type ledgerMetadata struct {
	Name        string `yaml:"name"`
	Feature     string `yaml:"feature"`
	Description string `yaml:"description"`
}

type ledgerBoundary struct {
	ID                  string   `yaml:"id"`
	Purpose             string   `yaml:"purpose"`
	Owner               string   `yaml:"owner"`
	Producers           []string `yaml:"producers"`
	Consumers           []string `yaml:"consumers"`
	AllowedData         []string `yaml:"allowed_data"`
	ProhibitedData      []string `yaml:"prohibited_data"`
	Authorization       string   `yaml:"authorization"`
	Audit               string   `yaml:"audit"`
	Observability       string   `yaml:"observability"`
	FailureBehavior     string   `yaml:"failure_behavior"`
	Versioning          string   `yaml:"versioning"`
	ReplacementPath     string   `yaml:"replacement_path"`
	MigrationPath       string   `yaml:"migration_path"`
	ReassessmentTrigger string   `yaml:"reassessment_trigger"`
}

func TestBoundaryLedgerStrictParseAndCategories(t *testing.T) {
	t.Parallel()

	path := filepath.Join(moduleRoot(t), BoundaryLedgerPath)
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", BoundaryLedgerPath, err)
	}

	var doc ledgerDoc
	dec := yaml.NewDecoder(bytes.NewReader(raw))
	dec.KnownFields(true)
	if err := dec.Decode(&doc); err != nil {
		t.Fatalf("strict YAML decode of %s: %v", BoundaryLedgerPath, err)
	}
	// Single-document ledger: a second decode must hit EOF.
	var extra any
	if err := dec.Decode(&extra); err == nil {
		t.Fatalf("%s: expected single YAML document, found additional document", BoundaryLedgerPath)
	}

	if doc.APIVersion == "" || doc.Kind == "" {
		t.Fatalf("ledger missing apiVersion/kind")
	}
	if strings.TrimSpace(doc.Metadata.Name) == "" {
		t.Fatalf("ledger metadata.name is required")
	}

	wantIDs := make([]string, 0, len(apimeta.AllBoundaries()))
	for _, b := range apimeta.AllBoundaries() {
		wantIDs = append(wantIDs, string(b))
	}
	if len(doc.Boundaries) != len(wantIDs) {
		t.Fatalf("want %d boundaries (%v), got %d", len(wantIDs), wantIDs, len(doc.Boundaries))
	}

	seen := make(map[string]int, len(doc.Boundaries))
	for i, entry := range doc.Boundaries {
		if !apimeta.Boundary(entry.ID).Valid() {
			t.Fatalf("boundaries[%d].id %q is not a Matrix C1 boundary", i, entry.ID)
		}
		seen[entry.ID]++
		if err := assertLedgerCategoriesPresent(entry); err != nil {
			t.Fatalf("boundaries[%d] id=%s: %v", i, entry.ID, err)
		}
	}
	for _, id := range wantIDs {
		if seen[id] != 1 {
			t.Errorf("boundary %q: want exactly one ledger entry, got %d", id, seen[id])
		}
	}
}

func assertLedgerCategoriesPresent(entry ledgerBoundary) error {
	// Map category name → non-empty check for F12-LEDGER-001.
	checks := map[string]bool{
		"purpose":              strings.TrimSpace(entry.Purpose) != "",
		"owner":                strings.TrimSpace(entry.Owner) != "",
		"producers":            len(entry.Producers) > 0 && allNonEmpty(entry.Producers),
		"consumers":            len(entry.Consumers) > 0 && allNonEmpty(entry.Consumers),
		"allowed_data":         len(entry.AllowedData) > 0 && allNonEmpty(entry.AllowedData),
		"prohibited_data":      len(entry.ProhibitedData) > 0 && allNonEmpty(entry.ProhibitedData),
		"authorization":        strings.TrimSpace(entry.Authorization) != "",
		"audit":                strings.TrimSpace(entry.Audit) != "",
		"observability":        strings.TrimSpace(entry.Observability) != "",
		"failure_behavior":     strings.TrimSpace(entry.FailureBehavior) != "",
		"versioning":           strings.TrimSpace(entry.Versioning) != "",
		"replacement_path":     strings.TrimSpace(entry.ReplacementPath) != "",
		"migration_path":       strings.TrimSpace(entry.MigrationPath) != "",
		"reassessment_trigger": strings.TrimSpace(entry.ReassessmentTrigger) != "",
	}
	for _, cat := range requiredLedgerCategories {
		ok, present := checks[cat]
		if !present {
			return fmt.Errorf("internal test bug: missing check for category %q", cat)
		}
		if !ok {
			return fmt.Errorf("F12-LEDGER-001 category %q is missing or empty", cat)
		}
	}
	return nil
}

func allNonEmpty(values []string) bool {
	for _, v := range values {
		if strings.TrimSpace(v) == "" {
			return false
		}
	}
	return true
}
