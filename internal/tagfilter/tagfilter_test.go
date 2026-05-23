package tagfilter_test

import (
	"testing"

	"github.com/cronwatch/internal/tagfilter"
)

func TestNew_EmptyTagReturnsError(t *testing.T) {
	_, err := tagfilter.New([]string{"prod", ""})
	if err == nil {
		t.Fatal("expected error for empty tag, got nil")
	}
}

func TestMatch_EmptyFilterMatchesAll(t *testing.T) {
	f, err := tagfilter.New(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !f.Match([]string{"prod", "critical"}) {
		t.Error("empty filter should match any job tags")
	}
	if !f.Match(nil) {
		t.Error("empty filter should match job with no tags")
	}
}

func TestMatch_AllRequiredTagsPresent(t *testing.T) {
	f, _ := tagfilter.New([]string{"prod", "critical"})
	if !f.Match([]string{"critical", "prod", "nightly"}) {
		t.Error("expected match when all required tags are present")
	}
}

func TestMatch_MissingRequiredTag(t *testing.T) {
	f, _ := tagfilter.New([]string{"prod", "critical"})
	if f.Match([]string{"prod"}) {
		t.Error("expected no match when a required tag is missing")
	}
}

func TestMatch_NoJobTags(t *testing.T) {
	f, _ := tagfilter.New([]string{"prod"})
	if f.Match(nil) {
		t.Error("expected no match when job has no tags")
	}
}

func TestEmpty_ReportsCorrectly(t *testing.T) {
	f, _ := tagfilter.New(nil)
	if !f.Empty() {
		t.Error("filter with no tags should be empty")
	}
	f2, _ := tagfilter.New([]string{"prod"})
	if f2.Empty() {
		t.Error("filter with tags should not be empty")
	}
}

func TestTags_ReturnsCopy(t *testing.T) {
	f, _ := tagfilter.New([]string{"prod", "critical"})
	tags := f.Tags()
	if len(tags) != 2 {
		t.Fatalf("expected 2 tags, got %d", len(tags))
	}
	// Mutating the returned slice must not affect the filter.
	tags[0] = "mutated"
	for _, got := range f.Tags() {
		if got == "mutated" {
			t.Error("Tags() returned a reference to internal state")
		}
	}
}
