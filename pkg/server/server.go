package server

import (
	"encoding/json"
	"github.com/BrunoRoese/socket/pkg/client"
	"github.com/BrunoRoese/socket/pkg/protocol"
	"github.com/BrunoRoese/socket/pkg/protocol/parser"
	"github.com/BrunoRoese/socket/pkg/server/handler"
	"github.com/BrunoRoese/socket/pkg/server/model"
	"log/slog"
)

type Service struct {
	Server        *model.Server
	ClientService *client.Service
}

var (
	Instance    *Service
	defaultPort = 8080

	requests  = make(chan *protocol.Request, 100)
	responses = make(chan struct {
		Source   string
		Response []byte
	}, 50)
)

func Init(ip string) (*Service, error) {
	var service Service

	server, err := model.GetServer()

	if err != nil {
		slog.Error("Error getting server", slog.String("error", err.Error()))
		return nil, err
	}

	service.Server = server

	Instance = &service

	service.ClientService = client.GetClientService()

	slog.Info("ServerService started", slog.String("ip", ip), slog.Int("port", defaultPort))

	service.startGeneralRoutine()
	service.startDiscoveryRoutine()
	service.responseRoutine()
	service.sendResponseRoutine()

	return &service, nil
}

func (s *Service) buildResponseJson(req *protocol.Request) ([]byte, error) {
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

func (s *Service) getClient(info protocol.Information) (*client.Client, error) {
	parsedSource, _, err := parser.ParseSource(info.Source)

	if err != nil {
		slog.Error("Error parsing source", slog.String("error", err.Error()))
		return nil, err
	}

	return s.ClientService.GetClientByIP(parsedSource), nil
}

func (s *Service) Close() {
	if err := s.Server.DiscoveryConn.Close(); err != nil {
		slog.Error("Error closing UDP connection", slog.String("error", err.Error()))
	}

	close(requests)
	close(responses)
	slog.Info("ServerService closed")
}
