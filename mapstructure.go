package xconf

import (
	"reflect"

	"github.com/sandwich-go/mapstructure"
)

// DecoderConfigOption A DecoderConfigOption can be passed to Unmarshal to configure mapstructure.DecoderConfig options
type DecoderConfigOption = func(*mapstructure.DecoderConfig)

// defaultDecoderConfig returns default mapsstructure.DecoderConfig with suppot
// of time.Duration values & string slices
func (x *XConf) defaultDecoderConfig(output interface{}) *mapstructure.DecoderConfig {
	c := &mapstructure.DecoderConfig{
		Metadata: nil,
		Result:   output,
		// If WeaklyTypedInput is true, the decoder will make the following
		// "weak" conversions:
		//
		//   - bools to string (true = "1", false = "0")
		//   - numbers to string (base 10)
		//   - bools to int/uint (true = 1, false = 0)
		//   - strings to int/uint (base implied by prefix)
		//   - int to bool (true if value != 0)
		//   - string to bool (accepts: 1, t, T, TRUE, true, True, 0, f, F,
		//     FALSE, false, False. Anything else is an error)
		//   - empty array = empty map and vice versa
		//   - negative numbers to overflowed uint values (base 10)
		//   - slice of maps to a merged map
		//   - single values are converted to slices if required. Each
		//     element is weakly decoded. For example: "4" can become []int{4}
		//     if the target type is an int slice.
		//
		WeaklyTypedInput: true,
		// ZeroFields, if set to true, will zero fields before writing them.
		// For example, a map will be emptied before decoded values are put in
		// it. If this is false, a map will be merged.
		ZeroFields: true,
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
