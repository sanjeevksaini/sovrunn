package apiref

import (
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
)

func TestTypedRefAlias(t *testing.T) {
	t.Parallel()

	ref := TypedRef{
		APIVersion: "fabric.sovrunn.io/v1alpha1",
		Kind:       "ResourcePool",
		Name:       "sovereign-pool-a",
		UID:        "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
	}
	base := apimeta.TypedRef(ref)
	if base.Name != ref.Name {
		t.Fatalf("TypedRef alias must be identical to apimeta.TypedRef")
	}
}

func TestValidateRefAllowedKind(t *testing.T) {
	t.Parallel()

	c := Constraint{
		AllowedKinds: []string{"ResourcePool"},
		Direction:    DirectionOutbound,
	}
	ref := TypedRef{
		APIVersion: "fabric.sovrunn.io/v1alpha1",
		Kind:       "ResourcePool",
		Name:       "sovereign-pool-a",
	}
	if issues := c.ValidateRef(ref, "/spec/resourcePoolRef"); len(issues) != 0 {
		t.Fatalf("allowed kind must pass, got %#v", issues)
	}
}

func TestValidateRefDisallowedKind(t *testing.T) {
	t.Parallel()

	c := Constraint{
		AllowedKinds: []string{"ResourcePool"},
		Direction:    DirectionOutbound,
	}
	ref := TypedRef{
		APIVersion: "fabric.sovrunn.io/v1alpha1",
		Kind:       "Project",
		Name:       "payments",
	}
	issues := c.ValidateRef(ref, "/spec/resourcePoolRef")
	if !hasCode(issues, CodeKindNotAllowed) {
		t.Fatalf("disallowed kind must yield %s, got %#v", CodeKindNotAllowed, issues)
	}
	if !hasPath(issues, "/spec/resourcePoolRef/kind") {
		t.Fatalf("kind issue path must target kind field, got %#v", issues)
	}
}

func TestValidateRefAllowedAndDisallowedScope(t *testing.T) {
	t.Parallel()

	c := Constraint{
		AllowedScopes: []apimeta.ScopeKind{apimeta.ScopeTenant, apimeta.ScopeProject},
		Direction:     DirectionInbound,
	}
	ok := TypedRef{
		APIVersion: "core.sovrunn.io/v1alpha1",
		Kind:       string(apimeta.ScopeTenant),
		Name:       "acme",
		UID:        "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb",
	}
	if issues := c.ValidateRef(ok, "/metadata/scopeRef"); len(issues) != 0 {
		t.Fatalf("allowed scope must pass, got %#v", issues)
	}

	bad := TypedRef{
		APIVersion: "core.sovrunn.io/v1alpha1",
		Kind:       string(apimeta.ScopeOrganization),
		Name:       "org-1",
		UID:        "cccccccccccccccccccccccccccccccc",
	}
	issues := c.ValidateRef(bad, "/metadata/scopeRef")
	if !hasCode(issues, CodeScopeNotAllowed) {
		t.Fatalf("disallowed scope must yield %s, got %#v", CodeScopeNotAllowed, issues)
	}
}

func TestValidateRefInvalidDirection(t *testing.T) {
	t.Parallel()

	c := Constraint{
		AllowedKinds: []string{"Project"},
		Direction:    Direction("Sideways"),
	}
	ref := TypedRef{
		APIVersion: "core.sovrunn.io/v1alpha1",
		Kind:       "Project",
		Name:       "payments",
	}
	issues := c.ValidateRef(ref, "/spec/projectRef")
	if !hasCode(issues, CodeDirectionInvalid) {
		t.Fatalf("invalid direction must yield %s, got %#v", CodeDirectionInvalid, issues)
	}
}

func TestCheckNameUIDAgreement(t *testing.T) {
	t.Parallel()

	const path = "/spec/targetRef"
	const uid = "dddddddddddddddddddddddddddddddd"

	agree := TypedRef{
		APIVersion: "core.sovrunn.io/v1alpha1",
		Kind:       "Project",
		Name:       "payments",
		UID:        uid,
	}
	if issues := CheckNameUIDAgreement(agree, path, "payments", uid); len(issues) != 0 {
		t.Fatalf("matching name/uid must pass, got %#v", issues)
	}

	// UID omitted: human-authored input; agreement not required yet.
	omitUID := TypedRef{
		APIVersion: "core.sovrunn.io/v1alpha1",
		Kind:       "Project",
		Name:       "payments",
	}
	if issues := CheckNameUIDAgreement(omitUID, path, "payments", uid); len(issues) != 0 {
		t.Fatalf("omitted uid must not fail agreement, got %#v", issues)
	}

	// Same UID, different name: stale rebinding / mismatch (F12-REF-002).
	mismatch := TypedRef{
		APIVersion: "core.sovrunn.io/v1alpha1",
		Kind:       "Project",
		Name:       "old-payments",
		UID:        uid,
	}
	issues := CheckNameUIDAgreement(mismatch, path, "payments", uid)
	if !hasCode(issues, CodeNameUIDMismatch) {
		t.Fatalf("name/uid mismatch must yield %s, got %#v", CodeNameUIDMismatch, issues)
	}
	if !hasPath(issues, path) {
		t.Fatalf("mismatch path = want %q, got %#v", path, issues)
	}

	// Same name, different UID.
	uidMismatch := TypedRef{
		APIVersion: "core.sovrunn.io/v1alpha1",
		Kind:       "Project",
		Name:       "payments",
		UID:        "eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee",
	}
	issues = CheckNameUIDAgreement(uidMismatch, path, "payments", uid)
	if !hasCode(issues, CodeNameUIDMismatch) {
		t.Fatalf("uid mismatch must yield %s, got %#v", CodeNameUIDMismatch, issues)
	}
}

func TestValidateRefProviderNativeRejection(t *testing.T) {
	t.Parallel()

	c := Constraint{
		AllowedKinds: []string{"ResourcePool", "AWS::RDS::DBInstance"},
		Direction:    DirectionOutbound,
	}

	cases := []struct {
		name string
		ref  TypedRef
	}{
		{
			name: "aws_arn_name",
			ref: TypedRef{
				APIVersion: "fabric.sovrunn.io/v1alpha1",
				Kind:       "ResourcePool",
				Name:       "arn:aws:rds:us-east-1:123456789012:db:prod",
			},
		},
		{
			name: "azure_resource_id",
			ref: TypedRef{
				APIVersion: "fabric.sovrunn.io/v1alpha1",
				Kind:       "ResourcePool",
				Name:       "/subscriptions/0000/resourceGroups/rg/providers/Microsoft.DBforPostgreSQL/servers/db1",
			},
		},
		{
			name: "gcp_resource_name",
			ref: TypedRef{
				APIVersion: "fabric.sovrunn.io/v1alpha1",
				Kind:       "ResourcePool",
				Name:       "projects/my-proj/locations/us-central1/instances/db1",
			},
		},
		{
			name: "cloudformation_kind",
			ref: TypedRef{
				APIVersion: "fabric.sovrunn.io/v1alpha1",
				Kind:       "AWS::RDS::DBInstance",
				Name:       "prod-db",
			},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			issues := c.ValidateRef(tc.ref, "/spec/resourcePoolRef")
			if !hasCode(issues, CodeProviderNativeID) {
				t.Fatalf("provider-native ref must yield %s, got %#v", CodeProviderNativeID, issues)
			}
		})
	}
}

func TestRefsValidateBounds(t *testing.T) {
	t.Parallel()

	c := Constraint{
		AllowedKinds: []string{"Project"},
		Direction:    DirectionBidirectional,
	}
	refs := make(Refs, DefaultMaxRefs+1)
	for i := range refs {
		refs[i] = TypedRef{
			APIVersion: "core.sovrunn.io/v1alpha1",
			Kind:       "Project",
			Name:       "p",
		}
	}
	issues := refs.Validate(c, "/spec/projectRefs", 0)
	if !hasCode(issues, CodeRefsExceedLimit) {
		t.Fatalf("over-limit Refs must yield %s, got %#v", CodeRefsExceedLimit, issues)
	}

	ok := Refs{
		{APIVersion: "core.sovrunn.io/v1alpha1", Kind: "Project", Name: "a"},
		{APIVersion: "core.sovrunn.io/v1alpha1", Kind: "Project", Name: "b"},
	}
	if issues := ok.Validate(c, "/spec/projectRefs", 8); len(issues) != 0 {
		t.Fatalf("in-bound Refs must pass, got %#v", issues)
	}

	disallowed := Refs{
		{APIVersion: "core.sovrunn.io/v1alpha1", Kind: "Tenant", Name: "t"},
	}
	issues = disallowed.Validate(c, "/spec/projectRefs", 8)
	if !hasCode(issues, CodeKindNotAllowed) {
		t.Fatalf("Refs element must apply kind constraint, got %#v", issues)
	}
	if !hasPath(issues, "/spec/projectRefs/0/kind") {
		t.Fatalf("Refs element path must include index, got %#v", issues)
	}
}

func TestDirectionValid(t *testing.T) {
	t.Parallel()

	for _, d := range []Direction{DirectionInbound, DirectionOutbound, DirectionBidirectional} {
		if !d.Valid() {
			t.Fatalf("Direction %q must be Valid", d)
		}
	}
	if Direction("").Valid() {
		t.Fatal("empty Direction must not be Valid")
	}
	if Direction("Both").Valid() {
		t.Fatal("unknown Direction must not be Valid")
	}
}

func hasCode(issues []RefIssue, code string) bool {
	for _, iss := range issues {
		if iss.Code == code {
			return true
		}
	}
	return false
}

func hasPath(issues []RefIssue, path string) bool {
	for _, iss := range issues {
		if iss.Path == path {
			return true
		}
	}
	return false
}
