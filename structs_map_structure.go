package xconf

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/sandwich-go/xconf/xfield"
)

// Struct Struct类型定义
type Struct struct {
	raw                 interface{}
	value               reflect.Value
	tagName             string
	tagNameDefaultValue string
	fieldTagConvertor   FieldTagConvertor
}

// NewStruct 构造Struct类型
func NewStruct(s interface{}, tagName, tagNameDefaultValue string, ff FieldTagConvertor) *Struct {
	return &Struct{
		raw:                 s,
		value:               strctVal(s),
		tagName:             tagName,
		tagNameDefaultValue: tagNameDefaultValue,
		fieldTagConvertor:   ff,
	}
}

func strctVal(s interface{}) reflect.Value {
	v := reflect.ValueOf(s)
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		panic("not struct")
	}
	return v
}

// StructFieldPathInfo field信息
type StructFieldPathInfo struct {
	TagListXConf  xfield.TagList
	Tag           reflect.StructTag
	FieldName     string
	FieldNameList []string
	DefaultGot    bool
	DefaultString string
}

// Map 返回数据及字段类型信息
func (s *Struct) Map() (map[string]interface{}, map[string]StructFieldPathInfo) {
	out := make(map[string]interface{})
	outPath := make(map[string]StructFieldPathInfo)
	s.fillMapStructure(out, outPath, "", nil)
	return out, outPath
}

func (s *Struct) fillMapStructure(out map[string]interface{}, outPath map[string]StructFieldPathInfo, prefix string, fieldNames []string) {
	if out == nil || outPath == nil {
		return
	}
	if prefix != "" {
		prefix += DefaultKeyDelim
	}

	fields := s.structFields()
	for _, field := range fields {
		name := field.Name
		val := s.value.FieldByName(name)
		tagVal, tagOpts := xfield.ParseTag(field.Tag.Get(s.tagName))
		if tagVal == "" {
			name = s.fieldTagConvertor(name)
		} else {
			name = tagVal
		}
		fullKey := prefix + name
		fileNameNow := append(fieldNames, field.Name)
		// TODO 指针类型且数据为nil,自动构造一个默认值便于后续分析，conf中的sub最好不要使用指针类型？
		if val.Kind() == reflect.Ptr && val.IsNil() {
			val = reflect.New(val.Type().Elem())
		}
		finalVal := s.nested(val, outPath, fullKey, fileNameNow)
		out[name] = finalVal
		defaultVal, defaultValGot := field.Tag.Lookup(s.tagNameDefaultValue)
		noConf := strings.HasPrefix(fullKey, "-") || strings.HasSuffix(fullKey, "-") || strings.Contains(fullKey, ".-.")
		if !noConf {
			outPath[fullKey] = StructFieldPathInfo{
				DefaultString: defaultVal,
				DefaultGot:    defaultValGot,
				FieldNameList: fileNameNow,
				TagListXConf:  tagOpts,
				Tag:           field.Tag,
				FieldName:     field.Name,
			}
		}
	}
}
func (s *Struct) structFields() []reflect.StructField {
	t := s.value.Type()
	var f []reflect.StructField
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		// 跳过私有字段
		if field.PkgPath != "" {
			continue
		}
		// 跳过omitted字段
		if tag := field.Tag.Get(s.tagName); tag == "-" {
			continue
		}
		f = append(f, field)
	}

	return f
}

// nested retrieves recursively all types for the given value and returns the
// nested value.
func (s *Struct) nested(val reflect.Value, outPath map[string]StructFieldPathInfo, prefix string, fieldNames []string) interface{} {
	var finalVal interface{}

	v := reflect.ValueOf(val.Interface())
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Struct:
		n := NewStruct(val.Interface(), s.tagName, s.tagNameDefaultValue, s.fieldTagConvertor)
		m := make(map[string]interface{})
		n.fillMapStructure(m, outPath, prefix, fieldNames)

		// do not add the converted value if there are no exported fields, ie:
		// time.Time
		if len(m) == 0 {
			finalVal = val.Interface()
		} else {
			finalVal = m
		}
	case reflect.Map:
		if len(val.MapKeys()) == 0 {
			finalVal = val.Interface()
		} else {
			m := make(map[string]interface{}, val.Len())
			for _, k := range val.MapKeys() {
				m[fmt.Sprintf("%v", k)] = s.nested(val.MapIndex(k), outPath, prefix, fieldNames)
			}
			finalVal = m
		}
	case reflect.Slice, reflect.Array:
		if val.Type().Kind() == reflect.Interface {
			finalVal = val.Interface()
			break
		}
		if val.Type().Elem().Kind() != reflect.Struct &&
			!(val.Type().Elem().Kind() == reflect.Ptr &&
				val.Type().Elem().Elem().Kind() == reflect.Struct) {
			finalVal = val.Interface()
			break
		}
		slices := make([]interface{}, val.Len())
		for x := 0; x < val.Len(); x++ {
			slices[x] = s.nested(val.Index(x), outPath, prefix, fieldNames)
		}
		finalVal = slices
	default:
		finalVal = val.Interface()
	}

	return finalVal
}
