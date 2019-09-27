package consul

import (
	"errors"
	"fmt"

	"github.com/hashicorp/consul/api"
)

// Client ...
type client struct {
	consul *api.Client
}

// NewConsulClient ...
func NewConsulClient(addr string) (Client, error) {
	config := api.DefaultConfig()
	config.Address = addr

	cli, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}

	return &client{
		consul: cli,
	}, nil
}

// GetService ...
func (c *client) GetService(service, tag string) ([]string, error) {
	passingOnly := true
	entries, _, err := c.consul.Health().Service(service, tag, passingOnly, nil)
	if err != nil {
		return nil, err
	}

	if len(entries) == 0 {
		return nil, errors.New("Service not found")
	}

	var instances []string
	for _, service := range entries {
		addr := service.Service.Address
		if len(addr) == 0 {
			addr = service.Node.Address
		}

		address := fmt.Sprintf("%s:%d", addr, service.Service.Port)
		instances = append(instances, address)
	}
	return instances, nil
}

// RegisterService ...
func (c *client) RegisterService(name, address string, tags []string, port int) error {
	register := &api.AgentServiceRegistration{
		ID:      name,
		Name:    name,
		Tags:    tags,
		Address: address,
		Port:    port,
	}

	return c.consul.Agent().ServiceRegister(register)
}

// DeRegisterService ...
func (c *client) DeRegisterService(id string) error {
	return c.consul.Agent().ServiceDeregister(id)
}
