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

	return "", nil
}

func GetUdpTimeout() time.Duration {
	return 5 * time.Second
}

func CreateConn() (*net.UDPConn, error) {
	lIp, err := GetLocalIp()

	if err != nil {
		slog.Error("Error getting local ip", slog.String("ip", lIp))
		return nil, err
	}

	conn, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.ParseIP(lIp),
		Port: 0,
	})

	if err != nil {
		slog.Error("Error creating connection", slog.String("ip", lIp))
		return nil, err
	}

	slog.Info("Connection created", slog.String("ip", lIp), slog.Int("port", conn.LocalAddr().(*net.UDPAddr).Port))

	return conn, nil
}
