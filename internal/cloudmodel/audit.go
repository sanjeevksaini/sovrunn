package cloudmodel

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/sanjeevksaini/sovrunn/internal/apiconform"
	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/cloudmodel/model"
)

// Audit category action names used in redacted AuditEvent records.
const (
	AuditActionCollectionCreate      = "collection.create"
	AuditActionResourcePatch         = "resource.patch"
	AuditActionParticipationCreate   = "participation.create"
	AuditActionParticipationAction   = "participation.action"
	AuditActionParticipationExpiry   = "participation.expire"
	AuditActionStatusWrite           = "status.write"
	AuditActionSystemOwnedFieldWrite = "system_owned_field.write"
	AuditActionForgedGrant           = "forged_grant"
	AuditActionRootDenial            = "cloudplatform.root_required"
	AuditActionListWithoutReadGrant  = "list.without_read_grant"
	AuditActionSafeDenial            = "authorization.safe_denial"

	// SchedulerSystemActorUID is the deterministic system actor for scheduler
	// expiry AuditEvents (ADH-2026-046).
	SchedulerSystemActorUID = "system:participation-expiry-scheduler"

	auditAPIVersion = "governance.sovrunn.io/v1alpha1"
	auditKind       = "AuditEvent"
)

// AuditAppender appends a redacted FEATURE-0013 AuditEvent. FEATURE-0015 does
// not own a durable AuditEvent store; the appender is injected (VS0-WRITER-011
// writer authority remains FEATURE-0013).
type AuditAppender interface {
	Append(ctx context.Context, event apiconform.AuditEvent) error
}

// MemoryAuditAppender is a process-local append-only log for tests and the
// Phase 2R in-memory composition. It is not a durable or external store.
type MemoryAuditAppender struct {
	mu     sync.Mutex
	events []apiconform.AuditEvent
	fail   error
}

// SetFail injects the next Append error (test helper). A nil error clears it.
func (a *MemoryAuditAppender) SetFail(err error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.fail = err
}

// Append stores a copy of the event or returns the injected failure.
func (a *MemoryAuditAppender) Append(ctx context.Context, event apiconform.AuditEvent) error {
	_ = ctx
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.fail != nil {
		err := a.fail
		return err
	}
	a.events = append(a.events, event)
	return nil
}

// Events returns a copy of appended events.
func (a *MemoryAuditAppender) Events() []apiconform.AuditEvent {
	a.mu.Lock()
	defer a.mu.Unlock()
	out := make([]apiconform.AuditEvent, len(a.events))
	copy(out, a.events)
	return out
}

// Len returns the number of appended events.
func (a *MemoryAuditAppender) Len() int {
	a.mu.Lock()
	defer a.mu.Unlock()
	return len(a.events)
}

// IdempotencyCompletion is the optional staged idempotency completion published
// only after a required AuditEvent append succeeds.
type IdempotencyCompletion struct {
	Namespace IdempotencyNamespace
	Digest    Digest
	Result    CompletedResult
}

// AuditedSuccess is a staged successful mutation prepared for atomic
// audit-before-publication (design DD-09).
type AuditedSuccess struct {
	Staged      *StagedOutcome
	Event       apiconform.AuditEvent
	Idempotency *IdempotencyCompletion
	// AfterPublish runs after successful Publish (for example scheduler Signal).
	AfterPublish func()
}

// AuditedDenial is a staged audited denial response.
type AuditedDenial struct {
	Event   apiconform.AuditEvent
	Problem *apiproblem.Problem
}

// PublicationCoordinator serializes under-lock recheck/staging (Task 4
// primitives), appends the required AuditEvent, then publishes the staged
// resource/lifecycle/idempotency outcome (design DD-09, §5.1 step 14).
type PublicationCoordinator struct {
	store *Store
	audit AuditAppender
	idemp *IdempotencyCoordinator
}

// PublicationCoordinatorConfig configures the coordinator.
type PublicationCoordinatorConfig struct {
	Store       *Store
	Audit       AuditAppender
	Idempotency *IdempotencyCoordinator
}

// NewPublicationCoordinator constructs the audit/mutation coordinator.
func NewPublicationCoordinator(cfg PublicationCoordinatorConfig) *PublicationCoordinator {
	return &PublicationCoordinator{
		store: cfg.Store,
		audit: cfg.Audit,
		idemp: cfg.Idempotency,
	}
}

// PublishSuccess appends the required AuditEvent before publishing the staged
// mutation and optional idempotency completion. On append failure it aborts the
// staged outcome, aborts any InFlight idempotency reservation, and returns
// INTERNAL_ERROR/500. DEPENDENCY_UNAVAILABLE is never used.
//
// The publication lock must already be held (BeginPublication).
func (c *PublicationCoordinator) PublishSuccess(ctx context.Context, success AuditedSuccess) *apiproblem.Problem {
	if c == nil || c.audit == nil {
		c.abortSuccess(success)
		return internalAuditError("audit appender is not configured")
	}
	if err := c.audit.Append(ctx, redactAuditEvent(success.Event)); err != nil {
		c.abortSuccess(success)
		return internalAuditError("required AuditEvent append failed")
	}
	if success.Staged != nil {
		c.store.Publish(success.Staged)
	}
	if success.Idempotency != nil && c.idemp != nil {
		c.idemp.Complete(success.Idempotency.Namespace, success.Idempotency.Digest, success.Idempotency.Result)
	}
	if success.AfterPublish != nil {
		success.AfterPublish()
	}
	return nil
}

// PublishDenial appends the required AuditEvent before returning the denial.
// On append failure it suppresses the denial and returns INTERNAL_ERROR/500.
func (c *PublicationCoordinator) PublishDenial(ctx context.Context, denial AuditedDenial) *apiproblem.Problem {
	if c == nil || c.audit == nil {
		return internalAuditError("audit appender is not configured")
	}
	if err := c.audit.Append(ctx, redactAuditEvent(denial.Event)); err != nil {
		return internalAuditError("required AuditEvent append failed")
	}
	if denial.Problem == nil {
		return apiproblem.New(apiproblem.CodeAuthorizationDenied)
	}
	return denial.Problem
}

// PublishSuccessAllowingPostAppendFailure is a test/support helper that appends
// first, then invokes publishFn. If publishFn fails after a successful append,
// the already-appended AuditEvent may remain (ADH-2026-048/049). Production
// PublishSuccess uses Store.Publish which does not fail after append.
func (c *PublicationCoordinator) PublishSuccessAllowingPostAppendFailure(
	ctx context.Context,
	success AuditedSuccess,
	publishFn func() error,
) *apiproblem.Problem {
	if c == nil || c.audit == nil {
		c.abortSuccess(success)
		return internalAuditError("audit appender is not configured")
	}
	if err := c.audit.Append(ctx, redactAuditEvent(success.Event)); err != nil {
		c.abortSuccess(success)
		return internalAuditError("required AuditEvent append failed")
	}
	if publishFn != nil {
		if err := publishFn(); err != nil {
			c.abortSuccess(success)
			return internalAuditError("publication failed after AuditEvent append")
		}
	} else if success.Staged != nil {
		c.store.Publish(success.Staged)
	}
	if success.Idempotency != nil && c.idemp != nil {
		c.idemp.Complete(success.Idempotency.Namespace, success.Idempotency.Digest, success.Idempotency.Result)
	}
	return nil
}

func (c *PublicationCoordinator) abortSuccess(success AuditedSuccess) {
	if c == nil {
		return
	}
	if success.Staged != nil && c.store != nil {
		c.store.Abort(success.Staged)
	}
	if success.Idempotency != nil && c.idemp != nil {
		c.idemp.Abort(success.Idempotency.Namespace)
	}
}

// ExpireParticipation is the audit-aware api-server expiry operation injected
// into the ParticipationExpiryScheduler (design §4.6). The participation action
// concurrency guard is held by the scheduler across recheck and this call.
func (c *PublicationCoordinator) ExpireParticipation(ctx context.Context, uid, observedVersion, schedulerExecutionID string) error {
	if c == nil || c.store == nil {
		return errors.New("publication coordinator is not configured")
	}
	c.store.BeginPublication()
	defer c.store.EndPublication()

	p, ok := c.store.participations[uid]
	if !ok {
		return nil
	}
	if p.Status.Phase != model.ParticipationPhasePending {
		return nil
	}
	if p.Metadata.ResourceVersion != observedVersion {
		return nil
	}
	next, prob := ApplyParticipationAction(p.Status, ParticipationActionExpire, "")
	if prob != nil {
		return nil
	}
	p.Status = next
	staged, prob := c.store.StageUpdateParticipation(p)
	if prob != nil {
		return errors.New(string(prob.Code))
	}

	event := NewRedactedAuditEvent(RedactedAuditInput{
		UID:       "audit-expiry-" + uid + "-" + schedulerExecutionID,
		RequestID: schedulerExecutionID,
		Action:    AuditActionParticipationExpiry,
		Outcome:   apiconform.AuditOutcomeSucceeded,
		Reason:    "participation.expired",
		Actor: apimeta.TypedRef{
			APIVersion: model.APIVersionCloudProviderParticipation,
			Kind:       "SystemActor",
			Name:       "participation-expiry-scheduler",
			UID:        SchedulerSystemActorUID,
		},
		Subject: apimeta.TypedRef{
			APIVersion: model.APIVersionCloudProviderParticipation,
			Kind:       model.KindCloudProviderParticipation,
			Name:       p.Metadata.Name,
			UID:        p.Metadata.UID,
		},
		SubjectResourceVersion: observedVersion,
	})

	if pubProb := c.PublishSuccess(ctx, AuditedSuccess{Staged: staged, Event: event}); pubProb != nil {
		return errors.New(string(pubProb.Code))
	}
	return nil
}

// RedactedAuditInput is the safe surface for building a FEATURE-0015 AuditEvent.
// Callers must not pass secrets, credentials, raw Idempotency-Key, request
// bodies, inaccessible-resource existence, or a complete idempotency namespace.
type RedactedAuditInput struct {
	UID                    string
	RequestID              string
	Action                 string
	Outcome                apiconform.AuditOutcome
	Reason                 string
	Actor                  apimeta.TypedRef
	Subject                apimeta.TypedRef
	SubjectResourceVersion string
}

// NewRedactedAuditEvent builds a redacted AuditEvent with request correlation.
func NewRedactedAuditEvent(in RedactedAuditInput) apiconform.AuditEvent {
	now := time.Now().UTC().Format(time.RFC3339)
	return apiconform.AuditEvent{
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
			Action:                 in.Action,
			Outcome:                in.Outcome,
			ReasonCode:             in.Reason,
		},
	}
}

func redactAuditEvent(event apiconform.AuditEvent) apiconform.AuditEvent {
	// Defensive copy: strip any accidental secret-bearing name fields from
	// actor/subject by keeping only UID-pinned identity. Request correlation
	// (requestId) is retained.
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

func internalAuditError(detail string) *apiproblem.Problem {
	return apiproblem.New(apiproblem.CodeInternalError).WithDetail(detail)
}
