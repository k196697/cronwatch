package jobstore_test

import (
	"testing"
	"time"

	"github.com/cronwatch/internal/jobstore"
)

func makeEntry(name string, state jobstore.State) jobstore.Entry {
	return jobstore.Entry{
		JobName:   name,
		State:     state,
		LastStart: time.Now(),
		PID:       1234,
	}
}

func TestSet_And_Get(t *testing.T) {
	s := jobstore.New()
	e := makeEntry("backup", jobstore.StateRunning)
	s.Set(e)

	got, ok := s.Get("backup")
	if !ok {
		t.Fatal("expected entry to exist")
	}
	if got.State != jobstore.StateRunning {
		t.Errorf("expected StateRunning, got %v", got.State)
	}
}

func TestGet_Missing(t *testing.T) {
	s := jobstore.New()
	_, ok := s.Get("nonexistent")
	if ok {
		t.Error("expected missing entry to return false")
	}
}

func TestSet_Overwrites(t *testing.T) {
	s := jobstore.New()
	s.Set(makeEntry("job", jobstore.StateRunning))
	s.Set(makeEntry("job", jobstore.StateSuccess))

	got, _ := s.Get("job")
	if got.State != jobstore.StateSuccess {
		t.Errorf("expected StateSuccess after overwrite, got %v", got.State)
	}
}

func TestAll_ReturnsCopy(t *testing.T) {
	s := jobstore.New()
	s.Set(makeEntry("a", jobstore.StateSuccess))
	s.Set(makeEntry("b", jobstore.StateFailed))

	all := s.All()
	if len(all) != 2 {
		t.Errorf("expected 2 entries, got %d", len(all))
	}

	// Mutating the returned slice must not affect the store.
	all[0].State = jobstore.StateTimedOut
	got, _ := s.Get(all[0].JobName)
	if got.State == jobstore.StateTimedOut {
		t.Error("All() returned a reference, not a copy")
	}
}

func TestDelete_RemovesEntry(t *testing.T) {
	s := jobstore.New()
	s.Set(makeEntry("cleanup", jobstore.StateSuccess))
	s.Delete("cleanup")

	_, ok := s.Get("cleanup")
	if ok {
		t.Error("expected entry to be deleted")
	}
}

func TestDelete_NonExistentIsNoop(t *testing.T) {
	s := jobstore.New()
	// Should not panic.
	s.Delete("ghost")
}
