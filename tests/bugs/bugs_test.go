package tests

import (
	"testing"

	"github.com/sandwich-go/xconf"
	. "github.com/smartystreets/goconvey/convey"
)

//go:generate optiongen --option_with_struct_name=false --new_func=NewNodeConfig --xconf=true --empty_composite_nil=true --usage_tag_name=usage
func NodeConfigOptionDeclareWithDefault() interface{} {
	return map[string]interface{}{
		"Master":   false,       // @MethodComment(是否为master节点，如果是，则为true，所有的节点，只有一个master节点)
		"Host":     "127.0.0.1", // @MethodComment(mysql节点host)
		"Port":     3306,        // @MethodComment(mysql节点port)
		"User":     "root",      // @MethodComment(mysql节点user)
		"Password": "",          // @MethodComment(mysql节点password)
	}
}

//go:generate optiongen --option_with_struct_name=false --new_func=NewConfig --xconf=true --empty_composite_nil=true --usage_tag_name=usage
func ConfigOptionDeclareWithDefault() interface{} {
	return map[string]interface{}{
		"Nodes": []*NodeConfig{}, // @MethodComment(所有的mysql节点信息)
	}
}

func TestMapSlice(t *testing.T) {
	Convey("not conf field", t, func(c C) {
		cc := NewConfig(WithNodes(NewNodeConfig(), NewNodeConfig(WithHost("127.0.0.2"))))
		x := xconf.New()
		err := x.Parse(cc)
		So(err, ShouldBeNil)
		saveTo := "./test.yaml"
		// x.SaveToFile(saveTo)
		ccNew := NewConfig()
		err = xconf.New().Parse(ccNew, xconf.WithFiles(saveTo))
		So(err, ShouldBeNil)
	})
}
