package xconf

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"reflect"

	"github.com/mitchellh/hashstructure/v2"
	"github.com/sandwich-go/mapstructure"
	"github.com/sandwich-go/xconf/xflag"
	"github.com/sandwich-go/xconf/xutil"
)

// XConf XConf struct
type XConf struct {
	zeroValPtrForLayout interface{}                    //只是用户获取结构信息，不存储数据信息
	cc                  *Options                       // 配置参数
	fieldPathInfoMap    map[string]StructFieldPathInfo // FieldPath到字段属性映射
	dataLatestCached    map[string]interface{}         // 缓存的最新数据
	dataMeta            map[string]interface{}         // todo meta数值,如有强烈的访问需求，meta值可以作为一份隐含的配置直接解析到对应的的struct
	dynamicUpdate       sync.Mutex
	kvs                 []*kvLoader
	updated             chan interface{}
	runningLogger       func(string)
	hasParsed           bool
	mapOnFieldUpdated   map[string]OnFieldUpdated
	changes             fieldChanges
	atomicSetFunc       func(interface{})
	optionUsage         string
	parseForMerge       bool
	valPtrForUsageDump  interface{}
}

// New 构造新的Xconf
func New(opts ...Option) *XConf { return NewWithConf(NewOptions(opts...)) }

// NewWithoutFlagEnv 构造新的Xconf,移除FlagSet和Environ解析
func NewWithoutFlagEnv(opts ...Option) *XConf {
	return New(append(opts, WithFlagSet(nil), WithEnviron())...)
}

// NewWithConf 由指定的配置构造XConf
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

// ErrorHandling 错误处理类型
type ErrorHandling int

const (
	// ContinueOnError 发生错误继续运行，Parse会返回错误
	ContinueOnError ErrorHandling = iota
	// ExitOnError 发生错误后退出
	ExitOnError
	// PanicOnError 发生错误后主动panic
	PanicOnError
)

// Latest 绑定当前XConf内缓存的数据到Parse时传入的类型中并以interface{}类型返回，需先调用Parse便于XConf确定配置类型
func (x *XConf) Latest() (interface{}, error) {
	if !x.hasParsed {
		return nil, errors.New("need Parse first")
	}
	zeroOne := reflect.New(reflect.ValueOf(x.zeroValPtrForLayout).Type().Elem()).Interface()
	return zeroOne, x.bindLatest(zeroOne)
}

// Copy 返回当前XConf的拷贝
func (x *XConf) Copy() *XConf { return NewWithConf(x.cc) }

// clean 不允许外部clean清空状态，配置的增量更新依赖于解析后的状态
func (x *XConf) clean() {
	x.fieldPathInfoMap = make(map[string]StructFieldPathInfo)
	x.dataLatestCached = make(map[string]interface{})
	x.changes.changed = make(map[string]*fieldValues)
	x.dataMeta = make(map[string]interface{})
}

// NotifyUpdate 通知更新
func (x *XConf) NotifyUpdate() <-chan interface{} { return x.updated }

func (x *XConf) defaultXFlagOptions() []xflag.Option {
	return []xflag.Option{
		xflag.WithFlatten(false),
		xflag.WithFlagSet(x.cc.FlagSet),
		xflag.WithKeyFormat(x.cc.FieldTagConvertor),
		xflag.WithTagName(x.cc.TagName),
		xflag.WithLogDebug(xflag.LogFunc(x.cc.LogDebug)),
		xflag.WithLogWarning(xflag.LogFunc(x.cc.LogWarning)),
		xflag.WithFlagCreateIgnoreFiledPath(x.cc.FlagCreateIgnoreFiledPath...),
		xflag.WithStringAlias(func(s string) string {
			if v, ok := x.cc.StringAlias[s]; ok {
				return v
			}
			if v, ok := x.cc.StringAliasFunc[s]; ok {
				return v(s)
			}
			return s
		}),
	}
}

// ZeroStructKeysTagList 获取参数s的空结构的Filed信息
func (x *XConf) ZeroStructKeysTagList(s interface{}) map[string]StructFieldPathInfo {
	_, v := NewStruct(
		reflect.New(reflect.ValueOf(s).Type().Elem()).Interface(),
		x.cc.TagName,
		x.cc.TagNameForDefaultValue,
		x.cc.FieldTagConvertor,
	).Map()
	return v
}

// StructMapStructure 获取传入的s的数据的map[string]interface{}
func (x *XConf) StructMapStructure(s interface{}) map[string]interface{} {
	v, _ := NewStruct(s, x.cc.TagName, x.cc.TagNameForDefaultValue, x.cc.FieldTagConvertor).Map()
	return v
}

func (x *XConf) runningLogData(name string, data map[string]interface{}) {
	x.runningLogger(fmt.Sprintf("===========================> Data %s\n%v\n", name, data))
}

// FieldMap 获取对象的Field对象Map
func FieldMap(valPtr interface{}, x *XConf) map[string]StructFieldPathInfo {
	// 如果应用层配置实现了XConfOptions
	if w, ok := valPtr.(XConfOptions); ok {
		x.runningLogger("apply config XConfOptions")
		x.cc.ApplyOption(w.XConfOptions()...)
	}
	// 获取bindto结构合法的FieldPath，并过滤合法的BindToFieldPath
	_, fieldsMap := NewStruct(
		reflect.New(reflect.ValueOf(valPtr).Type().Elem()).Interface(),
		x.cc.TagName,
		x.cc.TagNameForDefaultValue,
		x.cc.FieldTagConvertor,
	).Map()
	return fieldsMap
}

// FieldPathList 获取对象的FieldPath列表
func FieldPathList(valPtr interface{}, x *XConf) (ret []string) {
	for k := range FieldMap(valPtr, x) {
		ret = append(ret, k)
	}
	return ret
}

func (x *XConf) keysList() []string {
	var keys []string
	for k := range x.fieldPathInfoMap {
		keys = append(keys, k)
	}
	return keys
}

// mergeToDest 将指定的数据并入XConf缓存的最终数据字段dataLatestCached
func (x *XConf) mergeToDest(dataName string, data map[string]interface{}) error {
	x.runningLogData(dataName, data)
	x.runningLogger(fmt.Sprintf("----> merge src:%s dst:%s\n", dataName, "dest"))

	// 灰度发布初步支持
	grayLabelVal, ok := data[MetaKeyGrayLabel]
	if ok {
		if grayLabelStr, ok := grayLabelVal.(string); ok {
			grayLabelList := xutil.ToCleanStringSlice(grayLabelStr)
			if len(grayLabelList) > 0 && !xutil.ContainAtLeastOneEqualFold(grayLabelList, x.cc.AppLabelList) {
				x.cc.LogDebug(fmt.Sprintf("do not apply to local instance due to %v and %v", grayLabelList, x.cc.AppLabelList))
				return nil
			}
		}
	}
	// 剔除meta keys指定的数据,合并到dest的数据不需要包含meta值
	for _, metaKey := range metaKeyList {
		if v, ok := data[metaKey]; ok {
			x.dataMeta[metaKey] = v
		}
		delete(data, metaKey)
	}

	err := mergeMap("", 0, x.runningLogger, data, x.dataLatestCached, x.isLeafFieldPath, nil, &x.changes)
	return xutil.WrapIfErr(err, "got error:%w while merge data:%s to data: dst", err, dataName)
}

func (x *XConf) mergeMap(srcName string, dstName string, src map[string]interface{}, dst map[string]interface{}) error {
	x.runningLogger(fmt.Sprintf("----> merge src:%s dst:%s\n", srcName, dstName))
	return mergeMap("", 0, x.runningLogger, src, dst, x.isLeafFieldPath, nil, nil)
}

func (x *XConf) getOptionUsage(valPtr interface{}) (ret string) {
	if w, ok := valPtr.(GetOptionUsage); ok {
		ret = xutil.StringTrim(w.GetOptionUsage())
	}
	if ret == "" {
		return x.cc.OptionUsagePoweredBy
	}
	if x.cc.OptionUsagePoweredBy == "" {
		return ret
	}
	return ret + "\n" + x.cc.OptionUsagePoweredBy
}

// Merge 合并配置
func (x *XConf) Merge(opts ...Option) error {
	x.parseForMerge = true
	opts = append(opts, WithFlagSet(nil), WithEnviron())
	return x.Parse(nil, opts...)
}

func (x *XConf) parse(valPtr interface{}) (err error) {
	if x.cc.LogDebug == nil {
		x.cc.LogDebug = func(s string) {}
	}
	if x.cc.LogWarning == nil {
		x.cc.LogWarning = func(s string) {}
	}
	defer func() {
		if reason := recover(); reason != nil {
			err = fmt.Errorf("XConf parse got panic : %v", reason)
		}
	}()
	x.hasParsed = true

	if !x.parseForMerge {
		// 检测传入的数值是否为指针
		if reflect.ValueOf(valPtr).Kind() != reflect.Ptr {
			return errors.New("unsupported type, pass in as ptr")
		}
		x.optionUsage = x.getOptionUsage(valPtr)

		// 如果应用层配置实现了XConfOptions
		applyXConfOptions(valPtr, x)
	}

	if x.cc.FlagSet != nil && x.cc.ReplaceFlagSetUsage {
		x.valPtrForUsageDump = valPtr
		x.cc.FlagSet.Usage = x.Usage
	}

	if !x.parseForMerge {
		x.runningLogger(fmt.Sprintf("=========> Parse With\n%v\n", valPtr))
		// 保留结构信息
		x.zeroValPtrForLayout = reflect.New(reflect.ValueOf(valPtr).Type().Elem()).Interface()
		// 获取结构信息，后续的数据更新等依赖于该结构信息
		x.fieldPathInfoMap = x.ZeroStructKeysTagList(x.zeroValPtrForLayout)
		if reflect.DeepEqual(x.zeroValPtrForLayout, valPtr) && x.cc.ParseDefault {
			// 如果指定根据tag解析默认数值则进行一次解析操作,将解析得到的数值作为默认值
			// 如果input为空值，则不解析struct本身数值，struct中解析得到的是全量key-val的mapstructure，防止覆盖default
			xutil.PanicErr(x.updateDstDataWithParseDefault(valPtr))
		} else {
			// 如果input不为空，则进行解析，input值完全覆盖default内的值,不再解析default
			xutil.PanicErr(x.mergeToDest("struct_input", x.StructMapStructure(valPtr)))
		}
	}

	//flag指定的文件 直接覆盖 内部指定的文件, 独立解析flagset数据以获取files
	flagData, filesToParse, err := x.parseFlagFilesForXConf(valPtr, x.cc.FlagSet, x.cc.FlagArgs...)
	filesToParse = xutil.StringSliceWalk(filesToParse, xutil.StringSliceEmptyFilter)
	xutil.PanicErr(err)
	xutil.PanicErr(x.updateDstDataWithFiles(filesToParse...))
	xutil.PanicErr(x.updateDstDataWithReaders(x.cc.Readers...))
	xutil.PanicErr(x.commonUpdateDstData("flag", func() (map[string]interface{}, error) { return flagData, nil }))
	xutil.PanicErr(x.updateDstDataWithEnviron(x.cc.Environ...))
	if !x.parseForMerge {
		xutil.PanicErr(x.bindLatest(valPtr))
	}
	if w, ok := valPtr.(AtomicSetterProvider); ok {
		x.atomicSetFunc = w.AtomicSetFunc()
		x.atomicSetFunc(valPtr)
	}
	return nil
}

func (x *XConf) parseFlagFilesForXConf(valPtr interface{}, flagSet *flag.FlagSet, args ...string) (flagData map[string]interface{}, filesToParse []string, err error) {
	filesToParse = x.cc.Files
	if x.cc.FlagSet == nil {
		return
	}
	validKeys := x.keysList()
	if x.cc.ParseMetaKeyFlagFiles {
		if fv := x.cc.FlagSet.Lookup(MetaKeyFlagFiles); fv == nil {
			x.cc.FlagSet.String(MetaKeyFlagFiles, "", "xconf files provided by flag, file slice, split by ,")
		}
		validKeys = append(validKeys, MetaKeyFlagFiles)
	}
	flagData, err = xflagMapstructure(
		valPtr,
		validKeys,
		func(*xflag.Maker) []string { return x.cc.FlagArgs },
		append(x.defaultXFlagOptions(), xflag.WithFlagSet(x.cc.FlagSet))...)
	if err != nil {
		return
	}
	if x.cc.ParseMetaKeyFlagFiles {
		if v := flagData[MetaKeyFlagFiles]; v != nil {
			filesToParse = strings.Split(strings.Trim(v.(string), " "), ",")
			x.cc.LogDebug(fmt.Sprintf("config files changed from:%v to %v provided by FlagSet", x.cc.Files, filesToParse))
		}
	}

	delete(flagData, MetaKeyFlagFiles)
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
	xutil.PanicErrWithWrap(err, "got error:%w while decode using map structure", err)
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
	xutil.PanicErrWithWrap(err, "got error:%w while parse default value", err)
	if !defaultParsed {
		return
	}
	x.cc.LogDebug(fmt.Sprintf("Parse Default From Tag:%s", x.cc.TagNameForDefaultValue))
	xutil.PanicErr(x.mergeToDest("default_from_tag", dataDefault))
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
	if err != nil {
		return
	}
	if data == nil {
		return nil
	}
	xutil.PanicErrWithWrap(err, "got error while load %s,err :%w ", name, err)
	xutil.PanicErr(x.mergeToDest(name, data))
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
	// 根据前缀过滤
	if x.cc.EnvironPrefix != "" {
		environ = xutil.StringSliceWalk(environ, func(s string) (string, bool) {
			if strings.HasPrefix(strings.ToUpper(s), strings.ToUpper(x.cc.EnvironPrefix)) {
				return s, true
			}
			return s, false
		})
	}
	return x.commonUpdateDstData("env", func() (map[string]interface{}, error) {
		return xflagMapstructure(
			x.zeroValPtrForLayout,
			x.keysList(),
			func(xf *xflag.Maker) []string {
				return envBindToFlags(environ, xf.EnvKeysMapping(x.cc.EnvironPrefix, x.keysList()))
			},
			append(x.defaultXFlagOptions(), xflag.WithFlagSet(newFlagSetContinueOnError("Environ")))...)
	})
}

func (x *XConf) decode(data map[string]interface{}, valPtr interface{}) error {
	config := x.defaultDecoderConfig(valPtr)
	config.TagName = x.cc.TagName
	// config.Squash = x.cc.EnableSquash
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
			if xutil.ContainString(metaKeyList, v) {
				continue
			}
			// 逻辑层指定的移除的字段，报警
			if xutil.ContainString(x.cc.FieldPathRemoved, v) {
				deprecated = append(deprecated, v)
				continue
			}
			unused = append(unused, v)
		}
		if len(deprecated) != 0 {
			x.cc.LogWarning(fmt.Sprintf("!!! DEPRECATED FIELD, WILL REMOVE IN FUTURE. FIELDS: %s", strings.Join(deprecated, ",")))
		}
		if x.cc.ErrorUnused && len(unused) != 0 {
			return fmt.Errorf("!!! UNUSED FIELDS, SHOULD PAY ATTENTION. FIELDS: %s", strings.Join(unused, ","))
		}
	}
	return nil
}

// MustParse 解析配置到传入的参数中,如发生错误则直接panic
func (x *XConf) MustParse(valPtr interface{}, opts ...Option) {
	_ = x.Parse(valPtr, append(opts, WithErrorHandling(PanicOnError))...)
}

// Parse 解析配置到传入的参数中
func (x *XConf) Parse(valPtr interface{}, opts ...Option) error {
	x.cc.ApplyOption(opts...)
	err := x.parse(valPtr)
	if err == nil || IsErrHelp(err) {
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

// Hash 当前最新配置的Hash字符串，默认为DefaultInvalidHashString
func (x *XConf) Hash() (s string) {
	s = DefaultInvalidHashString
	v, err := x.Latest()
	if err != nil {
		x.cc.LogWarning(fmt.Sprintf("HashString got error:%s return default:%s", err.Error(), s))
		return
	}
	return x.HashStructure(v)
}

// HashStructure 返回指定数据的hash值
func (x *XConf) HashStructure(v interface{}) (s string) {
	s = DefaultInvalidHashString
	hashCode, err := hashstructure.Hash(v, hashstructure.FormatV2, &hashstructure.HashOptions{TagName: x.cc.TagName})
	if err != nil {
		x.cc.LogWarning(fmt.Sprintf("HashString got error:%s return default:%s", err.Error(), s))
		return
	}
	return fmt.Sprintf("%s%d", HashPrefix, hashCode)
}
