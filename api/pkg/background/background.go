package background

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	db "github.com/omarbdrn/simple-api/pkg/database"
	"github.com/omarbdrn/simple-api/pkg/rest"
	"github.com/omarbdrn/simple-api/pkg/server"
)

func RunService() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
		return
	}

	quitChan := make(chan struct{})
	if err := db.InitDB(); err != nil {
		log.Fatalf("Error initializing the database: %v", err)
	}

	mqServer := server.NewServer()

	go backgroundService(quitChan)
	go mqServer.ConnectMQ("api", "api", "simpleagent")
	go rest.RunRestServer(mqServer)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	_ = <-sigChan

	close(quitChan)
	defer mqServer.DisconnectMQ()

	connection, err := db.GetDB().DB()
	if err != nil {
	} else {
		defer connection.Close()
	}
}
