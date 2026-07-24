package apimeta

import "strings"

// TypeMeta identifies the contract of an externally exchanged object
// (F12-NAMING-001, F12-NAMING-002).
//
// apiVersion form: <domain>.sovrunn.io/{v1alpha1|v1beta1|v1}
// kind: singular PascalCase
type TypeMeta struct {
	APIVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
}

// Approved apiVersion maturity forms (F12-NAMING-001).
const (
	VersionV1Alpha1 = "v1alpha1"
	VersionV1Beta1  = "v1beta1"
	VersionV1       = "v1"
)

// ParseAPIVersion splits apiVersion into API group and version.
// Expected form: "<group>/<version>" (for example "fabric.sovrunn.io/v1alpha1").
// ok is false when the value is empty, lacks a single "/" separator, or has
// empty group/version segments.
func ParseAPIVersion(apiVersion string) (group, version string, ok bool) {
	if apiVersion == "" {
		return "", "", false
	}
	group, version, found := strings.Cut(apiVersion, "/")
	if !found || group == "" || version == "" {
		return "", "", false
	}
	if strings.Contains(version, "/") {
		return "", "", false
	}
	return group, version, true
}

// Group returns the API group portion of APIVersion, or "" if unparsable.
func (t TypeMeta) Group() string {
	group, _, ok := ParseAPIVersion(t.APIVersion)
	if !ok {
		return ""
	}
	return group
}

// Version returns the version portion of APIVersion, or "" if unparsable.
func (t TypeMeta) Version() string {
	_, version, ok := ParseAPIVersion(t.APIVersion)
	if !ok {
		return ""
	}
	return version
}

// IsKnownVersion reports whether version is one of the approved maturity forms
// from F12-NAMING-001 (v1alpha1, v1beta1, or v1).
func IsKnownVersion(version string) bool {
	switch version {
	case VersionV1Alpha1, VersionV1Beta1, VersionV1:
		return true
	default:
		return false
	}
}
