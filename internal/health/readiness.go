package health

import "sync/atomic"

// ReadinessState is an atomic boolean flag that tracks whether the server
// has completed initialization. Set to true by server.Start() after
// ListenAndServe begins; read by the readyz handler without locking.
type ReadinessState struct {
	ready atomic.Bool
}

// SetReady sets the readiness flag.
func (s *ReadinessState) SetReady(v bool) {
	s.ready.Store(v)
}

// IsReady returns whether the server has completed initialization.
func (s *ReadinessState) IsReady() bool {
	return s.ready.Load()
}
