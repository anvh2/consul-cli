package user

import (
	"context"
	"fmt"

	pbTransfer "github.com/anvh2/consul-cli/grpc-gen/transfer"
	pb "github.com/anvh2/consul-cli/grpc-gen/user"
	"github.com/anvh2/consul-cli/plugins/consul"
	rpc "github.com/anvh2/consul-cli/plugins/grpc"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc"
)

// Server ...
type Server struct {
	transferClient pbTransfer.TransferPointServiceClient
}

// NewServer ...
func NewServer() *Server {
	r, err := consul.NewResolver("TransferService", "DEV")
	if err != nil {
		fmt.Println(err)
	}

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	opts = append(opts, grpc.WithBalancer(grpc.RoundRobin(r))) // load balance with roundrobin strategy

	conn, err := grpc.Dial("", opts...)
	if err != nil {
		fmt.Println(err)
	}

	transferClient := pbTransfer.NewTransferPointServiceClient(conn)
	return &Server{
		transferClient: transferClient,
	}
}

// Run ...
func (s *Server) Run(port int) error {
	server := rpc.NewGrpcServer(s.registerServer)

	id, _ := uuid.NewV4()
	idstr := fmt.Sprintf("user-%s", id.String())
	config := consul.Config{
		ID:      idstr,
		Name:    "UserService",
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

// Transfer ...
func (s *Server) Transfer(ctx context.Context, req *pb.TransferRequest) (*pb.TransferResponse, error) {
	res, err := s.transferClient.TransferPoint(ctx, &pbTransfer.TransferRequest{
		ToID:   req.ToID,
		FromID: req.FromID,
		Amount: req.Amount,
	})

	if err != nil {
		// log
		return &pb.TransferResponse{
			Code:    -1,
			Message: err.Error(),
		}, nil
	}

	return &pb.TransferResponse{
		Code:    res.Code,
		Message: res.Message,
	}, nil
}
