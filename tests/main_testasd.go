package test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	common "onepass.app/facility/hts/common"
	"onepass.app/facility/hts/facility"
)

func TestSomething2qwe(t *testing.T) {
	assert := assert.New(t)
	facilityPath := os.Getenv("HTS_SVC_FACILITY")

	// Disable transport security is intentional
	opts := []grpc.DialOption{grpc.WithInsecure()}
	const deadlineSeconds = 5
	timeout := time.Duration(deadlineSeconds) * time.Second

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	connAccount, dialError := grpc.DialContext(ctx, facilityPath, opts...)
	if dialError != nil {
		panic(dialError)
	}
	facilityClient := facility.NewFacilityServiceClient(connAccount)
	_, err := facilityClient.Ping(context.Background(), &empty.Empty{})
	if err != nil {
		panic(err)
	}
	result, err := facilityClient.Ping(context.Background(), &empty.Empty{})
	expected := &common.Result{IsOk: true}
	assert.Nil(err)
	assert.Equal(result, expected)
}
