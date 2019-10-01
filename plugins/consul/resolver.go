package consul

import (
	"fmt"

	"go.uber.org/zap"
	"google.golang.org/grpc/naming"

	"github.com/hashicorp/consul/api"
)

// Resolver ...
type Resolver struct {
	client *api.Client

	service string

	tag string

	passingOnly bool

	done chan interface{}

	update chan []*naming.Update

	logger zap.Logger
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
		done:    make(chan interface{}),
		update:  make(chan []*naming.Update, 1),
	}

	instances, lastIndex, err := resolver.getInstances(0)
	if err != nil {
		// log
	}

	updates := resolver.updateInstances(instances, nil)
	if len(updates) > 0 {
		resolver.update <- updates
	}

	go resolver.worker(instances, lastIndex)

	return resolver, nil
}

// Resolve ...
func (r *Resolver) Resolve(target string) (naming.Watcher, error) {
	return r, nil
}

// Next ...
func (r *Resolver) Next() ([]*naming.Update, error) {
	return <-r.update, nil
}

// Close ...
func (r *Resolver) Close() {
	close(r.done)
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

// worker is background process
func (r *Resolver) worker(instances []string, lastIndex uint64) {
	var err error
	var newInstances []string

	for {
		select {
		case <-r.done:
			return
		default:
			newInstances, lastIndex, err = r.getInstances(lastIndex)
			if err != nil {
				// log
			}

			updatedInstances := r.updateInstances(instances, newInstances)
			if len(updatedInstances) > 0 {
				r.update <- updatedInstances
			}
			instances = newInstances
		}
	}
}

// updateInstance mix 2 array and truncat duplicate elements
func (r *Resolver) updateInstances(instances, newInstances []string) []*naming.Update {
	instances = append(instances, newInstances...)
	mapInstances := make(map[string]bool)

	for _, instance := range instances {
		mapInstances[instance] = true
	}

	var updates []*naming.Update
	for addr := range mapInstances {
		updates = append(updates, &naming.Update{
			Op:   naming.Add,
			Addr: addr,
		})
	}
	return updates
}
