package xconf

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/sandwich-go/xconf/xfield"
	"github.com/sandwich-go/xconf/xflag"
	"github.com/sandwich-go/xconf/xutil"
)

// Usage 打印usage信息
func (x *XConf) Usage() { x.UsageToWriter(os.Stderr, x.cc.FlagArgs...) }

func (x *XConf) UsageToWriter(w io.Writer, args ...string) {
	err := x.usageToWriter(w, args...)
	if err == nil {
		return
	}
	x.cc.LogWarning(fmt.Sprintf("UsageToWriter got error:%s", err.Error()))
}

func (x *XConf) usageToWriter(w io.Writer, args ...string) (err error) {
	parsedOptions := xflag.ParseArgsToMapStringString(args)
	val, got := parsedOptions["help"]
	if !got {
		val, got = parsedOptions["h"]
	}
	val = xutil.StringTrim(val)
	if got && strings.EqualFold(xutil.StringTrim(val), "xconf") {
		// 指定xconf_usage的FlagArgs为空，避免再次触发help逻辑
		xx := New(WithFlagSet(newFlagSetContinueOnError("xconf_usage")), WithFlagArgs(), WithErrorHandling(ContinueOnError))
		cc := NewOptions()
		xutil.PanicErr(xx.Parse(cc))
		return xx.usageLinesToWriter(w)
	}
	if got && strings.EqualFold(xutil.StringTrim(val), "yaml") {
		if x.valPtrForUsageDump == nil {
			return errors.New("usage for yaml got empty config input")
		}
		return x.SaveVarToWriterAsYAML(x.valPtrForUsageDump, w)
	}
	if got && strings.HasSuffix(val, string(ConfigTypeYAML)) { // 输出到文件
		defer func() {
			if err == nil {
				fmt.Println("\n🍺 save config file as yaml to: ", val)
			} else {
				fmt.Println("\n🚫 got error while save config file as yaml, err:", err.Error())
			}
			err = nil
		}()
		if x.valPtrForUsageDump == nil {
			return errors.New("usage for yaml file got config input")
		}
		bytesBuffer := bytes.NewBuffer([]byte{})
		err := x.SaveVarToWriterAsYAML(x.valPtrForUsageDump, bytesBuffer)
		if err != nil {
			return err
		}
		return xutil.FilePutContents(val, bytesBuffer.Bytes())
	}
	return x.usageLinesToWriter(w)
}

// usageLinesToWriter 打印usage信息到io.Writer
func (x *XConf) usageLinesToWriter(w io.Writer) error {
	optionUsageStr := x.optionUsage
	lines, magic, err := x.usageLines()
	if err != nil {
		return fmt.Errorf("Usage err: " + err.Error())
	}
	fmt.Fprintln(w, xutil.TableFormat(lines, magic, true, optionUsageStr))
	return nil
}

func (x *XConf) usageLines() ([]string, string, error) {
	magic := "\x00"
	var lineAll []string
	lineAll = append(lineAll, "FLAG"+magic+"ENV"+magic+"TYPE"+magic+"USAGE")
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
	sort.Strings(lineAll[1:])
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

	usageLines, magic, err := x.usageLines()
	if err != nil {
		x.cc.LogWarning("got error:" + err.Error())
		return
	}
	fmt.Fprintln(os.Stderr, xutil.TableFormat(usageLines, magic, true, lines...))

}
