package xconf

import "strings"

// SensitiveDataRedactedValue 敏感配置输出时使用的占位值。
const SensitiveDataRedactedValue = "[REDACTED]"

var sensitiveDataKeywords = [...]string{
	"token",
	"secret",
	"password",
}

// shouldRedactSensitiveData 判断字段路径是否需要脱敏。
// 匹配不区分大小写；通过 WithSensitiveDataRedaction(false) 可关闭。
func (x *XConf) shouldRedactSensitiveData(fieldPath string) bool {
	if x == nil || x.cc == nil || !x.cc.SensitiveDataRedaction {
		return false
	}
	return isSensitiveDataFieldPath(fieldPath)
}

func isSensitiveDataFieldPath(fieldPath string) bool {
	fieldName := strings.ToLower(fieldPath)
	if index := strings.LastIndex(fieldName, DefaultKeyDelim); index >= 0 {
		fieldName = fieldName[index+len(DefaultKeyDelim):]
	}
	for _, keyword := range sensitiveDataKeywords {
		if strings.Contains(fieldName, keyword) {
			return true
		}
	}
	return fieldName == "key" ||
		fieldName == "keys" ||
		strings.HasSuffix(fieldName, "_key") ||
		strings.HasSuffix(fieldName, "_keys") ||
		strings.Contains(fieldName, "_key_") ||
		strings.Contains(fieldName, "_keys_")
}

func (x *XConf) redactSensitiveMap(data map[string]interface{}) map[string]interface{} {
	return x.redactSensitiveMapWithPrefix(data, "")
}

func (x *XConf) redactSensitiveMapForLog(data map[string]interface{}) map[string]interface{} {
	if x == nil || x.cc == nil || !x.cc.SensitiveDataRedaction {
		return data
	}
	return x.redactSensitiveMap(data)
}

func (x *XConf) redactSensitiveMapWithPrefix(data map[string]interface{}, prefix string) map[string]interface{} {
	redacted := make(map[string]interface{}, len(data))
	for key, value := range data {
		fieldPath := joinSensitiveFieldPath(prefix, key)
		redacted[key] = x.redactSensitiveValue(fieldPath, value)
	}
	return redacted
}

func (x *XConf) redactSensitiveValue(fieldPath string, value interface{}) interface{} {
	if isSensitiveDataFieldPath(fieldPath) {
		return SensitiveDataRedactedValue
	}

	switch typedValue := value.(type) {
	case map[string]interface{}:
		return x.redactSensitiveMapWithPrefix(typedValue, fieldPath)
	case []interface{}:
		redacted := make([]interface{}, len(typedValue))
		for i, item := range typedValue {
			redacted[i] = x.redactSensitiveValue(fieldPath, item)
		}
		return redacted
	default:
		return value
	}
}

func joinSensitiveFieldPath(prefix string, key string) string {
	if prefix == "" {
		return key
	}
	if key == "" {
		return prefix
	}
	return prefix + DefaultKeyDelim + key
}
