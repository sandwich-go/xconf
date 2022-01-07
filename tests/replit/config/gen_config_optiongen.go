// Code generated by optiongen. DO NOT EDIT.
// optiongen: github.com/timestee/optiongen

package config

import (
	"sync/atomic"
	"time"
	"unsafe"
)

// Config struct
type Config struct {
	TypeBool              bool                     `xconf:"type_bool"`
	TypeString            string                   `xconf:"type_string"`
	TypeDuration          time.Duration            `xconf:"type_duration"`
	TypeFloat32           float32                  `xconf:"type_float32"`
	TypeFloat64           float32                  `xconf:"type_float64"`
	TypeInt               int                      `xconf:"type_int"`
	TypeUint              int                      `xconf:"type_uint"`
	TypeInt8              int8                     `xconf:"type_int8"`
	TypeUint8             uint8                    `xconf:"type_uint8"`
	TypeInt16             int16                    `xconf:"type_int16"`
	TypeUint16            uint16                   `xconf:"type_uint16"`
	TypeInt32             int32                    `xconf:"type_int32"`
	TypeUint32            uint32                   `xconf:"type_uint32"`
	TypeInt64             int64                    `xconf:"type_int64"`
	TypeUint64            uint64                   `xconf:"type_uint64"`
	TypeSliceInt          []int                    `xconf:"type_slice_int"`
	TypeSliceUint         []uint                   `xconf:"type_slice_uint"`
	TypeSliceInt8         []int8                   `xconf:"type_slice_int8"`
	TypeSliceUint8        []uint8                  `xconf:"type_slice_uint8"`
	TypeSliceInt16        []int16                  `xconf:"type_slice_int16"`
	TypeSliceUin16        []uint16                 `xconf:"type_slice_uin16"`
	TypeSliceInt32        []int32                  `xconf:"type_slice_int32"`
	TypeSliceUint32       []uint32                 `xconf:"type_slice_uint32"`
	TypeSliceInt64        []int64                  `xconf:"type_slice_int64"`
	TypeSliceUint64       []uint64                 `xconf:"type_slice_uint64"`
	TypeSliceString       []string                 `xconf:"type_slice_string"`
	TypeSliceFloat32      []float32                `xconf:"type_slice_float32"`
	TypeSliceFloat64      []float64                `xconf:"type_slice_float64"`
	TypeSliceDuratuon     []time.Duration          `xconf:"type_slice_duratuon"`
	TypeMapStringInt      map[string]int           `xconf:"type_map_string_int"`
	TypeMapIntString      map[int]string           `xconf:"type_map_int_string"`
	TypeMapStringString   map[string]string        `xconf:"type_map_string_string"`
	TypeMapIntInt         map[int]int              `xconf:"type_map_int_int"`
	TypeMapStringDuration map[string]time.Duration `xconf:"type_map_string_duration"`
	// annotation@Redis(getter=&#34;RedisVisitor&#34;)
	Redis         *Redis      `xconf:"redis"`
	ETCD          *ETCD       `xconf:"etcd"`
	TestInterface interface{} `xconf:"test_interface"`
}

// SetOption apply single option
func (cc *Config) SetOption(opt ConfigOption) {
	_ = opt(cc)
}

// ApplyOption apply mutiple options
func (cc *Config) ApplyOption(opts ...ConfigOption) {
	for _, opt := range opts {
		_ = opt(cc)
	}
}

// GetSetOption apply new option and return the old optuon
// sample:
// old := cc.GetSetOption(WithTimeout(time.Second))
// defer cc.SetOption(old)
func (cc *Config) GetSetOption(opt ConfigOption) ConfigOption {
	return opt(cc)
}

// ConfigOption option func
type ConfigOption func(cc *Config) ConfigOption

// WithTypeBool option func for TypeBool
func WithTypeBool(v bool) ConfigOption {
	return func(cc *Config) ConfigOption {
		previous := cc.TypeBool
		cc.TypeBool = v
		return WithTypeBool(previous)
	}
}

// WithTypeString option func for TypeString
func WithTypeString(v string) ConfigOption {
	return func(cc *Config) ConfigOption {
		previous := cc.TypeString
		cc.TypeString = v
		return WithTypeString(previous)
	}
}

// WithTypeDuration option func for TypeDuration
func WithTypeDuration(v time.Duration) ConfigOption {
	return func(cc *Config) ConfigOption {
		previous := cc.TypeDuration
		cc.TypeDuration = v
		return WithTypeDuration(previous)
	}
}

// WithTypeFloat32 option func for TypeFloat32
func WithTypeFloat32(v float32) ConfigOption {
	return func(cc *Config) ConfigOption {
		previous := cc.TypeFloat32
		cc.TypeFloat32 = v
		return WithTypeFloat32(previous)
	}
}

// WithTypeFloat64 option func for TypeFloat64
func WithTypeFloat64(v float32) ConfigOption {
	return func(cc *Config) ConfigOption {
		previous := cc.TypeFloat64
		cc.TypeFloat64 = v
		return WithTypeFloat64(previous)
	}
}

// WithTypeInt option func for TypeInt
func WithTypeInt(v int) ConfigOption {
	return func(cc *Config) ConfigOption {
		previous := cc.TypeInt
		cc.TypeInt = v
		return WithTypeInt(previous)
	}
}

// WithTypeUint option func for TypeUint
func WithTypeUint(v int) ConfigOption {
	return func(cc *Config) ConfigOption {
		previous := cc.TypeUint
		cc.TypeUint = v
		return WithTypeUint(previous)
	}
}

// WithTypeInt8 option func for TypeInt8
func WithTypeInt8(v int8) ConfigOption {
	return func(cc *Config) ConfigOption {
		previous := cc.TypeInt8
		cc.TypeInt8 = v
		return WithTypeInt8(previous)
	}
}

// WithTypeUint8 option func for TypeUint8
func WithTypeUint8(v uint8) ConfigOption {
	return func(cc *Config) ConfigOption {
		previous := cc.TypeUint8
		cc.TypeUint8 = v
		return WithTypeUint8(previous)
	}
}

// WithTypeInt16 option func for TypeInt16
func WithTypeInt16(v int16) ConfigOption {
	return func(cc *Config) ConfigOption {
		previous := cc.TypeInt16
		cc.TypeInt16 = v
		return WithTypeInt16(previous)
	}
}

// WithTypeUint16 option func for TypeUint16
func WithTypeUint16(v uint16) ConfigOption {
	return func(cc *Config) ConfigOption {
		previous := cc.TypeUint16
		cc.TypeUint16 = v
		return WithTypeUint16(previous)
	}
}

// WithTypeInt32 option func for TypeInt32
func WithTypeInt32(v int32) ConfigOption {
	return func(cc *Config) ConfigOption {
		previous := cc.TypeInt32
		cc.TypeInt32 = v
		return WithTypeInt32(previous)
	}
}

// WithTypeUint32 option func for TypeUint32
func WithTypeUint32(v uint32) ConfigOption {
	return func(cc *Config) ConfigOption {
		previous := cc.TypeUint32
		cc.TypeUint32 = v
		return WithTypeUint32(previous)
	}
}

// WithTypeInt64 option func for TypeInt64
func WithTypeInt64(v int64) ConfigOption {
	return func(cc *Config) ConfigOption {
		previous := cc.TypeInt64
		cc.TypeInt64 = v
		return WithTypeInt64(previous)
	}
}

// WithTypeUint64 option func for TypeUint64
func WithTypeUint64(v uint64) ConfigOption {
	return func(cc *Config) ConfigOption {
		previous := cc.TypeUint64
		cc.TypeUint64 = v
		return WithTypeUint64(previous)
	}
}

// WithTypeSliceInt option func for TypeSliceInt
func WithTypeSliceInt(v ...int) ConfigOption {
	return func(cc *Config) ConfigOption {
		previous := cc.TypeSliceInt
		cc.TypeSliceInt = v
		return WithTypeSliceInt(previous...)
	}
}

// WithTypeSliceUint option func for TypeSliceUint
func WithTypeSliceUint(v ...uint) ConfigOption {
	return func(cc *Config) ConfigOption {
		previous := cc.TypeSliceUint
		cc.TypeSliceUint = v
		return WithTypeSliceUint(previous...)
	}
}

// WithTypeSliceInt8 option func for TypeSliceInt8
func WithTypeSliceInt8(v ...int8) ConfigOption {
	return func(cc *Config) ConfigOption {
		previous := cc.TypeSliceInt8
		cc.TypeSliceInt8 = v
		return WithTypeSliceInt8(previous...)
	}
}

// WithTypeSliceUint8 option func for TypeSliceUint8
func WithTypeSliceUint8(v ...uint8) ConfigOption {
	return func(cc *Config) ConfigOption {
		previous := cc.TypeSliceUint8
		cc.TypeSliceUint8 = v
		return WithTypeSliceUint8(previous...)
	}
}

// WithTypeSliceInt16 option func for TypeSliceInt16
func WithTypeSliceInt16(v ...int16) ConfigOption {
	return func(cc *Config) ConfigOption {
		previous := cc.TypeSliceInt16
		cc.TypeSliceInt16 = v
		return WithTypeSliceInt16(previous...)
	}
}

// WithTypeSliceUin16 option func for TypeSliceUin16
func WithTypeSliceUin16(v ...uint16) ConfigOption {
	return func(cc *Config) ConfigOption {
		previous := cc.TypeSliceUin16
		cc.TypeSliceUin16 = v
		return WithTypeSliceUin16(previous...)
	}
}

// WithTypeSliceInt32 option func for TypeSliceInt32
func WithTypeSliceInt32(v ...int32) ConfigOption {
	return func(cc *Config) ConfigOption {
		previous := cc.TypeSliceInt32
		cc.TypeSliceInt32 = v
		return WithTypeSliceInt32(previous...)
	}
}

// WithTypeSliceUint32 option func for TypeSliceUint32
func WithTypeSliceUint32(v ...uint32) ConfigOption {
	return func(cc *Config) ConfigOption {
		previous := cc.TypeSliceUint32
		cc.TypeSliceUint32 = v
		return WithTypeSliceUint32(previous...)
	}
}

// WithTypeSliceInt64 option func for TypeSliceInt64
func WithTypeSliceInt64(v ...int64) ConfigOption {
	return func(cc *Config) ConfigOption {
		previous := cc.TypeSliceInt64
		cc.TypeSliceInt64 = v
		return WithTypeSliceInt64(previous...)
	}
}

// WithTypeSliceUint64 option func for TypeSliceUint64
func WithTypeSliceUint64(v ...uint64) ConfigOption {
	return func(cc *Config) ConfigOption {
		previous := cc.TypeSliceUint64
		cc.TypeSliceUint64 = v
		return WithTypeSliceUint64(previous...)
	}
}

// WithTypeSliceString option func for TypeSliceString
func WithTypeSliceString(v ...string) ConfigOption {
	return func(cc *Config) ConfigOption {
		previous := cc.TypeSliceString
		cc.TypeSliceString = v
		return WithTypeSliceString(previous...)
	}
}

// WithTypeSliceFloat32 option func for TypeSliceFloat32
func WithTypeSliceFloat32(v ...float32) ConfigOption {
	return func(cc *Config) ConfigOption {
		previous := cc.TypeSliceFloat32
		cc.TypeSliceFloat32 = v
		return WithTypeSliceFloat32(previous...)
	}
}

// WithTypeSliceFloat64 option func for TypeSliceFloat64
func WithTypeSliceFloat64(v ...float64) ConfigOption {
	return func(cc *Config) ConfigOption {
		previous := cc.TypeSliceFloat64
		cc.TypeSliceFloat64 = v
		return WithTypeSliceFloat64(previous...)
	}
}

// WithTypeSliceDuratuon option func for TypeSliceDuratuon
func WithTypeSliceDuratuon(v ...time.Duration) ConfigOption {
	return func(cc *Config) ConfigOption {
		previous := cc.TypeSliceDuratuon
		cc.TypeSliceDuratuon = v
		return WithTypeSliceDuratuon(previous...)
	}
}

// WithTypeMapStringInt option func for TypeMapStringInt
func WithTypeMapStringInt(v map[string]int) ConfigOption {
	return func(cc *Config) ConfigOption {
		previous := cc.TypeMapStringInt
		cc.TypeMapStringInt = v
		return WithTypeMapStringInt(previous)
	}
}

// WithTypeMapIntString option func for TypeMapIntString
func WithTypeMapIntString(v map[int]string) ConfigOption {
	return func(cc *Config) ConfigOption {
		previous := cc.TypeMapIntString
		cc.TypeMapIntString = v
		return WithTypeMapIntString(previous)
	}
}

// WithTypeMapStringString option func for TypeMapStringString
func WithTypeMapStringString(v map[string]string) ConfigOption {
	return func(cc *Config) ConfigOption {
		previous := cc.TypeMapStringString
		cc.TypeMapStringString = v
		return WithTypeMapStringString(previous)
	}
}

// WithTypeMapIntInt option func for TypeMapIntInt
func WithTypeMapIntInt(v map[int]int) ConfigOption {
	return func(cc *Config) ConfigOption {
		previous := cc.TypeMapIntInt
		cc.TypeMapIntInt = v
		return WithTypeMapIntInt(previous)
	}
}

// WithTypeMapStringDuration option func for TypeMapStringDuration
func WithTypeMapStringDuration(v map[string]time.Duration) ConfigOption {
	return func(cc *Config) ConfigOption {
		previous := cc.TypeMapStringDuration
		cc.TypeMapStringDuration = v
		return WithTypeMapStringDuration(previous)
	}
}

// WithRedis option func for Redis
func WithRedis(v *Redis) ConfigOption {
	return func(cc *Config) ConfigOption {
		previous := cc.Redis
		cc.Redis = v
		return WithRedis(previous)
	}
}

// WithETCD option func for ETCD
func WithETCD(v *ETCD) ConfigOption {
	return func(cc *Config) ConfigOption {
		previous := cc.ETCD
		cc.ETCD = v
		return WithETCD(previous)
	}
}

// WithTestInterface option func for TestInterface
func WithTestInterface(v interface{}) ConfigOption {
	return func(cc *Config) ConfigOption {
		previous := cc.TestInterface
		cc.TestInterface = v
		return WithTestInterface(previous)
	}
}

// NewConfig(opts... ConfigOption) new Config
func NewConfig(opts ...ConfigOption) *Config {
	cc := newDefaultConfig()

	for _, opt := range opts {
		_ = opt(cc)
	}
	if watchDogConfig != nil {
		watchDogConfig(cc)
	}
	return cc
}

// InstallConfigWatchDog the installed func will called when NewConfig(opts... ConfigOption)  called
func InstallConfigWatchDog(dog func(cc *Config)) {
	watchDogConfig = dog
}

// watchDogConfig global watch dog
var watchDogConfig func(cc *Config)

// newDefaultConfig new default Config
func newDefaultConfig() *Config {
	cc := &Config{}

	for _, opt := range [...]ConfigOption{
		WithTypeBool(false),
		WithTypeString("a"),
		WithTypeDuration(time.Second),
		WithTypeFloat32(32.32),
		WithTypeFloat64(64.64),
		WithTypeInt(32),
		WithTypeUint(32),
		WithTypeInt8(8),
		WithTypeUint8(8),
		WithTypeInt16(16),
		WithTypeUint16(16),
		WithTypeInt32(32),
		WithTypeUint32(32),
		WithTypeInt64(64),
		WithTypeUint64(64),
		WithTypeSliceInt([]int{1, 2, 3, 4}...),
		WithTypeSliceUint([]uint{1, 2, 3, 4}...),
		WithTypeSliceInt8([]int8{1, 2, 3, 4}...),
		WithTypeSliceUint8([]uint8{1, 2, 3, 4}...),
		WithTypeSliceInt16([]int16{1, 2, 3, 4}...),
		WithTypeSliceUin16([]uint16{1, 2, 3, 4}...),
		WithTypeSliceInt32([]int32{1, 2, 3, 4}...),
		WithTypeSliceUint32([]uint32{1, 2, 3, 4}...),
		WithTypeSliceInt64([]int64{1, 2, 3, 4}...),
		WithTypeSliceUint64([]uint64{1, 2, 3, 4}...),
		WithTypeSliceString([]string{"a", "b", "c"}...),
		WithTypeSliceFloat32([]float32{1.32, 2.32, 3.32, 4.32}...),
		WithTypeSliceFloat64([]float64{1.64, 2.64, 3.64, 4.64}...),
		WithTypeSliceDuratuon([]time.Duration{time.Second, time.Minute, time.Hour}...),
		WithTypeMapStringInt(map[string]int{"a": 1, "b": 2}),
		WithTypeMapIntString(map[int]string{1: "a", 2: "b"}),
		WithTypeMapStringString(map[string]string{"a": "a", "b": "b"}),
		WithTypeMapIntInt(map[int]int{1: 1, 2: 2}),
		WithTypeMapStringDuration(map[string]time.Duration{"read": time.Second, "write": time.Second * 5}),
		WithRedis(NewRedis()),
		WithETCD(NewETCD()),
		WithTestInterface(nil),
	} {
		_ = opt(cc)
	}

	return cc
}

// AtomicSetFunc used for XConf
func (cc *Config) AtomicSetFunc() func(interface{}) { return AtomicConfigSet }

// atomicConfig global *Config holder
var atomicConfig unsafe.Pointer

// AtomicConfigSet atomic setter for *Config
func AtomicConfigSet(update interface{}) {
	atomic.StorePointer(&atomicConfig, (unsafe.Pointer)(update.(*Config)))
}

// AtomicConfig return atomic *Config visitor
func AtomicConfig() ConfigVisitor {
	current := (*Config)(atomic.LoadPointer(&atomicConfig))
	if current == nil {
		atomic.CompareAndSwapPointer(&atomicConfig, nil, (unsafe.Pointer)(newDefaultConfig()))
		return (*Config)(atomic.LoadPointer(&atomicConfig))
	}
	return current
}

// all getter func
// GetTypeBool return struct field: TypeBool
func (cc *Config) GetTypeBool() bool { return cc.TypeBool }

// GetTypeString return struct field: TypeString
func (cc *Config) GetTypeString() string { return cc.TypeString }

// GetTypeDuration return struct field: TypeDuration
func (cc *Config) GetTypeDuration() time.Duration { return cc.TypeDuration }

// GetTypeFloat32 return struct field: TypeFloat32
func (cc *Config) GetTypeFloat32() float32 { return cc.TypeFloat32 }

// GetTypeFloat64 return struct field: TypeFloat64
func (cc *Config) GetTypeFloat64() float32 { return cc.TypeFloat64 }

// GetTypeInt return struct field: TypeInt
func (cc *Config) GetTypeInt() int { return cc.TypeInt }

// GetTypeUint return struct field: TypeUint
func (cc *Config) GetTypeUint() int { return cc.TypeUint }

// GetTypeInt8 return struct field: TypeInt8
func (cc *Config) GetTypeInt8() int8 { return cc.TypeInt8 }

// GetTypeUint8 return struct field: TypeUint8
func (cc *Config) GetTypeUint8() uint8 { return cc.TypeUint8 }

// GetTypeInt16 return struct field: TypeInt16
func (cc *Config) GetTypeInt16() int16 { return cc.TypeInt16 }

// GetTypeUint16 return struct field: TypeUint16
func (cc *Config) GetTypeUint16() uint16 { return cc.TypeUint16 }

// GetTypeInt32 return struct field: TypeInt32
func (cc *Config) GetTypeInt32() int32 { return cc.TypeInt32 }

// GetTypeUint32 return struct field: TypeUint32
func (cc *Config) GetTypeUint32() uint32 { return cc.TypeUint32 }

// GetTypeInt64 return struct field: TypeInt64
func (cc *Config) GetTypeInt64() int64 { return cc.TypeInt64 }

// GetTypeUint64 return struct field: TypeUint64
func (cc *Config) GetTypeUint64() uint64 { return cc.TypeUint64 }

// GetTypeSliceInt return struct field: TypeSliceInt
func (cc *Config) GetTypeSliceInt() []int { return cc.TypeSliceInt }

// GetTypeSliceUint return struct field: TypeSliceUint
func (cc *Config) GetTypeSliceUint() []uint { return cc.TypeSliceUint }

// GetTypeSliceInt8 return struct field: TypeSliceInt8
func (cc *Config) GetTypeSliceInt8() []int8 { return cc.TypeSliceInt8 }

// GetTypeSliceUint8 return struct field: TypeSliceUint8
func (cc *Config) GetTypeSliceUint8() []uint8 { return cc.TypeSliceUint8 }

// GetTypeSliceInt16 return struct field: TypeSliceInt16
func (cc *Config) GetTypeSliceInt16() []int16 { return cc.TypeSliceInt16 }

// GetTypeSliceUin16 return struct field: TypeSliceUin16
func (cc *Config) GetTypeSliceUin16() []uint16 { return cc.TypeSliceUin16 }

// GetTypeSliceInt32 return struct field: TypeSliceInt32
func (cc *Config) GetTypeSliceInt32() []int32 { return cc.TypeSliceInt32 }

// GetTypeSliceUint32 return struct field: TypeSliceUint32
func (cc *Config) GetTypeSliceUint32() []uint32 { return cc.TypeSliceUint32 }

// GetTypeSliceInt64 return struct field: TypeSliceInt64
func (cc *Config) GetTypeSliceInt64() []int64 { return cc.TypeSliceInt64 }

// GetTypeSliceUint64 return struct field: TypeSliceUint64
func (cc *Config) GetTypeSliceUint64() []uint64 { return cc.TypeSliceUint64 }

// GetTypeSliceString return struct field: TypeSliceString
func (cc *Config) GetTypeSliceString() []string { return cc.TypeSliceString }

// GetTypeSliceFloat32 return struct field: TypeSliceFloat32
func (cc *Config) GetTypeSliceFloat32() []float32 { return cc.TypeSliceFloat32 }

// GetTypeSliceFloat64 return struct field: TypeSliceFloat64
func (cc *Config) GetTypeSliceFloat64() []float64 { return cc.TypeSliceFloat64 }

// GetTypeSliceDuratuon return struct field: TypeSliceDuratuon
func (cc *Config) GetTypeSliceDuratuon() []time.Duration { return cc.TypeSliceDuratuon }

// GetTypeMapStringInt return struct field: TypeMapStringInt
func (cc *Config) GetTypeMapStringInt() map[string]int { return cc.TypeMapStringInt }

// GetTypeMapIntString return struct field: TypeMapIntString
func (cc *Config) GetTypeMapIntString() map[int]string { return cc.TypeMapIntString }

// GetTypeMapStringString return struct field: TypeMapStringString
func (cc *Config) GetTypeMapStringString() map[string]string { return cc.TypeMapStringString }

// GetTypeMapIntInt return struct field: TypeMapIntInt
func (cc *Config) GetTypeMapIntInt() map[int]int { return cc.TypeMapIntInt }

// GetTypeMapStringDuration return struct field: TypeMapStringDuration
func (cc *Config) GetTypeMapStringDuration() map[string]time.Duration {
	return cc.TypeMapStringDuration
}

// GetRedis return struct field: Redis
func (cc *Config) GetRedis() RedisVisitor { return cc.Redis }

// GetETCD return struct field: ETCD
func (cc *Config) GetETCD() *ETCD { return cc.ETCD }

// GetTestInterface return struct field: TestInterface
func (cc *Config) GetTestInterface() interface{} { return cc.TestInterface }

// ConfigVisitor visitor interface for Config
type ConfigVisitor interface {
	GetTypeBool() bool
	GetTypeString() string
	GetTypeDuration() time.Duration
	GetTypeFloat32() float32
	GetTypeFloat64() float32
	GetTypeInt() int
	GetTypeUint() int
	GetTypeInt8() int8
	GetTypeUint8() uint8
	GetTypeInt16() int16
	GetTypeUint16() uint16
	GetTypeInt32() int32
	GetTypeUint32() uint32
	GetTypeInt64() int64
	GetTypeUint64() uint64
	GetTypeSliceInt() []int
	GetTypeSliceUint() []uint
	GetTypeSliceInt8() []int8
	GetTypeSliceUint8() []uint8
	GetTypeSliceInt16() []int16
	GetTypeSliceUin16() []uint16
	GetTypeSliceInt32() []int32
	GetTypeSliceUint32() []uint32
	GetTypeSliceInt64() []int64
	GetTypeSliceUint64() []uint64
	GetTypeSliceString() []string
	GetTypeSliceFloat32() []float32
	GetTypeSliceFloat64() []float64
	GetTypeSliceDuratuon() []time.Duration
	GetTypeMapStringInt() map[string]int
	GetTypeMapIntString() map[int]string
	GetTypeMapStringString() map[string]string
	GetTypeMapIntInt() map[int]int
	GetTypeMapStringDuration() map[string]time.Duration
	GetRedis() RedisVisitor
	GetETCD() *ETCD
	GetTestInterface() interface{}
}
