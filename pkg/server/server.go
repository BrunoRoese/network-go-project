package server

import (
	"encoding/json"
	"github.com/BrunoRoese/socket/pkg/client"
	"github.com/BrunoRoese/socket/pkg/protocol"
	"log"
	"log/slog"
	"net"
)

type Server struct {
	UdpAddr       net.UDPAddr
	Conn          *net.UDPConn
	ClientService *client.Service
}

var Instance *Server

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
	go func() {
		for {
			slog.Info("Waiting for message")
			buffer := make([]byte, 1024)
			n, addr, err := s.Conn.ReadFromUDP(buffer)

			if err != nil {
				slog.Error("Error reading from UDP connection", slog.String("error", err.Error()))
				continue
			}

			foundClient := s.ClientService.GetClientByIP(addr.IP.String())

			if foundClient == nil {
				slog.Info("Client not found, adding to client list", slog.String("ip", addr.IP.String()))

				newClient := &client.Client{Ip: addr.IP.String(), Port: addr.Port}

				s.ClientService.AddClient(newClient)
			} else {
				slog.Info("Client found in client list", slog.String("ip", addr.IP.String()))
			}

			req, err := parseRequest(buffer[:n])

			if err != nil {
				slog.Error("Error parsing request", slog.String("error", err.Error()))
				continue
			}

			slog.Info("Received message", slog.String("from", addr.String()), slog.String("request", req.String()))
		}
	}()
}

func (s *Server) Close() {
	if err := s.Conn.Close(); err != nil {
		slog.Error("Error closing UDP connection", slog.String("error", err.Error()))
	}
	slog.Info("Server closed")
}

func parseRequest(n []byte) (*protocol.Request, error) {
	var req protocol.Request

	err := json.Unmarshal(n, &req)
	if err != nil {
		log.Printf("Error parsing request: %v", err)
		return nil, err
	}

	return &req, nil
}
