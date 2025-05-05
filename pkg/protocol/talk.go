package protocol

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"log/slog"
	"net"
)

type Talk struct{}

func (t *Talk) Name() string {
	return "TALK"
}

func (t *Talk) BuildRequest(headers map[string]string, body string, source net.UDPAddr) Request {
	request, err := json.Marshal(body)

	if err != nil {
		slog.Error("Error building talk request")
		return Request{}
	}

	return Request{
		Information: Information{
			Method: t.Name(),
			Id:     uuid.New(),
			Source: fmt.Sprintf("%s:%d", source.IP.String(), source.Port),
		},
		Headers: Header{
			XHeader:     headers,
			ContentType: "text/plain",
		},
		Body: string(request),
	}
}
