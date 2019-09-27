package consul

import (
	"log"

	"github.com/hashicorp/consul/api"
)

// Config ...
type Config struct {
	ID      string
	Name    string
	Tags    []string
	Address string
	Port    int
}

// Register -
func Register(config *Config) error {
	cli, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		return err
	}

	reg := &api.AgentServiceRegistration{
		ID:      config.ID,
		Name:    config.Name,
		Tags:    config.Tags,
		Address: config.Address,
		Port:    config.Port,
	}
	err = cli.Agent().ServiceRegister(reg)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

// DeRegister ...
func DeRegister(id string) error {
	cli, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		return err
	}
	cli.Agent().ServiceDeregister(id)
	return nil
}
