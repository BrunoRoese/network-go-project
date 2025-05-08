package service

import (
	"encoding/json"
	"errors"
	"github.com/BrunoRoese/socket/pkg/network"
	"github.com/BrunoRoese/socket/pkg/protocol"
	"log/slog"
	"net"
	"os"
)

type FileService struct {
	FilePath string
}

func (f *FileService) StartTransfer(ip string, filePath string) error {
	if err := validateFilePath(filePath); err != nil {
		return errors.New("file path is empty")
	}

	if err := validateIp(ip); err != nil {
		return errors.New("ip is empty")
	}

	fileReq := protocol.File{}

	req := fileReq.BuildRequest(nil, filePath, net.UDPAddr{IP: net.ParseIP(ip), Port: 8080})

	jsonReq, err := json.Marshal(req)

	if err != nil {
		slog.Error("Error marshaling jsonReq")
		return err
	}

	_, err = network.SendRequest(ip, 8080, jsonReq)
	if err != nil {
		return err
	}

	return nil
}

func validateIp(ip string) error {
	if ip == "" {
		return errors.New("ip is empty")
	}

	if net.ParseIP(ip) == nil {
		return errors.New("invalid ip address")
	}

	return nil
}

func validateFilePath(filePath string) error {
	if filePath == "" {
		return errors.New("file path is empty")
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return errors.New("file does not exist")
	}

	return nil
}
