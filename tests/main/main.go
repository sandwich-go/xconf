package main

import (
	"fmt"

	"github.com/sandwich-go/xconf"
	"github.com/sandwich-go/xconf/tests"
)

func main() {
	cc := tests.NewTestConfig()
	cc.Redis.RedisAddress = "127.0.0.1:6637"
	if err := xconf.Parse(cc, xconf.WithFiles("c1.yaml"), xconf.WithDebug(true)); err != nil {
		panic(err)
	}
	fmt.Println("cc.RedisAsPointer.RedisAddress ", cc.RedisAsPointer.RedisAddress)
	fmt.Println("cc.Redis.RedisAddress ", cc.Redis.RedisAddress)
	xconf.DumpInfo()
}
