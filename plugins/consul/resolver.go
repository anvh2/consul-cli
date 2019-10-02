package consul

import (
	"fmt"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc/naming"

	"github.com/hashicorp/consul/api"
)

// Resolver ...
type Resolver struct {
	client      *api.Client
	service     string
	tag         string
	passingOnly bool
	done        chan interface{}
	updatec     chan []*naming.Update
	logger      zap.Logger
}

// NewResolver ...
func NewResolver(service, tag string) (*Resolver, error) {
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		return nil, err
	}

	resolver := &Resolver{
		client:      client,
		service:     service,
		tag:         tag,
		passingOnly: true,
		done:        make(chan interface{}),
		updatec:     make(chan []*naming.Update, 1),
	}

	instances, lastIndex, err := resolver.getInstances(0)
	if err != nil {
		// log
	}

	updates := resolver.updateInstances(nil, instances)
	if len(updates) > 0 {
		resolver.updatec <- updates
	}

	go resolver.worker(instances, lastIndex)

	return resolver, nil
}

// WithLogger -
func (r *Resolver) WithLogger(logger zap.Logger) {
	r.logger = logger
}

// Resolve ...
func (r *Resolver) Resolve(target string) (naming.Watcher, error) {
	return r, nil
}

// Next ...
func (r *Resolver) Next() ([]*naming.Update, error) {
	return <-r.updatec, nil
}

// Close ...
func (r *Resolver) Close() {
	select {
	case <-r.done:
	default:
		close(r.done)
		close(r.updatec)
	}
}

func (r *Resolver) getInstances(lastIndex uint64) ([]string, uint64, error) {
	services, meta, err := r.client.Health().Service(r.service, r.tag, r.passingOnly, &api.QueryOptions{
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

// worker is background process, it query to consul and detect config change
func (r *Resolver) worker(instances []string, lastIndex uint64) {
	var err error
	var newInstances []string

	for {
		time.Sleep(5 * time.Second)
		select {
		case <-r.done:
			return
		default:
			newInstances, lastIndex, err = r.getInstances(lastIndex)
			if err != nil {
				// log
				continue
			}

			updatedInstances := r.updateInstances(instances, newInstances)
			if len(updatedInstances) > 0 {
				r.updatec <- updatedInstances
			}
			instances = newInstances
		}
	}
}

// updateInstance mix 2 array and truncat duplicate elements
func (r *Resolver) updateInstances(oldInstances, newInstances []string) []*naming.Update {
	oldAddr := make(map[string]bool, len(oldInstances))
	for _, instance := range oldInstances {
		oldAddr[instance] = true
	}

	newAddr := make(map[string]bool, len(newInstances))
	for _, instance := range newInstances {
		newAddr[instance] = true
	}

	var updates []*naming.Update
	for addr := range newAddr {
		if _, ok := oldAddr[addr]; !ok {
			updates = append(updates, &naming.Update{
				Op:   naming.Add,
				Addr: addr,
			})
		}
	}

	for addr := range oldAddr {
		if _, ok := newAddr[addr]; !ok {
			updates = append(updates, &naming.Update{
				Op:   naming.Delete,
				Addr: addr,
			})
		}
	}

	return updates
}
