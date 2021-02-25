package test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	common "onepass.app/facility/hts/common"
	"onepass.app/facility/hts/facility"
)

// Define the suite, and absorb the built-in basic suite
// functionality from testify - including a T() method which
// returns the current testing context
type ExampleTestSuite struct {
	suite.Suite
	Client facility.FacilityServiceClient
}

func (suite *ExampleTestSuite) SetupTest() {

	facilityPath := os.Getenv("HTS_SVC_FACILITY")
	opts := []grpc.DialOption{grpc.WithInsecure()}
	const deadlineSeconds = 5
	timeout := time.Duration(deadlineSeconds) * time.Second

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	connAccount, dialError := grpc.DialContext(ctx, facilityPath, opts...)
	if dialError != nil {
		panic(dialError)
	}
	suite.Client = facility.NewFacilityServiceClient(connAccount)
}

func (suite *ExampleTestSuite) TestPing() {
	assert := assert.New(suite.T())
	_, err := suite.Client.Ping(context.Background(), &empty.Empty{})
	if err != nil {
		panic(err)
	}
	result, err := suite.Client.Ping(context.Background(), &empty.Empty{})
	expected := &common.Result{IsOk: true}
	assert.Nil(err)
	assert.True(proto.Equal(result, expected))
}

func TestExampleTestSuite(t *testing.T) {
	suite.Run(t, new(ExampleTestSuite))
}
