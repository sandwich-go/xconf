# XCONF

Golang配置文件加载解析, [goconf](https://github.com/sandwich-go/xconf) v2，扩充了功能支持。
[![Build Status](https://travis-ci.org/sandwich-go/xconf.svg?branch=master)](https://travis-ci.org/sandwich-go/xconf)
[![GoDoc](https://godoc.org/github.com/sandwich-go/xconf?status.svg)](https://godoc.org/github.com/sandwich-go/xconf)
[![Go Report Card](https://goreportcard.com/badge/github.com/sandwich-go/xconf)](https://goreportcard.com/report/github.com/sandwich-go/xconf)[![Sourcegraph](https://sourcegraph.com/github.com/sandwich-go/xconf/-/badge.svg)](https://sourcegraph.com/github.com/sandwich-go/xconf?badge)

## 功能简介
- 支持默认值配置、解析
- 支持多种格式，内置JSON, TOML, YAML，FLAG, ENV支持，并可注册解码器扩展格式支持
- 支持多文件、多`io.Reader`数据加载，支持文件继承
- 支持由OS ENV变量数据加载配置
- 支持由命令行参数FLAGS加载数据
- 支持由远程URL加载配置数据
- 支持数据覆盖合并，加载多份数据时将按照加载文件的顺序按`FieldPath`自动合并
- 支持通过`${READ_TIMEOUT|5s}`、`${IP_ADDRESS}`等方式绑定Env参数
- 支持配置热加载、实时同步，内置ETCD、文件系统、内存热加载支持，支持异步更新通知
- 支持`WATCH`具体的`FieldPath`变动
- 支持导出配置到多种配置文件
- 支持配置HASH，便于比对配置一致性
- `FLAGS`、`ENV`、`FieldPath`支持复杂类型，支持自定义复杂类型扩展支持
- 支持配置访问秘钥
- 支持自定义基于Label的灰度更新


## 名词解释
- `FieldTag`
    - `xconf`在将配置由Strut与JSON, TOML, YAML，FLAG, ENV等转换时使用的字段别名,如:示例配置中的`HttpAddress`的`FieldTag`为`http_address`
    - 如果没有配置`xconf:"http_address"`,则默认会采用字段名的`SnakeCase`作为`FieldTag`，可以通过`xconf.WithFieldTagConvertor`指定为其他方案，如字段名的小写等，注意`FieldTag`策略必须与配置源中使用的字符串一致，否则会导致解析数据失败。
- `FieldPath`,由`FieldTag`组成的Field访问路径，如示例配置中的`Config.SubTest.HTTPAddress`的`FieldPath`为`config.sub_test.http_address`.
- `Leaf`,`xconf`中配置的最小单位，基础类型、slice类型都是最小单位，Struct不是配置的最小单位，会根据配置的属性字段进行赋值、覆盖。
    - 默认情况下map是`xconf`配置的最小单位，但是可以通过指定`notleaf`标签使map不作为最小单位，而是基于key进行合并.但是这种情况下map的Value依然是`xconf`中的最小单位，即使value是Struct也将作为配置合并的最小单位
	- 通过`xconf.WithMapMerge(true)`可以激活`MapMerge`模式，在这个模式下map及其Value都不再是配置的最小单位，配置的最小单位为基础类型和slice类型。


## 快速开始
### 定义配置结构
- 参考[xconf/tests/conf.go](https://github.com/sandwich-go/xconf/blob/master/tests/conf.go)使用[optiongen](https://github.com/timestee/optiongen)定义配置并指定`--xconf=true`以生成支持`xconf`需求的标签.
- 自定义结构,指定`xconf`需求的标签
```golang
type Server struct {
	Timeouts map[string]time.Duration `xconf:"timeouts"`
}

type SubTest struct {
	HTTPAddress string            `xconf:"http_address"`
	MapNotLeaf  map[string]int    `xconf:"map_not_leaf,notleaf"`
	Map2        map[string]int    `xconf:"map2"`
	Map3        map[string]int    `xconf:"map3"`
	Slice2      []int64           `xconf:"slice2"`
	Servers     map[string]Server `xconf:"servers,notleaf"`
}

type Config struct {
	HttpAddress     string          `xconf:"http_address"`
	Map1            map[string]int  `xconf:"map1"`
	MapNotLeaf      map[string]int  `xconf:"map_not_leaf,notleaf"`
	TimeDurations   []time.Duration `xconf:"time_durations"`
	Int64Slice      []int64         `xconf:"int64_slice"`
	Float64Slice    []float64       `xconf:"float64_slice"`
	Uin64Slice      []uint64        `xconf:"uin64_slice"`
	StringSlice     []string        `xconf:"string_slice"`
	ReadTimeout     time.Duration   `xconf:"read_timeout"`
	SubTest         SubTest         `xconf:"sub_test"`
}
```

### 从文件载入配置
以`yaml`格式为例(tests/)
```yaml
http_address: :3002
read_timeout: 100s
default_empty_map:
  test1: 1
map1:
  test1: 1000000
map_not_leaf:
  test1: 1000000
int64_slice:
- 1
- 2
sub_test:
  map2:
    ${IP_ADDRESS}: 2222
  map_not_leaf:
    test2222: 2222
  servers:
    s1:
      timeouts:
        read: ${READ_TIMEOUT|5s} 
```
> 参考:[tests/main/main.go](https://github.com/sandwich-go/xconf/blob/master/tests/main/main.go)，文件间的继承通过`xconf_inherit_files`指定,参考[tests/main/c2.toml](https://github.com/sandwich-go/xconf/blob/master/tests/main/c2.toml)
```golang
cc := NewTestConfig(
	xconf.WithFiles("c2.toml"), // 由指定的文件加载配置
	xconf.WithReaders(bytes.NewBuffer(yamlContents),bytes.NewBuffer(tomlContents),xconf.NewRemoteReader("http://127.0.0.1:9001/test.json", time.Duration(5)*time.Second)), // 由指定的reader加载配置
	xconf.WithFlagSet(flag.CommandLine), // 指定解析flag.CommandLine，默认值
	xconf.WithEnviron(os.Environ()), // 指定解析os.Environ()，默认值
)
xconf.Parse(cc)
```

## 可用选项
- `WithFiles` : 指定加载的文件，配置覆盖顺序依赖传入的文件顺序
- `WithReaders`: 指定加载的`io.Reader`，配置覆盖顺序依赖传入的`io.Reader`顺序。
- `WithFlagSet`: 指定解析的`FlagSet`,默认是全局的`flag.CommandLine`,如指定为`nil`则不会将参数自动创建到`FlagSet`，同样也不会解析`FlagSet`内数据。
- `WithFlagArgs`: 指定`FlagSet`解析的参数数据，默认为`os.Args[1:]`
- `WithFlagValueProvider`: `FlagSet`支持的类型有限，`xconf/xflag/vars`中扩展了部分类型，参考[Flag 与 Env支持]
- `WithEnviron`: 指定环境变量值
- `WithErrorHandling`:指定错误处理方式，同`flag.CommandLine`处理方式
- `WithLogDebug`: 指定debug日志输出
- `WithLogWarning`: 指定warn日志输出
- `WithFieldTagConvertor`: 当无法通过TagName获取`FieldTag`时，通过该方法转换，默认SnakeCase.
- `WithTagName`: `FieldTag`字段来源的Tag名，默认`xconf`
- `WithTagNameDefaultValue`: 默认值使用的Tag名称 ，默认`default`
- `WithParseDefault`:是否解析默认值，默认true，推荐使用[optiongen](https://github.com/timestee/optiongen)生成默认配置数据
- `WithDebug`: 调试模式，会输出详细的解析流程日志
- `WithDecoderConfigOption`: 调整mapstructure参数，`xconf`使用[mapstructure](https://github.com/mitchellh/mapstructure)进行类型转换
- `FieldPathDeprecated`: 弃用的配置，解析时不会报错，但会打印warning日志
- `ErrEnvBindNotExistWithoutDefault`: EnvBind时如果Env中不存在指定的key而且没有指定默认值时报错
- `FieldFlagSetCreateIgnore`: 指定的`FieldPath`或者类型名在没有Flag Provider的时候，不打印报警日志

## Flag 与 Env支持
- 支持Flag中通过`xconf_files`指定配置文件
- `xconf/xflag/vars`中扩展了部分类型如下:
	- float32,float64
	- int,int8,int16,int32,int64
	- uint,uint8,uint16,uint32,uint64
	- []float32,[]float64
	- []int,[]int8,[]int16,[]int32,[]int64
	- []uint,[]uint8,[]uint16,[]uint32,[]uint64
	- []string
	- []Duration
	- map[stirng][string],map[int][int],map[int64][int64],map[int64][string],map[stirng][int],map[stirng][int64],map[stirng][Duration]
- 扩展类型Slice与Map配置
   - slcie的定义方式为元素通过`vars.StringValueDelim`分割，默认为`,`，如:`--time_durations=5s,10s,100s`
   - map的定位方式为K、V通过`vars.StringValueDelim`分割，默认为`,`,如:`--sub_test.map_not_leaf=k1,1,k2,2,k3,3`
- 自定义扩展
	- 扩展需要实现`flag.Getter`接口,可以通过实现`Usage() string`实现自定义的Usage信息。
		```golang
		const JsnoPrefix = "json@"

		type serverProvider struct {
			s    string
			set  bool
			data *map[string]Server
		}

		func (sp *serverProvider) String() string {
			return sp.s
		}
		func (sp *serverProvider) Set(s string) error {
			sp.s = s
			if sp.set == false {
				*sp.data = make(map[string]Server)
			}
			if !strings.HasPrefix(s, JsnoPrefix) {
				return errors.New("server map need json data with prefix:" + JsnoPrefix)
			}
			s = strings.TrimPrefix(s, JsnoPrefix)
			return json.Unmarshal([]byte(s), sp.data)
		}
		func (sp *serverProvider) Get() interface{} {
			ret := make(map[string]interface{})
			for k, v := range *sp.data {
				ret[k] = v
			}
			return ret
		}
		func (sp *serverProvider) Usage() string {
			return fmt.Sprintf("server map, json format")
		}
		func newServerProvider(v interface{}) flag.Getter {
			return &serverProvider{data: v.(*map[string]Server)}
		}

		```
	- 注册扩展
		- `vars.SetProviderByFieldPath`通过`FieldPath`设定扩展
		- `vars.SetProviderByFieldType`通过字段类型名称设定扩展
	```golang
		cc := &Config{}
		jsonServer := `json@{"s1":{"timeouts":{"read":5000000000},"timeouts_not_leaf":{"write":5000000000}}}`
		x := xconf.New(
			xconf.WithFlagSet(flag.NewFlagSet("xconf-test", flag.ContinueOnError)),
			xconf.WithFlagArgs("--sub_test.servers="+jsonServer), //数据设定到flag中
			xconf.WithEnviron("sub_test_servers="+jsonServer), //数据设定到env中
		)
		vars.SetProviderByFieldPath("sub_test.servers", newServerProvider) // 根据字段FieldPath设定Provider
		vars.SetProviderByFieldType("map[string]Server", newServerProvider) // 根据字段类型名设定Provider
	```
- Keys
 通过`xconf.DumpInfo`获取配置支持的FLAG与ENV名称,如下所示.
 ```shell
	------------------------------------------------------
	#    FLAG                        ENV                        
	------------------------------------------------------
	1    default_empty_map           DEFAULT_EMPTY_MAP          
	2    float64_slice               FLOAT64_SLICE              
	3    http_address                HTTP_ADDRESS               
	4    int64_slice                 INT64_SLICE                
	5    map1                        MAP1                       
	6    map_not_leaf                MAP_NOT_LEAF               
	7    read_timeout                READ_TIMEOUT               
	8    string_slice                STRING_SLICE               
	9    sub_test.http_address       SUB_TEST_HTTP_ADDRESS      
	10   sub_test.map2               SUB_TEST_MAP2              
	11   sub_test.map3               SUB_TEST_MAP3              
	12   sub_test.map_not_leaf       SUB_TEST_MAP_NOT_LEAF      
	13   sub_test.slice2             SUB_TEST_SLICE2            
	14   time_durations              TIME_DURATIONS             
	15   uin64_slice                 UIN64_SLICE     
 ```
 
### ENV绑定
支持解析ENV变量名称,如下例：
```golang

var yamlTest2 = []byte(`
http_address: :3002
read_timeout: 100s
default_empty_map:
  test1: 1
map1:
  test1: 1000000
map_not_leaf:
  test1: 1000000
int64_slice:
- 1
- 2
sub_test:
  map2:
    ${IP_ADDRESS}: 2222
  map_not_leaf:
    test2222: 2222
  servers:
    s1:
      timeouts:
        read: ${READ_TIMEOUT|5s} 
`)

func TestEnvBind(t *testing.T) {
	Convey("env bind", t, func(c C) {
		cc := &Config{}
		x := xconf.NewWithoutFlagEnv()
		x.UpdateWithFieldPathValues("http_address", "${XCONF_HOST}:${XCONF_PORT}")
		err := x.Parse(cc)
		So(err, ShouldBeNil)
		So(cc.HttpAddress, ShouldEqual, "")
		host := "127.0.0.1"
		port := "9001"
		os.Setenv("XCONF_HOST", host)
		os.Setenv("XCONF_PORT", port)
		So(cc.HttpAddress, ShouldEqual, "")
		So(x.UpdateWithReader(bytes.NewBuffer(yamlTest2)), ShouldBeNil)
		So(x.UpdateWithFieldPathValues("http_address", "${XCONF_HOST}:${XCONF_PORT}"), ShouldBeNil)
		latest, err := x.Latest()
		So(err, ShouldBeNil)
		cc = latest.(*Config)
		So(cc.HttpAddress, ShouldEqual, host+":"+port)
		So(cc.SubTest.Servers["s1"].Timeouts["read"], ShouldEqual, time.Duration(5)*time.Second)
	})
}
```

### URL读取
```golang
cc := &Config{}
x := xconf.NewWithoutFlagEnv(xconf.WithReaders(xconf.NewRemoteReader("http://127.0.0.1:9001/test.yaml", time.Duration(1)*time.Second)))
```

## 动态更新

### 基于配置文件
```golang
	testBytesInMem := "memory_test_key"
	mem, err := xmem.New()
	panicErr(err)
	// xconf/kv提供了基于ETCD/FILE/MEMORY的更新机制
	// 可自行实现xconf的Loader接口或者对接到xmem,借助xmem的机制实现配置更新
	xconf.WatchUpdate(testBytesInMem, mem)
	updated := make(chan *Config, 1)
	go func() {
		for {
			select {
			case v := <-x.NotifyUpdate():
				updated <- v.(*Config)
			}
		}
	}()
```

### 基于`FieldPath`
```golang
err := xconf.WatchFieldPath("sub_test.http_address", func(from, to interface{}) {
	fmt.Printf("sub_test.http_address changed from %v to %v ", from, to)
})
panicErr(err)
```
可以通过以下方法实现基于文件或Buffer的更新，更新结果通过`xconf.NotifyUpdate`异步获取,或通过`xconf.Latest`同步获取。
- `UpdateWithFiles(files ...string) (err error)`
- `UpdateWithReader(readers ...io.Reader) (err error)`

可以通过以下方法实现基于`FieldPath`的配置更新，更新结果通过`xconf.NotifyUpdate`异步获取,或通过`xconf.Latest`同步获取。
- `UpdateWithFlagArgs(flagArgs ...string)  (err error)`
- `UpdateWithEnviron(environ ...string) (err error)`
- `UpdateWithFieldPathValues(kv ...string) (err error)`

### 绑定最新配置
```golang
xconf.Latest()
```

### Atomic自动更新
使用[optiongen](https://github.com/timestee/optiongen)定义配置并指定`--xconf=true`生成支持`XConf`的配置会默认生成`Atomic`更新支持：
```golang

func (cc *Config) AtomicSetFunc() func(interface{}) { return AtomicConfigSet }

var atomicConfig unsafe.Pointer

func AtomicConfigSet(update interface{}) {
	atomic.StorePointer(&atomicConfig, (unsafe.Pointer)(update.(*Config)))
}
func AtomicConfig() *Config {
	current := (*Config)(atomic.LoadPointer(&atomicConfig))
	if current == nil {
		atomic.CompareAndSwapPointer(&atomicConfig, nil, (unsafe.Pointer)(newDefaultConfig()))
		return (*Config)(atomic.LoadPointer(&atomicConfig))
	}
	return current
}

```

解析时提供`AtomicConfig()`即可，当配置更新的时候回自动回调`AtomicConfigSet`方法进行指针替换。
```golang
func TestAtomicVal(t *testing.T) {
	Convey("atomic val", t, func(c C) {
		x := xconf.NewWithoutFlagEnv()
		So(x.Parse(AtomicConfig()), ShouldBeNil)
		So(x.UpdateWithFieldPathValues("http_address", "10.10.10.10"), ShouldBeNil)
		So(AtomicConfig().HttpAddress, ShouldEqual, "10.10.10.10")
	})
}

```

