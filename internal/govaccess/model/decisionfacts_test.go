package model_test

import (
	"bytes"
	"testing"
	"time"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/model"
)

func TestApprovalDecisionFactsDeepCopyIsolation(t *testing.T) {
	t.Parallel()
	subj := []byte("subject")
	prop := []byte("proposal")
	principals := []model.PrincipalRef{{
		Issuer: "https://idp.example", Subject: "a", PrincipalType: model.PrincipalTypeHuman,
	}}
	f, ok := model.NewApprovalDecisionFacts(
		"pol-v1", subj, prop, "stage-1", principals, 2,
		model.ApprovalTerminalApproved, time.Date(2026, 9, 4, 12, 0, 0, 0, time.UTC),
	)
	if !ok || !f.Sealed() {
		t.Fatal("expected sealed facts")
	}
	subj[0] = 'X'
	prop[0] = 'Y'
	principals[0].Subject = "mutated"
	if bytes.Equal(f.SubjectDigest(), subj) || bytes.Equal(f.ProposalDigest(), prop) {
		t.Fatal("construction must deep-copy digests")
	}
	got := f.EligiblePrincipals()
	got[0].Subject = "mutated-again"
	if f.EligiblePrincipals()[0].Subject != "a" {
		t.Fatal("accessor must deep-copy principals")
	}
}

func TestExceptionDecisionFactsDeepCopyIsolation(t *testing.T) {
	t.Parallel()
	controls := []string{"c1"}
	evidence := []byte("approval-ev")
	f, ok := model.NewExceptionDecisionFacts(
		"ctrl-v1", "subj-1",
		apimeta.ScopeRef{TypedRef: apimeta.TypedRef{APIVersion: "core.sovrunn.io/v1alpha1", Kind: "Project", Name: "p1"}},
		time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC),
		"override-x", controls, evidence,
		model.ExceptionTerminalGrant, "approved", "proposal-1", "grant-1",
	)
	if !ok {
		t.Fatal("expected sealed facts")
	}
	controls[0] = "mutated"
	evidence[0] = 'Z'
	gotControls := f.CompensatingControls()
	gotControls[0] = "mutated-again"
	if f.CompensatingControls()[0] != "c1" {
		t.Fatal("accessor must deep-copy compensating controls")
	}
	if bytes.Equal(f.ApprovalEvidence(), evidence) {
		t.Fatal("construction must deep-copy approval evidence")
	}
}

func TestMutationEventFactsDeepCopyIsolation(t *testing.T) {
	t.Parallel()
	digest := []byte("after")
	f, ok := model.NewMutationEventFacts(
		"roleassignment.created",
		apimeta.TypedRef{APIVersion: "gov.sovrunn.io/v1alpha1", Kind: "RoleAssignment", Name: "ra-1"},
		"3", "create", digest,
	)
	if !ok {
		t.Fatal("expected sealed facts")
	}
	digest[0] = 'X'
	got := f.AfterStateDigest()
	got[0] = 'Y'
	if f.AfterStateDigest()[0] != 'a' {
		t.Fatal("accessor must deep-copy after-state digest")
	}
}

func TestExceptionDecisionFactsDenyRejectsGrantUID(t *testing.T) {
	t.Parallel()
	_, ok := model.NewExceptionDecisionFacts(
		"ctrl-v1", "subj-1",
		apimeta.ScopeRef{TypedRef: apimeta.TypedRef{APIVersion: "core.sovrunn.io/v1alpha1", Kind: "Project", Name: "p1"}},
		time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC),
		"override-x", nil, nil,
		model.ExceptionTerminalDeny, "denied", "proposal-1", "grant-1",
	)
	if ok {
		t.Fatal("Deny must reject grant UID")
	}
}
