package user

import (
	"context"
	"fmt"

	pbCounter "github.com/anvh2/consul-cli/grpc-gen/counter"
	pb "github.com/anvh2/consul-cli/grpc-gen/user"
	"github.com/anvh2/consul-cli/plugins/consul"
	rpc "github.com/anvh2/consul-cli/plugins/grpc"
	"google.golang.org/grpc"
)

// Server ...
type Server struct {
	counterClient pbCounter.CounterPointServiceClient
}

// NewServer ...
func NewServer() *Server {

	return &Server{}
}

// Run ...
func (s *Server) Run() error {
	server := rpc.NewGrpcServer(s.registerServer)

	port := 55216
	config := consul.Config{
		ID:      "user",
		Name:    "UserService",
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
	pb.RegisterUserServiceServer(server, s)
}

// Login ...
func (s *Server) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	return nil, nil
}

// Register ...
func (s *Server) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	return nil, nil
}
