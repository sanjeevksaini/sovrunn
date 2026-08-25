package policyeval

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"unicode/utf8"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/apiref"
)

// Request field bounds (VS0-SCHEMA-018 / CDG-F17-01).
const (
	minActionLen     = 1
	maxActionLen     = 63
	minRequestIDLen  = 1
	maxRequestIDLen  = 128
	minProfileRefs   = 1
	maxProfileRefs   = 32
	maxCandidateRefs = 64
)

// errRequestInvalid is the sole sanitized diagnostic for request decode and
// structural validation failures. It must not carry request values.
var errRequestInvalid = errors.New("policy evaluation request invalid")

// PolicyEvaluationRequest is the FEATURE-0017 engine-facing evaluation input
// envelope (VS0-SCHEMA-018). It is a TransientRequestResult-profile value with
// no TypeMeta, ObjectMeta, identity, scope, or status.
type PolicyEvaluationRequest struct {
	SubjectRef    apimeta.TypedRef   `json:"subjectRef"`
	Action        string             `json:"action"`
	ContextRef    *apimeta.TypedRef  `json:"contextRef,omitempty"`
	ProfileRefs   []apimeta.TypedRef `json:"profileRefs"`
	CandidateRefs []apimeta.TypedRef `json:"candidateRefs,omitempty"`
	RequestID     string             `json:"requestId"`
}

// policyEvaluationRequestDTO carries every root member as json.RawMessage so
// omitted versus explicit null can be distinguished before public construction.
type policyEvaluationRequestDTO struct {
	SubjectRef    json.RawMessage `json:"subjectRef"`
	Action        json.RawMessage `json:"action"`
	ContextRef    json.RawMessage `json:"contextRef"`
	ProfileRefs   json.RawMessage `json:"profileRefs"`
	CandidateRefs json.RawMessage `json:"candidateRefs"`
	RequestID     json.RawMessage `json:"requestId"`
}

// DecodePolicyEvaluationRequestJSON is the only byte decoder for
// PolicyEvaluationRequest. It rejects invalid UTF-8, lone/malformed UTF-16
// surrogate escapes, unknown and duplicate members, trailing JSON, explicit
// null where forbidden, bound violations, and structurally invalid references.
func DecodePolicyEvaluationRequestJSON(data []byte) (PolicyEvaluationRequest, error) {
	return decodePolicyEvaluationRequestJSON(data)
}

// UnmarshalJSON decodes through the same private strict decoder as
// DecodePolicyEvaluationRequestJSON so ordinary json.Unmarshal cannot bypass
// presence, unknown-field, or duplicate-member rules.
func (r *PolicyEvaluationRequest) UnmarshalJSON(data []byte) error {
	decoded, err := decodePolicyEvaluationRequestJSON(data)
	if err != nil {
		return err
	}
	*r = decoded
	return nil
}

// Validate performs presence/bound checks and FEATURE-0012 structural reference
// validation. Direct Go construction uses the same rules: nil ContextRef means
// absence; nil or empty CandidateRefs means the approved empty set.
//
// UID pinning, EffectiveGovernanceContext kind enforcement, and duplicate
// reference detection are deferred to a later validation step.
func (r PolicyEvaluationRequest) Validate() error {
	return validatePolicyEvaluationRequest(r)
}

func decodePolicyEvaluationRequestJSON(data []byte) (PolicyEvaluationRequest, error) {
	var zero PolicyEvaluationRequest
	if len(data) == 0 {
		return zero, errRequestInvalid
	}
	if err := lexicalPreflight(data); err != nil {
		return zero, err
	}
	if err := rejectDuplicateObjectMembers(data); err != nil {
		return zero, err
	}

	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	var dto policyEvaluationRequestDTO
	if err := dec.Decode(&dto); err != nil {
		return zero, errRequestInvalid
	}
	if _, err := dec.Token(); err != io.EOF {
		return zero, errRequestInvalid
	}

	req, err := dto.toRequest()
	if err != nil {
		return zero, err
	}
	if err := validatePolicyEvaluationRequest(req); err != nil {
		return zero, err
	}
	return req, nil
}

func (dto policyEvaluationRequestDTO) toRequest() (PolicyEvaluationRequest, error) {
	var zero PolicyEvaluationRequest

	subject, err := decodeTypedRefRaw(dto.SubjectRef, true)
	if err != nil {
		return zero, err
	}
	action, err := decodeRequiredString(dto.Action)
	if err != nil {
		return zero, err
	}
	requestID, err := decodeRequiredString(dto.RequestID)
	if err != nil {
		return zero, err
	}

	var contextRef *apimeta.TypedRef
	switch {
	case dto.ContextRef == nil:
		// Omitted: valid absence.
	case isJSONNull(dto.ContextRef):
		return zero, errRequestInvalid
	default:
		ref, err := decodeTypedRefRaw(dto.ContextRef, true)
		if err != nil {
			return zero, err
		}
		contextRef = &ref
	}

	if dto.ProfileRefs == nil || isJSONNull(dto.ProfileRefs) {
		return zero, errRequestInvalid
	}
	profileRefs, err := decodeTypedRefArrayRaw(dto.ProfileRefs)
	if err != nil {
		return zero, err
	}

	var candidateRefs []apimeta.TypedRef
	switch {
	case dto.CandidateRefs == nil:
		candidateRefs = make([]apimeta.TypedRef, 0)
	case isJSONNull(dto.CandidateRefs):
		return zero, errRequestInvalid
	default:
		candidateRefs, err = decodeTypedRefArrayRaw(dto.CandidateRefs)
		if err != nil {
			return zero, err
		}
	}

	return PolicyEvaluationRequest{
		SubjectRef:    subject,
		Action:        action,
		ContextRef:    contextRef,
		ProfileRefs:   profileRefs,
		CandidateRefs: candidateRefs,
		RequestID:     requestID,
	}, nil
}

func validatePolicyEvaluationRequest(r PolicyEvaluationRequest) error {
	actionLen := utf8.RuneCountInString(r.Action)
	if actionLen < minActionLen || actionLen > maxActionLen {
		return errRequestInvalid
	}
	idLen := utf8.RuneCountInString(r.RequestID)
	if idLen < minRequestIDLen || idLen > maxRequestIDLen {
		return errRequestInvalid
	}
	if len(r.ProfileRefs) < minProfileRefs || len(r.ProfileRefs) > maxProfileRefs {
		return errRequestInvalid
	}
	if len(r.CandidateRefs) > maxCandidateRefs {
		return errRequestInvalid
	}

	if issues := (apiref.Constraint{}).ValidateRef(r.SubjectRef, "/subjectRef"); len(issues) > 0 {
		return errRequestInvalid
	}
	if r.ContextRef != nil {
		if issues := (apiref.Constraint{}).ValidateRef(*r.ContextRef, "/contextRef"); len(issues) > 0 {
			return errRequestInvalid
		}
	}
	if issues := apiref.Refs(r.ProfileRefs).Validate(apiref.Constraint{}, "/profileRefs", maxProfileRefs); len(issues) > 0 {
		return errRequestInvalid
	}
	if issues := apiref.Refs(r.CandidateRefs).Validate(apiref.Constraint{}, "/candidateRefs", maxCandidateRefs); len(issues) > 0 {
		return errRequestInvalid
	}
	return nil
}

func decodeRequiredString(raw json.RawMessage) (string, error) {
	if raw == nil || isJSONNull(raw) {
		return "", errRequestInvalid
	}
	dec := json.NewDecoder(bytes.NewReader(raw))
	tok, err := dec.Token()
	if err != nil {
		return "", errRequestInvalid
	}
	s, ok := tok.(string)
	if !ok {
		return "", errRequestInvalid
	}
	if _, err := dec.Token(); err != io.EOF {
		return "", errRequestInvalid
	}
	return s, nil
}

func decodeTypedRefRaw(raw json.RawMessage, required bool) (apimeta.TypedRef, error) {
	var zero apimeta.TypedRef
	if raw == nil || isJSONNull(raw) {
		if required {
			return zero, errRequestInvalid
		}
		return zero, nil
	}
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.DisallowUnknownFields()
	var ref apimeta.TypedRef
	if err := dec.Decode(&ref); err != nil {
		return zero, errRequestInvalid
	}
	if _, err := dec.Token(); err != io.EOF {
		return zero, errRequestInvalid
	}
	return ref, nil
}

func decodeTypedRefArrayRaw(raw json.RawMessage) ([]apimeta.TypedRef, error) {
	dec := json.NewDecoder(bytes.NewReader(raw))
	tok, err := dec.Token()
	if err != nil {
		return nil, errRequestInvalid
	}
	delim, ok := tok.(json.Delim)
	if !ok || delim != '[' {
		return nil, errRequestInvalid
	}

	out := make([]apimeta.TypedRef, 0)
	for dec.More() {
		var elem json.RawMessage
		if err := dec.Decode(&elem); err != nil {
			return nil, errRequestInvalid
		}
		ref, err := decodeTypedRefRaw(elem, true)
		if err != nil {
			return nil, err
		}
		out = append(out, ref)
	}
	tok, err = dec.Token()
	if err != nil {
		return nil, errRequestInvalid
	}
	if delim, ok = tok.(json.Delim); !ok || delim != ']' {
		return nil, errRequestInvalid
	}
	if _, err := dec.Token(); err != io.EOF {
		return nil, errRequestInvalid
	}
	return out, nil
}

func isJSONNull(raw json.RawMessage) bool {
	return bytes.Equal(bytes.TrimSpace(raw), []byte("null"))
}

// lexicalPreflight rejects invalid UTF-8 and lone/malformed/unpaired UTF-16
// surrogate escapes before encoding/json can replace them with U+FFFD.
func lexicalPreflight(data []byte) error {
	if !utf8.Valid(data) {
		return errRequestInvalid
	}
	return rejectBadSurrogateEscapes(data)
}

func rejectBadSurrogateEscapes(data []byte) error {
	inString := false
	for i := 0; i < len(data); {
		c := data[i]
		if !inString {
			if c == '"' {
				inString = true
			}
			i++
			continue
		}
		if c == '"' {
			inString = false
			i++
			continue
		}
		if c != '\\' {
			i++
			continue
		}
		if i+1 >= len(data) {
			return errRequestInvalid
		}
		switch data[i+1] {
		case '"', '\\', '/', 'b', 'f', 'n', 'r', 't':
			i += 2
		case 'u':
			code, next, ok := parseJSONUnicodeEscape(data, i)
			if !ok {
				return errRequestInvalid
			}
			switch {
			case isHighSurrogate(code):
				if next+1 >= len(data) || data[next] != '\\' || data[next+1] != 'u' {
					return errRequestInvalid
				}
				low, next2, ok := parseJSONUnicodeEscape(data, next)
				if !ok || !isLowSurrogate(low) {
					return errRequestInvalid
				}
				i = next2
			case isLowSurrogate(code):
				return errRequestInvalid
			default:
				i = next
			}
		default:
			// Invalid escape sequences fail closed in the lexical preflight.
			return errRequestInvalid
		}
	}
	if inString {
		return errRequestInvalid
	}
	return nil
}

func parseJSONUnicodeEscape(data []byte, atSlash int) (rune, int, bool) {
	// Expect \uXXXX starting at atSlash.
	if atSlash+5 >= len(data) || data[atSlash] != '\\' || data[atSlash+1] != 'u' {
		return 0, 0, false
	}
	var code rune
	for j := 0; j < 4; j++ {
		h, ok := hexNibble(data[atSlash+2+j])
		if !ok {
			return 0, 0, false
		}
		code = code<<4 | rune(h)
	}
	return code, atSlash + 6, true
}

func hexNibble(b byte) (byte, bool) {
	switch {
	case b >= '0' && b <= '9':
		return b - '0', true
	case b >= 'a' && b <= 'f':
		return b - 'a' + 10, true
	case b >= 'A' && b <= 'F':
		return b - 'A' + 10, true
	default:
		return 0, false
	}
}

func isHighSurrogate(r rune) bool {
	return r >= 0xD800 && r <= 0xDBFF
}

func isLowSurrogate(r rune) bool {
	return r >= 0xDC00 && r <= 0xDFFF
}

// rejectDuplicateObjectMembers walks JSON tokens and rejects a repeated member
// name in any object (root or nested), before DTO decoding.
func rejectDuplicateObjectMembers(data []byte) error {
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()
	if err := scanJSONValueForDuplicates(dec); err != nil {
		return err
	}
	if _, err := dec.Token(); err != io.EOF {
		// Trailing content is also rejected by the DTO decode path; treat any
		// leftover token here as invalid so the preflight stays closed.
		if err == nil {
			return errRequestInvalid
		}
		if !errors.Is(err, io.EOF) {
			return errRequestInvalid
		}
	}
	return nil
}

func scanJSONValueForDuplicates(dec *json.Decoder) error {
	tok, err := dec.Token()
	if err != nil {
		return errRequestInvalid
	}
	switch t := tok.(type) {
	case json.Delim:
		switch t {
		case '{':
			return scanJSONObjectForDuplicates(dec)
		case '[':
			return scanJSONArrayForDuplicates(dec)
		default:
			return errRequestInvalid
		}
	case bool, float64, json.Number, string, nil:
		return nil
	default:
		return errRequestInvalid
	}
}

func scanJSONObjectForDuplicates(dec *json.Decoder) error {
	seen := make(map[string]struct{})
	for dec.More() {
		keyTok, err := dec.Token()
		if err != nil {
			return errRequestInvalid
		}
		key, ok := keyTok.(string)
		if !ok {
			return errRequestInvalid
		}
		if _, dup := seen[key]; dup {
			return errRequestInvalid
		}
		seen[key] = struct{}{}
		if err := scanJSONValueForDuplicates(dec); err != nil {
			return err
		}
	}
	tok, err := dec.Token()
	if err != nil {
		return errRequestInvalid
	}
	if delim, ok := tok.(json.Delim); !ok || delim != '}' {
		return errRequestInvalid
	}
	return nil
}

func scanJSONArrayForDuplicates(dec *json.Decoder) error {
	for dec.More() {
		if err := scanJSONValueForDuplicates(dec); err != nil {
			return err
		}
	}
	tok, err := dec.Token()
	if err != nil {
		return errRequestInvalid
	}
	if delim, ok := tok.(json.Delim); !ok || delim != ']' {
		return errRequestInvalid
	}
	return nil
}
