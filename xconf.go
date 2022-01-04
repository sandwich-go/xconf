package xconf

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"sync"

	"reflect"

	"github.com/mitchellh/hashstructure/v2"
	"github.com/mitchellh/mapstructure"
	"github.com/sandwich-go/xconf/xflag"
)

type XConf struct {
	zeroValPtrForLayout interface{}                    //只是用户获取结构信息，不存储数据信息
	cc                  *Options                       // 配置参数
	fieldPathInfoMap    map[string]StructFieldPathInfo // FieldPath到字段属性映射
	dataLatestCached    map[string]interface{}         // 缓存的最新数据
	dynamicUpdate       sync.Mutex
	kvs                 []*kvLoader
	updated             chan interface{}
	runningLogger       func(string)
	hasParsed           bool
	mapOnFieldUpdated   map[string]OnFieldUpdated
	changes             Changes
	atomicSetFunc       AtomicSetFunc
}

func New(opts ...Option) *XConf { return NewWithConf(NewOptions(opts...)) }
func NewWithoutFlagEnv(opts ...Option) *XConf {
	return New(append(opts, WithFlagSet(nil), WithEnviron())...)
}

func NewWithConf(cc *Options) *XConf {
	x := &XConf{cc: cc}
	x.updated = make(chan interface{}, 1)
	x.clean()
	x.mapOnFieldUpdated = make(map[string]OnFieldUpdated)
	x.runningLogger = func(s string) {
		if x.cc.Debug {
			fmt.Print(s)
		}
	}
	x.atomicSetFunc = func(interface{}) {}
	return x
}

type AtomicSetFunc = func(interface{})

func (x *XConf) Latest() (interface{}, error) {
	if !x.hasParsed {
		return nil, errors.New("need Parse first")
	}
	zeroOne := reflect.New(reflect.ValueOf(x.zeroValPtrForLayout).Type().Elem()).Interface()
	return zeroOne, x.bindLatest(zeroOne)
}
func (x *XConf) Copy() *XConf { return NewWithConf(x.cc) }

// 不允许外部clean清空状态，配置的增量更新依赖于解析后的状态
func (x *XConf) clean() {
	x.fieldPathInfoMap = make(map[string]StructFieldPathInfo)
	x.dataLatestCached = make(map[string]interface{})
	x.changes.Changed = make(map[string]*Values)
}

func (x *XConf) NotifyUpdate() <-chan interface{} { return x.updated }

func (x *XConf) defaultXFlagOptions() []xflag.Option {
	return []xflag.Option{
		xflag.WithFlatten(false),
		xflag.WithFlagSet(x.cc.FlagSet),
		xflag.WithKeyFormat(x.cc.FieldTagConvertor),
		xflag.WithTagName(x.cc.TagName),
		xflag.WithLogDebug(xflag.LogFunc(x.cc.LogDebug)),
		xflag.WithLogWarning(xflag.LogFunc(x.cc.LogWarning)),
		xflag.WithFlagSetIgnore(x.cc.FieldFlagSetCreateIgnore...),
	}
}

func (x *XConf) ZeroStructKeysTagList(s interface{}) map[string]StructFieldPathInfo {
	_, v := NewStruct(reflect.New(reflect.ValueOf(s).Type().Elem()).Interface(), x.cc.TagName, x.cc.TagNameDefaultValue, x.cc.FieldTagConvertor).Map()
	return v
}

func (x *XConf) StructMapStructure(s interface{}) map[string]interface{} {
	v, _ := NewStruct(s, x.cc.TagName, x.cc.TagNameDefaultValue, x.cc.FieldTagConvertor).Map()
	return v
}

func (x *XConf) runningLogData(name string, data map[string]interface{}) {
	x.runningLogger(fmt.Sprintf("===========================> Data %s\n%v\n", name, data))
}

func (x *XConf) keysList() []string {
	var keys []string
	for k := range x.fieldPathInfoMap {
		keys = append(keys, k)
	}
	return keys
}

func (x *XConf) mergeToDest(dataName string, data map[string]interface{}) error {
	x.runningLogData(dataName, data)
	x.runningLogger(fmt.Sprintf("----> merge src:%s dst:%s\n", dataName, "dest"))
	err := mergeMap("", 0, x.runningLogger, data, x.dataLatestCached, x.isLeafFieldPath, nil, &x.changes)
	return wrapIfErr(err, "got error:%w while merge data:%s to data: dst", err, dataName)
}

func (x *XConf) mergeMap(srcName string, dstName string, src map[string]interface{}, dst map[string]interface{}) error {
	x.runningLogger(fmt.Sprintf("----> merge src:%s dst:%s\n", srcName, dstName))
	return mergeMap("", 0, x.runningLogger, src, dst, x.isLeafFieldPath, nil, nil)
}

func (x *XConf) parse(valPtr interface{}) (err error) {
	defer func() {
		if reason := recover(); reason != nil {
			err = fmt.Errorf("%v", reason)
		}
	}()
	if x.hasParsed {
		return errors.New("has parsed befor")
	}
	x.hasParsed = true
	// 检测传入的数值是否为指针
	if reflect.ValueOf(valPtr).Kind() != reflect.Ptr {
		return errors.New("unsupported type, pass in as ptr")
	}
	// 保留结构信息
	x.zeroValPtrForLayout = reflect.New(reflect.ValueOf(valPtr).Type().Elem()).Interface()

	x.runningLogger(fmt.Sprintf("=========> Parse With\n%v\n", valPtr))

	// 获取结构信息，后续的数据更新等依赖于该结构信息
	x.fieldPathInfoMap = x.ZeroStructKeysTagList(x.zeroValPtrForLayout)

	if reflect.DeepEqual(x.zeroValPtrForLayout, valPtr) && x.cc.ParseDefault {
		// 如果指定根据tag解析默认数值则进行一次解析操作,将解析得到的数值作为默认值
		// 如果input为空值，则不解析struct本身数值，struct中解析得到的是全量key-val的mapstructure，防止覆盖default
		panicErr(x.updateDstDataWithParseDefault(valPtr))
	} else {
		// 如果input不为空，则进行解析，input值完全覆盖default内的值,不再解析default
		panicErr(x.mergeToDest("struct_input", x.StructMapStructure(valPtr)))
	}
	//flag指定的文件 直接覆盖 内部指定的文件, 独立解析flagset数据以获取files
	flagData, filesToParse, err := x.parseFlagFilesForXConf(x.cc.FlagSet, x.cc.FlagArgs...)
	panicErr(err)
	panicErr(x.updateDstDataWithFiles(filesToParse...))
	panicErr(x.updateDstDataWithReaders(x.cc.Readers...))
	panicErr(x.commonUpdateDstData("flag", func() (map[string]interface{}, error) { return flagData, nil }))
	// panicErr(x.updateDstDataWithFlagSet(x.cc.FlagSet, x.cc.FlagArgs...))
	panicErr(x.updateDstDataWithEnviron(x.cc.Environ...))
	panicErr(x.bindLatest(valPtr))
	if w, ok := valPtr.(interface{ AtomicSetFunc() func(interface{}) }); ok {
		x.cc.LogDebug("install AtomicSetFunc")
		x.atomicSetFunc = w.AtomicSetFunc()
		x.atomicSetFunc(valPtr)
	}
	return nil
}

func (x *XConf) parseFlagFilesForXConf(flagSet *flag.FlagSet, args ...string) (flagData map[string]interface{}, filesToParse []string, err error) {
	filesToParse = x.cc.Files
	if x.cc.FlagSet == nil || len(x.cc.FlagArgs) == 0 {
		return
	}
	if fv := x.cc.FlagSet.Lookup(MetaKeyFiles); fv == nil {
		x.cc.FlagSet.String(MetaKeyFiles, "", "xconf files provided by flag, file slice, split by `,`.")
	}
	flagData, err = xflagMapstructure(
		x.zeroValPtrForLayout,
		append(x.keysList(), MetaKeyFiles),
		func(*xflag.Maker) []string { return x.cc.FlagArgs },
		append(x.defaultXFlagOptions(), xflag.WithFlagSet(x.cc.FlagSet))...)

	if err != nil {
		return
	}

	if v := flagData[MetaKeyFiles]; v != nil {
		filesToParse = strings.Split(strings.Trim(v.(string), " "), ",")
		x.cc.LogDebug(fmt.Sprintf("config files changed from:%v to %v provided by FlagSet", x.cc.Files, filesToParse))
	}

	delete(flagData, MetaKeyFiles)
	return
}

func (x *XConf) bindLatest(valPtr interface{}) (err error) {
	defer func() {
		if reason := recover(); reason != nil {
			err = fmt.Errorf("%v", reason)
		}
	}()
	if reflect.ValueOf(valPtr).Kind() != reflect.Ptr {
		return errors.New("unsupported type, pass in as ptr")
	}
	err = x.decode(x.dataLatestCached, valPtr)
	panicErrWithWrap(err, "got error:%w while decode using map structure", err)
	return
}

func (x *XConf) updateDstDataWithParseDefault(valPtr interface{}) (err error) {
	defer func() {
		if reason := recover(); reason != nil {
			err = fmt.Errorf("%v", reason)
		}
	}()
	var dataDefault map[string]interface{}
	var defaultParsed bool
	dataDefault, defaultParsed, err = x.Copy().parseDefault(valPtr)
	panicErrWithWrap(err, "got error:%w while parse default value", err)
	if !defaultParsed {
		return
	}
	x.cc.LogWarning(fmt.Sprintf("Parse Default From Tag:%s", x.cc.TagNameDefaultValue))
	panicErr(x.mergeToDest("default_from_tag", dataDefault))
	return
}

func (x *XConf) isLeafFieldPath(fieldPath string) bool {
	if x.cc.MapMerge {
		return false
	}
	return isLeafFieldPath(x.fieldPathInfoMap, fieldPath)
}

func (x *XConf) commonUpdateDstData(name string, f func() (map[string]interface{}, error)) (err error) {
	defer func() {
		if reason := recover(); reason != nil {
			err = fmt.Errorf("%v", reason)
		}
	}()
	data, err := f()
	if data == nil {
		return nil
	}
	panicErrWithWrap(err, "got error while load %s,err :%w ", name, err)
	panicErr(x.mergeToDest(name, data))
	return
}

func (x *XConf) updateDstDataWithFiles(files ...string) (err error) {
	if len(files) == 0 {
		return
	}
	return x.commonUpdateDstData("files", func() (map[string]interface{}, error) {
		return x.loadFiles(files...)
	})
}

func (x *XConf) updateDstDataWithReaders(readers ...io.Reader) (err error) {
	if len(readers) == 0 {
		return
	}
	return x.commonUpdateDstData("readers", func() (map[string]interface{}, error) {
		return x.loadReaders(readers...)
	})
}

func (x *XConf) updateDstDataWithFlagSet(flagSet *flag.FlagSet, args ...string) (err error) {
	if flagSet == nil || len(args) == 0 {
		return
	}
	return x.commonUpdateDstData("flag", func() (map[string]interface{}, error) {
		return xflagMapstructure(
			x.zeroValPtrForLayout,
			x.keysList(),
			func(*xflag.Maker) []string { return args },
			append(x.defaultXFlagOptions(), xflag.WithFlagSet(flagSet))...)
	})
}

func (x *XConf) updateDstDataWithEnviron(environ ...string) (err error) {
	if len(environ) == 0 {
		return
	}
	return x.commonUpdateDstData("env", func() (map[string]interface{}, error) {
		return xflagMapstructure(
			x.zeroValPtrForLayout,
			x.keysList(),
			func(xf *xflag.Maker) []string { return EnvBindToFlags(environ, xf.EnvKeysMapping(x.keysList())) },
			append(x.defaultXFlagOptions(), xflag.WithFlagSet(newFlagSet("Environ")))...)
	})

}

func (x *XConf) decode(data map[string]interface{}, valPtr interface{}) error {
	config := x.defaultDecoderConfig(valPtr)
	config.TagName = x.cc.TagName
	// xconf默认使用的SnakeCase规则转换filedName
	config.MatchName = func(mapKey, fieldName string) bool {
		equal := strings.EqualFold(mapKey, fieldName)
		if equal {
			return true
		}
		return x.cc.FieldTagConvertor(fieldName) == mapKey
	}
	var metadata mapstructure.Metadata
	config.Metadata = &metadata
	config.ErrorUnused = false // 依赖于metadata.Unused做错误提示,过滤特定字段
	for _, opt := range x.cc.DecoderConfigOption {
		opt(config)
	}

	err := decode(data, config)
	if err != nil {
		return fmt.Errorf("got error:%w while decode using mapstructure", err)
	}
	if len(metadata.Unused) > 0 {
		var unused []string
		var deprecated []string
		for _, v := range metadata.Unused {
			// metadata中预留的key 用于做一些基础功能
			if containString(MetaKeyList, v) {
				continue
			}
			// 逻辑层指定的Deprecated字段，报警
			if containString(x.cc.FieldPathDeprecated, v) {
				deprecated = append(deprecated, v)
				continue
			}
			unused = append(unused, v)
		}
		if len(deprecated) != 0 {
			x.cc.LogWarning(fmt.Sprintf("!!! DEPRECATED FIELD, WILL REMOVE IN FUTURE. FIELDS: %s", strings.Join(deprecated, ",")))
		}
		if len(unused) != 0 {
			return fmt.Errorf("!!! UNUSED FIELDS, SHOULD PAY ATTENTION. FIELDS: %s", strings.Join(unused, ","))
		}
	}
	return nil
}

func (x *XConf) Parse(valPtr interface{}) error {
	err := x.parse(valPtr)
	if err == nil {
		return nil
	}
	switch x.cc.ErrorHandling {
	case ContinueOnError:
		return err
	case ExitOnError:
		os.Exit(2)
	case PanicOnError:
		panic(err)
	}
	return nil
}

func sstr(len int, s string) (ss string) {
	for i := 0; i < len; i++ {
		ss += s
	}
	return ss
}

func (x *XConf) DumpInfo() {
	if x.zeroValPtrForLayout == nil {
		fmt.Printf(" Parse Frist\n")
		return
	}
	opts := append(x.defaultXFlagOptions(),
		xflag.WithFlagSet(newFlagSet("dump_info")),
		xflag.WithLogWarning((func(string) {})),
	)
	xf := xflag.NewMaker(opts...)
	xf.Set(x.zeroValPtrForLayout)
	keysMapping := xf.EnvKeysMapping(x.keysList())
	var keys []string
	maxLen := 5
	for k := range keysMapping {
		keys = append(keys, k)
		if len(keysMapping[k]) > maxLen {
			maxLen = len(keysMapping[k])
		}
	}
	maxLen += 6
	fmtStr := fmt.Sprintf("%%-%ds %%-%ds %%-%ds\n", 4, maxLen, maxLen)
	fmt.Printf(sstr((maxLen)*2, "-") + "\n")
	fmt.Printf(fmtStr, "#", "FLAG", "ENV")
	fmt.Printf(sstr((maxLen)*2, "-") + "\n")
	sort.Strings(keys)
	for i, k := range keys {
		fmt.Printf(fmtStr, fmt.Sprintf("%d", i+1), keysMapping[k], k)
	}
	fmt.Printf(sstr((maxLen)*2, "-") + "\n")
	fmt.Printf("# DataDest: %v\n", x.dataLatestCached)
	fmt.Printf(sstr((maxLen)*2, "-") + "\n")
	hashCode := x.Hash()
	fmt.Printf("# Hash Local  : %s\n", hashCode)
	hashCenter := DefaultInvalidHashString
	if center := x.dataLatestCached[MetaKeyLatestHash]; center != nil {
		hashCenter = center.(string)
	}
	fmt.Printf("# Hash Center : %s\n", hashCenter)
	fmt.Printf(sstr((maxLen)*2, "-") + "\n")
}

func (x *XConf) Hash() (s string) {
	s = DefaultInvalidHashString
	v, err := x.Latest()
	if err != nil {
		x.cc.LogWarning(fmt.Sprintf("HashString got error:%s return default:%s", err.Error(), s))
		return
	}
	return x.HashStructure(v)
}

func (x *XConf) HashStructure(v interface{}) (s string) {
	s = DefaultInvalidHashString
	hashCode, err := hashstructure.Hash(v, hashstructure.FormatV2, &hashstructure.HashOptions{TagName: x.cc.TagName})
	if err != nil {
		x.cc.LogWarning(fmt.Sprintf("HashString got error:%s return default:%s", err.Error(), s))
		return
	}
	return fmt.Sprintf("%s%d", HashPrefix, hashCode)
}