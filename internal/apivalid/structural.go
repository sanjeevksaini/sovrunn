package apivalid

import "github.com/sanjeevksaini/sovrunn/internal/apiproblem"

// StructuralValidator is the layer-4 structural validation contract owned
// by apivalid. Implementations live in apiconform (using apischema) so
// apivalid never imports apischema directly.
//
// The Validate method returns both violations and an error. The error
// signals validator unavailability or configuration failure (e.g. schema
// not found, registry misconfiguration). Violations represent ordinary
// schema-mismatch findings. This separation makes failure distinguishable
// from clean validation:
//
//   - err != nil: structural validation was UNAVAILABLE; the pipeline MUST
//     stop at LayerStructural, set Result.Problem to a 500 INTERNAL_ERROR
//     Problem, set Result.Err to the internal cause, return no success
//     result, and MUST NOT execute layers 5 through 7.
//   - err == nil, len(violations) > 0: ordinary schema violations (422).
//   - err == nil, len(violations) == 0: instance is structurally valid.
//
// A nil StructuralValidator in Input at layer 4 is treated identically to
// a non-nil validator returning an error: the pipeline stops, sets
// Result.Problem to 500 INTERNAL_ERROR, sets Result.Err, and never
// executes layers 5 through 7.
//
// Primitive unit tests MAY call individual primitive functions directly or
// inject a deterministic stub StructuralValidator. No full pipeline
// invocation for an external object may omit the validator.
//
// Pipeline Result/Problem wiring lives in pipeline.go (task 6.4).
// Expanded structural fail-closed tests are task 6.6.
type StructuralValidator interface {
	Validate(instance any, schemaID string) ([]apiproblem.Violation, error)
}
