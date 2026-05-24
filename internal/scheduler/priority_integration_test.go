package scheduler_test

import (
	"testing"

	"cronwatch/internal/jobpriority"
)

// TestPriorityRegistry_HighBeatsLow verifies that a High-priority job
// reports a greater level than a Low-priority job, which the scheduler
// can use when ordering execution under contention.
func TestPriorityRegistry_HighBeatsLow(t *testing.T) {
	reg := jobpriority.New()
	_ = reg.Set("critical-backup", jobpriority.High)
	_ = reg.Set("weekly-report", jobpriority.Low)

	if reg.Get("critical-backup") <= reg.Get("weekly-report") {
		t.Error("expected High priority to exceed Low priority")
	}
}

// TestPriorityRegistry_DefaultFallback confirms that an unregistered job
// receives Normal priority, keeping scheduler behaviour predictable.
func TestPriorityRegistry_DefaultFallback(t *testing.T) {
	reg := jobpriority.New()
	_ = reg.Set("known-job", jobpriority.High)

	if reg.Get("unknown-job") != jobpriority.Normal {
		t.Error("expected Normal priority for unregistered job")
	}
}

// TestPriorityRegistry_AllEntriesVisible ensures that all set priorities
// are visible in the All() snapshot consumed by the scheduler at startup.
func TestPriorityRegistry_AllEntriesVisible(t *testing.T) {
	reg := jobpriority.New()
	jobs := map[string]jobpriority.Level{
		"job-a": jobpriority.Low,
		"job-b": jobpriority.Normal,
		"job-c": jobpriority.High,
	}
	for name, lvl := range jobs {
		if err := reg.Set(name, lvl); err != nil {
			t.Fatalf("Set(%s): %v", name, err)
		}
	}
	all := reg.All()
	for name, want := range jobs {
		if got, ok := all[name]; !ok || got != want {
			t.Errorf("All()[%s] = %d, want %d", name, got, want)
		}
	}
}
