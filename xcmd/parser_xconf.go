package xcmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/sandwich-go/xconf"
	"github.com/sandwich-go/xconf/xutil"
)

// ParserXConf xconf的Parser方法
func ParserXConf(ctx context.Context, cmd *Command, next Executer) error {
	cc := xconf.NewOptions(
		xconf.WithErrorHandling(xconf.ContinueOnError),
		xconf.WithFlagSet(cmd.FlagSet),
		xconf.WithFlagArgs(cmd.FlagArgs...))
	cc.ApplyOption(cmd.Config().GetXConfOption()...)

	x := xconf.NewWithConf(cc)
	keysList := xconf.FieldPathList(cmd.Bind(), x)

	bindFieldPath := cmd.bindFieldPath
	// keysList中的元素如果没有包含在GetBindFieldPath中，则为不允许通过flag覆盖的item
	var ignorePath []string
	// 检查GetBindFieldPath中的key是否合法
	var invalidKeys []string
	if len(bindFieldPath) > 0 {
		for _, k := range keysList {
			if !xutil.ContainStringEqualFold(bindFieldPath, k) {
				ignorePath = append(ignorePath, k)
			}
		}
		for _, v := range bindFieldPath {
			if !xutil.ContainString(keysList, v) {
				invalidKeys = append(invalidKeys, v)
			}
		}
	} else {
		ignorePath = keysList
	}

	if len(invalidKeys) > 0 {
		return cmd.wrapErr(fmt.Errorf("option BindFieldPath has invalid item: %s valid: %v", strings.Join(invalidKeys, ","), keysList))
	}
	// 更新忽略调的绑定字段，重新狗仔xconf实例
	cc.ApplyOption(xconf.WithFlagCreateIgnoreFiledPath(ignorePath...))
	x = xconf.NewWithConf(cc)
	// 更新FlagSet的Usage，使用xconf内置版本
	cmd.updateUsage(x)
	cc.FlagSet.Usage = cmd.usage
	err := x.Parse(cmd.bind)
	if err != nil {
		if IsErrHelp(err) {
			err = ErrHelp
		} else {
			err = fmt.Errorf("[ParserXConf] [%s] %s", cmd.name, err.Error())
		}
	}
	if err != nil {
		return err
	}
	return next(ctx, cmd)
}

// GenFieldPathStruct 生成filedPath struct
// todo 应随optiongen生成，手动指定FieldPath的时候可以防止出错，目前需要手动定义利用Command.Check检查
// type ConfigFieldPath struct {
// 	HttpAddress string
// 	Timeouts    string
// }

// func NewConfigFieldPath() *ConfigFieldPath {
// 	return &ConfigFieldPath{
// 		HttpAddress: "http_address",
// 		Timeouts:    "timeouts",
// 	}
// }
func GenFieldPathStruct(name string, fields map[string]xconf.StructFieldPathInfo) string {
	var lines []string
	structName := strings.Title(name) + "FieldPath"
	lines = append(lines, fmt.Sprintf("type %s struct{ ", structName))
	for _, v := range fields {
		lines = append(lines, fmt.Sprintf("	%s string", strings.Join(v.FieldNameList, "_")))
	}
	lines = append(lines, "}")
	lines = append(lines, fmt.Sprintf("func New%s() *%s { ", structName, structName))
	lines = append(lines, fmt.Sprintf("	return &%s{", structName))
	for k, v := range fields {
		lines = append(lines, fmt.Sprintf("		%s:\"%s\",", strings.Join(v.FieldNameList, "_"), k))
	}
	lines = append(lines, "	}")
	lines = append(lines, "}")
	return strings.Join(lines, "\n")
}
