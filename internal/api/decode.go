package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/apivalid"
	"github.com/sanjeevksaini/sovrunn/internal/cloudmodel/validate"
	"github.com/sanjeevksaini/sovrunn/internal/resources"
)

var (
	errBodyTooLarge       = errors.New("request body too large")
	errStatusFieldPresent = errors.New("status field present")
	errMalformedJSON      = errors.New("malformed JSON")
	errEmptyBody          = errors.New("request body is required")
	errUnknownField       = errors.New("unknown field")
)

// safeDecodeOrganization applies http.MaxBytesReader, detects whether the
// JSON request body contains the key "status", then decodes into the typed
// Organization struct using DisallowUnknownFields.
func safeDecodeOrganization(w http.ResponseWriter, r *http.Request) (resources.Organization, error) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		var maxErr *http.MaxBytesError
		if errors.As(err, &maxErr) {
			return resources.Organization{}, errBodyTooLarge
		}
		return resources.Organization{}, errMalformedJSON
	}

	if len(body) == 0 {
		return resources.Organization{}, errEmptyBody
	}

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(body, &raw); err != nil {
		return resources.Organization{}, errMalformedJSON
	}
	if _, ok := raw["status"]; ok {
		return resources.Organization{}, errStatusFieldPresent
	}

	dec := json.NewDecoder(bytes.NewReader(body))
	dec.DisallowUnknownFields()
	var org resources.Organization
	if err := dec.Decode(&org); err != nil {
		var syntaxErr *json.SyntaxError
		var typeErr *json.UnmarshalTypeError
		if errors.As(err, &syntaxErr) || errors.As(err, &typeErr) {
			return resources.Organization{}, errMalformedJSON
		}
		if errors.Is(err, io.EOF) {
			return resources.Organization{}, errEmptyBody
		}
		if isUnknownFieldError(err) {
			return resources.Organization{}, errUnknownField
		}
		return resources.Organization{}, errMalformedJSON
	}
	return org, nil
}

func isUnknownFieldError(err error) bool {
	msg := err.Error()
	return strings.Contains(msg, "unknown field")
}

// mapDecodeError maps decode errors to HTTP status and API error fields.
func mapDecodeError(err error) (int, string) {
	switch {
	case errors.Is(err, errBodyTooLarge):
		return http.StatusRequestEntityTooLarge, "request body too large"
	case errors.Is(err, errStatusFieldPresent):
		return http.StatusBadRequest, "status field is not allowed in request body"
	case errors.Is(err, errEmptyBody):
		return http.StatusBadRequest, "request body is required"
	case errors.Is(err, errUnknownField):
		return http.StatusBadRequest, "unknown field in request body"
	default:
		return http.StatusBadRequest, "malformed JSON"
	}
}

// ---------------------------------------------------------------------------
// FEATURE-0015 two-phase JSON decode (design §5.1.1)
// ---------------------------------------------------------------------------

// PhaseOneResult is the route-safe extraction produced by phase-one decode.
// It never persists, defaults, validates, or accepts the complete body.
type PhaseOneResult struct {
	Body              []byte
	Name              string
	Refs              map[string]apimeta.TypedRef // JSON path (dot form) -> TypedRef
	RawRoot           map[string]json.RawMessage
	HasBootstrapGrant bool
}

// DecodePhaseOne performs bounded, syntax-safe decoding used only to establish
// request safety before authorization and safe reference access.
// It enforces the inherited body-size limit, rejects malformed JSON and
// duplicate JSON members, and extracts only route-safe identifiers/references.
func DecodePhaseOne(data []byte, lim apivalid.Limits) (PhaseOneResult, *apiproblem.Problem) {
	var out PhaseOneResult
	if lim.MaxObjectBytes <= 0 {
		lim = apivalid.DefaultLimits()
	}
	if len(data) == 0 {
		return out, apiproblem.New(apiproblem.CodeMalformedRequest).WithDetail("request body is required")
	}
	if lim.MaxObjectBytes > 0 && len(data) > lim.MaxObjectBytes {
		return out, apiproblem.New(apiproblem.CodeRequestTooLarge).WithDetail("request body exceeds MaxObjectBytes")
	}

	// Reuse FEATURE-0012 duplicate/policy scan via DecodeJSON into a generic sink.
	var sink map[string]json.RawMessage
	pol := apivalid.FieldPolicy{
		Mode:              apivalid.ModeInternalObject,
		AllowStatus:       true,
		AllowSystemOwned:  true,
		AllowSpecMutation: true,
	}
	if prob := apivalid.DecodeJSON(data, lim, pol, &sink); prob != nil {
		// Phase one must not apply create field-policy denials; remap those
		// by re-scanning with a permissive policy already set. DecodeJSON
		// still rejects duplicates/malformed/unknown-for-dst. map[string]RawMessage
		// accepts any object members.
		return out, prob
	}
	out.Body = append([]byte(nil), data...)
	out.RawRoot = sink
	if _, ok := sink["bootstrapGrant"]; ok {
		out.HasBootstrapGrant = true
	}

	var generic map[string]any
	if err := json.Unmarshal(data, &generic); err != nil {
		return out, apiproblem.New(apiproblem.CodeMalformedRequest).WithDetail("malformed JSON")
	}
	out.Name = extractStringPath(generic, "metadata", "name")
	out.Refs = extractTypedRefs(generic)
	return out, nil
}

// DecodePhaseOneHTTP reads a bounded request body then runs DecodePhaseOne.
func DecodePhaseOneHTTP(w http.ResponseWriter, r *http.Request, lim apivalid.Limits) (PhaseOneResult, *apiproblem.Problem) {
	if lim.MaxObjectBytes <= 0 {
		lim = apivalid.DefaultLimits()
	}
	if r == nil {
		return PhaseOneResult{}, apiproblem.New(apiproblem.CodeMalformedRequest).WithDetail("request is required")
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
			return PhaseOneResult{}, apiproblem.New(apiproblem.CodeRequestTooLarge).WithDetail("request body exceeds MaxObjectBytes")
		}
		return PhaseOneResult{}, apiproblem.New(apiproblem.CodeMalformedRequest).WithDetail("failed to read request body")
	}
	return DecodePhaseOne(data, lim)
}

// DecodePhaseTwoStrict performs phase-two strict decode and closed-contract
// classification for a collection-create or PATCH body. An invalid or
// unclassified body must never be digested, stored, or replayed.
func DecodePhaseTwoStrict(kind validate.Kind, class validate.RequestClass, data []byte, dst any) *apiproblem.Problem {
	lim := apivalid.DefaultLimits()
	switch class {
	case validate.RequestClassCollectionCreate:
		if prob := validate.ClassifyCreateContract(kind, data); prob != nil {
			return prob
		}
		pol := apivalid.PolicyFor(apivalid.ModeCreateRequest)
		if prob := apivalid.DecodeJSON(data, lim, pol, dst); prob != nil {
			return prob
		}
		return nil
	case validate.RequestClassPatch:
		if prob := validate.ClassifyPatchContract(kind, data); prob != nil {
			return prob
		}
		// PATCH bodies are classified against the allowlist; typed decode of
		// the patch document itself is not required here because RFC 7396
		// merge operates on generic JSON (see ApplyMergePatch).
		return nil
	default:
		return apiproblem.New(apiproblem.CodeValidationFailed).WithDetail("unsupported decode classification")
	}
}

func extractStringPath(root map[string]any, parts ...string) string {
	cur := any(root)
	for _, p := range parts {
		m, ok := cur.(map[string]any)
		if !ok {
			return ""
		}
		cur, ok = m[p]
		if !ok {
			return ""
		}
	}
	s, _ := cur.(string)
	return s
}

func extractTypedRefs(root map[string]any) map[string]apimeta.TypedRef {
	out := make(map[string]apimeta.TypedRef)
	spec, _ := root["spec"].(map[string]any)
	if spec == nil {
		return out
	}
	for _, key := range []string{
		"cloudPlatformRef",
		"cloudProviderRef",
		"hostingLocationRef",
		"datacenterRef",
		"faultDomainRef",
	} {
		if ref, ok := parseTypedRef(spec[key]); ok {
			out["spec."+key] = ref
		}
	}
	return out
}

func parseTypedRef(v any) (apimeta.TypedRef, bool) {
	m, ok := v.(map[string]any)
	if !ok {
		return apimeta.TypedRef{}, false
	}
	ref := apimeta.TypedRef{
		APIVersion: asString(m["apiVersion"]),
		Kind:       asString(m["kind"]),
		Name:       asString(m["name"]),
		UID:        asString(m["uid"]),
	}
	if ref.APIVersion == "" && ref.Kind == "" && ref.Name == "" && ref.UID == "" {
		return apimeta.TypedRef{}, false
	}
	return ref, true
}

func asString(v any) string {
	s, _ := v.(string)
	return s
}
