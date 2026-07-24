package apivalid

import (
	"errors"
	"io"
	"mime"
	"net/http"
	"strings"

	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
)

// mediaKind classifies an accepted request Content-Type for decoder selection.
type mediaKind int

const (
	mediaUnsupported mediaKind = iota
	mediaJSON
	mediaYAML
)

// StrictDecode is the HTTP adapter for layer-1 HTTP/content/size checks
// (D-03, F12-VALIDATION-001(1), F12-ERROR-002). It:
//
//  1. selects a decoder by Content-Type — application/json, or
//     application/yaml (with application/x-yaml and text/yaml treated as
//     yaml); any other media type maps to 415 UNSUPPORTED_MEDIA_TYPE
//  2. bounds the body with http.MaxBytesReader using lim.MaxObjectBytes;
//     an oversized body maps to 400 REQUEST_TOO_LARGE (not 413)
//  3. delegates to DecodeJSON or DecodeYAML with PolicyFor(mode)
//
// StrictDecode never writes the HTTP response; callers map the returned
// Problem to the wire. Request bodies and field values are not logged
// (F12-SEC-003). RequestID correlation is left to the HTTP handler.
func StrictDecode(w http.ResponseWriter, r *http.Request, lim Limits, mode DecodeMode, dst any) *apiproblem.Problem {
	if r == nil {
		return decodeProblem(apiproblem.CodeMalformedRequest, "/", "request is required")
	}

	kind, prob := classifyRequestMediaType(r.Header.Get("Content-Type"))
	if prob != nil {
		return prob
	}

	body := r.Body
	if body == nil {
		body = http.NoBody
	}
	if lim.MaxObjectBytes > 0 {
		body = http.MaxBytesReader(w, body, int64(lim.MaxObjectBytes))
	}

	data, err := io.ReadAll(body)
	if err != nil {
		var maxErr *http.MaxBytesError
		if errors.As(err, &maxErr) {
			return decodeProblem(apiproblem.CodeRequestTooLarge, "/", "request body exceeds MaxObjectBytes")
		}
		return decodeProblem(apiproblem.CodeMalformedRequest, "/", "failed to read request body")
	}

	pol := PolicyFor(mode)
	switch kind {
	case mediaJSON:
		return DecodeJSON(data, lim, pol, dst)
	case mediaYAML:
		return DecodeYAML(data, lim, pol, dst)
	default:
		// classifyRequestMediaType already rejects unsupported types.
		return decodeProblem(apiproblem.CodeUnsupportedMediaType, "/", "unsupported media type")
	}
}

// classifyRequestMediaType maps a Content-Type header to an accepted decoder
// kind. Parameters (e.g. charset=utf-8) are ignored. Comparison is
// case-insensitive.
func classifyRequestMediaType(contentType string) (mediaKind, *apiproblem.Problem) {
	ct := strings.TrimSpace(contentType)
	if ct == "" {
		return mediaUnsupported, decodeProblem(
			apiproblem.CodeUnsupportedMediaType,
			"/",
			"Content-Type is required; accepted types are application/json and application/yaml",
		)
	}

	mediaType, _, err := mime.ParseMediaType(ct)
	if err != nil {
		return mediaUnsupported, decodeProblem(
			apiproblem.CodeUnsupportedMediaType,
			"/",
			"Content-Type could not be parsed",
		)
	}

	switch strings.ToLower(mediaType) {
	case "application/json":
		return mediaJSON, nil
	case "application/yaml", "application/x-yaml", "text/yaml":
		return mediaYAML, nil
	default:
		return mediaUnsupported, decodeProblem(
			apiproblem.CodeUnsupportedMediaType,
			"/",
			"unsupported media type; accepted types are application/json and application/yaml",
		)
	}
}
