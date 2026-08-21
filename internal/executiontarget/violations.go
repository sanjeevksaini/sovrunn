package executiontarget

import "github.com/sanjeevksaini/sovrunn/internal/apiproblem"

// Closed local FEATURE-0016 violation codes (design §5.2; REQ-F16-10).
// Exactly eight new codes appear only in violations[]. Inherited codes
// (VS0_STATUS_FIELD_WRITE, VS0_SYSTEM_OWNED_FIELD_WRITE,
// VS0_AUTHORIZATION_SAFE_DENIAL, VS0_IDEMPOTENCY_KEY_REUSE_MISMATCH) are
// reused from prior features and are intentionally not redefined here.
const (
	ViolationTargetRetired                           apiproblem.ViolationCode = "VS0_TARGET_RETIRED"
	ViolationTargetMaintenance                       apiproblem.ViolationCode = "VS0_TARGET_MAINTENANCE"
	ViolationTargetQualificationInProgress           apiproblem.ViolationCode = "VS0_TARGET_QUALIFICATION_IN_PROGRESS"
	ViolationTargetEpochStale                        apiproblem.ViolationCode = "VS0_TARGET_EPOCH_STALE"
	ViolationExecutionTargetParticipationUnavailable apiproblem.ViolationCode = "VS0_EXECUTION_TARGET_PARTICIPATION_UNAVAILABLE"
	ViolationExecutionTargetStackUnavailable         apiproblem.ViolationCode = "VS0_EXECUTION_TARGET_STACK_UNAVAILABLE"
	ViolationExecutionTargetScopeMismatch            apiproblem.ViolationCode = "VS0_EXECUTION_TARGET_SCOPE_MISMATCH"
	ViolationExecutionTargetViabilityStale           apiproblem.ViolationCode = "VS0_EXECUTION_TARGET_VIABILITY_STALE"
)

// AllLocalViolationCodes returns the closed F0016-local violation set in
// stable order. The length is exactly eight.
func AllLocalViolationCodes() []apiproblem.ViolationCode {
	return []apiproblem.ViolationCode{
		ViolationTargetRetired,
		ViolationTargetMaintenance,
		ViolationTargetQualificationInProgress,
		ViolationTargetEpochStale,
		ViolationExecutionTargetParticipationUnavailable,
		ViolationExecutionTargetStackUnavailable,
		ViolationExecutionTargetScopeMismatch,
		ViolationExecutionTargetViabilityStale,
	}
}
