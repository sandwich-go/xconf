package secconf

import (
	"bytes"
	"fmt"
	"io"
)

type secReader struct {
	wrapped       io.Reader
	cached        io.Reader
	secertKeyring io.Reader
}

func (s *secReader) Read(p []byte) (n int, err error) {
	if s.cached == nil {
		all, err := io.ReadAll(s.wrapped)
		if err != nil {
			return 0, fmt.Errorf("got err:%w while ReadAll from wrapped io.Reader", err)
		}
		// 解密
		allDecode, err := Decode(all, s.secertKeyring)
		if err != nil {
			return 0, fmt.Errorf("got err:%w while Decode", err)
		}
		s.cached = bytes.NewReader(allDecode)
	}
	return s.cached.Read(p)
}

func Reader(reader io.Reader, secertKeyring io.Reader) io.Reader {
	return &secReader{wrapped: reader, secertKeyring: secertKeyring}
}
