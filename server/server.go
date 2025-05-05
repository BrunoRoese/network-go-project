package server

import (
	"github.com/BrunoRoese/socket/client"
	"log/slog"
	"net"
)

type Server struct {
	udpAddr net.UDPAddr
	conn    *net.UDPConn
}

func Init(ip string, port int) (*Server, error) {
	var newServer Server

	newServer.udpAddr = net.UDPAddr{IP: net.ParseIP(ip), Port: port}

	slog.Info("Server initiating...", slog.String("ip", ip), slog.Int("port", port))

	conn, err := net.ListenUDP("udp", &newServer.udpAddr)

	if err != nil {
		slog.Error("Error starting server", slog.String("ip", ip), slog.Int("port", port), slog.String("error", err.Error()))
		return nil, err
	}

	newServer.conn = conn

	slog.Info("Server started", slog.String("ip", ip), slog.Int("port", port))

	return &newServer, nil
}

func (s *Server) StartListeningRoutine() {
	clientService := client.NewClientService("clients.json")

	go func() {
		for {
			slog.Info("Waiting for message")
			buffer := make([]byte, 1024)
			n, addr, err := s.conn.ReadFromUDP(buffer)

			if err != nil {
				slog.Error("Error reading from UDP connection", slog.String("error", err.Error()))
				continue
			}

			newClient := client.CreateClient(string(addr.IP), addr.Port)

			clientService.AddClient(newClient)
			slog.Info("Received message", slog.String("message", string(buffer[:n])), slog.String("from", addr.String()))
		}
	}()
}

func (s *Server) Close() {
	if err := s.conn.Close(); err != nil {
		slog.Error("Error closing UDP connection", slog.String("error", err.Error()))
	}
	slog.Info("Server closed")
}
