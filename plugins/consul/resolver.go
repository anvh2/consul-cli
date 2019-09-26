package consul

import (
	"fmt"

	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"
)

// Resolver ...
type Resolver struct {
	client *api.Client

	service string

	tag string

	passingOnly bool

	logger zap.Logger
}

// NewResolver ...
func NewResolver(service, tag string) (*Resolver, error) {
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		return nil, err
	}
	return &Resolver{
		client:      client,
		service:     service,
		tag:         tag,
		passingOnly: true,
	}, nil
}

// LookupServiceHost resolves a service name to a list of network addresses.
func (r *Resolver) LookupServiceHost(name string) ([]string, uint64, error) {
	services, meta, err := r.client.Health().Service(r.service, r.tag, r.passingOnly, &api.QueryOptions{})
	if err != nil {
		return nil, 0, err
	}

	var instances []string
	for _, service := range services {
		addr := service.Service.Address
		if len(addr) == 0 {
			addr = service.Node.Address
		}
		address := fmt.Sprintf("%s:%d", addr, service.Service.Port)
		instances = append(instances, address)
	}

	return instances, meta.LastIndex, nil
}

// LookupServices gets services
func (r *Resolver) LookupServices(name string) {

}
