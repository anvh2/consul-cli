package transfer

import (
	"context"

	pbCounter "github.com/anvh2/consul-cli/grpc-gen/counter"
	pb "github.com/anvh2/consul-cli/grpc-gen/transfer"
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
func (s *Server) Run(port int) error {
	server := rpc.NewGrpcServer(s.registerServer)

	// id, _ := uuid.NewV4()
	// idstr := fmt.Sprintf("transfer-%s", id.String())
	// config := consul.Config{
	// 	ID:      idstr,
	// 	Name:    "TransferService",
	// 	Tags:    []string{"DEV"},
	// 	Address: "127.0.0.1",
	// 	Port:    port,
	// }
	// err := server.RegisterWithConsul(&config)
	// if err != nil {
	// 	fmt.Println("Can't register service")
	// }
	// server.RegisterHealthCheck()

	// server.AddShutdownHook(func() {
	// 	server.DeRegisterFromConsul(idstr)
	// })

	return server.Run(port)
}

func (s *Server) registerServer(server *grpc.Server) {
	pb.RegisterTransferPointServiceServer(server, s)
}

// TransferPoint ...
func (s *Server) TransferPoint(ctx context.Context, req *pb.TransferRequest) (*pb.TransferResponse, error) {
	return nil, nil
}
