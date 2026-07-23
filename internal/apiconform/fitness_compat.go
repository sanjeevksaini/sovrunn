package apiconform

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/sanjeevksaini/sovrunn/internal/apicond"
	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/apiref"
	"github.com/sanjeevksaini/sovrunn/internal/apischema"
	"github.com/sanjeevksaini/sovrunn/internal/apivalid"
)

// Fitness check IDs for F12-VERIFY-001 checks implemented in this file
// (task 16.3: checks 10, 11, 12, 13). Aggregation across 1–15 is task 16.5.
const (
	FitnessCheckSchemaCompatibility           = "10"
	FitnessCheckSizesBounded                  = "11"
	FitnessCheckStableCodesAndJSONPointers    = "12"
	FitnessCheckGeneratedArtifactsMatchSchema = "13"
)

// Stable fitness finding codes for compatibility / limits / type-binding checks.
const (
	CodeFitnessBaselineIntegrityFailed  = "FITNESS_BASELINE_INTEGRITY_FAILED"
	CodeFitnessBaselineApprovalFailed   = "FITNESS_BASELINE_APPROVAL_FAILED"
	CodeFitnessBreakingChangeUndetected = "FITNESS_BREAKING_CHANGE_UNDETECTED"
	CodeFitnessUnapprovedBreakingChange = "FITNESS_UNAPPROVED_BREAKING_CHANGE"
	CodeFitnessSchemaDiffLoadFailed     = "FITNESS_SCHEMA_DIFF_LOAD_FAILED"
	CodeFitnessLimitsMismatch           = "FITNESS_LIMITS_MISMATCH"
	CodeFitnessSizeNotRejected          = "FITNESS_SIZE_NOT_REJECTED"
	CodeFitnessStableCodeMissing        = "FITNESS_STABLE_CODE_MISSING"
	CodeFitnessJSONPointerMissing       = "FITNESS_JSON_POINTER_MISSING"
	CodeFitnessTypeBindingMismatch      = "FITNESS_TYPE_BINDING_MISMATCH"
	CodeFitnessTypeBindingCoverage      = "FITNESS_TYPE_BINDING_COVERAGE"
	CodeFitnessDeliberateMismatchMissed = "FITNESS_DELIBERATE_MISMATCH_MISSED"
)

// BaselineSchemasDir is the repository-relative directory holding frozen
// schema snapshots, BASELINE_MANIFEST.json, and BASELINE_APPROVALS.json (D-11).
const BaselineSchemasDir = CanonicalSchemasDir + "/baseline"

// expectedDefaultLimits is the reviewed D-06 / F12-VALIDATION-007 table used
// by check 11. Values MUST stay aligned with apivalid.DefaultLimits().
var expectedDefaultLimits = apivalid.Limits{
	MaxObjectBytes:        1_048_576,
	MaxNestingDepth:       32,
	MaxLabels:             64,
	MaxLabelKeyChars:      63,
	MaxLabelValueChars:    253,
	MaxAnnotationsBytes:   262_144,
	MaxConditions:         32,
	MaxReferencesPerField: 64,
	MaxViolations:         100,
	DefaultPageSize:       50,
	MaxPageSize:           200,
}

// CheckSchemaCompatibilityDetectsBreaking implements F12-VERIFY-001 check 10:
// schema-diff detects breaking changes against api/schemas/baseline/*,
// VerifyBaselineIntegrity detects silent baseline tampering, and
// VerifyBaselineApproval rejects co-edited baseline+manifest without
// recorded approval evidence (D-11, F12-EVOLVE-002).
//
// Live repository state must pass integrity and approval. Current schemas
// compared to baseline must not introduce Breaking deltas without an
// approved baseline update. A built-in remove-field probe proves
// ClassifyChange still detects Breaking changes.
func CheckSchemaCompatibilityDetectsBreaking(moduleRoot string) []FitnessFinding {
	var findings []FitnessFinding

	baselineDir := filepath.Join(moduleRoot, BaselineSchemasDir)
	manifestPath := filepath.Join(baselineDir, apischema.BaselineManifestFileName)
	approvalsPath := filepath.Join(baselineDir, apischema.BaselineApprovalsFileName)

	if err := apischema.VerifyBaselineIntegrity(manifestPath, baselineDir); err != nil {
		findings = append(findings, FitnessFinding{
			Check:   FitnessCheckSchemaCompatibility,
			Schema:  BaselineSchemasDir,
			Path:    "/" + apischema.BaselineManifestFileName,
			Code:    CodeFitnessBaselineIntegrityFailed,
			Message: err.Error(),
		})
	}

	if err := apischema.VerifyBaselineApproval(approvalsPath, manifestPath, baselineDir); err != nil {
		findings = append(findings, FitnessFinding{
			Check:   FitnessCheckSchemaCompatibility,
			Schema:  BaselineSchemasDir,
			Path:    "/" + apischema.BaselineApprovalsFileName,
			Code:    CodeFitnessBaselineApprovalFailed,
			Message: err.Error(),
		})
	}

	// Built-in probe: remove-field MUST classify as Breaking (gate cannot be silent).
	oldProbe := []byte(`{"type":"object","properties":{"name":{"type":"string"},"label":{"type":"string"}}}`)
	newProbe := []byte(`{"type":"object","properties":{"name":{"type":"string"}}}`)
	probeChanges := apischema.ClassifyChange(oldProbe, newProbe)
	if !hasBreakingChange(probeChanges) {
		findings = append(findings, FitnessFinding{
			Check:   FitnessCheckSchemaCompatibility,
			Schema:  "classify-change-probe",
			Path:    "/properties/label",
			Code:    CodeFitnessBreakingChangeUndetected,
			Message: "ClassifyChange did not report Breaking for remove-field",
		})
	}

	findings = append(findings, baselineCurrentDiffFindings(moduleRoot)...)
	return sortFindings(findings)
}

// CheckSizesBounded implements F12-VERIFY-001 check 11: object, metadata,
// condition, violation, reference, and page sizes are bounded (D-06,
// F12-VALIDATION-007, F12-LIST-002). Over-limit inputs must reject with
// stable codes and JSON Pointer paths. Request bodies and field values are
// not logged (F12-SEC-003).
func CheckSizesBounded() []FitnessFinding {
	var findings []FitnessFinding

	got := apivalid.DefaultLimits()
	if got != expectedDefaultLimits {
		findings = append(findings, FitnessFinding{
			Check:   FitnessCheckSizesBounded,
			Schema:  "apivalid.DefaultLimits",
			Path:    "/",
			Code:    CodeFitnessLimitsMismatch,
			Message: fmt.Sprintf("DefaultLimits = %#v, want %#v", got, expectedDefaultLimits),
		})
	}
	if got.DefaultPageSize <= 0 || got.MaxPageSize < got.DefaultPageSize {
		findings = append(findings, FitnessFinding{
			Check:   FitnessCheckSizesBounded,
			Schema:  "apivalid.DefaultLimits",
			Path:    "/page",
			Code:    CodeFitnessLimitsMismatch,
			Message: fmt.Sprintf("page bounds invalid: default=%d max=%d", got.DefaultPageSize, got.MaxPageSize),
		})
	}
	if got.MaxObjectBytes <= 0 || got.MaxConditions <= 0 || got.MaxReferencesPerField <= 0 ||
		got.MaxViolations <= 0 || got.MaxLabels <= 0 || got.MaxAnnotationsBytes <= 0 {
		findings = append(findings, FitnessFinding{
			Check:   FitnessCheckSizesBounded,
			Schema:  "apivalid.DefaultLimits",
			Path:    "/",
			Code:    CodeFitnessLimitsMismatch,
			Message: "object/metadata/condition/violation/reference limits must be positive",
		})
	}

	findings = append(findings, checkObjectAndNestingBounds()...)
	findings = append(findings, checkMetadataConditionBounds()...)
	findings = append(findings, checkReferenceAndViolationBounds()...)

	return sortFindings(findings)
}

// CheckErrorsUseStableCodesAndJSONPointers implements F12-VERIFY-001 check 12:
// errors use registered stable ErrorCode / ViolationCode values and RFC 6901
// JSON Pointer field paths (F12-ERROR-001/003, F12-VALIDATION-006).
//
// Correlation fields (requestId/instance) are optional caller attachments and
// are intentionally not required here. Secrets, credentials, tokens, and
// inaccessible resource details MUST NOT appear in Problem/Violation messages.
func CheckErrorsUseStableCodesAndJSONPointers() []FitnessFinding {
	var findings []FitnessFinding

	codes := apiproblem.AllErrorCodes()
	if len(codes) == 0 {
		findings = append(findings, FitnessFinding{
			Check:   FitnessCheckStableCodesAndJSONPointers,
			Schema:  "apiproblem.ErrorCode",
			Path:    "/",
			Code:    CodeFitnessStableCodeMissing,
			Message: "AllErrorCodes is empty",
		})
	}
	for _, code := range codes {
		if !code.Valid() {
			findings = append(findings, FitnessFinding{
				Check:   FitnessCheckStableCodesAndJSONPointers,
				Schema:  "apiproblem.ErrorCode",
				Path:    "/" + string(code),
				Code:    CodeFitnessStableCodeMissing,
				Message: "ErrorCode is not registered as Valid",
			})
		}
		p := apiproblem.New(code)
		if p == nil || p.Code != code || p.Type == "" || p.Title == "" || p.Status == 0 {
			findings = append(findings, FitnessFinding{
				Check:   FitnessCheckStableCodesAndJSONPointers,
				Schema:  "apiproblem.ErrorCode",
				Path:    "/" + string(code),
				Code:    CodeFitnessStableCodeMissing,
				Message: "New(code) must populate type/title/status/code",
			})
		}
	}

	for _, code := range apiproblem.AllViolationCodes() {
		if !code.Valid() {
			findings = append(findings, FitnessFinding{
				Check:   FitnessCheckStableCodesAndJSONPointers,
				Schema:  "apiproblem.ViolationCode",
				Path:    "/" + string(code),
				Code:    CodeFitnessStableCodeMissing,
				Message: "ViolationCode is not registered as Valid",
			})
		}
	}

	findings = append(findings, checkDecodeErrorPointers()...)
	findings = append(findings, checkOperationScopeMismatchPointer()...)

	return sortFindings(findings)
}

// CheckGeneratedArtifactsMatchCanonicalSchema implements F12-VERIFY-001
// check 13: registered TypeBindings match canonical schemas via
// VerifyGoTypeAgainstSchema (D-01b). Fixture round-tripping is not sufficient
// proof. A deliberate mismatch probe proves the check rejects bad Go types.
func CheckGeneratedArtifactsMatchCanonicalSchema(moduleRoot string) []FitnessFinding {
	var findings []FitnessFinding

	if len(TypeBindings) == 0 {
		findings = append(findings, FitnessFinding{
			Check:   FitnessCheckGeneratedArtifactsMatchSchema,
			Schema:  "TypeBindings",
			Path:    "/",
			Code:    CodeFitnessTypeBindingCoverage,
			Message: "TypeBindings registry is empty",
		})
		return sortFindings(findings)
	}

	wantCount := len(externalCanonicalSchemaFiles) + len(commonSubSchemaFiles)
	if len(TypeBindings) != wantCount {
		findings = append(findings, FitnessFinding{
			Check:   FitnessCheckGeneratedArtifactsMatchSchema,
			Schema:  "TypeBindings",
			Path:    "/",
			Code:    CodeFitnessTypeBindingCoverage,
			Message: fmt.Sprintf("TypeBindings count=%d want %d (8 contracts + 9 _common)", len(TypeBindings), wantCount),
		})
	}

	seen := make(map[string]struct{}, len(TypeBindings))
	for _, b := range TypeBindings {
		if b.SchemaPath == "" || b.GoType == nil {
			findings = append(findings, FitnessFinding{
				Check:   FitnessCheckGeneratedArtifactsMatchSchema,
				Schema:  b.SchemaPath,
				Path:    "/",
				Code:    CodeFitnessTypeBindingCoverage,
				Message: "TypeBinding SchemaPath/GoType must be set",
			})
			continue
		}
		if _, dup := seen[b.SchemaPath]; dup {
			findings = append(findings, FitnessFinding{
				Check:   FitnessCheckGeneratedArtifactsMatchSchema,
				Schema:  b.SchemaPath,
				Path:    "/",
				Code:    CodeFitnessTypeBindingCoverage,
				Message: "duplicate TypeBinding SchemaPath",
			})
			continue
		}
		seen[b.SchemaPath] = struct{}{}

		schemaPath := filepath.Join(moduleRoot, filepath.FromSlash(b.SchemaPath))
		schema, err := os.ReadFile(schemaPath)
		if err != nil {
			findings = append(findings, FitnessFinding{
				Check:   FitnessCheckGeneratedArtifactsMatchSchema,
				Schema:  b.SchemaPath,
				Path:    "/",
				Code:    CodeFitnessSchemaLoadFailed,
				Message: err.Error(),
			})
			continue
		}
		if support := apischema.ValidateSchemaSupport(schema); len(support) > 0 {
			findings = append(findings, FitnessFinding{
				Check:   FitnessCheckGeneratedArtifactsMatchSchema,
				Schema:  b.SchemaPath,
				Path:    support[0].Path,
				Code:    support[0].Code,
				Message: support[0].Message,
			})
			continue
		}
		issues := apischema.VerifyGoTypeAgainstSchema(schema, b.GoType)
		for _, issue := range issues {
			findings = append(findings, FitnessFinding{
				Check:   FitnessCheckGeneratedArtifactsMatchSchema,
				Schema:  b.SchemaPath,
				Path:    issue.Path,
				Code:    CodeFitnessTypeBindingMismatch,
				Message: fmt.Sprintf("%s: %s", issue.Code, issue.Message),
			})
		}
	}

	// Deliberate mismatch probe against page.json (D-01b supporting proof).
	pageSchema, err := os.ReadFile(filepath.Join(moduleRoot, "api/schemas/_common/page.json"))
	if err != nil {
		findings = append(findings, FitnessFinding{
			Check:   FitnessCheckGeneratedArtifactsMatchSchema,
			Schema:  CanonicalSchemasDir + "/_common/page.json",
			Path:    "/",
			Code:    CodeFitnessSchemaLoadFailed,
			Message: err.Error(),
		})
	} else {
		type mismatchedPage struct {
			NextPageToken int `json:"nextPageToken,omitempty"`
		}
		if issues := apischema.VerifyGoTypeAgainstSchema(pageSchema, reflect.TypeOf(mismatchedPage{})); len(issues) == 0 {
			findings = append(findings, FitnessFinding{
				Check:   FitnessCheckGeneratedArtifactsMatchSchema,
				Schema:  CanonicalSchemasDir + "/_common/page.json",
				Path:    "/properties/nextPageToken",
				Code:    CodeFitnessDeliberateMismatchMissed,
				Message: "VerifyGoTypeAgainstSchema accepted a deliberate Go-type mismatch",
			})
		}
	}

	return sortFindings(findings)
}

func baselineCurrentDiffFindings(moduleRoot string) []FitnessFinding {
	var findings []FitnessFinding
	baselineDir := filepath.Join(moduleRoot, BaselineSchemasDir)
	schemasRoot := filepath.Join(moduleRoot, CanonicalSchemasDir)

	entries, err := listBaselineSchemaRelPaths(baselineDir)
	if err != nil {
		return []FitnessFinding{{
			Check:   FitnessCheckSchemaCompatibility,
			Schema:  BaselineSchemasDir,
			Path:    "/",
			Code:    CodeFitnessSchemaDiffLoadFailed,
			Message: err.Error(),
		}}
	}

	for _, rel := range entries {
		oldRaw, err := os.ReadFile(filepath.Join(baselineDir, filepath.FromSlash(rel)))
		if err != nil {
			findings = append(findings, FitnessFinding{
				Check:   FitnessCheckSchemaCompatibility,
				Schema:  BaselineSchemasDir + "/" + rel,
				Path:    "/",
				Code:    CodeFitnessSchemaDiffLoadFailed,
				Message: err.Error(),
			})
			continue
		}
		newRaw, err := os.ReadFile(filepath.Join(schemasRoot, filepath.FromSlash(rel)))
		if err != nil {
			findings = append(findings, FitnessFinding{
				Check:   FitnessCheckSchemaCompatibility,
				Schema:  CanonicalSchemasDir + "/" + rel,
				Path:    "/",
				Code:    CodeFitnessSchemaDiffLoadFailed,
				Message: err.Error(),
			})
			continue
		}
		for _, change := range apischema.ClassifyChange(oldRaw, newRaw) {
			if change.Class != apischema.ChangeBreaking {
				continue
			}
			findings = append(findings, FitnessFinding{
				Check:   FitnessCheckSchemaCompatibility,
				Schema:  CanonicalSchemasDir + "/" + rel,
				Path:    change.Path,
				Code:    CodeFitnessUnapprovedBreakingChange,
				Message: fmt.Sprintf("%s: %s", change.Kind, change.Message),
			})
		}
	}
	return findings
}

func listBaselineSchemaRelPaths(baselineDir string) ([]string, error) {
	var out []string
	err := filepath.WalkDir(baselineDir, func(path string, d os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if d.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(baselineDir, path)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)
		base := filepath.Base(rel)
		if base == apischema.BaselineManifestFileName || base == apischema.BaselineApprovalsFileName {
			return nil
		}
		if !strings.HasSuffix(rel, ".json") {
			return nil
		}
		out = append(out, rel)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return out, nil
}

func hasBreakingChange(changes []apischema.Change) bool {
	for _, c := range changes {
		if c.Class == apischema.ChangeBreaking {
			return true
		}
	}
	return false
}

func checkObjectAndNestingBounds() []FitnessFinding {
	var findings []FitnessFinding
	lim := apivalid.Limits{MaxObjectBytes: 8, MaxNestingDepth: 2}
	pol := apivalid.PolicyFor(apivalid.ModeCreateRequest)

	var oversized any
	overProb := apivalid.DecodeJSON(
		[]byte(`{"apiVersion":"v1","kind":"Project","metadata":{"name":"x"},"spec":{}}`),
		lim,
		pol,
		&oversized,
	)
	if overProb == nil || overProb.Code != apiproblem.CodeRequestTooLarge {
		code := apiproblem.ErrorCode("")
		if overProb != nil {
			code = overProb.Code
		}
		findings = append(findings, FitnessFinding{
			Check:   FitnessCheckSizesBounded,
			Schema:  "decode",
			Path:    "/",
			Code:    CodeFitnessSizeNotRejected,
			Message: fmt.Sprintf("oversized body Code=%q want REQUEST_TOO_LARGE", code),
		})
	}

	var nested any
	nestProb := apivalid.DecodeJSON([]byte(`{"a":{"b":{"c":1}}}`), lim, pol, &nested)
	if nestProb == nil || nestProb.Code != apiproblem.CodeRequestTooLarge {
		code := apiproblem.ErrorCode("")
		if nestProb != nil {
			code = nestProb.Code
		}
		findings = append(findings, FitnessFinding{
			Check:   FitnessCheckSizesBounded,
			Schema:  "decode",
			Path:    "/",
			Code:    CodeFitnessSizeNotRejected,
			Message: fmt.Sprintf("over-nested body Code=%q want REQUEST_TOO_LARGE", code),
		})
	}
	return findings
}

func checkMetadataConditionBounds() []FitnessFinding {
	var findings []FitnessFinding
	limits := apivalid.DefaultLimits()
	limits.MaxLabels = 2
	limits.MaxConditions = 1
	limits.MaxAnnotationsBytes = 8
	stage := apivalid.NewCommonSemantic(limits, false)

	obj := &fitnessSemanticStub{
		apiVersion: "platform.sovrunn.io/v1alpha1",
		kind:       "Project",
		name:       "demo",
		scope: &apimeta.ScopeRef{TypedRef: apimeta.TypedRef{
			APIVersion: "platform.sovrunn.io/v1alpha1",
			Kind:       string(apimeta.ScopeTenant),
			Name:       "tenant-a",
			UID:        "tenant-uid-1",
		}},
		labels:      map[string]string{"a": "1", "b": "2", "c": "3"},
		annotations: map[string]string{"k": "too-large-value"},
		conditions: []apicond.Condition{
			{Type: "Ready", Status: apicond.ConditionTrue, Reason: "Succeeded"},
			{Type: "Valid", Status: apicond.ConditionTrue, Reason: "Checked"},
		},
	}

	violations, err := stage.Validate(context.Background(), obj)
	if err != nil {
		return []FitnessFinding{{
			Check:   FitnessCheckSizesBounded,
			Schema:  "apivalid.CommonSemantic",
			Path:    "/",
			Code:    CodeFitnessSizeNotRejected,
			Message: fmt.Sprintf("semantic over-limit returned error: %v", err),
		}}
	}
	if !hasFitnessViolationField(violations, "/metadata/labels") ||
		!hasFitnessViolationCode(violations, apiproblem.ViolationOutOfRange) {
		findings = append(findings, FitnessFinding{
			Check:   FitnessCheckSizesBounded,
			Schema:  "apivalid.CommonSemantic",
			Path:    "/metadata/labels",
			Code:    CodeFitnessSizeNotRejected,
			Message: "over MaxLabels must yield OUT_OF_RANGE at /metadata/labels",
		})
	}
	if !hasFitnessViolationField(violations, "/metadata/annotations") {
		findings = append(findings, FitnessFinding{
			Check:   FitnessCheckSizesBounded,
			Schema:  "apivalid.CommonSemantic",
			Path:    "/metadata/annotations",
			Code:    CodeFitnessSizeNotRejected,
			Message: "over MaxAnnotationsBytes must yield violation at /metadata/annotations",
		})
	}
	if !hasFitnessViolationField(violations, "/status/conditions") {
		findings = append(findings, FitnessFinding{
			Check:   FitnessCheckSizesBounded,
			Schema:  "apivalid.CommonSemantic",
			Path:    "/status/conditions",
			Code:    CodeFitnessSizeNotRejected,
			Message: "over MaxConditions must yield violation at /status/conditions",
		})
	}
	return findings
}

func checkReferenceAndViolationBounds() []FitnessFinding {
	var findings []FitnessFinding

	// Reference collection bound (MaxReferencesPerField / DefaultMaxRefs).
	refs := make(apiref.Refs, apiref.DefaultMaxRefs+1)
	for i := range refs {
		refs[i] = apiref.TypedRef{
			APIVersion: "platform.sovrunn.io/v1alpha1",
			Kind:       "Project",
			Name:       fmt.Sprintf("p-%d", i),
		}
	}
	c := apiref.Constraint{AllowedKinds: []string{"Project"}}
	issues := refs.Validate(c, "/spec/projectRefs", 0)
	if !hasRefIssueCode(issues, apiref.CodeRefsExceedLimit) {
		findings = append(findings, FitnessFinding{
			Check:   FitnessCheckSizesBounded,
			Schema:  "apiref.Refs",
			Path:    "/spec/projectRefs",
			Code:    CodeFitnessSizeNotRejected,
			Message: "references exceeding MaxReferencesPerField must yield REFS_EXCEED_LIMIT",
		})
	}

	// Violation cap: MaxViolations bounds ordinary findings.
	limits := apivalid.DefaultLimits()
	limits.MaxLabels = 1
	limits.MaxConditions = 1
	limits.MaxViolations = 1
	stage := apivalid.NewCommonSemantic(limits, false)
	obj := &fitnessSemanticStub{
		apiVersion: "platform.sovrunn.io/v1alpha1",
		kind:       "Project",
		name:       "BAD NAME", // invalid resource name
		scope: &apimeta.ScopeRef{TypedRef: apimeta.TypedRef{
			APIVersion: "platform.sovrunn.io/v1alpha1",
			Kind:       string(apimeta.ScopeTenant),
			Name:       "tenant-a",
			UID:        "tenant-uid-1",
		}},
		labels: map[string]string{"a": "1", "b": "2"},
		conditions: []apicond.Condition{
			{Type: "Ready", Status: apicond.ConditionTrue, Reason: "Ok"},
			{Type: "Valid", Status: apicond.ConditionTrue, Reason: "Ok"},
		},
	}
	violations, err := stage.Validate(context.Background(), obj)
	if err != nil {
		findings = append(findings, FitnessFinding{
			Check:   FitnessCheckSizesBounded,
			Schema:  "apivalid.CommonSemantic",
			Path:    "/violations",
			Code:    CodeFitnessSizeNotRejected,
			Message: fmt.Sprintf("MaxViolations probe returned error: %v", err),
		})
	} else if len(violations) > limits.MaxViolations {
		findings = append(findings, FitnessFinding{
			Check:   FitnessCheckSizesBounded,
			Schema:  "apivalid.CommonSemantic",
			Path:    "/violations",
			Code:    CodeFitnessSizeNotRejected,
			Message: fmt.Sprintf("MaxViolations not enforced: got %d want <= %d", len(violations), limits.MaxViolations),
		})
	} else if len(violations) == 0 {
		findings = append(findings, FitnessFinding{
			Check:   FitnessCheckSizesBounded,
			Schema:  "apivalid.CommonSemantic",
			Path:    "/violations",
			Code:    CodeFitnessSizeNotRejected,
			Message: "MaxViolations probe expected at least one ordinary violation",
		})
	}
	return findings
}

func checkDecodeErrorPointers() []FitnessFinding {
	var findings []FitnessFinding
	lim := apivalid.DefaultLimits()
	pol := apivalid.PolicyFor(apivalid.ModeCreateRequest)

	type sample struct {
		APIVersion string `json:"apiVersion"`
		Kind       string `json:"kind"`
		Metadata   struct {
			Name string `json:"name"`
		} `json:"metadata"`
		Spec map[string]any `json:"spec"`
	}

	cases := []struct {
		name      string
		raw       string
		wantCode  apiproblem.ErrorCode
		wantField string
	}{
		{
			name:      "unknown-field",
			raw:       `{"apiVersion":"v1","kind":"Project","metadata":{"name":"demo"},"spec":{},"extraField":true}`,
			wantCode:  apiproblem.CodeUnknownField,
			wantField: "/extraField",
		},
		{
			name:      "duplicate-field",
			raw:       `{"apiVersion":"v1","kind":"Project","metadata":{"name":"a","name":"b"},"spec":{}}`,
			wantCode:  apiproblem.CodeDuplicateField,
			wantField: "/metadata/name",
		},
	}

	for _, tc := range cases {
		var dst sample
		prob := apivalid.DecodeJSON([]byte(tc.raw), lim, pol, &dst)
		if prob == nil {
			findings = append(findings, FitnessFinding{
				Check:   FitnessCheckStableCodesAndJSONPointers,
				Schema:  "decode/" + tc.name,
				Path:    tc.wantField,
				Code:    CodeFitnessStableCodeMissing,
				Message: "expected Problem, got nil",
			})
			continue
		}
		if !prob.Code.Valid() || prob.Code != tc.wantCode {
			findings = append(findings, FitnessFinding{
				Check:   FitnessCheckStableCodesAndJSONPointers,
				Schema:  "decode/" + tc.name,
				Path:    "/",
				Code:    CodeFitnessStableCodeMissing,
				Message: fmt.Sprintf("Code=%q Valid=%v want %q", prob.Code, prob.Code.Valid(), tc.wantCode),
			})
		}
		if len(prob.Violations) == 0 {
			findings = append(findings, FitnessFinding{
				Check:   FitnessCheckStableCodesAndJSONPointers,
				Schema:  "decode/" + tc.name,
				Path:    tc.wantField,
				Code:    CodeFitnessJSONPointerMissing,
				Message: "Problem must carry violations with JSON Pointer paths",
			})
			continue
		}
		v := prob.Violations[0]
		if !isRFC6901JSONPointer(v.Field) || v.Field != tc.wantField {
			findings = append(findings, FitnessFinding{
				Check:   FitnessCheckStableCodesAndJSONPointers,
				Schema:  "decode/" + tc.name,
				Path:    v.Field,
				Code:    CodeFitnessJSONPointerMissing,
				Message: fmt.Sprintf("violation Field=%q want JSON Pointer %q", v.Field, tc.wantField),
			})
		}
		if !v.Code.Valid() {
			findings = append(findings, FitnessFinding{
				Check:   FitnessCheckStableCodesAndJSONPointers,
				Schema:  "decode/" + tc.name,
				Path:    v.Field,
				Code:    CodeFitnessStableCodeMissing,
				Message: fmt.Sprintf("violation Code=%q is not registered", v.Code),
			})
		}
	}
	return findings
}

func checkOperationScopeMismatchPointer() []FitnessFinding {
	op := apimeta.ScopeIdentity{Kind: apimeta.ScopeTenant, UID: "tenant-a"}
	target := apimeta.ScopeIdentity{Kind: apimeta.ScopeProject, UID: "project-b"}
	v := apivalid.CheckOperationTargetScopeMatch(op, target)
	if v == nil {
		return []FitnessFinding{{
			Check:   FitnessCheckStableCodesAndJSONPointers,
			Schema:  "apivalid.CheckOperationTargetScopeMatch",
			Path:    "/metadata/scopeRef",
			Code:    CodeFitnessStableCodeMissing,
			Message: "scope mismatch must produce a Violation",
		}}
	}
	var findings []FitnessFinding
	if v.Code != apiproblem.ViolationOperationTargetScopeMismatch || !v.Code.Valid() {
		findings = append(findings, FitnessFinding{
			Check:   FitnessCheckStableCodesAndJSONPointers,
			Schema:  "apivalid.CheckOperationTargetScopeMatch",
			Path:    v.Field,
			Code:    CodeFitnessStableCodeMissing,
			Message: fmt.Sprintf("Code=%q want OPERATION_TARGET_SCOPE_MISMATCH", v.Code),
		})
	}
	if !isRFC6901JSONPointer(v.Field) || v.Field != "/metadata/scopeRef" {
		findings = append(findings, FitnessFinding{
			Check:   FitnessCheckStableCodesAndJSONPointers,
			Schema:  "apivalid.CheckOperationTargetScopeMatch",
			Path:    v.Field,
			Code:    CodeFitnessJSONPointerMissing,
			Message: "OPERATION_TARGET_SCOPE_MISMATCH must use JSON Pointer /metadata/scopeRef",
		})
	}
	return findings
}

func isRFC6901JSONPointer(p string) bool {
	if p == "" {
		return true // whole document
	}
	if !strings.HasPrefix(p, "/") {
		return false
	}
	// Reject unescaped '~' sequences that are not ~0/~1 (loose structural check).
	for i := 0; i < len(p); i++ {
		if p[i] != '~' {
			continue
		}
		if i+1 >= len(p) || (p[i+1] != '0' && p[i+1] != '1') {
			return false
		}
		i++
	}
	return true
}

func hasFitnessViolationField(vs []apiproblem.Violation, field string) bool {
	for _, v := range vs {
		if v.Field == field {
			return true
		}
	}
	return false
}

func hasFitnessViolationCode(vs []apiproblem.Violation, code apiproblem.ViolationCode) bool {
	for _, v := range vs {
		if v.Code == code {
			return true
		}
	}
	return false
}

// fitnessSemanticStub is a minimal SemanticCarrier for limit probes in check 11.
// It is fitness-gate only and must not log or serialize secret-bearing payloads.
type fitnessSemanticStub struct {
	apiVersion  string
	kind        string
	name        string
	scope       *apimeta.ScopeRef
	owner       *apimeta.OwnerRef
	labels      map[string]string
	annotations map[string]string
	conditions  []apicond.Condition
	phase       string
}

var _ apivalid.SemanticCarrier = (*fitnessSemanticStub)(nil)

func (r *fitnessSemanticStub) APIVersion() string                   { return r.apiVersion }
func (r *fitnessSemanticStub) Kind() string                         { return r.kind }
func (r *fitnessSemanticStub) ResourceName() string                 { return r.name }
func (r *fitnessSemanticStub) GetScopeRef() *apimeta.ScopeRef       { return r.scope }
func (r *fitnessSemanticStub) GetOwnerRef() *apimeta.OwnerRef       { return r.owner }
func (r *fitnessSemanticStub) Labels() map[string]string            { return r.labels }
func (r *fitnessSemanticStub) Annotations() map[string]string       { return r.annotations }
func (r *fitnessSemanticStub) Conditions() []apicond.Condition      { return r.conditions }
func (r *fitnessSemanticStub) Phase() string                        { return r.phase }
func (r *fitnessSemanticStub) Profile() (apimeta.Profile, bool)     { return "", false }
func (r *fitnessSemanticStub) Boundary() (apimeta.Boundary, bool)   { return "", false }
func (r *fitnessSemanticStub) Stability() (apimeta.Stability, bool) { return "", false }
func (r *fitnessSemanticStub) DataClassification() (apimeta.DataClassification, bool) {
	return "", false
}
