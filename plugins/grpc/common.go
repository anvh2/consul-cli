package rpc

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/anvh2/consul-cli/plugins/consul"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
)

// ShutdownHook -
type ShutdownHook func()

// BaseGrpcService ...
type BaseGrpcService struct {
	server       *grpc.Server
	port         int
	listener     net.Listener
	grpcRegister GrpcRegister
	hooks        []ShutdownHook
	done         chan error
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

	if s.server == nil {
		s.server = grpc.NewServer()
	}

	s.grpcRegister(s.server)
	fmt.Println("Server is running on port: ", s.port)
	go s.server.Serve(s.listener)

	sigs := make(chan os.Signal, 1)
	s.done = make(chan error, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		s.runHook()
		s.done <- s.Shutdown()
	}()
	err = <-s.done
	fmt.Println()
	return err
}

func (s *BaseGrpcService) runHook() {
	for _, hook := range s.hooks {
		defer hook()
	}
}

// Shutdown -
func (s *BaseGrpcService) Shutdown() error {
	if s.listener != nil {
		if err := s.listener.Close(); err != nil {
			return err
		}
		s.listener = nil
	}
	return nil
}

// AddShutdownHook -
func (s *BaseGrpcService) AddShutdownHook(fn ShutdownHook) {
	s.hooks = append(s.hooks, fn)
}

// RegisterWithConsul ...
func (s *BaseGrpcService) RegisterWithConsul(config *consul.Config) error {
	return consul.Register(config)
}

// DeRegisterFromConsul ...
func (s *BaseGrpcService) DeRegisterFromConsul(id string) {
	consul.DeRegister(id)
}

// RegisterHealthCheck -
func (s *BaseGrpcService) RegisterHealthCheck() {
	if s.server == nil {
		s.server = grpc.NewServer()
	}
	grpc_health_v1.RegisterHealthServer(s.server, &consul.HealthImpl{})
}
