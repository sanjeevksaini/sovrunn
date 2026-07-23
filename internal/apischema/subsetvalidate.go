package apischema

import (
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"
)

// SchemaIssue is a package-local diagnostic for schema-support or structural
// instance findings. apischema MUST NOT import apiproblem; the translation
// from SchemaIssue to apiproblem.Violation is owned by apiconform (D-01a, D-02).
type SchemaIssue struct {
	Path    string // RFC 6901 JSON Pointer
	Code    string // stable machine-readable code
	Message string // human-readable; must not carry secrets or provider detail
}

// Stable SchemaIssue codes for bounded-subset support scanning and instance
// validation (D-01a, F12-VALIDATION-006).
const (
	CodeUnsupportedKeyword = "SCHEMA_UNSUPPORTED_KEYWORD"
	CodeMalformedSchema    = "SCHEMA_MALFORMED"
	CodeUnresolvedRef      = "SCHEMA_UNRESOLVED_REF"
	CodeSchemaFalse        = "SCHEMA_FALSE"
	CodeTypeMismatch       = "TYPE_MISMATCH"
	CodeRequiredField      = "REQUIRED_FIELD"
	CodeEnumMismatch       = "ENUM_MISMATCH"
	CodeOutOfRange         = "OUT_OF_RANGE"
	CodePatternMismatch    = "PATTERN_MISMATCH"
	CodeUnknownField       = "UNKNOWN_FIELD"
)

// CoreSupportedKeywords is the explicit JSON Schema 2020-12 subset FEATURE-0012
// supports for structural validation. $defs is intentionally absent (prohibited);
// shared definitions use approved relative $ref under api/schemas/_common only.
var CoreSupportedKeywords = []string{
	"$schema",
	"$id",
	"$ref",
	"title",
	"description",
	"type",
	"properties",
	"required",
	"enum",
	"items",
	"additionalProperties",
	"minLength",
	"maxLength",
	"minimum",
	"maximum",
	"pattern",
	"default",
	"examples",
}

// RegisteredExtensionKeywords are the five x-sovrunn-* extension keywords
// registered by D-08. Unknown x-sovrunn-* names are unsupported and fail closed.
// Full vocabulary validation of extension values is owned by task 9.1
// (ReadAnnotations); this package only treats registered names as supported
// keywords and does not treat extension-object fields as schema keywords.
var RegisteredExtensionKeywords = []string{
	"x-sovrunn-profile",
	"x-sovrunn-boundary",
	"x-sovrunn-allowed-scopes",
	"x-sovrunn-stability",
	"x-sovrunn-field-policy",
}

// SupportedKeywords is the closed set of keywords ValidateSchemaSupport accepts.
// It is the union of CoreSupportedKeywords and RegisteredExtensionKeywords.
// Callers must not mutate the map; treat it as an immutable declared vocabulary.
var SupportedKeywords map[string]struct{}

var registeredExtensionSet map[string]struct{}

func init() {
	SupportedKeywords = make(map[string]struct{}, len(CoreSupportedKeywords)+len(RegisteredExtensionKeywords))
	for _, k := range CoreSupportedKeywords {
		SupportedKeywords[k] = struct{}{}
	}
	registeredExtensionSet = make(map[string]struct{}, len(RegisteredExtensionKeywords))
	for _, k := range RegisteredExtensionKeywords {
		SupportedKeywords[k] = struct{}{}
		registeredExtensionSet[k] = struct{}{}
	}
}

// IsSupportedKeyword reports whether key is in the declared supported subset.
func IsSupportedKeyword(key string) bool {
	_, ok := SupportedKeywords[key]
	return ok
}

// IsRegisteredExtension reports whether key is one of the five registered
// x-sovrunn-* extension keywords (D-08).
func IsRegisteredExtension(key string) bool {
	_, ok := registeredExtensionSet[key]
	return ok
}

// ValidateSchemaSupport scans a canonical JSON Schema document and returns
// issues for ANY actual schema-position keyword outside SupportedKeywords.
// It is FAIL-CLOSED: unsupported keywords are rejected, never ignored, so no
// constraint is silently unenforced (D-01a, F12-NAMING-005).
//
// The walker is context-aware:
//   - keys under "properties" are property identifiers, not keywords;
//   - fields inside registered extension objects (e.g. x-sovrunn-field-policy)
//     are extension-schema fields, not core JSON Schema keywords;
//   - only actual schema-position keywords outside the supported set produce
//     SCHEMA_UNSUPPORTED_KEYWORD issues.
//
// Nested schema values under properties, items, and additionalProperties are
// walked recursively. Document metadata keywords ($schema, $id, title,
// description) are accepted when present.
func ValidateSchemaSupport(schema []byte) []SchemaIssue {
	if len(schema) == 0 {
		return []SchemaIssue{{
			Path:    "/",
			Code:    CodeMalformedSchema,
			Message: "schema document is required",
		}}
	}

	var root any
	if err := json.Unmarshal(schema, &root); err != nil {
		return []SchemaIssue{{
			Path:    "/",
			Code:    CodeMalformedSchema,
			Message: "schema document is not valid JSON",
		}}
	}

	var issues []SchemaIssue
	walkSchema(root, "", &issues)
	sort.SliceStable(issues, func(i, j int) bool {
		if issues[i].Path != issues[j].Path {
			return issues[i].Path < issues[j].Path
		}
		return issues[i].Code < issues[j].Code
	})
	return issues
}

// walkSchema walks a JSON Schema node (object or boolean). Non-schema types at
// a schema position are malformed for FEATURE-0012 canonical documents.
func walkSchema(node any, path string, issues *[]SchemaIssue) {
	switch node.(type) {
	case bool:
		return
	case nil:
		*issues = append(*issues, SchemaIssue{
			Path:    pathOrRoot(path),
			Code:    CodeMalformedSchema,
			Message: "schema value must be an object or boolean",
		})
		return
	}

	obj, ok := node.(map[string]any)
	if !ok {
		*issues = append(*issues, SchemaIssue{
			Path:    pathOrRoot(path),
			Code:    CodeMalformedSchema,
			Message: "schema value must be an object or boolean",
		})
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

		if !IsSupportedKeyword(key) {
			*issues = append(*issues, SchemaIssue{
				Path:    childPath,
				Code:    CodeUnsupportedKeyword,
				Message: fmt.Sprintf("unsupported schema keyword %q", key),
			})
			// Do not walk into unsupported keyword values as schema contexts:
			// the keyword itself is the fail-closed finding.
			continue
		}

		switch {
		case key == "properties":
			walkPropertiesMap(val, childPath, issues)
		case key == "items":
			walkItems(val, childPath, issues)
		case key == "additionalProperties":
			if _, isBool := val.(bool); isBool {
				continue
			}
			walkSchema(val, childPath, issues)
		case IsRegisteredExtension(key):
			// Extension-object fields are not schema keywords. Do not recurse
			// into extension values as schema documents (task 9.1 owns value
			// shape/vocabulary checks).
			continue
		default:
			// Leaf supported keywords (type, required, enum, $ref, metadata,
			// bounds, default, examples, …): values are not schema objects.
		}
	}
}

// walkPropertiesMap treats map keys as property identifiers (never keywords)
// and walks each value as a nested schema.
func walkPropertiesMap(val any, path string, issues *[]SchemaIssue) {
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
		// propName is an identifier, even when it collides with a keyword name
		// such as "oneOf" or "$defs".
		walkSchema(obj[propName], joinPointer(path, propName), issues)
	}
}

func walkItems(val any, path string, issues *[]SchemaIssue) {
	switch v := val.(type) {
	case []any:
		for i, item := range v {
			walkSchema(item, joinPointer(path, strconv.Itoa(i)), issues)
		}
	default:
		walkSchema(v, path, issues)
	}
}

func joinPointer(base, token string) string {
	escaped := escapeJSONPointer(token)
	if base == "" || base == "/" {
		return "/" + escaped
	}
	return base + "/" + escaped
}

func escapeJSONPointer(s string) string {
	s = strings.ReplaceAll(s, "~", "~0")
	s = strings.ReplaceAll(s, "/", "~1")
	return s
}

func pathOrRoot(path string) string {
	if path == "" {
		return "/"
	}
	return path
}

// ValidateInstance structurally validates a decoded instance against a
// canonical JSON Schema document using only the FEATURE-0012 supported
// subset (D-01a, F12-VALIDATION-001(4), F12-VALIDATION-006).
//
// Callers MUST first pass ValidateSchemaSupport. This function also re-checks
// support and returns those issues if the schema is out of subset (fail-closed).
//
// Local $ref resolution is owned by apiconform (tasks 7a/8.2). An unresolved
// $ref in the schema document fails closed with SCHEMA_UNRESOLVED_REF.
//
// Typed Go values are normalized through a JSON round-trip so validation
// operates on JSON-compatible primitives (map/slice/string/number/bool/null).
func ValidateInstance(schema []byte, instance any) []SchemaIssue {
	if support := ValidateSchemaSupport(schema); len(support) > 0 {
		return support
	}

	var root any
	if err := json.Unmarshal(schema, &root); err != nil {
		// ValidateSchemaSupport already rejected malformed JSON; keep a
		// defensive path for empty/race cases.
		return []SchemaIssue{{
			Path:    "/",
			Code:    CodeMalformedSchema,
			Message: "schema document is not valid JSON",
		}}
	}

	normalized, normIssues := normalizeInstance(instance)
	if len(normIssues) > 0 {
		return normIssues
	}

	var issues []SchemaIssue
	validateNode(root, normalized, "", &issues)
	sort.SliceStable(issues, func(i, j int) bool {
		if issues[i].Path != issues[j].Path {
			return issues[i].Path < issues[j].Path
		}
		return issues[i].Code < issues[j].Code
	})
	return issues
}

func normalizeInstance(instance any) (any, []SchemaIssue) {
	switch instance.(type) {
	case nil, bool, string, float64, json.Number, map[string]any, []any:
		return instance, nil
	default:
		raw, err := json.Marshal(instance)
		if err != nil {
			return nil, []SchemaIssue{{
				Path:    "/",
				Code:    CodeTypeMismatch,
				Message: "instance could not be normalized for structural validation",
			}}
		}
		var out any
		if err := json.Unmarshal(raw, &out); err != nil {
			return nil, []SchemaIssue{{
				Path:    "/",
				Code:    CodeTypeMismatch,
				Message: "instance could not be normalized for structural validation",
			}}
		}
		return out, nil
	}
}

func validateNode(schema any, instance any, path string, issues *[]SchemaIssue) {
	switch s := schema.(type) {
	case bool:
		if !s {
			*issues = append(*issues, SchemaIssue{
				Path:    pathOrRoot(path),
				Code:    CodeSchemaFalse,
				Message: "schema rejects all values",
			})
		}
		return
	case nil:
		*issues = append(*issues, SchemaIssue{
			Path:    pathOrRoot(path),
			Code:    CodeMalformedSchema,
			Message: "schema value must be an object or boolean",
		})
		return
	}

	obj, ok := schema.(map[string]any)
	if !ok {
		*issues = append(*issues, SchemaIssue{
			Path:    pathOrRoot(path),
			Code:    CodeMalformedSchema,
			Message: "schema value must be an object or boolean",
		})
		return
	}

	if _, hasRef := obj["$ref"]; hasRef {
		*issues = append(*issues, SchemaIssue{
			Path:    pathOrRoot(path),
			Code:    CodeUnresolvedRef,
			Message: "schema $ref must be resolved before instance validation",
		})
		return
	}

	if typeVal, hasType := obj["type"]; hasType {
		if !instanceMatchesType(typeVal, instance) {
			*issues = append(*issues, SchemaIssue{
				Path:    pathOrRoot(path),
				Code:    CodeTypeMismatch,
				Message: "value does not match declared type",
			})
			return
		}
	}

	if enumVal, hasEnum := obj["enum"]; hasEnum {
		enumList, ok := enumVal.([]any)
		if !ok {
			*issues = append(*issues, SchemaIssue{
				Path:    joinPointer(path, "enum"),
				Code:    CodeMalformedSchema,
				Message: "enum must be an array",
			})
			return
		}
		if !enumContains(enumList, instance) {
			*issues = append(*issues, SchemaIssue{
				Path:    pathOrRoot(path),
				Code:    CodeEnumMismatch,
				Message: "value is not a member of enum",
			})
			return
		}
	}

	switch v := instance.(type) {
	case string:
		validateStringConstraints(obj, v, path, issues)
	case float64:
		validateNumberConstraints(obj, v, path, issues)
	case json.Number:
		f, err := v.Float64()
		if err != nil {
			*issues = append(*issues, SchemaIssue{
				Path:    pathOrRoot(path),
				Code:    CodeTypeMismatch,
				Message: "value does not match declared type",
			})
			return
		}
		validateNumberConstraints(obj, f, path, issues)
	case map[string]any:
		validateObjectConstraints(obj, v, path, issues)
	case []any:
		validateArrayConstraints(obj, v, path, issues)
	}
}

func instanceMatchesType(typeVal any, instance any) bool {
	switch t := typeVal.(type) {
	case string:
		return matchesSingleType(t, instance)
	case []any:
		if len(t) == 0 {
			return false
		}
		for _, item := range t {
			name, ok := item.(string)
			if !ok {
				return false
			}
			if matchesSingleType(name, instance) {
				return true
			}
		}
		return false
	default:
		return false
	}
}

func matchesSingleType(typeName string, instance any) bool {
	switch typeName {
	case "null":
		return instance == nil
	case "boolean":
		_, ok := instance.(bool)
		return ok
	case "string":
		_, ok := instance.(string)
		return ok
	case "object":
		_, ok := instance.(map[string]any)
		return ok
	case "array":
		_, ok := instance.([]any)
		return ok
	case "number":
		return isJSONNumber(instance)
	case "integer":
		return isJSONInteger(instance)
	default:
		return false
	}
}

func isJSONNumber(instance any) bool {
	switch v := instance.(type) {
	case float64:
		return !math.IsNaN(v) && !math.IsInf(v, 0)
	case json.Number:
		_, err := v.Float64()
		return err == nil
	default:
		return false
	}
}

func isJSONInteger(instance any) bool {
	switch v := instance.(type) {
	case float64:
		return !math.IsNaN(v) && !math.IsInf(v, 0) && v == math.Trunc(v)
	case json.Number:
		f, err := v.Float64()
		if err != nil {
			return false
		}
		return !math.IsNaN(f) && !math.IsInf(f, 0) && f == math.Trunc(f)
	default:
		return false
	}
}

func enumContains(enum []any, instance any) bool {
	for _, candidate := range enum {
		if jsonValuesEqual(candidate, instance) {
			return true
		}
	}
	return false
}

func jsonValuesEqual(a, b any) bool {
	if a == nil || b == nil {
		return a == nil && b == nil
	}
	af, aok := asFloat64(a)
	bf, bok := asFloat64(b)
	if aok && bok {
		return af == bf
	}
	return reflect.DeepEqual(a, b)
}

func asFloat64(v any) (float64, bool) {
	switch n := v.(type) {
	case float64:
		return n, true
	case json.Number:
		f, err := n.Float64()
		return f, err == nil
	default:
		return 0, false
	}
}

func validateStringConstraints(schema map[string]any, value string, path string, issues *[]SchemaIssue) {
	runeLen := utf8.RuneCountInString(value)

	if raw, ok := schema["minLength"]; ok {
		min, ok := asNonNegativeInt(raw)
		if !ok {
			*issues = append(*issues, SchemaIssue{
				Path:    joinPointer(path, "minLength"),
				Code:    CodeMalformedSchema,
				Message: "minLength must be a non-negative integer",
			})
		} else if runeLen < min {
			*issues = append(*issues, SchemaIssue{
				Path:    pathOrRoot(path),
				Code:    CodeOutOfRange,
				Message: fmt.Sprintf("string length must be >= %d", min),
			})
		}
	}

	if raw, ok := schema["maxLength"]; ok {
		max, ok := asNonNegativeInt(raw)
		if !ok {
			*issues = append(*issues, SchemaIssue{
				Path:    joinPointer(path, "maxLength"),
				Code:    CodeMalformedSchema,
				Message: "maxLength must be a non-negative integer",
			})
		} else if runeLen > max {
			*issues = append(*issues, SchemaIssue{
				Path:    pathOrRoot(path),
				Code:    CodeOutOfRange,
				Message: fmt.Sprintf("string length must be <= %d", max),
			})
		}
	}

	if raw, ok := schema["pattern"]; ok {
		pat, ok := raw.(string)
		if !ok {
			*issues = append(*issues, SchemaIssue{
				Path:    joinPointer(path, "pattern"),
				Code:    CodeMalformedSchema,
				Message: "pattern must be a string",
			})
			return
		}
		re, err := regexp.Compile(pat)
		if err != nil {
			*issues = append(*issues, SchemaIssue{
				Path:    joinPointer(path, "pattern"),
				Code:    CodeMalformedSchema,
				Message: "pattern is not a valid regular expression",
			})
			return
		}
		if !re.MatchString(value) {
			*issues = append(*issues, SchemaIssue{
				Path:    pathOrRoot(path),
				Code:    CodePatternMismatch,
				Message: "value does not match pattern",
			})
		}
	}
}

func validateNumberConstraints(schema map[string]any, value float64, path string, issues *[]SchemaIssue) {
	if raw, ok := schema["minimum"]; ok {
		min, ok := asFloat64(raw)
		if !ok {
			*issues = append(*issues, SchemaIssue{
				Path:    joinPointer(path, "minimum"),
				Code:    CodeMalformedSchema,
				Message: "minimum must be a number",
			})
		} else if value < min {
			*issues = append(*issues, SchemaIssue{
				Path:    pathOrRoot(path),
				Code:    CodeOutOfRange,
				Message: fmt.Sprintf("value must be >= %v", min),
			})
		}
	}

	if raw, ok := schema["maximum"]; ok {
		max, ok := asFloat64(raw)
		if !ok {
			*issues = append(*issues, SchemaIssue{
				Path:    joinPointer(path, "maximum"),
				Code:    CodeMalformedSchema,
				Message: "maximum must be a number",
			})
		} else if value > max {
			*issues = append(*issues, SchemaIssue{
				Path:    pathOrRoot(path),
				Code:    CodeOutOfRange,
				Message: fmt.Sprintf("value must be <= %v", max),
			})
		}
	}
}

func validateObjectConstraints(schema map[string]any, value map[string]any, path string, issues *[]SchemaIssue) {
	if raw, ok := schema["required"]; ok {
		req, ok := raw.([]any)
		if !ok {
			*issues = append(*issues, SchemaIssue{
				Path:    joinPointer(path, "required"),
				Code:    CodeMalformedSchema,
				Message: "required must be an array",
			})
		} else {
			for _, item := range req {
				name, ok := item.(string)
				if !ok {
					*issues = append(*issues, SchemaIssue{
						Path:    joinPointer(path, "required"),
						Code:    CodeMalformedSchema,
						Message: "required entries must be strings",
					})
					continue
				}
				if _, present := value[name]; !present {
					*issues = append(*issues, SchemaIssue{
						Path:    joinPointer(path, name),
						Code:    CodeRequiredField,
						Message: fmt.Sprintf("required property %q is missing", name),
					})
				}
			}
		}
	}

	var props map[string]any
	if raw, ok := schema["properties"]; ok {
		props, ok = raw.(map[string]any)
		if !ok {
			*issues = append(*issues, SchemaIssue{
				Path:    joinPointer(path, "properties"),
				Code:    CodeMalformedSchema,
				Message: "properties must be an object",
			})
			props = nil
		}
	}

	keys := make([]string, 0, len(value))
	for k := range value {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		childPath := joinPointer(path, key)
		if props != nil {
			if propSchema, ok := props[key]; ok {
				validateNode(propSchema, value[key], childPath, issues)
				continue
			}
		}

		if raw, ok := schema["additionalProperties"]; ok {
			switch ap := raw.(type) {
			case bool:
				if !ap {
					*issues = append(*issues, SchemaIssue{
						Path:    childPath,
						Code:    CodeUnknownField,
						Message: "additional property is not permitted",
					})
				}
			default:
				validateNode(ap, value[key], childPath, issues)
			}
		}
	}
}

func validateArrayConstraints(schema map[string]any, value []any, path string, issues *[]SchemaIssue) {
	raw, ok := schema["items"]
	if !ok {
		return
	}

	switch items := raw.(type) {
	case []any:
		for i, item := range value {
			if i >= len(items) {
				break
			}
			validateNode(items[i], item, joinPointer(path, strconv.Itoa(i)), issues)
		}
	default:
		for i, item := range value {
			validateNode(items, item, joinPointer(path, strconv.Itoa(i)), issues)
		}
	}
}

func asNonNegativeInt(v any) (int, bool) {
	f, ok := asFloat64(v)
	if !ok {
		return 0, false
	}
	if math.IsNaN(f) || math.IsInf(f, 0) || f < 0 || f != math.Trunc(f) {
		return 0, false
	}
	return int(f), true
}
