package policyeval

import "time"

// TimeSource supplies UTC wall time for evaluation chronology.
//
// Implementations used concurrently must be concurrency-safe. The evaluation
// boundary normalizes returned values with time.Time.UTC and serializes them
// with time.RFC3339Nano. This seam is local to FEATURE-0017 and does not reuse
// FEATURE-0016 Clock types.
type TimeSource interface {
	NowUTC() time.Time
}

// UTCTimeSource is the production TimeSource. It returns the current wall
// time normalized to UTC.
type UTCTimeSource struct{}

// NowUTC returns the current UTC time.
func (UTCTimeSource) NowUTC() time.Time {
	return time.Now().UTC()
}

// FixedTimeSource is an immutable, deterministic TimeSource for tests and
// replay. Repeated NowUTC calls return the same configured instant.
type FixedTimeSource struct {
	at time.Time
}

// NewFixedTimeSource returns a FixedTimeSource fixed at t, normalized to UTC.
func NewFixedTimeSource(t time.Time) FixedTimeSource {
	return FixedTimeSource{at: t.UTC()}
}

// NowUTC returns the configured UTC instant unchanged.
func (f FixedTimeSource) NowUTC() time.Time {
	return f.at
}
