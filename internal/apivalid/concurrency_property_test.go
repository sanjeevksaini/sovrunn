package apivalid

import (
	"fmt"
	"math/rand"
	"net/http"
	"testing"
	"unicode/utf8"

	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
)

// Deterministic seed for Property 10 reproducibility (F12-UPDATE-002).
const property10Seed int64 = 20260723

const property10Iterations = 100

// property10Alphabet is a printable opaque-token alphabet for generated
// resourceVersion / If-Match values. Comparison remains opaque string equality.
var property10Alphabet = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789._-:/+")

type property10Case struct {
	IfMatch                string
	CurrentResourceVersion string
}

// Feature: api-resource-naming-status-and-validation-standard, Property 10: Concurrency staleness
//
// For any pair of resource versions, CheckIfMatch returns a 412
// STALE_RESOURCE_VERSION Problem exactly when If-Match is non-empty and
// differs from the current resourceVersion; it returns nil when the values
// match or when no concurrency protection is required (absent If-Match).
//
// Validates: Requirements 4.11 (F12-UPDATE-002)
func TestProperty10_ConcurrencyStaleness(t *testing.T) {
	t.Parallel()

	rng := rand.New(rand.NewSource(property10Seed))
	for i := 0; i < property10Iterations; i++ {
		c := generateProperty10Case(rng, i)
		if err := checkProperty10Case(c, i); err != nil {
			t.Fatalf("property 10 failed at iteration %d (seed %d): %v", i, property10Seed, err)
		}
	}
}

func generateProperty10Case(rng *rand.Rand, iteration int) property10Case {
	// Force coverage buckets so each class of oracle outcome appears often.
	switch iteration % 5 {
	case 0:
		// Absent If-Match: no protection required regardless of current.
		return property10Case{
			IfMatch:                "",
			CurrentResourceVersion: property10MaybeVersion(rng),
		}
	case 1:
		// Exact match (including both empty: absent protection).
		v := property10MaybeVersion(rng)
		return property10Case{IfMatch: v, CurrentResourceVersion: v}
	case 2:
		// Non-empty If-Match against empty current → stale.
		return property10Case{
			IfMatch:                property10NonEmptyVersion(rng),
			CurrentResourceVersion: "",
		}
	case 3:
		// Distinct non-empty opaque versions → stale.
		a := property10NonEmptyVersion(rng)
		b := property10DistinctVersion(rng, a)
		return property10Case{IfMatch: a, CurrentResourceVersion: b}
	default:
		// Fully random pair (still deterministic via seed).
		return property10Case{
			IfMatch:                property10MaybeVersion(rng),
			CurrentResourceVersion: property10MaybeVersion(rng),
		}
	}
}

func property10MaybeVersion(rng *rand.Rand) string {
	if rng.Intn(5) == 0 {
		return ""
	}
	return property10NonEmptyVersion(rng)
}

func property10NonEmptyVersion(rng *rand.Rand) string {
	n := 1 + rng.Intn(48) // 1..48
	buf := make([]byte, n)
	for i := range buf {
		buf[i] = property10Alphabet[rng.Intn(len(property10Alphabet))]
	}
	return string(buf)
}

func property10DistinctVersion(rng *rand.Rand, other string) string {
	for {
		v := property10NonEmptyVersion(rng)
		if v != other {
			return v
		}
	}
}

func checkProperty10Case(c property10Case, iteration int) error {
	if !utf8.ValidString(c.IfMatch) || !utf8.ValidString(c.CurrentResourceVersion) {
		return fmt.Errorf("iteration %d: generator produced invalid UTF-8", iteration)
	}

	got := CheckIfMatch(c.IfMatch, c.CurrentResourceVersion)

	// Oracle: nil iff unprotected (absent If-Match) or exact opaque match.
	wantNil := c.IfMatch == "" || c.IfMatch == c.CurrentResourceVersion
	if wantNil {
		if got != nil {
			return fmt.Errorf("iteration %d: want nil for ifMatch=%q current=%q, got %#v",
				iteration, c.IfMatch, c.CurrentResourceVersion, got)
		}
		return nil
	}

	if got == nil {
		return fmt.Errorf("iteration %d: want 412 STALE_RESOURCE_VERSION for ifMatch=%q current=%q, got nil",
			iteration, c.IfMatch, c.CurrentResourceVersion)
	}
	if got.Status != http.StatusPreconditionFailed {
		return fmt.Errorf("iteration %d: status = %d, want %d (ifMatch=%q current=%q)",
			iteration, got.Status, http.StatusPreconditionFailed, c.IfMatch, c.CurrentResourceVersion)
	}
	if got.Code != apiproblem.CodeStaleResourceVersion {
		return fmt.Errorf("iteration %d: code = %q, want %q",
			iteration, got.Code, apiproblem.CodeStaleResourceVersion)
	}
	if got.Title != apiproblem.TitleFor(apiproblem.CodeStaleResourceVersion) {
		return fmt.Errorf("iteration %d: title = %q, want %q",
			iteration, got.Title, apiproblem.TitleFor(apiproblem.CodeStaleResourceVersion))
	}
	wantType := apiproblem.TypeURN(apiproblem.CodeStaleResourceVersion)
	if got.Type != wantType {
		return fmt.Errorf("iteration %d: type = %q, want %q", iteration, got.Type, wantType)
	}
	// Stable problem shape: no version leakage, no correlation fields filled here.
	if got.Detail != "" {
		return fmt.Errorf("iteration %d: detail must stay empty (no version leakage); got %q",
			iteration, got.Detail)
	}
	if got.RequestID != "" || got.Instance != "" {
		return fmt.Errorf("iteration %d: requestId/instance must stay empty for adopter correlation; got requestId=%q instance=%q",
			iteration, got.RequestID, got.Instance)
	}
	if len(got.Violations) != 0 {
		return fmt.Errorf("iteration %d: violations must be empty; got %#v", iteration, got.Violations)
	}

	// Determinism: same inputs always yield the same outcome class and code.
	again := CheckIfMatch(c.IfMatch, c.CurrentResourceVersion)
	if again == nil || again.Code != got.Code || again.Status != got.Status {
		return fmt.Errorf("iteration %d: CheckIfMatch non-deterministic: first=%#v second=%#v",
			iteration, got, again)
	}
	return nil
}
