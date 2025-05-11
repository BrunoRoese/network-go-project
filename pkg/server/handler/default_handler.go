package handler

import (
	"github.com/BrunoRoese/socket/pkg/network"
	"github.com/BrunoRoese/socket/pkg/protocol"
	"github.com/BrunoRoese/socket/pkg/server/model"
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

	server, _ := model.GetServer()

	udpAddr := net.UDPAddr{IP: net.ParseIP(localIp), Port: server.GeneralAddr.Port}

	headers := map[string]string{}

	headers["requestId"] = req.Information.Id.String()

	response := res.BuildRequest(headers, "OK", udpAddr)

	return &response
}
