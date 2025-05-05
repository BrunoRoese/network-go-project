package protocol

import (
	"fmt"
	"github.com/google/uuid"
	"strings"
)

type Information struct {
	Method string
	Id     uuid.UUID
	Source string
}

type Header struct {
	XHeader     map[string]string `json:"X-Header"`
	ContentType string            `json:"Content-Type"`
}

type Request struct {
	Information Information
	Headers     Header
	Body        string
}

func (i Information) String() string {
	return fmt.Sprintf("Method: %s\nSource: %s", i.Method, i.Source)
}

func (h Header) String() string {
	var builder strings.Builder

	for k, v := range h.XHeader {
		builder.WriteString(fmt.Sprintf("X-Header %s: %s\n", k, v))
	}

	builder.WriteString(fmt.Sprintf("Content-Type: %s", h.ContentType))

	return builder.String()
}

func (r Request) String() string {
	var builder strings.Builder

	builder.WriteString(r.Information.String())
	builder.WriteString("\n")
	builder.WriteString(r.Headers.String())
	builder.WriteString("\n")
	builder.WriteString("Body:\n")
	builder.WriteString(r.Body)

	return builder.String()
}
