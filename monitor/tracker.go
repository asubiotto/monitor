package monitor

import (
	"errors"
	"strings"
	"sync"
	"time"
)

// MonitorNode represents a section of our website. All MonitorNodes are kept
// track of in a HitList.
type MonitorNode struct {
	section string
	hits    int // How many times this section has been visited.
	index   int // Used by the priority queue.
}

// HitList is an array of sorted MonitorNodes. We do not use a priority queue
// because it is easier to get the n first elements from a sorted array. To
// get n elements at read time from the priority queue one must remove each
// element then reinsert after all elements have been read to produce the
// correct order.
type HitList []*MonitorNode

// Tracker is a wrapper struct around a HitList providing several
// middleware methods for keeping track of HTTP access logs.
type Tracker struct {
	hitList           *HitList
	sections          map[string]*MonitorNode
	traffic           int  // Total traffic in last trafficWindow.
	totalTraffic      int  // Total traffic since startTime.
	threshold         int  // Traffic threshold for traffic spikes.
	thresholdExceeded bool // Useful for reporting.
	startTime         time.Time
	mtx               sync.RWMutex
}

var (
	errMalformedLogEntry = errors.New("malformed log entry")

	tracker    *Tracker
	trackerMtx sync.RWMutex
)

// InsertSection inserts a new node into the HitList with the specified section
// and a hit counter of 1 and returns it.
func (hl *HitList) InsertSection(section string) *MonitorNode {
	node := &MonitorNode{
		section: section,
		index:   len(*hl),
	}
	*hl = append(*hl, node)
	hl.Update(node, 1) // Update takes the number of hits.
	return node
}

// Update modifies the hits of a MonitorNode and swaps until the list is sorted.
func (hl *HitList) Update(node *MonitorNode, hits int) {
	node.hits = hits
	curIndex := node.index

	list := *hl
	for {
		if curIndex == 0 || list[curIndex-1].hits >= list[curIndex].hits {
			break
		}

		list[curIndex-1], list[curIndex] = list[curIndex], list[curIndex-1]
		list[curIndex].index = curIndex
		list[curIndex-1].index = curIndex - 1
	}
}

// InitTracker initializes the global tracker with a specified traffic
// threshold. Not threadsafe.
func InitTracker(threshold int) {
	hl := make(HitList, 0)
	tracker = &Tracker{
		hitList:   &hl,
		sections:  make(map[string]*MonitorNode),
		threshold: threshold,
		startTime: time.Now(),
	}
}

// GetTracker initializes the global tracker if it has not been initialized yet
// and returns a pointer to it in any case.
func GetTracker() *Tracker {
	trackerMtx.Lock()
	defer trackerMtx.Unlock()

	if tracker == nil {
		InitTracker(trafficThreshold)
	}

	return tracker
}

// ProcessLogEntry consumes a logEntry of the form
// 	https://en.wikipedia.org/wiki/Common_Log_Format
// and updates the hits on that website section if it exists, otherwise starts
// tracking a new section.
func (t *Tracker) ProcessLogEntry(logEntry string) error {
	// Request line starts with "\"".
	values := strings.Split(logEntry, "\"")
	if len(values) < 2 {
		return errMalformedLogEntry
	}

	requestLine := values[1]

	// Find the first "/" after the start of request line.
	firstSlashIdx := 0
	for i, s := range requestLine {
		if s == '/' {
			firstSlashIdx = i
			break
		}
	}

	if firstSlashIdx == 0 || firstSlashIdx == len(requestLine)-1 {
		return errMalformedLogEntry
	}

	// Iterate until you hit a "\"" (end of request line), " " (end of url in
	// request line), or a "/", which will be the second "/".
	sectionEnd := firstSlashIdx
	for i, s := range requestLine[firstSlashIdx+1:] {
		if s == '"' || s == ' ' || s == '/' {
			sectionEnd = i + firstSlashIdx + 1
			break
		}
	}

	section := requestLine[firstSlashIdx:sectionEnd]
	t.upsertSection(section)

	return nil
}

// upsertSection adds a hit to the section if it exists or inserts it to start
// keeping track of it.
func (t *Tracker) upsertSection(section string) {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	if node, ok := t.sections[section]; ok {
		t.hitList.Update(node, node.hits+1)
	} else {
		t.sections[section] = t.hitList.InsertSection(section)
	}

	t.incrementTraffic(trafficWindow)
}

// incrementTraffic adds to total traffic checking for trafficThreshold
// violations and starts a timer to decrement traffic after tWindow. t.mtx
// should be held when calling this function.
func (t *Tracker) incrementTraffic(tWindow time.Duration) {
	t.traffic++
	t.totalTraffic++

	if ((t.traffic / len(*t.hitList)) > t.threshold) && !t.thresholdExceeded {
		// This traffic increment caused our threshold to be exceeded, therefore
		// report.
		t.thresholdExceeded = true
		reportTrafficSpike(t.traffic)
	}

	go func() {
		<-time.After(tWindow)
		t.mtx.Lock()
		defer t.mtx.Unlock()
		t.traffic--
		if ((t.traffic / len(*t.hitList)) <= t.threshold) && t.thresholdExceeded {
			// This traffic decrement caused our average traffic over all
			// sections to fall below our threshold.
			t.thresholdExceeded = false
			reportTrafficUnspike(t.traffic)
		}
	}()
}

// GetTopHits returns a slice of up to limit MonitorNodes with the highest
// number of hits.
func (t *Tracker) GetTopHits(limit int) []MonitorNode {
	t.mtx.RLock()
	defer t.mtx.RUnlock()

	if limit > len(*t.hitList) {
		limit = len(*t.hitList)
	}

	// Don't return pointers for safety.
	result := make([]MonitorNode, limit)
	for i := 0; i < limit; i++ {
		result[i] = *((*t.hitList)[i])
	}

	return result
}

// GetNumSections returns the total number of sections.
func (t *Tracker) GetNumSections() int {
	t.mtx.RLock()
	defer t.mtx.RUnlock()

	return len(*t.hitList)
}

// GetTotalTraffic returns the total number of hits.
func (t *Tracker) GetTotalTraffic() int {
	t.mtx.RLock()
	defer t.mtx.RUnlock()

	return t.totalTraffic
}

// GetRPS returns the requests per second in trafficWindow.
func (t *Tracker) GetRPS() float64 {
	t.mtx.RLock()
	defer t.mtx.RUnlock()

	return float64(t.totalTraffic) / time.Since(t.startTime).Seconds()
}
