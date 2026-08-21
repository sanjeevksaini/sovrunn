package executiontarget

import (
	"testing"
	"time"

	"github.com/sanjeevksaini/sovrunn/internal/executiontarget/model"
)

func TestSyntheticObserver_PresentFactsAndProvenance(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 8, 21, 12, 0, 0, 0, time.UTC)
	src := NewMapFixtureSource()
	src.Set("target-1", presentFixture("fx-rev-9", allSupportedTruths()))
	obs := NewSyntheticObserver(src, fixedObserverNow(now))

	got := obs.Observe("target-1")
	if got.MissingFixture || got.LogicalTimeout {
		t.Fatalf("unexpected classified outcome: %#v", got)
	}
	if got.ObserverID != model.ObserverID {
		t.Fatalf("observerID=%q", got.ObserverID)
	}
	if got.FactSchemaVersion != model.FactSchemaVersionV1 {
		t.Fatalf("factSchemaVersion=%q", got.FactSchemaVersion)
	}
	if got.FixtureRevision != "fx-rev-9" || got.ObserverRevision != "fx-rev-9" {
		t.Fatalf("fixture revision not retained: %#v", got)
	}
	if !got.ObservedAt.Equal(now) || !got.ExpiresAt.Equal(now.Add(FactFreshnessWindow)) {
		t.Fatalf("timestamps: observed=%v expires=%v", got.ObservedAt, got.ExpiresAt)
	}
	if len(got.Facts) != 4 {
		t.Fatalf("want 4 facts, got %d", len(got.Facts))
	}
	for _, f := range got.Facts {
		if f.Truth != model.FactSupported {
			t.Fatalf("fact %s=%s, want Supported", f.Name, f.Truth)
		}
	}
	if obs.ExternalCallCount() != 0 {
		t.Fatalf("externalCallCount=%d, want 0", obs.ExternalCallCount())
	}
}

func TestSyntheticObserver_MissingFixture_FourUnknownNoWait(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 8, 21, 12, 0, 0, 0, time.UTC)
	obs := NewSyntheticObserver(NewMapFixtureSource(), fixedObserverNow(now))

	start := time.Now()
	got := obs.Observe("missing-target")
	elapsed := time.Since(start)
	if elapsed > 50*time.Millisecond {
		t.Fatalf("missing fixture must not network-wait; elapsed=%s", elapsed)
	}
	if !got.MissingFixture {
		t.Fatal("want MissingFixture")
	}
	if !got.ObservedAt.Equal(now) || !got.ExpiresAt.Equal(now.Add(60*time.Second)) {
		t.Fatalf("want now/+60s, got %v / %v", got.ObservedAt, got.ExpiresAt)
	}
	assertFourUnknown(t, got.Facts)
	if obs.ExternalCallCount() != 0 {
		t.Fatalf("externalCallCount=%d", obs.ExternalCallCount())
	}
}

func TestSyntheticObserver_LogicalTimeout_FourUnknownNoWait(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 8, 21, 12, 0, 0, 0, time.UTC)
	src := NewMapFixtureSource()
	src.Set("target-1", logicalTimeoutFixture("fx-timeout"))
	obs := NewSyntheticObserver(src, fixedObserverNow(now))

	start := time.Now()
	got := obs.Observe("target-1")
	elapsed := time.Since(start)
	if elapsed > 50*time.Millisecond {
		t.Fatalf("logical timeout must not network-wait; elapsed=%s", elapsed)
	}
	if !got.LogicalTimeout {
		t.Fatal("want LogicalTimeout")
	}
	if !got.ObservedAt.Equal(now) || !got.ExpiresAt.Equal(now.Add(60*time.Second)) {
		t.Fatalf("want now/+60s, got %v / %v", got.ObservedAt, got.ExpiresAt)
	}
	assertFourUnknown(t, got.Facts)
	if got.FixtureRevision != "fx-timeout" {
		t.Fatalf("fixture revision=%q", got.FixtureRevision)
	}
	if obs.ExternalCallCount() != 0 {
		t.Fatalf("externalCallCount=%d", obs.ExternalCallCount())
	}
}

func TestSyntheticObserver_FixtureRevisionImmutableDuringObserve(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 8, 21, 12, 0, 0, 0, time.UTC)
	src := NewMapFixtureSource()
	src.Set("target-1", presentFixture("rev-a", allSupportedTruths()))

	started := make(chan struct{})
	release := make(chan struct{})
	blocking := &blockingFixtureSource{inner: src, started: started, release: release}
	obs := NewSyntheticObserver(blocking, fixedObserverNow(now))

	done := make(chan Observation, 1)
	go func() {
		done <- obs.Observe("target-1")
	}()
	<-started
	// Replace fixture while observation is in Lookup; captured revision must
	// remain the value returned by this Lookup (rev-a from inner at unblock).
	src.Set("target-1", presentFixture("rev-b", unsupportedVMTruths()))
	close(release)
	got := <-done
	if got.FixtureRevision != "rev-a" {
		t.Fatalf("selected fixture revision mutated mid-observe: %q", got.FixtureRevision)
	}
	if got.Facts[0].Truth != model.FactSupported {
		t.Fatalf("observation must retain selected fixture truths, got %#v", got.Facts)
	}
}

func TestSyntheticObserver_FaultModes(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 8, 21, 12, 0, 0, 0, time.UTC)
	src := NewMapFixtureSource()
	obs := NewSyntheticObserver(src, fixedObserverNow(now))

	src.Set("dup", duplicateFaultFixture("r1"))
	dup := obs.Observe("dup")
	if len(dup.Facts) != 5 {
		t.Fatalf("duplicate fixture should emit 5 facts, got %d", len(dup.Facts))
	}

	src.Set("miss", missingFaultFixture("r2", FactNameStorageObject))
	miss := obs.Observe("miss")
	if len(miss.Facts) != 3 {
		t.Fatalf("missing fixture should emit 3 facts, got %d", len(miss.Facts))
	}

	src.Set("bad", malformedFaultFixture("r3"))
	bad := obs.Observe("bad")
	if bad.FactSchemaVersion == model.FactSchemaVersionV1 {
		t.Fatal("malformed fixture must not emit factSchemaVersion=v1")
	}

	if obs.ExternalCallCount() != 0 {
		t.Fatalf("externalCallCount=%d", obs.ExternalCallCount())
	}
}
