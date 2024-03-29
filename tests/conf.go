package tests

import (
	"time"

	"github.com/sandwich-go/xconf/tests/redis"
)

// Server 测试配置
type Server struct {
	Timeouts map[string]time.Duration `xconf:"timeouts"`
}

// SubTest 测试配置
type SubTest struct {
	HTTPAddress string            `xconf:"http_address"`
	MapNotLeaf  map[string]int    `xconf:"map_not_leaf,notleaf"`
	Map2        map[string]int    `xconf:"map2"`
	Map3        map[string]int    `xconf:"map3"`
	Slice2      []int64           `xconf:"slice2"`
	Servers     map[string]Server `xconf:"servers,notleaf"`
}

// Redis 测试别名
type Redis = redis.Conf

// RedisTimeout 测试别名
type RedisTimeout = redis.Timeout

var optionUsage = `在这里描述一些应用级别的配置规则`

// ConfigOptionDeclareWithDefault go-lint
//
//go:generate optiongen --option_with_struct_name=false --new_func=NewTestConfig --xconf=true --empty_composite_nil=true --usage_tag_name=usage
func ConfigOptionDeclareWithDefault() interface{} {
	return map[string]interface{}{
		"OptionUsage": string(optionUsage),
		"HttpAddress": ":3001",
		"Map1": map[string]int{
			"test1": 100,
			"test2": 200,
		}, // annotation@Map1(comment="k,v使用,分割")
		// annotation@MapNotLeaf(xconf="map_not_leaf,notleaf",deprecated="使用Map1")
		"MapNotLeaf": map[string]int{
			"test1": 100,
			"test2": 200,
		},
		"ProcessCount":    int8(1),
		"MaxUint64":       uint64(0),
		"MaxInt":          int(0),
		"Bytes":           []byte(nil),
		"Int8":            int8(1),
		"TimeDurations":   []time.Duration([]time.Duration{time.Second, time.Second}), // @MethodComment(延迟队列)
		"DefaultEmptyMap": map[string]int{},
		"Int64Slice":      []int64{101, 202, 303},
		"Float64Slice":    []float64{101.191, 202.202, 303.303},
		"Uin64Slice":      []uint64{101, 202, 303},
		"StringSlice":     []string{"test1", "test2", "test3"},
		"ReadTimeout":     time.Duration(time.Second * time.Duration(5)),
		"SubTest":         SubTest(SubTest{}),
		"TestBool":        false,
		"TestBoolTrue":    true,
		"RedisAsPointer":  (*Redis)(&redis.Conf{}),
		"Redis":           (Redis)(redis.Conf{}),
		"RedisTimeout":    (*RedisTimeout)(&redis.Timeout{}),
	}
}
