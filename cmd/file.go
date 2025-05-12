package cmd

import (
	"github.com/BrunoRoese/socket/pkg/protocol/service"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
)

var (
	filePath string
)

var fileCmd = &cobra.Command{
	Use:   "sendfile",
	Short: "Indicates the start of a file transfer to a client",
	Long:  `Indicates the start of a file transfer to a client.`,
	Run: func(cmd *cobra.Command, args []string) {
		fileService := service.FileService{}

		err := fileService.StartTransfer(clientIp, filePath)

		if err != nil {
			slog.Error("Error starting file transfer", slog.String("error", err.Error()))
			return
		}
	},
}

func init() {
	fileCmd.PersistentFlags().StringVar(&filePath, "path", "", "The directory of the file to be sent")
	err := fileCmd.MarkPersistentFlagRequired("file")

	fileCmd.PersistentFlags().StringVar(&clientIp, "ip", "", "The client ip to send the file")
	err = fileCmd.MarkPersistentFlagRequired("ip")

	rootCmd.AddCommand(fileCmd)

	if err != nil {
		slog.Error("Error marking flag as required", slog.String("error", err.Error()))
		os.Exit(1)
	}
}
