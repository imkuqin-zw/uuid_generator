package config

import (
	"time"
	"flag"
	"io/ioutil"
	"gopkg.in/yaml.v2"
)

var (
	confPath string
	Conf *Config
)

type Config struct {
	Etcd *Etcd
	RpcServer *RpcServer
	ServiceDiscovery *ServiceDiscovery
	Locks map[string]*EtcdSync
}

type ServiceDiscovery struct {
	Name string
	Interval time.Duration
	Ttl	time.Duration
}

type Etcd struct {
	Name string
	Root string
	TimeOut time.Duration
	Addrs []string
}

type EtcdSync struct {
	Addrs []string
	Ttl time.Duration
	Root string
	Name string
}


type RpcServer struct {
	Proto string
	Addr string
}

type HttpServer struct {
	Addr string
}

func Init() (err error) {
	var buf []byte
	if buf, err = ioutil.ReadFile(confPath); err == nil {
		err = yaml.Unmarshal(buf, Conf)
	}
	return
}

func init()  {
	flag.StringVar(&confPath, "conf", "./default.yaml", "config path")
	Conf = &Config{}
}