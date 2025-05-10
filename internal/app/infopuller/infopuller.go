package infopuller

import (
	"context"
	"log/slog"
	"net"

	"google.golang.org/grpc"

	infopullerpb "github.com/kinofinder/proto/gen/go/infopuller"

	"infopuller/internal/utils/config"
)

type App struct {
	Server *grpc.Server

	Log *slog.Logger

	Config config.Config
}

func New(log *slog.Logger, config config.Config) *App {
	grpcs := grpc.NewServer()

	infopullerpb.RegisterInfoPullerServer(grpcs, &Handlers{
		Log: log,
	})

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

	// TODO: DEBUG LOG SERVER START

	err = a.Server.Serve(l)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) Shutdown() {
	a.Server.GracefulStop()

	// TODO: DEBUG LOG SERVER STOP
}

type Servicer interface {
	Random() (*Info, error)
}

type Handlers struct {
	infopullerpb.UnimplementedInfoPullerServer

	Service Servicer

	Log *slog.Logger
}

func (h *Handlers) Random(ctx context.Context, req *infopullerpb.RandomRequest) (*infopullerpb.RandomResponse, error) {
	return &infopullerpb.RandomResponse{}, nil
}

type Service struct {
	UnimplementedService

	Log *slog.Logger
}

type Info struct {
	Name        string              `json:"name"`
	Year        int32               `json:"year"`
	Description string              `json:"description"`
	Rating      []float64           `json:"rating"`
	Length      int32               `json:"movieLength"`
	Poster      map[string]string   `json:"poster"`
	Genres      []map[string]string `json:"genres"`
	Countries   []map[string]string `json:"countries"`
}

func (s *Service) Random() (*Info, error) {
	return &Info{}, nil
}

type UnimplementedService struct{}

func (u *UnimplementedService) Random() (*Info, error) {
	return &Info{}, nil
}
