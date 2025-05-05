package main

import (
	"github.com/BrunoRoese/socket/cmd"
	"github.com/BrunoRoese/socket/pkg/server"
)

var udpServer *server.Server

func main() {
	cmd.Execute()
}
