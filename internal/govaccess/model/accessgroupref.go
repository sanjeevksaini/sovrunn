package model

import "github.com/sanjeevksaini/sovrunn/internal/apimeta"

// Canonical AccessGroup identity values (VS0-SCHEMA-062).
const (
	APIVersionIAM      = "iam.sovrunn.io/v1alpha1"
	KindAccessGroup    = "AccessGroup"
	APIVersionGov      = "governance.sovrunn.io/v1alpha1"
	KindExceptionGrant = "ExceptionGrant"
)

// AccessGroupRef is a UID-pinned canonical reference to exactly one AccessGroup
// (REQ-F18-04; DD-03). It embeds apimeta.TypedRef pinned to Kind==AccessGroup
// and is not interchangeable with PrincipalRef.
type AccessGroupRef struct {
	apimeta.TypedRef
}

// Equal reports exact TypedRef field equality.
func (r AccessGroupRef) Equal(other AccessGroupRef) bool {
	return r.APIVersion == other.APIVersion &&
		r.Kind == other.Kind &&
		r.Name == other.Name &&
		r.UID == other.UID
}
