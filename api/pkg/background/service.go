package background

import (
	"fmt"
	"time"
)

func backgroundService(quitChan chan struct{}) {
	for {
		select {
		case <-quitChan:
			fmt.Println("Background service received a quit signal. Cleaning up...")
			return
		default:
			time.Sleep(5 * time.Second)
		}
	}
}
