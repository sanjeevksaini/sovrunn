package apimeta

// ScopeKind is the Matrix B vocabulary — the only valid scopeRef.kind values
// (F12-SCOPE-002, D-17). Exactly six values are permitted.
type ScopeKind string

const (
	ScopePlatform         ScopeKind = "Platform"
	ScopeOrganization     ScopeKind = "Organization"
	ScopeOrganizationUnit ScopeKind = "OrganizationUnit"
	ScopeTenant           ScopeKind = "Tenant"
	ScopeProject          ScopeKind = "Project"
	ScopeProvider         ScopeKind = "Provider"
)

// AllScopeKinds returns the closed Matrix B scope-kind set in stable order.
func AllScopeKinds() []ScopeKind {
	return []ScopeKind{
		ScopePlatform,
		ScopeOrganization,
		ScopeOrganizationUnit,
		ScopeTenant,
		ScopeProject,
		ScopeProvider,
	}
}

// Valid reports whether k is one of the six Matrix B scope kinds.
func (k ScopeKind) Valid() bool {
	switch k {
	case ScopePlatform,
		ScopeOrganization,
		ScopeOrganizationUnit,
		ScopeTenant,
		ScopeProject,
		ScopeProvider:
		return true
	default:
		return false
	}
}

// ScopeRef is the immutable primary security/governance ownership reference
// (F12-META-001, F12-SCOPE-002, F12-REF-001). It is not a location and not a
// lifecycle-containment reference.
//
// ScopeRef conforms to the common typed-reference contract by carrying
// apiVersion, kind, name, and optional immutable uid through the shared
// TypedRef base. Kind is additionally constrained to Matrix B ScopeKind
// values by validation; authorization resolves by uid, not name.
//
// Canonical platform scope (F12-SCOPE-002, D-16): the single canonical
// stored and emitted form of platform scope is an absent (nil) *ScopeRef.
// An explicit ScopeRef with Kind == "Platform" is an accepted input
// alternate only; NormalizeScope maps it to nil during layer-5 defaulting
// before identity, authorization, concurrency, persistence, and output
// processing. A nil/normalized platform scope is valid only when the
// schema's x-sovrunn-allowed-scopes includes "Platform".
type ScopeRef struct {
	// Anonymous embedding: encoding/json promotes the embedded fields
	// (apiVersion, kind, name, uid) with no tag.
	TypedRef
}

// OwnerRef expresses resource-local lifecycle containment only
// (F12-OWNER-001). It MUST NOT be used as a scopeRef.kind or as a
// security/governance scope. Like all references it uses the shared
// typed-reference base.
type OwnerRef struct {
	TypedRef // apiVersion, kind, name, optional uid; JSON promotes embedded fields
}

// PlatformScopeUID is the reserved platform-scope identity sentinel used by
// the "API group + kind + scope UID + name" uniqueness rule for
// platform-scoped resources (D-16). It is a fixed constant that can never
// collide with a generated uid (see IsGeneratedUIDFormat).
const PlatformScopeUID = "platform"

// NormalizeScope returns the canonical scope form: it maps an explicit
// Kind=="Platform" ScopeRef to nil and leaves all other scopes unchanged.
// It runs during layer-5 defaulting so identity, authorization,
// concurrency, persistence, and output all operate on the canonical
// representation (D-16).
func NormalizeScope(s *ScopeRef) *ScopeRef {
	if s == nil {
		return nil
	}
	if ScopeKind(s.Kind) == ScopePlatform {
		return nil
	}
	return s
}

// ScopeIdentity is a canonical value representation of a governance scope,
// usable for authorization comparison without requiring a full ScopeRef
// pointer (D-16, D-17). It avoids the nil-vs-non-nil ambiguity of *ScopeRef
// for platform scope and enables direct equality comparison.
type ScopeIdentity struct {
	Kind ScopeKind
	UID  string // PlatformScopeUID for platform; target scope UID otherwise
}

// CanonicalScopeIdentity converts a *ScopeRef to a ScopeIdentity:
//   - nil scopeRef (canonical platform) -> {ScopePlatform, PlatformScopeUID}
//   - explicit Kind==Platform (pre-normalization alternate) -> same platform identity
//   - non-platform scopeRef -> {ref.Kind, ref.UID}
func CanonicalScopeIdentity(s *ScopeRef) ScopeIdentity {
	if s == nil || ScopeKind(s.Kind) == ScopePlatform {
		return ScopeIdentity{Kind: ScopePlatform, UID: PlatformScopeUID}
	}
	return ScopeIdentity{Kind: ScopeKind(s.Kind), UID: s.UID}
}
