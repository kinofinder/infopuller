package logger

import (
	"log/slog"
	"os"

	"infopuller/internal/utils/config"
)

type Logger struct {
	*slog.Logger

	LogFile *os.File
}

func New(c config.Config) (*Logger, error) {
	var log *slog.Logger
	var file *os.File

	switch c.LogMode {
	case "local":
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))

	case "dev":
		var err error

		file, err = loadLogFile(c.LogDirectory, "infopuller.log.json")
		if err != nil {
			return nil, err
		}

		log = slog.New(slog.NewJSONHandler(nil, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))

	case "prod":
		// TODO: LOGGER FOR PROD

	default:
		log = slog.New(slog.DiscardHandler)
	}

	// TODO: DEBUG LOG LOGGER START

	return &Logger{
		Logger: log,

		LogFile: file,
	}, nil
}

func loadLogFile(dir string, name string) (*os.File, error) {
	err := os.MkdirAll(dir, 0666)
	if err != nil {
		return nil, err
	}

	file, err := os.Create(name)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func (l *Logger) Shutdown() {
	if l.LogFile != nil {
		l.LogFile.Close()
	}
}
