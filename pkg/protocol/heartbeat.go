package protocol

import (
	"fmt"
	"net"
)

type Heartbeat struct{}

func (h *Heartbeat) Name() string {
	return "HEARTBEAT"
}

func (h *Heartbeat) BuildRequest(headers map[string]string, body string, source net.UDPAddr) Request {
	return Request{
		Information: Information{
			Method: h.Name(),
			Source: fmt.Sprintf("%s:%d", source.IP.String(), source.Port),
		},
		Headers: Header{
			XHeader:     headers,
			ContentType: "text/plain",
		},
		Body: body,
	}
}
