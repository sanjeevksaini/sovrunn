package model

import "github.com/sanjeevksaini/sovrunn/internal/apimeta"

const (
	ApprovalRequirementNotRequired = "NotRequired"
	ApprovalRequirementRequired    = "Required"
)

// ApprovalRequirement is an immutable operation-local derivation result for
// approval routing. It is not a resource, route, writer, or event.
type ApprovalRequirement struct {
	sealed bool
	mode   string
	policy *apimeta.TypedRef
}

func NewApprovalRequirementNotRequired() ApprovalRequirement {
	return ApprovalRequirement{sealed: true, mode: ApprovalRequirementNotRequired}
}

func NewApprovalRequirementRequired(policy apimeta.TypedRef) (ApprovalRequirement, bool) {
	if policy.APIVersion == "" || policy.Kind == "" || policy.Name == "" || policy.UID == "" {
		return ApprovalRequirement{}, false
	}
	cp := policy
	return ApprovalRequirement{
		sealed: true,
		mode:   ApprovalRequirementRequired,
		policy: &cp,
	}, true
}

func (r ApprovalRequirement) Sealed() bool { return r.sealed }
func (r ApprovalRequirement) Mode() string { return r.mode }
func (r ApprovalRequirement) IsRequired() bool {
	return r.sealed && r.mode == ApprovalRequirementRequired
}
func (r ApprovalRequirement) PolicyRef() (apimeta.TypedRef, bool) {
	if !r.sealed || r.policy == nil {
		return apimeta.TypedRef{}, false
	}
	cp := *r.policy
	return cp, true
}

func (r ApprovalRequirement) Valid() bool {
	if !r.sealed {
		return false
	}
	switch r.mode {
	case ApprovalRequirementNotRequired:
		return r.policy == nil
	case ApprovalRequirementRequired:
		return r.policy != nil &&
			r.policy.APIVersion != "" &&
			r.policy.Kind != "" &&
			r.policy.Name != "" &&
			r.policy.UID != ""
	default:
		return false
	}
}
