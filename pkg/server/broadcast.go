package server

import (
	"github.com/BrunoRoese/socket/pkg/command"
	"github.com/BrunoRoese/socket/pkg/network"
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

	for range ticker.C {
		err := broadcast()
		if err == nil {
			slog.Info("Broadcast successful", slog.String("message", defaultBroadcastMessage))
		}
	}
}

func broadcast() error {
	slog.Info("Broadcasting to all IPs")

	output, err := command.HandleCommand("arp", "-a")

	if err != nil {
		slog.Error("Error executing command", slog.String("error", err.Error()))
		return err
	}

	listOfIps := extractIPs(output)

	for _, ip := range listOfIps {
		//slog.Info("Broadcasting to IP", slog.String("ip", ip))

		response, err := network.SendRequest(ip, defaultBroadcastPort, []byte(defaultBroadcastMessage))

		if err != nil {
			//slog.Error("Error sending broadcast", slog.String("ip", ip), slog.String("error", err.Error()))
			continue
		}

		slog.Info("Response from IP", slog.String("ip", ip), slog.String("response", response))
	}

	return nil
}

func extractIPs(arpData string) []string {
	re := regexp.MustCompile(`\((\d+\.\d+\.\d+\.\d+)\)`)

	matches := re.FindAllStringSubmatch(arpData, -1)

	var ips []string
	for _, match := range matches {
		ips = append(ips, match[1])
	}

	return ips
}
