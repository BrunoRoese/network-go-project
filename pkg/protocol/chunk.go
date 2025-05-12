package protocol

import (
	"github.com/google/uuid"
	"net"
)

type Chunk struct {
}

func (c *Chunk) BuildRequest(headers map[string]string, body string, source net.UDPAddr) Request {
	id, _ := uuid.Parse(headers["requestId"])

	return Request{
		Information: Information{
			Method: c.Name(),
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

func (c *Chunk) Name() string {
	return "CHUNK"
}
