package tests

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/sandwich-go/xconf"
	"github.com/sandwich-go/xconf/kv/xmem"
	"github.com/sandwich-go/xconf/xflag/vars"

	. "github.com/smartystreets/goconvey/convey"
)

func TestParseWithDefaultValue(t *testing.T) {
	defaultVal := NewTestConfig()
	cc := NewTestConfig()
	x := xconf.New(
		xconf.WithFlagSet(flag.NewFlagSet("TestParseWithDefaultValue", flag.ContinueOnError)),
		xconf.WithFlagArgs(),
	)
	Convey("just parse with default value", t, func(c C) {
		So(x.Parse(cc), ShouldBeNil)
		So(defaultVal, ShouldResemble, cc)
	})
}

var yamlTest = []byte(`
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
    test2222: 2222
  map_not_leaf:
    test2222: 2222
  servers:
    s1:
      timeouts:
        read: 20s 
`)

func TestOverideDefaultValue(t *testing.T) {
	Convey("parse with reader and overide default wieh map do not merge", t, func(c C) {
		cc := NewTestConfig()
		cc.SubTest.Map2 = make(map[string]int)
		cc.SubTest.MapNotLeaf = make(map[string]int)
		cc.SubTest.Map2["test1111"] = 1111
		cc.SubTest.MapNotLeaf["test1111"] = 1111
		cc.SubTest.Servers = make(map[string]Server)
		cc.SubTest.Servers["s1"] = Server{
			Timeouts: map[string]time.Duration{"read": time.Second * time.Duration(5)},
		}

		x := xconf.New(
			xconf.WithReaders(bytes.NewReader([]byte(yamlTest))),
			xconf.WithFlagSet(flag.NewFlagSet("TestOverideDefaultValue", flag.ContinueOnError)),
			xconf.WithFlagArgs(),
		)
		So(cc.Map1, ShouldResemble, map[string]int{"test1": 100, "test2": 200})
		So(cc.MapNotLeaf, ShouldResemble, map[string]int{"test1": 100, "test2": 200})
		So(cc.Int64Slice, ShouldResemble, []int64{101, 202, 303})
		So(x.Parse(cc), ShouldBeNil)

		So(cc.HttpAddress, ShouldEqual, ":3002")
		// map1作为叶子节点存在
		So(cc.Map1, ShouldResemble, map[string]int{"test1": 1000000})
		// MapNotLeaf指定了notleaf，覆盖的时候以key为单位
		So(cc.MapNotLeaf, ShouldResemble, map[string]int{"test1": 1000000, "test2": 200})
		So(cc.Int64Slice, ShouldResemble, []int64{1, 2})
		So(cc.DefaultEmptyMap, ShouldResemble, map[string]int{"test1": 1})
		So(cc.SubTest.Map2, ShouldResemble, map[string]int{"test2222": 2222})
		// MapNotLeaf指定了notleaf，覆盖的时候以key为单位
		So(cc.SubTest.MapNotLeaf, ShouldResemble, map[string]int{"test1111": 1111, "test2222": 2222})
		So(cc.SubTest.Servers["s1"].Timeouts, ShouldResemble, map[string]time.Duration{"read": time.Second * time.Duration(20)})

		// check latest
		latest, err := x.Latest()
		So(err, ShouldBeNil)
		latestConfig := latest.(*Config)
		So(latestConfig, ShouldResemble, cc)
	})

	Convey("with flag support", t, func(c C) {
		defaultVal := NewTestConfig()
		x := xconf.New(
			xconf.WithReaders(bytes.NewReader([]byte(yamlTest))),
			xconf.WithFlagArgs("--http_address=192.168.0.1", "--int64_slice=100,101", "--sub_test.map_not_leaf=k2,2222", "--sub_test.map2=k3,3333"),
			xconf.WithEnviron("read_timeout=20s", "map_not_leaf=k3,3333"),
		)
		cc := NewTestConfig()
		So(cc, ShouldResemble, defaultVal)
		cc.SubTest.MapNotLeaf = make(map[string]int)
		cc.SubTest.MapNotLeaf["k1"] = 1111
		cc.SubTest.Map2 = map[string]int{"test1": 11111}
		cc.MapNotLeaf = make(map[string]int)
		cc.MapNotLeaf["k1"] = 1111
		So(x.Parse(cc), ShouldBeNil)
		x.DumpInfo()
		So(cc.HttpAddress, ShouldEqual, "192.168.0.1")
		So(cc.Int64Slice, ShouldResemble, []int64{100, 101})
		So(cc.SubTest.MapNotLeaf, ShouldResemble, map[string]int{"k1": 1111, "k2": 2222, "test2222": 2222})
		So(cc.SubTest.Map2, ShouldResemble, map[string]int{"k3": 3333})

		So(cc.ReadTimeout, ShouldEqual, time.Duration(20)*time.Second)
		So(cc.MapNotLeaf, ShouldResemble, map[string]int{"k1": 1111, "k3": 3333, "test1": 1000000})
	})
}

func TestWatchUpdate(t *testing.T) {
	cc := NewTestConfig()
	cc.SubTest.Map2 = make(map[string]int)
	cc.SubTest.MapNotLeaf = make(map[string]int)
	cc.SubTest.Map2["test1111"] = 1111
	cc.SubTest.MapNotLeaf["test1111"] = 1111
	cc.SubTest.Servers = make(map[string]Server)
	cc.SubTest.Servers["s1"] = Server{
		Timeouts: map[string]time.Duration{"read": time.Second * time.Duration(5)},
	}
	x := xconf.New(
		xconf.WithReaders(bytes.NewReader([]byte(yamlTest))),
		xconf.WithFlagSet(flag.NewFlagSet("TestOverideDefaultValue", flag.ContinueOnError)),
		xconf.WithFlagArgs(),
		xconf.WithMapMerge(true))
	testBytesInMem := "memory_test_key"
	mem, _ := xmem.New()
	x.WatchUpdate(testBytesInMem, mem)
	updated := make(chan *Config, 1)
	go func() {
		for {
			v := <-x.NotifyUpdate()
			updated <- v.(*Config)
		}
	}()
	Convey("with flag support", t, func(c C) {

		So(x.Parse(cc), ShouldBeNil)
		cc.HttpAddress = "0.0.0.0"
		bytesBuffer := bytes.NewBuffer([]byte{})
		xconf.MustSaveVarToWriter(cc, xconf.ConfigTypeYAML, bytesBuffer)
		mem.Set(testBytesInMem, bytesBuffer.Bytes())

		gotUpdate := false
		var got *Config
		select {
		case got = <-updated:
			gotUpdate = true
		case <-time.After(time.Second * time.Duration(5)):
			gotUpdate = false
		}
		So(gotUpdate, ShouldBeTrue)
		So(got.HttpAddress, ShouldEqual, "0.0.0.0")

		Convey("watch filed path change", func(c C) {
			gotUpdate := false
			watchFieldPath := "sub_test.http_address"
			x.WatchFieldPath("sub_test.http_address", func(fieldPath string, from, to interface{}) {
				fmt.Printf("sub_test.http_address changed from %v to %v ", from, to)
			})
			to := "123.456.789.000"
			So(x.UpdateWithFieldPathValues(watchFieldPath, to), ShouldBeNil)
			var got *Config
			select {
			case got = <-updated:
				gotUpdate = true
				cc = got
			case <-time.After(time.Second * time.Duration(5)):
				gotUpdate = false
			}
			So(gotUpdate, ShouldBeTrue)
			So(got.SubTest.HTTPAddress, ShouldEqual, to)

			ccHash := x.Hash()

			//保存到字节流
			x2 := xconf.NewWithoutFlagEnv(xconf.WithReaders(bytes.NewReader(x.MustSaveToBytes(xconf.ConfigTypeYAML))))
			cc2 := NewTestConfig()
			So(x2.Parse(cc2), ShouldBeNil)
			fmt.Println(cc2)
			cc2Hash := x2.Hash()
			So(cc2, ShouldResemble, cc)
			So(cc2Hash, ShouldEqual, ccHash)
			So(x2.UpdateWithFieldPathValues("http_address", "10.10.10.10"), ShouldBeNil)
			cc2Hash = x2.Hash()
			So(cc2Hash, ShouldNotEqual, ccHash)

			latest, err := x.Latest()
			So(err, ShouldBeNil)
			latestConfig := latest.(*Config)
			So(latestConfig, ShouldResemble, got)
		})

	})
}

func TestMapMerge(t *testing.T) {
	cc := NewTestConfig()
	cc.SubTest.Map2 = make(map[string]int)
	cc.SubTest.MapNotLeaf = make(map[string]int)
	cc.SubTest.Map2["test1111"] = 1111
	cc.SubTest.MapNotLeaf["test1111"] = 1111
	cc.SubTest.Servers = make(map[string]Server)
	cc.SubTest.Servers["s1"] = Server{
		Timeouts: map[string]time.Duration{"read": time.Second * time.Duration(5)},
	}
	x := xconf.New(
		xconf.WithReaders(bytes.NewReader([]byte(yamlTest))),
		xconf.WithFlagSet(flag.NewFlagSet("TestOverideDefaultValue", flag.ContinueOnError)),
		xconf.WithFlagArgs(),
		xconf.WithMapMerge(true),
	)
	Convey("parse with reader and overide default wieh map merge", t, func(c C) {
		So(cc.Map1, ShouldResemble, map[string]int{"test1": 100, "test2": 200})
		So(cc.MapNotLeaf, ShouldResemble, map[string]int{"test1": 100, "test2": 200})
		So(cc.Int64Slice, ShouldResemble, []int64{101, 202, 303})
		So(x.Parse(cc), ShouldBeNil)

		So(cc.HttpAddress, ShouldEqual, ":3002")
		So(cc.Map1, ShouldResemble, map[string]int{"test1": 1000000, "test2": 200})
		So(cc.MapNotLeaf, ShouldResemble, map[string]int{"test1": 1000000, "test2": 200})
		So(cc.Int64Slice, ShouldResemble, []int64{1, 2})
		So(cc.DefaultEmptyMap, ShouldResemble, map[string]int{"test1": 1})
		So(cc.SubTest.Map2, ShouldResemble, map[string]int{"test2222": 2222, "test1111": 1111})
		// MapNotLeaf指定了notleaf，覆盖的时候以key为单位
		So(cc.SubTest.MapNotLeaf, ShouldResemble, map[string]int{"test1111": 1111, "test2222": 2222})
		So(cc.SubTest.Servers["s1"].Timeouts, ShouldResemble, map[string]time.Duration{"read": time.Second * time.Duration(20)})
	})
}

type confTestNestedSquash struct {
	Nested1 `xconf:",squash"`
	Nested2 `xconf:",squash"`
}
type confTestNestedSquashOff struct {
	Nested1 `xconf:"nested1"`
	Nested2 `xconf:"nested2"`
}

func TestSquash(t *testing.T) {
	Convey("TestSquash Enable", t, func(c C) {
		cc := &confTestNestedSquash{}
		cc.Nested1.Deadline = time.Now()
		cc.TimeoutMap = map[string]time.Duration{"read": time.Second}
		x := xconf.New(
			xconf.WithFlagSet(flag.NewFlagSet("suqash_anable", flag.ContinueOnError)),
			xconf.WithFlagArgs(),
			xconf.WithDebug(true),
			xconf.WithMapMerge(true),
		)
		So(x.Parse(cc), ShouldBeNil)
		So(strings.Contains(string(x.MustSaveToBytes(xconf.ConfigTypeYAML)), "nested"), ShouldBeFalse)
		So(strings.Contains(string(x.MustSaveToBytes(xconf.ConfigTypeJSON)), "nested"), ShouldBeFalse)
		So(strings.Contains(string(x.MustSaveToBytes(xconf.ConfigTypeTOML)), "nested"), ShouldBeFalse)
	})
	Convey("TestSquash Disable", t, func(c C) {
		cc := &confTestNestedSquashOff{}
		cc.Nested1.Deadline = time.Now()
		cc.Nested2.TimeoutMap = map[string]time.Duration{"read": time.Second}
		x := xconf.New(
			xconf.WithFlagSet(flag.NewFlagSet("suqash_disable", flag.ContinueOnError)),
			xconf.WithFlagArgs(),
			xconf.WithDebug(true),
			xconf.WithMapMerge(true),
		)
		So(x.Parse(cc), ShouldBeNil)
		So(strings.Contains(string(x.MustSaveToBytes(xconf.ConfigTypeYAML)), "nested"), ShouldBeTrue)
		So(strings.Contains(string(x.MustSaveToBytes(xconf.ConfigTypeJSON)), "nested"), ShouldBeTrue)
		So(strings.Contains(string(x.MustSaveToBytes(xconf.ConfigTypeTOML)), "nested"), ShouldBeTrue)
	})
}

type TestConf1 struct {
	HTTPAddress string   `xconf:"http_address" default:"0.0.0.0:0000"`
	Hosts       []string `flag:"hosts" cfg:"hosts" default:"127.0.0.0,127.0.0.1"`
	LogLevel    int      `default:"3"`
	BoolVar     bool     `default:"true"`
	IntSlice    []int64  `cfg:"int_slice"`
}

func TestParseDefault(t *testing.T) {
	Convey("parse default value", t, func(c C) {
		cc := &TestConf1{}
		x := xconf.New(
			xconf.WithFlagSet(flag.NewFlagSet("TestOverideDefaultValue", flag.ContinueOnError)),
			xconf.WithFlagArgs(),
			xconf.WithMapMerge(true),
		)
		So(x.Parse(cc), ShouldBeNil)
		So(cc.HTTPAddress, ShouldEqual, "0.0.0.0:0000")
		So(cc.LogLevel, ShouldEqual, 3)
		So(cc.BoolVar, ShouldEqual, true)
		So(cc.Hosts, ShouldResemble, []string{"127.0.0.0", "127.0.0.1"})

	})
}

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
func newServerProvider(v interface{}, alias func(string) string) flag.Getter {
	return &serverProvider{data: v.(*map[string]Server)}
}
func TestFlagProvider(t *testing.T) {
	Convey("flag provider", t, func(c C) {
		cc := &Config{}
		jsonServer := `json@{"s1":{"timeouts":{"read":5000000000},"timeouts_not_leaf":{"write":5000000000}}}`
		x := xconf.New(
			xconf.WithFlagSet(flag.NewFlagSet("TestFlagProvider", flag.ContinueOnError)),
			xconf.WithFlagArgs("--sub_test.servers="+jsonServer),
		)
		vars.SetProviderByFieldPath("sub_test.servers", newServerProvider)
		So(x.Parse(cc), ShouldBeNil)
		So(cc.SubTest.Servers["s1"].Timeouts, ShouldResemble, map[string]time.Duration{"read": time.Duration(5) * time.Second})
	})
}
func TestFlagProviderForEnv(t *testing.T) {
	Convey("flag provider for env", t, func(c C) {
		cc := &Config{}
		jsonServer := `json@{"s1":{"timeouts":{"read":5000000000},"timeouts_not_leaf":{"write":5000000000}}}`
		x := xconf.New(
			xconf.WithFlagSet(nil),
			xconf.WithEnviron("sub_test_servers="+jsonServer),
		)
		vars.SetProviderByFieldPath("sub_test.servers", newServerProvider)
		So(x.Parse(cc), ShouldBeNil)
		So(cc.SubTest.Servers["s1"].Timeouts, ShouldResemble, map[string]time.Duration{"read": time.Duration(5) * time.Second})
	})
}

func TestFlagProviderByType(t *testing.T) {
	Convey("flag provider with type", t, func(c C) {
		cc := &Config{}
		jsonServer := `json@{"s1":{"timeouts":{"read":5000000000},"timeouts_not_leaf":{"write":5000000000}}}`
		x := xconf.New(
			xconf.WithFlagSet(nil),
			xconf.WithEnviron("sub_test_servers="+jsonServer),
		)
		vars.SetProviderByFieldType("map[string]Server", newServerProvider)
		So(x.Parse(cc), ShouldBeNil)
		So(cc.SubTest.Servers["s1"].Timeouts, ShouldResemble, map[string]time.Duration{"read": time.Duration(5) * time.Second})
	})
}

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
    ${IP_ADDRESS|10.0.0.1}: 2222
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
		So(x.Parse(cc), ShouldBeNil)
		So(x.UpdateWithFieldPathValues("http_address", "${XCONF_HOST}:${XCONF_PORT}"), ShouldNotBeNil)
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
		So(x.UpdateWithFieldPathValues("http_address", "${XCONF_HOST_PORT}"), ShouldNotBeNil) // 解析会报错，XCONF_HOST_PORT不存在且为提供默认数值，Option指定ErrEnvBindNotExistWithoutDefault为true
		_, err = x.Latest()
		So(err, ShouldBeNil)
	})
}

func TestEnvBindWithPrefix(t *testing.T) {
	Convey("env bind", t, func(c C) {
		cc := &Config{}
		x := xconf.New(xconf.WithFlagSet(nil), xconf.WithEnvironPrefix("TEST_PREFIX_"))
		So(x.Parse(cc), ShouldBeNil)
		So(cc.HttpAddress, ShouldEqual, "")
		So(x.UpdateWithEnviron("http_address=10.0.0.1"), ShouldBeNil)
		latest, err := x.Latest()
		So(err, ShouldBeNil)
		So(latest.(*Config).HttpAddress, ShouldEqual, "")
		So(x.UpdateWithEnviron("TEST_PREFIX_http_address=10.0.0.1"), ShouldBeNil)
		latest, err = x.Latest()
		So(err, ShouldBeNil)
		So(latest.(*Config).HttpAddress, ShouldEqual, "10.0.0.1")
	})
}

func TestRemoteReaderWithURL(t *testing.T) {
	Convey("env bind", t, func(c C) {
		cc := &Config{}
		x := xconf.NewWithoutFlagEnv(xconf.WithReaders(xconf.NewRemoteReader("127.0.0.1:0001", time.Duration(1)*time.Second)))
		So(func() { _ = x.Parse(cc) }, ShouldPanic)
	})
}

func TestAtomicVal(t *testing.T) {
	Convey("atomic val", t, func(c C) {
		x := xconf.NewWithoutFlagEnv()
		So(x.Parse(AtomicConfig()), ShouldBeNil)
		So(x.UpdateWithFieldPathValues("http_address", "10.10.10.10"), ShouldBeNil)
		So(AtomicConfig().GetHttpAddress(), ShouldEqual, "10.10.10.10")
	})
}

func TestUpdateGray(t *testing.T) {
	Convey("gray update", t, func(c C) {
		hostName, _ := os.Hostname()
		x := xconf.NewWithoutFlagEnv(xconf.WithAppLabelList(hostName))

		var yamlTest3 = []byte(`
http_address: 120.0.0.0
xconf_gray_rule_label: "test-label"
`)
		So(x.Parse(AtomicConfig()), ShouldBeNil)
		So(x.UpdateWithFieldPathValues("http_address", "10.10.10.10"), ShouldBeNil)
		So(AtomicConfig().GetHttpAddress(), ShouldEqual, "10.10.10.10")
		So(x.UpdateWithReader(bytes.NewBuffer(yamlTest3)), ShouldBeNil)
		So(AtomicConfig().GetHttpAddress(), ShouldEqual, "10.10.10.10")
	})
}

func TestStringAlias(t *testing.T) {
	Convey("string alias", t, func(c C) {
		x := xconf.NewWithoutFlagEnv()
		So(x.Parse(AtomicConfig()), ShouldBeNil)
		So(x.UpdateWithFieldPathValues("process_count", "runtime.NumCPU"), ShouldBeNil)
		So(AtomicConfig().GetProcessCount(), ShouldEqual, runtime.NumCPU())
	})
}

type TestConf2 struct {
	LogLevel    int `xconf:"log_level"`
	NoConfValue int `xconf:"-"`
	private     int
}

func TestNotConfField(t *testing.T) {
	Convey("not conf field", t, func(c C) {
		x := xconf.NewWithoutFlagEnv(xconf.WithDebug(false))
		cc := &TestConf2{NoConfValue: 2, private: 1}
		So(x.Parse(cc), ShouldBeNil)
		So(cc.private, ShouldEqual, 1)
		So(cc.NoConfValue, ShouldEqual, 2)
		So(x.UpdateWithFieldPathValues("log_level", "2"), ShouldBeNil)
		ccInterface, err := x.Latest()
		So(err, ShouldBeNil)
		cc = ccInterface.(*TestConf2)
		// todo private字段, xconf内部无法获取，再次绑定到新值会导致private字段数据丢失
		// todo no conf字段，xconf内没有解析该字段，导致缓存的数据中没有这个字段，再次绑定时也会丢失该数据
		// todo 现在的应用场景，对于private，noconf字段，如果是在更新场景下，或者主动通过Latest再次绑定会丢失数据，建议绑定后再次走一遍private，noconf的赋值逻辑
		// todo Atomic自动更新逻辑下可以通过InstallCallbackOnAtomicConfigSet设定更新时的回调，再次为private，no conf field赋值
		So(cc.private, ShouldEqual, 0)
		So(cc.NoConfValue, ShouldEqual, 0)
		x.MustSaveToWriter(xconf.ConfigTypeYAML, os.Stderr)
	})
}

var y1 = []byte(`
map1:
  k1: "v100"
map2:
  k1: 1s
sub_conf1:
  field1: 1
  field2: test
  sub_conf2:
    field1: 1
    field2: test
  sub_conf3:
    field1: 1   
    field2: test
`)
var y2 = []byte(`
map1:
  k2: "v200"
`)

var y3 = []byte(`
#清空map1
map1:
#清空map2
map2:
#清空sub_conf1,应该也同步清空其子元素
sub_conf1:
`)

func TestMapDelete(t *testing.T) {
	Convey("map delete", t, func(c C) {
		type SubConf3 struct {
			Field1 int
			Field2 string
		}
		type SubConf2 struct {
			Field1 int
			Field2 string
		}
		type SubConf1 struct {
			Field1   int
			Field2   string
			SubConf2 SubConf2
			SubConf3 *SubConf3
		}
		type Conf1 struct {
			Map1     map[string]string
			Map2     map[string]time.Duration
			SubConf1 SubConf1
		}
		{
			c1 := &Conf1{Map1: make(map[string]string)}
			xx := xconf.New(xconf.WithFlagSet(nil), xconf.WithMapMerge(false),
				xconf.WithReaders(bytes.NewReader(y1), bytes.NewReader(y2)))
			So(xx.Parse(c1), ShouldBeNil)
			So(c1.Map1, ShouldResemble, map[string]string{"k2": "v200"})
			So(c1.Map2, ShouldResemble, map[string]time.Duration{"k1": time.Second})
			So(c1.SubConf1.Field1, ShouldEqual, 1)
			So(c1.SubConf1.Field2, ShouldEqual, "test")
			So(c1.SubConf1.SubConf2.Field1, ShouldEqual, 1)
			So(c1.SubConf1.SubConf2.Field2, ShouldEqual, "test")
			So(c1.SubConf1.SubConf3.Field1, ShouldEqual, 1)
			So(c1.SubConf1.SubConf3.Field2, ShouldEqual, "test")
		}
		{
			c1 := &Conf1{Map1: make(map[string]string)}
			xx := xconf.New(xconf.WithFlagSet(nil), xconf.WithMapMerge(true),
				xconf.WithReaders(bytes.NewReader(y1), bytes.NewReader(y2)))
			So(xx.Parse(c1), ShouldBeNil)
			So(c1.Map1, ShouldResemble, map[string]string{"k1": "v100", "k2": "v200"})
			So(c1.SubConf1.Field1, ShouldEqual, 1)
			So(c1.SubConf1.Field2, ShouldEqual, "test")
			So(c1.SubConf1.SubConf2.Field1, ShouldEqual, 1)
			So(c1.SubConf1.SubConf2.Field2, ShouldEqual, "test")
			So(c1.SubConf1.SubConf3.Field1, ShouldEqual, 1)
			So(c1.SubConf1.SubConf3.Field2, ShouldEqual, "test")
		}
		{
			c1 := &Conf1{Map1: make(map[string]string)}
			xx := xconf.New(xconf.WithFlagSet(nil), xconf.WithMapMerge(true),
				xconf.WithReaders(bytes.NewReader(y1), bytes.NewReader(y2), bytes.NewReader(y3)))
			So(xx.Parse(c1), ShouldBeNil)
			So(len(c1.Map1), ShouldBeZeroValue)
			So(c1.Map1, ShouldResemble, map[string]string{})
			So(c1.Map2, ShouldResemble, map[string]time.Duration{})
			So(c1.SubConf1.Field1, ShouldEqual, 0)
			So(c1.SubConf1.Field2, ShouldEqual, "")
			So(c1.SubConf1.SubConf2.Field1, ShouldEqual, 0)
			So(c1.SubConf1.SubConf2.Field2, ShouldEqual, "")
			So(c1.SubConf1.SubConf3.Field1, ShouldEqual, 0)
			So(c1.SubConf1.SubConf3.Field2, ShouldEqual, "")
		}
	})
}
