package server

import (
	"github.com/BrunoRoese/socket/pkg/client"
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

			for _, client := range s.ClientService.ClientList {
				if client.Ip == addr.IP.String() {
					slog.Info("Client already exists", slog.String("client", client.Ip))
					slog.Info("Message", slog.String("message", string(buffer[:n])))
					continue
				}
			}

			newClient := &client.Client{Ip: addr.IP.String(), Port: addr.Port}

			s.ClientService.AddClient(newClient)

			_, err = s.Conn.WriteToUDP([]byte("Hello, world!"), addr)
			if err != nil {
				return
			}

			slog.Info("Received message", slog.String("from", addr.String()))
		}
	}()
}

//func (s *Server) ackResponse(clientList []*client.Client) {
//	for _, client := range clientList {
//		if client.Ip == addr.IP.String() {
//			slog.Info("Client already exists", slog.String("client", client.Ip))
//			continue
//		}
//	}
//}

func (s *Server) Close() {
	if err := s.Conn.Close(); err != nil {
		slog.Error("Error closing UDP connection", slog.String("error", err.Error()))
	}
	slog.Info("Server closed")
}
