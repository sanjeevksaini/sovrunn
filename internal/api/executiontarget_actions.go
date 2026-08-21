package api

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"
	"unicode"

	"github.com/sanjeevksaini/sovrunn/internal/apiconform"
	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/cloudmodel"
	"github.com/sanjeevksaini/sovrunn/internal/cloudmodel/validate"
	"github.com/sanjeevksaini/sovrunn/internal/executiontarget"
	etmodel "github.com/sanjeevksaini/sovrunn/internal/executiontarget/model"
)

// FEATURE-0016 action ServeMux patterns (design §4.5; ADH-2026-058 clause 3).
const (
	PatternExecutionTargetQualify = "POST /apis/execution.sovrunn.io/v1alpha1/execution-targets/{uid}/actions/qualify"
	PatternExecutionTargetRetire  = "POST /apis/execution.sovrunn.io/v1alpha1/execution-targets/{uid}/actions/retire"
)

// ExecutionTargetActionKind selects qualify vs retire pipeline ownership.
type ExecutionTargetActionKind string

const (
	executionTargetActionQualify ExecutionTargetActionKind = "qualify"
	executionTargetActionRetire  ExecutionTargetActionKind = "retire"
)

// ExecutionTargetActionHandler serves one ExecutionTarget item-action route
// under /actions/qualify or /actions/retire (design §5.1; TASK-F16-11).
type ExecutionTargetActionHandler struct {
	Kind      ExecutionTargetActionKind
	Lifecycle *executiontarget.ExecutionTargetLifecycleService
	Cloud     *cloudmodel.Store
	Grants    GrantResolver
	Audit     executiontarget.AuditAppender
	// Fixtures is the non-blocking MapFixtureSource shared with the synthetic
	// observer. Qualify predicts the idempotency completion body from this
	// source without calling Observe, so Retire/Maintenance can win while a
	// test blocking observer holds the real observation open.
	Fixtures *executiontarget.MapFixtureSource
}

// NewExecutionTargetQualifyHandler constructs the qualify action handler.
func NewExecutionTargetQualifyHandler(
	lifecycle *executiontarget.ExecutionTargetLifecycleService,
	cloud *cloudmodel.Store,
	grants GrantResolver,
	audit executiontarget.AuditAppender,
) *ExecutionTargetActionHandler {
	return &ExecutionTargetActionHandler{
		Kind:      executionTargetActionQualify,
		Lifecycle: lifecycle,
		Cloud:     cloud,
		Grants:    grants,
		Audit:     audit,
	}
}

// NewExecutionTargetRetireHandler constructs the retire action handler.
func NewExecutionTargetRetireHandler(
	lifecycle *executiontarget.ExecutionTargetLifecycleService,
	cloud *cloudmodel.Store,
	grants GrantResolver,
	audit executiontarget.AuditAppender,
) *ExecutionTargetActionHandler {
	return &ExecutionTargetActionHandler{
		Kind:      executionTargetActionRetire,
		Lifecycle: lifecycle,
		Cloud:     cloud,
		Grants:    grants,
		Audit:     audit,
	}
}

// ServeHTTP dispatches POST only for the configured action.
func (h *ExecutionTargetActionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeMethodNotAllowed(w, r, "POST")
		return
	}
	uid := r.PathValue("uid")
	for {
		retry, finished := h.actionAttempt(w, r, uid)
		if finished || !retry {
			return
		}
	}
}

// actionAttempt runs one qualify/retire pipeline attempt. retry=true means the
// waiter woke and must re-enter at authentication.
func (h *ExecutionTargetActionHandler) actionAttempt(
	w http.ResponseWriter,
	r *http.Request,
	uid string,
) (retry, finished bool) {
	principal, prob := authenticatePrincipal(r, h.Grants)
	if prob != nil {
		writeProblem(w, r, prob)
		return false, true
	}
	ifMatch, grammarProb := validateExecutionTargetActionGrammar(r)
	if grammarProb != nil {
		writeProblem(w, r, grammarProb)
		return false, true
	}
	if h.Lifecycle == nil || h.Cloud == nil {
		writeProblem(w, r, apiproblem.New(apiproblem.CodeInternalError).WithDetail("executiontarget dependencies are not configured"))
		return false, true
	}
	if !apimeta.IsGeneratedUIDFormat(uid) {
		writeProblem(w, r, apiproblem.New(apiproblem.CodeMalformedRequest).WithDetail("path uid segment is malformed"))
		return false, true
	}

	et, ok := h.Lifecycle.Store().GetExecutionTarget(uid)
	if !ok {
		h.writeAuditedSafeDenial(w, r, principal)
		return false, true
	}

	scope := apimeta.ScopeIdentity{Kind: apimeta.ScopeCloudProvider, UID: et.CloudProviderScopeUID}
	grantAction := ActionExecutionTargetWrite
	grantDetail := "executiontarget.write grant is required"
	if h.Kind == executionTargetActionQualify {
		grantAction = ActionExecutionTargetQualify
		grantDetail = "executiontarget.qualify grant is required"
	}
	if !h.Grants.AuthorizeExact(grantAction, scope, "", "") {
		h.writeAuditedAuthorizationDenial(w, r, principal,
			apiproblem.New(apiproblem.CodeAuthorizationDenied).WithDetail(grantDetail),
		)
		return false, true
	}

	backing := h.resolveBackingViability(et)
	if h.Kind == executionTargetActionQualify {
		if admitProb := h.admitQualifyBacking(w, r, principal, et, &backing); admitProb {
			return false, true
		}
	}

	digest := executiontarget.DigestAction()
	ns := executiontarget.IdempotencyNamespace{
		PrincipalUID:     principal,
		RoutePattern:     h.routePattern(),
		CloudProviderUID: et.CloudProviderScopeUID,
		ActionTargetUID:  uid,
		Key:              r.Header.Get("Idempotency-Key"),
	}

	inspect := h.Lifecycle.ReserveOrInspectReplay(ns, digest)
	switch inspect.Outcome {
	case executiontarget.ReservationOutcomeReplay:
		writeExecutionTargetReplay(w, r, inspect.Replay)
		return false, true
	case executiontarget.ReservationOutcomeConflict:
		writeProblem(w, r, idempotencyKeyReuseMismatchProblem())
		return false, true
	case executiontarget.ReservationOutcomeInFlight:
		done, ok := h.Lifecycle.AttachWaiter(ns)
		if !ok {
			return true, false
		}
		if err := h.Lifecycle.AwaitReservation(r.Context(), done); err != nil {
			h.Lifecycle.DetachWaiter(ns)
			writeProblem(w, r, apiproblem.New(apiproblem.CodeInternalError).WithDetail("action waiter cancelled"))
			return false, true
		}
		return true, false
	case executiontarget.ReservationOutcomeReserved:
		// proceed
	default:
		writeProblem(w, r, apiproblem.New(apiproblem.CodeInternalError).WithDetail("unexpected reservation outcome"))
		return false, true
	}

	holdingReservation := true
	defer func() {
		if rec := recover(); rec != nil {
			if holdingReservation {
				h.Lifecycle.AbortReservation(ns)
				holdingReservation = false
			}
			panic(rec)
		}
	}()

	// Re-load target after reservation for current If-Match and lifecycle checks.
	live, ok := h.Lifecycle.Store().GetExecutionTarget(uid)
	if !ok {
		h.Lifecycle.AbortReservation(ns)
		holdingReservation = false
		h.writeAuditedSafeDenial(w, r, principal)
		return false, true
	}
	backing = h.resolveBackingViability(live)

	// Current If-Match precedes lifecycle evaluation (F16-74, F16-85).
	if prob := validate.CheckIfMatchCurrent(ifMatch, live.Metadata.ResourceVersion); prob != nil {
		h.Lifecycle.AbortReservation(ns)
		holdingReservation = false
		writeProblem(w, r, prob)
		return false, true
	}

	snap, snapOK := h.Lifecycle.CaptureProjectionSnapshot(uid, backing)
	if !snapOK {
		h.Lifecycle.AbortReservation(ns)
		holdingReservation = false
		h.writeAuditedSafeDenial(w, r, principal)
		return false, true
	}

	if h.Kind == executionTargetActionQualify {
		if live.Status.Lifecycle != etmodel.LifecycleActive {
			h.Lifecycle.AbortReservation(ns)
			holdingReservation = false
			writeProblem(w, r, targetRetiredProblem())
			return false, true
		}
		if snap.MaintenanceActive {
			h.Lifecycle.AbortReservation(ns)
			holdingReservation = false
			writeProblem(w, r, targetMaintenanceProblem())
			return false, true
		}
		return h.dispatchQualify(w, r, principal, uid, live, backing, ns, digest, &holdingReservation)
	}
	return h.dispatchRetire(w, r, principal, uid, live, backing, ns, digest, &holdingReservation)
}

func (h *ExecutionTargetActionHandler) dispatchQualify(
	w http.ResponseWriter,
	r *http.Request,
	principal, uid string,
	live etmodel.ExecutionTarget,
	backing executiontarget.BackingViability,
	ns executiontarget.IdempotencyNamespace,
	digest executiontarget.Digest,
	holding *bool,
) (retry, finished bool) {
	// Predict the completion body from the shared non-blocking fixture source.
	// Never call Lifecycle.Observer().Observe here: that path may block under
	// test and would race Retire/Maintenance wins before Qualifying opens.
	predictedBody, err := predictQualifyCompletionBody(uid, live, backing, h.Fixtures)
	if err != nil {
		h.Lifecycle.AbortReservation(ns)
		*holding = false
		writeProblem(w, r, apiproblem.New(apiproblem.CodeInternalError).WithDetail("failed to encode qualify response"))
		return false, true
	}

	factUID, err := apimeta.GenerateUID()
	if err != nil {
		h.Lifecycle.AbortReservation(ns)
		*holding = false
		writeProblem(w, r, apiproblem.New(apiproblem.CodeInternalError).WithDetail("failed to allocate fact set uid"))
		return false, true
	}
	resultUID, err := apimeta.GenerateUID()
	if err != nil {
		h.Lifecycle.AbortReservation(ns)
		*holding = false
		writeProblem(w, r, apiproblem.New(apiproblem.CodeInternalError).WithDetail("failed to allocate result uid"))
		return false, true
	}

	cloud := h.Cloud
	spec := live.Spec
	req := executiontarget.QualifyCommitRequest{
		TargetUID: uid,
		Backing:   backing,
		RefreshBacking: func() executiontarget.BackingViability {
			// Invoked under the DD-02 mutex + store publication lock. Must not
			// call F0016 Store getters (would re-enter the store mutex).
			return resolveBackingViabilityFromCloud(cloud, etmodel.ExecutionTarget{Spec: spec})
		},
		FactSetUID:  factUID,
		ResultUID:   resultUID,
		Actor:       actorRef(principal),
		RequestID:   requestID(r),
		AuditUID:    newAuditUID("et-qualify"),
		Idempotency: ns,
		Digest:      digest,
		Completion: executiontarget.CompletedResult{
			StatusCode: http.StatusOK,
			Body:       predictedBody,
		},
	}

	got, commitProb := h.Lifecycle.Qualify(r.Context(), req)
	*holding = false
	if commitProb != nil {
		// Lifecycle exclusively owns post-dispatch aborts; late AbortReservation
		// is a bounded no-op when already finalized.
		h.Lifecycle.AbortReservation(ns)
		writeProblem(w, r, commitProb)
		return false, true
	}

	// Return the same body stored in the idempotency completion so completed
	// replay matches the original successful response exactly (F16-106).
	w.Header().Set("ETag", got.Metadata.ResourceVersion)
	writeRawJSON(w, r, http.StatusOK, mediaJSON, predictedBody)
	return false, true
}

func (h *ExecutionTargetActionHandler) dispatchRetire(
	w http.ResponseWriter,
	r *http.Request,
	principal, uid string,
	live etmodel.ExecutionTarget,
	backing executiontarget.BackingViability,
	ns executiontarget.IdempotencyNamespace,
	digest executiontarget.Digest,
	holding *bool,
) (retry, finished bool) {
	if live.Status.Lifecycle != etmodel.LifecycleActive {
		h.Lifecycle.AbortReservation(ns)
		*holding = false
		writeProblem(w, r, targetRetiredProblem())
		return false, true
	}

	predicted := live
	predicted.Status.Lifecycle = etmodel.LifecycleRetired
	predicted.Status.Qualification = etmodel.QualificationUnqualified
	predicted.Status.FactSetRef = nil
	predicted.Status.QualificationResultRef = nil
	predicted.Metadata.ResourceVersion = nextOpaqueResourceVersion(live.Metadata.ResourceVersion)
	respBody, err := json.Marshal(executiontarget.Project(executiontarget.ProjectionSnapshot{
		Target:            predicted,
		MaintenanceActive: false,
		FactsFresh:        false,
		BackingViable:     backing.ParticipationEffectiveActive && backing.StackActive(),
	}))
	if err != nil {
		h.Lifecycle.AbortReservation(ns)
		*holding = false
		writeProblem(w, r, apiproblem.New(apiproblem.CodeInternalError).WithDetail("failed to encode retire response"))
		return false, true
	}

	req := executiontarget.RetireCommitRequest{
		TargetUID:               uid,
		ExpectedResourceVersion: live.Metadata.ResourceVersion,
		Actor:                   actorRef(principal),
		RequestID:               requestID(r),
		AuditUID:                newAuditUID("et-retire"),
		Idempotency:             ns,
		Digest:                  digest,
		Completion: executiontarget.CompletedResult{
			StatusCode: http.StatusOK,
			Body:       respBody,
		},
	}

	got, commitProb := h.Lifecycle.CommitRetire(r.Context(), req)
	*holding = false
	if commitProb != nil {
		h.Lifecycle.AbortReservation(ns)
		writeProblem(w, r, commitProb)
		return false, true
	}

	w.Header().Set("ETag", got.Metadata.ResourceVersion)
	writeRawJSON(w, r, http.StatusOK, mediaJSON, respBody)
	return false, true
}

func (h *ExecutionTargetActionHandler) routePattern() string {
	if h.Kind == executionTargetActionQualify {
		return PatternExecutionTargetQualify
	}
	return PatternExecutionTargetRetire
}

// admitQualifyBacking performs REQ-F16-02 safe backing/current-viability
// admission before replay/reservation. Returns true when a response was written.
func (h *ExecutionTargetActionHandler) admitQualifyBacking(
	w http.ResponseWriter,
	r *http.Request,
	principal string,
	et etmodel.ExecutionTarget,
	backing *executiontarget.BackingViability,
) bool {
	partUID := et.Spec.CloudProviderParticipationRef.UID
	stackUID := et.Spec.InfrastructureStackRef.UID
	part, partOK := h.Cloud.GetParticipation(partUID)
	stack, stackOK := h.Cloud.GetInfrastructureStack(stackUID)
	if !partOK || !stackOK {
		h.writeAuditedSafeDenial(w, r, principal)
		return true
	}

	*backing = executiontarget.BackingViability{
		ParticipationUID:             part.Metadata.UID,
		ParticipationEffectiveActive: participationEffectiveActive(part),
		ParticipationScopeUID:        part.Spec.CloudProviderRef.UID,
		StackUID:                     stack.Metadata.UID,
		StackPhase:                   string(stack.Status.Phase),
		StackScopeUID:                apimeta.CanonicalScopeIdentity(stack.Metadata.ScopeRef).UID,
		StackGeneration:              stack.Metadata.Generation,
	}

	if !backing.ParticipationEffectiveActive {
		writeProblem(w, r, apiproblem.New(apiproblem.CodeConflict).WithViolations([]apiproblem.Violation{{
			Field:   "/spec/cloudProviderParticipationRef/uid",
			Code:    executiontarget.ViolationExecutionTargetParticipationUnavailable,
			Message: "CloudProviderParticipation is not effective-Active",
		}}))
		return true
	}
	if !backing.StackActive() {
		writeProblem(w, r, apiproblem.New(apiproblem.CodeConflict).WithViolations([]apiproblem.Violation{{
			Field:   "/spec/infrastructureStackRef/uid",
			Code:    executiontarget.ViolationExecutionTargetStackUnavailable,
			Message: "InfrastructureStack is not Active",
		}}))
		return true
	}
	if backing.ParticipationScopeUID == "" || backing.StackScopeUID == "" ||
		backing.ParticipationScopeUID != backing.StackScopeUID {
		writeProblem(w, r, apiproblem.New(apiproblem.CodeValidationFailed).WithViolations([]apiproblem.Violation{{
			Field:   "/spec",
			Code:    executiontarget.ViolationExecutionTargetScopeMismatch,
			Message: "participation and stack must share the same CloudProvider scope",
		}}))
		return true
	}
	return false
}

func (h *ExecutionTargetActionHandler) resolveBackingViability(et etmodel.ExecutionTarget) executiontarget.BackingViability {
	return resolveBackingViabilityFromCloud(h.Cloud, et)
}

func resolveBackingViabilityFromCloud(cloud *cloudmodel.Store, et etmodel.ExecutionTarget) executiontarget.BackingViability {
	out := executiontarget.BackingViability{
		ParticipationUID: et.Spec.CloudProviderParticipationRef.UID,
		StackUID:         et.Spec.InfrastructureStackRef.UID,
	}
	if cloud == nil {
		return out
	}
	if part, ok := cloud.GetParticipation(et.Spec.CloudProviderParticipationRef.UID); ok {
		out.ParticipationEffectiveActive = participationEffectiveActive(part)
		out.ParticipationScopeUID = part.Spec.CloudProviderRef.UID
	}
	if stack, ok := cloud.GetInfrastructureStack(et.Spec.InfrastructureStackRef.UID); ok {
		out.StackPhase = string(stack.Status.Phase)
		out.StackScopeUID = apimeta.CanonicalScopeIdentity(stack.Metadata.ScopeRef).UID
		out.StackGeneration = stack.Metadata.Generation
	}
	return out
}

func validateExecutionTargetActionGrammar(r *http.Request) (ifMatch string, prob *apiproblem.Problem) {
	if prob := rejectExecutionTargetQuery(r); prob != nil {
		return "", prob
	}
	idemValues := r.Header.Values("Idempotency-Key")
	if len(idemValues) != 1 {
		if len(idemValues) == 0 {
			return "", apiproblem.New(apiproblem.CodeMalformedRequest).WithDetail("Idempotency-Key is required")
		}
		return "", apiproblem.New(apiproblem.CodeMalformedRequest).WithDetail("Idempotency-Key must appear exactly once")
	}
	if prob := validate.ValidateIdempotencyKey(idemValues[0]); prob != nil {
		return "", prob
	}
	ifMatch, prob = validateExecutionTargetStrongIfMatch(r)
	if prob != nil {
		return "", prob
	}
	if prob := requireApplicationJSON(r.Header.Get("Content-Type")); prob != nil {
		return "", prob
	}
	body, bodyProb := readExecutionTargetActionBody(r)
	if bodyProb != nil {
		return "", bodyProb
	}
	if len(body) != 0 {
		return "", apiproblem.New(apiproblem.CodeMalformedRequest).WithDetail("qualify/retire action body must be zero bytes")
	}
	return ifMatch, nil
}

func validateExecutionTargetStrongIfMatch(r *http.Request) (string, *apiproblem.Problem) {
	values := r.Header.Values("If-Match")
	if len(values) == 0 {
		return "", apiproblem.New(apiproblem.CodeStaleResourceVersion).WithDetail("If-Match is required")
	}
	if len(values) > 1 {
		return "", apiproblem.New(apiproblem.CodeStaleResourceVersion).WithDetail("If-Match must be exactly one strong opaque token")
	}
	raw := values[0]
	if strings.Contains(raw, ",") {
		return "", apiproblem.New(apiproblem.CodeStaleResourceVersion).WithDetail("If-Match must be exactly one strong opaque token")
	}
	ifMatch := strings.TrimSpace(raw)
	if ifMatch == "" {
		return "", apiproblem.New(apiproblem.CodeStaleResourceVersion).WithDetail("If-Match is required")
	}
	if ifMatch == "*" {
		return "", apiproblem.New(apiproblem.CodeStaleResourceVersion).WithDetail("If-Match wildcard is not accepted")
	}
	upper := strings.ToUpper(ifMatch)
	if strings.HasPrefix(upper, "W/") {
		return "", apiproblem.New(apiproblem.CodeStaleResourceVersion).WithDetail("If-Match weak validators are not accepted")
	}
	if strings.HasPrefix(ifMatch, "\"") || strings.HasSuffix(ifMatch, "\"") {
		return "", apiproblem.New(apiproblem.CodeStaleResourceVersion).WithDetail("If-Match must be a strong opaque resourceVersion")
	}
	for _, rr := range ifMatch {
		if rr > unicode.MaxASCII || !unicode.IsPrint(rr) || unicode.IsSpace(rr) {
			return "", apiproblem.New(apiproblem.CodeStaleResourceVersion).WithDetail("If-Match is malformed")
		}
	}
	return ifMatch, nil
}

func readExecutionTargetActionBody(r *http.Request) ([]byte, *apiproblem.Problem) {
	body := r.Body
	if body == nil {
		body = http.NoBody
	}
	// Bound reads; actions must be zero-byte so any non-empty payload fails.
	data, err := io.ReadAll(io.LimitReader(body, 64<<10))
	if err != nil {
		var maxErr *http.MaxBytesError
		if errors.As(err, &maxErr) {
			return nil, apiproblem.New(apiproblem.CodeRequestTooLarge).WithDetail("request body exceeds limit")
		}
		return nil, apiproblem.New(apiproblem.CodeMalformedRequest).WithDetail("failed to read request body")
	}
	return data, nil
}

func targetRetiredProblem() *apiproblem.Problem {
	return apiproblem.New(apiproblem.CodeConflict).WithViolations([]apiproblem.Violation{{
		Field:   "/status/lifecycle",
		Code:    executiontarget.ViolationTargetRetired,
		Message: "ExecutionTarget is not Active",
	}})
}

func targetMaintenanceProblem() *apiproblem.Problem {
	return apiproblem.New(apiproblem.CodeConflict).WithViolations([]apiproblem.Violation{{
		Field:   "/status/lifecycle",
		Code:    executiontarget.ViolationTargetMaintenance,
		Message: "ExecutionTarget is in Maintenance",
	}})
}

func qualificationStateFromOutcome(outcome etmodel.QualificationOutcome) etmodel.QualificationState {
	switch outcome {
	case etmodel.OutcomeQualified:
		return etmodel.QualificationQualified
	case etmodel.OutcomeRejected:
		return etmodel.QualificationRejected
	case etmodel.OutcomeIndeterminate:
		return etmodel.QualificationIndeterminate
	default:
		return etmodel.QualificationUnqualified
	}
}

// predictQualifyCompletionBody builds the idempotency completion JSON from a
// non-blocking fixture read that mirrors the synthetic observer's missing-fixture
// and present-fixture outcomes. Observer faults still dispatch into Qualify,
// which exclusively owns abort after commit dispatch begins.
func predictQualifyCompletionBody(
	uid string,
	live etmodel.ExecutionTarget,
	backing executiontarget.BackingViability,
	fixtures *executiontarget.MapFixtureSource,
) ([]byte, error) {
	obs := executiontarget.NewSyntheticObserver(fixtures, nil).Observe(uid)
	eval, evalErr := executiontarget.EvaluateObservation(obs, time.Now().UTC())

	predicted := live
	if evalErr == nil {
		predicted.Status.Qualification = qualificationStateFromOutcome(eval.Outcome)
		predicted.Status.ObservedGeneration = backing.StackGeneration
	}
	predicted.Metadata.ResourceVersion = nextOpaqueResourceVersion(live.Metadata.ResourceVersion)
	return json.Marshal(executiontarget.Project(executiontarget.ProjectionSnapshot{
		Target:            predicted,
		MaintenanceActive: false,
		FactsFresh:        evalErr == nil,
		BackingViable:     backing.ParticipationEffectiveActive && backing.StackActive(),
	}))
}

func (h *ExecutionTargetActionHandler) writeAuditedAuthorizationDenial(
	w http.ResponseWriter,
	r *http.Request,
	principal string,
	suppressed *apiproblem.Problem,
) {
	event := executiontarget.AssembleAuthorizationDenial(executiontarget.RedactedAuditInput{
		UID:       newAuditUID("et-action-authz-denied"),
		RequestID: requestID(r),
		Actor:     actorRef(principal),
	})
	h.appendDenialOrInternal(w, r, event, suppressed)
}

func (h *ExecutionTargetActionHandler) writeAuditedSafeDenial(
	w http.ResponseWriter,
	r *http.Request,
	principal string,
) {
	suppressed := apiproblem.New(apiproblem.CodeResourceNotFound).WithViolations([]apiproblem.Violation{{
		Code:    violationAuthorizationSafeDenial,
		Message: "resource not found",
	}})
	event := executiontarget.AssembleSafeDenial(executiontarget.RedactedAuditInput{
		UID:       newAuditUID("et-action-safe-denial"),
		RequestID: requestID(r),
		Actor:     actorRef(principal),
	})
	h.appendDenialOrInternal(w, r, event, suppressed)
}

func (h *ExecutionTargetActionHandler) appendDenialOrInternal(
	w http.ResponseWriter,
	r *http.Request,
	event apiconform.AuditEvent,
	suppressed *apiproblem.Problem,
) {
	if h.Audit == nil {
		writeProblem(w, r, apiproblem.New(apiproblem.CodeInternalError).WithDetail("audit appender is not configured"))
		return
	}
	if err := h.Audit.Append(r.Context(), event); err != nil {
		writeProblem(w, r, apiproblem.New(apiproblem.CodeInternalError).WithDetail("required AuditEvent append failed"))
		return
	}
	writeProblem(w, r, suppressed)
}
