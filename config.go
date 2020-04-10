package kvstore

import (
	"encoding/json"
	consulapi "github.com/hashicorp/consul/api"
	"github.com/pkg/errors"
	"github.com/struckoff/SFCFramework/curve"
	kvrouter_conf "github.com/struckoff/kvrouter/config"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	//force name of node instead of hostname or consul
	Name       *string `envconfig:"NAME"`
	Address    string  `envconfig:"ADDRERSS"`
	RPCAddress string  `envconfig:"RPC_ADDRESS"`
	//standalone, consul, kvrouter
	Mode     DiscoverMode `envconfig:"MODE"`
	Power    float64      `envconfig:"POWER"`
	Capacity float64      `envconfig:"CAPACITY"`
	DBpath   string       `envconfig:"DBPATH"`
	// TTL check config
	Health   HealthConfig
	KVRouter *KVRouterConfig
	// If config implies use of consul, this options will be taken from consul KV.
	// Otherwise it will be taken from config file.
	Balancer *kvrouter_conf.BalancerConfig
	Consul   *ConfigConsul
}

func (conf *Config) Prepare() error {
	switch conf.Mode {
	case StandaloneMode, KvrouterMode:
		if conf.Name == nil {
			name, err := os.Hostname()
			if err != nil {
				return err
			}
			conf.Name = &name
		}
	case ConsulMode:
		consul, err := consulapi.NewClient(&conf.Consul.Config)
		if err != nil {
			return err
		}
		if err := conf.fillConfigFromConsul(consul); err != nil {
			return err
		}
	default:
		return errors.New("wrong mode")
	}
	return nil
}

// fillConfigFromConsul take config options from consul  KV store
// List of options
//  - Balancer.Size
//  - Balancer.Dimensions
//  - Balancer.Curve
func (conf *Config) fillConfigFromConsul(consul *consulapi.Client) error {
	if conf.Name != nil {
		name, err := consul.Agent().NodeName()
		if err != nil {
			return err
		}
		conf.Name = &name
	}

	kvMap := make(map[string][]byte)
	kv := consul.KV()
	pairs, _, err := kv.List(conf.Consul.KVFolder, nil)
	if err != nil {
		return err
	}
	for _, pair := range pairs {
		pair.Key = strings.TrimLeft(pair.Key, conf.Consul.KVFolder)
		kvMap[strings.ToLower(pair.Key)] = pair.Value
	}
	var balConfig kvrouter_conf.BalancerConfig

	if val, ok := kvMap["size"]; ok {
		balConfig.Size, err = strconv.ParseUint(string(val), 10, 64)
		if err != nil {
			return err
		}
	}

	if val, ok := kvMap["dimensions"]; ok {
		balConfig.Dimensions, err = strconv.ParseUint(string(val), 10, 64)
		if err != nil {
			return err
		}
	}
	if val, ok := kvMap["curve"]; ok {
		if err := balConfig.Curve.UnmarshalJSON(val); err != nil {
			return err
		}
	}
	conf.Balancer = &balConfig
	return nil
}

type KVRouterConfig struct {
	Address string `envconfig:"KVSTORE_KVROUTER_ADDRESS"`
}

// TTL check config
type HealthConfig struct {
	// Use CheckInterval + CheckTimeout as interval setting for deadman switch.
	CheckInterval string //Default: 30s
	// TTL will be sent each time per CheckInterval
	CheckTimeout                   string //Default: 10s
	DeregisterCriticalServiceAfter string //Default: 10m
}

func (ct *HealthConfig) UnmarshalJSON(cb []byte) error {
	m := make(map[string]string)
	if err := json.Unmarshal(cb, &m); err != nil {
		return err
	}

	ct.CheckInterval = "30s"
	if val, ok := m["CheckInterval"]; ok {
		ct.CheckInterval = val
	}

	ct.CheckTimeout = "10s"
	if val, ok := m["CheckTimeout"]; ok {
		ct.CheckTimeout = val
	}

	ct.DeregisterCriticalServiceAfter = "10m"
	if val, ok := m["DeregisterCriticalServiceAfter"]; ok {
		ct.DeregisterCriticalServiceAfter = val
	}

	return nil
}

type CurveType struct {
	curve.CurveType
}

func (ct *CurveType) UnmarshalJSON(cb []byte) error {
	c := strings.ToLower(string(cb))
	c = strings.Trim(c, "\"")
	switch c {
	case "morton":
		ct.CurveType = curve.Morton
		return nil
	case "hilbert":
		ct.CurveType = curve.Hilbert
		return nil
	default:
		return errors.New("unknown curve type")
	}
}

type ConfigConsul struct {
	consulapi.Config
	Service string
	// Key prefix for config options stored in consul KV, c
	KVFolder string //Default: ConfigConsul.Service
}

func (ct *ConfigConsul) UnmarshalJSON(cb []byte) error {
	m := make(map[string]string)
	if err := json.Unmarshal(cb, &m); err != nil {
		return err
	}

	ct.Config = *consulapi.DefaultConfig()

	if val, ok := m["Address"]; ok {
		ct.Address = val
	}
	if val, ok := m["Scheme"]; ok {
		ct.Scheme = val
	}
	if val, ok := m["Datacenter"]; ok {
		ct.Datacenter = val
	}
	if val, ok := m["WaitTime"]; ok {
		d, err := time.ParseDuration(val)
		if err != nil {
			return err
		}
		ct.WaitTime = d
	}
	if val, ok := m["Token"]; ok {
		ct.Token = val
	}
	if val, ok := m["Namespace"]; ok {
		ct.Namespace = val
	}

	if val, ok := m["Service"]; ok {
		ct.Service = val
	}

	ct.KVFolder = ct.Service + "/"
	if val, ok := m["KVFolder"]; ok {
		val = strings.TrimRight(val, "/")
		ct.KVFolder = val
		if len(ct.KVFolder) > 0 {
			ct.KVFolder += "/"
		}

	}

	//if val, ok := m["Tags"]; ok {
	//	ct.Tags = val
	//}

	return nil
}

type DiscoverMode int

const (
	StandaloneMode DiscoverMode = iota
	KvrouterMode
	ConsulMode
)

func (dn *DiscoverMode) UnmarshalJSON(cb []byte) error {
	c := strings.ToLower(string(cb))
	c = strings.Trim(c, "\"")
	switch c {
	case "standalone":
		*dn = StandaloneMode
	case "kvrouter":
		*dn = KvrouterMode
	case "consul":
		*dn = ConsulMode
	default:
		return errors.New("wrong node mode")
	}
	return nil
}
