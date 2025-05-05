package client

import (
	"encoding/json"
	"log/slog"
	"os"
	"sync"
)

type Service struct {
	ClientList []*Client
	FilePath   string
}

var (
	instance *Service
	once     sync.Once
)

func GetClientService() *Service {
	return getClientService("resources/clients.json")
}

func getClientService(filePath string) *Service {
	once.Do(func() {
		instance = &Service{
			ClientList: []*Client{},
			FilePath:   filePath,
		}

		_ = instance.LoadFromFile()
	})
	return instance
}

func (c *Service) AddClient(client *Client) {
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

func (c *Service) LoadFromFile() error {
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

	slog.Info("Clients loaded from file", "clients", clients)
	c.ClientList = clients
	return nil
}

func (c *Service) RemoveClientByIP(ip string) error {
	for i, client := range c.ClientList {
		if client.Ip == ip {
			c.ClientList = append(c.ClientList[:i], c.ClientList[i+1:]...)

			err := c.saveToFile()
			if err != nil {
				slog.Error("Error saving to file after removal", slog.String("error", err.Error()))
				return err
			}

			slog.Info("Client removed successfully", slog.String("ip", ip))
			return nil
		}
	}

	slog.Info("Client not found", slog.String("ip", ip))
	return nil
}

func (c *Service) GetClientByIP(ip string) *Client {
	for _, client := range c.ClientList {
		if client.Ip == ip {
			slog.Info("Client found")
			return client
		}
	}

	return nil
}

func (c *Service) saveToFile() error {
	_, err := os.Stat(c.FilePath)

	var file *os.File

	if os.IsNotExist(err) {
		file, err = os.Create(c.FilePath)
		if err != nil {
			return err
		}
	} else {
		file, err = os.OpenFile(c.FilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			return err
		}
	}
	defer file.Close()

	slog.Info("Saving clients to file", "clientList", c.ClientList)

	encoder := json.NewEncoder(file)
	return encoder.Encode(c.ClientList)
}
