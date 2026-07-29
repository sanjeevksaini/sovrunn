package validate

import (
	"encoding/json"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/apivalid"
	"github.com/sanjeevksaini/sovrunn/internal/decision"
	"github.com/sanjeevksaini/sovrunn/internal/decision/bundle"
)

func orchestrationProfile() decision.DecisionProfile {
	p := validProfile()
	// Ensure obligations used by sample records are registered as supported
	// via the bundle KnownObligations default path when opts omit the map.
	return p
}

func mustOrchestrationView(t *testing.T, profiles ...decision.DecisionProfile) bundle.BundleView {
	t.Helper()
	if len(profiles) == 0 {
		profiles = []decision.DecisionProfile{orchestrationProfile()}
	}
	doc := bundle.Bundle{
		APIVersion: bundle.APIVersionDecisionProfileBundle,
		Kind:       bundle.KindDecisionProfileBundle,
		Spec: bundle.BundleSpec{
			Profiles: profiles,
			Obligations: []bundle.VersionedEntry{
				{ID: "audit-retain", Version: "1.0.0"},
				{ID: "notify-owner", Version: "1.0.0"},
			},
		},
	}
	raw, err := json.Marshal(doc)
	if err != nil {
		t.Fatalf("marshal bundle: %v", err)
	}
	loaded, prob := bundle.Load(raw, apivalid.ModeReadRepresentation)
	if prob != nil {
		t.Fatalf("bundle.Load = %#v", prob)
	}
	return loaded.View()
}

func sampleOrchestrationRecord() decision.DecisionRecord {
	return decision.DecisionRecord{
		TypeMeta: apimeta.TypeMeta{
			APIVersion: decision.APIVersionDecisionRecord,
			Kind:       decision.KindDecisionRecord,
		},
		Metadata: apimeta.ObjectMeta{
			Name: "placement-decision-123",
			UID:  "decision-123",
			ScopeRef: &apimeta.ScopeRef{
				TypedRef: apimeta.TypedRef{
					APIVersion: "core.sovrunn.io/v1alpha1",
					Kind:       string(apimeta.ScopeProject),
					Name:       "payments-production",
					UID:        "project-123",
				},
			},
		},
		Record: decision.DecisionBody{
			ProfileRef: decision.ProfileRef{Name: "PlacementDecision", Version: "1.0.0"},
			Form:       decision.DecisionFormSelection,
			Authority:  decision.DecisionAuthorityAuthoritative,
			Purpose:    "workload-placement",
			SubjectRefs: []apimeta.TypedRef{
				{
					APIVersion: "platform.sovrunn.io/v1alpha1",
					Kind:       KindServiceInstance,
					Name:       "postgres-primary",
					UID:        "service-456",
				},
			},
			Result: decision.DecisionResultBody{
				TypedResult: json.RawMessage(`{"selected":"pool-a"}`),
				Rationale: decision.DecisionRationale{
					ReasonCodes: []string{"CAPACITY_FIT"},
				},
				Obligations: []decision.Obligation{
					{ID: "audit-retain", Mandatory: true},
				},
			},
			SemanticIdentity: decision.SemanticDecisionIdentity{
				Basis:          "placement-v1",
				InputRefs:      []string{"input-1"},
				CanonicalScope: "Project:project-123",
			},
			Correlation: decision.Correlation{
				AuditEventRefs: []apimeta.TypedRef{
					{
						APIVersion: "governance.sovrunn.io/v1alpha1",
						Kind:       "AuditEvent",
						Name:       "audit-1",
						UID:        "audit-uid-1",
					},
				},
			},
			Finality: decision.FinalityFinal,
		},
	}
}

func sampleAuditEnvelope() AuditEventEnvelope {
	return AuditEventEnvelope{
		TypeMeta: apimeta.TypeMeta{
			APIVersion: "governance.sovrunn.io/v1alpha1",
			Kind:       "AuditEvent",
		},
		Metadata: apimeta.ObjectMeta{
			Name: "audit-1",
			UID:  "audit-uid-1",
			ScopeRef: &apimeta.ScopeRef{
				TypedRef: apimeta.TypedRef{
					APIVersion: "core.sovrunn.io/v1alpha1",
					Kind:       string(apimeta.ScopeProject),
					Name:       "payments-production",
					UID:        "project-123",
				},
			},
		},
		Record: AuditEventRecordFields{
			ActorRef: apimeta.TypedRef{
				APIVersion: "core.sovrunn.io/v1alpha1",
				Kind:       "User",
				Name:       "operator",
				UID:        "user-1",
			},
			RequestID: "req-1",
			SubjectRef: apimeta.TypedRef{
				APIVersion: "platform.sovrunn.io/v1alpha1",
				Kind:       KindServiceInstance,
				Name:       "postgres-primary",
				UID:        "service-456",
			},
			Action:     "decision.record",
			Outcome:    decision.AuditOutcomeSucceeded,
			ReasonCode: "RECORDED",
		},
	}
}

func sampleAuditLinkage() decision.AuditLinkage {
	return decision.AuditLinkage{
		DecisionRef: apimeta.TypedRef{
			APIVersion: decision.APIVersionDecisionRecord,
			Kind:       decision.KindDecisionRecord,
			Name:       "placement-decision-123",
			UID:        "decision-123",
		},
		Correlation: decision.Correlation{
			AuditEventRefs: []apimeta.TypedRef{
				{
					APIVersion: "governance.sovrunn.io/v1alpha1",
					Kind:       "AuditEvent",
					Name:       "audit-1",
					UID:        "audit-uid-1",
				},
			},
		},
	}
}

func assertOrchestrationViolation(t *testing.T, p *apiproblem.Problem, wantCode apiproblem.ViolationCode, wantField string) {
	t.Helper()
	if p == nil {
		t.Fatalf("expected Problem with %s at %s, got nil", wantCode, wantField)
	}
	if p.Code != apiproblem.CodeValidationFailed {
		t.Fatalf("top-level code = %q, want %q", p.Code, apiproblem.CodeValidationFailed)
	}
	if len(p.Violations) != 1 {
		t.Fatalf("violations len=%d, want 1: %#v", len(p.Violations), p.Violations)
	}
	if p.Violations[0].Code != wantCode {
		t.Fatalf("violations[0].code = %q, want %q", p.Violations[0].Code, wantCode)
	}
	if p.Violations[0].Field != wantField {
		t.Fatalf("violations[0].field = %q, want %q", p.Violations[0].Field, wantField)
	}
}

func assertNotDecisionViolation(t *testing.T, p *apiproblem.Problem) {
	t.Helper()
	if p == nil {
		t.Fatal("expected Problem, got nil")
	}
	if Valid(apiproblem.ViolationCode(p.Code)) {
		t.Fatalf("top-level code %q must not be a DECISION_* code", p.Code)
	}
	for _, v := range p.Violations {
		if Valid(v.Code) {
			t.Fatalf("decode/structural violations[].code %q must not be DECISION_*", v.Code)
		}
	}
}

func TestValidateDecisionRecord_HappyPath(t *testing.T) {
	t.Parallel()

	view := mustOrchestrationView(t)
	rec := sampleOrchestrationRecord()
	if p := ValidateDecisionRecord(rec, view); p != nil {
		t.Fatalf("ValidateDecisionRecord = %#v", p)
	}
}

func TestDecodeAndValidateDecisionRecord_HappyPath(t *testing.T) {
	t.Parallel()

	view := mustOrchestrationView(t)
	rec := sampleOrchestrationRecord()
	raw, err := json.Marshal(rec)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	out, prob := DecodeAndValidateDecisionRecord(raw, apivalid.ModeReadRepresentation, view)
	if prob != nil {
		t.Fatalf("DecodeAndValidateDecisionRecord = %#v", prob)
	}
	if out.Metadata.UID != "decision-123" {
		t.Fatalf("uid = %q", out.Metadata.UID)
	}
	if out.Record.Finality != decision.FinalityFinal {
		t.Fatalf("finality = %q", out.Record.Finality)
	}
}

func TestDecodeDecisionRecord_DuplicateFieldUsesInheritedTopLevelCode(t *testing.T) {
	t.Parallel()

	// Duplicate JSON key → FEATURE-0012 DUPLICATE_FIELD (pass 1), never DECISION_*.
	raw := []byte(`{"apiVersion":"governance.sovrunn.io/v1alpha1","apiVersion":"governance.sovrunn.io/v1alpha1","kind":"DecisionRecord"}`)
	_, prob := DecodeDecisionRecord(raw, apivalid.ModeReadRepresentation)
	if prob == nil {
		t.Fatal("expected decode Problem")
	}
	if prob.Code != apiproblem.CodeDuplicateField {
		t.Fatalf("top-level code = %q, want %q", prob.Code, apiproblem.CodeDuplicateField)
	}
	assertNotDecisionViolation(t, prob)
}

func TestDecodeAndValidateDecisionRecord_DecodeFailureShortCircuits(t *testing.T) {
	t.Parallel()

	view := mustOrchestrationView(t)
	raw := []byte(`{"apiVersion":"governance.sovrunn.io/v1alpha1","kind":"DecisionRecord","unknownField":true}`)
	_, prob := DecodeAndValidateDecisionRecord(raw, apivalid.ModeReadRepresentation, view)
	if prob == nil {
		t.Fatal("expected decode Problem")
	}
	if prob.Code != apiproblem.CodeUnknownField {
		t.Fatalf("top-level code = %q, want %q", prob.Code, apiproblem.CodeUnknownField)
	}
	assertNotDecisionViolation(t, prob)
}

func TestValidateDecisionRecord_StructuralBeforeScope(t *testing.T) {
	t.Parallel()

	view := mustOrchestrationView(t)
	rec := sampleOrchestrationRecord()
	rec.Record.Correlation.AuditEventRefs = nil // F13-AUDIT-001 structural
	rec.Metadata.ScopeRef = nil                 // would also fail scope if reached

	p := ValidateDecisionRecord(rec, view)
	if p == nil {
		t.Fatal("expected structural Problem")
	}
	if p.Code != apiproblem.CodeMalformedRequest {
		t.Fatalf("top-level code = %q, want %q (structural before scope)", p.Code, apiproblem.CodeMalformedRequest)
	}
	assertNotDecisionViolation(t, p)
}

func TestValidateDecisionRecord_ScopeBeforeProfile(t *testing.T) {
	t.Parallel()

	view := mustOrchestrationView(t)
	rec := sampleOrchestrationRecord()
	// Invalid ScopeKind → DECISION_SCOPE_INVALID at pass 2; profile is resolvable
	// so pass 3 would succeed if reached.
	rec.Metadata.ScopeRef = &apimeta.ScopeRef{
		TypedRef: apimeta.TypedRef{
			APIVersion: "core.sovrunn.io/v1alpha1",
			Kind:       KindServiceInstance, // never a ScopeKind
			Name:       "svc",
			UID:        "svc-1",
		},
	}
	rec.Record.SemanticIdentity.CanonicalScope = "ServiceInstance:svc-1"

	p := ValidateDecisionRecord(rec, view)
	assertOrchestrationViolation(t, p, CodeScopeInvalid, ptrMetadataScopeRefKind)
}

func TestValidateDecisionRecord_ProfileBeforeEvaluation(t *testing.T) {
	t.Parallel()

	view := mustOrchestrationView(t)
	rec := sampleOrchestrationRecord()
	rec.Record.ProfileRef = decision.ProfileRef{Name: "MissingProfile", Version: "9.9.9"}
	// Evaluation would fail if reached (unsupported type), but profile unknown
	// must short-circuit first.
	rec.Record.EvaluationResults = []decision.EvaluationResult{{
		Evaluator:    decision.EvaluatorIdentity{Type: "not-accepted", Version: "1.0.0"},
		ResultStatus: decision.EvaluationResultStatusSuccess,
		Result:       json.RawMessage(`{"ok":true}`),
		Timing: decision.TimingEnvelope{
			StartedAt:   "2026-01-01T00:00:00Z",
			CompletedAt: "2026-01-01T00:00:01Z",
		},
		EvaluatedAt: "2026-01-01T00:00:01Z",
		InputIntegrity: &decision.OpaqueIntegrityCarrier{
			State: "captured",
		},
	}}

	p := ValidateDecisionRecord(rec, view)
	assertOrchestrationViolation(t, p, CodeProfileUnknown, ptrProfileRefName)
}

func TestValidateDecisionRecord_EvaluationPointerPrefixed(t *testing.T) {
	t.Parallel()

	view := mustOrchestrationView(t)
	rec := sampleOrchestrationRecord()
	rec.Record.EvaluationResults = []decision.EvaluationResult{{
		Evaluator:    decision.EvaluatorIdentity{Type: "not-accepted", Version: "1.0.0"},
		ResultStatus: decision.EvaluationResultStatusSuccess,
		Result:       json.RawMessage(`{"ok":true}`),
		Timing: decision.TimingEnvelope{
			StartedAt:   "2026-01-01T00:00:00Z",
			CompletedAt: "2026-01-01T00:00:01Z",
		},
		EvaluatedAt: "2026-01-01T00:00:01Z",
		InputIntegrity: &decision.OpaqueIntegrityCarrier{
			State: "captured",
		},
	}}

	p := ValidateDecisionRecord(rec, view)
	assertOrchestrationViolation(t, p, CodeEvaluationTypeUnsupported, "/record/evaluationResults/0/evaluator/type")
}

func TestValidateDecisionRecord_ObligationPointerPrefixed(t *testing.T) {
	t.Parallel()

	view := mustOrchestrationView(t)
	rec := sampleOrchestrationRecord()
	rec.Record.Result.Obligations = []decision.Obligation{
		{ID: "unknown-obligation", Mandatory: true},
	}

	p := ValidateDecisionRecord(rec, view)
	assertOrchestrationViolation(t, p, CodeObligationUnknown, "/record/result/obligations/0/id")
}

func TestValidateDecisionRecord_TrustPresentEmptyRejected(t *testing.T) {
	t.Parallel()

	view := mustOrchestrationView(t)
	rec := sampleOrchestrationRecord()
	rec.Record.Trust = &decision.TrustCarrier{} // RID-08 present-but-empty

	p := ValidateDecisionRecord(rec, view)
	assertOrchestrationViolation(t, p, CodeTrustUnknown, "/record/trust")
}

func TestValidateDecisionRecord_RelationshipPass(t *testing.T) {
	t.Parallel()

	view := mustOrchestrationView(t)
	rec := sampleOrchestrationRecord()
	rec.Record.Relationship = &decision.DecisionRelationship{
		Kind: decision.RelationshipKind("AMENDS"),
		TargetRef: apimeta.TypedRef{
			APIVersion: decision.APIVersionDecisionRecord,
			Kind:       decision.KindDecisionRecord,
			Name:       "pred",
			UID:        "uid-pred",
		},
	}

	p := ValidateDecisionRecord(rec, view)
	assertOrchestrationViolation(t, p, CodeRelationshipKindInvalid, ptrRelationshipKind)
}

func TestValidateDecisionRecord_SensitivityCeiling(t *testing.T) {
	t.Parallel()

	view := mustOrchestrationView(t)
	rec := sampleOrchestrationRecord()
	// Declared sensitivity above profile ceiling (Internal).
	p := ValidateDecisionRecordWith(rec, view, DecisionRecordOptions{
		DeclaredSensitivity: decision.SensitivityRestricted,
	})
	assertOrchestrationViolation(t, p, CodeEvaluationLimitExceeded, ptrDeclaredSensitivity)
}

func TestValidateAuditEvent_HappyPath(t *testing.T) {
	t.Parallel()

	view := mustOrchestrationView(t)
	ev := sampleAuditEnvelope()
	ext := sampleAuditLinkage()
	if p := ValidateAuditEvent(ev, ext, view); p != nil {
		t.Fatalf("ValidateAuditEvent = %#v", p)
	}
}

func TestValidateAuditEvent_ScopeInvalid(t *testing.T) {
	t.Parallel()

	view := mustOrchestrationView(t)
	ev := sampleAuditEnvelope()
	ev.Metadata.ScopeRef = &apimeta.ScopeRef{
		TypedRef: apimeta.TypedRef{
			APIVersion: "core.sovrunn.io/v1alpha1",
			Kind:       KindServiceInstance,
			Name:       "svc",
			UID:        "svc-1",
		},
	}
	p := ValidateAuditEvent(ev, sampleAuditLinkage(), view)
	assertOrchestrationViolation(t, p, CodeScopeInvalid, ptrMetadataScopeRefKind)
}

func TestValidateAuditEvent_StructuralLinkageUID(t *testing.T) {
	t.Parallel()

	view := mustOrchestrationView(t)
	ext := sampleAuditLinkage()
	ext.DecisionRef.UID = ""
	p := ValidateAuditEvent(sampleAuditEnvelope(), ext, view)
	if p == nil {
		t.Fatal("expected structural Problem")
	}
	if p.Code != apiproblem.CodeMalformedRequest {
		t.Fatalf("top-level code = %q, want %q", p.Code, apiproblem.CodeMalformedRequest)
	}
	assertNotDecisionViolation(t, p)
}

func TestDecodeAuditEventEnvelope_MalformedUsesInheritedCode(t *testing.T) {
	t.Parallel()

	_, prob := DecodeAuditEventEnvelope([]byte(`{`), apivalid.ModeReadRepresentation)
	if prob == nil {
		t.Fatal("expected decode Problem")
	}
	if prob.Code != apiproblem.CodeMalformedRequest {
		t.Fatalf("top-level code = %q, want %q", prob.Code, apiproblem.CodeMalformedRequest)
	}
	assertNotDecisionViolation(t, prob)
}

func TestValidateRelationship_StillExportedPass(t *testing.T) {
	t.Parallel()

	// T-026 orchestration reuses the T-024 ValidateRelationship pass; ensure
	// the exported pass remains callable and is the relationship step.
	root := rootNode("root", "uid-root")
	p := ValidateRelationship(RelationshipInput{
		Relationship: &decision.DecisionRelationship{
			Kind:      decision.RelationshipKindCorrects,
			TargetRef: root.Ref,
		},
		Self:        decisionRef("new", "uid-new"),
		Established: []RelationshipNode{root},
	})
	if p != nil {
		t.Fatalf("ValidateRelationship = %#v", p)
	}
}

func TestPassOrderConstants(t *testing.T) {
	t.Parallel()

	if PassDecodeStructural != 1 || PassSensitivity != 9 {
		t.Fatalf("nine-pass constants out of order: decode=%d sensitivity=%d", PassDecodeStructural, PassSensitivity)
	}
	order := []int{
		PassDecodeStructural, PassScope, PassProfile, PassEvaluation,
		PassComposition, PassTrust, PassObligation, PassRelationship, PassSensitivity,
	}
	for i, n := range order {
		if n != i+1 {
			t.Fatalf("pass order[%d]=%d, want %d", i, n, i+1)
		}
	}
}
