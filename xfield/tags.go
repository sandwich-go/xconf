package xfield

import "strings"

// TagNotLeaf xconf中指定TagNotLeaf标志字段非叶子节点，用于map
const TagNotLeaf = "notleaf"

// TagDeprecated xconf中指定TagDeprecated标记字段将被弃用
const TagDeprecated = "deprecated"

// TagList tag列表，全量
type TagList []string

// Has 是否含有指定的标签，区分大小写
func (t TagList) Has(opt string) bool {
	for _, tagOpt := range t {
		if tagOpt == opt {
			return true
		}
	}
	return false
}

// HasIgnoreCase 是否含有指定的标签，不区分大小写
func (t TagList) HasIgnoreCase(opt string) bool {
	for _, tagOpt := range t {
		if strings.EqualFold(tagOpt, opt) {
			return true
		}
	}
	return false
}

// ParseTag 解析指定的tag
// tag is one of followings:
// ""
// "name"
// "name,opt"
// "name,opt,opt2"
// ",opt"
func ParseTag(tag string) (string, TagList) {
	res := strings.Split(tag, ",")
	return res[0], res
}
