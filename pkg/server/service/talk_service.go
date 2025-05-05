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
	var specifiedClient *client.Client

	for _, client := range clientService.ClientList {
		if client.Ip == ip {
			slog.Info("Client found, sending message")
			specifiedClient = client
			break
		}
	}

	if specifiedClient == nil {
		slog.Error("Client not found, stopping")
		return
	}

	serverUdpAddr := net.UDPAddr{IP: net.IP("192.168.0.11"), Port: 8080}

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
