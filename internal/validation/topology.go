package validation

import (
	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/resources"
)

// TopologyValue is an explicitly supplied in-memory projection of a
// FEATURE-0014 resource used by pure topology evaluation helpers
// (design §1.3, §3.1 I-5). It carries no store handle and implies no lookup.
//
// ProviderScopeUID is the resolved Provider scope UID for the value:
//   - Kind Provider: the resource's own metadata.uid
//   - descendants: metadata.scopeRef.uid when scopeRef.kind is Provider
//
// Isolation is expressed only through distinct ProviderScopeUID values;
// external operator names and native IDs are never consulted (F14-REQ-30).
type TopologyValue struct {
	Kind             string
	UID              string
	ProviderScopeUID string
}

// RequiredImmediateParentKind returns the single required immediate parent
// kind for childKind in the locked FEATURE-0014 hierarchy (F14-REQ-08,
// F14-AD-005). ok is false when childKind is not a hierarchy child
// (Provider has no topology parent; unknown kinds are rejected).
func RequiredImmediateParentKind(childKind string) (parentKind string, ok bool) {
	switch childKind {
	case resources.KindProviderLocation:
		return resources.KindProvider, true
	case resources.KindProviderDatacenter:
		return resources.KindProviderLocation, true
	case resources.KindDatacenterFailureDomain:
		return resources.KindProviderDatacenter, true
	case resources.KindInfrastructureStack:
		return resources.KindDatacenterFailureDomain, true
	default:
		return "", false
	}
}

// EvaluateScopeAndHierarchy reports whether parent and child agree on the
// same Provider scope UID and correct immediate-parent hierarchy ordering
// (F14-REQ-07, F14-REQ-08, F14-REQ-10, F14-REQ-13).
//
// Evaluation is pure and deterministic over the supplied values only: it
// performs no lookup, reads no store, and owns no state (design §1.3, I-5).
// Cross-Provider scope-UID mismatches and skipped/reordered hierarchy levels
// evaluate to rejection. Live rejection response shape and no-existence
// disclosure against authoritative state remain CONTRACT_ONLY / NO_TASK
// (design §3.1 C-2/C-8).
func EvaluateScopeAndHierarchy(parent, child TopologyValue) bool {
	required, ok := RequiredImmediateParentKind(child.Kind)
	if !ok {
		return false
	}
	if parent.Kind != required {
		return false
	}
	if parent.ProviderScopeUID == "" || child.ProviderScopeUID == "" {
		return false
	}
	return parent.ProviderScopeUID == child.ProviderScopeUID
}

// TopologyValueFromProvider projects a supplied Provider into a TopologyValue.
// The Provider's own UID is the Provider scope UID.
func TopologyValueFromProvider(p resources.Provider) TopologyValue {
	return TopologyValue{
		Kind:             p.Kind,
		UID:              p.Metadata.UID,
		ProviderScopeUID: p.Metadata.UID,
	}
}

// TopologyValueFromProviderLocation projects a supplied ProviderLocation.
func TopologyValueFromProviderLocation(loc resources.ProviderLocation) TopologyValue {
	return TopologyValue{
		Kind:             loc.Kind,
		UID:              loc.Metadata.UID,
		ProviderScopeUID: providerScopeUIDFromScopeRef(loc.Metadata.ScopeRef),
	}
}

// TopologyValueFromProviderDatacenter projects a supplied ProviderDatacenter.
func TopologyValueFromProviderDatacenter(dc resources.ProviderDatacenter) TopologyValue {
	return TopologyValue{
		Kind:             dc.Kind,
		UID:              dc.Metadata.UID,
		ProviderScopeUID: providerScopeUIDFromScopeRef(dc.Metadata.ScopeRef),
	}
}

// TopologyValueFromDatacenterFailureDomain projects a supplied
// DatacenterFailureDomain.
func TopologyValueFromDatacenterFailureDomain(fd resources.DatacenterFailureDomain) TopologyValue {
	return TopologyValue{
		Kind:             fd.Kind,
		UID:              fd.Metadata.UID,
		ProviderScopeUID: providerScopeUIDFromScopeRef(fd.Metadata.ScopeRef),
	}
}

// TopologyValueFromInfrastructureStack projects a supplied InfrastructureStack.
func TopologyValueFromInfrastructureStack(stack resources.InfrastructureStack) TopologyValue {
	return TopologyValue{
		Kind:             stack.Kind,
		UID:              stack.Metadata.UID,
		ProviderScopeUID: providerScopeUIDFromScopeRef(stack.Metadata.ScopeRef),
	}
}

// providerScopeUIDFromScopeRef returns the Provider scope UID from an
// explicitly supplied scopeRef. Non-Provider or absent scope yields empty,
// which cannot agree under EvaluateScopeAndHierarchy.
func providerScopeUIDFromScopeRef(scope *apimeta.ScopeRef) string {
	if scope == nil {
		return ""
	}
	if apimeta.ScopeKind(scope.Kind) != apimeta.ScopeCloudProvider {
		return ""
	}
	return scope.UID
}
