package config

import (
	"log/slog"
	"os"
)

func SetupLogger() {
	//logFilePath := "resources/app.log"
	//
	//if _, err := os.Stat(logFilePath); err == nil {
	//	err = os.Remove(logFilePath)
	//	if err != nil {
	//		slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
	//			Level: slog.LevelInfo,
	//		})))
	//		slog.Error("Failed to delete existing log file", "error", err)
	//		return
	//	}
	//}
	//
	//logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	//if err != nil {
	//	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
	//		Level: slog.LevelInfo,
	//	})))
	//	return
	//}

	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})))
}
