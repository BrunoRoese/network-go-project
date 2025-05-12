package service

import (
	"encoding/json"
	"github.com/BrunoRoese/socket/pkg/client"
	"github.com/BrunoRoese/socket/pkg/network"
	"github.com/BrunoRoese/socket/pkg/protocol"
	"github.com/BrunoRoese/socket/pkg/protocol/parser"
	"github.com/BrunoRoese/socket/pkg/protocol/validator"
	"github.com/google/uuid"
	"log/slog"
	"net"
)

type FileService struct {
	FilePath   string
	conn       *net.UDPConn
	clientAddr *net.UDPAddr
	currentId  uuid.UUID

	currentChunk chan int
	stopSending  chan bool
}

func (s *FileService) StartTransfer(ip string, filePath string) error {
	if err := validator.Validate(ip, filePath); err != nil {
		slog.Error("Error validating input", slog.String("error", err.Error()))
		return err
	}
	s.FilePath = filePath
	s.stopSending = make(chan bool)

	specifiedClient := client.FindByIp(ip)

	if specifiedClient == nil {
		slog.Error("Client not found, stopping")
		return nil
	}

	conn, err := network.CreateConn()

	if err != nil {
		slog.Error("Error creating connection", slog.String("error", err.Error()))
		return err
	}
	defer conn.Close()

	s.conn = conn

	fileContent, err := parser.ParseFile(s.FilePath)

	if err != nil {
		slog.Error("Error parsing file", slog.String("error", err.Error()))
		return err
	}

	err = s.signalStart(specifiedClient)

	if err != nil {
		slog.Error("Error signaling start", slog.String("error", err.Error()))
		s.stopSending <- true
	}

	s.startRoutines(fileContent)

	<-s.stopSending

	s.close()

	return nil
}

func (s *FileService) signalStart(specifiedClient *client.Client) error {
	req := (&protocol.File{}).BuildRequest(nil, s.FilePath, *s.conn.LocalAddr().(*net.UDPAddr))

	s.currentId = req.Information.Id

	jsonReq, err := json.Marshal(req)

	if err != nil {
		slog.Error("Error marshaling jsonReq")
		return err
	}

	_, _ = network.SendRequest(specifiedClient.Ip, specifiedClient.Port, jsonReq)

	return nil
}

func (s *FileService) startRoutines(fileContent []string) {
	s.currentChunk = make(chan int)

	s.startSendingRoutine(fileContent)
	s.startListeningRoutine()
}

func (s *FileService) startSendingRoutine(fileContent []string) {
	go func(fileContent []string) {
		for chunk := range s.currentChunk {
			slog.Info("Sending chunk", "chunk", chunk)

			if chunk >= 0 && chunk < len(fileContent) {
				currentChunk := fileContent[chunk]
				slog.Info("Sending chunk", "chunk", currentChunk)

				headers := map[string]string{
					"X-Chunk":      string(rune(chunk)),
					"X-Chunk-Size": string(rune(len(currentChunk))),
					"requestId":    s.currentId.String(),
				}

				res, err := parser.ParseProtocol(&protocol.Chunk{}, s.conn, currentChunk, headers)

				if err != nil {
					slog.Error("Error marshalling request", slog.String("error", err.Error()))
					continue
				}

				_, _ = network.SendRequest(s.clientAddr.IP.String(), s.clientAddr.Port, res)
			} else {
				slog.Error("Chunk index out of bounds", "chunk", chunk)
				s.stopSending <- true
			}
		}
	}(fileContent)
}

func (s *FileService) startListeningRoutine() {
	slog.Info("Starting listening routine")
	currentChunk := 0
	go func() {
		for {
			buffer := make([]byte, 1024)
			n, addr, err := s.conn.ReadFromUDPAddrPort(buffer)
			if err != nil {
				slog.Error("Error reading from UDP", slog.String("error", err.Error()))
				return
			}

			slog.Info("Response received", "data", string(buffer[:n]), "from", addr.String())
			req, err := parser.ParseRequest(buffer[:n])

			if err != nil {
				slog.Error("Error parsing request", slog.String("error", err.Error()))
				return
			}

			if req.Information.Id != s.currentId {
				slog.Error("Invalid request ID", slog.String("expected", s.currentId.String()), slog.String("received", req.Information.Id.String()))
				return
			}

			ip, port, err := parser.ParseSource(req.Information.Source)

			if err != nil {
				slog.Error("Error parsing source", slog.String("error", err.Error()))
				return
			}

			slog.Info("Parsed source", "ip", ip, "port", port)

			s.clientAddr = &net.UDPAddr{
				IP:   net.ParseIP(ip),
				Port: port,
			}

			currentChunk++

			s.currentChunk <- currentChunk
		}
	}()
}

func (s *FileService) close() {
	close(s.currentChunk)
	close(s.stopSending)
}
