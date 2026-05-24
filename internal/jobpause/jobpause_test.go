package jobpause

import (
	"testing"
	"time"
)

func TestPause_And_IsPaused(t *testing.T) {
	s := New()
	if err := s.Pause("backup", 0); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !s.IsPaused("backup") {
		t.Error("expected job to be paused")
	}
}

func TestResume_ClearsPause(t *testing.T) {
	s := New()
	_ = s.Pause("backup", 0)
	if err := s.Resume("backup"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.IsPaused("backup") {
		t.Error("expected job to no longer be paused after resume")
	}
}

func TestIsPaused_UnknownJob(t *testing.T) {
	s := New()
	if s.IsPaused("nonexistent") {
		t.Error("expected unknown job to not be paused")
	}
}

func TestPause_TimedExpiry(t *testing.T) {
	s := New()
	_ = s.Pause("cleanup", 10*time.Millisecond)
	if !s.IsPaused("cleanup") {
		t.Fatal("expected job to be paused immediately after Pause")
	}
	time.Sleep(20 * time.Millisecond)
	if s.IsPaused("cleanup") {
		t.Error("expected timed pause to have expired")
	}
}

func TestPause_EmptyNameReturnsError(t *testing.T) {
	s := New()
	if err := s.Pause("", 0); err == nil {
		t.Error("expected error for empty job name")
	}
}

func TestResume_EmptyNameReturnsError(t *testing.T) {
	s := New()
	if err := s.Resume(""); err == nil {
		t.Error("expected error for empty job name")
	}
}

func TestAll_ReturnsCopy(t *testing.T) {
	s := New()
	_ = s.Pause("job1", 0)
	_ = s.Pause("job2", 0)

	all := s.All()
	if len(all) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(all))
	}
	// Mutating the copy must not affect the store.
	delete(all, "job1")
	if !s.IsPaused("job1") {
		t.Error("store should not be affected by mutating the All() copy")
	}
}
