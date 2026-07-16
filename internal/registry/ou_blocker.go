package registry

import "context"

// OUChildBlockerChecker implements ChildBlockerChecker for OrganizationUnit.
// It queries the OrganizationUnitRegistry to determine whether any
// OrganizationUnits reference the Organization being deleted.
type OUChildBlockerChecker struct {
	ouRegistry OrganizationUnitRegistryIface
}

// NewOUChildBlockerChecker constructs the blocker.
func NewOUChildBlockerChecker(ouReg OrganizationUnitRegistryIface) *OUChildBlockerChecker {
	return &OUChildBlockerChecker{ouRegistry: ouReg}
}

// BlockedByChildren returns a BlockedBy entry identifying OrganizationUnit
// as the blocking resource kind when one or more OrganizationUnits
// reference orgName. It returns an empty result when none reference it.
func (c *OUChildBlockerChecker) BlockedByChildren(
	ctx context.Context, orgName string,
) ([]BlockedBy, error) {
	count, err := c.ouRegistry.CountByOrganization(ctx, orgName)
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return []BlockedBy{{Kind: "OrganizationUnit", Count: count}}, nil
	}
	return nil, nil
}
