package model

import (
	"github.com/BrunoRoese/socket/pkg/client"
	"github.com/BrunoRoese/socket/pkg/network"
	"log/slog"
	"net"
	"sync"
)

type Server struct {
	DiscoveryAddr net.UDPAddr
	GeneralAddr   net.UDPAddr
	DiscoveryConn *net.UDPConn
	GeneralConn   *net.UDPConn
	ClientService *client.Service
}

var (
	instance *Server
	once     sync.Once
)

func GetServer() (*Server, error) {
	var err error
	once.Do(func() {
		localIp, err := network.GetLocalIp()
		if err != nil {
			slog.Error("Error getting local IP", slog.String("error", err.Error()))
			return
		}
		discoveryConn, connErr := net.ListenUDP("udp", &net.UDPAddr{Port: 8080, IP: net.ParseIP(localIp)})
		if connErr != nil {
			err = connErr
			return
		}

		generalConn, connErr := net.ListenUDP("udp", &net.UDPAddr{Port: 0, IP: net.ParseIP(localIp)})
		if connErr != nil {
			err = connErr
			return
		}

		clientService := client.GetClientService()

		instance = &Server{
			DiscoveryAddr: *discoveryConn.LocalAddr().(*net.UDPAddr),
			GeneralAddr:   *generalConn.LocalAddr().(*net.UDPAddr),
			DiscoveryConn: discoveryConn,
			GeneralConn:   generalConn,
			ClientService: clientService,
		}
	})
	return instance, err
}
