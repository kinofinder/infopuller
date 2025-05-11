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
	ErrInternal = status.Error(codes.Internal, "failed while invoking a method")

	ErrClientFailed = fmt.Errorf("client failed")
	ErrUnmarshaling = fmt.Errorf("unmarshaling failed")
)

type App struct {
	Server *grpc.Server

	Log *slog.Logger

	Config config.Config
}

func New(log *slog.Logger, client *client.Client, c config.Config) *App {
	const op = "infopuller.New()"

	log.Debug(
		"initiallizing infopuller",
		slog.String("op", op),
	)

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
	const op = "infopuller.Run()"

	a.Log.Debug(
		"starting infopuller server",
		slog.String("op", op),
	)

	l, err := net.Listen(a.Config.Network, a.Config.Address)
	if err != nil {
		a.Log.Debug(
			"failed to listen",
			slog.String("op", op),
			slog.Any("err", err),
		)

		return err
	}

	err = a.Server.Serve(l)
	if err != nil {
		a.Log.Debug(
			"failed to serve",
			slog.String("op", op),
			slog.Any("err", err),
		)

		return err
	}

	return nil
}

func (a *App) Shutdown() {
	const op = "infopuller.Shutdown()"

	a.Log.Debug(
		"stopping infopuller server",
		slog.String("op", op),
	)

	a.Server.GracefulStop()
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
	const op = "grpc.Random()"

	info, err := h.Service.Random()
	if err != nil {
		h.Log.Error(
			"failure in random method",
			slog.String("op", op),
			slog.Any("err", err),
		)

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
		Type:        info.Type,
		Year:        info.Year,
		Description: info.Description,
		Length:      info.Length,
		Poster:      info.Poster["url"],
		Genres:      genres,
		Countries:   countries,
	}
}

type Clienter interface {
	Random() ([]byte, error)
}

type Service struct {
	UnimplementedService

	Client Clienter

	Log *slog.Logger
}

type Info struct {
	Name        string              `json:"name"`
	Type        string              `json:"type"`
	Year        int32               `json:"year"`
	Description string              `json:"description"`
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
