package tests

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"

	infopullerpb "github.com/kinofinder/proto/gen/go/infopuller"

	"infopuller/tests/suite"
)

func TestRandom_Functional(t *testing.T) {
	cases := []struct {
		name         string
		expectedResp bool
		expectedErr  error
		expectedCode codes.Code
	}{
		{
			name:         "happy case",
			expectedResp: true,
			expectedErr:  nil,
			expectedCode: codes.OK,
		},
	}

	suite := suite.New(t)
	go suite.App.Run()

	for _, cs := range cases {
		suite.T.Run(cs.name, func(t *testing.T) {
			resp, err := suite.Client.Random(context.Background(), &infopullerpb.RandomRequest{})
			if err != nil {
				assert.Equal(t, cs.expectedErr, err)
			}

			if cs.expectedResp {
				assert.NotEqual(t, "", resp.String())
			} else {
				assert.Equal(t, "", resp.String())
			}
		})
	}
}
