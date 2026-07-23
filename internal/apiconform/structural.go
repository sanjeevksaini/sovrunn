package apiconform

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/apischema"
	"github.com/sanjeevksaini/sovrunn/internal/apivalid"
)

// ErrStructuralValidator is returned when the StructuralValidator adapter
// cannot run (nil receiver, misconfigured registry/resolver, missing schema,
// or $ref resolution failure). Callers map this to LayerStructural fail-closed
// 500 INTERNAL_ERROR behavior (D-01a, D-04).
var ErrStructuralValidator = errors.New("apiconform: structural validator unavailable")

// StructuralValidator adapts SchemaRegistry + RefResolver + apischema into
// the apivalid.StructuralValidator interface (D-01a, D-02, D-04;
// F12-VALIDATION-001(4), F12-VALIDATION-006).
//
// Validate loads a schema by ID, resolves approved local $ref values via the
// configured RefResolver, runs ValidateSchemaSupport and ValidateInstance,
// and translates SchemaIssue findings to apiproblem.Violation values.
// Registry, configuration, and ref-resolution failures return an error so the
// pipeline fails closed at LayerStructural.
type StructuralValidator struct {
	cfg StructuralValidatorConfig
}

// NewStructuralValidator builds a StructuralValidator from an explicit
// StructuralValidatorConfig (task 7a.3). A config with a nil SchemaRegistry
// or nil RefResolver is rejected.
func NewStructuralValidator(cfg StructuralValidatorConfig) (*StructuralValidator, error) {
	if cfg.Registry() == nil {
		return nil, fmt.Errorf("%w: nil schema registry", ErrStructuralValidator)
	}
	if cfg.RefResolver() == nil {
		return nil, fmt.Errorf("%w: nil ref resolver", ErrStructuralValidator)
	}
	return &StructuralValidator{cfg: cfg}, nil
}

// Validate implements apivalid.StructuralValidator.
//
// Outcomes:
//   - err != nil: structural validation unavailable (missing schema, nil/misconfigured
//     registry or resolver, or $ref resolution failure)
//   - err == nil, len(violations) > 0: ordinary schema / support findings (422)
//   - err == nil, len(violations) == 0: instance is structurally valid
func (v *StructuralValidator) Validate(instance any, schemaID string) ([]apiproblem.Violation, error) {
	if v == nil {
		return nil, fmt.Errorf("%w: nil structural validator", ErrStructuralValidator)
	}
	registry := v.cfg.Registry()
	resolver := v.cfg.RefResolver()
	if registry == nil {
		return nil, fmt.Errorf("%w: nil schema registry", ErrStructuralValidator)
	}
	if resolver == nil {
		return nil, fmt.Errorf("%w: nil ref resolver", ErrStructuralValidator)
	}

	raw, err := registry.Load(schemaID)
	if err != nil {
		return nil, fmt.Errorf("%w: load %q: %w", ErrStructuralValidator, schemaID, err)
	}

	resolved, err := resolveSchemaRefs(resolver, schemaID, raw)
	if err != nil {
		return nil, fmt.Errorf("%w: resolve refs in %q: %w", ErrStructuralValidator, schemaID, err)
	}

	if support := apischema.ValidateSchemaSupport(resolved); len(support) > 0 {
		return SchemaIssuesToViolations(support), nil
	}

	issues := apischema.ValidateInstance(resolved, instance)
	return SchemaIssuesToViolations(issues), nil
}

// resolveSchemaRefs inlines all local $ref values in a schema document using
// the configured RefResolver. Any Resolve failure fails closed (returned as
// error). Sibling keywords alongside $ref are discarded; FEATURE-0012
// ValidateInstance rejects unresolved $ref nodes, so inlining must produce a
// fully expanded schema tree.
func resolveSchemaRefs(resolver RefResolver, schemaID string, raw []byte) ([]byte, error) {
	if resolver == nil {
		return nil, fmt.Errorf("%w: nil ref resolver", ErrStructuralValidator)
	}
	var root any
	if err := json.Unmarshal(raw, &root); err != nil {
		return nil, fmt.Errorf("%w: schema %q is not valid JSON: %v", ErrStructuralValidator, schemaID, err)
	}
	inlined, err := inlineRefs(resolver, schemaID, root, []string{schemaID})
	if err != nil {
		return nil, err
	}
	out, err := json.Marshal(inlined)
	if err != nil {
		return nil, fmt.Errorf("%w: marshal resolved schema %q: %v", ErrStructuralValidator, schemaID, err)
	}
	return out, nil
}

func inlineRefs(resolver RefResolver, baseID string, node any, stack []string) (any, error) {
	switch n := node.(type) {
	case map[string]any:
		if refVal, hasRef := n["$ref"]; hasRef {
			ref, ok := refVal.(string)
			if !ok {
				return nil, fmt.Errorf("%w: non-string $ref in schema %q", ErrStructuralValidator, baseID)
			}
			targetID, targetBody, err := resolver.Resolve(baseID, ref)
			if err != nil {
				return nil, err
			}
			if len(stack) >= DefaultMaxRefDepth {
				return nil, fmt.Errorf("%w: maxDepth=%d at %q", ErrRefDepthExceeded, DefaultMaxRefDepth, targetID)
			}
			for _, seen := range stack {
				if seen == targetID {
					return nil, fmt.Errorf("%w: %q", ErrRefCycle, targetID)
				}
			}
			var target any
			if err := json.Unmarshal(targetBody, &target); err != nil {
				return nil, fmt.Errorf("%w: target schema %q is not valid JSON: %v", ErrStructuralValidator, targetID, err)
			}
			next := append(append([]string(nil), stack...), targetID)
			return inlineRefs(resolver, targetID, target, next)
		}

		out := make(map[string]any, len(n))
		for k, child := range n {
			inlined, err := inlineRefs(resolver, baseID, child, stack)
			if err != nil {
				return nil, err
			}
			out[k] = inlined
		}
		return out, nil

	case []any:
		out := make([]any, len(n))
		for i, child := range n {
			inlined, err := inlineRefs(resolver, baseID, child, stack)
			if err != nil {
				return nil, err
			}
			out[i] = inlined
		}
		return out, nil

	default:
		return node, nil
	}
}

var _ apivalid.StructuralValidator = (*StructuralValidator)(nil)

// SchemaIssuesToViolations translates package-local apischema.SchemaIssue
// values into apiproblem.Violation values (F12-VALIDATION-006, D-01a, D-02).
//
// Mapping is field-for-field:
//
//	SchemaIssue.Path    → Violation.Field  (RFC 6901 JSON Pointer)
//	SchemaIssue.Code    → Violation.Code   (stable machine contract)
//	SchemaIssue.Message → Violation.Message (informational; must not carry secrets)
//
// The returned slice is newly allocated so callers cannot mutate the input
// through the result. A nil or empty input yields nil.
//
// Translation lives in apiconform (not apischema) so apischema never imports
// apiproblem, preserving the D-02 import-direction boundary.
func SchemaIssuesToViolations(issues []apischema.SchemaIssue) []apiproblem.Violation {
	if len(issues) == 0 {
		return nil
	}
	out := make([]apiproblem.Violation, len(issues))
	for i, issue := range issues {
		out[i] = apiproblem.Violation{
			Field:   issue.Path,
			Code:    apiproblem.ViolationCode(issue.Code),
			Message: issue.Message,
		}
	}
	return out
}
