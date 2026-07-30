package validation

import (
	"github.com/sanjeevksaini/sovrunn/internal/apicond"
	"github.com/sanjeevksaini/sovrunn/internal/resources"
)

// TopologyComplete condition type and stable reason codes (DD-03, F14-REQ-14).
// TopologyComplete is an administrative topology fact only; it carries no
// capability, capacity, placement, or readiness meaning.
const (
	ConditionTypeTopologyComplete = "TopologyComplete"
	ConditionTypeValid            = "Valid"

	ReasonLeafValid      = "LeafValid"
	ReasonPathComplete   = "PathComplete"
	ReasonPathIncomplete = "PathIncomplete"
)

// CompletenessValue is an explicitly supplied in-memory projection of a
// FEATURE-0014 resource used by pure topology-completeness helpers
// (design §1.3, §3.1 I-5, DD-03). It carries no store handle and implies
// no lookup.
//
// ParentUID is the immediate topology parent UID:
//   - Kind Provider: empty (root)
//   - Kind ProviderLocation: Provider scope UID (metadata.scopeRef)
//   - other descendants: typed immediate-parent reference UID
//
// Valid is the supplied registered-validity fact (Valid condition True).
// Registered validity is distinct from topology completeness (F14-AD-007,
// F14-REQ-12).
type CompletenessValue struct {
	Kind             string
	UID              string
	ProviderScopeUID string
	ParentUID        string
	Valid            bool
}

// TopologyCompletenessEvaluation is the deterministic outcome of evaluating
// TopologyComplete and immediate child existence over a supplied set
// (F14-REQ-14, F14-REQ-26, F14-REQ-29).
type TopologyCompletenessEvaluation struct {
	// Complete is true when TopologyComplete would be True.
	Complete bool
	// Reason is the stable PascalCase TopologyComplete reason code.
	Reason string
	// HasChildren is true when at least one immediate topology child exists
	// in the supplied set (contract basis for deletion rejection; F14-REQ-26).
	HasChildren bool
}

// EvaluateTopologyCompleteness deterministically evaluates whether subject has
// an unbroken required child path down to at least one Valid
// InfrastructureStack over the explicitly supplied set (F14-REQ-14, DD-03).
//
// Rules:
//   - A leaf InfrastructureStack is complete iff it is itself Valid.
//   - A zero-child registered parent is incomplete (Valid may still be true).
//   - Hierarchy skips and cross-Provider links never contribute a complete path.
//   - Completeness carries no capability/capacity/placement/readiness meaning.
//
// Evaluation is pure, order-independent, and free of lookup/persistence
// (design §1.3, §7.3; live recomputation remains CONTRACT_ONLY / NO_TASK).
func EvaluateTopologyCompleteness(subject CompletenessValue, set []CompletenessValue) TopologyCompletenessEvaluation {
	children := indexImmediateChildren(set)
	hasChildren := hasImmediateChild(subject, children)
	complete := topologyPathComplete(subject, children, map[string]bool{})
	reason := ReasonPathIncomplete
	if complete {
		if subject.Kind == resources.KindInfrastructureStack {
			reason = ReasonLeafValid
		} else {
			reason = ReasonPathComplete
		}
	}
	return TopologyCompletenessEvaluation{
		Complete:    complete,
		Reason:      reason,
		HasChildren: hasChildren,
	}
}

// HasImmediateTopologyChildren reports whether any value in the supplied set
// is a correct immediate topology child of parent (same Provider scope UID and
// locked hierarchy ordering). This is the contract basis for deletion
// rejection while children exist (F14-REQ-26, F14-AD-020). Live deletion
// execution and authoritative child query remain CONTRACT_ONLY / NO_TASK.
func HasImmediateTopologyChildren(parent CompletenessValue, set []CompletenessValue) bool {
	return hasImmediateChild(parent, indexImmediateChildren(set))
}

// CompletenessValueFromProvider projects a supplied Provider.
func CompletenessValueFromProvider(p resources.Provider) CompletenessValue {
	tv := TopologyValueFromProvider(p)
	return CompletenessValue{
		Kind:             tv.Kind,
		UID:              tv.UID,
		ProviderScopeUID: tv.ProviderScopeUID,
		ParentUID:        "",
		Valid: conditionStatusTrueAtGeneration(
			p.Status.Conditions,
			ConditionTypeValid,
			p.Metadata.Generation,
			p.Status.ObservedGeneration,
		),
	}
}

// CompletenessValueFromProviderLocation projects a supplied ProviderLocation.
// Immediate parent is the Provider identified by metadata.scopeRef (DD-08).
func CompletenessValueFromProviderLocation(loc resources.ProviderLocation) CompletenessValue {
	tv := TopologyValueFromProviderLocation(loc)
	return CompletenessValue{
		Kind:             tv.Kind,
		UID:              tv.UID,
		ProviderScopeUID: tv.ProviderScopeUID,
		ParentUID:        tv.ProviderScopeUID,
		Valid: conditionStatusTrueAtGeneration(
			loc.Status.Conditions,
			ConditionTypeValid,
			loc.Metadata.Generation,
			loc.Status.ObservedGeneration,
		),
	}
}

// CompletenessValueFromProviderDatacenter projects a supplied ProviderDatacenter.
func CompletenessValueFromProviderDatacenter(dc resources.ProviderDatacenter) CompletenessValue {
	tv := TopologyValueFromProviderDatacenter(dc)
	return CompletenessValue{
		Kind:             tv.Kind,
		UID:              tv.UID,
		ProviderScopeUID: tv.ProviderScopeUID,
		ParentUID:        dc.Spec.ProviderLocationRef.UID,
		Valid: conditionStatusTrueAtGeneration(
			dc.Status.Conditions,
			ConditionTypeValid,
			dc.Metadata.Generation,
			dc.Status.ObservedGeneration,
		),
	}
}

// CompletenessValueFromDatacenterFailureDomain projects a supplied
// DatacenterFailureDomain.
func CompletenessValueFromDatacenterFailureDomain(fd resources.DatacenterFailureDomain) CompletenessValue {
	tv := TopologyValueFromDatacenterFailureDomain(fd)
	return CompletenessValue{
		Kind:             tv.Kind,
		UID:              tv.UID,
		ProviderScopeUID: tv.ProviderScopeUID,
		ParentUID:        fd.Spec.ProviderDatacenterRef.UID,
		Valid: conditionStatusTrueAtGeneration(
			fd.Status.Conditions,
			ConditionTypeValid,
			fd.Metadata.Generation,
			fd.Status.ObservedGeneration,
		),
	}
}

// CompletenessValueFromInfrastructureStack projects a supplied InfrastructureStack.
func CompletenessValueFromInfrastructureStack(stack resources.InfrastructureStack) CompletenessValue {
	tv := TopologyValueFromInfrastructureStack(stack)
	return CompletenessValue{
		Kind:             tv.Kind,
		UID:              tv.UID,
		ProviderScopeUID: tv.ProviderScopeUID,
		ParentUID:        stack.Spec.DatacenterFailureDomainRef.UID,
		Valid: conditionStatusTrueAtGeneration(
			stack.Status.Conditions,
			ConditionTypeValid,
			stack.Metadata.Generation,
			stack.Status.ObservedGeneration,
		),
	}
}

// topologyPathComplete reports whether an unbroken required path to a Valid
// InfrastructureStack exists beneath v (or v itself when it is the leaf).
// visiting prevents cycles; evaluation remains deterministic for DAGs and
// rejects cyclic graphs without non-termination.
func topologyPathComplete(v CompletenessValue, children map[string][]CompletenessValue, visiting map[string]bool) bool {
	if v.UID != "" {
		if visiting[v.UID] {
			return false
		}
		visiting[v.UID] = true
		defer delete(visiting, v.UID)
	}

	if v.Kind == resources.KindInfrastructureStack {
		return v.Valid
	}
	if _, ok := requiredImmediateChildKind(v.Kind); !ok {
		return false
	}

	for _, child := range children[v.UID] {
		if !isImmediateTopologyChild(v, child) {
			continue
		}
		if topologyPathComplete(child, children, visiting) {
			return true
		}
	}
	return false
}

func hasImmediateChild(parent CompletenessValue, children map[string][]CompletenessValue) bool {
	for _, child := range children[parent.UID] {
		if isImmediateTopologyChild(parent, child) {
			return true
		}
	}
	return false
}

// isImmediateTopologyChild reports whether child is a locked-hierarchy,
// same-Provider immediate child of parent. Skips and cross-Provider links
// evaluate false (F14-REQ-08, F14-REQ-10, F14-REQ-13).
func isImmediateTopologyChild(parent, child CompletenessValue) bool {
	if parent.UID == "" || child.ParentUID != parent.UID {
		return false
	}
	return EvaluateScopeAndHierarchy(
		TopologyValue{
			Kind:             parent.Kind,
			UID:              parent.UID,
			ProviderScopeUID: parent.ProviderScopeUID,
		},
		TopologyValue{
			Kind:             child.Kind,
			UID:              child.UID,
			ProviderScopeUID: child.ProviderScopeUID,
		},
	)
}

// indexImmediateChildren groups supplied values by ParentUID. Grouping is
// order-independent for existence/completeness boolean outcomes.
func indexImmediateChildren(set []CompletenessValue) map[string][]CompletenessValue {
	out := make(map[string][]CompletenessValue, len(set))
	for _, v := range set {
		if v.ParentUID == "" {
			continue
		}
		out[v.ParentUID] = append(out[v.ParentUID], v)
	}
	return out
}

// requiredImmediateChildKind is the inverse of RequiredImmediateParentKind
// for the locked five-level hierarchy.
func requiredImmediateChildKind(parentKind string) (childKind string, ok bool) {
	switch parentKind {
	case resources.KindProvider:
		return resources.KindProviderLocation, true
	case resources.KindProviderLocation:
		return resources.KindProviderDatacenter, true
	case resources.KindProviderDatacenter:
		return resources.KindDatacenterFailureDomain, true
	case resources.KindDatacenterFailureDomain:
		return resources.KindInfrastructureStack, true
	default:
		return "", false
	}
}

// conditionStatusTrueAtGeneration accepts a True condition only when both the
// resource status and the condition observed the current desired-state
// generation. A stale current-fact condition must not make a topology complete
// after concurrent change (design §7.3/§8; F14-REQ-29).
func conditionStatusTrueAtGeneration(
	conds []apicond.Condition,
	condType string,
	metadataGeneration int64,
	statusObservedGeneration int64,
) bool {
	if statusObservedGeneration != metadataGeneration {
		return false
	}
	c, ok := apicond.GetCondition(conds, condType)
	return ok &&
		c.Status == apicond.ConditionTrue &&
		c.ObservedGeneration == metadataGeneration
}
