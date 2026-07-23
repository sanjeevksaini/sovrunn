package apimeta

import "testing"

func TestMatrixAProfiles(t *testing.T) {
	t.Parallel()

	want := []Profile{
		"ManagedResource",
		"ObservedExternalResource",
		"VersionedDefinition",
		"ImmutableRecord",
		"LongRunningOperation",
		"TransientRequestResult",
		"EmbeddedValue",
		"ListEnvelope",
	}
	got := AllProfiles()
	if len(got) != len(want) {
		t.Fatalf("AllProfiles len=%d, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("AllProfiles[%d]=%q, want %q", i, got[i], want[i])
		}
		if !got[i].Valid() {
			t.Fatalf("profile %q should be Valid", got[i])
		}
	}
	if Profile("UnknownProfile").Valid() {
		t.Fatal("unknown profile must not be Valid")
	}
}

func TestMatrixC1Boundaries(t *testing.T) {
	t.Parallel()

	want := []Boundary{
		"customer-facing",
		"operator-facing",
		"internal-engine-facing",
		"adapter-facing",
		"plugin-facing",
		"governance-only",
	}
	got := AllBoundaries()
	if len(got) != len(want) {
		t.Fatalf("AllBoundaries len=%d, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("AllBoundaries[%d]=%q, want %q", i, got[i], want[i])
		}
		if !got[i].Valid() {
			t.Fatalf("boundary %q should be Valid", got[i])
		}
	}
	if Boundary("public").Valid() {
		t.Fatal("unknown boundary must not be Valid")
	}
}

func TestStabilityVocabulary(t *testing.T) {
	t.Parallel()

	want := []Stability{"alpha", "beta", "stable"}
	got := AllStabilities()
	if len(got) != len(want) {
		t.Fatalf("AllStabilities len=%d, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("AllStabilities[%d]=%q, want %q", i, got[i], want[i])
		}
		if !got[i].Valid() {
			t.Fatalf("stability %q should be Valid", got[i])
		}
	}
	if Stability("ga").Valid() {
		t.Fatal("unknown stability must not be Valid")
	}
}

func TestDataClassificationVocabulary(t *testing.T) {
	t.Parallel()

	want := []DataClassification{
		"Public",
		"Customer-visible",
		"Tenant-confidential",
		"Operator-confidential",
		"Internal",
		"Sensitive",
		"Secret-reference-only",
	}
	got := AllDataClassifications()
	if len(got) != len(want) {
		t.Fatalf("AllDataClassifications len=%d, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("AllDataClassifications[%d]=%q, want %q", i, got[i], want[i])
		}
		if !got[i].Valid() {
			t.Fatalf("classification %q should be Valid", got[i])
		}
	}
	if DataClassification("Top-secret").Valid() {
		t.Fatal("unknown classification must not be Valid")
	}
}

func TestObjectMetaFieldPresence(t *testing.T) {
	t.Parallel()

	// Compile-time / structural smoke: ObjectMeta carries the F12-META-001 subset
	// including a ScopeRef pointer slot completed in task 2.2.
	meta := ObjectMeta{
		Name:        "payments-production",
		DisplayName: "Payments Production",
		ScopeRef:    &ScopeRef{},
		Labels:      map[string]string{"env": "prod"},
		Annotations: map[string]string{"sovrunn.io/note": "demo"},
		Generation:  1,
	}
	if meta.Name == "" || meta.ScopeRef == nil {
		t.Fatal("ObjectMeta smoke fields must be set")
	}
}
