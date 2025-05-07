package config

import (
	"log/slog"
	"os"
)

func SetupLogger() {
	logFile, err := os.OpenFile("resources/app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

	if err != nil {
		slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})))
		return
	}

	slog.SetDefault(slog.New(slog.NewJSONHandler(logFile, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})))
}
