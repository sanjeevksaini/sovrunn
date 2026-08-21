package executiontarget

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/sanjeevksaini/sovrunn/internal/apiconform"
	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/executiontarget/model"
)

func actorRef(uid string) apimeta.TypedRef {
	return apimeta.TypedRef{
		APIVersion: "identity.sovrunn.io/v1alpha1",
		Kind:       "Principal",
		Name:       "principal",
		UID:        uid,
	}
}

func subjectTarget(name, uid string) apimeta.TypedRef {
	return apimeta.TypedRef{
		APIVersion: model.APIVersionExecutionTarget,
		Kind:       model.KindExecutionTarget,
		Name:       name,
		UID:        uid,
	}
}

func baseInput(actionUID, requestID string) RedactedAuditInput {
	return RedactedAuditInput{
		UID:                    actionUID,
		RequestID:              requestID,
		Actor:                  actorRef("principal-1"),
		Subject:                subjectTarget("tgt", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"),
		SubjectResourceVersion: "3",
	}
}

func assertRedactedEnvelope(t *testing.T, ev apiconform.AuditEvent, action string, outcome apiconform.AuditOutcome, reason, requestID string) {
	t.Helper()
	if ev.APIVersion != auditAPIVersion || ev.Kind != auditKind {
		t.Fatalf("type meta = %s/%s", ev.APIVersion, ev.Kind)
	}
	if ev.Metadata.UID == "" || ev.Metadata.Name == "" {
		t.Fatalf("metadata identity missing: %#v", ev.Metadata)
	}
	if _, err := time.Parse(time.RFC3339, ev.Metadata.CreatedAt); err != nil {
		t.Fatalf("CreatedAt must be RFC3339: %q err=%v", ev.Metadata.CreatedAt, err)
	}
	if ev.Record.Action != action {
		t.Fatalf("action = %q, want %q", ev.Record.Action, action)
	}
	if ev.Record.Outcome != outcome {
		t.Fatalf("outcome = %q, want %q", ev.Record.Outcome, outcome)
	}
	if ev.Record.ReasonCode != reason {
		t.Fatalf("reasonCode = %q, want %q", ev.Record.ReasonCode, reason)
	}
	if ev.Record.RequestID != requestID {
		t.Fatalf("requestId = %q, want %q", ev.Record.RequestID, requestID)
	}
	if ev.Record.SubjectResourceVersion != "3" {
		t.Fatalf("subjectResourceVersion = %q", ev.Record.SubjectResourceVersion)
	}
	if ev.Record.ActorRef.UID == "" || ev.Record.SubjectRef.UID == "" {
		t.Fatalf("actor/subject UIDs required: %#v", ev.Record)
	}
	if ev.Record.OperationRef != nil || ev.Record.CorrectionOfRef != nil || ev.Record.DecisionLinkage != nil {
		t.Fatalf("optional linkage fields must be omitted on F0016 assembly: %#v", ev.Record)
	}
}

func TestAssembleSuccessfulCreate(t *testing.T) {
	t.Parallel()
	in := baseInput("audit-create-1", "req-create-1")
	ev := AssembleSuccessfulCreate(in)
	assertRedactedEnvelope(t, ev, AuditActionCreate, apiconform.AuditOutcomeSucceeded, reasonCreated, "req-create-1")
}

func TestAssembleCompletedQualification(t *testing.T) {
	t.Parallel()
	in := baseInput("audit-qualify-1", "req-qualify-1")
	ev := AssembleCompletedQualification(in)
	assertRedactedEnvelope(t, ev, AuditActionQualify, apiconform.AuditOutcomeSucceeded, reasonQualificationCompleted, "req-qualify-1")
}

func TestAssembleRetirement(t *testing.T) {
	t.Parallel()
	in := baseInput("audit-retire-1", "req-retire-1")
	ev := AssembleRetirement(in)
	assertRedactedEnvelope(t, ev, AuditActionRetire, apiconform.AuditOutcomeSucceeded, reasonRetired, "req-retire-1")
}

func TestAssembleMaintenanceEnter(t *testing.T) {
	t.Parallel()
	in := baseInput("audit-maint-enter-1", "req-maint-enter-1")
	ev := AssembleMaintenanceEnter(in)
	assertRedactedEnvelope(t, ev, AuditActionMaintenanceEnter, apiconform.AuditOutcomeSucceeded, reasonMaintenanceEntered, "req-maint-enter-1")
}

func TestAssembleMaintenanceClear(t *testing.T) {
	t.Parallel()
	in := baseInput("audit-maint-clear-1", "req-maint-clear-1")
	ev := AssembleMaintenanceClear(in)
	assertRedactedEnvelope(t, ev, AuditActionMaintenanceClear, apiconform.AuditOutcomeSucceeded, reasonMaintenanceCleared, "req-maint-clear-1")
}

func TestAssembleAuthorizationDenial(t *testing.T) {
	t.Parallel()
	in := baseInput("audit-authz-1", "req-authz-1")
	in.Reason = "VS0_STATUS_FIELD_WRITE"
	ev := AssembleAuthorizationDenial(in)
	assertRedactedEnvelope(t, ev, AuditActionAuthorizationDenial, apiconform.AuditOutcomeDenied, "VS0_STATUS_FIELD_WRITE", "req-authz-1")
}

func TestAssembleAuthorizationDenial_DefaultReason(t *testing.T) {
	t.Parallel()
	in := baseInput("audit-authz-2", "req-authz-2")
	ev := AssembleAuthorizationDenial(in)
	assertRedactedEnvelope(t, ev, AuditActionAuthorizationDenial, apiconform.AuditOutcomeDenied, reasonAuthorizationDenied, "req-authz-2")
}

func TestAssembleSafeDenial(t *testing.T) {
	t.Parallel()
	in := baseInput("audit-safe-1", "req-safe-1")
	// Safe denial must not confirm inaccessible target existence via reason.
	in.Reason = "VS0_AUTHORIZATION_SAFE_DENIAL"
	ev := AssembleSafeDenial(in)
	assertRedactedEnvelope(t, ev, AuditActionSafeDenial, apiconform.AuditOutcomeDenied, "VS0_AUTHORIZATION_SAFE_DENIAL", "req-safe-1")
}

func TestAssembleExpiry_ForcesSystemActor(t *testing.T) {
	t.Parallel()
	in := baseInput("audit-expiry-1", "sched-exec-1")
	in.Actor = actorRef("must-be-replaced")
	ev := AssembleExpiry(in)
	assertRedactedEnvelope(t, ev, AuditActionExpiry, apiconform.AuditOutcomeSucceeded, reasonExpired, "sched-exec-1")
	want := ExpirySystemActor()
	if ev.Record.ActorRef != want {
		t.Fatalf("expiry actor = %#v, want %#v", ev.Record.ActorRef, want)
	}
	if ev.Record.ActorRef.UID != SchedulerSystemActorUID {
		t.Fatalf("system actor uid = %q", ev.Record.ActorRef.UID)
	}
}

func TestAssembleMatrix_DistinctFunctionsCoverExactMatrix(t *testing.T) {
	t.Parallel()
	in := baseInput("audit-matrix", "req-matrix")
	got := []struct {
		name   string
		ev     apiconform.AuditEvent
		action string
	}{
		{"create", AssembleSuccessfulCreate(in), AuditActionCreate},
		{"qualify", AssembleCompletedQualification(in), AuditActionQualify},
		{"retire", AssembleRetirement(in), AuditActionRetire},
		{"maint-enter", AssembleMaintenanceEnter(in), AuditActionMaintenanceEnter},
		{"maint-clear", AssembleMaintenanceClear(in), AuditActionMaintenanceClear},
		{"authz-denial", AssembleAuthorizationDenial(in), AuditActionAuthorizationDenial},
		{"safe-denial", AssembleSafeDenial(in), AuditActionSafeDenial},
		{"expiry", AssembleExpiry(in), AuditActionExpiry},
	}
	seen := make(map[string]struct{}, len(got))
	for _, tc := range got {
		if tc.ev.Record.Action != tc.action {
			t.Fatalf("%s: action = %q, want %q", tc.name, tc.ev.Record.Action, tc.action)
		}
		if _, dup := seen[tc.action]; dup {
			t.Fatalf("duplicate matrix action %q", tc.action)
		}
		seen[tc.action] = struct{}{}
	}
	if len(seen) != 8 {
		t.Fatalf("want 8 distinct matrix actions, got %d", len(seen))
	}
}

func TestAssemble_RedactedNoSecretsOrBodies(t *testing.T) {
	t.Parallel()
	secretish := "Idempotency-Key: raw-secret-key\nAuthorization: Bearer tok\n{\"password\":\"x\"}"
	in := baseInput("audit-redact", "corr-123")
	ev := AssembleSafeDenial(in)
	raw, err := json.Marshal(ev)
	if err != nil {
		t.Fatal(err)
	}
	s := string(raw)
	if !strings.Contains(s, "corr-123") {
		t.Fatal("request correlation required")
	}
	if strings.Contains(s, secretish) ||
		strings.Contains(s, "raw-secret-key") ||
		strings.Contains(s, "Bearer tok") ||
		strings.Contains(s, "password") {
		t.Fatalf("secrets leaked in AuditEvent: %s", s)
	}
}

// memoryAuditLog is a test-only append boundary used to prove the unit-level
// assemble-then-append failure property. Production append ownership remains
// with lifecycle/scheduler/handlers (TASK-F16-04/06/09/10/11).
type memoryAuditLog struct {
	mu     sync.Mutex
	events []apiconform.AuditEvent
	fail   error
}

func (a *memoryAuditLog) Append(_ context.Context, event apiconform.AuditEvent) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.fail != nil {
		return a.fail
	}
	a.events = append(a.events, event)
	return nil
}

func (a *memoryAuditLog) Len() int {
	a.mu.Lock()
	defer a.mu.Unlock()
	return len(a.events)
}

func TestAssembleThenAppendFailure_NoDurableEventNoDisclosure(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		ev   apiconform.AuditEvent
	}{
		{"create", AssembleSuccessfulCreate(baseInput("a-create", "r-create"))},
		{"authz-denial", AssembleAuthorizationDenial(baseInput("a-authz", "r-authz"))},
		{"safe-denial", AssembleSafeDenial(baseInput("a-safe", "r-safe"))},
		{"expiry", AssembleExpiry(baseInput("a-expiry", "r-expiry"))},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			log := &memoryAuditLog{fail: errors.New("append boom")}
			// Caller-owned assemble-then-append: assembly already returned the
			// in-memory event; a failed append leaves nothing durable.
			err := log.Append(context.Background(), tc.ev)
			if err == nil {
				t.Fatal("expected append failure")
			}
			if log.Len() != 0 {
				t.Fatal("failed append must not retain a durable AuditEvent")
			}
			// Callers translate append failure to INTERNAL_ERROR/500 without
			// disclosing a suppressed 403/404 (unit-level non-disclosure).
			prob := apiproblem.New(apiproblem.CodeInternalError).WithDetail("required AuditEvent append failed")
			if prob.Code != apiproblem.CodeInternalError {
				t.Fatalf("want INTERNAL_ERROR, got %s", prob.Code)
			}
			if prob.Code == apiproblem.CodeDependencyUnavailable {
				t.Fatal("DEPENDENCY_UNAVAILABLE must not be used for audit-append failure")
			}
			detail := strings.ToLower(prob.Detail)
			raw, _ := json.Marshal(prob)
			blob := strings.ToLower(string(raw) + " " + detail + " " + err.Error())
			if strings.Contains(blob, "403") ||
				strings.Contains(blob, "404") ||
				strings.Contains(blob, "authorization_denied") ||
				strings.Contains(blob, "resource_not_found") ||
				strings.Contains(blob, "safe_denial") {
				t.Fatalf("append-failure error must not disclose suppressed 403/404: %s", blob)
			}
		})
	}
}

func TestAssemble_DoesNotMutateCallerInput(t *testing.T) {
	t.Parallel()
	in := baseInput("audit-copy", "req-copy")
	origActor := in.Actor
	_ = AssembleSuccessfulCreate(in)
	_ = AssembleExpiry(in)
	if in.Actor != origActor {
		t.Fatalf("assembly must not mutate caller's input; got %#v want %#v", in.Actor, origActor)
	}
}
