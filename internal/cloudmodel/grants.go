package cloudmodel

import (
	"encoding/json"
	"net/http"

	"github.com/sanjeevksaini/sovrunn/internal/apiconform"
	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
)

// Exact server-resolved bootstrap action names (design §6.1; ADH-2026-045).
// Suspend/resume action constants live in lifecycle.go and are reused here.
const (
	ActionCloudPlatformWrite = "cloudplatform.write"
	ActionCloudPlatformRead  = "cloudplatform.read"
	ActionCloudProviderWrite = "cloudprovider.write"
	ActionCloudProviderRead  = "cloudprovider.read"
	ActionTopologyWrite      = "topology.write"
	ActionTopologyRead       = "topology.read"
	ActionParticipationRead  = "participation.read"

	ActionParticipationRequestProvider        = "participation.request.provider"
	ActionParticipationAcceptPlatform         = "participation.accept.platform"
	ActionParticipationRejectPlatform         = "participation.reject.platform"
	ActionParticipationRequestReleasePlatform = "participation.request-release.platform"
	ActionParticipationWithdrawProvider       = "participation.withdraw.provider"
	ActionParticipationAcceptReleaseProvider  = "participation.accept-release.provider"
	ActionParticipationDeclineReleaseProvider = "participation.decline-release.provider"
)

// HeaderBootstrapGrant is the only reserved forged-grant HTTP header
// (ADH-2026-054; VS0-CF-F15-22). Any present value is ignored and denied.
const HeaderBootstrapGrant = "X-Sovrunn-Bootstrap-Grant"

// BodyMemberBootstrapGrant is the only reserved forged-grant top-level JSON
// member (ADH-2026-054; VS0-CF-F15-22). Any present value is ignored and denied.
const BodyMemberBootstrapGrant = "bootstrapGrant"

// PlatformRootScope is the deployment Platform-root ScopeIdentity used by
// cloudplatform/cloudprovider write and read grants.
var PlatformRootScope = apimeta.ScopeIdentity{
	Kind: apimeta.ScopePlatform,
	UID:  apimeta.PlatformScopeUID,
}

// Grant is a server-resolved, deterministic bootstrap grant (ADH-2026-045).
// Grants are never accepted from headers or bodies. ResourceUID, when set,
// narrows the grant to one target resource (PATCH/GET narrowing).
// RelatedCloudPlatformUID is set only for participation.request.provider.
type Grant struct {
	Action                  string
	Scope                   apimeta.ScopeIdentity
	ResourceUID             string
	RelatedCloudPlatformUID string
}

// ScopedActionGrant converts a Grant to the lifecycle suspend/resume matching
// surface (design §6.1 exact-one resolution).
func (g Grant) ScopedActionGrant() ScopedActionGrant {
	return ScopedActionGrant{Action: g.Action, ScopeUID: g.Scope.UID}
}

// ParticipationRelation binds a CloudProvider UID to a CloudPlatform UID for
// participation.request.provider grants.
type ParticipationRelation struct {
	CloudPlatformUID string
	CloudProviderUID string
}

// SuspendResumeGrantSpec selects which suspend/resume action grants the
// bootstrap principal receives. Callers must ensure that, for any concrete
// participation pair evaluated at request time, exactly one matching grant
// is present; zero or multiple matches fail closed (design §6.1).
type SuspendResumeGrantSpec struct {
	PlatformUIDs []string
	ProviderUIDs []string
}

// BootstrapGrantConfig configures the process-local bootstrap grant set.
// FEATURE-0015 persists no roles, memberships, or assignments.
type BootstrapGrantConfig struct {
	PrincipalID            string
	CloudPlatformUIDs      []string
	CloudProviderUIDs      []string
	ParticipationRelations []ParticipationRelation
	SuspendResume          SuspendResumeGrantSpec
}

// BootstrapGrantResolver is the server-resolved bootstrap-grant source
// (design DD-07 / §6.1 grant file; REQ-F15-17). It holds only in-process
// grant descriptors — never persisted roles, memberships, or assignments.
type BootstrapGrantResolver struct {
	principalID string
	grants      []Grant
}

// NewBootstrapGrantResolver builds the deterministic bootstrap grant set.
// Header and body grant claims are never consulted.
func NewBootstrapGrantResolver(cfg BootstrapGrantConfig) *BootstrapGrantResolver {
	grants := make([]Grant, 0, 16)

	// Platform-root write/read for CloudPlatform and CloudProvider.
	grants = append(grants,
		Grant{Action: ActionCloudPlatformWrite, Scope: PlatformRootScope},
		Grant{Action: ActionCloudPlatformRead, Scope: PlatformRootScope},
		Grant{Action: ActionCloudProviderWrite, Scope: PlatformRootScope},
		Grant{Action: ActionCloudProviderRead, Scope: PlatformRootScope},
	)

	for _, providerUID := range cfg.CloudProviderUIDs {
		if providerUID == "" {
			continue
		}
		scope := apimeta.ScopeIdentity{Kind: apimeta.ScopeCloudProvider, UID: providerUID}
		grants = append(grants,
			Grant{Action: ActionTopologyWrite, Scope: scope},
			Grant{Action: ActionTopologyRead, Scope: scope},
			Grant{Action: ActionParticipationWithdrawProvider, Scope: scope},
			Grant{Action: ActionParticipationAcceptReleaseProvider, Scope: scope},
			Grant{Action: ActionParticipationDeclineReleaseProvider, Scope: scope},
		)
	}

	for _, platformUID := range cfg.CloudPlatformUIDs {
		if platformUID == "" {
			continue
		}
		scope := apimeta.ScopeIdentity{Kind: apimeta.ScopeCloudPlatform, UID: platformUID}
		grants = append(grants,
			Grant{Action: ActionParticipationAcceptPlatform, Scope: scope},
			Grant{Action: ActionParticipationRejectPlatform, Scope: scope},
			Grant{Action: ActionParticipationRequestReleasePlatform, Scope: scope},
			Grant{Action: ActionParticipationRead, Scope: scope},
		)
	}

	for _, rel := range cfg.ParticipationRelations {
		if rel.CloudProviderUID == "" || rel.CloudPlatformUID == "" {
			continue
		}
		grants = append(grants, Grant{
			Action:                  ActionParticipationRequestProvider,
			Scope:                   apimeta.ScopeIdentity{Kind: apimeta.ScopeCloudProvider, UID: rel.CloudProviderUID},
			RelatedCloudPlatformUID: rel.CloudPlatformUID,
		})
	}

	for _, platformUID := range cfg.SuspendResume.PlatformUIDs {
		if platformUID == "" {
			continue
		}
		scope := apimeta.ScopeIdentity{Kind: apimeta.ScopeCloudPlatform, UID: platformUID}
		grants = append(grants,
			Grant{Action: ActionSuspendPlatform, Scope: scope},
			Grant{Action: ActionResumePlatform, Scope: scope},
		)
	}
	for _, providerUID := range cfg.SuspendResume.ProviderUIDs {
		if providerUID == "" {
			continue
		}
		scope := apimeta.ScopeIdentity{Kind: apimeta.ScopeCloudProvider, UID: providerUID}
		grants = append(grants,
			Grant{Action: ActionSuspendProvider, Scope: scope},
			Grant{Action: ActionResumeProvider, Scope: scope},
		)
	}

	return &BootstrapGrantResolver{
		principalID: cfg.PrincipalID,
		grants:      grants,
	}
}

// PrincipalID returns the configured bootstrap principal identifier.
func (r *BootstrapGrantResolver) PrincipalID() string {
	if r == nil {
		return ""
	}
	return r.principalID
}

// Grants returns a copy of the current server-resolved grants. Callers must
// not treat the result as persisted roles, memberships, or assignments.
func (r *BootstrapGrantResolver) Grants() []Grant {
	if r == nil || len(r.grants) == 0 {
		return nil
	}
	out := make([]Grant, len(r.grants))
	copy(out, r.grants)
	return out
}

// PersistsRolesMembershipsOrAssignments reports whether this resolver stores
// FEATURE-0018 governance objects. FEATURE-0015 always returns false.
func (r *BootstrapGrantResolver) PersistsRolesMembershipsOrAssignments() bool {
	return false
}

// CoarseLookup returns grants whose action equals the route-required action
// (design §5.1 step 5). Scope and resourceUID restrictions are not evaluated.
// Header/body claims never expand or substitute for this result.
func (r *BootstrapGrantResolver) CoarseLookup(action string) []Grant {
	if r == nil || action == "" {
		return nil
	}
	var out []Grant
	for _, g := range r.grants {
		if g.Action == action {
			out = append(out, g)
		}
	}
	return out
}

// AuthorizeExact evaluates exact scope and optional resourceUID restriction
// (design §5.1 step 7). An empty resourceUID means collection-create/LIST
// (no target narrowing). A grant with a non-empty ResourceUID authorizes only
// that target. participation.request.provider additionally requires the
// RelatedCloudPlatformUID relation when relatedPlatformUID is non-empty.
func (r *BootstrapGrantResolver) AuthorizeExact(
	action string,
	scope apimeta.ScopeIdentity,
	resourceUID string,
	relatedPlatformUID string,
) bool {
	if r == nil || action == "" {
		return false
	}
	for _, g := range r.grants {
		if g.Action != action {
			continue
		}
		if g.Scope != scope {
			continue
		}
		if g.ResourceUID != "" && g.ResourceUID != resourceUID {
			continue
		}
		if action == ActionParticipationRequestProvider {
			if relatedPlatformUID != "" && g.RelatedCloudPlatformUID != relatedPlatformUID {
				continue
			}
		}
		return true
	}
	return false
}

// NarrowGrant returns a copy of the matching unrestricted grant with
// ResourceUID set for PATCH/GET target narrowing. It does not mutate the
// resolver's stored grants.
func (r *BootstrapGrantResolver) NarrowGrant(action string, scope apimeta.ScopeIdentity, resourceUID string) (Grant, bool) {
	if r == nil || resourceUID == "" {
		return Grant{}, false
	}
	for _, g := range r.grants {
		if g.Action != action || g.Scope != scope {
			continue
		}
		if g.ResourceUID != "" && g.ResourceUID != resourceUID {
			continue
		}
		out := g
		out.ResourceUID = resourceUID
		return out, true
	}
	return Grant{}, false
}

// ResolveSuspendResume resolves exactly one matching current scoped
// suspend/resume grant via the lifecycle fail-closed helper (design §6.1).
func (r *BootstrapGrantResolver) ResolveSuspendResume(
	action ParticipationAction,
	platformUID, providerUID string,
) (HoldParty, *apiproblem.Problem) {
	var scoped []ScopedActionGrant
	if r != nil {
		for _, g := range r.grants {
			switch g.Action {
			case ActionSuspendPlatform, ActionSuspendProvider, ActionResumePlatform, ActionResumeProvider:
				scoped = append(scoped, g.ScopedActionGrant())
			}
		}
	}
	return ResolveSuspendResumeParty(action, platformUID, providerUID, scoped)
}

// HasForgedGrantHeader reports whether the reserved forged-grant header is
// present with any value. The claim value is never trusted or returned.
func HasForgedGrantHeader(h http.Header) bool {
	if h == nil {
		return false
	}
	_, ok := h[http.CanonicalHeaderKey(HeaderBootstrapGrant)]
	return ok
}

// HasForgedGrantBodyMember reports whether a syntactically valid duplicate-free
// decoded object root contains the reserved top-level bootstrapGrant member.
// Callers must only invoke this after phase-one syntax/duplicate checks succeed.
// The claim value is never trusted, parsed as authorization, or disclosed.
func HasForgedGrantBodyMember(root map[string]json.RawMessage) bool {
	if root == nil {
		return false
	}
	_, ok := root[BodyMemberBootstrapGrant]
	return ok
}

// DenyForgedGrant returns AUTHORIZATION_DENIED for a forged-grant carrier.
// The claim itself is ignored and never disclosed (VS0-CF-F15-22).
func DenyForgedGrant() *apiproblem.Problem {
	return apiproblem.New(apiproblem.CodeAuthorizationDenied).WithDetail(
		"forged bootstrap grant claim is denied",
	)
}

// EvaluateForgedGrantCarriers enforces ADH-2026-054 forged-grant detection for
// the grant resolver. Header presence is checked first. Body-member detection
// is only considered when headerPresent is false (callers that already decoded
// a valid duplicate-free body). Both paths ignore claim values and return
// AUTHORIZATION_DENIED. A nil problem means no forged carrier was present.
func EvaluateForgedGrantCarriers(headerPresent, bodyMemberPresent bool) *apiproblem.Problem {
	if headerPresent || bodyMemberPresent {
		return DenyForgedGrant()
	}
	return nil
}

// NewForgedGrantDenial builds the audited AUTHORIZATION_DENIED outcome for a
// forged-grant attempt. The forged claim value must not appear in the event.
func NewForgedGrantDenial(actor apimeta.TypedRef, requestID, auditUID string) AuditedDenial {
	return AuditedDenial{
		Problem: DenyForgedGrant(),
		Event: NewRedactedAuditEvent(RedactedAuditInput{
			UID:       auditUID,
			RequestID: requestID,
			Action:    AuditActionForgedGrant,
			Outcome:   apiconform.AuditOutcomeDenied,
			Reason:    "forged_grant_denied",
			Actor:     actor,
		}),
	}
}
