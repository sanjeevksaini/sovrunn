package isocodes

// IsAssignedAlpha2 reports whether code is an officially assigned ISO-3166-1
// alpha-2 value in AssignedAlpha2. The check is exact and case-sensitive:
// lowercase or mixed-case input is rejected rather than normalized.
func IsAssignedAlpha2(code string) bool {
	_, ok := AssignedAlpha2[code]
	return ok
}
