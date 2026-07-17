package registry

import "context"

// OUChildBlocker is injected into the OrganizationUnit delete path in a later
// task. It reports child resources that block deleting a specific
// OrganizationUnit.
type OUChildBlocker interface {
	BlockedByOUChildren(ctx context.Context, orgName, ouName string) ([]BlockedBy, error)
}

// TenantChildBlockerChecker implements OUChildBlocker for Tenant resources.
// It queries the TenantRegistry to determine whether any Tenants reference the
// OrganizationUnit being deleted.
type TenantChildBlockerChecker struct {
	tenantRegistry TenantRegistryIface
}

// NewTenantChildBlockerChecker constructs the Tenant blocker.
func NewTenantChildBlockerChecker(reg TenantRegistryIface) *TenantChildBlockerChecker {
	return &TenantChildBlockerChecker{tenantRegistry: reg}
}

// BlockedByOUChildren returns a BlockedBy entry identifying Tenant as the
// blocking resource kind when one or more Tenants reference the parent
// OrganizationUnit. It returns an empty result when none reference it.
func (c *TenantChildBlockerChecker) BlockedByOUChildren(
	ctx context.Context, orgName, ouName string,
) ([]BlockedBy, error) {
	count, err := c.tenantRegistry.CountByOrganizationUnit(ctx, orgName, ouName)
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return []BlockedBy{{Kind: "Tenant", Count: count}}, nil
	}
	return nil, nil
}
