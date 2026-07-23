package apischema

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
)

// Stable SchemaIssue codes for x-sovrunn-* annotation parsing (D-08,
// F12-NAMING-006, F12-SEC-001).
const (
	CodeAnnotationMissing       = "SCHEMA_ANNOTATION_MISSING"
	CodeAnnotationInvalid       = "SCHEMA_ANNOTATION_INVALID"
	CodeUnknownExtension        = "SCHEMA_UNKNOWN_EXTENSION"
	CodeFieldPolicyInvalid      = "SCHEMA_FIELD_POLICY_INVALID"
	CodeFieldPolicyUnknownField = "SCHEMA_FIELD_POLICY_UNKNOWN_FIELD"
)

// Extension keyword names (D-08). Exactly five registered extensions exist;
// unknown x-sovrunn-* names fail closed.
const (
	ExtProfile       = "x-sovrunn-profile"
	ExtBoundary      = "x-sovrunn-boundary"
	ExtAllowedScopes = "x-sovrunn-allowed-scopes"
	ExtStability     = "x-sovrunn-stability"
	ExtFieldPolicy   = "x-sovrunn-field-policy"
	extPrefixSovrunn = "x-sovrunn-"
)

// Required schema-level annotation keywords (F12-NAMING-006).
var requiredSchemaAnnotations = []string{
	ExtProfile,
	ExtBoundary,
	ExtAllowedScopes,
	ExtStability,
}

// Exact field-policy object keys (D-08, F12-SEC-001). No inheritance algorithm
// is applied in FEATURE-0012; when present, the object must carry exactly these
// eight fields and no others.
var requiredFieldPolicyKeys = []string{
	"classification",
	"authorizedWriter",
	"authorizedReaders",
	"mutability",
	"retention",
	"redaction",
	"residency",
	"auditRequired",
}

// Field-policy controlled vocabularies beyond apimeta.DataClassification.
// Writer/reader/mutability values align with Matrix C2 ownership classes and
// Matrix C1 boundary audiences; retention/redaction/residency are closed sets
// used by FEATURE-0012 schema annotations.
const (
	WriterCreator     = "creator"
	WriterSystem      = "system"
	WriterSpecOwner   = "spec-owner"
	WriterStatusOwner = "status-owner"

	ReaderCustomer   = "customer"
	ReaderOperator   = "operator"
	ReaderInternal   = "internal"
	ReaderAdapter    = "adapter"
	ReaderPlugin     = "plugin"
	ReaderGovernance = "governance"
	ReaderSystem     = "system"

	MutabilityImmutable  = "immutable"
	MutabilityMutable    = "mutable"
	MutabilityAppendOnly = "append-only"
	MutabilitySystemOnly = "system-only"

	RetentionNone     = "none"
	RetentionStandard = "standard"
	RetentionExtended = "extended"
	RetentionArchival = "archival"

	RedactionNone   = "none"
	RedactionRedact = "redact"
	RedactionOmit   = "omit"
	RedactionHash   = "hash"

	ResidencyAny        = "any"
	ResidencyRestricted = "restricted"
)

var (
	validWriters = map[string]struct{}{
		WriterCreator:     {},
		WriterSystem:      {},
		WriterSpecOwner:   {},
		WriterStatusOwner: {},
	}
	validReaders = map[string]struct{}{
		ReaderCustomer:   {},
		ReaderOperator:   {},
		ReaderInternal:   {},
		ReaderAdapter:    {},
		ReaderPlugin:     {},
		ReaderGovernance: {},
		ReaderSystem:     {},
	}
	validMutabilities = map[string]struct{}{
		MutabilityImmutable:  {},
		MutabilityMutable:    {},
		MutabilityAppendOnly: {},
		MutabilitySystemOnly: {},
	}
	validRetentions = map[string]struct{}{
		RetentionNone:     {},
		RetentionStandard: {},
		RetentionExtended: {},
		RetentionArchival: {},
	}
	validRedactions = map[string]struct{}{
		RedactionNone:   {},
		RedactionRedact: {},
		RedactionOmit:   {},
		RedactionHash:   {},
	}
	validResidencies = map[string]struct{}{
		ResidencyAny:        {},
		ResidencyRestricted: {},
	}
	requiredFieldPolicyKeySet map[string]struct{}
)

func init() {
	requiredFieldPolicyKeySet = make(map[string]struct{}, len(requiredFieldPolicyKeys))
	for _, k := range requiredFieldPolicyKeys {
		requiredFieldPolicyKeySet[k] = struct{}{}
	}
}

// SchemaMeta is the validated machine-readable metadata extracted from a
// canonical schema's registered x-sovrunn-* extensions (D-08, F12-NAMING-006).
type SchemaMeta struct {
	Profile       apimeta.Profile
	Boundary      apimeta.Boundary
	AllowedScopes []apimeta.ScopeKind
	Stability     apimeta.Stability
	// FieldPolicies maps the JSON Pointer of the schema object that declared
	// x-sovrunn-field-policy (typically a property schema) to the validated
	// policy. FEATURE-0012 does not inherit field policy across nesting.
	FieldPolicies map[string]FieldPolicyMeta
}

// FieldPolicyMeta is the strictly validated property-level
// x-sovrunn-field-policy object (D-08, F12-SEC-001).
type FieldPolicyMeta struct {
	Classification    apimeta.DataClassification
	AuthorizedWriter  string
	AuthorizedReaders []string
	Mutability        string
	Retention         string
	Redaction         string
	Residency         string
	AuditRequired     bool
}

// ReadAnnotations parses and validates the five registered x-sovrunn-* schema
// extensions (D-08). It:
//   - requires schema-level profile, boundary, allowed-scopes, and stability;
//   - validates those values against apimeta controlled vocabularies;
//   - strictly validates any property-level x-sovrunn-field-policy object
//     (exactly the eight declared fields; controlled vocabularies; no unknown
//     policy fields);
//   - fails closed on unknown x-sovrunn-* extensions (never silently ignored).
//
// Keys under "properties" remain property identifiers, not extension keywords.
// apischema MUST NOT import apiproblem; findings are package-local SchemaIssue.
func ReadAnnotations(schema []byte) (SchemaMeta, []SchemaIssue) {
	var zero SchemaMeta
	if len(schema) == 0 {
		return zero, []SchemaIssue{{
			Path:    "/",
			Code:    CodeMalformedSchema,
			Message: "schema document is required",
		}}
	}

	var root any
	if err := json.Unmarshal(schema, &root); err != nil {
		return zero, []SchemaIssue{{
			Path:    "/",
			Code:    CodeMalformedSchema,
			Message: "schema document is not valid JSON",
		}}
	}

	obj, ok := root.(map[string]any)
	if !ok {
		return zero, []SchemaIssue{{
			Path:    "/",
			Code:    CodeMalformedSchema,
			Message: "schema document must be a JSON object",
		}}
	}

	meta := SchemaMeta{
		FieldPolicies: make(map[string]FieldPolicyMeta),
	}
	var issues []SchemaIssue

	parseSchemaLevelAnnotations(obj, &meta, &issues)
	walkAnnotationSchema(obj, "", &meta, &issues)

	sort.SliceStable(issues, func(i, j int) bool {
		if issues[i].Path != issues[j].Path {
			return issues[i].Path < issues[j].Path
		}
		return issues[i].Code < issues[j].Code
	})

	if len(issues) > 0 {
		// Return partial meta for diagnostics; callers must treat any issue
		// as failure (fail-closed).
		return meta, issues
	}
	return meta, nil
}

func parseSchemaLevelAnnotations(obj map[string]any, meta *SchemaMeta, issues *[]SchemaIssue) {
	for _, key := range requiredSchemaAnnotations {
		if _, ok := obj[key]; !ok {
			*issues = append(*issues, SchemaIssue{
				Path:    joinPointer("", key),
				Code:    CodeAnnotationMissing,
				Message: fmt.Sprintf("required schema annotation %q is missing", key),
			})
		}
	}

	if raw, ok := obj[ExtProfile]; ok {
		if s, ok := asNonEmptyString(raw); ok {
			p := apimeta.Profile(s)
			if !p.Valid() {
				*issues = append(*issues, SchemaIssue{
					Path:    joinPointer("", ExtProfile),
					Code:    CodeAnnotationInvalid,
					Message: fmt.Sprintf("invalid x-sovrunn-profile value %q", s),
				})
			} else {
				meta.Profile = p
			}
		} else {
			*issues = append(*issues, SchemaIssue{
				Path:    joinPointer("", ExtProfile),
				Code:    CodeAnnotationInvalid,
				Message: "x-sovrunn-profile must be a non-empty string",
			})
		}
	}

	if raw, ok := obj[ExtBoundary]; ok {
		if s, ok := asNonEmptyString(raw); ok {
			b := apimeta.Boundary(s)
			if !b.Valid() {
				*issues = append(*issues, SchemaIssue{
					Path:    joinPointer("", ExtBoundary),
					Code:    CodeAnnotationInvalid,
					Message: fmt.Sprintf("invalid x-sovrunn-boundary value %q", s),
				})
			} else {
				meta.Boundary = b
			}
		} else {
			*issues = append(*issues, SchemaIssue{
				Path:    joinPointer("", ExtBoundary),
				Code:    CodeAnnotationInvalid,
				Message: "x-sovrunn-boundary must be a non-empty string",
			})
		}
	}

	if raw, ok := obj[ExtStability]; ok {
		if s, ok := asNonEmptyString(raw); ok {
			st := apimeta.Stability(s)
			if !st.Valid() {
				*issues = append(*issues, SchemaIssue{
					Path:    joinPointer("", ExtStability),
					Code:    CodeAnnotationInvalid,
					Message: fmt.Sprintf("invalid x-sovrunn-stability value %q", s),
				})
			} else {
				meta.Stability = st
			}
		} else {
			*issues = append(*issues, SchemaIssue{
				Path:    joinPointer("", ExtStability),
				Code:    CodeAnnotationInvalid,
				Message: "x-sovrunn-stability must be a non-empty string",
			})
		}
	}

	if raw, ok := obj[ExtAllowedScopes]; ok {
		scopes, scopeIssues := parseAllowedScopes(raw, joinPointer("", ExtAllowedScopes))
		*issues = append(*issues, scopeIssues...)
		if len(scopeIssues) == 0 {
			meta.AllowedScopes = scopes
		}
	}
}

func parseAllowedScopes(raw any, path string) ([]apimeta.ScopeKind, []SchemaIssue) {
	arr, ok := raw.([]any)
	if !ok {
		return nil, []SchemaIssue{{
			Path:    path,
			Code:    CodeAnnotationInvalid,
			Message: "x-sovrunn-allowed-scopes must be a non-empty array of scope kinds",
		}}
	}
	if len(arr) == 0 {
		return nil, []SchemaIssue{{
			Path:    path,
			Code:    CodeAnnotationInvalid,
			Message: "x-sovrunn-allowed-scopes must be a non-empty array of scope kinds",
		}}
	}

	seen := make(map[apimeta.ScopeKind]struct{}, len(arr))
	out := make([]apimeta.ScopeKind, 0, len(arr))
	var issues []SchemaIssue
	for i, item := range arr {
		itemPath := joinPointer(path, strconv.Itoa(i))
		s, ok := asNonEmptyString(item)
		if !ok {
			issues = append(issues, SchemaIssue{
				Path:    itemPath,
				Code:    CodeAnnotationInvalid,
				Message: "allowed scope entry must be a non-empty string",
			})
			continue
		}
		kind := apimeta.ScopeKind(s)
		if !kind.Valid() {
			issues = append(issues, SchemaIssue{
				Path:    itemPath,
				Code:    CodeAnnotationInvalid,
				Message: fmt.Sprintf("invalid scope kind %q", s),
			})
			continue
		}
		if _, dup := seen[kind]; dup {
			issues = append(issues, SchemaIssue{
				Path:    itemPath,
				Code:    CodeAnnotationInvalid,
				Message: fmt.Sprintf("duplicate allowed scope %q", s),
			})
			continue
		}
		seen[kind] = struct{}{}
		out = append(out, kind)
	}
	return out, issues
}

// walkAnnotationSchema walks schema objects looking for unknown x-sovrunn-*
// extensions and property-level field-policy objects. Property map keys are
// identifiers, not keywords.
func walkAnnotationSchema(node any, path string, meta *SchemaMeta, issues *[]SchemaIssue) {
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
		val := obj[key]
		childPath := joinPointer(path, key)

		switch {
		case isUnknownSovrunnExtension(key):
			*issues = append(*issues, SchemaIssue{
				Path:    childPath,
				Code:    CodeUnknownExtension,
				Message: fmt.Sprintf("unknown x-sovrunn extension %q; registered extensions only", key),
			})
			continue

		case key == ExtFieldPolicy:
			policy, policyIssues := parseFieldPolicy(val, childPath)
			*issues = append(*issues, policyIssues...)
			if len(policyIssues) == 0 {
				meta.FieldPolicies[pathOrRoot(path)] = policy
			}
			continue

		case key == ExtProfile || key == ExtBoundary || key == ExtAllowedScopes || key == ExtStability:
			// Schema-level values are validated in parseSchemaLevelAnnotations
			// when present at the document root. Nested copies are rejected so
			// profile/boundary/scope/stability remain document-level metadata.
			if path != "" && path != "/" {
				*issues = append(*issues, SchemaIssue{
					Path:    childPath,
					Code:    CodeAnnotationInvalid,
					Message: fmt.Sprintf("%s is only valid at the schema document root", key),
				})
			}
			continue

		case key == "properties":
			walkAnnotationProperties(val, childPath, meta, issues)
			continue

		case key == "items":
			walkAnnotationItems(val, childPath, meta, issues)
			continue

		case key == "additionalProperties":
			if _, isBool := val.(bool); isBool {
				continue
			}
			walkAnnotationSchema(val, childPath, meta, issues)
			continue
		}
	}
}

func walkAnnotationProperties(val any, path string, meta *SchemaMeta, issues *[]SchemaIssue) {
	obj, ok := val.(map[string]any)
	if !ok {
		*issues = append(*issues, SchemaIssue{
			Path:    pathOrRoot(path),
			Code:    CodeMalformedSchema,
			Message: "properties must be an object",
		})
		return
	}
	keys := make([]string, 0, len(obj))
	for k := range obj {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, propName := range keys {
		// Property identifiers are never treated as extension keywords, even
		// when they collide with x-sovrunn-* names.
		walkAnnotationSchema(obj[propName], joinPointer(path, propName), meta, issues)
	}
}

func walkAnnotationItems(val any, path string, meta *SchemaMeta, issues *[]SchemaIssue) {
	switch v := val.(type) {
	case []any:
		for i, item := range v {
			walkAnnotationSchema(item, joinPointer(path, strconv.Itoa(i)), meta, issues)
		}
	default:
		walkAnnotationSchema(v, path, meta, issues)
	}
}

func isUnknownSovrunnExtension(key string) bool {
	if len(key) < len(extPrefixSovrunn) || key[:len(extPrefixSovrunn)] != extPrefixSovrunn {
		return false
	}
	return !IsRegisteredExtension(key)
}

func parseFieldPolicy(raw any, path string) (FieldPolicyMeta, []SchemaIssue) {
	var zero FieldPolicyMeta
	obj, ok := raw.(map[string]any)
	if !ok {
		return zero, []SchemaIssue{{
			Path:    path,
			Code:    CodeFieldPolicyInvalid,
			Message: "x-sovrunn-field-policy must be an object",
		}}
	}

	var issues []SchemaIssue
	for key := range obj {
		if _, ok := requiredFieldPolicyKeySet[key]; !ok {
			issues = append(issues, SchemaIssue{
				Path:    joinPointer(path, key),
				Code:    CodeFieldPolicyUnknownField,
				Message: fmt.Sprintf("unknown x-sovrunn-field-policy field %q", key),
			})
		}
	}
	for _, key := range requiredFieldPolicyKeys {
		if _, ok := obj[key]; !ok {
			issues = append(issues, SchemaIssue{
				Path:    joinPointer(path, key),
				Code:    CodeFieldPolicyInvalid,
				Message: fmt.Sprintf("x-sovrunn-field-policy missing required field %q", key),
			})
		}
	}
	if len(issues) > 0 {
		sort.SliceStable(issues, func(i, j int) bool {
			if issues[i].Path != issues[j].Path {
				return issues[i].Path < issues[j].Path
			}
			return issues[i].Code < issues[j].Code
		})
		return zero, issues
	}

	out := FieldPolicyMeta{}

	if s, ok := asNonEmptyString(obj["classification"]); ok {
		c := apimeta.DataClassification(s)
		if !c.Valid() {
			issues = append(issues, SchemaIssue{
				Path:    joinPointer(path, "classification"),
				Code:    CodeFieldPolicyInvalid,
				Message: fmt.Sprintf("invalid classification %q", s),
			})
		} else {
			out.Classification = c
		}
	} else {
		issues = append(issues, SchemaIssue{
			Path:    joinPointer(path, "classification"),
			Code:    CodeFieldPolicyInvalid,
			Message: "classification must be a non-empty string",
		})
	}

	if s, ok := asNonEmptyString(obj["authorizedWriter"]); ok {
		if _, ok := validWriters[s]; !ok {
			issues = append(issues, SchemaIssue{
				Path:    joinPointer(path, "authorizedWriter"),
				Code:    CodeFieldPolicyInvalid,
				Message: fmt.Sprintf("invalid authorizedWriter %q", s),
			})
		} else {
			out.AuthorizedWriter = s
		}
	} else {
		issues = append(issues, SchemaIssue{
			Path:    joinPointer(path, "authorizedWriter"),
			Code:    CodeFieldPolicyInvalid,
			Message: "authorizedWriter must be a non-empty string",
		})
	}

	readers, readerIssues := parseAuthorizedReaders(obj["authorizedReaders"], joinPointer(path, "authorizedReaders"))
	issues = append(issues, readerIssues...)
	out.AuthorizedReaders = readers

	if s, ok := asNonEmptyString(obj["mutability"]); ok {
		if _, ok := validMutabilities[s]; !ok {
			issues = append(issues, SchemaIssue{
				Path:    joinPointer(path, "mutability"),
				Code:    CodeFieldPolicyInvalid,
				Message: fmt.Sprintf("invalid mutability %q", s),
			})
		} else {
			out.Mutability = s
		}
	} else {
		issues = append(issues, SchemaIssue{
			Path:    joinPointer(path, "mutability"),
			Code:    CodeFieldPolicyInvalid,
			Message: "mutability must be a non-empty string",
		})
	}

	if s, ok := asNonEmptyString(obj["retention"]); ok {
		if _, ok := validRetentions[s]; !ok {
			issues = append(issues, SchemaIssue{
				Path:    joinPointer(path, "retention"),
				Code:    CodeFieldPolicyInvalid,
				Message: fmt.Sprintf("invalid retention %q", s),
			})
		} else {
			out.Retention = s
		}
	} else {
		issues = append(issues, SchemaIssue{
			Path:    joinPointer(path, "retention"),
			Code:    CodeFieldPolicyInvalid,
			Message: "retention must be a non-empty string",
		})
	}

	if s, ok := asNonEmptyString(obj["redaction"]); ok {
		if _, ok := validRedactions[s]; !ok {
			issues = append(issues, SchemaIssue{
				Path:    joinPointer(path, "redaction"),
				Code:    CodeFieldPolicyInvalid,
				Message: fmt.Sprintf("invalid redaction %q", s),
			})
		} else {
			out.Redaction = s
		}
	} else {
		issues = append(issues, SchemaIssue{
			Path:    joinPointer(path, "redaction"),
			Code:    CodeFieldPolicyInvalid,
			Message: "redaction must be a non-empty string",
		})
	}

	if s, ok := asNonEmptyString(obj["residency"]); ok {
		if _, ok := validResidencies[s]; !ok {
			issues = append(issues, SchemaIssue{
				Path:    joinPointer(path, "residency"),
				Code:    CodeFieldPolicyInvalid,
				Message: fmt.Sprintf("invalid residency %q", s),
			})
		} else {
			out.Residency = s
		}
	} else {
		issues = append(issues, SchemaIssue{
			Path:    joinPointer(path, "residency"),
			Code:    CodeFieldPolicyInvalid,
			Message: "residency must be a non-empty string",
		})
	}

	if b, ok := obj["auditRequired"].(bool); ok {
		out.AuditRequired = b
	} else {
		issues = append(issues, SchemaIssue{
			Path:    joinPointer(path, "auditRequired"),
			Code:    CodeFieldPolicyInvalid,
			Message: "auditRequired must be a boolean",
		})
	}

	if len(issues) > 0 {
		sort.SliceStable(issues, func(i, j int) bool {
			if issues[i].Path != issues[j].Path {
				return issues[i].Path < issues[j].Path
			}
			return issues[i].Code < issues[j].Code
		})
		return zero, issues
	}
	return out, nil
}

func parseAuthorizedReaders(raw any, path string) ([]string, []SchemaIssue) {
	arr, ok := raw.([]any)
	if !ok {
		return nil, []SchemaIssue{{
			Path:    path,
			Code:    CodeFieldPolicyInvalid,
			Message: "authorizedReaders must be a non-empty array of strings",
		}}
	}
	if len(arr) == 0 {
		return nil, []SchemaIssue{{
			Path:    path,
			Code:    CodeFieldPolicyInvalid,
			Message: "authorizedReaders must be a non-empty array of strings",
		}}
	}

	seen := make(map[string]struct{}, len(arr))
	out := make([]string, 0, len(arr))
	var issues []SchemaIssue
	for i, item := range arr {
		itemPath := joinPointer(path, strconv.Itoa(i))
		s, ok := asNonEmptyString(item)
		if !ok {
			issues = append(issues, SchemaIssue{
				Path:    itemPath,
				Code:    CodeFieldPolicyInvalid,
				Message: "authorizedReaders entry must be a non-empty string",
			})
			continue
		}
		if _, ok := validReaders[s]; !ok {
			issues = append(issues, SchemaIssue{
				Path:    itemPath,
				Code:    CodeFieldPolicyInvalid,
				Message: fmt.Sprintf("invalid authorizedReaders value %q", s),
			})
			continue
		}
		if _, dup := seen[s]; dup {
			issues = append(issues, SchemaIssue{
				Path:    itemPath,
				Code:    CodeFieldPolicyInvalid,
				Message: fmt.Sprintf("duplicate authorizedReaders value %q", s),
			})
			continue
		}
		seen[s] = struct{}{}
		out = append(out, s)
	}
	return out, issues
}

func asNonEmptyString(v any) (string, bool) {
	s, ok := v.(string)
	if !ok || s == "" {
		return "", false
	}
	return s, true
}
