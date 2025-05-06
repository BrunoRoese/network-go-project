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
	clientService := client.GetClientService()

	if clientService.ClientList == nil || len(clientService.ClientList) == 0 {
		err := Discover()

		if err != nil {
			slog.Error("Error discovering IPs", slog.String("error", err.Error()))
			return
		}
	}

	specifiedClient := clientService.GetClientByIP(ip)

	if specifiedClient == nil {
		slog.Error("Client not found, stopping")
		return
	}

	localIp, err := network.GetLocalIp()

	if err != nil {
		slog.Error("Error getting local IP", slog.String("error", err.Error()))
		return
	}

	serverUdpAddr := net.UDPAddr{IP: net.ParseIP(localIp), Port: 8080}

	talk := protocol.Talk{}

	request := talk.BuildRequest(nil, message, serverUdpAddr)

	jsonRequest, err := json.Marshal(request)

	if err != nil {
		slog.Error("Error marshalling request", slog.String("error", err.Error()))
		return
	}

	_, err = network.SendRequest(ip, 8080, jsonRequest)

	if err != nil {
		slog.Error("Error sending request", slog.String("ip", ip), slog.String("error", err.Error()))
		return
	}

	slog.Info("Request sent!", "Request", string(jsonRequest))
}
