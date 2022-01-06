package kv

import (
	"context"
	"crypto/md5"
	"fmt"
	"sync"
)

// todo Loder实现Reader接口完全对接到io.Reader，将远程的首次加载流程直接对接到xconf的WithReader
// Loader kv加载基础接口
type Loader interface {
	// Name Loader名称
	Name() string
	// Get 主动获取指定confPath的数据
	Get(ctx context.Context, confPath string) ([]byte, error)
	// Watch Watch指定的confPath，数据发生变化会回调onContentChange
	Watch(ctx context.Context, confPath string, onContentChange ContentChange)
	// Close 关闭
	Close(ctx context.Context) error
}

// loaderImplement Loader特有逻辑对接，基础逻辑落实在Common对应接口中
type loaderImplement interface {
	CloseImplement(ctx context.Context) error
	GetImplement(ctx context.Context, confPath string) ([]byte, error)
	WatchImplement(ctx context.Context, confPath string, onContentChange ContentChange)
}

// Common 基础的kv.Loader实现，实现基础逻辑，自定义的Loader对接到loaderImplement接口即可
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

// IsChanged 指定的配置是否发生变化
func (c *Common) IsChanged(confPath string, data []byte) bool {
	c.Lock()
	defer c.Unlock()
	hash := md5.New()
	hash.Write(data)
	md5Str := string(hash.Sum(nil))
	if v, ok := c.fileMap[confPath]; ok && v == md5Str {
		return false
	}
	c.fileMap[confPath] = md5Str
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
