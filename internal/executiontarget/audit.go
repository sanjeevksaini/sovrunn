package executiontarget

import (
	"time"

	"github.com/sanjeevksaini/sovrunn/internal/apiconform"
	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/executiontarget/model"
)

// Audit category action names used in redacted FEATURE-0013 AuditEvent records
// for the FEATURE-0016 audit matrix (REQ-F16-09; design §6.3; ADH-2026-058
// clause 8).
const (
	AuditActionCreate              = "executiontarget.create"
	AuditActionQualify             = "executiontarget.qualify"
	AuditActionRetire              = "executiontarget.retire"
	AuditActionMaintenanceEnter    = "executiontarget.maintenance.enter"
	AuditActionMaintenanceClear    = "executiontarget.maintenance.clear"
	AuditActionAuthorizationDenial = "authorization.denied"
	AuditActionSafeDenial          = "authorization.safe_denial"
	AuditActionExpiry              = "executiontarget.expire"

	// SchedulerSystemActorUID is the deterministic system actor for expiry
	// AuditEvents assembled for the TASK-F16-06 scheduler.
	SchedulerSystemActorUID = "system:executiontarget-expiry-scheduler"

	auditAPIVersion = "governance.sovrunn.io/v1alpha1"
	auditKind       = "AuditEvent"

	reasonCreated                = "created"
	reasonQualificationCompleted = "qualification.completed"
	reasonRetired                = "retired"
	reasonMaintenanceEntered     = "maintenance.entered"
	reasonMaintenanceCleared     = "maintenance.cleared"
	reasonAuthorizationDenied    = "authorization.denied"
	reasonSafeDenial             = "authorization.safe_denial"
	reasonExpired                = "executiontarget.expired"
)

// RedactedAuditInput is the safe surface for assembling a FEATURE-0016
// AuditEvent. Callers must not pass secrets, credentials, raw Idempotency-Key,
// request bodies, inaccessible-resource existence, or a complete idempotency
// namespace.
//
// Assembly functions are build-only: they construct and return an in-memory
// AuditEvent. Exactly one append owner (lifecycle service, expiry scheduler, or
// owning HTTP handler) performs exactly one inherited append attempt after
// assembly. Assembly never calls, owns, or retries append.
type RedactedAuditInput struct {
	UID                    string
	RequestID              string
	Reason                 string
	Actor                  apimeta.TypedRef
	Subject                apimeta.TypedRef
	SubjectResourceVersion string
}

// AssembleSuccessfulCreate builds the redacted AuditEvent for a successful
// ExecutionTarget create (REQ-F16-09 audit matrix).
func AssembleSuccessfulCreate(in RedactedAuditInput) apiconform.AuditEvent {
	return newRedactedAuditEvent(in, AuditActionCreate, apiconform.AuditOutcomeSucceeded, reasonOrDefault(in.Reason, reasonCreated))
}

// AssembleCompletedQualification builds the redacted AuditEvent for a completed
// qualification conclusion publication (REQ-F16-09 audit matrix).
func AssembleCompletedQualification(in RedactedAuditInput) apiconform.AuditEvent {
	return newRedactedAuditEvent(in, AuditActionQualify, apiconform.AuditOutcomeSucceeded, reasonOrDefault(in.Reason, reasonQualificationCompleted))
}

// AssembleRetirement builds the redacted AuditEvent for a successful retirement
// (REQ-F16-09 audit matrix).
func AssembleRetirement(in RedactedAuditInput) apiconform.AuditEvent {
	return newRedactedAuditEvent(in, AuditActionRetire, apiconform.AuditOutcomeSucceeded, reasonOrDefault(in.Reason, reasonRetired))
}

// AssembleMaintenanceEnter builds the redacted AuditEvent for a successful
// Maintenance entry (REQ-F16-09 audit matrix).
func AssembleMaintenanceEnter(in RedactedAuditInput) apiconform.AuditEvent {
	return newRedactedAuditEvent(in, AuditActionMaintenanceEnter, apiconform.AuditOutcomeSucceeded, reasonOrDefault(in.Reason, reasonMaintenanceEntered))
}

// AssembleMaintenanceClear builds the redacted AuditEvent for a successful
// Maintenance clear (REQ-F16-09 audit matrix).
func AssembleMaintenanceClear(in RedactedAuditInput) apiconform.AuditEvent {
	return newRedactedAuditEvent(in, AuditActionMaintenanceClear, apiconform.AuditOutcomeSucceeded, reasonOrDefault(in.Reason, reasonMaintenanceCleared))
}

// AssembleAuthorizationDenial builds the redacted AuditEvent for an
// authenticated authorization denial (REQ-F16-09 audit matrix).
func AssembleAuthorizationDenial(in RedactedAuditInput) apiconform.AuditEvent {
	return newRedactedAuditEvent(in, AuditActionAuthorizationDenial, apiconform.AuditOutcomeDenied, reasonOrDefault(in.Reason, reasonAuthorizationDenied))
}

// AssembleSafeDenial builds the redacted AuditEvent for an authenticated safe
// denial (REQ-F16-09 audit matrix). The event must never disclose target or
// backing-reference existence.
func AssembleSafeDenial(in RedactedAuditInput) apiconform.AuditEvent {
	return newRedactedAuditEvent(in, AuditActionSafeDenial, apiconform.AuditOutcomeDenied, reasonOrDefault(in.Reason, reasonSafeDenial))
}

// AssembleExpiry builds the redacted AuditEvent for fact/target expiry
// (REQ-F16-09; sole background audit exception). The TASK-F16-06 scheduler
// calls this on each retry attempt and owns append timing/cadence; this
// function only assembles the event struct and forces the deterministic
// system actor.
func AssembleExpiry(in RedactedAuditInput) apiconform.AuditEvent {
	out := in
	out.Actor = ExpirySystemActor()
	return newRedactedAuditEvent(out, AuditActionExpiry, apiconform.AuditOutcomeSucceeded, reasonOrDefault(out.Reason, reasonExpired))
}

// ExpirySystemActor returns the deterministic system actor TypedRef used for
// expiry AuditEvents.
func ExpirySystemActor() apimeta.TypedRef {
	return apimeta.TypedRef{
		APIVersion: model.APIVersionExecutionTarget,
		Kind:       "SystemActor",
		Name:       "executiontarget-expiry-scheduler",
		UID:        SchedulerSystemActorUID,
	}
}

func newRedactedAuditEvent(
	in RedactedAuditInput,
	action string,
	outcome apiconform.AuditOutcome,
	reason string,
) apiconform.AuditEvent {
	now := time.Now().UTC().Format(time.RFC3339)
	ev := apiconform.AuditEvent{
		TypeMeta: apimeta.TypeMeta{APIVersion: auditAPIVersion, Kind: auditKind},
		Metadata: apimeta.ObjectMeta{
			Name:      in.UID,
			UID:       in.UID,
			CreatedAt: now,
		},
		Record: apiconform.AuditEventRecord{
			ActorRef:               in.Actor,
			RequestID:              in.RequestID,
			SubjectRef:             in.Subject,
			SubjectResourceVersion: in.SubjectResourceVersion,
			Action:                 action,
			Outcome:                outcome,
			ReasonCode:             reason,
		},
	}
	return redactAuditEvent(ev)
}

func redactAuditEvent(event apiconform.AuditEvent) apiconform.AuditEvent {
	// Defensive copy: keep only UID-pinned identity on actor/subject. Request
	// correlation (requestId) is retained; no secret-bearing payload fields
	// exist on the closed AuditEvent shape.
	out := event
	out.Record.ActorRef = uidPinned(event.Record.ActorRef)
	out.Record.SubjectRef = uidPinned(event.Record.SubjectRef)
	return out
}

func uidPinned(ref apimeta.TypedRef) apimeta.TypedRef {
	return apimeta.TypedRef{
		APIVersion: ref.APIVersion,
		Kind:       ref.Kind,
		Name:       ref.Name,
		UID:        ref.UID,
	}
}

func reasonOrDefault(reason, fallback string) string {
	if reason == "" {
		return fallback
	}
	return reason
}
