package service

import (
	"encoding/json"
	"github.com/BrunoRoese/socket/pkg/client"
	"github.com/BrunoRoese/socket/pkg/command"
	"github.com/BrunoRoese/socket/pkg/network"
	"github.com/BrunoRoese/socket/pkg/protocol"
	"github.com/BrunoRoese/socket/pkg/server"
	"github.com/BrunoRoese/socket/pkg/server/handler"
	"log/slog"
	"regexp"
	"sync"
	"time"
)

func Broadcast() {
	ticker := time.NewTicker(network.GetUdpTimeout())
	defer ticker.Stop()

	err := Discover()

	if err != nil {
		slog.Error("Error discovering IPs", slog.String("error", err.Error()))
	}

	if client.GetClientService().ClientList == nil || len(client.GetClientService().ClientList) == 0 {
		slog.Info("No clients found, trying discover again")
		<-ticker.C
		Broadcast()
	}

	for range ticker.C {
		if client.GetClientService().ClientList == nil || len(client.GetClientService().ClientList) == 0 {
			slog.Info("No clients found, trying discover again")
			break
		}
		broadcast()
	}

	Broadcast()
}

func Discover() error {
	slog.Info("Broadcasting to all IPs")

	output, err := command.HandleCommand("arp", "-a")
	if err != nil {
		slog.Error("Error executing command", slog.String("error", err.Error()))
		return err
	}

	listOfIps := extractIPs(output)
	localIp, err := network.GetLocalIp()
	if err != nil {
		slog.Error("Error retrieving local IP", slog.String("error", err.Error()))
		return err
	}

	var wg sync.WaitGroup
	for _, ip := range listOfIps {
		if ip == localIp {
			continue
		}

		wg.Add(1)
		go func(ip string) {
			defer wg.Done()

			jsonRequest, err := buildHeartbeatReq()
			if err != nil {
				slog.Error("Error building request", slog.String("ip", ip), slog.String("error", err.Error()))
				return
			}

			network.SendRequest(ip, server.Instance.Server.DiscoveryAddr.Port, jsonRequest)
		}(ip)
	}

	wg.Wait()
	return nil
}

func broadcast() {
	slog.Info("Broadcasting to discovered IPs")

	clientService := client.GetClientService()

	for _, c := range clientService.ClientList {
		jsonRequest, err := buildHeartbeatReq()

		if err != nil {
			slog.Info("Error building request", slog.String("error", err.Error()))
			continue
		}

		slog.Info("Sending heartbeat to", slog.String("ip", c.Ip))
		_, err = network.SendRequest(c.Ip, server.Instance.Server.DiscoveryAddr.Port, jsonRequest)

		handler.IncrementByIp(c.Ip)
	}
}

func buildHeartbeatReq() ([]byte, error) {
	heartbeat := protocol.Heartbeat{}
	request := heartbeat.BuildRequest(nil, "", server.Instance.Server.GeneralAddr)

	jsonRequest, err := json.Marshal(request)

	if err != nil {
		slog.Info("Error marshalling request", slog.String("error", err.Error()))
		return nil, err
	}

	return jsonRequest, nil
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
