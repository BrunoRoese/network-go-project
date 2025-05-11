package server

import (
	"encoding/json"
	"github.com/BrunoRoese/socket/pkg/client"
	"github.com/BrunoRoese/socket/pkg/network"
	"github.com/BrunoRoese/socket/pkg/protocol"
	"github.com/BrunoRoese/socket/pkg/protocol/parser"
	"github.com/BrunoRoese/socket/pkg/server/handler"
	"log/slog"
	"net"
)

type Server struct {
	UdpAddr       net.UDPAddr
	Conn          *net.UDPConn
	ClientService *client.Service
}

var (
	Instance *Server

	requests  = make(chan *protocol.Request, 100)
	responses = make(chan struct {
		Source   string
		Response []byte
	}, 50)
)

func Init(ip string, port int) (*Server, error) {
	var newServer Server

	newServer.UdpAddr = net.UDPAddr{IP: net.ParseIP(ip), Port: port}

	slog.Info("Server initiating...", slog.String("ip", ip), slog.Int("port", port))

	conn, err := net.ListenUDP("udp", &newServer.UdpAddr)

	if err != nil {
		slog.Error("Error starting server", slog.String("ip", ip), slog.Int("port", port), slog.String("error", err.Error()))
		return nil, err
	}

	newServer.Conn = conn

	Instance = &newServer

	newServer.ClientService = client.GetClientService()

	slog.Info("Server started", slog.String("ip", ip), slog.Int("port", port))

	return &newServer, nil
}

func (s *Server) StartListeningRoutine() {
	s.responseRoutine()
	s.sendResponseRoutine()
	go func() {
		for {
			req, addr, err := s.parseRequest()

			if err != nil {
				slog.Error("Error handling request", slog.String("error", err.Error()))
				continue
			}

			foundClient, err := s.getClient(req.Information)
			if err != nil {
				slog.Error(err.Error())
				continue
			}

			if foundClient == nil && s.ClientService.HandleNewClient(req) != nil {
				slog.Error("Error handling new client", slog.String("error", err.Error()))
				continue
			}

			slog.Info("Received message", slog.String("from", addr.String()), slog.String("request", req.String()))

			if req.Information.Method == "ACK" {
				handler.ZeroByIp(foundClient.Ip)
				continue
			}

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

			slog.Info("Sending response", slog.String("response", string(response)))
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
			slog.Info("Sending response", slog.String("ip", ip), slog.String("response", string(res.Response)))
			_, err = network.SendRequest(ip, 8080, res.Response)
			if err != nil {
				slog.Error("Error sending response", slog.String("ip", ip), slog.String("error", err.Error()))
				return
			}
		}
	}()
}

func (s *Server) buildResponseJson(req *protocol.Request) ([]byte, error) {
	reqFunc := handler.GetRequestType(req)

	response := reqFunc(req)

	jsonRequest, err := json.Marshal(response)

	if err != nil {
		slog.Error("Error marshalling request", slog.String("error", err.Error()))
		return nil, err
	}

	return jsonRequest, nil
}

func (s *Server) Close() {
	if err := s.Conn.Close(); err != nil {
		slog.Error("Error closing UDP connection", slog.String("error", err.Error()))
	}

	close(requests)
	close(responses)
	slog.Info("Server closed")
}
