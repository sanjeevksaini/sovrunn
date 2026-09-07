package main

import (
	"go/parser"
	"go/token"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"
	"testing/quick"
)

func TestValidFixturesPass(t *testing.T) {
	t.Parallel()

	diags, err := run(testdataPath(t, "valid"))
	if err != nil {
		t.Fatalf("run(valid) error = %v", err)
	}
	if len(diags) != 0 {
		t.Fatalf("expected no diagnostics, got %+v", diags)
	}
}

func TestInvalidFixturesFailPerRule(t *testing.T) {
	t.Parallel()

	cases := []struct {
		dir  string
		rule string
	}{
		{"rule001", rule001},
		{"rule002", rule002},
		{"rule003", rule003},
		{"rule004", rule004},
		{"rule005", rule005},
		{"rule006", rule006},
		{"rule007", rule007},
		{"rule008", rule008},
		{"rule009", rule009},
		{"rule010", rule010},
		{"rule011", rule011},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.dir, func(t *testing.T) {
			t.Parallel()

			diags, err := run(testdataPath(t, filepath.Join("invalid", tc.dir)))
			if err != nil {
				t.Fatalf("run(%s) error = %v", tc.dir, err)
			}
			if len(diags) == 0 {
				t.Fatalf("expected diagnostics for %s", tc.rule)
			}
			if diags[0].RuleID != tc.rule {
				t.Fatalf("first rule = %s, want %s", diags[0].RuleID, tc.rule)
			}
		})
	}
}

func TestCheckerUsesOnlyStdlib(t *testing.T) {
	t.Parallel()

	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime caller")
	}
	mainPath := filepath.Join(filepath.Dir(thisFile), "main.go")

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, mainPath, nil, parser.ImportsOnly)
	if err != nil {
		t.Fatalf("parse imports: %v", err)
	}
	for _, imp := range file.Imports {
		path := strings.Trim(imp.Path.Value, `"`)
		if strings.Contains(path, ".") {
			t.Fatalf("non-stdlib import found: %s", path)
		}
	}
}

func TestPropertyDeterministicDiagnostics(t *testing.T) {
	t.Parallel()

	dirs := []string{
		testdataPath(t, "valid"),
		testdataPath(t, filepath.Join("invalid", "rule001")),
		testdataPath(t, filepath.Join("invalid", "rule003")),
		testdataPath(t, filepath.Join("invalid", "rule007")),
		testdataPath(t, filepath.Join("invalid", "rule011")),
	}

	prop := func(n uint8) bool {
		dir := dirs[int(n)%len(dirs)]
		first, err := run(dir)
		if err != nil {
			return false
		}
		second, err := run(dir)
		if err != nil {
			return false
		}
		return reflect.DeepEqual(first, second)
	}

	if err := quick.Check(prop, &quick.Config{MaxCount: 16}); err != nil {
		t.Fatal(err)
	}
}

func testdataPath(t *testing.T, suffix string) string {
	t.Helper()
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime caller")
	}
	return filepath.Join(filepath.Dir(thisFile), "testdata", suffix)
}
