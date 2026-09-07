package fixtures

import (
	"errors"
	"testing"
)

func TestVS0CFF1850ProhibitedConceptRejectedBeforeExecution(t *testing.T) {
	t.Parallel()

	set := FixtureWithProhibitedFutureOwnedResource()
	exec := NewFakeExecutor()

	_, err := exec.Execute(set)
	if err == nil {
		t.Fatal("expected prohibited concept rejection")
	}
	if !errors.Is(err, ErrProhibitedConcept) {
		t.Fatalf("error = %v, want ErrProhibitedConcept", err)
	}
	if exec.ExternalCallCount() != 0 {
		t.Fatalf("externalCallCount = %d, want 0", exec.ExternalCallCount())
	}
}

func TestNewProhibitedConceptValidation(t *testing.T) {
	t.Parallel()

	if _, err := NewProhibitedConcept(ProhibitedFutureField, "metadata.futureField"); err != nil {
		t.Fatalf("valid concept error = %v", err)
	}
	if _, err := NewProhibitedConcept(ProhibitedCategory("bad"), "x"); err == nil {
		t.Fatal("expected category validation error")
	}
}
