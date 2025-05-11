package service

import (
	"github.com/BrunoRoese/socket/pkg/network"
	"github.com/BrunoRoese/socket/pkg/protocol/parser"
	"github.com/BrunoRoese/socket/pkg/protocol/validator"
	"log/slog"
	"net"
	"os"
)

type FileService struct {
	FilePath string
	conn     *net.UDPConn

	currentChunk chan int
}

func (s *FileService) StartTransfer(ip string, filePath string) error {
	if err := validator.Validate(ip, filePath); err != nil {
		slog.Error("Error validating input", slog.String("error", err.Error()))
		return err
	}
	s.FilePath = filePath

	conn, err := network.CreateConn()

	if err != nil {
		slog.Error("Error creating connection", slog.String("error", err.Error()))
		return err
	}
	defer conn.Close()

	s.conn = conn

	fileSize, err := s.getFileSize()

	if err != nil {
		slog.Error("Error getting file size", slog.String("error", err.Error()))
		return err
	}

	slog.Info("File size", "size", fileSize)

	fileContent, err := parser.ParseFile(s.FilePath)

	if err != nil {
		slog.Error("Error parsing file", slog.String("error", err.Error()))
		return err
	}

	slog.Info("File content", "content", len(fileContent))

	return nil
	//fileReq := protocol.File{}
	//
	//req := fileReq.BuildRequest(nil, filePath, net.UDPAddr{IP: net.ParseIP(ip), Port: 8080})
	//
	//jsonReq, err := json.Marshal(req)
	//
	//if err != nil {
	//	slog.Error("Error marshaling jsonReq")
	//	return err
	//}
	//
	//_, err = network.SendRequest(ip, 8080, jsonReq)
	//if err != nil {
	//	return err
	//}
	//
	//return nil
}

func (s *FileService) getFileSize() (int64, error) {
	slog.Info("FilePath", "filePath", s.FilePath)
	fileInfo, err := os.Stat(s.FilePath)

	if err != nil {
		return 0, err
	}

	return fileInfo.Size(), nil
}

func (s *FileService) startSendingRoutine(fileContent []string) {
	go func(fileContent []string) {
		for chunk := range s.currentChunk {
			slog.Info("Sending chunk", "chunk", chunk)

			currentChunk := fileContent[chunk]

			slog.Info("Sending chunk", "chunk", currentChunk)
		}
	}(fileContent)
}

func (s *FileService) startListeningRoutine() {
	go func() {
		buffer := make([]byte, 1024)
		n, addr, err := s.conn.ReadFromUDPAddrPort(buffer)
		if err != nil {
			slog.Error("Error reading from UDP", slog.String("error", err.Error()))
			return
		}

		slog.Info("Response received", "data", string(buffer[:n]), "from", addr.String())
	}()
}
