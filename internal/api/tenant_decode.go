package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/sanjeevksaini/sovrunn/internal/resources"
)

// safeDecodeTenant applies http.MaxBytesReader, detects whether the JSON
// request body contains the key "status", then decodes into the typed Tenant
// struct using DisallowUnknownFields. It reuses the sentinel errors and
// unknown-field detection defined in decode.go so error mapping stays
// consistent with the existing resource decoders.
func safeDecodeTenant(w http.ResponseWriter, r *http.Request) (resources.Tenant, error) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		var maxErr *http.MaxBytesError
		if errors.As(err, &maxErr) {
			return resources.Tenant{}, errBodyTooLarge
		}
		return resources.Tenant{}, errMalformedJSON
	}

	if len(body) == 0 {
		return resources.Tenant{}, errEmptyBody
	}

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(body, &raw); err != nil {
		return resources.Tenant{}, errMalformedJSON
	}
	if _, ok := raw["status"]; ok {
		return resources.Tenant{}, errStatusFieldPresent
	}

	dec := json.NewDecoder(bytes.NewReader(body))
	dec.DisallowUnknownFields()
	var t resources.Tenant
	if err := dec.Decode(&t); err != nil {
		var syntaxErr *json.SyntaxError
		var typeErr *json.UnmarshalTypeError
		if errors.As(err, &syntaxErr) || errors.As(err, &typeErr) {
			return resources.Tenant{}, errMalformedJSON
		}
		if errors.Is(err, io.EOF) {
			return resources.Tenant{}, errEmptyBody
		}
		if isUnknownFieldError(err) {
			return resources.Tenant{}, errUnknownField
		}
		return resources.Tenant{}, errMalformedJSON
	}
	return t, nil
}
