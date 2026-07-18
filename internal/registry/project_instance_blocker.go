package registry

import "context"

// ProjectInstanceBlocker reports ServiceInstances that prevent deleting a
// specific Project.
type ProjectInstanceBlocker interface {
	BlockedByProjectInstances(
		ctx context.Context,
		orgName, ouName, tenantName, projectName string,
	) ([]BlockedBy, error)
}

// ServiceInstanceProjectBlockerChecker implements ProjectInstanceBlocker by
// counting ServiceInstances under the Project's four-part governance identity.
type ServiceInstanceProjectBlockerChecker struct {
	serviceInstanceRegistry ServiceInstanceRegistryIface
}

// NewServiceInstanceProjectBlockerChecker constructs the Project instance
// blocker.
func NewServiceInstanceProjectBlockerChecker(
	reg ServiceInstanceRegistryIface,
) *ServiceInstanceProjectBlockerChecker {
	return &ServiceInstanceProjectBlockerChecker{serviceInstanceRegistry: reg}
}

// BlockedByProjectInstances returns a BlockedBy entry when one or more
// ServiceInstances exist under the requested Project. It returns nil when the
// project has no ServiceInstances.
func (c *ServiceInstanceProjectBlockerChecker) BlockedByProjectInstances(
	ctx context.Context,
	orgName, ouName, tenantName, projectName string,
) ([]BlockedBy, error) {
	count, err := c.serviceInstanceRegistry.CountByProject(
		ctx, orgName, ouName, tenantName, projectName,
	)
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return []BlockedBy{{Kind: "ServiceInstance", Count: count}}, nil
	}
	return nil, nil
}
