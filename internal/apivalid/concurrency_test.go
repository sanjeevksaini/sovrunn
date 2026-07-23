package apivalid

import (
	"net/http"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
)

func TestCheckIfMatchAbsentReturnsNil(t *testing.T) {
	t.Parallel()

	if got := CheckIfMatch("", "rv-1"); got != nil {
		t.Fatalf("absent If-Match: got %#v, want nil", got)
	}
	if got := CheckIfMatch("", ""); got != nil {
		t.Fatalf("absent If-Match with empty current: got %#v, want nil", got)
	}
}

func TestCheckIfMatchMatchReturnsNil(t *testing.T) {
	t.Parallel()

	const version = "opaque-rv-42"
	if got := CheckIfMatch(version, version); got != nil {
		t.Fatalf("matching If-Match: got %#v, want nil", got)
	}
}

func TestCheckIfMatchStaleReturns412(t *testing.T) {
	t.Parallel()

	got := CheckIfMatch("rv-stale", "rv-current")
	if got == nil {
		t.Fatal("stale If-Match: got nil, want 412 STALE_RESOURCE_VERSION Problem")
	}
	if got.Status != http.StatusPreconditionFailed {
		t.Fatalf("status = %d, want %d", got.Status, http.StatusPreconditionFailed)
	}
	if got.Code != apiproblem.CodeStaleResourceVersion {
		t.Fatalf("code = %q, want %q", got.Code, apiproblem.CodeStaleResourceVersion)
	}
	if got.Title != apiproblem.TitleFor(apiproblem.CodeStaleResourceVersion) {
		t.Fatalf("title = %q, want %q", got.Title, apiproblem.TitleFor(apiproblem.CodeStaleResourceVersion))
	}
	wantType := apiproblem.TypeURN(apiproblem.CodeStaleResourceVersion)
	if got.Type != wantType {
		t.Fatalf("type = %q, want %q", got.Type, wantType)
	}
	if got.Detail != "" {
		t.Fatalf("detail must stay empty (no version leakage); got %q", got.Detail)
	}
	if got.RequestID != "" {
		t.Fatalf("requestId must stay empty for caller correlation; got %q", got.RequestID)
	}
	if len(got.Violations) != 0 {
		t.Fatalf("violations must be empty; got %#v", got.Violations)
	}
}

func TestCheckIfMatchStaleWhenCurrentEmpty(t *testing.T) {
	t.Parallel()

	// Non-empty If-Match against empty current is still a mismatch.
	got := CheckIfMatch("rv-1", "")
	if got == nil {
		t.Fatal("non-empty If-Match vs empty current: got nil, want stale Problem")
	}
	if got.Code != apiproblem.CodeStaleResourceVersion {
		t.Fatalf("code = %q, want %q", got.Code, apiproblem.CodeStaleResourceVersion)
	}
	if got.Status != http.StatusPreconditionFailed {
		t.Fatalf("status = %d, want %d", got.Status, http.StatusPreconditionFailed)
	}
}
