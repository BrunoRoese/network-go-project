package client

import (
	"encoding/json"
	"log/slog"
	"os"
)

func GetListFromFile() *[]Client {
	data, err := os.ReadFile("resources/clients.json")
	if err != nil {
		if os.IsNotExist(err) {
			return &[]Client{}
		}
		return nil
	}

	var clients []Client
	err = json.Unmarshal(data, &clients)
	if err != nil {
		return nil
	}

	slog.Info("Clients loaded from file", "clients", clients)
	return &clients
}

func FindByIp(ip string) *Client {
	clientList := *GetListFromFile()

	for _, c := range clientList {
		if c.Ip == ip {
			slog.Info("Client found")
			return &c
		}
	}

	slog.Info("Client not found")
	return nil
}
