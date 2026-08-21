package api

import (
	"encoding/json"
	"mime"
	"net/http"
	"sort"
	"strings"

	"github.com/sanjeevksaini/sovrunn/internal/apiconform"
	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/apivalid"
	"github.com/sanjeevksaini/sovrunn/internal/cloudmodel"
	cmmodel "github.com/sanjeevksaini/sovrunn/internal/cloudmodel/model"
	"github.com/sanjeevksaini/sovrunn/internal/cloudmodel/validate"
	"github.com/sanjeevksaini/sovrunn/internal/executiontarget"
	etmodel "github.com/sanjeevksaini/sovrunn/internal/executiontarget/model"
)

// FEATURE-0016 collection create idempotency route pattern (design §4.5; DD-05).
const PatternExecutionTargetCreate = "POST /apis/execution.sovrunn.io/v1alpha1/execution-targets"

// FEATURE-0016 server-resolved grant actions (ADH-2026-058 clause 3).
const (
	ActionExecutionTargetWrite   = "executiontarget.write"
	ActionExecutionTargetRead    = "executiontarget.read"
	ActionExecutionTargetQualify = "executiontarget.qualify"
)

// ExecutionTargetCollectionHandler serves ExecutionTarget collection POST (create)
// and GET (LIST) per design §5.1.1.
type ExecutionTargetCollectionHandler struct {
	Lifecycle *executiontarget.ExecutionTargetLifecycleService
	Cloud     *cloudmodel.Store
	Grants    GrantResolver
	Audit     executiontarget.AuditAppender
}

// NewExecutionTargetCollectionHandler constructs the collection handler.
func NewExecutionTargetCollectionHandler(
	lifecycle *executiontarget.ExecutionTargetLifecycleService,
	cloud *cloudmodel.Store,
	grants GrantResolver,
	audit executiontarget.AuditAppender,
) *ExecutionTargetCollectionHandler {
	return &ExecutionTargetCollectionHandler{
		Lifecycle: lifecycle,
		Cloud:     cloud,
		Grants:    grants,
		Audit:     audit,
	}
}

// ServeHTTP dispatches GET (LIST) and POST (create).
func (h *ExecutionTargetCollectionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.list(w, r)
	case http.MethodPost:
		h.create(w, r)
	default:
		writeMethodNotAllowed(w, r, "GET, POST")
	}
}

func (h *ExecutionTargetCollectionHandler) list(w http.ResponseWriter, r *http.Request) {
	// LIST ignores Idempotency-Key and never touches idempotency state.
	principal, prob := authenticatePrincipal(r, h.Grants)
	if prob != nil {
		writeProblem(w, r, prob)
		return
	}
	if prob := rejectExecutionTargetQuery(r); prob != nil {
		writeProblem(w, r, prob)
		return
	}
	if len(h.Grants.CoarseLookup(ActionExecutionTargetRead)) == 0 {
		h.writeAuditedAuthorizationDenial(w, r, principal,
			apiproblem.New(apiproblem.CodeAuthorizationDenied).WithDetail("executiontarget.read grant is required"),
		)
		return
	}
	if h.Lifecycle == nil {
		writeProblem(w, r, apiproblem.New(apiproblem.CodeInternalError).WithDetail("lifecycle service is not configured"))
		return
	}

	targets := h.Lifecycle.Store().ListExecutionTargets()
	items := make([]executiontarget.SafeExecutionTarget, 0, len(targets))
	for _, et := range targets {
		scope := apimeta.ScopeIdentity{Kind: apimeta.ScopeCloudProvider, UID: et.CloudProviderScopeUID}
		if !h.Grants.AuthorizeExact(ActionExecutionTargetRead, scope, "", "") {
			continue
		}
		backing := h.resolveBackingViability(et)
		snap, ok := h.Lifecycle.CaptureProjectionSnapshot(et.Metadata.UID, backing)
		if !ok {
			continue
		}
		items = append(items, executiontarget.Project(snap))
	}
	writeJSONSuccess(w, r, http.StatusOK, items)
}

func (h *ExecutionTargetCollectionHandler) create(w http.ResponseWriter, r *http.Request) {
	var (
		cachedBody []byte
		bodyCached bool
	)
	for {
		retry, finished := h.createAttempt(w, r, &cachedBody, &bodyCached)
		if finished || !retry {
			return
		}
		// Waiter woke after owner abort/complete: re-enter at authentication
		// with the cached create body (request Body is already consumed).
	}
}

// createAttempt runs one create pipeline attempt. retry=true means the caller
// must re-enter from authentication after a waiter wake. finished=true means
// a response was written.
func (h *ExecutionTargetCollectionHandler) createAttempt(
	w http.ResponseWriter,
	r *http.Request,
	cachedBody *[]byte,
	bodyCached *bool,
) (retry, finished bool) {
	principal, prob := authenticatePrincipal(r, h.Grants)
	if prob != nil {
		writeProblem(w, r, prob)
		return false, true
	}
	if prob := validateExecutionTargetCreateGrammar(r); prob != nil {
		writeProblem(w, r, prob)
		return false, true
	}
	if h.Lifecycle == nil || h.Cloud == nil {
		writeProblem(w, r, apiproblem.New(apiproblem.CodeInternalError).WithDetail("executiontarget dependencies are not configured"))
		return false, true
	}

	var phaseOne PhaseOneResult
	if bodyCached != nil && *bodyCached && cachedBody != nil {
		phaseOne, prob = DecodePhaseOne(*cachedBody, apivalid.DefaultLimits())
	} else {
		phaseOne, prob = DecodePhaseOneHTTP(w, r, apivalid.DefaultLimits())
		if prob == nil && cachedBody != nil && bodyCached != nil {
			*cachedBody = append([]byte(nil), phaseOne.Body...)
			*bodyCached = true
		}
	}
	if prob != nil {
		writeProblem(w, r, prob)
		return false, true
	}

	participationUID, stackUID, extractProb := extractExecutionTargetCreateUIDs(phaseOne)
	if extractProb != nil {
		writeProblem(w, r, extractProb)
		return false, true
	}

	var (
		part             cmmodel.CloudProviderParticipation
		stack            cmmodel.InfrastructureStack
		cloudProviderUID string
		providerScope    apimeta.ScopeIdentity
		backing          executiontarget.BackingViability
	)

	if participationUID == "" || stackUID == "" {
		// Missing refs cannot derive scope or perform safe access. Require a
		// coarse write grant, then let classification/graph return 422.
		if len(h.Grants.CoarseLookup(ActionExecutionTargetWrite)) == 0 {
			h.writeAuditedAuthorizationDenial(w, r, principal,
				apiproblem.New(apiproblem.CodeAuthorizationDenied).WithDetail("executiontarget.write grant is required"),
			)
			return false, true
		}
	} else {
		var partOK bool
		part, partOK = h.Cloud.GetParticipation(participationUID)
		if !partOK {
			h.writeAuditedSafeDenial(w, r, principal)
			return false, true
		}
		cloudProviderUID = part.Spec.CloudProviderRef.UID
		if cloudProviderUID == "" {
			h.writeAuditedSafeDenial(w, r, principal)
			return false, true
		}
		providerScope = apimeta.ScopeIdentity{Kind: apimeta.ScopeCloudProvider, UID: cloudProviderUID}
		if !h.Grants.AuthorizeExact(ActionExecutionTargetWrite, providerScope, "", "") {
			h.writeAuditedAuthorizationDenial(w, r, principal,
				apiproblem.New(apiproblem.CodeAuthorizationDenied).WithDetail("executiontarget.write grant is required"),
			)
			return false, true
		}

		var stackOK bool
		stack, stackOK = h.Cloud.GetInfrastructureStack(stackUID)
		if !stackOK {
			h.writeAuditedSafeDenial(w, r, principal)
			return false, true
		}
	}

	classProb, clientStatusAttempted := classifyExecutionTargetCreate(phaseOne.Body)
	if classProb != nil {
		h.writeClassifiedCreateDenial(w, r, principal, classProb)
		return false, true
	}

	name := strings.TrimSpace(phaseOne.Name)
	targetClass, graphProb := validateExecutionTargetCreateGraph(name, part, stack, phaseOne.Body)
	if graphProb != nil {
		writeProblem(w, r, graphProb)
		return false, true
	}

	if cloudProviderUID == "" {
		cloudProviderUID = part.Spec.CloudProviderRef.UID
	}
	backing = executiontarget.BackingViability{
		ParticipationUID:             part.Metadata.UID,
		ParticipationEffectiveActive: participationEffectiveActive(part),
		ParticipationScopeUID:        cloudProviderUID,
		StackUID:                     stack.Metadata.UID,
		StackPhase:                   string(stack.Status.Phase),
		StackScopeUID:                apimeta.CanonicalScopeIdentity(stack.Metadata.ScopeRef).UID,
		StackGeneration:              stack.Metadata.Generation,
	}

	digest, err := executiontarget.DigestCreate(json.RawMessage(phaseOne.Body))
	if err != nil {
		writeProblem(w, r, apiproblem.New(apiproblem.CodeInternalError).WithDetail("failed to compute idempotency digest"))
		return false, true
	}
	ns := executiontarget.IdempotencyNamespace{
		PrincipalUID:     principal,
		RoutePattern:     PatternExecutionTargetCreate,
		CloudProviderUID: cloudProviderUID,
		ActionTargetUID:  "",
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
			writeProblem(w, r, apiproblem.New(apiproblem.CodeInternalError).WithDetail("create waiter cancelled"))
			return false, true
		}
		return true, false
	case executiontarget.ReservationOutcomeReserved:
		// proceed to commit
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

	uid, err := apimeta.GenerateUID()
	if err != nil {
		h.Lifecycle.AbortReservation(ns)
		holdingReservation = false
		writeProblem(w, r, apiproblem.New(apiproblem.CodeInternalError).WithDetail("failed to allocate uid"))
		return false, true
	}

	partRef := apimeta.TypedRef{
		APIVersion: cmmodel.APIVersionCloudProviderParticipation,
		Kind:       cmmodel.KindCloudProviderParticipation,
		Name:       part.Metadata.Name,
		UID:        part.Metadata.UID,
	}
	stackRef := apimeta.TypedRef{
		APIVersion: cmmodel.APIVersionInfrastructureStack,
		Kind:       cmmodel.KindInfrastructureStack,
		Name:       stack.Metadata.Name,
		UID:        stack.Metadata.UID,
	}

	proposed := etmodel.NewCreateProposal(
		name, uid, cloudProviderUID, partRef, stackRef, stack.Metadata.Generation,
	)
	if targetClass != "" {
		proposed.Spec.TargetClass = targetClass
	}
	proj := executiontarget.Project(executiontarget.ProjectionSnapshot{
		Target:            proposed,
		MaintenanceActive: false,
		FactsFresh:        false,
		BackingViable:     backing.ParticipationEffectiveActive && backing.StackActive(),
	})
	respBody, err := json.Marshal(proj)
	if err != nil {
		h.Lifecycle.AbortReservation(ns)
		holdingReservation = false
		writeProblem(w, r, apiproblem.New(apiproblem.CodeInternalError).WithDetail("failed to encode create response"))
		return false, true
	}

	commitReq := executiontarget.CreateCommitRequest{
		Name:                  name,
		UID:                   uid,
		CloudProviderScopeUID: cloudProviderUID,
		ParticipationRef:      partRef,
		StackRef:              stackRef,
		StackGeneration:       stack.Metadata.Generation,
		TargetClass:           targetClass,
		ClientStatusAttempted: clientStatusAttempted,
		Actor:                 actorRef(principal),
		RequestID:             requestID(r),
		AuditUID:              newAuditUID("et-create"),
		Idempotency:           ns,
		Digest:                digest,
		Completion: executiontarget.CompletedResult{
			StatusCode: http.StatusCreated,
			Body:       respBody,
		},
	}

	created, commitProb := h.Lifecycle.CommitCreate(r.Context(), commitReq)
	holdingReservation = false
	if commitProb != nil {
		writeProblem(w, r, commitProb)
		return false, true
	}

	// Return the same body stored in the idempotency completion so completed
	// replay matches the original successful response exactly (F16-31).
	writeExecutionTargetCreated(w, r, created.Metadata.ResourceVersion, respBody)
	return false, true
}

func validateExecutionTargetCreateGrammar(r *http.Request) *apiproblem.Problem {
	if prob := rejectExecutionTargetQuery(r); prob != nil {
		return prob
	}
	if strings.TrimSpace(r.Header.Get("If-Match")) != "" {
		return apiproblem.New(apiproblem.CodeMalformedRequest).WithDetail("If-Match is not accepted on executiontarget create")
	}
	if prob := requireApplicationJSON(r.Header.Get("Content-Type")); prob != nil {
		return prob
	}
	return validate.ValidateIdempotencyKey(r.Header.Get("Idempotency-Key"))
}

func rejectExecutionTargetQuery(r *http.Request) *apiproblem.Problem {
	if r == nil || r.URL == nil {
		return nil
	}
	if r.URL.RawQuery != "" {
		return apiproblem.New(apiproblem.CodeMalformedRequest).WithDetail("query parameters are not accepted")
	}
	return nil
}

func requireApplicationJSON(contentType string) *apiproblem.Problem {
	ct := strings.TrimSpace(contentType)
	if ct == "" {
		return apiproblem.New(apiproblem.CodeUnsupportedMediaType).WithDetail("Content-Type is required; accepted type is application/json")
	}
	mediaType, _, err := mime.ParseMediaType(ct)
	if err != nil {
		return apiproblem.New(apiproblem.CodeUnsupportedMediaType).WithDetail("Content-Type could not be parsed")
	}
	if strings.ToLower(mediaType) != "application/json" {
		return apiproblem.New(apiproblem.CodeUnsupportedMediaType).WithDetail("unsupported media type; accepted type is application/json")
	}
	return nil
}

func extractExecutionTargetCreateUIDs(phaseOne PhaseOneResult) (participationUID, stackUID string, prob *apiproblem.Problem) {
	var root map[string]any
	if err := json.Unmarshal(phaseOne.Body, &root); err != nil {
		return "", "", apiproblem.New(apiproblem.CodeMalformedRequest).WithDetail("malformed JSON")
	}
	spec, _ := root["spec"].(map[string]any)
	if spec == nil {
		return "", "", nil
	}
	if ref, ok := parseTypedRef(spec["cloudProviderParticipationRef"]); ok {
		participationUID = ref.UID
	}
	if ref, ok := parseTypedRef(spec["infrastructureStackRef"]); ok {
		stackUID = ref.UID
	}
	return participationUID, stackUID, nil
}

func classifyExecutionTargetCreate(body []byte) (prob *apiproblem.Problem, clientStatusAttempted bool) {
	var root map[string]any
	if err := json.Unmarshal(body, &root); err != nil {
		return apiproblem.New(apiproblem.CodeMalformedRequest).WithDetail("malformed JSON"), false
	}
	paths := collectExecutionTargetPaths(root, "")
	sort.Strings(paths)

	allowed := map[string]struct{}{
		"metadata.name":                          {},
		"spec.cloudProviderParticipationRef.uid": {},
		"spec.infrastructureStackRef.uid":        {},
		"spec.targetClass":                       {},
	}
	structural := map[string]struct{}{
		"metadata":                           {},
		"spec":                               {},
		"spec.cloudProviderParticipationRef": {},
		"spec.infrastructureStackRef":        {},
	}

	var (
		statusPaths  []string
		systemOwned  []string
		scopePaths   []string
		unknownPaths []string
	)
	for _, p := range paths {
		switch {
		case p == "status" || strings.HasPrefix(p, "status."):
			statusPaths = append(statusPaths, p)
			clientStatusAttempted = true
		case p == "metadata.scopeRef" || strings.HasPrefix(p, "metadata.scopeRef."):
			scopePaths = append(scopePaths, p)
			clientStatusAttempted = true
		case isExecutionTargetSystemOwnedPath(p):
			systemOwned = append(systemOwned, p)
			clientStatusAttempted = true
		case p == "apiVersion" || p == "kind":
			unknownPaths = append(unknownPaths, p)
		default:
			if _, ok := allowed[p]; ok {
				continue
			}
			if _, ok := structural[p]; ok {
				continue
			}
			if isAllowedExecutionTargetPath(p, allowed, structural) {
				continue
			}
			unknownPaths = append(unknownPaths, p)
		}
	}

	if len(statusPaths) > 0 {
		return apiproblem.New(apiproblem.CodeAuthorizationDenied).WithViolations([]apiproblem.Violation{{
			Field:   "/status",
			Code:    validate.ViolationStatusFieldWrite,
			Message: "status is api-server-owned and cannot be client-authored",
		}}), true
	}
	if len(systemOwned) > 0 {
		return apiproblem.New(apiproblem.CodeAuthorizationDenied).WithViolations([]apiproblem.Violation{{
			Field:   "/" + strings.ReplaceAll(systemOwned[0], ".", "/"),
			Code:    validate.ViolationSystemOwnedFieldWrite,
			Message: "api-server-owned metadata field cannot be client-authored",
		}}), true
	}
	if len(scopePaths) > 0 {
		return apiproblem.New(apiproblem.CodeAuthorizationDenied).WithViolations([]apiproblem.Violation{{
			Field:   "/metadata/scopeRef",
			Code:    validate.ViolationSystemOwnedFieldWrite,
			Message: "metadata.scopeRef cannot be client-authored",
		}}), true
	}
	if len(unknownPaths) > 0 {
		return apiproblem.New(apiproblem.CodeUnknownField).WithViolations([]apiproblem.Violation{{
			Field:   "/" + strings.ReplaceAll(unknownPaths[0], ".", "/"),
			Code:    apiproblem.ViolationUnknownField,
			Message: "unknown field",
		}}), clientStatusAttempted
	}
	return nil, clientStatusAttempted
}

func isExecutionTargetSystemOwnedPath(path string) bool {
	owned := []string{"uid", "generation", "resourceVersion", "createdAt", "updatedAt", "deletionTimestamp", "ownerRefs"}
	for _, f := range owned {
		if path == "metadata."+f || strings.HasPrefix(path, "metadata."+f+".") {
			return true
		}
	}
	return false
}

func isAllowedExecutionTargetPath(path string, allowed, structural map[string]struct{}) bool {
	if _, ok := allowed[path]; ok {
		return true
	}
	if _, ok := structural[path]; ok {
		return true
	}
	for a := range allowed {
		if strings.HasPrefix(path, a+".") {
			return true
		}
	}
	return false
}

func collectExecutionTargetPaths(v any, prefix string) []string {
	switch t := v.(type) {
	case map[string]any:
		if len(t) == 0 && prefix != "" {
			return []string{prefix}
		}
		keys := make([]string, 0, len(t))
		for k := range t {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		var out []string
		for _, k := range keys {
			child := k
			if prefix != "" {
				child = prefix + "." + k
			}
			out = append(out, collectExecutionTargetPaths(t[k], child)...)
		}
		return out
	case []any:
		if prefix != "" {
			return []string{prefix}
		}
		return nil
	default:
		if prefix == "" {
			return nil
		}
		return []string{prefix}
	}
}

func validateExecutionTargetCreateGraph(
	name string,
	part cmmodel.CloudProviderParticipation,
	stack cmmodel.InfrastructureStack,
	body []byte,
) (etmodel.TargetClass, *apiproblem.Problem) {
	if name == "" {
		return "", apiproblem.New(apiproblem.CodeValidationFailed).WithViolations([]apiproblem.Violation{{
			Field:   "/metadata/name",
			Code:    apiproblem.ViolationOutOfRange,
			Message: "metadata.name is required",
		}})
	}
	if part.Metadata.UID == "" {
		return "", apiproblem.New(apiproblem.CodeValidationFailed).WithViolations([]apiproblem.Violation{{
			Field:   "/spec/cloudProviderParticipationRef/uid",
			Code:    apiproblem.ViolationOutOfRange,
			Message: "participation uid is required",
		}})
	}
	if stack.Metadata.UID == "" {
		return "", apiproblem.New(apiproblem.CodeValidationFailed).WithViolations([]apiproblem.Violation{{
			Field:   "/spec/infrastructureStackRef/uid",
			Code:    apiproblem.ViolationOutOfRange,
			Message: "stack uid is required",
		}})
	}

	var root map[string]any
	_ = json.Unmarshal(body, &root)
	spec, _ := root["spec"].(map[string]any)
	rawClass, _ := spec["targetClass"].(string)
	targetClass := etmodel.TargetClass(rawClass)
	if targetClass == "" {
		targetClass = etmodel.TargetClassSyntheticIaaS
	}
	if targetClass != etmodel.TargetClassSyntheticIaaS {
		return "", apiproblem.New(apiproblem.CodeValidationFailed).WithViolations([]apiproblem.Violation{{
			Field:   "/spec/targetClass",
			Code:    apiproblem.ViolationOutOfRange,
			Message: "targetClass must be synthetic-iaas",
		}})
	}

	if !participationEffectiveActive(part) {
		return "", apiproblem.New(apiproblem.CodeConflict).WithViolations([]apiproblem.Violation{{
			Field:   "/spec/cloudProviderParticipationRef/uid",
			Code:    executiontarget.ViolationExecutionTargetParticipationUnavailable,
			Message: "CloudProviderParticipation is not effective-Active",
		}})
	}
	if stack.Status.Phase != cmmodel.InfrastructureStackPhaseActive {
		return "", apiproblem.New(apiproblem.CodeConflict).WithViolations([]apiproblem.Violation{{
			Field:   "/spec/infrastructureStackRef/uid",
			Code:    executiontarget.ViolationExecutionTargetStackUnavailable,
			Message: "InfrastructureStack is not Active",
		}})
	}

	partScope := part.Spec.CloudProviderRef.UID
	stackScope := apimeta.CanonicalScopeIdentity(stack.Metadata.ScopeRef).UID
	if partScope == "" || stackScope == "" || partScope != stackScope {
		return "", apiproblem.New(apiproblem.CodeValidationFailed).WithViolations([]apiproblem.Violation{{
			Field:   "/spec",
			Code:    executiontarget.ViolationExecutionTargetScopeMismatch,
			Message: "participation and stack must share the same CloudProvider scope",
		}})
	}
	return targetClass, nil
}

func participationEffectiveActive(part cmmodel.CloudProviderParticipation) bool {
	return part.Status.Phase == cmmodel.ParticipationPhaseActive &&
		!part.Status.PlatformSuspended &&
		!part.Status.ProviderSuspended
}

func (h *ExecutionTargetCollectionHandler) resolveBackingViability(et etmodel.ExecutionTarget) executiontarget.BackingViability {
	out := executiontarget.BackingViability{
		ParticipationUID: et.Spec.CloudProviderParticipationRef.UID,
		StackUID:         et.Spec.InfrastructureStackRef.UID,
	}
	if h.Cloud == nil {
		return out
	}
	if part, ok := h.Cloud.GetParticipation(et.Spec.CloudProviderParticipationRef.UID); ok {
		out.ParticipationEffectiveActive = participationEffectiveActive(part)
		out.ParticipationScopeUID = part.Spec.CloudProviderRef.UID
	}
	if stack, ok := h.Cloud.GetInfrastructureStack(et.Spec.InfrastructureStackRef.UID); ok {
		out.StackPhase = string(stack.Status.Phase)
		out.StackScopeUID = apimeta.CanonicalScopeIdentity(stack.Metadata.ScopeRef).UID
		out.StackGeneration = stack.Metadata.Generation
	}
	return out
}

func (h *ExecutionTargetCollectionHandler) writeAuditedAuthorizationDenial(
	w http.ResponseWriter,
	r *http.Request,
	principal string,
	suppressed *apiproblem.Problem,
) {
	event := executiontarget.AssembleAuthorizationDenial(executiontarget.RedactedAuditInput{
		UID:       newAuditUID("et-authz-denied"),
		RequestID: requestID(r),
		Actor:     actorRef(principal),
	})
	h.appendDenialOrInternal(w, r, event, suppressed)
}

func (h *ExecutionTargetCollectionHandler) writeAuditedSafeDenial(
	w http.ResponseWriter,
	r *http.Request,
	principal string,
) {
	suppressed := apiproblem.New(apiproblem.CodeResourceNotFound).WithViolations([]apiproblem.Violation{{
		Code:    violationAuthorizationSafeDenial,
		Message: "resource not found",
	}})
	event := executiontarget.AssembleSafeDenial(executiontarget.RedactedAuditInput{
		UID:       newAuditUID("et-safe-denial"),
		RequestID: requestID(r),
		Actor:     actorRef(principal),
	})
	h.appendDenialOrInternal(w, r, event, suppressed)
}

func (h *ExecutionTargetCollectionHandler) writeClassifiedCreateDenial(
	w http.ResponseWriter,
	r *http.Request,
	principal string,
	prob *apiproblem.Problem,
) {
	if prob == nil {
		return
	}
	if prob.Code == apiproblem.CodeAuthorizationDenied &&
		(hasViolation(prob, validate.ViolationStatusFieldWrite) || hasViolation(prob, validate.ViolationSystemOwnedFieldWrite)) {
		h.writeAuditedAuthorizationDenial(w, r, principal, prob)
		return
	}
	writeProblem(w, r, prob)
}

func (h *ExecutionTargetCollectionHandler) appendDenialOrInternal(
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
		// Never disclose the suppressed outcome on append failure.
		writeProblem(w, r, apiproblem.New(apiproblem.CodeInternalError).WithDetail("required AuditEvent append failed"))
		return
	}
	writeProblem(w, r, suppressed)
}

func writeExecutionTargetCreated(w http.ResponseWriter, r *http.Request, resourceVersion string, body []byte) {
	w.Header().Set("ETag", resourceVersion)
	writeRawJSON(w, r, http.StatusCreated, mediaJSON, body)
}

func writeExecutionTargetReplay(w http.ResponseWriter, r *http.Request, replay *executiontarget.CompletedResult) {
	if replay == nil {
		writeProblem(w, r, apiproblem.New(apiproblem.CodeInternalError).WithDetail("missing replay payload"))
		return
	}
	// Fresh correlation only; never echo Idempotency-Key, stored correlation,
	// credentials, or request headers.
	writeRawJSON(w, r, replay.StatusCode, mediaJSON, replay.Body)
}

func idempotencyKeyReuseMismatchProblem() *apiproblem.Problem {
	return apiproblem.New(apiproblem.CodeConflict).
		WithDetail("Idempotency-Key was reused with a different request digest").
		WithViolations([]apiproblem.Violation{{
			Field:   "/headers/Idempotency-Key",
			Code:    cloudmodel.ViolationIdempotencyKeyReuseMismatch,
			Message: "idempotency key reused with a different canonical digest",
		}})
}
