package integration

import (
	"context"
	"os"
	"testing"

	pb "github.com/anvh2/consul-cli/grpc-gen/counter"
	"github.com/anvh2/consul-cli/plugins/consul"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

var client pb.CounterPointServiceClient

func TestMain(m *testing.M) {
	r, err := consul.NewResolver("CounterService", "DEV")
	if err != nil {
		panic(err)
	}

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	opts = append(opts, grpc.WithBalancer(grpc.RoundRobin(r)))

	conn, err := grpc.Dial("", opts...)
	if err != nil {
		panic(err)
	}
	client = pb.NewCounterPointServiceClient(conn)

	os.Exit(m.Run())
}

func TestIncreasePoint(t *testing.T) {
	res, err := client.IncreasePoint(context.Background(), &pb.IncreaseRequest{
		Data: &pb.PointData{
			UserID: 1,
			Amount: 100,
		},
	})

	assert.Nil(t, err)
	assert.Equal(t, int64(1), res.Code)
}
