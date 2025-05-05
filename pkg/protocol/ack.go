package protocol

import (
	"github.com/google/uuid"
	"net"
)

type ACK struct{}

func (h *ACK) Name() string {
	return "ACK"
}

func (h *ACK) BuildRequest(headers map[string]string, body string, source net.UDPAddr) Request {
	requestId := parseUUID(headers["requestId"])

	if requestId == uuid.Nil {
		return Request{}
	}

	return Request{
		Information: Information{
			Method: h.Name(),
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

func parseUUID(id string) uuid.UUID {
	requestId, err := uuid.Parse(id)

	if err != nil {
		return uuid.Nil
	}

	return requestId
}
