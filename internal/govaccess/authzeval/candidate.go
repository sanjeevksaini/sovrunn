package authzeval

import (
	"crypto/sha256"
	"encoding/json"
	"sort"

	"github.com/sanjeevksaini/sovrunn/internal/govaccess/evidence"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/model"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/operation"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/state"
)

// AuthorizationCandidate is the sealed pure-algebra output that is not yet a
// published AuthorizationResult.
type AuthorizationCandidate struct {
	sealed bool

	input              model.AuthorizationInput
	decision           string
	reasonCodes        []string
	contributingGrants []model.RoleAssignmentSpec
	digest             evidence.CandidateDigest
}

func newAuthorizationCandidate(
	input model.AuthorizationInput,
	decision string,
	reasonCodes []string,
	contributing []model.RoleAssignmentSpec,
) (AuthorizationCandidate, operation.MechanicalFailure) {
	if !input.Sealed() || !modelDecisionValid(decision) || len(reasonCodes) == 0 {
		return AuthorizationCandidate{}, operation.NewMechanicalFailure("authorization_candidate_invalid")
	}
	reasons := append([]string(nil), reasonCodes...)
	sort.Strings(reasons)
	payload := struct {
		RequestID string   `json:"requestId"`
		Decision  string   `json:"decision"`
		Reasons   []string `json:"reasons"`
	}{
		RequestID: input.RequestID(),
		Decision:  decision,
		Reasons:   reasons,
	}
	raw, err := json.Marshal(payload)
	if err != nil {
		return AuthorizationCandidate{}, operation.NewMechanicalFailure("authorization_candidate_marshal")
	}
	sum := sha256.Sum256(raw)
	digest, fail := evidence.NewCandidateDigest(sum[:])
	if fail.Reason() != "" {
		return AuthorizationCandidate{}, fail
	}
	return AuthorizationCandidate{
		sealed:             true,
		input:              input,
		decision:           decision,
		reasonCodes:        reasons,
		contributingGrants: append([]model.RoleAssignmentSpec(nil), contributing...),
		digest:             digest,
	}, operation.MechanicalFailure{}
}

func (c AuthorizationCandidate) Sealed() bool { return c.sealed }
func (c AuthorizationCandidate) Input() model.AuthorizationInput {
	return c.input
}
func (c AuthorizationCandidate) Decision() string { return c.decision }
func (c AuthorizationCandidate) ReasonCodes() []string {
	return append([]string(nil), c.reasonCodes...)
}
func (c AuthorizationCandidate) ContributingGrants() []model.RoleAssignmentSpec {
	return append([]model.RoleAssignmentSpec(nil), c.contributingGrants...)
}
func (c AuthorizationCandidate) CandidateDigest() evidence.CandidateDigest { return c.digest }

// OperationCandidate carries the authzeval candidate plus currentness claim for
// publication-time revalidation.
type OperationCandidate struct {
	sealed      bool
	candidate   AuthorizationCandidate
	currentness state.AuthorizationCurrentnessClaim
}

func NewOperationCandidate(
	candidate AuthorizationCandidate,
	currentness state.AuthorizationCurrentnessClaim,
) (OperationCandidate, operation.MechanicalFailure) {
	if !candidate.sealed {
		return OperationCandidate{}, operation.NewMechanicalFailure("operation_candidate_unsealed")
	}
	if !currentness.Sealed() {
		return OperationCandidate{}, operation.NewMechanicalFailure("operation_candidate_currentness_unsealed")
	}
	return OperationCandidate{
		sealed:      true,
		candidate:   candidate,
		currentness: currentness,
	}, operation.MechanicalFailure{}
}

func (c OperationCandidate) Sealed() bool { return c.sealed }
func (c OperationCandidate) Candidate() AuthorizationCandidate {
	return c.candidate
}
func (c OperationCandidate) Currentness() state.AuthorizationCurrentnessClaim {
	return c.currentness
}

func modelDecisionValid(v string) bool {
	return v == model.AuthorizationDecisionAllow || v == model.AuthorizationDecisionDeny
}
