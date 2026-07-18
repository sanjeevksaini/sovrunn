package health

import "testing"

func TestReadinessState_DefaultNotReady(t *testing.T) {
	var rs ReadinessState
	if rs.IsReady() {
		t.Fatal("default readiness should be false")
	}
}

func TestReadinessState_SetReadyTrue(t *testing.T) {
	var rs ReadinessState
	rs.SetReady(true)
	if !rs.IsReady() {
		t.Fatal("expected ready after SetReady(true)")
	}
}

func TestReadinessState_SetReadyFalse(t *testing.T) {
	var rs ReadinessState
	rs.SetReady(true)
	rs.SetReady(false)
	if rs.IsReady() {
		t.Fatal("expected not ready after SetReady(false)")
	}
}

func TestReadinessState_ZeroValueReason(t *testing.T) {
	var rs ReadinessState
	if got := rs.Reason(); got != ReasonInitializing {
		t.Fatalf("Reason() = %q, want %q", got, ReasonInitializing)
	}
}

func TestReadinessState_NewDefaultsInitializing(t *testing.T) {
	rs := NewReadinessState()
	if rs.IsReady() {
		t.Fatal("NewReadinessState should not be ready")
	}
	if got := rs.Reason(); got != ReasonInitializing {
		t.Fatalf("Reason() = %q, want %q", got, ReasonInitializing)
	}
}

func TestReadinessState_SetReadyTrueClearsReason(t *testing.T) {
	rs := NewReadinessState()
	rs.SetReady(true)
	if got := rs.Reason(); got != "" {
		t.Fatalf("Reason() = %q, want empty", got)
	}
}

func TestReadinessState_SetReadyFalseDefaultsReason(t *testing.T) {
	var rs ReadinessState
	rs.SetReady(true)
	rs.SetReady(false)
	if got := rs.Reason(); got != ReasonInitializing {
		t.Fatalf("Reason() = %q, want %q", got, ReasonInitializing)
	}
}

func TestReadinessState_SetShuttingDown(t *testing.T) {
	var rs ReadinessState
	rs.SetShuttingDown()
	if rs.IsReady() {
		t.Fatal("expected not ready after SetShuttingDown")
	}
	if got := rs.Reason(); got != ReasonShuttingDown {
		t.Fatalf("Reason() = %q, want %q", got, ReasonShuttingDown)
	}
}

func TestReadinessState_SetInitializing(t *testing.T) {
	var rs ReadinessState
	rs.SetReady(true)
	rs.SetInitializing()
	if rs.IsReady() {
		t.Fatal("expected not ready after SetInitializing")
	}
	if got := rs.Reason(); got != ReasonInitializing {
		t.Fatalf("Reason() = %q, want %q", got, ReasonInitializing)
	}
}

func TestReadinessState_SetShuttingDownThenSetReady(t *testing.T) {
	var rs ReadinessState
	rs.SetShuttingDown()
	rs.SetReady(true)
	if got := rs.Reason(); got != "" {
		t.Fatalf("Reason() = %q, want empty", got)
	}
}
