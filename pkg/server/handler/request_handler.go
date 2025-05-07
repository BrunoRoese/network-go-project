package handler

import "github.com/BrunoRoese/socket/pkg/protocol"

func GetRequestType(req *protocol.Protocol) string {
	switch (*req).(type) {
	case *protocol.ACK:
		return "ACK"
	case *protocol.Heartbeat:
		return "Heartbeat"
	default:
		return "Unknown"
	}
}
