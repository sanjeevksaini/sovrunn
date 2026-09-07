package govaccess

import (
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"testing/quick"
	"time"

	"github.com/sanjeevksaini/sovrunn/internal/govaccess/clock"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/evidence"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/state"
)

func TestNewRootWiresDependencies(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 9, 7, 12, 0, 0, 0, time.UTC)
	sentinelSeam := struct{ Name string }{Name: "policy-seam"}
	root, err := NewRoot(clock.NewFixedClock(now), sentinelSeam)
	if err != nil {
		t.Fatalf("NewRoot() error = %v", err)
	}

	if root.StateStore == nil {
		t.Fatal("state store must be wired")
	}
	if root.MutationCoordinator == nil {
		t.Fatal("uow coordinator must be wired")
	}
	if root.PolicySeam != sentinelSeam {
		t.Fatal("policy seam value should be wired from constructor argument")
	}

	if _, ok := root.BundleView.Profile(evidence.ProfileAuthorizationDecision, evidence.ProfileVersionV1); !ok {
		t.Fatal("authorization-decision/v1 profile missing")
	}
	if _, ok := root.BundleView.Profile(evidence.ProfileApprovalDecision, evidence.ProfileVersionV1); !ok {
		t.Fatal("approval-decision/v1 profile missing")
	}
	if _, ok := root.BundleView.Profile(evidence.ProfileExceptionDecision, evidence.ProfileVersionV1); !ok {
		t.Fatal("exception-decision/v1 profile missing")
	}
	if _, ok := root.BundleView.Profile("unexpected-profile", evidence.ProfileVersionV1); ok {
		t.Fatal("unexpected fourth profile must not exist")
	}

	root.StateStore.SetResourceVersion("dependency-resource", "1")
	deps, err := state.NewAuthorizationDependencyVersionSet(root.StateStore.SnapshotVersion(), []state.ExpectedVersionPredicate{
		{ResourceUID: "dependency-resource", ResourceVersion: "1"},
	})
	if err != nil {
		t.Fatalf("dependency set: %v", err)
	}
	digest, fail := evidence.NewCandidateDigest([]byte("candidate-digest"))
	if fail.Reason() != "" {
		t.Fatalf("candidate digest: %s", fail.Reason())
	}
	claim, err := state.NewAuthorizationCurrentnessClaim(digest, deps)
	if err != nil {
		t.Fatalf("currentness claim: %v", err)
	}
	tx, err := root.StateStore.BeginEvidenceTransaction(claim)
	if err != nil {
		t.Fatalf("begin evidence tx: %v", err)
	}
	defer tx.Abort()
	if got := tx.PublicationContext().PublicationInstant(); !got.Equal(now) {
		t.Fatalf("publication instant = %v, want %v", got, now)
	}
}

func TestRootDefinesNoProfileAndNoSecondRegistry(t *testing.T) {
	t.Parallel()

	for name, f := range parseRootPackageFiles(t) {
		if name == "foundation.go" {
			continue
		}
		for _, imp := range f.Imports {
			path := strings.Trim(imp.Path.Value, `"`)
			if strings.HasSuffix(path, "/internal/decision") {
				t.Fatalf("%s imports decision package directly", name)
			}
		}
		ast.Inspect(f, func(n ast.Node) bool {
			// Root/bundle must not define DecisionProfile literals.
			if lit, ok := n.(*ast.CompositeLit); ok {
				if sel, ok := lit.Type.(*ast.SelectorExpr); ok && sel.Sel.Name == "DecisionProfile" {
					t.Fatalf("%s constructs decision.DecisionProfile", name)
				}
			}
			// Root/bundle must not create another VersionKey map registry.
			if mt, ok := n.(*ast.MapType); ok {
				if sel, ok := mt.Key.(*ast.SelectorExpr); ok && sel.Sel.Name == "VersionKey" {
					t.Fatalf("%s defines secondary bundle registry map", name)
				}
			}
			return true
		})
	}
}

func TestRootHoldsNoMutableState(t *testing.T) {
	t.Parallel()

	for name, f := range parseRootPackageFiles(t) {
		ast.Inspect(f, func(n ast.Node) bool {
			ts, ok := n.(*ast.TypeSpec)
			if !ok || ts.Name.Name != "Root" {
				return true
			}
			st, ok := ts.Type.(*ast.StructType)
			if !ok {
				t.Fatalf("%s Root must be a struct", name)
			}
			for _, field := range st.Fields.List {
				switch ft := field.Type.(type) {
				case *ast.MapType:
					t.Fatalf("%s Root must not own map state", name)
				case *ast.SelectorExpr:
					if ft.Sel.Name == "Mutex" || ft.Sel.Name == "RWMutex" {
						t.Fatalf("%s Root must not own mutex state", name)
					}
				}
			}
			return false
		})
	}
}

func TestRootOpensNoTransaction(t *testing.T) {
	t.Parallel()

	forbidden := map[string]bool{
		"InspectOrReserveCallerMutation": true,
		"BeginControllerMutation":        true,
		"BeginEvidenceTransaction":       true,
		"Begin":                          true,
		"Commit":                         true,
		"Seal":                           true,
	}
	for name, f := range parseRootPackageFiles(t) {
		ast.Inspect(f, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}
			if sel, ok := call.Fun.(*ast.SelectorExpr); ok && forbidden[sel.Sel.Name] {
				t.Fatalf("%s opens transaction path via %s", name, sel.Sel.Name)
			}
			return true
		})
	}
}

func TestPropertyRootConstructsNoDomainValue(t *testing.T) {
	t.Parallel()

	forbidden := map[string]bool{
		"NewPreparedAccessGroupIntent":       true,
		"NewPreparedMembershipIntent":        true,
		"NewPreparedRoleDefinitionIntent":    true,
		"NewPreparedRoleAssignmentIntent":    true,
		"newPreparedRoleAssignmentIntent":    true,
		"NewPreparedPrivilegedIntent":        true,
		"NewPreparedReviewIntent":            true,
		"NewPreparedApprovalIntent":          true,
		"NewPreparedExceptionIntent":         true,
		"NewPreparedGovernanceProfileIntent": true,
		"NewMutationEventFacts":              true,
		"NewResultMaterial":                  true,
	}

	prop := func(_ uint8) bool {
		for _, f := range parseRootPackageFiles(t) {
			violation := false
			ast.Inspect(f, func(n ast.Node) bool {
				call, ok := n.(*ast.CallExpr)
				if !ok {
					return true
				}
				switch fn := call.Fun.(type) {
				case *ast.Ident:
					if forbidden[fn.Name] {
						violation = true
						return false
					}
				case *ast.SelectorExpr:
					if forbidden[fn.Sel.Name] {
						violation = true
						return false
					}
				}
				return true
			})
			if violation {
				return false
			}
		}
		return true
	}

	if err := quick.Check(prop, &quick.Config{MaxCount: 8}); err != nil {
		t.Fatal(err)
	}
}

func parseRootPackageFiles(t *testing.T) map[string]*ast.File {
	t.Helper()

	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime caller")
	}
	dir := filepath.Dir(thisFile)
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dir, nil, 0)
	if err != nil {
		t.Fatal(err)
	}

	out := map[string]*ast.File{}
	for _, pkg := range pkgs {
		for name, f := range pkg.Files {
			if strings.HasSuffix(name, "_test.go") {
				continue
			}
			if strings.HasSuffix(name, "root.go") || strings.HasSuffix(name, "bundle.go") || strings.HasSuffix(name, "foundation.go") {
				out[name] = f
			}
		}
	}
	return out
}
