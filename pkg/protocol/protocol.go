package protocol

import "net"

type Protocol interface {
	Name() string
	BuildRequest(headers map[string]string, body string, source net.UDPAddr) Request
}
