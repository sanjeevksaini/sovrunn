package idempotency_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/idempotency"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/model"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/operation"
)

type fixedClock struct {
	now time.Time
}

func (c fixedClock) NowUTC() time.Time { return c.now.UTC() }

func mkLookupKey(t *testing.T, key string) idempotency.IdempotencyLookupKey {
	t.Helper()
	actor := model.PrincipalRef{
		Issuer:        "https://idp.example",
		Subject:       "user-1",
		PrincipalType: model.PrincipalTypeHuman,
	}
	out, fail := idempotency.NewIdempotencyLookupKey(actor, "roleassignment.grant", key)
	if fail.Reason() != "" {
		t.Fatalf("lookup key: %s", fail.Reason())
	}
	return out
}

func mkBinding(t *testing.T, digest string) idempotency.RequestBinding {
	t.Helper()
	out, fail := idempotency.NewRequestBinding(
		apimeta.TypedRef{APIVersion: model.APIVersionIAM, Kind: "RoleAssignment", Name: "ra-1", UID: "ra-1"},
		apimeta.ScopeRef{TypedRef: apimeta.TypedRef{APIVersion: "core.sovrunn.io/v1alpha1", Kind: "Project", Name: "p1", UID: "p1"}},
		[]byte(digest),
	)
	if fail.Reason() != "" {
		t.Fatalf("request binding: %s", fail.Reason())
	}
	return out
}

func mkPrepared(t *testing.T, key idempotency.IdempotencyLookupKey, binding idempotency.RequestBinding) idempotency.PreparedCompletedResult {
	t.Helper()
	result, fail := operation.NewResultMaterial([]byte(`{"ok":true}`))
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	pub, fail := operation.NewPublicationBinding([]byte("publication"))
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	mut, fail := operation.NewMutationDigest([]byte("mutation"))
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	cb, fail := operation.NewCompletionBinding(pub, mut)
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	prepared, fail := idempotency.StageCompleted(key, binding, result, cb)
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	return prepared
}

func TestInFlightTableReserveLookupComplete(t *testing.T) {
	t.Parallel()
	table := idempotency.NewReservationTable()
	key := mkLookupKey(t, "k-1")
	binding := mkBinding(t, "digest-1")
	res, fail := idempotency.NewInFlightReservation(key, binding, []byte("owner-token"))
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	ok, wait, fail := table.Reserve(res)
	if fail.Reason() != "" || !ok || wait == nil {
		t.Fatalf("reserve ok=%v wait=%v fail=%s", ok, wait, fail.Reason())
	}
	lookup, fail := table.Lookup(key)
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	if in, _, found := lookup.InFlight(); !found || !in.Binding().Equal(binding) {
		t.Fatal("expected in-flight reservation")
	}

	prepared := mkPrepared(t, key, binding)
	_, retainUntil, fail := idempotency.CompleteLifecycle(prepared, fixedClock{now: time.Date(2026, 9, 7, 11, 0, 0, 0, time.UTC)})
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	if fail := table.Complete(prepared, retainUntil); fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	lookup, fail = table.Lookup(key)
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	if rec, _, found := lookup.Completed(); !found || !rec.Binding().Equal(binding) {
		t.Fatal("expected completed record")
	}
}

func TestReplayDetectionComparesBindingAfterLookup(t *testing.T) {
	t.Parallel()
	table := idempotency.NewReservationTable()
	key := mkLookupKey(t, "k-2")
	b1 := mkBinding(t, "digest-1")
	b2 := mkBinding(t, "digest-2")
	res, fail := idempotency.NewInFlightReservation(key, b1, []byte("owner-token"))
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	if ok, _, fail := table.Reserve(res); !ok || fail.Reason() != "" {
		t.Fatal("reserve should succeed")
	}
	lookup, fail := table.Lookup(key)
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	decision, fail := idempotency.EvaluateReplayRecheck(lookup, b2)
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	if decision.Kind() != idempotency.ReplayDecisionConflict {
		t.Fatalf("kind=%s", decision.Kind())
	}
}

func TestWaitOnInFlightReservation(t *testing.T) {
	t.Parallel()
	table := idempotency.NewReservationTable()
	key := mkLookupKey(t, "k-3")
	binding := mkBinding(t, "digest-3")
	res, fail := idempotency.NewInFlightReservation(key, binding, []byte("owner-token"))
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	if ok, _, fail := table.Reserve(res); !ok || fail.Reason() != "" {
		t.Fatal("reserve should succeed")
	}

	lookup, fail := table.Lookup(key)
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	decision, fail := idempotency.EvaluateReplayRecheck(lookup, binding)
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	if decision.Kind() != idempotency.ReplayDecisionWait {
		t.Fatalf("kind=%s", decision.Kind())
	}
	wait, ok := decision.WaitChannel()
	if !ok || wait == nil {
		t.Fatal("wait channel missing")
	}

	done := make(chan struct{})
	go func() {
		defer close(done)
		<-wait
	}()

	prepared := mkPrepared(t, key, binding)
	_, retainUntil, fail := idempotency.CompleteLifecycle(prepared, fixedClock{now: time.Date(2026, 9, 7, 11, 30, 0, 0, time.UTC)})
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	if fail := table.Complete(prepared, retainUntil); fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("waiter not released")
	}
}

func TestCompleteLifecycleUsesInjectedClockHorizon(t *testing.T) {
	t.Parallel()
	key := mkLookupKey(t, "k-4")
	binding := mkBinding(t, "digest-4")
	prepared := mkPrepared(t, key, binding)
	now := time.Date(2026, 9, 7, 12, 0, 0, 0, time.UTC)
	_, retainUntil, fail := idempotency.CompleteLifecycle(prepared, fixedClock{now: now})
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	if got, want := retainUntil, now.Add(24*time.Hour); !got.Equal(want) {
		t.Fatalf("retainUntil=%s want=%s", got, want)
	}
}

func TestStageCompletedPreservesNeutralCompletionBinding(t *testing.T) {
	t.Parallel()
	key := mkLookupKey(t, "k-5")
	binding := mkBinding(t, "digest-5")
	result, fail := operation.NewResultMaterial([]byte(`{"ok":true}`))
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	pub, fail := operation.NewPublicationBinding([]byte("pub"))
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	mut, fail := operation.NewMutationDigest([]byte("mut"))
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	completion, fail := operation.NewCompletionBinding(pub, mut)
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	prepared, fail := idempotency.StageCompleted(key, binding, result, completion)
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	if !prepared.Sealed() || !prepared.Binding().Equal(binding) || !prepared.LookupKey().Equal(key) {
		t.Fatal("prepared completed result mismatch")
	}
	if !prepared.Result().Equal(result) || !prepared.Completion().Sealed() {
		t.Fatal("result/completion not preserved")
	}
}

func TestPropertyIdempotencyBehaviorNoEvidenceImportAndNoValueTypeRedefinition(t *testing.T) {
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
	redefined := map[string]bool{
		"IdempotencyLookupKey":    false,
		"RequestBinding":          false,
		"InFlightReservation":     false,
		"PreparedCompletedResult": false,
		"CompletedRecord":         false,
	}
	for _, pkg := range pkgs {
		for name, file := range pkg.Files {
			base := filepath.Base(name)
			if strings.HasSuffix(base, "_test.go") {
				continue
			}
			for _, imp := range file.Imports {
				path := strings.Trim(imp.Path.Value, `"`)
				if path == "github.com/sanjeevksaini/sovrunn/internal/govaccess/evidence" {
					t.Fatalf("forbidden evidence import in %s", name)
				}
			}
			ast.Inspect(file, func(n ast.Node) bool {
				ts, ok := n.(*ast.TypeSpec)
				if !ok {
					return true
				}
				if _, found := redefined[ts.Name.Name]; found && base != "lookupkey.go" && base != "binding.go" && base != "reservation.go" && base != "completedresult.go" {
					redefined[ts.Name.Name] = true
				}
				return true
			})
		}
	}
	for name, bad := range redefined {
		if bad {
			t.Fatalf("value type %s was redefined outside Task-2 files", name)
		}
	}
}
