package server

import (
	"github.com/BrunoRoese/socket/pkg/network"
	"github.com/BrunoRoese/socket/pkg/protocol"
	"github.com/BrunoRoese/socket/pkg/protocol/parser"
	"github.com/BrunoRoese/socket/pkg/server/handler"
	"github.com/BrunoRoese/socket/pkg/server/model"
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

			req, err := parser.ParseLargeRequest(buffer[:n])

			if err != nil {
				slog.Error("Error handling request", slog.String("error", err.Error()))
				continue
			}

			if req.Information.Method == "FILE" {
				slog.Info("File request received, creating new connection")
				newConn, err := network.CreateConn()

				if err != nil {
					slog.Error("Error creating new connection", slog.String("error", err.Error()))
					continue
				}

				s.Server.FileAddrMap[req.Information.Id.String()] = newConn

				s.startFileSavingRoutine(newConn)
			}

			slog.Info("Received message", slog.String("request", req.String()))

			requests <- req
		}
	}()
}

func (s *Service) startDiscoveryRoutine() {
	go func() {
		for {
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
			} else {
				slog.Info("Client found, updating")
				s.ClientService.UpdateClient(foundClient)
			}

			if req.Information.Method != "ACK" && req.Information.Method != "HEARTBEAT" {
				slog.Error("Invalid method", slog.String("method", req.Information.Method))
				continue //TODO: RETURN NACK
			}

			if req.Information.Method == "ACK" && req.Information.Id == uuid.Nil {
				slog.Info("ACK request received for hb, skipping response")
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
			go func(innerReq *protocol.Request) {
				response, err := s.buildResponseJson(innerReq)
				if err != nil {
					slog.Error("Error marshalling response", slog.String("error", err.Error()))
					return
				}

				if response == nil {
					slog.Info("Null response received, skipping")
					return
				}

				responses <- model.Response{Source: req.Information.Source, Res: response, Method: req.Information.Method}
			}(req)
		}
	}()
}

func (s *Service) sendResponseRoutine() {
	go func() {
		for res := range responses {
			go func(response model.Response) {
				ip, port, err := parser.ParseSource(res.Source)
				if err != nil {
					slog.Error("Error parsing source", slog.String("error", err.Error()))
					return
				}
				if res.Method == "HEARTBEAT" {
					slog.Info("Heartbeat response received, rewriting port to 8080")
					port = 8080
				}
				_, _ = network.SendRequest(ip, port, res.Res)
			}(res)
		}
	}()
}
