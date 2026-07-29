package validate

import (
	"strconv"
	"strings"

	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/decision"
	"github.com/sanjeevksaini/sovrunn/internal/decision/bundle"
	"github.com/sanjeevksaini/sovrunn/internal/decision/graph"
)

// RFC 6901 pointers for the composition validation pass (design §9.1 step 5).
// Strategy pointers are record-rooted. Graph failures use the deterministic
// pointer-base rule from design §7.8 (graph-local, optionally prefixed).
const (
	ptrCompositionStrategyName    = "/record/composition/strategyName"
	ptrCompositionStrategyVersion = "/record/composition/strategyVersion"
	ptrCompositionGraphRef        = "/record/composition/graphRef"
	ptrCompositionGraphVersion    = "/record/composition/graphVersion"
	ptrCompositionInputRefs       = "/record/composition/inputRefs"
)

// CompositionInput is the input to ValidateComposition (design §9.1 step 5;
// F13-COMP-001…004; AD-006).
//
// Composition may be absent on a DecisionRecord. When StrategyName/Version and
// graph identity are all empty and ResolvedGraph is nil, the pass is a no-op.
//
// PointerBase implements design §7.8:
//   - empty → standalone graph-local pointers (/nodes/0/id)
//   - "/spec/graphs/{index}" → bundle-embedded pointers
//     (/spec/graphs/{index}/nodes/0/id)
type CompositionInput struct {
	Composition decision.CompositionRef
	Profile     decision.DecisionProfile

	// Bundle enables graph reference/version resolution via ResolveGraph
	// before Build (closure 4). Required when Composition.GraphRef is set
	// and ResolvedGraph is nil.
	Bundle *bundle.Bundle

	// ResolvedGraph is an already-resolved Definition (standalone artifact
	// or pre-resolved bundle entry). When set, Bundle resolution is skipped.
	ResolvedGraph *graph.Definition

	// PointerBase is the RFC 6901 prefix preserved and prepended to every
	// graph-local pointer. Validators never rewrite or discard the base.
	PointerBase string
}

// ValidateComposition runs the FEATURE-0013 composition validation pass
// (design §7.8, §9.1 step 5, §9.2; F13-COMP-001…004; AD-006).
//
// Checks (fail closed, short-circuit on first blocking violation):
//  1. composition strategy accepted by the governing profile
//  2. graph reference/version resolution in the declarative bundle (when needed)
//  3. in-memory graph.Build under effective inherited bounds
//  4. structural input-conflict checks against input-kind nodes
//
// Graph version/reference mismatch is owned by bundle.ResolveGraph before
// Build. Build performs no registry lookup (closure 4).
func ValidateComposition(in CompositionInput) *apiproblem.Problem {
	hasStrategy := strings.TrimSpace(in.Composition.StrategyName) != "" ||
		strings.TrimSpace(in.Composition.StrategyVersion) != ""
	hasGraphRef := strings.TrimSpace(in.Composition.GraphRef) != "" ||
		strings.TrimSpace(in.Composition.GraphVersion) != ""
	hasGraph := in.ResolvedGraph != nil || hasGraphRef

	if !hasStrategy && !hasGraph {
		return nil
	}

	// A composition graph is always paired with a versioned strategy
	// (F13-COMP-001/002). Standalone graph checks use ValidateGraph.
	if hasStrategy || hasGraph {
		if prob := checkCompositionStrategy(in.Composition, in.Profile); prob != nil {
			return prob
		}
	}

	if !hasGraph {
		return nil
	}

	def, prob := resolveCompositionGraph(in)
	if prob != nil {
		return prob
	}

	bounds := EffectiveGraphBounds(in.Profile.Spec.Limits, def.Bounds)
	built, buildProb := graph.Build(def, bounds)
	if buildProb != nil {
		return prefixGraphProblem(buildProb, in.PointerBase)
	}
	return checkCompositionInputConflicts(built, def, in.Composition.InputRefs, in.PointerBase)
}

// ValidateGraph validates an already-resolved graph Definition under effective
// Bounds with the deterministic pointer-base rule (design §7.8).
//
// pointerBase is preserved and prepended to every graph-local RFC 6901 pointer
// returned by graph.Build. Empty pointerBase yields standalone root-relative
// pointers such as /nodes/0/id.
func ValidateGraph(def graph.Definition, bounds graph.Bounds, pointerBase string) *apiproblem.Problem {
	_, prob := graph.Build(def, bounds)
	return prefixGraphProblem(prob, pointerBase)
}

// EffectiveGraphBounds returns the most restrictive applicable graph ceilings
// from FEATURE-0012 platform outer bounds, DecisionProfile limits, and the
// definition-declared Bounds (design §1.2 / §7.8). Zero means unset.
func EffectiveGraphBounds(profile decision.ProfileLimits, def graph.Bounds) graph.Bounds {
	plat := InheritedPlatformCeilings()
	return graph.Bounds{
		MaxNodes:        EffectiveIntLimit(0, profile.MaxNodes, def.MaxNodes),
		MaxEdges:        EffectiveIntLimit(0, profile.MaxEdges, def.MaxEdges),
		MaxDepth:        EffectiveIntLimit(plat.MaxDepth, profile.MaxDepth, def.MaxDepth),
		MaxWidth:        EffectiveIntLimit(0, profile.MaxWidth, def.MaxWidth),
		MaxFanOut:       EffectiveIntLimit(0, profile.MaxFanOut, def.MaxFanOut),
		MaxPayloadBytes: EffectiveInt64Limit(plat.MaxPayloadBytes, profile.MaxPayloadBytes, def.MaxPayloadBytes),
	}
}

func checkCompositionStrategy(c decision.CompositionRef, p decision.DecisionProfile) *apiproblem.Problem {
	name := strings.TrimSpace(c.StrategyName)
	version := strings.TrimSpace(c.StrategyVersion)

	if name == "" {
		return compositionProblem(CodeCompositionStrategyUnsupported, ptrCompositionStrategyName,
			"composition strategyName is required and must be accepted by the governing profile")
	}
	if version == "" {
		return compositionProblem(CodeCompositionStrategyUnsupported, ptrCompositionStrategyVersion,
			"composition strategyVersion is required and must be accepted by the governing profile")
	}
	if !strategyAccepted(name, version, p.Spec.AcceptedCompositionStrategies) {
		return compositionProblem(CodeCompositionStrategyUnsupported, ptrCompositionStrategyName,
			"composition strategy is not accepted by the governing decision profile")
	}
	return nil
}

func strategyAccepted(name, version string, accepted []string) bool {
	for _, a := range accepted {
		a = strings.TrimSpace(a)
		if a == "" {
			continue
		}
		if a == name {
			return true
		}
		if a == name+"@"+version || a == name+":"+version {
			return true
		}
	}
	return false
}

func resolveCompositionGraph(in CompositionInput) (graph.Definition, *apiproblem.Problem) {
	if in.ResolvedGraph != nil {
		return *in.ResolvedGraph, nil
	}

	ref := strings.TrimSpace(in.Composition.GraphRef)
	version := strings.TrimSpace(in.Composition.GraphVersion)
	if ref == "" {
		return graph.Definition{}, compositionProblem(CodeCompositionGraphInvalid, ptrCompositionGraphRef,
			"composition graphRef is required when a graph is declared")
	}
	if version == "" {
		return graph.Definition{}, compositionProblem(CodeCompositionGraphInvalid, ptrCompositionGraphVersion,
			"composition graphVersion is required when a graph is declared")
	}
	if in.Bundle == nil {
		return graph.Definition{}, compositionProblem(CodeCompositionGraphInvalid, ptrCompositionGraphRef,
			"composition graph resolution requires a loaded declarative bundle")
	}

	def, prob := in.Bundle.ResolveGraph(bundle.GraphRef{Ref: ref, Version: version})
	if prob == nil {
		return def, nil
	}
	return graph.Definition{}, remapResolveGraphProblem(prob)
}

func remapResolveGraphProblem(prob *apiproblem.Problem) *apiproblem.Problem {
	if prob == nil || len(prob.Violations) == 0 {
		return prob
	}
	field := prob.Violations[0].Field
	switch field {
	case "/ref":
		field = ptrCompositionGraphRef
	case "/version":
		field = ptrCompositionGraphVersion
	default:
		field = ptrCompositionGraphRef
	}
	msg := prob.Violations[0].Message
	if msg == "" {
		msg = prob.Detail
	}
	return compositionProblem(CodeCompositionGraphInvalid, field, msg)
}

func checkCompositionInputConflicts(g graph.Graph, def graph.Definition, inputRefs []string, pointerBase string) *apiproblem.Problem {
	if len(inputRefs) == 0 {
		return nil
	}

	seen := make(map[string]int, len(inputRefs))
	for i, raw := range inputRefs {
		ref := strings.TrimSpace(raw)
		ptr := ptrCompositionInputRefs + "/" + strconv.Itoa(i)
		if ref == "" {
			return compositionProblem(CodeCompositionInputConflict, ptr,
				"composition inputRefs entries must be non-empty")
		}
		if prev, dup := seen[ref]; dup {
			_ = prev
			if idx := nodeIndex(def, ref); idx >= 0 {
				return compositionProblem(CodeCompositionInputConflict,
					joinPointerBase(pointerBase, "/nodes/"+strconv.Itoa(idx)),
					"composition inputRefs contain a conflicting duplicate input")
			}
			return compositionProblem(CodeCompositionInputConflict, ptr,
				"composition inputRefs contain a conflicting duplicate input")
		}
		seen[ref] = i

		n, ok := g.Node(ref)
		if !ok {
			// Opaque external input identities are permitted; only graph node
			// references are checked for kind conflicts.
			continue
		}
		if n.Kind != graph.NodeKindInput {
			idx := nodeIndex(def, ref)
			return compositionProblem(CodeCompositionInputConflict,
				joinPointerBase(pointerBase, "/nodes/"+strconv.Itoa(idx)),
				"composition inputRefs must reference input-kind nodes")
		}
	}
	return nil
}

func nodeIndex(def graph.Definition, id string) int {
	for i, n := range def.Nodes {
		if n.ID == id {
			return i
		}
	}
	return -1
}

func prefixGraphProblem(prob *apiproblem.Problem, pointerBase string) *apiproblem.Problem {
	if prob == nil {
		return nil
	}
	if pointerBase == "" || len(prob.Violations) == 0 {
		return prob
	}
	violations := make([]apiproblem.Violation, len(prob.Violations))
	for i, v := range prob.Violations {
		violations[i] = apiproblem.Violation{
			Field:   joinPointerBase(pointerBase, v.Field),
			Code:    v.Code,
			Message: v.Message,
		}
	}
	out := apiproblem.New(apiproblem.CodeValidationFailed).
		WithDetail(prob.Detail).
		WithViolations(violations)
	return out
}

func joinPointerBase(base, local string) string {
	if base == "" {
		return local
	}
	if local == "" || local == "/" {
		return strings.TrimSuffix(base, "/")
	}
	if !strings.HasPrefix(local, "/") {
		local = "/" + local
	}
	return strings.TrimSuffix(base, "/") + local
}

func compositionProblem(code apiproblem.ViolationCode, field, message string) *apiproblem.Problem {
	return apiproblem.New(apiproblem.CodeValidationFailed).
		WithDetail(message).
		WithViolations([]apiproblem.Violation{{
			Field:   field,
			Code:    code,
			Message: message,
		}})
}
