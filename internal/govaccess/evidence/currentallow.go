package evidence

import (
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/operation"
)

// CurrentAllowProof is the unforgeable mutation gate for authorization-bearing
// publication (DD-12; design §5.5).
type CurrentAllowProof struct {
	sealed    bool
	candidate CandidateDigest
	admission AdmissionDigest
	context   PublicationContext
}

// RouteDenyToEvidenceOnly is the internal direction for a finalized Deny
// candidate. It is not a CurrentAllowProof and cannot authorize mutation.
type RouteDenyToEvidenceOnly struct {
	sealed    bool
	candidate CandidateDigest
	context   PublicationContext
}

// Sealed reports whether r was produced by RequireCurrentAllow.
func (r RouteDenyToEvidenceOnly) Sealed() bool { return r.sealed }

// CandidateDigest returns the sealed candidate digest.
func (r RouteDenyToEvidenceOnly) CandidateDigest() CandidateDigest { return r.candidate }

// PublicationContext returns the sealed publication context.
func (r RouteDenyToEvidenceOnly) PublicationContext() PublicationContext { return r.context }

// RequireCurrentAllowOutcome is a closed result of RequireCurrentAllow.
type RequireCurrentAllowOutcome struct {
	kind  requireKind
	allow CurrentAllowProof
	deny  RouteDenyToEvidenceOnly
	fail  operation.MechanicalFailure
}

type requireKind uint8

const (
	requireNone requireKind = iota
	requireAllow
	requireDeny
	requireFail
)

// RequireCurrentAllow converts only a finalized Allow into CurrentAllowProof.
// A finalized Deny yields RouteDenyToEvidenceOnly (design §5.4/§5.5).
func RequireCurrentAllow(
	finalized FinalizedAuthorizationCandidate,
	admission AuthorizationAdmissionProof,
	ctx PublicationContext,
) RequireCurrentAllowOutcome {
	if !finalized.sealed {
		return RequireCurrentAllowOutcome{kind: requireFail, fail: operation.NewMechanicalFailure("require_current_allow_finalized_unsealed")}
	}
	if !admission.sealed {
		return RequireCurrentAllowOutcome{kind: requireFail, fail: operation.NewMechanicalFailure("require_current_allow_admission_unsealed")}
	}
	if !ctx.sealed {
		return RequireCurrentAllowOutcome{kind: requireFail, fail: operation.NewMechanicalFailure("require_current_allow_context_unsealed")}
	}
	if !finalized.admission.admission.Equal(admission.admission) {
		return RequireCurrentAllowOutcome{kind: requireFail, fail: operation.NewMechanicalFailure("require_current_allow_admission_mismatch")}
	}
	if !finalized.context.Equal(ctx) || !admission.context.Equal(ctx) {
		return RequireCurrentAllowOutcome{kind: requireFail, fail: operation.NewMechanicalFailure("require_current_allow_context_mismatch")}
	}
	switch finalized.Outcome() {
	case AuthorizationAllow:
		return RequireCurrentAllowOutcome{
			kind: requireAllow,
			allow: CurrentAllowProof{
				sealed:    true,
				candidate: finalized.descriptor.digest,
				admission: admission.admission,
				context:   ctx,
			},
		}
	case AuthorizationDeny:
		return RequireCurrentAllowOutcome{
			kind: requireDeny,
			deny: RouteDenyToEvidenceOnly{
				sealed:    true,
				candidate: finalized.descriptor.digest,
				context:   ctx,
			},
		}
	default:
		return RequireCurrentAllowOutcome{kind: requireFail, fail: operation.NewMechanicalFailure("require_current_allow_outcome_invalid")}
	}
}

// AsAllow returns the CurrentAllowProof when the outcome is Allow.
func (o RequireCurrentAllowOutcome) AsAllow() (CurrentAllowProof, bool) {
	if o.kind != requireAllow {
		return CurrentAllowProof{}, false
	}
	return o.allow, true
}

// AsDenyRoute returns the RouteDenyToEvidenceOnly when the outcome is Deny.
func (o RequireCurrentAllowOutcome) AsDenyRoute() (RouteDenyToEvidenceOnly, bool) {
	if o.kind != requireDeny {
		return RouteDenyToEvidenceOnly{}, false
	}
	return o.deny, true
}

// AsMechanicalFailure returns the failure when RequireCurrentAllow failed.
func (o RequireCurrentAllowOutcome) AsMechanicalFailure() (operation.MechanicalFailure, bool) {
	if o.kind != requireFail {
		return operation.MechanicalFailure{}, false
	}
	return o.fail, true
}

// Sealed reports whether p was produced by RequireCurrentAllow for Allow.
func (p CurrentAllowProof) Sealed() bool { return p.sealed }

// CandidateDigest returns the sealed candidate digest.
func (p CurrentAllowProof) CandidateDigest() CandidateDigest { return p.candidate }

// AdmissionDigest returns the sealed admission digest.
func (p CurrentAllowProof) AdmissionDigest() AdmissionDigest { return p.admission }

// PublicationContext returns the sealed publication context.
func (p CurrentAllowProof) PublicationContext() PublicationContext { return p.context }

// Equal reports sealed equality.
func (p CurrentAllowProof) Equal(other CurrentAllowProof) bool {
	return p.sealed && other.sealed &&
		p.candidate.Equal(other.candidate) &&
		p.admission.Equal(other.admission) &&
		p.context.Equal(other.context)
}
