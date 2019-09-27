package transfer

import (
	"context"
	"fmt"

	pb "github.com/anvh2/consul-cli/grpc-gen/transfer"
	"github.com/anvh2/consul-cli/plugins/consul"
	rpc "github.com/anvh2/consul-cli/plugins/grpc"
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

	port := 55217
	config := consul.Config{
		ID:      "transfer",
		Name:    "TransferService",
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
	pb.RegisterTransferPointServiceServer(server, s)
}

// TransferPoint ...
func (s *Server) TransferPoint(ctx context.Context, req *pb.TransferRequest) (*pb.TransferResponse, error) {
	return nil, nil
}
