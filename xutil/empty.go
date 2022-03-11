package xutil

import (
	"reflect"
	"time"
)

// apiString is used for type assert api for String().
type apiString interface {
	String() string
}

// apiInterfaces is used for type assert api for Interfaces.
type apiInterfaces interface {
	Interfaces() []interface{}
}

// apiMapStrAny is the interface support for converting struct parameter to map.
type apiMapStrAny interface {
	MapStrAny() map[string]interface{}
}

type apiTime interface {
	Date() (year int, month time.Month, day int)
	IsZero() bool
}

// IsEmpty [影响性能] 返回value是否为空，如value为0,nil,false,"",len(slice/map/chan) == 0会返回true，否则返回false
func IsEmpty(value interface{}) bool {
	if value == nil {
		return true
	}
	switch value := value.(type) {
	case int:
		return value == 0
	case int8:
		return value == 0
	case int16:
		return value == 0
	case int32:
		return value == 0
	case int64:
		return value == 0
	case uint:
		return value == 0
	case uint8:
		return value == 0
	case uint16:
		return value == 0
	case uint32:
		return value == 0
	case uint64:
		return value == 0
	case float32:
		return value == 0
	case float64:
		return value == 0
	case bool:
		return !value
	case string:
		return value == ""
	case []byte:
		return len(value) == 0
	case []rune:
		return len(value) == 0
	case []int:
		return len(value) == 0
	case []string:
		return len(value) == 0
	case []float32:
		return len(value) == 0
	case []float64:
		return len(value) == 0
	case map[string]interface{}:
		return len(value) == 0
	default:
		if f, ok := value.(apiTime); ok {
			if f == nil {
				return true
			}
			return f.IsZero()
		}
		if f, ok := value.(apiString); ok {
			if f == nil {
				return true
			}
			return f.String() == ""
		}
		if f, ok := value.(apiInterfaces); ok {
			if f == nil {
				return true
			}
			return len(f.Interfaces()) == 0
		}
		if f, ok := value.(apiMapStrAny); ok {
			if f == nil {
				return true
			}
			return len(f.MapStrAny()) == 0
		}
		var rv reflect.Value
		if v, ok := value.(reflect.Value); ok {
			rv = v
		} else {
			rv = reflect.ValueOf(value)
		}
		switch rv.Kind() {
		case reflect.Bool:
			return !rv.Bool()
		case reflect.Int,
			reflect.Int8,
			reflect.Int16,
			reflect.Int32,
			reflect.Int64:
			return rv.Int() == 0
		case reflect.Uint,
			reflect.Uint8,
			reflect.Uint16,
			reflect.Uint32,
			reflect.Uint64,
			reflect.Uintptr:
			return rv.Uint() == 0
		case reflect.Float32,
			reflect.Float64:
			return rv.Float() == 0
		case reflect.String:
			return rv.Len() == 0
		case reflect.Struct:
			for i := 0; i < rv.NumField(); i++ {
				if !IsEmpty(rv) {
					return false
				}
			}
			return true
		case reflect.Chan,
			reflect.Map,
			reflect.Slice,
			reflect.Array:
			return rv.Len() == 0
		case reflect.Func,
			reflect.Ptr,
			reflect.Interface,
			reflect.UnsafePointer:
			if rv.IsNil() {
				return true
			}
		}
	}
	return false
}
