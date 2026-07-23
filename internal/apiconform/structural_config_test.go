package apiconform

import (
	"errors"
	"testing"
)

func TestNewStructuralValidatorConfig_ValidDependencies(t *testing.T) {
	t.Parallel()

	reg, err := NewMemorySchemaRegistry(map[string][]byte{
		"api/schemas/_common/type-meta.json": []byte(`{"type":"object"}`),
	})
	if err != nil {
		t.Fatalf("NewMemorySchemaRegistry: %v", err)
	}
	resolver, err := NewLocalRefResolver(reg, DefaultMaxRefDepth)
	if err != nil {
		t.Fatalf("NewLocalRefResolver: %v", err)
	}

	cfg, err := NewStructuralValidatorConfig(reg, resolver)
	if err != nil {
		t.Fatalf("NewStructuralValidatorConfig: %v", err)
	}
	if cfg.Registry() != reg {
		t.Fatalf("Registry() = %v, want constructed registry", cfg.Registry())
	}
	if cfg.RefResolver() != resolver {
		t.Fatalf("RefResolver() = %v, want constructed resolver", cfg.RefResolver())
	}
}

func TestNewStructuralValidatorConfig_NilRegistryRejected(t *testing.T) {
	t.Parallel()

	reg, err := NewMemorySchemaRegistry(map[string][]byte{
		"api/schemas/_common/type-meta.json": []byte(`{"type":"object"}`),
	})
	if err != nil {
		t.Fatalf("NewMemorySchemaRegistry: %v", err)
	}
	resolver, err := NewLocalRefResolver(reg, DefaultMaxRefDepth)
	if err != nil {
		t.Fatalf("NewLocalRefResolver: %v", err)
	}

	_, err = NewStructuralValidatorConfig(nil, resolver)
	if !errors.Is(err, ErrStructuralConfig) {
		t.Fatalf("nil registry: err = %v, want ErrStructuralConfig", err)
	}
}

func TestNewStructuralValidatorConfig_NilResolverRejected(t *testing.T) {
	t.Parallel()

	reg, err := NewMemorySchemaRegistry(map[string][]byte{
		"api/schemas/_common/type-meta.json": []byte(`{"type":"object"}`),
	})
	if err != nil {
		t.Fatalf("NewMemorySchemaRegistry: %v", err)
	}

	_, err = NewStructuralValidatorConfig(reg, nil)
	if !errors.Is(err, ErrStructuralConfig) {
		t.Fatalf("nil resolver: err = %v, want ErrStructuralConfig", err)
	}
}

func TestNewStructuralValidatorConfig_NilBothRejected(t *testing.T) {
	t.Parallel()

	_, err := NewStructuralValidatorConfig(nil, nil)
	if !errors.Is(err, ErrStructuralConfig) {
		t.Fatalf("nil both: err = %v, want ErrStructuralConfig", err)
	}
}

func TestStructuralValidatorConfig_ZeroValueHasNilDeps(t *testing.T) {
	t.Parallel()

	var cfg StructuralValidatorConfig
	if cfg.Registry() != nil {
		t.Fatalf("zero Registry() = %v, want nil", cfg.Registry())
	}
	if cfg.RefResolver() != nil {
		t.Fatalf("zero RefResolver() = %v, want nil", cfg.RefResolver())
	}
}
