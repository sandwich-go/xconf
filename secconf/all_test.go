package secconf

import (
	"testing"
)

var encodingTests = []struct {
	in, out string
}{
	{"secret", "secret"},
}

func TestEncoding(t *testing.T) {
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
