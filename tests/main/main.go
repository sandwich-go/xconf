package main

import (
	"fmt"

	"github.com/sandwich-go/xconf"
	"github.com/sandwich-go/xconf/tests"
)

func main() {
	cc := tests.NewTestConfig()
	cc.HttpAddress = "127.0.0.1"
	cc.Redis.RedisAddress = "127.0.0.1:6637"
	xx := xconf.New(
		xconf.WithFiles("c1.yaml"),
		xconf.WithDebug(true),
		xconf.WithEnvironPrefix("test_prefix_"),
	)
	if err := xx.Parse(cc); err != nil {
		panic(err)
	}
	fmt.Println("cc.RedisAsPointer.RedisAddress ", cc.RedisAsPointer.RedisAddress)
	fmt.Println("cc.Redis.RedisAddress ", cc.Redis.RedisAddress)
	fmt.Println("cc.ProcessCount ", cc.ProcessCount)
	fmt.Println("cc.MaxUint64 ", cc.MaxUint64)
	fmt.Println("cc.MaxInt ", cc.MaxInt)
	fmt.Println("cc.ReadTimeout ", cc.ReadTimeout)
	xx.Usage()
}
