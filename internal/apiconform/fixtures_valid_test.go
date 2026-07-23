package apiconform

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/apivalid"
)

// validFixtureCase is one Task 14.1 positive fixture: ModeReadRepresentation
// decode into the conformance contract type, then StructuralValidator against
// the canonical schema. Operation variants also assert D-17 scope shape.
type validFixtureCase struct {
	file     string
	schemaID string
	newDst   func() any
	// checkScope is optional; used for Operation D-17 scope assertions.
	checkScope func(t *testing.T, dst any)
}

func TestValidJSONFixturesDecodeAndValidate(t *testing.T) {
	t.Parallel()

	root := moduleRoot(t)
	reg, err := NewRepositorySchemaRegistry(filepath.Join(root, CanonicalSchemasDir))
	if err != nil {
		t.Fatalf("NewRepositorySchemaRegistry: %v", err)
	}
	resolver, err := NewLocalRefResolver(reg, DefaultMaxRefDepth)
	if err != nil {
		t.Fatalf("NewLocalRefResolver: %v", err)
	}
	cfg, err := NewStructuralValidatorConfig(reg, resolver)
	if err != nil {
		t.Fatalf("NewStructuralValidatorConfig: %v", err)
	}
	v, err := NewStructuralValidator(cfg)
	if err != nil {
		t.Fatalf("NewStructuralValidator: %v", err)
	}

	lim := apivalid.DefaultLimits()
	pol := apivalid.PolicyFor(apivalid.ModeReadRepresentation)

	cases := []validFixtureCase{
		{
			file:     "project.json",
			schemaID: CanonicalSchemasDir + "/project.json",
			newDst:   func() any { return &Project{} },
		},
		{
			file:     "resource-pool.json",
			schemaID: CanonicalSchemasDir + "/resource-pool.json",
			newDst:   func() any { return &ResourcePool{} },
		},
		{
			file:     "discovered-database.json",
			schemaID: CanonicalSchemasDir + "/discovered-database.json",
			newDst:   func() any { return &DiscoveredDatabase{} },
		},
		{
			file:     "plugin-definition.json",
			schemaID: CanonicalSchemasDir + "/plugin-definition.json",
			newDst:   func() any { return &PluginDefinition{} },
		},
		{
			file:     "adapter-configuration.json",
			schemaID: CanonicalSchemasDir + "/adapter-configuration.json",
			newDst:   func() any { return &AdapterConfiguration{} },
		},
		{
			file:     "placement-evaluation-request.json",
			schemaID: CanonicalSchemasDir + "/placement-evaluation-request.json",
			newDst:   func() any { return &PlacementEvaluationRequest{} },
		},
		{
			file:     "audit-event.json",
			schemaID: CanonicalSchemasDir + "/audit-event.json",
			newDst:   func() any { return &AuditEvent{} },
		},
		{
			file:       "operation.json",
			schemaID:   CanonicalSchemasDir + "/operation.json",
			newDst:     func() any { return &Operation{} },
			checkScope: assertOperationPlatformScope,
		},
		{
			file:       "operation-platform.json",
			schemaID:   CanonicalSchemasDir + "/operation.json",
			newDst:     func() any { return &Operation{} },
			checkScope: assertOperationPlatformScope,
		},
		{
			file:       "operation-organization.json",
			schemaID:   CanonicalSchemasDir + "/operation.json",
			newDst:     func() any { return &Operation{} },
			checkScope: assertOperationScopeKind(apimeta.ScopeOrganization),
		},
		{
			file:       "operation-organizationunit.json",
			schemaID:   CanonicalSchemasDir + "/operation.json",
			newDst:     func() any { return &Operation{} },
			checkScope: assertOperationScopeKind(apimeta.ScopeOrganizationUnit),
		},
		{
			file:       "operation-tenant.json",
			schemaID:   CanonicalSchemasDir + "/operation.json",
			newDst:     func() any { return &Operation{} },
			checkScope: assertOperationScopeKind(apimeta.ScopeTenant),
		},
		{
			file:       "operation-project.json",
			schemaID:   CanonicalSchemasDir + "/operation.json",
			newDst:     func() any { return &Operation{} },
			checkScope: assertOperationScopeKind(apimeta.ScopeProject),
		},
		{
			file:       "operation-provider.json",
			schemaID:   CanonicalSchemasDir + "/operation.json",
			newDst:     func() any { return &Operation{} },
			checkScope: assertOperationScopeKind(apimeta.ScopeProvider),
		},
	}

	if len(cases) != 14 {
		t.Fatalf("expected 14 Task 14.1 fixtures (7 non-Operation + 7 Operation family), got %d", len(cases))
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.file, func(t *testing.T) {
			t.Parallel()

			raw, err := os.ReadFile(filepath.Join(root, ConformanceFixturesDir, tc.file))
			if err != nil {
				t.Fatalf("read fixture: %v", err)
			}
			dst := tc.newDst()
			if prob := apivalid.DecodeJSON(raw, lim, pol, dst); prob != nil {
				t.Fatalf("DecodeJSON ModeReadRepresentation: code=%s detail=%s violations=%v",
					prob.Code, prob.Detail, prob.Violations)
			}
			violations, err := v.Validate(dst, tc.schemaID)
			if err != nil {
				t.Fatalf("Validate(%s): %v", tc.schemaID, err)
			}
			if len(violations) != 0 {
				t.Fatalf("unexpected structural violations for %s: %#v", tc.file, violations)
			}
			if tc.checkScope != nil {
				tc.checkScope(t, dst)
			}
		})
	}
}

func assertOperationPlatformScope(t *testing.T, dst any) {
	t.Helper()
	op, ok := dst.(*Operation)
	if !ok {
		t.Fatalf("dst type %T, want *Operation", dst)
	}
	if op.Metadata.ScopeRef != nil {
		t.Fatalf("Platform Operation must use canonical nil scopeRef, got %+v", op.Metadata.ScopeRef)
	}
	id := apimeta.CanonicalScopeIdentity(op.Metadata.ScopeRef)
	if id.Kind != apimeta.ScopePlatform || id.UID != apimeta.PlatformScopeUID {
		t.Fatalf("CanonicalScopeIdentity=%+v, want Platform/%q", id, apimeta.PlatformScopeUID)
	}
	if op.Spec.TargetRef.APIVersion == "" || op.Spec.TargetRef.Kind == "" || op.Spec.TargetRef.Name == "" {
		t.Fatalf("Platform Operation targetRef incomplete: %+v", op.Spec.TargetRef)
	}
}

func assertOperationScopeKind(want apimeta.ScopeKind) func(t *testing.T, dst any) {
	return func(t *testing.T, dst any) {
		t.Helper()
		op, ok := dst.(*Operation)
		if !ok {
			t.Fatalf("dst type %T, want *Operation", dst)
		}
		if op.Metadata.ScopeRef == nil {
			t.Fatalf("%s Operation must carry non-nil scopeRef matching target governance scope", want)
		}
		if apimeta.ScopeKind(op.Metadata.ScopeRef.Kind) != want {
			t.Fatalf("scopeRef.kind=%q, want %q", op.Metadata.ScopeRef.Kind, want)
		}
		if op.Metadata.ScopeRef.UID == "" {
			t.Fatalf("%s Operation scopeRef.uid must be non-empty", want)
		}
		id := apimeta.CanonicalScopeIdentity(op.Metadata.ScopeRef)
		if id.Kind != want || id.UID != op.Metadata.ScopeRef.UID {
			t.Fatalf("CanonicalScopeIdentity=%+v, want {%s %s}", id, want, op.Metadata.ScopeRef.UID)
		}
		// D-17 fixture invariant for these variants: targetRef uid equals
		// scopeRef uid when the target is the scope object itself, or for
		// Provider the scopeRef identifies the target's Provider governance
		// scope (target may be a different kind under that Provider).
		if want != apimeta.ScopeProvider && op.Spec.TargetRef.UID != "" &&
			op.Spec.TargetRef.UID != op.Metadata.ScopeRef.UID {
			t.Fatalf("D-17: targetRef.uid %q must match scopeRef.uid %q for %s fixture",
				op.Spec.TargetRef.UID, op.Metadata.ScopeRef.UID, want)
		}
		if want == apimeta.ScopeProvider && op.Metadata.ScopeRef.UID == "" {
			t.Fatalf("Provider Operation scopeRef.uid must identify Provider governance scope")
		}
	}
}
