package protocol

import (
	"github.com/google/uuid"
	"net"
)

type File struct{}

func (f *File) Name() string {
	return "FILE"
}

func (f *File) BuildRequest(headers map[string]string, body string, source net.UDPAddr) Request {
	return Request{
		Information: Information{
			Method: f.Name(),
			Id:     uuid.New(),
			Source: source.String(),
		},
		Headers: Header{
			XHeader:     headers,
			ContentType: "text/plain",
		},
		Body: body,
	}
}
