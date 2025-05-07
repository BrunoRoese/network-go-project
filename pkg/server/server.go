package server

import (
	"errors"
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
			req, err := parser.ParseRequest(buffer[:n])

			if err != nil {
				slog.Error("Error parsing request", slog.String("error", err.Error()))
				continue
			}

			if foundClient == nil {
				err := handleNewClient(s, addr, req)

				if err != nil {
					slog.Error("Error handling new client", slog.String("error", err.Error()))
					continue
				}
			} else {
				slog.Info("Client found in client list", slog.String("ip", addr.IP.String()))
			}

			reqFunc := handler.GetRequestType(req)

			reqFunc(req)

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

func handleNewClient(s *Server, addr *net.UDPAddr, req *protocol.Request) error {
	slog.Info("Client not found, adding to client list", slog.String("ip", addr.IP.String()))

	ip, port, err := parser.ParseSource(req.Information.Source)

	if ip, err = network.GetLocalIp(); ip == "" && err == nil {
		slog.Info("Client is local, using local IP", slog.String("ip", ip))
		return errors.New("client is local")
	}

	if err != nil {
		slog.Error("Error getting source parts", slog.String("error", err.Error()))
		return err
	}

	newClient := &client.Client{Ip: ip, Port: port}

	s.ClientService.AddClient(newClient)

	return nil
}
