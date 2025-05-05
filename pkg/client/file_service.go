package client

import (
	"encoding/json"
	"log/slog"
	"os"
)

func (c *ClientService) saveToFile() error {
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
