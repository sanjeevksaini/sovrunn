package registry

import (
	"context"

	"github.com/sanjeevksaini/sovrunn/internal/resources"
)

// CapabilityLookupImpl implements CapabilityLookup by querying a
// CapabilityRegistryIface. A Capability is considered active when
// status.phase == "Active" AND spec.supported == true.
type CapabilityLookupImpl struct {
	capRegistry CapabilityRegistryIface
}

// NewCapabilityLookup constructs a CapabilityLookupImpl backed by reg.
func NewCapabilityLookup(reg CapabilityRegistryIface) *CapabilityLookupImpl {
	return &CapabilityLookupImpl{capRegistry: reg}
}

// Compile-time check that CapabilityLookupImpl satisfies CapabilityLookup.
var _ CapabilityLookup = (*CapabilityLookupImpl)(nil)

// HasActiveCapabilityForServiceClass reports whether at least one Capability
// for serviceClassRef is both Active and Supported. It returns false when
// none match; registry errors are propagated unchanged.
func (l *CapabilityLookupImpl) HasActiveCapabilityForServiceClass(
	ctx context.Context, serviceClassRef string,
) (bool, error) {
	caps, err := l.capRegistry.ListCapabilities(ctx, "", serviceClassRef)
	if err != nil {
		return false, err
	}
	for _, c := range caps {
		if c.Status.Phase == resources.PhaseActive && c.Spec.Supported {
			return true, nil
		}
	}
	return false, nil
}
