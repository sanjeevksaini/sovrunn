package apiconform

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFitnessCheckLaterFeatureRuntimeAbsent(t *testing.T) {
	t.Parallel()

	findings := CheckLaterFeatureRuntimeAbsent(moduleRoot(t))
	if len(findings) != 0 {
		t.Fatalf("check 14 failed: %#v", findings)
	}
}

func TestFitnessCheckLaterFeatureRuntimeAbsentServiceTypeFails(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	pkgDir := filepath.Join(root, "internal", "apimeta")
	if err := os.MkdirAll(pkgDir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	src := `package apimeta

type PlacementEngine struct{}
`
	if err := os.WriteFile(filepath.Join(pkgDir, "engine.go"), []byte(src), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	// Minimal stubs so other grammar package scans do not fail on missing dirs.
	for _, pkg := range fitnessCoreGrammarPackages {
		dir := filepath.Join(root, "internal", pkg)
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatalf("mkdir %s: %v", pkg, err)
		}
		if pkg == "apimeta" {
			continue
		}
		if err := os.WriteFile(filepath.Join(dir, "stub.go"), []byte("package "+pkg+"\n"), 0o644); err != nil {
			t.Fatalf("write stub: %v", err)
		}
	}
	if err := os.MkdirAll(filepath.Join(root, "internal", "server"), 0o755); err != nil {
		t.Fatalf("mkdir server: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "internal", "server", "server.go"),
		[]byte("package server\n"), 0o644); err != nil {
		t.Fatalf("write server: %v", err)
	}

	findings := CheckLaterFeatureRuntimeAbsent(root)
	if !hasFitnessFinding(findings, FitnessCheckLaterFeatureRuntimeAbsent,
		"internal/apimeta/engine.go", "/PlacementEngine", CodeFitnessRuntimeServicePresent) {
		t.Fatalf("expected PlacementEngine runtime service finding, got %#v", findings)
	}
}

func TestFitnessCheckLaterFeatureRuntimeAbsentForbiddenPackageFails(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	for _, pkg := range fitnessCoreGrammarPackages {
		dir := filepath.Join(root, "internal", pkg)
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatalf("mkdir %s: %v", pkg, err)
		}
		if err := os.WriteFile(filepath.Join(dir, "stub.go"), []byte("package "+pkg+"\n"), 0o644); err != nil {
			t.Fatalf("write stub: %v", err)
		}
	}
	if err := os.MkdirAll(filepath.Join(root, "internal", "policy"), 0o755); err != nil {
		t.Fatalf("mkdir policy: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(root, "internal", "server"), 0o755); err != nil {
		t.Fatalf("mkdir server: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "internal", "server", "server.go"),
		[]byte("package server\n"), 0o644); err != nil {
		t.Fatalf("write server: %v", err)
	}

	findings := CheckLaterFeatureRuntimeAbsent(root)
	if !hasFitnessFinding(findings, FitnessCheckLaterFeatureRuntimeAbsent,
		"internal/policy", "/", CodeFitnessForbiddenRuntimePackage) {
		t.Fatalf("expected forbidden policy package finding, got %#v", findings)
	}
}

func TestFitnessCheckLaterFeatureRuntimeAbsentAPIsRouteFails(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	for _, pkg := range fitnessCoreGrammarPackages {
		dir := filepath.Join(root, "internal", pkg)
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatalf("mkdir %s: %v", pkg, err)
		}
		if err := os.WriteFile(filepath.Join(dir, "stub.go"), []byte("package "+pkg+"\n"), 0o644); err != nil {
			t.Fatalf("write stub: %v", err)
		}
	}
	if err := os.MkdirAll(filepath.Join(root, "internal", "server"), 0o755); err != nil {
		t.Fatalf("mkdir server: %v", err)
	}
	serverSrc := `package server

import "net/http"

func register(mux *http.ServeMux) {
	mux.Handle("/apis/core.sovrunn.io/v1alpha1/projects", http.HandlerFunc(nil))
}
`
	if err := os.WriteFile(filepath.Join(root, "internal", "server", "server.go"),
		[]byte(serverSrc), 0o644); err != nil {
		t.Fatalf("write server: %v", err)
	}

	findings := CheckLaterFeatureRuntimeAbsent(root)
	if !hasFitnessFinding(findings, FitnessCheckLaterFeatureRuntimeAbsent,
		"internal/server/server.go", "/", CodeFitnessRuntimeRoutePresent) {
		t.Fatalf("expected /apis/ route finding, got %#v", findings)
	}
}

func TestFitnessCheckLaterFeatureRuntimeAbsentGrammarRouteFails(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	for _, pkg := range fitnessCoreGrammarPackages {
		dir := filepath.Join(root, "internal", pkg)
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatalf("mkdir %s: %v", pkg, err)
		}
		src := "package " + pkg + "\n"
		if pkg == "apiconform" {
			src = `package apiconform

import "net/http"

func register(mux *http.ServeMux) {
	mux.HandleFunc("/v1/evil", func(http.ResponseWriter, *http.Request) {})
}
`
		}
		if err := os.WriteFile(filepath.Join(dir, "stub.go"), []byte(src), 0o644); err != nil {
			t.Fatalf("write stub: %v", err)
		}
	}
	if err := os.MkdirAll(filepath.Join(root, "internal", "server"), 0o755); err != nil {
		t.Fatalf("mkdir server: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "internal", "server", "server.go"),
		[]byte("package server\n"), 0o644); err != nil {
		t.Fatalf("write server: %v", err)
	}

	findings := CheckLaterFeatureRuntimeAbsent(root)
	if !hasFitnessFinding(findings, FitnessCheckLaterFeatureRuntimeAbsent,
		"internal/apiconform/stub.go", "/", CodeFitnessRuntimeRoutePresent) {
		t.Fatalf("expected grammar route registration finding, got %#v", findings)
	}
}

func TestFitnessCheckExceptionsRequireApprovedHandoff(t *testing.T) {
	t.Parallel()

	findings := CheckExceptionsRequireApprovedHandoff(moduleRoot(t))
	if len(findings) != 0 {
		t.Fatalf("check 15 failed: %#v", findings)
	}
}

func TestFitnessCheckExceptionsRequireApprovedHandoffUnknownExceptionFails(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	handoffDir := filepath.Join(root, ArchitectureHandoffsDir)
	if err := os.MkdirAll(handoffDir, 0o755); err != nil {
		t.Fatalf("mkdir handoffs: %v", err)
	}
	for _, id := range requiredFeature0012Handoffs {
		body := "---\nstatus: Approved\n---\n\n# " + id + "\n\n- Approval status: Approved\n"
		if err := os.WriteFile(filepath.Join(handoffDir, id+"-test.md"), []byte(body), 0o644); err != nil {
			t.Fatalf("write handoff: %v", err)
		}
	}
	featureDir := filepath.Join(root, "docs", "features")
	if err := os.MkdirAll(featureDir, 0o755); err != nil {
		t.Fatalf("mkdir features: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, Feature0012FeatureDocPath),
		[]byte("controlling_handoff: ADH-2026-012\n"), 0o644); err != nil {
		t.Fatalf("write feature doc: %v", err)
	}
	reportDir := filepath.Join(root, "docs", "api")
	if err := os.MkdirAll(reportDir, 0o755); err != nil {
		t.Fatalf("mkdir api docs: %v", err)
	}
	report := "# Phase 1 Compatibility Report\n\n| EX-UNAPPROVED-FOO | silent divergence |\n"
	if err := os.WriteFile(filepath.Join(root, Phase1CompatibilityReportPath), []byte(report), 0o644); err != nil {
		t.Fatalf("write report: %v", err)
	}

	findings := CheckExceptionsRequireApprovedHandoff(root)
	if !hasFitnessFinding(findings, FitnessCheckExceptionsRequireApprovedHandoff,
		Phase1CompatibilityReportPath, "/EX-UNAPPROVED-FOO", CodeFitnessExceptionWithoutHandoff) {
		t.Fatalf("expected unapproved exception finding, got %#v", findings)
	}
}

func TestFitnessCheckExceptionsRequireApprovedHandoffMissingRequiredFails(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	handoffDir := filepath.Join(root, ArchitectureHandoffsDir)
	if err := os.MkdirAll(handoffDir, 0o755); err != nil {
		t.Fatalf("mkdir handoffs: %v", err)
	}
	// Only ADH-2026-012 present; ADH-2026-013 missing.
	body := "---\nstatus: Approved\n---\n\n- Approval status: Approved\n"
	if err := os.WriteFile(filepath.Join(handoffDir, "ADH-2026-012-test.md"), []byte(body), 0o644); err != nil {
		t.Fatalf("write handoff: %v", err)
	}
	featureDir := filepath.Join(root, "docs", "features")
	if err := os.MkdirAll(featureDir, 0o755); err != nil {
		t.Fatalf("mkdir features: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, Feature0012FeatureDocPath),
		[]byte("controlling_handoff: ADH-2026-012\n"), 0o644); err != nil {
		t.Fatalf("write feature doc: %v", err)
	}
	reportDir := filepath.Join(root, "docs", "api")
	if err := os.MkdirAll(reportDir, 0o755); err != nil {
		t.Fatalf("mkdir api docs: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, Phase1CompatibilityReportPath),
		[]byte("# report\n| EX-P1-ROUTE | coexistence |\n"), 0o644); err != nil {
		t.Fatalf("write report: %v", err)
	}

	findings := CheckExceptionsRequireApprovedHandoff(root)
	if !hasFitnessFinding(findings, FitnessCheckExceptionsRequireApprovedHandoff,
		ArchitectureHandoffsDir, "/ADH-2026-013", CodeFitnessHandoffMissing) {
		t.Fatalf("expected missing ADH-2026-013 finding, got %#v", findings)
	}
}

func TestFitnessRuntimeChecksInventory(t *testing.T) {
	t.Parallel()

	want := []string{
		FitnessCheckLaterFeatureRuntimeAbsent,
		FitnessCheckExceptionsRequireApprovedHandoff,
	}
	for _, id := range want {
		if id == "" {
			t.Fatal("empty fitness check id")
		}
	}
	if ArchitectureHandoffsDir == "" || Feature0012FeatureDocPath == "" {
		t.Fatal("path constants must be set")
	}
	if len(requiredFeature0012Handoffs) != 2 {
		t.Fatalf("requiredFeature0012Handoffs len=%d want 2", len(requiredFeature0012Handoffs))
	}
	if len(phase1CoexistenceExceptions) != 9 {
		t.Fatalf("phase1CoexistenceExceptions len=%d want 9", len(phase1CoexistenceExceptions))
	}
}

func TestParseHandoffApprovalStatus(t *testing.T) {
	t.Parallel()

	yamlFront := "---\nstatus: Approved\n---\n\n# Title\n"
	if got := parseHandoffApprovalStatus(yamlFront); got != "Approved" {
		t.Fatalf("yaml front status=%q want Approved", got)
	}
	mdMeta := "# Title\n\n- Approval status: Approved\n"
	if got := parseHandoffApprovalStatus(mdMeta); got != "Approved" {
		t.Fatalf("md status=%q want Approved", got)
	}
	proposed := "---\nstatus: Proposed\n---\n"
	if got := parseHandoffApprovalStatus(proposed); got != "Proposed" {
		t.Fatalf("proposed status=%q want Proposed", got)
	}
}
