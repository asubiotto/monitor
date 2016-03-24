// monitor provides monitoring of HTTP access logs through the console.
package monitor

import (
	"fmt"
	"time"
)

func reportMetrics(trafficThreshold int, finishReporting chan struct{}) {
	tracker := GetTracker()
	for {
		select {
		case <-time.After(reportInterval):
			fmt.Println("Reporting")
			for _, section := range tracker.GetTopHits(reportLimit) {
				fmt.Println(section.section, "has", section.hits, "hits")
			}
		case <-finishReporting:
			return
		}
	}
}
