package handler

import (
	"github.com/BrunoRoese/socket/pkg/client"
	"log/slog"
)

var requestsMap = make(map[string]int)

func GetByIp(ip string) int {
	if count, ok := requestsMap[ip]; ok {
		return count
	}

	return 0
}

func ZeroByIp(ip string) {
	if _, ok := requestsMap[ip]; ok {
		delete(requestsMap, ip)
		slog.Info("Request count reset", slog.String("ip", ip))
	} else {
		slog.Info("Request count already zero", slog.String("ip", ip))
	}
}

func IncrementByIp(ip string) {
	if count, ok := requestsMap[ip]; ok {
		requestsMap[ip] = count + 1
	} else {
		requestsMap[ip] = 1
	}

	if requestsMap[ip] > 4 {
		delete(requestsMap, ip)

		clientService := client.GetClientService()

		err := clientService.RemoveClientByIP(ip)

		if err != nil {
			slog.Error("Error removing client by IP, rolling back", slog.String("ip", ip), slog.String("error", err.Error()))

			requestsMap[ip] = 5
		} else {
			slog.Info("Client removed successfully", slog.String("ip", ip))
		}
	} else {
		slog.Info("Incremented request count", slog.String("ip", ip), slog.Int("count", requestsMap[ip]))
	}
}
