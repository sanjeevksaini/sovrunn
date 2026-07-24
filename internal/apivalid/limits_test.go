package apivalid

import "testing"

func TestDefaultLimitsMatchDesignTable(t *testing.T) {
	t.Parallel()

	got := DefaultLimits()
	want := Limits{
		MaxObjectBytes:        1_048_576, // 1 MiB
		MaxNestingDepth:       32,
		MaxLabels:             64,
		MaxLabelKeyChars:      63,
		MaxLabelValueChars:    253,
		MaxAnnotationsBytes:   262_144, // 256 KiB
		MaxConditions:         32,
		MaxReferencesPerField: 64,
		MaxViolations:         100,
		DefaultPageSize:       50,
		MaxPageSize:           200,
	}
	if got != want {
		t.Fatalf("DefaultLimits() = %#v, want %#v", got, want)
	}

	// Equivalence checks against the design rationale units.
	if got.MaxObjectBytes != 1<<20 {
		t.Errorf("MaxObjectBytes = %d, want 1<<20 (1 MiB)", got.MaxObjectBytes)
	}
	if got.MaxAnnotationsBytes != 256<<10 {
		t.Errorf("MaxAnnotationsBytes = %d, want 256<<10 (256 KiB)", got.MaxAnnotationsBytes)
	}
}
