package validate

import (
	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
)

// RFC 6901 pointers for the scope validation pass (design §9.1 step 2, §7.9).
const (
	ptrMetadataScopeRef     = "/metadata/scopeRef"
	ptrMetadataScopeRefKind = "/metadata/scopeRef/kind"
	ptrMetadataScopeRefUID  = "/metadata/scopeRef/uid"
	ptrCanonicalScope       = "/record/semanticIdentity/canonicalScope"
	ptrTopLevelScopeRef     = "/scopeRef"
)

// KindServiceInstance is the typed subject kind for ServiceInstance resources.
// It is never a ScopeKind (F13-SCOPE-004; AD-015).
const KindServiceInstance = "ServiceInstance"

// ScopeInput is the input to ValidateScope (design §8 / §9.1 step 2).
//
// Metadata.ScopeRef is the sole governance-scope authority. Allowed and
// PlatformAllowed carry the governing profile's scope membership. Optional
// fields enable sole-authority conflict detection, ServiceInstance-as-subject
// checks, and semanticIdentity.canonicalScope agreement.
type ScopeInput struct {
	Metadata        apimeta.ObjectMeta
	Allowed         []apimeta.ScopeKind
	PlatformAllowed bool

	// ParallelScopeFields lists RFC 6901 pointers to any second/parallel scope
	// sources detected outside metadata.scopeRef (F13-SCOPE-005; edge 2).
	// Empty means sole authority holds. A typical decode-time finding is
	// "/scopeRef" for a prohibited top-level scopeRef.
	ParallelScopeFields []string

	// SubjectRefs are typed subjects from the decision body (F13-SCOPE-004/008).
	SubjectRefs []apimeta.TypedRef

	// CanonicalScope is record.semanticIdentity.canonicalScope (F13-SCALE-002).
	// When CheckCanonicalScope is true, it must equal CanonicalScopeDescriptor
	// computed from the normalized metadata.scopeRef.
	CanonicalScope      string
	CheckCanonicalScope bool
}

// CanonicalScopeDescriptor returns the derived canonical-scope descriptor
// string from metadata.scopeRef (design §7.9; F13-SCALE-002).
//
// Forms:
//   - Platform (nil / normalized absent scopeRef): "Platform"
//   - Non-Platform: "{ScopeKind}:{uid}"
//
// The descriptor is derived evidence only and never an independent scope
// authority.
func CanonicalScopeDescriptor(scope *apimeta.ScopeRef) string {
	id := apimeta.CanonicalScopeIdentity(apimeta.NormalizeScope(scope))
	if id.Kind == apimeta.ScopePlatform {
		return string(apimeta.ScopePlatform)
	}
	return string(id.Kind) + ":" + id.UID
}

// ValidateScope runs the FEATURE-0013 scope validation pass (design §9.1
// step 2; F13-SCOPE-001…008; F13-SCALE-002; AD-015).
//
// Checks (fail closed, short-circuit on first blocking violation):
//  1. sole-authority — no second/parallel scope source
//  2. canonical Platform absence (NormalizeScope)
//  3. required scope when Platform is not permitted
//  4. six-value ScopeKind vocabulary (ServiceInstance-as-scope → invalid)
//  5. non-Platform uid required
//  6. profile allowed-scope membership (impermissible widening)
//  7. ServiceInstance subject requires Project metadata.scopeRef
//  8. semanticIdentity.canonicalScope matches computed descriptor
//
// The design §8 summary signature is Metadata + Allowed + PlatformAllowed;
// ParallelScopeFields, SubjectRefs, and CanonicalScope are the same pass's
// additional inputs required by §9.1 / §7.9.
func ValidateScope(in ScopeInput) *apiproblem.Problem {
	if p := checkParallelScopeSources(in.ParallelScopeFields); p != nil {
		return p
	}

	normalized := apimeta.NormalizeScope(in.Metadata.ScopeRef)

	if p := checkScopeRequired(normalized, in.PlatformAllowed); p != nil {
		return p
	}

	if normalized != nil {
		if p := checkScopeKindAndUID(normalized); p != nil {
			return p
		}
		if p := checkScopeAllowed(normalized, in.Allowed); p != nil {
			return p
		}
	}

	if p := checkServiceInstanceSubject(normalized, in.SubjectRefs); p != nil {
		return p
	}

	if in.CheckCanonicalScope {
		if p := checkCanonicalScope(normalized, in.CanonicalScope); p != nil {
			return p
		}
	}

	return nil
}

// ValidateScopeMeta is the design §8 convenience form for the core
// Metadata/Allowed/PlatformAllowed checks without subject or canonical-scope
// extras. Prefer ValidateScope when exercising the full §9.1 scope pass.
func ValidateScopeMeta(m apimeta.ObjectMeta, allowed []apimeta.ScopeKind, platformAllowed bool) *apiproblem.Problem {
	return ValidateScope(ScopeInput{
		Metadata:        m,
		Allowed:         allowed,
		PlatformAllowed: platformAllowed,
	})
}

func checkParallelScopeSources(fields []string) *apiproblem.Problem {
	if len(fields) == 0 {
		return nil
	}
	field := fields[0]
	if field == "" {
		field = ptrTopLevelScopeRef
	}
	return scopeProblem(CodeScopeConflict, field,
		"scope must be declared only through metadata.scopeRef; parallel or second scope source is prohibited")
}

func checkScopeRequired(normalized *apimeta.ScopeRef, platformAllowed bool) *apiproblem.Problem {
	if normalized != nil {
		return nil
	}
	if platformAllowed {
		return nil
	}
	return scopeProblem(CodeScopeRequired, ptrMetadataScopeRef,
		"metadata.scopeRef is required when Platform scope is not permitted")
}

func checkScopeKindAndUID(scope *apimeta.ScopeRef) *apiproblem.Problem {
	kind := apimeta.ScopeKind(scope.Kind)
	if !kind.Valid() {
		return scopeProblem(CodeScopeInvalid, ptrMetadataScopeRefKind,
			"metadata.scopeRef.kind must be one of the six FEATURE-0012 governance ScopeKind values; ServiceInstance is a subject, not a ScopeKind")
	}
	if scope.UID == "" {
		return scopeProblem(CodeScopeInvalid, ptrMetadataScopeRefUID,
			"non-Platform metadata.scopeRef requires uid")
	}
	return nil
}

func checkScopeAllowed(scope *apimeta.ScopeRef, allowed []apimeta.ScopeKind) *apiproblem.Problem {
	kind := apimeta.ScopeKind(scope.Kind)
	if containsScopeKind(allowed, kind) {
		return nil
	}
	return scopeProblem(CodeScopeWidening, ptrMetadataScopeRefKind,
		"metadata.scopeRef.kind is not permitted by the governing profile allowedScopes")
}

func checkServiceInstanceSubject(scope *apimeta.ScopeRef, subjects []apimeta.TypedRef) *apiproblem.Problem {
	hasServiceInstance := false
	for _, s := range subjects {
		if s.Kind == KindServiceInstance {
			hasServiceInstance = true
			break
		}
	}
	if !hasServiceInstance {
		return nil
	}
	// F13-SCOPE-004/008: ServiceInstance subject requires Project scopeRef.
	if scope == nil || apimeta.ScopeKind(scope.Kind) != apimeta.ScopeProject || scope.UID == "" {
		return scopeProblem(CodeScopeMismatch, ptrMetadataScopeRef,
			"ServiceInstance subject requires metadata.scopeRef set to the containing Project")
	}
	return nil
}

func checkCanonicalScope(scope *apimeta.ScopeRef, supplied string) *apiproblem.Problem {
	want := CanonicalScopeDescriptor(scope)
	if supplied == want {
		return nil
	}
	return scopeProblem(CodeScopeMismatch, ptrCanonicalScope,
		"semanticIdentity.canonicalScope must match the descriptor derived from metadata.scopeRef")
}

func containsScopeKind(set []apimeta.ScopeKind, v apimeta.ScopeKind) bool {
	for _, k := range set {
		if k == v {
			return true
		}
	}
	return false
}

func scopeProblem(code apiproblem.ViolationCode, field, message string) *apiproblem.Problem {
	return apiproblem.New(apiproblem.CodeValidationFailed).
		WithDetail(message).
		WithViolations([]apiproblem.Violation{{
			Field:   field,
			Code:    code,
			Message: message,
		}})
}
