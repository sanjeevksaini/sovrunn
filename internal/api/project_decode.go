package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/sanjeevksaini/sovrunn/internal/resources"
)

// safeDecodeProject applies http.MaxBytesReader, detects whether the JSON
// request body contains the key "status", then decodes into the typed Project
// struct using DisallowUnknownFields. It reuses the sentinel errors and
// unknown-field detection defined in decode.go so error mapping stays
// consistent with the existing resource decoders.
func safeDecodeProject(w http.ResponseWriter, r *http.Request) (resources.Project, error) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		var maxErr *http.MaxBytesError
		if errors.As(err, &maxErr) {
			return resources.Project{}, errBodyTooLarge
		}
		return resources.Project{}, errMalformedJSON
	}

	if len(body) == 0 {
		return resources.Project{}, errEmptyBody
	}

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(body, &raw); err != nil {
		return resources.Project{}, errMalformedJSON
	}
	if _, ok := raw["status"]; ok {
		return resources.Project{}, errStatusFieldPresent
	}

	dec := json.NewDecoder(bytes.NewReader(body))
	dec.DisallowUnknownFields()
	var p resources.Project
	if err := dec.Decode(&p); err != nil {
		var syntaxErr *json.SyntaxError
		var typeErr *json.UnmarshalTypeError
		if errors.As(err, &syntaxErr) || errors.As(err, &typeErr) {
			return resources.Project{}, errMalformedJSON
		}
		if errors.Is(err, io.EOF) {
			return resources.Project{}, errEmptyBody
		}
		if isUnknownFieldError(err) {
			return resources.Project{}, errUnknownField
		}
		return resources.Project{}, errMalformedJSON
	}
	return p, nil
}
