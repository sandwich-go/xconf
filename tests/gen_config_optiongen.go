// Code generated by optiongen. DO NOT EDIT.
// optiongen: github.com/timestee/optiongen

package tests

import (
	"sync/atomic"
	"time"
	"unsafe"
)

type Config struct {
	HttpAddress     string          `xconf:"http_address"`
	Map1            map[string]int  `xconf:"map1"`
	MapNotLeaf      map[string]int  `xconf:"map_not_leaf,notleaf" usage:"k,v使用,分割, 测试特殊符号：\"test\""`
	TimeDurations   []time.Duration `xconf:"time_durations" usage:"延迟队列"`
	DefaultEmptyMap map[string]int  `xconf:"default_empty_map"`
	Int64Slice      []int64         `xconf:"int64_slice"`
	Float64Slice    []float64       `xconf:"float64_slice"`
	Uin64Slice      []uint64        `xconf:"uin64_slice"`
	StringSlice     []string        `xconf:"string_slice"`
	ReadTimeout     time.Duration   `xconf:"read_timeout"`
	SubTest         SubTest         `xconf:"sub_test"`
	TestBool        bool            `xconf:"test_bool"`
}

func (cc *Config) SetOption(opt ConfigOption) {
	_ = opt(cc)
}

func (cc *Config) ApplyOption(opts ...ConfigOption) {
	for _, opt := range opts {
		_ = opt(cc)
	}
}

func (cc *Config) GetSetOption(opt ConfigOption) ConfigOption {
	return opt(cc)
}

type ConfigOption func(cc *Config) ConfigOption

func WithHttpAddress(v string) ConfigOption {
	return func(cc *Config) ConfigOption {
		previous := cc.HttpAddress
		cc.HttpAddress = v
		return WithHttpAddress(previous)
	}
}

func WithMap1(v map[string]int) ConfigOption {
	return func(cc *Config) ConfigOption {
		previous := cc.Map1
		cc.Map1 = v
		return WithMap1(previous)
	}
}

// k,v使用,分割, 测试特殊符号：&#34;test&#34;
func WithMapNotLeaf(v map[string]int) ConfigOption {
	return func(cc *Config) ConfigOption {
		previous := cc.MapNotLeaf
		cc.MapNotLeaf = v
		return WithMapNotLeaf(previous)
	}
}

// 延迟队列
func WithTimeDurations(v ...time.Duration) ConfigOption {
	return func(cc *Config) ConfigOption {
		previous := cc.TimeDurations
		cc.TimeDurations = v
		return WithTimeDurations(previous...)
	}
}

func WithDefaultEmptyMap(v map[string]int) ConfigOption {
	return func(cc *Config) ConfigOption {
		previous := cc.DefaultEmptyMap
		cc.DefaultEmptyMap = v
		return WithDefaultEmptyMap(previous)
	}
}

func WithInt64Slice(v ...int64) ConfigOption {
	return func(cc *Config) ConfigOption {
		previous := cc.Int64Slice
		cc.Int64Slice = v
		return WithInt64Slice(previous...)
	}
}

func WithFloat64Slice(v ...float64) ConfigOption {
	return func(cc *Config) ConfigOption {
		previous := cc.Float64Slice
		cc.Float64Slice = v
		return WithFloat64Slice(previous...)
	}
}

func WithUin64Slice(v ...uint64) ConfigOption {
	return func(cc *Config) ConfigOption {
		previous := cc.Uin64Slice
		cc.Uin64Slice = v
		return WithUin64Slice(previous...)
	}
}

func WithStringSlice(v ...string) ConfigOption {
	return func(cc *Config) ConfigOption {
		previous := cc.StringSlice
		cc.StringSlice = v
		return WithStringSlice(previous...)
	}
}

func WithReadTimeout(v time.Duration) ConfigOption {
	return func(cc *Config) ConfigOption {
		previous := cc.ReadTimeout
		cc.ReadTimeout = v
		return WithReadTimeout(previous)
	}
}

func WithSubTest(v SubTest) ConfigOption {
	return func(cc *Config) ConfigOption {
		previous := cc.SubTest
		cc.SubTest = v
		return WithSubTest(previous)
	}
}

func WithTestBool(v bool) ConfigOption {
	return func(cc *Config) ConfigOption {
		previous := cc.TestBool
		cc.TestBool = v
		return WithTestBool(previous)
	}
}

func NewTestConfig(opts ...ConfigOption) *Config {
	cc := newDefaultConfig()

	for _, opt := range opts {
		_ = opt(cc)
	}
	if watchDogConfig != nil {
		watchDogConfig(cc)
	}
	return cc
}

func InstallConfigWatchDog(dog func(cc *Config)) {
	watchDogConfig = dog
}

var watchDogConfig func(cc *Config)

func newDefaultConfig() *Config {

	cc := &Config{}

	for _, opt := range [...]ConfigOption{
		WithHttpAddress(":3001"),
		WithMap1(map[string]int{"test1": 100, "test2": 200}),
		WithMapNotLeaf(map[string]int{"test1": 100, "test2": 200}),
		WithTimeDurations([]time.Duration{time.Second, time.Second}...),
		WithDefaultEmptyMap(nil),
		WithInt64Slice([]int64{101, 202, 303}...),
		WithFloat64Slice([]float64{101.191, 202.202, 303.303}...),
		WithUin64Slice([]uint64{101, 202, 303}...),
		WithStringSlice([]string{"test1", "test2", "test3"}...),
		WithReadTimeout(time.Second * time.Duration(5)),
		WithSubTest(SubTest{}),
		WithTestBool(false),
	} {
		_ = opt(cc)
	}

	return cc
}

func (cc *Config) AtomicSetFunc() func(interface{}) { return AtomicConfigSet }

var atomicConfig unsafe.Pointer

func AtomicConfigSet(update interface{}) {
	atomic.StorePointer(&atomicConfig, (unsafe.Pointer)(update.(*Config)))
}

func AtomicConfig() ConfigInterface {
	current := (*Config)(atomic.LoadPointer(&atomicConfig))
	if current == nil {
		atomic.CompareAndSwapPointer(&atomicConfig, nil, (unsafe.Pointer)(newDefaultConfig()))
		return (*Config)(atomic.LoadPointer(&atomicConfig))
	}
	return current
}

// all getter func
func (cc *Config) GetHttpAddress() string             { return cc.HttpAddress }
func (cc *Config) GetMap1() map[string]int            { return cc.Map1 }
func (cc *Config) GetMapNotLeaf() map[string]int      { return cc.MapNotLeaf }
func (cc *Config) GetTimeDurations() []time.Duration  { return cc.TimeDurations }
func (cc *Config) GetDefaultEmptyMap() map[string]int { return cc.DefaultEmptyMap }
func (cc *Config) GetInt64Slice() []int64             { return cc.Int64Slice }
func (cc *Config) GetFloat64Slice() []float64         { return cc.Float64Slice }
func (cc *Config) GetUin64Slice() []uint64            { return cc.Uin64Slice }
func (cc *Config) GetStringSlice() []string           { return cc.StringSlice }
func (cc *Config) GetReadTimeout() time.Duration      { return cc.ReadTimeout }
func (cc *Config) GetSubTest() SubTest                { return cc.SubTest }
func (cc *Config) GetTestBool() bool                  { return cc.TestBool }

// interface for Config
type ConfigInterface interface {
	GetHttpAddress() string
	GetMap1() map[string]int
	GetMapNotLeaf() map[string]int
	GetTimeDurations() []time.Duration
	GetDefaultEmptyMap() map[string]int
	GetInt64Slice() []int64
	GetFloat64Slice() []float64
	GetUin64Slice() []uint64
	GetStringSlice() []string
	GetReadTimeout() time.Duration
	GetSubTest() SubTest
	GetTestBool() bool
}
