package config

import ()

type Config struct {
	Address    string `envconfig:"ADDRESS"`
	RPCAddress string `envconfig:"RPC_ADDRESS"`
	Balancer   *BalancerConfig
}
