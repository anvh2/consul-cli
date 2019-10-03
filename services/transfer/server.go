package transfer

import (
	"context"
	"fmt"

	"google.golang.org/grpc/balancer/grpclb/grpc_lb_v1"

	pbCounter "github.com/anvh2/consul-cli/grpc-gen/counter"
	pb "github.com/anvh2/consul-cli/grpc-gen/transfer"
	"github.com/anvh2/consul-cli/plugins/consul"
	rpc "github.com/anvh2/consul-cli/plugins/grpc"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc"
)

// Server ...
type Server struct {
	counterClient pbCounter.CounterPointServiceClient
	lbClient      grpc_lb_v1.LoadBalancer_BalanceLoadClient
}

// NewServer ...
func NewServer() *Server {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())

	// port of load balancer server
	conn, err := grpc.Dial(":55210", opts...)
	if err != nil {
		// log
		fmt.Println(err)
	}

	lbClient, err := grpc_lb_v1.NewLoadBalancerClient(conn).BalanceLoad(context.Background())
	if err != nil {
		// log
		fmt.Println(err)
	}

	return &Server{
		lbClient: lbClient,
	}
}

// Run ...
func (s *Server) Run(port int) error {
	server := rpc.NewGrpcServer(s.registerServer)

	id, _ := uuid.NewV4()
	idstr := fmt.Sprintf("transfer-%s", id.String())
	config := consul.Config{
		ID:      idstr,
		Name:    "TransferService",
		Tags:    []string{"DEV"},
		Address: "127.0.0.1",
		Port:    port,
	}
	err := server.RegisterWithConsul(&config)
	if err != nil {
		fmt.Println("Can't register service")
	}
	server.RegisterHealthCheck()

	server.AddShutdownHook(func() {
		server.DeRegisterFromConsul(idstr)
	})

	return server.Run(port)
}

func (s *Server) registerServer(server *grpc.Server) {
	pb.RegisterTransferPointServiceServer(server, s)
}

// TransferPoint ...
func (s *Server) TransferPoint(ctx context.Context, req *pb.TransferRequest) (*pb.TransferResponse, error) {
	return nil, nil
}
