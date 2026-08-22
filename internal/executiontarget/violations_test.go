package executiontarget

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
)

func TestLocalViolationCodes_ExactlyEightClosed(t *testing.T) {
	t.Parallel()

	codes := AllLocalViolationCodes()
	if len(codes) != 8 {
		t.Fatalf("want exactly 8 local violation codes, got %d: %#v", len(codes), codes)
	}

	want := []apiproblem.ViolationCode{
		"VS0_TARGET_RETIRED",
		"VS0_TARGET_MAINTENANCE",
		"VS0_TARGET_QUALIFICATION_IN_PROGRESS",
		"VS0_TARGET_EPOCH_STALE",
		"VS0_EXECUTION_TARGET_PARTICIPATION_UNAVAILABLE",
		"VS0_EXECUTION_TARGET_STACK_UNAVAILABLE",
		"VS0_EXECUTION_TARGET_SCOPE_MISMATCH",
		"VS0_EXECUTION_TARGET_VIABILITY_STALE",
	}
	for i, c := range want {
		if codes[i] != c {
			t.Fatalf("codes[%d]=%q, want %q", i, codes[i], c)
		}
	}

	seen := make(map[apiproblem.ViolationCode]struct{}, len(codes))
	for _, c := range codes {
		if _, dup := seen[c]; dup {
			t.Fatalf("duplicate local violation code %q", c)
		}
		seen[c] = struct{}{}
		if string(c) == "" {
			t.Fatal("empty violation code")
		}
	}
}

func TestLocalViolationCodes_NamedConstantsMatch(t *testing.T) {
	t.Parallel()

	checks := []struct {
		got  apiproblem.ViolationCode
		want string
	}{
		{ViolationTargetRetired, "VS0_TARGET_RETIRED"},
		{ViolationTargetMaintenance, "VS0_TARGET_MAINTENANCE"},
		{ViolationTargetQualificationInProgress, "VS0_TARGET_QUALIFICATION_IN_PROGRESS"},
		{ViolationTargetEpochStale, "VS0_TARGET_EPOCH_STALE"},
		{ViolationExecutionTargetParticipationUnavailable, "VS0_EXECUTION_TARGET_PARTICIPATION_UNAVAILABLE"},
		{ViolationExecutionTargetStackUnavailable, "VS0_EXECUTION_TARGET_STACK_UNAVAILABLE"},
		{ViolationExecutionTargetScopeMismatch, "VS0_EXECUTION_TARGET_SCOPE_MISMATCH"},
		{ViolationExecutionTargetViabilityStale, "VS0_EXECUTION_TARGET_VIABILITY_STALE"},
	}
	for _, tc := range checks {
		if string(tc.got) != tc.want {
			t.Fatalf("constant=%q, want %q", tc.got, tc.want)
		}
	}
}

func TestLocalViolationCodes_InheritedNotRedefined(t *testing.T) {
	t.Parallel()

	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	srcPath := filepath.Join(filepath.Dir(thisFile), "violations.go")
	raw, err := os.ReadFile(srcPath)
	if err != nil {
		t.Fatalf("read violations.go: %v", err)
	}
	src := string(raw)

	inherited := []string{
		"VS0_STATUS_FIELD_WRITE",
		"VS0_SYSTEM_OWNED_FIELD_WRITE",
		"VS0_AUTHORIZATION_SAFE_DENIAL",
		"VS0_IDEMPOTENCY_KEY_REUSE_MISMATCH",
	}
	for _, code := range inherited {
		if strings.Contains(src, `"`+code+`"`) {
			t.Fatalf("inherited violation code %s must not be redefined in violations.go", code)
		}
		for _, local := range AllLocalViolationCodes() {
			if string(local) == code {
				t.Fatalf("inherited violation code %s must not appear in AllLocalViolationCodes", code)
			}
		}
	}
}
