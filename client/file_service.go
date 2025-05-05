package client

import (
	"encoding/json"
	"os"
)

func (c *ClientService) SaveToFile() error {
	_, err := os.Stat(c.filePath)

	if err == nil && os.IsNotExist(err) {
		file, err := os.Create(c.filePath)
		if err != nil {
			return err
		}
		defer file.Close()

		encoder := json.NewEncoder(file)

		encoder.Encode(c.clientList)

		return nil
	}

	file, err := os.OpenFile(c.filePath, os.O_APPEND|os.O_CREATE|os.O_TRUNC, 0644)

	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)

	return encoder.Encode(c.clientList)
}
