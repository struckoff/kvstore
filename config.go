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
	//Name        string
	Address     string
	RPCAddress  string
	Power       float64
	Capacity    float64
	DBpath      string
	Entrypoints []string  //If not empty node tries to connect to each entrypoint, send its meta and receive cluster info
	Dimensions  uint64    //Amount of space filling curve dimensions
	Size        uint64    //Size of space filling curve
	Curve       CurveType //Space filling curve type
	Consul      *ConfigConsul
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
	Service                        string
	CheckInterval                  string
	CheckTimeout                   string
	DeregisterCriticalServiceAfter string
}

func (ct *ConfigConsul) UnmarshalJSON(cb []byte) error {
	//defconf := consulapi.DefaultConfig()
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
	ct.DeregisterCriticalServiceAfter = "600s"
	if val, ok := m["DeregisterCriticalServiceAfter"]; ok {
		ct.DeregisterCriticalServiceAfter = val
	}

	//if val, ok := m["Tags"]; ok {
	//	ct.Service = val
	//}

	return nil
}
