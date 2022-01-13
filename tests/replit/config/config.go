package config

import (
	"time"
)

// Timeouts 非optiongen生成测试
type Timeouts struct {
	ReadTimeout  time.Duration `xconf:"read_timeout" default:"5s"`
	WriteTimeout time.Duration `xconf:"write_timeout" default:"10s"`
	ConnTimeout  time.Duration `xconf:"conn_timeout" default:"20s"`
}

// ETCDOptionDeclareWithDefault go-lint
//go:generate optiongen --option_with_struct_name=true --xconf=true --empty_composite_nil=true --usage_tag_name=usage --xconf=true
func ETCDOptionDeclareWithDefault() interface{} {
	return map[string]interface{}{
		"Endpoints":       []string{"10.0.0.1", "10.0.0.2"},
		"TimeoutsPointer": (*Timeouts)(&Timeouts{}),
	}
}

// RedisOptionDeclareWithDefault go-lint
//go:generate optiongen --option_with_struct_name=true --xconf=true --empty_composite_nil=true --usage_tag_name=usage --xconf=true
func RedisOptionDeclareWithDefault() interface{} {
	return map[string]interface{}{
		"Endpoints":      []string{"192.168.0.1", "192.168.0.2"},
		"Cluster":        true,
		"TimeoutsStruct": (Timeouts)(Timeouts{}),
	}
}

// ConfigOptionDeclareWithDefault go-lint
//go:generate optiongen --option_with_struct_name=false  --xconf=true --empty_composite_nil=true --usage_tag_name=usage --xconf=true
func ConfigOptionDeclareWithDefault() interface{} {
	return map[string]interface{}{
		"TypeBool":     false,
		"TypeString":   "a",
		"TypeDuration": time.Duration(time.Second),

		"TypeFloat32": float32(32.32),
		"TypeFloat64": float32(64.64),

		"TypeInt":    32,
		"TypeUint":   32,
		"TypeInt8":   int8(8),
		"TypeUint8":  uint8(8),
		"TypeInt16":  int16(16),
		"TypeUint16": uint16(16),
		"TypeInt32":  int32(32),
		"TypeUint32": uint32(32),
		"TypeInt64":  int64(64),
		"TypeUint64": uint64(64),

		"TypeSliceInt":      []int{1, 2, 3, 4},
		"TypeSliceUint":     []uint{1, 2, 3, 4},
		"TypeSliceInt8":     []int8{1, 2, 3, 4},
		"TypeSliceUint8":    []uint8{1, 2, 3, 4},
		"TypeSliceInt16":    []int16{1, 2, 3, 4},
		"TypeSliceUin16":    []uint16{1, 2, 3, 4},
		"TypeSliceInt32":    []int32{1, 2, 3, 4},
		"TypeSliceUint32":   []uint32{1, 2, 3, 4},
		"TypeSliceInt64":    []int64{1, 2, 3, 4},
		"TypeSliceUint64":   []uint64{1, 2, 3, 4},
		"TypeSliceString":   []string{"a", "b", "c"},
		"TypeSliceFloat32":  []float32{1.32, 2.32, 3.32, 4.32},
		"TypeSliceFloat64":  []float64{1.64, 2.64, 3.64, 4.64},
		"TypeSliceDuratuon": []time.Duration([]time.Duration{time.Second, time.Minute, time.Hour}),
		// annotation@TypeMapStringIntNotLeaf(xconf="type_map_string_int_not_leaf,notleaf")
		"TypeMapStringIntNotLeaf": map[string]int{"a": 1, "b": 2},
		"TypeMapStringInt":        map[string]int{"a": 1, "b": 2},
		"TypeMapIntString":        map[int]string{1: "a", 2: "b"},
		"TypeMapStringString":     map[string]string{"a": "a", "b": "b"},
		"TypeMapIntInt":           map[int]int{1: 1, 2: 2},
		"TypeMapStringDuration":   map[string]time.Duration(map[string]time.Duration{"read": time.Second, "write": time.Second * 5}),
		// annotation@Redis(getter="RedisVisitor")
		"Redis":         (*Redis)(NewRedis()),
		"ETCD":          (*ETCD)(NewETCD()),
		"TestInterface": (interface{})(nil),
	}
}
