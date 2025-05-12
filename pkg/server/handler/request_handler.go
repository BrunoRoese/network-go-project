package handler

import (
	"github.com/BrunoRoese/socket/pkg/protocol"
	"log/slog"
)

func GetRequestType(req *protocol.Request) func(request *protocol.Request) *protocol.Request {
	slog.Info("request method", slog.String("method", req.Information.Method))
	switch req.Information.Method {
	case "HEARTBEAT":
		slog.Info("HEARTBEAT received")
		return HandleHeartbeatReq
	case "FILE":
		slog.Info("FILE received")
		return HandleFileReq
	case "END":
		slog.Info("END received")
		return HandleEndReq
	case "NACK":
		slog.Info("NACK received")
		return
	default:
		slog.Info("Default request received")
		return HandleDefaultReq
	}
}
