package apivalid

import (
	"context"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/apiref"
)

// CallerContext is the minimal request-scoped identity/scope the authorizer
// needs. It is provided by the adopting feature, not resolved here
// (D-04, D-17, F12-SEC-004).
type CallerContext struct {
	Scopes []apimeta.ScopeIdentity // scopes the caller is entitled to
}

// Decision is the outcome of an authorization check. Denials are mapped
// uniformly so inaccessible objects are never disclosed (F12-SEC-004).
type Decision int

const (
	// Allow permits the operation to proceed to subsequent layer-8 checks.
	Allow Decision = iota
	// DenyNotDisclosed is a cross-scope denial that MUST be indistinguishable
	// from "not found" in the client-facing response.
	DenyNotDisclosed
	// DenyKnown is an in-scope authorization denial where existence is
	// already known to the caller.
	DenyKnown
)

// ScopeAuthorizer is implemented by adopting features. FEATURE-0012 ships the
// interface and the safe-denial mapping only — no concrete authorizer. When
// the target scope is derivable from the request without reading the object,
// adopters MUST call Authorize BEFORE lookup so a cross-scope denial performs
// no existence-dependent work (D-04, F12-SEC-004).
type ScopeAuthorizer interface {
	Authorize(ctx context.Context, caller CallerContext, target apiref.TypedRef, targetScope apimeta.ScopeIdentity) Decision
}

// AuthorizedResolver is the required contract when the target scope is only
// knowable AFTER reading the object. Adopters implement Resolve so that a
// missing object and a present-but-unauthorized object return the SAME uniform
// unavailable outcome (found=false with no leaked detail), hiding whether
// resolution or authorization failed. Implementations MUST NOT branch
// observably (timing, side effects, logs, audit) between the two cases
// (D-04, F12-SEC-004, F12-IMPL-002).
type AuthorizedResolver interface {
	Resolve(ctx context.Context, caller CallerContext, target apiref.TypedRef) (obj any, found bool)
}

// AuthorizedTargetScopeResolver resolves the canonical governance scope of
// an Operation's target reference WITH authorization. It is the adopter-owned
// contract that layer 8 uses for Operation target-scope equality when the
// target scope requires lookup (path B).
//
// Binding semantics:
//   - available=false is the single uniform result for BOTH:
//     (a) target absent;
//     (b) target present but unauthorized.
//     Those two cases MUST follow the same path and map through SafeDenial.
//     No target scope or mismatch detail is returned when available=false.
//   - available=true means an authorized canonical ScopeIdentity is
//     available and the caller may proceed to CheckOperationTargetScopeMatch.
//
// Implementations MUST NOT branch observably (timing, side effects, logs,
// audit) between the absent and unauthorized cases (D-17, F12-SCOPE-002).
type AuthorizedTargetScopeResolver interface {
	ResolveAuthorizedTargetScope(
		ctx context.Context,
		caller CallerContext,
		target apiref.TypedRef,
	) (scope apimeta.ScopeIdentity, available bool)
}

// SafeDenial maps a Decision to a uniform Problem. DenyNotDisclosed always
// maps to an identical 404 RESOURCE_NOT_FOUND (same code, title, and detail)
// so that "exists but inaccessible" and "does not exist" are indistinguishable
// in the RESPONSE (F12-SCOPE-002, F12-SEC-004). DenyKnown maps to 403
// AUTHORIZATION_DENIED. Allow returns nil.
//
// NOTE: this response mapping alone is NOT sufficient for
// no-existence-disclosure — control-flow/timing equivalence is the adopter's
// responsibility via authorize-before-lookup or AuthorizedResolver. This is a
// path/response-equivalence guarantee, not a perfect constant-time guarantee.
//
// SafeDenial is owned and tested under internal/apivalid/authz.go, not
// internal/apiproblem. RequestID/Instance are left empty so callers may attach
// correlation identifiers without changing the stable denial shape; secrets,
// credentials, tokens, and inaccessible resource details MUST NOT be added.
func SafeDenial(d Decision) *apiproblem.Problem {
	switch d {
	case Allow:
		return nil
	case DenyKnown:
		return apiproblem.New(apiproblem.CodeAuthorizationDenied)
	case DenyNotDisclosed:
		return apiproblem.New(apiproblem.CodeResourceNotFound)
	default:
		// Fail closed without disclosing existence for unrecognized values.
		return apiproblem.New(apiproblem.CodeResourceNotFound)
	}
}

const operationTargetScopeRefPointer = "/metadata/scopeRef"

// CheckOperationTargetScopeMatch compares the Operation's canonical
// ScopeIdentity with the resolved target's canonical ScopeIdentity
// (D-17, F12-SCOPE-002, F12-REF-001).
//
// Returns nil when scopes match. When they differ (kind or UID), returns a
// Violation with code OPERATION_TARGET_SCOPE_MISMATCH and JSON Pointer
// /metadata/scopeRef.
//
// Callers MUST invoke this only after authorized target-scope resolution
// succeeds. Unauthorized or unavailable targets MUST use SafeDenial and MUST
// NOT reach this comparator (no existence or mismatch disclosure).
//
// The Message is intentionally generic: it does not embed scope kind/UID
// values, secrets, or inaccessible-resource detail.
func CheckOperationTargetScopeMatch(
	opScope apimeta.ScopeIdentity,
	targetScope apimeta.ScopeIdentity,
) *apiproblem.Violation {
	if opScope == targetScope {
		return nil
	}
	return &apiproblem.Violation{
		Field:   operationTargetScopeRefPointer,
		Code:    apiproblem.ViolationOperationTargetScopeMismatch,
		Message: "operation scopeRef does not match the target governance scope",
	}
}
