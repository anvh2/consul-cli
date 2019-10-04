package rpc

import (
	"fmt"
	"sync"
	"time"

	duration "github.com/golang/protobuf/ptypes/duration"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/grpclb/grpc_lb_v1"
)

// rpcStats is same as lbmpb.ClientStats, except that numCallsDropped is a map
// instead of a slice.
type rpcStats struct {
	// Only access the following fields atomically.
	numCallsStarted                        int64
	numCallsFinished                       int64
	numCallsFinishedWithClientFailedToSend int64
	numCallsFinishedKnownReceived          int64

	mu sync.Mutex
	// map load_balance_token -> num_calls_dropped
	numCallsDropped map[string]int64
}

// RemoteBalancer -
type RemoteBalancer struct {
	beServerName string
	ServerList   chan *grpc_lb_v1.ServerList
	statsDura    time.Duration
	quitc        chan struct{}
	rpcStats     *rpcStats
}

func newRPCStats() *rpcStats {
	return &rpcStats{
		numCallsDropped: make(map[string]int64),
	}
}

// NewRemoteBalancerServer -
func NewRemoteBalancerServer(beServerName string) *RemoteBalancer {
	return &RemoteBalancer{
		beServerName: beServerName,
		ServerList:   make(chan *grpc_lb_v1.ServerList, 1),
		quitc:        make(chan struct{}),
		rpcStats:     newRPCStats(),
	}
}

// BalanceLoad -
func (lbs *RemoteBalancer) BalanceLoad(stream grpc_lb_v1.LoadBalancer_BalanceLoadServer) error {
	req, err := stream.Recv()
	if err != nil {
		return err
	}
	initReq := req.GetInitialRequest()
	if initReq.Name != lbs.beServerName {
		// Log
		return err
	}
	resp := &grpc_lb_v1.LoadBalanceResponse{
		LoadBalanceResponseType: &grpc_lb_v1.LoadBalanceResponse_InitialResponse{
			InitialResponse: &grpc_lb_v1.InitialLoadBalanceResponse{
				ClientStatsReportInterval: &duration.Duration{
					Seconds: int64(lbs.statsDura.Seconds()),
					Nanos:   int32(lbs.statsDura.Nanoseconds() - int64(lbs.statsDura.Seconds()*1e9)),
				},
			},
		},
	}
	if err := stream.Send(resp); err != nil {
		// Log
		return err
	}

	// receiver
	go func() {
		for {
			var (
				req *grpc_lb_v1.LoadBalanceRequest
				err error
			)
			if req, err = stream.Recv(); err != nil {
				return
			}
			fmt.Println("RECEIVED REQUEST:", req)
			lbs.rpcStats.merge(req.GetClientStats())
		}
	}()
	// sender
	for {
		select {
		case v := <-lbs.ServerList:
			resp = &grpc_lb_v1.LoadBalanceResponse{
				LoadBalanceResponseType: &grpc_lb_v1.LoadBalanceResponse_ServerList{
					ServerList: v,
				},
			}
		case <-stream.Context().Done():
			return stream.Context().Err()
		}
		if err := stream.Send(resp); err != nil {
			return err
		}
	}
}

func (s *rpcStats) merge(cs *grpc_lb_v1.ClientStats) {
	// atomic.AddInt64(&s.numCallsStarted, cs.NumCallsStarted)
	// atomic.AddInt64(&s.numCallsFinished, cs.NumCallsFinished)
	// atomic.AddInt64(&s.numCallsFinishedWithClientFailedToSend, cs.NumCallsFinishedWithClientFailedToSend)
	// atomic.AddInt64(&s.numCallsFinishedKnownReceived, cs.NumCallsFinishedKnownReceived)
	// s.mu.Lock()
	// for _, perToken := range cs.CallsFinishedWithDrop {
	// 	s.numCallsDropped[perToken.LoadBalanceToken] += perToken.NumCalls
	// }
	// s.mu.Unlock()
}

// BalanceLoadClientStream -
type BalanceLoadClientStream struct {
	grpc.ClientStream
}

// Send -
func (x *BalanceLoadClientStream) Send(m *grpc_lb_v1.LoadBalanceResponse) error {
	return x.ClientStream.SendMsg(m)
}

// Recv -
func (x *BalanceLoadClientStream) Recv() (*grpc_lb_v1.LoadBalanceRequest, error) {
	m := new(grpc_lb_v1.LoadBalanceRequest)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}
