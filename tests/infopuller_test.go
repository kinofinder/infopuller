package tests

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	infopullerpb "github.com/kinofinder/proto/gen/go/infopuller"

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

	suite := suite.New(t)
	go suite.App.Run()

	for _, cs := range cases {
		suite.T.Run(cs.name, func(t *testing.T) {
			resp, err := suite.Client.Random(context.Background(), &infopullerpb.RandomRequest{})
			st, ok := status.FromError(err)
			if ok {
				assert.Equal(t, cs.expectedCode, st)
			}

			if cs.expectedResp {
				assert.NotEqual(t, "", resp.String())
			} else {
				assert.Equal(t, "", resp.String())
			}
		})
	}
}
