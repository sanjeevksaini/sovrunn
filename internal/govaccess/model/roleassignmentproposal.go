package model

import (
	"time"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
)

// RoleAssignmentProposal is the immutable operation-local proposed grant input
// used by grantintent. It is not a resource, route, writer, or event.
type RoleAssignmentProposal struct {
	sealed bool

	proposalUID     string
	resourceUID     string
	resourceName    string
	scopeRef        apimeta.ScopeRef
	spec            RoleAssignmentSpec
	expectedVersion string
	nextVersion     string
	privileged      bool
	policyRef       *apimeta.TypedRef
	nextReviewDueAt *time.Time
}

func NewRoleAssignmentProposal(
	proposalUID string,
	resourceUID string,
	resourceName string,
	scopeRef apimeta.ScopeRef,
	spec RoleAssignmentSpec,
	expectedVersion string,
	nextVersion string,
	privileged bool,
	policyRef *apimeta.TypedRef,
	nextReviewDueAt *time.Time,
) (RoleAssignmentProposal, bool) {
	if proposalUID == "" || resourceUID == "" || resourceName == "" {
		return RoleAssignmentProposal{}, false
	}
	if scopeRef.Kind == "" || !apimeta.ScopeKind(scopeRef.Kind).Valid() {
		return RoleAssignmentProposal{}, false
	}
	if expectedVersion == "" || nextVersion == "" {
		return RoleAssignmentProposal{}, false
	}
	if !validAssignmentList([]RoleAssignmentSpec{spec}) {
		return RoleAssignmentProposal{}, false
	}
	if privileged {
		if !spec.RoleHolderRef.IsPrincipal() || spec.Validity != AssignmentValidityTimeBound {
			return RoleAssignmentProposal{}, false
		}
		if policyRef == nil || policyRef.UID == "" || policyRef.Kind == "" || policyRef.APIVersion == "" || policyRef.Name == "" {
			return RoleAssignmentProposal{}, false
		}
	}
	if !privileged && policyRef != nil {
		return RoleAssignmentProposal{}, false
	}
	out := RoleAssignmentProposal{
		sealed:          true,
		proposalUID:     proposalUID,
		resourceUID:     resourceUID,
		resourceName:    resourceName,
		scopeRef:        scopeRef,
		spec:            copyRoleAssignmentSpecForTask4(spec),
		expectedVersion: expectedVersion,
		nextVersion:     nextVersion,
		privileged:      privileged,
	}
	if policyRef != nil {
		cp := *policyRef
		out.policyRef = &cp
	}
	if nextReviewDueAt != nil {
		t := nextReviewDueAt.UTC()
		out.nextReviewDueAt = &t
	}
	return out, true
}

func (p RoleAssignmentProposal) Sealed() bool               { return p.sealed }
func (p RoleAssignmentProposal) ProposalUID() string        { return p.proposalUID }
func (p RoleAssignmentProposal) ResourceUID() string        { return p.resourceUID }
func (p RoleAssignmentProposal) ResourceName() string       { return p.resourceName }
func (p RoleAssignmentProposal) ScopeRef() apimeta.ScopeRef { return p.scopeRef }
func (p RoleAssignmentProposal) Spec() RoleAssignmentSpec {
	return copyRoleAssignmentSpecForTask4(p.spec)
}
func (p RoleAssignmentProposal) ExpectedVersion() string { return p.expectedVersion }
func (p RoleAssignmentProposal) NextVersion() string     { return p.nextVersion }
func (p RoleAssignmentProposal) Privileged() bool        { return p.privileged }
func (p RoleAssignmentProposal) PolicyRef() (apimeta.TypedRef, bool) {
	if p.policyRef == nil {
		return apimeta.TypedRef{}, false
	}
	cp := *p.policyRef
	return cp, true
}
func (p RoleAssignmentProposal) NextReviewDueAt() *time.Time {
	if p.nextReviewDueAt == nil {
		return nil
	}
	t := p.nextReviewDueAt.UTC()
	return &t
}
