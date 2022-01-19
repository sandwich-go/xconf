package xcmd

import (
	"fmt"
	"strings"

	"github.com/sandwich-go/xconf"
)

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
