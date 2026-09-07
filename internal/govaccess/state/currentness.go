package state

import (
	"errors"
	"fmt"

	"github.com/sanjeevksaini/sovrunn/internal/govaccess/evidence"
)

var errInvalid = func(reason string) error {
	return fmt.Errorf("state: %s", reason)
}

// AuthorizationCurrentnessClaim binds a candidate digest to its dependency set.
type AuthorizationCurrentnessClaim struct {
	sealed     bool
	candidate  evidence.CandidateDigest
	dependency AuthorizationDependencyVersionSet
}

// NewAuthorizationCurrentnessClaim seals a currentness claim.
func NewAuthorizationCurrentnessClaim(
	candidate evidence.CandidateDigest,
	dependency AuthorizationDependencyVersionSet,
) (AuthorizationCurrentnessClaim, error) {
	if !candidate.Sealed() {
		return AuthorizationCurrentnessClaim{}, errInvalid("currentness_candidate_unsealed")
	}
	if !dependency.sealed {
		return AuthorizationCurrentnessClaim{}, errInvalid("currentness_dependency_unsealed")
	}
	return AuthorizationCurrentnessClaim{
		sealed:     true,
		candidate:  candidate,
		dependency: dependency,
	}, nil
}

// Sealed reports whether c was package-constructed.
func (c AuthorizationCurrentnessClaim) Sealed() bool { return c.sealed }

// CandidateDigest returns the sealed candidate digest.
func (c AuthorizationCurrentnessClaim) CandidateDigest() evidence.CandidateDigest {
	return c.candidate
}

// DependencySet returns the sealed dependency-version set.
func (c AuthorizationCurrentnessClaim) DependencySet() AuthorizationDependencyVersionSet {
	return c.dependency
}

// ParticipantOwner is the closed semantic-owner identity constant.
type ParticipantOwner string

const (
	RoleAssignmentOwner    ParticipantOwner = "roleassign"
	MembershipOwner        ParticipantOwner = "membership"
	RoleDefinitionOwner    ParticipantOwner = "roledefinition"
	AccessGroupOwner       ParticipantOwner = "accessgroup"
	ApprovalOwner          ParticipantOwner = "approval"
	PrivilegedOwner        ParticipantOwner = "privileged"
	ReviewOwner            ParticipantOwner = "review"
	ExceptionOwner         ParticipantOwner = "exception"
	GovernanceProfileOwner ParticipantOwner = "governanceprofile"
)

// Valid reports whether o is a registered participant owner.
func (o ParticipantOwner) Valid() bool {
	switch o {
	case RoleAssignmentOwner, MembershipOwner, RoleDefinitionOwner, AccessGroupOwner,
		ApprovalOwner, PrivilegedOwner, ReviewOwner, ExceptionOwner, GovernanceProfileOwner:
		return true
	default:
		return false
	}
}

// ParticipantID is a stable participant identity.
type ParticipantID string

// IntentDigest is an immutable pre-publication intent digest.
type IntentDigest []byte

// ResultRole is OriginatingResult | SupportingResult.
type ResultRole string

const (
	OriginatingResult ResultRole = "OriginatingResult"
	SupportingResult  ResultRole = "SupportingResult"
)

// Valid reports whether r is a closed ResultRole.
func (r ResultRole) Valid() bool {
	switch r {
	case OriginatingResult, SupportingResult:
		return true
	default:
		return false
	}
}

// ParticipantClaim is a sealed publication-time-free participant claim.
type ParticipantClaim struct {
	sealed   bool
	owner    ParticipantOwner
	id       ParticipantID
	intent   []byte
	versions ExpectedVersionSet
	role     ResultRole
}

// NewParticipantClaim seals a participant claim.
func NewParticipantClaim(
	owner ParticipantOwner,
	id ParticipantID,
	intent IntentDigest,
	versions ExpectedVersionSet,
	role ResultRole,
) (ParticipantClaim, error) {
	if !owner.Valid() {
		return ParticipantClaim{}, errInvalid("participant_owner_invalid")
	}
	if id == "" {
		return ParticipantClaim{}, errInvalid("participant_id_empty")
	}
	if len(intent) == 0 {
		return ParticipantClaim{}, errInvalid("participant_intent_empty")
	}
	if !versions.sealed {
		return ParticipantClaim{}, errInvalid("participant_versions_unsealed")
	}
	if !role.Valid() {
		return ParticipantClaim{}, errInvalid("participant_role_invalid")
	}
	return ParticipantClaim{
		sealed:   true,
		owner:    owner,
		id:       id,
		intent:   copyBytes(intent),
		versions: versions,
		role:     role,
	}, nil
}

// Sealed reports whether c was package-constructed.
func (c ParticipantClaim) Sealed() bool { return c.sealed }

// Owner returns the sealed owner constant.
func (c ParticipantClaim) Owner() ParticipantOwner { return c.owner }

// ID returns the sealed participant id.
func (c ParticipantClaim) ID() ParticipantID { return c.id }

// IntentDigest returns a defensive copy of the intent digest.
func (c ParticipantClaim) IntentDigest() []byte { return copyBytes(c.intent) }

// Versions returns the sealed expected version set.
func (c ParticipantClaim) Versions() ExpectedVersionSet { return c.versions }

// Role returns the sealed result role.
func (c ParticipantClaim) Role() ResultRole { return c.role }

// Equal reports sealed claim equality.
func (c ParticipantClaim) Equal(other ParticipantClaim) bool {
	if !c.sealed || !other.sealed {
		return false
	}
	return c.owner == other.owner &&
		c.id == other.id &&
		c.role == other.role &&
		digestEqual(c.intent, other.intent) &&
		digestEqual(c.versions.digest, other.versions.digest)
}

// MutationMode is CallerMutationMode | ControllerMutationMode.
type MutationMode string

const (
	CallerMutationMode     MutationMode = "CallerMutation"
	ControllerMutationMode MutationMode = "ControllerMutation"
)

// MutationAuthorityMode is AuthorizationBearing | AutomaticMutationOnly.
type MutationAuthorityMode string

const (
	AuthorizationBearing  MutationAuthorityMode = "AuthorizationBearing"
	AutomaticMutationOnly MutationAuthorityMode = "AutomaticMutationOnly"
)

// Valid reports whether m is a closed MutationAuthorityMode.
func (m MutationAuthorityMode) Valid() bool {
	switch m {
	case AuthorizationBearing, AutomaticMutationOnly:
		return true
	default:
		return false
	}
}

// MutationAdmission is a sealed mode-specific admission value.
type MutationAdmission struct {
	sealed       bool
	mode         MutationMode
	authority    MutationAuthorityMode
	participants []ParticipantClaim
	currentness  AuthorizationCurrentnessClaim
	hasCurrent   bool
}

// NewCallerMutationAdmission seals a caller authorization-bearing admission.
func NewCallerMutationAdmission(
	claims []ParticipantClaim,
	currentness AuthorizationCurrentnessClaim,
) (MutationAdmission, error) {
	if err := validateClaims(claims, true); err != nil {
		return MutationAdmission{}, err
	}
	if !currentness.sealed {
		return MutationAdmission{}, errInvalid("caller_admission_currentness_required")
	}
	return MutationAdmission{
		sealed:       true,
		mode:         CallerMutationMode,
		authority:    AuthorizationBearing,
		participants: append([]ParticipantClaim(nil), claims...),
		currentness:  currentness,
		hasCurrent:   true,
	}, nil
}

// NewAuthorizedControllerMutationAdmission seals an authorization-bearing controller admission.
func NewAuthorizedControllerMutationAdmission(
	claims []ParticipantClaim,
	currentness AuthorizationCurrentnessClaim,
) (MutationAdmission, error) {
	if err := validateClaims(claims, false); err != nil {
		return MutationAdmission{}, err
	}
	if !currentness.sealed {
		return MutationAdmission{}, errInvalid("controller_authz_currentness_required")
	}
	return MutationAdmission{
		sealed:       true,
		mode:         ControllerMutationMode,
		authority:    AuthorizationBearing,
		participants: append([]ParticipantClaim(nil), claims...),
		currentness:  currentness,
		hasCurrent:   true,
	}, nil
}

// NewAutomaticControllerMutationAdmission seals an automatic mutation-only admission.
func NewAutomaticControllerMutationAdmission(claims []ParticipantClaim) (MutationAdmission, error) {
	if err := validateClaims(claims, false); err != nil {
		return MutationAdmission{}, err
	}
	return MutationAdmission{
		sealed:       true,
		mode:         ControllerMutationMode,
		authority:    AutomaticMutationOnly,
		participants: append([]ParticipantClaim(nil), claims...),
	}, nil
}

func validateClaims(claims []ParticipantClaim, requireOriginating bool) error {
	if len(claims) == 0 {
		return errInvalid("admission_participants_empty")
	}
	originating := 0
	seen := map[ParticipantID]bool{}
	for _, c := range claims {
		if !c.sealed {
			return errInvalid("admission_participant_unsealed")
		}
		if seen[c.id] {
			return errInvalid("admission_participant_duplicate")
		}
		seen[c.id] = true
		if c.role == OriginatingResult {
			originating++
		}
	}
	if requireOriginating {
		if originating != 1 {
			return errInvalid("admission_originating_required")
		}
	} else if originating != 0 {
		return errInvalid("admission_originating_forbidden")
	}
	return nil
}

// ErrConflict is returned when version/CAS/reservation checks fail.
var ErrConflict = errors.New("state: conflict")

// ErrNoOp is returned for registered effective-once controller no-op.
var ErrNoOp = errors.New("state: noop")

// ErrRetryEvaluation is returned when evidence-only begin cannot proceed.
var ErrRetryEvaluation = errors.New("state: retry_evaluation")
