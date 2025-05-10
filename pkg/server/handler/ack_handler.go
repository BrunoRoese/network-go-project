package handler

import (
	"github.com/BrunoRoese/socket/pkg/protocol"
	"github.com/BrunoRoese/socket/pkg/protocol/parser"
	"github.com/google/uuid"
	"log/slog"
)

func HandleAckReq(request *protocol.Request) *protocol.Request {
	if request.Information.Id == uuid.Nil {
		slog.Info("Broadcast ACK received", slog.String("source", request.Information.Source), slog.String("id", "0"))

		ip, _, err := parser.ParseSource(request.Information.Source)

		if err != nil {
			slog.Error("Error parsing source", slog.String("error", err.Error()))
			return nil
		}

		ZeroByIp(ip)
	}

	return nil
}
