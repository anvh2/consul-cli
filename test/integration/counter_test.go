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
	conn, err := grpc.Dial(":", grpc.WithInsecure())
	if err != nil {
		fmt.Println("Can't connect to server", err)
	}

	return conn
}

func TestConn(t *testing.T) {
	client := pb.NewCounterPointServiceClient(getConnDev())

	_, err := client.IncreasePoint(context.Background(), &pb.IncreaseRequest{})
	assert.Nil(t, err)
}
