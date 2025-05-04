package cmd

import (
	"github.com/BrunoRoese/socket/config"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "socket",
	Short: "This is the socket application, start this and it will expose a port for the udp connection",
	Long:  `By exposing the port, it will be possible to send and receive data from the network.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.socket.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	config.SetupLogger()
}
