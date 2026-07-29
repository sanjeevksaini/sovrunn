package validate

import (
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/decision"
)

// RFC 6901 pointers for the sensitivity/projection validation pass
// (design §9.1 step 9, §12; closure 10). Pass-local pointers for input-bound
// fields may be prefixed by orchestration (T-026). Profile projection pointers
// are rooted at /spec/projectionRules.
const (
	ptrDeclaredSensitivity = "/declaredSensitivity"
	ptrClassifications     = "/classifications"
	ptrContentCategories   = "/contentCategories"
	ptrSensitivityCeiling  = "/spec/sensitivityCeiling"
	ptrProjectionRules     = "/spec/projectionRules"

	// Schema-aligned structural identifier ceilings for content-category and
	// audience identifiers (api/schemas/decision-profile.json).
	sensitivityIDMaxLen     = 253
	projectionPointerMaxLen = 512
)

// SensitivityInput is the input to ValidateSensitivity (design §9.1 step 9,
// §12; F13-SEC-002…007; AD-019; closure 10).
//
// Design §8 does not list a one-line signature for this pass; the inputs
// required by architecture §12.4 and closure 10 are the governing profile,
// classification/sensitivity carriers, captured content-category evidence,
// prohibited-category rules, and the profile-declared canonical view used to
// resolve projection pointers.
//
// Projection conformance is structural only: no redaction execution or
// projected-output generation is performed (F13-SEC-007; §27.9(10)).
type SensitivityInput struct {
	Profile decision.DecisionProfile

	// Classifications are FEATURE-0012 DataClassification values that apply to
	// the validated content. The sensitivity floor is the most restrictive
	// mapped floor (architecture §12.4). Empty skips the floor check.
	Classifications []apimeta.DataClassification

	// DeclaredSensitivity is the FEATURE-0013 sensitivity declared for the
	// content. When Classifications is non-empty it must meet the derived
	// floor (raise-only; F13-SEC-003). When non-empty it must not exceed the
	// profile sensitivity ceiling (F13-SEC-005).
	DeclaredSensitivity decision.Sensitivity

	// ContentCategories are typed content-category identifiers present on the
	// captured/submitted content (structural evidence only; F13-SEC-004/006).
	ContentCategories []string

	// ProhibitedCategories are profile-declared prohibited content-category
	// identifiers (F13-SEC-004). Exposure of a prohibited category fails closed.
	// Profile.Spec.ContentCategories, when non-empty, is the allowed set.
	ProhibitedCategories []string

	// CanonicalView is the set of RFC 6901 pointers that exist in the
	// profile-declared canonical schema/view (closure 10). Projection
	// Includes/Excludes must resolve into this set; they are never resolved
	// against fields that happen to appear in one example record.
	CanonicalView []string

	// AllowedAudiences, when non-empty, constrains ProjectionRule.Audience
	// membership (audience mismatch). When empty, each rule's non-empty
	// Audience is accepted as profile-declared.
	AllowedAudiences []string
}

// ValidateSensitivity runs the FEATURE-0013 sensitivity and projection
// validation pass (design §9.1 step 9, §12; F13-SEC-002…007; AD-019).
//
// Checks (fail closed, short-circuit on first blocking violation):
//  1. profile sensitivity ceiling vocabulary
//  2. classification → sensitivity floor (unmapped/unknown/below-floor/
//     contradictory); raise-only — floor cannot be lowered
//  3. declared/captured sensitivity vs profile ceiling
//  4. content-category allowed/prohibited membership
//  5. projection rules: audience, pointer syntax/target against CanonicalView,
//     include/exclude overlap (equality + ancestor/descendant containment)
//
// Failures map to existing profile/evaluation codes (architecture §27.9(10));
// no DECISION_SENSITIVITY_* code is introduced. No content-scanning,
// redaction, or projected-output generation is performed.
func ValidateSensitivity(in SensitivityInput) *apiproblem.Problem {
	ceiling := in.Profile.Spec.SensitivityCeiling
	if !ceiling.Valid() {
		return sensitivityProblem(CodeProfileSchemaInvalid, ptrSensitivityCeiling,
			"spec.sensitivityCeiling must be a closed Sensitivity value")
	}

	if prob := checkSensitivityFloorAndCeiling(in, ceiling); prob != nil {
		return prob
	}
	if prob := checkContentCategories(in); prob != nil {
		return prob
	}
	if prob := checkProjectionRules(in); prob != nil {
		return prob
	}
	return nil
}

func checkSensitivityFloorAndCeiling(in SensitivityInput, ceiling decision.Sensitivity) *apiproblem.Problem {
	declared := in.DeclaredSensitivity
	hasDeclared := declared != ""
	hasClasses := len(in.Classifications) > 0

	if !hasClasses && !hasDeclared {
		return nil
	}

	if hasDeclared && !declared.Valid() {
		return sensitivityProblem(CodeProfileSchemaInvalid, ptrDeclaredSensitivity,
			"declared sensitivity must be a closed Sensitivity value")
	}

	if hasClasses {
		floor, ok := decision.MostRestrictiveFloor(in.Classifications...)
		if !ok {
			// Locate the first unmapped/unknown classification for the pointer.
			for i, c := range in.Classifications {
				if _, mapped := decision.SensitivityFloor(c); !mapped {
					return sensitivityProblem(CodeProfileSchemaInvalid, classificationPtr(i),
						"unknown or unmapped DataClassification fails closed against the sensitivity floor")
				}
			}
			return sensitivityProblem(CodeProfileSchemaInvalid, ptrClassifications,
				"unknown or unmapped DataClassification fails closed against the sensitivity floor")
		}

		if !hasDeclared {
			return sensitivityProblem(CodeProfileSchemaInvalid, ptrDeclaredSensitivity,
				"declared sensitivity is required when classifications apply")
		}

		// Raise-only: declared must meet the classification-derived floor.
		// Declaring below the floor is lowering and fails closed (F13-SEC-003/07/08).
		if !decision.MeetsFloor(declared, floor) {
			return sensitivityProblem(CodeProfileSchemaInvalid, ptrDeclaredSensitivity,
				"declared sensitivity is below the DataClassification-derived floor; the floor cannot be lowered")
		}

		// Contradictory profile: ceiling below the applicable floor.
		if ceiling.Rank() < floor.Rank() {
			return sensitivityProblem(CodeProfileSchemaInvalid, ptrSensitivityCeiling,
				"sensitivity ceiling is below the applicable classification-derived floor")
		}
	}

	if hasDeclared {
		// F13-SEC-005 / F13-SEC-01: captured/declared sensitivity above ceiling.
		if declared.Rank() > ceiling.Rank() {
			return sensitivityProblem(CodeEvaluationLimitExceeded, ptrDeclaredSensitivity,
				"declared sensitivity exceeds the profile sensitivity ceiling")
		}
	}
	return nil
}

func checkContentCategories(in SensitivityInput) *apiproblem.Problem {
	allowed := categorySet(in.Profile.Spec.ContentCategories)
	prohibited := categorySet(in.ProhibitedCategories)

	// Contradictory profile rules: a category cannot be both allowed and prohibited.
	for cat := range allowed {
		if _, bad := prohibited[cat]; bad {
			return sensitivityProblem(CodeProfileSchemaInvalid, "/spec/contentCategories",
				"content category cannot be both allowed and prohibited")
		}
	}

	for i, raw := range in.ContentCategories {
		cat := strings.TrimSpace(raw)
		if cat == "" || utf8.RuneCountInString(cat) > sensitivityIDMaxLen || cat != raw {
			return sensitivityProblem(CodeEvaluationResultInvalid, contentCategoryPtr(i),
				"content category identifier is structurally invalid")
		}
		if _, bad := prohibited[cat]; bad {
			return sensitivityProblem(CodeEvaluationResultInvalid, contentCategoryPtr(i),
				"prohibited content category exposure is rejected")
		}
		if len(allowed) > 0 {
			if _, ok := allowed[cat]; !ok {
				return sensitivityProblem(CodeEvaluationResultInvalid, contentCategoryPtr(i),
					"content category is not in the profile-declared allowed set")
			}
		}
	}
	return nil
}

func checkProjectionRules(in SensitivityInput) *apiproblem.Problem {
	rules := in.Profile.Spec.ProjectionRules
	if len(rules) == 0 {
		return nil
	}

	canonical := categorySet(in.CanonicalView)
	audiences := categorySet(in.AllowedAudiences)
	requireCanonical := projectionRulesReferencePointers(rules)
	if requireCanonical && len(canonical) == 0 {
		return sensitivityProblem(CodeProfileSchemaInvalid, ptrProjectionRules,
			"projection pointers require a profile-declared canonical view")
	}

	for i, rule := range rules {
		audience := strings.TrimSpace(rule.Audience)
		if audience == "" || utf8.RuneCountInString(audience) > sensitivityIDMaxLen || audience != rule.Audience {
			return sensitivityProblem(CodeProfileSchemaInvalid, projectionRulePtr(i, "audience"),
				"projection rule audience is required and must be a non-empty structural identifier")
		}
		if len(audiences) > 0 {
			if _, ok := audiences[audience]; !ok {
				return sensitivityProblem(CodeProfileSchemaInvalid, projectionRulePtr(i, "audience"),
					"projection rule audience is not in the profile-declared allowed audience set")
			}
		}

		for j, raw := range rule.Includes {
			if prob := checkProjectionPointer(i, "includes", j, raw, canonical, requireCanonical); prob != nil {
				return prob
			}
		}
		for j, raw := range rule.Excludes {
			if prob := checkProjectionPointer(i, "excludes", j, raw, canonical, requireCanonical); prob != nil {
				return prob
			}
		}

		// Overlap = equality + ancestor/descendant containment (closure 10).
		for ji, inc := range rule.Includes {
			for _, exc := range rule.Excludes {
				if pointersOverlap(inc, exc) {
					return sensitivityProblem(CodeProfileSchemaInvalid, projectionRulePtr(i, "includes")+"/"+strconv.Itoa(ji),
						"projection includes/excludes overlap (equality or containment) is invalid")
				}
			}
		}
	}
	return nil
}

func projectionRulesReferencePointers(rules []decision.ProjectionRule) bool {
	for _, rule := range rules {
		if len(rule.Includes) > 0 || len(rule.Excludes) > 0 {
			return true
		}
	}
	return false
}

func checkProjectionPointer(ruleIdx int, side string, ptrIdx int, raw string, canonical map[string]struct{}, requireCanonical bool) *apiproblem.Problem {
	field := projectionRulePtr(ruleIdx, side) + "/" + strconv.Itoa(ptrIdx)
	if raw == "" {
		// Empty JSON Pointer is the whole document; must be declared in the
		// canonical view when pointers are resolved against it (closure 10).
		if requireCanonical {
			if _, ok := canonical[""]; !ok {
				return sensitivityProblem(CodeProfileSchemaInvalid, field,
					"projection pointer is not present in the profile-declared canonical view")
			}
		}
		return nil
	}
	if !validJSONPointer(raw) {
		return sensitivityProblem(CodeProfileSchemaInvalid, field,
			"projection pointer must be a valid RFC 6901 JSON Pointer")
	}
	if utf8.RuneCountInString(raw) > projectionPointerMaxLen {
		return sensitivityProblem(CodeProfileSchemaInvalid, field,
			"projection pointer exceeds structural length ceiling")
	}
	if requireCanonical {
		if _, ok := canonical[raw]; !ok {
			return sensitivityProblem(CodeProfileSchemaInvalid, field,
				"projection pointer is not present in the profile-declared canonical view")
		}
	}
	return nil
}

// pointersOverlap reports whether a and b are equal or one contains the other
// as an RFC 6901 ancestor/descendant (closure 10).
func pointersOverlap(a, b string) bool {
	if a == b {
		return true
	}
	return pointerContains(a, b) || pointerContains(b, a)
}

// pointerContains reports whether ancestor contains descendant (or equals it)
// under RFC 6901 path-segment boundaries. The empty pointer (whole document)
// contains every pointer.
func pointerContains(ancestor, descendant string) bool {
	if ancestor == descendant {
		return true
	}
	if ancestor == "" {
		return true
	}
	if !strings.HasPrefix(descendant, ancestor) {
		return false
	}
	if len(descendant) == len(ancestor) {
		return true
	}
	return descendant[len(ancestor)] == '/'
}

// validJSONPointer reports whether p is a well-formed RFC 6901 JSON Pointer.
// The empty string (whole document) is valid. Non-empty pointers must start
// with '/' and use only well-formed ~0 / ~1 escapes.
func validJSONPointer(p string) bool {
	if p == "" {
		return true
	}
	if p[0] != '/' {
		return false
	}
	for i := 0; i < len(p); i++ {
		if p[i] != '~' {
			continue
		}
		if i+1 >= len(p) {
			return false
		}
		switch p[i+1] {
		case '0', '1':
			i++
		default:
			return false
		}
	}
	return true
}

func categorySet(values []string) map[string]struct{} {
	if len(values) == 0 {
		return nil
	}
	out := make(map[string]struct{}, len(values))
	for _, v := range values {
		out[v] = struct{}{}
	}
	return out
}

func classificationPtr(index int) string {
	return ptrClassifications + "/" + strconv.Itoa(index)
}

func contentCategoryPtr(index int) string {
	return ptrContentCategories + "/" + strconv.Itoa(index)
}

func projectionRulePtr(index int, field string) string {
	base := ptrProjectionRules + "/" + strconv.Itoa(index)
	if field == "" {
		return base
	}
	return base + "/" + field
}

func sensitivityProblem(code apiproblem.ViolationCode, field, message string) *apiproblem.Problem {
	return apiproblem.New(apiproblem.CodeValidationFailed).
		WithDetail(message).
		WithViolations([]apiproblem.Violation{{
			Field:   field,
			Code:    code,
			Message: message,
		}})
}
