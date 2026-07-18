package registry

import (
	"context"
	"errors"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/resources"
)

func TestCapabilityRegistry_CountByPluginAcrossSeededCapabilities(t *testing.T) {
	reg := NewCapabilityRegistry()
	ctx := context.Background()
	for _, capability := range []resources.Capability{
		sampleCapability("postgres-provision", "postgres-basic", "datastore.postgresql"),
		sampleCapability("postgres-backup", "postgres-basic", "datastore.postgresql"),
		sampleCapability("redis-provision", "redis-basic", "cache.redis"),
	} {
		if _, err := reg.CreateCapability(ctx, capability); err != nil {
			t.Fatalf("CreateCapability(%q) error = %v", capability.Metadata.Name, err)
		}
	}

	for _, test := range []struct {
		pluginName string
		want       int
	}{
		{pluginName: "postgres-basic", want: 2},
		{pluginName: "redis-basic", want: 1},
		{pluginName: "missing", want: 0},
	} {
		count, err := reg.CountByPlugin(ctx, test.pluginName)
		if err != nil {
			t.Fatalf("CountByPlugin(%q) error = %v", test.pluginName, err)
		}
		if count != test.want {
			t.Errorf("CountByPlugin(%q) = %d, want %d", test.pluginName, count, test.want)
		}
	}
}

func TestCapabilityChildBlockerChecker_BlocksWithCapabilityCount(t *testing.T) {
	reg := NewCapabilityRegistry()
	ctx := context.Background()
	for _, capability := range []resources.Capability{
		sampleCapability("postgres-provision", "postgres-basic", "datastore.postgresql"),
		sampleCapability("postgres-backup", "postgres-basic", "datastore.postgresql"),
		sampleCapability("postgres-restore", "postgres-basic", "datastore.postgresql"),
		sampleCapability("redis-provision", "redis-basic", "cache.redis"),
	} {
		if _, err := reg.CreateCapability(ctx, capability); err != nil {
			t.Fatalf("CreateCapability(%q) error = %v", capability.Metadata.Name, err)
		}
	}

	blocker := NewCapabilityChildBlockerChecker(reg)
	blockers, err := blocker.BlockedByPluginChildren(ctx, "postgres-basic")
	if err != nil {
		t.Fatalf("BlockedByPluginChildren() error = %v", err)
	}
	if len(blockers) != 1 {
		t.Fatalf("blockers = %+v, want exactly one entry", blockers)
	}
	if blockers[0].Kind != resources.CapabilityKind {
		t.Errorf("Kind = %q, want %q", blockers[0].Kind, resources.CapabilityKind)
	}
	if blockers[0].Count != 3 {
		t.Errorf("Count = %d, want 3", blockers[0].Count)
	}
}

func TestCapabilityChildBlockerChecker_ReturnsNilWithNoChildren(t *testing.T) {
	blocker := NewCapabilityChildBlockerChecker(NewCapabilityRegistry())

	blockers, err := blocker.BlockedByPluginChildren(context.Background(), "postgres-basic")
	if err != nil {
		t.Fatalf("BlockedByPluginChildren() error = %v", err)
	}
	if blockers != nil {
		t.Errorf("blockers = %+v, want nil", blockers)
	}
}

func TestCapabilityChildBlockerChecker_PropagatesRegistryError(t *testing.T) {
	wantErr := errors.New("count failed")
	blocker := NewCapabilityChildBlockerChecker(failingCapabilityCountRegistry{err: wantErr})

	blockers, err := blocker.BlockedByPluginChildren(context.Background(), "postgres-basic")
	if !errors.Is(err, wantErr) {
		t.Fatalf("error = %v, want %v", err, wantErr)
	}
	if blockers != nil {
		t.Errorf("blockers = %+v, want nil", blockers)
	}
}

type failingCapabilityCountRegistry struct {
	CapabilityRegistryIface
	err error
}

func (r failingCapabilityCountRegistry) CountByPlugin(
	ctx context.Context, pluginName string,
) (int, error) {
	return 0, r.err
}
