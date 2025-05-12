package protocol

import (
	"net"
)

type NACK struct {
}

func (n *NACK) Name() string {
	return "ACK"
}

func (n *NACK) BuildRequest(headers map[string]string, body string, source net.UDPAddr) Request {
	requestId := parseUUID(headers["requestId"])

	delete(headers, "requestId")

	return Request{
		Information: Information{
			Method: n.Name(),
			Id:     requestId,
			Source: source.String(),
		},
		Headers: Header{
			XHeader:     headers,
			ContentType: "text/plain",
		},
		Body: body,
	}
}
