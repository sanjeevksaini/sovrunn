package registry

import "context"

// TenantChildBlocker is injected into the Tenant delete path in a later task.
// It reports child resources that block deleting a specific Tenant.
type TenantChildBlocker interface {
	BlockedByTenantChildren(ctx context.Context, orgName, ouName, tenantName string) ([]BlockedBy, error)
}

// ProjectChildBlockerChecker implements TenantChildBlocker for Project
// resources. It queries the ProjectRegistry to determine whether any Projects
// reference the Tenant being deleted.
type ProjectChildBlockerChecker struct {
	projectRegistry ProjectRegistryIface
}

// NewProjectChildBlockerChecker constructs the Project blocker.
func NewProjectChildBlockerChecker(reg ProjectRegistryIface) *ProjectChildBlockerChecker {
	return &ProjectChildBlockerChecker{projectRegistry: reg}
}

// BlockedByTenantChildren returns a BlockedBy entry identifying Project as the
// blocking resource kind when one or more Projects reference the parent
// Tenant. It returns an empty result when none reference it.
func (c *ProjectChildBlockerChecker) BlockedByTenantChildren(
	ctx context.Context, orgName, ouName, tenantName string,
) ([]BlockedBy, error) {
	count, err := c.projectRegistry.CountByTenant(ctx, orgName, ouName, tenantName)
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return []BlockedBy{{Kind: "Project", Count: count}}, nil
	}
	return nil, nil
}
