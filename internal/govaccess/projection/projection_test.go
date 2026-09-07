package projection

import "testing"

func TestAllJourneysHasApprovedSix(t *testing.T) {
	journeys := AllJourneys()
	if len(journeys) != 6 {
		t.Fatalf("journey count=%d want 6", len(journeys))
	}
	seen := map[Journey]bool{}
	for _, j := range journeys {
		if seen[j] {
			t.Fatalf("duplicate journey %q", j)
		}
		seen[j] = true
	}
}

func TestProjectProgressiveDisclosurePerJourney(t *testing.T) {
	t.Parallel()
	input := map[string]string{
		"action":          "read",
		"target":          "target-1",
		"result":          "Allow",
		"status":          "Succeeded",
		"assignee":        "alice",
		"role":            "operator",
		"scope":           "project-a",
		"expiresAt":       "2026-09-08T00:00:00Z",
		"requestId":       "req-1",
		"requestedRole":   "db-admin",
		"justification":   "incident response",
		"duration":        "2h",
		"subject":         "workload-1",
		"decision":        "Approve",
		"reason":          "required",
		"reviewId":        "rev-1",
		"holder":          "group-1",
		"recommendation":  "Retain",
		"control":         "change-window",
		"effectiveFrom":   "2026-09-07T00:00:00Z",
		"effectiveUntil":  "2026-09-08T00:00:00Z",
		"policyDigest":    "hidden",
		"issuer":          "hidden",
		"subjectToken":    "hidden",
		"controllerState": "hidden",
	}
	for _, journey := range AllJourneys() {
		view, ok := ViewForJourney(journey)
		if !ok {
			t.Fatalf("missing view for %q", journey)
		}
		got, err := Project(journey, input)
		if err != nil {
			t.Fatalf("project(%s): %v", journey, err)
		}
		if len(got) == 0 {
			t.Fatalf("project(%s): expected non-empty projection", journey)
		}
		allowed := map[string]struct{}{}
		for _, key := range view.Fields {
			allowed[key] = struct{}{}
		}
		for k := range got {
			if _, ok := allowed[k]; !ok {
				t.Fatalf("project(%s): leaked non-journey field %q", journey, k)
			}
		}
		if _, ok := got["policyDigest"]; ok {
			t.Fatalf("project(%s): leaked sensitive field policyDigest", journey)
		}
		if _, ok := got["controllerState"]; ok {
			t.Fatalf("project(%s): leaked controller field controllerState", journey)
		}
	}
}

func TestRedactDropsSensitiveAndInternalKeys(t *testing.T) {
	t.Parallel()
	got := Redact(map[string]string{
		"action":        "read",
		"target":        "resource-x",
		"token":         "secret-value",
		"private_key":   "pk",
		"auditDigest":   "sensitive",
		"internalState": "controller-only",
	})
	if _, ok := got["token"]; ok {
		t.Fatal("token must be redacted")
	}
	if _, ok := got["private_key"]; ok {
		t.Fatal("private_key must be redacted")
	}
	if _, ok := got["auditDigest"]; ok {
		t.Fatal("auditDigest must be redacted")
	}
	if _, ok := got["internalState"]; ok {
		t.Fatal("internalState must be redacted")
	}
	if got["action"] == "" || got["target"] == "" {
		t.Fatal("non-sensitive fields must remain")
	}
}

func TestSafeDenialProjectionIsGeneric(t *testing.T) {
	t.Parallel()
	got := SafeDenialProjection()
	if got.Code != "RESOURCE_NOT_FOUND" {
		t.Fatalf("code=%s", got.Code)
	}
	if got.Message == "" {
		t.Fatal("message must be present")
	}
	if got.Message == "resource exists but forbidden" {
		t.Fatal("safe denial must not disclose existence")
	}
}

func TestProjectNoControllerStateProperty(t *testing.T) {
	t.Parallel()
	for _, journey := range AllJourneys() {
		got, err := Project(journey, map[string]string{
			"action":          "update",
			"target":          "x",
			"status":          "ok",
			"resourceVersion": "7",
			"controllerState": "internal",
			"auditLink":       "must-hide",
			"correlationId":   "must-hide",
		})
		if err != nil {
			t.Fatalf("project(%s): %v", journey, err)
		}
		for _, forbidden := range []string{"resourceVersion", "controllerState", "auditLink", "correlationId"} {
			if _, ok := got[forbidden]; ok {
				t.Fatalf("project(%s): leaked %s", journey, forbidden)
			}
		}
	}
}
