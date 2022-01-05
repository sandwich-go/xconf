package tests

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
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

	x.Parse(cc)
	Convey("just parse with default value", t, func(c C) {
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
		x.Parse(cc)

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
			xconf.WithDebug(true),
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
		x.Parse(cc)
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
			select {
			case v := <-x.NotifyUpdate():
				updated <- v.(*Config)
			}
		}
	}()
	x.Parse(cc)
	cc.HttpAddress = "0.0.0.0"
	bytesBuffer := bytes.NewBuffer([]byte{})
	xconf.MustSaveVarToWriter(cc, xconf.ConfigTypeYAML, bytesBuffer)

	mem.Set(testBytesInMem, bytesBuffer.Bytes())
	Convey("with flag support", t, func(c C) {
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
			x.WatchFieldPath("sub_test.http_address", func(from, to interface{}) {
				fmt.Printf("sub_test.http_address changed from %v to %v ", from, to)
			})
			to := "123.456.789.000"
			err := x.UpdateWithFieldPathValues(watchFieldPath, to)
			So(err, ShouldBeNil)
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
			So(err, ShouldBeNil)

			//保存到字节流
			x2 := xconf.NewWithoutFlagEnv(xconf.WithReaders(bytes.NewReader(x.MustAsBytes(xconf.ConfigTypeYAML))))
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
		x.Parse(cc)

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
func newServerProvider(v interface{}) flag.Getter {
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
		err := x.Parse(cc)
		So(err, ShouldBeNil)
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
		err := x.Parse(cc)
		So(err, ShouldBeNil)
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
		err := x.Parse(cc)
		So(err, ShouldBeNil)
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
		So(x.UpdateWithFieldPathValues("http_address", "${XCONF_HOST_PORT}"), ShouldNotBeNil) // 解析会报错，XCONF_HOST_PORT不存在且为提供默认数值，Option指定ErrEnvBindNotExistWithoutDefault为true
		// todo update逻辑应该先保护本地绑定逻辑，当更新会引起绑定逻辑失败时，应该提前在update时返回错误，不能将错误扩散到绑定阶段
		_, err = x.Latest()
		So(err, ShouldNotBeNil)
	})
}

func TestRemoteReaderWithURL(t *testing.T) {
	Convey("env bind", t, func(c C) {
		cc := &Config{}
		x := xconf.NewWithoutFlagEnv(xconf.WithReaders(xconf.NewRemoteReader("127.0.0.1:0001", time.Duration(1)*time.Second)))
		So(func() { x.Parse(cc) }, ShouldPanic)
	})
}

func TestAtomicVal(t *testing.T) {
	Convey("atomic val", t, func(c C) {
		x := xconf.NewWithoutFlagEnv()
		So(x.Parse(AtomicConfig()), ShouldBeNil)
		So(x.UpdateWithFieldPathValues("http_address", "10.10.10.10"), ShouldBeNil)
		So(AtomicConfig().HttpAddress, ShouldEqual, "10.10.10.10")
	})
}

func TestUpdateGray(t *testing.T) {
	Convey("gray update", t, func(c C) {
		hostName, _ := os.Hostname()
		x := xconf.NewWithoutFlagEnv(xconf.WithAppLabelList(hostName))

		var yamlTest3 = []byte(fmt.Sprintf(`
http_address: 120.0.0.0
xconf_gray_rule_label: "%s"
`, hostName))

		fmt.Println("yamlTest3yamlTest3yamlTest3 ", string(yamlTest3))

		So(x.Parse(AtomicConfig()), ShouldBeNil)
		So(x.UpdateWithFieldPathValues("http_address", "10.10.10.10"), ShouldBeNil)
		So(AtomicConfig().HttpAddress, ShouldEqual, "10.10.10.10")
		So(x.UpdateWithReader(bytes.NewBuffer(yamlTest3)), ShouldBeNil)
		So(AtomicConfig().HttpAddress, ShouldEqual, "10.10.10.10")
	})
}
