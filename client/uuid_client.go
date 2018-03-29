package main

import (
	"github.com/imkuqin-zw/uuid_generator/common/service_discovery/etcd"
	"google.golang.org/grpc"
	"github.com/golang/glog"
	"time"
	"github.com/imkuqin-zw/uuid_generator/rpc/protobuf"
	"golang.org/x/net/context"
	"fmt"
)

func main() {
	r := etcd.NewResolver("uuid_server")
	b := grpc.RoundRobin(r)
	conn, err := grpc.Dial("127.0.0.1:2379", grpc.WithInsecure(), grpc.WithBalancer(b))
	if err != nil {
		glog.Error(err)
		panic(err)
	}
	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case <- ticker.C:
			client := protobuf.NewGeneratorClient(conn)
			//resp, err := client.GetUUID(context.Background(), &protobuf.SnowflakeNullReq{})
			//if err != nil {
			//	fmt.Println(err.Error())
			//} else {
			//	fmt.Println(resp.Uuid)
			//}
			resp, err := client.Next(context.Background(), &protobuf.SnowflakeKey{Name:"user_id"})
			if err != nil {
				fmt.Println(err.Error())
			} else {
				fmt.Println(resp.Value)
			}
			break
		}
	}
}