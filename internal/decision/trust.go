package decision

// TrustCarrier is the algorithm-agile integrity/signature and structural
// trust-state carrier (design §7.9; F13-TRUST-001; AD-022).
//
// Fields are structural identifiers and opaque trust-state tokens only.
// No canonicalization, digest computation, signing, or cryptographic
// verification is selected or executed in FEATURE-0013 (ADR-F13-002).
//
// RID-08 (closure 9): a present-but-structurally-empty optional TrustCarrier
// is invalid. An absent optional carrier is valid. A present-but-empty
// carrier never satisfies provenance, required trust, identity, or evidence.
type TrustCarrier struct {
	State                  string `json:"state"`                            // structural trust-state; required when carrier present
	CanonicalizationMethod string `json:"canonicalizationMethod,omitempty"` // structural id; no execution
	DigestAlgorithm        string `json:"digestAlgorithm,omitempty"`        // structural id; no computation
	CoveredFields          string `json:"coveredFields,omitempty"`          // covered-field descriptor
	SignatureAlgorithm     string `json:"signatureAlgorithm,omitempty"`     // structural id; no signing/verification
}

// IsStructurallyEmpty reports whether every TrustCarrier field is unset.
// Used to express the RID-08 present-but-empty rule as a pure structural
// predicate; validation mapping belongs to ValidateTrust (T-022).
func (c TrustCarrier) IsStructurallyEmpty() bool {
	return c.State == "" &&
		c.CanonicalizationMethod == "" &&
		c.DigestAlgorithm == "" &&
		c.CoveredFields == "" &&
		c.SignatureAlgorithm == ""
}
