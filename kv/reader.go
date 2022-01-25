package kv

import (
	"bytes"
	"context"
	"io"
	"time"
)

type reader struct {
	Getter
	p       string
	timeout time.Duration
	reader  io.Reader
}

// NewReader 返回一个kv reader
func NewReader(getter Getter, path string, timeout time.Duration) io.Reader {
	return &reader{Getter: getter, p: path, timeout: timeout}
}

func (r *reader) do() error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()
	bb, err := r.Getter.Get(ctx, r.p)
	if err != nil {
		return err
	}
	r.reader = bytes.NewReader(bb)
	return nil
}

// Read 实现io.Reader接口
func (r *reader) Read(p []byte) (n int, err error) {
	if r.reader == nil {
		err = r.do()
		if err != nil {
			return 0, err
		}
	}
	return r.reader.Read(p)
}
