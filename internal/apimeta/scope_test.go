package apimeta

import (
	"encoding/json"
	"testing"
)

func TestMatrixBScopeKinds(t *testing.T) {
	t.Parallel()

	want := []ScopeKind{
		"Platform",
		"Organization",
		"OrganizationUnit",
		"Tenant",
		"Project",
		"Provider",
	}
	got := AllScopeKinds()
	if len(got) != len(want) {
		t.Fatalf("AllScopeKinds len=%d, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("AllScopeKinds[%d]=%q, want %q", i, got[i], want[i])
		}
		if !got[i].Valid() {
			t.Fatalf("scope kind %q should be Valid", got[i])
		}
	}
	if ScopeKind("Cluster").Valid() {
		t.Fatal("unknown scope kind must not be Valid")
	}
}

func TestNormalizeScope(t *testing.T) {
	t.Parallel()

	if got := NormalizeScope(nil); got != nil {
		t.Fatalf("NormalizeScope(nil)=%v, want nil", got)
	}

	platform := &ScopeRef{TypedRef: TypedRef{
		APIVersion: "core.sovrunn.io/v1alpha1",
		Kind:       string(ScopePlatform),
		Name:       "platform",
		UID:        PlatformScopeUID,
	}}
	if got := NormalizeScope(platform); got != nil {
		t.Fatalf("NormalizeScope(Platform)=%v, want nil", got)
	}

	tenant := &ScopeRef{TypedRef: TypedRef{
		APIVersion: "core.sovrunn.io/v1alpha1",
		Kind:       string(ScopeTenant),
		Name:       "acme",
		UID:        "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
	}}
	if got := NormalizeScope(tenant); got != tenant {
		t.Fatalf("NormalizeScope(Tenant) must leave non-platform scope unchanged")
	}

	// Idempotent: normalizing nil (already canonical platform) stays nil.
	if got := NormalizeScope(NormalizeScope(platform)); got != nil {
		t.Fatalf("NormalizeScope idempotent: got %v, want nil", got)
	}
}

func TestCanonicalScopeIdentity(t *testing.T) {
	t.Parallel()

	got := CanonicalScopeIdentity(nil)
	if got.Kind != ScopePlatform || got.UID != PlatformScopeUID {
		t.Fatalf("CanonicalScopeIdentity(nil)=%+v, want Platform/%q", got, PlatformScopeUID)
	}

	explicitPlatform := &ScopeRef{TypedRef: TypedRef{Kind: string(ScopePlatform)}}
	got = CanonicalScopeIdentity(explicitPlatform)
	if got.Kind != ScopePlatform || got.UID != PlatformScopeUID {
		t.Fatalf("CanonicalScopeIdentity(Platform)=%+v, want Platform/%q", got, PlatformScopeUID)
	}

	// Normalize then canonicalize yields the same platform identity.
	got = CanonicalScopeIdentity(NormalizeScope(explicitPlatform))
	if got.Kind != ScopePlatform || got.UID != PlatformScopeUID {
		t.Fatalf("after NormalizeScope: %+v, want Platform/%q", got, PlatformScopeUID)
	}

	const tenantUID = "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
	tenant := &ScopeRef{TypedRef: TypedRef{
		Kind: string(ScopeTenant),
		UID:  tenantUID,
	}}
	got = CanonicalScopeIdentity(tenant)
	if got.Kind != ScopeTenant || got.UID != tenantUID {
		t.Fatalf("CanonicalScopeIdentity(Tenant)=%+v, want Tenant/%q", got, tenantUID)
	}
}

func TestTypedRefJSONPromotion(t *testing.T) {
	t.Parallel()

	ref := ScopeRef{TypedRef: TypedRef{
		APIVersion: "core.sovrunn.io/v1alpha1",
		Kind:       string(ScopeProject),
		Name:       "payments",
		UID:        "cccccccccccccccccccccccccccccccc",
	}}
	b, err := json.Marshal(ref)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var raw map[string]any
	if err := json.Unmarshal(b, &raw); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	for _, key := range []string{"apiVersion", "kind", "name", "uid"} {
		if _, ok := raw[key]; !ok {
			t.Fatalf("expected promoted JSON field %q in %s", key, string(b))
		}
	}
	if _, ok := raw["TypedRef"]; ok {
		t.Fatalf("TypedRef must not appear as a nested JSON object: %s", string(b))
	}

	owner := OwnerRef{TypedRef: TypedRef{
		APIVersion: "core.sovrunn.io/v1alpha1",
		Kind:       "Project",
		Name:       "payments",
	}}
	ob, err := json.Marshal(owner)
	if err != nil {
		t.Fatalf("owner marshal: %v", err)
	}
	var ownerRaw map[string]any
	if err := json.Unmarshal(ob, &ownerRaw); err != nil {
		t.Fatalf("owner unmarshal: %v", err)
	}
	if _, ok := ownerRaw["apiVersion"]; !ok {
		t.Fatalf("OwnerRef must promote TypedRef fields: %s", string(ob))
	}
}

func TestPlatformScopeUIDConstant(t *testing.T) {
	t.Parallel()
	if PlatformScopeUID != "platform" {
		t.Fatalf("PlatformScopeUID=%q, want %q", PlatformScopeUID, "platform")
	}
}
