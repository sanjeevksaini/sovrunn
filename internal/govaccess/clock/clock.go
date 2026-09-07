// Package clock is the FEATURE-0018 deterministic UTC time leaf (DD-05; REQ-F18-22).
//
// Phase 2R supplies FixedClock and a manually advanced FixtureClock only.
// No UTCClock, time.Now, or production wall-clock wiring is implemented here.
package clock

import (
	"sync"
	"time"
)

// Clock supplies authoritative UTC wall time for FEATURE-0018 owners.
// Concurrent implementations must be concurrency-safe.
type Clock interface {
	NowUTC() time.Time
}

// FixedClock is an immutable, deterministic Clock. Repeated NowUTC calls
// return the same configured instant until a new FixedClock is constructed.
type FixedClock struct {
	at time.Time
}

// NewFixedClock returns a FixedClock fixed at t, normalized to UTC.
func NewFixedClock(t time.Time) FixedClock {
	return FixedClock{at: t.UTC()}
}

// NowUTC returns the configured UTC instant unchanged.
func (f FixedClock) NowUTC() time.Time {
	return f.at
}

// FixtureClock is a mutex-guarded, manually advanced Clock for tests and
// deterministic fixtures. It never waits on wall time.
type FixtureClock struct {
	mu  sync.Mutex
	now time.Time
}

// NewFixtureClock returns a FixtureClock fixed at t (normalized to UTC).
func NewFixtureClock(t time.Time) *FixtureClock {
	return &FixtureClock{now: t.UTC()}
}

// NowUTC returns the current fixture time.
func (c *FixtureClock) NowUTC() time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.now
}

// Set replaces the fixture time with t (UTC).
func (c *FixtureClock) Set(t time.Time) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.now = t.UTC()
}

// Advance moves the fixture clock forward by d.
func (c *FixtureClock) Advance(d time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.now = c.now.Add(d)
}
