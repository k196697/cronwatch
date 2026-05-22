package metrics

import (
	"testing"
	"time"
)

func TestRecord_TracksSuccessAndFailure(t *testing.T) {
	c := New()
	c.Record("backup", true, 2*time.Second)
	c.Record("backup", false, 3*time.Second)

	m, ok := c.Get("backup")
	if !ok {
		t.Fatal("expected metrics for 'backup'")
	}
	if m.TotalRuns != 2 {
		t.Errorf("TotalRuns: want 2, got %d", m.TotalRuns)
	}
	if m.SuccessCount != 1 {
		t.Errorf("SuccessCount: want 1, got %d", m.SuccessCount)
	}
	if m.FailureCount != 1 {
		t.Errorf("FailureCount: want 1, got %d", m.FailureCount)
	}
}

func TestRecord_ComputesAvgDuration(t *testing.T) {
	c := New()
	c.Record("job", true, 2*time.Second)
	c.Record("job", true, 4*time.Second)

	m, _ := c.Get("job")
	if m.AvgDuration != 3*time.Second {
		t.Errorf("AvgDuration: want 3s, got %v", m.AvgDuration)
	}
}

func TestGet_UnknownJob(t *testing.T) {
	c := New()
	_, ok := c.Get("nonexistent")
	if ok {
		t.Error("expected false for unknown job")
	}
}

func TestAll_ReturnsCopy(t *testing.T) {
	c := New()
	c.Record("alpha", true, time.Second)
	c.Record("beta", false, time.Second)

	all := c.All()
	if len(all) != 2 {
		t.Errorf("All: want 2 entries, got %d", len(all))
	}

	// Mutating the returned slice must not affect the collector.
	all[0].TotalRuns = 999
	for _, name := range []string{"alpha", "beta"} {
		m, _ := c.Get(name)
		if m.TotalRuns == 999 {
			t.Errorf("All() returned a reference, not a copy")
		}
	}
}

func TestRecord_SetsLastRun(t *testing.T) {
	c := New()
	before := time.Now()
	c.Record("job", true, time.Millisecond)
	after := time.Now()

	m, _ := c.Get("job")
	if m.LastRun.Before(before) || m.LastRun.After(after) {
		t.Errorf("LastRun %v not within expected range [%v, %v]", m.LastRun, before, after)
	}
}
