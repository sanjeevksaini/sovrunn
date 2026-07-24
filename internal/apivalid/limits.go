package apivalid

// Limits bounds decode-time, validation, and list paging checks
// (F12-VALIDATION-007, F12-LIST-002, D-06).
//
// Zero on a field means that check is not enforced, so callers may pass
// partial Limits in tests. Production callers should use DefaultLimits.
type Limits struct {
	MaxObjectBytes        int // 1_048_576 (1 MiB)
	MaxNestingDepth       int // 32
	MaxLabels             int // 64
	MaxLabelKeyChars      int // 63
	MaxLabelValueChars    int // 253
	MaxAnnotationsBytes   int // 262_144 (256 KiB)
	MaxConditions         int // 32
	MaxReferencesPerField int // 64
	MaxViolations         int // 100
	DefaultPageSize       int // 50
	MaxPageSize           int // 200
}

// DefaultLimits returns the reviewed platform default Limits (D-06).
// Values are overridable via validated configuration; they are not
// provider-specific.
func DefaultLimits() Limits {
	return Limits{
		MaxObjectBytes:        1_048_576, // 1 MiB; matches Phase 1 MaxBytesReader (1<<20)
		MaxNestingDepth:       32,
		MaxLabels:             64,
		MaxLabelKeyChars:      63,
		MaxLabelValueChars:    253,
		MaxAnnotationsBytes:   262_144, // 256 KiB total
		MaxConditions:         32,
		MaxReferencesPerField: 64,
		MaxViolations:         100,
		DefaultPageSize:       50,
		MaxPageSize:           200,
	}
}
