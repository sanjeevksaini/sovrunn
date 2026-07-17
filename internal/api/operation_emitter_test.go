package api

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/sanjeevksaini/sovrunn/internal/registry"
	"github.com/sanjeevksaini/sovrunn/internal/resources"
)

// stubOperationRegistry is a configurable registry.OperationRegistryIface used
// to drive collision and error paths deterministically.
type stubOperationRegistry struct {
	createCalls int
	// failFirst returns ErrAlreadyExists for the first failFirst calls, then
	// succeeds and captures the Operation.
	failFirst int
	// alwaysErr, when non-nil, is returned from every CreateOperation call.
	alwaysErr error
	captured  resources.Operation
}

func (s *stubOperationRegistry) CreateOperation(ctx context.Context, op resources.Operation) (resources.Operation, error) {
	s.createCalls++
	if s.alwaysErr != nil {
		return resources.Operation{}, s.alwaysErr
	}
	if s.createCalls <= s.failFirst {
		return resources.Operation{}, registry.ErrAlreadyExists
	}
	s.captured = op
	return op, nil
}

func (s *stubOperationRegistry) GetOperation(ctx context.Context, id string) (resources.Operation, error) {
	return resources.Operation{}, registry.ErrNotFound
}

func (s *stubOperationRegistry) ListOperations(ctx context.Context) ([]resources.Operation, error) {
	return []resources.Operation{}, nil
}

// failingEmitter is an OperationEmitter whose Emit always returns an error.
type failingEmitter struct {
	calls int
}

func (f *failingEmitter) Emit(ctx context.Context, spec resources.OperationSpec) error {
	f.calls++
	return errors.New("emit failed")
}

func sampleSpec() resources.OperationSpec {
	return resources.OperationSpec{
		Type:                 resources.OpCreateProject,
		ResourceKind:         resources.ProjectKind,
		ResourceName:         "project-a",
		OrganizationName:     "org-a",
		OrganizationUnitName: "ou-a",
		TenantName:           "tenant-a",
		ProjectName:          "project-a",
		RequestID:            "req-123",
	}
}

func TestNewOperationID(t *testing.T) {
	const samples = 100
	seen := make(map[string]struct{}, samples)
	for i := 0; i < samples; i++ {
		id, err := newOperationID()
		if err != nil {
			t.Fatalf("newOperationID() error = %v", err)
		}
		if len(id) != 32 {
			t.Fatalf("id %q length = %d, want 32", id, len(id))
		}
		for _, c := range id {
			if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f')) {
				t.Fatalf("id %q contains non-lowercase-hex char %q", id, c)
			}
		}
		if _, dup := seen[id]; dup {
			t.Fatalf("duplicate id generated: %q", id)
		}
		seen[id] = struct{}{}
	}
}

func TestEmit_StoresWellFormedOperation(t *testing.T) {
	reg := registry.NewOperationRegistry()
	emitter := NewRegistryEmitter(reg, nil)
	ctx := context.Background()

	if err := emitter.Emit(ctx, sampleSpec()); err != nil {
		t.Fatalf("Emit() error = %v", err)
	}

	items, err := reg.ListOperations(ctx)
	if err != nil {
		t.Fatalf("ListOperations() error = %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("stored operations = %d, want 1", len(items))
	}
	op := items[0]

	if op.APIVersion != resources.OperationAPIVersion {
		t.Errorf("APIVersion = %q, want %q", op.APIVersion, resources.OperationAPIVersion)
	}
	if op.Kind != resources.OperationKind {
		t.Errorf("Kind = %q, want %q", op.Kind, resources.OperationKind)
	}
	if op.Metadata.Name == "" {
		t.Error("Metadata.Name is empty, want generated ID")
	}
	if op.Spec.Actor != "system" {
		t.Errorf("Spec.Actor = %q, want system", op.Spec.Actor)
	}
	if op.Spec.Type != resources.OpCreateProject || op.Spec.ResourceKind != resources.ProjectKind ||
		op.Spec.ResourceName != "project-a" {
		t.Errorf("Spec identity not preserved: %+v", op.Spec)
	}
	if op.Spec.OrganizationName != "org-a" || op.Spec.OrganizationUnitName != "ou-a" ||
		op.Spec.TenantName != "tenant-a" || op.Spec.ProjectName != "project-a" {
		t.Errorf("Spec parent refs not preserved: %+v", op.Spec)
	}
	if op.Spec.RequestID != "req-123" {
		t.Errorf("Spec.RequestID = %q, want req-123", op.Spec.RequestID)
	}
	if op.Status.Phase != resources.OperationPhaseSucceeded {
		t.Errorf("Status.Phase = %q, want Succeeded", op.Status.Phase)
	}
	if op.Status.CreatedAt == "" {
		t.Error("Status.CreatedAt is empty")
	}
	if op.Status.UpdatedAt != op.Status.CreatedAt {
		t.Errorf("Status.UpdatedAt = %q, want == CreatedAt %q", op.Status.UpdatedAt, op.Status.CreatedAt)
	}
	if op.Status.CompletedAt != op.Status.CreatedAt {
		t.Errorf("Status.CompletedAt = %q, want == CreatedAt %q", op.Status.CompletedAt, op.Status.CreatedAt)
	}
	if _, err := time.Parse(time.RFC3339, op.Status.CreatedAt); err != nil {
		t.Errorf("Status.CreatedAt %q is not RFC3339: %v", op.Status.CreatedAt, err)
	}
}

func TestEmit_ForcesActorSystem(t *testing.T) {
	reg := registry.NewOperationRegistry()
	emitter := NewRegistryEmitter(reg, nil)
	ctx := context.Background()

	spec := sampleSpec()
	spec.Actor = "malicious-user"
	if err := emitter.Emit(ctx, spec); err != nil {
		t.Fatalf("Emit() error = %v", err)
	}

	items, err := reg.ListOperations(ctx)
	if err != nil {
		t.Fatalf("ListOperations() error = %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("stored operations = %d, want 1", len(items))
	}
	if items[0].Spec.Actor != "system" {
		t.Errorf("Spec.Actor = %q, want system", items[0].Spec.Actor)
	}
}

func TestEmit_NilAndUnavailable(t *testing.T) {
	ctx := context.Background()
	spec := sampleSpec()

	var nilEmitter *registryEmitter
	if err := nilEmitter.Emit(ctx, spec); err == nil {
		t.Error("nil *registryEmitter Emit() error = nil, want non-nil")
	} else if !errors.Is(err, errOperationEmitterUnavailable) {
		t.Errorf("nil Emit() error = %v, want errOperationEmitterUnavailable", err)
	}

	if err := NewRegistryEmitter(nil, nil).Emit(ctx, spec); err == nil {
		t.Error("NewRegistryEmitter(nil, nil).Emit() error = nil, want non-nil")
	} else if !errors.Is(err, errOperationEmitterUnavailable) {
		t.Errorf("nil-registry Emit() error = %v, want errOperationEmitterUnavailable", err)
	}
}

func TestEmit_CollisionRetrySucceeds(t *testing.T) {
	stub := &stubOperationRegistry{failFirst: 2}
	emitter := NewRegistryEmitter(stub, nil)
	ctx := context.Background()

	if err := emitter.Emit(ctx, sampleSpec()); err != nil {
		t.Fatalf("Emit() error = %v, want nil after retry", err)
	}
	if stub.createCalls != 3 {
		t.Errorf("CreateOperation calls = %d, want 3 (2 collisions + 1 success)", stub.createCalls)
	}
	if stub.captured.Spec.Actor != "system" {
		t.Errorf("captured Actor = %q, want system", stub.captured.Spec.Actor)
	}
	if stub.captured.Status.Phase != resources.OperationPhaseSucceeded {
		t.Errorf("captured Phase = %q, want Succeeded", stub.captured.Status.Phase)
	}
	if stub.captured.Metadata.Name == "" {
		t.Error("captured Metadata.Name is empty")
	}
}

func TestEmit_CollisionExhaustion(t *testing.T) {
	stub := &stubOperationRegistry{alwaysErr: registry.ErrAlreadyExists}
	emitter := NewRegistryEmitter(stub, nil)
	ctx := context.Background()

	err := emitter.Emit(ctx, sampleSpec())
	if !errors.Is(err, errOperationIDExhausted) {
		t.Fatalf("Emit() error = %v, want errOperationIDExhausted", err)
	}
	if stub.createCalls != 5 {
		t.Errorf("CreateOperation calls = %d, want 5", stub.createCalls)
	}
}

func TestEmit_NonCollisionError(t *testing.T) {
	boom := errors.New("boom")
	stub := &stubOperationRegistry{alwaysErr: boom}
	emitter := NewRegistryEmitter(stub, nil)
	ctx := context.Background()

	err := emitter.Emit(ctx, sampleSpec())
	if !errors.Is(err, boom) {
		t.Fatalf("Emit() error = %v, want boom", err)
	}
	if stub.createCalls != 1 {
		t.Errorf("CreateOperation calls = %d, want 1 (no retry on non-collision error)", stub.createCalls)
	}
}

func TestEmitOperation_NilSafe(t *testing.T) {
	// Must not panic with a nil emitter.
	emitOperation(context.Background(), nil, sampleSpec())
}

func TestEmitOperation_SwallowsErrors(t *testing.T) {
	f := &failingEmitter{}
	// Must not panic or surface an error.
	emitOperation(context.Background(), f, sampleSpec())
	if f.calls != 1 {
		t.Errorf("failingEmitter.calls = %d, want 1", f.calls)
	}
}
