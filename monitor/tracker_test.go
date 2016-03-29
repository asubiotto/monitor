package monitor

import (
	"testing"
	"time"
)

func TestTrackerOrder(t *testing.T) {
	InitTracker(3)
	tracker.upsertSection("first")
	tracker.upsertSection("second")
	hits := tracker.GetTopHits(10)
	if len(hits) != 2 {
		t.Fatalf("Expected 2 hits got %d\n", len(hits))
	}

	if hits[0].section != "first" {
		t.Errorf("Expected %q, got %s\n", "first", hits[0].section)
	}

	if hits[1].section != "second" {
		t.Errorf("Expected %q, got %s\n", "second", hits[1].section)
	}

	if hits[0].index != 0 || hits[1].index != 1 {
		t.Errorf("Wrong index order\n")
	}

	// Move second before first.
	tracker.upsertSection("second")
	hits = tracker.GetTopHits(10)

	if len(hits) != 2 {
		t.Fatalf("Expected 2 hits got %d\n", len(hits))
	}

	if hits[0].section != "second" {
		t.Errorf("Expected %q, got %s\n", "second", hits[0].section)
	}

	if hits[1].section != "first" {
		t.Errorf("Expected %q, got %s\n", "first", hits[1].section)
	}

	if hits[0].index != 0 || hits[1].index != 1 {
		t.Errorf("Wrong index order\n")
	}
}

func TestTrafficSpike(t *testing.T) {
	threshold := 10
	InitTracker(threshold)
	// Use a custom tWindow so this test doesn't have to run for minutes.
	tWindow := time.Second * 2

	// Insert a section so we don't panic when we calculate average traffic.
	// Inserting through the hl bypasses the increment traffic logic of the
	// tracker.
	tracker.hitList.InsertSection("/")

	for i := 0; i < threshold; i++ {
		tracker.incrementTraffic(tWindow)
		if tracker.thresholdExceeded {
			t.Errorf("Threshold exceeded early\n")
		}
	}
	tracker.incrementTraffic(tWindow)
	if !tracker.thresholdExceeded {
		t.Errorf("Threshold should be exceeded\n")
	}

	if tracker.traffic != threshold+1 {
		t.Errorf("Expected traffic to be %d, got %d\n", threshold+1, tracker.traffic)
	}

	<-time.After(tWindow + time.Second)

	if tracker.traffic != 0 {
		t.Errorf("Expected traffic in trafficWindow to have gone back to 0\n")
	}

	if tracker.totalTraffic != threshold+1 {
		t.Errorf("Expected traffic to be %d, got %d\n", threshold+1, tracker.totalTraffic)
	}
}
