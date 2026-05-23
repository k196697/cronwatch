package joblock

import (
	"sync"
	"testing"
)

func TestTryAcquire_SucceedsWhenFree(t *testing.T) {
	l := New()
	ok, release, err := l.TryAcquire("backup")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected lock to be acquired")
	}
	if release == nil {
		t.Fatal("expected non-nil release func")
	}
	release()
}

func TestTryAcquire_FailsWhenHeld(t *testing.T) {
	l := New()
	ok, release, _ := l.TryAcquire("backup")
	if !ok {
		t.Fatal("first acquire should succeed")
	}
	defer release()

	ok2, release2, err := l.TryAcquire("backup")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok2 {
		t.Fatal("second acquire should fail while lock is held")
	}
	if release2 != nil {
		t.Fatal("release func should be nil on failure")
	}
}

func TestTryAcquire_SucceedsAfterRelease(t *testing.T) {
	l := New()
	_, release, _ := l.TryAcquire("sync")
	release()

	ok, release2, _ := l.TryAcquire("sync")
	if !ok {
		t.Fatal("expected lock to be re-acquired after release")
	}
	release2()
}

func TestTryAcquire_IndependentJobs(t *testing.T) {
	l := New()
	ok1, r1, _ := l.TryAcquire("job-a")
	ok2, r2, _ := l.TryAcquire("job-b")
	defer r1()
	defer r2()
	if !ok1 || !ok2 {
		t.Fatal("independent jobs should not block each other")
	}
}

func TestTryAcquire_EmptyNameReturnsError(t *testing.T) {
	l := New()
	_, _, err := l.TryAcquire("")
	if err == nil {
		t.Fatal("expected error for empty job name")
	}
}

func TestIsLocked_ReflectsState(t *testing.T) {
	l := New()
	if l.IsLocked("x") {
		t.Fatal("unknown job should not be locked")
	}
	_, release, _ := l.TryAcquire("x")
	if !l.IsLocked("x") {
		t.Fatal("job should be locked after acquire")
	}
	release()
	if l.IsLocked("x") {
		t.Fatal("job should not be locked after release")
	}
}

func TestTryAcquire_ConcurrentSafety(t *testing.T) {
	l := New()
	var wg sync.WaitGroup
	acquired := 0
	var mu sync.Mutex

	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ok, release, _ := l.TryAcquire("shared")
			if ok {
				mu.Lock()
				acquired++
				mu.Unlock()
				release()
			}
		}()
	}
	wg.Wait()
	if acquired == 0 {
		t.Fatal("expected at least one goroutine to acquire the lock")
	}
}
