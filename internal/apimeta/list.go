package apimeta

// ListEnvelope is the paginated collection response profile (F12-LIST-001).
//
// Embedding note (F12-NAMING-002 / F12-NAMING-005 correctness): the standard
// library encoding/json does NOT honor a `json:",inline"` tag. An anonymous
// embedded struct without a JSON tag already promotes its fields
// (apiVersion, kind) to the enclosing JSON object, which is the intended
// flat representation. DecodeYAML does not perform direct yaml.v3 typed
// decoding: accepted YAML is normalized to JSON and decoded through
// DecodeJSON, so YAML struct tags are neither required nor used by the
// typed decoder.
type ListEnvelope[T any] struct {
	TypeMeta      // JSON fields promoted through anonymous embedding
	Items    []T  `json:"items"`
	Page     Page `json:"page"`
}

// Page carries pagination metadata for ListEnvelope responses.
//
// NextPageToken is opaque: it MUST NOT expose database offsets, provider
// details, or authorization context (F12-LIST-002, F12-LIST-003).
// Total count is intentionally optional and omitted by default (F12-LIST-002).
type Page struct {
	NextPageToken string `json:"nextPageToken,omitempty"`
}
