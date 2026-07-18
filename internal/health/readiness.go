package health

import "sync/atomic"

// Reason constants for not-ready ReadinessState.
const (
	ReasonInitializing = "initializing"
	ReasonShuttingDown = "shutting_down"
)

// ReadinessState is an atomic flag that tracks whether the server has
// completed initialization, plus an optional reason when not ready.
// The zero value is fully usable: not ready with reason "initializing".
type ReadinessState struct {
	ready  atomic.Bool
	reason atomic.Pointer[string]
}

// NewReadinessState returns a ReadinessState that is not ready with
// reason ReasonInitializing.
func NewReadinessState() *ReadinessState {
	s := &ReadinessState{}
	s.SetInitializing()
	return s
}

// SetReady sets the readiness flag. When v is true, the reason pointer
// is cleared. When v is false, the reason is left unchanged.
func (s *ReadinessState) SetReady(v bool) {
	s.ready.Store(v)
	if v {
		s.reason.Store(nil)
	}
}

// SetInitializing marks the state not ready with ReasonInitializing.
func (s *ReadinessState) SetInitializing() {
	r := ReasonInitializing
	s.reason.Store(&r)
	s.ready.Store(false)
}

// SetShuttingDown marks the state not ready with ReasonShuttingDown.
func (s *ReadinessState) SetShuttingDown() {
	r := ReasonShuttingDown
	s.reason.Store(&r)
	s.ready.Store(false)
}

// IsReady returns whether the server has completed initialization.
func (s *ReadinessState) IsReady() bool {
	return s.ready.Load()
}

// Reason returns the not-ready reason. When ready it returns "".
// When not ready and the reason pointer is nil, it returns
// ReasonInitializing (zero-value safe).
func (s *ReadinessState) Reason() string {
	if s.IsReady() {
		return ""
	}
	if p := s.reason.Load(); p != nil {
		return *p
	}
	return ReasonInitializing
}
