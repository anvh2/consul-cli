package integration

import (
	"context"
	"fmt"
	"testing"

	pb "github.com/anvh2/consul-cli/grpc-gen/counter"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

func getConnDev() *grpc.ClientConn {
	conn, err := grpc.Dial(":55215", grpc.WithInsecure())
	if err != nil {
		fmt.Println("Can't connect to server", err)
	}

	return conn
}

var client = pb.NewCounterPointServiceClient(getConnDev())

func TestIncreasePoint(t *testing.T) {
	_, err := client.IncreasePoint(context.Background(), &pb.IncreaseRequest{
		Data: &pb.PointData{
			UserID: 1,
			Amount: 100,
		},
	})

	assert.Nil(t, err)
}
