package decision

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
)

func sampleDecisionRecord() DecisionRecord {
	return DecisionRecord{
		TypeMeta: apimeta.TypeMeta{
			APIVersion: APIVersionDecisionRecord,
			Kind:       KindDecisionRecord,
		},
		Metadata: apimeta.ObjectMeta{
			Name: "placement-decision-123",
			UID:  "decision-123",
			ScopeRef: &apimeta.ScopeRef{
				TypedRef: apimeta.TypedRef{
					APIVersion: "core.sovrunn.io/v1alpha1",
					Kind:       "Project",
					Name:       "payments-production",
					UID:        "project-123",
				},
			},
		},
		Record: DecisionBody{
			ProfileRef: ProfileRef{Name: "PlacementDecision", Version: "1.0.0"},
			Form:       DecisionFormSelection,
			Authority:  DecisionAuthorityAuthoritative,
			Purpose:    "workload-placement",
			SubjectRefs: []apimeta.TypedRef{
				{
					APIVersion: "platform.sovrunn.io/v1alpha1",
					Kind:       "ServiceInstance",
					Name:       "postgres-primary",
					UID:        "service-456",
				},
			},
			Result: DecisionResultBody{
				TypedResult: json.RawMessage(`{"selected":"pool-a"}`),
				Rationale: DecisionRationale{
					ReasonCodes: []string{"CAPACITY_FIT"},
					Reasons:     []string{"selected pool meets capacity"},
				},
				Obligations: []Obligation{
					{ID: "audit-retain", Mandatory: true},
				},
			},
			RetryKey: "retry-abc",
			SemanticIdentity: SemanticDecisionIdentity{
				Basis:          "placement-v1",
				InputRefs:      []string{"input-1"},
				CanonicalScope: "Project:project-123",
			},
			Correlation: Correlation{
				AuditEventRefs: []apimeta.TypedRef{
					{
						APIVersion: APIVersionDecisionRecord,
						Kind:       "AuditEvent",
						Name:       "audit-1",
						UID:        "audit-uid-1",
					},
				},
				TraceRef: "trace-xyz",
			},
			Finality: FinalityFinal,
		},
	}
}

func TestDecisionRecordJSONRoundTrip(t *testing.T) {
	t.Parallel()

	in := sampleDecisionRecord()
	data, err := json.Marshal(in)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var out DecisionRecord
	if err := json.Unmarshal(data, &out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if out.APIVersion != APIVersionDecisionRecord || out.Kind != KindDecisionRecord {
		t.Fatalf("TypeMeta = {%q,%q}, want {%q,%q}", out.APIVersion, out.Kind, APIVersionDecisionRecord, KindDecisionRecord)
	}
	if out.Metadata.Name != "placement-decision-123" || out.Metadata.UID != "decision-123" {
		t.Fatalf("metadata identity mismatch: %+v", out.Metadata)
	}
	if out.Metadata.ScopeRef == nil || out.Metadata.ScopeRef.UID != "project-123" {
		t.Fatalf("metadata.scopeRef must remain sole scope authority, got %#v", out.Metadata.ScopeRef)
	}
	if out.Record.ProfileRef != (ProfileRef{Name: "PlacementDecision", Version: "1.0.0"}) {
		t.Fatalf("profileRef = %+v", out.Record.ProfileRef)
	}
	if out.Record.Form != DecisionFormSelection || out.Record.Authority != DecisionAuthorityAuthoritative {
		t.Fatalf("form/authority = %q/%q", out.Record.Form, out.Record.Authority)
	}
	if out.Record.Finality != FinalityFinal {
		t.Fatalf("finality = %q, want FINAL", out.Record.Finality)
	}
	if len(out.Record.Correlation.AuditEventRefs) != 1 {
		t.Fatalf("auditEventRefs len=%d, want 1", len(out.Record.Correlation.AuditEventRefs))
	}
	if string(out.Record.Result.TypedResult) != `{"selected":"pool-a"}` {
		t.Fatalf("typedResult = %s", out.Record.Result.TypedResult)
	}
	if len(out.Record.Result.Obligations) != 1 || !out.Record.Result.Obligations[0].Mandatory {
		t.Fatalf("obligations = %+v", out.Record.Result.Obligations)
	}
	if out.Record.SemanticIdentity.CanonicalScope != "Project:project-123" {
		t.Fatalf("semanticIdentity.canonicalScope = %q", out.Record.SemanticIdentity.CanonicalScope)
	}
}

func TestDecisionRecordOmitemptyOptionalRefs(t *testing.T) {
	t.Parallel()

	in := sampleDecisionRecord()
	// Ensure optional singular refs/objects stay unset.
	in.Record.Adjudication = nil
	in.Record.Composition = nil
	in.Record.Relationship = nil
	in.Record.Sovereignty = nil
	in.Record.Trust = nil
	in.Record.Correlation.OperationRef = nil
	in.Record.RetryKey = ""
	in.Record.Result.Rationale = DecisionRationale{}
	in.Record.Result.Obligations = []Obligation{}
	in.Record.SubjectRefs = nil

	data, err := json.Marshal(in)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	raw := string(data)

	for _, key := range []string{
		`"adjudication"`,
		`"composition"`,
		`"relationship"`,
		`"sovereignty"`,
		`"trust"`,
		`"operationRef"`,
		`"retryKey"`,
		`"subjectRefs"`,
		`"evaluationResults"`,
		`"reasonCodes"`,
		`"reasons"`,
		`"alternatives"`,
		`"warnings"`,
		`"correctiveSuggestions"`,
	} {
		if strings.Contains(raw, key) {
			t.Fatalf("optional field %s must be omitted when unset; json=%s", key, raw)
		}
	}

	// Required correlation.auditEventRefs must remain present.
	if !strings.Contains(raw, `"auditEventRefs"`) {
		t.Fatalf("required auditEventRefs missing; json=%s", raw)
	}
	if !strings.Contains(raw, `"semanticIdentity"`) {
		t.Fatalf("required semanticIdentity missing; json=%s", raw)
	}
	if !strings.Contains(raw, `"finality"`) {
		t.Fatalf("required finality missing; json=%s", raw)
	}
	if !strings.Contains(raw, `"obligations"`) {
		t.Fatalf("required-present obligations missing; json=%s", raw)
	}
}

func TestDecisionRecordRequiredFieldsPresent(t *testing.T) {
	t.Parallel()

	in := sampleDecisionRecord()
	data, err := json.Marshal(in)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var top map[string]json.RawMessage
	if err := json.Unmarshal(data, &top); err != nil {
		t.Fatalf("unmarshal top: %v", err)
	}
	for _, key := range []string{"apiVersion", "kind", "metadata", "record"} {
		if _, ok := top[key]; !ok {
			t.Fatalf("missing top-level required field %q", key)
		}
	}

	var record map[string]json.RawMessage
	if err := json.Unmarshal(top["record"], &record); err != nil {
		t.Fatalf("unmarshal record: %v", err)
	}
	for _, key := range []string{
		"profileRef",
		"form",
		"authority",
		"purpose",
		"result",
		"semanticIdentity",
		"correlation",
		"finality",
	} {
		if _, ok := record[key]; !ok {
			t.Fatalf("missing record required field %q", key)
		}
	}

	var corr map[string]json.RawMessage
	if err := json.Unmarshal(record["correlation"], &corr); err != nil {
		t.Fatalf("unmarshal correlation: %v", err)
	}
	refsRaw, ok := corr["auditEventRefs"]
	if !ok {
		t.Fatal("correlation.auditEventRefs must be present")
	}
	var refs []apimeta.TypedRef
	if err := json.Unmarshal(refsRaw, &refs); err != nil {
		t.Fatalf("unmarshal auditEventRefs: %v", err)
	}
	if len(refs) < 1 {
		t.Fatal("correlation.auditEventRefs must be non-empty (minItems: 1)")
	}

	var result map[string]json.RawMessage
	if err := json.Unmarshal(record["result"], &result); err != nil {
		t.Fatalf("unmarshal result: %v", err)
	}
	for _, key := range []string{"typedResult", "rationale", "obligations"} {
		if _, ok := result[key]; !ok {
			t.Fatalf("missing result required field %q", key)
		}
	}
}

func TestDecisionRecordOptionalRefsRoundTrip(t *testing.T) {
	t.Parallel()

	outcome := AdjudicationOutcomeAllowed
	in := sampleDecisionRecord()
	in.Record.Adjudication = &outcome
	in.Record.Composition = &CompositionRef{
		StrategyName:    "first-applicable",
		StrategyVersion: "1.0.0",
		GraphRef:        "placement-graph",
		GraphVersion:    "1.0.0",
		InputRefs:       []string{"eval-1"},
	}
	in.Record.Relationship = &DecisionRelationship{
		Kind: RelationshipKindSupersedes,
		TargetRef: apimeta.TypedRef{
			APIVersion: APIVersionDecisionRecord,
			Kind:       KindDecisionRecord,
			Name:       "prior-decision",
			UID:        "decision-000",
		},
	}
	in.Record.Sovereignty = &SovereigntyCarrier{
		Data:        "residency:in-jurisdiction",
		Legal:       "policy:export-restricted",
		AIModel:     "ai:disabled",
		Operational: "break-glass:local",
	}
	opRef := apimeta.TypedRef{
		APIVersion: "core.sovrunn.io/v1alpha1",
		Kind:       "Operation",
		Name:       "op-1",
		UID:        "op-uid-1",
	}
	in.Record.Correlation.OperationRef = &opRef
	in.Record.Result.Obligations[0].Payload = json.RawMessage(`{"retainDays":365}`)

	data, err := json.Marshal(in)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var out DecisionRecord
	if err := json.Unmarshal(data, &out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if out.Record.Adjudication == nil || *out.Record.Adjudication != AdjudicationOutcomeAllowed {
		t.Fatalf("adjudication = %#v", out.Record.Adjudication)
	}
	if out.Record.Composition == nil || out.Record.Composition.StrategyName != "first-applicable" {
		t.Fatalf("composition = %#v", out.Record.Composition)
	}
	if out.Record.Relationship == nil || out.Record.Relationship.Kind != RelationshipKindSupersedes {
		t.Fatalf("relationship = %#v", out.Record.Relationship)
	}
	if out.Record.Sovereignty == nil || out.Record.Sovereignty.Data != "residency:in-jurisdiction" {
		t.Fatalf("sovereignty = %#v", out.Record.Sovereignty)
	}
	if out.Record.Correlation.OperationRef == nil || out.Record.Correlation.OperationRef.UID != "op-uid-1" {
		t.Fatalf("operationRef = %#v", out.Record.Correlation.OperationRef)
	}
	if string(out.Record.Result.Obligations[0].Payload) != `{"retainDays":365}` {
		t.Fatalf("obligation payload = %s", out.Record.Result.Obligations[0].Payload)
	}
}

func TestDecisionRecordNoSpecStatus(t *testing.T) {
	t.Parallel()

	data, err := json.Marshal(sampleDecisionRecord())
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	raw := string(data)
	if strings.Contains(raw, `"spec"`) || strings.Contains(raw, `"status"`) {
		t.Fatalf("DecisionRecord must not carry mutable spec/status; json=%s", raw)
	}
}
