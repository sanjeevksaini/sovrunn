package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

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
