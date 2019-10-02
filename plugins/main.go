package main

import (
	"fmt"

	"github.com/anvh2/consul-cli/plugins/consul"
	"github.com/hashicorp/consul/api"
)

var service = "TestService"

func registerService() string {
	id := "test1"
	config := &consul.Config{
		ID:      id,
		Name:    service,
		Tags:    []string{"DEV"},
		Address: "127.0.0.1",
		Port:    8001,
	}
	consul.Register(config)

	return id
}

func deregisterService() {
	err := consul.DeRegister("test1")
	if err != nil {

	}

	broadCastFailedService()
}

type healthCheck struct {
	oldInstances []string
	newInstances []string
}

// helthCheck checks helth of service by giving name
func helthCheck(service string, oldInstances []string) {
	// cli, err := api.NewClient(api.DefaultConfig())
	// if err != nil {

	// }

	// services, _, err := cli.Health().Service(service, "", true, &api.QueryOptions{})
	// if err != nil {

	// }

	// var newInstances []string
	// for _, service := range services {

	// }
}

func broadCastFailedService() {
	cli, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		fmt.Println(err)
	}

	health, meta, err := cli.Health().Checks("UserSerivce", &api.QueryOptions{})

	if err != nil {
		panic(err)
	}

	fmt.Println(health)
	fmt.Println(meta)
}

func getInstances() ([]string, error) {
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		return nil, err
	}

	services, _, err := client.Health().Service("CounterService", "DEV", true, &api.QueryOptions{})
	if err != nil {
		return nil, err
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

	return instances, nil
}

func main() {
	// registerService()
	// registerService()

	// go helthCheck(service, []string{})

	// cli, _ := api.NewClient(api.DefaultConfig())

	// cli.Agent().ServiceDeregister("user-f04411d5-13b4-4c68-9ad2-e6abfd459371")

	for {
		// time.Sleep(10 * time.Millisecond)
		instances, err := getInstances()
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(instances)
	}
}
