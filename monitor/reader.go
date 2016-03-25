package monitor

import (
	"log"

	"github.com/hpcloud/tail"
)

// readStream tails a file of HTTP access logs, parsing them and storing them
// to be tracked until it receives a message on finishReading.
func readStream(tail *tail.Tail, finishReading chan struct{}) {
	tracker := GetTracker()
	for {
		select {
		case line := <-tail.Lines:
			// Lines expected to follow this form:
			// 	https://en.wikipedia.org/wiki/Common_Log_Format
			if err := tracker.ProcessLogEntry(line.Text); err != nil {
				// Notify if there is an error, but keep on reading other log
				// entries.
				log.Println(err)
			}
		case <-finishReading:
			tail.Stop()
			return
		}
	}
}
