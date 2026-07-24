package apiconform

import (
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Grammar packages whose import direction is enforced by FEATURE-0012 (D-02).
var grammarPackages = []string{
	"apimeta",
	"apiref",
	"apicond",
	"apiproblem",
	"apivalid",
	"apischema",
	"apiconform",
}

// modulePath is the repository module path from go.mod.
const modulePath = "github.com/sanjeevksaini/sovrunn"

// approvedYAML is the only approved non-stdlib third-party dependency (D-14).
const approvedYAML = "gopkg.in/yaml.v3"

// providerSDKPrefixes are banned in all grammar packages (provider neutrality).
var providerSDKPrefixes = []string{
	"k8s.io/",
	"sigs.k8s.io/",
	"github.com/kubernetes/",
	"github.com/aws/",
	"github.com/Azure/",
	"github.com/microsoft/azure-",
	"cloud.google.com/",
	"google.golang.org/api/",
	"google.golang.org/genproto",
	"github.com/googleapis/",
	"github.com/oracle/oci-go-sdk",
	"github.com/digitalocean/",
	"github.com/hetznercloud/",
	"github.com/linode/",
	"github.com/vultr/",
	"github.com/IBM/",
	"github.com/hashicorp/terraform-provider",
	"github.com/pulumi/",
}

// allowedGrammarImports encodes the D-02 import-direction matrix.
// Standard library imports are always allowed. gopkg.in/yaml.v3 is allowed
// only for packages that the design uses for YAML decoding/ledger work.
var allowedGrammarImports = map[string]map[string]struct{}{
	"apimeta":    {},
	"apiref":     {"apimeta": {}},
	"apicond":    {},
	"apiproblem": {},
	"apivalid": {
		"apimeta":    {},
		"apiref":     {},
		"apicond":    {},
		"apiproblem": {},
	},
	"apischema": {
		"apimeta": {},
	},
	"apiconform": {
		"apimeta":    {},
		"apiref":     {},
		"apicond":    {},
		"apiproblem": {},
		"apivalid":   {},
		"apischema":  {},
	},
}

var yamlAllowedPackages = map[string]struct{}{
	"apivalid":   {},
	"apiconform": {},
}

func TestGrammarPackageImportBoundaries(t *testing.T) {
	t.Parallel()

	root := moduleRoot(t)
	for _, pkg := range grammarPackages {
		pkg := pkg
		t.Run(pkg, func(t *testing.T) {
			t.Parallel()
			dir := filepath.Join(root, "internal", pkg)
			imports := packageImports(t, dir)
			assertImportDirection(t, pkg, imports)
		})
	}
}

func assertImportDirection(t *testing.T, pkg string, imports []importRef) {
	t.Helper()

	allowed := allowedGrammarImports[pkg]
	if allowed == nil {
		t.Fatalf("missing allowed-import matrix entry for %s", pkg)
	}

	for _, imp := range imports {
		path := imp.path

		if isStdlib(path) {
			continue
		}

		if path == approvedYAML {
			if _, ok := yamlAllowedPackages[pkg]; !ok {
				t.Errorf("%s:%s imports %q; yaml.v3 is only allowed in apivalid and apiconform",
					imp.file, imp.path, path)
			}
			continue
		}

		if isProviderSDK(path) {
			t.Errorf("%s:%s imports provider SDK %q (provider neutrality violated)",
				imp.file, path, path)
			continue
		}

		if isForbiddenRuntime(path) {
			t.Errorf("%s:%s imports forbidden runtime package %q (must not import internal/api or internal/server)",
				imp.file, path, path)
			continue
		}

		grammarName, ok := grammarImportName(path)
		if !ok {
			t.Errorf("%s:%s imports disallowed non-grammar path %q", imp.file, path, path)
			continue
		}

		if grammarName == pkg {
			// Same-package import path should not appear, but ignore if it does.
			continue
		}

		if _, ok := allowed[grammarName]; !ok {
			t.Errorf("%s:%s imports grammar package %q which is outside the allowed set for %s (allowed: %v)",
				imp.file, path, grammarName, pkg, keys(allowed))
		}

		// Explicit negatives called out by the task.
		if pkg == "apivalid" && grammarName == "apischema" {
			t.Errorf("%s:%s: apivalid MUST NOT import apischema", imp.file, path)
		}
		if pkg == "apischema" && grammarName == "apiproblem" {
			t.Errorf("%s:%s: apischema MUST NOT import apiproblem", imp.file, path)
		}
	}
}

type importRef struct {
	file string
	path string
}

func packageImports(t *testing.T, dir string) []importRef {
	t.Helper()

	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("read package dir %s: %v", dir, err)
	}

	fset := token.NewFileSet()
	var out []importRef
	for _, ent := range entries {
		if ent.IsDir() || !strings.HasSuffix(ent.Name(), ".go") {
			continue
		}
		path := filepath.Join(dir, ent.Name())
		file, err := parser.ParseFile(fset, path, nil, parser.ImportsOnly)
		if err != nil {
			t.Fatalf("parse %s: %v", path, err)
		}
		for _, is := range file.Imports {
			if is.Path == nil {
				continue
			}
			impPath := strings.Trim(is.Path.Value, `"`)
			out = append(out, importRef{file: ent.Name(), path: impPath})
		}
	}
	return out
}

func moduleRoot(t *testing.T) string {
	t.Helper()

	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	dir := wd
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatalf("go.mod not found walking up from %s", wd)
		}
		dir = parent
	}
}

func isStdlib(importPath string) bool {
	if importPath == "" {
		return false
	}
	// Stdlib paths have no dot in the first path element (e.g. "fmt", "net/http").
	first, _, _ := strings.Cut(importPath, "/")
	return !strings.Contains(first, ".")
}

func grammarImportName(importPath string) (string, bool) {
	prefix := modulePath + "/internal/"
	if !strings.HasPrefix(importPath, prefix) {
		return "", false
	}
	rest := strings.TrimPrefix(importPath, prefix)
	name, _, _ := strings.Cut(rest, "/")
	for _, pkg := range grammarPackages {
		if name == pkg {
			return name, true
		}
	}
	return "", false
}

func isForbiddenRuntime(importPath string) bool {
	prefix := modulePath + "/internal/"
	if !strings.HasPrefix(importPath, prefix) {
		// Also catch relative-looking forms if ever present.
		return strings.Contains(importPath, "/internal/api") ||
			strings.Contains(importPath, "/internal/server") ||
			importPath == "internal/api" ||
			importPath == "internal/server"
	}
	rest := strings.TrimPrefix(importPath, prefix)
	name, _, _ := strings.Cut(rest, "/")
	return name == "api" || name == "server"
}

func isProviderSDK(importPath string) bool {
	for _, p := range providerSDKPrefixes {
		if strings.HasPrefix(importPath, p) {
			return true
		}
	}
	return false
}

func keys(m map[string]struct{}) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}
