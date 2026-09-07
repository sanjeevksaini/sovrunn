package idempotency

import (
	"time"

	"github.com/sanjeevksaini/sovrunn/internal/govaccess/operation"
)

const completedRetentionHorizon = 24 * time.Hour

// Clock is the deterministic injected time source for retention computation.
type Clock interface {
	NowUTC() time.Time
}

// StageCompleted constructs Task-2 PreparedCompletedResult over neutral
// operation.ResultMaterial and operation.CompletionBinding.
func StageCompleted(
	lookupKey IdempotencyLookupKey,
	binding RequestBinding,
	result operation.ResultMaterial,
	completion operation.CompletionBinding,
) (PreparedCompletedResult, operation.MechanicalFailure) {
	return NewPreparedCompletedResult(lookupKey, binding, result, completion)
}

// CompleteLifecycle transitions prepared result to completed record and returns
// the 24h retention horizon from the injected clock.
func CompleteLifecycle(
	prepared PreparedCompletedResult,
	clk Clock,
) (CompletedRecord, time.Time, operation.MechanicalFailure) {
	if !prepared.Sealed() {
		return CompletedRecord{}, time.Time{}, operation.NewMechanicalFailure("lifecycle_prepared_unsealed")
	}
	if clk == nil {
		return CompletedRecord{}, time.Time{}, operation.NewMechanicalFailure("lifecycle_clock_nil")
	}
	now := clk.NowUTC()
	if now.IsZero() {
		return CompletedRecord{}, time.Time{}, operation.NewMechanicalFailure("lifecycle_clock_zero")
	}
	retainUntil := now.UTC().Add(completedRetentionHorizon)
	record, fail := NewCompletedRecord(prepared.LookupKey(), prepared.Binding(), prepared.Result(), prepared.Completion())
	if fail.Reason() != "" {
		return CompletedRecord{}, time.Time{}, fail
	}
	return record, retainUntil, operation.MechanicalFailure{}
}

// CompletedRetentionHorizon exposes the fixed retention window for tests.
func CompletedRetentionHorizon() time.Duration {
	return completedRetentionHorizon
}
