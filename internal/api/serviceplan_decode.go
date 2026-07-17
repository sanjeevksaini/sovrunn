package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/sanjeevksaini/sovrunn/internal/resources"
)

// safeDecodeServicePlan applies http.MaxBytesReader, detects whether the JSON
// request body contains the key "status", then decodes into the typed
// ServicePlan struct using DisallowUnknownFields. It reuses the sentinel
// errors and unknown-field detection defined in decode.go so error mapping
// stays consistent with the existing resource decoders.
func safeDecodeServicePlan(w http.ResponseWriter, r *http.Request) (resources.ServicePlan, error) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		var maxErr *http.MaxBytesError
		if errors.As(err, &maxErr) {
			return resources.ServicePlan{}, errBodyTooLarge
		}
		return resources.ServicePlan{}, errMalformedJSON
	}

	if len(body) == 0 {
		return resources.ServicePlan{}, errEmptyBody
	}

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(body, &raw); err != nil {
		return resources.ServicePlan{}, errMalformedJSON
	}
	if _, ok := raw["status"]; ok {
		return resources.ServicePlan{}, errStatusFieldPresent
	}

	dec := json.NewDecoder(bytes.NewReader(body))
	dec.DisallowUnknownFields()
	var sp resources.ServicePlan
	if err := dec.Decode(&sp); err != nil {
		var syntaxErr *json.SyntaxError
		var typeErr *json.UnmarshalTypeError
		if errors.As(err, &syntaxErr) || errors.As(err, &typeErr) {
			return resources.ServicePlan{}, errMalformedJSON
		}
		if errors.Is(err, io.EOF) {
			return resources.ServicePlan{}, errEmptyBody
		}
		if isUnknownFieldError(err) {
			return resources.ServicePlan{}, errUnknownField
		}
		return resources.ServicePlan{}, errMalformedJSON
	}
	return sp, nil
}
