package consul

import (
	"context"
	"sync"

	"go.uber.org/zap"
	"google.golang.org/grpc/health/grpc_health_v1"
)

// Health -
type Health struct {
	mu   sync.Mutex
	endp []*Endpoint

	logger zap.Logger
}

// Endpoint -
type Endpoint struct {
}

// ServiceStatus -
type ServiceStatus int64

var (
	// UNKNOWN -
	UNKNOWN ServiceStatus
	// SERVING -
	SERVING ServiceStatus = 1
	// NOTSERVING -
	NOTSERVING ServiceStatus = 2
	// SERVICEUNKNOWN -
	SERVICEUNKNOWN ServiceStatus = 3
)

// HealthImpl -
type HealthImpl struct{}

// Check -
func (h *HealthImpl) Check(ctx context.Context, req *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	return &grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	}, nil
}

// Watch -
func (h *HealthImpl) Watch(req *grpc_health_v1.HealthCheckRequest, w grpc_health_v1.Health_WatchServer) error {
	return nil
}
