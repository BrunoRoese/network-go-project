package server

import (
	"github.com/BrunoRoese/socket/pkg/network"
	"github.com/BrunoRoese/socket/pkg/protocol/parser"
	"github.com/BrunoRoese/socket/pkg/server/handler"
	"github.com/google/uuid"
	"log/slog"
)

func (s *Service) startGeneralRoutine() {
	go func() {
		for {
			slog.Info("Waiting for message")
			buffer := make([]byte, 1024)
			n, _, err := s.Server.GeneralConn.ReadFromUDPAddrPort(buffer)

			if err != nil {
				slog.Error("Error reading from UDP connection", slog.String("error", err.Error()))
				continue
			}

			req, err := parser.ParseRequest(buffer[:n])

			if err != nil {
				slog.Error("Error handling request", slog.String("error", err.Error()))
				continue
			}

			slog.Info("Received message", slog.String("request", req.String()))

			requests <- req
		}
	}()
}

func (s *Service) startDiscoveryRoutine() {
	go func() {
		for {
			//slog.Info("Waiting for message")
			buffer := make([]byte, 1024)
			n, addr, err := s.Server.DiscoveryConn.ReadFromUDPAddrPort(buffer)

			if err != nil {
				slog.Error("Error reading from UDP connection", slog.String("error", err.Error()))
				continue
			}

			req, err := parser.ParseRequest(buffer[:n])

			if err != nil {
				slog.Error("Error handling request", slog.String("error", err.Error()))
				continue
			}

			foundClient, err := s.getClient(req.Information)
			if err != nil {
				slog.Error(err.Error())
				continue
			}

			if foundClient == nil {
				if handleErr := s.ClientService.HandleNewClient(req); handleErr != nil {
					slog.Error("Error handling new client", slog.String("error", handleErr.Error()))
					continue
				}
			}

			if req.Information.Method != "ACK" && req.Information.Method != "HEARTBEAT" {
				slog.Error("Invalid method", slog.String("method", req.Information.Method))
				continue //TODO: RETURN NACK
			}

			if req.Information.Method == "ACK" && req.Information.Id == uuid.Nil {
				slog.Info("ACK request received for hb, skipping response")
				if foundClient != nil {
					s.ClientService.UpdateClient(foundClient)
				}
				handler.HandleAckReq(req)
				continue
			}

			slog.Info("Received message", slog.String("from", addr.String()), slog.String("request", req.String()))

			requests <- req
		}
	}()
}

func (s *Service) responseRoutine() {
	go func() {
		for req := range requests {
			response, err := s.buildResponseJson(req)
			if err != nil {
				slog.Error("Error marshalling response", slog.String("error", err.Error()))
				continue
			}

			if response == nil {
				slog.Info("Null response received, skipping")
				continue
			}

			responses <- struct {
				Source   string
				Response []byte
				Method   string
			}{Source: req.Information.Source, Response: response, Method: req.Information.Method}
		}
	}()
}

func (s *Service) sendResponseRoutine() {
	go func() {
		for res := range responses {
			ip, port, err := parser.ParseSource(res.Source)
			if err != nil {
				slog.Error("Error parsing source", slog.String("error", err.Error()))
				continue
			}
			if res.Method == "HEARTBEAT" {
				slog.Info("Heartbeat response received, rewriting port to 8080")
				port = 8080
			}
			_, err = network.SendRequest(ip, port, res.Response)
			if err != nil {
				slog.Error("Failed to send response", slog.String("ip", ip), slog.String("error", err.Error()))
				continue
			}
		}
	}()
}
