package apivalid

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"gopkg.in/yaml.v3"
)

// systemOwnedMetadataKeys are Matrix C2 system-only ObjectMeta fields
// rejected when AllowSystemOwned is false (D-15, F12-META-002).
var systemOwnedMetadataKeys = map[string]struct{}{
	"uid":             {},
	"generation":      {},
	"resourceVersion": {},
	"createdAt":       {},
	"updatedAt":       {},
}

// DecodeJSON is a pure JSON decoder: encoding/json with DisallowUnknownFields
// plus a token-scan duplicate-key detector, applying FieldPolicy and Limits.
// It has no HTTP dependency (D-03, D-15).
//
// On failure it returns an *apiproblem.Problem with a stable ErrorCode and an
// RFC 6901 JSON Pointer on Violations[].Field (F12-VALIDATION-006).
// On success it returns nil and writes into dst.
//
// Request bodies and field values are not logged (F12-SEC-003).
func DecodeJSON(data []byte, lim Limits, pol FieldPolicy, dst any) *apiproblem.Problem {
	if len(data) == 0 {
		return decodeProblem(apiproblem.CodeMalformedRequest, "/", "request body is required")
	}
	if lim.MaxObjectBytes > 0 && len(data) > lim.MaxObjectBytes {
		return decodeProblem(apiproblem.CodeRequestTooLarge, "/", "request body exceeds MaxObjectBytes")
	}

	if prob := scanJSONDuplicatesAndPolicy(data, lim, pol); prob != nil {
		return prob
	}

	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil {
		return mapJSONDecodeError(err)
	}

	// Reject trailing tokens after the first JSON value.
	if _, err := dec.Token(); err != io.EOF {
		if err == nil {
			return decodeProblem(apiproblem.CodeMalformedRequest, "/", "trailing data after JSON value")
		}
		if !errors.Is(err, io.EOF) {
			return decodeProblem(apiproblem.CodeMalformedRequest, "/", "malformed JSON after value")
		}
	}
	return nil
}

// DecodeYAML is a pure YAML decoder that accepts only a strict JSON-compatible
// YAML subset (D-03a). It never performs direct yaml.v3 typed decoding into the
// destination, never depends on YAML struct tags, and never enables
// KnownFields(true). The ordered pipeline is:
//
//  1. yaml.Node safety parsing (syntax tree only)
//  2. Reject YAML-only constructs (aliases, anchors, merge keys, custom/explicit
//     tags, non-finite numbers, timestamp/binary coercions, multiple documents,
//     non-string mapping keys)
//  3. Explicit yaml.Node duplicate-key detection
//  4. Normalize the accepted YAML node to a JSON-compatible value
//  5. Marshal that value to JSON bytes
//  6. Decode those bytes through DecodeJSON (unknown-field rejection,
//     FieldPolicy, same destination type, same stable codes and JSON Pointers)
//
// Request bodies and field values are not logged (F12-SEC-003).
func DecodeYAML(data []byte, lim Limits, pol FieldPolicy, dst any) *apiproblem.Problem {
	if len(data) == 0 {
		return decodeProblem(apiproblem.CodeMalformedRequest, "/", "request body is required")
	}
	if lim.MaxObjectBytes > 0 && len(data) > lim.MaxObjectBytes {
		return decodeProblem(apiproblem.CodeRequestTooLarge, "/", "request body exceeds MaxObjectBytes")
	}

	// 1. yaml.Node safety parsing (syntax tree only); require exactly one document.
	dec := yaml.NewDecoder(bytes.NewReader(data))
	var doc yaml.Node
	if err := dec.Decode(&doc); err != nil {
		if errors.Is(err, io.EOF) {
			return decodeProblem(apiproblem.CodeMalformedRequest, "/", "YAML document is required")
		}
		return decodeProblem(apiproblem.CodeMalformedRequest, "/", "malformed YAML")
	}
	var trailing yaml.Node
	if err := dec.Decode(&trailing); err == nil {
		return decodeProblem(apiproblem.CodeMalformedRequest, "/", "multiple YAML documents are not permitted")
	} else if !errors.Is(err, io.EOF) {
		return decodeProblem(apiproblem.CodeMalformedRequest, "/", "malformed YAML after document")
	}

	root := yamlDocumentRoot(&doc)
	if root == nil {
		return decodeProblem(apiproblem.CodeMalformedRequest, "/", "YAML document is required")
	}

	// 2. Reject YAML-only constructs.
	if prob := rejectYAMLOnlyConstructs(root, ""); prob != nil {
		return prob
	}

	// 3. Explicit yaml.Node duplicate-key detection pass.
	if prob := rejectYAMLDuplicateKeys(root, ""); prob != nil {
		return prob
	}

	// 4. Normalize the accepted YAML node to a JSON-compatible value.
	normalized, prob := normalizeYAMLNode(root)
	if prob != nil {
		return prob
	}

	// 5. Marshal that normalized value to JSON bytes.
	jsonBytes, err := json.Marshal(normalized)
	if err != nil {
		return decodeProblem(apiproblem.CodeMalformedRequest, "/", "YAML could not be normalized to JSON")
	}

	// 6. Pass those JSON bytes through the same DecodeJSON path.
	return DecodeJSON(jsonBytes, lim, pol, dst)
}

func yamlDocumentRoot(doc *yaml.Node) *yaml.Node {
	if doc == nil {
		return nil
	}
	if doc.Kind == yaml.DocumentNode {
		if len(doc.Content) != 1 {
			return nil
		}
		return doc.Content[0]
	}
	return doc
}

// rejectYAMLOnlyConstructs rejects aliases, anchors, merge keys, custom/explicit
// tags, non-finite numbers, YAML-only timestamp/binary coercions, and non-string
// mapping keys (D-03a step 2).
func rejectYAMLOnlyConstructs(n *yaml.Node, path string) *apiproblem.Problem {
	if n == nil {
		return decodeProblem(apiproblem.CodeMalformedRequest, pathOrRoot(path), "malformed YAML node")
	}
	if n.Kind == yaml.AliasNode || n.Alias != nil {
		return decodeProblem(apiproblem.CodeMalformedRequest, pathOrRoot(path), "YAML aliases are not permitted")
	}
	if n.Anchor != "" {
		return decodeProblem(apiproblem.CodeMalformedRequest, pathOrRoot(path), "YAML anchors are not permitted")
	}
	if n.Style&yaml.TaggedStyle != 0 {
		return decodeProblem(apiproblem.CodeMalformedRequest, pathOrRoot(path), "YAML explicit tags are not permitted")
	}

	short := n.ShortTag()
	switch n.Kind {
	case yaml.DocumentNode:
		if len(n.Content) != 1 {
			return decodeProblem(apiproblem.CodeMalformedRequest, "/", "YAML document is required")
		}
		return rejectYAMLOnlyConstructs(n.Content[0], path)
	case yaml.MappingNode:
		if short != "!!map" {
			return decodeProblem(apiproblem.CodeMalformedRequest, pathOrRoot(path), "unsupported YAML mapping tag")
		}
		if len(n.Content)%2 != 0 {
			return decodeProblem(apiproblem.CodeMalformedRequest, pathOrRoot(path), "malformed YAML mapping")
		}
		for i := 0; i < len(n.Content); i += 2 {
			key := n.Content[i]
			val := n.Content[i+1]
			if key == nil {
				return decodeProblem(apiproblem.CodeMalformedRequest, pathOrRoot(path), "malformed YAML mapping key")
			}
			if key.Kind == yaml.AliasNode || key.Alias != nil {
				return decodeProblem(apiproblem.CodeMalformedRequest, pathOrRoot(path), "YAML aliases are not permitted")
			}
			if key.Anchor != "" {
				return decodeProblem(apiproblem.CodeMalformedRequest, pathOrRoot(path), "YAML anchors are not permitted")
			}
			if key.Style&yaml.TaggedStyle != 0 {
				return decodeProblem(apiproblem.CodeMalformedRequest, pathOrRoot(path), "YAML explicit tags are not permitted")
			}
			keyTag := key.ShortTag()
			if keyTag == "!!merge" || (key.Kind == yaml.ScalarNode && key.Value == "<<" && keyTag != "!!str") {
				childPath := path + "/" + escapeJSONPointer("<<")
				return decodeProblem(apiproblem.CodeMalformedRequest, pathOrRoot(childPath), "YAML merge keys are not permitted")
			}
			if key.Kind != yaml.ScalarNode || keyTag != "!!str" {
				return decodeProblem(apiproblem.CodeMalformedRequest, pathOrRoot(path), "YAML mapping keys must be strings")
			}
			childPath := path + "/" + escapeJSONPointer(key.Value)
			if prob := rejectYAMLOnlyConstructs(val, childPath); prob != nil {
				return prob
			}
		}
		return nil
	case yaml.SequenceNode:
		if short != "!!seq" {
			return decodeProblem(apiproblem.CodeMalformedRequest, pathOrRoot(path), "unsupported YAML sequence tag")
		}
		for i, child := range n.Content {
			childPath := fmt.Sprintf("%s/%d", path, i)
			if prob := rejectYAMLOnlyConstructs(child, childPath); prob != nil {
				return prob
			}
		}
		return nil
	case yaml.ScalarNode:
		switch short {
		case "!!null", "!!bool", "!!int", "!!str":
			return nil
		case "!!float":
			if isNonFiniteYAMLFloat(n.Value) {
				return decodeProblem(apiproblem.CodeMalformedRequest, pathOrRoot(path), "non-finite YAML numbers are not permitted")
			}
			return nil
		case "!!timestamp":
			return decodeProblem(apiproblem.CodeMalformedRequest, pathOrRoot(path), "YAML timestamp coercions are not permitted")
		case "!!binary":
			return decodeProblem(apiproblem.CodeMalformedRequest, pathOrRoot(path), "YAML binary coercions are not permitted")
		case "!!merge":
			return decodeProblem(apiproblem.CodeMalformedRequest, pathOrRoot(path), "YAML merge keys are not permitted")
		default:
			return decodeProblem(apiproblem.CodeMalformedRequest, pathOrRoot(path), "YAML custom tags are not permitted")
		}
	default:
		return decodeProblem(apiproblem.CodeMalformedRequest, pathOrRoot(path), "unsupported YAML node kind")
	}
}

func isNonFiniteYAMLFloat(v string) bool {
	switch strings.ToLower(strings.TrimSpace(v)) {
	case ".nan", ".inf", "+.inf", "-.inf":
		return true
	default:
		return false
	}
}

// rejectYAMLDuplicateKeys walks mappings and rejects duplicate string keys
// with DUPLICATE_FIELD and an RFC 6901 JSON Pointer (D-03a step 3).
func rejectYAMLDuplicateKeys(n *yaml.Node, path string) *apiproblem.Problem {
	if n == nil {
		return nil
	}
	switch n.Kind {
	case yaml.DocumentNode:
		if len(n.Content) == 1 {
			return rejectYAMLDuplicateKeys(n.Content[0], path)
		}
		return nil
	case yaml.MappingNode:
		seen := make(map[string]struct{}, len(n.Content)/2)
		for i := 0; i+1 < len(n.Content); i += 2 {
			key := n.Content[i]
			val := n.Content[i+1]
			childPath := path + "/" + escapeJSONPointer(key.Value)
			if _, dup := seen[key.Value]; dup {
				return decodeProblem(apiproblem.CodeDuplicateField, childPath, "duplicate field")
			}
			seen[key.Value] = struct{}{}
			if prob := rejectYAMLDuplicateKeys(val, childPath); prob != nil {
				return prob
			}
		}
		return nil
	case yaml.SequenceNode:
		for i, child := range n.Content {
			childPath := fmt.Sprintf("%s/%d", path, i)
			if prob := rejectYAMLDuplicateKeys(child, childPath); prob != nil {
				return prob
			}
		}
		return nil
	default:
		return nil
	}
}

// normalizeYAMLNode converts an accepted YAML node into a JSON-compatible
// Go value (maps with string keys, slices, strings, bools, nil, json.Number).
func normalizeYAMLNode(n *yaml.Node) (any, *apiproblem.Problem) {
	if n == nil {
		return nil, decodeProblem(apiproblem.CodeMalformedRequest, "/", "malformed YAML node")
	}
	switch n.Kind {
	case yaml.DocumentNode:
		if len(n.Content) != 1 {
			return nil, decodeProblem(apiproblem.CodeMalformedRequest, "/", "YAML document is required")
		}
		return normalizeYAMLNode(n.Content[0])
	case yaml.MappingNode:
		out := make(map[string]any, len(n.Content)/2)
		for i := 0; i+1 < len(n.Content); i += 2 {
			key := n.Content[i]
			val, prob := normalizeYAMLNode(n.Content[i+1])
			if prob != nil {
				return nil, prob
			}
			out[key.Value] = val
		}
		return out, nil
	case yaml.SequenceNode:
		out := make([]any, 0, len(n.Content))
		for _, child := range n.Content {
			val, prob := normalizeYAMLNode(child)
			if prob != nil {
				return nil, prob
			}
			out = append(out, val)
		}
		return out, nil
	case yaml.ScalarNode:
		switch n.ShortTag() {
		case "!!null":
			return nil, nil
		case "!!bool":
			switch strings.ToLower(n.Value) {
			case "true":
				return true, nil
			case "false":
				return false, nil
			default:
				return nil, decodeProblem(apiproblem.CodeMalformedRequest, "/", "malformed YAML bool")
			}
		case "!!int", "!!float":
			return json.Number(n.Value), nil
		case "!!str":
			return n.Value, nil
		default:
			return nil, decodeProblem(apiproblem.CodeMalformedRequest, "/", "unsupported YAML scalar during normalize")
		}
	default:
		return nil, decodeProblem(apiproblem.CodeMalformedRequest, "/", "unsupported YAML node during normalize")
	}
}

func scanJSONDuplicatesAndPolicy(data []byte, lim Limits, pol FieldPolicy) *apiproblem.Problem {
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()
	return scanJSONValue(dec, "", 0, lim, pol)
}

func scanJSONValue(dec *json.Decoder, path string, depth int, lim Limits, pol FieldPolicy) *apiproblem.Problem {
	tok, err := dec.Token()
	if err != nil {
		if errors.Is(err, io.EOF) {
			return decodeProblem(apiproblem.CodeMalformedRequest, pathOrRoot(path), "unexpected end of JSON input")
		}
		return decodeProblem(apiproblem.CodeMalformedRequest, pathOrRoot(path), "malformed JSON")
	}

	switch t := tok.(type) {
	case json.Delim:
		switch t {
		case '{':
			return scanJSONObject(dec, path, depth+1, lim, pol)
		case '[':
			return scanJSONArray(dec, path, depth+1, lim, pol)
		default:
			return decodeProblem(apiproblem.CodeMalformedRequest, pathOrRoot(path), "malformed JSON delimiter")
		}
	case bool, float64, json.Number, string, nil:
		return nil
	default:
		return decodeProblem(apiproblem.CodeMalformedRequest, pathOrRoot(path), "malformed JSON token")
	}
}

func scanJSONObject(dec *json.Decoder, path string, depth int, lim Limits, pol FieldPolicy) *apiproblem.Problem {
	if lim.MaxNestingDepth > 0 && depth > lim.MaxNestingDepth {
		return decodeProblem(apiproblem.CodeRequestTooLarge, pathOrRoot(path), "JSON nesting depth exceeds MaxNestingDepth")
	}

	seen := make(map[string]struct{})
	for dec.More() {
		keyTok, err := dec.Token()
		if err != nil {
			return decodeProblem(apiproblem.CodeMalformedRequest, pathOrRoot(path), "malformed JSON object key")
		}
		key, ok := keyTok.(string)
		if !ok {
			return decodeProblem(apiproblem.CodeMalformedRequest, pathOrRoot(path), "JSON object key must be a string")
		}

		childPath := path + "/" + escapeJSONPointer(key)
		if _, dup := seen[key]; dup {
			return decodeProblem(apiproblem.CodeDuplicateField, childPath, "duplicate field")
		}
		seen[key] = struct{}{}

		if prob := checkFieldPolicy(path, key, childPath, pol); prob != nil {
			return prob
		}
		if prob := scanJSONValue(dec, childPath, depth, lim, pol); prob != nil {
			return prob
		}
	}

	// Consume closing '}'.
	tok, err := dec.Token()
	if err != nil {
		return decodeProblem(apiproblem.CodeMalformedRequest, pathOrRoot(path), "malformed JSON object")
	}
	if delim, ok := tok.(json.Delim); !ok || delim != '}' {
		return decodeProblem(apiproblem.CodeMalformedRequest, pathOrRoot(path), "malformed JSON object close")
	}
	return nil
}

func scanJSONArray(dec *json.Decoder, path string, depth int, lim Limits, pol FieldPolicy) *apiproblem.Problem {
	if lim.MaxNestingDepth > 0 && depth > lim.MaxNestingDepth {
		return decodeProblem(apiproblem.CodeRequestTooLarge, pathOrRoot(path), "JSON nesting depth exceeds MaxNestingDepth")
	}

	index := 0
	for dec.More() {
		childPath := fmt.Sprintf("%s/%d", path, index)
		if prob := scanJSONValue(dec, childPath, depth, lim, pol); prob != nil {
			return prob
		}
		index++
	}

	tok, err := dec.Token()
	if err != nil {
		return decodeProblem(apiproblem.CodeMalformedRequest, pathOrRoot(path), "malformed JSON array")
	}
	if delim, ok := tok.(json.Delim); !ok || delim != ']' {
		return decodeProblem(apiproblem.CodeMalformedRequest, pathOrRoot(path), "malformed JSON array close")
	}
	return nil
}

func checkFieldPolicy(parentPath, key, childPath string, pol FieldPolicy) *apiproblem.Problem {
	switch parentPath {
	case "":
		if key == "status" && !pol.AllowStatus {
			return decodeProblem(apiproblem.CodeValidationFailed, childPath, "status is not permitted for this decode mode")
		}
		if key == "spec" && !pol.AllowSpecMutation {
			return decodeProblem(apiproblem.CodeValidationFailed, childPath, "spec mutation is not permitted for this decode mode")
		}
	case "/metadata":
		if _, system := systemOwnedMetadataKeys[key]; system && !pol.AllowSystemOwned {
			return decodeProblem(apiproblem.CodeValidationFailed, childPath, "system-owned metadata field is not permitted for this decode mode")
		}
	}
	return nil
}

func mapJSONDecodeError(err error) *apiproblem.Problem {
	if err == nil {
		return nil
	}
	if errors.Is(err, io.EOF) {
		return decodeProblem(apiproblem.CodeMalformedRequest, "/", "unexpected end of JSON input")
	}

	var syntaxErr *json.SyntaxError
	if errors.As(err, &syntaxErr) {
		return decodeProblem(apiproblem.CodeMalformedRequest, "/", "malformed JSON")
	}

	var typeErr *json.UnmarshalTypeError
	if errors.As(err, &typeErr) {
		ptr := "/"
		if typeErr.Field != "" {
			ptr = goFieldPathToJSONPointer(typeErr.Field)
		}
		return decodeProblem(apiproblem.CodeMalformedRequest, ptr, "JSON value type mismatch")
	}

	msg := err.Error()
	const unknownPrefix = "json: unknown field "
	if strings.HasPrefix(msg, unknownPrefix) {
		field := strings.Trim(msg[len(unknownPrefix):], `"`)
		return decodeProblem(apiproblem.CodeUnknownField, "/"+escapeJSONPointer(field), "unknown field")
	}

	return decodeProblem(apiproblem.CodeMalformedRequest, "/", "malformed JSON")
}

func decodeProblem(code apiproblem.ErrorCode, field, message string) *apiproblem.Problem {
	vcode := violationCodeFor(code)
	return apiproblem.New(code).
		WithDetail(message).
		WithViolations([]apiproblem.Violation{{
			Field:   field,
			Code:    vcode,
			Message: message,
		}})
}

func violationCodeFor(code apiproblem.ErrorCode) apiproblem.ViolationCode {
	switch code {
	case apiproblem.CodeUnknownField:
		return apiproblem.ViolationUnknownField
	case apiproblem.CodeDuplicateField:
		return apiproblem.ViolationDuplicateField
	default:
		// Carry the problem-level stable code on the violation when no
		// dedicated ViolationCode exists (e.g. MALFORMED_REQUEST,
		// VALIDATION_FAILED, REQUEST_TOO_LARGE).
		return apiproblem.ViolationCode(code)
	}
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

// goFieldPathToJSONPointer converts encoding/json UnmarshalTypeError.Field
// (dot-separated Go field names) into a best-effort RFC 6901 pointer.
func goFieldPathToJSONPointer(field string) string {
	if field == "" {
		return "/"
	}
	parts := strings.Split(field, ".")
	var b strings.Builder
	for _, p := range parts {
		b.WriteByte('/')
		b.WriteString(escapeJSONPointer(p))
	}
	return b.String()
}
