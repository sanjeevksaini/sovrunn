package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const (
	rule001 = "F18-ARCH-001_UNAUTHORIZED_COMMIT_CALLER"
	rule002 = "F18-ARCH-002_ROLEASSIGNMENT_WRITER"
	rule003 = "F18-ARCH-003_EVIDENCE_IMPORT"
	rule004 = "F18-ARCH-004_EDITOR_ESCAPE"
	rule005 = "F18-ARCH-005_PACKAGE_DIRECTION"
	rule006 = "F18-ARCH-006_TRIGGER_CALLER"
	rule007 = "F18-ARCH-007_POLICY_SEAM"
	rule008 = "F18-ARCH-008_EVIDENCE_PUBLISHER"
	rule009 = "F18-ARCH-009_PUBLICATION_FINALIZATION"
	rule010 = "F18-ARCH-010_APPLIED_CHANGE_RECEIPT"
	rule011 = "F18-ARCH-011_AUTHORIZATION_CURRENTNESS"
)

type diagnostic struct {
	RuleID  string
	File    string
	Line    int
	Message string
}

func main() {
	target := flag.String("path", "internal/govaccess", "path to check")
	flag.Parse()

	diags, err := run(*target)
	if err != nil {
		fmt.Fprintf(os.Stderr, "feature0018-architecture-check: %v\n", err)
		os.Exit(2)
	}
	if len(diags) == 0 {
		fmt.Println("feature0018-architecture-check: ok")
		return
	}
	for _, d := range diags {
		fmt.Printf("%s %s:%d %s\n", d.RuleID, d.File, d.Line, d.Message)
	}
	os.Exit(1)
}

func run(target string) ([]diagnostic, error) {
	abs, err := filepath.Abs(target)
	if err != nil {
		return nil, err
	}
	files, err := collectGoFiles(abs)
	if err != nil {
		return nil, err
	}
	fset := token.NewFileSet()
	diags := make([]diagnostic, 0)
	for _, file := range files {
		parsed, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
		if err != nil {
			return nil, fmt.Errorf("parse %s: %w", file, err)
		}
		info := typeCheckFile(fset, parsed)
		pkg := parsed.Name.Name
		diags = append(diags, applyRules(fset, parsed, info, pkg)...)
	}
	sortDiagnostics(diags)
	return diags, nil
}

func collectGoFiles(root string) ([]string, error) {
	out := make([]string, 0, 128)
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			base := filepath.Base(path)
			if base == ".git" || base == "vendor" || base == "node_modules" {
				return fs.SkipDir
			}
			return nil
		}
		if strings.HasSuffix(path, "_test.go") {
			return nil
		}
		if strings.HasSuffix(path, ".go") {
			out = append(out, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Strings(out)
	return out, nil
}

func applyRules(fset *token.FileSet, file *ast.File, info *types.Info, pkg string) []diagnostic {
	out := make([]diagnostic, 0, 16)
	out = append(out, checkRule001(fset, file, pkg)...)
	out = append(out, checkRule002(fset, file, pkg)...)
	out = append(out, checkRule003(fset, file, pkg)...)
	out = append(out, checkRule004(fset, file, pkg)...)
	out = append(out, checkRule005(fset, file, pkg)...)
	out = append(out, checkRule006(fset, file, info, pkg)...)
	out = append(out, checkRule007(fset, file, pkg)...)
	out = append(out, checkRule008(fset, file, pkg)...)
	out = append(out, checkRule009(fset, file, pkg)...)
	out = append(out, checkRule010(fset, file, pkg)...)
	out = append(out, checkRule011(fset, file, pkg)...)
	return out
}

func checkRule001(fset *token.FileSet, file *ast.File, pkg string) []diagnostic {
	if pkg == "uow" {
		return nil
	}
	forbidden := map[string]bool{
		"InspectOrReserveCallerMutation": true,
		"BeginControllerMutation":        true,
		"BeginEvidenceTransaction":       true,
		"Commit":                         true,
		"Seal":                           true,
	}
	return checkSelectorCalls(fset, file, forbidden, rule001)
}

func checkRule002(fset *token.FileSet, file *ast.File, pkg string) []diagnostic {
	if pkg == "roleassign" {
		return nil
	}
	out := checkCallNames(fset, file, map[string]bool{
		"newPreparedRoleAssignmentIntent":  true,
		"NewPreparedRoleAssignmentIntent":  true,
		"NewFinalizedRoleAssignmentChange": true,
	}, rule002)
	ast.Inspect(file, func(n ast.Node) bool {
		lit, ok := n.(*ast.CompositeLit)
		if !ok {
			return true
		}
		if sel, ok := lit.Type.(*ast.SelectorExpr); ok {
			if sel.Sel.Name == "PreparedRoleAssignmentIntent" || sel.Sel.Name == "FinalizedRoleAssignmentChange" {
				out = append(out, newDiag(fset, lit.Pos(), rule002, "forbidden roleassign concrete construction"))
			}
		}
		return true
	})
	return out
}

func checkRule003(fset *token.FileSet, file *ast.File, pkg string) []diagnostic {
	out := make([]diagnostic, 0, 4)
	if pkg == "evidence" {
		for _, imp := range file.Imports {
			path := strings.Trim(imp.Path.Value, `"`)
			if strings.HasSuffix(path, "/state") || strings.HasSuffix(path, "/authzeval") || strings.HasSuffix(path, "/uow") || path == "github.com/sanjeevksaini/sovrunn/internal/govaccess" {
				out = append(out, newDiag(fset, imp.Pos(), rule003, "evidence forbidden import"))
			}
		}
	}
	if pkg != "authzeval" {
		out = append(out, checkCallNames(fset, file, map[string]bool{
			"NewAuthorizationCandidateDescriptor": true,
		}, rule003)...)
	}
	if pkg != "uow" {
		out = append(out, checkCallNames(fset, file, map[string]bool{
			"RequireCurrentAllow": true,
		}, rule003)...)
	}
	return out
}

func checkRule004(fset *token.FileSet, file *ast.File, pkg string) []diagnostic {
	out := make([]diagnostic, 0, 2)
	ast.Inspect(file, func(n ast.Node) bool {
		switch v := n.(type) {
		case *ast.StructType:
			for _, field := range v.Fields.List {
				if selectorName(field.Type) == "StateEditor" {
					out = append(out, newDiag(fset, field.Pos(), rule004, "state editor must not escape in struct fields"))
				}
			}
		case *ast.FuncDecl:
			if v.Type.Params == nil {
				return true
			}
			for _, field := range v.Type.Params.List {
				if selectorName(field.Type) == "StateEditor" && v.Name.Name != "ApplyTo" {
					out = append(out, newDiag(fset, field.Pos(), rule004, "state editor accepted outside ApplyTo"))
				}
			}
		}
		return true
	})
	return out
}

func checkRule005(fset *token.FileSet, file *ast.File, pkg string) []diagnostic {
	out := make([]diagnostic, 0, 4)
	ownerPkgs := map[string]bool{
		"accessgroup": true, "membership": true, "roledefinition": true, "roleassign": true,
		"privileged": true, "review": true, "approval": true, "exception": true, "governanceprofile": true,
	}
	for _, imp := range file.Imports {
		path := strings.Trim(imp.Path.Value, `"`)
		if ownerPkgs[pkg] && (strings.HasSuffix(path, "/uow") || strings.HasSuffix(path, "/command") || strings.HasSuffix(path, "/idempotency") || path == "github.com/sanjeevksaini/sovrunn/internal/govaccess") {
			out = append(out, newDiag(fset, imp.Pos(), rule005, "owner package imports forbidden dependency"))
		}
		if pkg == "state" && (strings.HasSuffix(path, "/uow") || strings.HasSuffix(path, "/command") || strings.HasSuffix(path, "/authzeval") || ownerPkgs[pathBase(path)]) {
			out = append(out, newDiag(fset, imp.Pos(), rule005, "state package import direction violation"))
		}
		if pkg == "idempotency" && (strings.HasSuffix(path, "/state") || strings.HasSuffix(path, "/evidence")) {
			out = append(out, newDiag(fset, imp.Pos(), rule005, "idempotency import direction violation"))
		}
		if pkg == "evidence" && strings.HasSuffix(path, "/state") {
			out = append(out, newDiag(fset, imp.Pos(), rule005, "evidence import direction violation"))
		}
	}
	ast.Inspect(file, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}
		if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
			if x, ok := sel.X.(*ast.Ident); ok && x.Name == "time" && sel.Sel.Name == "Now" {
				out = append(out, newDiag(fset, call.Pos(), rule005, "time.Now is forbidden"))
			}
		}
		return true
	})
	return out
}

func checkRule006(fset *token.FileSet, file *ast.File, info *types.Info, pkg string) []diagnostic {
	out := make([]diagnostic, 0, 4)
	ast.Inspect(file, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}
		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok {
			return true
		}
		recv := receiverTypeString(info, sel)
		switch {
		case strings.HasSuffix(recv, ".GrantPort") && sel.Sel.Name == "Submit" && pkg != "grantintent":
			out = append(out, newDiag(fset, call.Pos(), rule006, "GrantPort.Submit caller must be grantintent"))
		case strings.HasSuffix(recv, ".ActivationPort") && sel.Sel.Name == "Submit" && pkg != "privileged":
			out = append(out, newDiag(fset, call.Pos(), rule006, "ActivationPort.Submit caller must be privileged"))
		case strings.HasSuffix(recv, ".ReviewRemediationPort") && sel.Sel.Name == "Submit" && pkg != "review":
			out = append(out, newDiag(fset, call.Pos(), rule006, "ReviewRemediationPort.Submit caller must be review"))
		case strings.HasSuffix(recv, ".DirectRevokePort") && sel.Sel.Name == "Prepare" && pkg != "command":
			out = append(out, newDiag(fset, call.Pos(), rule006, "DirectRevokePort.Prepare caller must be command"))
		}
		return true
	})
	out = append(out, checkCallNames(fset, file, map[string]bool{
		"NewGrantTriggerIntent":             pkg != "grantintent",
		"NewActivationTriggerIntent":        pkg != "privileged",
		"NewReviewRemediationTriggerIntent": pkg != "review",
	}, rule006)...)
	return out
}

func checkRule007(fset *token.FileSet, file *ast.File, pkg string) []diagnostic {
	out := make([]diagnostic, 0, 3)
	if pkg == "roleassign" {
		for _, imp := range file.Imports {
			path := strings.Trim(imp.Path.Value, `"`)
			if strings.HasSuffix(path, "/approvalreq") {
				out = append(out, newDiag(fset, imp.Pos(), rule007, "roleassign must not import approvalreq"))
			}
		}
	}
	if pkg != "privileged" {
		out = append(out, checkCallNames(fset, file, map[string]bool{"EvaluatePrivileged": true}, rule007)...)
	}
	return out
}

func checkRule008(fset *token.FileSet, file *ast.File, pkg string) []diagnostic {
	out := make([]diagnostic, 0, 8)
	if pkg != "uow" {
		out = append(out, checkCallNames(fset, file, map[string]bool{
			"RequireCurrentAllow":            true,
			"CompleteAuthorizationCandidate": true,
			"PrepareMutationCarrierSet":      true,
			"PrepareAuthorizationCarrierSet": true,
			"ReleaseAfterPublication":        true,
		}, rule008)...)
	}
	if pkg == "govaccess" {
		out = append(out, checkSelectorCalls(fset, file, map[string]bool{
			"InspectOrReserveCallerMutation": true,
			"BeginControllerMutation":        true,
			"BeginEvidenceTransaction":       true,
			"CoordinateCallerMutation":       true,
			"CoordinateControllerMutation":   true,
			"CoordinateEvidenceOnly":         true,
			"Commit":                         true,
			"Seal":                           true,
			"Abort":                          true,
		}, rule008)...)
	}
	return out
}

func checkRule009(fset *token.FileSet, file *ast.File, pkg string) []diagnostic {
	if pkg == "uow" {
		return nil
	}
	return checkCallNames(fset, file, map[string]bool{"FinalizeAt": true}, rule009)
}

func checkRule010(fset *token.FileSet, file *ast.File, pkg string) []diagnostic {
	out := make([]diagnostic, 0, 4)
	if pkg != "state" {
		out = append(out, checkCallNames(fset, file, map[string]bool{
			"NewAppliedChangeLink":     true,
			"NewAppliedConclusionLink": true,
		}, rule010)...)
	}
	if pkg != "evidence" {
		out = append(out, checkCallNames(fset, file, map[string]bool{
			"NewCompletionBinding": true,
		}, rule010)...)
	}
	return out
}

func checkRule011(fset *token.FileSet, file *ast.File, pkg string) []diagnostic {
	out := make([]diagnostic, 0, 3)
	if pkg != "authzeval" {
		out = append(out, checkCallNames(fset, file, map[string]bool{
			"NewAuthorizationDependencyVersionSet": true,
			"NewAuthorizationCurrentnessClaim":     true,
		}, rule011)...)
	}
	if pkg != "uow" {
		out = append(out, checkCallNames(fset, file, map[string]bool{
			"FinalizeCandidateAt": true,
		}, rule011)...)
	}
	return out
}

func checkCallNames(fset *token.FileSet, file *ast.File, names map[string]bool, ruleID string) []diagnostic {
	out := make([]diagnostic, 0, 2)
	ast.Inspect(file, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}
		switch fn := call.Fun.(type) {
		case *ast.Ident:
			if names[fn.Name] {
				out = append(out, newDiag(fset, call.Pos(), ruleID, "forbidden call "+fn.Name))
			}
		case *ast.SelectorExpr:
			if names[fn.Sel.Name] {
				out = append(out, newDiag(fset, call.Pos(), ruleID, "forbidden call "+fn.Sel.Name))
			}
		}
		return true
	})
	return out
}

func checkSelectorCalls(fset *token.FileSet, file *ast.File, names map[string]bool, ruleID string) []diagnostic {
	out := make([]diagnostic, 0, 2)
	ast.Inspect(file, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}
		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok {
			return true
		}
		if names[sel.Sel.Name] {
			out = append(out, newDiag(fset, call.Pos(), ruleID, "forbidden selector call "+sel.Sel.Name))
		}
		return true
	})
	return out
}

func typeCheckFile(fset *token.FileSet, file *ast.File) *types.Info {
	info := &types.Info{
		Types:      map[ast.Expr]types.TypeAndValue{},
		Defs:       map[*ast.Ident]types.Object{},
		Uses:       map[*ast.Ident]types.Object{},
		Selections: map[*ast.SelectorExpr]*types.Selection{},
	}
	cfg := &types.Config{
		Importer: permissiveImporter{stdlib: importer.Default()},
		Error:    func(error) {},
	}
	_, _ = cfg.Check(file.Name.Name, fset, []*ast.File{file}, info)
	return info
}

type permissiveImporter struct {
	stdlib types.Importer
}

func (p permissiveImporter) Import(path string) (*types.Package, error) {
	if isStdlibImport(path) {
		return p.stdlib.Import(path)
	}
	return types.NewPackage(path, pathBase(path)), nil
}

func isStdlibImport(path string) bool {
	first := path
	if idx := strings.Index(path, "/"); idx >= 0 {
		first = path[:idx]
	}
	return !strings.Contains(first, ".")
}

func receiverTypeString(info *types.Info, sel *ast.SelectorExpr) string {
	if info == nil {
		return ""
	}
	s := info.Selections[sel]
	if s == nil || s.Recv() == nil {
		return ""
	}
	return s.Recv().String()
}

func sortDiagnostics(diags []diagnostic) {
	sort.Slice(diags, func(i, j int) bool {
		if diags[i].RuleID != diags[j].RuleID {
			return diags[i].RuleID < diags[j].RuleID
		}
		if diags[i].File != diags[j].File {
			return diags[i].File < diags[j].File
		}
		if diags[i].Line != diags[j].Line {
			return diags[i].Line < diags[j].Line
		}
		return diags[i].Message < diags[j].Message
	})
}

func newDiag(fset *token.FileSet, pos token.Pos, ruleID, message string) diagnostic {
	p := fset.Position(pos)
	return diagnostic{
		RuleID:  ruleID,
		File:    p.Filename,
		Line:    p.Line,
		Message: message,
	}
}

func pathBase(path string) string {
	path = strings.TrimRight(path, "/")
	if idx := strings.LastIndex(path, "/"); idx >= 0 {
		return path[idx+1:]
	}
	return path
}

func selectorName(expr ast.Expr) string {
	sel, ok := expr.(*ast.SelectorExpr)
	if !ok {
		return ""
	}
	return sel.Sel.Name
}
