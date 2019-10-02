package consul

import (
	"fmt"
	"time"

	"github.com/hashicorp/consul/api"
)

// Config ...
type Config struct {
	ID              string
	Name            string
	Tags            []string
	Address         string
	Port            int
	Interval        time.Duration
	DeRegisterAfter time.Duration
}

// DefaultConfig returns Config with default Interval and DeregisterAfter
func DefaultConfig() *Config {
	return &Config{
		Interval:        time.Duration(10) * time.Second,
		DeRegisterAfter: time.Duration(1) * time.Minute,
	}
}

// Register -
func Register(config *Config) error {
	cli, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		return err
	}

	// check service health per 10s
	interval := time.Duration(10) * time.Second
	// check if health is cretical after 1 minute, this service will deregister frorm consul
	deregister := time.Duration(1) * time.Minute

	reg := &api.AgentServiceRegistration{
		ID:      config.ID,
		Name:    config.Name,
		Tags:    config.Tags,
		Address: config.Address,
		Port:    config.Port,
		Check: &api.AgentServiceCheck{
			GRPC:                           fmt.Sprintf("%v:%v/%v", config.Address, config.Port, config.Name),
			Interval:                       interval.String(),
			DeregisterCriticalServiceAfter: deregister.String(),
		},
	}

	err = cli.Agent().ServiceRegister(reg)
	if err != nil {
		return err
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
