package secconf

import (
	"bytes"
	"testing"
)

var encodingTests = []struct {
	in, out string
}{
	{"secret", "secret"},
}

func TestEncodingStandard(t *testing.T) {
	xxteaKey := []byte("xconf")
	encoder := StandardChainEncode(NewEncoderXXTEA(xxteaKey))
	decoder := StandardChainDecode(NewDecoderXXTEA(xxteaKey))

	for _, tt := range encodingTests {
		encoded, err := encoder.Apply([]byte(tt.in))
		if err != nil {
			t.Errorf(err.Error())
		}
		decoded, err := decoder.Apply(encoded)
		if err != nil {
			t.Errorf(err.Error())
		}
		if tt.out != string(decoded) {
			t.Errorf("want %s, got %s", tt.out, decoded)
		}
	}
}

func TestEncodingWithPrefix(t *testing.T) {
	magicPrefix := []byte("$xconf$")
	xxteaKey := []byte("xconf")
	encoder := CodecFrom(EncoderGZip, NewEncoderXXTEA(xxteaKey), EncoderBase64, NewEncoderMagic(magicPrefix))
	decoder := CodecFrom(NewDecoderMagic(magicPrefix), DecoderBase64, NewDecoderXXTEA(xxteaKey), DecoderGZip)

	for _, tt := range encodingTests {
		encoded, err := encoder.Apply([]byte(tt.in))
		if err != nil {
			t.Errorf(err.Error())
		}
		if !bytes.HasPrefix(encoded, magicPrefix) {
			t.Errorf("encoded data should have magic prefix:" + string(magicPrefix))
		}
		decoded, err := decoder.Apply(encoded)
		if err != nil {
			t.Errorf(err.Error())
		}
		if tt.out != string(decoded) {
			t.Errorf("want %s, got %s", tt.out, decoded)
		}
	}
}
