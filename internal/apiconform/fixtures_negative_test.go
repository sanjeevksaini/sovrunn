package apiconform

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/apicond"
	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/apiref"
	"github.com/sanjeevksaini/sovrunn/internal/apischema"
	"github.com/sanjeevksaini/sovrunn/internal/apivalid"
)

// ConformanceNegativeFixturesDir holds Task 14.3 invalid fixtures
// (F12-VALIDATION-002/006, F12-OWNER-001, F12-REF-002, F12-SCOPE-002, D-17).
const ConformanceNegativeFixturesDir = "tests/conformance/fixtures/negative"

// negativeSemanticCarrier adapts fixture identity fields to CommonSemantic
// without requiring domain types to implement SemanticCarrier (task 14.3 only).
type negativeSemanticCarrier struct {
	apiVersion string
	kind       string
	name       string
	scope      *apimeta.ScopeRef
	owner      *apimeta.OwnerRef
}

func (c *negativeSemanticCarrier) APIVersion() string              { return c.apiVersion }
func (c *negativeSemanticCarrier) Kind() string                    { return c.kind }
func (c *negativeSemanticCarrier) ResourceName() string            { return c.name }
func (c *negativeSemanticCarrier) GetScopeRef() *apimeta.ScopeRef  { return c.scope }
func (c *negativeSemanticCarrier) GetOwnerRef() *apimeta.OwnerRef  { return c.owner }
func (c *negativeSemanticCarrier) Labels() map[string]string       { return nil }
func (c *negativeSemanticCarrier) Annotations() map[string]string  { return nil }
func (c *negativeSemanticCarrier) Conditions() []apicond.Condition { return nil }
func (c *negativeSemanticCarrier) Phase() string                   { return "" }
func (c *negativeSemanticCarrier) Profile() (apimeta.Profile, bool) {
	return "", false
}
func (c *negativeSemanticCarrier) Boundary() (apimeta.Boundary, bool) {
	return "", false
}
func (c *negativeSemanticCarrier) Stability() (apimeta.Stability, bool) {
	return "", false
}
func (c *negativeSemanticCarrier) DataClassification() (apimeta.DataClassification, bool) {
	return "", false
}

// TestNegativeFixturesRejectedWithStableCodeAndPointer verifies Task 14.3:
// each negative fixture is rejected with the expected stable code and JSON
// Pointer (design testing strategy negative suite; D-17).
func TestNegativeFixturesRejectedWithStableCodeAndPointer(t *testing.T) {
	t.Parallel()

	root := moduleRoot(t)
	dir := filepath.Join(root, ConformanceNegativeFixturesDir)
	lim := apivalid.DefaultLimits()
	readPol := apivalid.PolicyFor(apivalid.ModeReadRepresentation)
	createPol := apivalid.PolicyFor(apivalid.ModeCreateRequest)

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
	structural, err := NewStructuralValidator(cfg)
	if err != nil {
		t.Fatalf("NewStructuralValidator: %v", err)
	}

	type decodeDst struct {
		APIVersion string         `json:"apiVersion"`
		Kind       string         `json:"kind"`
		Metadata   map[string]any `json:"metadata"`
		Spec       map[string]any `json:"spec"`
		Status     map[string]any `json:"status,omitempty"`
	}

	cases := []struct {
		name      string
		wantCode  string
		wantField string
		run       func(t *testing.T) (code, field string)
	}{
		{
			name:      "unknown-field.json",
			wantCode:  string(apiproblem.CodeUnknownField),
			wantField: "/extraField",
			run: func(t *testing.T) (string, string) {
				t.Helper()
				raw := mustReadNegative(t, dir, "unknown-field.json")
				var dst decodeDst
				prob := apivalid.DecodeJSON(raw, lim, readPol, &dst)
				return problemCodeField(t, prob)
			},
		},
		{
			name:      "duplicate-key.json",
			wantCode:  string(apiproblem.CodeDuplicateField),
			wantField: "/metadata/name",
			run: func(t *testing.T) (string, string) {
				t.Helper()
				raw := mustReadNegative(t, dir, "duplicate-key.json")
				var dst decodeDst
				prob := apivalid.DecodeJSON(raw, lim, readPol, &dst)
				return problemCodeField(t, prob)
			},
		},
		{
			name:      "duplicate-key.yaml",
			wantCode:  string(apiproblem.CodeDuplicateField),
			wantField: "/metadata/name",
			run: func(t *testing.T) (string, string) {
				t.Helper()
				raw := mustReadNegative(t, dir, "duplicate-key.yaml")
				var dst decodeDst
				prob := apivalid.DecodeYAML(raw, lim, readPol, &dst)
				return problemCodeField(t, prob)
			},
		},
		{
			name:      "unauthorized-status-create.json",
			wantCode:  string(apiproblem.CodeValidationFailed),
			wantField: "/status",
			run: func(t *testing.T) (string, string) {
				t.Helper()
				raw := mustReadNegative(t, dir, "unauthorized-status-create.json")
				var dst decodeDst
				prob := apivalid.DecodeJSON(raw, lim, createPol, &dst)
				return problemCodeField(t, prob)
			},
		},
		{
			name:      "name-uid-mismatch.json",
			wantCode:  apiref.CodeNameUIDMismatch,
			wantField: "/metadata/scopeRef",
			run: func(t *testing.T) (string, string) {
				t.Helper()
				raw := mustReadNegative(t, dir, "name-uid-mismatch.json")
				var ref apiref.TypedRef
				if err := json.Unmarshal(raw, &ref); err != nil {
					t.Fatalf("unmarshal TypedRef: %v", err)
				}
				// Resolved identity shares the name but not the uid (F12-REF-002).
				issues := apiref.CheckNameUIDAgreement(
					ref,
					"/metadata/scopeRef",
					"acme-tenant",
					"b2c3d4e5f60718293a4b5c6d7e8f901a",
				)
				if len(issues) == 0 {
					t.Fatal("expected REF_NAME_UID_MISMATCH, got none")
				}
				return issues[0].Code, issues[0].Path
			},
		},
		{
			name:      "invalid-scope-kind.json",
			wantCode:  "ENUM_MISMATCH",
			wantField: "/metadata/scopeRef/kind",
			run: func(t *testing.T) (string, string) {
				t.Helper()
				raw := mustReadNegative(t, dir, "invalid-scope-kind.json")
				var dst Project
				if prob := apivalid.DecodeJSON(raw, lim, readPol, &dst); prob != nil {
					t.Fatalf("DecodeJSON: %#v", prob)
				}
				violations, err := structural.Validate(&dst, CanonicalSchemasDir+"/project.json")
				if err != nil {
					t.Fatalf("Validate: %v", err)
				}
				if len(violations) == 0 {
					t.Fatal("expected ENUM_MISMATCH violation, got none")
				}
				return string(violations[0].Code), violations[0].Field
			},
		},
		{
			name:      "ownerref-as-scoperef.json",
			wantCode:  string(apivalid.ViolationScopeRefRequired),
			wantField: "/metadata/scopeRef",
			run: func(t *testing.T) (string, string) {
				t.Helper()
				raw := mustReadNegative(t, dir, "ownerref-as-scoperef.json")
				var surf struct {
					APIVersion string             `json:"apiVersion"`
					Kind       string             `json:"kind"`
					Metadata   apimeta.ObjectMeta `json:"metadata"`
					OwnerRef   *apimeta.OwnerRef  `json:"ownerRef"`
				}
				if err := json.Unmarshal(raw, &surf); err != nil {
					t.Fatalf("unmarshal: %v", err)
				}
				stage := apivalid.NewCommonSemantic(lim, false)
				carrier := &negativeSemanticCarrier{
					apiVersion: surf.APIVersion,
					kind:       surf.Kind,
					name:       surf.Metadata.Name,
					scope:      surf.Metadata.ScopeRef,
					owner:      surf.OwnerRef,
				}
				if carrier.owner == nil {
					t.Fatal("fixture must carry ownerRef")
				}
				if carrier.scope != nil {
					t.Fatal("fixture must omit scopeRef so ownerRef cannot substitute")
				}
				violations, err := stage.Validate(context.Background(), carrier)
				if err != nil {
					t.Fatalf("Validate: %v", err)
				}
				if len(violations) == 0 {
					t.Fatal("expected SCOPE_REF_REQUIRED, got none")
				}
				return string(violations[0].Code), violations[0].Field
			},
		},
		{
			name:      "nil-scoperef-non-platform.json",
			wantCode:  string(apivalid.ViolationScopeRefRequired),
			wantField: "/metadata/scopeRef",
			run: func(t *testing.T) (string, string) {
				t.Helper()
				raw := mustReadNegative(t, dir, "nil-scoperef-non-platform.json")
				var dst Project
				if prob := apivalid.DecodeJSON(raw, lim, readPol, &dst); prob != nil {
					t.Fatalf("DecodeJSON: %#v", prob)
				}
				if dst.Metadata.ScopeRef != nil {
					t.Fatal("fixture must omit scopeRef")
				}
				stage := apivalid.NewCommonSemantic(lim, false)
				carrier := &negativeSemanticCarrier{
					apiVersion: dst.APIVersion,
					kind:       dst.Kind,
					name:       dst.Metadata.Name,
					scope:      dst.Metadata.ScopeRef,
				}
				violations, err := stage.Validate(context.Background(), carrier)
				if err != nil {
					t.Fatalf("Validate: %v", err)
				}
				if len(violations) == 0 {
					t.Fatal("expected SCOPE_REF_REQUIRED, got none")
				}
				return string(violations[0].Code), violations[0].Field
			},
		},
		{
			name:      "oversized-body.json",
			wantCode:  string(apiproblem.CodeRequestTooLarge),
			wantField: "/",
			run: func(t *testing.T) (string, string) {
				t.Helper()
				raw := mustReadNegative(t, dir, "oversized-body.json")
				tiny := apivalid.Limits{MaxObjectBytes: 32, MaxNestingDepth: 32}
				var dst decodeDst
				prob := apivalid.DecodeJSON(raw, tiny, createPol, &dst)
				return problemCodeField(t, prob)
			},
		},
		{
			name:      "over-nested.json",
			wantCode:  string(apiproblem.CodeRequestTooLarge),
			wantField: "/a/b",
			run: func(t *testing.T) (string, string) {
				t.Helper()
				raw := mustReadNegative(t, dir, "over-nested.json")
				shallow := apivalid.Limits{MaxObjectBytes: 1 << 20, MaxNestingDepth: 2}
				var dst map[string]any
				prob := apivalid.DecodeJSON(raw, shallow, readPol, &dst)
				return problemCodeField(t, prob)
			},
		},
		{
			name:      "unsupported-media-type",
			wantCode:  string(apiproblem.CodeUnsupportedMediaType),
			wantField: "/",
			run: func(t *testing.T) (string, string) {
				t.Helper()
				ct := strings.TrimSpace(string(mustReadNegative(t, dir, "unsupported-media-type.content-type")))
				body := mustReadNegative(t, dir, "unsupported-media-type.body.json")
				req := httptest.NewRequest(http.MethodPost, "/apis/core.sovrunn.io/v1alpha1/projects", strings.NewReader(string(body)))
				req.Header.Set("Content-Type", ct)
				rec := httptest.NewRecorder()
				var dst decodeDst
				prob := apivalid.StrictDecode(rec, req, lim, apivalid.ModeCreateRequest, &dst)
				return problemCodeField(t, prob)
			},
		},
		{
			name:      "unversioned-route.path",
			wantCode:  apischema.CodeRouteUnversioned,
			wantField: "/",
			run: func(t *testing.T) (string, string) {
				t.Helper()
				path := strings.TrimSpace(string(mustReadNegative(t, dir, "unversioned-route.path")))
				err := apischema.ValidateRoute(path)
				if err == nil {
					t.Fatal("expected ROUTE_UNVERSIONED, got nil")
				}
				re, ok := err.(*apischema.RouteError)
				if !ok {
					t.Fatalf("error type %T, want *apischema.RouteError", err)
				}
				// Route grammar has no body pointer; report document root.
				return re.Code, "/"
			},
		},
		{
			name:      "yaml-alias.yaml",
			wantCode:  string(apiproblem.CodeMalformedRequest),
			wantField: "/metadata/name",
			run: func(t *testing.T) (string, string) {
				t.Helper()
				raw := mustReadNegative(t, dir, "yaml-alias.yaml")
				if !strings.Contains(string(raw), "*") {
					t.Fatal("alias fixture must contain a YAML alias marker")
				}
				var dst decodeDst
				prob := apivalid.DecodeYAML(raw, lim, readPol, &dst)
				return problemCodeField(t, prob)
			},
		},
		{
			name:      "yaml-anchor.yaml",
			wantCode:  string(apiproblem.CodeMalformedRequest),
			wantField: "/metadata/name",
			run: func(t *testing.T) (string, string) {
				t.Helper()
				raw := mustReadNegative(t, dir, "yaml-anchor.yaml")
				var dst decodeDst
				prob := apivalid.DecodeYAML(raw, lim, readPol, &dst)
				return problemCodeField(t, prob)
			},
		},
		{
			name:      "yaml-merge-key.yaml",
			wantCode:  string(apiproblem.CodeMalformedRequest),
			wantField: "/metadata/<<",
			run: func(t *testing.T) (string, string) {
				t.Helper()
				raw := mustReadNegative(t, dir, "yaml-merge-key.yaml")
				var dst decodeDst
				prob := apivalid.DecodeYAML(raw, lim, readPol, &dst)
				return problemCodeField(t, prob)
			},
		},
		{
			name:      "yaml-custom-tag.yaml",
			wantCode:  string(apiproblem.CodeMalformedRequest),
			wantField: "/metadata/name",
			run: func(t *testing.T) (string, string) {
				t.Helper()
				raw := mustReadNegative(t, dir, "yaml-custom-tag.yaml")
				var dst decodeDst
				prob := apivalid.DecodeYAML(raw, lim, readPol, &dst)
				return problemCodeField(t, prob)
			},
		},
		{
			name:      "yaml-multiple-docs.yaml",
			wantCode:  string(apiproblem.CodeMalformedRequest),
			wantField: "/",
			run: func(t *testing.T) (string, string) {
				t.Helper()
				raw := mustReadNegative(t, dir, "yaml-multiple-docs.yaml")
				var dst decodeDst
				prob := apivalid.DecodeYAML(raw, lim, readPol, &dst)
				return problemCodeField(t, prob)
			},
		},
		{
			name:      "yaml-non-string-key.yaml",
			wantCode:  string(apiproblem.CodeMalformedRequest),
			wantField: "/metadata",
			run: func(t *testing.T) (string, string) {
				t.Helper()
				raw := mustReadNegative(t, dir, "yaml-non-string-key.yaml")
				var dst decodeDst
				prob := apivalid.DecodeYAML(raw, lim, readPol, &dst)
				return problemCodeField(t, prob)
			},
		},
		{
			name:      "yaml-non-finite.yaml",
			wantCode:  string(apiproblem.CodeMalformedRequest),
			wantField: "/spec/weight",
			run: func(t *testing.T) (string, string) {
				t.Helper()
				raw := mustReadNegative(t, dir, "yaml-non-finite.yaml")
				var dst decodeDst
				prob := apivalid.DecodeYAML(raw, lim, readPol, &dst)
				return problemCodeField(t, prob)
			},
		},
		{
			name:      "operation-scope-target-kind-mismatch.json",
			wantCode:  string(apiproblem.ViolationOperationTargetScopeMismatch),
			wantField: "/metadata/scopeRef",
			run: func(t *testing.T) (string, string) {
				t.Helper()
				raw := mustReadNegative(t, dir, "operation-scope-target-kind-mismatch.json")
				var op Operation
				if prob := apivalid.DecodeJSON(raw, lim, readPol, &op); prob != nil {
					t.Fatalf("DecodeJSON: %#v", prob)
				}
				opScope := apimeta.CanonicalScopeIdentity(op.Metadata.ScopeRef)
				// Target is a Project: resolved governance scope kind differs from Tenant.
				targetScope := apimeta.ScopeIdentity{
					Kind: apimeta.ScopeProject,
					UID:  op.Spec.TargetRef.UID,
				}
				v := apivalid.CheckOperationTargetScopeMatch(opScope, targetScope)
				if v == nil {
					t.Fatal("expected OPERATION_TARGET_SCOPE_MISMATCH, got nil")
				}
				return string(v.Code), v.Field
			},
		},
		{
			name:      "operation-scope-target-uid-mismatch.json",
			wantCode:  string(apiproblem.ViolationOperationTargetScopeMismatch),
			wantField: "/metadata/scopeRef",
			run: func(t *testing.T) (string, string) {
				t.Helper()
				raw := mustReadNegative(t, dir, "operation-scope-target-uid-mismatch.json")
				var op Operation
				if prob := apivalid.DecodeJSON(raw, lim, readPol, &op); prob != nil {
					t.Fatalf("DecodeJSON: %#v", prob)
				}
				opScope := apimeta.CanonicalScopeIdentity(op.Metadata.ScopeRef)
				targetScope := apimeta.ScopeIdentity{
					Kind: apimeta.ScopeTenant,
					UID:  op.Spec.TargetRef.UID,
				}
				v := apivalid.CheckOperationTargetScopeMatch(opScope, targetScope)
				if v == nil {
					t.Fatal("expected OPERATION_TARGET_SCOPE_MISMATCH, got nil")
				}
				return string(v.Code), v.Field
			},
		},
	}

	if len(cases) != 21 {
		t.Fatalf("expected 21 Task 14.3 negative cases, got %d", len(cases))
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			gotCode, gotField := tc.run(t)
			if gotCode != tc.wantCode {
				t.Fatalf("code = %q, want %q", gotCode, tc.wantCode)
			}
			if gotField != tc.wantField {
				t.Fatalf("field = %q, want %q", gotField, tc.wantField)
			}
			if !strings.HasPrefix(gotField, "/") {
				t.Fatalf("field %q is not an RFC 6901 JSON Pointer", gotField)
			}
		})
	}
}

func mustReadNegative(t *testing.T, dir, name string) []byte {
	t.Helper()
	raw, err := os.ReadFile(filepath.Join(dir, name))
	if err != nil {
		t.Fatalf("read %s: %v", name, err)
	}
	if len(strings.TrimSpace(string(raw))) == 0 {
		t.Fatalf("fixture %s is empty", name)
	}
	return raw
}

func problemCodeField(t *testing.T, prob *apiproblem.Problem) (code, field string) {
	t.Helper()
	if prob == nil {
		t.Fatal("expected Problem, got nil")
	}
	if len(prob.Violations) == 0 {
		t.Fatalf("expected Violations on Problem %#v", prob)
	}
	return string(prob.Code), prob.Violations[0].Field
}
