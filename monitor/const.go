// monitor provides monitoring of HTTP access logs through the console.
package monitor

import "time"

const (
	// Reporter reports metrics every reportInterval.
	reportInterval = 10 * time.Second

	// Reporter shows only the top reportLimit sections every reportInterval.
	reportLimit = 10

	// Reporter repots spikes in traffic that occurred in trafficWindow.
	trafficWindow = 2 * time.Minute
)
