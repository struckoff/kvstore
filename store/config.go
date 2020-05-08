package store

import (
	"encoding/json"
	"github.com/kelseyhightower/envconfig"
	"github.com/struckoff/kvstore/router/config"
	"github.com/struckoff/kvstore/router/rpcapi"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/pkg/errors"
	"github.com/struckoff/SFCFramework/curve"
)

const (
	envPrefix = "KVSTORE"
)

var (
	ignoreEnvs = []string{
		"GEO_XXX_UNRECOGNIZED",
		"GEO_XXX_SIZECACHE",
	}
)

/******* ENVIRONMENT OPTIONS ******
KVSTORE_NAME
KVSTORE_ADDRESS
KVSTORE_RPC_ADDRESS
KVSTORE_GEO_LONGITUDE
KVSTORE_GEO_LATITUDE
KVSTORE_MODE
KVSTORE_POWER
KVSTORE_CAPACITY
KVSTORE_DBPATH
KVSTORE_HEALTH_INTERVAL
KVSTORE_HEALTH_TIMEOUT
KVSTORE_HEALTH_DEREGISTER_CRITICAL_SERVICE_AFTER
KVSTORE_KVROUTER_ADDRESS
KVSTORE_BALANCER_MODE
KVSTORE_BALANCER_SFC_DIMENSIONS
KVSTORE_BALANCER_SFC_SIZE
KVSTORE_BALANCER_SFC_CURVE_CURVETYPE
KVSTORE_BALANCER_NODE_HASH
KVSTORE_BALANCER_DATA_MODE
KVSTORE_CONSUL_ADDRESS
KVSTORE_CONSUL_SCHEME
KVSTORE_CONSUL_DATACENTER
KVSTORE_CONSUL_TRANSPORT_PROXY
KVSTORE_CONSUL_TRANSPORT_DIALCONTEXT
KVSTORE_CONSUL_TRANSPORT_DIAL
KVSTORE_CONSUL_TRANSPORT_DIALTLSCONTEXT
KVSTORE_CONSUL_TRANSPORT_DIALTLS
KVSTORE_CONSUL_TRANSPORT_TLSCLIENTCONFIG_RAND
KVSTORE_CONSUL_TRANSPORT_TLSCLIENTCONFIG_TIME
KVSTORE_CONSUL_TRANSPORT_TLSCLIENTCONFIG_CERTIFICATES
KVSTORE_CONSUL_TRANSPORT_TLSCLIENTCONFIG_NAMETOCERTIFICATE
KVSTORE_CONSUL_TRANSPORT_TLSCLIENTCONFIG_GETCERTIFICATE
KVSTORE_CONSUL_TRANSPORT_TLSCLIENTCONFIG_GETCLIENTCERTIFICATE
KVSTORE_CONSUL_TRANSPORT_TLSCLIENTCONFIG_GETCONFIGFORCLIENT
KVSTORE_CONSUL_TRANSPORT_TLSCLIENTCONFIG_VERIFYPEERCERTIFICATE
KVSTORE_CONSUL_TRANSPORT_TLSCLIENTCONFIG_NEXTPROTOS
KVSTORE_CONSUL_TRANSPORT_TLSCLIENTCONFIG_SERVERNAME
KVSTORE_CONSUL_TRANSPORT_TLSCLIENTCONFIG_CLIENTAUTH
KVSTORE_CONSUL_TRANSPORT_TLSCLIENTCONFIG_INSECURESKIPVERIFY
KVSTORE_CONSUL_TRANSPORT_TLSCLIENTCONFIG_CIPHERSUITES
KVSTORE_CONSUL_TRANSPORT_TLSCLIENTCONFIG_PREFERSERVERCIPHERSUITES
KVSTORE_CONSUL_TRANSPORT_TLSCLIENTCONFIG_SESSIONTICKETSDISABLED
KVSTORE_CONSUL_TRANSPORT_TLSCLIENTCONFIG_SESSIONTICKETKEY
KVSTORE_CONSUL_TRANSPORT_TLSCLIENTCONFIG_CLIENTSESSIONCACHE
KVSTORE_CONSUL_TRANSPORT_TLSCLIENTCONFIG_MINVERSION
KVSTORE_CONSUL_TRANSPORT_TLSCLIENTCONFIG_MAXVERSION
KVSTORE_CONSUL_TRANSPORT_TLSCLIENTCONFIG_CURVEPREFERENCES
KVSTORE_CONSUL_TRANSPORT_TLSCLIENTCONFIG_DYNAMICRECORDSIZINGDISABLED
KVSTORE_CONSUL_TRANSPORT_TLSCLIENTCONFIG_RENEGOTIATION
KVSTORE_CONSUL_TRANSPORT_TLSCLIENTCONFIG_KEYLOGWRITER
KVSTORE_CONSUL_TRANSPORT_TLSHANDSHAKETIMEOUT
KVSTORE_CONSUL_TRANSPORT_DISABLEKEEPALIVES
KVSTORE_CONSUL_TRANSPORT_DISABLECOMPRESSION
KVSTORE_CONSUL_TRANSPORT_MAXIDLECONNS
KVSTORE_CONSUL_TRANSPORT_MAXIDLECONNSPERHOST
KVSTORE_CONSUL_TRANSPORT_MAXCONNSPERHOST
KVSTORE_CONSUL_TRANSPORT_IDLECONNTIMEOUT
KVSTORE_CONSUL_TRANSPORT_RESPONSEHEADERTIMEOUT
KVSTORE_CONSUL_TRANSPORT_EXPECTCONTINUETIMEOUT
KVSTORE_CONSUL_TRANSPORT_TLSNEXTPROTO
KVSTORE_CONSUL_TRANSPORT_PROXYCONNECTHEADER
KVSTORE_CONSUL_TRANSPORT_MAXRESPONSEHEADERBYTES
KVSTORE_CONSUL_TRANSPORT_WRITEBUFFERSIZE
KVSTORE_CONSUL_TRANSPORT_READBUFFERSIZE
KVSTORE_CONSUL_TRANSPORT_FORCEATTEMPTHTTP2
KVSTORE_CONSUL_HTTPCLIENT_TRANSPORT
KVSTORE_CONSUL_HTTPCLIENT_CHECKREDIRECT
KVSTORE_CONSUL_HTTPCLIENT_JAR
KVSTORE_CONSUL_HTTPCLIENT_TIMEOUT
KVSTORE_CONSUL_HTTPAUTH_USERNAME
KVSTORE_CONSUL_HTTPAUTH_PASSWORD
KVSTORE_CONSUL_WAITTIME
KVSTORE_CONSUL_TOKEN
KVSTORE_CONSUL_TOKENFILE
KVSTORE_CONSUL_NAMESPACE
KVSTORE_CONSUL_TLSCONFIG_ADDRESS
KVSTORE_CONSUL_TLSCONFIG_CAFILE
KVSTORE_CONSUL_TLSCONFIG_CAPATH
KVSTORE_CONSUL_TLSCONFIG_CAPEM
KVSTORE_CONSUL_TLSCONFIG_CERTFILE
KVSTORE_CONSUL_TLSCONFIG_CERTPEM
KVSTORE_CONSUL_TLSCONFIG_KEYFILE
KVSTORE_CONSUL_TLSCONFIG_KEYPEM
KVSTORE_CONSUL_TLSCONFIG_INSECURESKIPVERIFY
KVSTORE_CONSUL_SERVICE
KVSTORE_CONSUL_KVFOLDER
*********************************************/

type Config struct {
	//force name of node instead of hostname or consul
	Name       *string         `envconfig:"NAME"`
	Address    string          `envconfig:"ADDRESS"`
	RPCAddress string          `envconfig:"RPC_ADDRESS"`
	Geo        *rpcapi.GeoData `envconfig:"GEO"`
	//standalone, consul, kvrouter
	Mode     DiscoverMode `envconfig:"MODE"`
	Power    float64      `envconfig:"POWER"`
	Capacity float64      `envconfig:"CAPACITY"`
	DBpath   string       `envconfig:"DBPATH"`
	// TTL check config
	Health   HealthConfig
	KVRouter *KVRouterConfig `envconfig:"KVROUTER"`
	// If config implies use of consul, this options will be taken from consul KV.
	// Otherwise it will be taken from config file.
	Balancer *config.BalancerConfig `envconfig:"BALANCER"`
	Consul   *ConfigConsul
}

func ReadConfig(cfgPath string) (Config, error) {
	configFile, err := os.Open(cfgPath)
	if err != nil {
		return Config{}, errors.Wrap(err, "failed to open config file")
	}
	defer configFile.Close()
	var conf Config
	if err := json.NewDecoder(configFile).Decode(&conf); err != nil {
		return Config{}, errors.Wrap(err, "failed to parse config file")
	}
	for _, key := range ignoreEnvs {
		if len(envPrefix) > 0 {
			key = envPrefix + "_" + key
		}
		if err := os.Unsetenv(key); err != nil {
			return Config{}, err
		}
	}
	if err := envconfig.Process(envPrefix, &conf); err != nil {
		return Config{}, err
	}

	if err := conf.Prepare(); err != nil {
		return Config{}, err
	}
	return conf, nil
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
//  nodehasher 			   -> Balancer.NodeHash
//  balancermode 		   -> Balancer.Mode
//  datamode 			   -> Balancer.DataMode
//  sfc.size 	   		   -> Balancer.SFC.Size
//  sfc.dimensions 		   -> Balancer.SFC.Dimensions
//  sfc.curve 	   		   -> Balancer.SFC.Curve
func (conf *Config) fillConfigFromConsul(consul *consulapi.Client) error {
	logform := "%s: found in consul with value \"%s\""
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
		pair.Key = strings.TrimPrefix(pair.Key, conf.Consul.KVFolder)
		kvMap[strings.ToLower(pair.Key)] = pair.Value
	}
	//var balConfig config.BalancerConfig

	if val, ok := kvMap["balancermode"]; ok {
		log.Printf(logform, "balancermode", val)
		if err := conf.Balancer.Mode.UnmarshalJSON(val); err != nil {
			return err
		}
	}

	if val, ok := kvMap["nodehasher"]; ok {
		log.Printf(logform, "nodehasher", val)
		if err := conf.Balancer.NodeHash.UnmarshalJSON(val); err != nil {
			return err
		}
	}
	if val, ok := kvMap["datamode"]; ok {
		log.Printf(logform, "datamode", val)
		if err := conf.Balancer.DataMode.UnmarshalJSON(val); err != nil {
			return err
		}
	}
	if val, ok := kvMap["sfc.size"]; ok {
		log.Printf(logform, "sfc.size", val)
		if conf.Balancer.SFC == nil {
			conf.Balancer.SFC = &config.SFCConfig{}
		}
		conf.Balancer.SFC.Size, err = strconv.ParseUint(string(val), 10, 64)
		if err != nil {
			return err
		}
	}
	if val, ok := kvMap["sfc.dimensions"]; ok {
		log.Printf(logform, "sfc.dimensions", val)
		if conf.Balancer.SFC == nil {
			conf.Balancer.SFC = &config.SFCConfig{}
		}
		conf.Balancer.SFC.Dimensions, err = strconv.ParseUint(string(val), 10, 64)
		if err != nil {
			return err
		}
	}
	if val, ok := kvMap["sfc.curve"]; ok {
		log.Printf(logform, "sfc.curve", val)
		if conf.Balancer.SFC == nil {
			conf.Balancer.SFC = &config.SFCConfig{}
		}
		if err := conf.Balancer.SFC.Curve.UnmarshalJSON(val); err != nil {
			return err
		}
	}
	//if val, ok := kvMap["ring.load"]; ok {
	//	log.Printf(logform, "ring.load", val)
	//	if balConfig.Ring == nil {
	//		balConfig.Ring = &consistent.Config{}
	//	}
	//	balConfig.Ring.Load, err = strconv.ParseFloat(string(val), 64)
	//	if err != nil {
	//		return err
	//	}
	//}
	//if val, ok := kvMap["ring.partitioncount"]; ok {
	//	log.Printf(logform, "ring.partitionCount", val)
	//	if balConfig.Ring == nil {
	//		balConfig.Ring = &consistent.Config{}
	//	}
	//	v, err := strconv.ParseInt(string(val), 10, 64)
	//	if err != nil {
	//		return err
	//	}
	//	balConfig.Ring.PartitionCount = int(v)
	//}
	//
	//if val, ok := kvMap["ring.replicationfactor"]; ok {
	//	log.Printf(logform, "ring.replicationFactor", val)
	//	if balConfig.Ring == nil {
	//		balConfig.Ring = &consistent.Config{}
	//	}
	//	v, err := strconv.ParseInt(string(val), 10, 64)
	//	if err != nil {
	//		return err
	//	}
	//	balConfig.Ring.ReplicationFactor = int(v)
	//}

	//conf.Balancer = &balConfig
	return nil
}

type KVRouterConfig struct {
	Address string `envconfig:"ADDRESS"`
}

// TTL check config
type HealthConfig struct {
	// Use CheckInterval + CheckTimeout as interval setting for deadman switch.
	//Default: 30s
	CheckInterval string `envconfig:"INTERVAL"`
	// TTL will be sent each time per CheckInterval
	//Default: 10s
	CheckTimeout string `envconfig:"TIMEOUT"`
	//Default: 10m
	DeregisterCriticalServiceAfter string `envconfig:"DEREGISTER_CRITICAL_SERVICE_AFTER"`
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
	StandaloneMode DiscoverMode = iota + 1
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
