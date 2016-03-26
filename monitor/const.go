// monitor provides monitoring of HTTP access logs through the console.
package monitor

import "time"

const (
	// How often top section hits should be reported.
	reportInterval = 10 * time.Second

	// How many top section hits should be shown.
	reportLimit = 10

	// If total traffic averaged over the number of sections goes over
	// trafficThreshold, report a traffic spike.
	trafficWindow    = 10 * time.Second
	trafficThreshold = 10
)
