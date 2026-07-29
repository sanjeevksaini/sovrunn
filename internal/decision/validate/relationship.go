package validate

import (
	"strings"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/decision"
)

// RFC 6901 pointers for the relationship validation pass (design §9.1 step 8, §9.3).
const (
	ptrRelationship       = "/record/relationship"
	ptrRelationshipKind   = "/record/relationship/kind"
	ptrRelationshipTarget = "/record/relationship/targetRef"
)

// RelationshipNode is one established DecisionRecord identity plus its optional
// append-only relationship edge (design §9.3). A nil Relationship marks a root.
// Effective state is never stored on the node; it is projection-only.
type RelationshipNode struct {
	Ref          apimeta.TypedRef
	Relationship *decision.DecisionRelationship
}

// RelationshipInput is the input to ValidateRelationship (design §8 / §9.1 step 8,
// §9.3; F13-IMMUT-002/004/005; AD-004, AD-034; closure 5).
//
// The design §8 summary signature is
// ValidateRelationship(r DecisionRelationship, chain []TypedRef, maxChain int).
// Self and Established are the same pass's additional inputs required by the
// deterministic conflict matrix (resolvable targets, cycle/chain walks,
// same-target duplicates, terminal revocation, and current-state forks).
type RelationshipInput struct {
	// Relationship is the candidate link on the new record (orientation:
	// new record → predecessor). Nil means the relationship is absent and
	// the pass is a no-op.
	Relationship *decision.DecisionRelationship

	// Self is the new record's uid-pinned identity (self-reference and cycle).
	Self apimeta.TypedRef

	// Established is the already-accepted append-only relationship graph.
	// Every resolvable predecessor must appear here (roots included with a
	// nil Relationship).
	Established []RelationshipNode

	// MaxChain is profile ValidityRules.MaxChainDepth. When <= 0, no
	// chain-depth ceiling is enforced (optional profile field). When > 0,
	// the candidate's chain depth (edge count from the new record to the
	// root) must not exceed MaxChain.
	MaxChain int
}

// ValidateRelationship runs the FEATURE-0013 relationship validation pass
// (design §9.1 step 8, §9.3; F13-IMMUT-002/004/005; AD-004, AD-034; closure 5).
//
// Deterministic algorithm order (short-circuits at the first failing step):
//
//	precondition. kind vocabulary → DECISION_RELATIONSHIP_KIND_INVALID
//	1. target present/resolvable/uid-pinned/not self-referential
//	   → DECISION_RELATIONSHIP_TARGET_INVALID
//	2. cycle check → DECISION_RELATIONSHIP_CYCLE
//	3. chain length → DECISION_RELATIONSHIP_CHAIN_LIMIT_EXCEEDED
//	4. same-target different authoritative kinds → DECISION_RELATIONSHIP_CONFLICT
//	5. terminal revocation conflict → DECISION_RELATIONSHIP_CONFLICT
//	6. fork / current-state conflict → DECISION_RELATIONSHIP_CONFLICT
//	7. projection — computes EFFECTIVE/SUPERSEDED/REVOKED; never stored
//
// Append-only: predecessor records are never mutated. Effective state is a
// projection only and is never written back to a canonical record.
func ValidateRelationship(in RelationshipInput) *apiproblem.Problem {
	if in.Relationship == nil {
		return nil
	}
	r := *in.Relationship

	// Precondition: closed RelationshipKind vocabulary (design §9.3).
	if !r.Kind.Valid() {
		return relationshipProblem(CodeRelationshipKindInvalid, ptrRelationshipKind,
			"relationship kind is not in the closed CORRECTS/SUPERSEDES/REVOKES vocabulary")
	}

	// Step 1: target validation.
	if prob := checkRelationshipTarget(r.TargetRef, in.Self, in.Established); prob != nil {
		return prob
	}

	byUID := indexRelationshipNodes(in.Established)

	// Step 2: cycle check (new record → … → new record).
	if relationshipCycle(in.Self.UID, r.TargetRef.UID, byUID) {
		return relationshipProblem(CodeRelationshipCycle, ptrRelationshipTarget,
			"relationship chain forms a cycle")
	}

	// Step 3: chain length.
	depth := 1 + relationshipChainDepth(r.TargetRef.UID, byUID)
	if in.MaxChain > 0 && depth > in.MaxChain {
		return relationshipProblem(CodeRelationshipChainLimitExceeded, ptrRelationship,
			"relationship chain depth exceeds the declared maximum")
	}

	targetUID := r.TargetRef.UID

	// Step 4: same-target different authoritative kinds (CORRECTS ↔ SUPERSEDES).
	if isAuthoritativeKind(r.Kind) {
		for _, succ := range currentSuccessors(targetUID, byUID, "") {
			if isAuthoritativeKind(succ.kind) && succ.kind != r.Kind {
				return relationshipProblem(CodeRelationshipConflict, ptrRelationship,
					"incompatible authoritative relationship kinds on the same target")
			}
		}
	}

	// Step 5: terminal revocation — any later relationship targeting a
	// revoked target or its revoked chain conflicts.
	if onRevokedChain(targetUID, byUID) {
		return relationshipProblem(CodeRelationshipConflict, ptrRelationshipTarget,
			"relationship targets a revoked target or revoked chain")
	}

	// Step 6: fork / current-state conflict.
	// REVOKES of a live target is valid terminal behavior and does not
	// conflict with existing CORRECTS/SUPERSEDES successors of that target.
	// CORRECTS/SUPERSEDES must not create more than one current successor.
	if isAuthoritativeKind(r.Kind) {
		currents := currentSuccessors(targetUID, byUID, "")
		// Exclude same-kind currents that are already present: adding another
		// current same-kind successor yields >1 effective head.
		currentCount := 1 // the candidate becomes a current successor
		for _, succ := range currents {
			if isAuthoritativeKind(succ.kind) {
				currentCount++
			}
		}
		if currentCount > 1 {
			return relationshipProblem(CodeRelationshipConflict, ptrRelationship,
				"relationship would create more than one effective/current successor")
		}
	}

	// Step 7: projection only — never written back.
	_ = ProjectEffectiveState(in.Self, append(in.Established, RelationshipNode{
		Ref:          in.Self,
		Relationship: &r,
	}))
	return nil
}

// ProjectEffectiveState computes the projection-only effective state for ref
// against an append-only relationship graph (design §9.3 step 7; F13-IMMUT-003).
// It never mutates nodes and never writes state back to a canonical record.
//
// Rules (any matching successor edge counts; historical heads remain
// SUPERSEDED/REVOKED even when a later record supersedes that successor):
//   - any REVOKES edge targeting ref → REVOKED
//   - else any CORRECTS/SUPERSEDES edge targeting ref → SUPERSEDED
//   - else EFFECTIVE
//
// "Current successor" filtering is used only by the conflict-matrix steps,
// not by this projection.
func ProjectEffectiveState(ref apimeta.TypedRef, graph []RelationshipNode) decision.EffectiveState {
	uid := strings.TrimSpace(ref.UID)
	if uid == "" {
		return decision.EffectiveStateEffective
	}
	revoked := false
	superseded := false
	for _, n := range graph {
		if n.Relationship == nil {
			continue
		}
		if strings.TrimSpace(n.Relationship.TargetRef.UID) != uid {
			continue
		}
		switch n.Relationship.Kind {
		case decision.RelationshipKindRevokes:
			revoked = true
		case decision.RelationshipKindCorrects, decision.RelationshipKindSupersedes:
			superseded = true
		}
	}
	if revoked {
		return decision.EffectiveStateRevoked
	}
	if superseded {
		return decision.EffectiveStateSuperseded
	}
	return decision.EffectiveStateEffective
}

func checkRelationshipTarget(target, self apimeta.TypedRef, established []RelationshipNode) *apiproblem.Problem {
	if !typedRefUIDPinned(target) {
		return relationshipProblem(CodeRelationshipTargetInvalid, ptrRelationshipTarget,
			"relationship targetRef requires apiVersion, kind, name, and uid pinning")
	}
	if strings.TrimSpace(target.Kind) != "" && target.Kind != decision.KindDecisionRecord {
		return relationshipProblem(CodeRelationshipTargetInvalid, ptrRelationshipTarget,
			"relationship targetRef kind must be DecisionRecord")
	}
	if !typedRefUIDPinned(self) {
		return relationshipProblem(CodeRelationshipTargetInvalid, ptrRelationshipTarget,
			"relationship self identity requires uid pinning for self-reference checks")
	}
	if refsMatchUID(target, self) {
		return relationshipProblem(CodeRelationshipTargetInvalid, ptrRelationshipTarget,
			"relationship targetRef must not be self-referential")
	}
	if !relationshipTargetResolvable(target, established) {
		return relationshipProblem(CodeRelationshipTargetInvalid, ptrRelationshipTarget,
			"relationship targetRef is unresolvable")
	}
	return nil
}

func relationshipTargetResolvable(target apimeta.TypedRef, established []RelationshipNode) bool {
	want := strings.TrimSpace(target.UID)
	for _, n := range established {
		if strings.TrimSpace(n.Ref.UID) == want {
			return true
		}
	}
	return false
}

func indexRelationshipNodes(nodes []RelationshipNode) map[string]RelationshipNode {
	out := make(map[string]RelationshipNode, len(nodes))
	for _, n := range nodes {
		uid := strings.TrimSpace(n.Ref.UID)
		if uid == "" {
			continue
		}
		out[uid] = n
	}
	return out
}

func relationshipCycle(selfUID, targetUID string, byUID map[string]RelationshipNode) bool {
	seen := make(map[string]struct{})
	cur := strings.TrimSpace(targetUID)
	for cur != "" {
		if cur == strings.TrimSpace(selfUID) {
			return true
		}
		if _, dup := seen[cur]; dup {
			return true
		}
		seen[cur] = struct{}{}
		n, ok := byUID[cur]
		if !ok || n.Relationship == nil {
			return false
		}
		cur = strings.TrimSpace(n.Relationship.TargetRef.UID)
	}
	return false
}

func relationshipChainDepth(uid string, byUID map[string]RelationshipNode) int {
	depth := 0
	seen := make(map[string]struct{})
	cur := strings.TrimSpace(uid)
	for cur != "" {
		if _, dup := seen[cur]; dup {
			return depth
		}
		seen[cur] = struct{}{}
		n, ok := byUID[cur]
		if !ok || n.Relationship == nil {
			return depth
		}
		depth++
		cur = strings.TrimSpace(n.Relationship.TargetRef.UID)
	}
	return depth
}

type relSuccessor struct {
	uid  string
	kind decision.RelationshipKind
}

// currentSuccessors returns non-superseded, non-revoked records whose
// relationship targets targetUID (design §9.3 "current successor").
// excludeUID skips a UID (unused by callers today; reserved for symmetry).
func currentSuccessors(targetUID string, byUID map[string]RelationshipNode, excludeUID string) []relSuccessor {
	targetUID = strings.TrimSpace(targetUID)
	excludeUID = strings.TrimSpace(excludeUID)
	var out []relSuccessor
	for uid, n := range byUID {
		if uid == excludeUID || n.Relationship == nil {
			continue
		}
		if strings.TrimSpace(n.Relationship.TargetRef.UID) != targetUID {
			continue
		}
		if recordSupersededOrRevoked(uid, byUID) {
			continue
		}
		out = append(out, relSuccessor{uid: uid, kind: n.Relationship.Kind})
	}
	return out
}

// recordSupersededOrRevoked reports whether any relationship targets uid with
// CORRECTS, SUPERSEDES, or REVOKES (making it non-current as a successor).
func recordSupersededOrRevoked(uid string, byUID map[string]RelationshipNode) bool {
	uid = strings.TrimSpace(uid)
	for _, n := range byUID {
		if n.Relationship == nil {
			continue
		}
		if strings.TrimSpace(n.Relationship.TargetRef.UID) != uid {
			continue
		}
		switch n.Relationship.Kind {
		case decision.RelationshipKindCorrects, decision.RelationshipKindSupersedes, decision.RelationshipKindRevokes:
			return true
		}
	}
	return false
}

func onRevokedChain(uid string, byUID map[string]RelationshipNode) bool {
	seen := make(map[string]struct{})
	cur := strings.TrimSpace(uid)
	for cur != "" {
		if _, dup := seen[cur]; dup {
			return false
		}
		seen[cur] = struct{}{}
		if hasRevokesLinkTo(cur, byUID) {
			return true
		}
		n, ok := byUID[cur]
		if !ok || n.Relationship == nil {
			return false
		}
		cur = strings.TrimSpace(n.Relationship.TargetRef.UID)
	}
	return false
}

func hasRevokesLinkTo(uid string, byUID map[string]RelationshipNode) bool {
	uid = strings.TrimSpace(uid)
	for _, n := range byUID {
		if n.Relationship == nil {
			continue
		}
		if n.Relationship.Kind == decision.RelationshipKindRevokes &&
			strings.TrimSpace(n.Relationship.TargetRef.UID) == uid {
			return true
		}
	}
	return false
}

func isAuthoritativeKind(k decision.RelationshipKind) bool {
	return k == decision.RelationshipKindCorrects || k == decision.RelationshipKindSupersedes
}

func refsMatchUID(a, b apimeta.TypedRef) bool {
	au := strings.TrimSpace(a.UID)
	bu := strings.TrimSpace(b.UID)
	return au != "" && au == bu
}

func relationshipProblem(code apiproblem.ViolationCode, field, message string) *apiproblem.Problem {
	return apiproblem.New(apiproblem.CodeValidationFailed).
		WithDetail(message).
		WithViolations([]apiproblem.Violation{{
			Field:   field,
			Code:    code,
			Message: message,
		}})
}
