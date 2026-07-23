package apiconform

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
)

func TestBoundaryLedgerStrictParseAndCategories(t *testing.T) {
	t.Parallel()

	path := filepath.Join(moduleRoot(t), BoundaryLedgerPath)
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", BoundaryLedgerPath, err)
	}

	doc, err := ParseBoundaryLedgerYAML(raw)
	if err != nil {
		t.Fatalf("strict YAML decode of %s: %v", BoundaryLedgerPath, err)
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

func TestBoundaryLedgerMarkdownGenerator(t *testing.T) {
	t.Parallel()

	path := filepath.Join(moduleRoot(t), BoundaryLedgerPath)
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", BoundaryLedgerPath, err)
	}

	md1, err := GenerateBoundaryLedgerMarkdown(raw)
	if err != nil {
		t.Fatalf("GenerateBoundaryLedgerMarkdown: %v", err)
	}
	md2, err := GenerateBoundaryLedgerMarkdown(raw)
	if err != nil {
		t.Fatalf("GenerateBoundaryLedgerMarkdown second call: %v", err)
	}
	if !bytes.Equal(md1, md2) {
		t.Fatalf("Markdown generation is not deterministic: outputs differ")
	}
	if len(md1) == 0 {
		t.Fatalf("generated Markdown is empty")
	}
	if !bytes.HasSuffix(md1, []byte("\n")) {
		t.Fatalf("generated Markdown must end with a trailing newline")
	}
	if bytes.HasSuffix(md1, []byte("\n\n")) {
		t.Fatalf("generated Markdown must end with exactly one trailing newline")
	}
	if bytes.Contains(md1, []byte("\r")) {
		t.Fatalf("generated Markdown must use LF endings only")
	}

	text := string(md1)
	if !strings.HasPrefix(text, "# Boundary Ledger\n") {
		t.Fatalf("Markdown missing H1 title")
	}
	if !strings.Contains(text, BoundaryLedgerPath) {
		t.Fatalf("Markdown must reference source path %s", BoundaryLedgerPath)
	}
	for _, b := range apimeta.AllBoundaries() {
		heading := "### " + string(b)
		if !strings.Contains(text, heading) {
			t.Errorf("Markdown missing boundary heading %q", heading)
		}
	}
	for _, cat := range requiredLedgerCategories {
		title := categoryTitle[cat]
		if title == "" {
			t.Fatalf("missing categoryTitle for %q", cat)
		}
		if !strings.Contains(text, "#### "+title) {
			t.Errorf("Markdown missing category heading %q", title)
		}
	}

	// Render via parsed document must match GenerateBoundaryLedgerMarkdown.
	doc, err := ParseBoundaryLedgerYAML(raw)
	if err != nil {
		t.Fatalf("ParseBoundaryLedgerYAML: %v", err)
	}
	if got := RenderBoundaryLedgerMarkdown(doc); !bytes.Equal(got, md1) {
		t.Fatalf("RenderBoundaryLedgerMarkdown diverged from GenerateBoundaryLedgerMarkdown")
	}
}

func assertLedgerCategoriesPresent(entry BoundaryLedgerEntry) error {
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
