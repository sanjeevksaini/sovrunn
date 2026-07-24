package apimeta

// ObjectMeta is the applicable subset of common metadata for persistent
// resources (F12-META-001). Ownership and mutability are enforced by
// validation (F12-META-002, Matrix C2), not by Go type constraints.
//
// Field ownership / mutability (F12-META-002):
//   - Name: creator-authored, immutable after create; lowercase kebab-case identity
//   - UID: Sovrunn-authored, immutable, globally unique, opaque, never reused
//   - DisplayName: owner-authored, mutable; human-readable; not identity
//   - ScopeRef: creator on create, validated, normally immutable; security/governance scope
//   - Labels: authorized owners, bounded; MUST NOT contain secrets
//   - Annotations: namespaced authorized owners, bounded; MUST NOT contain secrets
//   - Generation: system-only; advances when desired state changes
//   - ResourceVersion: system-only, opaque; advances on any stored representation change
//   - CreatedAt: system-only, UTC RFC 3339
//   - UpdatedAt: system-only, UTC RFC 3339
type ObjectMeta struct {
	Name            string            `json:"name"`
	UID             string            `json:"uid,omitempty"`
	DisplayName     string            `json:"displayName,omitempty"`
	ScopeRef        *ScopeRef         `json:"scopeRef,omitempty"`
	Labels          map[string]string `json:"labels,omitempty"`
	Annotations     map[string]string `json:"annotations,omitempty"`
	Generation      int64             `json:"generation,omitempty"`
	ResourceVersion string            `json:"resourceVersion,omitempty"`
	CreatedAt       string            `json:"createdAt,omitempty"`
	UpdatedAt       string            `json:"updatedAt,omitempty"`
}
