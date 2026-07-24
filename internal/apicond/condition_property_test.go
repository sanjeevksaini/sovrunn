package apicond

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

// Deterministic seed for Property 8 reproducibility (F12-STATUS-002/003).
const property8Seed int64 = 20260723

const property8Iterations = 100

var (
	property8Types = []string{
		"Ready",
		"Valid",
		"Available",
		"Progressing",
		"Degraded",
		"Synced",
	}
	property8Reasons = []string{
		"Pending",
		"Succeeded",
		"Failed",
		"Waiting",
		"Reconciling",
		"Unknown",
		"Step42",
	}
	property8Statuses = []ConditionStatus{
		ConditionTrue,
		ConditionFalse,
		ConditionUnknown,
	}
)

type property8Upsert struct {
	Type               string
	Status             ConditionStatus
	Reason             string
	Message            string
	ObservedGeneration int64
}

// Feature: api-resource-naming-status-and-validation-standard, Property 8: Condition transition semantics
//
// For any sequence of condition upserts, SetCondition advances lastTransitionTime
// iff status changed (or the type is new); other conditions are unchanged; status
// stays True/False/Unknown; type/reason stay PascalCase; types never accumulate
// as event history (at most one entry per type).
//
// Validates: Requirements 4.8 (F12-STATUS-002, F12-STATUS-003)
func TestProperty8_ConditionTransitionSemantics(t *testing.T) {
	t.Parallel()

	rng := rand.New(rand.NewSource(property8Seed))
	for i := 0; i < property8Iterations; i++ {
		seq := generateProperty8Sequence(rng)
		if err := checkProperty8Sequence(seq, i); err != nil {
			t.Fatalf("property 8 failed at iteration %d (seed %d): %v", i, property8Seed, err)
		}
	}
}

func generateProperty8Sequence(rng *rand.Rand) []property8Upsert {
	n := 1 + rng.Intn(24) // 1..24 upserts
	seq := make([]property8Upsert, n)
	for i := range seq {
		seq[i] = property8Upsert{
			Type:               property8Types[rng.Intn(len(property8Types))],
			Status:             property8Statuses[rng.Intn(len(property8Statuses))],
			Reason:             property8Reasons[rng.Intn(len(property8Reasons))],
			Message:            fmt.Sprintf("msg-%d-%d", rng.Intn(1000), i),
			ObservedGeneration: int64(rng.Intn(50) + 1),
		}
	}
	return seq
}

func checkProperty8Sequence(seq []property8Upsert, iteration int) error {
	base := time.Date(2026, 7, 23, 0, 0, 0, 0, time.UTC)
	var conds []Condition

	for step, op := range seq {
		now := base.Add(time.Duration(step) * time.Minute)
		nowStr := now.UTC().Format(time.RFC3339)

		prevByType := make(map[string]Condition, len(conds))
		for _, c := range conds {
			prevByType[c.Type] = c
		}
		prevLen := len(conds)

		incoming := Condition{
			Type:               op.Type,
			Status:             op.Status,
			Reason:             op.Reason,
			Message:            op.Message,
			ObservedGeneration: op.ObservedGeneration,
		}
		if !incoming.Valid() {
			return fmt.Errorf("iteration %d step %d: generated upsert must be Valid: %#v", iteration, step, incoming)
		}

		conds = SetCondition(conds, incoming, now)

		// Current facts only: at most one entry per type (F12-STATUS-002).
		seen := make(map[string]struct{}, len(conds))
		for _, c := range conds {
			if _, dup := seen[c.Type]; dup {
				return fmt.Errorf("iteration %d step %d: duplicate type %q (history accumulated): %#v", iteration, step, c.Type, conds)
			}
			seen[c.Type] = struct{}{}

			if !c.Status.Valid() {
				return fmt.Errorf("iteration %d step %d: invalid status %q on %#v", iteration, step, c.Status, c)
			}
			if !IsPascalCase(c.Type) || !IsPascalCase(c.Reason) {
				return fmt.Errorf("iteration %d step %d: type/reason must stay PascalCase: %#v", iteration, step, c)
			}
			if !c.Valid() {
				return fmt.Errorf("iteration %d step %d: stored condition must remain Valid: %#v", iteration, step, c)
			}
		}

		got, ok := GetCondition(conds, op.Type)
		if !ok {
			return fmt.Errorf("iteration %d step %d: upserted type %q missing after SetCondition", iteration, step, op.Type)
		}
		if got.Status != op.Status || got.Reason != op.Reason || got.Message != op.Message || got.ObservedGeneration != op.ObservedGeneration {
			return fmt.Errorf("iteration %d step %d: upserted fields not applied: got %#v want status=%q reason=%q message=%q obs=%d",
				iteration, step, got, op.Status, op.Reason, op.Message, op.ObservedGeneration)
		}

		prev, existed := prevByType[op.Type]
		statusChanged := !existed || prev.Status != op.Status
		if statusChanged {
			if got.LastTransitionTime != nowStr {
				return fmt.Errorf("iteration %d step %d: LastTransitionTime must advance on status change/new: got %q want %q (prev=%#v)",
					iteration, step, got.LastTransitionTime, nowStr, prev)
			}
		} else if got.LastTransitionTime != prev.LastTransitionTime {
			return fmt.Errorf("iteration %d step %d: LastTransitionTime must not advance when status unchanged: got %q want %q",
				iteration, step, got.LastTransitionTime, prev.LastTransitionTime)
		}

		// Other conditions must be byte-for-byte unchanged.
		for typ, before := range prevByType {
			if typ == op.Type {
				continue
			}
			after, ok := GetCondition(conds, typ)
			if !ok {
				return fmt.Errorf("iteration %d step %d: unrelated type %q removed", iteration, step, typ)
			}
			if after != before {
				return fmt.Errorf("iteration %d step %d: unrelated condition changed: before=%#v after=%#v", iteration, step, before, after)
			}
		}

		wantLen := prevLen
		if !existed {
			wantLen++
		}
		if len(conds) != wantLen {
			return fmt.Errorf("iteration %d step %d: len=%d want %d (existed=%v)", iteration, step, len(conds), wantLen, existed)
		}
	}

	return nil
}
