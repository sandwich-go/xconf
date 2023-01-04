# XCONF

[![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/sandwich-go/xconf/ci.yml?style=flat-square)](https://github.com/sandwich-go/xconf/actions?query=workflow%3ACI)
![Go Version](https://img.shields.io/badge/go%20version-%3E=1.14-61CFDD.svg?style=flat-square)
[![Go Report Card](https://goreportcard.com/badge/github.com/sandwich-go/xconf)](https://goreportcard.com/report/github.com/sandwich-go/xconf)
[![GoDoc](https://godoc.org/github.com/sandwich-go/xconf?status.svg)](https://godoc.org/github.com/sandwich-go/xconf)

[README | 中文](README.md)

Golang configuration file loading parsing, [goconf](https://github.com//timestee/goconf) v2, extended feature support

Run XConf Example: [![run on repl.it](https://repl.it/badge/github/timestee/XConf-example)](https://repl.it//@timestee/XConf-example#main.go)

Run XCmd Example: [![run on repl.it](https://replit.com/badge/github/timestee/XCmd-example)](https://replit.com/@timestee/XCmd-example#main.go)

## Function Introduction
- Support default value configuration, parsing
- Support multiple formats, built-in JSON, TOML, YAML, FLAG, ENV support, and can register decoder to extend format support
- Supports multi-file, multi-`io.Reader` data loading, file inheritance support
- Support data loading configuration by OS ENV variables
- Support loading data by command line parameter FLAGS
- Support loading configuration data by remote URL
- Support data overwrite merge, when loading multiple copies of data will be automatically merged by `FieldPath` in the order of the loaded files
- Support binding Env parameters by `${READ_TIMEOUT|5s}`, `${IP_ADDRESS}`, etc.
- Support configuration hotload, real-time synchronization, built-in memory hotload support, support asynchronous update notification, support [xconf-providers](https://github.com/sandwich-go/xconf-providers): ETCD, file system.
- Support `WATCH` specific `FieldPath` changes
- Support export configuration to multiple configuration files
- Support configuration HASH, easy to compare configuration consistency
- `FLAGS`, `ENV`, `FieldPath` support complex types, support for custom complex type extension support
- Support configuration of access secret key
- Support custom grayscale update based on Label
- Support numeric aliases, such as `math.MaxInt`,`runtime.NumCPU`
- Support ",squash" to mention the fields of sub-structures to the parent structure for configuration expansion

## Explanation of terms
- `FieldTag`
    - `xconf` is the field alias used when converting the configuration from Strut to JSON, TOML, YAML, FLAG, ENV, etc. For example: the `FieldTag` of `HttpAddress` in the example configuration is `http_address`.
    - If `xconf: "http_address"` is not configured, the default field name of `SnakeCase` will be used as the `FieldTag`, which can be specified by `xconf.WithFieldTagConvertor` for other programs, such as lowercase field names, etc. Note that the `FieldTag` strategy must be be consistent with the string used in the configuration source, otherwise it will cause the parsing data to fail.
- `FieldPath`, the Field access path composed by `FieldTag`, such as the `Config.SubTest.HTTPAddress` in the sample configuration `FieldPath` for `config.sub_test.http_address`.
- `Leaf`,`xconf` in the configuration of the minimum unit, the base type, slice type are the minimum unit, Struct is not the minimum unit of the configuration, will be assigned, override according to the configuration of the property field.
    - By default map is the minimum unit configured in `xconf`, but you can specify the `notleaf` tag so that map is not the minimum unit, but is merged based on key. But in this case the value of map is still the minimum unit in `xconf`, even if the value is a Struct, it will be the minimum unit for configuration merging
	- With `xconf.WithMapMerge(true)` you can activate the `MapMerge` mode, in which both map and its value are no longer the minimum unit of the configuration, but the minimum unit of the configuration is the base type and the slice type.

## Quick start
### Define the configuration structure
- Refer to [xconf/tests/conf.go](https://github.com/sandwich-go/xconf/blob/master/tests/conf.go) to use [optiongen](https://github.com/timestee/ optiongen) to define the configuration and specify `--xconf=true` to generate tags that support `xconf` requirements.
- Customize the structure to specify `xconf`-required tags
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

### Load configuration from file
Take the `yaml` format as an example (tests/)
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
> Reference:[test/main/main.go](https://github.com/sandwich-go/xconf/blob/master/tests/main/main.go), inheritance between files is specified by `xconf_inherit_files`, reference [ test/main/c2.toml](https://github.com/sandwich-go/xconf/blob/master/tests/main/c2.toml)
```golang
cc := NewTestConfig(
	xconf.WithFiles("c2.toml"),
	xconf.WithReaders(bytes.NewBuffer(yamlContents),bytes.NewBuffer(tomlContents),xconf.NewRemoteReader("http://127.0.0.1:9001/test.json", time.Duration(5)*time.Second)),
	xconf.WithFlagSet(flag.CommandLine),
	xconf.WithEnviron(os.Environ()),
)
xconf.Parse(cc)
```

### Configuration deposit file
```golang
// SaveToFile dumps the built-in parsed data to a file, selecting the codec according to the file suffix.
func SaveToFile(fileName string) error
// SaveToWriter dumps the built-in parsed data to the writer, type ct
func SaveToWriter(ct ConfigType, writer io.Writer) error 

// SaveVarToFile writes the external valPtr to the fileName, selecting the codec according to the file suffix.
func SaveVarToFile(valPtr interface{}, fileName string) error 

// SaveVarToWriter writes the external valPtr to the writer, type ct
func SaveVarToWriter(valPtr interface{}, ct ConfigType, writer io.Writer) error 

// MustSaveToFile dump the built-in parsed data to a file, choose the codec according to the file suffix, if an error occurs it will panic
func MustSaveToFile(f string) 
// MustSaveToWriter dumps the built-in parsed data to the writer, specify the ConfigType, if an error occurs, it will panic.
func MustSaveToWriter(ct ConfigType, writer io.) 

// MustSaveVarToFile write external valPtr to fileName, select codec according to file suffix
func MustSaveVarToFile(v interface{}, f string) 

// MustSaveVarToWriter writes the external valPtr to writer, type ct
func MustSaveVarToWriter(v interface{}, ct ConfigType, w io.Writer) 

// MustSaveToBytes returns the built-in parsed data as a byte stream, ConfigType must be specified
func MustSaveToBytes(ct ConfigType) []byte { return xx.MustSaveToBytes(ct) }

// SaveVarToWriterAsYAML parses the built-in parsed data to yaml with comment
func SaveVarToWriterAsYAML(valPtr interface{}, writer io.Writer) error
```

## Available options
- `WithFiles` : specifies the files to be loaded, the configuration override order depends on the incoming file order
- `WithReaders`: specifies the loaded `io.Reader`, the configuration override order depends on the incoming `io.Reader` order.
- CommandLine`, if `nil` is specified, the parameters will not be automatically created to `FlagSet`, and the data in `FlagSet` will not be parsed.
- `WithFlagArgs`: Specify the parameter data to be parsed by `FlagSet`, default is `os.Args[1:]`.
- `WithFlagValueProvider`: `FlagSet` supports limited types, some types are extended in `xconf/xflag/vars`, see [Flag and Env support].
- `WithEnviron`: specifies the value of the environment variable
- `WithErrorHandling`: Specify the error handling method, same as `flag.
- `WithLogDebug`: Specify the debug log output
- `WithLogWarning`: Specify the warn log output
- `WithFieldTagConvertor`: This method converts `FieldTag` when it cannot be obtained by TagName, default SnakeCase.
- `WithTagName`: Tag name of the source of the `FieldTag` field, default `xconf`.
- `WithTagNameDefaultValue`: The Tag name used for the default value, default `default`.
- `WithParseDefault`: whether to parse the default value, default true, recommended to use [optiongen](https://github.com/timestee/optiongen) to generate the default configuration data
- `WithDebug`: debug mode, will output detailed log of parsing process
- `WithDecoderConfigOption`: adjust the mapstructure parameter, `xconf` uses [mapstructure](https://github.com/mitchellh/mapstructure) for type conversion
- `FieldPathDeprecated`: deprecated configuration, no error will be reported when parsing, but a warning log will be printed
- `ErrEnvBindNotExistWithoutDefault`: Error when EnvBind if the specified key does not exist in Env and no default value is specified
- `FieldFlagSetCreateIgnore`: The specified `FieldPath` or type name will not print the warning log when there is no Flag Provider.

## Flag and Env Support
- Support for specifying configuration files in Flag via `xconf_files`
- `xconf/xflag/vars` extends some of the types as follows:
    - float32,float64
    - int,int8,int16,int32,int64
    - uint,uint8,uint16,uint32,uint64
    - []float32,[]float64
    - []int,[]int8,[]int16,[]int32,[]int64
    - []uint,[]uint8,[]uint16,[]uint32,[]uint64
    - []string
    - []Duration
    - map[stirng]string,map[int]int,map[int64]int64,map[int64]string,map[stirng]int,map[stirng]int64,map[stirng]Duration
- Extended type Slice and Map configuration
   - slcie is defined in such a way that elements are split by `vars.StringValueDelim`, the default is `,`, for example:`--time_durations=5s,10s,100s`
   - map is positioned as K,V split by `vars.StringValueDelim`, default is `,`, e.g.:`--sub_test.map_not_leaf=k1,1,k2,2,k3,3`
- Custom extensions
    - The extension needs to implement the `flag.Getter` interface, which can be used to implement custom Usage information by implementing the `Usage() string`.
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
    - Registering extensions
        - `vars.SetProviderByFieldPath` set extension by `FieldPath`
        - `vars.SetProviderByFieldType` sets extensions by field type name
    ```golang
        cc := &Config{}
        jsonServer := `json@{"s1":{"timeouts":{"read":5000000000},"timeouts_not_leaf":{"write":5000000000}}}`
        x := xconf.New(
            xconf.WithFlagSet(flag.NewFlagSet("xconf-test", flag.ContinueOnError)),
            xconf.WithFlagArgs("--sub_test.servers="+jsonServer),
            xconf.WithEnviron("sub_test_servers="+jsonServer),
        )
        vars.SetProviderByFieldPath("sub_test.servers", newServerProvider)
        vars.SetProviderByFieldType("map[string]Server", newServerProvider)
    ```
- Keys
 `xconf.DumpInfo` to get the FLAG and ENV names supported by the configuration, as shown below, where Y is the Option configuration item, D is the Deprecated field, and M is the xconf internal field.
 ```shell
------------------------------------------------------------------------------------------------------------------------------------------------------------------------
FLAG                              ENV                                         TYPE            USAGE
------------------------------------------------------------------------------------------------------------------------------------------------------------------------
--default_empty_map               TEST_PREFIX_DEFAULT_EMPTY_MAP               map[string]int  |Y| xconf/xflag/vars, key and value split by ,
--float64_slice                   TEST_PREFIX_FLOAT64_SLICE                   []float64       |Y| xconf/xflag/vars, value split by , (default [101.191 202.202 303.303])
--http_address                    TEST_PREFIX_HTTP_ADDRESS                    string          |Y| http_address (default "127.0.0.1")
--int64_slice                     TEST_PREFIX_INT64_SLICE                     []int64         |Y| xconf/xflag/vars, value split by , (default [101 202 303])
--int8                            TEST_PREFIX_INT8                            int8            |Y| int8 (default 1)
--map1                            TEST_PREFIX_MAP1                            map[string]int  |Y| k,v使用,分割 (default map[test1:100 test2:200])
--map_not_leaf                    TEST_PREFIX_MAP_NOT_LEAF                    map[string]int  |D| Deprecated: 使用Map1 (default map[test1:100 test2:200])
--max_int                         TEST_PREFIX_MAX_INT                         int             |Y| max_int (default 0)
--max_uint64                      TEST_PREFIX_MAX_UINT64                      uint64          |Y| max_uint64 (default 0)
--option_usage                    TEST_PREFIX_OPTION_USAGE                    string          |Y| option_usage (default "Some application-level configuration rules are described here")
--process_count                   TEST_PREFIX_PROCESS_COUNT                   int8            |Y| process_count (default 1)
--read_timeout                    TEST_PREFIX_READ_TIMEOUT                    Duration        |Y| read_timeout (default 5s)
--redis.redis_address             TEST_PREFIX_REDIS_REDIS_ADDRESS             string          |Y| redis.redis_address (default "127.0.0.1:6637")
--redis_as_pointer.redis_address  TEST_PREFIX_REDIS_AS_POINTER_REDIS_ADDRESS  string          |Y| redis_as_pointer.redis_address
--redis_timeout.read_timeout      TEST_PREFIX_REDIS_TIMEOUT_READ_TIMEOUT      Duration        |Y| redis_timeout.read_timeout (default 0s)
--string_slice                    TEST_PREFIX_STRING_SLICE                    []string        |Y| xconf/xflag/vars, value split by , (default [test1 test2 test3])
--sub_test.http_address           TEST_PREFIX_SUB_TEST_HTTP_ADDRESS           string          |Y| sub_test.http_address
--sub_test.map2                   TEST_PREFIX_SUB_TEST_MAP2                   map[string]int  |Y| xconf/xflag/vars, key and value split by ,
--sub_test.map3                   TEST_PREFIX_SUB_TEST_MAP3                   map[string]int  |Y| xconf/xflag/vars, key and value split by ,
--sub_test.map_not_leaf           TEST_PREFIX_SUB_TEST_MAP_NOT_LEAF           map[string]int  |Y| xconf/xflag/vars, key and value split by ,
--sub_test.slice2                 TEST_PREFIX_SUB_TEST_SLICE2                 []int64         |Y| xconf/xflag/vars, value split by ,
--test_bool                       TEST_PREFIX_TEST_BOOL                       bool            |Y| test_bool (default false)
--test_bool_true                  TEST_PREFIX_TEST_BOOL_TRUE                  bool            |Y| test_bool_true (default true)
--time_durations                  TEST_PREFIX_TIME_DURATIONS                  []Duration      |Y| 延迟队列 (default [1s 1s])
--uin64_slice                     TEST_PREFIX_UIN64_SLICE                     []uint64        |Y| xconf/xflag/vars, value split by , (default [101 202 303])
--xconf_flag_files                TEST_PREFIX_XCONF_FLAG_FILES                string          |M| xconf files provided by flag, file slice, split by ,
------------------------------------------------------------------------------------------------------------------------------------------------------------------------
Some application-level configuration rules are described here
------------------------------------------------------------------------------------------------------------------------------------------------------------------------
 ```
 
### ENV binding
Support resolving ENV variable names, as in the following example.
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
		So(x.UpdateWithFieldPathValues("http_address", "${XCONF_HOST}:${XCONF_PORT}"),ShouldBeNil)
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

### URL read
```golang
cc := &Config{}
x := xconf.NewWithoutFlagEnv(xconf.WithReaders(xconf.NewRemoteReader("http://127.0.0.1:9001/test.yaml", time.Duration(1)*time.Second)))
```

## Dynamic updates

### Configuration file based
```golang
	testBytesInMem := "memory_test_key"
	mem, err := xmem.New()
	panicErr(err)
	// xconf/kv provides an update mechanism based on ETCD/FILE/MEMORY
	// You can implement xconf's loader interface or interface to xmem to update the configuration with xmem's mechanism
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

### Based on `FieldPath`
```golang
err := xconf.WatchFieldPath("sub_test.http_address", func(from, to interface{}) {
	fmt.Printf("sub_test.http_address changed from %v to %v ", from, to)
})
panicErr(err)
```
File-based or Buffer-based updates can be implemented by the following methods, and the update results can be obtained asynchronously via `xconf.NotifyUpdate`, or synchronously via `xconf.Latest`.
- `UpdateWithFiles(files ...string) (err error)`
- `UpdateWithReader(readers ...io.Reader) (err error)`

File-based updates can be implemented by the following methods, the update result is obtained asynchronously by `xconf.NotifyUpdate`, or synchronously by `xconf.Latest`.
- `UpdateWithFlagArgs(flagArgs ...string)  (err error)`
- `UpdateWithEnviron(environ ...string) (err error)`
- `UpdateWithFieldPathValues(kv ...string) (err error)`

### Bind the latest configuration
```golang
xconf.Latest()
```

### Atomic auto-update
Using [optiongen](https://github.com/timestee/optiongen) to define a configuration and specifying `--xconf=true` to generate a configuration with `XConf` support generates `Atomic` update support by default for.
```golang

func (cc *Config) AtomicSetFunc() func(interface{}) { return AtomicConfigSet }

var atomicConfig unsafe.Pointer

func AtomicConfigSet(update interface{}) {
	atomic.StorePointer(&atomicConfig, (unsafe.Pointer)(update.(*Config)))
}
func AtomicConfig() ConfigVisitor {
	current := (*Config)(atomic.LoadPointer(&atomicConfig))
	if current == nil {
		atomic.CompareAndSwapPointer(&atomicConfig, nil, (unsafe.Pointer)(newDefaultConfig()))
		return (*Config)(atomic.LoadPointer(&atomicConfig))
	}
	return current
}

```

Just provide `AtomicConfig()` when parsing and automatically call back the `AtomicConfigSet` method for pointer replacement when the configuration is updated.
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

## Usage examples
### Loading encrypted configuration by URL
```golang
package main

import (
	"time"

	"github.com/sandwich-go/xconf"
	"github.com/sandwich-go/xconf/secconf"
	"github.com/sandwich-go/xconf/tests"
)

func main() {
	urlReader := xconf.NewRemoteReader("127.0.0.1:9001", time.Duration(1)*time.Second)
	key, _ := xconf.ParseEnvValue("${XXXTEA_KEY}|1dxz29pew", false)
	urlReaderSec := secconf.Reader(urlReader, secconf.StandardChainDecode(secconf.NewDecoderXXTEA([]byte(key))))
	xconf.Parse(tests.AtomicConfig(), xconf.WithReaders(urlReaderSec))
}
```

## Usage restrictions
### Private fields
If a private field is defined in the configuration or is hidden from XConf (specified as `-` by the xconf tag), the following usage restrictions apply when using the `Dynamic Update` feature.
- Bind the latest configuration in the active call to `Latest`.
- `Atomic` active binding mode, the configuration update

### Flag
- The configuration fields automatically created and defined in `FlagSet` are limited to the types supported by xconf/xflag.
- Complex types such as: "map[string][]time.Durtaion", "map[string]*Server", etc. cannot be created automatically and will have WARNGING logs printed, and these fields can be actively ignored by `WithFlagCreateIgnoreFiledPath`.
- Fields that cannot be created automatically in `FlagSet` cannot get the information and default values of the fields through `--help` or `Usage()`.

XConf cannot cache private and hidden field data according to `Parse`. In order to prevent possible data multi-processing access problems between the logical layer accessing the configuration and configuration update, when `Atomic` is passively updated or `Latest` is actively called to bind, the incoming structure constructs a new configuration structure, resulting in the data obtained at this time will not contain private fields and hidden fields.

When using the `Dynamic Update` feature, it is recommended that the private data or hidden fields be reassigned after the `Latest` call or in the set `InstallCallbackOnAtomicXXXXXXXXXSet` callback logic.

## xcmd Command Line Support
xcmd relies on xconf to automatically create, bind, and parse flag parameters, and supports custom flags, middleware, and subcommands. Reference: [xcmd/main/main.go](https://github.com/sandwich-go/xconf/blob/master/xcmd/main/main.go)

## help command extension
- `--help=yaml`
  - Print the current parsed configuration to the terminal in `-yaml` format
- `--help=. /test.yaml`
  - Print the currently parsed configuration in `-yaml` format to the specified file, which will be created automatically if the file does not exist

> Since the help command truncates the configuration parsing process, the output of the help extension command is the content of the incoming structure itself (the default value), not the content of the specified file, FLAG, ENV, etc.
