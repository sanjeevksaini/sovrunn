package authzeval

import "github.com/sanjeevksaini/sovrunn/internal/govaccess/model"

// PrivilegedRoleEligibleGrant reports whether a privileged grant keeps the
// exact direct-principal + time-bound constraint.
func PrivilegedRoleEligibleGrant(spec model.RoleAssignmentSpec) bool {
	return spec.RoleHolderRef.IsPrincipal() && spec.Validity == model.AssignmentValidityTimeBound
}
