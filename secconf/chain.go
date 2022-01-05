package secconf

// gip => encrypt => base64
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
