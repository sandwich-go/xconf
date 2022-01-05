package secconf

import (
	"bytes"
	"fmt"
	"io"
)

type secReader struct {
	wrapped io.Reader
	cached  io.Reader
	cf      Codec
}

func (s *secReader) Read(p []byte) (n int, err error) {
	if s.cached == nil {
		all, err := io.ReadAll(s.wrapped)
		if err != nil {
			return 0, fmt.Errorf("got err:%w while ReadAll from wrapped io.Reader", err)
		}
		// 解密
		allDecode, err := s.cf.Apply(all)
		if err != nil {
			return 0, fmt.Errorf("got err:%w while Decode", err)
		}
		s.cached = bytes.NewReader(allDecode)
	}
	return s.cached.Read(p)
}

func Reader(reader io.Reader, cf Codec) io.Reader { return &secReader{wrapped: reader, cf: cf} }
