package apimeta

import (
	"crypto/rand"
	"encoding/hex"
)

const (
	// uidByteLen is the opaque 128-bit UID entropy size (D-07, F12-META-004).
	uidByteLen = 16
	// uidHexLen is the encoded string length (16 bytes → 32 lowercase hex chars).
	uidHexLen = uidByteLen * 2
)

// GenerateUID returns a new opaque 128-bit collision-resistant UID encoded as
// 32 lowercase hexadecimal characters (D-07, F12-META-004).
//
// Collision resistance is not a collision-proof guarantee: the probability of
// collision is astronomically low but non-zero. Adopting storage MUST perform
// a uniqueness/collision check when persisting a new object and MUST reject
// or regenerate on a detected collision. Clients treat the value as opaque.
//
// PlatformScopeUID ("platform") is a reserved sentinel and is never a valid
// value produced by this function (see IsGeneratedUIDFormat).
func GenerateUID() (string, error) {
	var b [uidByteLen]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", err
	}
	return hex.EncodeToString(b[:]), nil
}

// IsGeneratedUIDFormat reports whether s matches the opaque 128-bit hex
// encoding produced by GenerateUID. Reserved sentinels such as
// PlatformScopeUID are intentionally not valid generated UID formats.
func IsGeneratedUIDFormat(s string) bool {
	if len(s) != uidHexLen {
		return false
	}
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch {
		case c >= '0' && c <= '9':
		case c >= 'a' && c <= 'f':
		default:
			return false
		}
	}
	return true
}
