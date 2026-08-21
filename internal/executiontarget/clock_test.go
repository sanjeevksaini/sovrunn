package executiontarget

import (
	"testing"
	"time"
)

func TestMonotonicClock_NowUTC(t *testing.T) {
	t.Parallel()

	var c MonotonicClock
	before := time.Now().UTC().Add(-time.Second)
	got := c.Now()
	after := time.Now().UTC().Add(time.Second)
	if got.Location() != time.UTC {
		t.Fatalf("location=%v", got.Location())
	}
	if got.Before(before) || got.After(after) {
		t.Fatalf("now=%v out of expected window", got)
	}
}

func TestFakeClock_DeterministicAdvance(t *testing.T) {
	t.Parallel()

	start := time.Date(2026, 8, 21, 10, 0, 0, 0, time.UTC)
	c := NewFakeClock(start)
	if !c.Now().Equal(start) {
		t.Fatalf("initial=%v", c.Now())
	}
	c.Advance(time.Second)
	if !c.Now().Equal(start.Add(time.Second)) {
		t.Fatalf("after advance=%v", c.Now())
	}
	c.Set(start.Add(5 * time.Second))
	if !c.Now().Equal(start.Add(5 * time.Second)) {
		t.Fatalf("after set=%v", c.Now())
	}
}

func TestFakeClock_ConcurrentSafe(t *testing.T) {
	t.Parallel()

	c := NewFakeClock(time.Date(2026, 8, 21, 0, 0, 0, 0, time.UTC))
	done := make(chan struct{})
	go func() {
		for i := 0; i < 100; i++ {
			c.Advance(time.Millisecond)
			_ = c.Now()
		}
		close(done)
	}()
	for i := 0; i < 100; i++ {
		c.Advance(time.Millisecond)
		_ = c.Now()
	}
	<-done
}
