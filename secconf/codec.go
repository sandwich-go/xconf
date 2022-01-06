package secconf

// Codec 编解码器接口
type Codec interface {
	Apply(data []byte) ([]byte, error)
}

// CodecFunc 编辑吗函数类型
type CodecFunc func(data []byte) ([]byte, error)

// Apply CodecFunc 实现Codec接口
func (f CodecFunc) Apply(data []byte) ([]byte, error) { return f(data) }

// CodecFuncChain 编解码器链
type CodecFuncChain []CodecFunc

// Apply CodecFuncChain实现Codec接口
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

// CodecFrom 由给定的CodecFunc列表构造Codec
func CodecFrom(codec ...CodecFunc) Codec { return CodecFuncChain(codec) }

// StandardChainEncode 默认编码器链：gzip => encrypt => base64
func StandardChainEncode(encrypt CodecFunc) Codec {
	return CodecFuncChain([]CodecFunc{
		EncoderGZip,
		encrypt,
		EncoderBase64,
	})
}

// StandardChainDecode 默认解码器链：base64 => decrypt => gzip
func StandardChainDecode(decrypt CodecFunc) Codec {
	return CodecFuncChain([]CodecFunc{
		DecoderBase64,
		decrypt,
		DecoderGZip,
	})
}
