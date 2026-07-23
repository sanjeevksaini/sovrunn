package apiconform

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

// TestBoundaryLedgerMarkdownSynchronized enforces D-12 / F12-LEDGER-001:
// docs/api/boundary-ledger.yaml is the sole source of truth and
// docs/api/BOUNDARY_LEDGER.md must be its byte-identical regenerable derivative.
func TestBoundaryLedgerMarkdownSynchronized(t *testing.T) {
	t.Parallel()

	root := moduleRoot(t)
	yamlPath := filepath.Join(root, BoundaryLedgerPath)
	mdPath := filepath.Join(root, BoundaryLedgerMarkdownPath)

	raw, err := os.ReadFile(yamlPath)
	if err != nil {
		t.Fatalf("read %s: %v", BoundaryLedgerPath, err)
	}
	want, err := GenerateBoundaryLedgerMarkdown(raw)
	if err != nil {
		t.Fatalf("GenerateBoundaryLedgerMarkdown: %v", err)
	}
	got, err := os.ReadFile(mdPath)
	if err != nil {
		t.Fatalf("read %s: %v (run with SOVRUNN_WRITE_BOUNDARY_LEDGER=1 to regenerate)", BoundaryLedgerMarkdownPath, err)
	}
	if err := checkBoundaryLedgerMarkdownSync(want, got); err != nil {
		t.Fatalf("%s is stale relative to %s: %v\nRegenerate with: SOVRUNN_WRITE_BOUNDARY_LEDGER=1 go test ./internal/apiconform -run TestWriteBoundaryLedgerMarkdown",
			BoundaryLedgerMarkdownPath, BoundaryLedgerPath, err)
	}
}

// TestBoundaryLedgerMarkdownStaleDetection proves the sync checker fails when
// the derivative Markdown does not match the generator output byte-for-byte.
func TestBoundaryLedgerMarkdownStaleDetection(t *testing.T) {
	t.Parallel()

	root := moduleRoot(t)
	raw, err := os.ReadFile(filepath.Join(root, BoundaryLedgerPath))
	if err != nil {
		t.Fatalf("read %s: %v", BoundaryLedgerPath, err)
	}
	want, err := GenerateBoundaryLedgerMarkdown(raw)
	if err != nil {
		t.Fatalf("GenerateBoundaryLedgerMarkdown: %v", err)
	}
	if err := checkBoundaryLedgerMarkdownSync(want, want); err != nil {
		t.Fatalf("identical bytes must pass sync check: %v", err)
	}

	stale := append([]byte(nil), want...)
	stale = append(stale, []byte("\n<!-- deliberately stale -->\n")...)
	if err := checkBoundaryLedgerMarkdownSync(want, stale); err == nil {
		t.Fatal("deliberately stale Markdown must fail sync check")
	}

	truncated := append([]byte(nil), want...)
	if len(truncated) == 0 {
		t.Fatal("generated Markdown is empty")
	}
	truncated[0] ^= 0x01
	if err := checkBoundaryLedgerMarkdownSync(want, truncated); err == nil {
		t.Fatal("byte-mutated Markdown must fail sync check")
	}

	if err := checkBoundaryLedgerMarkdownSync(want, nil); err == nil {
		t.Fatal("nil/empty Markdown must fail sync check when expected is non-empty")
	}
}

// TestWriteBoundaryLedgerMarkdown regenerates docs/api/BOUNDARY_LEDGER.md from
// the YAML source. Opt-in via SOVRUNN_WRITE_BOUNDARY_LEDGER=1 so normal test
// runs stay read-only.
func TestWriteBoundaryLedgerMarkdown(t *testing.T) {
	if os.Getenv("SOVRUNN_WRITE_BOUNDARY_LEDGER") != "1" {
		t.Skip("set SOVRUNN_WRITE_BOUNDARY_LEDGER=1 to regenerate " + BoundaryLedgerMarkdownPath)
	}

	root := moduleRoot(t)
	raw, err := os.ReadFile(filepath.Join(root, BoundaryLedgerPath))
	if err != nil {
		t.Fatalf("read %s: %v", BoundaryLedgerPath, err)
	}
	md, err := GenerateBoundaryLedgerMarkdown(raw)
	if err != nil {
		t.Fatalf("GenerateBoundaryLedgerMarkdown: %v", err)
	}
	out := filepath.Join(root, BoundaryLedgerMarkdownPath)
	if err := os.WriteFile(out, md, 0o644); err != nil {
		t.Fatalf("write %s: %v", BoundaryLedgerMarkdownPath, err)
	}
	t.Logf("wrote %s (%d bytes)", BoundaryLedgerMarkdownPath, len(md))
}

func checkBoundaryLedgerMarkdownSync(want, got []byte) error {
	if bytes.Equal(want, got) {
		return nil
	}
	if len(want) != len(got) {
		return fmt.Errorf("length mismatch: generated=%d on-disk=%d", len(want), len(got))
	}
	for i := 0; i < len(want) && i < len(got); i++ {
		if want[i] != got[i] {
			return fmt.Errorf("first differing byte at offset %d: generated=0x%02x on-disk=0x%02x", i, want[i], got[i])
		}
	}
	return fmt.Errorf("byte content mismatch")
}
