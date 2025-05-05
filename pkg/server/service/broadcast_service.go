package service

import (
	"encoding/json"
	"github.com/BrunoRoese/socket/pkg/client"
	"github.com/BrunoRoese/socket/pkg/command"
	"github.com/BrunoRoese/socket/pkg/network"
	"github.com/BrunoRoese/socket/pkg/protocol"
	"github.com/BrunoRoese/socket/pkg/server"
	"log/slog"
	"regexp"
	"time"
)

var (
	defaultBroadcastPort = 8080

	defaultBroadcastMessage = "Hello, world!"
)

func Broadcast() {
	ticker := time.NewTicker(network.GetUdpTimeout())
	defer ticker.Stop()

	err := Discover()

	if err != nil {
		slog.Error("Error discovering IPs", slog.String("error", err.Error()))
	}

	ipNumberOfErrorsMap := map[string]int{}

	for range ticker.C {
		broadcast(ipNumberOfErrorsMap)
	}
}

func Discover() error {
	slog.Info("Broadcasting to all IPs")

	output, err := command.HandleCommand("arp", "-a")

	if err != nil {
		slog.Error("Error executing command", slog.String("error", err.Error()))
		return err
	}

	listOfIps := extractIPs(output)

	for _, ip := range listOfIps {
		slog.Info("Sending broadcast to", slog.String("ip", ip))
		_, err := network.SendRequest(ip, defaultBroadcastPort, []byte(defaultBroadcastMessage))

		if err != nil {
			//slog.Error("Error sending broadcast", slog.String("ip", ip), slog.String("error", err.Error()))
			continue
		}
	}

	return nil
}

func broadcast(ipNumberOfErrorsMap map[string]int) {
	slog.Info("Broadcasting to discovered IPs")

	clientService := client.GetClientService()

	for _, client := range clientService.ClientList {
		heartbeat := protocol.Heartbeat{}
		request := heartbeat.BuildRequest(nil, "", server.Instance.UdpAddr)

		slog.Info("Sending heartbeat to", slog.String("ip", client.Ip))

		jsonRequest, err := json.Marshal(request)

		if err != nil {
			slog.Info("Error marshalling request", slog.String("error", err.Error()))
			continue
		}

		slog.Info("Heartbeat request", slog.String("request", string(jsonRequest)))
		_, err = network.SendRequest(client.Ip, 8080, jsonRequest)

		if err != nil {
			ipNumberOfErrorsMap[client.Ip]++
		} else {
			ipNumberOfErrorsMap[client.Ip] = 0
		}

		if ipNumberOfErrorsMap[client.Ip] > 4 {
			err = clientService.RemoveClientByIP(client.Ip)

			if err != nil {
				slog.Error("Error removing client", slog.String("ip", client.Ip), slog.String("error", err.Error()))
			} else {
				slog.Info("Client removed", slog.String("ip", client.Ip))
			}
		}
	}
}

func extractIPs(arpData string) []string {
	re := regexp.MustCompile(`\b\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}\b`)

	matches := re.FindAllString(arpData, -1)

	ipSet := make(map[string]struct{})
	for _, ip := range matches {
		ipSet[ip] = struct{}{}
	}

	var ips []string
	for ip := range ipSet {
		ips = append(ips, ip)
	}

	return ips
}
