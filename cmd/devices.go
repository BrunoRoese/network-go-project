package cmd

import (
	"github.com/BrunoRoese/socket/pkg/client"
	"github.com/spf13/cobra"
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
			cmd.Printf("IP: %s, Port: %d", client.Ip, client.Port)
		}
	},
}

func init() {
	rootCmd.AddCommand(devicesCmd)
}
