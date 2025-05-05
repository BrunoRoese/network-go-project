package cmd

import (
	"github.com/BrunoRoese/socket/pkg/server/service"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
)

var (
	clientIp string
	talkMsg  string
)

var talkCmd = &cobra.Command{
	Use:   "talk",
	Short: "Send a message to a specific ip location",
	Long:  `Sends a message to a specific ip location by providing a uuid to identify the request`,
	Run:   runTalk,
}

func runTalk(cmd *cobra.Command, args []string) {
	service.Talk(clientIp, talkMsg)
}

func init() {
	talkCmd.PersistentFlags().StringVar(&clientIp, "ip", "", "Your network IP address")
	err := talkCmd.MarkPersistentFlagRequired("ip")

	talkCmd.PersistentFlags().StringVar(&talkMsg, "msg", "Hello, world!", "Your message")
	err = talkCmd.MarkPersistentFlagRequired("msg")

	rootCmd.AddCommand(talkCmd)

	if err != nil {
		slog.Error("Error marking flag as required", slog.String("error", err.Error()))
		os.Exit(1)
	}
}
