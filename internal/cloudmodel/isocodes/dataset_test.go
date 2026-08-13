package isocodes

import "testing"

func TestAssignedAsOfPinned(t *testing.T) {
	t.Parallel()
	if AssignedAsOf != "2026-08-12" {
		t.Fatalf("AssignedAsOf=%q, want 2026-08-12", AssignedAsOf)
	}
}

func TestIsAssignedAlpha2KnownAssigned(t *testing.T) {
	t.Parallel()
	for _, code := range []string{"US", "IN", "GB", "DE", "JP", "AU", "BR", "ZA"} {
		if !IsAssignedAlpha2(code) {
			t.Fatalf("IsAssignedAlpha2(%q)=false, want true", code)
		}
	}
}

func TestIsAssignedAlpha2RejectsUnassignedAndLowercase(t *testing.T) {
	t.Parallel()
	for _, code := range []string{"ZZ", "AA", "QM", "XZ", "XK", "us", "Us", "uS", "", "USA", "U"} {
		if IsAssignedAlpha2(code) {
			t.Fatalf("IsAssignedAlpha2(%q)=true, want false (no normalization)", code)
		}
	}
}

func TestAssignedAlpha2SizeAndUppercaseOnly(t *testing.T) {
	t.Parallel()
	if got := len(AssignedAlpha2); got != 249 {
		t.Fatalf("AssignedAlpha2 len=%d, want 249", got)
	}
	for code := range AssignedAlpha2 {
		if len(code) != 2 {
			t.Fatalf("code %q length=%d, want 2", code, len(code))
		}
		for i := 0; i < 2; i++ {
			c := code[i]
			if c < 'A' || c > 'Z' {
				t.Fatalf("code %q contains non-uppercase ASCII letter", code)
			}
		}
	}
}
