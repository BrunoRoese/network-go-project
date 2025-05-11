package service

import (
	"encoding/json"
	"github.com/BrunoRoese/socket/pkg/client"
	"github.com/BrunoRoese/socket/pkg/network"
	"github.com/BrunoRoese/socket/pkg/protocol"
	"log/slog"
	"net"
)

func Talk(ip string, message string) {
	clientList := *client.GetListFromFile()

	var specifiedClient *client.Client
	for _, ci := range clientList {
		if ci.Ip == ip {
			slog.Info("Client found")
			specifiedClient = &ci
		}
	}

	if specifiedClient == nil {
		slog.Error("Client not found, stopping")
		return
	}

	localIp, err := network.GetLocalIp()

	if err != nil {
		slog.Error("Error getting local IP", slog.String("error", err.Error()))
		return
	}

	serverUdpAddr := net.UDPAddr{IP: net.ParseIP(localIp), Port: 0}

	talk := protocol.Talk{}

	request := talk.BuildRequest(nil, message, serverUdpAddr)

	jsonRequest, err := json.Marshal(request)

	if err != nil {
		slog.Error("Error marshalling request", slog.String("error", err.Error()))
		return
	}

	network.SendRequest(specifiedClient.Ip, specifiedClient.Port, jsonRequest)

	conn, err := net.ListenUDP("udp", &serverUdpAddr)
	if err != nil {
		return
	}
	defer conn.Close()

	buffer := make([]byte, 1024)
	n, addr, err := conn.ReadFromUDP(buffer)
	if err != nil {
		slog.Error("Error reading from UDP", slog.String("error", err.Error()))
		return
	}

	slog.Info("Response received", "data", string(buffer[:n]), "from", addr.String())
}
