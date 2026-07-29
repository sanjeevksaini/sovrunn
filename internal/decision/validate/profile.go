package validate

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/apivalid"
	"github.com/sanjeevksaini/sovrunn/internal/decision"
)

// RFC 6901 pointers for the profile validation pass (design §7.3, §7.7, §9.1 step 3).
const (
	ptrProfileRef            = "/record/profileRef"
	ptrProfileRefName        = "/record/profileRef/name"
	ptrProfileRefVersion     = "/record/profileRef/version"
	ptrSpecLifecycleStatus   = "/spec/lifecycleStatus"
	ptrSpecName              = "/spec/name"
	ptrSpecFamily            = "/spec/family"
	ptrSpecVersion           = "/spec/version"
	ptrSpecOwner             = "/spec/owner"
	ptrSpecPrimaryForm       = "/spec/primaryForm"
	ptrSpecAuthorityLevels   = "/spec/allowedAuthorityLevels"
	ptrSpecFailPosture       = "/spec/failureBehavior/failPosture"
	ptrSpecSecurityException = "/spec/failureBehavior/securityExceptionRef"
	ptrSpecExceptionUID      = "/spec/failureBehavior/securityExceptionRef/approvalRef/uid"
	ptrSpecExceptionOwner    = "/spec/failureBehavior/securityExceptionRef/ownerRef"
	ptrSpecExceptionScope    = "/spec/failureBehavior/securityExceptionRef/exceptionScopeRef"
	ptrSpecExceptionUntil    = "/spec/failureBehavior/securityExceptionRef/effectiveUntil"
	ptrSpecExceptionControls = "/spec/failureBehavior/securityExceptionRef/compensatingControls"
	ptrSpecExceptionModes    = "/spec/failureBehavior/securityExceptionRef/coveredFailureModes"
	ptrSpecLimits            = "/spec/limits"
)

// LifecycleStatusActive is the only lifecycle status accepted as active
// (F13-PROF-007). Any other non-empty value is treated as inactive and
// rejected with DECISION_PROFILE_INACTIVE.
const LifecycleStatusActive = "Active"

// ProfileInput is the input to ValidateProfile (design §8 / §9.1 step 3).
//
// The design §8 summary signature is ValidateProfile(p DecisionProfile).
// Lookup fields enable F13-PROF-007 unknown/version-unsupported codes when
// resolving a record profileRef against a derivative registry view; DecisionTime
// enables SecurityExceptionRef expiry-vs-decision-time checks (design §7.7).
type ProfileInput struct {
	Profile decision.DecisionProfile

	// CheckLookup enables registry-resolution codes (F13-PROF-007).
	// When true:
	//   - !IDExists → DECISION_PROFILE_UNKNOWN
	//   - IDExists && !Found → DECISION_PROFILE_VERSION_UNSUPPORTED
	//   - Found → validate Profile (including inactive lifecycle)
	CheckLookup bool
	Ref         decision.ProfileRef
	Found       bool // exact (id, version) present in the registry view
	IDExists    bool // any version of Ref.Name exists

	// DecisionTime is the captured/effective decision time (UTC RFC3339).
	// Empty skips the "before captured/effective decision time" expiry check.
	DecisionTime string
}

// PlatformLimitCeilings are the FEATURE-0012 absolute outer bounds that
// constrain comparable DecisionProfile limit fields (F13-PROF-005; design §1.2).
// Graph-only ceilings (nodes/edges/width/fan-out/calls/…) have no FEATURE-0012
// numeric outer bound and must not invent a FEATURE-0013 default (architecture §9).
type PlatformLimitCeilings struct {
	MaxPayloadBytes int64
	MaxDepth        int
	MaxFieldCount   int
}

// InheritedPlatformCeilings returns the FEATURE-0012-derived outer bounds used
// by profile limit hierarchy checks. Values come from apivalid.DefaultLimits:
// MaxObjectBytes → MaxPayloadBytes, MaxNestingDepth → MaxDepth, MaxLabels →
// MaxFieldCount.
func InheritedPlatformCeilings() PlatformLimitCeilings {
	lim := apivalid.DefaultLimits()
	return PlatformLimitCeilings{
		MaxPayloadBytes: int64(lim.MaxObjectBytes),
		MaxDepth:        lim.MaxNestingDepth,
		MaxFieldCount:   lim.MaxLabels,
	}
}

// EffectiveIntLimit returns the most restrictive positive bound among
// platform, profile, and optional evaluator values (design §1.2). Zero means
// unset and is ignored. When no positive value is present, returns 0.
func EffectiveIntLimit(platform, profile, evaluator int) int {
	effective := 0
	for _, v := range []int{platform, profile, evaluator} {
		if v <= 0 {
			continue
		}
		if effective == 0 || v < effective {
			effective = v
		}
	}
	return effective
}

// EffectiveInt64Limit is EffectiveIntLimit for int64 ceilings.
func EffectiveInt64Limit(platform, profile, evaluator int64) int64 {
	effective := int64(0)
	for _, v := range []int64{platform, profile, evaluator} {
		if v <= 0 {
			continue
		}
		if effective == 0 || v < effective {
			effective = v
		}
	}
	return effective
}

// ValidateProfile runs the FEATURE-0013 profile validation pass (design §9.1
// step 3; F13-PROF-001…008; F13-EVAL-02/03/04; AD-029, AD-040, AD-043).
//
// Checks (fail closed, short-circuit on first blocking violation):
//  1. optional registry lookup (unknown / version-unsupported)
//  2. identity and required schema fields
//  3. lifecycle active status
//  4. failure posture and SecurityExceptionRef structural evidence
//  5. inherited limit hierarchy (FEATURE-0012 outer → profile → evaluator narrow-only)
//
// SecurityExceptionRef is structural evidence only; no approval workflow,
// issuance, revocation, or runtime enforcement is performed.
func ValidateProfile(in ProfileInput) *apiproblem.Problem {
	if in.CheckLookup {
		if p := checkProfileLookup(in); p != nil {
			return p
		}
	}

	profile := in.Profile
	if p := checkProfileIdentity(profile); p != nil {
		return p
	}
	if p := checkProfileLifecycle(profile); p != nil {
		return p
	}
	if p := checkProfileSchema(profile); p != nil {
		return p
	}
	if p := checkFailureBehavior(profile, in.DecisionTime); p != nil {
		return p
	}
	if p := checkProfileLimits(profile); p != nil {
		return p
	}
	return nil
}

// ValidateProfileDocument is the design §8 convenience form that validates a
// DecisionProfile document without registry lookup.
func ValidateProfileDocument(p decision.DecisionProfile) *apiproblem.Problem {
	return ValidateProfile(ProfileInput{Profile: p})
}

func checkProfileLookup(in ProfileInput) *apiproblem.Problem {
	if !in.IDExists {
		return profileProblem(CodeProfileUnknown, ptrProfileRefName,
			"decision profile is unknown; accepting arbitrary profile URIs is prohibited")
	}
	if !in.Found {
		return profileProblem(CodeProfileVersionUnsupported, ptrProfileRefVersion,
			"decision profile version is unsupported; no fallback to another version is permitted")
	}
	return nil
}

func checkProfileIdentity(p decision.DecisionProfile) *apiproblem.Problem {
	if p.APIVersion != "" && p.APIVersion != decision.APIVersionDecisionProfile {
		return profileProblem(CodeProfileSchemaInvalid, "/apiVersion",
			"DecisionProfile apiVersion must be governance.sovrunn.io/v1alpha1")
	}
	if p.Kind != "" && p.Kind != decision.KindDecisionProfile {
		return profileProblem(CodeProfileSchemaInvalid, "/kind",
			"DecisionProfile kind must be DecisionProfile")
	}
	if strings.TrimSpace(p.Spec.Name) == "" {
		return profileProblem(CodeProfileSchemaInvalid, ptrSpecName,
			"spec.name is required and must be non-empty")
	}
	if strings.TrimSpace(p.Spec.Family) == "" {
		return profileProblem(CodeProfileSchemaInvalid, ptrSpecFamily,
			"spec.family is required and must be non-empty")
	}
	if strings.TrimSpace(p.Spec.Version) == "" {
		return profileProblem(CodeProfileSchemaInvalid, ptrSpecVersion,
			"spec.version is required and must be non-empty")
	}
	if strings.TrimSpace(p.Spec.Owner) == "" {
		return profileProblem(CodeProfileSchemaInvalid, ptrSpecOwner,
			"spec.owner is required and must be non-empty")
	}
	if strings.TrimSpace(p.Spec.LifecycleStatus) == "" {
		return profileProblem(CodeProfileSchemaInvalid, ptrSpecLifecycleStatus,
			"spec.lifecycleStatus is required and must be non-empty")
	}
	return nil
}

func checkProfileLifecycle(p decision.DecisionProfile) *apiproblem.Problem {
	if p.Spec.LifecycleStatus == LifecycleStatusActive {
		return nil
	}
	return profileProblem(CodeProfileInactive, ptrSpecLifecycleStatus,
		"decision profile is inactive; only lifecycleStatus Active is accepted")
}

func checkProfileSchema(p decision.DecisionProfile) *apiproblem.Problem {
	if !p.Spec.PrimaryForm.Valid() {
		return profileProblem(CodeProfileSchemaInvalid, ptrSpecPrimaryForm,
			"spec.primaryForm must be a closed DecisionForm value")
	}
	for i, facet := range p.Spec.SecondaryFacets {
		if !facet.Valid() {
			return profileProblem(CodeProfileSchemaInvalid,
				fmt.Sprintf("/spec/secondaryFacets/%d", i),
				"spec.secondaryFacets entries must be closed DecisionForm values")
		}
	}
	if len(p.Spec.AllowedAuthorityLevels) == 0 {
		return profileProblem(CodeProfileSchemaInvalid, ptrSpecAuthorityLevels,
			"spec.allowedAuthorityLevels is required and must be non-empty")
	}
	for i, auth := range p.Spec.AllowedAuthorityLevels {
		if !auth.Valid() {
			return profileProblem(CodeProfileSchemaInvalid,
				fmt.Sprintf("/spec/allowedAuthorityLevels/%d", i),
				"spec.allowedAuthorityLevels entries must be closed DecisionAuthority values")
		}
	}
	if strings.TrimSpace(p.Spec.InputSchemaRef) == "" {
		return profileProblem(CodeProfileSchemaInvalid, "/spec/inputSchemaRef",
			"spec.inputSchemaRef is required and must be non-empty")
	}
	if strings.TrimSpace(p.Spec.TypedResultSchemaRef) == "" {
		return profileProblem(CodeProfileSchemaInvalid, "/spec/typedResultSchemaRef",
			"spec.typedResultSchemaRef is required and must be non-empty")
	}
	if strings.TrimSpace(p.Spec.RationaleSchemaRef) == "" {
		return profileProblem(CodeProfileSchemaInvalid, "/spec/rationaleSchemaRef",
			"spec.rationaleSchemaRef is required and must be non-empty")
	}
	if strings.TrimSpace(p.Spec.ObligationSchemaRef) == "" {
		return profileProblem(CodeProfileSchemaInvalid, "/spec/obligationSchemaRef",
			"spec.obligationSchemaRef is required and must be non-empty")
	}
	if strings.TrimSpace(p.Spec.IdempotencyIdentity) == "" {
		return profileProblem(CodeProfileSchemaInvalid, "/spec/idempotencyIdentity",
			"spec.idempotencyIdentity is required and must be non-empty")
	}
	if !p.Spec.SensitivityCeiling.Valid() {
		return profileProblem(CodeProfileSchemaInvalid, "/spec/sensitivityCeiling",
			"spec.sensitivityCeiling must be a closed Sensitivity value")
	}
	if !p.Spec.FailureBehavior.FailPosture.Valid() {
		return profileProblem(CodeProfileSchemaInvalid, ptrSpecFailPosture,
			"spec.failureBehavior.failPosture must be FAIL_CLOSED, REQUIRES_APPROVAL, or FAIL_OPEN")
	}
	return nil
}

func checkFailureBehavior(p decision.DecisionProfile, decisionTime string) *apiproblem.Problem {
	fb := p.Spec.FailureBehavior
	securityClass := isSecurityClassFamily(p.Spec.Family)

	if securityClass && fb.FailPosture == decision.FailPostureFailOpen && fb.SecurityExceptionRef == nil {
		return profileProblem(CodeProfileSchemaInvalid, ptrSpecSecurityException,
			"security-class fail-open requires a SecurityExceptionRef structural evidence carrier")
	}
	if fb.FailPosture == decision.FailPostureFailOpen && fb.SecurityExceptionRef == nil {
		return profileProblem(CodeProfileSchemaInvalid, ptrSpecSecurityException,
			"fail-open requires a SecurityExceptionRef structural evidence carrier")
	}
	if fb.SecurityExceptionRef == nil {
		return nil
	}
	return checkSecurityExceptionRef(p, fb.SecurityExceptionRef, decisionTime)
}

func checkSecurityExceptionRef(p decision.DecisionProfile, ref *decision.SecurityExceptionRef, decisionTime string) *apiproblem.Problem {
	if ref == nil {
		return profileProblem(CodeProfileSchemaInvalid, ptrSpecSecurityException,
			"SecurityExceptionRef is malformed or missing")
	}

	if !typedRefUIDPinned(ref.ApprovalRef) {
		field := ptrSpecExceptionUID
		if ref.ApprovalRef.UID != "" && (ref.ApprovalRef.APIVersion == "" || ref.ApprovalRef.Kind == "" || ref.ApprovalRef.Name == "") {
			field = "/spec/failureBehavior/securityExceptionRef/approvalRef"
		}
		return profileProblem(CodeProfileSchemaInvalid, field,
			"SecurityExceptionRef.approvalRef requires apiVersion, kind, name, and uid pinning")
	}
	if !typedRefUIDPinned(ref.OwnerRef) {
		return profileProblem(CodeProfileSchemaInvalid, ptrSpecExceptionOwner,
			"SecurityExceptionRef.ownerRef requires apiVersion, kind, name, and uid pinning")
	}
	if !typedRefUIDPinned(ref.ApprovingAuthorityRef) {
		return profileProblem(CodeProfileSchemaInvalid,
			"/spec/failureBehavior/securityExceptionRef/approvingAuthorityRef",
			"SecurityExceptionRef.approvingAuthorityRef requires apiVersion, kind, name, and uid pinning")
	}

	// Owner mismatch vs governing profile owner identity (design §7.7).
	if ref.OwnerRef.Name != p.Spec.Owner {
		return profileProblem(CodeProfileSchemaInvalid, ptrSpecExceptionOwner,
			"SecurityExceptionRef.ownerRef must match the governing profile owner")
	}

	if p := checkExceptionScope(p, ref.ExceptionScopeRef); p != nil {
		return p
	}

	if len(ref.CompensatingControls) == 0 || !nonEmptyStrings(ref.CompensatingControls) {
		return profileProblem(CodeProfileSchemaInvalid, ptrSpecExceptionControls,
			"SecurityExceptionRef.compensatingControls must be a non-empty list of control identifiers")
	}
	if len(ref.CoveredFailureModes) == 0 || !nonEmptyStrings(ref.CoveredFailureModes) {
		return profileProblem(CodeProfileSchemaInvalid, ptrSpecExceptionModes,
			"SecurityExceptionRef.coveredFailureModes must be a non-empty list")
	}
	for i, mode := range ref.CoveredFailureModes {
		if !containsString(p.Spec.FailureBehavior.DeclaredFailureModes, mode) {
			return profileProblem(CodeProfileSchemaInvalid,
				ptrSpecExceptionModes+"/"+strconv.Itoa(i),
				"SecurityExceptionRef.coveredFailureModes entry is not a profile-declared failure mode")
		}
	}

	if strings.TrimSpace(ref.AuditTreatment) == "" {
		return profileProblem(CodeProfileSchemaInvalid,
			"/spec/failureBehavior/securityExceptionRef/auditTreatment",
			"SecurityExceptionRef.auditTreatment is required")
	}
	if strings.TrimSpace(ref.ReassessmentTrigger) == "" {
		return profileProblem(CodeProfileSchemaInvalid,
			"/spec/failureBehavior/securityExceptionRef/reassessmentTrigger",
			"SecurityExceptionRef.reassessmentTrigger is required")
	}
	if strings.TrimSpace(ref.Purpose) == "" {
		return profileProblem(CodeProfileSchemaInvalid,
			"/spec/failureBehavior/securityExceptionRef/purpose",
			"SecurityExceptionRef.purpose is required")
	}

	from, errFrom := time.Parse(time.RFC3339, ref.EffectiveFrom)
	until, errUntil := time.Parse(time.RFC3339, ref.EffectiveUntil)
	if errFrom != nil || errUntil != nil {
		return profileProblem(CodeProfileSchemaInvalid, ptrSpecExceptionUntil,
			"SecurityExceptionRef effective interval must use UTC RFC3339 timestamps")
	}
	if !until.After(from) {
		return profileProblem(CodeProfileSchemaInvalid, ptrSpecExceptionUntil,
			"SecurityExceptionRef.effectiveUntil must be after effectiveFrom")
	}
	if decisionTime != "" {
		at, err := time.Parse(time.RFC3339, decisionTime)
		if err != nil {
			return profileProblem(CodeProfileSchemaInvalid, ptrSpecExceptionUntil,
				"captured/effective decision time must be UTC RFC3339 for exception expiry checks")
		}
		if at.After(until) || at.Equal(until) {
			return profileProblem(CodeProfileSchemaInvalid, ptrSpecExceptionUntil,
				"SecurityExceptionRef is expired relative to captured/effective decision time")
		}
		if at.Before(from) {
			return profileProblem(CodeProfileSchemaInvalid, ptrSpecExceptionUntil,
				"SecurityExceptionRef is not yet effective at captured/effective decision time")
		}
	}
	return nil
}

func checkExceptionScope(p decision.DecisionProfile, exceptionScope *apimeta.ScopeRef) *apiproblem.Problem {
	governing := apimeta.NormalizeScope(p.Metadata.ScopeRef)
	normalizedException := apimeta.NormalizeScope(exceptionScope)

	if normalizedException != nil {
		kind := apimeta.ScopeKind(normalizedException.Kind)
		if !kind.Valid() {
			return profileProblem(CodeProfileSchemaInvalid, ptrSpecExceptionScope,
				"SecurityExceptionRef.exceptionScopeRef.kind must be a FEATURE-0012 ScopeKind")
		}
		if normalizedException.UID == "" {
			return profileProblem(CodeProfileSchemaInvalid, ptrSpecExceptionScope,
				"non-Platform SecurityExceptionRef.exceptionScopeRef requires uid")
		}
		if !containsScopeKind(p.Spec.ScopeSemantics.AllowedScopes, kind) && (kind != apimeta.ScopePlatform || !p.Spec.ScopeSemantics.PlatformAllowed) {
			return profileProblem(CodeProfileSchemaInvalid, ptrSpecExceptionScope,
				"SecurityExceptionRef.exceptionScopeRef is not compatible with governing profile allowedScopes")
		}
	}

	govID := apimeta.CanonicalScopeIdentity(governing)
	excID := apimeta.CanonicalScopeIdentity(normalizedException)
	if govID != excID {
		return profileProblem(CodeProfileSchemaInvalid, ptrSpecExceptionScope,
			"SecurityExceptionRef.exceptionScopeRef must be compatible with the governing profile metadata.scopeRef")
	}
	return nil
}

func checkProfileLimits(p decision.DecisionProfile) *apiproblem.Problem {
	plat := InheritedPlatformCeilings()
	lim := p.Spec.Limits

	if p := checkLimitNonNegative(lim, ptrSpecLimits); p != nil {
		return p
	}

	// FEATURE-0012 absolute outer bounds for comparable fields.
	if lim.MaxPayloadBytes > 0 && lim.MaxPayloadBytes > plat.MaxPayloadBytes {
		return profileProblem(CodeProfileLimitInvalid, ptrSpecLimits+"/maxPayloadBytes",
			"spec.limits.maxPayloadBytes exceeds the FEATURE-0012 MaxObjectBytes outer bound")
	}
	if lim.MaxDepth > 0 && lim.MaxDepth > plat.MaxDepth {
		return profileProblem(CodeProfileLimitInvalid, ptrSpecLimits+"/maxDepth",
			"spec.limits.maxDepth exceeds the FEATURE-0012 MaxNestingDepth outer bound")
	}
	if lim.MaxFieldCount > 0 && lim.MaxFieldCount > plat.MaxFieldCount {
		return profileProblem(CodeProfileLimitInvalid, ptrSpecLimits+"/maxFieldCount",
			"spec.limits.maxFieldCount exceeds the FEATURE-0012 MaxLabels outer bound")
	}

	// Mandatory finite values when composition is declared (architecture §9 /
	// F13-PROF-006): missing mandatory bounds fail closed. No FEATURE-0013
	// default is invented for unset optional domains.
	if len(p.Spec.AcceptedCompositionStrategies) > 0 {
		if lim.MaxNodes <= 0 {
			return profileProblem(CodeProfileLimitInvalid, ptrSpecLimits+"/maxNodes",
				"spec.limits.maxNodes is a mandatory positive bound when composition strategies are declared")
		}
		if lim.MaxEdges <= 0 {
			return profileProblem(CodeProfileLimitInvalid, ptrSpecLimits+"/maxEdges",
				"spec.limits.maxEdges is a mandatory positive bound when composition strategies are declared")
		}
		if lim.MaxDepth <= 0 {
			return profileProblem(CodeProfileLimitInvalid, ptrSpecLimits+"/maxDepth",
				"spec.limits.maxDepth is a mandatory positive bound when composition strategies are declared")
		}
		if lim.MaxWidth <= 0 {
			return profileProblem(CodeProfileLimitInvalid, ptrSpecLimits+"/maxWidth",
				"spec.limits.maxWidth is a mandatory positive bound when composition strategies are declared")
		}
		if lim.MaxFanOut <= 0 {
			return profileProblem(CodeProfileLimitInvalid, ptrSpecLimits+"/maxFanOut",
				"spec.limits.maxFanOut is a mandatory positive bound when composition strategies are declared")
		}
		if lim.MaxPayloadBytes <= 0 {
			return profileProblem(CodeProfileLimitInvalid, ptrSpecLimits+"/maxPayloadBytes",
				"spec.limits.maxPayloadBytes is a mandatory positive bound when composition strategies are declared")
		}
	}

	for i, et := range p.Spec.AcceptedEvaluationTypes {
		if et.Type == "" {
			return profileProblem(CodeProfileSchemaInvalid,
				fmt.Sprintf("/spec/acceptedEvaluationTypes/%d/type", i),
				"acceptedEvaluationTypes.type is required")
		}
		if et.Limits == nil {
			continue
		}
		base := fmt.Sprintf("/spec/acceptedEvaluationTypes/%d/limits", i)
		if p := checkLimitNonNegative(*et.Limits, base); p != nil {
			return p
		}
		if p := checkEvaluatorNarrows(lim, *et.Limits, plat, base); p != nil {
			return p
		}
	}
	return nil
}

func checkEvaluatorNarrows(profile, eval decision.ProfileLimits, plat PlatformLimitCeilings, base string) *apiproblem.Problem {
	// Evaluator may only narrow; never widen profile or FEATURE-0012 ceilings.
	checks := []struct {
		name     string
		eval     int64
		profile  int64
		platform int64
	}{
		{"maxNodes", int64(eval.MaxNodes), int64(profile.MaxNodes), 0},
		{"maxEdges", int64(eval.MaxEdges), int64(profile.MaxEdges), 0},
		{"maxDepth", int64(eval.MaxDepth), int64(profile.MaxDepth), int64(plat.MaxDepth)},
		{"maxWidth", int64(eval.MaxWidth), int64(profile.MaxWidth), 0},
		{"maxFanOut", int64(eval.MaxFanOut), int64(profile.MaxFanOut), 0},
		{"maxPayloadBytes", eval.MaxPayloadBytes, profile.MaxPayloadBytes, plat.MaxPayloadBytes},
		{"maxFieldCount", int64(eval.MaxFieldCount), int64(profile.MaxFieldCount), int64(plat.MaxFieldCount)},
		{"maxCalls", int64(eval.MaxCalls), int64(profile.MaxCalls), 0},
		{"maxConcurrency", int64(eval.MaxConcurrency), int64(profile.MaxConcurrency), 0},
		{"maxRetries", int64(eval.MaxRetries), int64(profile.MaxRetries), 0},
		{"maxLatencyMs", eval.MaxLatencyMs, profile.MaxLatencyMs, 0},
		{"maxJurisdictions", int64(eval.MaxJurisdictions), int64(profile.MaxJurisdictions), 0},
	}
	for _, c := range checks {
		if c.eval <= 0 {
			continue
		}
		if c.platform > 0 && c.eval > c.platform {
			return profileProblem(CodeProfileLimitInvalid, base+"/"+c.name,
				"evaluator limit exceeds the FEATURE-0012 outer bound")
		}
		if c.profile <= 0 {
			return profileProblem(CodeProfileLimitInvalid, base+"/"+c.name,
				"evaluator limit is incomparable or missing a governing profile mandatory bound")
		}
		if c.eval > c.profile {
			return profileProblem(CodeProfileLimitInvalid, base+"/"+c.name,
				"evaluator registration may only narrow profile limits; widening is prohibited")
		}
	}

	if eval.MaxCostBudget != "" {
		if profile.MaxCostBudget == "" {
			return profileProblem(CodeProfileLimitInvalid, base+"/maxCostBudget",
				"evaluator maxCostBudget is incomparable without a governing profile cost budget")
		}
		if eval.MaxCostBudget != profile.MaxCostBudget {
			return profileProblem(CodeProfileLimitInvalid, base+"/maxCostBudget",
				"evaluator maxCostBudget must equal the profile opaque budget tag; opaque tags are not ordered")
		}
	}
	return nil
}

func checkLimitNonNegative(lim decision.ProfileLimits, base string) *apiproblem.Problem {
	switch {
	case lim.MaxNodes < 0:
		return profileProblem(CodeProfileLimitInvalid, base+"/maxNodes", "limit must not be negative")
	case lim.MaxEdges < 0:
		return profileProblem(CodeProfileLimitInvalid, base+"/maxEdges", "limit must not be negative")
	case lim.MaxDepth < 0:
		return profileProblem(CodeProfileLimitInvalid, base+"/maxDepth", "limit must not be negative")
	case lim.MaxWidth < 0:
		return profileProblem(CodeProfileLimitInvalid, base+"/maxWidth", "limit must not be negative")
	case lim.MaxFanOut < 0:
		return profileProblem(CodeProfileLimitInvalid, base+"/maxFanOut", "limit must not be negative")
	case lim.MaxPayloadBytes < 0:
		return profileProblem(CodeProfileLimitInvalid, base+"/maxPayloadBytes", "limit must not be negative")
	case lim.MaxFieldCount < 0:
		return profileProblem(CodeProfileLimitInvalid, base+"/maxFieldCount", "limit must not be negative")
	case lim.MaxCalls < 0:
		return profileProblem(CodeProfileLimitInvalid, base+"/maxCalls", "limit must not be negative")
	case lim.MaxConcurrency < 0:
		return profileProblem(CodeProfileLimitInvalid, base+"/maxConcurrency", "limit must not be negative")
	case lim.MaxRetries < 0:
		return profileProblem(CodeProfileLimitInvalid, base+"/maxRetries", "limit must not be negative")
	case lim.MaxLatencyMs < 0:
		return profileProblem(CodeProfileLimitInvalid, base+"/maxLatencyMs", "limit must not be negative")
	case lim.MaxJurisdictions < 0:
		return profileProblem(CodeProfileLimitInvalid, base+"/maxJurisdictions", "limit must not be negative")
	}
	return nil
}

func isSecurityClassFamily(family string) bool {
	switch strings.ToLower(strings.TrimSpace(family)) {
	case "security", "authorization", "admission", "sovereignty", "compliance",
		"privileged-access", "data-movement", "break-glass":
		return true
	default:
		return false
	}
}

func typedRefUIDPinned(ref apimeta.TypedRef) bool {
	return ref.APIVersion != "" && ref.Kind != "" && ref.Name != "" && ref.UID != ""
}

func nonEmptyStrings(vals []string) bool {
	if len(vals) == 0 {
		return false
	}
	for _, v := range vals {
		if strings.TrimSpace(v) == "" {
			return false
		}
	}
	return true
}

func containsString(set []string, v string) bool {
	for _, s := range set {
		if s == v {
			return true
		}
	}
	return false
}

func profileProblem(code apiproblem.ViolationCode, field, message string) *apiproblem.Problem {
	return apiproblem.New(apiproblem.CodeValidationFailed).
		WithDetail(message).
		WithViolations([]apiproblem.Violation{{
			Field:   field,
			Code:    code,
			Message: message,
		}})
}
