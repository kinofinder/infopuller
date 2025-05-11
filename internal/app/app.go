package app

import (
	"fmt"
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
)

type App struct {
	InfoPuller *infopuller.App
	Client     *client.Client

	Logger *logger.Logger

	Config config.Config
}

func New() (*App, error) {
	config := config.New()

	logger, err := logger.New(config)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrLoggerInitialization, err)
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
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	errChan := make(chan error, 1)

	// TODO: DEBUG LOG APP START

	go func() {
		err := a.InfoPuller.Run()
		if err != nil {
			errChan <- err
		}
	}()

	select {
	case <-sigChan:
		// TODO: INFO LOG SIGNAL HERE

	case err := <-errChan:
		// TODO: LOG ERROR HERE

		fmt.Println(err) // TEMP SOLUTION
	}

	a.shutdown()

	// TODO: DEBUG LOG APP STOP
}

func (a *App) shutdown() {
	a.Client.Shutdown()
	a.InfoPuller.Shutdown()
	a.Logger.Shutdown()
}
