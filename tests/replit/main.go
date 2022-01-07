package main

import (
	"bytes"
	"fmt"
	"os"

	"github.com/sandwich-go/xconf"
	"github.com/sandwich-go/xconf/tests/replit/config"
)

func panicErr(err error) {
	if err != nil {
		panic(fmt.Sprintf("### replit panic ### %v", err))
	}
}

func main() {
	os.Setenv("TEST_ENV_HOST", "test_env_host_127.0.0.1")
	// parse default config,c1.yaml inherit from c1.toml
	err := xconf.Parse(config.AtomicConfig(), xconf.WithFiles("./c1.yaml"),
		xconf.WithFlagArgs(
			"--type_string=type_string_value_from_flag",
			"--type_slice_duratuon=1s,2s,3s,4s,5s",
			"--redis.timeouts_struct.read_timeout=200000s",
			"--type_map_string_int=${TEST_ENV_HOST},1",
		),
		xconf.WithDebug(false),
	)
	panicErr(err)

	fmt.Println(config.AtomicConfig().GetETCD())                                  // {[10.0.0.1 10.0.0.2 10.0.0.3 10.0.0.4]
	fmt.Println(config.AtomicConfig().GetRedis())                                 // {[192.168.0.1 192.168.0.2] true {16m40s 16m40s 16m40s}}
	fmt.Println(config.AtomicConfig().GetTypeString())                            // type_string_value_from_flag
	fmt.Println(config.AtomicConfig().GetTypeSliceDuratuon())                     // [1s 2s 3s 4s 5s]
	fmt.Println(config.AtomicConfig().GetRedis().GetTimeoutsStruct().ReadTimeout) // 16m40s
	fmt.Println(config.AtomicConfig().GetTypeMapStringInt())                      // map[test_env_host_127.0.0.1:1]

	// save to bytes and parse again
	bb := xconf.MustSaveToBytes(xconf.ConfigTypeYAML)

	// upate with reader
	panicErr(xconf.UpdateWithReader(bytes.NewReader(bb)))
	// upate with filed path
	panicErr(xconf.UpdateWithFieldPathValues("redis.timeouts_struct.read_timeout", "1s", "redis.timeouts_struct.write_timeout", "120s"))
	fmt.Println(config.AtomicConfig().GetRedis().GetTimeoutsStruct().ReadTimeout)  // 1s
	fmt.Println(config.AtomicConfig().GetRedis().GetTimeoutsStruct().WriteTimeout) // 2m0s

	empteOne := &config.Config{ETCD: &config.ETCD{}, Redis: &config.Redis{}}
	x1 := xconf.NewWithoutFlagEnv(xconf.WithReaders(bytes.NewBuffer(bb)))
	panicErr(x1.Parse(empteOne))
	fmt.Println("empty config etcd : ", empteOne.GetETCD())  // {[10.0.0.1 10.0.0.2 10.0.0.3 10.0.0.4]
	fmt.Println("empty config redis: ", empteOne.GetRedis()) // {[192.168.0.1 192.168.0.2] true {16m40s 16m40s 16m40s}}

}
