package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/sanjeevksaini/sovrunn/internal/resources"
)

// safeDecodeServiceClass applies http.MaxBytesReader, detects whether the JSON
// request body contains the key "status", then decodes into the typed
// ServiceClass struct using DisallowUnknownFields. It reuses the sentinel
// errors and unknown-field detection defined in decode.go so error mapping
// stays consistent with the existing resource decoders.
func safeDecodeServiceClass(w http.ResponseWriter, r *http.Request) (resources.ServiceClass, error) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		var maxErr *http.MaxBytesError
		if errors.As(err, &maxErr) {
			return resources.ServiceClass{}, errBodyTooLarge
		}
		return resources.ServiceClass{}, errMalformedJSON
	}

	if len(body) == 0 {
		return resources.ServiceClass{}, errEmptyBody
	}

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(body, &raw); err != nil {
		return resources.ServiceClass{}, errMalformedJSON
	}
	if _, ok := raw["status"]; ok {
		return resources.ServiceClass{}, errStatusFieldPresent
	}

	dec := json.NewDecoder(bytes.NewReader(body))
	dec.DisallowUnknownFields()
	var sc resources.ServiceClass
	if err := dec.Decode(&sc); err != nil {
		var syntaxErr *json.SyntaxError
		var typeErr *json.UnmarshalTypeError
		if errors.As(err, &syntaxErr) || errors.As(err, &typeErr) {
			return resources.ServiceClass{}, errMalformedJSON
		}
		if errors.Is(err, io.EOF) {
			return resources.ServiceClass{}, errEmptyBody
		}
		if isUnknownFieldError(err) {
			return resources.ServiceClass{}, errUnknownField
		}
		return resources.ServiceClass{}, errMalformedJSON
	}
	return sc, nil
}
