// monitor provides monitoring of HTTP access logs through the console.
package monitor

import (
	"log"
	"os"
	"os/signal"

	"github.com/hpcloud/tail"
)

// MonitorHTTPAccessLogs consumes a path to a file of HTTP access logs as well
// as a trafficThreshold that indicates a minimum amount of traffic that must
// occur in a window of time (trafficWindow) for a traffic spike to be reported.
// It spins off a reader that tails the input file as well as a reporter that
// reports key data every reportInterval and traffic spikes.
func MonitorHTTPAccessLogs(path string, trafficThreshold int) {
	log.Println("Starting monitor")

	// Tail file to take care of any potential errors before we spin up
	// goroutines.
	tail, err := tail.TailFile(path, tail.Config{Follow: true})
	if err != nil {
		log.Fatal(err)
	}
	defer tail.Cleanup()

	// Initialize the tracker.
	InitTracker(trafficThreshold)

	// Channels to indicate to the reader and reporter that they should clean
	// up and exit.
	finishReading := make(chan struct{})
	finishReporting := make(chan struct{})

	// Channel to catch SIGINTs. Capacity of 1 to block writers.
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)

	go readStream(tail, finishReading)
	go reportMetrics(finishReporting)

	log.Println("Ctrl-C to quit")
	for _ = range sigint {
		log.Println("Caught signal, shutting down...")
		finishReading <- struct{}{}
		finishReporting <- struct{}{}
		break
	}

	log.Println("Done monitoring")
}
