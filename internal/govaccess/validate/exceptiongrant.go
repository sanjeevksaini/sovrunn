package validate

import (
	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/model"
)

// ValidateExceptionGrant validates ExceptionGrant structural/local semantic rules
// (VS0-SCHEMA-067; REQ-F18-17).
func ValidateExceptionGrant(g model.ExceptionGrant) *apiproblem.Problem {
	if prob := requireTypeMeta(g.TypeMeta, model.APIVersionGov, model.KindExceptionGrant, ""); prob != nil {
		return prob
	}
	if prob := requireObjectName(g.Metadata); prob != nil {
		return prob
	}
	if prob := ValidateScopeApplicability(ContractExceptionGrant, g.Metadata.ScopeRef); prob != nil {
		return prob
	}
	if !g.Record.Effect.Valid() {
		return invalidEnum("/record/effect")
	}
	if prob := requireNonEmpty(g.Record.ControlRef.UID, "/record/controlRef/uid"); prob != nil {
		return prob
	}
	if prob := requireNonEmpty(g.Record.SubjectUID, "/record/subjectUid"); prob != nil {
		return prob
	}
	if g.Record.NotBefore.IsZero() || g.Record.ExpiresAt.IsZero() {
		return missingField("/record/notBefore")
	}
	if !g.Record.ExpiresAt.After(g.Record.NotBefore) {
		return validationFailed("/record/expiresAt", "expiresAt must be after notBefore", apiproblem.ViolationOutOfRange)
	}
	return nil
}
