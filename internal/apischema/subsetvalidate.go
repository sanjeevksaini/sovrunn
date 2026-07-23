package apischema

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// SchemaIssue is a package-local diagnostic for schema-support or (later)
// structural instance findings. apischema MUST NOT import apiproblem; the
// translation from SchemaIssue to apiproblem.Violation is owned by
// apiconform (D-01a, D-02).
type SchemaIssue struct {
	Path    string // RFC 6901 JSON Pointer
	Code    string // stable machine-readable code
	Message string // human-readable; must not carry secrets or provider detail
}

// Stable SchemaIssue codes for bounded-subset support scanning (D-01a).
const (
	CodeUnsupportedKeyword = "SCHEMA_UNSUPPORTED_KEYWORD"
	CodeMalformedSchema    = "SCHEMA_MALFORMED"
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
