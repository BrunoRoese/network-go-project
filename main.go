package main

import (
	"errors"
	"flag"
	"github.com/BrunoRoese/socket/config"
	"github.com/BrunoRoese/socket/server"
	"log/slog"
	"os"
)

var udpServer *server.Server

func main() {
	config.SetupLogger()

	shutdown := make(chan os.Signal, 1)

	ip, port, err := handleFlags()

	if err != nil {
		slog.Error("Error handling flags", slog.String("error", err.Error()))
		return
	}

	startUpServer(ip, port)

	server.Broadcast()

	<-shutdown

	udpServer.Close()

	slog.Info("Shutting down server")
}

func handleFlags() (ip string, port int, err error) {
	flag.StringVar(&ip, "ip", "", "IP address to bind connection")

	flag.Parse()

	if ip == "" {
		return "", 0, errors.New("IP address not provided")
	}

	return ip, 8080, nil
}

func startUpServer(ip string, port int) {
	udpServer, err := server.Init(ip, port)

	if err != nil {
		slog.Error("Error starting server, stopping application", slog.String("error", err.Error()))
		return
	}

	udpServer.StartListeningRoutine()

	slog.Info("Server started")
}
