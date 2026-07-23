package apicond

import (
	"time"
)

// ConditionStatus is the closed status vocabulary for a Condition
// (F12-STATUS-003): True, False, or Unknown.
type ConditionStatus string

const (
	ConditionTrue    ConditionStatus = "True"
	ConditionFalse   ConditionStatus = "False"
	ConditionUnknown ConditionStatus = "Unknown"
)

// Valid reports whether s is one of True, False, or Unknown.
func (s ConditionStatus) Valid() bool {
	switch s {
	case ConditionTrue, ConditionFalse, ConditionUnknown:
		return true
	default:
		return false
	}
}

// Condition is a stable, machine-readable current-fact observation.
// It is NOT event history (F12-STATUS-002/004).
//
// Ownership: status producers own Condition values; clients must not treat
// message as a machine contract. type and reason are stable PascalCase
// identifiers; message is human-readable and informational only.
type Condition struct {
	Type               string          `json:"type"`
	Status             ConditionStatus `json:"status"`
	Reason             string          `json:"reason"`
	Message            string          `json:"message,omitempty"`
	ObservedGeneration int64           `json:"observedGeneration"`
	LastTransitionTime string          `json:"lastTransitionTime"`
}

// Valid reports whether c satisfies F12-STATUS-003 grammar rules:
// status is True/False/Unknown, and type/reason are non-empty PascalCase
// machine identifiers. Message is informational and is not validated.
func (c Condition) Valid() bool {
	return c.Status.Valid() && IsPascalCase(c.Type) && IsPascalCase(c.Reason)
}

// IsPascalCase reports whether s is a non-empty stable PascalCase machine
// identifier: it starts with an uppercase ASCII letter and continues with
// ASCII letters or digits only (no separators).
func IsPascalCase(s string) bool {
	if s == "" {
		return false
	}
	for i, r := range s {
		if i == 0 {
			if r < 'A' || r > 'Z' {
				return false
			}
			continue
		}
		if (r < 'A' || r > 'Z') && (r < 'a' || r > 'z') && (r < '0' || r > '9') {
			return false
		}
	}
	return true
}

// SetCondition upserts cond into conds by Type. LastTransitionTime advances
// to now (UTC RFC3339) if and only if the condition is new or its Status
// changed; otherwise the previous LastTransitionTime is preserved
// (F12-STATUS-003). Other fields (reason, message, observedGeneration) are
// always replaced by cond's values. Conditions remain current facts only:
// each type appears at most once (F12-STATUS-002).
//
// The returned slice is a new slice; the input slice is not mutated.
func SetCondition(conds []Condition, cond Condition, now time.Time) []Condition {
	nowStr := now.UTC().Format(time.RFC3339)
	for i := range conds {
		if conds[i].Type != cond.Type {
			continue
		}
		if conds[i].Status == cond.Status {
			cond.LastTransitionTime = conds[i].LastTransitionTime
		} else {
			cond.LastTransitionTime = nowStr
		}
		out := make([]Condition, len(conds))
		copy(out, conds)
		out[i] = cond
		return out
	}

	cond.LastTransitionTime = nowStr
	out := make([]Condition, len(conds)+1)
	copy(out, conds)
	out[len(conds)] = cond
	return out
}

// GetCondition returns the condition with the given type, if present.
// Conditions are keyed by type; at most one entry exists per type.
func GetCondition(conds []Condition, condType string) (Condition, bool) {
	for _, c := range conds {
		if c.Type == condType {
			return c, true
		}
	}
	return Condition{}, false
}
