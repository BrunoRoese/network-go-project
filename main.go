package main

import (
	"github.com/BrunoRoese/socket/config"
	"github.com/BrunoRoese/socket/server"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	config.SetupLogger()

	shutdown := make(chan os.Signal, 1)
	udpServer, err := server.Init("127.0.0.1", 8080)

	signal.Notify(shutdown, syscall.SIGTERM, syscall.SIGINT)

	if err != nil {
		slog.Error("Error starting server, stopping application", slog.String("error", err.Error()))
		return
	}

	udpServer.StartListeningRoutine()

	<-shutdown

	udpServer.Close()

	slog.Info("Shutting down server")
}
