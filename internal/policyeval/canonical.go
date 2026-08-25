package policyeval

import (
	"crypto/sha256"
	"encoding/hex"
	"sort"
	"unicode/utf8"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
)

// Exact v1 schema separator embedded in the normalized logical object
// (CDG-F17-02). No external digest prefix or suffix is used.
const policyEvaluationRequestSchemaV1 = "sovrunn.policy-evaluation-request/v1"

// digestPolicyEvaluationRequest returns the lowercase 64-character hex
// SHA-256 of the RFC 8785 JCS bytes of the v1 normalized logical object.
// requestId is excluded. Invalid UTF-8 or non-scalar Unicode in any
// serialized string yields errRequestInvalid and no digest.
func digestPolicyEvaluationRequest(req PolicyEvaluationRequest) (string, error) {
	canonical, err := canonicalPolicyEvaluationRequestBytes(req)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(canonical)
	return hex.EncodeToString(sum[:]), nil
}

// canonicalPolicyEvaluationRequestBytes deep-copies req, sorts reference
// lists by normalized identity, and emits the domain-bounded RFC 8785 JCS
// v1 serialization with no insignificant whitespace, BOM, trailing newline,
// or external prefix/suffix.
func canonicalPolicyEvaluationRequestBytes(req PolicyEvaluationRequest) ([]byte, error) {
	snap := snapshotPolicyEvaluationRequest(req)
	sortTypedRefs(snap.ProfileRefs)
	sortTypedRefs(snap.CandidateRefs)

	buf := make([]byte, 0, 256)
	buf = append(buf, '{')

	var err error
	buf, err = appendJCSMember(buf, "action", snap.Action)
	if err != nil {
		return nil, err
	}
	buf = append(buf, ',')

	buf, err = appendJCSRefArrayMember(buf, "candidateRefs", snap.CandidateRefs)
	if err != nil {
		return nil, err
	}

	if snap.ContextRef != nil {
		buf = append(buf, ',')
		buf, err = appendJCSRefMember(buf, "contextRef", *snap.ContextRef)
		if err != nil {
			return nil, err
		}
	}

	buf = append(buf, ',')
	buf, err = appendJCSRefArrayMember(buf, "profileRefs", snap.ProfileRefs)
	if err != nil {
		return nil, err
	}

	buf = append(buf, ',')
	buf, err = appendJCSMember(buf, "schema", policyEvaluationRequestSchemaV1)
	if err != nil {
		return nil, err
	}

	buf = append(buf, ',')
	buf, err = appendJCSRefMember(buf, "subjectRef", snap.SubjectRef)
	if err != nil {
		return nil, err
	}

	buf = append(buf, '}')
	return buf, nil
}

func snapshotPolicyEvaluationRequest(req PolicyEvaluationRequest) PolicyEvaluationRequest {
	out := PolicyEvaluationRequest{
		SubjectRef: req.SubjectRef,
		Action:     req.Action,
		RequestID:  req.RequestID,
	}
	if req.ContextRef != nil {
		cp := *req.ContextRef
		out.ContextRef = &cp
	}
	if req.ProfileRefs != nil {
		out.ProfileRefs = append([]apimeta.TypedRef(nil), req.ProfileRefs...)
	} else {
		out.ProfileRefs = []apimeta.TypedRef{}
	}
	if req.CandidateRefs != nil {
		out.CandidateRefs = append([]apimeta.TypedRef(nil), req.CandidateRefs...)
	} else {
		// Omitted and empty candidateRefs are the same empty set (CDG-F17-02).
		out.CandidateRefs = []apimeta.TypedRef{}
	}
	return out
}

func sortTypedRefs(refs []apimeta.TypedRef) {
	sort.SliceStable(refs, func(i, j int) bool {
		return typedRefLess(refs[i], refs[j])
	})
}

// typedRefLess compares normalized reference identity tuples
// (apiVersion, kind, name, uid-or-empty) with bytewise UTF-8 lexical order.
func typedRefLess(a, b apimeta.TypedRef) bool {
	if a.APIVersion != b.APIVersion {
		return a.APIVersion < b.APIVersion
	}
	if a.Kind != b.Kind {
		return a.Kind < b.Kind
	}
	if a.Name != b.Name {
		return a.Name < b.Name
	}
	return a.UID < b.UID
}

func appendJCSMember(buf []byte, key, value string) ([]byte, error) {
	var err error
	buf, err = appendJCSString(buf, key)
	if err != nil {
		return nil, err
	}
	buf = append(buf, ':')
	buf, err = appendJCSString(buf, value)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func appendJCSRefMember(buf []byte, key string, ref apimeta.TypedRef) ([]byte, error) {
	var err error
	buf, err = appendJCSString(buf, key)
	if err != nil {
		return nil, err
	}
	buf = append(buf, ':')
	buf, err = appendJCSRef(buf, ref)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func appendJCSRefArrayMember(buf []byte, key string, refs []apimeta.TypedRef) ([]byte, error) {
	var err error
	buf, err = appendJCSString(buf, key)
	if err != nil {
		return nil, err
	}
	buf = append(buf, ':')
	buf, err = appendJCSRefArray(buf, refs)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func appendJCSRefArray(buf []byte, refs []apimeta.TypedRef) ([]byte, error) {
	buf = append(buf, '[')
	for i, ref := range refs {
		if i > 0 {
			buf = append(buf, ',')
		}
		var err error
		buf, err = appendJCSRef(buf, ref)
		if err != nil {
			return nil, err
		}
	}
	buf = append(buf, ']')
	return buf, nil
}

// appendJCSRef emits a normalized reference object with keys in fixed JCS
// order: apiVersion, kind, name, and uid only when present (non-empty).
func appendJCSRef(buf []byte, ref apimeta.TypedRef) ([]byte, error) {
	buf = append(buf, '{')

	var err error
	buf, err = appendJCSMember(buf, "apiVersion", ref.APIVersion)
	if err != nil {
		return nil, err
	}
	buf = append(buf, ',')
	buf, err = appendJCSMember(buf, "kind", ref.Kind)
	if err != nil {
		return nil, err
	}
	buf = append(buf, ',')
	buf, err = appendJCSMember(buf, "name", ref.Name)
	if err != nil {
		return nil, err
	}
	if ref.UID != "" {
		buf = append(buf, ',')
		buf, err = appendJCSMember(buf, "uid", ref.UID)
		if err != nil {
			return nil, err
		}
	}

	buf = append(buf, '}')
	return buf, nil
}

// appendJCSString emits a JSON string per the FEATURE-0017 domain-bounded
// RFC 8785 JCS subset (D-03): reject invalid UTF-8 and non-scalar Unicode;
// no Unicode normalization; escape quote, reverse-solidus, and controls;
// emit all other scalar values as original UTF-8 (including HTML-sensitive
// characters, U+2028, U+2029, and supplementary code points).
func appendJCSString(buf []byte, s string) ([]byte, error) {
	if !utf8.ValidString(s) {
		return nil, errRequestInvalid
	}
	buf = append(buf, '"')
	for i := 0; i < len(s); {
		r, size := utf8.DecodeRuneInString(s[i:])
		if r == utf8.RuneError && size == 1 {
			return nil, errRequestInvalid
		}
		// Non-scalar Unicode (UTF-16 surrogates) must never appear.
		if r >= 0xD800 && r <= 0xDFFF {
			return nil, errRequestInvalid
		}
		switch {
		case r == '"':
			buf = append(buf, '\\', '"')
		case r == '\\':
			buf = append(buf, '\\', '\\')
		case r == '\b':
			buf = append(buf, '\\', 'b')
		case r == '\t':
			buf = append(buf, '\\', 't')
		case r == '\n':
			buf = append(buf, '\\', 'n')
		case r == '\f':
			buf = append(buf, '\\', 'f')
		case r == '\r':
			buf = append(buf, '\\', 'r')
		case r < 0x20:
			buf = append(buf, '\\', 'u', '0', '0', lowerHex(byte(r>>4)), lowerHex(byte(r&0xf)))
		default:
			buf = append(buf, s[i:i+size]...)
		}
		i += size
	}
	buf = append(buf, '"')
	return buf, nil
}

func lowerHex(v byte) byte {
	if v < 10 {
		return '0' + v
	}
	return 'a' + (v - 10)
}
