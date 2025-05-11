package app

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"infopuller/internal/app/infopuller"
	"infopuller/internal/client"
	"infopuller/internal/lib/logger"
	"infopuller/internal/utils/config"
)

var (
	ErrLoggerInitialization = fmt.Errorf("logger failed to initialize")
	ErrConfigLoading        = fmt.Errorf("failed to load config")
)

type App struct {
	InfoPuller *infopuller.App
	Client     *client.Client

	Logger *logger.Logger

	Config config.Config
}

func New() (*App, error) {
	config, err := config.New()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrLoggerInitialization, err)
	}

	logger, err := logger.New(config)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrConfigLoading, err)
	}

	client := client.New(logger.Logger, config)

	infopuller := infopuller.New(logger.Logger, client, config)

	return &App{
		InfoPuller: infopuller,
		Client:     client,

		Logger: logger,

		Config: config,
	}, nil
}

func (a *App) Run() {
	const op = "app.Run()"

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	errChan := make(chan error, 1)

	a.Logger.Info(
		"starting an app",
		slog.String("op", op),
		slog.Any("c", a.Config),
	)

	go func() {
		err := a.InfoPuller.Run()
		if err != nil {
			errChan <- err
		}
	}()

	select {
	case sig := <-sigChan:
		a.Logger.Info(
			"received a signal to shutdown",
			slog.String("op", op),
			slog.Any("sig", sig),
		)

	case err := <-errChan:
		a.Logger.Error(
			"fatal error happened while running",
			slog.String("op", op),
			slog.Any("err", err),
		)
	}

	a.Logger.Info(
		"shutting down the app",
		slog.String("op", op),
	)

	a.shutdown()
}

func (a *App) shutdown() {
	a.Client.Shutdown()
	a.InfoPuller.Shutdown()
	a.Logger.Shutdown()
}
