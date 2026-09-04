package idempotency

import (
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/operation"
)

// InFlightReservation is a sealed value carrying the full RequestBinding for
// an in-flight caller mutation (DD-11). Table operations are Task-6-owned.
type InFlightReservation struct {
	sealed     bool
	lookupKey  IdempotencyLookupKey
	binding    RequestBinding
	ownerToken []byte
}

// NewInFlightReservation seals a structurally valid reservation value.
// Owner token is deep-copied. Lifecycle transitions are not performed here.
func NewInFlightReservation(
	lookupKey IdempotencyLookupKey,
	binding RequestBinding,
	ownerToken []byte,
) (InFlightReservation, operation.MechanicalFailure) {
	if !lookupKey.sealed {
		return InFlightReservation{}, operation.NewMechanicalFailure("in_flight_reservation_lookup_unsealed")
	}
	if !binding.sealed {
		return InFlightReservation{}, operation.NewMechanicalFailure("in_flight_reservation_binding_unsealed")
	}
	if len(ownerToken) == 0 {
		return InFlightReservation{}, operation.NewMechanicalFailure("in_flight_reservation_owner_token_empty")
	}
	cp := make([]byte, len(ownerToken))
	copy(cp, ownerToken)
	return InFlightReservation{
		sealed:     true,
		lookupKey:  lookupKey,
		binding:    binding,
		ownerToken: cp,
	}, operation.MechanicalFailure{}
}

// Sealed reports whether the reservation was constructed through NewInFlightReservation.
func (r InFlightReservation) Sealed() bool { return r.sealed }

// LookupKey returns the sealed lookup key.
func (r InFlightReservation) LookupKey() IdempotencyLookupKey { return r.lookupKey }

// Binding returns the full sealed RequestBinding.
func (r InFlightReservation) Binding() RequestBinding { return r.binding }

// OwnerToken returns a defensive copy of the owner token.
func (r InFlightReservation) OwnerToken() []byte {
	if !r.sealed || len(r.ownerToken) == 0 {
		return nil
	}
	cp := make([]byte, len(r.ownerToken))
	copy(cp, r.ownerToken)
	return cp
}
