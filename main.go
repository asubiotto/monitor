package main

import (
	"flag"
	"log"

	"github.com/asubiotto/monitor/monitor"
)

var (
	path             = flag.String("path", "", "path to HTTP access logs file")
	trafficThreshold = flag.Int("threshold", 10, "threshold used when alerting traffic spikes")
)

// Given a path and an optional trafficThreshold, starts up a monitor to monitor
// HTTP access logs. Example run:
// 	go run main.go -path=to/file -threshold=10
func main() {
	flag.Parse()
	if *path == "" {
		log.Fatal("Please specify a path to the HTTP access logs file")
	}
	monitor.MonitorHTTPAccessLogs(*path, *trafficThreshold)
}
