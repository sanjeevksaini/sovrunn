package apiconform

import (
	"errors"
	"fmt"
)

// ErrStructuralConfig is returned when StructuralValidatorConfig construction
// rejects a missing dependency (nil SchemaRegistry or nil RefResolver).
var ErrStructuralConfig = errors.New("apiconform: structural validator config invalid")

// StructuralValidatorConfig is an immutable configuration value that binds a
// SchemaRegistry and RefResolver for the StructuralValidator adapter
// (D-01a, D-04; F12-VALIDATION-001(4), F12-VALIDATION-004).
//
// It holds no process-global mutable state. Each value is fixed after
// construction; task 8.2 consumes this configuration explicitly. Registry,
// configuration, or ref-resolution failures ultimately cause the adapter to
// return an error so the pipeline fails closed at LayerStructural.
type StructuralValidatorConfig struct {
	registry SchemaRegistry
	resolver RefResolver
}

// NewStructuralValidatorConfig builds an immutable StructuralValidatorConfig.
// A nil SchemaRegistry or nil RefResolver is rejected.
func NewStructuralValidatorConfig(registry SchemaRegistry, resolver RefResolver) (StructuralValidatorConfig, error) {
	if registry == nil {
		return StructuralValidatorConfig{}, fmt.Errorf("%w: nil schema registry", ErrStructuralConfig)
	}
	if resolver == nil {
		return StructuralValidatorConfig{}, fmt.Errorf("%w: nil ref resolver", ErrStructuralConfig)
	}
	return StructuralValidatorConfig{
		registry: registry,
		resolver: resolver,
	}, nil
}

// Registry returns the configured SchemaRegistry.
func (c StructuralValidatorConfig) Registry() SchemaRegistry {
	return c.registry
}

// RefResolver returns the configured RefResolver.
func (c StructuralValidatorConfig) RefResolver() RefResolver {
	return c.resolver
}
