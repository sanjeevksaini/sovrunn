package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/sanjeevksaini/sovrunn/internal/resources"
)

// safeDecodeServiceInstance applies http.MaxBytesReader, detects whether the JSON
// request body contains the key "status", then decodes into the typed
// ServiceInstance struct using DisallowUnknownFields. It reuses the sentinel
// errors and unknown-field detection defined in decode.go so error mapping
// stays consistent with the existing resource decoders.
func safeDecodeServiceInstance(w http.ResponseWriter, r *http.Request) (resources.ServiceInstance, error) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		var maxErr *http.MaxBytesError
		if errors.As(err, &maxErr) {
			return resources.ServiceInstance{}, errBodyTooLarge
		}
		return resources.ServiceInstance{}, errMalformedJSON
	}

	if len(body) == 0 {
		return resources.ServiceInstance{}, errEmptyBody
	}

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(body, &raw); err != nil {
		return resources.ServiceInstance{}, errMalformedJSON
	}
	if _, ok := raw["status"]; ok {
		return resources.ServiceInstance{}, errStatusFieldPresent
	}

	dec := json.NewDecoder(bytes.NewReader(body))
	dec.DisallowUnknownFields()
	var si resources.ServiceInstance
	if err := dec.Decode(&si); err != nil {
		var syntaxErr *json.SyntaxError
		var typeErr *json.UnmarshalTypeError
		if errors.As(err, &syntaxErr) || errors.As(err, &typeErr) {
			return resources.ServiceInstance{}, errMalformedJSON
		}
		if errors.Is(err, io.EOF) {
			return resources.ServiceInstance{}, errEmptyBody
		}
		if isUnknownFieldError(err) {
			return resources.ServiceInstance{}, errUnknownField
		}
		return resources.ServiceInstance{}, errMalformedJSON
	}
	return si, nil
}
