// Package govaccess is the FEATURE-0018 composition root (DD-01).
//
// Task 1 wires only foundation leaves. Semantic owners, state, uow, evidence,
// and HTTP routes are composed by later tasks. This package holds no live state
// and performs no semantic logic.
package govaccess

import (
	"errors"

	"github.com/sanjeevksaini/sovrunn/internal/govaccess/clock"
)

// ErrNilClock is returned when composition is attempted without a Clock.
var ErrNilClock = errors.New("govaccess: clock is required")

// Foundation is the Task-1 wiring surface: deterministic clock injection only.
type Foundation struct {
	Clock clock.Clock
}

// NewFoundation wires foundation leaves. Child semantic packages are composed
// by later tasks; this function constructs no records and opens no transactions.
func NewFoundation(clk clock.Clock) (Foundation, error) {
	if clk == nil {
		return Foundation{}, ErrNilClock
	}
	return Foundation{Clock: clk}, nil
}
