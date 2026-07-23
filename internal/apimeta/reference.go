package apimeta

// TypedRef is the common typed-reference base (F12-REF-001): apiVersion,
// kind, name, and optional immutable uid. It is stdlib-only and lives in
// apimeta so both scope/owner references and apiref aliases embed it without
// an import cycle (D-16).
type TypedRef struct {
	APIVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
	Name       string `json:"name"`
	UID        string `json:"uid,omitempty"` // optional immutable; must agree with name
}
