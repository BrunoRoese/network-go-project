package client

import "log/slog"

type ClientService struct {
	clientList []*Client
	filePath   string
}

func NewClientService(filepath string) *ClientService {
	return &ClientService{
		clientList: make([]*Client, 0),
		filePath:   filepath,
	}
}

func (c *ClientService) AddClient(client *Client) {
	for _, fromList := range c.clientList {
		if fromList.ip == client.ip {
			slog.Info("Client already registered, skipping")
			return
		}
	}

	c.clientList = append(c.clientList, client)

	err := c.SaveToFile()

	if err != nil {
		slog.Error("Error saving to file", slog.String("error", err.Error()))
	}
}
