package clock_test

import (
	"testing"
	"testing/quick"
	"time"

	"github.com/sanjeevksaini/sovrunn/internal/govaccess/clock"
)

func TestFixedClockReturnsConfiguredInstant(t *testing.T) {
	t.Parallel()
	fixed := time.Date(2026, 8, 28, 12, 0, 0, 0, time.UTC)
	clk := clock.NewFixedClock(fixed)
	for i := 0; i < 5; i++ {
		got := clk.NowUTC()
		if !got.Equal(fixed) {
			t.Fatalf("NowUTC() = %v, want %v", got, fixed)
		}
		if got.Location() != time.UTC {
			t.Fatalf("NowUTC location = %v, want UTC", got.Location())
		}
	}
}

func TestFixedClockNormalizesToUTC(t *testing.T) {
	t.Parallel()
	loc := time.FixedZone("UTC+5:30", 5*3600+30*60)
	in := time.Date(2026, 8, 28, 17, 30, 0, 0, loc)
	clk := clock.NewFixedClock(in)
	got := clk.NowUTC()
	want := in.UTC()
	if !got.Equal(want) {
		t.Fatalf("NowUTC() = %v, want %v", got, want)
	}
}

func TestFixtureClockManualAdvancement(t *testing.T) {
	t.Parallel()
	start := time.Date(2026, 8, 28, 12, 0, 0, 0, time.UTC)
	clk := clock.NewFixtureClock(start)
	if got := clk.NowUTC(); !got.Equal(start) {
		t.Fatalf("initial NowUTC() = %v, want %v", got, start)
	}
	clk.Advance(90 * time.Second)
	want := start.Add(90 * time.Second)
	if got := clk.NowUTC(); !got.Equal(want) {
		t.Fatalf("after Advance NowUTC() = %v, want %v", got, want)
	}
	next := start.Add(2 * time.Hour)
	clk.Set(next)
	if got := clk.NowUTC(); !got.Equal(next) {
		t.Fatalf("after Set NowUTC() = %v, want %v", got, next)
	}
}

func TestPropertyFixedClockStableUntilReconstructed(t *testing.T) {
	t.Parallel()
	fn := func(sec int64) bool {
		if sec < 0 {
			sec = -sec
		}
		fixed := time.Unix(sec%1_000_000_000, 0).UTC()
		clk := clock.NewFixedClock(fixed)
		for i := 0; i < 8; i++ {
			if !clk.NowUTC().Equal(fixed) {
				return false
			}
		}
		return true
	}
	if err := quick.Check(fn, &quick.Config{MaxCount: 64}); err != nil {
		t.Fatal(err)
	}
}

func TestPropertyFixtureClockStableUntilAdvanced(t *testing.T) {
	t.Parallel()
	fn := func(sec int64, steps uint8) bool {
		if sec < 0 {
			sec = -sec
		}
		start := time.Unix(sec%1_000_000_000, 0).UTC()
		clk := clock.NewFixtureClock(start)
		for i := 0; i < 4; i++ {
			if !clk.NowUTC().Equal(start) {
				return false
			}
		}
		n := int(steps%5) + 1
		for i := 0; i < n; i++ {
			clk.Advance(time.Second)
		}
		want := start.Add(time.Duration(n) * time.Second)
		return clk.NowUTC().Equal(want)
	}
	if err := quick.Check(fn, &quick.Config{MaxCount: 64}); err != nil {
		t.Fatal(err)
	}
}
