package consul

import (
	"fmt"

	"google.golang.org/grpc/naming"

	"github.com/hashicorp/consul/api"
)

// Resolver ...
type Resolver struct {
	client *api.Client

	service string

	tag string

	// passingOnly bool

	// logger zap.Logger
}

// NewResolver ...
func NewResolver(service, tag string) (*Resolver, error) {
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		return nil, err
	}

	resolver := &Resolver{
		client:  client,
		service: service,
		tag:     tag,
	}

	return resolver, nil
}

// Resolve ...
func (r *Resolver) Resolve(target string) (naming.Watcher, error) {
	return r, nil
}

// Next ...
func (r *Resolver) Next() ([]*naming.Update, error) {
	return nil, nil
}

// Close ...
func (r *Resolver) Close() {

}

func (r *Resolver) getInstances(lastIndex uint64) ([]string, uint64, error) {
	services, meta, err := r.client.Health().Service(r.service, r.tag, true, &api.QueryOptions{
		WaitIndex: lastIndex,
	})
	if err != nil {
		return nil, lastIndex, err
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

// LookupServiceHost resolves a service name to a list of network addresses.
func (r *Resolver) LookupServiceHost(service, tag string) ([]string, uint64, error) {
	services, meta, err := r.client.Health().Service(service, tag, true, &api.QueryOptions{})
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
func LookupServices(service, tag string) ([]string, error) {
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		return nil, err
	}

	services, _, err := client.Health().Service(service, tag, true, &api.QueryOptions{})
	if err != nil {
		return nil, err
	}

	var hosts []string
	for _, service := range services {
		addr := service.Service.Address
		if len(addr) == 0 {
			addr = service.Node.Address
		}
		address := fmt.Sprintf("%s:%d", addr, service.Service.Port)
		hosts = append(hosts, address)
	}

	return hosts, nil
}

// BroadCastFailedService ...
func BroadCastFailedService() {

}
