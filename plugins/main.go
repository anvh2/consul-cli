package main

import "github.com/hashicorp/consul/api"

func main() {
	cli, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		return
	}

	cli.KV().Put(&api.KVPair{Key: "anvh2", Value: []byte("An")}, nil)
}

func runWatch() {
	
}
