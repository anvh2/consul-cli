package transfer

import (
	"context"
	"fmt"

	pbCounter "github.com/anvh2/consul-cli/grpc-gen/counter"
	pb "github.com/anvh2/consul-cli/grpc-gen/transfer"
	"github.com/anvh2/consul-cli/plugins/consul"
	rpc "github.com/anvh2/consul-cli/plugins/grpc"
	"github.com/hashicorp/consul/api"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc"
)

// Server ...
type Server struct {
	counterClient pbCounter.CounterPointServiceClient
	consulAgent   *api.Agent
}

// NewServer ...
func NewServer() *Server {
	r, err := consul.NewResolver("CounterService", "DEV")
	if err != nil {
		fmt.Println(err)
	}

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	opts = append(opts, grpc.WithBalancer(grpc.RoundRobin(r)))

	conn, err := grpc.Dial("", opts...)
	if err != nil {
		fmt.Println(err)
	}

	counterClient := pbCounter.NewCounterPointServiceClient(conn)

	return &Server{
		counterClient: counterClient,
	}
}

// Run ...
func (s *Server) Run(port int) error {
	server := rpc.NewGrpcServer(s.registerServer)

	id, _ := uuid.NewV4()
	config := consul.Config{
		ID:      fmt.Sprintf("transfer-%s", id.String()),
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
	_, err := s.counterClient.IncreasePoint(ctx, &pbCounter.IncreaseRequest{
		Data: &pbCounter.PointData{
			UserID: req.ToID,
			Amount: req.Amount,
		},
	})
	if err != nil {
		return &pb.TransferResponse{
			Code:    -1,
			Message: err.Error(),
		}, err
	}

	_, err = s.counterClient.DecreasePoint(ctx, &pbCounter.DecreaseRequest{
		UserID: req.FromID,
		Amount: req.Amount,
	})
	if err != nil {
		return &pb.TransferResponse{
			Code:    -1,
			Message: err.Error(),
		}, err
	}

	return &pb.TransferResponse{
		Code:    1,
		Message: "OK",
	}, nil
}
