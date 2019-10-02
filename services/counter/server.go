package counter

import (
	"context"
	"fmt"

	pb "github.com/anvh2/consul-cli/grpc-gen/counter"
	"github.com/anvh2/consul-cli/plugins/consul"
	rpc "github.com/anvh2/consul-cli/plugins/grpc"
	"github.com/anvh2/consul-cli/storages/mysql"
	uuid "github.com/satori/go.uuid"
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
func (s *Server) Run(port int) error {
	server := rpc.NewGrpcServer(s.registerServer)

	id, _ := uuid.NewV4()
	idstr := fmt.Sprintf("counter-%s", id.String())
	config := consul.Config{
		ID:      idstr,
		Name:    "CounterService",
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
	pb.RegisterCounterPointServiceServer(server, s)
}

// IncreasePoint ...
func (s *Server) IncreasePoint(ctx context.Context, req *pb.IncreaseRequest) (*pb.IncreaseResponse, error) {
	err := s.counterDb.IncreasePoint(req.Data)
	if err != nil {
		return &pb.IncreaseResponse{
			Code:    -1,
			Message: err.Error(),
		}, nil
	}

	return &pb.IncreaseResponse{
		Code:    1,
		Message: "OK",
		Data:    req.Data,
	}, nil
}

// DecreasePoint ...
func (s *Server) DecreasePoint(ctx context.Context, req *pb.DecreaseRequest) (*pb.DecreaseResponse, error) {
	return &pb.DecreaseResponse{}, nil
}
