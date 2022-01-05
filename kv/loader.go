package kv

import (
	"context"
	"crypto/md5"
	"fmt"
	"sync"
)

// todo Loder实现Reader接口完全对接到io.Reader，将远程的首次加载流程直接对接到xconf的WithReader
type Loader interface {
	Get(ctx context.Context, confPath string) ([]byte, error)
	Watch(ctx context.Context, confPath string, onContentChange ContentChange)
	Name() string
	Close(ctx context.Context) error
}

type loaderImplement interface {
	CloseImplement(ctx context.Context) error
	GetImplement(ctx context.Context, confPath string) ([]byte, error)
	WatchImplement(ctx context.Context, confPath string, onContentChange ContentChange)
}

type Common struct {
	name string
	Done chan struct{}
	CC   *Options

	implement loaderImplement

	sync.Mutex
	fileMap map[string]string
}

func New(name string, implement loaderImplement, opts ...Option) *Common {
	return &Common{implement: implement, fileMap: make(map[string]string), CC: NewOptions(opts...)}
}

func (c *Common) Close(ctx context.Context) error {
	close(c.Done)
	return c.implement.CloseImplement(ctx)
}

func (c *Common) Get(ctx context.Context, confPath string) ([]byte, error) {
	data, err := c.implement.GetImplement(ctx, confPath)
	if err != nil {
		return nil, err
	}
	return c.decode(data)
}

func (c *Common) CheckOnWatchError(watchError WatchError) {
	if c.CC.OnWatchError == nil {
		c.CC.OnWatchError = watchError
	}
}

func (c *Common) IsChanged(name string, data []byte) bool {
	c.Lock()
	defer c.Unlock()
	hash := md5.New()
	hash.Write(data)
	md5Str := string(hash.Sum(nil))
	if v, ok := c.fileMap[name]; ok && v == md5Str {
		return false
	}
	c.fileMap[name] = md5Str
	return true
}

func (c *Common) decode(in []byte) ([]byte, error) {
	if c.CC.Decoder == nil {
		return in, nil
	}
	dataOut, err := c.CC.Decoder.Apply(in)
	if err != nil {
		return nil, fmt.Errorf("got error :%w while decode using secconf", err)
	}
	return dataOut, nil
}

func (c *Common) Watch(ctx context.Context, confPath string, onContentChange ContentChange) {
	if c.CC.OnWatchError == nil {
		c.CC.OnWatchError = func(string, string, error) {}
	}
	c.implement.WatchImplement(ctx, confPath, func(name string, confPath string, data []byte) {
		out, err := c.decode(data)
		if err == nil {
			onContentChange(name, confPath, out)
			return
		}
		c.CC.OnWatchError(name, confPath, fmt.Errorf("got error :%w while decode using secconf", err))
	})
}

func (c *Common) Name() string { return c.name }
