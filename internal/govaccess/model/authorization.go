package model

import (
	"sort"
	"time"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
)

const (
	AuthorizationDecisionAllow = "Allow"
	AuthorizationDecisionDeny  = "Deny"

	PolicyOutcomeAllow            = "Allow"
	PolicyOutcomeDeny             = "Deny"
	PolicyOutcomeIndeterminate    = "Indeterminate"
	PolicyOutcomeRequiresApproval = "RequiresApproval"
)

// AuthorizationInput is the immutable operation-local evaluator input used by
// Task-4 authorization algebra. It is not a resource, route, writer, or event.
type AuthorizationInput struct {
	sealed bool

	actor       PrincipalRef
	action      string
	scopeRef    apimeta.ScopeRef
	targetRef   apimeta.TypedRef
	evaluatedAt time.Time
	requestID   string

	directAssignments []RoleAssignmentSpec
	groupAssignments  []RoleAssignmentSpec
	activeGroupUIDs   []string
	roleActions       map[string][]string
	privilegedRoleUID []string
	suspendedRoleUID  []string

	scopeMembershipCurrent bool
	guardrailAllows        bool
	policyOutcome          string
}

// NewAuthorizationInput seals a deterministic evaluator input and deep-copies
// all reference-typed fields.
func NewAuthorizationInput(
	actor PrincipalRef,
	action string,
	scopeRef apimeta.ScopeRef,
	targetRef apimeta.TypedRef,
	evaluatedAt time.Time,
	requestID string,
	directAssignments []RoleAssignmentSpec,
	groupAssignments []RoleAssignmentSpec,
	activeGroupUIDs []string,
	roleActions map[string][]string,
	privilegedRoleUID []string,
	suspendedRoleUID []string,
	scopeMembershipCurrent bool,
	guardrailAllows bool,
	policyOutcome string,
) (AuthorizationInput, bool) {
	if actor.Issuer == "" || actor.Subject == "" || !actor.PrincipalType.Valid() {
		return AuthorizationInput{}, false
	}
	if action == "" || requestID == "" || evaluatedAt.IsZero() {
		return AuthorizationInput{}, false
	}
	if scopeRef.Kind == "" || !apimeta.ScopeKind(scopeRef.Kind).Valid() {
		return AuthorizationInput{}, false
	}
	if targetRef.APIVersion == "" || targetRef.Kind == "" || targetRef.Name == "" {
		return AuthorizationInput{}, false
	}
	if !validPolicyOutcome(policyOutcome) {
		return AuthorizationInput{}, false
	}
	if !validRoleActionMap(roleActions) {
		return AuthorizationInput{}, false
	}
	if !validAssignmentList(directAssignments) || !validAssignmentList(groupAssignments) {
		return AuthorizationInput{}, false
	}
	return AuthorizationInput{
		sealed:                 true,
		actor:                  actor,
		action:                 action,
		scopeRef:               scopeRef,
		targetRef:              targetRef,
		evaluatedAt:            evaluatedAt.UTC(),
		requestID:              requestID,
		directAssignments:      copyRoleAssignmentSpecsForTask4(directAssignments),
		groupAssignments:       copyRoleAssignmentSpecsForTask4(groupAssignments),
		activeGroupUIDs:        copyStringsForTask4(activeGroupUIDs),
		roleActions:            copyRoleActionMapForTask4(roleActions),
		privilegedRoleUID:      copyStringsForTask4(privilegedRoleUID),
		suspendedRoleUID:       copyStringsForTask4(suspendedRoleUID),
		scopeMembershipCurrent: scopeMembershipCurrent,
		guardrailAllows:        guardrailAllows,
		policyOutcome:          policyOutcome,
	}, true
}

func (i AuthorizationInput) Sealed() bool { return i.sealed }
func (i AuthorizationInput) Actor() PrincipalRef {
	return i.actor
}
func (i AuthorizationInput) Action() string { return i.action }
func (i AuthorizationInput) ScopeRef() apimeta.ScopeRef {
	return i.scopeRef
}
func (i AuthorizationInput) TargetRef() apimeta.TypedRef {
	return i.targetRef
}
func (i AuthorizationInput) EvaluatedAt() time.Time { return i.evaluatedAt }
func (i AuthorizationInput) RequestID() string      { return i.requestID }
func (i AuthorizationInput) DirectAssignments() []RoleAssignmentSpec {
	return copyRoleAssignmentSpecsForTask4(i.directAssignments)
}
func (i AuthorizationInput) GroupAssignments() []RoleAssignmentSpec {
	return copyRoleAssignmentSpecsForTask4(i.groupAssignments)
}
func (i AuthorizationInput) ActiveGroupUIDs() []string {
	return copyStringsForTask4(i.activeGroupUIDs)
}
func (i AuthorizationInput) RoleActions() map[string][]string {
	return copyRoleActionMapForTask4(i.roleActions)
}
func (i AuthorizationInput) PrivilegedRoleUIDs() []string {
	return copyStringsForTask4(i.privilegedRoleUID)
}
func (i AuthorizationInput) SuspendedRoleUIDs() []string {
	return copyStringsForTask4(i.suspendedRoleUID)
}
func (i AuthorizationInput) ScopeMembershipCurrent() bool { return i.scopeMembershipCurrent }
func (i AuthorizationInput) GuardrailAllows() bool        { return i.guardrailAllows }
func (i AuthorizationInput) PolicyOutcome() string        { return i.policyOutcome }

// AuthorizationResult is the immutable operation-local evaluator output. It is
// non-bearer and cannot authorize a separate request.
type AuthorizationResult struct {
	sealed bool

	decision     string
	reasonCodes  []string
	requestID    string
	action       string
	targetRef    apimeta.TypedRef
	evaluatedAt  time.Time
	releasedAt   time.Time
	contributing []RoleAssignmentSpec
}

// NewAuthorizationResult seals and deep-copies an evaluator decision value.
func NewAuthorizationResult(
	decision string,
	reasonCodes []string,
	requestID string,
	action string,
	targetRef apimeta.TypedRef,
	evaluatedAt time.Time,
	releasedAt time.Time,
	contributing []RoleAssignmentSpec,
) (AuthorizationResult, bool) {
	if !validAuthorizationDecision(decision) {
		return AuthorizationResult{}, false
	}
	if action == "" || requestID == "" || evaluatedAt.IsZero() || releasedAt.IsZero() {
		return AuthorizationResult{}, false
	}
	if targetRef.APIVersion == "" || targetRef.Kind == "" || targetRef.Name == "" {
		return AuthorizationResult{}, false
	}
	if len(reasonCodes) == 0 {
		return AuthorizationResult{}, false
	}
	reasons := copyStringsForTask4(reasonCodes)
	sort.Strings(reasons)
	for _, code := range reasons {
		if code == "" {
			return AuthorizationResult{}, false
		}
	}
	if !validAssignmentList(contributing) {
		return AuthorizationResult{}, false
	}
	return AuthorizationResult{
		sealed:       true,
		decision:     decision,
		reasonCodes:  reasons,
		requestID:    requestID,
		action:       action,
		targetRef:    targetRef,
		evaluatedAt:  evaluatedAt.UTC(),
		releasedAt:   releasedAt.UTC(),
		contributing: copyRoleAssignmentSpecsForTask4(contributing),
	}, true
}

func (r AuthorizationResult) Sealed() bool                { return r.sealed }
func (r AuthorizationResult) Decision() string            { return r.decision }
func (r AuthorizationResult) ReasonCodes() []string       { return copyStringsForTask4(r.reasonCodes) }
func (r AuthorizationResult) RequestID() string           { return r.requestID }
func (r AuthorizationResult) Action() string              { return r.action }
func (r AuthorizationResult) TargetRef() apimeta.TypedRef { return r.targetRef }
func (r AuthorizationResult) EvaluatedAt() time.Time      { return r.evaluatedAt }
func (r AuthorizationResult) ReleasedAt() time.Time       { return r.releasedAt }
func (r AuthorizationResult) ContributingAssignments() []RoleAssignmentSpec {
	return copyRoleAssignmentSpecsForTask4(r.contributing)
}

func validAuthorizationDecision(v string) bool {
	return v == AuthorizationDecisionAllow || v == AuthorizationDecisionDeny
}

func validPolicyOutcome(v string) bool {
	switch v {
	case PolicyOutcomeAllow, PolicyOutcomeDeny, PolicyOutcomeIndeterminate, PolicyOutcomeRequiresApproval:
		return true
	default:
		return false
	}
}

func validRoleActionMap(in map[string][]string) bool {
	for k, v := range in {
		if k == "" || len(v) == 0 {
			return false
		}
		for _, action := range v {
			if action == "" {
				return false
			}
		}
	}
	return true
}

func validAssignmentList(assignments []RoleAssignmentSpec) bool {
	for _, spec := range assignments {
		if spec.RoleDefinitionRef.UID == "" || spec.RoleDefinitionVersion == "" || !spec.Validity.Valid() {
			return false
		}
		if spec.RoleHolderRef.Kind == "" {
			return false
		}
	}
	return true
}

func copyRoleActionMapForTask4(in map[string][]string) map[string][]string {
	if len(in) == 0 {
		return nil
	}
	out := make(map[string][]string, len(in))
	for k, v := range in {
		out[k] = copyStringsForTask4(v)
	}
	return out
}

func copyRoleAssignmentSpecsForTask4(in []RoleAssignmentSpec) []RoleAssignmentSpec {
	if len(in) == 0 {
		return nil
	}
	out := make([]RoleAssignmentSpec, len(in))
	for idx := range in {
		out[idx] = copyRoleAssignmentSpecForTask4(in[idx])
	}
	return out
}

func copyRoleAssignmentSpecForTask4(in RoleAssignmentSpec) RoleAssignmentSpec {
	out := in
	if in.NotBefore != nil {
		t := in.NotBefore.UTC()
		out.NotBefore = &t
	}
	if in.ExpiresAt != nil {
		t := in.ExpiresAt.UTC()
		out.ExpiresAt = &t
	}
	if in.ResourceRef != nil {
		ref := *in.ResourceRef
		out.ResourceRef = &ref
	}
	if in.MembershipRef != nil {
		ref := *in.MembershipRef
		out.MembershipRef = &ref
	}
	if in.ResponsiblePartyRef != nil {
		ref := *in.ResponsiblePartyRef
		out.ResponsiblePartyRef = &ref
	}
	if in.AccessReviewRuleRef != nil {
		ref := *in.AccessReviewRuleRef
		out.AccessReviewRuleRef = &ref
	}
	if in.RoleHolderRef.Principal != nil {
		p := *in.RoleHolderRef.Principal
		out.RoleHolderRef.Principal = &p
	}
	if in.RoleHolderRef.AccessGroup != nil {
		g := *in.RoleHolderRef.AccessGroup
		out.RoleHolderRef.AccessGroup = &g
	}
	return out
}

func copyStringsForTask4(in []string) []string {
	if len(in) == 0 {
		return nil
	}
	out := make([]string, len(in))
	copy(out, in)
	return out
}
