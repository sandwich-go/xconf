package bugs

import (
	"bytes"
	"testing"

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
