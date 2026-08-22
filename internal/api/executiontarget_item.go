package api

import (
	"context"
	"net/http"

	"github.com/sanjeevksaini/sovrunn/internal/apiconform"
	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/cloudmodel"
	"github.com/sanjeevksaini/sovrunn/internal/executiontarget"
	etmodel "github.com/sanjeevksaini/sovrunn/internal/executiontarget/model"
)

// FEATURE-0016 item GET ServeMux pattern (design §4.5; ADH-2026-058 clause 3).
const PatternExecutionTargetItemGet = "GET /apis/execution.sovrunn.io/v1alpha1/execution-targets/{uid}"

// ExecutionTargetItemHandler serves ExecutionTarget item GET per design §5.1.
// Idempotency-Key is ignored: the handler never creates, waits on, reads, or
// replays an idempotency record (REQ-F16-04 / REQ-F16-06).
type ExecutionTargetItemHandler struct {
	Lifecycle *executiontarget.ExecutionTargetLifecycleService
	Cloud     *cloudmodel.Store
	Grants    GrantResolver
	Audit     executiontarget.AuditAppender
}

// NewExecutionTargetItemHandler constructs the item GET handler.
func NewExecutionTargetItemHandler(
	lifecycle *executiontarget.ExecutionTargetLifecycleService,
	cloud *cloudmodel.Store,
	grants GrantResolver,
	audit executiontarget.AuditAppender,
) *ExecutionTargetItemHandler {
	return &ExecutionTargetItemHandler{
		Lifecycle: lifecycle,
		Cloud:     cloud,
		Grants:    grants,
		Audit:     audit,
	}
}

// ServeHTTP dispatches GET only. Other methods → 405 (transport guard also
// rejects disallowed methods before this handler for F0016 paths).
func (h *ExecutionTargetItemHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	uid := r.PathValue("uid")
	switch r.Method {
	case http.MethodGet:
		h.get(w, r, uid)
	default:
		writeMethodNotAllowed(w, r, "GET")
	}
}

func (h *ExecutionTargetItemHandler) get(w http.ResponseWriter, r *http.Request, uid string) {
	// Ordered pipeline (design §5.1 item GET): authentication → phase-one
	// path-UID safe resolution → derived-scope read authorization → safe
	// access → safe projection. Query and malformed UID are rejected before
	// safe resolution. Idempotency-Key is never inspected.

	principal, prob := authenticatePrincipal(r, h.Grants)
	if prob != nil {
		writeProblem(w, r, prob)
		return
	}
	if prob := rejectExecutionTargetQuery(r); prob != nil {
		writeProblem(w, r, prob)
		return
	}
	if !apimeta.IsGeneratedUIDFormat(uid) {
		writeProblem(w, r, apiproblem.New(apiproblem.CodeMalformedRequest).WithDetail("path uid segment is malformed"))
		return
	}
	if h.Lifecycle == nil {
		writeProblem(w, r, apiproblem.New(apiproblem.CodeInternalError).WithDetail("lifecycle service is not configured"))
		return
	}

	// Phase-one path-UID safe resolution: read only the target needed to derive
	// CloudProvider scope. Missing/inaccessible target → audited safe 404.
	et, ok := h.Lifecycle.Store().GetExecutionTarget(uid)
	if !ok {
		h.writeAuditedSafeDenial(w, r, principal)
		return
	}

	// Derived-scope read authorization after successful UID-only resolution.
	scope := apimeta.ScopeIdentity{Kind: apimeta.ScopeCloudProvider, UID: et.CloudProviderScopeUID}
	if !h.Grants.AuthorizeExact(ActionExecutionTargetRead, scope, "", "") {
		h.writeAuditedAuthorizationDenial(w, r, principal,
			apiproblem.New(apiproblem.CodeAuthorizationDenied).WithDetail("executiontarget.read grant is required"),
		)
		return
	}

	backing := h.resolveBackingViability(et)
	snap, ok := h.Lifecycle.CaptureProjectionSnapshot(uid, backing)
	if !ok {
		h.writeAuditedSafeDenial(w, r, principal)
		return
	}
	proj := executiontarget.Project(snap)
	w.Header().Set("ETag", proj.Metadata.ResourceVersion)
	writeJSONSuccess(w, r, http.StatusOK, proj)
}

func (h *ExecutionTargetItemHandler) resolveBackingViability(et etmodel.ExecutionTarget) executiontarget.BackingViability {
	if h.Cloud == nil || h.Grants == nil {
		return executiontarget.BackingViability{
			ParticipationUID: et.Spec.CloudProviderParticipationRef.UID,
			StackUID:         et.Spec.InfrastructureStackRef.UID,
		}
	}
	access := executiontarget.NewBackingAccessProvider(h.Cloud).Access(
		context.Background(), "", ActionExecutionTargetRead,
		et.Spec.CloudProviderParticipationRef.UID,
		et.Spec.InfrastructureStackRef.UID,
		h.Grants,
	)
	if access.Disposition != executiontarget.BackingAllowed {
		return executiontarget.BackingViability{
			ParticipationUID: et.Spec.CloudProviderParticipationRef.UID,
			StackUID:         et.Spec.InfrastructureStackRef.UID,
		}
	}
	return access.View.AsViability()
}

func (h *ExecutionTargetItemHandler) writeAuditedAuthorizationDenial(
	w http.ResponseWriter,
	r *http.Request,
	principal string,
	suppressed *apiproblem.Problem,
) {
	event := executiontarget.AssembleAuthorizationDenial(executiontarget.RedactedAuditInput{
		UID:       newAuditUID("et-item-authz-denied"),
		RequestID: requestID(r),
		Actor:     actorRef(principal),
	})
	h.appendDenialOrInternal(w, r, event, suppressed)
}

func (h *ExecutionTargetItemHandler) writeAuditedSafeDenial(
	w http.ResponseWriter,
	r *http.Request,
	principal string,
) {
	suppressed := apiproblem.New(apiproblem.CodeResourceNotFound).WithViolations([]apiproblem.Violation{{
		Code:    violationAuthorizationSafeDenial,
		Message: "resource not found",
	}})
	event := executiontarget.AssembleSafeDenial(executiontarget.RedactedAuditInput{
		UID:       newAuditUID("et-item-safe-denial"),
		RequestID: requestID(r),
		Actor:     actorRef(principal),
	})
	h.appendDenialOrInternal(w, r, event, suppressed)
}

// appendDenialOrInternal performs exactly one inherited AuditEvent append
// attempt. Append failure returns only INTERNAL_ERROR/500 with no durable
// AuditEvent and no disclosure of the suppressed 403/404 (F16-102, F16-112).
func (h *ExecutionTargetItemHandler) appendDenialOrInternal(
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
