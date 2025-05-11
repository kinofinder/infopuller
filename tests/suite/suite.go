package suite

import (
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	infopullerpb "github.com/kinofinder/proto/gen/go/infopuller"

	"infopuller/internal/app"
	"infopuller/internal/app/infopuller"
	"infopuller/internal/lib/logger"
	"infopuller/internal/utils/config"
)

type Suite struct {
	*testing.T

	App    *app.App
	Client infopullerpb.InfoPullerClient
}

func New(t *testing.T, service infopuller.Servicer, log *logger.Logger, c *config.Config) *Suite {
	t.Helper()
	t.Parallel()

	grpcServer := grpc.NewServer()

	infopullerpb.RegisterInfoPullerServer(grpcServer, &infopuller.Handlers{
		Service: service,

		Log: log.Logger,

		Config: c,
	})

	infopuller := &infopuller.App{
		Server: grpcServer,

		Log: log.Logger,

		Config: c,
	}

	app := &app.App{
		InfoPuller: infopuller,

		Logger: log,

		Config: c,
	}

	conn, err := grpc.NewClient(c.ServerAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	grpcClient := infopullerpb.NewInfoPullerClient(conn)

	t.Cleanup(func() {
		app.InfoPuller.Shutdown()
	})

	return &Suite{
		T: t,

		App:    app,
		Client: grpcClient,
	}
}
