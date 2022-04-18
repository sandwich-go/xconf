package bugs

import (
	"bytes"
	"flag"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/sandwich-go/xconf"
	. "github.com/smartystreets/goconvey/convey"
)

var yamlTest = []byte(`
sm: 
  - foo : bar
  - bar : baz
`)

func TestBug_6(t *testing.T) {
	Convey("wrong type should got error", t, func(c C) {
		type Conf struct {
			Result map[string]string `xconf:"sm"`
		}
		// under MapMerge=false, map is leaf node
		cc := &Conf{Result: make(map[string]string)}
		{
			cc.Result = make(map[string]string)
			So(xconf.New(xconf.WithFlagSet(nil), xconf.WithReaders(bytes.NewReader(yamlTest))).Parse(cc), ShouldBeNil)
			So(cc.Result, ShouldResemble, map[string]string{"foo": "bar", "bar": "baz"})
		}
		{
			cc.Result = map[string]string{"test": "test"}
			So(xconf.New(xconf.WithFlagSet(nil), xconf.WithReaders(bytes.NewReader(yamlTest))).Parse(cc), ShouldBeNil)
			So(cc.Result, ShouldResemble, map[string]string{"foo": "bar", "bar": "baz"})
		}

		// under MapMerge=true, map is not leaf  node
		{
			cc.Result = map[string]string{}
			So(xconf.New(xconf.WithMapMerge(true), xconf.WithErrorHandling(xconf.ContinueOnError), xconf.WithFlagSet(nil), xconf.WithReaders(bytes.NewReader(yamlTest))).Parse(cc), ShouldBeNil)
			So(cc.Result, ShouldResemble, map[string]string{"foo": "bar", "bar": "baz"})
		}
		{
			// when the result have default element, we do not merge the element between map and slice
			cc.Result = map[string]string{"test": "test"}
			So(xconf.New(xconf.WithMapMerge(true), xconf.WithErrorHandling(xconf.ContinueOnError), xconf.WithFlagSet(nil), xconf.WithReaders(bytes.NewReader(yamlTest))).Parse(cc), ShouldNotBeNil)
		}
	})
}

func TestBug_7(t *testing.T) {
	Convey("should skip private field", t, func(c C) {
		type Server2 struct {
			Timeouts map[string]time.Duration
		}
		type Server struct {
			Timeouts      map[string]time.Duration `xconf:"timeouts"`
			privateServer Server2
		}
		ss := &Server{Timeouts: map[string]time.Duration{"read": time.Second}}
		ss.privateServer.Timeouts = map[string]time.Duration{"read": time.Second}
		So(xconf.SaveVarToWriterAsYAML(ss, os.Stderr), ShouldBeNil)
	})
}

type Nested1 struct {
	Deadline         time.Time `xconf:"deadline"`
	DeadlineAsSecond int
}
type Nested2 struct {
	TimeoutMap map[string]time.Duration `xconf:"timeout_map"`
}
type confTestNestedSquash struct {
	Nested1 `xconf:",squash"`
	Nested2 `xconf:",squash"`
}
type confTestNestedSquashOff struct {
	Nested1 `xconf:"nested1"`
	Nested2 `xconf:"nested2"`
}

func TestBug_8(t *testing.T) {
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
