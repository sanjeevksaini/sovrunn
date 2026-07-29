package decision_test

import (
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// FEATURE-0013 T-033 / design §5 / architecture §27.7:
// assert the internal/decision package DAG is acyclic and free of prohibited
// network/database/queue/workflow/crypto/provider/runtime imports.

const modulePath = "github.com/sanjeevksaini/sovrunn"

// decisionPackages are the FEATURE-0013 domain packages under internal/decision.
// Keys are package-relative IDs used by the allowed-import matrix.
var decisionPackages = []string{
	"decision",
	"decision/graph",
	"decision/bundle",
	"decision/compose",
	"decision/validate",
}

// allowedDecisionImports encodes design §5 (arrows = allowed import direction).
// Subsets of the allowed set are valid; imports outside the set fail closed.
// Stdlib imports are handled separately. No third-party imports are permitted
// in the decision tree (yaml.v3 is consumed only via apivalid).
var allowedDecisionImports = map[string]map[string]struct{}{
	"decision": {
		modulePath + "/internal/apimeta": {},
	},
	"decision/graph": {
		modulePath + "/internal/apiproblem": {},
		modulePath + "/internal/decision":   {},
	},
	"decision/bundle": {
		modulePath + "/internal/apiproblem":     {},
		modulePath + "/internal/apivalid":       {},
		modulePath + "/internal/apischema":      {},
		modulePath + "/internal/decision":       {},
		modulePath + "/internal/decision/graph": {},
	},
	"decision/compose": {
		modulePath + "/internal/apimeta":        {},
		modulePath + "/internal/apiproblem":     {},
		modulePath + "/internal/decision":       {},
		modulePath + "/internal/decision/graph": {},
	},
	"decision/validate": {
		modulePath + "/internal/apimeta":          {},
		modulePath + "/internal/apiproblem":       {},
		modulePath + "/internal/apivalid":         {},
		modulePath + "/internal/decision":         {},
		modulePath + "/internal/decision/bundle":  {},
		modulePath + "/internal/decision/graph":   {},
		modulePath + "/internal/decision/compose": {},
	},
}

// providerSDKPrefixes are banned (AD-023 provider neutrality).
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

// bannedExactStdlib are stdlib roots that would introduce network/database/
// crypto/runtime coupling into the FEATURE-0013 contract packages.
// net/http and net/http/httptest are allowed only as FEATURE-0012 StrictDecode
// adapters (no dial/listen); bare "net" and other net/* clients are banned.
var bannedExactStdlib = map[string]struct{}{
	"database/sql":        {},
	"database/sql/driver": {},
	"plugin":              {},
	"syscall":             {},
	"net":                 {},
	"net/rpc":             {},
	"net/smtp":            {},
	"net/mail":            {},
	"crypto":              {},
}

// bannedStdlibPrefixes catch cryptographic and SQL driver subtrees.
var bannedStdlibPrefixes = []string{
	"crypto/",
	"database/",
}

// bannedModulePrefixes catch common queue/workflow/crypto/provider modules.
var bannedModulePrefixes = []string{
	"golang.org/x/crypto",
	"github.com/redis/",
	"github.com/go-redis/",
	"github.com/segmentio/kafka-go",
	"github.com/IBM/sarama",
	"github.com/Shopify/sarama",
	"github.com/nats-io/",
	"github.com/rabbitmq/",
	"github.com/streadway/amqp",
	"go.temporal.io/",
	"go.uber.org/cadence",
	"github.com/jackc/pgx",
	"github.com/lib/pq",
	"go.mongodb.org/",
	"github.com/golang/crypto",
}

// allowedNetHTTPAdapters are the only net/* imports permitted (StrictDecode).
var allowedNetHTTPAdapters = map[string]struct{}{
	"net/http":          {},
	"net/http/httptest": {},
}

type importRef struct {
	file string
	path string
}

func TestDecisionPackageImportDAG(t *testing.T) {
	t.Parallel()

	root := moduleRoot(t)
	edges := make(map[string]map[string]struct{}, len(decisionPackages))

	for _, pkg := range decisionPackages {
		dir := filepath.Join(root, "internal", filepath.FromSlash(pkg))
		imports := packageImports(t, dir)
		assertDecisionImports(t, pkg, imports)

		edges[pkg] = map[string]struct{}{}
		for _, imp := range imports {
			if id, ok := decisionPackageID(imp.path); ok {
				edges[pkg][id] = struct{}{}
			}
		}
	}
	if cycle := findCycle(edges); cycle != "" {
		t.Fatalf("decision package DAG must be acyclic; found cycle involving %s", cycle)
	}
}

func TestDecisionPackagesHaveNoProhibitedRuntimeImports(t *testing.T) {
	t.Parallel()

	root := moduleRoot(t)
	for _, pkg := range decisionPackages {
		pkg := pkg
		t.Run(pkg, func(t *testing.T) {
			t.Parallel()
			dir := filepath.Join(root, "internal", filepath.FromSlash(pkg))
			for _, imp := range packageImports(t, dir) {
				if reason := prohibitedReason(imp.path); reason != "" {
					t.Errorf("%s:%s imports prohibited path %q (%s)",
						imp.file, pkg, imp.path, reason)
				}
			}
		})
	}
}

func TestAPIDoesNotImportServer(t *testing.T) {
	t.Parallel()

	root := moduleRoot(t)
	apiDir := filepath.Join(root, "internal", "api")
	entries, err := os.ReadDir(apiDir)
	if err != nil {
		if os.IsNotExist(err) {
			t.Skip("internal/api not present")
		}
		t.Fatalf("read internal/api: %v", err)
	}
	fset := token.NewFileSet()
	for _, ent := range entries {
		if ent.IsDir() || !strings.HasSuffix(ent.Name(), ".go") {
			continue
		}
		path := filepath.Join(apiDir, ent.Name())
		file, err := parser.ParseFile(fset, path, nil, parser.ImportsOnly)
		if err != nil {
			t.Fatalf("parse %s: %v", path, err)
		}
		for _, is := range file.Imports {
			impPath := strings.Trim(is.Path.Value, `"`)
			if impPath == modulePath+"/internal/server" ||
				strings.HasPrefix(impPath, modulePath+"/internal/server/") ||
				impPath == "internal/server" {
				t.Errorf("%s imports internal/server (prohibited)", ent.Name())
			}
		}
	}
}

func assertDecisionImports(t *testing.T, pkg string, imports []importRef) {
	t.Helper()

	allowed := allowedDecisionImports[pkg]
	if allowed == nil {
		t.Fatalf("missing allowed-import matrix entry for %s", pkg)
	}

	for _, imp := range imports {
		path := imp.path

		if isStdlib(path) {
			if _, ok := allowedNetHTTPAdapters[path]; ok {
				continue
			}
			if reason := prohibitedReason(path); reason != "" {
				t.Errorf("%s:%s imports prohibited stdlib %q (%s)",
					imp.file, path, path, reason)
			}
			continue
		}

		if reason := prohibitedReason(path); reason != "" {
			t.Errorf("%s:%s imports prohibited path %q (%s)",
				imp.file, path, path, reason)
			continue
		}

		if _, ok := allowed[path]; !ok {
			t.Errorf("%s:%s imports %q outside the design §5 allowed set for %s",
				imp.file, path, path, pkg)
		}
	}
}

func prohibitedReason(importPath string) string {
	if _, ok := allowedNetHTTPAdapters[importPath]; ok {
		return ""
	}
	if _, ok := bannedExactStdlib[importPath]; ok {
		return "banned stdlib runtime/network/crypto/database import"
	}
	for _, p := range bannedStdlibPrefixes {
		if strings.HasPrefix(importPath, p) {
			return "banned stdlib prefix " + p
		}
	}
	for _, p := range providerSDKPrefixes {
		if strings.HasPrefix(importPath, p) {
			return "provider SDK (provider neutrality)"
		}
	}
	for _, p := range bannedModulePrefixes {
		if importPath == p || strings.HasPrefix(importPath, p+"/") {
			return "banned network/database/queue/workflow/crypto module"
		}
	}
	if importPath == modulePath+"/internal/api" ||
		strings.HasPrefix(importPath, modulePath+"/internal/api/") ||
		importPath == modulePath+"/internal/server" ||
		strings.HasPrefix(importPath, modulePath+"/internal/server/") {
		return "forbidden runtime package"
	}
	if importPath == modulePath+"/internal/apiconform" ||
		strings.HasPrefix(importPath, modulePath+"/internal/apiconform/") {
		return "apiconform must not be imported by decision (one-way conformance adapter only)"
	}
	// Any non-stdlib, non-module-internal third-party path is prohibited.
	if !isStdlib(importPath) && !strings.HasPrefix(importPath, modulePath+"/") {
		return "third-party dependency not permitted in decision tree"
	}
	return ""
}

func decisionPackageID(importPath string) (string, bool) {
	prefix := modulePath + "/internal/"
	if !strings.HasPrefix(importPath, prefix) {
		return "", false
	}
	rest := strings.TrimPrefix(importPath, prefix)
	switch {
	case rest == "decision":
		return "decision", true
	case strings.HasPrefix(rest, "decision/"):
		return rest, true
	default:
		return "", false
	}
}

func findCycle(edges map[string]map[string]struct{}) string {
	const (
		white = 0
		gray  = 1
		black = 2
	)
	color := make(map[string]int, len(edges))
	var visit func(string) string
	visit = func(n string) string {
		color[n] = gray
		for m := range edges[n] {
			switch color[m] {
			case gray:
				return n + " -> " + m
			case white:
				if c := visit(m); c != "" {
					return c
				}
			}
		}
		color[n] = black
		return ""
	}
	for n := range edges {
		if color[n] == white {
			if c := visit(n); c != "" {
				return c
			}
		}
	}
	return ""
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
		// Production import graph only; test files may import test helpers.
		if strings.HasSuffix(ent.Name(), "_test.go") {
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
	first, _, _ := strings.Cut(importPath, "/")
	return !strings.Contains(first, ".")
}
