package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/sanjeevksaini/sovrunn/internal/resources"
)

// safeDecodePlugin applies http.MaxBytesReader, detects whether the JSON
// request body contains the key "status", then decodes into the typed Plugin
// struct using DisallowUnknownFields. It reuses the shared decoder sentinel
// errors so error mapping stays consistent with existing resource decoders.
func safeDecodePlugin(w http.ResponseWriter, r *http.Request) (resources.Plugin, error) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		var maxErr *http.MaxBytesError
		if errors.As(err, &maxErr) {
			return resources.Plugin{}, errBodyTooLarge
		}
		return resources.Plugin{}, errMalformedJSON
	}

	if len(body) == 0 {
		return resources.Plugin{}, errEmptyBody
	}

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(body, &raw); err != nil {
		return resources.Plugin{}, errMalformedJSON
	}
	if _, ok := raw["status"]; ok {
		return resources.Plugin{}, errStatusFieldPresent
	}

	dec := json.NewDecoder(bytes.NewReader(body))
	dec.DisallowUnknownFields()
	var plugin resources.Plugin
	if err := dec.Decode(&plugin); err != nil {
		var syntaxErr *json.SyntaxError
		var typeErr *json.UnmarshalTypeError
		if errors.As(err, &syntaxErr) || errors.As(err, &typeErr) {
			return resources.Plugin{}, errMalformedJSON
		}
		if errors.Is(err, io.EOF) {
			return resources.Plugin{}, errEmptyBody
		}
		if isUnknownFieldError(err) {
			return resources.Plugin{}, errUnknownField
		}
		return resources.Plugin{}, errMalformedJSON
	}
	return plugin, nil
}
