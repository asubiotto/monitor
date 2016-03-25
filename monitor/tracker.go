package monitor

import (
	"errors"
	"strings"
	"sync"
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
	hitList  *HitList
	sections map[string]*MonitorNode
	mtx      sync.RWMutex
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

// GetTracker initializes the global tracker if it has not been initialized yet
// and returns a pointer to it in any case.
func GetTracker() *Tracker {
	trackerMtx.Lock()
	defer trackerMtx.Unlock()

	if tracker == nil {
		hl := make(HitList, 0)
		tracker = &Tracker{
			hitList:  &hl,
			sections: make(map[string]*MonitorNode),
		}
	}

	return tracker
}

// ProcessLogEntry consumes a logEntry of the form
// 	https://en.wikipedia.org/wiki/Common_Log_Format
// and updates the hits on that website section if it exists, otherwise starts
// tracking a new section.
func (t *Tracker) ProcessLogEntry(logEntry string) error {
	// HTTP access logs seem to vary quite a lot depending on the server
	// logging the accesses. The unique element across all access logs is that
	// the HTTP method is always contained and occurs after the first \" char.
	values := strings.Split(logEntry, "\"")
	if len(values) < 2 {
		print(logEntry)
		return errMalformedLogEntry
	}

	pathComponents := strings.Split(values[1], "/")

	var section string
	if len(pathComponents) > 1 {
		section = pathComponents[1]
	} else {
		print(logEntry)
		return errMalformedLogEntry
	}

	if len(pathComponents[1]) < 1 {
		// There was nothing after the first slash.
		return errMalformedLogEntry
	}

	// TODO(asubiotto): Modify this depending on what Ryan says.
	if string(pathComponents[1][0]) == " " {
		// If there was a space, our section was the root.
		section = "/"
	}

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
