package monitor

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"
)

func reportMetrics(finishReporting chan struct{}) {
	tracker := GetTracker()
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 0, '\t', 0)
	for {
		select {
		case <-time.After(reportInterval):
			fmt.Fprintln(w, "section\thits")
			for _, section := range tracker.GetTopHits(reportLimit) {
				fmt.Fprintln(
					w,
					fmt.Sprintf("%s\t%d", section.section, section.hits),
				)
			}
			w.Flush()
			fmt.Println()
		case <-finishReporting:
			return
		}
	}
}
