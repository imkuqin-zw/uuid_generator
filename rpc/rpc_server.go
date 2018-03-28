package rpc

import (
	"github.com/imkuqin-zw/id_generator/common/etcd"
	"github.com/imkuqin-zw/id_generator/config"
	"net"
	sd "github.com/imkuqin-zw/id_generator/common/service_discovery/etcd"
	"strings"
	"google.golang.org/grpc"
	"golang.org/x/net/context"
	"github.com/imkuqin-zw/id_generator/rpc/protobuf"
	"github.com/imkuqin-zw/id_generator/common/snowflake"
)

const (
	UUID_QUEUE = 1024 // uuid process queue
)

type RPCServer struct {
	etcd       *etcd.Etcd
	machineID uint64 // 10-bit machine id
	chProc    chan chan uint64
}

func (s *RPCServer) init() (err error) {
	s.chProc = make(chan chan uint64, UUID_QUEUE)
	if s.etcd, err = etcd.NewEtcd(config.Conf.Etcd); err != nil {
		return
	}
	if s.machineID, err = s.etcd.GetMachineID(); err != nil {
		return
	}

	go snowflake.CreateUUID(s.chProc, s.machineID)
	return
}

func (s *RPCServer) Next(cxt context.Context, req *protobuf.SnowflakeKey) (*protobuf.SnowflakeVal, error) {
	val, err := s.etcd.GetNextByName(req.Name)
	if err != nil {
		return nil, err
	}

	return &protobuf.SnowflakeVal{Value:val}, nil
}

func (s *RPCServer) GetUUID(context.Context, *protobuf.SnowflakeNullReq) (*protobuf.SnowflakeUUID, error) {
	req := make(chan uint64, 1)
	s.chProc <- req
	return &protobuf.SnowflakeUUID{<-req}, nil
}

func Run() {
	lis, err := net.Listen(config.Conf.RpcServer.Proto, config.Conf.RpcServer.Addr)
	if err != nil {
		panic(err)
	}
	target := strings.Join(config.Conf.Etcd.Addrs, ",")
	err = sd.Register(config.Conf.ServiceDiscovery.Name, config.Conf.RpcServer.Addr,
		target, config.Conf.ServiceDiscovery.Interval, config.Conf.ServiceDiscovery.Ttl)
	if err != nil {
		panic(err)
	}

	rpcServer := &RPCServer{}
	if err = rpcServer.init(); err != nil {
		panic(err)
	}
	s := grpc.NewServer()
	protobuf.RegisterGeneratorServer(s, rpcServer)
	s.Serve(lis)
}