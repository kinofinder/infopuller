package tests

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	infopullerpb "github.com/kinofinder/proto/gen/go/infopuller"

	"infopuller/internal/app/infopuller"
	"infopuller/internal/lib/logger"
	"infopuller/internal/utils/config"
	"infopuller/tests/suite"
)

func TestRandom_Functional(t *testing.T) {
	cases := []struct {
		name         string
		expectedResp bool
		expectedCode codes.Code
	}{
		{
			name:         "happy case",
			expectedResp: true,
			expectedCode: codes.OK,
		},
	}

	os.Setenv("CONFIG_LOCATION", "test.env")
	defer os.Unsetenv("CONFIG_LOCATION")

	config, err := config.New()
	assert.NoError(t, err)

	logger, err := logger.New(config)
	assert.NoError(t, err)

	suite := suite.New(t, &infopuller.UnimplementedService{}, logger, config)

	go suite.App.Run()

	time.Sleep(time.Second * 2)

	for _, cs := range cases {
		suite.T.Run(cs.name, func(t *testing.T) {
			_, err := suite.Client.Random(context.Background(), &infopullerpb.RandomRequest{})
			st, ok := status.FromError(err)
			if ok {
				assert.Equal(t, cs.expectedCode, st.Code())
			}
		})
	}
}
