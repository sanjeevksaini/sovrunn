package apicond

import (
	"testing"
	"time"
)

func TestConditionStatusValid(t *testing.T) {
	t.Parallel()

	for _, s := range []ConditionStatus{ConditionTrue, ConditionFalse, ConditionUnknown} {
		if !s.Valid() {
			t.Fatalf("%q must be valid", s)
		}
	}
	if ConditionStatus("true").Valid() {
		t.Fatal("lowercase true must be invalid")
	}
	if ConditionStatus("").Valid() {
		t.Fatal("empty status must be invalid")
	}
}

func TestIsPascalCase(t *testing.T) {
	t.Parallel()

	valid := []string{"Valid", "ValidationSucceeded", "Available", "A", "Ready42"}
	for _, s := range valid {
		if !IsPascalCase(s) {
			t.Fatalf("%q must be PascalCase", s)
		}
	}

	invalid := []string{
		"",
		"valid",
		"validationSucceeded",
		"Validation-Succeeded",
		"Validation_Succeeded",
		" Validation",
		"Valid ation",
		"42Ready",
	}
	for _, s := range invalid {
		if IsPascalCase(s) {
			t.Fatalf("%q must not be PascalCase", s)
		}
	}
}

func TestConditionValid(t *testing.T) {
	t.Parallel()

	ok := Condition{
		Type:   "Valid",
		Status: ConditionTrue,
		Reason: "ValidationSucceeded",
	}
	if !ok.Valid() {
		t.Fatal("well-formed condition must be Valid")
	}

	badType := ok
	badType.Type = "valid"
	if badType.Valid() {
		t.Fatal("non-PascalCase type must be invalid")
	}

	badReason := ok
	badReason.Reason = "not-pascal"
	if badReason.Valid() {
		t.Fatal("non-PascalCase reason must be invalid")
	}

	badStatus := ok
	badStatus.Status = ConditionStatus("Maybe")
	if badStatus.Valid() {
		t.Fatal("invalid status must be invalid")
	}

	// Message is informational; empty or free-form does not affect Valid.
	withMsg := ok
	withMsg.Message = "any human text; not a contract"
	if !withMsg.Valid() {
		t.Fatal("message must not affect Valid")
	}
}

func TestSetConditionTransitionTimeInvariant(t *testing.T) {
	t.Parallel()

	t0 := time.Date(2026, 7, 11, 0, 0, 0, 0, time.UTC)
	t1 := t0.Add(time.Hour)
	t2 := t1.Add(time.Hour)

	conds := SetCondition(nil, Condition{
		Type:               "Ready",
		Status:             ConditionFalse,
		Reason:             "Pending",
		ObservedGeneration: 1,
	}, t0)
	got, ok := GetCondition(conds, "Ready")
	if !ok {
		t.Fatal("condition must exist after first SetCondition")
	}
	firstTime := got.LastTransitionTime
	if firstTime != t0.Format(time.RFC3339) {
		t.Fatalf("first LastTransitionTime = %q, want %q", firstTime, t0.Format(time.RFC3339))
	}

	// Same status: LastTransitionTime must not advance; reason/message may update.
	conds = SetCondition(conds, Condition{
		Type:               "Ready",
		Status:             ConditionFalse,
		Reason:             "StillPending",
		Message:            "waiting",
		ObservedGeneration: 2,
	}, t1)
	got, ok = GetCondition(conds, "Ready")
	if !ok {
		t.Fatal("condition must still exist")
	}
	if got.LastTransitionTime != firstTime {
		t.Fatalf("LastTransitionTime advanced on same status: got %q want %q", got.LastTransitionTime, firstTime)
	}
	if got.Reason != "StillPending" || got.Message != "waiting" || got.ObservedGeneration != 2 {
		t.Fatalf("non-status fields must update on upsert: %#v", got)
	}

	// Status change: LastTransitionTime advances to now.
	conds = SetCondition(conds, Condition{
		Type:               "Ready",
		Status:             ConditionTrue,
		Reason:             "Succeeded",
		ObservedGeneration: 2,
	}, t2)
	got, ok = GetCondition(conds, "Ready")
	if !ok {
		t.Fatal("condition must still exist after status change")
	}
	want := t2.Format(time.RFC3339)
	if got.LastTransitionTime != want {
		t.Fatalf("LastTransitionTime after status change = %q, want %q", got.LastTransitionTime, want)
	}
}

func TestSetConditionConditionsNotHistory(t *testing.T) {
	t.Parallel()

	t0 := time.Date(2026, 7, 11, 12, 0, 0, 0, time.UTC)
	var conds []Condition

	// Repeated upserts of the same type must leave exactly one entry.
	for i, status := range []ConditionStatus{ConditionUnknown, ConditionFalse, ConditionTrue, ConditionTrue, ConditionFalse} {
		conds = SetCondition(conds, Condition{
			Type:               "Ready",
			Status:             status,
			Reason:             "Step",
			ObservedGeneration: int64(i + 1),
		}, t0.Add(time.Duration(i)*time.Minute))
	}
	if len(conds) != 1 {
		t.Fatalf("conditions must not accumulate history; len=%d want 1, got %#v", len(conds), conds)
	}
	if conds[0].Type != "Ready" || conds[0].Status != ConditionFalse {
		t.Fatalf("final Ready condition = %#v", conds[0])
	}

	// Distinct types coexist as current facts, not as an event log.
	conds = SetCondition(conds, Condition{
		Type:   "Valid",
		Status: ConditionTrue,
		Reason: "ValidationSucceeded",
	}, t0.Add(10*time.Minute))
	conds = SetCondition(conds, Condition{
		Type:   "Available",
		Status: ConditionFalse,
		Reason: "NotYet",
	}, t0.Add(11*time.Minute))
	if len(conds) != 3 {
		t.Fatalf("distinct types must each appear once; len=%d want 3, got %#v", len(conds), conds)
	}

	// Updating Available must not append another Available row.
	conds = SetCondition(conds, Condition{
		Type:   "Available",
		Status: ConditionTrue,
		Reason: "ReadyNow",
	}, t0.Add(12*time.Minute))
	if len(conds) != 3 {
		t.Fatalf("upsert must replace, not append; len=%d want 3", len(conds))
	}
	got, ok := GetCondition(conds, "Available")
	if !ok || got.Status != ConditionTrue || got.Reason != "ReadyNow" {
		t.Fatalf("Available upsert failed: ok=%v got=%#v", ok, got)
	}
}

func TestGetConditionMissing(t *testing.T) {
	t.Parallel()

	_, ok := GetCondition(nil, "Ready")
	if ok {
		t.Fatal("GetCondition on empty slice must miss")
	}
	conds := SetCondition(nil, Condition{
		Type:   "Ready",
		Status: ConditionTrue,
		Reason: "Ok",
	}, time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC))
	_, ok = GetCondition(conds, "Valid")
	if ok {
		t.Fatal("GetCondition must miss absent type")
	}
}

func TestSetConditionDoesNotMutateInput(t *testing.T) {
	t.Parallel()

	t0 := time.Date(2026, 7, 11, 0, 0, 0, 0, time.UTC)
	orig := []Condition{{
		Type:               "Ready",
		Status:             ConditionFalse,
		Reason:             "Pending",
		LastTransitionTime: t0.Format(time.RFC3339),
	}}
	origCopy := append([]Condition(nil), orig...)

	_ = SetCondition(orig, Condition{
		Type:   "Ready",
		Status: ConditionTrue,
		Reason: "Succeeded",
	}, t0.Add(time.Hour))

	if orig[0].Status != origCopy[0].Status || orig[0].Reason != origCopy[0].Reason {
		t.Fatalf("SetCondition must not mutate input slice elements: got %#v", orig[0])
	}
}
