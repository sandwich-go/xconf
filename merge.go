package xconf

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/sandwich-go/xconf/xfield"
	"github.com/sandwich-go/xconf/xutil"
)

// 当前key指向的结构是否采用sub key级别的merge方案
// map的val只能作为叶子节点存在，不允许notleaf，避免复杂的配置覆盖逻辑
// 要求fieldPathMap是基于`空`结构数据生成，否则对MAP结构信息判断会造成影响
func isLeafFieldPath(fieldPathMap map[string]StructFieldPathInfo, fieldPath string) bool {
	prefix := fieldPath + DefaultKeyDelim //避免两个字段但是拥有共同前缀，如TypeMapStringInt TypeMapStringIntNotLeaf
	for k := range fieldPathMap {
		if strings.HasPrefix(k, prefix) {
			// k是fieldPath的子节点，则fieldPath不指向叶子节点
			return false
		}
	}

	// 没有找到fieldPath.为前缀的节点，可能是一个基础类型，slice或者是一个map，map可能是一个非叶子节点，通过是否配置了xfield.TagNotLeaf标签决定
	if field, ok := fieldPathMap[fieldPath]; ok && field.TagListXConf.HasIgnoreCase(xfield.TagNotLeaf) {
		return false
	}

	return true
}

type fieldValues struct {
	fieldPath string
	from      interface{}
	to        interface{}
}
type fieldChanges struct {
	changed map[string]*fieldValues
}

func (c *fieldChanges) Set(fieldPath string, from, to interface{}) {
	if c.changed == nil {
		c.changed = make(map[string]*fieldValues)
	}
	c.changed[fieldPath] = &fieldValues{
		fieldPath: fieldPath,
		from:      from,
		to:        to,
	}
}

func mergeMap(
	prefix string,
	depth int,
	logger LogFunc,
	src, dst map[string]interface{},
	isLeafFieldPath func(fieldPath string) bool,
	itgt map[interface{}]interface{},
	changes *fieldChanges) error {

	indent := ""
	for i := 0; i < depth; i++ {
		indent += "   "
	}
	if prefix != "" {
		logger(fmt.Sprintf("%s----> merge prefix: %s\n", indent, prefix))
	}
	indentNow := indent + "      "
	depth += 2
	if prefix != "" {
		prefix += DefaultKeyDelim
	}
	for srcKey, srcVal := range src {
		fieldPath := prefix + srcKey
		dstKey := keyExists(srcKey, dst)
		if dstKey == "" {
			srcValCleaned := valClean(srcVal)
			dst[srcKey] = srcValCleaned
			if changes != nil {
				changes.Set(fieldPath, nil, srcValCleaned)
			}
			if itgt != nil {
				itgt[srcKey] = srcValCleaned
			}
			logger(fmt.Sprintf("%s%s srcKey:%s not in dst, overide.\n", indentNow, fieldPath, srcKey))
			continue
		}

		dstVal, ok := dst[dstKey]
		if !ok {
			srcValCleaned := valClean(srcVal)
			dst[srcKey] = srcValCleaned
			if changes != nil {
				changes.Set(fieldPath, nil, srcValCleaned)
			}
			if itgt != nil {
				itgt[srcKey] = srcValCleaned
			}
			logger(fmt.Sprintf("%s%sdstKey:%s in dst no value, overide\n", indentNow, fieldPath, srcKey))
			continue
		}
		if isLeafFieldPath(fieldPath) {
			dst[dstKey] = srcVal
			if changes != nil {
				changes.Set(fieldPath, nil, srcVal)
			}
			if itgt != nil {
				itgt[dstKey] = srcVal
			}
			logger(fmt.Sprintf("%s%s is leaf key, overide.\n", indentNow, fieldPath))
		} else {
			var mergeErr error
			switch dstValType := dstVal.(type) {
			case map[interface{}]interface{}:
				logger(fmt.Sprintf("%s%s dstVal is map[interface{}]interface{}, deep merge.\n", indentNow, fieldPath))
				tsv := srcVal.(map[interface{}]interface{})
				ssv := castToMapStringInterface(tsv)
				stv := castToMapStringInterface(dstValType)
				mergeErr = mergeMap(fieldPath, depth, logger, ssv, stv, isLeafFieldPath, dstValType, changes)
			case map[string]interface{}:
				switch srcValType := srcVal.(type) {
				case map[interface{}]interface{}:
					logger(fmt.Sprintf("%s%s dstVal: map[string]interface{} srcVal:map[interface{}]interface{}, deep merge\n", indentNow, fieldPath))
					ssv := castToMapStringInterface(srcValType)
					mergeErr = mergeMap(fieldPath, depth, logger, ssv, dstValType, isLeafFieldPath, nil, changes)
				case map[string]interface{}:
					logger(fmt.Sprintf("%s%s dstVal: map[string]interface{} srcVal:map[string]interface{}, deep merge\n", indentNow, fieldPath))
					mergeErr = mergeMap(fieldPath, depth, logger, srcValType, dstValType, isLeafFieldPath, nil, changes)
				default:
					// 如果dest是map结构但是src不是map结构, 检测srcVal是否为空，不为空则返回错误
					if !xutil.IsEmpty(srcVal) {
						return fmt.Errorf("dst is map but src not map not empty,got:%v , while merge:%s", srcVal, fieldPath)
					}
					srcVal = make(map[string]interface{})
					dst[dstKey] = srcVal
					if changes != nil {
						changes.Set(fieldPath, nil, srcVal)
					}
					if itgt != nil {
						itgt[dstKey] = srcVal
					}
					logger(fmt.Sprintf("%s%s srcVal type:%v,overide\n", indentNow, fieldPath, reflect.TypeOf(srcValType)))
				}
			default:
				dst[dstKey] = srcVal
				if changes != nil {
					changes.Set(fieldPath, nil, srcVal)
				}
				if itgt != nil {
					itgt[dstKey] = srcVal
				}
				logger(fmt.Sprintf("%s%s dstVal type:%v,overide\n", indentNow, fieldPath, reflect.TypeOf(dstVal)))
			}
			if mergeErr != nil {
				return fmt.Errorf("got err:%w while merge:%s", mergeErr, fieldPath)
			}
		}
	}
	return nil
}

func valClean(src interface{}) interface{} {
	switch srcType := src.(type) {
	case map[interface{}]interface{}:
		tgt := map[string]interface{}{}
		for k, v := range srcType {
			tgt[fmt.Sprintf("%v", k)] = valClean(v)
		}
		return tgt
	case map[string]interface{}:
		tgt := map[string]interface{}{}
		for k, v := range srcType {
			tgt[k] = valClean(v)
		}
		return tgt
	default:
		return srcType
	}
}

func castToMapStringInterface(
	src map[interface{}]interface{}) map[string]interface{} {
	tgt := map[string]interface{}{}
	for k, v := range src {
		tgt[fmt.Sprintf("%v", k)] = v
	}
	return tgt
}

func keyExists(k string, m map[string]interface{}) string {
	lk := strings.ToLower(k)
	for mk := range m {
		lmk := strings.ToLower(mk)
		if lmk == lk {
			return mk
		}
	}
	return ""
}

func isNil(value interface{}, traceSource ...bool) bool {
	if value == nil {
		return true
	}
	var rv reflect.Value
	if v, ok := value.(reflect.Value); ok {
		rv = v
	} else {
		rv = reflect.ValueOf(value)
	}
	switch rv.Kind() {
	case reflect.Chan,
		reflect.Map,
		reflect.Slice,
		reflect.Func,
		reflect.Interface,
		reflect.UnsafePointer:
		return !rv.IsValid() || rv.IsNil()

	case reflect.Ptr:
		if len(traceSource) > 0 && traceSource[0] {
			for rv.Kind() == reflect.Ptr {
				rv = rv.Elem()
			}
			if !rv.IsValid() {
				return true
			}
			if rv.Kind() == reflect.Ptr {
				return rv.IsNil()
			}
		} else {
			return !rv.IsValid() || rv.IsNil()
		}
	}
	return false
}
