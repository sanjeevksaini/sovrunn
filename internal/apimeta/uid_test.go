package apimeta

import (
	"testing"
)

func TestGenerateUID(t *testing.T) {
	t.Parallel()

	seen := make(map[string]struct{}, 64)
	for i := 0; i < 64; i++ {
		uid, err := GenerateUID()
		if err != nil {
			t.Fatalf("GenerateUID: %v", err)
		}
		if !IsGeneratedUIDFormat(uid) {
			t.Fatalf("GenerateUID returned non-opaque format %q", uid)
		}
		if uid == PlatformScopeUID {
			t.Fatal("GenerateUID must never return PlatformScopeUID")
		}
		if _, dup := seen[uid]; dup {
			t.Fatalf("GenerateUID collision in test sample: %q", uid)
		}
		seen[uid] = struct{}{}
	}
}

func TestPlatformScopeUIDNotGeneratedFormat(t *testing.T) {
	t.Parallel()

	if IsGeneratedUIDFormat(PlatformScopeUID) {
		t.Fatalf("PlatformScopeUID %q must not be a valid generated uid format", PlatformScopeUID)
	}
	if IsGeneratedUIDFormat("") {
		t.Fatal("empty string must not be a valid generated uid format")
	}
	if IsGeneratedUIDFormat("ABCDEF0123456789ABCDEF0123456789") {
		t.Fatal("uppercase hex must not match generated uid format")
	}
	if IsGeneratedUIDFormat("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa") {
		t.Fatal("31-char string must not match generated uid format")
	}
	if !IsGeneratedUIDFormat("0123456789abcdef0123456789abcdef") {
		t.Fatal("32 lowercase hex chars must match generated uid format")
	}
}
