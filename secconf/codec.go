package secconf

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

// gzip => encrypt => base64
func StandardChainEncode(encrypt CodecFunc) Codec {
	return CodecFuncChain([]CodecFunc{
		EncoderGZip,
		encrypt,
		EncoderBase64,
	})
}

// base64 => decrypt => gzip
func StandardChainDecode(decrypt CodecFunc) Codec {
	return CodecFuncChain([]CodecFunc{
		DecoderBase64,
		decrypt,
		DecoderGZip,
	})
}
