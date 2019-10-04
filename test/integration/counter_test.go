package integration

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	pb "github.com/anvh2/consul-cli/grpc-gen/counter"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc/resolver/manual"
	"google.golang.org/grpc/status"
)

var (
	client pb.CounterPointServiceClient
	conn   *grpc.ClientConn
	err    error
)

type test struct {
	servers   []*grpc.Server
	addresses []string
}

func TestMain(m *testing.M) {
	r, cleanup := manual.GenerateAndRegisterManualResolver()
	defer cleanup()

	conn, err := grpc.Dial(r.Scheme()+":///test.server", grpc.WithInsecure(), grpc.WithBalancerName(roundrobin.Name))
	if err != nil {
		fmt.Println(err)
	}

	counterClient := pb.NewCounterPointServiceClient(conn)

	// The first RPC should fail because there's no address.
	ctx, cancle := context.WithTimeout(context.Background(), time.Millisecond)
	defer cancle()
	if _, err := counterClient.IncreasePoint(ctx, &pb.IncreaseRequest{}); err == nil || status.Code(err) != codes.DeadlineExceeded {
		fmt.Printf("EmptyCall() = _, %v, want _, DeadlineExceeded\n", err)
	}

	var resolvedAddrs []resolver.Address
	for i := 0; i < 3; i++ {
		resolvedAddrs = append(resolvedAddrs, resolver.Address{Addr: test.addresses[i]})
	}

	r.UpdateState(resolver.State{Addresses: resolvedAddrs})

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
