package validate

import (
	"strconv"

	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/apivalid"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/model"
)

// ValidateRoleDefinition validates RoleDefinition structural/local semantic rules
// (VS0-SCHEMA-022; REQ-F18-07).
func ValidateRoleDefinition(r model.RoleDefinition) *apiproblem.Problem {
	if prob := requireTypeMeta(r.TypeMeta, model.APIVersionIAM, RoleDefinitionKind, ""); prob != nil {
		return prob
	}
	if prob := requireObjectName(r.Metadata); prob != nil {
		return prob
	}
	if prob := ValidateScopeApplicability(ContractRoleDefinition, r.Metadata.ScopeRef); prob != nil {
		return prob
	}
	if prob := requireNonEmpty(r.Spec.Version, "/spec/version"); prob != nil {
		return prob
	}
	if len(r.Spec.Actions) == 0 {
		return missingField("/spec/actions")
	}
	for i, action := range r.Spec.Actions {
		if action == "" || action == "*" {
			return validationFailed("/spec/actions/"+strconv.Itoa(i), "actions must be non-empty registered identifiers without wildcards", apivalid.ViolationInvalidEnum)
		}
	}
	if !r.Spec.PublicationState.Valid() {
		return invalidEnum("/spec/publicationState")
	}
	if r.Spec.BaseClassification != "" && !r.Spec.BaseClassification.Valid() {
		return invalidEnum("/spec/baseClassification")
	}
	return nil
}
