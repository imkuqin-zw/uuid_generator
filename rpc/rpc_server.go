package rpc

import (
	"github.com/imkuqin-zw/uuid_generator/common/etcd"
	"github.com/imkuqin-zw/uuid_generator/config"
	"net"
	sd "github.com/imkuqin-zw/uuid_generator/common/service_discovery/etcd"
	"strings"
	"google.golang.org/grpc"
	"golang.org/x/net/context"
	"github.com/imkuqin-zw/uuid_generator/rpc/protobuf"
	"github.com/imkuqin-zw/uuid_generator/common/snowflake"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

const (
	UUID_QUEUE = 1024
)

type RpcServer struct {
	etcd       *etcd.Etcd
	machineID uint64 // 10-bit machine id
	chProc    chan chan uint64
}

func (s *RpcServer) init() (err error) {
	s.chProc = make(chan chan uint64, UUID_QUEUE)
	if s.etcd, err = etcd.NewEtcd(config.Conf.Etcd, config.Conf.Locks); err != nil {
		return
	}
	if s.machineID, err = s.etcd.GetMachineID(); err != nil {
		return
	}

	go snowflake.CreateUUID(s.chProc, s.machineID)
	return
}

func (s *RpcServer) Next(cxt context.Context, req *protobuf.SnowflakeKey) (*protobuf.SnowflakeVal, error) {
	val, err := s.etcd.GetNextByName(req.Name)
	if err != nil {
		return nil, err
	}

	return &protobuf.SnowflakeVal{Value:val}, nil
}

func (s *RpcServer) GetUUID(context.Context, *protobuf.SnowflakeNullReq) (*protobuf.SnowflakeUUID, error) {
	req := make(chan uint64, 1)
	s.chProc <- req
	return &protobuf.SnowflakeUUID{Uuid: <-req}, nil
}

func Run() {
	fmt.Printf("server started at %s\r\n", config.Conf.RpcServer.Addr)
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

	rpcServer := &RpcServer{}
	if err = rpcServer.init(); err != nil {
		panic(err)
	}
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGQUIT)
	go func() {
		s := <-ch
		fmt.Printf("receive signal %v\r\n", s)
		sd.UnRegister()
		os.Exit(0)
	}()
	s := grpc.NewServer()
	protobuf.RegisterGeneratorServer(s, rpcServer)
	if err = s.Serve(lis); err != nil {
		panic(err)
	}
}