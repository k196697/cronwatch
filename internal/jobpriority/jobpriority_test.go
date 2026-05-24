package jobpriority_test

import (
	"testing"

	"cronwatch/internal/jobpriority"
)

func TestSet_And_Get(t *testing.T) {
	r := jobpriority.New()
	if err := r.Set("backup", jobpriority.High); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := r.Get("backup"); got != jobpriority.High {
		t.Errorf("expected High(%d), got %d", jobpriority.High, got)
	}
}

func TestGet_DefaultsToNormal(t *testing.T) {
	r := jobpriority.New()
	if got := r.Get("unknown"); got != jobpriority.Normal {
		t.Errorf("expected Normal(%d), got %d", jobpriority.Normal, got)
	}
}

func TestSet_EmptyJobReturnsError(t *testing.T) {
	r := jobpriority.New()
	if err := r.Set("", jobpriority.Low); err == nil {
		t.Fatal("expected error for empty job name")
	}
}

func TestSet_InvalidLevelReturnsError(t *testing.T) {
	r := jobpriority.New()
	for _, lvl := range []jobpriority.Level{0, 11, -1} {
		if err := r.Set("job", lvl); err == nil {
			t.Errorf("expected error for level %d", lvl)
		}
	}
}

func TestDelete_RemovesEntry(t *testing.T) {
	r := jobpriority.New()
	_ = r.Set("cleanup", jobpriority.Low)
	r.Delete("cleanup")
	if got := r.Get("cleanup"); got != jobpriority.Normal {
		t.Errorf("expected Normal after delete, got %d", got)
	}
}

func TestAll_ReturnsCopy(t *testing.T) {
	r := jobpriority.New()
	_ = r.Set("a", jobpriority.Low)
	_ = r.Set("b", jobpriority.High)
	all := r.All()
	if len(all) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(all))
	}
	// Mutating the copy must not affect the registry.
	all["a"] = jobpriority.High
	if r.Get("a") != jobpriority.Low {
		t.Error("mutation of All() snapshot affected registry")
	}
}

func TestSet_Overwrites(t *testing.T) {
	r := jobpriority.New()
	_ = r.Set("sync", jobpriority.Low)
	_ = r.Set("sync", jobpriority.High)
	if got := r.Get("sync"); got != jobpriority.High {
		t.Errorf("expected High after overwrite, got %d", got)
	}
}
