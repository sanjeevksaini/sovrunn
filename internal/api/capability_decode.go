package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/sanjeevksaini/sovrunn/internal/resources"
)

// safeDecodeCapability applies http.MaxBytesReader, detects whether the JSON
// request body contains the key "status", then decodes into the typed Capability
// struct using DisallowUnknownFields. It reuses the shared decoder sentinel
// errors so error mapping stays consistent with existing resource decoders.
func safeDecodeCapability(w http.ResponseWriter, r *http.Request) (resources.Capability, error) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		var maxErr *http.MaxBytesError
		if errors.As(err, &maxErr) {
			return resources.Capability{}, errBodyTooLarge
		}
		return resources.Capability{}, errMalformedJSON
	}

	if len(body) == 0 {
		return resources.Capability{}, errEmptyBody
	}

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(body, &raw); err != nil {
		return resources.Capability{}, errMalformedJSON
	}
	if _, ok := raw["status"]; ok {
		return resources.Capability{}, errStatusFieldPresent
	}

	dec := json.NewDecoder(bytes.NewReader(body))
	dec.DisallowUnknownFields()
	var capability resources.Capability
	if err := dec.Decode(&capability); err != nil {
		var syntaxErr *json.SyntaxError
		var typeErr *json.UnmarshalTypeError
		if errors.As(err, &syntaxErr) || errors.As(err, &typeErr) {
			return resources.Capability{}, errMalformedJSON
		}
		if errors.Is(err, io.EOF) {
			return resources.Capability{}, errEmptyBody
		}
		if isUnknownFieldError(err) {
			return resources.Capability{}, errUnknownField
		}
		return resources.Capability{}, errMalformedJSON
	}
	return capability, nil
}
