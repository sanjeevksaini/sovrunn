package policyseam_test

import (
	"context"
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/decision"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/policyseam"
	"github.com/sanjeevksaini/sovrunn/internal/policyeval"
)

type fakeBoundary struct {
	result policyeval.PolicyEvaluationResult
	non    policyeval.NonResult
	err    error
}

func (f fakeBoundary) Evaluate(context.Context, policyeval.PolicyEvaluationRequest) (policyeval.PolicyEvaluationResult, decision.TimingEnvelope, policyeval.NonResult, error) {
	fixed := time.Date(2026, 9, 7, 10, 0, 0, 0, time.UTC).Format(time.RFC3339Nano)
	return f.result, decision.TimingEnvelope{
		StartedAt:   fixed,
		CompletedAt: fixed,
	}, f.non, f.err
}

func TestIndeterminateMappedToDeny(t *testing.T) {
	t.Parallel()
	port, fail := policyseam.NewPort(fakeBoundary{
		result: policyeval.PolicyEvaluationResult{
			Outcome:     policyeval.OutcomeIndeterminate,
			ReasonCodes: []string{"UNKNOWN_INPUT"},
			InputDigest: strings.Repeat("a", 64),
			EvaluatedAt: time.Date(2026, 9, 7, 10, 0, 0, 0, time.UTC).Format(time.RFC3339Nano),
		},
	})
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	req := policyeval.PolicyEvaluationRequest{
		SubjectRef:  apimeta.TypedRef{APIVersion: "a/v1", Kind: "PrivilegedAccessRequest", Name: "p", UID: "p1"},
		Action:      "privilegedaccessrequest.submit",
		ProfileRefs: []apimeta.TypedRef{{APIVersion: "a/v1", Kind: "ApprovalPolicy", Name: "pol", UID: "pol1"}},
		RequestID:   "req",
	}
	out, fail := port.EvaluatePrivileged(context.Background(), req)
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	if out.Outcome != "Deny" || len(out.ReasonCodes) != 1 || out.ReasonCodes[0] != "POLICY_INDETERMINATE_MAPPED_TO_DENY" {
		t.Fatalf("unexpected output: %#v", out)
	}
}

func TestPolicyseamIsSolePolicyevalImporter(t *testing.T) {
	t.Parallel()
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("caller")
	}
	govaccessDir := filepath.Clean(filepath.Join(filepath.Dir(thisFile), ".."))
	allowed := map[string]bool{
		"internal/govaccess/policyseam": true,
	}
	err := filepath.WalkDir(govaccessDir, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if d.IsDir() || !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}
		file, err := parser.ParseFile(token.NewFileSet(), path, nil, parser.AllErrors)
		if err != nil {
			return err
		}
		rel := filepath.ToSlash(strings.TrimPrefix(path, govaccessDir+"/"))
		pkgPath := "internal/govaccess/" + filepath.ToSlash(filepath.Dir(rel))
		for _, imp := range file.Imports {
			importPath := strings.Trim(imp.Path.Value, `"`)
			if strings.Contains(importPath, "/internal/policyeval") && !allowed[pkgPath] {
				t.Fatalf("package %s imports policyeval in %s", pkgPath, filepath.Base(path))
			}
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}
