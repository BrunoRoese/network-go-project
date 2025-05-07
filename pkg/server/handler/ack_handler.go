package handler

import (
	"github.com/BrunoRoese/socket/pkg/client"
	"github.com/BrunoRoese/socket/pkg/protocol"
	"github.com/BrunoRoese/socket/pkg/protocol/parser"
	"github.com/BrunoRoese/socket/pkg/server"
	"github.com/google/uuid"
	"log/slog"
)

func HandleAckReq(request *protocol.Request) {
	if request.Information.Id == uuid.Nil {
		slog.Info("Broadcast ACK received", slog.String("source", request.Information.Source), slog.String("id", "0"))

		ip, _, err := parser.ParseSource(request.Information.Source)

		if err != nil {
			slog.Error("Error parsing source", slog.String("error", err.Error()))
			return
		}

		handleCounter(ip)
	}
}

func handleCounter(ip string) {
	counter := server.GetByIp(ip)
	clientServiceSingleton := client.GetClientService()

	if counter > 4 {
		err := clientServiceSingleton.RemoveClientByIP(ip)

		if err != nil {
			slog.Error("Error removing ip: ", slog.String("ip", ip), slog.String("error", err.Error()))
			return
		}

		slog.Info("Ip deleted successfully", slog.String("ip", ip))
		return
	}

	slog.Info("Incrementing by ip: ", slog.String("ip", ip))
	server.IncrementByIp(ip)
}
