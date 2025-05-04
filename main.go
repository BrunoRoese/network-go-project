package main

import (
	"errors"
	"flag"
	"github.com/BrunoRoese/socket/cmd"
	"github.com/BrunoRoese/socket/server"
	"log/slog"
	"os"
)

var udpServer *server.Server

func main() {
	cmd.Execute()

	shutdown := make(chan os.Signal, 1)

	ip, port, err := handleFlags()

	if err != nil {
		slog.Error("Error handling flags", slog.String("error", err.Error()))
		return
	}

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
