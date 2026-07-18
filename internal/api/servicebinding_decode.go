package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/sanjeevksaini/sovrunn/internal/resources"
)

// safeDecodeServiceBinding applies http.MaxBytesReader, detects whether the JSON
// request body contains the key "status", then decodes into the typed
// ServiceBinding struct using DisallowUnknownFields. It reuses the sentinel
// errors and unknown-field detection defined in decode.go so error mapping
// stays consistent with the existing resource decoders.
func safeDecodeServiceBinding(w http.ResponseWriter, r *http.Request) (resources.ServiceBinding, error) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		var maxErr *http.MaxBytesError
		if errors.As(err, &maxErr) {
			return resources.ServiceBinding{}, errBodyTooLarge
		}
		return resources.ServiceBinding{}, errMalformedJSON
	}

	if len(body) == 0 {
		return resources.ServiceBinding{}, errEmptyBody
	}

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(body, &raw); err != nil {
		return resources.ServiceBinding{}, errMalformedJSON
	}
	if _, ok := raw["status"]; ok {
		return resources.ServiceBinding{}, errStatusFieldPresent
	}

	dec := json.NewDecoder(bytes.NewReader(body))
	dec.DisallowUnknownFields()
	var sb resources.ServiceBinding
	if err := dec.Decode(&sb); err != nil {
		var syntaxErr *json.SyntaxError
		var typeErr *json.UnmarshalTypeError
		if errors.As(err, &syntaxErr) || errors.As(err, &typeErr) {
			return resources.ServiceBinding{}, errMalformedJSON
		}
		if errors.Is(err, io.EOF) {
			return resources.ServiceBinding{}, errEmptyBody
		}
		if isUnknownFieldError(err) {
			return resources.ServiceBinding{}, errUnknownField
		}
		return resources.ServiceBinding{}, errMalformedJSON
	}
	return sb, nil
}
