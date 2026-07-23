package apivalid

import (
	"context"
	"errors"
	"fmt"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/apiref"
)

// Layer is one ordered validation stage (F12-VALIDATION-001, D-04).
type Layer int

const (
	LayerHTTPContent       Layer = iota + 1 // 1 HTTP/content/size
	LayerDecode                             // 2 safe decode
	LayerFieldHygiene                       // 3 duplicate/unknown-field rejection
	LayerStructural                         // 4 structural schema validation
	LayerDefaulting                         // 5 deterministic defaulting
	LayerSemantic                           // 6 semantic validation
	LayerReference                          // 7 STRUCTURAL reference/kind/scope validation
	LayerAuthorization                      // 8 caller-specific authz + no-existence-disclosure
	LayerCapabilityRuntime                  // 9 later-feature runtime — RESERVED
)

// String returns a stable label for the layer.
func (l Layer) String() string {
	switch l {
	case LayerHTTPContent:
		return "LayerHTTPContent"
	case LayerDecode:
		return "LayerDecode"
	case LayerFieldHygiene:
		return "LayerFieldHygiene"
	case LayerStructural:
		return "LayerStructural"
	case LayerDefaulting:
		return "LayerDefaulting"
	case LayerSemantic:
		return "LayerSemantic"
	case LayerReference:
		return "LayerReference"
	case LayerAuthorization:
		return "LayerAuthorization"
	case LayerCapabilityRuntime:
		return "LayerCapabilityRuntime"
	default:
		return fmt.Sprintf("Layer(%d)", int(l))
	}
}

// Result keeps failure classes distinguishable (F12-VALIDATION-004).
//
// Problem is the safe client-facing failure contract (e.g. 500 INTERNAL_ERROR
// when a StructuralValidator is nil/unavailable or layer-8 is misconfigured).
// Err is internal diagnostic context and MUST NOT be serialized or exposed
// to callers. Secrets, credentials, tokens, and inaccessible resource
// details MUST NOT appear in Problem or Err messages.
//
// Binding rules:
//   - On a nil StructuralValidator or validator error at layer 4:
//     FailedAt = LayerStructural; Problem = 500 INTERNAL_ERROR;
//     Err = internal cause; Violations empty; layers 5–7 do not execute.
//   - Ordinary structural (or later-layer) violations populate Violations,
//     set FailedAt, leave Problem and Err nil (422 mapping is adopter-owned).
//   - Successful validation: Violations empty, Problem nil, Err nil,
//     FailedAt zero.
type Result struct {
	Violations []apiproblem.Violation
	FailedAt   Layer
	Problem    *apiproblem.Problem // safe client-facing failure; nil on success or ordinary violations
	Err        error               // internal diagnostic; MUST NOT be serialized or exposed
}

// Input carries everything the pipeline requires for a single validation
// invocation (D-04, D-17).
//
// StructuralValidator MUST be non-nil for any pipeline run that processes an
// external object; a nil validator at layer 4 causes the pipeline to stop
// with Result.Problem set to 500 INTERNAL_ERROR and Result.Err recording
// the internal cause.
//
// Stages holds Defaulting (layer 5), Semantic (layer 6), and Reference
// (layer 7). After structural success, Validate invokes each stage in
// order. A nil required stage or stage-internal error fails closed at that
// layer with 500 INTERNAL_ERROR. Ordinary semantic/reference violations
// populate Result.Violations, set FailedAt, leave Problem/Err nil, and
// stop before later layers. Deterministic no-op stages are still invoked
// and must not be silently skipped (D-04; task 6.5e).
//
// LAYER-8 CONFIGURATION MATRIX:
//
// When OperationScope is non-nil, exactly one of Path A or Path B MUST be
// completely configured. Configuring both TargetScope and TargetScopeResolver
// is invalid. Configuring neither is invalid. Missing TargetRef is invalid.
// Missing Caller is invalid. An incomplete path configuration is invalid.
// Every invalid layer-8 configuration stops at LayerAuthorization with
// Result.FailedAt = LayerAuthorization, Result.Problem = 500 INTERNAL_ERROR,
// Result.Err = internal configuration cause, no target lookup, no success
// result, and no silent skip.
//
// Path A — authoritative target scope derivable without lookup:
//
//	Required: OperationScope, TargetRef, TargetScope, Authorizer, Caller.
//	ScopeAuthorizer runs before any target lookup. Denial maps through
//	SafeDenial. Allow permits CheckOperationTargetScopeMatch.
//
// Path B — target scope requires authorized lookup:
//
//	Required: OperationScope, TargetRef, TargetScopeResolver, Caller.
//	available=false maps through SafeDenial. available=true permits
//	CheckOperationTargetScopeMatch.
//
// Generic non-Operation validation (OperationScope is nil and no layer-8
// capability was requested) MAY stop successfully after layer 7.
type Input struct {
	Data       []byte
	Mode       DecodeMode
	SchemaID   string
	Validator  StructuralValidator // required; nil = unavailable = 500
	Authorizer ScopeAuthorizer     // required for Path A; nil when using Path B only
	Caller     *CallerContext      // required when OperationScope is non-nil
	Dst        any                 // decode target / structural instance

	// Stages holds Defaulting (layer 5), Semantic (layer 6), and Reference
	// (layer 7). Required for full external-object validation after
	// structural success; nil slots fail closed at the corresponding layer.
	Stages StageSet

	// Operation target-scope equality (layer 8).
	OperationScope      *apimeta.ScopeIdentity        // Operation's canonical scope for the equality check
	TargetRef           *apiref.TypedRef              // Operation.targetRef; required when OperationScope is non-nil
	TargetScope         *apimeta.ScopeIdentity        // pre-derived authoritative target scope (Path A)
	TargetScopeResolver AuthorizedTargetScopeResolver // combined authorized resolver (Path B)
}

// Validate runs the nine-layer validation pipeline (D-04, D-17,
// F12-VALIDATION-001/004/005, F12-SCOPE-002).
//
// Layer behavior:
//
//	1–3: When len(Data) > 0, decode Data into Dst with PolicyFor(Mode) and
//	   lim (size, safe decode, field hygiene). When Data is empty, Dst is
//	   treated as an already-decoded instance (offline / stub path).
//	4: StructuralValidator fail-closed (nil/error → 500; ordinary
//	   violations → FailedAt LayerStructural; success continues).
//	5–7: Input.Stages — Defaulting, Semantic, Reference in order; nil
//	   required stage or stage error → 500 at that layer; ordinary
//	   violations stop before later layers; no-op stages are invoked.
//	8: Operation target-scope equality via the configuration matrix.
//	9: Reserved; never executed.
//
// Request bodies, field values, secrets, credentials, tokens, private keys,
// and inaccessible resource details are not logged here.
func Validate(ctx context.Context, in Input, lim Limits) Result {
	if err := ctx.Err(); err != nil {
		return internalFailure(LayerDecode, err)
	}

	// Layers 1–3: optional decode when request bytes are provided.
	if len(in.Data) > 0 {
		if lim.MaxObjectBytes > 0 && len(in.Data) > lim.MaxObjectBytes {
			return Result{
				FailedAt: LayerHTTPContent,
				Problem:  apiproblem.New(apiproblem.CodeRequestTooLarge),
			}
		}
		if in.Dst == nil {
			return internalFailure(LayerDecode, errors.New("nil decode destination"))
		}
		if prob := DecodeJSON(in.Data, lim, PolicyFor(in.Mode), in.Dst); prob != nil {
			return Result{
				FailedAt:   decodeFailedAt(prob),
				Problem:    prob,
				Violations: append([]apiproblem.Violation(nil), prob.Violations...),
			}
		}
	}

	// Layer 4: structural schema validation (fail-closed).
	if res, stop := runStructural(in); stop {
		return res
	}

	// Layers 5–7: defaulting → semantic → reference (D-04; task 6.5e).
	if res, stop := runStages(ctx, in); stop {
		return res
	}

	// Layer 8: Operation target-scope equality (configuration matrix).
	if in.OperationScope != nil {
		return runLayer8(ctx, in)
	}

	// Generic non-Operation validation stops successfully after layer 7.
	// Layer 9 is reserved.
	return Result{}
}

// runStages invokes layers 5–7 from Input.Stages after structural success.
// Defaulting.Apply returns the object passed to Semantic and Reference.
// A nil required stage or stage error fails closed (500). Ordinary
// violations stop before later layers. Stages are never silently skipped.
func runStages(ctx context.Context, in Input) (Result, bool) {
	if err := ctx.Err(); err != nil {
		return internalFailure(LayerDefaulting, err), true
	}

	// Layer 5: deterministic defaulting.
	if in.Stages.Defaulting == nil {
		return internalFailure(LayerDefaulting, errors.New("nil defaulting stage")), true
	}
	defaulted, err := in.Stages.Defaulting.Apply(ctx, in.Dst)
	if err != nil {
		return internalFailure(LayerDefaulting, err), true
	}

	// Layer 6: semantic validation on the defaulted object.
	if in.Stages.Semantic == nil {
		return internalFailure(LayerSemantic, errors.New("nil semantic stage")), true
	}
	violations, err := in.Stages.Semantic.Validate(ctx, defaulted)
	if err != nil {
		return internalFailure(LayerSemantic, err), true
	}
	if len(violations) > 0 {
		return Result{
			Violations: append([]apiproblem.Violation(nil), violations...),
			FailedAt:   LayerSemantic,
		}, true
	}

	// Layer 7: structural reference/kind/scope validation on the defaulted object.
	if in.Stages.Reference == nil {
		return internalFailure(LayerReference, errors.New("nil reference stage")), true
	}
	violations, err = in.Stages.Reference.Validate(ctx, defaulted)
	if err != nil {
		return internalFailure(LayerReference, err), true
	}
	if len(violations) > 0 {
		return Result{
			Violations: append([]apiproblem.Violation(nil), violations...),
			FailedAt:   LayerReference,
		}, true
	}

	return Result{}, false
}

func runStructural(in Input) (Result, bool) {
	if in.Validator == nil {
		return Result{
			FailedAt: LayerStructural,
			Problem:  apiproblem.New(apiproblem.CodeInternalError),
			Err:      errors.New("nil validator"),
		}, true
	}

	violations, err := in.Validator.Validate(in.Dst, in.SchemaID)
	if err != nil {
		return Result{
			FailedAt: LayerStructural,
			Problem:  apiproblem.New(apiproblem.CodeInternalError),
			Err:      err,
		}, true
	}
	if len(violations) > 0 {
		return Result{
			Violations: append([]apiproblem.Violation(nil), violations...),
			FailedAt:   LayerStructural,
		}, true
	}
	return Result{}, false
}

func runLayer8(ctx context.Context, in Input) Result {
	// Configuration checks first: no Authorize / ResolveAuthorizedTargetScope
	// until the matrix is valid (no silent skip, no target lookup).
	if in.TargetRef == nil {
		return internalFailure(LayerAuthorization, errors.New("layer-8: missing TargetRef"))
	}
	if in.Caller == nil {
		return internalFailure(LayerAuthorization, errors.New("layer-8: missing Caller"))
	}

	hasTargetScope := in.TargetScope != nil
	hasResolver := in.TargetScopeResolver != nil

	switch {
	case !hasTargetScope && !hasResolver:
		return internalFailure(LayerAuthorization, errors.New("layer-8: neither TargetScope nor TargetScopeResolver configured"))
	case hasTargetScope && hasResolver:
		return internalFailure(LayerAuthorization, errors.New("layer-8: both TargetScope and TargetScopeResolver configured"))
	case hasTargetScope && in.Authorizer == nil:
		return internalFailure(LayerAuthorization, errors.New("layer-8: incomplete Path A (Authorizer nil)"))
	}

	opScope := *in.OperationScope
	target := *in.TargetRef

	if hasTargetScope {
		// Path A: authorize before lookup, then pure scope match.
		decision := in.Authorizer.Authorize(ctx, *in.Caller, target, *in.TargetScope)
		if decision != Allow {
			return Result{
				FailedAt: LayerAuthorization,
				Problem:  SafeDenial(decision),
			}
		}
		if v := CheckOperationTargetScopeMatch(opScope, *in.TargetScope); v != nil {
			return Result{
				Violations: []apiproblem.Violation{*v},
				FailedAt:   LayerAuthorization,
			}
		}
		return Result{}
	}

	// Path B: combined authorized target-scope resolution.
	targetScope, available := in.TargetScopeResolver.ResolveAuthorizedTargetScope(ctx, *in.Caller, target)
	if !available {
		return Result{
			FailedAt: LayerAuthorization,
			Problem:  SafeDenial(DenyNotDisclosed),
		}
	}
	if v := CheckOperationTargetScopeMatch(opScope, targetScope); v != nil {
		return Result{
			Violations: []apiproblem.Violation{*v},
			FailedAt:   LayerAuthorization,
		}
	}
	return Result{}
}

func internalFailure(at Layer, cause error) Result {
	return Result{
		FailedAt: at,
		Problem:  apiproblem.New(apiproblem.CodeInternalError),
		Err:      cause,
	}
}

func decodeFailedAt(prob *apiproblem.Problem) Layer {
	if prob == nil {
		return LayerDecode
	}
	switch prob.Code {
	case apiproblem.CodeRequestTooLarge:
		return LayerHTTPContent
	case apiproblem.CodeUnknownField, apiproblem.CodeDuplicateField:
		return LayerFieldHygiene
	default:
		return LayerDecode
	}
}
