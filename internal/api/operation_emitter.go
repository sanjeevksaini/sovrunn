package api

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"log"
	"time"

	"github.com/sanjeevksaini/sovrunn/internal/registry"
	"github.com/sanjeevksaini/sovrunn/internal/requestctx"
	"github.com/sanjeevksaini/sovrunn/internal/resources"
)

var (
	errOperationIDExhausted        = errors.New("operation id generation exhausted")
	errOperationEmitterUnavailable = errors.New("operation emitter unavailable")
)

// OperationEmitter records a control-plane Operation after a successful
// mutating action.
type OperationEmitter interface {
	Emit(ctx context.Context, spec resources.OperationSpec) error
}

// newOperationID returns a URL-safe, path-segment-safe opaque token using
// crypto/rand: 16 random bytes encoded as 32 lowercase hex characters.
func newOperationID() (string, error) {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", err
	}
	return hex.EncodeToString(b[:]), nil
}

// registryEmitter adapts the storage-only OperationRegistry to the
// OperationEmitter interface. logger is optional and may be nil.
type registryEmitter struct {
	registry registry.OperationRegistryIface
	logger   *log.Logger
}

// NewRegistryEmitter constructs a registry-backed OperationEmitter.
func NewRegistryEmitter(reg registry.OperationRegistryIface, logger *log.Logger) *registryEmitter {
	return &registryEmitter{registry: reg, logger: logger}
}

// Emit generates an Operation ID, forces server-owned fields, and stores the
// Operation. ID collisions are retried up to five times.
func (e *registryEmitter) Emit(ctx context.Context, spec resources.OperationSpec) error {
	if e == nil || e.registry == nil {
		return errOperationEmitterUnavailable
	}

	spec.Actor = "system"
	now := time.Now().UTC().Format(time.RFC3339)

	for attempt := 0; attempt < 5; attempt++ {
		id, err := newOperationID()
		if err != nil {
			return err
		}

		op := resources.Operation{
			APIVersion: resources.OperationAPIVersion,
			Kind:       resources.OperationKind,
			Metadata:   resources.Metadata{Name: id},
			Spec:       spec,
			Status: resources.OperationStatus{
				Phase:       resources.OperationPhaseSucceeded,
				CreatedAt:   now,
				UpdatedAt:   now,
				CompletedAt: now,
			},
		}

		if _, err := e.registry.CreateOperation(ctx, op); err != nil {
			if errors.Is(err, registry.ErrAlreadyExists) {
				continue
			}
			return err
		}
		return nil
	}

	return errOperationIDExhausted
}

// emitOperation records an Operation but never affects the primary handler
// response. It is nil-safe and swallows emitter errors.
func emitOperation(ctx context.Context, emitter OperationEmitter, spec resources.OperationSpec) {
	if emitter == nil {
		return
	}
	_ = emitter.Emit(ctx, spec)
}

// requestIDFromContext returns the request ID stored by server middleware via
// the shared requestctx package. It returns an empty string when absent.
func requestIDFromContext(ctx context.Context) string {
	return requestctx.RequestIDFromContext(ctx)
}
