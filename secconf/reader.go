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

// Read secReader实现io.Reader接口
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

// Reader 用给定的Codec将给定的reader封装成新的io.Reader，数据读取流程中会经给定的Codec进行编解码
func Reader(reader io.Reader, cf Codec) io.Reader { return &secReader{wrapped: reader, cf: cf} }
