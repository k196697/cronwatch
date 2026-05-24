package joblog

import (
	"strings"
	"testing"
)

func TestNew_InvalidLimit(t *testing.T) {
	_, err := New(0)
	if err == nil {
		t.Fatal("expected error for limit=0")
	}
}

func TestAppend_And_Get(t *testing.T) {
	s, _ := New(10)
	if err := s.Append("backup", "info", "started"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	entries := s.Get("backup")
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Message != "started" {
		t.Errorf("unexpected message: %s", entries[0].Message)
	}
}

func TestAppend_RespectsLimit(t *testing.T) {
	s, _ := New(3)
	for i := 0; i < 5; i++ {
		_ = s.Append("job", "info", "msg")
	}
	if got := len(s.Get("job")); got != 3 {
		t.Errorf("expected 3 entries, got %d", got)
	}
}

func TestAppend_EmptyJobReturnsError(t *testing.T) {
	s, _ := New(5)
	if err := s.Append("", "info", "hello"); err == nil {
		t.Fatal("expected error for empty job name")
	}
}

func TestAppend_EmptyLevelReturnsError(t *testing.T) {
	s, _ := New(5)
	if err := s.Append("job", "", "hello"); err == nil {
		t.Fatal("expected error for empty level")
	}
}

func TestGet_UnknownJobReturnsEmpty(t *testing.T) {
	s, _ := New(5)
	if entries := s.Get("ghost"); len(entries) != 0 {
		t.Errorf("expected empty slice, got %d entries", len(entries))
	}
}

func TestGet_ReturnsCopy(t *testing.T) {
	s, _ := New(5)
	_ = s.Append("job", "info", "original")
	got := s.Get("job")
	got[0].Message = "mutated"
	again := s.Get("job")
	if again[0].Message == "mutated" {
		t.Error("Get should return a copy, not a reference")
	}
}

func TestClear_RemovesEntries(t *testing.T) {
	s, _ := New(5)
	_ = s.Append("job", "error", "boom")
	s.Clear("job")
	if entries := s.Get("job"); len(entries) != 0 {
		t.Errorf("expected empty after clear, got %d", len(entries))
	}
}

func TestAppend_SetsTimestamp(t *testing.T) {
	s, _ := New(5)
	_ = s.Append("job", "warn", "late")
	e := s.Get("job")[0]
	if e.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}

func TestAppend_StoresLevel(t *testing.T) {
	s, _ := New(5)
	_ = s.Append("job", "error", "failed")
	e := s.Get("job")[0]
	if !strings.EqualFold(e.Level, "error") {
		t.Errorf("unexpected level: %s", e.Level)
	}
}
