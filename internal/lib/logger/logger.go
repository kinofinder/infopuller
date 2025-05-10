package logger

import (
	"log"
	"log/slog"
	"os"
)

type Logger struct {
	Log *slog.Logger
}

func New() *Logger {
	var log *log.Logger

	switch {
	case "local":
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
	}

	return &Logger{
		Log: log,
	}
}
