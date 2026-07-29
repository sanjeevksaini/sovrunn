package validate

import (
	"encoding/json"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/decision"
)

// RFC 6901 pointers for the obligation validation pass (design §9.1 step 7).
// Pointers are rooted at the obligations array; orchestration (T-026) may
// prefix them with /record/result when validating an embedded DecisionRecord.
const (
	ptrObligations = "/obligations"

	// Schema-aligned structural identifier ceiling
	// (api/schemas/decision-record.json obligations.items.properties.id).
	obligationIDMaxLen = 253
)

// ValidateObligations runs the FEATURE-0013 obligation validation pass
// (design §9.1 step 7, §9.2; F13-OBLIG-001/002/003; AD-032, AD-043).
//
// supported is the enforcement-point capability map keyed by obligation
// vocabulary identity:
//   - key absent → unknown vocabulary identity
//   - key present, value false → known but unsupported
//   - key present, value true → known and supported
//
// Checks (fail closed, short-circuit on first blocking violation), per
// obligation in index order:
//  1. structural id/payload → DECISION_OBLIGATION_INVALID
//  2. vocabulary membership → DECISION_OBLIGATION_UNKNOWN
//  3. mandatory support → DECISION_OBLIGATION_UNSUPPORTED_MANDATORY
//
// An unsupported mandatory obligation makes an ALLOWED decision
// unenforceable (F13-OBLIG-002; edge 14). Optional unsupported obligations
// are skipped. An empty obligations list is valid (design §7.6).
func ValidateObligations(obs []decision.Obligation, supported map[string]bool) *apiproblem.Problem {
	for i, ob := range obs {
		if prob := checkObligationStructure(i, ob); prob != nil {
			return prob
		}

		id := strings.TrimSpace(ob.ID)
		ok, known := supported[id]
		if !known {
			return obligationProblem(CodeObligationUnknown, obligationPtr(i, "id"),
				"obligation vocabulary identity is unknown")
		}
		if ob.Mandatory && !ok {
			// F13-OBLIG-002/003; edge 14: unsupported mandatory → allow unenforceable.
			return obligationProblem(CodeObligationUnsupportedMandatory, obligationPtr(i, "id"),
				"mandatory obligation is not supported by the enforcement point; ALLOWED decision is unenforceable")
		}
	}
	return nil
}

// checkObligationStructure validates id requiredness/ceilings and payload
// structural well-formedness (F13-OBLIG-003). Payload is optional; when
// present it must be non-empty valid JSON. No schema-content evaluation is
// performed (profile-declared opaque payload).
func checkObligationStructure(i int, ob decision.Obligation) *apiproblem.Problem {
	id := strings.TrimSpace(ob.ID)
	if id == "" {
		return obligationProblem(CodeObligationInvalid, obligationPtr(i, "id"),
			"obligation id is required and must be non-empty")
	}
	if utf8.RuneCountInString(id) > obligationIDMaxLen {
		return obligationProblem(CodeObligationInvalid, obligationPtr(i, "id"),
			"obligation id exceeds structural identifier length")
	}
	// Reject whitespace-padded identity that would otherwise diverge from the
	// trimmed vocabulary key used for supported-map lookup.
	if id != ob.ID {
		return obligationProblem(CodeObligationInvalid, obligationPtr(i, "id"),
			"obligation id must not carry leading or trailing whitespace")
	}

	if len(ob.Payload) == 0 {
		return nil
	}
	if !json.Valid(ob.Payload) {
		return obligationProblem(CodeObligationInvalid, obligationPtr(i, "payload"),
			"obligation payload must be valid JSON when present")
	}
	return nil
}

func obligationPtr(index int, field string) string {
	base := ptrObligations + "/" + strconv.Itoa(index)
	if field == "" {
		return base
	}
	return base + "/" + field
}

func obligationProblem(code apiproblem.ViolationCode, field, message string) *apiproblem.Problem {
	return apiproblem.New(apiproblem.CodeValidationFailed).
		WithDetail(message).
		WithViolations([]apiproblem.Violation{{
			Field:   field,
			Code:    code,
			Message: message,
		}})
}
