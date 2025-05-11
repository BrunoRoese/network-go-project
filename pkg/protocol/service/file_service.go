package service

import (
	"encoding/json"
	"github.com/BrunoRoese/socket/pkg/network"
	"github.com/BrunoRoese/socket/pkg/protocol"
	"github.com/BrunoRoese/socket/pkg/protocol/validator"
	"log/slog"
	"net"
)

type FileService struct {
	FilePath string
}

func (f *FileService) StartTransfer(ip string, filePath string) error {
	if err := validator.Validate(ip, filePath); err != nil {
		slog.Error("Error validating input", slog.String("error", err.Error()))
		return err
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
