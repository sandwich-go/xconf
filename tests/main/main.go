package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/sandwich-go/xconf"
	"github.com/sandwich-go/xconf/tests"
)

func main() {
	cc := tests.NewTestConfig()
	cc.SubTest.Map2 = make(map[string]int)
	cc.SubTest.Map2["test1"] = 1
	cc.DefaultEmptyMap = make(map[string]int)
	cc.DefaultEmptyMap["test1"] = 1
	cc.SubTest.Servers = make(map[string]tests.Server)
	cc.SubTest.Servers["s1"] = tests.Server{
		Timeouts: map[string]time.Duration{"read": time.Second * time.Duration(5)},
	}
	xconf.Parse(cc, xconf.WithDebug(true))
	x := xconf.New(xconf.WithDebug(true), xconf.WithFiles("c2.toml"), xconf.WithFlagSet(flag.NewFlagSet("test", flag.ContinueOnError)))
	err := x.Parse(cc)
	fmt.Println("result :", cc, err)
	x.DumpInfo()
}
