package cmd

import (
	"github.com/BrunoRoese/socket/pkg/client"
	"github.com/spf13/cobra"
	"time"
)

var devicesCmd = &cobra.Command{
	Use:   "devices",
	Short: "Get's a list of all devices listed and discovered.",
	Long:  `This will return a list of all devices and informations.`,
	Run: func(cmd *cobra.Command, args []string) {
		client.GetClientService()

		clientList := client.GetClientService().ClientList

		if len(clientList) == 0 {
			cmd.Println("No devices found.")
			return
		}

		cmd.Println("Devices found:")
		for _, client := range clientList {
			timeDiff := time.Now().Unix() - client.LastHeartbeat
			cmd.Printf("IP: %s, Port: %d, Seconds since last HB: %d", client.Ip, client.Port, timeDiff)
		}
	},
}

func init() {
	rootCmd.AddCommand(devicesCmd)
}
