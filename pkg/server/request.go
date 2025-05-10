package server

import (
	"github.com/BrunoRoese/socket/pkg/client"
	"github.com/BrunoRoese/socket/pkg/protocol"
	"github.com/BrunoRoese/socket/pkg/protocol/parser"
	"log/slog"
	"net"
)

func (s *Server) parseRequest() (*protocol.Request, *net.UDPAddr, error) {
	slog.Info("Waiting for message")
	buffer := make([]byte, 1024)
	n, addr, err := s.Conn.ReadFromUDP(buffer)

	if err != nil {
		slog.Error("Error reading from UDP connection", slog.String("error", err.Error()))
		return nil, nil, err
	}

	req, err := parser.ParseRequest(buffer[:n])

	return req, addr, err
}

func (s *Server) getClient(info protocol.Information) (*client.Client, error) {
	parsedSource, _, err := parser.ParseSource(info.Source)

	if err != nil {
		slog.Error("Error parsing source", slog.String("error", err.Error()))
		return nil, err
	}

	return s.ClientService.GetClientByIP(parsedSource), nil
}
