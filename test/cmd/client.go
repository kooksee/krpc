package main

import (
	"fmt"
	"github.com/kooksee/krpc/client"
	"github.com/kooksee/krpc/test"
	"github.com/tendermint/go-amino"
)

func main() {
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
