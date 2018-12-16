package main

import (
	"fmt"
	"net/http"
	"github.com/tendermint/go-amino"
	"github.com/kooksee/golog"
	"github.com/kooksee/krpc/server"
	"github.com/kooksee/krpc/test"
)

var routes = map[string]*krpcs.RPCFunc{
	"hello_world": krpcs.NewRPCFunc(HelloWorld, "name,num"),
}

func HelloWorld(name string, num int) (krpctest.Result, error) {
	return krpctest.Result{fmt.Sprintf("hi %s %d", name, num)}, nil
}

func main() {
	cfg := golog.DefaultConfig()
	cfg.Service = "test krpc"
	cfg.InitLog()

	krpcs.Init()

	mux := http.NewServeMux()
	cdc := amino.NewCodec()
	krpcs.RegisterRPCFuncs(mux, routes, cdc)
	krpcs.StartHTTPServer("tcp://0.0.0.0:8008", mux, krpcs.Config{})
}
