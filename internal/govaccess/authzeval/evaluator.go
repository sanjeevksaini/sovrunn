package authzeval

import (
	"strings"
	"time"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/model"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/operation"
)

// Evaluator performs pure deterministic authorization algebra.
type Evaluator struct{}

func NewEvaluator() Evaluator { return Evaluator{} }

func (Evaluator) Evaluate(input model.AuthorizationInput) (AuthorizationCandidate, operation.MechanicalFailure) {
	if !input.Sealed() {
		return AuthorizationCandidate{}, operation.NewMechanicalFailure("authorization_input_unsealed")
	}
	decision, reasons, grants := evaluateDecision(input)
	return newAuthorizationCandidate(input, decision, reasons, grants)
}

func evaluateDecision(input model.AuthorizationInput) (string, []string, []model.RoleAssignmentSpec) {
	if !input.GuardrailAllows() {
		return model.AuthorizationDecisionDeny, []string{"GUARDRAIL_DENY"}, nil
	}
	switch input.PolicyOutcome() {
	case model.PolicyOutcomeDeny:
		return model.AuthorizationDecisionDeny, []string{"POLICY_DENY"}, nil
	case model.PolicyOutcomeIndeterminate:
		return model.AuthorizationDecisionDeny, []string{"POLICY_INDETERMINATE_MAPPED_TO_DENY"}, nil
	case model.PolicyOutcomeRequiresApproval:
		return model.AuthorizationDecisionDeny, []string{"POLICY_REQUIRES_APPROVAL"}, nil
	}

	if membershipRequiredForScope(apimeta.ScopeKind(input.ScopeRef().Kind)) && !input.ScopeMembershipCurrent() {
		return model.AuthorizationDecisionDeny, []string{"MEMBERSHIP_NOT_CURRENT"}, nil
	}

	privileged := toSet(input.PrivilegedRoleUIDs())
	suspended := toSet(input.SuspendedRoleUIDs())
	actions := input.RoleActions()

	var contributing []model.RoleAssignmentSpec
	for _, spec := range input.DirectAssignments() {
		if !assignmentApplicable(spec, input, actions, privileged, suspended, true) {
			continue
		}
		contributing = append(contributing, spec)
	}
	for _, spec := range input.GroupAssignments() {
		if !assignmentApplicable(spec, input, actions, privileged, suspended, false) {
			continue
		}
		contributing = append(contributing, spec)
	}
	if len(contributing) == 0 {
		return model.AuthorizationDecisionDeny, []string{"NO_APPLICABLE_ASSIGNMENT"}, nil
	}
	return model.AuthorizationDecisionAllow, []string{"MATCHED_ASSIGNMENT"}, contributing
}

func assignmentApplicable(
	spec model.RoleAssignmentSpec,
	input model.AuthorizationInput,
	roleActions map[string][]string,
	privileged map[string]struct{},
	suspended map[string]struct{},
	requireDirect bool,
) bool {
	if _, found := suspended[spec.RoleDefinitionRef.UID]; found {
		return false
	}
	if _, isPrivileged := privileged[spec.RoleDefinitionRef.UID]; isPrivileged {
		if spec.Validity != model.AssignmentValidityTimeBound || !spec.RoleHolderRef.IsPrincipal() {
			return false
		}
	}
	if !assignmentMatchesHolder(spec, input, requireDirect) {
		return false
	}
	if !validAt(spec, input.EvaluatedAt()) {
		return false
	}
	if spec.ResourceRef != nil && !typedRefEqual(*spec.ResourceRef, input.TargetRef()) {
		return false
	}
	key := spec.RoleDefinitionRef.UID + "@" + spec.RoleDefinitionVersion
	allowedActions, ok := roleActions[key]
	if !ok {
		return false
	}
	for _, action := range allowedActions {
		if action == input.Action() {
			return true
		}
	}
	return false
}

func assignmentMatchesHolder(spec model.RoleAssignmentSpec, input model.AuthorizationInput, requireDirect bool) bool {
	if requireDirect {
		if !spec.RoleHolderRef.IsPrincipal() {
			return false
		}
		pr := spec.RoleHolderRef.Principal
		return pr != nil && pr.Equal(input.Actor())
	}
	if !spec.RoleHolderRef.IsAccessGroup() {
		return false
	}
	ref := spec.RoleHolderRef.AccessGroup
	if ref == nil {
		return false
	}
	for _, uid := range input.ActiveGroupUIDs() {
		if uid == ref.UID {
			return true
		}
	}
	return false
}

func validAt(spec model.RoleAssignmentSpec, at time.Time) bool {
	switch spec.Validity {
	case model.AssignmentValidityStanding:
		return true
	case model.AssignmentValidityTimeBound:
		if spec.NotBefore == nil || spec.ExpiresAt == nil {
			return false
		}
		return (at.Equal(spec.NotBefore.UTC()) || at.After(spec.NotBefore.UTC())) &&
			at.Before(spec.ExpiresAt.UTC())
	default:
		return false
	}
}

func typedRefEqual(a, b apimeta.TypedRef) bool {
	return a.APIVersion == b.APIVersion && a.Kind == b.Kind && a.Name == b.Name && a.UID == b.UID
}

func toSet(in []string) map[string]struct{} {
	out := make(map[string]struct{}, len(in))
	for _, v := range in {
		v = strings.TrimSpace(v)
		if v != "" {
			out[v] = struct{}{}
		}
	}
	return out
}

func membershipRequiredForScope(kind apimeta.ScopeKind) bool {
	switch kind {
	case apimeta.ScopeOrganization, apimeta.ScopeOrganizationUnit, apimeta.ScopeTenant, apimeta.ScopeProject, apimeta.ScopeCloudProvider:
		return true
	default:
		return false
	}
}
