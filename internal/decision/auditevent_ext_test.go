package decision

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
)

func sampleAuditLinkage() AuditLinkage {
	opRef := apimeta.TypedRef{
		APIVersion: "core.sovrunn.io/v1alpha1",
		Kind:       "Operation",
		Name:       "op-1",
		UID:        "op-uid-1",
	}
	obligationRef := apimeta.TypedRef{
		APIVersion: APIVersionDecisionRecord,
		Kind:       "Obligation",
		Name:       "audit-retain",
		UID:        "obligation-uid-1",
	}
	return AuditLinkage{
		DecisionRef: apimeta.TypedRef{
			APIVersion: APIVersionDecisionRecord,
			Kind:       KindDecisionRecord,
			Name:       "placement-decision-123",
			UID:        "decision-123",
		},
		AuditObligationRef: &obligationRef,
		Correlation: Correlation{
			AuditEventRefs: []apimeta.TypedRef{
				{
					APIVersion: APIVersionDecisionRecord,
					Kind:       "AuditEvent",
					Name:       "audit-1",
					UID:        "audit-uid-1",
				},
			},
			OperationRef: &opRef,
			TraceRef:     "trace-xyz",
		},
	}
}

func TestAuditLinkageJSONRoundTrip(t *testing.T) {
	t.Parallel()

	in := sampleAuditLinkage()
	data, err := json.Marshal(in)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var out AuditLinkage
	if err := json.Unmarshal(data, &out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if out.DecisionRef != in.DecisionRef {
		t.Fatalf("decisionRef = %+v, want %+v", out.DecisionRef, in.DecisionRef)
	}
	if out.AuditObligationRef == nil || *out.AuditObligationRef != *in.AuditObligationRef {
		t.Fatalf("auditObligationRef = %#v, want %#v", out.AuditObligationRef, in.AuditObligationRef)
	}
	if !reflect.DeepEqual(out.Correlation, in.Correlation) {
		t.Fatalf("correlation = %+v, want %+v", out.Correlation, in.Correlation)
	}
}

func TestAuditLinkageReusesCorrelationCarrier(t *testing.T) {
	t.Parallel()

	// Correlation on AuditLinkage must be the same package type used by
	// DecisionRecord (FEATURE-0012/0013 correlation carrier), not a parallel type.
	linkageField, ok := reflect.TypeOf(AuditLinkage{}).FieldByName("Correlation")
	if !ok {
		t.Fatal("AuditLinkage.Correlation field missing")
	}
	recordField, ok := reflect.TypeOf(DecisionBody{}).FieldByName("Correlation")
	if !ok {
		t.Fatal("DecisionBody.Correlation field missing")
	}
	if linkageField.Type != recordField.Type {
		t.Fatalf("AuditLinkage.Correlation type %v must equal DecisionBody.Correlation type %v",
			linkageField.Type, recordField.Type)
	}
	if linkageField.Type != reflect.TypeOf(Correlation{}) {
		t.Fatalf("AuditLinkage.Correlation must be Correlation, got %v", linkageField.Type)
	}

	// Round-trip proves the shared carrier fields (auditEventRefs, operationRef, traceRef).
	in := sampleAuditLinkage()
	data, err := json.Marshal(in.Correlation)
	if err != nil {
		t.Fatalf("marshal correlation: %v", err)
	}
	var corr Correlation
	if err := json.Unmarshal(data, &corr); err != nil {
		t.Fatalf("unmarshal correlation: %v", err)
	}
	if len(corr.AuditEventRefs) != 1 || corr.AuditEventRefs[0].UID != "audit-uid-1" {
		t.Fatalf("auditEventRefs = %+v", corr.AuditEventRefs)
	}
	if corr.OperationRef == nil || corr.OperationRef.UID != "op-uid-1" {
		t.Fatalf("operationRef = %#v", corr.OperationRef)
	}
	if corr.TraceRef != "trace-xyz" {
		t.Fatalf("traceRef = %q", corr.TraceRef)
	}
}

func TestAuditLinkageOmitemptyOptionalObligationRef(t *testing.T) {
	t.Parallel()

	in := sampleAuditLinkage()
	in.AuditObligationRef = nil

	data, err := json.Marshal(in)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	raw := string(data)
	if strings.Contains(raw, `"auditObligationRef"`) {
		t.Fatalf("optional auditObligationRef must be omitted when unset; json=%s", raw)
	}
	if !strings.Contains(raw, `"decisionRef"`) || !strings.Contains(raw, `"correlation"`) {
		t.Fatalf("required fields must be present; json=%s", raw)
	}

	var out AuditLinkage
	if err := json.Unmarshal(data, &out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if out.AuditObligationRef != nil {
		t.Fatalf("auditObligationRef = %#v, want nil", out.AuditObligationRef)
	}
}

func TestAuditLinkageDoesNotDuplicateBaseAuditEventFields(t *testing.T) {
	t.Parallel()

	// Extension type must carry only linkage fields; base envelope fields stay
	// on FEATURE-0012 AuditEvent (closure 7).
	typ := reflect.TypeOf(AuditLinkage{})
	forbidden := []string{
		"APIVersion", "Kind", "TypeMeta", "Metadata", "Record",
		"ActorRef", "RequestID", "SubjectRef", "SubjectResourceVersion",
		"Action", "Outcome", "ReasonCode", "CorrectionOfRef",
	}
	for _, name := range forbidden {
		if _, ok := typ.FieldByName(name); ok {
			t.Fatalf("AuditLinkage must not duplicate base AuditEvent field %q", name)
		}
	}

	required := []string{"DecisionRef", "AuditObligationRef", "Correlation"}
	for _, name := range required {
		if _, ok := typ.FieldByName(name); !ok {
			t.Fatalf("AuditLinkage missing required field %q", name)
		}
	}
	if typ.NumField() != 3 {
		t.Fatalf("AuditLinkage field count = %d, want 3 (no base-field duplication)", typ.NumField())
	}
}

func TestAuditOutcomeRemainsClosedSet(t *testing.T) {
	t.Parallel()

	// T-007 preserves the closed Succeeded/Denied/Failed set; no extension value.
	got := AllAuditOutcomes()
	want := []AuditOutcome{
		AuditOutcomeSucceeded,
		AuditOutcomeDenied,
		AuditOutcomeFailed,
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("AllAuditOutcomes() = %v, want %v", got, want)
	}
}
