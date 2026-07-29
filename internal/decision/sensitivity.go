package decision

import "github.com/sanjeevksaini/sovrunn/internal/apimeta"

// Sensitivity vocabulary (PUBLIC < INTERNAL < CONFIDENTIAL < RESTRICTED) is
// defined in vocab.go (T-001; F13-SEC-001). This file owns the FEATURE-0012
// DataClassification → FEATURE-0013 sensitivity-floor map and raise-only /
// fail-closed structural helpers (design §12; architecture §12.4; F13-SEC-002/003;
// AD-019).
//
// The two dimensions coexist: sensitivity never replaces, redefines, erases, or
// lowers DataClassification. Mapping is structural only; no content-scanning,
// redaction, or projection runtime is implemented here.

// SensitivityFloor returns the canonical minimum Sensitivity for classification
// c per design §12 / architecture §12.4:
//
//	Public                 → PUBLIC
//	Customer-visible       → CONFIDENTIAL
//	Tenant-confidential    → CONFIDENTIAL
//	Operator-confidential  → CONFIDENTIAL
//	Internal               → INTERNAL
//	Sensitive              → RESTRICTED
//	Secret-reference-only  → RESTRICTED
//
// ok is false for unknown or unmapped classifications (fail-closed).
func SensitivityFloor(c apimeta.DataClassification) (Sensitivity, bool) {
	switch c {
	case apimeta.ClassPublic:
		return SensitivityPublic, true
	case apimeta.ClassCustomerVisible,
		apimeta.ClassTenantConfidential,
		apimeta.ClassOperatorConfidential:
		return SensitivityConfidential, true
	case apimeta.ClassInternal:
		return SensitivityInternal, true
	case apimeta.ClassSensitive,
		apimeta.ClassSecretReferenceOnly:
		return SensitivityRestricted, true
	default:
		return "", false
	}
}

// Rank returns the ascending sensitivity order index:
// PUBLIC=0, INTERNAL=1, CONFIDENTIAL=2, RESTRICTED=3.
// Invalid values return -1.
func (s Sensitivity) Rank() int {
	switch s {
	case SensitivityPublic:
		return 0
	case SensitivityInternal:
		return 1
	case SensitivityConfidential:
		return 2
	case SensitivityRestricted:
		return 3
	default:
		return -1
	}
}

// MeetsFloor reports whether declared is at least as restrictive as floor
// (raise-only). Both values must be Valid; otherwise the check fails closed.
func MeetsFloor(declared, floor Sensitivity) bool {
	if !declared.Valid() || !floor.Valid() {
		return false
	}
	return declared.Rank() >= floor.Rank()
}

// ApplyRaisedFloor returns raised when it meets or exceeds baseFloor
// (raise-only; F13-SEC-003). It never lowers the floor. ok is false when
// either value is invalid or raised is below baseFloor (below-floor /
// contradictory fail-closed).
func ApplyRaisedFloor(baseFloor, raised Sensitivity) (Sensitivity, bool) {
	if !MeetsFloor(raised, baseFloor) {
		return "", false
	}
	return raised, true
}

// MostRestrictiveSensitivity returns the most restrictive Valid sensitivity
// among values. ok is false when values is empty or any value is invalid
// (unknown/contradictory fail-closed).
func MostRestrictiveSensitivity(values ...Sensitivity) (Sensitivity, bool) {
	if len(values) == 0 {
		return "", false
	}
	best := values[0]
	if !best.Valid() {
		return "", false
	}
	for _, s := range values[1:] {
		if !s.Valid() {
			return "", false
		}
		if s.Rank() > best.Rank() {
			best = s
		}
	}
	return best, true
}

// MostRestrictiveFloor returns the most restrictive SensitivityFloor among
// classifications. When multiple classifications apply, validation uses the
// most restrictive result (architecture §12.4). ok is false when
// classifications is empty or any classification is unknown/unmapped.
func MostRestrictiveFloor(classes ...apimeta.DataClassification) (Sensitivity, bool) {
	if len(classes) == 0 {
		return "", false
	}
	floors := make([]Sensitivity, 0, len(classes))
	for _, c := range classes {
		floor, ok := SensitivityFloor(c)
		if !ok {
			return "", false
		}
		floors = append(floors, floor)
	}
	return MostRestrictiveSensitivity(floors...)
}

// CheckClassificationSensitivity verifies that declared meets the sensitivity
// floor derived from classification (raise-only). It fails closed for
// unmapped/unknown classification, invalid declared sensitivity, or
// below-floor / contradictory declared values (F13-SEC-003).
func CheckClassificationSensitivity(c apimeta.DataClassification, declared Sensitivity) bool {
	floor, ok := SensitivityFloor(c)
	if !ok {
		return false
	}
	return MeetsFloor(declared, floor)
}
