package network

import (
	"fmt"
	"log/slog"
	"net"
	"time"
)

func SendRequest(ip string, port int, data []byte) (string, error) {
	if localIp, _ = GetLocalIp(); localIp == ip {
		slog.Error("Cannot send request to self", slog.String("ip", ip))
		return "", fmt.Errorf("cannot send request to self")
	}

	udpAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		return "", err
	}

	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	err = conn.SetDeadline(time.Now().Add(GetUdpTimeout()))
	if err != nil {
		return "", err
	}

	_, err = conn.Write(data)
	if err != nil {
		slog.Error("Error sending data", slog.String("error", err.Error()))
		return "", err
	}

	buffer := make([]byte, 1024)
	_, _, err = conn.ReadFromUDP(buffer)
	if err != nil {
		return "", err
	}

	return "", nil
}

func GetUdpTimeout() time.Duration {
	return 5 * time.Second
}
