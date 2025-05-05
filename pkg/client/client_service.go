package client

import "log/slog"

type ClientService struct {
	ClientList []*Client
	FilePath   string
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
