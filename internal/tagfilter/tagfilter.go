// Package tagfilter provides job filtering by tag labels.
package tagfilter

import "fmt"

// Filter holds a set of required tags and matches jobs against them.
type Filter struct {
	tags map[string]struct{}
}

// New creates a Filter that matches jobs containing all of the given tags.
// Returns an error if any tag is empty.
func New(tags []string) (*Filter, error) {
	set := make(map[string]struct{}, len(tags))
	for _, t := range tags {
		if t == "" {
			return nil, fmt.Errorf("tagfilter: tag must not be empty")
		}
		set[t] = struct{}{}
	}
	return &Filter{tags: set}, nil
}

// Match reports whether jobTags satisfies the filter.
// An empty filter matches every job.
func (f *Filter) Match(jobTags []string) bool {
	if len(f.tags) == 0 {
		return true
	}
	present := make(map[string]struct{}, len(jobTags))
	for _, t := range jobTags {
		present[t] = struct{}{}
	}
	for required := range f.tags {
		if _, ok := present[required]; !ok {
			return false
		}
	}
	return true
}

// Tags returns the set of required tags as a slice.
func (f *Filter) Tags() []string {
	out := make([]string, 0, len(f.tags))
	for t := range f.tags {
		out = append(out, t)
	}
	return out
}

// Empty reports whether the filter has no required tags.
func (f *Filter) Empty() bool {
	return len(f.tags) == 0
}
