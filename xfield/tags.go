package xfield

import "strings"

const TagNotLeaf = "notleaf"

type TagList []string

func (t TagList) Has(opt string) bool {
	for _, tagOpt := range t {
		if tagOpt == opt {
			return true
		}
	}
	return false
}

func (t TagList) HasIgnoreCase(opt string) bool {
	for _, tagOpt := range t {
		if strings.EqualFold(tagOpt, opt) {
			return true
		}
	}
	return false
}

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
