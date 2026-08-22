package executiontarget

import (
	"sync"
	"time"
)

// Clock is the injected time source for FEATURE-0016 expiry freshness and
// once-per-clock-second audit retry (DD-06; REQ-F16-09).
type Clock interface {
	Now() time.Time
}

// MonotonicClock is the production clock. time.Now retains a monotonic reading
// for duration measurements while UTC wall time drives expiresAt comparisons.
type MonotonicClock struct{}

// Now returns the current UTC time with an attached monotonic reading.
func (MonotonicClock) Now() time.Time {
	return time.Now().UTC()
}

// FakeClock is a deterministic, mutex-guarded clock for tests. It never waits
// on wall time; callers Advance or Set explicitly.
type FakeClock struct {
	mu  sync.Mutex
	now time.Time
}

// NewFakeClock returns a FakeClock fixed at t (normalized to UTC).
func NewFakeClock(t time.Time) *FakeClock {
	return &FakeClock{now: t.UTC()}
}

// Now returns the current fake time.
func (c *FakeClock) Now() time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.now
}

// Set replaces the fake time with t (UTC).
func (c *FakeClock) Set(t time.Time) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.now = t.UTC()
}

// Advance moves the fake clock forward by d.
func (c *FakeClock) Advance(d time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.now = c.now.Add(d)
}
