package protocol

import (
	"github.com/google/uuid"
	"net"
)

type End struct{}

func (e *End) BuildRequest(headers map[string]string, body string, source net.UDPAddr) Request {
	id, _ := uuid.Parse(headers["requestId"])

	return Request{
		Information: Information{
			Method: e.Name(),
			Id:     id,
			Source: source.String(),
		},
		Headers: Header{
			XHeader:     headers,
			ContentType: "text/plain",
		},
		Body: body,
	}
}

func (e *End) Name() string {
	return "END"
}
