package conformance

import (
	"fmt"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/validate"
)

func TestVS0_CF_F18_LocalRegistryCoverage(t *testing.T) {
	for i := 1; i <= 37; i++ {
		id := fmt.Sprintf("VS0-CF-F18-%02d", i)
		if _, ok := f18CaseRegistry[id]; !ok {
			t.Fatalf("missing local conformance case %s", id)
		}
	}
	for i := 39; i <= 49; i++ {
		id := fmt.Sprintf("VS0-CF-F18-%02d", i)
		if _, ok := f18CaseRegistry[id]; !ok {
			t.Fatalf("missing local conformance case %s", id)
		}
	}
	for i := 51; i <= 53; i++ {
		id := fmt.Sprintf("VS0-CF-F18-%02d", i)
		if _, ok := f18CaseRegistry[id]; !ok {
			t.Fatalf("missing local conformance case %s", id)
		}
	}
	for _, excluded := range []string{"VS0-CF-F18-38", "VS0-CF-F18-50", "VS0-CF-F18-54"} {
		if _, ok := f18CaseRegistry[excluded]; ok {
			t.Fatalf("%s must not be a runtime local case", excluded)
		}
	}
}

func TestVS0_CF_F01(t *testing.T) {
	// inherited authentication presence check: AUTH_REQUIRED/401
	p := apiproblem.New(apiproblem.CodeAuthRequired)
	if p.Code != apiproblem.CodeAuthRequired || p.Status != 401 {
		t.Fatalf("expected AUTH_REQUIRED/401, got %#v", p)
	}
}

func TestVS0_CF_F02(t *testing.T) {
	// inherited safe denial check: no resource enumeration detail
	p := validate.SafeDenyInaccessible()
	if p == nil || p.Code != apiproblem.CodeResourceNotFound || p.Status != 404 {
		t.Fatalf("expected RESOURCE_NOT_FOUND/404, got %#v", p)
	}
	if p.Detail != "" {
		t.Fatal("safe denial must not disclose inaccessible target details")
	}
}

func TestVS0_CF_X01(t *testing.T) {
	// inherited cross-Organization isolation check through safe visibility resolver
	p := validate.ResolveReferenceVisibility(true, false)
	if p == nil || p.Code != apiproblem.CodeResourceNotFound {
		t.Fatalf("expected safe-denial for inaccessible cross-organization target, got %#v", p)
	}
}

func TestVS0_CF_X02(t *testing.T) {
	// inherited cross-Project isolation check through safe visibility resolver
	p := validate.ResolveReferenceVisibility(false, true)
	if p == nil || p.Code != apiproblem.CodeResourceNotFound {
		t.Fatalf("expected safe-denial for inaccessible cross-project target, got %#v", p)
	}
}
