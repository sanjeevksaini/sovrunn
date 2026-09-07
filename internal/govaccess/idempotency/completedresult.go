package idempotency

import (
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/operation"
)

// PreparedCompletedResult is a sealed value binding a completed caller result
// to its lookup key and full RequestBinding. Structural sealing only; the
// lifecycle constructor PrepareCompletedResult is Task-6-owned (DD-11).
type PreparedCompletedResult struct {
	sealed     bool
	lookupKey  IdempotencyLookupKey
	binding    RequestBinding
	result     operation.ResultMaterial
	completion operation.CompletionBinding
}

// NewPreparedCompletedResult seals a structurally valid prepared completed
// result. It performs no reservation-table or UOW coordination.
func NewPreparedCompletedResult(
	lookupKey IdempotencyLookupKey,
	binding RequestBinding,
	result operation.ResultMaterial,
	completion operation.CompletionBinding,
) (PreparedCompletedResult, operation.MechanicalFailure) {
	if !lookupKey.sealed {
		return PreparedCompletedResult{}, operation.NewMechanicalFailure("prepared_completed_result_lookup_unsealed")
	}
	if !binding.sealed {
		return PreparedCompletedResult{}, operation.NewMechanicalFailure("prepared_completed_result_binding_unsealed")
	}
	if !result.Sealed() {
		return PreparedCompletedResult{}, operation.NewMechanicalFailure("prepared_completed_result_material_unsealed")
	}
	if !completion.Sealed() {
		return PreparedCompletedResult{}, operation.NewMechanicalFailure("prepared_completed_result_completion_unsealed")
	}
	return PreparedCompletedResult{
		sealed:     true,
		lookupKey:  lookupKey,
		binding:    binding,
		result:     result,
		completion: completion,
	}, operation.MechanicalFailure{}
}

// Sealed reports whether the value was constructed through NewPreparedCompletedResult.
func (r PreparedCompletedResult) Sealed() bool { return r.sealed }

// LookupKey returns the sealed lookup key.
func (r PreparedCompletedResult) LookupKey() IdempotencyLookupKey { return r.lookupKey }

// Binding returns the full sealed RequestBinding.
func (r PreparedCompletedResult) Binding() RequestBinding { return r.binding }

// Result returns the sealed ResultMaterial.
func (r PreparedCompletedResult) Result() operation.ResultMaterial { return r.result }

// Completion returns the sealed CompletionBinding.
func (r PreparedCompletedResult) Completion() operation.CompletionBinding { return r.completion }

// CompletedRecord is a sealed committed completed-result value stored under
// IdempotencyLookupKey and carrying the full RequestBinding (DD-11).
type CompletedRecord struct {
	sealed     bool
	lookupKey  IdempotencyLookupKey
	binding    RequestBinding
	result     operation.ResultMaterial
	completion operation.CompletionBinding
}

// NewCompletedRecord seals a structurally valid completed record.
func NewCompletedRecord(
	lookupKey IdempotencyLookupKey,
	binding RequestBinding,
	result operation.ResultMaterial,
	completion operation.CompletionBinding,
) (CompletedRecord, operation.MechanicalFailure) {
	if !lookupKey.sealed {
		return CompletedRecord{}, operation.NewMechanicalFailure("completed_record_lookup_unsealed")
	}
	if !binding.sealed {
		return CompletedRecord{}, operation.NewMechanicalFailure("completed_record_binding_unsealed")
	}
	if !result.Sealed() {
		return CompletedRecord{}, operation.NewMechanicalFailure("completed_record_material_unsealed")
	}
	if !completion.Sealed() {
		return CompletedRecord{}, operation.NewMechanicalFailure("completed_record_completion_unsealed")
	}
	return CompletedRecord{
		sealed:     true,
		lookupKey:  lookupKey,
		binding:    binding,
		result:     result,
		completion: completion,
	}, operation.MechanicalFailure{}
}

// Sealed reports whether the record was constructed through NewCompletedRecord.
func (r CompletedRecord) Sealed() bool { return r.sealed }

// LookupKey returns the sealed lookup key.
func (r CompletedRecord) LookupKey() IdempotencyLookupKey { return r.lookupKey }

// Binding returns the full sealed RequestBinding.
func (r CompletedRecord) Binding() RequestBinding { return r.binding }

// Result returns the sealed ResultMaterial.
func (r CompletedRecord) Result() operation.ResultMaterial { return r.result }

// Completion returns the sealed CompletionBinding.
func (r CompletedRecord) Completion() operation.CompletionBinding { return r.completion }
