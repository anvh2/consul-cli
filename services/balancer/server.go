package balancer

import (
	"context"

	"google.golang.org/grpc/balancer/grpclb/grpc_lb_v1"
	"google.golang.org/grpc/metadata"
)

// LoadBalancerServer -
type LoadBalancerServer struct {
	serverList grpc_lb_v1.ServerList
}

// NewLoadBalancerServer -
func NewLoadBalancerServer() *LoadBalancerServer {
	return &LoadBalancerServer{}
}

// BalanceLoad -
func (lbs *LoadBalancerServer) BalanceLoad(grpc_lb_v1.LoadBalancer_BalanceLoadServer) error {
	return nil
}

// LoadBalancer -
type LoadBalancer struct{}

// Send -
func (lb *LoadBalancer) Send(*grpc_lb_v1.LoadBalanceResponse) error {
	return nil
}

// Recv -
func (lb *LoadBalancer) Recv() (*grpc_lb_v1.LoadBalanceRequest, error) {
	return nil, nil
}

// SetHeader -
func (lb *LoadBalancer) SetHeader(metadata.MD) error {
	return nil
}

// SendHeader -
func (lb *LoadBalancer) SendHeader(metadata.MD) error {
	return nil
}

// SetTrailer -
func (lb *LoadBalancer) SetTrailer(metadata.MD) {

}

// Context -
func (lb *LoadBalancer) Context() context.Context {
	return nil
}

// SendMsg -
func (lb *LoadBalancer) SendMsg(m interface{}) error {
	return nil
}

// RecvMsg -
func (lb *LoadBalancer) RecvMsg(m interface{}) error {
	return nil
}
