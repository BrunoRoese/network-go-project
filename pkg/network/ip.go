package network

import "net"

var localIp string

func GetLocalIp() (string, error) {
	if localIp == "" {
		conn, err := net.Dial("udp", "8.8.8.8:80")

		if err != nil {
			return "", err
		}
		defer conn.Close()

		localAddr := conn.LocalAddr().(*net.UDPAddr)

		localIp = localAddr.IP.String()

		return localAddr.IP.String(), nil
	}

	return localIp, nil
}
