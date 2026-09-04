// Package idempotency holds FEATURE-0018 caller-keyed idempotency value types
// (DD-11; design §3.4, §5.4).
//
// Task 2 owns the sealed value-type contract and pure structural validation
// only. Reservation-table operations, replay handling, lifecycle transitions,
// PrepareCompletedResult, and UOW coordination are owned by Task 6.
package idempotency

import (
	"unicode/utf8"

	"github.com/sanjeevksaini/sovrunn/internal/apivalid"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/model"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/operation"
)

// IdempotencyLookupKey is the caller-supplied lookup identity. It is
// deliberately separate from RequestBinding (DD-11).
type IdempotencyLookupKey struct {
	sealed               bool
	actorRef             model.PrincipalRef
	operation            string
	clientIdempotencyKey string
}

// NewIdempotencyLookupKey seals a structurally valid lookup key.
func NewIdempotencyLookupKey(
	actorRef model.PrincipalRef,
	operationName string,
	clientIdempotencyKey string,
) (IdempotencyLookupKey, operation.MechanicalFailure) {
	if actorRef.Issuer == "" || actorRef.Subject == "" || !actorRef.PrincipalType.Valid() {
		return IdempotencyLookupKey{}, operation.NewMechanicalFailure("idempotency_lookup_key_actor_invalid")
	}
	if operationName == "" {
		return IdempotencyLookupKey{}, operation.NewMechanicalFailure("idempotency_lookup_key_operation_empty")
	}
	if clientIdempotencyKey == "" {
		return IdempotencyLookupKey{}, operation.NewMechanicalFailure("idempotency_lookup_key_client_key_empty")
	}
	maxKey := apivalid.DefaultLimits().MaxLabelValueChars
	if maxKey > 0 && utf8.RuneCountInString(clientIdempotencyKey) > maxKey {
		return IdempotencyLookupKey{}, operation.NewMechanicalFailure("idempotency_lookup_key_client_key_too_long")
	}
	return IdempotencyLookupKey{
		sealed:               true,
		actorRef:             actorRef,
		operation:            operationName,
		clientIdempotencyKey: clientIdempotencyKey,
	}, operation.MechanicalFailure{}
}

// Sealed reports whether the key was constructed through NewIdempotencyLookupKey.
func (k IdempotencyLookupKey) Sealed() bool { return k.sealed }

// ActorRef returns the sealed actor principal.
func (k IdempotencyLookupKey) ActorRef() model.PrincipalRef { return k.actorRef }

// Operation returns the sealed operation name.
func (k IdempotencyLookupKey) Operation() string { return k.operation }

// ClientIdempotencyKey returns the sealed client idempotency key.
func (k IdempotencyLookupKey) ClientIdempotencyKey() string { return k.clientIdempotencyKey }

// Equal reports exact sealed field equality.
func (k IdempotencyLookupKey) Equal(other IdempotencyLookupKey) bool {
	if !k.sealed || !other.sealed {
		return false
	}
	return k.actorRef.Equal(other.actorRef) &&
		k.operation == other.operation &&
		k.clientIdempotencyKey == other.clientIdempotencyKey
}
