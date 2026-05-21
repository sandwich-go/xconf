package xconf

import "strings"

// readIgnoreFieldsFromMeta 从 dataMeta 中读取 MetaKeyIgnoreFields 白名单
// 兼容 yaml 中的多种写法：单字符串、字符串数组、interface 数组
func readIgnoreFieldsFromMeta(meta map[string]interface{}) []string {
	if meta == nil {
		return nil
	}
	v, ok := meta[MetaKeyIgnoreFields]
	if !ok {
		return nil
	}
	return toStringSliceLoose(v)
}

// toStringSliceLoose 把任意常见的 yaml/json 反序列化结果转成 []string
// 支持: string / []string / []interface{}（元素可为 string）
func toStringSliceLoose(v interface{}) []string {
	switch val := v.(type) {
	case nil:
		return nil
	case string:
		if val == "" {
			return nil
		}
		return []string{val}
	case []string:
		out := make([]string, 0, len(val))
		for _, s := range val {
			if s != "" {
				out = append(out, s)
			}
		}
		return out
	case []interface{}:
		out := make([]string, 0, len(val))
		for _, it := range val {
			if s, ok := it.(string); ok && s != "" {
				out = append(out, s)
			}
		}
		return out
	}
	return nil
}

// mergeIgnoreFields 把 newOnes 合并进 base，做并集去重，返回新切片
// 顺序：base 在前，newOnes 中新出现的追加在后
func mergeIgnoreFields(base, newOnes []string) []string {
	if len(newOnes) == 0 {
		return base
	}
	seen := make(map[string]bool, len(base)+len(newOnes))
	out := make([]string, 0, len(base)+len(newOnes))
	for _, s := range base {
		if s == "" || seen[s] {
			continue
		}
		seen[s] = true
		out = append(out, s)
	}
	for _, s := range newOnes {
		if s == "" || seen[s] {
			continue
		}
		seen[s] = true
		out = append(out, s)
	}
	return out
}

// matchFieldPath 判断字段路径 path 是否命中白名单 patterns 中的任意一项
// 命中规则（按 . 边界）：
//
//  1. 精确匹配：path == pattern
//  2. pattern 是 path 的祖先：strings.HasPrefix(path, pattern + ".")
//     即白名单覆盖整棵子树，例如 pattern="plugin.foo" 命中 "plugin.foo.bar"
//  3. path 是 pattern 的祖先：strings.HasPrefix(pattern, path + ".")
//     适配 mapstructure 在顶层 key 缺失时只汇报顶层 key 的特性，
//     例如 yaml 写了 plugin.foo.bar 而 struct 中无 Plugin 字段，
//     mapstructure 仅会在 Unused 中记录 "plugin"。此时白名单写 "plugin.foo"
//     仍应放行该顶层 key（因为它一定来自配置中的 plugin.foo 子树）
//
// 例：pattern = "plugin.foo"
//
//	命中: "plugin.foo"、"plugin.foo.bar"、"plugin.foo.x.y"、"plugin"
//	不命中: "plugin.foobar"、"plugin.fo"、"pluginx"
func matchFieldPath(path string, patterns []string) bool {
	if path == "" || len(patterns) == 0 {
		return false
	}
	for _, p := range patterns {
		if p == "" {
			continue
		}
		if path == p {
			return true
		}
		// pattern 是 path 的祖先（白名单覆盖整棵子树）
		if strings.HasPrefix(path, p+".") {
			return true
		}
		// path 是 pattern 的祖先（白名单写得更细，但 mapstructure 只汇报到 path 层级）
		if strings.HasPrefix(p, path+".") {
			return true
		}
	}
	return false
}
