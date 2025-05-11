package parser

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/BrunoRoese/socket/pkg/protocol"
	"log"
	"net"
	"strconv"
	"strings"
)

func ParseRequest(data []byte) (*protocol.Request, error) {
	var req protocol.Request

	err := json.Unmarshal(data, &req)

	if err != nil {
		log.Printf("Error parsing request: %v", err)
		return nil, err
	}

	return &req, nil
}

func ParseSource(source string) (string, int, error) {
	parts := strings.Split(source, ":")
	if len(parts) != 2 {
		return "", 0, fmt.Errorf("invalid source format: %s", source)
	}

	ip := parts[0]
	port, err := strconv.Atoi(parts[1])
	if err != nil {
		return "", 0, fmt.Errorf("invalid port: %s", parts[1])
	}

	return ip, port, nil
}

func ParseProtocol(protocol protocol.Protocol, conn *net.UDPConn, message string) ([]byte, error) {
	serverUdpAddr := conn.LocalAddr().(*net.UDPAddr)

	request := protocol.BuildRequest(nil, message, *serverUdpAddr)

	jsonRequest, err := json.Marshal(request)

	if err != nil {
		return nil, errors.New("error marshalling request")
	}

	return jsonRequest, nil
}
