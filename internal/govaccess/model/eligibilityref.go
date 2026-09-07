package model

// EligibilityRef is an exact-Human eligibility entry containing exactly one
// Human PrincipalRef (REQ-F18-12/14/16; SEC-06; DD-03). Roles, assignments,
// groups, Membership, and claims never populate or satisfy it.
type EligibilityRef struct {
	Principal PrincipalRef `json:"principal"`
}
