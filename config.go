package kvstore

import (
	"encoding/json"
	consulapi "github.com/hashicorp/consul/api"
	"github.com/pkg/errors"
	"github.com/struckoff/SFCFramework/curve"
	"strings"
	"time"
)

type Config struct {
	Address    string
	RPCAddress string
	Power      float64
	Capacity   float64
	DBpath     string
	Balancer   *BalancerConfig
	Consul     *ConfigConsul
}

// If config implies use of consul, this options will be taken from consul KV.
// Otherwise it will be taken from config file.
type BalancerConfig struct {
	Dimensions uint64    //Amount of space filling curve dimensions
	Size       uint64    //Size of space filling curve
	Curve      CurveType //Space filling curve type
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
	// Use CheckInterval + CheckTimeout as interval setting for deadman switch.
	// TTL will be sent each time per CheckInterval
	CheckInterval                  string //Default: 30s
	CheckTimeout                   string //Default: 10s
	DeregisterCriticalServiceAfter string //Default: 10m
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

	ct.KVFolder = ct.Service + "/"
	if val, ok := m["KVFolder"]; ok {
		val = strings.TrimRight(val, "/")
		ct.KVFolder = val
		if len(ct.KVFolder) > 0 {
			ct.KVFolder += "/"
		}
		ct.DeregisterCriticalServiceAfter = val
	}

	//if val, ok := m["Tags"]; ok {
	//	ct.Tags = val
	//}

	return nil
}
