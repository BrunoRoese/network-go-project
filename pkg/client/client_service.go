package client

import (
	"encoding/json"
	"log/slog"
	"os"
	"sync"
)

type ClientService struct {
	ClientList []*Client
	FilePath   string
}

var (
	instance *ClientService
	once     sync.Once
)

func GetClientService() *ClientService {
	return getClientService("resources/clients.json")
}

func getClientService(filePath string) *ClientService {
	once.Do(func() {
		instance = &ClientService{
			ClientList: []*Client{},
			FilePath:   filePath,
		}

		_ = instance.LoadFromFile()
	})
	return instance
}

func (c *ClientService) AddClient(client *Client) {
	for _, fromList := range c.ClientList {
		if fromList.Ip == client.Ip {
			slog.Info("Client already registered, skipping")
			return
		}
	}

	slog.Info("Client not registered, adding", "client", client)

	c.ClientList = append(c.ClientList, client)

	slog.Info("Client added", "clientList", c.ClientList)

	err := c.saveToFile()

	if err != nil {
		slog.Error("Error saving to file", slog.String("error", err.Error()))
	}
}

func (c *ClientService) LoadFromFile() error {
	data, err := os.ReadFile(c.FilePath)
	if err != nil {
		if os.IsNotExist(err) {
			c.ClientList = []*Client{}
			return nil
		}
		return err
	}

	var clients []*Client
	err = json.Unmarshal(data, &clients)
	if err != nil {
		return err
	}

	c.ClientList = clients
	return nil
}
