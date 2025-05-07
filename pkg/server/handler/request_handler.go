package handler

import "github.com/BrunoRoese/socket/pkg/protocol"

func GetRequestType(req *protocol.Protocol) func(request *protocol.Request) {
	switch (*req).(type) {
	case *protocol.ACK:
		return HandleAckReq
	case *protocol.Heartbeat:
		return nil
	default:
		return nil
	}
}
