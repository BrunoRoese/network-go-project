package handler

import (
	"github.com/BrunoRoese/socket/pkg/network"
	"github.com/BrunoRoese/socket/pkg/protocol"
	"github.com/google/uuid"
	"log/slog"
	"net"
)

func HandleHeartbeatReq(request *protocol.Request) {
	response := protocol.ACK{}

	responseId := request.Information.Id

	if responseId != uuid.Nil {
		slog.Info("Heartbeat request received", slog.String("requestId", responseId.String()))
		return
	}

	headers := map[string]string{}

	headers["requestId"] = responseId.String()

	localIp, err := network.GetLocalIp()

	if err != nil {
		slog.Error("Error getting local IP", slog.String("error", err.Error()))
		return
	}

	response.BuildRequest(headers, "OK", net.UDPAddr{IP: net.ParseIP(localIp), Port: 8080})
}
