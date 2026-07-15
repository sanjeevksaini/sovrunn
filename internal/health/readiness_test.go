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
