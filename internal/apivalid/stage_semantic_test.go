package apivalid

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/apicond"
	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
)

// Compile-time interface conformance for the test resource stub.
var _ SemanticCarrier = (*stubSemanticResource)(nil)

type stubSemanticResource struct {
	apiVersion   string
	kind         string
	name         string
	scope        *apimeta.ScopeRef
	owner        *apimeta.OwnerRef
	labels       map[string]string
	annotations  map[string]string
	conditions   []apicond.Condition
	phase        string
	profile      apimeta.Profile
	profileSet   bool
	boundary     apimeta.Boundary
	boundarySet  bool
	stability    apimeta.Stability
	stabilitySet bool
	dataClass    apimeta.DataClassification
	dataClassSet bool
}

func (r *stubSemanticResource) APIVersion() string              { return r.apiVersion }
func (r *stubSemanticResource) Kind() string                    { return r.kind }
func (r *stubSemanticResource) ResourceName() string            { return r.name }
func (r *stubSemanticResource) GetScopeRef() *apimeta.ScopeRef  { return r.scope }
func (r *stubSemanticResource) GetOwnerRef() *apimeta.OwnerRef  { return r.owner }
func (r *stubSemanticResource) Labels() map[string]string       { return r.labels }
func (r *stubSemanticResource) Annotations() map[string]string  { return r.annotations }
func (r *stubSemanticResource) Conditions() []apicond.Condition { return r.conditions }
func (r *stubSemanticResource) Phase() string                   { return r.phase }

func (r *stubSemanticResource) Profile() (apimeta.Profile, bool) {
	return r.profile, r.profileSet
}
func (r *stubSemanticResource) Boundary() (apimeta.Boundary, bool) {
	return r.boundary, r.boundarySet
}
func (r *stubSemanticResource) Stability() (apimeta.Stability, bool) {
	return r.stability, r.stabilitySet
}
func (r *stubSemanticResource) DataClassification() (apimeta.DataClassification, bool) {
	return r.dataClass, r.dataClassSet
}

func validSemanticStub() *stubSemanticResource {
	return &stubSemanticResource{
		apiVersion: "fabric.sovrunn.io/v1alpha1",
		kind:       "Project",
		name:       "payments-production",
		scope: &apimeta.ScopeRef{TypedRef: apimeta.TypedRef{
			APIVersion: "tenancy.sovrunn.io/v1alpha1",
			Kind:       string(apimeta.ScopeTenant),
			Name:       "acme",
			UID:        "tenant-uid-1",
		}},
	}
}

func TestCommonSemanticValidObjectPasses(t *testing.T) {
	t.Parallel()

	stage := NewCommonSemantic(DefaultLimits(), false)
	violations, err := stage.Validate(context.Background(), validSemanticStub())
	if err != nil {
		t.Fatalf("Validate: unexpected error: %v", err)
	}
	if len(violations) != 0 {
		t.Fatalf("valid object: violations=%#v, want none", violations)
	}
}

func TestCommonSemanticInvalidNameViolation(t *testing.T) {
	t.Parallel()

	stage := NewCommonSemantic(DefaultLimits(), false)
	obj := validSemanticStub()
	obj.name = "INVALID_NAME"

	violations, err := stage.Validate(context.Background(), obj)
	if err != nil {
		t.Fatalf("Validate: unexpected error: %v", err)
	}
	if !hasViolationCode(violations, ViolationInvalidResourceName) {
		t.Fatalf("invalid name must yield %s, got %#v", ViolationInvalidResourceName, violations)
	}
	if !hasViolationField(violations, "/metadata/name") {
		t.Fatalf("invalid name must point at /metadata/name, got %#v", violations)
	}
}

func TestCommonSemanticInvalidEnumViolation(t *testing.T) {
	t.Parallel()

	stage := NewCommonSemantic(DefaultLimits(), false)
	obj := validSemanticStub()
	obj.scope.Kind = "NotAScopeKind"
	obj.profile = apimeta.Profile("NotAProfile")
	obj.profileSet = true

	violations, err := stage.Validate(context.Background(), obj)
	if err != nil {
		t.Fatalf("Validate: unexpected error: %v", err)
	}
	if !hasViolationCode(violations, ViolationInvalidEnum) {
		t.Fatalf("invalid enum must yield %s, got %#v", ViolationInvalidEnum, violations)
	}
	if !hasViolationField(violations, "/metadata/scopeRef/kind") {
		t.Fatalf("invalid scope kind must point at /metadata/scopeRef/kind, got %#v", violations)
	}
}

func TestCommonSemanticOwnerRefMissingRequiredScopeRefFails(t *testing.T) {
	t.Parallel()

	// Non-platform resource: primary scopeRef required. ownerRef present must
	// not satisfy the requirement (F12-OWNER-001).
	stage := NewCommonSemantic(DefaultLimits(), false)
	obj := validSemanticStub()
	obj.scope = nil
	obj.owner = &apimeta.OwnerRef{TypedRef: apimeta.TypedRef{
		APIVersion: "tenancy.sovrunn.io/v1alpha1",
		Kind:       "Tenant",
		Name:       "acme",
		UID:        "tenant-uid-1",
	}}

	violations, err := stage.Validate(context.Background(), obj)
	if err != nil {
		t.Fatalf("Validate: unexpected error: %v", err)
	}
	if !hasViolationCode(violations, ViolationScopeRefRequired) {
		t.Fatalf("missing required scopeRef must yield %s, got %#v", ViolationScopeRefRequired, violations)
	}
	if !hasViolationField(violations, "/metadata/scopeRef") {
		t.Fatalf("missing scopeRef must point at /metadata/scopeRef, got %#v", violations)
	}
}

func TestCommonSemanticOwnerRefNeverSubstitutesForScopeAuthorization(t *testing.T) {
	t.Parallel()

	// Platform-allowed: nil scopeRef is canonical Platform. ownerRef pointing
	// at a Tenant must NOT become the governance scope for authorization.
	stage := NewCommonSemantic(DefaultLimits(), true)
	obj := validSemanticStub()
	obj.kind = "PluginDefinition"
	obj.scope = nil
	obj.owner = &apimeta.OwnerRef{TypedRef: apimeta.TypedRef{
		APIVersion: "tenancy.sovrunn.io/v1alpha1",
		Kind:       string(apimeta.ScopeTenant),
		Name:       "acme",
		UID:        "tenant-uid-1",
	}}

	violations, err := stage.Validate(context.Background(), obj)
	if err != nil {
		t.Fatalf("Validate: unexpected error: %v", err)
	}
	if len(violations) != 0 {
		t.Fatalf("platform + ownerRef must pass semantic checks, got %#v", violations)
	}

	// Authorization identity is scopeRef-derived only (D-16/D-17).
	got := apimeta.CanonicalScopeIdentity(obj.GetScopeRef())
	want := apimeta.ScopeIdentity{Kind: apimeta.ScopePlatform, UID: apimeta.PlatformScopeUID}
	if got != want {
		t.Fatalf("CanonicalScopeIdentity(scopeRef)=%#v, want %#v (ownerRef must not substitute)", got, want)
	}
	if got.Kind == apimeta.ScopeTenant || got.UID == "tenant-uid-1" {
		t.Fatal("authorization must never derive governance scope from ownerRef")
	}
}

func TestCommonSemanticOwnerAndScopeSameTargetAllowed(t *testing.T) {
	t.Parallel()

	stage := NewCommonSemantic(DefaultLimits(), false)
	obj := validSemanticStub()
	// Same target identity on ownerRef and scopeRef is allowed; replacement
	// / governance misuse is the prohibited behavior, not identity overlap.
	obj.owner = &apimeta.OwnerRef{TypedRef: apimeta.TypedRef{
		APIVersion: obj.scope.APIVersion,
		Kind:       obj.scope.Kind,
		Name:       obj.scope.Name,
		UID:        obj.scope.UID,
	}}

	violations, err := stage.Validate(context.Background(), obj)
	if err != nil {
		t.Fatalf("Validate: unexpected error: %v", err)
	}
	if len(violations) != 0 {
		t.Fatalf("same-target ownerRef+scopeRef must be allowed, got %#v", violations)
	}
}

func TestCommonSemanticOverLimitViolation(t *testing.T) {
	t.Parallel()

	limits := DefaultLimits()
	limits.MaxLabels = 2
	limits.MaxConditions = 1
	limits.MaxAnnotationsBytes = 8
	stage := NewCommonSemantic(limits, false)

	obj := validSemanticStub()
	obj.labels = map[string]string{"a": "1", "b": "2", "c": "3"}
	obj.annotations = map[string]string{"k": "too-large-value"}
	obj.conditions = []apicond.Condition{
		{Type: "Ready", Status: apicond.ConditionTrue, Reason: "Succeeded"},
		{Type: "Valid", Status: apicond.ConditionTrue, Reason: "Checked"},
	}

	violations, err := stage.Validate(context.Background(), obj)
	if err != nil {
		t.Fatalf("Validate: unexpected error: %v", err)
	}
	if !hasViolationCode(violations, apiproblem.ViolationOutOfRange) {
		t.Fatalf("over-limit must yield %s, got %#v", apiproblem.ViolationOutOfRange, violations)
	}
	if !hasViolationField(violations, "/metadata/labels") {
		t.Fatalf("over MaxLabels must point at /metadata/labels, got %#v", violations)
	}
	if !hasViolationField(violations, "/metadata/annotations") {
		t.Fatalf("over MaxAnnotationsBytes must point at /metadata/annotations, got %#v", violations)
	}
	if !hasViolationField(violations, "/status/conditions") {
		t.Fatalf("over MaxConditions must point at /status/conditions, got %#v", violations)
	}
}

func TestCommonSemanticConditionPascalCaseAndPhaseCoherence(t *testing.T) {
	t.Parallel()

	stage := NewCommonSemantic(DefaultLimits(), false)
	obj := validSemanticStub()
	obj.phase = "True" // condition-status vocabulary — incoherent as phase
	obj.conditions = []apicond.Condition{
		{Type: "not-pascal", Status: apicond.ConditionTrue, Reason: "also-bad"},
		{Type: "Ready", Status: apicond.ConditionTrue, Reason: "Ok"},
		{Type: "Ready", Status: apicond.ConditionFalse, Reason: "Dup"},
	}

	violations, err := stage.Validate(context.Background(), obj)
	if err != nil {
		t.Fatalf("Validate: unexpected error: %v", err)
	}
	if !hasViolationCode(violations, ViolationInvalidCondition) {
		t.Fatalf("non-PascalCase condition must yield %s, got %#v", ViolationInvalidCondition, violations)
	}
	if !hasViolationCode(violations, ViolationPhaseConditionIncoherent) {
		t.Fatalf("phase/condition coherence must yield %s, got %#v", ViolationPhaseConditionIncoherent, violations)
	}
}

func TestCommonSemanticInternalFaultReturnsError(t *testing.T) {
	t.Parallel()

	stage := NewCommonSemantic(DefaultLimits(), false)

	_, err := stage.Validate(context.Background(), nil)
	if !errors.Is(err, ErrSemanticInternal) {
		t.Fatalf("nil object: err=%v, want ErrSemanticInternal", err)
	}

	var typedNil *stubSemanticResource
	_, err = stage.Validate(context.Background(), typedNil)
	if !errors.Is(err, ErrSemanticInternal) {
		t.Fatalf("typed nil: err=%v, want ErrSemanticInternal", err)
	}

	var nilStage *CommonSemantic
	_, err = nilStage.Validate(context.Background(), validSemanticStub())
	if !errors.Is(err, ErrSemanticInternal) {
		t.Fatalf("nil stage: err=%v, want ErrSemanticInternal", err)
	}
}

func TestCommonSemanticNoOpForUnknownKind(t *testing.T) {
	t.Parallel()

	stage := NewCommonSemantic(DefaultLimits(), false)
	unknown := &unknownKindObject{Name: "x"}
	violations, err := stage.Validate(context.Background(), unknown)
	if err != nil {
		t.Fatalf("unknown type: unexpected error: %v", err)
	}
	if len(violations) != 0 {
		t.Fatalf("unknown type must no-op, got %#v", violations)
	}
}

func TestCommonSemanticInvalidAPIVersionAndKind(t *testing.T) {
	t.Parallel()

	stage := NewCommonSemantic(DefaultLimits(), false)
	obj := validSemanticStub()
	obj.apiVersion = "example.com/v2"
	obj.kind = "notPascal"

	violations, err := stage.Validate(context.Background(), obj)
	if err != nil {
		t.Fatalf("Validate: unexpected error: %v", err)
	}
	if !hasViolationCode(violations, ViolationInvalidAPIVersion) {
		t.Fatalf("bad apiVersion must yield %s, got %#v", ViolationInvalidAPIVersion, violations)
	}
	if !hasViolationCode(violations, ViolationInvalidKind) {
		t.Fatalf("bad kind must yield %s, got %#v", ViolationInvalidKind, violations)
	}
}

func TestCommonSemanticMessagesDoNotEmbedSecrets(t *testing.T) {
	t.Parallel()

	stage := NewCommonSemantic(DefaultLimits(), false)
	obj := validSemanticStub()
	secret := "super-secret-token-value"
	obj.name = "Bad_Name"
	obj.annotations = map[string]string{"x": secret}
	// Force annotation over-limit with a huge payload that includes the secret
	// so we can assert messages never echo annotation contents.
	obj.annotations["x"] = strings.Repeat("a", DefaultLimits().MaxAnnotationsBytes+1)

	violations, err := stage.Validate(context.Background(), obj)
	if err != nil {
		t.Fatalf("Validate: unexpected error: %v", err)
	}
	for _, v := range violations {
		if strings.Contains(v.Message, secret) || strings.Contains(v.Message, "aaaa") {
			t.Fatalf("violation message must not embed field payloads: %q", v.Message)
		}
	}
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
