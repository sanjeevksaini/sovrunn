package registry

import "context"

// PluginChildBlocker is injected into the Plugin delete path.
// It reports child resources that block deleting a specific Plugin.
type PluginChildBlocker interface {
	BlockedByPluginChildren(ctx context.Context, pluginName string) ([]BlockedBy, error)
}

// CapabilityChildBlockerChecker implements PluginChildBlocker for
// Capability resources. It queries the CapabilityRegistry to determine
// whether any Capabilities reference the Plugin being deleted.
type CapabilityChildBlockerChecker struct {
	capabilityRegistry CapabilityRegistryIface
}

// NewCapabilityChildBlockerChecker constructs the Capability blocker.
func NewCapabilityChildBlockerChecker(reg CapabilityRegistryIface) *CapabilityChildBlockerChecker {
	return &CapabilityChildBlockerChecker{capabilityRegistry: reg}
}

// BlockedByPluginChildren returns a BlockedBy entry identifying
// Capability as the blocking resource kind when one or more Capabilities
// reference the parent Plugin. It returns nil when none reference it.
func (c *CapabilityChildBlockerChecker) BlockedByPluginChildren(
	ctx context.Context, pluginName string,
) ([]BlockedBy, error) {
	count, err := c.capabilityRegistry.CountByPlugin(ctx, pluginName)
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return []BlockedBy{{Kind: "Capability", Count: count}}, nil
	}
	return nil, nil
}
