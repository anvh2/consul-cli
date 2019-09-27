package transfer

import (
	"context"
	"fmt"
	"log"
	"net"
	"strconv"

	pb "github.com/anvh2/consul-cli/grpc-gen/transfer"
	"github.com/hashicorp/consul/api"
	"github.com/olivere/randport"
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
	addr := fmt.Sprintf("127.0.0.1:%d", randport.Get())
	address, portstr, err := net.SplitHostPort(addr)
	if err != nil {
		log.Fatal(err)
	}
	port, err := strconv.Atoi(portstr)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Port in use: ", port)

	cli, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		fmt.Println(err)
	}

	reg := &api.AgentServiceRegistration{
		ID:      "transfer1",
		Name:    "Transfer Service",
		Tags:    []string{"Dev", "Test"},
		Address: address,
		Port:    port,
	}

	err = cli.Agent().ServiceRegister(reg)
	if err != nil {
		fmt.Println(err)
	}
	defer cli.Agent().ServiceDeregister(reg.ID)

	// create a listener on TCP port 7777
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// create a gRPC server object
	grpcServer := grpc.NewServer()
	// attach the Ping service to the server
	pb.RegisterTransferPointServiceServer(grpcServer, s)

	return grpcServer.Serve(lis)
}

// TransferPoint ...
func (s *Server) TransferPoint(ctx context.Context, req *pb.TransferRequest) (*pb.TransferResponse, error) {
	return nil, nil
}
