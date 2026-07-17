package registry

import (
	"context"
	"errors"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/resources"
)

func testOperation(id string, createdAt string) resources.Operation {
	return resources.Operation{
		APIVersion: resources.OperationAPIVersion,
		Kind:       resources.OperationKind,
		Metadata: resources.Metadata{
			Name:        id,
			Labels:      map[string]string{"env": "test"},
			Annotations: map[string]string{"note": "x"},
		},
		Spec: resources.OperationSpec{
			Type:                 resources.OpCreateProject,
			ResourceKind:         resources.ProjectKind,
			ResourceName:         "project-a",
			OrganizationName:     "org-a",
			OrganizationUnitName: "ou-a",
			TenantName:           "tenant-a",
			ProjectName:          "project-a",
			Actor:                "system",
		},
		Status: resources.OperationStatus{
			Phase:       resources.OperationPhaseSucceeded,
			CreatedAt:   createdAt,
			UpdatedAt:   createdAt,
			CompletedAt: createdAt,
		},
	}
}

func TestCreateOperation_MissingID(t *testing.T) {
	reg := NewOperationRegistry()
	ctx := context.Background()

	op := testOperation("", "2026-01-01T00:00:00Z")
	_, err := reg.CreateOperation(ctx, op)
	if !errors.Is(err, ErrMissingOperationID) {
		t.Fatalf("CreateOperation() error = %v, want ErrMissingOperationID", err)
	}

	items, err := reg.ListOperations(ctx)
	if err != nil {
		t.Fatalf("ListOperations() error = %v", err)
	}
	if items == nil {
		t.Fatal("ListOperations() returned nil slice, want non-nil empty slice")
	}
	if len(items) != 0 {
		t.Errorf("ListOperations() len = %d, want 0 (operation must not be stored)", len(items))
	}
}

func TestCreateOperation_Stores(t *testing.T) {
	reg := NewOperationRegistry()
	ctx := context.Background()

	op := testOperation("op-1", "2026-01-01T00:00:00Z")
	created, err := reg.CreateOperation(ctx, op)
	if err != nil {
		t.Fatalf("CreateOperation() error = %v", err)
	}
	if created.Metadata.Name != "op-1" || created.APIVersion != resources.OperationAPIVersion ||
		created.Kind != resources.OperationKind {
		t.Errorf("created identity = %+v, want id=op-1 apiVersion/kind set", created)
	}
	if created.Spec != op.Spec {
		t.Errorf("created.Spec = %+v, want %+v", created.Spec, op.Spec)
	}
	if created.Status != op.Status {
		t.Errorf("created.Status = %+v, want %+v", created.Status, op.Status)
	}

	got, err := reg.GetOperation(ctx, "op-1")
	if err != nil {
		t.Fatalf("GetOperation() error = %v", err)
	}
	if got.Metadata.Name != "op-1" || got.Spec != op.Spec || got.Status != op.Status {
		t.Errorf("got = %+v, want stored operation equal to %+v", got, op)
	}
}

func TestCreateOperation_Duplicate(t *testing.T) {
	reg := NewOperationRegistry()
	ctx := context.Background()

	original := testOperation("op-1", "2026-01-01T00:00:00Z")
	original.Spec.ResourceName = "original"
	if _, err := reg.CreateOperation(ctx, original); err != nil {
		t.Fatalf("first CreateOperation() error = %v", err)
	}

	dup := testOperation("op-1", "2026-02-02T00:00:00Z")
	dup.Spec.ResourceName = "changed"
	dup.Spec.Type = resources.OpDeleteProject
	_, err := reg.CreateOperation(ctx, dup)
	if !errors.Is(err, ErrAlreadyExists) {
		t.Fatalf("duplicate CreateOperation() error = %v, want ErrAlreadyExists", err)
	}

	got, err := reg.GetOperation(ctx, "op-1")
	if err != nil {
		t.Fatalf("GetOperation() error = %v", err)
	}
	if got.Spec.ResourceName != "original" || got.Spec.Type != resources.OpCreateProject ||
		got.Status.CreatedAt != "2026-01-01T00:00:00Z" {
		t.Errorf("stored operation was overwritten: %+v", got)
	}
}

func TestGetOperation_NotFound(t *testing.T) {
	reg := NewOperationRegistry()
	ctx := context.Background()

	if _, err := reg.CreateOperation(ctx, testOperation("op-1", "2026-01-01T00:00:00Z")); err != nil {
		t.Fatalf("CreateOperation() error = %v", err)
	}

	if _, err := reg.GetOperation(ctx, "op-1"); err != nil {
		t.Fatalf("GetOperation(existing) error = %v", err)
	}
	if _, err := reg.GetOperation(ctx, "missing"); !errors.Is(err, ErrNotFound) {
		t.Errorf("GetOperation(missing) error = %v, want ErrNotFound", err)
	}
}

func TestListOperations_EmptyNonNil(t *testing.T) {
	reg := NewOperationRegistry()
	items, err := reg.ListOperations(context.Background())
	if err != nil {
		t.Fatalf("ListOperations() error = %v", err)
	}
	if items == nil {
		t.Fatal("ListOperations() returned nil, want non-nil empty slice")
	}
	if len(items) != 0 {
		t.Errorf("ListOperations() len = %d, want 0", len(items))
	}
}

func TestListOperations_Sorted(t *testing.T) {
	reg := NewOperationRegistry()
	ctx := context.Background()

	// Insert out of order. "b" and "c" share a createdAt to exercise the
	// metadata.name tie-breaker.
	inputs := []resources.Operation{
		testOperation("c", "2026-01-02T00:00:00Z"),
		testOperation("a", "2026-01-03T00:00:00Z"),
		testOperation("b", "2026-01-02T00:00:00Z"),
		testOperation("d", "2026-01-01T00:00:00Z"),
	}
	for _, op := range inputs {
		if _, err := reg.CreateOperation(ctx, op); err != nil {
			t.Fatalf("CreateOperation(%s) error = %v", op.Metadata.Name, err)
		}
	}

	items, err := reg.ListOperations(ctx)
	if err != nil {
		t.Fatalf("ListOperations() error = %v", err)
	}
	wantOrder := []string{"d", "b", "c", "a"}
	if len(items) != len(wantOrder) {
		t.Fatalf("ListOperations() len = %d, want %d", len(items), len(wantOrder))
	}
	for i, want := range wantOrder {
		if items[i].Metadata.Name != want {
			t.Errorf("item[%d] = %q, want %q (sort by createdAt asc, then name asc)", i, items[i].Metadata.Name, want)
		}
	}
}

func TestOperationRegistry_DeepCopyImmutability(t *testing.T) {
	reg := NewOperationRegistry()
	ctx := context.Background()

	created, err := reg.CreateOperation(ctx, testOperation("op-1", "2026-01-01T00:00:00Z"))
	if err != nil {
		t.Fatalf("CreateOperation() error = %v", err)
	}

	// Mutate the copy returned by Create.
	created.Metadata.Labels["env"] = "mutated"
	created.Metadata.Annotations["note"] = "mutated"

	got, err := reg.GetOperation(ctx, "op-1")
	if err != nil {
		t.Fatalf("GetOperation() error = %v", err)
	}
	if got.Metadata.Labels["env"] != "test" || got.Metadata.Annotations["note"] != "x" {
		t.Errorf("Create return shares maps with store: labels=%v annotations=%v", got.Metadata.Labels, got.Metadata.Annotations)
	}

	// Mutate the copy returned by Get.
	got.Metadata.Labels["env"] = "mutated-again"
	got.Metadata.Annotations["note"] = "mutated-again"

	items, err := reg.ListOperations(ctx)
	if err != nil {
		t.Fatalf("ListOperations() error = %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("ListOperations() len = %d, want 1", len(items))
	}
	if items[0].Metadata.Labels["env"] != "test" || items[0].Metadata.Annotations["note"] != "x" {
		t.Errorf("Get return shares maps with store: labels=%v annotations=%v", items[0].Metadata.Labels, items[0].Metadata.Annotations)
	}

	// Mutate the copy returned by List.
	items[0].Metadata.Labels["env"] = "list-mutated"
	items[0].Metadata.Annotations["note"] = "list-mutated"

	after, err := reg.GetOperation(ctx, "op-1")
	if err != nil {
		t.Fatalf("GetOperation() error = %v", err)
	}
	if after.Metadata.Labels["env"] != "test" || after.Metadata.Annotations["note"] != "x" {
		t.Errorf("List return shares maps with store: labels=%v annotations=%v", after.Metadata.Labels, after.Metadata.Annotations)
	}
}
