package cmd

import (
	"github.com/BrunoRoese/socket/pkg/network"
	"github.com/BrunoRoese/socket/pkg/protocol/service"
	"github.com/BrunoRoese/socket/pkg/server"
	"github.com/spf13/cobra"
	"log/slog"
	"os"
)

var (
	serveCmd = &cobra.Command{
		Use:   "serve",
		Short: "Start the UDP server",
		Long:  `This command starts the UDP server that listens for incoming connections`,
		Run:   run,
	}
)

func run(cmd *cobra.Command, args []string) {
	shutdown := make(chan os.Signal, 1)
	ip, err := network.GetLocalIp()

	if err != nil {
		slog.Error("IP address not provided")
		os.Exit(1)
	}

	udpServer := startUpServer(ip)

	if udpServer == nil {
		slog.Error("Error starting server, stopping application")
		os.Exit(1)
	}

	slog.Info("Broadcasting to all IPs")

	service.Broadcast()

	<-shutdown

	udpServer.Close()

	slog.Info("Shutting down server")
}

func startUpServer(ip string) *server.Service {
	udpServer, err := server.Init(ip)

	if err != nil {
		slog.Error("Error starting server, stopping application", slog.String("error", err.Error()))
		return nil
	}

	slog.Info("ServerService started")

	return udpServer
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
