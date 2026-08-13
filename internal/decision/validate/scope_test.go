package validate

import (
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
)

func projectScope(uid string) *apimeta.ScopeRef {
	return &apimeta.ScopeRef{
		TypedRef: apimeta.TypedRef{
			APIVersion: "core.sovrunn.io/v1alpha1",
			Kind:       string(apimeta.ScopeProject),
			Name:       "payments-production",
			UID:        uid,
		},
	}
}

func tenantScope(uid string) *apimeta.ScopeRef {
	return &apimeta.ScopeRef{
		TypedRef: apimeta.TypedRef{
			APIVersion: "core.sovrunn.io/v1alpha1",
			Kind:       string(apimeta.ScopeTenant),
			Name:       "acme",
			UID:        uid,
		},
	}
}

func serviceInstanceSubject(uid string) apimeta.TypedRef {
	return apimeta.TypedRef{
		APIVersion: "platform.sovrunn.io/v1alpha1",
		Kind:       KindServiceInstance,
		Name:       "postgres-primary",
		UID:        uid,
	}
}

func assertScopeViolation(t *testing.T, p *apiproblem.Problem, wantCode apiproblem.ViolationCode, wantField string) {
	t.Helper()
	if p == nil {
		t.Fatalf("expected Problem with %s at %s, got nil", wantCode, wantField)
	}
	if p.Code != apiproblem.CodeValidationFailed {
		t.Fatalf("top-level code = %q, want %q", p.Code, apiproblem.CodeValidationFailed)
	}
	if p.Status != 422 {
		t.Fatalf("status = %d, want 422", p.Status)
	}
	if p.Type != "urn:sovrunn:problem:validation-failed" {
		t.Fatalf("type = %q, want urn:sovrunn:problem:validation-failed", p.Type)
	}
	if len(p.Violations) != 1 {
		t.Fatalf("violations len=%d, want 1: %#v", len(p.Violations), p.Violations)
	}
	v := p.Violations[0]
	if v.Code != wantCode {
		t.Fatalf("violations[0].code = %q, want %q", v.Code, wantCode)
	}
	if v.Field != wantField {
		t.Fatalf("violations[0].field = %q, want %q", v.Field, wantField)
	}
	if v.Message == "" {
		t.Fatal("violations[0].message must be non-empty (redactable)")
	}
}

func TestValidateScope_PlatformAllowedNilOK(t *testing.T) {
	t.Parallel()

	p := ValidateScope(ScopeInput{
		Metadata:            apimeta.ObjectMeta{Name: "platform-decision"},
		PlatformAllowed:     true,
		Allowed:             []apimeta.ScopeKind{apimeta.ScopePlatform},
		CheckCanonicalScope: true,
		CanonicalScope:      "Platform",
	})
	if p != nil {
		t.Fatalf("canonical Platform must be accepted, got %#v", p)
	}
}

func TestValidateScope_ExplicitPlatformNormalizesWithoutConflict(t *testing.T) {
	t.Parallel()

	// F13-SCOPE-006 / edge 12: explicit Platform input normalizes to absent;
	// normalized form must not trigger DECISION_SCOPE_CONFLICT.
	explicit := &apimeta.ScopeRef{
		TypedRef: apimeta.TypedRef{
			APIVersion: "core.sovrunn.io/v1alpha1",
			Kind:       string(apimeta.ScopePlatform),
			Name:       "platform",
			UID:        apimeta.PlatformScopeUID,
		},
	}
	p := ValidateScope(ScopeInput{
		Metadata:            apimeta.ObjectMeta{Name: "platform-decision", ScopeRef: explicit},
		PlatformAllowed:     true,
		CheckCanonicalScope: true,
		CanonicalScope:      "Platform",
	})
	if p != nil {
		t.Fatalf("explicit Platform must normalize without conflict, got %#v", p)
	}
}

func TestValidateScope_RequiredWhenPlatformNotPermitted(t *testing.T) {
	t.Parallel()

	p := ValidateScope(ScopeInput{
		Metadata:        apimeta.ObjectMeta{Name: "project-decision"},
		PlatformAllowed: false,
		Allowed:         []apimeta.ScopeKind{apimeta.ScopeProject},
	})
	assertScopeViolation(t, p, CodeScopeRequired, ptrMetadataScopeRef)
}

func TestValidateScope_InvalidOutOfVocabularyIncludingServiceInstance(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		kind string
	}{
		{name: "ServiceInstance-as-scope", kind: KindServiceInstance},
		{name: "unknown-kind", kind: "Cluster"},
		{name: "empty-kind", kind: ""},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			p := ValidateScope(ScopeInput{
				Metadata: apimeta.ObjectMeta{
					Name: "bad-scope",
					ScopeRef: &apimeta.ScopeRef{
						TypedRef: apimeta.TypedRef{
							APIVersion: "core.sovrunn.io/v1alpha1",
							Kind:       tc.kind,
							Name:       "x",
							UID:        "uid-1",
						},
					},
				},
				Allowed: apimeta.AllScopeKinds(),
			})
			assertScopeViolation(t, p, CodeScopeInvalid, ptrMetadataScopeRefKind)
		})
	}
}

func TestValidateScope_InvalidMissingUID(t *testing.T) {
	t.Parallel()

	p := ValidateScope(ScopeInput{
		Metadata: apimeta.ObjectMeta{
			Name: "project-decision",
			ScopeRef: &apimeta.ScopeRef{
				TypedRef: apimeta.TypedRef{
					APIVersion: "core.sovrunn.io/v1alpha1",
					Kind:       string(apimeta.ScopeProject),
					Name:       "payments-production",
				},
			},
		},
		Allowed: []apimeta.ScopeKind{apimeta.ScopeProject},
	})
	assertScopeViolation(t, p, CodeScopeInvalid, ptrMetadataScopeRefUID)
}

func TestValidateScope_ConflictParallelScopeSource(t *testing.T) {
	t.Parallel()

	p := ValidateScope(ScopeInput{
		Metadata: apimeta.ObjectMeta{
			Name:     "project-decision",
			ScopeRef: projectScope("project-123"),
		},
		Allowed:             []apimeta.ScopeKind{apimeta.ScopeProject},
		ParallelScopeFields: []string{ptrTopLevelScopeRef},
	})
	assertScopeViolation(t, p, CodeScopeConflict, ptrTopLevelScopeRef)
}

func TestValidateScope_MismatchServiceInstanceSubjectNotUnderProject(t *testing.T) {
	t.Parallel()

	p := ValidateScope(ScopeInput{
		Metadata: apimeta.ObjectMeta{
			Name:     "tenant-scoped-with-si",
			ScopeRef: tenantScope("tenant-1"),
		},
		Allowed: []apimeta.ScopeKind{
			apimeta.ScopeTenant,
			apimeta.ScopeProject,
		},
		SubjectRefs: []apimeta.TypedRef{serviceInstanceSubject("service-456")},
	})
	assertScopeViolation(t, p, CodeScopeMismatch, ptrMetadataScopeRef)
}

func TestValidateScope_ServiceInstanceSubjectUnderProjectOK(t *testing.T) {
	t.Parallel()

	p := ValidateScope(ScopeInput{
		Metadata: apimeta.ObjectMeta{
			Name:     "si-under-project",
			ScopeRef: projectScope("project-123"),
		},
		Allowed:             []apimeta.ScopeKind{apimeta.ScopeProject},
		SubjectRefs:         []apimeta.TypedRef{serviceInstanceSubject("service-456")},
		CheckCanonicalScope: true,
		CanonicalScope:      "Project:project-123",
	})
	if p != nil {
		t.Fatalf("ServiceInstance subject under Project must be accepted, got %#v", p)
	}
}

func TestValidateScope_MismatchCanonicalScopeDescriptor(t *testing.T) {
	t.Parallel()

	p := ValidateScope(ScopeInput{
		Metadata: apimeta.ObjectMeta{
			Name:     "project-decision",
			ScopeRef: projectScope("project-123"),
		},
		Allowed:             []apimeta.ScopeKind{apimeta.ScopeProject},
		CheckCanonicalScope: true,
		CanonicalScope:      "Tenant:tenant-9", // derived-descriptor mismatch
	})
	assertScopeViolation(t, p, CodeScopeMismatch, ptrCanonicalScope)
}

func TestValidateScope_WideningDisallowedScope(t *testing.T) {
	t.Parallel()

	// Valid Matrix B kind, but not in the governing profile allowedScopes.
	p := ValidateScope(ScopeInput{
		Metadata: apimeta.ObjectMeta{
			Name: "org-decision",
			ScopeRef: &apimeta.ScopeRef{
				TypedRef: apimeta.TypedRef{
					APIVersion: "core.sovrunn.io/v1alpha1",
					Kind:       string(apimeta.ScopeOrganization),
					Name:       "acme",
					UID:        "org-1",
				},
			},
		},
		Allowed: []apimeta.ScopeKind{apimeta.ScopeProject}, // Organization is a widening
	})
	assertScopeViolation(t, p, CodeScopeWidening, ptrMetadataScopeRefKind)
}

func TestValidateScope_AllSevenGovernanceScopesAccepted(t *testing.T) {
	t.Parallel()

	nonPlatform := []struct {
		kind apimeta.ScopeKind
		uid  string
	}{
		{apimeta.ScopeOrganization, "org-1"},
		{apimeta.ScopeOrganizationUnit, "ou-1"},
		{apimeta.ScopeTenant, "tenant-1"},
		{apimeta.ScopeProject, "project-1"},
		{apimeta.ScopeCloudPlatform, "cloud-platform-1"},
		{apimeta.ScopeCloudProvider, "cloud-provider-1"},
	}
	for _, tc := range nonPlatform {
		tc := tc
		t.Run(string(tc.kind), func(t *testing.T) {
			t.Parallel()
			scope := &apimeta.ScopeRef{
				TypedRef: apimeta.TypedRef{
					APIVersion: "core.sovrunn.io/v1alpha1",
					Kind:       string(tc.kind),
					Name:       "n",
					UID:        tc.uid,
				},
			}
			p := ValidateScope(ScopeInput{
				Metadata:            apimeta.ObjectMeta{Name: "rec", ScopeRef: scope},
				Allowed:             apimeta.AllScopeKinds(),
				CheckCanonicalScope: true,
				CanonicalScope:      CanonicalScopeDescriptor(scope),
			})
			if p != nil {
				t.Fatalf("%s scope must be accepted, got %#v", tc.kind, p)
			}
		})
	}
}

func TestCanonicalScopeDescriptor(t *testing.T) {
	t.Parallel()

	if got := CanonicalScopeDescriptor(nil); got != "Platform" {
		t.Fatalf("nil = %q, want Platform", got)
	}
	explicitPlatform := &apimeta.ScopeRef{
		TypedRef: apimeta.TypedRef{Kind: string(apimeta.ScopePlatform), UID: apimeta.PlatformScopeUID},
	}
	if got := CanonicalScopeDescriptor(explicitPlatform); got != "Platform" {
		t.Fatalf("explicit Platform = %q, want Platform", got)
	}
	if got := CanonicalScopeDescriptor(projectScope("project-123")); got != "Project:project-123" {
		t.Fatalf("project = %q, want Project:project-123", got)
	}
}

func TestValidateScopeMeta_MatchesDesignSignature(t *testing.T) {
	t.Parallel()

	p := ValidateScopeMeta(
		apimeta.ObjectMeta{Name: "x"},
		[]apimeta.ScopeKind{apimeta.ScopeProject},
		false,
	)
	assertScopeViolation(t, p, CodeScopeRequired, ptrMetadataScopeRef)
}

func TestValidateScope_ConflictShortCircuitsBeforeRequired(t *testing.T) {
	t.Parallel()

	// Parallel source present and scope absent: sole-authority conflict wins.
	p := ValidateScope(ScopeInput{
		Metadata:            apimeta.ObjectMeta{Name: "x"},
		PlatformAllowed:     false,
		ParallelScopeFields: []string{ptrTopLevelScopeRef},
	})
	assertScopeViolation(t, p, CodeScopeConflict, ptrTopLevelScopeRef)
}

func TestValidateScope_NoSecondScopeAuthorityInvented(t *testing.T) {
	t.Parallel()

	// OwnerRef must not become a scope substitute; ValidateScope only reads
	// metadata.scopeRef. Presence of ownerRef alone does not satisfy scope.
	p := ValidateScope(ScopeInput{
		Metadata: apimeta.ObjectMeta{
			Name: "x",
			// ScopeRef intentionally absent; ownerRef is not modeled on ObjectMeta
			// as a governance scope. Platform not allowed → REQUIRED on scopeRef.
		},
		PlatformAllowed: false,
		Allowed:         []apimeta.ScopeKind{apimeta.ScopeProject},
	})
	assertScopeViolation(t, p, CodeScopeRequired, ptrMetadataScopeRef)
}
