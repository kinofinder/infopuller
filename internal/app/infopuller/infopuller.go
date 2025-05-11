package infopuller

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	infopullerpb "github.com/kinofinder/proto/gen/go/infopuller"

	"infopuller/internal/client"
	"infopuller/internal/utils/config"
)

var (
	ErrInternal = status.Error(codes.Internal, "internal error")

	ErrClientFailed = fmt.Errorf("client failed")
	ErrUnmarshaling = fmt.Errorf("unmarshaling failed")
)

type App struct {
	Server *grpc.Server

	Log *slog.Logger

	Config config.Config
}

func New(log *slog.Logger, client *client.Client, c config.Config) *App {
	grpcs := grpc.NewServer()

	infopullerpb.RegisterInfoPullerServer(grpcs, &Handlers{
		Service: &Service{
			Client: client,

			Log: log,
		},

		Log: log,
	})

	return &App{
		Server: grpcs,

		Log: log,

		Config: c,
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
	info, err := h.Service.Random()
	if err != nil {
		// TODO: LOG ERROR

		return nil, ErrInternal
	}

	return buildResponse(info), nil
}

func buildResponse(info *Info) *infopullerpb.RandomResponse {
	var genres []*infopullerpb.Genre

	for _, val := range info.Genres {
		genre, ok := val["name"]
		if ok {
			genres = append(genres, &infopullerpb.Genre{Name: genre})
		}
	}

	var countries []*infopullerpb.Country

	for _, val := range info.Countries {
		country, ok := val["name"]
		if ok {
			countries = append(countries, &infopullerpb.Country{Name: country})
		}
	}

	return &infopullerpb.RandomResponse{
		Name:        info.Name,
		Year:        info.Year,
		Description: info.Description,
		Length:      info.Length,
		Poster:      info.Poster["url"],
		Genres:      genres,
		Countries:   countries,
	}
}

type Service struct {
	UnimplementedService

	Client *client.Client

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
	data, err := s.Client.Random()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrClientFailed, err)
	}

	var info *Info

	err = json.Unmarshal(data, &info)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrUnmarshaling, err)
	}

	return info, nil
}

type UnimplementedService struct{}

func (u *UnimplementedService) Random() (*Info, error) {
	return &Info{}, nil
}
