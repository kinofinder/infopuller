package suite

import (
	"os"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	infopullerpb "github.com/kinofinder/proto/gen/go/infopuller"

	"infopuller/internal/app"
	"infopuller/internal/app/infopuller"
	"infopuller/internal/client"
	"infopuller/internal/lib/logger"
	"infopuller/internal/utils/config"
)

type Suite struct {
	*testing.T

	App    *app.App
	Client infopullerpb.InfoPullerClient
}

func New(t *testing.T) *Suite {
	t.Helper()
	t.Parallel()

	os.Setenv("CONFIG_LOCATION", "test.env")
	defer os.Unsetenv("CONFIG_LOCATION")

	config, err := config.New()
	if err != nil {
		panic(err)
	}

	logger, err := logger.New(config)
	if err != nil {
		panic(err)
	}

	client := client.New(logger.Logger, config)

	app := &app.App{
		InfoPuller: infopuller.New(logger.Logger, client, config),
		Client:     client,

		Logger: logger,

		Config: config,
	}

	conn, err := grpc.NewClient(config.ServerAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	grpcClient := infopullerpb.NewInfoPullerClient(conn)

	t.Cleanup(func() {
		app.Client.Shutdown()
		app.InfoPuller.Shutdown()
		app.Logger.Shutdown()
	})

	return &Suite{
		T: t,

		App:    app,
		Client: grpcClient,
	}
}
