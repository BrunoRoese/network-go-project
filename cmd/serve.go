package cmd

import (
	"github.com/BrunoRoese/socket/pkg/server"
	"github.com/BrunoRoese/socket/pkg/server/service"
	"github.com/spf13/cobra"
	"log/slog"
	"os"
)

var (
	ip string

	udpServer *server.Server

	serveCmd = &cobra.Command{
		Use:   "serve",
		Short: "Start the UDP server",
		Long:  `This command starts the UDP server that listens for incoming connections`,
		Run:   run,
	}
)

func run(cmd *cobra.Command, args []string) {
	shutdown := make(chan os.Signal, 1)

	if ip == "" {
		slog.Error("IP address not provided")
		os.Exit(1)
	}

	startUpServer(ip, 8080)

	slog.Info("Broadcasting to all IPs")

	service.Broadcast()

	<-shutdown

	udpServer.Close()

	slog.Info("Shutting down server")
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

func init() {
	serveCmd.PersistentFlags().StringVar(&ip, "ip", "", "Your network IP address")
	err := serveCmd.MarkPersistentFlagRequired("ip")

	rootCmd.AddCommand(serveCmd)

	if err != nil {
		slog.Error("Error marking flag as required", slog.String("error", err.Error()))
		os.Exit(1)
	}
}
