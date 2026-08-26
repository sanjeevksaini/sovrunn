package policyeval

import (
	"testing"
	"time"
)

func TestUTCTimeSource_NowUTC(t *testing.T) {
	t.Parallel()

	var ts UTCTimeSource
	before := time.Now().UTC().Add(-time.Second)
	got := ts.NowUTC()
	after := time.Now().UTC().Add(time.Second)

	if got.Location() != time.UTC {
		t.Fatalf("location=%v want=UTC", got.Location())
	}
	if got.Before(before) || got.After(after) {
		t.Fatalf("now=%v out of expected window [%v, %v]", got, before, after)
	}
}

func TestFixedTimeSource_DeterministicAndNormalized(t *testing.T) {
	t.Parallel()

	loc := time.FixedZone("EST", -5*3600)
	supplied := time.Date(2026, 8, 25, 12, 30, 45, 123456789, loc)
	want := supplied.UTC()

	ts := NewFixedTimeSource(supplied)
	first := ts.NowUTC()
	second := ts.NowUTC()

	if first.Location() != time.UTC {
		t.Fatalf("location=%v want=UTC", first.Location())
	}
	if !first.Equal(want) {
		t.Fatalf("got=%v want=%v", first, want)
	}
	if !second.Equal(first) {
		t.Fatalf("repeated NowUTC changed: first=%v second=%v", first, second)
	}
}

func TestTimeSource_RFC3339NanoSerialization(t *testing.T) {
	t.Parallel()

	fixedAt := time.Date(2026, 8, 25, 15, 4, 5, 987654321, time.UTC)
	fixed := NewFixedTimeSource(fixedAt)
	gotFixed := fixed.NowUTC().Format(time.RFC3339Nano)
	wantFixed := fixedAt.Format(time.RFC3339Nano)
	if gotFixed != wantFixed {
		t.Fatalf("fixed RFC3339Nano=%q want=%q", gotFixed, wantFixed)
	}
	if _, err := time.Parse(time.RFC3339Nano, gotFixed); err != nil {
		t.Fatalf("fixed RFC3339Nano parse: %v", err)
	}

	var utc UTCTimeSource
	gotUTC := utc.NowUTC().Format(time.RFC3339Nano)
	parsed, err := time.Parse(time.RFC3339Nano, gotUTC)
	if err != nil {
		t.Fatalf("utc RFC3339Nano parse: %v", err)
	}
	if parsed.Location() != time.UTC {
		t.Fatalf("parsed location=%v want=UTC", parsed.Location())
	}
}
