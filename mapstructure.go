package xconf

import (
	"reflect"

	"github.com/mitchellh/mapstructure"
)

// DecoderConfigOption A DecoderConfigOption can be passed to Unmarshal to configure mapstructure.DecoderConfig options
type DecoderConfigOption = func(*mapstructure.DecoderConfig)

// defaultDecoderConfig returns default mapsstructure.DecoderConfig with suppot
// of time.Duration values & string slices
func (x *XConf) defaultDecoderConfig(output interface{}) *mapstructure.DecoderConfig {
	c := &mapstructure.DecoderConfig{
		Metadata:         nil,
		Result:           output,
		WeaklyTypedInput: true,
		ZeroFields:       true,
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			StringAlias(x.cc.StringAlias, x.cc.StringAliasFunc),
			parseEnvVarStringHookFunc(x.cc.EnvBindShouldErrorWhenFailed),
			mapstructure.StringToTimeDurationHookFunc(),
			mapstructure.StringToSliceHookFunc(","),
		),
	}
	return c
}

func StringAlias(alias map[string]string, aliasFunc map[string]func(string) string) mapstructure.DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}
		if v, ok := alias[data.(string)]; ok {
			return v, nil
		}
		if v, ok := aliasFunc[data.(string)]; ok {
			return v(data.(string)), nil
		}
		return data, nil
	}
}

func parseEnvVarStringHookFunc(envBindShouldErrorWhenFailed bool) mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}
		ret, err := ParseEnvValue(data.(string))
		if !envBindShouldErrorWhenFailed {
			err = nil
		}
		return ret, err
	}
}

// A wrapper around mapstructure.Decode that mimics the WeakDecode functionality
func decode(input interface{}, config *mapstructure.DecoderConfig) error {
	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}
	// return Copy(config.Result, input)
	return decoder.Decode(input)
}
