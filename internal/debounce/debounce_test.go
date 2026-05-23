package debounce

import (
	"testing"
	"time"
)

func fixedNow(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestAllow_FirstCallAlwaysPasses(t *testing.T) {
	d := New(5 * time.Second)
	if !d.Allow("job-a") {
		t.Fatal("expected first call to be allowed")
	}
}

func TestAllow_SuppressedWithinWindow(t *testing.T) {
	now := time.Now()
	d := New(10 * time.Second)
	d.nowFunc = fixedNow(now)

	d.Allow("job-a")

	d.nowFunc = fixedNow(now.Add(5 * time.Second))
	if d.Allow("job-a") {
		t.Fatal("expected alert to be suppressed within window")
	}
}

func TestAllow_PassesAfterWindowExpires(t *testing.T) {
	now := time.Now()
	d := New(10 * time.Second)
	d.nowFunc = fixedNow(now)

	d.Allow("job-a")

	d.nowFunc = fixedNow(now.Add(11 * time.Second))
	if !d.Allow("job-a") {
		t.Fatal("expected alert to pass after window expires")
	}
}

func TestAllow_IndependentKeys(t *testing.T) {
	now := time.Now()
	d := New(10 * time.Second)
	d.nowFunc = fixedNow(now)

	d.Allow("job-a")

	d.nowFunc = fixedNow(now.Add(1 * time.Second))
	if !d.Allow("job-b") {
		t.Fatal("expected job-b to be independent of job-a")
	}
}

func TestReset_AllowsImmediateRetrigger(t *testing.T) {
	now := time.Now()
	d := New(10 * time.Second)
	d.nowFunc = fixedNow(now)

	d.Allow("job-a")
	d.Reset("job-a")

	if !d.Allow("job-a") {
		t.Fatal("expected allow after reset")
	}
}

func TestResetAll_ClearsAllKeys(t *testing.T) {
	now := time.Now()
	d := New(10 * time.Second)
	d.nowFunc = fixedNow(now)

	d.Allow("job-a")
	d.Allow("job-b")
	d.ResetAll()

	if !d.Allow("job-a") || !d.Allow("job-b") {
		t.Fatal("expected both keys to be cleared after ResetAll")
	}
}
