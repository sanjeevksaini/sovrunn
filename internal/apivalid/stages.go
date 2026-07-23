package apivalid

import (
	"context"

	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
)

// DefaultingStage is the layer-5 deterministic defaulting contract owned by
// apivalid (D-04; F12-VALIDATION-004). Apply returns the object that all
// later layers MUST use. Concrete implementations may live in apivalid or
// apiconform; they own immutable trusted rule configuration and MUST NOT
// accept arbitrary caller-supplied defaulting rules at invocation time.
//
// A non-nil error means defaulting was unavailable or encountered an
// internal fault: the pipeline MUST fail closed at LayerDefaulting with
// Result.Problem = 500 INTERNAL_ERROR, Result.Err = the internal cause,
// and MUST NOT execute layers 6–7 (or later).
type DefaultingStage interface {
	Apply(ctx context.Context, object any) (objectAfterDefaults any, err error)
}

// ValidationStage is the shared contract for layer-6 semantic validation and
// layer-7 structural reference/kind/scope validation (D-04;
// F12-VALIDATION-004, F12-VALIDATION-005). Both Semantic and Reference slots
// in StageSet use this interface. Implementations receive the defaulted
// object from DefaultingStage.Apply.
//
// Return semantics mirror StructuralValidator's violation/error split:
//
//   - err != nil: stage unavailable or internal fault; pipeline MUST stop at
//     the current layer (LayerSemantic or LayerReference), set
//     Result.Problem to 500 INTERNAL_ERROR, set Result.Err to the cause,
//     and MUST NOT execute any later layer.
//   - err == nil, len(violations) > 0: ordinary findings; populate
//     Result.Violations, set FailedAt to the current layer, leave Problem
//     and Err nil, and stop before any later layer executes.
//   - err == nil, len(violations) == 0: stage passed; continue.
//
// Reference-stage construction receives trusted reference constraints and
// allowed scopes; those are not accepted as arbitrary caller input on
// Validate.
type ValidationStage interface {
	Validate(ctx context.Context, object any) ([]apiproblem.Violation, error)
}

// StageSet holds the layer 5–7 stage implementations carried by Input.Stages.
//
// Binding rules (D-04; validation pipeline layers 5–7):
//
//   - Defaulting.Apply returns the object used by all later layers;
//     Semantic.Validate and Reference.Validate receive that defaulted object.
//   - Stage implementations own immutable trusted rule configuration.
//   - Reference-stage construction receives trusted reference constraints
//     and allowed scopes; these are not accepted as arbitrary caller input.
//   - A missing required stage (nil slot where the contract requires one)
//     fails closed at its corresponding layer with Result.Problem = 500
//     INTERNAL_ERROR, Result.Err = internal cause, and no later layer
//     execution. FailedAt MUST be LayerDefaulting, LayerSemantic, or
//     LayerReference as appropriate.
//   - An explicitly constructed deterministic no-op stage is still invoked
//     and is allowed only when the contract declares that layer inapplicable;
//     the pipeline MUST NOT silently skip the invocation.
//   - Stage errors set Result.Problem to 500 INTERNAL_ERROR and Result.Err
//     to the internal cause.
//   - Ordinary semantic/reference findings populate Result.Violations, set
//     FailedAt to the current layer, leave Problem and Err nil, and stop
//     before any later layer executes.
//
// Pipeline invocation rule: a full external-object pipeline MUST NOT
// silently omit a requested layer. A missing required stage implementation
// or stage-internal error fails closed at that layer with Result.Problem =
// 500 INTERNAL_ERROR, Result.Err = internal cause, and no later layer
// execution.
//
// Package boundaries: these interfaces are apivalid-owned; concrete
// implementations may live in apivalid or apiconform. apivalid MUST NOT
// import apischema.
//
// Input (including Stages) and Layer/Result wiring live in pipeline.go
// (task 6.4). Stage invocation into layers 5–7 is owned by task 6.5e;
// ordering/fail-closed tests are tasks 6.7a.
type StageSet struct {
	Defaulting DefaultingStage
	Semantic   ValidationStage
	Reference  ValidationStage
}
