package handler

import (
	"github.com/BrunoRoese/socket/pkg/protocol"
	"github.com/google/uuid"
	"log/slog"
)

func HandleAckReq(request *protocol.Request) {
	if request.Information.Id == uuid.Nil() {
		slog.Info("Broadcast ACK received", slog.String("source", request.Information.Source), slog.String("id", "0"))

	}
}
