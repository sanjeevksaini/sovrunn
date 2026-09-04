package evidence

import (
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/operation"
)

// AuthorizationOutcome is the closed Allow | Deny candidate outcome.
type AuthorizationOutcome string

const (
	AuthorizationAllow AuthorizationOutcome = "Allow"
	AuthorizationDeny  AuthorizationOutcome = "Deny"
)

// Valid reports whether o is a closed AuthorizationOutcome.
func (o AuthorizationOutcome) Valid() bool {
	switch o {
	case AuthorizationAllow, AuthorizationDeny:
		return true
	default:
		return false
	}
}

// AuthorizationCandidateDescriptor is a sealed evidence-owned candidate
// descriptor constructed only by this package.
type AuthorizationCandidateDescriptor struct {
	sealed   bool
	digest   CandidateDigest
	outcome  AuthorizationOutcome
	action   string
	targetID string
}

// NewAuthorizationCandidateDescriptor seals a candidate descriptor.
func NewAuthorizationCandidateDescriptor(
	digest CandidateDigest,
	outcome AuthorizationOutcome,
	action string,
	targetID string,
) (AuthorizationCandidateDescriptor, operation.MechanicalFailure) {
	if !digest.sealed {
		return AuthorizationCandidateDescriptor{}, operation.NewMechanicalFailure("authz_candidate_digest_unsealed")
	}
	if !outcome.Valid() {
		return AuthorizationCandidateDescriptor{}, operation.NewMechanicalFailure("authz_candidate_outcome_invalid")
	}
	if action == "" || targetID == "" {
		return AuthorizationCandidateDescriptor{}, operation.NewMechanicalFailure("authz_candidate_fields_empty")
	}
	return AuthorizationCandidateDescriptor{
		sealed:   true,
		digest:   digest,
		outcome:  outcome,
		action:   action,
		targetID: targetID,
	}, operation.MechanicalFailure{}
}

// Sealed reports whether d was package-constructed.
func (d AuthorizationCandidateDescriptor) Sealed() bool { return d.sealed }

// Digest returns the sealed candidate digest.
func (d AuthorizationCandidateDescriptor) Digest() CandidateDigest { return d.digest }

// Outcome returns the sealed Allow|Deny outcome.
func (d AuthorizationCandidateDescriptor) Outcome() AuthorizationOutcome { return d.outcome }

// Action returns the sealed action.
func (d AuthorizationCandidateDescriptor) Action() string { return d.action }

// TargetID returns the sealed target identity.
func (d AuthorizationCandidateDescriptor) TargetID() string { return d.targetID }

// AuthorizationAdmissionProof binds candidate digest, dependency versions, and
// publication context (design §5.4).
type AuthorizationAdmissionProof struct {
	sealed     bool
	candidate  CandidateDigest
	dependency DependencyVersionDigest
	context    PublicationContext
	admission  AdmissionDigest
}

// BindAuthorizationAdmission seals an admission proof.
func BindAuthorizationAdmission(
	candidate CandidateDigest,
	dependency DependencyVersionDigest,
	ctx PublicationContext,
) (AuthorizationAdmissionProof, operation.MechanicalFailure) {
	if !candidate.sealed {
		return AuthorizationAdmissionProof{}, operation.NewMechanicalFailure("admission_candidate_unsealed")
	}
	if !dependency.sealed {
		return AuthorizationAdmissionProof{}, operation.NewMechanicalFailure("admission_dependency_unsealed")
	}
	if !ctx.sealed {
		return AuthorizationAdmissionProof{}, operation.NewMechanicalFailure("admission_context_unsealed")
	}
	h := sha256Concat(candidate.value, dependency.value, []byte(ctx.publicationInstant.UTC().Format("2006-01-02T15:04:05.999999999Z07:00")))
	return AuthorizationAdmissionProof{
		sealed:     true,
		candidate:  candidate,
		dependency: dependency,
		context:    ctx,
		admission:  AdmissionDigest{sealed: true, value: h},
	}, operation.MechanicalFailure{}
}

// Sealed reports whether p was constructed through BindAuthorizationAdmission.
func (p AuthorizationAdmissionProof) Sealed() bool { return p.sealed }

// CandidateDigest returns the sealed candidate digest.
func (p AuthorizationAdmissionProof) CandidateDigest() CandidateDigest { return p.candidate }

// DependencyVersionDigest returns the sealed dependency digest.
func (p AuthorizationAdmissionProof) DependencyVersionDigest() DependencyVersionDigest {
	return p.dependency
}

// PublicationContext returns the sealed publication context.
func (p AuthorizationAdmissionProof) PublicationContext() PublicationContext { return p.context }

// AdmissionDigest returns the sealed admission digest.
func (p AuthorizationAdmissionProof) AdmissionDigest() AdmissionDigest { return p.admission }

// FinalizedAuthorizationCandidate binds a descriptor to an admission proof and
// publication context after semantic revalidation.
type FinalizedAuthorizationCandidate struct {
	sealed     bool
	descriptor AuthorizationCandidateDescriptor
	admission  AuthorizationAdmissionProof
	context    PublicationContext
}

// NewFinalizedAuthorizationCandidate seals a finalized candidate.
func NewFinalizedAuthorizationCandidate(
	descriptor AuthorizationCandidateDescriptor,
	admission AuthorizationAdmissionProof,
	ctx PublicationContext,
) (FinalizedAuthorizationCandidate, operation.MechanicalFailure) {
	if !descriptor.sealed {
		return FinalizedAuthorizationCandidate{}, operation.NewMechanicalFailure("finalized_candidate_descriptor_unsealed")
	}
	if !admission.sealed {
		return FinalizedAuthorizationCandidate{}, operation.NewMechanicalFailure("finalized_candidate_admission_unsealed")
	}
	if !ctx.sealed {
		return FinalizedAuthorizationCandidate{}, operation.NewMechanicalFailure("finalized_candidate_context_unsealed")
	}
	if !descriptor.digest.Equal(admission.candidate) {
		return FinalizedAuthorizationCandidate{}, operation.NewMechanicalFailure("finalized_candidate_digest_mismatch")
	}
	if !admission.context.Equal(ctx) {
		return FinalizedAuthorizationCandidate{}, operation.NewMechanicalFailure("finalized_candidate_context_mismatch")
	}
	return FinalizedAuthorizationCandidate{
		sealed:     true,
		descriptor: descriptor,
		admission:  admission,
		context:    ctx,
	}, operation.MechanicalFailure{}
}

// Sealed reports whether c was package-constructed.
func (c FinalizedAuthorizationCandidate) Sealed() bool { return c.sealed }

// Descriptor returns the sealed descriptor.
func (c FinalizedAuthorizationCandidate) Descriptor() AuthorizationCandidateDescriptor {
	return c.descriptor
}

// AdmissionProof returns the sealed admission proof.
func (c FinalizedAuthorizationCandidate) AdmissionProof() AuthorizationAdmissionProof {
	return c.admission
}

// PublicationContext returns the sealed publication context.
func (c FinalizedAuthorizationCandidate) PublicationContext() PublicationContext { return c.context }

// Outcome returns the sealed Allow|Deny outcome.
func (c FinalizedAuthorizationCandidate) Outcome() AuthorizationOutcome {
	return c.descriptor.outcome
}
