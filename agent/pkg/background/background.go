package background

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/omarbdrn/simple_agent/internal/radio"
	"github.com/omarbdrn/simple_agent/pkg/server"
)

func RunService() {
	quitChan := make(chan struct{})

	mqServer := server.NewListener()

	go backgroundService(quitChan)
	go mqServer.ConnectMQ("agent", "agent", "simpleagent")
	go radio.StartRadio()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	_ = <-sigChan

	close(quitChan)
	defer mqServer.DisconnectMQ()
}
