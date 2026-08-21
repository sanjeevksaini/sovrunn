package executiontarget

import (
	"sync/atomic"
	"time"

	"github.com/sanjeevksaini/sovrunn/internal/executiontarget/model"
)

// Registered observer identity and revision (design §4.4; REQ-F16-07).
const (
	RegisteredObserverRevision = "1"
	FactFreshnessWindow        = 60 * time.Second
)

// Named fact identifiers emitted by sovrunn.synthetic-iaas-observer/v1.
const (
	FactNameComputeVM      = "compute.vm"
	FactNameStorageBlock   = "storage.block"
	FactNameStorageObject  = "storage.object"
	FactNameNetworkPrivate = "network.private"
)

// FixtureMode classifies a target-bound synthetic fixture outcome.
type FixtureMode int

const (
	// FixturePresent yields the fixture's declared facts.
	FixturePresent FixtureMode = iota
	// FixtureLogicalTimeout is a fixture-declared logical timeout resolved
	// against the injected clock with no elapsed network wait.
	FixtureLogicalTimeout
	// FixtureFaultDuplicate emits a duplicate named fact for observer-fault tests.
	FixtureFaultDuplicate
	// FixtureFaultMissing omits one required named fact.
	FixtureFaultMissing
	// FixtureFaultMalformed emits malformed provenance or an invalid truth.
	FixtureFaultMalformed
)

// TargetBoundFixture is an in-memory, target-bound observer fixture. It is
// external-effect-free and never triggers network I/O.
type TargetBoundFixture struct {
	// Revision is selected once for the qualification attempt and retained as
	// FactSet provenance (ObserverRevision).
	Revision string
	Mode     FixtureMode
	// Truths are used when Mode is FixturePresent. Zero values are Unknown.
	Truths model.FactSetTruths
	// MissingName selects which required fact to omit under FixtureFaultMissing.
	MissingName string
}

// FixtureSource resolves a target-bound fixture. A miss yields the classified
// missing-fixture outcome (four Unknown facts at now/+60s).
type FixtureSource interface {
	Lookup(targetUID string) (TargetBoundFixture, bool)
}

// MapFixtureSource is a process-local fixture map used by tests and the
// synthetic observer. It is safe for concurrent read after publication.
type MapFixtureSource struct {
	fixtures map[string]TargetBoundFixture
}

// NewMapFixtureSource constructs an empty fixture source.
func NewMapFixtureSource() *MapFixtureSource {
	return &MapFixtureSource{fixtures: make(map[string]TargetBoundFixture)}
}

// Set registers or replaces a target-bound fixture.
func (m *MapFixtureSource) Set(targetUID string, fx TargetBoundFixture) {
	if m.fixtures == nil {
		m.fixtures = make(map[string]TargetBoundFixture)
	}
	m.fixtures[targetUID] = fx
}

// Delete removes a target-bound fixture.
func (m *MapFixtureSource) Delete(targetUID string) {
	delete(m.fixtures, targetUID)
}

// Lookup returns a copy of the target-bound fixture when present.
func (m *MapFixtureSource) Lookup(targetUID string) (TargetBoundFixture, bool) {
	fx, ok := m.fixtures[targetUID]
	return fx, ok
}

// ProposedFact is a single named fact proposal from the synthetic observer.
type ProposedFact struct {
	Name  string
	Truth model.FactTruth
}

// Observation is the observer's fact proposal for one qualification attempt.
// The selected fixture revision is immutable for the attempt.
type Observation struct {
	TargetUID         string
	FixtureRevision   string
	ObserverID        string
	ObserverRevision  string
	FactSchemaVersion string
	ObservedAt        time.Time
	ExpiresAt         time.Time
	Facts             []ProposedFact
	MissingFixture    bool
	LogicalTimeout    bool
}

// SyntheticObserver is sovrunn.synthetic-iaas-observer/v1: a target-bound,
// clock-driven, external-effect-free fixture reader (DD-03).
type SyntheticObserver struct {
	source FixtureSource
	now    func() time.Time
	// externalCallCount remains 0 for every invocation (Risk 2).
	externalCallCount atomic.Int64
}

// NewSyntheticObserver constructs the sole FEATURE-0016 observer.
func NewSyntheticObserver(source FixtureSource, now func() time.Time) *SyntheticObserver {
	if now == nil {
		now = func() time.Time { return time.Now().UTC() }
	}
	if source == nil {
		source = NewMapFixtureSource()
	}
	return &SyntheticObserver{source: source, now: now}
}

// ExternalCallCount returns the number of external/network calls attempted.
// Always zero for this synthetic observer.
func (o *SyntheticObserver) ExternalCallCount() int64 {
	return o.externalCallCount.Load()
}

// Observe reads the target-bound in-memory fixture and proposes facts. It never
// performs network I/O or blocks on a real call. A missing fixture and a
// fixture-declared logical timeout each yield four Unknown facts at now/+60s.
func (o *SyntheticObserver) Observe(targetUID string) Observation {
	// Deliberately no network wait and no external call increment.
	now := o.now().UTC()
	expires := now.Add(FactFreshnessWindow)

	fx, ok := o.source.Lookup(targetUID)
	if !ok {
		return Observation{
			TargetUID:         targetUID,
			FixtureRevision:   RegisteredObserverRevision,
			ObserverID:        model.ObserverID,
			ObserverRevision:  RegisteredObserverRevision,
			FactSchemaVersion: model.FactSchemaVersionV1,
			ObservedAt:        now,
			ExpiresAt:         expires,
			Facts:             fourUnknownFacts(),
			MissingFixture:    true,
		}
	}

	// Capture revision once; fixture replacement after this point cannot affect
	// this observation (immutability for the attempt).
	revision := fx.Revision
	if revision == "" {
		revision = RegisteredObserverRevision
	}

	switch fx.Mode {
	case FixtureLogicalTimeout:
		return Observation{
			TargetUID:         targetUID,
			FixtureRevision:   revision,
			ObserverID:        model.ObserverID,
			ObserverRevision:  revision,
			FactSchemaVersion: model.FactSchemaVersionV1,
			ObservedAt:        now,
			ExpiresAt:         expires,
			Facts:             fourUnknownFacts(),
			LogicalTimeout:    true,
		}
	case FixtureFaultDuplicate:
		facts := factsFromTruths(fx.Truths)
		facts = append(facts, ProposedFact{Name: FactNameComputeVM, Truth: model.FactSupported})
		return Observation{
			TargetUID:         targetUID,
			FixtureRevision:   revision,
			ObserverID:        model.ObserverID,
			ObserverRevision:  revision,
			FactSchemaVersion: model.FactSchemaVersionV1,
			ObservedAt:        now,
			ExpiresAt:         expires,
			Facts:             facts,
		}
	case FixtureFaultMissing:
		missing := fx.MissingName
		if missing == "" {
			missing = FactNameNetworkPrivate
		}
		return Observation{
			TargetUID:         targetUID,
			FixtureRevision:   revision,
			ObserverID:        model.ObserverID,
			ObserverRevision:  revision,
			FactSchemaVersion: model.FactSchemaVersionV1,
			ObservedAt:        now,
			ExpiresAt:         expires,
			Facts:             factsOmitting(fx.Truths, missing),
		}
	case FixtureFaultMalformed:
		return Observation{
			TargetUID:         targetUID,
			FixtureRevision:   revision,
			ObserverID:        model.ObserverID,
			ObserverRevision:  revision,
			FactSchemaVersion: "malformed",
			ObservedAt:        now,
			ExpiresAt:         expires,
			Facts:             factsFromTruths(fx.Truths),
		}
	default: // FixturePresent
		return Observation{
			TargetUID:         targetUID,
			FixtureRevision:   revision,
			ObserverID:        model.ObserverID,
			ObserverRevision:  revision,
			FactSchemaVersion: model.FactSchemaVersionV1,
			ObservedAt:        now,
			ExpiresAt:         expires,
			Facts:             factsFromTruths(fx.Truths),
		}
	}
}

func fourUnknownFacts() []ProposedFact {
	return []ProposedFact{
		{Name: FactNameComputeVM, Truth: model.FactUnknown},
		{Name: FactNameStorageBlock, Truth: model.FactUnknown},
		{Name: FactNameStorageObject, Truth: model.FactUnknown},
		{Name: FactNameNetworkPrivate, Truth: model.FactUnknown},
	}
}

func factsFromTruths(t model.FactSetTruths) []ProposedFact {
	vm := t.Compute.VM
	if vm == "" {
		vm = model.FactUnknown
	}
	block := t.Storage.Block
	if block == "" {
		block = model.FactUnknown
	}
	object := t.Storage.Object
	if object == "" {
		object = model.FactUnknown
	}
	private := t.Network.Private
	if private == "" {
		private = model.FactUnknown
	}
	return []ProposedFact{
		{Name: FactNameComputeVM, Truth: vm},
		{Name: FactNameStorageBlock, Truth: block},
		{Name: FactNameStorageObject, Truth: object},
		{Name: FactNameNetworkPrivate, Truth: private},
	}
}

func factsOmitting(t model.FactSetTruths, omit string) []ProposedFact {
	all := factsFromTruths(t)
	out := make([]ProposedFact, 0, 3)
	for _, f := range all {
		if f.Name == omit {
			continue
		}
		out = append(out, f)
	}
	return out
}
