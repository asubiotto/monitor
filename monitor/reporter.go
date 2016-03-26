package monitor

import (
	"fmt"
	"log"
	"os"
	"text/tabwriter"
	"time"
)

// reportMetrics starts reporting metrics about section hits and other metrics.
func reportMetrics(finishReporting chan struct{}) {
	tracker := GetTracker()
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 0, '\t', 0)
	for {
		select {
		case <-time.After(reportInterval):
			fmt.Fprintln(w, "-")
			fmt.Fprintln(w, "Section Hits")
			fmt.Fprintln(w, "-")
			for _, section := range tracker.GetTopHits(reportLimit) {
				fmt.Fprintln(
					w,
					fmt.Sprintf("%s\t%d", section.section, section.hits),
				)
			}
			fmt.Fprintln(w, "-")
			fmt.Fprintln(w, "Other Metrics")
			fmt.Fprintln(w, "-")
			fmt.Fprintln(w, fmt.Sprintf("%s\t%.2f", "rps", tracker.GetRPS()))
			fmt.Fprintln(
				w,
				fmt.Sprintf("%s\t%d", "tot sections", tracker.GetNumSections()),
			)
			fmt.Fprintln(
				w,
				fmt.Sprintf("%s\t%d", "tot hits", tracker.GetTotalTraffic()),
			)
			w.Flush()
			fmt.Println()
		case <-finishReporting:
			return
		}
	}
}

// reportTrafficSpike reports a traffic spike and the total number of hits in
// the last two minutes.
func reportTrafficSpike(hits int) {
	log.Println("High traffic generated an alert - hits =", hits)
}

// reportTrafficUnspike reports that traffic on average has fallen below the
// threshold.
func reportTrafficUnspike(hits int) {
	log.Println("High traffic alert recovered - hits =", hits)
}
