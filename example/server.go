package example

import (
	"github.com/kooksee/krpc/server"
	"github.com/tendermint/go-amino"
	"net/http"
)

// Define some routes
var Routes = map[string]*krpcs.RPCFunc{
	"echo":            krpcs.NewRPCFunc(EchoResult, "arg"),
	"echo_bytes":      krpcs.NewRPCFunc(EchoBytesResult, "arg"),
	"echo_data_bytes": krpcs.NewRPCFunc(EchoDataBytesResult, "arg"),
	"echo_int":        krpcs.NewRPCFunc(EchoIntResult, "arg"),
}

const (
	tcpAddr = "tcp://0.0.0.0:8088"
)

var RoutesCdc = amino.NewCodec()

func setup() {
	mux := http.NewServeMux()
	krpcs.RegisterRPCFuncs(mux, Routes, RoutesCdc)
	krpcs.StartHTTPServer(tcpAddr, mux, krpcs.Config{})
	select {}
}
