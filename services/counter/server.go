package counter

import (
	"context"
	"net"

	pb "github.com/anvh2/consul-cli/grpc-gen/counter"
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
	ports := []int{55210, 55211, 55212}
	ips := []net.IP{[]byte(""), []byte(""), []byte("")}
	server, err := rpc.NewServerWithExternalLoadBalancer(s.registerServer, 3, ports, ips)
	if err != nil {
		return err
	}

	server.Run()

	return nil
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
