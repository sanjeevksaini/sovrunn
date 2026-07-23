package apivalid

import (
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
)

func TestPolicyForMatrixC2(t *testing.T) {
	cases := []struct {
		name              string
		mode              DecodeMode
		allowStatus       bool
		allowSystemOwned  bool
		allowSpecMutation bool
	}{
		{
			name:              "create-request-customer",
			mode:              ModeCreateRequest,
			allowStatus:       false,
			allowSystemOwned:  false,
			allowSpecMutation: true,
		},
		{
			name:              "replace-request-customer",
			mode:              ModeReplaceRequest,
			allowStatus:       false,
			allowSystemOwned:  false,
			allowSpecMutation: true,
		},
		{
			name:              "status-update-controller",
			mode:              ModeStatusUpdate,
			allowStatus:       true,
			allowSystemOwned:  true,
			allowSpecMutation: false,
		},
		{
			name:              "internal-object",
			mode:              ModeInternalObject,
			allowStatus:       true,
			allowSystemOwned:  true,
			allowSpecMutation: true,
		},
		{
			name:              "read-representation",
			mode:              ModeReadRepresentation,
			allowStatus:       true,
			allowSystemOwned:  true,
			allowSpecMutation: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			pol := PolicyFor(tc.mode)
			if pol.Mode != tc.mode {
				t.Fatalf("Mode = %v, want %v", pol.Mode, tc.mode)
			}
			if pol.AllowStatus != tc.allowStatus {
				t.Fatalf("AllowStatus = %v, want %v", pol.AllowStatus, tc.allowStatus)
			}
			if pol.AllowSystemOwned != tc.allowSystemOwned {
				t.Fatalf("AllowSystemOwned = %v, want %v", pol.AllowSystemOwned, tc.allowSystemOwned)
			}
			if pol.AllowSpecMutation != tc.allowSpecMutation {
				t.Fatalf("AllowSpecMutation = %v, want %v", pol.AllowSpecMutation, tc.allowSpecMutation)
			}
		})
	}
}

func TestPolicyForUnknownModeFailsClosed(t *testing.T) {
	pol := PolicyFor(DecodeMode(99))
	if pol.AllowStatus || pol.AllowSystemOwned || pol.AllowSpecMutation {
		t.Fatalf("unknown mode must fail closed, got %#v", pol)
	}
	if pol.Mode != DecodeMode(99) {
		t.Fatalf("Mode = %v, want 99", pol.Mode)
	}
}

func TestPolicyForCustomerModesRejectStatusAndSystem(t *testing.T) {
	const withStatus = `{"apiVersion":"v1","kind":"Project","metadata":{"name":"x"},"spec":{},"status":{"phase":"Ready"}}`
	const withUID = `{"apiVersion":"v1","kind":"Project","metadata":{"name":"x","uid":"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},"spec":{}}`

	for _, mode := range []DecodeMode{ModeCreateRequest, ModeReplaceRequest} {
		t.Run(mode.String()+"/status", func(t *testing.T) {
			var dst decodeSample
			prob := DecodeJSON([]byte(withStatus), testLimits, PolicyFor(mode), &dst)
			if prob == nil {
				t.Fatal("expected status rejection, got nil")
			}
			if prob.Code != apiproblem.CodeValidationFailed {
				t.Fatalf("Code = %q, want %q", prob.Code, apiproblem.CodeValidationFailed)
			}
			if prob.Violations[0].Field != "/status" {
				t.Fatalf("Field = %q, want /status", prob.Violations[0].Field)
			}
		})
		t.Run(mode.String()+"/system", func(t *testing.T) {
			var dst decodeSample
			prob := DecodeJSON([]byte(withUID), testLimits, PolicyFor(mode), &dst)
			if prob == nil {
				t.Fatal("expected system-owned rejection, got nil")
			}
			if prob.Code != apiproblem.CodeValidationFailed {
				t.Fatalf("Code = %q, want %q", prob.Code, apiproblem.CodeValidationFailed)
			}
			if prob.Violations[0].Field != "/metadata/uid" {
				t.Fatalf("Field = %q, want /metadata/uid", prob.Violations[0].Field)
			}
		})
	}
}

func TestPolicyForInternalAndReadModesAcceptStatusAndSystem(t *testing.T) {
	const raw = `{"apiVersion":"v1","kind":"Project","metadata":{"name":"x","uid":"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa","resourceVersion":"3","generation":1,"createdAt":"2026-01-01T00:00:00Z","updatedAt":"2026-01-02T00:00:00Z"},"spec":{"displayName":"Demo"},"status":{"phase":"Ready"}}`

	for _, mode := range []DecodeMode{ModeStatusUpdate, ModeInternalObject, ModeReadRepresentation} {
		t.Run(mode.String(), func(t *testing.T) {
			pol := PolicyFor(mode)
			var dst decodeSample
			// Status-update rejects spec; strip spec for that mode.
			body := raw
			if mode == ModeStatusUpdate {
				body = `{"apiVersion":"v1","kind":"Project","metadata":{"name":"x","uid":"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa","resourceVersion":"3"},"status":{"phase":"Ready"}}`
			}
			if prob := DecodeJSON([]byte(body), testLimits, pol, &dst); prob != nil {
				t.Fatalf("mode %v must accept status/system fields: %#v", mode, prob)
			}
			if dst.Status["phase"] != "Ready" {
				t.Fatalf("status not decoded under %v: %#v", mode, dst.Status)
			}
			if dst.Metadata.UID == "" {
				t.Fatalf("uid not decoded under %v", mode)
			}
		})
	}
}

func TestPolicyForStatusUpdateRejectsSpec(t *testing.T) {
	const raw = `{"apiVersion":"v1","kind":"Project","metadata":{"name":"x"},"spec":{"displayName":"Nope"},"status":{"phase":"Ready"}}`
	var dst decodeSample
	prob := DecodeJSON([]byte(raw), testLimits, PolicyFor(ModeStatusUpdate), &dst)
	if prob == nil {
		t.Fatal("expected spec rejection under ModeStatusUpdate, got nil")
	}
	if prob.Code != apiproblem.CodeValidationFailed {
		t.Fatalf("Code = %q, want %q", prob.Code, apiproblem.CodeValidationFailed)
	}
	if prob.Violations[0].Field != "/spec" {
		t.Fatalf("Field = %q, want /spec", prob.Violations[0].Field)
	}
}
