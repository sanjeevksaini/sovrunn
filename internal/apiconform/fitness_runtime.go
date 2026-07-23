package apiconform

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Fitness check IDs for F12-VERIFY-001 checks implemented in this file
// (task 16.4: checks 14, 15). Aggregation across 1–15 is task 16.5.
const (
	FitnessCheckLaterFeatureRuntimeAbsent        = "14"
	FitnessCheckExceptionsRequireApprovedHandoff = "15"
)

// Stable fitness finding codes for runtime-absence / exception-governance checks.
const (
	CodeFitnessRuntimeServicePresent   = "FITNESS_RUNTIME_SERVICE_PRESENT"
	CodeFitnessRuntimeRoutePresent     = "FITNESS_RUNTIME_ROUTE_PRESENT"
	CodeFitnessForbiddenRuntimeImport  = "FITNESS_FORBIDDEN_RUNTIME_IMPORT"
	CodeFitnessForbiddenRuntimePackage = "FITNESS_FORBIDDEN_RUNTIME_PACKAGE"
	CodeFitnessHandoffMissing          = "FITNESS_HANDOFF_MISSING"
	CodeFitnessHandoffNotApproved      = "FITNESS_HANDOFF_NOT_APPROVED"
	CodeFitnessExceptionWithoutHandoff = "FITNESS_EXCEPTION_WITHOUT_HANDOFF"
	CodeFitnessTraceabilityMissing     = "FITNESS_TRACEABILITY_MISSING"
)

// ArchitectureHandoffsDir is the repository-relative directory for Architecture
// Decision Handoff records (ADH-YYYY-NNN-*.md).
const ArchitectureHandoffsDir = "docs/reviews/architecture-decision-handoffs"

// Feature0012FeatureDocPath is the repository-relative FEATURE-0012 feature file
// used for controlling-handoff traceability (F12-TRACE-001 / check 15).
const Feature0012FeatureDocPath = "docs/features/FEATURE-0012-api-resource-naming-status-and-validation-standard.md"

// requiredFeature0012Handoffs are the controlling ADHs for FEATURE-0012
// (ADH-2026-012 Extend approval; ADH-2026-013 Operation scopes clarification).
var requiredFeature0012Handoffs = []string{
	"ADH-2026-012",
	"ADH-2026-013",
}

// phase1CoexistenceExceptions maps documented Phase 1 EX-P1-* exception IDs in
// docs/api/PHASE1_COMPATIBILITY_REPORT.md to the approved controlling handoff
// that authorized coexistence without rewrite (ADH-2026-012 / F12-COMPAT-003).
var phase1CoexistenceExceptions = map[string]string{
	"EX-P1-ROUTE":          "ADH-2026-012",
	"EX-P1-ERROR-ENVELOPE": "ADH-2026-012",
	"EX-P1-FIELD-PATH":     "ADH-2026-012",
	"EX-P1-METADATA":       "ADH-2026-012",
	"EX-P1-STATUS":         "ADH-2026-012",
	"EX-P1-REFS":           "ADH-2026-012",
	"EX-P1-LIST":           "ADH-2026-012",
	"EX-P1-CONCURRENCY":    "ADH-2026-012",
	"EX-P1-DECODE":         "ADH-2026-012",
}

// forbiddenLaterFeatureInternalPackages are top-level internal/ directory names
// that would constitute later-feature runtime services (F12-IMPL-002). Phase 1
// packages (api/server/resources/…) are intentionally not listed.
var forbiddenLaterFeatureInternalPackages = map[string]struct{}{
	"provider":        {},
	"providers":       {},
	"policy":          {},
	"policies":        {},
	"placement":       {},
	"provisioning":    {},
	"provisioner":     {},
	"pluginruntime":   {},
	"pluginexec":      {},
	"pluginexecution": {},
	"auditservice":    {},
	"auditsvc":        {},
	"auditengine":     {},
}

// runtimeServiceTypeName matches type names that look like later-feature
// execution services inside grammar packages. Contract types such as
// PluginDefinition / AuditEvent / PlacementEvaluationRequest do not match.
var runtimeServiceTypeName = regexp.MustCompile(
	`(?i)^(Provider|Plugin|Policy|Placement|Audit|Provisioning).*(Service|Engine|Executor|Runner|Controller)$|` +
		`(?i).*(PolicyEngine|PlacementEngine|PluginExecutor|Provisioner)$`,
)

var exceptionIDPattern = regexp.MustCompile(`\bEX-[A-Z0-9]+(?:-[A-Z0-9]+)*\b`)
var adhIDPattern = regexp.MustCompile(`\bADH-\d{4}-\d{3}\b`)
var adhFilenamePattern = regexp.MustCompile(`^(ADH-\d{4}-\d{3})(?:-.*)?\.md$`)

// CheckLaterFeatureRuntimeAbsent implements F12-VERIFY-001 check 14:
// later-feature runtime behavior is absent (F12-IMPL-002). Grammar packages
// must not host provider/plugin/policy/placement/audit/provisioning services,
// must not import internal/api or internal/server, must not register HTTP
// routes, and FEATURE-0012 must not add /apis/... runtime routes.
//
// This check is gate/test-only. It does not execute providers, plugins,
// policy engines, placement, audit emitters, or provisioning. Request bodies,
// secrets, and credentials are never logged.
func CheckLaterFeatureRuntimeAbsent(moduleRoot string) []FitnessFinding {
	var findings []FitnessFinding

	findings = append(findings, checkForbiddenLaterFeaturePackages(moduleRoot)...)
	findings = append(findings, checkGrammarForbiddenRuntimeImports(moduleRoot)...)
	findings = append(findings, checkGrammarRuntimeServiceTypes(moduleRoot)...)
	findings = append(findings, checkGrammarHTTPRouteRegistration(moduleRoot)...)
	findings = append(findings, checkServerNoFeature0012APIsRoutes(moduleRoot)...)

	return sortFindings(findings)
}

// CheckExceptionsRequireApprovedHandoff implements F12-VERIFY-001 check 15:
// exceptions require an approved architecture handoff (F12-SCOPESTD-004,
// F12-COMPAT-003, F12-TRACE-001). Controlling FEATURE-0012 ADHs must exist and
// be Approved; every EX-* ID in the Phase 1 compatibility report must map to
// an approved ADH (Phase 1 coexistence EX-P1-* → ADH-2026-012).
func CheckExceptionsRequireApprovedHandoff(moduleRoot string) []FitnessFinding {
	var findings []FitnessFinding

	handoffs, err := loadApprovedHandoffs(moduleRoot)
	if err != nil {
		return sortFindings([]FitnessFinding{{
			Check:   FitnessCheckExceptionsRequireApprovedHandoff,
			Schema:  ArchitectureHandoffsDir,
			Path:    "/",
			Code:    CodeFitnessHandoffMissing,
			Message: err.Error(),
		}})
	}

	for _, id := range requiredFeature0012Handoffs {
		status, ok := handoffs[id]
		if !ok {
			findings = append(findings, FitnessFinding{
				Check:   FitnessCheckExceptionsRequireApprovedHandoff,
				Schema:  ArchitectureHandoffsDir,
				Path:    "/" + id,
				Code:    CodeFitnessHandoffMissing,
				Message: fmt.Sprintf("required FEATURE-0012 controlling handoff %s is missing", id),
			})
			continue
		}
		if !isApprovedHandoffStatus(status) {
			findings = append(findings, FitnessFinding{
				Check:   FitnessCheckExceptionsRequireApprovedHandoff,
				Schema:  ArchitectureHandoffsDir,
				Path:    "/" + id,
				Code:    CodeFitnessHandoffNotApproved,
				Message: fmt.Sprintf("required handoff %s status=%q want Approved", id, status),
			})
		}
	}

	findings = append(findings, checkFeatureDocHandoffTraceability(moduleRoot)...)
	findings = append(findings, checkCompatibilityExceptionsHaveHandoffs(moduleRoot, handoffs)...)

	return sortFindings(findings)
}

func checkForbiddenLaterFeaturePackages(moduleRoot string) []FitnessFinding {
	internalDir := filepath.Join(moduleRoot, "internal")
	entries, err := os.ReadDir(internalDir)
	if err != nil {
		return []FitnessFinding{{
			Check:   FitnessCheckLaterFeatureRuntimeAbsent,
			Schema:  "internal",
			Path:    "/",
			Code:    CodeFitnessForbiddenRuntimePackage,
			Message: err.Error(),
		}}
	}
	var findings []FitnessFinding
	for _, ent := range entries {
		if !ent.IsDir() {
			continue
		}
		name := strings.ToLower(ent.Name())
		if _, banned := forbiddenLaterFeatureInternalPackages[name]; banned {
			findings = append(findings, FitnessFinding{
				Check:   FitnessCheckLaterFeatureRuntimeAbsent,
				Schema:  "internal/" + ent.Name(),
				Path:    "/",
				Code:    CodeFitnessForbiddenRuntimePackage,
				Message: fmt.Sprintf("later-feature runtime package %q is prohibited under F12-IMPL-002", ent.Name()),
			})
		}
	}
	return findings
}

func checkGrammarForbiddenRuntimeImports(moduleRoot string) []FitnessFinding {
	var findings []FitnessFinding
	for _, pkg := range fitnessCoreGrammarPackages {
		dir := filepath.Join(moduleRoot, "internal", pkg)
		imports, err := listGoImportPaths(dir)
		if err != nil {
			findings = append(findings, FitnessFinding{
				Check:   FitnessCheckLaterFeatureRuntimeAbsent,
				Schema:  "internal/" + pkg,
				Path:    "/",
				Code:    CodeFitnessSchemaLoadFailed,
				Message: err.Error(),
			})
			continue
		}
		for _, imp := range imports {
			if isFitnessForbiddenRuntimeImport(imp) {
				findings = append(findings, FitnessFinding{
					Check:   FitnessCheckLaterFeatureRuntimeAbsent,
					Schema:  "internal/" + pkg,
					Path:    "/",
					Code:    CodeFitnessForbiddenRuntimeImport,
					Message: fmt.Sprintf("grammar package imports forbidden runtime path %q", imp),
				})
			}
		}
	}
	return findings
}

func checkGrammarRuntimeServiceTypes(moduleRoot string) []FitnessFinding {
	var findings []FitnessFinding
	fset := token.NewFileSet()
	for _, pkg := range fitnessCoreGrammarPackages {
		dir := filepath.Join(moduleRoot, "internal", pkg)
		entries, err := os.ReadDir(dir)
		if err != nil {
			findings = append(findings, FitnessFinding{
				Check:   FitnessCheckLaterFeatureRuntimeAbsent,
				Schema:  "internal/" + pkg,
				Path:    "/",
				Code:    CodeFitnessSchemaLoadFailed,
				Message: err.Error(),
			})
			continue
		}
		for _, ent := range entries {
			if ent.IsDir() || !strings.HasSuffix(ent.Name(), ".go") || strings.HasSuffix(ent.Name(), "_test.go") {
				continue
			}
			path := filepath.Join(dir, ent.Name())
			file, err := parser.ParseFile(fset, path, nil, 0)
			if err != nil {
				findings = append(findings, FitnessFinding{
					Check:   FitnessCheckLaterFeatureRuntimeAbsent,
					Schema:  "internal/" + pkg + "/" + ent.Name(),
					Path:    "/",
					Code:    CodeFitnessSchemaLoadFailed,
					Message: err.Error(),
				})
				continue
			}
			for _, decl := range file.Decls {
				gen, ok := decl.(*ast.GenDecl)
				if !ok || gen.Tok != token.TYPE {
					continue
				}
				for _, spec := range gen.Specs {
					ts, ok := spec.(*ast.TypeSpec)
					if !ok || ts.Name == nil {
						continue
					}
					name := ts.Name.Name
					if runtimeServiceTypeName.MatchString(name) {
						findings = append(findings, FitnessFinding{
							Check:   FitnessCheckLaterFeatureRuntimeAbsent,
							Schema:  "internal/" + pkg + "/" + ent.Name(),
							Path:    "/" + name,
							Code:    CodeFitnessRuntimeServicePresent,
							Message: fmt.Sprintf("later-feature runtime service type %q is prohibited in grammar packages", name),
						})
					}
				}
			}
		}
	}
	return findings
}

func checkGrammarHTTPRouteRegistration(moduleRoot string) []FitnessFinding {
	var findings []FitnessFinding
	fset := token.NewFileSet()
	for _, pkg := range fitnessCoreGrammarPackages {
		dir := filepath.Join(moduleRoot, "internal", pkg)
		entries, err := os.ReadDir(dir)
		if err != nil {
			findings = append(findings, FitnessFinding{
				Check:   FitnessCheckLaterFeatureRuntimeAbsent,
				Schema:  "internal/" + pkg,
				Path:    "/",
				Code:    CodeFitnessSchemaLoadFailed,
				Message: err.Error(),
			})
			continue
		}
		for _, ent := range entries {
			if ent.IsDir() || !strings.HasSuffix(ent.Name(), ".go") || strings.HasSuffix(ent.Name(), "_test.go") {
				continue
			}
			path := filepath.Join(dir, ent.Name())
			file, err := parser.ParseFile(fset, path, nil, 0)
			if err != nil {
				findings = append(findings, FitnessFinding{
					Check:   FitnessCheckLaterFeatureRuntimeAbsent,
					Schema:  "internal/" + pkg + "/" + ent.Name(),
					Path:    "/",
					Code:    CodeFitnessSchemaLoadFailed,
					Message: err.Error(),
				})
				continue
			}
			ast.Inspect(file, func(n ast.Node) bool {
				call, ok := n.(*ast.CallExpr)
				if !ok {
					return true
				}
				sel, ok := call.Fun.(*ast.SelectorExpr)
				if !ok || sel.Sel == nil {
					return true
				}
				switch sel.Sel.Name {
				case "Handle", "HandleFunc", "ListenAndServe", "ListenAndServeTLS":
					findings = append(findings, FitnessFinding{
						Check:   FitnessCheckLaterFeatureRuntimeAbsent,
						Schema:  "internal/" + pkg + "/" + ent.Name(),
						Path:    "/",
						Code:    CodeFitnessRuntimeRoutePresent,
						Message: fmt.Sprintf("grammar package must not register runtime HTTP routes (%s)", sel.Sel.Name),
					})
				}
				return true
			})
		}
	}
	return findings
}

func checkServerNoFeature0012APIsRoutes(moduleRoot string) []FitnessFinding {
	serverPath := filepath.Join(moduleRoot, "internal", "server", "server.go")
	raw, err := os.ReadFile(serverPath)
	if err != nil {
		return []FitnessFinding{{
			Check:   FitnessCheckLaterFeatureRuntimeAbsent,
			Schema:  "internal/server/server.go",
			Path:    "/",
			Code:    CodeFitnessSchemaLoadFailed,
			Message: err.Error(),
		}}
	}
	if strings.Contains(string(raw), "/apis/") {
		return []FitnessFinding{{
			Check:   FitnessCheckLaterFeatureRuntimeAbsent,
			Schema:  "internal/server/server.go",
			Path:    "/",
			Code:    CodeFitnessRuntimeRoutePresent,
			Message: "FEATURE-0012 must not register /apis/... runtime HTTP routes (grammar/conformance only)",
		}}
	}
	return nil
}

func checkFeatureDocHandoffTraceability(moduleRoot string) []FitnessFinding {
	path := filepath.Join(moduleRoot, Feature0012FeatureDocPath)
	raw, err := os.ReadFile(path)
	if err != nil {
		return []FitnessFinding{{
			Check:   FitnessCheckExceptionsRequireApprovedHandoff,
			Schema:  Feature0012FeatureDocPath,
			Path:    "/",
			Code:    CodeFitnessTraceabilityMissing,
			Message: err.Error(),
		}}
	}
	text := string(raw)
	var findings []FitnessFinding
	if !strings.Contains(text, "ADH-2026-012") {
		findings = append(findings, FitnessFinding{
			Check:   FitnessCheckExceptionsRequireApprovedHandoff,
			Schema:  Feature0012FeatureDocPath,
			Path:    "/controlling_handoff",
			Code:    CodeFitnessTraceabilityMissing,
			Message: "FEATURE-0012 feature doc must cite controlling handoff ADH-2026-012",
		})
	}
	return findings
}

func checkCompatibilityExceptionsHaveHandoffs(moduleRoot string, handoffs map[string]string) []FitnessFinding {
	reportPath := filepath.Join(moduleRoot, Phase1CompatibilityReportPath)
	raw, err := os.ReadFile(reportPath)
	if err != nil {
		return []FitnessFinding{{
			Check:   FitnessCheckExceptionsRequireApprovedHandoff,
			Schema:  Phase1CompatibilityReportPath,
			Path:    "/",
			Code:    CodeFitnessSchemaLoadFailed,
			Message: err.Error(),
		}}
	}

	seen := map[string]struct{}{}
	var findings []FitnessFinding
	report := string(raw)
	for _, loc := range exceptionIDPattern.FindAllStringIndex(report, -1) {
		id := report[loc[0]:loc[1]]
		// Skip wildcard shorthand such as "EX-P1" inside "EX-P1-*".
		rest := report[loc[1]:]
		if strings.HasPrefix(rest, "-*") || strings.HasPrefix(rest, "*") {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}

		if approving, ok := phase1CoexistenceExceptions[id]; ok {
			status, present := handoffs[approving]
			if !present {
				findings = append(findings, FitnessFinding{
					Check:   FitnessCheckExceptionsRequireApprovedHandoff,
					Schema:  Phase1CompatibilityReportPath,
					Path:    "/" + id,
					Code:    CodeFitnessHandoffMissing,
					Message: fmt.Sprintf("exception %s requires approving handoff %s which is missing", id, approving),
				})
				continue
			}
			if !isApprovedHandoffStatus(status) {
				findings = append(findings, FitnessFinding{
					Check:   FitnessCheckExceptionsRequireApprovedHandoff,
					Schema:  Phase1CompatibilityReportPath,
					Path:    "/" + id,
					Code:    CodeFitnessHandoffNotApproved,
					Message: fmt.Sprintf("exception %s approving handoff %s is not Approved", id, approving),
				})
			}
			continue
		}

		// Unknown EX-* IDs must cite an Approved ADH somewhere in the report.
		cited := findCitedApprovedHandoffForException(report, id, handoffs)
		if cited == "" {
			findings = append(findings, FitnessFinding{
				Check:   FitnessCheckExceptionsRequireApprovedHandoff,
				Schema:  Phase1CompatibilityReportPath,
				Path:    "/" + id,
				Code:    CodeFitnessExceptionWithoutHandoff,
				Message: fmt.Sprintf("exception %s requires an approved architecture handoff (F12-SCOPESTD-004)", id),
			})
		}
	}
	return findings
}

// findCitedApprovedHandoffForException looks for ADH IDs on the same line as
// the exception ID (table row / inline citation). Returns the first Approved
// ADH ID found, or empty if none.
func findCitedApprovedHandoffForException(report, exceptionID string, handoffs map[string]string) string {
	for _, line := range strings.Split(report, "\n") {
		if !strings.Contains(line, exceptionID) {
			continue
		}
		for _, adh := range adhIDPattern.FindAllString(line, -1) {
			if isApprovedHandoffStatus(handoffs[adh]) {
				return adh
			}
		}
	}
	return ""
}

func loadApprovedHandoffs(moduleRoot string) (map[string]string, error) {
	dir := filepath.Join(moduleRoot, ArchitectureHandoffsDir)
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	out := make(map[string]string)
	for _, ent := range entries {
		if ent.IsDir() {
			continue
		}
		name := ent.Name()
		m := adhFilenamePattern.FindStringSubmatch(name)
		if m == nil {
			continue
		}
		if strings.Contains(strings.ToUpper(name), "EXAMPLE") {
			continue
		}
		id := m[1]
		raw, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			return nil, err
		}
		out[id] = parseHandoffApprovalStatus(string(raw))
	}
	return out, nil
}

func parseHandoffApprovalStatus(content string) string {
	// Prefer YAML front matter status when present.
	if strings.HasPrefix(strings.TrimSpace(content), "---") {
		parts := strings.SplitN(content, "---", 3)
		if len(parts) >= 3 {
			for _, line := range strings.Split(parts[1], "\n") {
				line = strings.TrimSpace(line)
				lower := strings.ToLower(line)
				if strings.HasPrefix(lower, "status:") {
					_, after, ok := strings.Cut(line, ":")
					if ok {
						return strings.TrimSpace(after)
					}
				}
			}
		}
	}
	for _, line := range strings.Split(content, "\n") {
		trimmed := strings.TrimSpace(line)
		lower := strings.ToLower(trimmed)
		if strings.Contains(lower, "approval status:") {
			_, after, ok := strings.Cut(trimmed, ":")
			if ok {
				return strings.TrimSpace(after)
			}
		}
	}
	return ""
}

func isApprovedHandoffStatus(status string) bool {
	return strings.EqualFold(strings.TrimSpace(status), "Approved")
}

func isFitnessForbiddenRuntimeImport(importPath string) bool {
	const prefix = "github.com/sanjeevksaini/sovrunn/internal/"
	if strings.HasPrefix(importPath, prefix) {
		rest := strings.TrimPrefix(importPath, prefix)
		name, _, _ := strings.Cut(rest, "/")
		return name == "api" || name == "server"
	}
	return strings.Contains(importPath, "/internal/api") ||
		strings.Contains(importPath, "/internal/server") ||
		importPath == "internal/api" ||
		importPath == "internal/server"
}
