package model

// PrincipalType is the closed principal category discriminator (VS0-SCHEMA-020;
// REQ-F18-03; DD-03).
type PrincipalType string

const (
	PrincipalTypeHuman    PrincipalType = "Human"
	PrincipalTypeWorkload PrincipalType = "Workload"
	PrincipalTypeSystem   PrincipalType = "System"
)

// AllPrincipalTypes returns the closed PrincipalType set in stable order.
func AllPrincipalTypes() []PrincipalType {
	return []PrincipalType{
		PrincipalTypeHuman,
		PrincipalTypeWorkload,
		PrincipalTypeSystem,
	}
}

// Valid reports whether t is a closed PrincipalType value.
func (t PrincipalType) Valid() bool {
	switch t {
	case PrincipalTypeHuman, PrincipalTypeWorkload, PrincipalTypeSystem:
		return true
	default:
		return false
	}
}

// PrincipalRef is the FEATURE-0018-owned already-authenticated issuer/subject
// identity value (VS0-SCHEMA-020; REQ-F18-03; DD-03). It is not a TypedRef and
// carries no apiVersion/kind/name/uid resource-reference grammar. Email is
// never an identity field.
type PrincipalRef struct {
	Issuer        string        `json:"issuer"`
	Subject       string        `json:"subject"`
	PrincipalType PrincipalType `json:"principalType"`
}

// Equal reports exact issuer/subject/type identity equality.
func (p PrincipalRef) Equal(other PrincipalRef) bool {
	return p.Issuer == other.Issuer &&
		p.Subject == other.Subject &&
		p.PrincipalType == other.PrincipalType
}
