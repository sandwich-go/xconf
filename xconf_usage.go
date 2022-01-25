package xconf

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/sandwich-go/xconf/xfield"
	"github.com/sandwich-go/xconf/xflag"
	"github.com/sandwich-go/xconf/xutil"
)

// Usage 打印usage信息
func (x *XConf) Usage() { x.UsageToWriter(os.Stderr, x.cc.FlagArgs...) }

func (x *XConf) UsageToWriter(w io.Writer, args ...string) {
	parsedOptions := xflag.ParseArgsToMapStringString(args)
	val, got := parsedOptions["help"]
	if !got {
		val, got = parsedOptions["h"]
	}
	var err error
	if got && strings.EqualFold(xutil.StringTrim(val), "xconf") {
		// 指定xconf_usage的FlagArgs为空，避免再次触发help逻辑
		xx := New(WithFlagSet(newFlagSetContinueOnError("xconf_usage")), WithFlagArgs(), WithErrorHandling(ContinueOnError))
		cc := NewOptions()
		xutil.PanicErr(xx.Parse(cc))
		err = xx.usageLinesToWriter(w)
	} else {
		err = x.usageLinesToWriter(w)
	}
	if err != nil {
		x.cc.LogWarning(fmt.Sprintf("UsageToWriter got error:%s", err.Error()))
	}
}

// usageLinesToWriter 打印usage信息到io.Writer
func (x *XConf) usageLinesToWriter(w io.Writer) error {
	using := x.zeroValPtrForLayout
	optionUsageStr := x.optionUsage
	if using == nil {
		return errors.New("usageToWriter: should parse first")
	}
	lines, magic, err := x.usageLines(using)
	if err != nil {
		return fmt.Errorf("Usage err: " + err.Error())
	}
	fmt.Fprintln(w, xutil.TableFormat(lines, magic, true, optionUsageStr))
	return nil
}

func (x *XConf) usageLines(valPtr interface{}) ([]string, string, error) {
	magic := "\x00"
	var lineAll []string
	lineAll = append(lineAll, "FLAG"+"\x00"+"ENV"+"\x00"+"TYPE"+"\x00"+"USAGE")
	allFlag := xflag.GetFlagInfo(x.cc.FlagSet)
	for _, v := range allFlag.List {
		line := fmt.Sprintf("--%s", v.Name)
		line += magic
		tag := FlagTypeStr(x, v.Name)
		if tag == "-" {
			// - 脱离xconf的tag, flag只是我们操作的原子单位，无法将数据附加到flag，再次更新
			// M xconf原子tag，但通过环境变量设置的意义不大，考虑移除这部分对环境变量的支持
			line += "-"
		} else {
			line += xflag.FlagToEnvUppercase(x.cc.EnvironPrefix, v.Name)
		}
		line += magic
		line += v.TypeName
		line += magic
		usage := ""
		if info, ok := x.fieldPathInfoMap[v.Name]; ok {
			usage = info.Tag.Get("usage")
		}
		if usage == "" {
			usage = v.Usage
		}
		if !xflag.IsZeroValue(v.Flag, v.DefValue) {
			if v.TypeName == "string" {
				usage += fmt.Sprintf(" (default %q)", v.DefValue)
			} else {
				usage += fmt.Sprintf(" (default %s)", v.DefValue)
			}
		}

		line += fmt.Sprintf("|%s| %s", tag, usage)
		lineAll = append(lineAll, line)
	}
	return lineAll, magic, nil
}

// FlagTypeStr 获取子弹标记，Y代表xconf解析管理的配置，M标识xconf内置配置，D标识Deprecated，-表示为非xconf管理的配置
func FlagTypeStr(x *XConf, name string) (tag string) {
	v, ok := x.fieldPathInfoMap[name]
	if !ok {
		if xutil.ContainStringEqualFold(metaKeyList, name) {
			return "M"
		}
		return "-"
	}
	if v.TagListXConf.HasIgnoreCase(xfield.TagDeprecated) {
		return "D"
	}
	return "Y"
}

// DumpInfo 打印调试信息
func (x *XConf) DumpInfo() {
	var lines []string
	lines = append(lines, fmt.Sprintf("# FieldPath: \n%v", x.keysList()))
	lines = append(lines, fmt.Sprintf("# DataDest: \n%v", x.dataLatestCached))
	lines = append(lines, fmt.Sprintf("# DataMeta: \n%v", x.dataMeta))
	hashCode := x.Hash()
	lines = append(lines, fmt.Sprintf("# Hash Local  : %s", hashCode))
	hashCenter := DefaultInvalidHashString
	if center := x.dataMeta[MetaKeyLatestHash]; center != nil {
		hashCenter = center.(string)
	}
	lines = append(lines, fmt.Sprintf("# Hash Center : %s", hashCenter))

	usageLines, magic, err := x.usageLines(x.zeroValPtrForLayout)
	if err != nil {
		x.cc.LogWarning("got error:" + err.Error())
		return
	}
	fmt.Fprintln(os.Stderr, xutil.TableFormat(usageLines, magic, true, lines...))

}
