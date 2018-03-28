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
	HttpServer *HttpServer
	ServiceDiscovery *ServiceDiscovery
}

type ServiceDiscovery struct {
	Name string
	Interval time.Duration
	Ttl	time.Duration
}

type Etcd struct {
	Addrs []string
	Root string
	Name string
	TimeOut time.Duration
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
}