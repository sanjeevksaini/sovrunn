package projection

import "fmt"

// Journey identifies one approved progressive-disclosure user flow (REQ-F18-24).
type Journey string

const (
	JourneyRoutineOperation  Journey = "routine-operation"
	JourneyAccessAssignment  Journey = "access-assignment"
	JourneyPrivilegedRequest Journey = "privileged-request"
	JourneyApproverDecision  Journey = "approver-decision"
	JourneyAccessReview      Journey = "access-review"
	JourneyExceptionRequest  Journey = "exception-request"
)

// View is the minimal end-customer field projection for one journey.
type View struct {
	Journey Journey
	Risk    string
	Fields  []string
}

var journeyCatalog = map[Journey]View{
	JourneyRoutineOperation: {
		Journey: JourneyRoutineOperation,
		Risk:    "low",
		Fields:  []string{"action", "target", "result", "status"},
	},
	JourneyAccessAssignment: {
		Journey: JourneyAccessAssignment,
		Risk:    "medium",
		Fields:  []string{"assignee", "role", "scope", "status", "expiresAt"},
	},
	JourneyPrivilegedRequest: {
		Journey: JourneyPrivilegedRequest,
		Risk:    "high",
		Fields:  []string{"requestId", "requestedRole", "justification", "duration", "status", "expiresAt"},
	},
	JourneyApproverDecision: {
		Journey: JourneyApproverDecision,
		Risk:    "high",
		Fields:  []string{"requestId", "subject", "decision", "reason", "status"},
	},
	JourneyAccessReview: {
		Journey: JourneyAccessReview,
		Risk:    "high",
		Fields:  []string{"reviewId", "holder", "role", "scope", "recommendation", "status"},
	},
	JourneyExceptionRequest: {
		Journey: JourneyExceptionRequest,
		Risk:    "high",
		Fields:  []string{"requestId", "control", "justification", "effectiveFrom", "effectiveUntil", "status"},
	},
}

// AllJourneys returns the closed six-journey set in stable order.
func AllJourneys() []Journey {
	return []Journey{
		JourneyRoutineOperation,
		JourneyAccessAssignment,
		JourneyPrivilegedRequest,
		JourneyApproverDecision,
		JourneyAccessReview,
		JourneyExceptionRequest,
	}
}

// ViewForJourney returns the approved field projection for one journey.
func ViewForJourney(j Journey) (View, bool) {
	v, ok := journeyCatalog[j]
	if !ok {
		return View{}, false
	}
	return v, true
}

// Project produces the end-customer view by redacting sensitive fields and
// allowing only the journey's approved minimal field set.
func Project(j Journey, input map[string]string) (map[string]string, error) {
	view, ok := ViewForJourney(j)
	if !ok {
		return nil, fmt.Errorf("unknown journey: %s", j)
	}
	return whitelist(Redact(input), view.Fields), nil
}

func whitelist(in map[string]string, allow []string) map[string]string {
	allowed := make(map[string]struct{}, len(allow))
	for _, key := range allow {
		allowed[key] = struct{}{}
	}
	out := make(map[string]string, len(allow))
	for k, v := range in {
		if _, ok := allowed[k]; ok {
			out[k] = v
		}
	}
	return out
}
