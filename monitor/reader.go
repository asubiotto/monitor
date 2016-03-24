// monitor provides monitoring of HTTP access logs through the console.
package monitor

import (
	"fmt"

	"github.com/hpcloud/tail"
)

func readStream(tail *tail.Tail, finishReading chan struct{}) {
	for {
		select {
		case line := <-tail.Lines:
			fmt.Println("Got line", line.Text)
		case <-finishReading:
			tail.Stop()
			return
		}
	}
}
