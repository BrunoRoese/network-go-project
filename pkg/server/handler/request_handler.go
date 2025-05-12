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
	default:
		slog.Info("Default request received")
		return HandleDefaultReq
	}
}
