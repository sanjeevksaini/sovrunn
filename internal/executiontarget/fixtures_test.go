package executiontarget

import (
	"testing"
	"time"

	"github.com/sanjeevksaini/sovrunn/internal/executiontarget/model"
)

// Test-only target-bound fixtures for sovrunn.synthetic-iaas-observer/v1.

func allSupportedTruths() model.FactSetTruths {
	return model.FactSetTruths{
		Compute: model.FactComputeTruths{VM: model.FactSupported},
		Storage: model.FactStorageTruths{
			Block:  model.FactSupported,
			Object: model.FactSupported,
		},
		Network: model.FactNetworkTruths{Private: model.FactSupported},
	}
}

func unsupportedVMTruths() model.FactSetTruths {
	t := allSupportedTruths()
	t.Compute.VM = model.FactUnsupported
	return t
}

func unknownBlockTruths() model.FactSetTruths {
	t := allSupportedTruths()
	t.Storage.Block = model.FactUnknown
	return t
}

func presentFixture(revision string, truths model.FactSetTruths) TargetBoundFixture {
	return TargetBoundFixture{
		Revision: revision,
		Mode:     FixturePresent,
		Truths:   truths,
	}
}

func logicalTimeoutFixture(revision string) TargetBoundFixture {
	return TargetBoundFixture{
		Revision: revision,
		Mode:     FixtureLogicalTimeout,
	}
}

func duplicateFaultFixture(revision string) TargetBoundFixture {
	return TargetBoundFixture{
		Revision: revision,
		Mode:     FixtureFaultDuplicate,
		Truths:   allSupportedTruths(),
	}
}

func missingFaultFixture(revision, missing string) TargetBoundFixture {
	return TargetBoundFixture{
		Revision:    revision,
		Mode:        FixtureFaultMissing,
		Truths:      allSupportedTruths(),
		MissingName: missing,
	}
}

func malformedFaultFixture(revision string) TargetBoundFixture {
	return TargetBoundFixture{
		Revision: revision,
		Mode:     FixtureFaultMalformed,
		Truths:   allSupportedTruths(),
	}
}

// blockingFixtureSource blocks Lookup until release is closed. Used to hold
// observation open while Retire-wins commits.
type blockingFixtureSource struct {
	inner   FixtureSource
	started chan struct{}
	release chan struct{}
}

func (b *blockingFixtureSource) Lookup(targetUID string) (TargetBoundFixture, bool) {
	// Snapshot before blocking so mid-flight fixture replacement cannot change
	// the selected revision for this observation attempt.
	var (
		fx TargetBoundFixture
		ok bool
	)
	if b.inner != nil {
		fx, ok = b.inner.Lookup(targetUID)
	}
	select {
	case <-b.started:
	default:
		close(b.started)
	}
	<-b.release
	return fx, ok
}

func fixedObserverNow(t time.Time) func() time.Time {
	return func() time.Time { return t.UTC() }
}

func assertFourUnknown(t *testing.T, facts []ProposedFact) {
	t.Helper()
	if len(facts) != 4 {
		t.Fatalf("want 4 facts, got %d", len(facts))
	}
	want := map[string]model.FactTruth{
		FactNameComputeVM:      model.FactUnknown,
		FactNameStorageBlock:   model.FactUnknown,
		FactNameStorageObject:  model.FactUnknown,
		FactNameNetworkPrivate: model.FactUnknown,
	}
	for _, f := range facts {
		if want[f.Name] != f.Truth {
			t.Fatalf("fact %s=%s, want Unknown", f.Name, f.Truth)
		}
		delete(want, f.Name)
	}
	if len(want) != 0 {
		t.Fatalf("missing facts: %#v", want)
	}
}
