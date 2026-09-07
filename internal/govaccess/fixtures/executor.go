package fixtures

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
)

// ExecutionResult is the deterministic synthetic executor output.
type ExecutionResult struct {
	DeterministicDigest string
	ResourceCount       int
}

// FakeExecutor is a deterministic in-memory executor with zero external I/O.
type FakeExecutor struct {
	externalCallCount int
}

// NewFakeExecutor constructs a deterministic fake executor.
func NewFakeExecutor() *FakeExecutor {
	return &FakeExecutor{}
}

// ExternalCallCount returns the external side-effect sentinel.
func (e *FakeExecutor) ExternalCallCount() int {
	if e == nil {
		return 0
	}
	return e.externalCallCount
}

// Execute validates fixtures before use, performs zero external I/O, and
// returns deterministic results from canonicalized in-memory data only.
func (e *FakeExecutor) Execute(set FixtureSet) (ExecutionResult, error) {
	if err := Validate(set); err != nil {
		return ExecutionResult{}, err
	}
	canonical := canonicalize(set)
	sum := sha256.Sum256([]byte(canonical))
	return ExecutionResult{
		DeterministicDigest: hex.EncodeToString(sum[:]),
		ResourceCount:       len(set.Resources),
	}, nil
}

func canonicalize(set FixtureSet) string {
	resources := make([]Resource, 0, len(set.Resources))
	for _, r := range set.Resources {
		cp := copyResource(r)
		sort.Strings(cp.RefUIDs)
		resources = append(resources, cp)
	}
	sort.Slice(resources, func(i, j int) bool {
		if resources[i].Kind != resources[j].Kind {
			return resources[i].Kind < resources[j].Kind
		}
		return resources[i].UID < resources[j].UID
	})
	lines := make([]string, 0, len(resources))
	for _, r := range resources {
		lines = append(lines, fmt.Sprintf("%s|%s|%s", r.Kind, r.UID, strings.Join(r.RefUIDs, ",")))
	}
	return strings.Join(lines, "\n")
}
