package jobqueue

import (
	"testing"
)

func TestNew_InvalidMaxSize(t *testing.T) {
	_, err := New(0)
	if err == nil {
		t.Fatal("expected error for maxSize=0")
	}
}

func TestPush_And_Pop_Order(t *testing.T) {
	q, _ := New(10)
	_ = q.Push(Entry{JobName: "low", Priority: 1})
	_ = q.Push(Entry{JobName: "high", Priority: 10})
	_ = q.Push(Entry{JobName: "mid", Priority: 5})

	e, ok := q.Pop()
	if !ok || e.JobName != "high" {
		t.Fatalf("expected high-priority job first, got %+v", e)
	}
	e, _ = q.Pop()
	if e.JobName != "mid" {
		t.Fatalf("expected mid-priority job second, got %+v", e)
	}
}

func TestPop_EmptyQueue(t *testing.T) {
	q, _ := New(5)
	_, ok := q.Pop()
	if ok {
		t.Fatal("expected false on empty queue")
	}
}

func TestPush_RejectsFullQueue(t *testing.T) {
	q, _ := New(2)
	_ = q.Push(Entry{JobName: "a", Priority: 1})
	_ = q.Push(Entry{JobName: "b", Priority: 2})
	err := q.Push(Entry{JobName: "c", Priority: 3})
	if err == nil {
		t.Fatal("expected error when queue is full")
	}
}

func TestPush_RejectsEmptyJobName(t *testing.T) {
	q, _ := New(5)
	err := q.Push(Entry{JobName: "", Priority: 1})
	if err == nil {
		t.Fatal("expected error for empty job name")
	}
}

func TestLen_TracksSize(t *testing.T) {
	q, _ := New(10)
	if q.Len() != 0 {
		t.Fatal("expected initial length 0")
	}
	_ = q.Push(Entry{JobName: "x", Priority: 1})
	if q.Len() != 1 {
		t.Fatalf("expected length 1, got %d", q.Len())
	}
}

func TestDrain_ReturnsAllAndClears(t *testing.T) {
	q, _ := New(10)
	_ = q.Push(Entry{JobName: "a", Priority: 3})
	_ = q.Push(Entry{JobName: "b", Priority: 1})
	_ = q.Push(Entry{JobName: "c", Priority: 2})

	out := q.Drain()
	if len(out) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(out))
	}
	if out[0].JobName != "a" {
		t.Fatalf("expected highest priority first, got %s", out[0].JobName)
	}
	if q.Len() != 0 {
		t.Fatal("expected queue to be empty after drain")
	}
}
