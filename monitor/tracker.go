// monitor provides monitoring of HTTP access logs through the console.
package monitor

import (
	"container/heap"
	"errors"
	"fmt"
	"strings"
	"sync"
)

// MonitorNode represents a section of our website. All MonitorNodes are kept
// track of in a PriorityQueue.
type MonitorNode struct {
	section string
	hits    int // How many times this section has been visited.
	index   int // Used by the priority queue.
}

// A PriorityQueue implements heap.Interface and holds MonitorNodes.
type PriorityQueue []*MonitorNode

// Tracker is a wrapper struct around a PriorityQueue providing several
// middleware methods for keeping track of HTTP access logs.
type Tracker struct {
	pq *PriorityQueue
}

var (
	errMalformedLogEntry = errors.New("malformed log entry")
	sections             = make(map[string]*MonitorNode)
	tracker              *Tracker
	trackerMtx           sync.RWMutex
)

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, hits so we use greater
	// than here.
	return pq[i].hits > pq[j].hits
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	node := x.(*MonitorNode)
	node.index = n
	*pq = append(*pq, node)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	node := old[n-1]
	node.index = -1 // for safety
	*pq = old[0 : n-1]
	return node
}

// update modifies the hits of a MonitorNode.
func (pq *PriorityQueue) update(node *MonitorNode, hits int) {
	node.hits = hits
	heap.Fix(pq, node.index)
}

// GetTracker initializes the global tracker if it has not been initialized yet
// and returns a pointer to it in any case.
func GetTracker() *Tracker {
	if tracker == nil {
		trackerMtx.Lock()
		defer trackerMtx.Unlock()

		pq := make(PriorityQueue, 0)
		heap.Init(&pq)
		tracker = &Tracker{pq: &pq}
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
	// the request line is contained within quotes and always contains the
	// HTTP method.
	values := strings.Split(logEntry, "\"")
	if len(values) != 3 {
		return errMalformedLogEntry
	}

	pathComponents := strings.Split(values[1], "/")

	var section string
	if len(pathComponents) > 1 {
		section = pathComponents[1]
	} else {
		fmt.Println(pathComponents)
		fmt.Println(logEntry)
		return errMalformedLogEntry
	}

	if section == "" {
		// We split "/" into ["", ""].
		section = "/"
	}

	// Check for existence of the MonitorNode associated with a specific section
	// and update accordingly.
	trackerMtx.Lock()
	defer trackerMtx.Unlock()
	if node, ok := sections[section]; ok {
		t.pq.update(node, node.hits+1)
	} else {
		node = &MonitorNode{
			section: section,
			hits:    1,
		}
		sections[section] = node
		heap.Push(t.pq, node)
	}

	return nil
}

// GetTopHits returns a slice of up to limit MonitorNodes with the highest
// number of hits.
func (t *Tracker) GetTopHits(limit int) []MonitorNode {
	trackerMtx.RLock()
	defer trackerMtx.RUnlock()

	// Copy priority queue.
	pq := *t.pq[:]
	result := []MonitorNode{}

	for ; limit > 0; limit-- {

	}

	if limit == 0 {
		return []*MonitorNode{}
	} else if limit >= len(*t.pq) {
		return (*t.pq)[:]
	}

	return (*t.pq)[len(*t.pq)-limit:]*/
	return result
}
