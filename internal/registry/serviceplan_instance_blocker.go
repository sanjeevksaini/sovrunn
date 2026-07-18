package registry

import "context"

// ServicePlanInstanceBlocker reports ServiceInstances that prevent deleting a
// specific ServicePlan.
type ServicePlanInstanceBlocker interface {
	BlockedByServicePlanInstances(
		ctx context.Context,
		serviceClassName, planName string,
	) ([]BlockedBy, error)
}

// ServiceInstancePlanBlockerChecker implements ServicePlanInstanceBlocker by
// counting ServiceInstances that reference the ServicePlan's composite
// ServiceClass/name identity.
type ServiceInstancePlanBlockerChecker struct {
	serviceInstanceRegistry ServiceInstanceRegistryIface
}

// NewServiceInstancePlanBlockerChecker constructs the ServicePlan instance
// blocker.
func NewServiceInstancePlanBlockerChecker(
	reg ServiceInstanceRegistryIface,
) *ServiceInstancePlanBlockerChecker {
	return &ServiceInstancePlanBlockerChecker{serviceInstanceRegistry: reg}
}

// BlockedByServicePlanInstances returns a BlockedBy entry when one or more
// ServiceInstances reference the requested ServicePlan. It returns nil when
// the plan has no referencing ServiceInstances.
func (c *ServiceInstancePlanBlockerChecker) BlockedByServicePlanInstances(
	ctx context.Context,
	serviceClassName, planName string,
) ([]BlockedBy, error) {
	count, err := c.serviceInstanceRegistry.CountByServicePlan(ctx, serviceClassName, planName)
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return []BlockedBy{{Kind: "ServiceInstance", Count: count}}, nil
	}
	return nil, nil
}
