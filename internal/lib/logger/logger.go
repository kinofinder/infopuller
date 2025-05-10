package logger

import (
	"log/slog"
	"os"

	"infopuller/internal/utils/config"
)

type Logger struct {
	Log *slog.Logger
}

func New(c config.Config) *Logger {
	var log *slog.Logger

	switch c.LogMode {
	case "local":
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))

	case "dev":
		log = slog.New(slog.NewJSONHandler(nil, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))

	case "prod":
		// TODO: LOGGER FOR PROD
	}

	return &Logger{
		Log: log,
	}
}
