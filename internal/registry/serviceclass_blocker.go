package registry

import "context"

// ServiceClassChildBlocker is injected into the ServiceClass delete path.
// It reports child resources that block deleting a specific ServiceClass.
type ServiceClassChildBlocker interface {
	BlockedByServiceClassChildren(ctx context.Context, serviceClassName string) ([]BlockedBy, error)
}

// ServicePlanChildBlockerChecker implements ServiceClassChildBlocker for
// ServicePlan resources. It queries the ServicePlanRegistry to determine
// whether any ServicePlans reference the ServiceClass being deleted.
type ServicePlanChildBlockerChecker struct {
	servicePlanRegistry ServicePlanRegistryIface
}

// NewServicePlanChildBlockerChecker constructs the ServicePlan blocker.
func NewServicePlanChildBlockerChecker(reg ServicePlanRegistryIface) *ServicePlanChildBlockerChecker {
	return &ServicePlanChildBlockerChecker{servicePlanRegistry: reg}
}

// BlockedByServiceClassChildren returns a BlockedBy entry identifying
// ServicePlan as the blocking resource kind when one or more ServicePlans
// reference the parent ServiceClass. It returns nil when none reference it.
// Lifecycle state (including Retired) does not exempt a ServicePlan from
// blocking its parent's deletion.
func (c *ServicePlanChildBlockerChecker) BlockedByServiceClassChildren(
	ctx context.Context, serviceClassName string,
) ([]BlockedBy, error) {
	count, err := c.servicePlanRegistry.CountByServiceClass(ctx, serviceClassName)
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return []BlockedBy{{Kind: "ServicePlan", Count: count}}, nil
	}
	return nil, nil
}
