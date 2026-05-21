package xconf

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
)

func (x *XConf) loadFiles(files ...string) (map[string]interface{}, error) {
	finalData := make(map[string]interface{})
	for _, file := range files {
		x.cc.LogDebug(fmt.Sprintf("load file: %s", file))
		data, err := x.loadFile(file)
		if err != nil {
			return finalData, fmt.Errorf("got error while load file:%s err:%w", file, err)
		}
		// 对 xconf_ignore_fields 做特殊处理：本次加载与已聚合的白名单做并集去重
		mergedIgnore := mergeIgnoreFields(
			toStringSliceLoose(finalData[MetaKeyIgnoreFields]),
			toStringSliceLoose(data[MetaKeyIgnoreFields]),
		)
		err = x.mergeMap("file:"+file, "file:tmp", data, finalData)
		if err != nil {
			return finalData, fmt.Errorf("got error while merge file:%s err:%w", file, err)
		}
		if len(mergedIgnore) > 0 {
			finalData[MetaKeyIgnoreFields] = mergedIgnore
		}
	}
	return finalData, nil
}

func (x *XConf) loadFile(file string) (map[string]interface{}, error) {
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("got error while read file:%s err:%w", file, err)
	}
	unmarshal := GetDecodeFunc(filepath.Ext(file))
	data := make(map[string]interface{})
	if err = unmarshal(bytes, data); err != nil {
		return nil, fmt.Errorf("load %s with error %s", file, err.Error())
	}
	// 检测是否继承文件
	inherit, ok := data[MetaKeyInheritFiles]
	if !ok {
		inherit, ok = data[MetaKeyInheritFilesDeprecatedFromGoconf]
	}
	if !ok {
		return data, nil
	}
	basePath := filepath.Dir(file) + "/"
	var inheritData map[string]interface{}
	var inheritErr error
	var inheritFiles []string
	switch it := inherit.(type) {
	case string:
		inheritFiles = append(inheritFiles, basePath+it)
	case []interface{}:
		for _, fi := range it {
			inheritFiles = append(inheritFiles, basePath+fi.(string))
		}
	}
	inheritData, inheritErr = x.loadFiles(inheritFiles...)

	if inheritErr != nil {
		return data, fmt.Errorf("got error:%w while inherit file:%v", inheritErr, inheritFiles)
	}
	// 对 xconf_ignore_fields 做特殊处理：base 与 inherit 的白名单做并集去重
	// 而不是按 mergeMap 的"叶子节点直接覆盖"语义
	mergedIgnore := mergeIgnoreFields(
		toStringSliceLoose(inheritData[MetaKeyIgnoreFields]),
		toStringSliceLoose(data[MetaKeyIgnoreFields]),
	)
	// 本文件内容覆盖继承而来的数据
	mergeErr := x.mergeMap("file:"+file, "inherited:"+strings.Join(inheritFiles, ","), data, inheritData)
	if mergeErr != nil {
		return data, fmt.Errorf("got error:%w while merge file:%v", mergeErr, inheritFiles)
	}
	// 把并集后的白名单写回（覆盖 mergeMap 直接覆盖的结果）
	if len(mergedIgnore) > 0 {
		inheritData[MetaKeyIgnoreFields] = mergedIgnore
	}
	return inheritData, nil
}
