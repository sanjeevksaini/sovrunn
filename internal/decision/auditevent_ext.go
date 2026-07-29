package decision

import "github.com/sanjeevksaini/sovrunn/internal/apimeta"

// AuditLinkage is the FEATURE-0013-owned AuditEvent payload/linkage extension
// (design §7.4, §7.9; F13-AUDIT-001/002/003; F13-COMPAT-009; AD-015, AD-026).
//
// The FEATURE-0012 AuditEvent base envelope (apiconform.AuditEvent /
// api/schemas/audit-event.json) is preserved. This type carries only the
// additive decision-linkage fields and must not duplicate base AuditEvent
// fields (closure 7). Audit outcomes remain the closed AuditOutcome
// vocabulary (Succeeded/Denied/Failed).
//
// Correlation reuses the package Correlation carrier defined for
// DecisionRecord; it is never redefined here.
type AuditLinkage struct {
	DecisionRef        apimeta.TypedRef  `json:"decisionRef"`                  // uid required at validation
	AuditObligationRef *apimeta.TypedRef `json:"auditObligationRef,omitempty"` // optional obligation identity
	Correlation        Correlation       `json:"correlation"`                  // reused carrier; never redefined
}
