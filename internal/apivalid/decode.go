package apivalid

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
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
