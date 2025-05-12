package handler

import (
	"github.com/BrunoRoese/socket/pkg/protocol"
	"github.com/BrunoRoese/socket/pkg/server/model"
	"github.com/google/uuid"
	"log/slog"
	"net"
)

func HandleNackRe(request *protocol.Request) *protocol.Request {
	responseId := request.Information.Id

	if responseId == uuid.Nil {
		slog.Info("File request received with null id", slog.String("requestId", responseId.String()))
		return nil
	}

	headers := map[string]string{}

	headers["requestId"] = responseId.String()

	server, _ := model.GetServer()

	res := (&protocol.End{}).BuildRequest(headers, "ERROR", *server.FileAddrMap[request.Information.Id.String()].LocalAddr().(*net.UDPAddr))

	return &res
}
