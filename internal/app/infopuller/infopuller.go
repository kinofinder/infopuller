package infopuller

import (
	"google.golang.org/grpc"

	infopullerpb "github.com/kinofinder/proto/gen/go/infopuller"
)

type App struct {
}

func New() *App {
	grpcs := grpc.NewServer()

	infopullerpb.RegisterInfoPullerServer(grpcs, &Handlers{})

	return &App{}
}

func (a *App) Run() {}

func (a *App) Shutdown() {}

type Handlers struct {
	infopullerpb.UnimplementedInfoPullerServer
}
