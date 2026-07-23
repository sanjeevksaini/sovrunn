package apiconform

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/apicond"
	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/apiref"
	"github.com/sanjeevksaini/sovrunn/internal/apischema"
	"github.com/sanjeevksaini/sovrunn/internal/apivalid"
)

// Task 14.5 conformance suite (F12-FIXTURE-002, F12-VALIDATION-002,
// F12-SEC-003/004, F12-SCOPE-002; D-03a, D-15, D-16, D-17).

const property4CompleteSeed int64 = 20260723145
const property4CompleteIterations = 100

// conformanceFixtureCase binds one canonical fixture family to its schema and
// Go contract type for positive decode/validate/annotation proofs.
type conformanceFixtureCase struct {
	base     string // fixture basename without extension
	schema   string // file under api/schemas/
	schemaID string
	newDst   func() any
}

func conformanceFixtureCases() []conformanceFixtureCase {
	return []conformanceFixtureCase{
		{base: "project", schema: "project.json", schemaID: CanonicalSchemasDir + "/project.json", newDst: func() any { return &Project{} }},
		{base: "resource-pool", schema: "resource-pool.json", schemaID: CanonicalSchemasDir + "/resource-pool.json", newDst: func() any { return &ResourcePool{} }},
		{base: "discovered-database", schema: "discovered-database.json", schemaID: CanonicalSchemasDir + "/discovered-database.json", newDst: func() any { return &DiscoveredDatabase{} }},
		{base: "plugin-definition", schema: "plugin-definition.json", schemaID: CanonicalSchemasDir + "/plugin-definition.json", newDst: func() any { return &PluginDefinition{} }},
		{base: "adapter-configuration", schema: "adapter-configuration.json", schemaID: CanonicalSchemasDir + "/adapter-configuration.json", newDst: func() any { return &AdapterConfiguration{} }},
		{base: "placement-evaluation-request", schema: "placement-evaluation-request.json", schemaID: CanonicalSchemasDir + "/placement-evaluation-request.json", newDst: func() any { return &PlacementEvaluationRequest{} }},
		{base: "audit-event", schema: "audit-event.json", schemaID: CanonicalSchemasDir + "/audit-event.json", newDst: func() any { return &AuditEvent{} }},
		{base: "operation", schema: "operation.json", schemaID: CanonicalSchemasDir + "/operation.json", newDst: func() any { return &Operation{} }},
		{base: "operation-platform", schema: "operation.json", schemaID: CanonicalSchemasDir + "/operation.json", newDst: func() any { return &Operation{} }},
		{base: "operation-organization", schema: "operation.json", schemaID: CanonicalSchemasDir + "/operation.json", newDst: func() any { return &Operation{} }},
		{base: "operation-organizationunit", schema: "operation.json", schemaID: CanonicalSchemasDir + "/operation.json", newDst: func() any { return &Operation{} }},
		{base: "operation-tenant", schema: "operation.json", schemaID: CanonicalSchemasDir + "/operation.json", newDst: func() any { return &Operation{} }},
		{base: "operation-project", schema: "operation.json", schemaID: CanonicalSchemasDir + "/operation.json", newDst: func() any { return &Operation{} }},
		{base: "operation-provider", schema: "operation.json", schemaID: CanonicalSchemasDir + "/operation.json", newDst: func() any { return &Operation{} }},
	}
}

func TestConformancePositiveFixturesDecodeValidateAnnotations(t *testing.T) {
	t.Parallel()

	root := moduleRoot(t)
	v := mustCanonicalStructuralValidator(t, root)
	lim := apivalid.DefaultLimits()
	pol := apivalid.PolicyFor(apivalid.ModeReadRepresentation)
	dir := filepath.Join(root, ConformanceFixturesDir)
	schemasDir := filepath.Join(root, CanonicalSchemasDir)

	cases := conformanceFixtureCases()
	if len(cases) != 14 {
		t.Fatalf("expected 14 fixtures (7 non-Operation + 7 Operation family), got %d", len(cases))
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.base, func(t *testing.T) {
			t.Parallel()

			schemaRaw, err := os.ReadFile(filepath.Join(schemasDir, tc.schema))
			if err != nil {
				t.Fatalf("read schema: %v", err)
			}
			meta, issues := apischema.ReadAnnotations(schemaRaw)
			if len(issues) != 0 {
				t.Fatalf("ReadAnnotations(%s): %#v", tc.schema, issues)
			}
			if meta.Profile == "" || meta.Boundary == "" || meta.Stability == "" || len(meta.AllowedScopes) == 0 {
				t.Fatalf("schema annotations incomplete: %#v", meta)
			}

			jsonRaw, err := os.ReadFile(filepath.Join(dir, tc.base+".json"))
			if err != nil {
				t.Fatalf("read JSON: %v", err)
			}
			yamlRaw, err := os.ReadFile(filepath.Join(dir, tc.base+".yaml"))
			if err != nil {
				t.Fatalf("read YAML: %v", err)
			}

			fromJSON := tc.newDst()
			if prob := apivalid.DecodeJSON(jsonRaw, lim, pol, fromJSON); prob != nil {
				t.Fatalf("DecodeJSON: code=%s detail=%s violations=%v", prob.Code, prob.Detail, prob.Violations)
			}
			violations, err := v.Validate(fromJSON, tc.schemaID)
			if err != nil {
				t.Fatalf("Validate JSON: %v", err)
			}
			if len(violations) != 0 {
				t.Fatalf("structural violations (JSON): %#v", violations)
			}

			fromYAML := tc.newDst()
			if prob := apivalid.DecodeYAML(yamlRaw, lim, pol, fromYAML); prob != nil {
				t.Fatalf("DecodeYAML: code=%s detail=%s violations=%v", prob.Code, prob.Detail, prob.Violations)
			}
			violations, err = v.Validate(fromYAML, tc.schemaID)
			if err != nil {
				t.Fatalf("Validate YAML: %v", err)
			}
			if len(violations) != 0 {
				t.Fatalf("structural violations (YAML): %#v", violations)
			}
			if !reflect.DeepEqual(fromJSON, fromYAML) {
				t.Fatalf("JSON/YAML typed values diverge:\nJSON=%#v\nYAML=%#v", fromJSON, fromYAML)
			}
		})
	}
}

func TestConformanceOperationAwareFieldPolicy(t *testing.T) {
	t.Parallel()

	root := moduleRoot(t)
	raw, err := os.ReadFile(filepath.Join(root, ConformanceFixturesDir, "project.json"))
	if err != nil {
		t.Fatalf("read project.json: %v", err)
	}
	lim := apivalid.DefaultLimits()

	type decodeDst struct {
		APIVersion string         `json:"apiVersion"`
		Kind       string         `json:"kind"`
		Metadata   map[string]any `json:"metadata"`
		Spec       map[string]any `json:"spec"`
		Status     map[string]any `json:"status,omitempty"`
	}

	// Task 14.5: same object accepted under internal/read; rejected under
	// create/replace when it carries status/system-owned fields (D-15).
	// ModeStatusUpdate is not used here: full fixtures also carry spec, which
	// that mode intentionally rejects (AllowSpecMutation=false).
	acceptModes := []apivalid.DecodeMode{
		apivalid.ModeReadRepresentation,
		apivalid.ModeInternalObject,
	}
	rejectModes := []apivalid.DecodeMode{
		apivalid.ModeCreateRequest,
		apivalid.ModeReplaceRequest,
	}

	for _, mode := range acceptModes {
		mode := mode
		t.Run("accept/"+mode.String(), func(t *testing.T) {
			t.Parallel()
			var dst decodeDst
			if prob := apivalid.DecodeJSON(raw, lim, apivalid.PolicyFor(mode), &dst); prob != nil {
				t.Fatalf("mode %s must accept status/system fields: %#v", mode, prob)
			}
			if dst.Status == nil || dst.Status["phase"] == nil {
				t.Fatal("expected status.phase after accept-mode decode")
			}
		})
	}

	for _, mode := range rejectModes {
		mode := mode
		t.Run("reject/"+mode.String(), func(t *testing.T) {
			t.Parallel()
			var dst decodeDst
			prob := apivalid.DecodeJSON(raw, lim, apivalid.PolicyFor(mode), &dst)
			if prob == nil {
				t.Fatalf("mode %s must reject status/system-owned fields on project fixture", mode)
			}
			if prob.Code != apiproblem.CodeValidationFailed {
				t.Fatalf("code=%q, want VALIDATION_FAILED", prob.Code)
			}
			if len(prob.Violations) == 0 {
				t.Fatalf("expected field violations, got %#v", prob)
			}
			field := prob.Violations[0].Field
			if !strings.HasPrefix(field, "/status") && !strings.HasPrefix(field, "/metadata/") {
				t.Fatalf("unexpected rejection field %q", field)
			}
		})
	}
}

// Feature: api-resource-naming-status-and-validation-standard, Property 4 (complete):
// Canonical platform scope with schema x-sovrunn-allowed-scopes.
//
// For any object whose schema allows Platform, absent and explicit Platform
// scopeRef normalize to the identical canonical form and share
// CanonicalScopeIdentity (Platform/PlatformScopeUID). For any object whose
// schema does not allow Platform, nil/normalized platform scope is rejected.
//
// Validates: Requirements 4.4, 4.5 (F12-SCOPE-002, F12-REF-001; D-16)
func TestProperty4_CanonicalPlatformScopeComplete(t *testing.T) {
	t.Parallel()

	root := moduleRoot(t)
	platformAllowed := mustSchemaAllowedScopes(t, root, "plugin-definition.json")
	platformDisallowed := mustSchemaAllowedScopes(t, root, "project.json")

	if !containsScopeKind(platformAllowed, apimeta.ScopePlatform) {
		t.Fatalf("plugin-definition allowed-scopes must include Platform, got %v", platformAllowed)
	}
	if containsScopeKind(platformDisallowed, apimeta.ScopePlatform) {
		t.Fatalf("project allowed-scopes must NOT include Platform, got %v", platformDisallowed)
	}

	rng := rand.New(rand.NewSource(property4CompleteSeed))
	for i := 0; i < property4CompleteIterations; i++ {
		if err := checkProperty4CompleteIteration(rng, i, platformAllowed, platformDisallowed); err != nil {
			t.Fatalf("property 4 (complete) failed at iteration %d (seed %d): %v",
				i, property4CompleteSeed, err)
		}
	}
}

func checkProperty4CompleteIteration(
	rng *rand.Rand,
	iteration int,
	platformAllowed, platformDisallowed []apimeta.ScopeKind,
) error {
	explicit := &apimeta.ScopeRef{TypedRef: apimeta.TypedRef{
		APIVersion: "platform.sovrunn.io/v1alpha1",
		Kind:       string(apimeta.ScopePlatform),
		Name:       fmt.Sprintf("platform-%d", rng.Intn(10000)),
		UID:        apimeta.PlatformScopeUID,
	}}
	normalizedExplicit := apimeta.NormalizeScope(explicit)
	normalizedNil := apimeta.NormalizeScope(nil)

	if normalizedExplicit != nil {
		return fmt.Errorf("NormalizeScope(explicit Platform) = %#v, want nil", normalizedExplicit)
	}
	if normalizedNil != nil {
		return fmt.Errorf("NormalizeScope(nil) = %#v, want nil", normalizedNil)
	}
	if apimeta.NormalizeScope(normalizedExplicit) != nil {
		return fmt.Errorf("NormalizeScope must be idempotent on nil")
	}

	idNil := apimeta.CanonicalScopeIdentity(nil)
	idExplicit := apimeta.CanonicalScopeIdentity(explicit)
	idNormalized := apimeta.CanonicalScopeIdentity(normalizedExplicit)
	want := apimeta.ScopeIdentity{Kind: apimeta.ScopePlatform, UID: apimeta.PlatformScopeUID}
	if idNil != want || idExplicit != want || idNormalized != want {
		return fmt.Errorf("identity mismatch: nil=%#v explicit=%#v normalized=%#v want=%#v",
			idNil, idExplicit, idNormalized, want)
	}
	if idNil != idExplicit || idNil != idNormalized {
		return fmt.Errorf("platform forms must share identity tuple")
	}

	// Schema annotation: Platform allowed → nil accepted via CommonReference.
	if err := assertAllowedScopeOutcome(platformAllowed, nil, true, iteration, "platform-allowed/nil"); err != nil {
		return err
	}
	if err := assertAllowedScopeOutcome(platformAllowed, normalizedExplicit, true, iteration, "platform-allowed/normalized"); err != nil {
		return err
	}

	// Schema annotation: Platform not allowed → nil rejected.
	if err := assertAllowedScopeOutcome(platformDisallowed, nil, false, iteration, "platform-disallowed/nil"); err != nil {
		return err
	}
	if err := assertAllowedScopeOutcome(platformDisallowed, normalizedExplicit, false, iteration, "platform-disallowed/normalized"); err != nil {
		return err
	}

	// Non-platform scope still accepted under project (Tenant) annotation when matching.
	tenantUID := fmt.Sprintf("%032x", rng.Uint64()^uint64(iteration))
	tenant := &apimeta.ScopeRef{TypedRef: apimeta.TypedRef{
		APIVersion: "tenancy.sovrunn.io/v1alpha1",
		Kind:       string(apimeta.ScopeTenant),
		Name:       fmt.Sprintf("tenant-%d", iteration),
		UID:        tenantUID[:32],
	}}
	if got := apimeta.NormalizeScope(tenant); got != tenant {
		return fmt.Errorf("non-platform NormalizeScope mutated pointer")
	}
	if err := assertAllowedScopeOutcome(platformDisallowed, tenant, true, iteration, "tenant-allowed"); err != nil {
		return err
	}
	return nil
}

func assertAllowedScopeOutcome(
	allowed []apimeta.ScopeKind,
	scope *apimeta.ScopeRef,
	wantAccept bool,
	iteration int,
	label string,
) error {
	stage := apivalid.NewCommonReference(apivalid.ReferenceConfig{
		AllowedScopes: append([]apimeta.ScopeKind(nil), allowed...),
		Fields: []apivalid.RefField{{
			Path: "/spec/targetRef",
			Constraint: apiref.Constraint{
				AllowedKinds: []string{"PluginDefinition"},
				Direction:    apiref.DirectionOutbound,
			},
		}},
	}, apivalid.DefaultLimits())

	obj := &conformanceRefCarrier{
		scope: scope,
		singular: map[string]apiref.TypedRef{
			"/spec/targetRef": {
				APIVersion: "plugin.sovrunn.io/v1alpha1",
				Kind:       "PluginDefinition",
				Name:       fmt.Sprintf("demo-%d", iteration),
			},
		},
	}
	violations, err := stage.Validate(context.Background(), obj)
	if err != nil {
		return fmt.Errorf("%s: CommonReference error: %v", label, err)
	}
	hasReject := hasViolationCode(violations, apiproblem.ViolationCode(apiref.CodeScopeNotAllowed))
	if wantAccept && hasReject {
		return fmt.Errorf("%s: expected accept, got REF_SCOPE_NOT_ALLOWED (%#v)", label, violations)
	}
	if !wantAccept && !hasReject {
		return fmt.Errorf("%s: expected REF_SCOPE_NOT_ALLOWED, got %#v", label, violations)
	}
	return nil
}

func TestConformanceOperationTargetScopeEquality(t *testing.T) {
	t.Parallel()

	root := moduleRoot(t)
	loader, err := NewFixtureLoader(root)
	if err != nil {
		t.Fatalf("NewFixtureLoader: %v", err)
	}
	v := mustCanonicalStructuralValidator(t, root)
	lim := apivalid.DefaultLimits()
	readPol := apivalid.PolicyFor(apivalid.ModeReadRepresentation)

	variants := []string{
		"operation.json",
		"operation-platform.json",
		"operation-organization.json",
		"operation-organizationunit.json",
		"operation-tenant.json",
		"operation-project.json",
		"operation-provider.json",
	}
	for _, file := range variants {
		file := file
		t.Run("match/"+file, func(t *testing.T) {
			t.Parallel()
			raw, err := loader.Read(file)
			if err != nil {
				t.Fatalf("read: %v", err)
			}
			var op Operation
			if prob := apivalid.DecodeJSON(raw, lim, readPol, &op); prob != nil {
				t.Fatalf("DecodeJSON: %#v", prob)
			}
			violations, err := v.Validate(&op, CanonicalSchemasDir+"/operation.json")
			if err != nil {
				t.Fatalf("Validate: %v", err)
			}
			if len(violations) != 0 {
				t.Fatalf("structural violations: %#v", violations)
			}
			opScope := apimeta.CanonicalScopeIdentity(op.Metadata.ScopeRef)
			targetScope := resolveOperationTargetGovernanceScope(&op)
			if mismatch := apivalid.CheckOperationTargetScopeMatch(opScope, targetScope); mismatch != nil {
				t.Fatalf("D-17 match failed: %#v (op=%+v target=%+v)", mismatch, opScope, targetScope)
			}
			if file == "operation.json" || file == "operation-platform.json" {
				if op.Metadata.ScopeRef != nil {
					t.Fatalf("Platform fixture must use canonical nil scopeRef, got %+v", op.Metadata.ScopeRef)
				}
				if opScope.UID != apimeta.PlatformScopeUID {
					t.Fatalf("Platform identity UID=%q, want %q", opScope.UID, apimeta.PlatformScopeUID)
				}
			}
		})
	}

	for _, file := range []string{
		"negative/operation-scope-target-kind-mismatch.json",
		"negative/operation-scope-target-uid-mismatch.json",
	} {
		file := file
		t.Run("mismatch/"+file, func(t *testing.T) {
			t.Parallel()
			if err := assertOperationScopeMismatchFixture(loader, file); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestConformanceJSONYAMLEquivalenceCanonical(t *testing.T) {
	t.Parallel()

	root := moduleRoot(t)
	v := mustCanonicalStructuralValidator(t, root)
	lim := apivalid.DefaultLimits()
	pol := apivalid.PolicyFor(apivalid.ModeReadRepresentation)
	dir := filepath.Join(root, ConformanceFixturesDir)

	for _, tc := range conformanceFixtureCases() {
		tc := tc
		t.Run(tc.base, func(t *testing.T) {
			t.Parallel()

			jsonRaw, err := os.ReadFile(filepath.Join(dir, tc.base+".json"))
			if err != nil {
				t.Fatalf("read JSON: %v", err)
			}
			yamlRaw, err := os.ReadFile(filepath.Join(dir, tc.base+".yaml"))
			if err != nil {
				t.Fatalf("read YAML: %v", err)
			}

			fromJSON := tc.newDst()
			if prob := apivalid.DecodeJSON(jsonRaw, lim, pol, fromJSON); prob != nil {
				t.Fatalf("DecodeJSON: %#v", prob)
			}
			fromYAML := tc.newDst()
			if prob := apivalid.DecodeYAML(yamlRaw, lim, pol, fromYAML); prob != nil {
				t.Fatalf("DecodeYAML: %#v", prob)
			}
			if !reflect.DeepEqual(fromJSON, fromYAML) {
				t.Fatalf("typed values diverge")
			}

			vj, err := v.Validate(fromJSON, tc.schemaID)
			if err != nil {
				t.Fatalf("Validate JSON: %v", err)
			}
			vy, err := v.Validate(fromYAML, tc.schemaID)
			if err != nil {
				t.Fatalf("Validate YAML: %v", err)
			}
			if !reflect.DeepEqual(vj, vy) {
				t.Fatalf("structural violation sets diverge:\nJSON=%#v\nYAML=%#v", vj, vy)
			}
			if len(vj) != 0 {
				t.Fatalf("unexpected violations: %#v", vj)
			}
		})
	}
}

func TestConformanceBoundaryLimits(t *testing.T) {
	t.Parallel()

	lim := apivalid.DefaultLimits()

	t.Run("MaxLabels", func(t *testing.T) {
		t.Parallel()
		labels := make(map[string]string, lim.MaxLabels+1)
		for i := 0; i <= lim.MaxLabels; i++ {
			labels[fmt.Sprintf("k%d", i)] = "v"
		}
		carrier := &conformanceSemanticCarrier{
			apiVersion: "core.sovrunn.io/v1alpha1",
			kind:       "Project",
			name:       "payments-production",
			scope: &apimeta.ScopeRef{TypedRef: apimeta.TypedRef{
				APIVersion: "tenancy.sovrunn.io/v1alpha1",
				Kind:       string(apimeta.ScopeTenant),
				Name:       "acme-tenant",
				UID:        "b2c3d4e5f60718293a4b5c6d7e8f901a",
			}},
			labels: labels,
		}
		violations, err := apivalid.NewCommonSemantic(lim, false).Validate(context.Background(), carrier)
		if err != nil {
			t.Fatalf("Validate: %v", err)
		}
		if !hasViolationCode(violations, apiproblem.ViolationOutOfRange) {
			t.Fatalf("expected OUT_OF_RANGE for MaxLabels, got %#v", violations)
		}
		if !hasViolationField(violations, "/metadata/labels") {
			t.Fatalf("expected /metadata/labels, got %#v", violations)
		}
	})

	t.Run("MaxConditions", func(t *testing.T) {
		t.Parallel()
		conds := make([]apicond.Condition, lim.MaxConditions+1)
		for i := range conds {
			conds[i] = apicond.Condition{
				Type:               fmt.Sprintf("Ready%d", i),
				Status:             apicond.ConditionTrue,
				Reason:             "Reconciled",
				LastTransitionTime: "2026-07-01T12:00:00Z",
			}
		}
		carrier := &conformanceSemanticCarrier{
			apiVersion: "core.sovrunn.io/v1alpha1",
			kind:       "Project",
			name:       "payments-production",
			scope: &apimeta.ScopeRef{TypedRef: apimeta.TypedRef{
				APIVersion: "tenancy.sovrunn.io/v1alpha1",
				Kind:       string(apimeta.ScopeTenant),
				Name:       "acme-tenant",
				UID:        "b2c3d4e5f60718293a4b5c6d7e8f901a",
			}},
			conditions: conds,
			phase:      "Active",
		}
		violations, err := apivalid.NewCommonSemantic(lim, false).Validate(context.Background(), carrier)
		if err != nil {
			t.Fatalf("Validate: %v", err)
		}
		if !hasViolationCode(violations, apiproblem.ViolationOutOfRange) {
			t.Fatalf("expected OUT_OF_RANGE for MaxConditions, got %#v", violations)
		}
		if !hasViolationField(violations, "/status/conditions") {
			t.Fatalf("expected /status/conditions, got %#v", violations)
		}
	})

	t.Run("MaxObjectBytes", func(t *testing.T) {
		t.Parallel()
		tiny := apivalid.Limits{MaxObjectBytes: 32, MaxNestingDepth: 32}
		raw, err := os.ReadFile(filepath.Join(moduleRoot(t), ConformanceFixturesDir, "project.json"))
		if err != nil {
			t.Fatalf("read: %v", err)
		}
		var dst Project
		prob := apivalid.DecodeJSON(raw, tiny, apivalid.PolicyFor(apivalid.ModeReadRepresentation), &dst)
		if prob == nil || prob.Code != apiproblem.CodeRequestTooLarge {
			t.Fatalf("expected REQUEST_TOO_LARGE, got %#v", prob)
		}
	})

	t.Run("MaxPageSize", func(t *testing.T) {
		t.Parallel()
		if lim.DefaultPageSize <= 0 || lim.DefaultPageSize > lim.MaxPageSize || lim.MaxPageSize != 200 {
			t.Fatalf("list bounds unexpected: default=%d max=%d", lim.DefaultPageSize, lim.MaxPageSize)
		}
	})
}

func TestConformanceSecuritySecretLikeAndSafeDenial(t *testing.T) {
	t.Parallel()

	root := moduleRoot(t)
	loader, err := NewFixtureLoader(root)
	if err != nil {
		t.Fatalf("NewFixtureLoader: %v", err)
	}

	t.Run("fixtures_reject_secret_like_metadata", func(t *testing.T) {
		t.Parallel()
		for _, tc := range conformanceFixtureCases() {
			raw, err := loader.Read(tc.base + ".json")
			if err != nil {
				t.Fatalf("%s: %v", tc.base, err)
			}
			if hit := findSecretLikeInMetadataJSON(raw); hit != "" {
				t.Fatalf("%s embeds secret-like metadata token %q", tc.base, hit)
			}
		}
		// Injected secret-like label must be detected by the same scanner.
		poison := []byte(`{"metadata":{"labels":{"password":"x"},"annotations":{}}}`)
		if hit := findSecretLikeInMetadataJSON(poison); hit == "" {
			t.Fatal("scanner must reject injected secret-like label key")
		}
	})

	t.Run("safe_denial_path_response_equivalence", func(t *testing.T) {
		t.Parallel()

		want := apivalid.SafeDenial(apivalid.DenyNotDisclosed)
		absent := apiproblem.New(apiproblem.CodeResourceNotFound)
		wantJSON, err := json.Marshal(want)
		if err != nil {
			t.Fatalf("marshal SafeDenial: %v", err)
		}
		absentJSON, err := json.Marshal(absent)
		if err != nil {
			t.Fatalf("marshal absent: %v", err)
		}
		if !bytes.Equal(wantJSON, absentJSON) {
			t.Fatalf("SafeDenial not byte-identical to absent 404:\ndenied=%s\nabsent=%s", wantJSON, absentJSON)
		}

		var project Project
		if err := loader.DecodeJSON("project.json", &project); err != nil {
			t.Fatalf("DecodeJSON project: %v", err)
		}
		opScope := apimeta.CanonicalScopeIdentity(project.Metadata.ScopeRef)
		foreign := apimeta.ScopeIdentity{Kind: apimeta.ScopeTenant, UID: "ffffffffffffffffffffffffffffffff"}
		target := apiref.TypedRef{
			APIVersion: project.APIVersion,
			Kind:       project.Kind,
			Name:       project.Metadata.Name,
			UID:        project.Metadata.UID,
		}
		caller := apivalid.CallerContext{Scopes: []apimeta.ScopeIdentity{opScope}}
		// Contract types are not SemanticCarrier/ReferenceCarrier; use
		// deterministic no-op stages so layer 8 remains the focus.
		stages := apivalid.StageSet{
			Defaulting: &conformanceNoopDefaulting{},
			Semantic:   &conformanceNoopValidation{},
			Reference:  &conformanceNoopValidation{},
		}

		authCalls := &atomic.Int32{}
		resA := apivalid.Validate(context.Background(), apivalid.Input{
			Validator:      mustCanonicalStructuralValidator(t, root),
			SchemaID:       CanonicalSchemasDir + "/project.json",
			Dst:            &project,
			Mode:           apivalid.ModeReadRepresentation,
			Stages:         stages,
			OperationScope: &opScope,
			TargetRef:      &target,
			TargetScope:    &foreign,
			Authorizer:     conformanceScopeAuthorizer{decision: apivalid.DenyNotDisclosed, calls: authCalls},
			Caller:         &caller,
		}, apivalid.DefaultLimits())
		if authCalls.Load() != 1 {
			t.Fatalf("Path A Authorizer calls=%d, want 1 (authorize-before-lookup)", authCalls.Load())
		}
		assertSafeDenialResult(t, resA, want, "path A")

		resolverCalls := &atomic.Int32{}
		traceAbsent := []string{}
		traceUnauthorized := []string{}
		runPathB := func(latentExists bool, trace *[]string) apivalid.Result {
			return apivalid.Validate(context.Background(), apivalid.Input{
				Validator:      mustCanonicalStructuralValidator(t, root),
				SchemaID:       CanonicalSchemasDir + "/project.json",
				Dst:            &project,
				Mode:           apivalid.ModeReadRepresentation,
				Stages:         stages,
				OperationScope: &opScope,
				TargetRef:      &target,
				TargetScopeResolver: conformanceTargetScopeResolver{
					available:    false,
					latentExists: latentExists,
					calls:        resolverCalls,
					trace:        trace,
				},
				Caller: &caller,
			}, apivalid.DefaultLimits())
		}
		resAbsent := runPathB(false, &traceAbsent)
		resUnauthorized := runPathB(true, &traceUnauthorized)
		if resolverCalls.Load() != 2 {
			t.Fatalf("Path B resolver calls=%d, want 2", resolverCalls.Load())
		}
		assertSafeDenialResult(t, resAbsent, want, "path B absent")
		assertSafeDenialResult(t, resUnauthorized, want, "path B unauthorized")
		if !reflect.DeepEqual(traceAbsent, traceUnauthorized) {
			t.Fatalf("Path B control-flow traces diverge: absent=%v unauthorized=%v",
				traceAbsent, traceUnauthorized)
		}

		// Combined AuthorizedResolver: absent and unauthorized share outcome.
		absentTrace := []string{}
		unauthTrace := []string{}
		arAbsent := conformanceAuthorizedResolver{found: false, latentExists: false, trace: &absentTrace}
		arUnauth := conformanceAuthorizedResolver{found: false, latentExists: true, trace: &unauthTrace}
		objA, foundA := arAbsent.Resolve(context.Background(), caller, target)
		objB, foundB := arUnauth.Resolve(context.Background(), caller, target)
		if foundA || foundB || objA != nil || objB != nil {
			t.Fatalf("AuthorizedResolver must return uniform unavailable; got found=%v/%v obj=%#v/%#v",
				foundA, foundB, objA, objB)
		}
		if !reflect.DeepEqual(absentTrace, unauthTrace) {
			t.Fatalf("AuthorizedResolver traces diverge: %v vs %v", absentTrace, unauthTrace)
		}
		pA := apivalid.SafeDenial(apivalid.DenyNotDisclosed)
		pB := apivalid.SafeDenial(apivalid.DenyNotDisclosed)
		ja, _ := json.Marshal(pA)
		jb, _ := json.Marshal(pB)
		if !bytes.Equal(ja, jb) || !bytes.Equal(ja, wantJSON) {
			t.Fatalf("AuthorizedResolver SafeDenial bytes diverge")
		}
	})
}

// ---------------------------------------------------------------------------
// Test helpers (Task 14.5)
// ---------------------------------------------------------------------------

func mustCanonicalStructuralValidator(t *testing.T, root string) *StructuralValidator {
	t.Helper()
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
	return v
}

func mustSchemaAllowedScopes(t *testing.T, root, schemaFile string) []apimeta.ScopeKind {
	t.Helper()
	raw, err := os.ReadFile(filepath.Join(root, CanonicalSchemasDir, schemaFile))
	if err != nil {
		t.Fatalf("read schema %s: %v", schemaFile, err)
	}
	meta, issues := apischema.ReadAnnotations(raw)
	if len(issues) != 0 {
		t.Fatalf("ReadAnnotations(%s): %#v", schemaFile, issues)
	}
	out := make([]apimeta.ScopeKind, len(meta.AllowedScopes))
	copy(out, meta.AllowedScopes)
	return out
}

func containsScopeKind(scopes []apimeta.ScopeKind, want apimeta.ScopeKind) bool {
	for _, s := range scopes {
		if s == want {
			return true
		}
	}
	return false
}

func hasViolationCode(vs []apiproblem.Violation, code apiproblem.ViolationCode) bool {
	for _, v := range vs {
		if v.Code == code {
			return true
		}
	}
	return false
}

func hasViolationField(vs []apiproblem.Violation, field string) bool {
	for _, v := range vs {
		if v.Field == field {
			return true
		}
	}
	return false
}

func assertSafeDenialResult(t *testing.T, res apivalid.Result, want *apiproblem.Problem, label string) {
	t.Helper()
	if res.FailedAt != apivalid.LayerAuthorization {
		t.Fatalf("%s: FailedAt=%v, want LayerAuthorization", label, res.FailedAt)
	}
	if res.Problem == nil {
		t.Fatalf("%s: Problem is nil", label)
	}
	got, err := json.Marshal(res.Problem)
	if err != nil {
		t.Fatalf("%s: marshal Problem: %v", label, err)
	}
	wantJSON, err := json.Marshal(want)
	if err != nil {
		t.Fatalf("%s: marshal want: %v", label, err)
	}
	if !bytes.Equal(got, wantJSON) {
		t.Fatalf("%s: SafeDenial not byte-identical:\ngot=%s\nwant=%s", label, got, wantJSON)
	}
	if len(res.Violations) != 0 {
		t.Fatalf("%s: must not disclose mismatch violations %#v", label, res.Violations)
	}
}

// secretLikeTokens are scanned case-insensitively against metadata label and
// annotation keys/values (F12-SEC-003). Composite phrases only — plain "key"
// is intentionally excluded.
var secretLikeTokens = []string{
	"password", "secret", "token", "credential",
	"apikey", "accesskey", "secretkey", "privatekey", "private_key",
	"connectionstring", "secretvalue",
}

func findSecretLikeInMetadataJSON(raw []byte) string {
	var top map[string]json.RawMessage
	if err := json.Unmarshal(raw, &top); err != nil {
		return ""
	}
	metaRaw, ok := top["metadata"]
	if !ok {
		return ""
	}
	var meta map[string]json.RawMessage
	if err := json.Unmarshal(metaRaw, &meta); err != nil {
		return ""
	}
	for _, field := range []string{"labels", "annotations"} {
		fieldRaw, ok := meta[field]
		if !ok {
			continue
		}
		var kv map[string]string
		if err := json.Unmarshal(fieldRaw, &kv); err != nil {
			continue
		}
		for k, v := range kv {
			if hit := secretLikeHit(k); hit != "" {
				return hit
			}
			if hit := secretLikeHit(v); hit != "" {
				return hit
			}
		}
	}
	return ""
}

func secretLikeHit(s string) string {
	lower := strings.ToLower(s)
	for _, tok := range secretLikeTokens {
		if strings.Contains(lower, tok) {
			return tok
		}
	}
	return ""
}

// conformanceRefCarrier is a minimal ReferenceCarrier for Property 4 complete.
type conformanceRefCarrier struct {
	scope    *apimeta.ScopeRef
	singular map[string]apiref.TypedRef
}

func (c *conformanceRefCarrier) GetScopeRef() *apimeta.ScopeRef { return c.scope }
func (c *conformanceRefCarrier) RefAt(path string) (apiref.TypedRef, bool) {
	ref, ok := c.singular[path]
	return ref, ok
}
func (c *conformanceRefCarrier) RefsAt(string) (apiref.Refs, bool) { return nil, false }

var _ apivalid.ReferenceCarrier = (*conformanceRefCarrier)(nil)

// conformanceSemanticCarrier adapts limit-boundary objects to SemanticCarrier.
type conformanceSemanticCarrier struct {
	apiVersion string
	kind       string
	name       string
	scope      *apimeta.ScopeRef
	owner      *apimeta.OwnerRef
	labels     map[string]string
	annot      map[string]string
	conditions []apicond.Condition
	phase      string
}

func (c *conformanceSemanticCarrier) APIVersion() string              { return c.apiVersion }
func (c *conformanceSemanticCarrier) Kind() string                    { return c.kind }
func (c *conformanceSemanticCarrier) ResourceName() string            { return c.name }
func (c *conformanceSemanticCarrier) GetScopeRef() *apimeta.ScopeRef  { return c.scope }
func (c *conformanceSemanticCarrier) GetOwnerRef() *apimeta.OwnerRef  { return c.owner }
func (c *conformanceSemanticCarrier) Labels() map[string]string       { return c.labels }
func (c *conformanceSemanticCarrier) Annotations() map[string]string  { return c.annot }
func (c *conformanceSemanticCarrier) Conditions() []apicond.Condition { return c.conditions }
func (c *conformanceSemanticCarrier) Phase() string                   { return c.phase }
func (c *conformanceSemanticCarrier) Profile() (apimeta.Profile, bool) {
	return "", false
}
func (c *conformanceSemanticCarrier) Boundary() (apimeta.Boundary, bool) {
	return "", false
}
func (c *conformanceSemanticCarrier) Stability() (apimeta.Stability, bool) {
	return "", false
}
func (c *conformanceSemanticCarrier) DataClassification() (apimeta.DataClassification, bool) {
	return "", false
}

var _ apivalid.SemanticCarrier = (*conformanceSemanticCarrier)(nil)

type conformanceNoopDefaulting struct{}

func (conformanceNoopDefaulting) Apply(_ context.Context, object any) (any, error) {
	return object, nil
}

type conformanceNoopValidation struct{}

func (conformanceNoopValidation) Validate(context.Context, any) ([]apiproblem.Violation, error) {
	return nil, nil
}

type conformanceScopeAuthorizer struct {
	decision apivalid.Decision
	calls    *atomic.Int32
}

func (a conformanceScopeAuthorizer) Authorize(
	context.Context,
	apivalid.CallerContext,
	apiref.TypedRef,
	apimeta.ScopeIdentity,
) apivalid.Decision {
	if a.calls != nil {
		a.calls.Add(1)
	}
	return a.decision
}

type conformanceTargetScopeResolver struct {
	available    bool
	latentExists bool
	calls        *atomic.Int32
	trace        *[]string
}

func (r conformanceTargetScopeResolver) ResolveAuthorizedTargetScope(
	context.Context,
	apivalid.CallerContext,
	apiref.TypedRef,
) (apimeta.ScopeIdentity, bool) {
	if r.calls != nil {
		r.calls.Add(1)
	}
	if r.trace != nil {
		// Latent existence must not change the side-effect trace.
		*r.trace = append(*r.trace, "resolve-unavailable")
	}
	_ = r.latentExists
	return apimeta.ScopeIdentity{}, r.available
}

type conformanceAuthorizedResolver struct {
	found        bool
	latentExists bool
	trace        *[]string
}

func (r conformanceAuthorizedResolver) Resolve(
	context.Context,
	apivalid.CallerContext,
	apiref.TypedRef,
) (any, bool) {
	if r.trace != nil {
		*r.trace = append(*r.trace, "resolve-uniform-unavailable")
	}
	_ = r.latentExists
	return nil, r.found
}

var (
	_ apivalid.ScopeAuthorizer               = conformanceScopeAuthorizer{}
	_ apivalid.AuthorizedTargetScopeResolver = conformanceTargetScopeResolver{}
	_ apivalid.AuthorizedResolver            = conformanceAuthorizedResolver{}
	_ apivalid.DefaultingStage               = conformanceNoopDefaulting{}
	_ apivalid.ValidationStage               = conformanceNoopValidation{}
)
