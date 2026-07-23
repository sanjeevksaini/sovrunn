package apiconform

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/apischema"
	"github.com/sanjeevksaini/sovrunn/internal/apivalid"
)

// Fitness check IDs for F12-VERIFY-001 checks implemented in this file
// (task 16.1: checks 1, 1a, 3, 4, 9). Aggregation across 1–15 is task 16.5.
const (
	FitnessCheckExternalSchemaAnnotations     = "1"
	FitnessCheckFieldPolicyCoverage           = "1a"
	FitnessCheckMutableFieldOwnership         = "3"
	FitnessCheckUnknownAndDuplicateFields     = "4"
	FitnessCheckPublishedDefinitionsImmutable = "9"
)

// Stable fitness finding codes for schema/metadata/ownership checks.
const (
	CodeFitnessSchemaLoadFailed          = "FITNESS_SCHEMA_LOAD_FAILED"
	CodeFitnessAnnotationMissing         = "FITNESS_ANNOTATION_MISSING"
	CodeFitnessAnnotationInvalid         = "FITNESS_ANNOTATION_INVALID"
	CodeFitnessFieldPolicyMissing        = "FIELD_POLICY_MISSING"
	CodeFitnessOwnershipMissing          = "FITNESS_OWNERSHIP_MISSING"
	CodeFitnessOwnershipConflict         = "FITNESS_OWNERSHIP_CONFLICT"
	CodeFitnessConditionOwnerInvalid     = "FITNESS_CONDITION_OWNER_INVALID"
	CodeFitnessUnknownFieldNotRejected   = "FITNESS_UNKNOWN_FIELD_NOT_REJECTED"
	CodeFitnessDuplicateFieldNotRejected = "FITNESS_DUPLICATE_FIELD_NOT_REJECTED"
	CodeFitnessPublishedFieldMutable     = "FITNESS_PUBLISHED_FIELD_MUTABLE"
	CodeFitnessImmutableRecordMutable    = "FITNESS_IMMUTABLE_RECORD_MUTABLE"
)

// FitnessFinding is one fail-closed finding from an executable fitness check
// (architecture §8.2; F12-VERIFY-001).
type FitnessFinding struct {
	Check   string
	Schema  string
	Path    string
	Code    string
	Message string
}

// String returns a stable diagnostic representation (not logged at runtime;
// fitness checks are gate/test-only).
func (f FitnessFinding) String() string {
	return fmt.Sprintf("check=%s schema=%s path=%s code=%s message=%s",
		f.Check, f.Schema, f.Path, f.Code, f.Message)
}

// externalCanonicalSchemaFiles is the Matrix D eight-contract inventory.
// These are the "external schemas" that must declare profile/boundary/
// stability/allowed-scopes (F12-VERIFY-001 check 1).
var externalCanonicalSchemaFiles = []string{
	"project.json",
	"resource-pool.json",
	"discovered-database.json",
	"plugin-definition.json",
	"adapter-configuration.json",
	"placement-evaluation-request.json",
	"operation.json",
	"audit-event.json",
}

// commonSubSchemaFiles is the _common shared sub-schema inventory used by
// field-policy coverage (check 1a) and ownership (check 3).
var commonSubSchemaFiles = []string{
	"type-meta.json",
	"object-meta.json",
	"typed-ref.json",
	"scope-ref.json",
	"owner-ref.json",
	"condition.json",
	"problem.json",
	"violation.json",
	"page.json",
}

// ExternalCanonicalSchemaFiles returns a copy of the eight Matrix D contract
// schema filenames under api/schemas/.
func ExternalCanonicalSchemaFiles() []string {
	out := make([]string, len(externalCanonicalSchemaFiles))
	copy(out, externalCanonicalSchemaFiles)
	return out
}

// CommonSubSchemaFiles returns a copy of the _common sub-schema filenames.
func CommonSubSchemaFiles() []string {
	out := make([]string, len(commonSubSchemaFiles))
	copy(out, commonSubSchemaFiles)
	return out
}

// CheckExternalSchemaAnnotations implements F12-VERIFY-001 check 1:
// every external (canonical Matrix D) schema declares profile, boundary,
// stability, and allowed scopes with valid controlled vocabularies (D-08).
func CheckExternalSchemaAnnotations(schemasRoot string) []FitnessFinding {
	var findings []FitnessFinding
	for _, name := range externalCanonicalSchemaFiles {
		schemaID := CanonicalSchemasDir + "/" + name
		body, err := os.ReadFile(filepath.Join(schemasRoot, name))
		if err != nil {
			findings = append(findings, FitnessFinding{
				Check:   FitnessCheckExternalSchemaAnnotations,
				Schema:  schemaID,
				Path:    "/",
				Code:    CodeFitnessSchemaLoadFailed,
				Message: err.Error(),
			})
			continue
		}
		meta, issues := apischema.ReadAnnotations(body)
		if len(issues) != 0 {
			for _, issue := range issues {
				code := CodeFitnessAnnotationInvalid
				if issue.Code == apischema.CodeAnnotationMissing {
					code = CodeFitnessAnnotationMissing
				}
				findings = append(findings, FitnessFinding{
					Check:   FitnessCheckExternalSchemaAnnotations,
					Schema:  schemaID,
					Path:    issue.Path,
					Code:    code,
					Message: issue.Message,
				})
			}
			continue
		}
		if meta.Profile == "" || meta.Boundary == "" || meta.Stability == "" || len(meta.AllowedScopes) == 0 {
			findings = append(findings, FitnessFinding{
				Check:   FitnessCheckExternalSchemaAnnotations,
				Schema:  schemaID,
				Path:    "/",
				Code:    CodeFitnessAnnotationMissing,
				Message: "profile, boundary, stability, and allowed-scopes must all be present",
			})
		}
	}
	return sortFindings(findings)
}

// CheckFieldPolicyCoverage implements F12-VERIFY-001 check 1a / F12-SEC-001:
// every property crossing an API boundary explicitly declares
// x-sovrunn-field-policy with exactly the eight required fields; controlled
// values are valid; unknown policy fields and unknown x-sovrunn-* extensions
// fail closed. FEATURE-0012 uses no inheritance algorithm.
func CheckFieldPolicyCoverage(schemasRoot string) []FitnessFinding {
	var findings []FitnessFinding

	for _, name := range externalCanonicalSchemaFiles {
		schemaID := CanonicalSchemasDir + "/" + name
		body, err := os.ReadFile(filepath.Join(schemasRoot, name))
		if err != nil {
			findings = append(findings, FitnessFinding{
				Check:   FitnessCheckFieldPolicyCoverage,
				Schema:  schemaID,
				Path:    "/",
				Code:    CodeFitnessSchemaLoadFailed,
				Message: err.Error(),
			})
			continue
		}
		findings = append(findings, fieldPolicyFindingsForSchema(schemaID, body)...)
		// External schemas also require registered root annotations; unknown
		// x-sovrunn-* extensions fail closed via ReadAnnotations.
		_, annIssues := apischema.ReadAnnotations(body)
		for _, issue := range annIssues {
			if issue.Code == apischema.CodeUnknownExtension ||
				issue.Code == apischema.CodeFieldPolicyUnknownField ||
				issue.Code == apischema.CodeFieldPolicyInvalid {
				findings = append(findings, FitnessFinding{
					Check:   FitnessCheckFieldPolicyCoverage,
					Schema:  schemaID,
					Path:    issue.Path,
					Code:    issue.Code,
					Message: issue.Message,
				})
			}
		}
	}

	commonRoot := filepath.Join(schemasRoot, "_common")
	for _, name := range commonSubSchemaFiles {
		schemaID := CanonicalSchemasDir + "/_common/" + name
		body, err := os.ReadFile(filepath.Join(commonRoot, name))
		if err != nil {
			findings = append(findings, FitnessFinding{
				Check:   FitnessCheckFieldPolicyCoverage,
				Schema:  schemaID,
				Path:    "/",
				Code:    CodeFitnessSchemaLoadFailed,
				Message: err.Error(),
			})
			continue
		}
		findings = append(findings, fieldPolicyFindingsForSchema(schemaID, body)...)
		// _common schemas omit document-level profile annotations; still
		// fail closed on unknown x-sovrunn-* extensions.
		for _, issue := range scanUnknownSovrunnExtensions(body) {
			findings = append(findings, FitnessFinding{
				Check:   FitnessCheckFieldPolicyCoverage,
				Schema:  schemaID,
				Path:    issue.Path,
				Code:    issue.Code,
				Message: issue.Message,
			})
		}
	}

	return sortFindings(dedupeFindings(findings))
}

// CheckFieldPolicyCompleteness asserts every named property under "properties"
// explicitly declares a complete x-sovrunn-field-policy (D-08, F12-SEC-001).
// FEATURE-0012 uses no inheritance algorithm: policies must be present on each
// property schema object in the document under test.
func CheckFieldPolicyCompleteness(schema []byte) []apischema.SchemaIssue {
	var root any
	if err := json.Unmarshal(schema, &root); err != nil {
		return []apischema.SchemaIssue{{
			Path:    "/",
			Code:    apischema.CodeMalformedSchema,
			Message: "schema document is not valid JSON",
		}}
	}
	var issues []apischema.SchemaIssue
	walkFieldPolicyCompleteness(root, "", &issues)
	sort.SliceStable(issues, func(i, j int) bool {
		if issues[i].Path != issues[j].Path {
			return issues[i].Path < issues[j].Path
		}
		return issues[i].Code < issues[j].Code
	})
	return issues
}

// CheckMutableFieldOwnership implements F12-VERIFY-001 check 3 / F12-OWNER-002:
// every mutable field and every condition has exactly one authoritative writer
// declared via x-sovrunn-field-policy.authorizedWriter.
func CheckMutableFieldOwnership(schemasRoot string) []FitnessFinding {
	var findings []FitnessFinding

	scan := func(schemaID string, body []byte, isConditionSchema bool) {
		findings = append(findings, ownershipFindingsForSchema(schemaID, body, isConditionSchema)...)
	}

	for _, name := range externalCanonicalSchemaFiles {
		schemaID := CanonicalSchemasDir + "/" + name
		body, err := os.ReadFile(filepath.Join(schemasRoot, name))
		if err != nil {
			findings = append(findings, FitnessFinding{
				Check:   FitnessCheckMutableFieldOwnership,
				Schema:  schemaID,
				Path:    "/",
				Code:    CodeFitnessSchemaLoadFailed,
				Message: err.Error(),
			})
			continue
		}
		scan(schemaID, body, false)
	}
	commonRoot := filepath.Join(schemasRoot, "_common")
	for _, name := range commonSubSchemaFiles {
		schemaID := CanonicalSchemasDir + "/_common/" + name
		body, err := os.ReadFile(filepath.Join(commonRoot, name))
		if err != nil {
			findings = append(findings, FitnessFinding{
				Check:   FitnessCheckMutableFieldOwnership,
				Schema:  schemaID,
				Path:    "/",
				Code:    CodeFitnessSchemaLoadFailed,
				Message: err.Error(),
			})
			continue
		}
		scan(schemaID, body, name == "condition.json")
	}
	return sortFindings(findings)
}

// CheckUnknownAndDuplicateFieldsFail implements F12-VERIFY-001 check 4:
// unknown and duplicate fields fail under strict DecodeJSON with stable codes
// and JSON Pointer paths (D-03, F12-VALIDATION-002/006).
//
// This check exercises the decode grammar directly; it does not log request
// bodies or field values (F12-SEC-003).
func CheckUnknownAndDuplicateFieldsFail() []FitnessFinding {
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

	var unknownDst sample
	unknownProb := apivalid.DecodeJSON(
		[]byte(`{"apiVersion":"core.sovrunn.io/v1alpha1","kind":"Project","metadata":{"name":"demo"},"spec":{},"extraField":true}`),
		lim,
		pol,
		&unknownDst,
	)
	if unknownProb == nil {
		findings = append(findings, FitnessFinding{
			Check:   FitnessCheckUnknownAndDuplicateFields,
			Schema:  "decode",
			Path:    "/extraField",
			Code:    CodeFitnessUnknownFieldNotRejected,
			Message: "DecodeJSON accepted an unknown field",
		})
	} else if unknownProb.Code != apiproblem.CodeUnknownField {
		findings = append(findings, FitnessFinding{
			Check:   FitnessCheckUnknownAndDuplicateFields,
			Schema:  "decode",
			Path:    "/",
			Code:    CodeFitnessUnknownFieldNotRejected,
			Message: fmt.Sprintf("DecodeJSON unknown-field Code=%q want %q", unknownProb.Code, apiproblem.CodeUnknownField),
		})
	} else if len(unknownProb.Violations) == 0 || unknownProb.Violations[0].Field != "/extraField" {
		field := ""
		if len(unknownProb.Violations) > 0 {
			field = unknownProb.Violations[0].Field
		}
		findings = append(findings, FitnessFinding{
			Check:   FitnessCheckUnknownAndDuplicateFields,
			Schema:  "decode",
			Path:    field,
			Code:    CodeFitnessUnknownFieldNotRejected,
			Message: "DecodeJSON unknown-field JSON Pointer must be /extraField",
		})
	}

	var dupDst sample
	dupProb := apivalid.DecodeJSON(
		[]byte(`{"apiVersion":"core.sovrunn.io/v1alpha1","kind":"Project","metadata":{"name":"a","name":"b"},"spec":{}}`),
		lim,
		pol,
		&dupDst,
	)
	if dupProb == nil {
		findings = append(findings, FitnessFinding{
			Check:   FitnessCheckUnknownAndDuplicateFields,
			Schema:  "decode",
			Path:    "/metadata/name",
			Code:    CodeFitnessDuplicateFieldNotRejected,
			Message: "DecodeJSON accepted a duplicate field",
		})
	} else if dupProb.Code != apiproblem.CodeDuplicateField {
		findings = append(findings, FitnessFinding{
			Check:   FitnessCheckUnknownAndDuplicateFields,
			Schema:  "decode",
			Path:    "/",
			Code:    CodeFitnessDuplicateFieldNotRejected,
			Message: fmt.Sprintf("DecodeJSON duplicate-field Code=%q want %q", dupProb.Code, apiproblem.CodeDuplicateField),
		})
	} else if len(dupProb.Violations) == 0 || dupProb.Violations[0].Field != "/metadata/name" {
		field := ""
		if len(dupProb.Violations) > 0 {
			field = dupProb.Violations[0].Field
		}
		findings = append(findings, FitnessFinding{
			Check:   FitnessCheckUnknownAndDuplicateFields,
			Schema:  "decode",
			Path:    field,
			Code:    CodeFitnessDuplicateFieldNotRejected,
			Message: "DecodeJSON duplicate-field JSON Pointer must be /metadata/name",
		})
	}

	return sortFindings(findings)
}

// CheckPublishedDefinitionsImmutable implements F12-VERIFY-001 check 9:
// published VersionedDefinition contract fields are immutable, and
// ImmutableRecord payload fields are immutable or append-only (never mutable).
func CheckPublishedDefinitionsImmutable(schemasRoot string) []FitnessFinding {
	var findings []FitnessFinding
	for _, name := range externalCanonicalSchemaFiles {
		schemaID := CanonicalSchemasDir + "/" + name
		body, err := os.ReadFile(filepath.Join(schemasRoot, name))
		if err != nil {
			findings = append(findings, FitnessFinding{
				Check:   FitnessCheckPublishedDefinitionsImmutable,
				Schema:  schemaID,
				Path:    "/",
				Code:    CodeFitnessSchemaLoadFailed,
				Message: err.Error(),
			})
			continue
		}
		meta, issues := apischema.ReadAnnotations(body)
		if len(issues) != 0 {
			// Annotation failures are owned by check 1; skip profile-specific
			// immutability when annotations are not trustworthy.
			continue
		}
		switch meta.Profile {
		case apimeta.ProfileVersionedDefinition:
			findings = append(findings, versionedDefinitionImmutabilityFindings(schemaID, body)...)
		case apimeta.ProfileImmutableRecord:
			findings = append(findings, immutableRecordMutabilityFindings(schemaID, body)...)
		}
	}
	return sortFindings(findings)
}

func fieldPolicyFindingsForSchema(schemaID string, body []byte) []FitnessFinding {
	issues := CheckFieldPolicyCompleteness(body)
	out := make([]FitnessFinding, 0, len(issues))
	for _, issue := range issues {
		code := issue.Code
		if code == "FIELD_POLICY_MISSING" {
			code = CodeFitnessFieldPolicyMissing
		}
		out = append(out, FitnessFinding{
			Check:   FitnessCheckFieldPolicyCoverage,
			Schema:  schemaID,
			Path:    issue.Path,
			Code:    code,
			Message: issue.Message,
		})
	}
	return out
}

func walkFieldPolicyCompleteness(node any, path string, issues *[]apischema.SchemaIssue) {
	obj, ok := node.(map[string]any)
	if !ok {
		return
	}

	if rawProps, hasProps := obj["properties"]; hasProps {
		props, ok := rawProps.(map[string]any)
		if !ok {
			*issues = append(*issues, apischema.SchemaIssue{
				Path:    joinFitnessPointer(path, "properties"),
				Code:    apischema.CodeMalformedSchema,
				Message: "properties must be an object",
			})
			return
		}
		names := make([]string, 0, len(props))
		for name := range props {
			names = append(names, name)
		}
		sort.Strings(names)
		for _, name := range names {
			propPath := joinFitnessPointer(joinFitnessPointer(path, "properties"), name)
			propSchema, ok := props[name].(map[string]any)
			if !ok {
				*issues = append(*issues, apischema.SchemaIssue{
					Path:    propPath,
					Code:    apischema.CodeMalformedSchema,
					Message: "property schema must be an object",
				})
				continue
			}
			rawPolicy, hasPolicy := propSchema[apischema.ExtFieldPolicy]
			if !hasPolicy {
				*issues = append(*issues, apischema.SchemaIssue{
					Path:    propPath,
					Code:    CodeFitnessFieldPolicyMissing,
					Message: "boundary-crossing property missing explicit x-sovrunn-field-policy",
				})
			} else if rawPolicy == nil {
				*issues = append(*issues, apischema.SchemaIssue{
					Path:    joinFitnessPointer(propPath, apischema.ExtFieldPolicy),
					Code:    apischema.CodeFieldPolicyInvalid,
					Message: "x-sovrunn-field-policy must be an object",
				})
			} else {
				// Reuse ReadAnnotations vocabulary checks by wrapping a tiny
				// annotated document around this property schema.
				probe := map[string]any{
					"type":                     "object",
					"x-sovrunn-profile":        "EmbeddedValue",
					"x-sovrunn-boundary":       "customer-facing",
					"x-sovrunn-allowed-scopes": []any{"Platform"},
					"x-sovrunn-stability":      "stable",
					"properties": map[string]any{
						"probe": propSchema,
					},
				}
				raw, err := json.Marshal(probe)
				if err != nil {
					*issues = append(*issues, apischema.SchemaIssue{
						Path:    propPath,
						Code:    apischema.CodeMalformedSchema,
						Message: "could not marshal field-policy probe",
					})
				} else {
					_, annIssues := apischema.ReadAnnotations(raw)
					for _, issue := range annIssues {
						mapped := issue
						const probePrefix = "/properties/probe"
						if len(issue.Path) >= len(probePrefix) && issue.Path[:len(probePrefix)] == probePrefix {
							mapped.Path = propPath + issue.Path[len(probePrefix):]
						}
						*issues = append(*issues, mapped)
					}
				}
			}
			walkFieldPolicyCompleteness(propSchema, propPath, issues)
		}
	}

	if rawItems, hasItems := obj["items"]; hasItems {
		walkFieldPolicyCompleteness(rawItems, joinFitnessPointer(path, "items"), issues)
	}
	if rawAP, hasAP := obj["additionalProperties"]; hasAP {
		if _, isBool := rawAP.(bool); !isBool {
			walkFieldPolicyCompleteness(rawAP, joinFitnessPointer(path, "additionalProperties"), issues)
		}
	}
}

func ownershipFindingsForSchema(schemaID string, body []byte, isConditionSchema bool) []FitnessFinding {
	var root any
	if err := json.Unmarshal(body, &root); err != nil {
		return []FitnessFinding{{
			Check:   FitnessCheckMutableFieldOwnership,
			Schema:  schemaID,
			Path:    "/",
			Code:    apischema.CodeMalformedSchema,
			Message: "schema document is not valid JSON",
		}}
	}
	var findings []FitnessFinding
	walkOwnership(root, "", schemaID, isConditionSchema, &findings)
	return findings
}

func walkOwnership(node any, path, schemaID string, underCondition bool, findings *[]FitnessFinding) {
	obj, ok := node.(map[string]any)
	if !ok {
		return
	}

	if rawProps, hasProps := obj["properties"]; hasProps {
		props, ok := rawProps.(map[string]any)
		if !ok {
			return
		}
		names := make([]string, 0, len(props))
		for name := range props {
			names = append(names, name)
		}
		sort.Strings(names)
		for _, name := range names {
			propPath := joinFitnessPointer(joinFitnessPointer(path, "properties"), name)
			propSchema, ok := props[name].(map[string]any)
			if !ok {
				continue
			}
			isConditionProp := underCondition || name == "conditions" || strings.HasSuffix(propPath, "/conditions")
			if rawPolicy, hasPolicy := propSchema[apischema.ExtFieldPolicy]; hasPolicy {
				*findings = append(*findings, ownershipFindingsForPolicy(schemaID, propPath, rawPolicy, isConditionProp)...)
			} else if isMutableCandidate(propSchema) || isConditionProp {
				// Mutable/condition properties without an explicit owner fail closed.
				*findings = append(*findings, FitnessFinding{
					Check:   FitnessCheckMutableFieldOwnership,
					Schema:  schemaID,
					Path:    propPath,
					Code:    CodeFitnessOwnershipMissing,
					Message: "mutable field or condition missing authorizedWriter via x-sovrunn-field-policy",
				})
			}
			walkOwnership(propSchema, propPath, schemaID, isConditionProp, findings)
		}
	}

	if rawItems, hasItems := obj["items"]; hasItems {
		walkOwnership(rawItems, joinFitnessPointer(path, "items"), schemaID, underCondition, findings)
	}
	if rawAP, hasAP := obj["additionalProperties"]; hasAP {
		if _, isBool := rawAP.(bool); !isBool {
			walkOwnership(rawAP, joinFitnessPointer(path, "additionalProperties"), schemaID, underCondition, findings)
		}
	}
}

func ownershipFindingsForPolicy(schemaID, propPath string, rawPolicy any, isCondition bool) []FitnessFinding {
	pol, ok := rawPolicy.(map[string]any)
	if !ok || pol == nil {
		return []FitnessFinding{{
			Check:   FitnessCheckMutableFieldOwnership,
			Schema:  schemaID,
			Path:    joinFitnessPointer(propPath, apischema.ExtFieldPolicy),
			Code:    CodeFitnessOwnershipMissing,
			Message: "x-sovrunn-field-policy must be an object declaring authorizedWriter",
		}}
	}

	writerRaw, hasWriter := pol["authorizedWriter"]
	writer, writerOK := writerRaw.(string)
	writer = strings.TrimSpace(writer)
	if !hasWriter || !writerOK || writer == "" {
		return []FitnessFinding{{
			Check:   FitnessCheckMutableFieldOwnership,
			Schema:  schemaID,
			Path:    joinFitnessPointer(joinFitnessPointer(propPath, apischema.ExtFieldPolicy), "authorizedWriter"),
			Code:    CodeFitnessOwnershipMissing,
			Message: "authorizedWriter is required (exactly one owner)",
		}}
	}

	// authorizedWriter must be a single string owner — arrays are a conflict.
	if _, isArr := writerRaw.([]any); isArr {
		return []FitnessFinding{{
			Check:   FitnessCheckMutableFieldOwnership,
			Schema:  schemaID,
			Path:    joinFitnessPointer(joinFitnessPointer(propPath, apischema.ExtFieldPolicy), "authorizedWriter"),
			Code:    CodeFitnessOwnershipConflict,
			Message: "authorizedWriter must name exactly one owner, not a list",
		}}
	}

	mutability, _ := pol["mutability"].(string)
	needsOwner := mutability == apischema.MutabilityMutable ||
		mutability == apischema.MutabilitySystemOnly ||
		mutability == apischema.MutabilityAppendOnly ||
		isCondition
	if needsOwner && writer == "" {
		return []FitnessFinding{{
			Check:   FitnessCheckMutableFieldOwnership,
			Schema:  schemaID,
			Path:    joinFitnessPointer(joinFitnessPointer(propPath, apischema.ExtFieldPolicy), "authorizedWriter"),
			Code:    CodeFitnessOwnershipMissing,
			Message: "mutable field or condition requires exactly one authorizedWriter",
		}}
	}

	if isCondition && writer != apischema.WriterStatusOwner {
		return []FitnessFinding{{
			Check:   FitnessCheckMutableFieldOwnership,
			Schema:  schemaID,
			Path:    joinFitnessPointer(joinFitnessPointer(propPath, apischema.ExtFieldPolicy), "authorizedWriter"),
			Code:    CodeFitnessConditionOwnerInvalid,
			Message: fmt.Sprintf("condition owner must be %q, got %q", apischema.WriterStatusOwner, writer),
		}}
	}
	return nil
}

func isMutableCandidate(propSchema map[string]any) bool {
	rawPolicy, ok := propSchema[apischema.ExtFieldPolicy]
	if !ok {
		// Without a policy we cannot prove mutability; ownership check only
		// flags missing policy for properties that declare conditions above.
		return false
	}
	pol, ok := rawPolicy.(map[string]any)
	if !ok {
		return true
	}
	mutability, _ := pol["mutability"].(string)
	return mutability == apischema.MutabilityMutable ||
		mutability == apischema.MutabilitySystemOnly ||
		mutability == apischema.MutabilityAppendOnly
}

func versionedDefinitionImmutabilityFindings(schemaID string, body []byte) []FitnessFinding {
	var root map[string]any
	if err := json.Unmarshal(body, &root); err != nil {
		return []FitnessFinding{{
			Check:   FitnessCheckPublishedDefinitionsImmutable,
			Schema:  schemaID,
			Path:    "/",
			Code:    apischema.CodeMalformedSchema,
			Message: "schema document is not valid JSON",
		}}
	}
	props, _ := root["properties"].(map[string]any)
	spec, _ := props["spec"].(map[string]any)
	specProps, _ := spec["properties"].(map[string]any)
	if len(specProps) == 0 {
		return []FitnessFinding{{
			Check:   FitnessCheckPublishedDefinitionsImmutable,
			Schema:  schemaID,
			Path:    "/properties/spec/properties",
			Code:    CodeFitnessPublishedFieldMutable,
			Message: "VersionedDefinition spec must declare published contract fields",
		}}
	}

	// Published contract fields (everything under spec except publication
	// lifecycle controls) must be immutable. publicationState may remain
	// mutable so Draft → Published transitions are expressible.
	var findings []FitnessFinding
	names := make([]string, 0, len(specProps))
	for name := range specProps {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		if name == "publicationState" {
			continue
		}
		propPath := "/properties/spec/properties/" + name
		propSchema, ok := specProps[name].(map[string]any)
		if !ok {
			continue
		}
		pol, _ := propSchema[apischema.ExtFieldPolicy].(map[string]any)
		mutability, _ := pol["mutability"].(string)
		if mutability != apischema.MutabilityImmutable {
			findings = append(findings, FitnessFinding{
				Check:   FitnessCheckPublishedDefinitionsImmutable,
				Schema:  schemaID,
				Path:    joinFitnessPointer(joinFitnessPointer(propPath, apischema.ExtFieldPolicy), "mutability"),
				Code:    CodeFitnessPublishedFieldMutable,
				Message: fmt.Sprintf("published VersionedDefinition field %q must be immutable, got %q", name, mutability),
			})
		}
	}
	return findings
}

func immutableRecordMutabilityFindings(schemaID string, body []byte) []FitnessFinding {
	var root any
	if err := json.Unmarshal(body, &root); err != nil {
		return []FitnessFinding{{
			Check:   FitnessCheckPublishedDefinitionsImmutable,
			Schema:  schemaID,
			Path:    "/",
			Code:    apischema.CodeMalformedSchema,
			Message: "schema document is not valid JSON",
		}}
	}
	var findings []FitnessFinding
	walkImmutableRecordPolicies(root, "", schemaID, &findings)
	return findings
}

func walkImmutableRecordPolicies(node any, path, schemaID string, findings *[]FitnessFinding) {
	obj, ok := node.(map[string]any)
	if !ok {
		return
	}
	if rawProps, hasProps := obj["properties"]; hasProps {
		props, ok := rawProps.(map[string]any)
		if !ok {
			return
		}
		names := make([]string, 0, len(props))
		for name := range props {
			names = append(names, name)
		}
		sort.Strings(names)
		for _, name := range names {
			propPath := joinFitnessPointer(joinFitnessPointer(path, "properties"), name)
			propSchema, ok := props[name].(map[string]any)
			if !ok {
				continue
			}
			// metadata may carry limited system mutability for identity
			// bookkeeping; record payload fields must not be mutable.
			if name != "metadata" && name != "apiVersion" && name != "kind" {
				if pol, ok := propSchema[apischema.ExtFieldPolicy].(map[string]any); ok {
					mutability, _ := pol["mutability"].(string)
					if mutability == apischema.MutabilityMutable {
						*findings = append(*findings, FitnessFinding{
							Check:   FitnessCheckPublishedDefinitionsImmutable,
							Schema:  schemaID,
							Path:    joinFitnessPointer(joinFitnessPointer(propPath, apischema.ExtFieldPolicy), "mutability"),
							Code:    CodeFitnessImmutableRecordMutable,
							Message: fmt.Sprintf("ImmutableRecord field %q must not be mutable", name),
						})
					}
				}
			}
			walkImmutableRecordPolicies(propSchema, propPath, schemaID, findings)
		}
	}
	if rawItems, hasItems := obj["items"]; hasItems {
		walkImmutableRecordPolicies(rawItems, joinFitnessPointer(path, "items"), schemaID, findings)
	}
}

func scanUnknownSovrunnExtensions(schema []byte) []apischema.SchemaIssue {
	var root any
	if err := json.Unmarshal(schema, &root); err != nil {
		return []apischema.SchemaIssue{{
			Path:    "/",
			Code:    apischema.CodeMalformedSchema,
			Message: "schema document is not valid JSON",
		}}
	}
	var issues []apischema.SchemaIssue
	walkUnknownExtensions(root, "", &issues)
	sort.SliceStable(issues, func(i, j int) bool {
		if issues[i].Path != issues[j].Path {
			return issues[i].Path < issues[j].Path
		}
		return issues[i].Code < issues[j].Code
	})
	return issues
}

func walkUnknownExtensions(node any, path string, issues *[]apischema.SchemaIssue) {
	obj, ok := node.(map[string]any)
	if !ok {
		return
	}
	keys := make([]string, 0, len(obj))
	for k := range obj {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, key := range keys {
		childPath := joinFitnessPointer(path, key)
		if strings.HasPrefix(key, "x-sovrunn-") && !apischema.IsRegisteredExtension(key) {
			*issues = append(*issues, apischema.SchemaIssue{
				Path:    childPath,
				Code:    apischema.CodeUnknownExtension,
				Message: fmt.Sprintf("unknown x-sovrunn-* extension %q", key),
			})
			continue
		}
		if key == "properties" {
			props, ok := obj[key].(map[string]any)
			if !ok {
				continue
			}
			propNames := make([]string, 0, len(props))
			for name := range props {
				propNames = append(propNames, name)
			}
			sort.Strings(propNames)
			for _, name := range propNames {
				// Property names are never extension keywords.
				walkUnknownExtensions(props[name], joinFitnessPointer(childPath, name), issues)
			}
			continue
		}
		if key == apischema.ExtFieldPolicy {
			// Field-policy object keys are validated elsewhere; do not treat
			// policy field names as schema extension keywords.
			continue
		}
		walkUnknownExtensions(obj[key], childPath, issues)
	}
}

func joinFitnessPointer(base, key string) string {
	if base == "" || base == "/" {
		return "/" + escapeFitnessPointerToken(key)
	}
	return base + "/" + escapeFitnessPointerToken(key)
}

func escapeFitnessPointerToken(s string) string {
	r := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '~':
			r = append(r, '~', '0')
		case '/':
			r = append(r, '~', '1')
		default:
			r = append(r, s[i])
		}
	}
	return string(r)
}

func sortFindings(findings []FitnessFinding) []FitnessFinding {
	sort.SliceStable(findings, func(i, j int) bool {
		if findings[i].Check != findings[j].Check {
			return findings[i].Check < findings[j].Check
		}
		if findings[i].Schema != findings[j].Schema {
			return findings[i].Schema < findings[j].Schema
		}
		if findings[i].Path != findings[j].Path {
			return findings[i].Path < findings[j].Path
		}
		return findings[i].Code < findings[j].Code
	})
	return findings
}

func dedupeFindings(findings []FitnessFinding) []FitnessFinding {
	if len(findings) == 0 {
		return findings
	}
	type key struct{ check, schema, path, code, message string }
	seen := make(map[key]struct{}, len(findings))
	out := make([]FitnessFinding, 0, len(findings))
	for _, f := range findings {
		k := key{f.Check, f.Schema, f.Path, f.Code, f.Message}
		if _, ok := seen[k]; ok {
			continue
		}
		seen[k] = struct{}{}
		out = append(out, f)
	}
	return out
}
