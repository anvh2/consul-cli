package counter

import (
	"context"
	"fmt"

	pb "github.com/anvh2/consul-cli/grpc-gen/counter"
	"github.com/anvh2/consul-cli/plugins/consul"
	rpc "github.com/anvh2/consul-cli/plugins/grpc"
	"github.com/anvh2/consul-cli/storages/mysql"
	"google.golang.org/grpc"
)

// Server ...
type Server struct {
	counterDb *mysql.CounterDb
}

// NewServer ...
func NewServer(counterDb *mysql.CounterDb) *Server {
	return &Server{
		counterDb: counterDb,
	}
}

// Run ...
func (s *Server) Run() error {
	server := rpc.NewGrpcServer(s.registerServer)

	port := 55215
	config := consul.Config{
		ID:      "counter",
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
	s.counterDb.IncreasePoint(req.Data)
	return nil, nil
}

// DecreasePoint ...
func (s *Server) DecreasePoint(ctx context.Context, req *pb.DecreaseRequest) (*pb.DecreaseResponse, error) {
	return nil, nil
}
