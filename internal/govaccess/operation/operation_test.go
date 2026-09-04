package operation_test

import (
	"bytes"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/govaccess/operation"
)

func TestNewResultMaterialRejectsEmpty(t *testing.T) {
	t.Parallel()
	if _, fail := operation.NewResultMaterial(nil); fail.Reason() == "" {
		t.Fatal("expected mechanical failure for nil bytes")
	}
	if _, fail := operation.NewResultMaterial([]byte{}); fail.Reason() == "" {
		t.Fatal("expected mechanical failure for empty bytes")
	}
}

func TestResultMaterialImmutability(t *testing.T) {
	t.Parallel()
	src := []byte("canonical-result")
	mat, fail := operation.NewResultMaterial(src)
	if fail.Reason() != "" {
		t.Fatalf("unexpected failure: %s", fail.Reason())
	}
	src[0] = 'X'
	got := mat.Bytes()
	if !bytes.Equal(got, []byte("canonical-result")) {
		t.Fatalf("mutated input leaked into ResultMaterial: %q", got)
	}
	got[0] = 'Y'
	if !bytes.Equal(mat.Bytes(), []byte("canonical-result")) {
		t.Fatal("Bytes() must return a defensive copy")
	}
	if !mat.Sealed() {
		t.Fatal("expected sealed ResultMaterial")
	}
}

func TestAppliedChangeLinkSealedConstruction(t *testing.T) {
	t.Parallel()
	p, _ := operation.NewParticipantBinding([]byte("p"))
	pub, _ := operation.NewPublicationBinding([]byte("pub"))
	shadow, _ := operation.NewShadowDeltaDigest([]byte("s"))
	index, _ := operation.NewIndexDeltaDigest([]byte("i"))
	result, _ := operation.NewResultMaterial([]byte("r"))

	link, fail := operation.NewAppliedChangeLink(p, pub, shadow, index, result, true)
	if fail.Reason() != "" {
		t.Fatalf("unexpected failure: %s", fail.Reason())
	}
	if !link.Sealed() || !link.HasResult() {
		t.Fatal("expected sealed link with result")
	}
	_, fail = operation.NewAppliedChangeLink(operation.ParticipantBinding{}, pub, shadow, index, result, true)
	if fail.Reason() == "" {
		t.Fatal("expected mechanical failure for unsealed participant")
	}
}

func TestAppliedConclusionAndCompletionBinding(t *testing.T) {
	t.Parallel()
	p, _ := operation.NewParticipantBinding([]byte("p"))
	pub, _ := operation.NewPublicationBinding([]byte("pub"))
	link, fail := operation.NewAppliedConclusionLink(p, pub)
	if fail.Reason() != "" || !link.Sealed() {
		t.Fatalf("conclusion link failed: %s", fail.Reason())
	}
	mut, _ := operation.NewMutationDigest([]byte("m"))
	bind, fail := operation.NewCompletionBinding(pub, mut)
	if fail.Reason() != "" || !bind.Sealed() {
		t.Fatalf("completion binding failed: %s", fail.Reason())
	}
	if _, fail := operation.NewCompletionBinding(operation.PublicationBinding{}, mut); fail.Reason() == "" {
		t.Fatal("expected failure for unsealed publication")
	}
}

func TestFinalizationOutcomeExactlyOneVariant(t *testing.T) {
	t.Parallel()
	type change struct{ ID string }
	type domain struct{ Code string }

	fin := operation.Finalized[change, domain](change{ID: "c1"})
	if fin.KindName() != "Finalized" {
		t.Fatalf("kind = %s", fin.KindName())
	}
	if v, ok := fin.AsFinalized(); !ok || v.ID != "c1" {
		t.Fatal("expected Finalized variant")
	}
	if _, ok := fin.AsDomain(); ok {
		t.Fatal("Finalized must not expose Domain")
	}
	if _, ok := fin.AsMechanicalFailure(); ok {
		t.Fatal("Finalized must not expose MechanicalFailure")
	}

	dom := operation.Domain[change, domain](domain{Code: "Deny"})
	if dom.KindName() != "Domain" {
		t.Fatalf("kind = %s", dom.KindName())
	}
	if _, ok := dom.AsFinalized(); ok {
		t.Fatal("Domain must not expose Finalized")
	}

	fail := operation.Failure[change, domain](operation.NewMechanicalFailure("invariant"))
	if fail.KindName() != "MechanicalFailure" {
		t.Fatalf("kind = %s", fail.KindName())
	}
	if f, ok := fail.AsMechanicalFailure(); !ok || f.Reason() != "invariant" {
		t.Fatal("expected MechanicalFailure variant")
	}
}
