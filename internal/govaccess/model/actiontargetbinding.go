package model

import "github.com/sanjeevksaini/sovrunn/internal/apimeta"

// ActionTargetBinding is a closed discriminated union of ExactResource,
// CreateParent, or ScopeOnly (REQ-F18-06; design §4.3). Exactly one variant
// applies; mixed/missing/empty/wildcard/unknown/ambiguous bindings fail closed.
type ActionTargetBinding struct {
	Mode            ActionTargetMode  `json:"mode"`
	Action          string            `json:"action"`
	TargetKind      string            `json:"targetKind,omitempty"`
	ParentScopeKind apimeta.ScopeKind `json:"parentScopeKind,omitempty"`
	ScopeKind       apimeta.ScopeKind `json:"scopeKind,omitempty"`
}

// IsExactResource reports whether ExactResource is the sole active variant.
func (b ActionTargetBinding) IsExactResource() bool {
	return b.Mode == ActionTargetExactResource &&
		b.TargetKind != "" &&
		b.ParentScopeKind == "" &&
		b.ScopeKind == ""
}

// IsCreateParent reports whether CreateParent is the sole active variant.
func (b ActionTargetBinding) IsCreateParent() bool {
	return b.Mode == ActionTargetCreateParent &&
		b.ParentScopeKind != "" &&
		b.TargetKind == "" &&
		b.ScopeKind == ""
}

// IsScopeOnly reports whether ScopeOnly is the sole active variant.
func (b ActionTargetBinding) IsScopeOnly() bool {
	return b.Mode == ActionTargetScopeOnly &&
		b.ScopeKind != "" &&
		b.TargetKind == "" &&
		b.ParentScopeKind == ""
}
