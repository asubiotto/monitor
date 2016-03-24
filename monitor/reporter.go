// monitor provides monitoring of HTTP access logs through the console.
package monitor

import (
	"fmt"
	"time"
)

func reportMetrics(trafficThreshold int, finishReporting chan struct{}) {
	for {
		select {
		case <-time.After(reportInterval):
			fmt.Println("Reporting")
		case <-finishReporting:
			return
		}
	}
}
