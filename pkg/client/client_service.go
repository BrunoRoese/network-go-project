package client

import (
	"encoding/json"
	"errors"
	"github.com/BrunoRoese/socket/pkg/network"
	"github.com/BrunoRoese/socket/pkg/protocol"
	"github.com/BrunoRoese/socket/pkg/protocol/parser"
	"log/slog"
	"os"
	"sync"
	"time"
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

		_ = instance.restartFile()
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

	client.LastHeartbeat = time.Now().Unix()
	c.ClientList = append(c.ClientList, client)

	slog.Info("Client added", "clientList", c.ClientList)

	err := c.saveToFile()

	if err != nil {
		slog.Error("Error saving to file", slog.String("error", err.Error()))
	}
}

func (c *Service) UpdateClient(client *Client) {
	for _, fromList := range c.ClientList {
		if fromList.Ip == client.Ip {
			slog.Info("Client already registered, updating")
			fromList.LastHeartbeat = time.Now().Unix()
		}
	}

	slog.Info("Client added", "clientList", c.ClientList)

	err := c.saveToFile()

	if err != nil {
		slog.Error("Error saving to file", slog.String("error", err.Error()))
	}
}

func (c *Service) restartFile() error {
	if _, err := os.Stat(c.FilePath); err == nil {
		if err := os.Remove(c.FilePath); err != nil {
			return err
		}
	}

	c.ClientList = []*Client{}
	slog.Info("File deleted and client list reset")
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

func (c *Service) HandleNewClient(req *protocol.Request) error {
	slog.Info("Client not found, adding to client list")

	ip, port, err := parser.ParseSource(req.Information.Source)

	if localIp, err := network.GetLocalIp(); localIp == ip && err == nil {
		slog.Info("Client is local, using local IP", slog.String("ip", ip))
		return errors.New("client is local")
	}

	if err != nil {
		slog.Error("Error getting source parts", slog.String("error", err.Error()))
		return err
	}

	slog.Info("Parsed source", slog.String("ip", ip), slog.Int("port", port))

	newClient := &Client{Ip: ip, Port: port, LastHeartbeat: time.Now().Unix()}

	c.AddClient(newClient)

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
