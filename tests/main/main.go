package main

import (
	"fmt"

	"github.com/sandwich-go/xconf"
	"github.com/sandwich-go/xconf/tests"
)

func main() {
	cc := tests.NewTestConfig()
	cc.Redis.RedisAddress = "127.0.0.1:6637"
	xx := xconf.New(xconf.WithFiles("c1.yaml"), xconf.WithDebug(false))
	if err := xx.Parse(cc); err != nil {
		panic(err)
	}
	fmt.Println("cc.RedisAsPointer.RedisAddress ", cc.RedisAsPointer.RedisAddress)
	fmt.Println("cc.Redis.RedisAddress ", cc.Redis.RedisAddress)
	xx.Usage()
}
