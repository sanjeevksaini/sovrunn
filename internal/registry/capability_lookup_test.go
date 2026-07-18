package registry

import (
	"context"
	"errors"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/resources"
)

func TestCapabilityLookup_NoCapabilitiesReturnsFalse(t *testing.T) {
	lookup := NewCapabilityLookup(NewCapabilityRegistry())

	ok, err := lookup.HasActiveCapabilityForServiceClass(context.Background(), "datastore.postgresql")
	if err != nil {
		t.Fatalf("HasActiveCapabilityForServiceClass() error = %v", err)
	}
	if ok {
		t.Errorf("HasActiveCapabilityForServiceClass() = true, want false")
	}
}

func TestCapabilityLookup_InactivePhaseReturnsFalse(t *testing.T) {
	reg := NewCapabilityRegistry()
	ctx := context.Background()
	c := sampleCapability("postgres-provision", "postgres-basic", "datastore.postgresql")
	c.Status.Phase = resources.PhaseInactive
	if _, err := reg.CreateCapability(ctx, c); err != nil {
		t.Fatalf("CreateCapability() error = %v", err)
	}

	lookup := NewCapabilityLookup(reg)
	ok, err := lookup.HasActiveCapabilityForServiceClass(ctx, "datastore.postgresql")
	if err != nil {
		t.Fatalf("HasActiveCapabilityForServiceClass() error = %v", err)
	}
	if ok {
		t.Errorf("HasActiveCapabilityForServiceClass() = true, want false for inactive phase")
	}
}

func TestCapabilityLookup_UnsupportedReturnsFalse(t *testing.T) {
	reg := NewCapabilityRegistry()
	ctx := context.Background()
	c := sampleCapability("postgres-provision", "postgres-basic", "datastore.postgresql")
	c.Spec.Supported = false
	if _, err := reg.CreateCapability(ctx, c); err != nil {
		t.Fatalf("CreateCapability() error = %v", err)
	}

	lookup := NewCapabilityLookup(reg)
	ok, err := lookup.HasActiveCapabilityForServiceClass(ctx, "datastore.postgresql")
	if err != nil {
		t.Fatalf("HasActiveCapabilityForServiceClass() error = %v", err)
	}
	if ok {
		t.Errorf("HasActiveCapabilityForServiceClass() = true, want false for unsupported")
	}
}

func TestCapabilityLookup_ActiveAndSupportedReturnsTrue(t *testing.T) {
	reg := NewCapabilityRegistry()
	ctx := context.Background()
	c := sampleCapability("postgres-provision", "postgres-basic", "datastore.postgresql")
	if _, err := reg.CreateCapability(ctx, c); err != nil {
		t.Fatalf("CreateCapability() error = %v", err)
	}

	lookup := NewCapabilityLookup(reg)
	ok, err := lookup.HasActiveCapabilityForServiceClass(ctx, "datastore.postgresql")
	if err != nil {
		t.Fatalf("HasActiveCapabilityForServiceClass() error = %v", err)
	}
	if !ok {
		t.Errorf("HasActiveCapabilityForServiceClass() = false, want true")
	}
}

func TestCapabilityLookup_MultipleOnlyOneActiveSupportedReturnsTrue(t *testing.T) {
	reg := NewCapabilityRegistry()
	ctx := context.Background()

	inactive := sampleCapability("postgres-inactive", "postgres-basic", "datastore.postgresql")
	inactive.Status.Phase = resources.PhaseInactive
	unsupported := sampleCapability("postgres-unsupported", "postgres-basic", "datastore.postgresql")
	unsupported.Spec.Supported = false
	active := sampleCapability("postgres-active", "postgres-basic", "datastore.postgresql")

	for _, c := range []resources.Capability{inactive, unsupported, active} {
		if _, err := reg.CreateCapability(ctx, c); err != nil {
			t.Fatalf("CreateCapability(%q) error = %v", c.Metadata.Name, err)
		}
	}

	lookup := NewCapabilityLookup(reg)
	ok, err := lookup.HasActiveCapabilityForServiceClass(ctx, "datastore.postgresql")
	if err != nil {
		t.Fatalf("HasActiveCapabilityForServiceClass() error = %v", err)
	}
	if !ok {
		t.Errorf("HasActiveCapabilityForServiceClass() = false, want true when one is active+supported")
	}
}

func TestCapabilityLookup_PropagatesRegistryError(t *testing.T) {
	wantErr := errors.New("list failed")
	lookup := NewCapabilityLookup(failingCapabilityListRegistry{err: wantErr})

	ok, err := lookup.HasActiveCapabilityForServiceClass(context.Background(), "datastore.postgresql")
	if !errors.Is(err, wantErr) {
		t.Fatalf("error = %v, want %v", err, wantErr)
	}
	if ok {
		t.Errorf("HasActiveCapabilityForServiceClass() = true, want false on error")
	}
}

type failingCapabilityListRegistry struct {
	CapabilityRegistryIface
	err error
}

func (r failingCapabilityListRegistry) ListCapabilities(
	ctx context.Context, pluginRef, serviceClassRef string,
) ([]resources.Capability, error) {
	return nil, r.err
}
