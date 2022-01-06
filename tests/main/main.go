package main

import (
	"flag"
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
	if err := xconf.Parse(cc, xconf.WithDebug(true)); err != nil {
		panic(err)
	}
	x := xconf.New(xconf.WithDebug(true), xconf.WithFiles("c2.toml"), xconf.WithFlagSet(flag.NewFlagSet("test", flag.ContinueOnError)))
	if err := x.Parse(cc); err != nil {
		panic(err)
	}
	x.DumpInfo()
}
