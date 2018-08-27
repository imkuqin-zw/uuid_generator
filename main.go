package main

import (
	"flag"
	"github.com/imkuqin-zw/uuid_generator/config"
	"github.com/imkuqin-zw/uuid_generator/rpc"
	"github.com/imkuqin-zw/uuid_generator/gateway"
)

var mode string

func main() {
	if err := config.Init(); err != nil {
		panic(err)
	}

	switch mode {
	case "rpc":
		rpc.Run()
	case "http":
		gateway.Run()
	default:
		panic("mode error")
	}
}

func init() {
	flag.StringVar(&mode, "mode", "rpc", "run mode")
	flag.Parse()
}