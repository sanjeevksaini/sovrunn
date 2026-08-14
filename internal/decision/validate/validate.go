package validate

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/apivalid"
	"github.com/sanjeevksaini/sovrunn/internal/decision"
	"github.com/sanjeevksaini/sovrunn/internal/decision/bundle"
	"github.com/sanjeevksaini/sovrunn/internal/decision/graph"
)

// Fixed nine-pass order (design §9.1; RID-05). Passes short-circuit on the
// first blocking failure. Decode/structural failures use inherited FEATURE-0012
// top-level Problem codes (not DECISION_*). Semantic passes emit
// VALIDATION_FAILED with exactly one closed DECISION_* violations[].code.
const (
	PassDecodeStructural = 1
	PassScope            = 2
	PassProfile          = 3
	PassEvaluation       = 4
	PassComposition      = 5
	PassTrust            = 6
	PassObligation       = 7
	PassRelationship     = 8
	PassSensitivity      = 9
)

// DecisionRecordOptions carries optional orchestration context that cannot be
// derived from DecisionRecord + BundleView alone (design §8 / §9.1).
//
// Zero value is valid: relationship/sensitivity extras are skipped or use
// fail-closed defaults; obligations default to bundle-known identities;
// trust requiredness defaults from profile EvidenceSemantics.RequireEvidence.
type DecisionRecordOptions struct {
	// ParallelScopeFields lists decode-time parallel/second scope sources
	// (F13-SCOPE-005). Empty means sole authority holds.
	ParallelScopeFields []string

	// Established is the append-only relationship graph for pass 8.
	Established []RelationshipNode

	// SupportedObligations maps obligation id → supported. When nil, identities
	// present in the bundle obligations registry are treated as supported.
	SupportedObligations map[string]bool

	// TrustRequired overrides profile-derived trust requiredness. Nil means
	// derive from profile EvidenceSemantics.RequireEvidence.
	TrustRequired      *bool
	ExpectedTrustState string

	// Sensitivity / projection extras (pass 9). Empty classifications and
	// declared sensitivity skip the floor/ceiling content checks; projection
	// rules still require CanonicalView when they declare pointers.
	Classifications      []apimeta.DataClassification
	DeclaredSensitivity  decision.Sensitivity
	ContentCategories    []string
	ProhibitedCategories []string
	CanonicalView        []string
	AllowedAudiences     []string
}

// AuditEventEnvelope carries FEATURE-0012 AuditEvent base fields needed by
// ValidateAuditEvent without importing apiconform (design §5 DAG: validate must
// not import apiconform; apiconform calls validate one-directionally in T-027).
//
// Canonical linkage ownership remains decision.AuditLinkage (closure 7), passed
// separately to ValidateAuditEvent.
type AuditEventEnvelope struct {
	apimeta.TypeMeta
	Metadata apimeta.ObjectMeta     `json:"metadata"`
	Record   AuditEventRecordFields `json:"record"`
}

// AuditEventRecordFields is the FEATURE-0012 AuditEvent record payload subset
// required for structural and scope orchestration (F13-AUDIT-002/003).
type AuditEventRecordFields struct {
	ActorRef   apimeta.TypedRef      `json:"actorRef"`
	RequestID  string                `json:"requestId"`
	SubjectRef apimeta.TypedRef      `json:"subjectRef"`
	Action     string                `json:"action"`
	Outcome    decision.AuditOutcome `json:"outcome"`
	ReasonCode string                `json:"reasonCode"`
}

// AuditEventOptions carries optional AuditEvent orchestration context.
type AuditEventOptions struct {
	ParallelScopeFields []string
}

// DecodeDecisionRecord decodes a DecisionRecord from local bytes using
// apivalid.StrictDecode (design §9.1 pass 1; RID-09). Failures return
// FEATURE-0012 top-level codes (MALFORMED_REQUEST, UNKNOWN_FIELD,
// DUPLICATE_FIELD, REQUEST_TOO_LARGE, …) — never a DECISION_* code.
func DecodeDecisionRecord(data []byte, mode apivalid.DecodeMode) (decision.DecisionRecord, *apiproblem.Problem) {
	var rec decision.DecisionRecord
	if prob := strictDecodeBytes(data, mode, &rec); prob != nil {
		return decision.DecisionRecord{}, prob
	}
	return rec, nil
}

// DecodeAuditEventEnvelope decodes an AuditEventEnvelope from local bytes using
// apivalid.StrictDecode (design §9.1 pass 1; RID-09). Failures return
// FEATURE-0012 top-level codes — never a DECISION_* code.
func DecodeAuditEventEnvelope(data []byte, mode apivalid.DecodeMode) (AuditEventEnvelope, *apiproblem.Problem) {
	var ev AuditEventEnvelope
	if prob := strictDecodeBytes(data, mode, &ev); prob != nil {
		return AuditEventEnvelope{}, prob
	}
	return ev, nil
}

// DecodeAndValidateDecisionRecord runs the fixed nine-pass order including
// decode/structural (design §9.1). On success it returns the decoded
// DecisionRecord (F13-OBJ-007 synchronous completion). Rejection returns
// inherited Problem Details. No pending-decision envelope is produced.
func DecodeAndValidateDecisionRecord(
	data []byte,
	mode apivalid.DecodeMode,
	view bundle.BundleView,
) (decision.DecisionRecord, *apiproblem.Problem) {
	return DecodeAndValidateDecisionRecordWith(data, mode, view, DecisionRecordOptions{})
}

// DecodeAndValidateDecisionRecordWith is DecodeAndValidateDecisionRecord with
// optional orchestration context.
func DecodeAndValidateDecisionRecordWith(
	data []byte,
	mode apivalid.DecodeMode,
	view bundle.BundleView,
	opts DecisionRecordOptions,
) (decision.DecisionRecord, *apiproblem.Problem) {
	rec, prob := DecodeDecisionRecord(data, mode)
	if prob != nil {
		return decision.DecisionRecord{}, prob
	}
	if prob := ValidateDecisionRecordWith(rec, view, opts); prob != nil {
		return decision.DecisionRecord{}, prob
	}
	return rec, nil
}

// ValidateDecisionRecord runs passes 1 (typed structural) through 9 against an
// already-decoded DecisionRecord (design §8 / §9.1; F13-OBJ-007; F13-ERR-001;
// F13-AUDIT-001; AD-009). Optional orchestration context defaults to zero value.
//
// Pass order (short-circuit on first blocking failure):
//  1. decode/structural (typed structural; byte decode is DecodeDecisionRecord)
//  2. scope
//  3. profile
//  4. evaluation
//  5. composition
//  6. trust
//  7. obligation
//  8. relationship
//  9. sensitivity/projection
//
// Completion returns nil (caller already holds the DecisionRecord). Accepted
// async work remains the FEATURE-0012 Operation reference on Correlation —
// FEATURE-0013 produces no PendingDecision envelope.
func ValidateDecisionRecord(rec decision.DecisionRecord, view bundle.BundleView) *apiproblem.Problem {
	return ValidateDecisionRecordWith(rec, view, DecisionRecordOptions{})
}

// ValidateDecisionRecordWith is ValidateDecisionRecord with optional
// orchestration context (relationship chain, obligation support, sensitivity).
func ValidateDecisionRecordWith(rec decision.DecisionRecord, view bundle.BundleView, opts DecisionRecordOptions) *apiproblem.Problem {
	// Pass 1 — typed structural (FEATURE-0012 top-level codes; not DECISION_*).
	if prob := checkDecisionRecordStructural(rec); prob != nil {
		return prob
	}

	ref := rec.Record.ProfileRef
	profile, found := view.Profile(ref.Name, ref.Version)
	idExists := view.ProfileIDExists(ref.Name)

	// Pass 2 — Scope. Profile scope membership is applied when the profile is
	// resolvable; otherwise the seven-value vocabulary and sole-authority checks
	// still run (profile unknown is owned by pass 3).
	allowed, platformAllowed := governanceScopesForRecord(profile, found)
	if prob := ValidateScope(ScopeInput{
		Metadata:            rec.Metadata,
		Allowed:             allowed,
		PlatformAllowed:     platformAllowed,
		ParallelScopeFields: opts.ParallelScopeFields,
		SubjectRefs:         rec.Record.SubjectRefs,
		CanonicalScope:      rec.Record.SemanticIdentity.CanonicalScope,
		CheckCanonicalScope: true,
	}); prob != nil {
		return prob
	}

	// Pass 3 — Profile (lookup + document).
	if prob := ValidateProfile(ProfileInput{
		Profile:      profile,
		CheckLookup:  true,
		Ref:          ref,
		Found:        found,
		IDExists:     idExists,
		DecisionTime: decisionTimeHint(rec),
	}); prob != nil {
		return prob
	}

	recScope := apimeta.CanonicalScopeIdentity(apimeta.NormalizeScope(rec.Metadata.ScopeRef))

	// Pass 4 — Evaluation (each embedded capture).
	for i, er := range rec.Record.EvaluationResults {
		if prob := ValidateEvaluation(er, recScope, profile); prob != nil {
			return prefixProblemFields(prob, "/record/evaluationResults/"+strconv.Itoa(i))
		}
	}

	// Pass 5 — Composition.
	if rec.Record.Composition != nil {
		compIn := CompositionInput{
			Composition: *rec.Record.Composition,
			Profile:     profile,
		}
		if def, ok := resolveCompositionFromView(view, *rec.Record.Composition); ok {
			compIn.ResolvedGraph = &def
		} else if strings.TrimSpace(rec.Record.Composition.GraphRef) != "" {
			// Graph declared but not resolvable from the view → graph invalid
			// at the composition pointer (bundle resolution ownership).
			return compositionProblem(CodeCompositionGraphInvalid, ptrCompositionGraphRef,
				"composition graph reference/version not found in bundle; no version fallback")
		}
		if prob := ValidateComposition(compIn); prob != nil {
			return prob
		}
	}

	// Pass 6 — Trust.
	trustRequired := profile.Spec.EvidenceSemantics.RequireEvidence
	if opts.TrustRequired != nil {
		trustRequired = *opts.TrustRequired
	}
	if prob := ValidateTrust(TrustInput{
		Carrier:       rec.Record.Trust,
		Required:      trustRequired,
		ExpectedState: opts.ExpectedTrustState,
	}); prob != nil {
		return prefixProblemFields(prob, "/record")
	}

	// Pass 7 — Obligations.
	supported := opts.SupportedObligations
	if supported == nil {
		supported = view.KnownObligations()
	}
	if prob := ValidateObligations(rec.Record.Result.Obligations, supported); prob != nil {
		return prefixProblemFields(prob, "/record/result")
	}

	// Pass 8 — Relationship.
	if prob := ValidateRelationship(RelationshipInput{
		Relationship: rec.Record.Relationship,
		Self: apimeta.TypedRef{
			APIVersion: rec.APIVersion,
			Kind:       rec.Kind,
			Name:       rec.Metadata.Name,
			UID:        rec.Metadata.UID,
		},
		Established: opts.Established,
		MaxChain:    profile.Spec.ValidityRules.MaxChainDepth,
	}); prob != nil {
		return prob
	}

	// Pass 9 — Sensitivity / projection.
	if prob := ValidateSensitivity(SensitivityInput{
		Profile:              profile,
		Classifications:      opts.Classifications,
		DeclaredSensitivity:  opts.DeclaredSensitivity,
		ContentCategories:    opts.ContentCategories,
		ProhibitedCategories: opts.ProhibitedCategories,
		CanonicalView:        opts.CanonicalView,
		AllowedAudiences:     opts.AllowedAudiences,
	}); prob != nil {
		return prob
	}

	return nil
}

// ValidateAuditEvent runs the fixed nine-pass orchestration for an AuditEvent
// base envelope plus FEATURE-0013 AuditLinkage (design §8 / §9.1; F13-AUDIT-001;
// F13-SCOPE-001; AD-009).
//
// Inapplicable decision-domain passes (profile document limits, evaluation,
// composition, trust, obligation vocabulary, relationship matrix, sensitivity
// content) are no-ops. Scope and linkage structural checks always run.
//
// Design §8 names the first argument apiconform.AuditEvent; this package uses
// AuditEventEnvelope instead so validate does not import apiconform (design §5
// DAG; T-027 calls validators one-directionally from apiconform).
func ValidateAuditEvent(ev AuditEventEnvelope, ext decision.AuditLinkage, view bundle.BundleView) *apiproblem.Problem {
	return ValidateAuditEventWith(ev, ext, view, AuditEventOptions{})
}

// ValidateAuditEventWith is ValidateAuditEvent with optional orchestration context.
func ValidateAuditEventWith(ev AuditEventEnvelope, ext decision.AuditLinkage, view bundle.BundleView, opts AuditEventOptions) *apiproblem.Problem {
	_ = view // reserved for later profile-governed audit semantics; no registry lookup required now

	// Pass 1 — typed structural (FEATURE-0012 top-level codes).
	if prob := checkAuditEventStructural(ev, ext); prob != nil {
		return prob
	}

	// Pass 2 — Scope (six governance scopes; Platform permitted).
	subjects := []apimeta.TypedRef{}
	if ev.Record.SubjectRef.Kind != "" || ev.Record.SubjectRef.Name != "" || ev.Record.SubjectRef.UID != "" {
		subjects = append(subjects, ev.Record.SubjectRef)
	}
	if prob := ValidateScope(ScopeInput{
		Metadata:            ev.Metadata,
		Allowed:             allGovernanceScopes(),
		PlatformAllowed:     true,
		ParallelScopeFields: opts.ParallelScopeFields,
		SubjectRefs:         subjects,
	}); prob != nil {
		return prob
	}

	// Passes 3–9: no-op for AuditEvent contract-only validation in FEATURE-0013.
	return nil
}

// DecodeAndValidateAuditEvent runs decode then ValidateAuditEvent (nine-pass
// order with decode as pass 1).
func DecodeAndValidateAuditEvent(
	data []byte,
	mode apivalid.DecodeMode,
	ext decision.AuditLinkage,
	view bundle.BundleView,
) (AuditEventEnvelope, *apiproblem.Problem) {
	return DecodeAndValidateAuditEventWith(data, mode, ext, view, AuditEventOptions{})
}

// DecodeAndValidateAuditEventWith is DecodeAndValidateAuditEvent with options.
func DecodeAndValidateAuditEventWith(
	data []byte,
	mode apivalid.DecodeMode,
	ext decision.AuditLinkage,
	view bundle.BundleView,
	opts AuditEventOptions,
) (AuditEventEnvelope, *apiproblem.Problem) {
	ev, prob := DecodeAuditEventEnvelope(data, mode)
	if prob != nil {
		return AuditEventEnvelope{}, prob
	}
	if prob := ValidateAuditEventWith(ev, ext, view, opts); prob != nil {
		return AuditEventEnvelope{}, prob
	}
	return ev, nil
}

func strictDecodeBytes(data []byte, mode apivalid.DecodeMode, dst any) *apiproblem.Problem {
	contentType := sniffContentType(data)
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(data))
	req.Header.Set("Content-Type", contentType)
	return apivalid.StrictDecode(httptest.NewRecorder(), req, apivalid.DefaultLimits(), mode, dst)
}

func sniffContentType(data []byte) string {
	trimmed := bytes.TrimSpace(data)
	if len(trimmed) > 0 && (trimmed[0] == '{' || trimmed[0] == '[') {
		return "application/json"
	}
	return "application/yaml"
}

func checkDecisionRecordStructural(rec decision.DecisionRecord) *apiproblem.Problem {
	if rec.APIVersion != "" && rec.APIVersion != decision.APIVersionDecisionRecord {
		return structuralProblem(apiproblem.CodeMalformedRequest, "/apiVersion",
			"DecisionRecord apiVersion must be governance.sovrunn.io/v1alpha1")
	}
	if rec.Kind != "" && rec.Kind != decision.KindDecisionRecord {
		return structuralProblem(apiproblem.CodeMalformedRequest, "/kind",
			"DecisionRecord kind must be DecisionRecord")
	}
	if strings.TrimSpace(rec.Metadata.Name) == "" {
		return structuralProblem(apiproblem.CodeMalformedRequest, "/metadata/name",
			"metadata.name is required")
	}
	if strings.TrimSpace(rec.Record.ProfileRef.Name) == "" || strings.TrimSpace(rec.Record.ProfileRef.Version) == "" {
		return structuralProblem(apiproblem.CodeMalformedRequest, "/record/profileRef",
			"record.profileRef name and version are required")
	}
	if !rec.Record.Form.Valid() {
		return structuralProblem(apiproblem.CodeMalformedRequest, "/record/form",
			"record.form must be a closed DecisionForm value")
	}
	if !rec.Record.Authority.Valid() {
		return structuralProblem(apiproblem.CodeMalformedRequest, "/record/authority",
			"record.authority must be a closed DecisionAuthority value")
	}
	if strings.TrimSpace(rec.Record.Purpose) == "" {
		return structuralProblem(apiproblem.CodeMalformedRequest, "/record/purpose",
			"record.purpose is required")
	}
	if strings.TrimSpace(rec.Record.SemanticIdentity.Basis) == "" ||
		len(rec.Record.SemanticIdentity.InputRefs) == 0 ||
		strings.TrimSpace(rec.Record.SemanticIdentity.CanonicalScope) == "" {
		return structuralProblem(apiproblem.CodeMalformedRequest, "/record/semanticIdentity",
			"record.semanticIdentity basis, inputRefs, and canonicalScope are required")
	}
	// F13-AUDIT-001: every meaningful decision links to ≥1 AuditEvent.
	if len(rec.Record.Correlation.AuditEventRefs) == 0 {
		return structuralProblem(apiproblem.CodeMalformedRequest, "/record/correlation/auditEventRefs",
			"correlation.auditEventRefs is required and must be non-empty")
	}
	// Finality is FINAL only on a persisted DecisionRecord (design §7.5).
	if rec.Record.Finality != decision.FinalityFinal {
		return structuralProblem(apiproblem.CodeMalformedRequest, "/record/finality",
			"record.finality must be FINAL")
	}
	return nil
}

func checkAuditEventStructural(ev AuditEventEnvelope, ext decision.AuditLinkage) *apiproblem.Problem {
	if ev.APIVersion != "" && ev.APIVersion != "governance.sovrunn.io/v1alpha1" {
		return structuralProblem(apiproblem.CodeMalformedRequest, "/apiVersion",
			"AuditEvent apiVersion must be governance.sovrunn.io/v1alpha1")
	}
	if ev.Kind != "" && ev.Kind != "AuditEvent" {
		return structuralProblem(apiproblem.CodeMalformedRequest, "/kind",
			"AuditEvent kind must be AuditEvent")
	}
	if strings.TrimSpace(ev.Metadata.Name) == "" {
		return structuralProblem(apiproblem.CodeMalformedRequest, "/metadata/name",
			"metadata.name is required")
	}
	if !ev.Record.Outcome.Valid() {
		return structuralProblem(apiproblem.CodeMalformedRequest, "/record/outcome",
			"record.outcome must be a closed AuditOutcome value")
	}
	// Linkage: decisionRef uid-pinned; correlation.auditEventRefs non-empty.
	if strings.TrimSpace(ext.DecisionRef.APIVersion) == "" ||
		strings.TrimSpace(ext.DecisionRef.Kind) == "" ||
		strings.TrimSpace(ext.DecisionRef.Name) == "" {
		return structuralProblem(apiproblem.CodeMalformedRequest, "/record/decisionLinkage/decisionRef",
			"decisionLinkage.decisionRef requires apiVersion, kind, and name")
	}
	if strings.TrimSpace(ext.DecisionRef.UID) == "" {
		return structuralProblem(apiproblem.CodeMalformedRequest, "/record/decisionLinkage/decisionRef/uid",
			"decisionLinkage.decisionRef.uid is required")
	}
	if len(ext.Correlation.AuditEventRefs) == 0 {
		return structuralProblem(apiproblem.CodeMalformedRequest, "/record/decisionLinkage/correlation/auditEventRefs",
			"decisionLinkage.correlation.auditEventRefs is required and must be non-empty")
	}
	return nil
}

func structuralProblem(code apiproblem.ErrorCode, field, message string) *apiproblem.Problem {
	p := apiproblem.New(code).WithDetail(message)
	// Attach a field pointer when the top-level code supports violations.
	// Decode/structural failures must not use a DECISION_* violations[].code.
	if code == apiproblem.CodeValidationFailed ||
		code == apiproblem.CodeMalformedRequest ||
		code == apiproblem.CodeUnknownField ||
		code == apiproblem.CodeDuplicateField {
		return p.WithViolations([]apiproblem.Violation{{
			Field:   field,
			Code:    apiproblem.ViolationOutOfRange,
			Message: message,
		}})
	}
	return p
}

func governanceScopesForRecord(profile decision.DecisionProfile, found bool) ([]apimeta.ScopeKind, bool) {
	if !found {
		// Profile not yet resolved: allow the seven-value vocabulary so pass 2
		// can still reject invalid kinds; widening against profile is deferred
		// to pass 3 unknown/version failures.
		return allGovernanceScopes(), true
	}
	allowed := append([]apimeta.ScopeKind(nil), profile.Spec.ScopeSemantics.AllowedScopes...)
	if profile.Spec.ScopeSemantics.PlatformAllowed {
		if !containsScopeKind(allowed, apimeta.ScopePlatform) {
			allowed = append(allowed, apimeta.ScopePlatform)
		}
	}
	return allowed, profile.Spec.ScopeSemantics.PlatformAllowed
}

func allGovernanceScopes() []apimeta.ScopeKind {
	return []apimeta.ScopeKind{
		apimeta.ScopePlatform,
		apimeta.ScopeOrganization,
		apimeta.ScopeOrganizationUnit,
		apimeta.ScopeTenant,
		apimeta.ScopeProject,
		apimeta.ScopeCloudPlatform,
		apimeta.ScopeCloudProvider,
	}
}

func decisionTimeHint(rec decision.DecisionRecord) string {
	// Prefer the latest evaluation completedAt as the captured decision time
	// for SecurityExceptionRef expiry checks; empty skips the check.
	for i := len(rec.Record.EvaluationResults) - 1; i >= 0; i-- {
		if t := strings.TrimSpace(rec.Record.EvaluationResults[i].Timing.CompletedAt); t != "" {
			return t
		}
		if t := strings.TrimSpace(rec.Record.EvaluationResults[i].EvaluatedAt); t != "" {
			return t
		}
	}
	return ""
}

func resolveCompositionFromView(view bundle.BundleView, c decision.CompositionRef) (graph.Definition, bool) {
	ref := strings.TrimSpace(c.GraphRef)
	version := strings.TrimSpace(c.GraphVersion)
	if ref == "" || version == "" {
		return graph.Definition{}, false
	}
	return view.Graph(ref, version)
}

func prefixProblemFields(prob *apiproblem.Problem, pointerBase string) *apiproblem.Problem {
	if prob == nil || pointerBase == "" || len(prob.Violations) == 0 {
		return prob
	}
	violations := make([]apiproblem.Violation, len(prob.Violations))
	for i, v := range prob.Violations {
		violations[i] = apiproblem.Violation{
			Field:   joinPointerBase(pointerBase, v.Field),
			Code:    v.Code,
			Message: v.Message,
		}
	}
	return apiproblem.New(prob.Code).
		WithDetail(prob.Detail).
		WithViolations(violations)
}
