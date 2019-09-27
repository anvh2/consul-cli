package rpc

import (
	"fmt"
	"net"

	"github.com/anvh2/consul-cli/plugins/consul"
	"google.golang.org/grpc"
)

// BaseGrpcService ...
type BaseGrpcService struct {
	port         int
	listener     net.Listener
	grpcRegister GrpcRegister
}

// GrpcRegister ...
type GrpcRegister func(server *grpc.Server)

// NewGrpcServer ...
func NewGrpcServer(grpcRegister GrpcRegister) *BaseGrpcService {
	return &BaseGrpcService{
		grpcRegister: grpcRegister,
	}
}

// Run ...
func (s *BaseGrpcService) Run(port int) error {
	s.port = port
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		return err
	}
	s.listener = lis

	grpcServer := grpc.NewServer()

	return grpcServer.Serve(s.listener)
}

// RegisterWithConsul ...
func (s *BaseGrpcService) RegisterWithConsul(config *consul.Config) error {
	return consul.Register(config)
}

// DeRegisterFromConsul ...
func (s *BaseGrpcService) DeRegisterFromConsul(id string) {
	consul.DeRegister(id)
}
