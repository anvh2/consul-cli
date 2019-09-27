package consul

import (
	"fmt"

	"github.com/hashicorp/consul/api"
)

// NewClient ...
func NewClient() {
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		panic(err)
	}

	kv := client.KV()

	p := &api.KVPair{Key: "REDIS_MAXCLIENTS", Value: []byte("1000")}
	_, err = kv.Put(p, nil)
	if err != nil {
		panic(err)
	}

	pair, _, err := kv.Get("REDIS_MAXCLIENTS", nil)
	if err != nil {
		panic(err)
	}

	fmt.Printf("MAX_CLIENTS: %s\n", pair.Value)
}

// RegisterWithConsul ...
func RegisterWithConsul(id, name, address string, port int) {
	cli, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		fmt.Println(err)
	}

	reg := &api.AgentServiceRegistration{
		ID:      id,
		Name:    name,
		Tags:    []string{"dev"},
		Address: address,
		Port:    port,
	}

	err = cli.Agent().ServiceRegister(reg)
	if err != nil {
		fmt.Println(err)
	}
	defer cli.Agent().ServiceDeregister(reg.ID)
}
