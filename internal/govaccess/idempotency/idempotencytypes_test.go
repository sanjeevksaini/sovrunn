package idempotency_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/idempotency"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/model"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/operation"
)

func testActor() model.PrincipalRef {
	return model.PrincipalRef{
		Issuer:        "https://idp.example",
		Subject:       "user-1",
		PrincipalType: model.PrincipalTypeHuman,
	}
}

func testBinding(t *testing.T) idempotency.RequestBinding {
	t.Helper()
	b, fail := idempotency.NewRequestBinding(
		apimeta.TypedRef{APIVersion: "gov.sovrunn.io/v1alpha1", Kind: "RoleAssignment", Name: "ra-1", UID: "uid-1"},
		apimeta.ScopeRef{TypedRef: apimeta.TypedRef{APIVersion: "core.sovrunn.io/v1alpha1", Kind: "Project", Name: "p1", UID: "p-uid"}},
		[]byte("digest-1"),
	)
	if fail.Reason() != "" {
		t.Fatalf("NewRequestBinding: %s", fail.Reason())
	}
	return b
}

func TestIdempotencyLookupKeyConstruction(t *testing.T) {
	t.Parallel()
	k, fail := idempotency.NewIdempotencyLookupKey(testActor(), "roleassignment.create", "client-key-1")
	if fail.Reason() != "" {
		t.Fatalf("unexpected failure: %s", fail.Reason())
	}
	if !k.Sealed() || k.Operation() != "roleassignment.create" || k.ClientIdempotencyKey() != "client-key-1" {
		t.Fatalf("unexpected key values")
	}
	_, fail = idempotency.NewIdempotencyLookupKey(testActor(), "", "k")
	if fail.Reason() == "" {
		t.Fatal("expected empty operation failure")
	}
}

func TestRequestBindingDistinctFromLookupKey(t *testing.T) {
	t.Parallel()
	k, fail := idempotency.NewIdempotencyLookupKey(testActor(), "op", "key")
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	b := testBinding(t)
	if len(b.CanonicalRequestDigest()) == 0 {
		t.Fatal("binding must carry digest")
	}
	if k.ActorRef().Subject == "" {
		t.Fatal("lookup must carry actor")
	}
	if k.ClientIdempotencyKey() == string(b.CanonicalRequestDigest()) {
		t.Fatal("lookup key and request binding must remain distinct value types")
	}
}

func TestInFlightReservationPreparedCompletedResultCompletedRecord(t *testing.T) {
	t.Parallel()
	k, fail := idempotency.NewIdempotencyLookupKey(testActor(), "op", "key")
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	b := testBinding(t)
	res, fail := idempotency.NewInFlightReservation(k, b, []byte("owner-token"))
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	if !res.Binding().Equal(b) || !res.LookupKey().Equal(k) {
		t.Fatal("reservation must carry full binding and lookup")
	}

	rm, fail := operation.NewResultMaterial([]byte(`{"ok":true}`))
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	pub, fail := operation.NewPublicationBinding([]byte("pub"))
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	md, fail := operation.NewMutationDigest([]byte("mut"))
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	cb, fail := operation.NewCompletionBinding(pub, md)
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}

	prep, fail := idempotency.NewPreparedCompletedResult(k, b, rm, cb)
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	if !prep.Binding().Equal(b) {
		t.Fatal("prepared result must carry full RequestBinding")
	}
	rec, fail := idempotency.NewCompletedRecord(k, b, rm, cb)
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	if !rec.Binding().Equal(b) {
		t.Fatal("completed record must carry full RequestBinding")
	}
}

func TestIdempotencyValueTypesDeterministic(t *testing.T) {
	t.Parallel()
	a1, f1 := idempotency.NewIdempotencyLookupKey(testActor(), "op", "key")
	a2, f2 := idempotency.NewIdempotencyLookupKey(testActor(), "op", "key")
	if f1.Reason() != "" || f2.Reason() != "" || !a1.Equal(a2) {
		t.Fatalf("deterministic construction failed: %q %q", f1.Reason(), f2.Reason())
	}
	b1 := testBinding(t)
	b2, fail := idempotency.NewRequestBinding(b1.ExactTarget(), b1.ScopeRef(), b1.CanonicalRequestDigest())
	if fail.Reason() != "" || !b1.Equal(b2) {
		t.Fatal("binding construction not deterministic")
	}
}

func TestIdempotencyPackageImportsOnlyAllowedLeaves(t *testing.T) {
	t.Parallel()
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("caller")
	}
	dir := filepath.Dir(thisFile)
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dir, nil, 0)
	if err != nil {
		t.Fatal(err)
	}
	forbidden := []string{
		"github.com/sanjeevksaini/sovrunn/internal/govaccess/state",
		"github.com/sanjeevksaini/sovrunn/internal/govaccess/evidence",
		"github.com/sanjeevksaini/sovrunn/internal/govaccess/uow",
	}
	for _, pkg := range pkgs {
		for name, f := range pkg.Files {
			if strings.HasSuffix(name, "_test.go") {
				continue
			}
			for _, imp := range f.Imports {
				path := strings.Trim(imp.Path.Value, `"`)
				for _, bad := range forbidden {
					if path == bad {
						t.Fatalf("forbidden import %s in %s", path, name)
					}
				}
			}
			ast.Inspect(f, func(n ast.Node) bool {
				if id, ok := n.(*ast.Ident); ok {
					switch id.Name {
					case "InFlightTable", "PrepareCompletedResult":
						t.Fatalf("behavioral symbol %s must not appear in Task-2 value-type files", id.Name)
					}
				}
				return true
			})
		}
	}
}
