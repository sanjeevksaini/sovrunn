package apiconform

import (
	"bytes"
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

// Boundary ledger paths (D-12, F12-LEDGER-001).
// YAML is the sole source of truth; Markdown is a regenerable derivative.
const (
	BoundaryLedgerPath         = "docs/api/boundary-ledger.yaml"
	BoundaryLedgerMarkdownPath = "docs/api/BOUNDARY_LEDGER.md"
)

// F12-LEDGER-001 categories in stable render order. allowed_data and
// prohibited_data together satisfy the "allowed/prohibited data" category.
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

// categoryTitle maps YAML field names to fixed Markdown headings.
var categoryTitle = map[string]string{
	"purpose":              "Purpose",
	"owner":                "Owner",
	"producers":            "Producers",
	"consumers":            "Consumers",
	"allowed_data":         "Allowed data",
	"prohibited_data":      "Prohibited data",
	"authorization":        "Authorization",
	"audit":                "Audit",
	"observability":        "Observability",
	"failure_behavior":     "Failure behavior",
	"versioning":           "Versioning",
	"replacement_path":     "Replacement path",
	"migration_path":       "Migration path",
	"reassessment_trigger": "Reassessment trigger",
}

// BoundaryLedger is the strict typed shape for docs/api/boundary-ledger.yaml.
// Unknown top-level or entry fields fail KnownFields decoding.
type BoundaryLedger struct {
	APIVersion string                `yaml:"apiVersion"`
	Kind       string                `yaml:"kind"`
	Metadata   BoundaryLedgerMeta    `yaml:"metadata"`
	Boundaries []BoundaryLedgerEntry `yaml:"boundaries"`
}

// BoundaryLedgerMeta is ledger document metadata.
type BoundaryLedgerMeta struct {
	Name        string `yaml:"name"`
	Feature     string `yaml:"feature"`
	Description string `yaml:"description"`
}

// BoundaryLedgerEntry is one Matrix C1 boundary ledger record.
type BoundaryLedgerEntry struct {
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

// ParseBoundaryLedgerYAML strictly decodes a single-document boundary ledger.
func ParseBoundaryLedgerYAML(raw []byte) (BoundaryLedger, error) {
	var doc BoundaryLedger
	dec := yaml.NewDecoder(bytes.NewReader(raw))
	dec.KnownFields(true)
	if err := dec.Decode(&doc); err != nil {
		return BoundaryLedger{}, fmt.Errorf("decode boundary ledger: %w", err)
	}
	var extra any
	if err := dec.Decode(&extra); err == nil {
		return BoundaryLedger{}, fmt.Errorf("boundary ledger must be a single YAML document")
	}
	return doc, nil
}

// GenerateBoundaryLedgerMarkdown parses YAML bytes and renders deterministic
// Markdown for docs/api/BOUNDARY_LEDGER.md. Same YAML always yields
// byte-identical Markdown (LF endings, stable section order, trailing newline).
func GenerateBoundaryLedgerMarkdown(raw []byte) ([]byte, error) {
	doc, err := ParseBoundaryLedgerYAML(raw)
	if err != nil {
		return nil, err
	}
	return RenderBoundaryLedgerMarkdown(doc), nil
}

// RenderBoundaryLedgerMarkdown renders a parsed ledger to deterministic Markdown.
// Output is fully determined by the decoded document fields and their order.
func RenderBoundaryLedgerMarkdown(doc BoundaryLedger) []byte {
	var b strings.Builder

	writeLine(&b, "# Boundary Ledger")
	writeLine(&b, "")
	writeLine(&b, "<!-- Generated from "+BoundaryLedgerPath+". Do not edit by hand. -->")
	writeLine(&b, "")
	writeLine(&b, "Machine-readable source of truth: `"+BoundaryLedgerPath+"`.")
	writeLine(&b, "This Markdown file is a regenerable human view (D-12, F12-LEDGER-001).")
	writeLine(&b, "")

	writeLine(&b, "## Document")
	writeLine(&b, "")
	writeLine(&b, "| Field | Value |")
	writeLine(&b, "| --- | --- |")
	writeLine(&b, "| apiVersion | "+escapeTableCell(strings.TrimSpace(doc.APIVersion))+" |")
	writeLine(&b, "| kind | "+escapeTableCell(strings.TrimSpace(doc.Kind))+" |")
	writeLine(&b, "| name | "+escapeTableCell(strings.TrimSpace(doc.Metadata.Name))+" |")
	writeLine(&b, "| feature | "+escapeTableCell(strings.TrimSpace(doc.Metadata.Feature))+" |")
	writeLine(&b, "")

	if desc := normalizeProse(doc.Metadata.Description); desc != "" {
		writeLine(&b, "## Description")
		writeLine(&b, "")
		writeLine(&b, desc)
		writeLine(&b, "")
	}

	writeLine(&b, "## Boundaries")
	writeLine(&b, "")
	if len(doc.Boundaries) == 0 {
		writeLine(&b, "_No boundaries declared._")
		writeLine(&b, "")
	} else {
		for i, entry := range doc.Boundaries {
			if i > 0 {
				writeLine(&b, "")
			}
			renderBoundaryEntry(&b, entry)
		}
	}

	// Canonical generated Markdown ends with exactly one LF, without a
	// trailing blank line. This keeps generated files compatible with
	// git diff --check while preserving deterministic output.
	out := strings.TrimRight(b.String(), "\r\n") + "\n"
	return []byte(out)
}

func renderBoundaryEntry(b *strings.Builder, entry BoundaryLedgerEntry) {
	id := strings.TrimSpace(entry.ID)
	if id == "" {
		id = "(unnamed)"
	}
	writeLine(b, "### "+id)
	writeLine(b, "")

	for _, cat := range requiredLedgerCategories {
		title := categoryTitle[cat]
		writeLine(b, "#### "+title)
		writeLine(b, "")
		switch cat {
		case "purpose":
			writeProse(b, entry.Purpose)
		case "owner":
			writeProse(b, entry.Owner)
		case "producers":
			writeList(b, entry.Producers)
		case "consumers":
			writeList(b, entry.Consumers)
		case "allowed_data":
			writeList(b, entry.AllowedData)
		case "prohibited_data":
			writeList(b, entry.ProhibitedData)
		case "authorization":
			writeProse(b, entry.Authorization)
		case "audit":
			writeProse(b, entry.Audit)
		case "observability":
			writeProse(b, entry.Observability)
		case "failure_behavior":
			writeProse(b, entry.FailureBehavior)
		case "versioning":
			writeProse(b, entry.Versioning)
		case "replacement_path":
			writeProse(b, entry.ReplacementPath)
		case "migration_path":
			writeProse(b, entry.MigrationPath)
		case "reassessment_trigger":
			writeProse(b, entry.ReassessmentTrigger)
		}
		writeLine(b, "")
	}
}

func writeProse(b *strings.Builder, text string) {
	prose := normalizeProse(text)
	if prose == "" {
		writeLine(b, "_Not specified._")
		return
	}
	writeLine(b, prose)
}

func writeList(b *strings.Builder, items []string) {
	wrote := false
	for _, item := range items {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}
		writeLine(b, "- "+item)
		wrote = true
	}
	if !wrote {
		writeLine(b, "_None._")
	}
}

func writeLine(b *strings.Builder, line string) {
	b.WriteString(line)
	b.WriteByte('\n')
}

// normalizeProse trims and collapses internal whitespace so folded YAML
// scalars render as a single stable paragraph.
func normalizeProse(text string) string {
	fields := strings.Fields(strings.TrimSpace(text))
	if len(fields) == 0 {
		return ""
	}
	return strings.Join(fields, " ")
}

func escapeTableCell(s string) string {
	s = strings.ReplaceAll(s, "|", `\|`)
	s = strings.ReplaceAll(s, "\n", " ")
	return s
}
