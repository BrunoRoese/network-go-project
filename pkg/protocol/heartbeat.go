package protocol

import (
	"fmt"
	"github.com/BrunoRoese/socket/pkg/server"
	"strings"
)

type Heartbeat struct{}

func (h *Heartbeat) Name() string {
	return "HEARTBEAT"
}

func (h *Heartbeat) BuildRequest(headers map[string]string, body string) string {
	var builder strings.Builder
	udpAddr := server.GetServer().UdpAddr

	builder.WriteString(fmt.Sprintf("Method: %s\n", h.Name()))
	builder.WriteString(fmt.Sprintf("Source: %s:%s", udpAddr.IP.String(), udpAddr.Port))

	for k, v := range headers {
		builder.WriteString(fmt.Sprintf("\n%s: %s", k, v))
	}

	builder.WriteString(body)

	return builder.String()
}
