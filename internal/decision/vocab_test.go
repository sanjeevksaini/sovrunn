package decision

import "testing"

func TestDecisionForms(t *testing.T) {
	t.Parallel()

	want := []DecisionForm{
		"ADJUDICATION",
		"SELECTION",
		"RANKING",
		"CLASSIFICATION",
		"RESOLUTION",
		"ALLOCATION",
		"PLAN",
		"ASSESSMENT",
		"RECOMMENDATION",
	}
	got := AllDecisionForms()
	if len(got) != len(want) {
		t.Fatalf("AllDecisionForms len=%d, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("AllDecisionForms[%d]=%q, want %q", i, got[i], want[i])
		}
		if !got[i].Valid() {
			t.Fatalf("DecisionForm %q should be Valid", got[i])
		}
	}
	if DecisionForm("UNKNOWN").Valid() {
		t.Fatal("unknown DecisionForm must not be Valid")
	}
}

func TestDecisionAuthorities(t *testing.T) {
	t.Parallel()

	want := []DecisionAuthority{
		"AUTHORITATIVE",
		"ADVISORY",
		"RECOMMENDATION",
		"SIMULATION",
	}
	got := AllDecisionAuthorities()
	if len(got) != len(want) {
		t.Fatalf("AllDecisionAuthorities len=%d, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("AllDecisionAuthorities[%d]=%q, want %q", i, got[i], want[i])
		}
		if !got[i].Valid() {
			t.Fatalf("DecisionAuthority %q should be Valid", got[i])
		}
	}
	if DecisionAuthority("ENFORCEMENT").Valid() {
		t.Fatal("unknown DecisionAuthority must not be Valid")
	}
}

func TestAdjudicationOutcomes(t *testing.T) {
	t.Parallel()

	want := []AdjudicationOutcome{
		"ALLOWED",
		"DENIED",
		"REQUIRES_APPROVAL",
	}
	got := AllAdjudicationOutcomes()
	if len(got) != len(want) {
		t.Fatalf("AllAdjudicationOutcomes len=%d, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("AllAdjudicationOutcomes[%d]=%q, want %q", i, got[i], want[i])
		}
		if !got[i].Valid() {
			t.Fatalf("AdjudicationOutcome %q should be Valid", got[i])
		}
	}
	if AdjudicationOutcome("PENDING").Valid() {
		t.Fatal("unknown AdjudicationOutcome must not be Valid")
	}
}

func TestAuditOutcomes(t *testing.T) {
	t.Parallel()

	want := []AuditOutcome{
		"Succeeded",
		"Denied",
		"Failed",
	}
	got := AllAuditOutcomes()
	if len(got) != len(want) {
		t.Fatalf("AllAuditOutcomes len=%d, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("AllAuditOutcomes[%d]=%q, want %q", i, got[i], want[i])
		}
		if !got[i].Valid() {
			t.Fatalf("AuditOutcome %q should be Valid", got[i])
		}
	}
	if AuditOutcome("SUCCESS").Valid() {
		t.Fatal("unknown AuditOutcome must not be Valid")
	}
	if AuditOutcome("SUCCEEDED").Valid() {
		t.Fatal("uppercase SUCCEEDED must not be Valid; canonical form is Succeeded")
	}
}

func TestRelationshipKinds(t *testing.T) {
	t.Parallel()

	want := []RelationshipKind{
		"CORRECTS",
		"SUPERSEDES",
		"REVOKES",
	}
	got := AllRelationshipKinds()
	if len(got) != len(want) {
		t.Fatalf("AllRelationshipKinds len=%d, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("AllRelationshipKinds[%d]=%q, want %q", i, got[i], want[i])
		}
		if !got[i].Valid() {
			t.Fatalf("RelationshipKind %q should be Valid", got[i])
		}
	}
	if RelationshipKind("AMENDS").Valid() {
		t.Fatal("unknown RelationshipKind must not be Valid")
	}
}

func TestSensitivities(t *testing.T) {
	t.Parallel()

	want := []Sensitivity{
		"PUBLIC",
		"INTERNAL",
		"CONFIDENTIAL",
		"RESTRICTED",
	}
	got := AllSensitivities()
	if len(got) != len(want) {
		t.Fatalf("AllSensitivities len=%d, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("AllSensitivities[%d]=%q, want %q", i, got[i], want[i])
		}
		if !got[i].Valid() {
			t.Fatalf("Sensitivity %q should be Valid", got[i])
		}
	}
	if Sensitivity("TOP_SECRET").Valid() {
		t.Fatal("unknown Sensitivity must not be Valid")
	}
}

func TestFinalities(t *testing.T) {
	t.Parallel()

	want := []Finality{"FINAL"}
	got := AllFinalities()
	if len(got) != len(want) {
		t.Fatalf("AllFinalities len=%d, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("AllFinalities[%d]=%q, want %q", i, got[i], want[i])
		}
		if !got[i].Valid() {
			t.Fatalf("Finality %q should be Valid", got[i])
		}
	}
	if Finality("PENDING").Valid() {
		t.Fatal("PENDING Finality must not be Valid")
	}
	if Finality("NON_FINAL").Valid() {
		t.Fatal("NON_FINAL Finality must not be Valid")
	}
}

func TestEffectiveStates(t *testing.T) {
	t.Parallel()

	want := []EffectiveState{
		"EFFECTIVE",
		"SUPERSEDED",
		"REVOKED",
	}
	got := AllEffectiveStates()
	if len(got) != len(want) {
		t.Fatalf("AllEffectiveStates len=%d, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("AllEffectiveStates[%d]=%q, want %q", i, got[i], want[i])
		}
		if !got[i].Valid() {
			t.Fatalf("EffectiveState %q should be Valid", got[i])
		}
	}
	if EffectiveState("PENDING").Valid() {
		t.Fatal("unknown EffectiveState must not be Valid")
	}
}
