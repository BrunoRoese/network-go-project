package server

import (
	"encoding/json"
	"github.com/BrunoRoese/socket/pkg/client"
	"github.com/BrunoRoese/socket/pkg/protocol"
	"github.com/BrunoRoese/socket/pkg/protocol/parser"
	"github.com/BrunoRoese/socket/pkg/server/handler"
	"log/slog"
	"net"
)

type Server struct {
	DiscoveryAddr net.UDPAddr
	GeneralAddr   net.UDPAddr
	DiscoveryConn *net.UDPConn
	GeneralConn   *net.UDPConn
	ClientService *client.Service
}

var (
	Instance    *Server
	defaultPort = 8080

	requests  = make(chan *protocol.Request, 100)
	responses = make(chan struct {
		Source   string
		Response []byte
	}, 50)
)

func Init(ip string) (*Server, error) {
	var newServer Server

	newServer.DiscoveryAddr = net.UDPAddr{IP: net.ParseIP(ip), Port: defaultPort}

	slog.Info("Server initiating...", slog.String("ip", ip), slog.Int("port", defaultPort))

	conn, err := net.ListenUDP("udp", &newServer.DiscoveryAddr)

	if err != nil {
		slog.Error("Error starting server", slog.String("ip", ip), slog.Int("port", defaultPort), slog.String("error", err.Error()))
		return nil, err
	}

	newServer.DiscoveryConn = conn

	Instance = &newServer

	newServer.ClientService = client.GetClientService()

	slog.Info("Server started", slog.String("ip", ip), slog.Int("port", defaultPort))

	newServer.startDiscoveryRoutine()
	newServer.responseRoutine()
	newServer.sendResponseRoutine()

	return &newServer, nil
}

func (s *Server) buildResponseJson(req *protocol.Request) ([]byte, error) {
	reqFunc := handler.GetRequestType(req)

	response := reqFunc(req)

	if response == nil {
		slog.Info("Handler returned null response, skipping")
		return nil, nil
	}

	jsonRequest, err := json.Marshal(response)

	if err != nil {
		slog.Error("Error marshalling request", slog.String("error", err.Error()))
		return nil, err
	}

	return jsonRequest, nil
}

func (s *Server) getClient(info protocol.Information) (*client.Client, error) {
	parsedSource, _, err := parser.ParseSource(info.Source)

	if err != nil {
		slog.Error("Error parsing source", slog.String("error", err.Error()))
		return nil, err
	}

	return s.ClientService.GetClientByIP(parsedSource), nil
}

func (s *Server) Close() {
	if err := s.DiscoveryConn.Close(); err != nil {
		slog.Error("Error closing UDP connection", slog.String("error", err.Error()))
	}

	close(requests)
	close(responses)
	slog.Info("Server closed")
}
