package xconf

import (
	"reflect"

	"github.com/mitchellh/mapstructure"
)

// A DecoderConfigOption can be passed to Unmarshal to configure mapstructure.DecoderConfig options
type DecoderConfigOption func(*mapstructure.DecoderConfig)

// defaultDecoderConfig returns default mapsstructure.DecoderConfig with suppot
// of time.Duration values & string slices
func (x *XConf) defaultDecoderConfig(output interface{}) *mapstructure.DecoderConfig {
	c := &mapstructure.DecoderConfig{
		Metadata:         nil,
		Result:           output,
		WeaklyTypedInput: true,
		ZeroFields:       true,
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			ParseEnvVarStringHookFunc(x.cc.ErrEnvBindNotExistWithoutDefault),
			mapstructure.StringToTimeDurationHookFunc(),
			mapstructure.StringToSliceHookFunc(","),
		),
	}
	return c
}

func ParseEnvVarStringHookFunc(errEnvBindNotExistWithoutDefault bool) mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}

		return ParseEnvValue(data.(string), errEnvBindNotExistWithoutDefault)
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
