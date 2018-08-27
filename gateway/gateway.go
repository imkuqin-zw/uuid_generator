package gateway

import (
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	"net/http"
	"github.com/imkuqin-zw/uuid_generator/rpc/protobuf"
	"github.com/imkuqin-zw/uuid_generator/config"
	"golang.org/x/net/context"
	"fmt"
	"strings"
	"github.com/imkuqin-zw/uuid_generator/common/service_discovery/etcd"
)

func Run() {
	fmt.Printf("endpoint server at %s; server started at %s\r\n ", config.Conf.RpcServer.Addr, config.Conf.HttpServer.Addr)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	mux := runtime.NewServeMux()
	r := etcd.NewResolver("uuid_server")
	b := grpc.RoundRobin(r)
	opts := []grpc.DialOption{grpc.WithInsecure(), grpc.WithBalancer(b)}
	target := strings.Join(config.Conf.Etcd.Addrs, ",")
	err := protobuf.RegisterGeneratorHandlerFromEndpoint(ctx, mux, target, opts)
	if err != nil {
		panic(err)
	}
	fmt.Println(http.ListenAndServe(config.Conf.HttpServer.Addr, mux))
}