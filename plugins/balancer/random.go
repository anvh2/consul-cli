package balancer

import (
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/resolver"
)

var (
	random = "random"
)

type randomBuilder struct{}

// Build -
func (*randomBuilder) Build(map[resolver.Address]balancer.SubConn) balancer.Picker {
	return nil
}
