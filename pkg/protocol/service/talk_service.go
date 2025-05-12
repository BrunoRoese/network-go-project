package service

import (
	"github.com/BrunoRoese/socket/pkg/client"
	"github.com/BrunoRoese/socket/pkg/network"
	"github.com/BrunoRoese/socket/pkg/protocol"
	"github.com/BrunoRoese/socket/pkg/protocol/parser"
	"log/slog"
)

func Talk(ip string, message string) {
	specifiedClient := client.FindByIp(ip)

	if specifiedClient == nil {
		slog.Error("Client not found, stopping")
		return
	}

	conn, err := network.CreateConn()
	if err != nil {
		slog.Error("Error creating UDP connection", slog.String("error", err.Error()))
		return
	}
	defer conn.Close()

	jsonRequest, err := parser.ParseProtocol(&protocol.Talk{}, conn, message)

	if err != nil {
		slog.Error("Error marshalling request", slog.String("error", err.Error()))
		return
	}

	_, _ = network.SendRequest(specifiedClient.Ip, specifiedClient.Port, jsonRequest)

	buffer := make([]byte, 1024)
	n, addr, err := conn.ReadFromUDP(buffer)
	if err != nil {
		slog.Error("Error reading from UDP", slog.String("error", err.Error()))
		return
	}

	slog.Info("Response received", "data", string(buffer[:n]), "from", addr.String())
}
