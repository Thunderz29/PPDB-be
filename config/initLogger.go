package config

import (
	"log/slog"
	"os"
)

func InitLogger() {
	// Create text handler or json handler
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	// Set default slog output
	logger := slog.New(handler)
	slog.SetDefault(logger)

	slog.Info("Slog logger initialized")
}