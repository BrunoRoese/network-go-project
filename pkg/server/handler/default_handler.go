package handler

import (
	"github.com/BrunoRoese/socket/pkg/network"
	"github.com/BrunoRoese/socket/pkg/protocol"
	"log/slog"
	"net"
)

func HandleDefaultReq(req *protocol.Request) *protocol.Request {
	res := protocol.ACK{}

	localIp, err := network.GetLocalIp()

	if err != nil {
		slog.Error("Error getting local IP", slog.String("error", err.Error()))
		return nil
	}

	udpAddr := net.UDPAddr{IP: net.ParseIP(localIp), Port: 8080}

	headers := map[string]string{}

	headers["requestId"] = req.Information.Id.String()

	response := res.BuildRequest(headers, "OK", udpAddr)

	return &response
}
