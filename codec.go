package xconf

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v2"
)

type ConfigType string

const ConfigTypeTOML ConfigType = ".toml"
const ConfigTypeJSON ConfigType = ".json"
const ConfigTypeYAML ConfigType = ".yaml"

type DecodeFunc func([]byte, map[string]interface{}) error
type EncodeFunc func(v map[string]interface{}) ([]byte, error)

var decoderMap = make(map[string]DecodeFunc)
var encorderMap = make(map[string]EncodeFunc)

func extClean(e string) string {
	e = strings.ToLower(e)
	if !strings.HasPrefix(e, ".") {
		e = "." + e
	}
	return e
}

func RegisterFileCodec(ct ConfigType, d DecodeFunc, e EncodeFunc) {
	ext := string(ct)
	ext = extClean(ext)
	decoderMap[ext] = d
	encorderMap[ext] = e
}

func GetDecodeFunc(ext string) DecodeFunc {
	got, ok := decoderMap[extClean(ext)]
	if !ok {
		got = loopDocode
	}
	return got
}

func loopDocode(buf []byte, data map[string]interface{}) error {
	var errs []string
	for name, v := range decoderMap {
		err := v(buf, data)
		if err != nil {
			errs = append(errs, fmt.Sprintf(" #codec:%s err:%s", name, err.Error()))
			continue
		}
		_, ok := data["inherit_files"]
		if ok {
			return fmt.Errorf("codec:%s load from bytes do not support inherit_files", name)
		}
		return nil
	}
	return fmt.Errorf("codec not found, %s", strings.Join(errs, " "))
}

func GetEncodeFunc(ext string) EncodeFunc {
	got, ok := encorderMap[extClean(ext)]
	if !ok {
		got = func(v map[string]interface{}) ([]byte, error) {
			return nil, fmt.Errorf("can not find encoder with key:%s", ext)
		}
	}
	return got
}

func init() {
	RegisterFileCodec(ConfigTypeTOML, func(bytes []byte, data map[string]interface{}) error {
		_, err := toml.Decode(string(bytes), &data)
		return err
	}, func(v map[string]interface{}) ([]byte, error) {
		var bb bytes.Buffer
		encoder := toml.NewEncoder(&bb)
		err := encoder.Encode(v)
		return bb.Bytes(), err
	})
	RegisterFileCodec(ConfigTypeJSON, func(bytes []byte, data map[string]interface{}) error {
		return json.Unmarshal(bytes, &data)
	}, func(v map[string]interface{}) ([]byte, error) {
		return json.MarshalIndent(v, "", "    ")
	})
	RegisterFileCodec(ConfigTypeYAML, func(bytes []byte, data map[string]interface{}) error {
		return yaml.Unmarshal(bytes, &data)
	}, func(v map[string]interface{}) ([]byte, error) {
		return yaml.Marshal(v)
	})
}
