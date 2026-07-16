package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/sanjeevksaini/sovrunn/internal/resources"
)

// safeDecodeOrganizationUnit applies http.MaxBytesReader, detects whether
// the JSON request body contains the key "status", then decodes into the
// typed OrganizationUnit struct using DisallowUnknownFields. It reuses the
// sentinel errors and unknown-field detection defined in decode.go so error
// mapping stays consistent with safeDecodeOrganization.
func safeDecodeOrganizationUnit(w http.ResponseWriter, r *http.Request) (resources.OrganizationUnit, error) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		var maxErr *http.MaxBytesError
		if errors.As(err, &maxErr) {
			return resources.OrganizationUnit{}, errBodyTooLarge
		}
		return resources.OrganizationUnit{}, errMalformedJSON
	}

	if len(body) == 0 {
		return resources.OrganizationUnit{}, errEmptyBody
	}

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(body, &raw); err != nil {
		return resources.OrganizationUnit{}, errMalformedJSON
	}
	if _, ok := raw["status"]; ok {
		return resources.OrganizationUnit{}, errStatusFieldPresent
	}

	dec := json.NewDecoder(bytes.NewReader(body))
	dec.DisallowUnknownFields()
	var ou resources.OrganizationUnit
	if err := dec.Decode(&ou); err != nil {
		var syntaxErr *json.SyntaxError
		var typeErr *json.UnmarshalTypeError
		if errors.As(err, &syntaxErr) || errors.As(err, &typeErr) {
			return resources.OrganizationUnit{}, errMalformedJSON
		}
		if errors.Is(err, io.EOF) {
			return resources.OrganizationUnit{}, errEmptyBody
		}
		if isUnknownFieldError(err) {
			return resources.OrganizationUnit{}, errUnknownField
		}
		return resources.OrganizationUnit{}, errMalformedJSON
	}
	return ou, nil
}
