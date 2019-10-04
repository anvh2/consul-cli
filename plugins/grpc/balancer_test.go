package rpc

import (
	"context"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"google.golang.org/grpc/balancer/grpclb/grpc_lb_v1"
	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc/resolver/manual"

	"google.golang.org/grpc"

	pb "github.com/anvh2/consul-cli/grpc-gen/echo"
	"github.com/stretchr/testify/assert"
)

type echoTestServer struct{}

func (s *echoTestServer) Echo(ctx context.Context, req *pb.EchoRequest) (*pb.EchoResponse, error) {
	return &pb.EchoResponse{}, nil
}

var (
	lbServer   *GrpcServiceWithLB
	echoServer = &echoTestServer{}
	echoClient pb.EchoServiceClient
	Address    []string
)

func TestMain(m *testing.M) {
	var err error
	lbServer, err = NewServerWithExternalLoadBalancer(echoServer.registerService, "EchoService", 2, []int{55210, 55211})
	if err != nil {
		panic(err)
	}

	go lbServer.Run()
	time.Sleep(1 * time.Second)
	os.Exit(m.Run())
}

func (s *echoTestServer) registerService(server *grpc.Server) {
	pb.RegisterEchoServiceServer(server, s)
}

func TestRPCRemoteBalancerServer(t *testing.T) {
	conn, err := grpc.Dial(":55220", grpc.WithInsecure())
	defer conn.Close()
	assert.Nil(t, err)
	client := grpc_lb_v1.NewLoadBalancerClient(conn)

	stream, err := client.BalanceLoad(context.Background())
	assert.Nil(t, err)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		err := stream.Send(&grpc_lb_v1.LoadBalanceRequest{
			LoadBalanceRequestType: &grpc_lb_v1.LoadBalanceRequest_InitialRequest{
				InitialRequest: &grpc_lb_v1.InitialLoadBalanceRequest{
					Name: "EchoService",
				},
			},
		})
		assert.Nil(t, err)
	}()

	go func() {
		defer wg.Done()
		res, err := stream.Recv()
		assert.Nil(t, err)
		fmt.Println("RES:", res)
	}()

	wg.Wait()
}

func TestCallGRPC_LB(t *testing.T) {
	r, cleanup := manual.GenerateAndRegisterManualResolver()
	defer cleanup()
	// fmt.Println("RESULT: ", lbServer.beName)
	be := &grpc_lb_v1.Server{
		IpAddress:        []byte("localhost"),
		Port:             int32(lbServer.bePorts[0]),
		LoadBalanceToken: "iamatoken",
	}

	var bes []*grpc_lb_v1.Server
	bes = append(bes, be)
	serverList := &grpc_lb_v1.ServerList{
		Servers: bes,
	}

	lbServer.remoteBalancer.ServerList <- serverList

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	conn, err := grpc.DialContext(ctx, r.Scheme()+":///EchoService", grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial to the backend %v", err)
	}
	defer conn.Close()

	echoClient = pb.NewEchoServiceClient(conn)

	r.UpdateState(resolver.State{
		Addresses: []resolver.Address{
			{
				Addr:       lbServer.lbAddr,
				Type:       resolver.GRPCLB,
				ServerName: "EchoService",
			},
		},
	})

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err = echoClient.Echo(ctx, &pb.EchoRequest{})

	fmt.Println(err)
}
