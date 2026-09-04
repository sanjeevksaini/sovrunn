package model

// RoleHolderKind is the closed RoleHolderRef discriminator (REQ-F18-04; DD-03).
type RoleHolderKind string

const (
	RoleHolderPrincipal   RoleHolderKind = "Principal"
	RoleHolderAccessGroup RoleHolderKind = "AccessGroup"
)

// AllRoleHolderKinds returns the closed RoleHolderKind set in stable order.
func AllRoleHolderKinds() []RoleHolderKind {
	return []RoleHolderKind{RoleHolderPrincipal, RoleHolderAccessGroup}
}

// Valid reports whether k is a closed RoleHolderKind value.
func (k RoleHolderKind) Valid() bool {
	switch k {
	case RoleHolderPrincipal, RoleHolderAccessGroup:
		return true
	default:
		return false
	}
}

// RoleHolderRef is a closed discriminated union of exactly one PrincipalRef or
// one AccessGroupRef (REQ-F18-04/08; DD-03). Mixed, empty, or unknown holders
// fail closed at validation.
type RoleHolderRef struct {
	Kind        RoleHolderKind  `json:"kind"`
	Principal   *PrincipalRef   `json:"principal,omitempty"`
	AccessGroup *AccessGroupRef `json:"accessGroup,omitempty"`
}

// IsPrincipal reports whether the Principal variant is the sole active variant.
func (r RoleHolderRef) IsPrincipal() bool {
	return r.Kind == RoleHolderPrincipal && r.Principal != nil && r.AccessGroup == nil
}

// IsAccessGroup reports whether the AccessGroup variant is the sole active variant.
func (r RoleHolderRef) IsAccessGroup() bool {
	return r.Kind == RoleHolderAccessGroup && r.AccessGroup != nil && r.Principal == nil
}
