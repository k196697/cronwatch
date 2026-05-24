package jobqueue

import (
	"errors"
	"sync"
)

// Entry represents a pending job execution request.
type Entry struct {
	JobName  string
	Priority int
}

// Queue is a bounded, thread-safe priority job queue.
type Queue struct {
	mu      sync.Mutex
	items   []Entry
	maxSize int
}

// New creates a Queue with the given maximum capacity.
func New(maxSize int) (*Queue, error) {
	if maxSize <= 0 {
		return nil, errors.New("jobqueue: maxSize must be greater than zero")
	}
	return &Queue{maxSize: maxSize, items: make([]Entry, 0, maxSize)}, nil
}

// Push adds an entry to the queue. Returns an error if the queue is full.
func (q *Queue) Push(e Entry) error {
	if e.JobName == "" {
		return errors.New("jobqueue: job name must not be empty")
	}
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.items) >= q.maxSize {
		return errors.New("jobqueue: queue is full")
	}
	q.items = append(q.items, e)
	q.sort()
	return nil
}

// Pop removes and returns the highest-priority entry.
// Returns false if the queue is empty.
func (q *Queue) Pop() (Entry, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.items) == 0 {
		return Entry{}, false
	}
	e := q.items[0]
	q.items = q.items[1:]
	return e, true
}

// Len returns the current number of items in the queue.
func (q *Queue) Len() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.items)
}

// Drain removes and returns all entries in priority order.
func (q *Queue) Drain() []Entry {
	q.mu.Lock()
	defer q.mu.Unlock()
	out := make([]Entry, len(q.items))
	copy(out, q.items)
	q.items = q.items[:0]
	return out
}

// sort performs an insertion sort by descending priority (in-place).
func (q *Queue) sort() {
	s := q.items
	for i := 1; i < len(s); i++ {
		for j := i; j > 0 && s[j].Priority > s[j-1].Priority; j-- {
			s[j], s[j-1] = s[j-1], s[j]
		}
	}
}
