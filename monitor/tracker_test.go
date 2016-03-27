package monitor

import "testing"

func TestTrackerOrder(t *testing.T) {
	InitTracker()
	tracker.upsertSection("first")
	tracker.upsertSection("second")
	hits := GetTopHits()
	if len(hits) != 2 {
		t.Fatalf("Expected 2 hits got %d\n", len(hits))
	}

	if hits[0].section != "first" {
		t.Errorf("Expected %q, got %s", "first", hits[0].section)
	}

	if hits[1].section != "second" {
		t.Errorf("Expected %q, got %s", "second", hits[1].section)
	}

	// Move second before first.
	tracker.upsertSection("second")
	if hits[0].section != "second" {
		t.Errorf("Expected %q, got %s", "second", hits[0].section)
	}

	if hits[1].section != "first" {
		t.Errorf("Expected %q, got %s", "first", hits[1].section)
	}
}
