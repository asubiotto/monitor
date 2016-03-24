// monitor provides monitoring of HTTP access logs through the console.
package monitor

import "time"

const (
	// Reporter reports metrics every reportInterval.
	reportInterval = 10 * time.Second

	// Reporter repots spikes in traffic that occurred in trafficWindow.
	trafficWindow = 2 * time.Minute
)
