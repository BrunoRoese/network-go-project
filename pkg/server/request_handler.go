package server

import (
	"github.com/BrunoRoese/socket/pkg/protocol"
)

func GetRequestType(req *protocol.Request) func(request *protocol.Request) {
	switch req.Information.Method {
	case "ACK":
		return HandleAckReq
	case "HEARTBEAT":
		return nil
	default:
		return nil
	}
}
