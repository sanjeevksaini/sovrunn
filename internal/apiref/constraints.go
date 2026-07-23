package apiref

import (
	"strings"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
)

// Direction is the allowed relationship direction for a reference field
// (F12-REF-003/004).
type Direction string

const (
	DirectionInbound       Direction = "Inbound"
	DirectionOutbound      Direction = "Outbound"
	DirectionBidirectional Direction = "Bidirectional"
)

// Valid reports whether d is one of the three approved directions.
func (d Direction) Valid() bool {
	switch d {
	case DirectionInbound, DirectionOutbound, DirectionBidirectional:
		return true
	default:
		return false
	}
}

// Constraint restricts a reference field's allowed kinds, scopes, and
// direction. Public schemas SHOULD expose domain-specific aliases built on
// this base (F12-REF-004). ValidateRef returns package-local RefIssue values;
// apiref MUST NOT import apiproblem.
type Constraint struct {
	AllowedKinds  []string
	AllowedScopes []apimeta.ScopeKind
	Direction     Direction // Inbound, Outbound, or Bidirectional
}

// RefIssue is a package-local diagnostic for a reference constraint failure
// (field path + stable code + message). Translation to apiproblem.Violation
// is owned by apivalid.
type RefIssue struct {
	Path    string // RFC 6901 JSON Pointer
	Code    string // stable machine-readable code
	Message string // human-readable; must not carry secrets or provider detail
}

// ValidateRef enforces allowed kinds/scopes/direction, required TypedRef
// fields, and provider-native-id rejection (F12-REF-001/003). Name/uid
// agreement against a resolved identity is CheckNameUIDAgreement
// (F12-REF-002); this method stays store-free for offline structural use.
func (c Constraint) ValidateRef(ref TypedRef, path string) []RefIssue {
	var issues []RefIssue

	if c.Direction != "" && !c.Direction.Valid() {
		issues = append(issues, RefIssue{
			Path:    path,
			Code:    CodeDirectionInvalid,
			Message: "reference direction must be Inbound, Outbound, or Bidirectional",
		})
	}

	if ref.APIVersion == "" {
		issues = append(issues, RefIssue{
			Path:    joinPointer(path, "apiVersion"),
			Code:    CodeMissingAPIVersion,
			Message: "reference apiVersion is required",
		})
	}
	if ref.Kind == "" {
		issues = append(issues, RefIssue{
			Path:    joinPointer(path, "kind"),
			Code:    CodeMissingKind,
			Message: "reference kind is required",
		})
	}
	if ref.Name == "" {
		issues = append(issues, RefIssue{
			Path:    joinPointer(path, "name"),
			Code:    CodeMissingName,
			Message: "reference name is required",
		})
	}

	if ref.Kind != "" && len(c.AllowedKinds) > 0 && !containsString(c.AllowedKinds, ref.Kind) {
		issues = append(issues, RefIssue{
			Path:    joinPointer(path, "kind"),
			Code:    CodeKindNotAllowed,
			Message: "reference kind is not in the allowed set for this field",
		})
	}

	if ref.Kind != "" && len(c.AllowedScopes) > 0 {
		sk := apimeta.ScopeKind(ref.Kind)
		if !containsScope(c.AllowedScopes, sk) {
			issues = append(issues, RefIssue{
				Path:    joinPointer(path, "kind"),
				Code:    CodeScopeNotAllowed,
				Message: "reference scope kind is not in the allowed set for this field",
			})
		}
	}

	if looksProviderNative(ref) {
		issues = append(issues, RefIssue{
			Path:    path,
			Code:    CodeProviderNativeID,
			Message: "provider-native identifiers must not act as core references",
		})
	}

	return issues
}

func joinPointer(base, field string) string {
	if base == "" {
		return "/" + field
	}
	if base == "/" {
		return "/" + field
	}
	return base + "/" + field
}

func containsString(set []string, v string) bool {
	for _, s := range set {
		if s == v {
			return true
		}
	}
	return false
}

func containsScope(set []apimeta.ScopeKind, v apimeta.ScopeKind) bool {
	for _, s := range set {
		if s == v {
			return true
		}
	}
	return false
}

// looksProviderNative reports whether the ref carries a provider-native
// identifier shape that must not act as a core TypedRef (F12-REF-003).
// Patterns cover common cloud resource ID forms; messages stay generic.
func looksProviderNative(ref TypedRef) bool {
	if providerNativeKind(ref.Kind) {
		return true
	}
	if providerNativeID(ref.Name) {
		return true
	}
	if providerNativeID(ref.UID) {
		return true
	}
	return false
}

func providerNativeKind(kind string) bool {
	if kind == "" {
		return false
	}
	// CloudFormation / ARM-style type names are never Sovrunn core kinds.
	if strings.Contains(kind, "::") {
		return true
	}
	lower := strings.ToLower(kind)
	if strings.HasPrefix(lower, "microsoft.") {
		return true
	}
	if strings.HasPrefix(lower, "google.") || strings.HasPrefix(lower, "compute.googleapis.com/") {
		return true
	}
	return false
}

func providerNativeID(s string) bool {
	if s == "" {
		return false
	}
	lower := strings.ToLower(s)
	switch {
	case strings.HasPrefix(lower, "arn:aws:"),
		strings.HasPrefix(lower, "arn:aws-cn:"),
		strings.HasPrefix(lower, "arn:aws-us-gov:"):
		return true
	case strings.Contains(lower, "/subscriptions/") && strings.Contains(lower, "/providers/"):
		return true
	case strings.HasPrefix(lower, "projects/") && strings.Contains(lower, "/locations/"):
		return true
	case strings.HasPrefix(lower, "oci."):
		return true
	default:
		return false
	}
}
