package apischema

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strings"
)

// Stable SchemaIssue codes for derivative Go-type / schema consistency
// checks (D-01b, F12-VERIFY-001(13)).
const (
	CodeGoTypeMismatch                 = "GO_TYPE_MISMATCH"
	CodeGoFieldMissing                 = "GO_FIELD_MISSING"
	CodeGoFieldExtra                   = "GO_FIELD_EXTRA"
	CodeGoRequiredMismatch             = "GO_REQUIRED_MISMATCH"
	CodeGoAdditionalPropertiesMismatch = "GO_ADDITIONAL_PROPERTIES_MISMATCH"
	CodeGoEnumTypeMismatch             = "GO_ENUM_TYPE_MISMATCH"
)

// TypeBinding maps one canonical schema path to its derivative Go type.
// Bindings are generic: SchemaPath identifies the schema document and GoType
// is any reflect.Type. Concrete contract types are registered by callers
// (apiconform) so apischema never imports them (D-01b).
//
// Fixture round-tripping is supporting evidence only; VerifyGoTypeAgainstSchema
// is the authoritative consistency check.
type TypeBinding struct {
	SchemaPath string       // e.g. api/schemas/project.json
	GoType     reflect.Type // e.g. reflect.TypeOf(Project{})
}

// VerifyGoTypeAgainstSchema verifies, via reflection, that a derivative Go
// type matches its canonical schema for the SUPPORTED subset (D-01b).
//
// It checks at least:
//   - property names and JSON tags
//   - required vs optional (required → no omitempty / non-pointer;
//     optional → omitempty or pointer)
//   - primitive types (string / number / integer / boolean)
//   - arrays (slices) and maps (objects with additionalProperties)
//   - embedded (promoted) fields
//   - enum-backed named types (named string/int types whose schema declares enum)
//   - additionalProperties behavior (open maps vs closed structs)
//
// Property schemas that carry $ref are treated as opaque nested structs:
// the referenced schema is verified through its own TypeBinding, not by
// resolving the reference here. Callers MUST first pass ValidateSchemaSupport
// so only supported-subset schemas are checked.
//
// Returns package-local SchemaIssue values (NOT apiproblem.Violation).
func VerifyGoTypeAgainstSchema(schema []byte, goType reflect.Type) []SchemaIssue {
	if len(schema) == 0 {
		return []SchemaIssue{{
			Path:    "/",
			Code:    CodeMalformedSchema,
			Message: "schema document is required",
		}}
	}
	if goType == nil {
		return []SchemaIssue{{
			Path:    "/",
			Code:    CodeGoTypeMismatch,
			Message: "Go type is required",
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
	verifySchemaNode(root, unwrapType(goType), "", 0, &issues)
	sort.SliceStable(issues, func(i, j int) bool {
		if issues[i].Path != issues[j].Path {
			return issues[i].Path < issues[j].Path
		}
		return issues[i].Code < issues[j].Code
	})
	return issues
}

const maxTypeBindingDepth = 32

type goJSONField struct {
	jsonName  string
	typ       reflect.Type
	omitempty bool
	path      string // JSON Pointer to this field
}

func verifySchemaNode(schema any, goType reflect.Type, path string, depth int, issues *[]SchemaIssue) {
	if depth > maxTypeBindingDepth {
		*issues = append(*issues, SchemaIssue{
			Path:    jsonPointer(path),
			Code:    CodeGoTypeMismatch,
			Message: "schema/Go type nesting exceeds maximum depth",
		})
		return
	}

	switch s := schema.(type) {
	case bool:
		if !s {
			*issues = append(*issues, SchemaIssue{
				Path:    jsonPointer(path),
				Code:    CodeSchemaFalse,
				Message: "schema is false; no Go type can satisfy it",
			})
		}
		return
	case map[string]any:
		verifyObjectSchema(s, goType, path, depth, issues)
	default:
		*issues = append(*issues, SchemaIssue{
			Path:    jsonPointer(path),
			Code:    CodeMalformedSchema,
			Message: "schema node must be an object or boolean",
		})
	}
}

func verifyObjectSchema(schema map[string]any, goType reflect.Type, path string, depth int, issues *[]SchemaIssue) {
	goType = unwrapType(goType)

	// $ref-only (plus metadata/extensions) schemas are opaque nested objects.
	if ref, ok := schema["$ref"].(string); ok && strings.TrimSpace(ref) != "" {
		if !isStructKind(goType) {
			*issues = append(*issues, SchemaIssue{
				Path:    jsonPointer(path),
				Code:    CodeGoTypeMismatch,
				Message: fmt.Sprintf("$ref %q requires a struct Go type, got %s", ref, typeName(goType)),
			})
		}
		return
	}

	schemaType, hasType := schemaTypeString(schema)
	_, hasEnum := schema["enum"]
	_, hasProperties := schema["properties"]
	apRaw, hasAP := schema["additionalProperties"]

	switch {
	case hasType && schemaType == "object", hasProperties, hasAP && !hasType:
		verifyObjectGoType(schema, goType, path, depth, issues, hasProperties, hasAP, apRaw)
	case hasType && schemaType == "array":
		verifyArrayGoType(schema, goType, path, depth, issues)
	case hasType || hasEnum:
		verifyPrimitiveGoType(schema, goType, path, schemaType, hasEnum, issues)
	default:
		// Empty/unconstrained schema accepts any Go type.
	}
}

func verifyObjectGoType(schema map[string]any, goType reflect.Type, path string, depth int, issues *[]SchemaIssue, hasProperties, hasAP bool, apRaw any) {
	props := map[string]any{}
	if raw, ok := schema["properties"].(map[string]any); ok {
		props = raw
		hasProperties = true
	}

	apFalse := hasAP && apRaw == false
	apOpen := hasAP && apRaw == true
	apSchema, apIsSchema := apRaw.(map[string]any)
	apBoolSchema, apIsBool := apRaw.(bool)
	if apIsBool && apBoolSchema {
		apOpen = true
	}

	// Map form: object with additionalProperties schema/true and no declared properties.
	if (!hasProperties || len(props) == 0) && (apOpen || apIsSchema) {
		if goType.Kind() != reflect.Map {
			*issues = append(*issues, SchemaIssue{
				Path:    jsonPointer(path),
				Code:    CodeGoAdditionalPropertiesMismatch,
				Message: fmt.Sprintf("additionalProperties object requires map Go type, got %s", typeName(goType)),
			})
			return
		}
		if goType.Key().Kind() != reflect.String {
			*issues = append(*issues, SchemaIssue{
				Path:    jsonPointer(path),
				Code:    CodeGoAdditionalPropertiesMismatch,
				Message: fmt.Sprintf("map key type must be string, got %s", typeName(goType.Key())),
			})
			return
		}
		if apIsSchema {
			verifySchemaNode(apSchema, goType.Elem(), path, depth+1, issues)
		}
		return
	}

	// Closed or property-bearing object → struct (including anonymous embeds).
	if goType.Kind() != reflect.Struct {
		*issues = append(*issues, SchemaIssue{
			Path:    jsonPointer(path),
			Code:    CodeGoTypeMismatch,
			Message: fmt.Sprintf("object schema requires struct Go type, got %s", typeName(goType)),
		})
		return
	}

	if apOpen {
		*issues = append(*issues, SchemaIssue{
			Path:    jsonPointer(path),
			Code:    CodeGoAdditionalPropertiesMismatch,
			Message: "additionalProperties true with declared properties requires an open map representation, not a closed struct",
		})
		return
	}
	if apIsSchema && len(props) > 0 {
		*issues = append(*issues, SchemaIssue{
			Path:    jsonPointer(path),
			Code:    CodeGoAdditionalPropertiesMismatch,
			Message: "struct Go types cannot represent property maps with additionalProperties schemas; use a map type binding",
		})
		return
	}
	if hasAP && !apFalse && !apOpen && !apIsSchema {
		*issues = append(*issues, SchemaIssue{
			Path:    jsonPointer(path),
			Code:    CodeGoAdditionalPropertiesMismatch,
			Message: "additionalProperties value is not a boolean or schema object",
		})
		return
	}

	fields := collectJSONFields(goType, path)
	required := requiredSet(schema)

	for name, propSchema := range props {
		field, ok := fields[name]
		if !ok {
			*issues = append(*issues, SchemaIssue{
				Path:    jsonPointer(joinPath(path, name)),
				Code:    CodeGoFieldMissing,
				Message: fmt.Sprintf("schema property %q has no matching Go JSON field", name),
			})
			continue
		}
		delete(fields, name)

		_, isRequired := required[name]
		verifyRequiredOptional(field, isRequired, joinPath(path, name), issues)
		verifySchemaNode(propSchema, field.typ, joinPath(path, name), depth+1, issues)
	}

	for name, field := range fields {
		*issues = append(*issues, SchemaIssue{
			Path:    jsonPointer(field.path),
			Code:    CodeGoFieldExtra,
			Message: fmt.Sprintf("Go JSON field %q is not declared in schema properties (additionalProperties false)", name),
		})
	}
}

func verifyArrayGoType(schema map[string]any, goType reflect.Type, path string, depth int, issues *[]SchemaIssue) {
	if goType.Kind() != reflect.Slice && goType.Kind() != reflect.Array {
		*issues = append(*issues, SchemaIssue{
			Path:    jsonPointer(path),
			Code:    CodeGoTypeMismatch,
			Message: fmt.Sprintf("array schema requires slice/array Go type, got %s", typeName(goType)),
		})
		return
	}
	if items, ok := schema["items"]; ok {
		verifySchemaNode(items, goType.Elem(), path, depth+1, issues)
	}
}

func verifyPrimitiveGoType(schema map[string]any, goType reflect.Type, path, schemaType string, hasEnum bool, issues *[]SchemaIssue) {
	goType = unwrapType(goType)
	underlying := goType
	for underlying.Kind() == reflect.Ptr {
		underlying = underlying.Elem()
	}

	if hasEnum {
		// Enum-backed fields may be named string/int types or their underlying kinds.
		if !enumCompatible(underlying, schemaType) {
			*issues = append(*issues, SchemaIssue{
				Path:    jsonPointer(path),
				Code:    CodeGoEnumTypeMismatch,
				Message: fmt.Sprintf("enum field requires string/integer-compatible Go type, got %s", typeName(goType)),
			})
			return
		}
		// Named types are accepted; plain underlying kinds are also accepted
		// (e.g. TypedRef.Kind string with a scope-kind enum).
		return
	}

	if schemaType == "" {
		return
	}

	if !primitiveCompatible(underlying, schemaType) {
		*issues = append(*issues, SchemaIssue{
			Path:    jsonPointer(path),
			Code:    CodeGoTypeMismatch,
			Message: fmt.Sprintf("schema type %q is incompatible with Go type %s", schemaType, typeName(goType)),
		})
	}
}

func verifyRequiredOptional(field goJSONField, required bool, path string, issues *[]SchemaIssue) {
	isPointer := field.typ.Kind() == reflect.Ptr

	if required {
		if field.omitempty || isPointer {
			*issues = append(*issues, SchemaIssue{
				Path:    jsonPointer(path),
				Code:    CodeGoRequiredMismatch,
				Message: fmt.Sprintf("required property %q must be a non-pointer Go field without omitempty", lastSegment(path)),
			})
		}
		return
	}

	if !field.omitempty && !isPointer {
		*issues = append(*issues, SchemaIssue{
			Path:    jsonPointer(path),
			Code:    CodeGoRequiredMismatch,
			Message: fmt.Sprintf("optional property %q must use omitempty or a pointer Go field", lastSegment(path)),
		})
	}
}

func collectJSONFields(t reflect.Type, basePath string) map[string]goJSONField {
	out := make(map[string]goJSONField)
	collectJSONFieldsInto(t, basePath, out, nil)
	return out
}

func collectJSONFieldsInto(t reflect.Type, basePath string, out map[string]goJSONField, seen map[reflect.Type]struct{}) {
	t = unwrapType(t)
	if t.Kind() != reflect.Struct {
		return
	}
	if seen == nil {
		seen = make(map[reflect.Type]struct{})
	}
	if _, ok := seen[t]; ok {
		return
	}
	seen[t] = struct{}{}

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if f.PkgPath != "" && !f.Anonymous {
			// Unexported non-anonymous fields are invisible to encoding/json.
			continue
		}

		tag := f.Tag.Get("json")
		name, omitempty, skip := parseJSONTag(tag)
		if skip {
			continue
		}

		ft := f.Type
		// Anonymous embed with no explicit JSON name: promote fields
		// (encoding/json behavior; json:",inline" is NOT honored).
		if f.Anonymous && name == "" && isStructKind(ft) {
			collectJSONFieldsInto(ft, basePath, out, seen)
			continue
		}
		if name == "" {
			name = f.Name
		}
		// First declaration wins for duplicate JSON names.
		if _, exists := out[name]; exists {
			continue
		}
		out[name] = goJSONField{
			jsonName:  name,
			typ:       ft,
			omitempty: omitempty,
			path:      joinPath(basePath, name),
		}
	}
}

func parseJSONTag(tag string) (name string, omitempty bool, skip bool) {
	if tag == "-" {
		return "", false, true
	}
	if tag == "" {
		return "", false, false
	}
	parts := strings.Split(tag, ",")
	name = parts[0]
	for _, opt := range parts[1:] {
		if opt == "omitempty" {
			omitempty = true
		}
	}
	return name, omitempty, false
}

func schemaTypeString(schema map[string]any) (string, bool) {
	switch v := schema["type"].(type) {
	case string:
		return v, true
	case []any:
		// Unsupported multi-type in FEATURE-0012 subset; treat as absent.
		return "", false
	default:
		return "", false
	}
}

func enumCompatible(t reflect.Type, schemaType string) bool {
	switch t.Kind() {
	case reflect.String:
		return schemaType == "" || schemaType == "string"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return schemaType == "" || schemaType == "integer" || schemaType == "number"
	default:
		return false
	}
}

func primitiveCompatible(t reflect.Type, schemaType string) bool {
	switch schemaType {
	case "string":
		return t.Kind() == reflect.String
	case "boolean":
		return t.Kind() == reflect.Bool
	case "integer":
		switch t.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return true
		default:
			return false
		}
	case "number":
		switch t.Kind() {
		case reflect.Float32, reflect.Float64,
			reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return true
		default:
			return false
		}
	case "null":
		return t.Kind() == reflect.Ptr || t.Kind() == reflect.Interface || t.Kind() == reflect.Map || t.Kind() == reflect.Slice
	default:
		return false
	}
}

func unwrapType(t reflect.Type) reflect.Type {
	for t != nil && t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t
}

func isStructKind(t reflect.Type) bool {
	t = unwrapType(t)
	return t != nil && t.Kind() == reflect.Struct
}

func typeName(t reflect.Type) string {
	if t == nil {
		return "<nil>"
	}
	if t.Name() != "" {
		if t.PkgPath() != "" {
			return t.PkgPath() + "." + t.Name()
		}
		return t.String()
	}
	return t.String()
}

func joinPath(base, name string) string {
	if base == "" {
		return "/" + escapePointerToken(name)
	}
	if strings.HasSuffix(base, "/") {
		return base + escapePointerToken(name)
	}
	return base + "/" + escapePointerToken(name)
}

func jsonPointer(path string) string {
	if path == "" {
		return "/"
	}
	return path
}

func escapePointerToken(s string) string {
	s = strings.ReplaceAll(s, "~", "~0")
	s = strings.ReplaceAll(s, "/", "~1")
	return s
}

func lastSegment(path string) string {
	path = strings.TrimPrefix(path, "/")
	if path == "" {
		return ""
	}
	parts := strings.Split(path, "/")
	seg := parts[len(parts)-1]
	seg = strings.ReplaceAll(seg, "~1", "/")
	seg = strings.ReplaceAll(seg, "~0", "~")
	return seg
}
