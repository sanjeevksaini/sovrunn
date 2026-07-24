package apiconform

import (
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/apischema"
)

func TestTypeBindingsCoverCanonicalAndCommonSchemas(t *testing.T) {
	t.Parallel()

	want := make([]string, 0, len(commonSchemaFiles)+len(canonicalSchemaFiles))
	for _, name := range commonSchemaFiles {
		want = append(want, "api/schemas/_common/"+name)
	}
	for _, name := range canonicalSchemaFiles {
		want = append(want, "api/schemas/"+name)
	}
	sort.Strings(want)

	got := make([]string, 0, len(TypeBindings))
	seen := make(map[string]reflect.Type, len(TypeBindings))
	for _, b := range TypeBindings {
		if b.SchemaPath == "" {
			t.Fatal("TypeBinding SchemaPath must not be empty")
		}
		if b.GoType == nil {
			t.Fatalf("TypeBinding %q GoType must not be nil", b.SchemaPath)
		}
		if prev, dup := seen[b.SchemaPath]; dup {
			t.Fatalf("duplicate TypeBinding for %q (%s and %s)", b.SchemaPath, typeName(prev), typeName(b.GoType))
		}
		seen[b.SchemaPath] = b.GoType
		got = append(got, b.SchemaPath)
	}
	sort.Strings(got)

	if len(got) != len(want) {
		t.Fatalf("TypeBindings count = %d, want %d\ngot:  %v\nwant: %v", len(got), len(want), got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("TypeBindings[%d] = %q, want %q\ngot:  %v\nwant: %v", i, got[i], want[i], got, want)
		}
	}
}

func TestVerifyGoTypeAgainstSchemaForAllTypeBindings(t *testing.T) {
	t.Parallel()

	root := moduleRoot(t)
	for _, binding := range TypeBindings {
		binding := binding
		t.Run(binding.SchemaPath, func(t *testing.T) {
			t.Parallel()

			schemaPath := filepath.Join(root, filepath.FromSlash(binding.SchemaPath))
			schema, err := os.ReadFile(schemaPath)
			if err != nil {
				t.Fatalf("read schema %s: %v", schemaPath, err)
			}

			supportIssues := apischema.ValidateSchemaSupport(schema)
			if len(supportIssues) > 0 {
				t.Fatalf("ValidateSchemaSupport failed for %s: %v", binding.SchemaPath, supportIssues)
			}

			issues := apischema.VerifyGoTypeAgainstSchema(schema, binding.GoType)
			if len(issues) > 0 {
				var b strings.Builder
				for _, issue := range issues {
					b.WriteString("\n  ")
					b.WriteString(issue.Path)
					b.WriteString(" ")
					b.WriteString(issue.Code)
					b.WriteString(": ")
					b.WriteString(issue.Message)
				}
				t.Fatalf("VerifyGoTypeAgainstSchema failed for %s → %s:%s",
					binding.SchemaPath, binding.GoType.PkgPath(), binding.GoType.Name()+b.String())
			}
		})
	}
}

func TestTypeBindingsRejectDeliberateMismatch(t *testing.T) {
	t.Parallel()

	schema, err := os.ReadFile(filepath.Join(moduleRoot(t), "api/schemas/_common/page.json"))
	if err != nil {
		t.Fatalf("read page schema: %v", err)
	}

	type mismatchedPage struct {
		NextPageToken int `json:"nextPageToken,omitempty"`
	}
	issues := apischema.VerifyGoTypeAgainstSchema(schema, reflect.TypeOf(mismatchedPage{}))
	if len(issues) == 0 {
		t.Fatal("expected deliberate Go-type mismatch to be rejected")
	}
}

func typeName(t reflect.Type) string {
	if t == nil {
		return "<nil>"
	}
	if t.Name() != "" {
		if t.PkgPath() != "" {
			return t.PkgPath() + "." + t.Name()
		}
		return t.Name()
	}
	return t.String()
}
