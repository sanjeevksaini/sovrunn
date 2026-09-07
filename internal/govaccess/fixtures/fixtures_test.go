package fixtures

import (
	"errors"
	"testing"
	"testing/quick"
)

func TestBuilderIsImmutableAndDeepCopied(t *testing.T) {
	t.Parallel()

	base := NewBuilder()
	a := base.WithResource(Resource{Kind: KindAccessGroup, UID: "ag-1", RefUIDs: []string{"rd-1"}})
	b := a.WithResource(Resource{Kind: KindRoleDefinition, UID: "rd-1"})

	aSet := a.Build()
	bSet := b.Build()
	if len(aSet.Resources) != 1 {
		t.Fatalf("a resources = %d, want 1", len(aSet.Resources))
	}
	if len(bSet.Resources) != 2 {
		t.Fatalf("b resources = %d, want 2", len(bSet.Resources))
	}

	aSet.Resources[0].RefUIDs[0] = "changed"
	again := a.Build()
	if again.Resources[0].RefUIDs[0] != "rd-1" {
		t.Fatal("builder did not deep copy nested slices")
	}
}

func TestValidatorReferentialClosure(t *testing.T) {
	t.Parallel()

	okSet := NewBuilder().
		WithResource(Resource{Kind: KindAccessGroup, UID: "ag-1"}).
		WithResource(Resource{Kind: KindRoleDefinition, UID: "rd-1"}).
		WithResource(Resource{Kind: KindRoleAssignment, UID: "ra-1", RefUIDs: []string{"ag-1", "rd-1"}}).
		Build()
	if err := Validate(okSet); err != nil {
		t.Fatalf("Validate(okSet) error = %v", err)
	}

	badSet := NewBuilder().
		WithResource(Resource{Kind: KindRoleAssignment, UID: "ra-2", RefUIDs: []string{"missing"}}).
		Build()
	if err := Validate(badSet); err == nil {
		t.Fatal("expected missing reference error")
	}
}

func TestExecutorDeterministicAndZeroExternalIO(t *testing.T) {
	t.Parallel()

	set := NewBuilder().
		WithResource(Resource{Kind: KindAccessGroup, UID: "ag-1"}).
		WithResource(Resource{Kind: KindRoleDefinition, UID: "rd-1"}).
		WithResource(Resource{Kind: KindRoleAssignment, UID: "ra-1", RefUIDs: []string{"rd-1", "ag-1"}}).
		Build()

	exec := NewFakeExecutor()
	first, err := exec.Execute(set)
	if err != nil {
		t.Fatalf("Execute(first) error = %v", err)
	}
	second, err := exec.Execute(set)
	if err != nil {
		t.Fatalf("Execute(second) error = %v", err)
	}
	if first != second {
		t.Fatalf("results differ: first=%+v second=%+v", first, second)
	}
	if exec.ExternalCallCount() != 0 {
		t.Fatalf("externalCallCount = %d, want 0", exec.ExternalCallCount())
	}
}

func TestPropertySameFixtureSameExecutionResult(t *testing.T) {
	t.Parallel()

	prop := func(n uint8) bool {
		set := NewBuilder().
			WithResource(Resource{Kind: KindAccessGroup, UID: "ag-1"}).
			WithResource(Resource{Kind: KindRoleDefinition, UID: "rd-1"}).
			WithResource(Resource{Kind: KindRoleAssignment, UID: "ra-1", RefUIDs: []string{"rd-1", "ag-1"}}).
			Build()
		if n%2 == 0 {
			set = NewBuilder().
				WithResource(Resource{Kind: KindRoleDefinition, UID: "rd-1"}).
				WithResource(Resource{Kind: KindRoleAssignment, UID: "ra-1", RefUIDs: []string{"ag-1", "rd-1"}}).
				WithResource(Resource{Kind: KindAccessGroup, UID: "ag-1"}).
				Build()
		}
		exec := NewFakeExecutor()
		a, errA := exec.Execute(set)
		b, errB := exec.Execute(set)
		if errA != nil || errB != nil {
			return false
		}
		return a == b && exec.ExternalCallCount() == 0
	}

	if err := quick.Check(prop, &quick.Config{MaxCount: 32}); err != nil {
		t.Fatal(err)
	}
}

func TestUnknownKindRejected(t *testing.T) {
	t.Parallel()

	set := NewBuilder().WithResource(Resource{Kind: ResourceKind("CloudEnrollment"), UID: "ce-1"}).Build()
	err := Validate(set)
	if err == nil {
		t.Fatal("expected unknown kind rejection")
	}
	if errors.Is(err, ErrProhibitedConcept) {
		t.Fatalf("expected ErrInvalidSet, got %v", err)
	}
}
