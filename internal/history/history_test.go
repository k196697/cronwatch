package history

import (
	"testing"
	"time"
)

func makeResult(job string, success bool) Result {
	return Result{
		JobName:   job,
		StartedAt: time.Now(),
		Duration:  100 * time.Millisecond,
		Success:   success,
		Output:    "output",
	}
}

func TestRecord_StoresResult(t *testing.T) {
	s := New(5)
	s.Record(makeResult("backup", true))

	r, ok := s.Latest("backup")
	if !ok {
		t.Fatal("expected a result")
	}
	if r.JobName != "backup" || !r.Success {
		t.Errorf("unexpected result: %+v", r)
	}
}

func TestLatest_MissingJob(t *testing.T) {
	s := New(5)
	_, ok := s.Latest("nonexistent")
	if ok {
		t.Error("expected false for unknown job")
	}
}

func TestRecord_RespectsLimit(t *testing.T) {
	s := New(3)
	for i := 0; i < 6; i++ {
		s.Record(makeResult("job", i%2 == 0))
	}

	all := s.All("job")
	if len(all) != 3 {
		t.Errorf("expected 3 entries, got %d", len(all))
	}
}

func TestAll_ReturnsCopy(t *testing.T) {
	s := New(5)
	s.Record(makeResult("job", true))

	all := s.All("job")
	all[0].JobName = "mutated"

	r, _ := s.Latest("job")
	if r.JobName == "mutated" {
		t.Error("All() should return a copy, not a reference")
	}
}

func TestConsecutiveFailures_Mixed(t *testing.T) {
	s := New(10)
	s.Record(makeResult("job", true))
	s.Record(makeResult("job", false))
	s.Record(makeResult("job", false))
	s.Record(makeResult("job", false))

	if n := s.ConsecutiveFailures("job"); n != 3 {
		t.Errorf("expected 3 consecutive failures, got %d", n)
	}
}

func TestConsecutiveFailures_AllSuccess(t *testing.T) {
	s := New(5)
	s.Record(makeResult("job", true))
	s.Record(makeResult("job", true))

	if n := s.ConsecutiveFailures("job"); n != 0 {
		t.Errorf("expected 0, got %d", n)
	}
}

func TestConsecutiveFailures_Empty(t *testing.T) {
	s := New(5)
	if n := s.ConsecutiveFailures("job"); n != 0 {
		t.Errorf("expected 0 for unknown job, got %d", n)
	}
}
