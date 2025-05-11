package server

import (
	"github.com/BrunoRoese/socket/pkg/network"
	"github.com/BrunoRoese/socket/pkg/protocol/parser"
	"github.com/BrunoRoese/socket/pkg/server/handler"
	"log/slog"
)

func (s *Server) startDiscoveryRoutine() {
	go func() {
		for {
			slog.Info("Waiting for message")
			buffer := make([]byte, 1024)
			n, addr, err := s.DiscoveryConn.ReadFromUDP(buffer)

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

			if req.Information.Method == "ACK" {
				slog.Info("ACK request received, skipping response")
				handler.HandleAckReq(req)
				continue
			}

			slog.Info("Received message", slog.String("from", addr.String()), slog.String("request", req.String()))

			requests <- req
		}
	}()
}

func (s *Server) responseRoutine() {
	go func() {
		for req := range requests {
			slog.Info("Handling request", slog.String("request", req.String()))
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
			}{Source: req.Information.Source, Response: response}
		}
	}()
}

func (s *Server) sendResponseRoutine() {
	go func() {
		for res := range responses {
			ip, _, err := parser.ParseSource(res.Source)
			if err != nil {
				slog.Error("Error parsing source", slog.String("error", err.Error()))
				continue
			}
			slog.Info("Parsed IP", slog.String("ip", ip))
			slog.Info("Sending response", slog.String("ip", ip), slog.String("response", string(res.Response)))
			_, err = network.SendRequest(ip, 8080, res.Response)
			if err != nil {
				slog.Error("Failed to send response", slog.String("ip", ip), slog.String("error", err.Error()))
				continue
			}
		}
	}()
}
