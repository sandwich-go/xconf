package xflag

import (
	"flag"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/sandwich-go/xconf/xfield"
	"github.com/sandwich-go/xconf/xflag/vars"
)

type Maker struct {
	cc *Options
	fs *flag.FlagSet
}

func NewMaker(opts ...Option) *Maker {
	cc := NewOptions(opts...)
	return &Maker{
		cc: cc,
		fs: cc.FlagSet,
	}
}

func ParseArgs(obj interface{}, args []string, opts ...Option) ([]string, error) {
	fm := NewMaker(opts...)
	return fm.ParseArgs(obj, args)
}

func (fm *Maker) FlagKeys() []string {
	var keys []string
	fm.fs.VisitAll(func(ff *flag.Flag) {
		keys = append(keys, ff.Name)
	})
	return keys
}

func containsString(s []string, v string) bool {
	for _, vv := range s {
		if vv == v {
			return true
		}
	}
	return false
}

func (fm *Maker) EnvKeysMapping(validKeys []string) map[string]string {
	keyMap := make(map[string]string)
	fm.fs.VisitAll(func(ff *flag.Flag) {
		if !containsString(validKeys, ff.Name) {
			return
		}
		keyMap[strings.ToUpper(strings.ReplaceAll(ff.Name, ".", "_"))] = ff.Name
	})
	return keyMap
}

func (fm *Maker) PrintDefaults() {
	fm.fs.PrintDefaults()
}
func (fm *Maker) FlagSet() *flag.FlagSet    { return fm.cc.FlagSet }
func (fm *Maker) Parse(args []string) error { return fm.cc.FlagSet.Parse(args) }
func (fm *Maker) Set(obj interface{}) error {
	v := reflect.ValueOf(obj)
	if v.Kind() != reflect.Ptr {
		return fmt.Errorf("top level object must be a pointer. %v is passed", v.Type())
	}
	if v.IsNil() {
		return fmt.Errorf("top level object cannot be nil")
	}
	switch e := v.Elem(); e.Kind() {
	case reflect.Struct:
		fm.enumerateAndCreate("", nil, e, "")
	case reflect.Interface:
		if e.Elem().Kind() == reflect.Ptr {
			fm.enumerateAndCreate("", nil, e, "")
		} else {
			return fmt.Errorf("interface must have pointer underlying type. %v is passed", v.Type())
		}
	default:
		return fmt.Errorf("object must be a pointer to struct or interface. %v is passed", v.Type())
	}
	return nil
}

// ParseArgs parses the arguments based on the FlagMaker's setting.
func (fm *Maker) ParseArgs(obj interface{}, args []string) ([]string, error) {
	err := fm.Set(obj)
	if err != nil {
		return args, err
	}
	err = fm.fs.Parse(args)
	return fm.fs.Args(), err
}
func (fm *Maker) warningCanNotCreate(path string, typeStr string) {
	if containsString(fm.cc.FlagSetIgnore, path) {
		return
	}
	if containsString(fm.cc.FlagSetIgnore, typeStr) {
		return
	}
	fm.cc.LogWarning(fmt.Sprintf("xflag(%s): got unsupported type, not create to FlagSet, path: %s type_str:%s", fm.cc.Name, path, typeStr))
}
func usage(provider flag.Getter, prefix string, usageFromTag string) string {
	if usageFromTag != "" {
		return usageFromTag
	}
	if u, ok := provider.(interface{ Usage() string }); ok {
		return u.Usage()
	}
	return prefix
}
func (fm *Maker) enumerateAndCreate(prefix string, tags xfield.TagList, value reflect.Value, usageFromTag string) {
	switch value.Kind() {
	case
		// do no create flag for these types
		reflect.Uintptr,
		reflect.UnsafePointer,
		reflect.Array,
		reflect.Chan,
		reflect.Func:
		fm.warningCanNotCreate(prefix, reflect.TypeOf(value.Interface()).Name())
		return
	case reflect.Map:
		keyName := reflect.TypeOf(value.Interface()).Key().Name()
		valName := reflect.TypeOf(value.Interface()).Elem().Name()
		typeName := fmt.Sprintf("map[%s]%s", keyName, valName)
		provider, ok := fm.cc.FlagValueProvider(prefix, typeName, value.Addr().Interface())
		if !ok {
			fm.warningCanNotCreate(prefix, typeName)
			return
		}
		fm.fs.Var(provider, prefix, usage(provider, prefix, usageFromTag))
		return
	case reflect.Slice:
		typeName := fmt.Sprintf("[]%s", reflect.TypeOf(value.Interface()).Elem().Name())
		provider, ok := fm.cc.FlagValueProvider(prefix, typeName, value.Addr().Interface())
		if !ok {
			fm.warningCanNotCreate(prefix, typeName)
			return
		}
		fm.fs.Var(provider, prefix, usage(provider, prefix, usageFromTag))
		return
	case
		// Basic value types
		reflect.String,
		reflect.Bool,
		reflect.Float32, reflect.Float64,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		fm.defineFlag(prefix, value, usageFromTag)
		return
	case reflect.Interface:
		if value.IsNil() {
			fm.warningCanNotCreate(prefix, reflect.TypeOf(value.Interface()).Name())
			return
		}
		fm.enumerateAndCreate(prefix, tags, value.Elem(), usageFromTag)
		return
	case reflect.Ptr:
		if value.IsNil() {
			value.Set(reflect.New(value.Type().Elem()))
		}
		fm.enumerateAndCreate(prefix, tags, value.Elem(), usageFromTag)
		return
	case reflect.Struct:
		// keep going
	default:
		fm.warningCanNotCreate(prefix, reflect.TypeOf(value.Interface()).Name())
		return
	}

	numFields := value.NumField()
	tt := value.Type()

	for i := 0; i < numFields; i++ {
		stField := tt.Field(i)
		// Skip unexported fields, as only exported fields can be set. This is similar to how json and yaml work.
		if stField.PkgPath != "" && !stField.Anonymous {
			continue
		}
		if stField.Anonymous && fm.getUnderlyingType(stField.Type).Kind() != reflect.Struct {
			continue
		}
		field := value.Field(i)
		optName, tags := fm.getName(stField)
		usage := stField.Tag.Get(fm.cc.UsageTagName)
		if len(prefix) > 0 && !fm.cc.Flatten {
			optName = prefix + "." + optName
		}
		fm.enumerateAndCreate(optName, tags, field, usage)
	}
}
func (fm *Maker) getName(field reflect.StructField) (string, xfield.TagList) {
	name, tags := xfield.ParseTag(field.Tag.Get(fm.cc.TagName))
	if len(name) == 0 {
		if field.Anonymous {
			name = fm.getUnderlyingType(field.Type).Name()
		} else {
			name = field.Name
		}
	}
	return fm.cc.KeyFormat(name), tags
}

func (fm *Maker) getUnderlyingType(ttype reflect.Type) reflect.Type {
	// this only deals with *T unnamed type, other unnamed types, e.g. []int, struct{}
	// will return empty string.
	if ttype.Kind() == reflect.Ptr {
		return fm.getUnderlyingType(ttype.Elem())
	}
	return ttype
}

var (
	stringPtrType  = reflect.TypeOf((*string)(nil))
	boolPtrType    = reflect.TypeOf((*bool)(nil))
	float32PtrType = reflect.TypeOf((*float32)(nil))
	float64PtrType = reflect.TypeOf((*float64)(nil))
	intPtrType     = reflect.TypeOf((*int)(nil))
	int8PtrType    = reflect.TypeOf((*int8)(nil))
	int16PtrType   = reflect.TypeOf((*int16)(nil))
	int32PtrType   = reflect.TypeOf((*int32)(nil))
	int64PtrType   = reflect.TypeOf((*int64)(nil))
	uintPtrType    = reflect.TypeOf((*uint)(nil))
	uint8PtrType   = reflect.TypeOf((*uint8)(nil))
	uint16PtrType  = reflect.TypeOf((*uint16)(nil))
	uint32PtrType  = reflect.TypeOf((*uint32)(nil))
	uint64PtrType  = reflect.TypeOf((*uint64)(nil))
)

func (fm *Maker) defineFlag(name string, value reflect.Value, usageFromTag string) {
	usage := usageFromTag
	if usage == "" {
		usage = name
	}
	// v must be scalar, otherwise panic
	ptrValue := value.Addr()
	switch value.Kind() {
	case reflect.String:
		v := ptrValue.Convert(stringPtrType).Interface().(*string)
		fm.fs.StringVar(v, name, value.String(), usage)
	case reflect.Bool:
		v := ptrValue.Convert(boolPtrType).Interface().(*bool)
		fm.fs.BoolVar(v, name, value.Bool(), usage)
	case reflect.Int:
		v := ptrValue.Convert(intPtrType).Interface().(*int)
		fm.fs.IntVar(v, name, int(value.Int()), usage)
	case reflect.Int8:
		v := ptrValue.Convert(int8PtrType).Interface().(*int8)
		fm.fs.Var(vars.NewInt8(v), name, usage)
	case reflect.Int16:
		v := ptrValue.Convert(int16PtrType).Interface().(*int16)
		fm.fs.Var(vars.NewInt16(v), name, usage)
	case reflect.Int32:
		v := ptrValue.Convert(int32PtrType).Interface().(*int32)
		fm.fs.Var(vars.NewInt32(v), name, usage)
	case reflect.Int64:
		switch v := ptrValue.Interface().(type) {
		case *int64:
			fm.fs.Int64Var(v, name, value.Int(), usage)
		case *time.Duration:
			fm.fs.DurationVar(v, name, value.Interface().(time.Duration), usage)
		default:
			vv := ptrValue.Convert(int64PtrType).Interface().(*int64)
			fm.fs.Int64Var(vv, name, value.Int(), usage)
		}
	case reflect.Float32:
		v := ptrValue.Convert(float32PtrType).Interface().(*float32)
		fm.fs.Var(vars.NewFloat32(v), name, usage)
	case reflect.Float64:
		v := ptrValue.Convert(float64PtrType).Interface().(*float64)
		fm.fs.Float64Var(v, name, value.Float(), usage)
	case reflect.Uint:
		v := ptrValue.Convert(uintPtrType).Interface().(*uint)
		fm.fs.UintVar(v, name, uint(value.Uint()), usage)
	case reflect.Uint8:
		v := ptrValue.Convert(uint8PtrType).Interface().(*uint8)
		fm.fs.Var(vars.NewUint8(v), name, usage)
	case reflect.Uint16:
		v := ptrValue.Convert(uint16PtrType).Interface().(*uint16)
		fm.fs.Var(vars.NewUint16(v), name, usage)
	case reflect.Uint32:
		v := ptrValue.Convert(uint32PtrType).Interface().(*uint32)
		fm.fs.Var(vars.NewUint32(v), name, usage)
	case reflect.Uint64:
		v := ptrValue.Convert(uint64PtrType).Interface().(*uint64)
		fm.fs.Uint64Var(v, name, value.Uint(), usage)
	}
}
