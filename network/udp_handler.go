package network

import (
	"fmt"
	"net"
	"time"
)

func SendRequest(ip string, port int, data []byte) (string, error) {
	udpAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		return "", err
	}

	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	err = conn.SetDeadline(time.Now().Add(time.Second * 5))
	if err != nil {
		return "", err
	}

	_, err = conn.Write(data)
	if err != nil {
		return "", err
	}

	buffer := make([]byte, 1024)
	n, _, err := conn.ReadFromUDP(buffer)
	if err != nil {
		return "", err
	}

	return string(buffer[:n]), nil
}
