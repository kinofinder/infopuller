package infopuller

import (
	"infopuller/internal/utils/config"
	"log/slog"
	"net"

	"google.golang.org/grpc"

	infopullerpb "github.com/kinofinder/proto/gen/go/infopuller"
)

type App struct {
	Server *grpc.Server

	Log *slog.Logger

	Config config.Config
}

func New(log *slog.Logger, config config.Config) *App {
	grpcs := grpc.NewServer()

	infopullerpb.RegisterInfoPullerServer(grpcs, &Handlers{})

	return &App{
		Server: grpcs,

		Log: log,

		Config: config,
	}
}

func (a *App) Run() error {
	l, err := net.Listen(a.Config.Network, a.Config.Address)
	if err != nil {
		return err
	}

	err = a.Server.Serve(l)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) Shutdown() {
	a.Server.GracefulStop()
}

type Handlers struct {
	infopullerpb.UnimplementedInfoPullerServer

	Log *slog.Logger
}
