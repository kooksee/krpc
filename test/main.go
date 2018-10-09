package main

import (
	"fmt"
	"net/http"
	"github.com/tendermint/go-amino"
	"github.com/kooksee/krpc/server"
	"github.com/kooksee/golog"
)

var routes = map[string]*rpcserver.RPCFunc{
	"hello_world": rpcserver.NewRPCFunc(HelloWorld, "name,num"),
}

func HelloWorld(name string, num int) (Result, error) {
	return Result{fmt.Sprintf("hi %s %d", name, num)}, nil
}

type Result struct {
	Result string
}

func main() {
	cfg := golog.DefaultConfig()
	cfg.Service = "test krpc"
	cfg.InitLog()

	mux := http.NewServeMux()
	cdc := amino.NewCodec()
	rpcserver.RegisterRPCFuncs(mux, routes, cdc)
	rpcserver.StartHTTPServer("tcp://0.0.0.0:8008", mux, rpcserver.Config{})
	select {}
}
