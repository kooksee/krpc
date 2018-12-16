package main

import (
	"github.com/kooksee/golog"
	"fmt"
	"github.com/tendermint/go-amino"
	"github.com/kooksee/krpc/client"
	"github.com/kooksee/krpc/test"
)

func main() {
	cfg := golog.DefaultConfig()
	cfg.Service = "test krpc"
	//cfg.IsDebug=false
	cfg.InitLog()

	krpcc.Init()
	c := krpcc.NewJSONRPCClient("tcp://0.0.0.0:8008")

	result := krpctest.Result{}

	cdc := amino.NewCodec()
	cdc.RegisterConcrete(result, "rrr", nil)
	c.SetCodec(cdc)

	if err := c.Call("hello_world", krpcc.M{"name": "hello", "num": 345}, &result); err != nil {
		panic(err.Error())
	}

	fmt.Println(result.Result)
}
