// Code generated by optiongen. DO NOT EDIT.
// optiongen: github.com/timestee/optiongen

package tests

import (
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/sandwich-go/xconf/tests/redis"
)

// Config should use NewTestConfig to initialize it
type Config struct {
	OptionUsage string         `xconf:"option_usage"`
	HttpAddress string         `xconf:"http_address"`
	Map1        map[string]int `xconf:"map1" usage:"k,v使用,分割"`
	// annotation@Map1(comment="k,v使用,分割")
	// annotation@MapNotLeaf(xconf="map_not_leaf,notleaf",deprecated="使用Map1")
	MapNotLeaf      map[string]int  `xconf:"map_not_leaf,notleaf,deprecated" usage:"Deprecated: 使用Map1"`
	ProcessCount    int8            `xconf:"process_count"`
	MaxUint64       uint64          `xconf:"max_uint64"`
	MaxInt          int             `xconf:"max_int"`
	Int8            int8            `xconf:"int8"`
	TimeDurations   []time.Duration `xconf:"time_durations" usage:"延迟队列"`
	DefaultEmptyMap map[string]int  `xconf:"default_empty_map"`
	Int64Slice      []int64         `xconf:"int64_slice"`
	Float64Slice    []float64       `xconf:"float64_slice"`
	Uin64Slice      []uint64        `xconf:"uin64_slice"`
	StringSlice     []string        `xconf:"string_slice"`
	ReadTimeout     time.Duration   `xconf:"read_timeout"`
	SubTest         SubTest         `xconf:"sub_test"`
	TestBool        bool            `xconf:"test_bool"`
	TestBoolTrue    bool            `xconf:"test_bool_true"`
	RedisAsPointer  *Redis          `xconf:"redis_as_pointer"`
	Redis           Redis           `xconf:"redis"`
	RedisTimeout    *RedisTimeout   `xconf:"redis_timeout"`
}

// NewTestConfig new Config
func NewTestConfig(opts ...ConfigOption) *Config {
	cc := newDefaultConfig()
	for _, opt := range opts {
		opt(cc)
	}
	if watchDogConfig != nil {
		watchDogConfig(cc)
	}
	return cc
}

// ApplyOption apply mutiple new option and return the old ones
// sample:
// old := cc.ApplyOption(WithTimeout(time.Second))
// defer cc.ApplyOption(old...)
func (cc *Config) ApplyOption(opts ...ConfigOption) []ConfigOption {
	var previous []ConfigOption
	for _, opt := range opts {
		previous = append(previous, opt(cc))
	}
	return previous
}

// ConfigOption option func
type ConfigOption func(cc *Config) ConfigOption

// WithOptionUsage option func for filed OptionUsage
func WithOptionUsage(v string) ConfigOption {
	return func(cc *Config) ConfigOption {
		previous := cc.OptionUsage
		cc.OptionUsage = v
		return WithOptionUsage(previous)
	}
}

// WithHttpAddress option func for filed HttpAddress
func WithHttpAddress(v string) ConfigOption {
	return func(cc *Config) ConfigOption {
		previous := cc.HttpAddress
		cc.HttpAddress = v
		return WithHttpAddress(previous)
	}
}

// WithMap1 k,v使用,分割
func WithMap1(v map[string]int) ConfigOption {
	return func(cc *Config) ConfigOption {
		previous := cc.Map1
		cc.Map1 = v
		return WithMap1(previous)
	}
}

// WithMapNotLeaf option func for filed MapNotLeaf
//
// Deprecated: 使用Map1
func WithMapNotLeaf(v map[string]int) ConfigOption {
	return func(cc *Config) ConfigOption {
		previous := cc.MapNotLeaf
		cc.MapNotLeaf = v
		return WithMapNotLeaf(previous)
	}
}

// WithProcessCount option func for filed ProcessCount
func WithProcessCount(v int8) ConfigOption {
	return func(cc *Config) ConfigOption {
		previous := cc.ProcessCount
		cc.ProcessCount = v
		return WithProcessCount(previous)
	}
}

// WithMaxUint64 option func for filed MaxUint64
func WithMaxUint64(v uint64) ConfigOption {
	return func(cc *Config) ConfigOption {
		previous := cc.MaxUint64
		cc.MaxUint64 = v
		return WithMaxUint64(previous)
	}
}

// WithMaxInt option func for filed MaxInt
func WithMaxInt(v int) ConfigOption {
	return func(cc *Config) ConfigOption {
		previous := cc.MaxInt
		cc.MaxInt = v
		return WithMaxInt(previous)
	}
}

// WithInt8 option func for filed Int8
func WithInt8(v int8) ConfigOption {
	return func(cc *Config) ConfigOption {
		previous := cc.Int8
		cc.Int8 = v
		return WithInt8(previous)
	}
}

// WithTimeDurations 延迟队列
func WithTimeDurations(v ...time.Duration) ConfigOption {
	return func(cc *Config) ConfigOption {
		previous := cc.TimeDurations
		cc.TimeDurations = v
		return WithTimeDurations(previous...)
	}
}

// WithDefaultEmptyMap option func for filed DefaultEmptyMap
func WithDefaultEmptyMap(v map[string]int) ConfigOption {
	return func(cc *Config) ConfigOption {
		previous := cc.DefaultEmptyMap
		cc.DefaultEmptyMap = v
		return WithDefaultEmptyMap(previous)
	}
}

// WithInt64Slice option func for filed Int64Slice
func WithInt64Slice(v ...int64) ConfigOption {
	return func(cc *Config) ConfigOption {
		previous := cc.Int64Slice
		cc.Int64Slice = v
		return WithInt64Slice(previous...)
	}
}

// WithFloat64Slice option func for filed Float64Slice
func WithFloat64Slice(v ...float64) ConfigOption {
	return func(cc *Config) ConfigOption {
		previous := cc.Float64Slice
		cc.Float64Slice = v
		return WithFloat64Slice(previous...)
	}
}

// WithUin64Slice option func for filed Uin64Slice
func WithUin64Slice(v ...uint64) ConfigOption {
	return func(cc *Config) ConfigOption {
		previous := cc.Uin64Slice
		cc.Uin64Slice = v
		return WithUin64Slice(previous...)
	}
}

// WithStringSlice option func for filed StringSlice
func WithStringSlice(v ...string) ConfigOption {
	return func(cc *Config) ConfigOption {
		previous := cc.StringSlice
		cc.StringSlice = v
		return WithStringSlice(previous...)
	}
}

// WithReadTimeout option func for filed ReadTimeout
func WithReadTimeout(v time.Duration) ConfigOption {
	return func(cc *Config) ConfigOption {
		previous := cc.ReadTimeout
		cc.ReadTimeout = v
		return WithReadTimeout(previous)
	}
}

// WithSubTest option func for filed SubTest
func WithSubTest(v SubTest) ConfigOption {
	return func(cc *Config) ConfigOption {
		previous := cc.SubTest
		cc.SubTest = v
		return WithSubTest(previous)
	}
}

// WithTestBool option func for filed TestBool
func WithTestBool(v bool) ConfigOption {
	return func(cc *Config) ConfigOption {
		previous := cc.TestBool
		cc.TestBool = v
		return WithTestBool(previous)
	}
}

// WithTestBoolTrue option func for filed TestBoolTrue
func WithTestBoolTrue(v bool) ConfigOption {
	return func(cc *Config) ConfigOption {
		previous := cc.TestBoolTrue
		cc.TestBoolTrue = v
		return WithTestBoolTrue(previous)
	}
}

// WithRedisAsPointer option func for filed RedisAsPointer
func WithRedisAsPointer(v *Redis) ConfigOption {
	return func(cc *Config) ConfigOption {
		previous := cc.RedisAsPointer
		cc.RedisAsPointer = v
		return WithRedisAsPointer(previous)
	}
}

// WithRedis option func for filed Redis
func WithRedis(v Redis) ConfigOption {
	return func(cc *Config) ConfigOption {
		previous := cc.Redis
		cc.Redis = v
		return WithRedis(previous)
	}
}

// WithRedisTimeout option func for filed RedisTimeout
func WithRedisTimeout(v *RedisTimeout) ConfigOption {
	return func(cc *Config) ConfigOption {
		previous := cc.RedisTimeout
		cc.RedisTimeout = v
		return WithRedisTimeout(previous)
	}
}

// InstallConfigWatchDog the installed func will called when NewTestConfig  called
func InstallConfigWatchDog(dog func(cc *Config)) { watchDogConfig = dog }

// watchDogConfig global watch dog
var watchDogConfig func(cc *Config)

// newDefaultConfig new default Config
func newDefaultConfig() *Config {
	cc := &Config{}

	for _, opt := range [...]ConfigOption{
		WithOptionUsage(optionUsage),
		WithHttpAddress(":3001"),
		WithMap1(map[string]int{"test1": 100, "test2": 200}),
		WithMapNotLeaf(map[string]int{"test1": 100, "test2": 200}),
		WithProcessCount(1),
		WithMaxUint64(0),
		WithMaxInt(0),
		WithInt8(1),
		WithTimeDurations([]time.Duration{time.Second, time.Second}...),
		WithDefaultEmptyMap(nil),
		WithInt64Slice([]int64{101, 202, 303}...),
		WithFloat64Slice([]float64{101.191, 202.202, 303.303}...),
		WithUin64Slice([]uint64{101, 202, 303}...),
		WithStringSlice([]string{"test1", "test2", "test3"}...),
		WithReadTimeout(time.Second * time.Duration(5)),
		WithSubTest(SubTest{}),
		WithTestBool(false),
		WithTestBoolTrue(true),
		WithRedisAsPointer(&redis.Conf{}),
		WithRedis(redis.Conf{}),
		WithRedisTimeout(&redis.Timeout{}),
	} {
		opt(cc)
	}

	return cc
}

// AtomicSetFunc used for XConf
func (cc *Config) AtomicSetFunc() func(interface{}) { return AtomicConfigSet }

// atomicConfig global *Config holder
var atomicConfig unsafe.Pointer

// onAtomicConfigSet global call back when  AtomicConfigSet called by XConf.
// use ConfigInterface.ApplyOption to modify the updated cc
// if passed in cc not valid, then return false, cc will not set to atomicConfig
var onAtomicConfigSet func(cc ConfigInterface) bool

// InstallCallbackOnAtomicConfigSet install callback
func InstallCallbackOnAtomicConfigSet(callback func(cc ConfigInterface) bool) {
	onAtomicConfigSet = callback
}

// AtomicConfigSet atomic setter for *Config
func AtomicConfigSet(update interface{}) {
	cc := update.(*Config)
	if onAtomicConfigSet != nil && !onAtomicConfigSet(cc) {
		return
	}
	atomic.StorePointer(&atomicConfig, (unsafe.Pointer)(cc))
}

// AtomicConfig return atomic *ConfigVisitor
func AtomicConfig() ConfigVisitor {
	current := (*Config)(atomic.LoadPointer(&atomicConfig))
	if current == nil {
		defaultOne := newDefaultConfig()
		if watchDogConfig != nil {
			watchDogConfig(defaultOne)
		}
		atomic.CompareAndSwapPointer(&atomicConfig, nil, (unsafe.Pointer)(defaultOne))
		return (*Config)(atomic.LoadPointer(&atomicConfig))
	}
	return current
}

// all getter func
func (cc *Config) GetOptionUsage() string  { return cc.OptionUsage }
func (cc *Config) GetHttpAddress() string  { return cc.HttpAddress }
func (cc *Config) GetMap1() map[string]int { return cc.Map1 }

// GetMapNotLeaf visitor func for filed MapNotLeaf
//
// Deprecated: 使用Map1
func (cc *Config) GetMapNotLeaf() map[string]int      { return cc.MapNotLeaf }
func (cc *Config) GetProcessCount() int8              { return cc.ProcessCount }
func (cc *Config) GetMaxUint64() uint64               { return cc.MaxUint64 }
func (cc *Config) GetMaxInt() int                     { return cc.MaxInt }
func (cc *Config) GetInt8() int8                      { return cc.Int8 }
func (cc *Config) GetTimeDurations() []time.Duration  { return cc.TimeDurations }
func (cc *Config) GetDefaultEmptyMap() map[string]int { return cc.DefaultEmptyMap }
func (cc *Config) GetInt64Slice() []int64             { return cc.Int64Slice }
func (cc *Config) GetFloat64Slice() []float64         { return cc.Float64Slice }
func (cc *Config) GetUin64Slice() []uint64            { return cc.Uin64Slice }
func (cc *Config) GetStringSlice() []string           { return cc.StringSlice }
func (cc *Config) GetReadTimeout() time.Duration      { return cc.ReadTimeout }
func (cc *Config) GetSubTest() SubTest                { return cc.SubTest }
func (cc *Config) GetTestBool() bool                  { return cc.TestBool }
func (cc *Config) GetTestBoolTrue() bool              { return cc.TestBoolTrue }
func (cc *Config) GetRedisAsPointer() *Redis          { return cc.RedisAsPointer }
func (cc *Config) GetRedis() Redis                    { return cc.Redis }
func (cc *Config) GetRedisTimeout() *RedisTimeout     { return cc.RedisTimeout }

// ConfigVisitor visitor interface for Config
type ConfigVisitor interface {
	GetOptionUsage() string
	GetHttpAddress() string
	GetMap1() map[string]int
	// GetMapNotLeaf visitor func for filed MapNotLeaf
	//
	// Deprecated: 使用Map1
	GetMapNotLeaf() map[string]int
	GetProcessCount() int8
	GetMaxUint64() uint64
	GetMaxInt() int
	GetInt8() int8
	GetTimeDurations() []time.Duration
	GetDefaultEmptyMap() map[string]int
	GetInt64Slice() []int64
	GetFloat64Slice() []float64
	GetUin64Slice() []uint64
	GetStringSlice() []string
	GetReadTimeout() time.Duration
	GetSubTest() SubTest
	GetTestBool() bool
	GetTestBoolTrue() bool
	GetRedisAsPointer() *Redis
	GetRedis() Redis
	GetRedisTimeout() *RedisTimeout
}

// ConfigInterface visitor + ApplyOption interface for Config
type ConfigInterface interface {
	ConfigVisitor
	ApplyOption(...ConfigOption) []ConfigOption
}
