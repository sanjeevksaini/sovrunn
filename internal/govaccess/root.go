package govaccess

import (
	"fmt"

	"github.com/sanjeevksaini/sovrunn/internal/decision/bundle"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/authzeval"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/clock"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/command"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/evidence"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/exception"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/grantintent"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/membership"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/privileged"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/review"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/roleassign"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/state"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/uow"
)

// Root is the FEATURE-0018 composition output: dependency wiring only.
type Root struct {
	Foundation Foundation

	BundleView bundle.BundleView

	StateStore            *state.Store
	AuthorizationEval     authzeval.Evaluator
	PolicySeam            any
	MutationCoordinator   uow.MutationCoordinatorPort
	GrantPort             roleassign.GrantPort
	ActivationPort        roleassign.ActivationPort
	ReviewRemediationPort roleassign.ReviewRemediationPort
	DirectRevokePort      roleassign.DirectRevokePort

	RoleGrantRequirementDeriver grantintent.RoleGrantRequirementDeriver
	MembershipDeriver           membership.ApprovalRequirementDeriver
	PrivilegedDeriver           privileged.ApprovalRequirementDeriver
	ExceptionDeriver            exception.ApprovalRequirementDeriver

	RoleAssignmentRevokeService command.RoleAssignmentRevokeService
	PrivilegedActivationService privileged.ActivationService
	ReviewRemediationService    review.RemediationService
}

// NewRoot wires clock, state, policy seam, uow, domain ports/derivers, command
// services, and the FEATURE-0013 BundleView. It constructs no records and opens
// no transaction.
func NewRoot(clk clock.Clock, policySeam any) (Root, error) {
	foundation, err := NewFoundation(clk)
	if err != nil {
		return Root{}, err
	}

	view, prob := ComposeFeature0013BundleView()
	if prob != nil {
		return Root{}, fmt.Errorf("govaccess: compose feature-0013 bundle: %s", prob.Code)
	}

	store := state.NewStore(foundation.Clock)
	evaluator := authzeval.NewEvaluator()
	coordinator := uow.NewCoordinator(store)

	grantPort := roleassign.NewGrantPort()
	activationPort := roleassign.NewActivationPort()
	remediationPort := roleassign.NewReviewRemediationPort()
	directRevokePort := roleassign.NewDirectRevokePort()

	return Root{
		Foundation: foundation,
		BundleView: view,

		StateStore:            store,
		AuthorizationEval:     evaluator,
		PolicySeam:            policySeam,
		MutationCoordinator:   coordinator,
		GrantPort:             grantPort,
		ActivationPort:        activationPort,
		ReviewRemediationPort: remediationPort,
		DirectRevokePort:      directRevokePort,

		RoleGrantRequirementDeriver: grantintent.RoleGrantRequirementDeriver{},
		MembershipDeriver:           membership.ApprovalRequirementDeriver{},
		PrivilegedDeriver:           privileged.ApprovalRequirementDeriver{},
		ExceptionDeriver:            exception.ApprovalRequirementDeriver{},

		RoleAssignmentRevokeService: command.NewRoleAssignmentRevokeService(evaluator, directRevokePort, coordinator),
		PrivilegedActivationService: privileged.NewActivationService(activationPort),
		ReviewRemediationService:    review.NewRemediationService(remediationPort),
	}, nil
}

// LocalProfileBundleBytes exposes evidence-owned strict local bundle bytes for
// deterministic tests without defining or publishing profiles in root.
func (r Root) LocalProfileBundleBytes() []byte {
	return evidence.LocalDecisionProfileBundleBytes()
}
