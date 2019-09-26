package counter

import (
	"context"
	"fmt"

	pb "github.com/anvh2/consul-demo/grpc-gen/counter"
	"github.com/anvh2/consul-demo/plugins/consul"
	rpc "github.com/anvh2/consul-demo/plugins/grpc"
	"google.golang.org/grpc"
)

// Server ...
type Server struct{}

// NewServer ...
func NewServer() *Server {
	return &Server{}
}

// Run ...
func (s *Server) Run() error {
	server := rpc.NewGrpcServer(s.registerServer)

	port := 55216
	config := consul.Config{
		ID:      "counter1",
		Name:    "CounterService",
		Tags:    []string{"DEV"},
		Address: "127.0.0.1",
		Port:    port,
	}
	err := server.RegisterWithConsul(&config)
	if err != nil {
		fmt.Println("Can't register service")
	}

	return server.Run(port)
}

func (s *Server) registerServer(server *grpc.Server) {
	pb.RegisterCounterPointServiceServer(server, s)
}

// IncreasePoint ...
func (s *Server) IncreasePoint(ctx context.Context, req *pb.IncreaseRequest) (*pb.IncreaseResponse, error) {
	return nil, nil
}

// DecreasePoint ...
func (s *Server) DecreasePoint(ctx context.Context, req *pb.DecreaseRequest) (*pb.DecreaseResponse, error) {
	return nil, nil
}
