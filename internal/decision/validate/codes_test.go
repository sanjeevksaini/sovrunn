package validate

import (
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
)

// architectureOwnedClosedRegistry is the exact closed public violations[].code
// set from architecture §15.2 / F13-ERR-002. Tests bind constants to this list;
// no invented, renamed, or aliased code is permitted.
var architectureOwnedClosedRegistry = []apiproblem.ViolationCode{
	"DECISION_PROFILE_UNKNOWN",
	"DECISION_PROFILE_VERSION_UNSUPPORTED",
	"DECISION_PROFILE_INACTIVE",
	"DECISION_PROFILE_SCHEMA_INVALID",
	"DECISION_PROFILE_LIMIT_INVALID",
	"DECISION_EVALUATION_TYPE_UNSUPPORTED",
	"DECISION_EVALUATION_RESULT_INVALID",
	"DECISION_EVALUATION_SCOPE_MISMATCH",
	"DECISION_EVALUATION_LIMIT_EXCEEDED",
	"DECISION_COMPOSITION_STRATEGY_UNSUPPORTED",
	"DECISION_COMPOSITION_GRAPH_INVALID",
	"DECISION_COMPOSITION_CYCLE",
	"DECISION_COMPOSITION_LIMIT_EXCEEDED",
	"DECISION_COMPOSITION_INPUT_CONFLICT",
	"DECISION_OBLIGATION_UNKNOWN",
	"DECISION_OBLIGATION_INVALID",
	"DECISION_OBLIGATION_UNSUPPORTED_MANDATORY",
	"DECISION_TRUST_REQUIRED",
	"DECISION_TRUST_UNKNOWN",
	"DECISION_TRUST_EXPIRED",
	"DECISION_TRUST_REVOKED",
	"DECISION_TRUST_MISMATCH",
	"DECISION_SCOPE_REQUIRED",
	"DECISION_SCOPE_INVALID",
	"DECISION_SCOPE_CONFLICT",
	"DECISION_SCOPE_MISMATCH",
	"DECISION_SCOPE_WIDENING",
	"DECISION_RELATIONSHIP_KIND_INVALID",
	"DECISION_RELATIONSHIP_TARGET_INVALID",
	"DECISION_RELATIONSHIP_CYCLE",
	"DECISION_RELATIONSHIP_CHAIN_LIMIT_EXCEEDED",
	"DECISION_RELATIONSHIP_CONFLICT",
}

func TestClosedViolationCodeRegistry(t *testing.T) {
	t.Parallel()

	const wantCount = 32
	if len(architectureOwnedClosedRegistry) != wantCount {
		t.Fatalf("architecture registry len=%d, want %d", len(architectureOwnedClosedRegistry), wantCount)
	}

	got := AllCodes()
	if len(got) != wantCount {
		t.Fatalf("AllCodes len=%d, want %d", len(got), wantCount)
	}
	if len(got) != len(architectureOwnedClosedRegistry) {
		t.Fatalf("AllCodes len=%d, architecture registry len=%d", len(got), len(architectureOwnedClosedRegistry))
	}

	seen := make(map[apiproblem.ViolationCode]struct{}, len(got))
	for i, want := range architectureOwnedClosedRegistry {
		if got[i] != want {
			t.Fatalf("AllCodes[%d]=%q, want architecture-owned %q", i, got[i], want)
		}
		if !Valid(got[i]) {
			t.Fatalf("%q from AllCodes must be Valid", got[i])
		}
		if _, dup := seen[got[i]]; dup {
			t.Fatalf("duplicate code %q in AllCodes", got[i])
		}
		seen[got[i]] = struct{}{}
	}
}

func TestClosedViolationCodeConstants(t *testing.T) {
	t.Parallel()

	checks := []struct {
		name string
		got  apiproblem.ViolationCode
		want apiproblem.ViolationCode
	}{
		{"CodeProfileUnknown", CodeProfileUnknown, "DECISION_PROFILE_UNKNOWN"},
		{"CodeProfileVersionUnsupported", CodeProfileVersionUnsupported, "DECISION_PROFILE_VERSION_UNSUPPORTED"},
		{"CodeProfileInactive", CodeProfileInactive, "DECISION_PROFILE_INACTIVE"},
		{"CodeProfileSchemaInvalid", CodeProfileSchemaInvalid, "DECISION_PROFILE_SCHEMA_INVALID"},
		{"CodeProfileLimitInvalid", CodeProfileLimitInvalid, "DECISION_PROFILE_LIMIT_INVALID"},
		{"CodeEvaluationTypeUnsupported", CodeEvaluationTypeUnsupported, "DECISION_EVALUATION_TYPE_UNSUPPORTED"},
		{"CodeEvaluationResultInvalid", CodeEvaluationResultInvalid, "DECISION_EVALUATION_RESULT_INVALID"},
		{"CodeEvaluationScopeMismatch", CodeEvaluationScopeMismatch, "DECISION_EVALUATION_SCOPE_MISMATCH"},
		{"CodeEvaluationLimitExceeded", CodeEvaluationLimitExceeded, "DECISION_EVALUATION_LIMIT_EXCEEDED"},
		{"CodeCompositionStrategyUnsupported", CodeCompositionStrategyUnsupported, "DECISION_COMPOSITION_STRATEGY_UNSUPPORTED"},
		{"CodeCompositionGraphInvalid", CodeCompositionGraphInvalid, "DECISION_COMPOSITION_GRAPH_INVALID"},
		{"CodeCompositionCycle", CodeCompositionCycle, "DECISION_COMPOSITION_CYCLE"},
		{"CodeCompositionLimitExceeded", CodeCompositionLimitExceeded, "DECISION_COMPOSITION_LIMIT_EXCEEDED"},
		{"CodeCompositionInputConflict", CodeCompositionInputConflict, "DECISION_COMPOSITION_INPUT_CONFLICT"},
		{"CodeObligationUnknown", CodeObligationUnknown, "DECISION_OBLIGATION_UNKNOWN"},
		{"CodeObligationInvalid", CodeObligationInvalid, "DECISION_OBLIGATION_INVALID"},
		{"CodeObligationUnsupportedMandatory", CodeObligationUnsupportedMandatory, "DECISION_OBLIGATION_UNSUPPORTED_MANDATORY"},
		{"CodeTrustRequired", CodeTrustRequired, "DECISION_TRUST_REQUIRED"},
		{"CodeTrustUnknown", CodeTrustUnknown, "DECISION_TRUST_UNKNOWN"},
		{"CodeTrustExpired", CodeTrustExpired, "DECISION_TRUST_EXPIRED"},
		{"CodeTrustRevoked", CodeTrustRevoked, "DECISION_TRUST_REVOKED"},
		{"CodeTrustMismatch", CodeTrustMismatch, "DECISION_TRUST_MISMATCH"},
		{"CodeScopeRequired", CodeScopeRequired, "DECISION_SCOPE_REQUIRED"},
		{"CodeScopeInvalid", CodeScopeInvalid, "DECISION_SCOPE_INVALID"},
		{"CodeScopeConflict", CodeScopeConflict, "DECISION_SCOPE_CONFLICT"},
		{"CodeScopeMismatch", CodeScopeMismatch, "DECISION_SCOPE_MISMATCH"},
		{"CodeScopeWidening", CodeScopeWidening, "DECISION_SCOPE_WIDENING"},
		{"CodeRelationshipKindInvalid", CodeRelationshipKindInvalid, "DECISION_RELATIONSHIP_KIND_INVALID"},
		{"CodeRelationshipTargetInvalid", CodeRelationshipTargetInvalid, "DECISION_RELATIONSHIP_TARGET_INVALID"},
		{"CodeRelationshipCycle", CodeRelationshipCycle, "DECISION_RELATIONSHIP_CYCLE"},
		{"CodeRelationshipChainLimitExceeded", CodeRelationshipChainLimitExceeded, "DECISION_RELATIONSHIP_CHAIN_LIMIT_EXCEEDED"},
		{"CodeRelationshipConflict", CodeRelationshipConflict, "DECISION_RELATIONSHIP_CONFLICT"},
	}

	if len(checks) != 32 {
		t.Fatalf("constant checks len=%d, want 32", len(checks))
	}
	for _, tc := range checks {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if tc.got != tc.want {
				t.Fatalf("%s=%q, want %q", tc.name, tc.got, tc.want)
			}
			if !Valid(tc.got) {
				t.Fatalf("%s must be Valid", tc.name)
			}
		})
	}
}

func TestClosedViolationCodeRejectsUnknown(t *testing.T) {
	t.Parallel()

	unknowns := []apiproblem.ViolationCode{
		"",
		"NOT_A_REAL_CODE",
		"DECISION_PROFILE_ALIASED",
		"VALIDATION_FAILED",
		"UNKNOWN_FIELD",
		"DECISION_PENDING",
		"DECISION_SCOPE_OWNER_DERIVED",
	}
	for _, c := range unknowns {
		if Valid(c) {
			t.Fatalf("unknown code %q must not be Valid", c)
		}
	}
}

func TestClosedViolationCodeNoFEATURE0012Merge(t *testing.T) {
	t.Parallel()

	// FEATURE-0012 owns top-level ErrorCode / shared ViolationCode values.
	// FEATURE-0013 DECISION_* codes must not collide with or absorb them.
	f12Codes := []apiproblem.ViolationCode{
		apiproblem.ViolationUnknownField,
		apiproblem.ViolationDuplicateField,
		apiproblem.ViolationOutOfRange,
		apiproblem.ViolationOperationTargetScopeMismatch,
	}
	for _, c := range f12Codes {
		if Valid(c) {
			t.Fatalf("FEATURE-0012 code %q must not be in FEATURE-0013 closed registry", c)
		}
	}
}
