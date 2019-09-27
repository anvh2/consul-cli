package rpc

import (
	"crypto/tls"

	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/credentials"
)

// NewGrpcClient ...
func NewGrpcClient(address string) (*grpc.ClientConn, error) {
	var opts []grpc.DialOption
	tlsCreds := credentials.NewTLS(&tls.Config{InsecureSkipVerify: true})
	opts = append(opts, grpc.WithTransportCredentials(tlsCreds))
	opts = append(opts, grpc.WithBalancerName(roundrobin.Name))

	return grpc.Dial(address, opts...)
}
