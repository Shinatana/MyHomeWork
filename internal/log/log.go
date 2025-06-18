package log

import (
	"log/slog"
	"os"
)

func NewLog(format, level string) *slog.Logger {
	var slogOpts slog.HandlerOptions
	switch level {
	case "debug":
		slogOpts = slog.HandlerOptions{Level: slog.LevelDebug}
	case "error":
		slogOpts = slog.HandlerOptions{Level: slog.LevelError}
	case "warn":
		slogOpts = slog.HandlerOptions{Level: slog.LevelWarn}
	case "info":
		slogOpts = slog.HandlerOptions{Level: slog.LevelInfo}
	default:
		slogOpts = slog.HandlerOptions{Level: slog.LevelWarn}
	}

	if format == "text" {
		return slog.New(slog.NewTextHandler(os.Stdout, &slogOpts))
	}

	return slog.New(slog.NewJSONHandler(os.Stdout, &slogOpts))
}
