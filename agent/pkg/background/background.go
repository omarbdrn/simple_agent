package background

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/omarbdrn/simple_agent/internal/radio"
)

func RunService() {
	quitChan := make(chan struct{})

	go backgroundService(quitChan)
	go radio.StartRadio()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	_ = <-sigChan

	close(quitChan)
}
