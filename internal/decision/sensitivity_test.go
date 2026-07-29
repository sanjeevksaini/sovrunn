package decision

import (
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
)

func TestSensitivityFloorMapExact(t *testing.T) {
	t.Parallel()

	want := map[apimeta.DataClassification]Sensitivity{
		apimeta.ClassPublic:               SensitivityPublic,
		apimeta.ClassCustomerVisible:      SensitivityConfidential,
		apimeta.ClassTenantConfidential:   SensitivityConfidential,
		apimeta.ClassOperatorConfidential: SensitivityConfidential,
		apimeta.ClassInternal:             SensitivityInternal,
		apimeta.ClassSensitive:            SensitivityRestricted,
		apimeta.ClassSecretReferenceOnly:  SensitivityRestricted,
	}

	got := AllDataClassificationFloors(t)
	if len(got) != len(want) {
		t.Fatalf("mapped classifications = %d, want %d", len(got), len(want))
	}
	for class, floor := range want {
		gotFloor, ok := SensitivityFloor(class)
		if !ok {
			t.Fatalf("SensitivityFloor(%q) unexpectedly unmapped", class)
		}
		if gotFloor != floor {
			t.Fatalf("SensitivityFloor(%q) = %q, want %q", class, gotFloor, floor)
		}
		if got[class] != floor {
			t.Fatalf("AllDataClassificationFloors[%q] = %q, want %q", class, got[class], floor)
		}
	}

	// Every FEATURE-0012 classification must be mapped (no silent gaps).
	for _, class := range apimeta.AllDataClassifications() {
		if _, ok := want[class]; !ok {
			t.Fatalf("FEATURE-0012 classification %q missing from floor map", class)
		}
	}
}

// AllDataClassificationFloors is a test helper that materializes the full floor map.
func AllDataClassificationFloors(t *testing.T) map[apimeta.DataClassification]Sensitivity {
	t.Helper()
	out := make(map[apimeta.DataClassification]Sensitivity, len(apimeta.AllDataClassifications()))
	for _, class := range apimeta.AllDataClassifications() {
		floor, ok := SensitivityFloor(class)
		if !ok {
			t.Fatalf("SensitivityFloor(%q) must succeed for every FEATURE-0012 classification", class)
		}
		out[class] = floor
	}
	return out
}

func TestSensitivityFloorUnmappedFailsClosed(t *testing.T) {
	t.Parallel()

	for _, class := range []apimeta.DataClassification{
		"",
		"Top-secret",
		"PUBLIC",
		"public",
		"Unknown",
	} {
		if floor, ok := SensitivityFloor(class); ok {
			t.Fatalf("SensitivityFloor(%q) = %q, want fail-closed", class, floor)
		}
	}
}

func TestSensitivityFloorRaiseOnly(t *testing.T) {
	t.Parallel()

	// Public floor is PUBLIC: raising to any higher core value is allowed.
	for _, raised := range AllSensitivities() {
		got, ok := ApplyRaisedFloor(SensitivityPublic, raised)
		if !ok {
			t.Fatalf("raising PUBLIC floor to %q must succeed", raised)
		}
		if got != raised {
			t.Fatalf("ApplyRaisedFloor(PUBLIC, %q) = %q", raised, got)
		}
	}

	// Customer-visible floor is CONFIDENTIAL: may raise to RESTRICTED, not lower.
	floor, ok := SensitivityFloor(apimeta.ClassCustomerVisible)
	if !ok || floor != SensitivityConfidential {
		t.Fatalf("Customer-visible floor = %q ok=%v, want CONFIDENTIAL", floor, ok)
	}
	if got, ok := ApplyRaisedFloor(floor, SensitivityRestricted); !ok || got != SensitivityRestricted {
		t.Fatalf("raise CONFIDENTIAL→RESTRICTED = %q ok=%v", got, ok)
	}
	if got, ok := ApplyRaisedFloor(floor, SensitivityConfidential); !ok || got != SensitivityConfidential {
		t.Fatalf("exact floor CONFIDENTIAL must be accepted; got %q ok=%v", got, ok)
	}
	for _, below := range []Sensitivity{SensitivityPublic, SensitivityInternal} {
		if got, ok := ApplyRaisedFloor(floor, below); ok {
			t.Fatalf("lowering CONFIDENTIAL floor to %q must fail; got %q", below, got)
		}
	}
}

func TestSensitivityBelowFloorFailsClosed(t *testing.T) {
	t.Parallel()

	cases := []struct {
		class    apimeta.DataClassification
		declared Sensitivity
	}{
		{apimeta.ClassCustomerVisible, SensitivityPublic},
		{apimeta.ClassCustomerVisible, SensitivityInternal},
		{apimeta.ClassTenantConfidential, SensitivityPublic},
		{apimeta.ClassOperatorConfidential, SensitivityInternal},
		{apimeta.ClassInternal, SensitivityPublic},
		{apimeta.ClassSensitive, SensitivityPublic},
		{apimeta.ClassSensitive, SensitivityInternal},
		{apimeta.ClassSensitive, SensitivityConfidential},
		{apimeta.ClassSecretReferenceOnly, SensitivityConfidential},
	}
	for _, tc := range cases {
		if CheckClassificationSensitivity(tc.class, tc.declared) {
			t.Fatalf("CheckClassificationSensitivity(%q, %q) must fail below-floor", tc.class, tc.declared)
		}
		if MeetsFloor(tc.declared, mustFloor(t, tc.class)) {
			t.Fatalf("MeetsFloor(%q, floor(%q)) must be false", tc.declared, tc.class)
		}
	}
}

func TestSensitivityAtOrAboveFloorAccepted(t *testing.T) {
	t.Parallel()

	cases := []struct {
		class    apimeta.DataClassification
		declared Sensitivity
	}{
		{apimeta.ClassPublic, SensitivityPublic},
		{apimeta.ClassPublic, SensitivityRestricted},
		{apimeta.ClassInternal, SensitivityInternal},
		{apimeta.ClassInternal, SensitivityConfidential},
		{apimeta.ClassCustomerVisible, SensitivityConfidential},
		{apimeta.ClassCustomerVisible, SensitivityRestricted},
		{apimeta.ClassSensitive, SensitivityRestricted},
		{apimeta.ClassSecretReferenceOnly, SensitivityRestricted},
	}
	for _, tc := range cases {
		if !CheckClassificationSensitivity(tc.class, tc.declared) {
			t.Fatalf("CheckClassificationSensitivity(%q, %q) must succeed", tc.class, tc.declared)
		}
	}
}

func TestSensitivityUnknownDeclaredFailsClosed(t *testing.T) {
	t.Parallel()

	if CheckClassificationSensitivity(apimeta.ClassPublic, Sensitivity("TOP_SECRET")) {
		t.Fatal("unknown declared Sensitivity must fail closed")
	}
	if MeetsFloor(Sensitivity("TOP_SECRET"), SensitivityPublic) {
		t.Fatal("MeetsFloor with unknown declared must fail closed")
	}
	if MeetsFloor(SensitivityPublic, Sensitivity("TOP_SECRET")) {
		t.Fatal("MeetsFloor with unknown floor must fail closed")
	}
	if _, ok := ApplyRaisedFloor(SensitivityPublic, Sensitivity("TOP_SECRET")); ok {
		t.Fatal("ApplyRaisedFloor with unknown raised must fail closed")
	}
}

func TestSensitivityUnmappedClassificationCheckFails(t *testing.T) {
	t.Parallel()

	if CheckClassificationSensitivity(apimeta.DataClassification("Top-secret"), SensitivityRestricted) {
		t.Fatal("unmapped classification must fail closed even with Valid declared Sensitivity")
	}
}

func TestMostRestrictiveFloor(t *testing.T) {
	t.Parallel()

	got, ok := MostRestrictiveFloor(
		apimeta.ClassPublic,          // PUBLIC
		apimeta.ClassInternal,        // INTERNAL
		apimeta.ClassCustomerVisible, // CONFIDENTIAL
	)
	if !ok || got != SensitivityConfidential {
		t.Fatalf("MostRestrictiveFloor = %q ok=%v, want CONFIDENTIAL", got, ok)
	}

	got, ok = MostRestrictiveFloor(
		apimeta.ClassCustomerVisible, // CONFIDENTIAL
		apimeta.ClassSensitive,       // RESTRICTED
	)
	if !ok || got != SensitivityRestricted {
		t.Fatalf("MostRestrictiveFloor = %q ok=%v, want RESTRICTED", got, ok)
	}

	if _, ok := MostRestrictiveFloor(); ok {
		t.Fatal("empty classification list must fail closed")
	}
	if _, ok := MostRestrictiveFloor(apimeta.ClassPublic, apimeta.DataClassification("Unknown")); ok {
		t.Fatal("any unmapped classification must fail closed")
	}
}

func TestMostRestrictiveSensitivity(t *testing.T) {
	t.Parallel()

	got, ok := MostRestrictiveSensitivity(
		SensitivityPublic,
		SensitivityConfidential,
		SensitivityInternal,
	)
	if !ok || got != SensitivityConfidential {
		t.Fatalf("MostRestrictiveSensitivity = %q ok=%v, want CONFIDENTIAL", got, ok)
	}
	if _, ok := MostRestrictiveSensitivity(); ok {
		t.Fatal("empty sensitivity list must fail closed")
	}
	if _, ok := MostRestrictiveSensitivity(SensitivityPublic, Sensitivity("NOPE")); ok {
		t.Fatal("any invalid sensitivity must fail closed")
	}
}

func TestSensitivityRankOrder(t *testing.T) {
	t.Parallel()

	order := AllSensitivities()
	for i := 1; i < len(order); i++ {
		if order[i-1].Rank() >= order[i].Rank() {
			t.Fatalf("rank order broken: %q(%d) >= %q(%d)",
				order[i-1], order[i-1].Rank(), order[i], order[i].Rank())
		}
	}
	if Sensitivity("NOPE").Rank() != -1 {
		t.Fatal("invalid Sensitivity Rank must be -1")
	}
}

func mustFloor(t *testing.T, class apimeta.DataClassification) Sensitivity {
	t.Helper()
	floor, ok := SensitivityFloor(class)
	if !ok {
		t.Fatalf("SensitivityFloor(%q) unexpectedly failed", class)
	}
	return floor
}
