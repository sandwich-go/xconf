package secconf

import "errors"

type Codec interface {
	Apply(data []byte) ([]byte, error)
}

type CodecFuncChain []CodecFunc

func (c CodecFuncChain) Apply(data []byte) (out []byte, err error) {
	out = data
	for _, v := range c {
		out, err = v.Apply(out)
		if err != nil {
			return
		}
	}
	return
}

type CodecFunc func(data []byte) ([]byte, error)

func (f CodecFunc) Apply(data []byte) ([]byte, error) { return f(data) }

func newInvalidCodec(name string) Codec {
	return CodecFunc(func(data []byte) ([]byte, error) {
		return data, errors.New("invalid codec func for " + name)
	})
}

var decoderRegister = make(map[string]Codec)
var encoderRegister = make(map[string]Codec)

func RegisterDecoder(name string, cf Codec) { decoderRegister[name] = cf }
func RegisterEncoder(name string, cf Codec) { encoderRegister[name] = cf }

func Decoder(name string) Codec {
	if cf, ok := decoderRegister[name]; ok {
		return cf
	}
	return newInvalidCodec(name)
}

func Encoder(name string) Codec {
	if cf, ok := encoderRegister[name]; ok {
		return cf
	}
	return newInvalidCodec(name)
}
