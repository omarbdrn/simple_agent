package background

import (
	"fmt"
	// "runtime"
	// "time"
)

func backgroundService(quitChan chan struct{}) {
	for {
		select {
		case <-quitChan:
			fmt.Println("Background service received a quit signal. Cleaning up...")
			return
		default:
			// runtime.GC() // Forcing Garbage Collection
			// time.Sleep(5 * time.Second)
		}
	}
}
