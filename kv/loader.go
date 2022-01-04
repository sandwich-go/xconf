package kv

import (
	"context"
	"crypto/md5"
	"sync"
)

type Loader interface {
	Get(ctx context.Context, confPath string) ([]byte, error)
	Watch(ctx context.Context, confPath string, onContentChange func(string, []byte))
	Name() string
	Close(ctx context.Context) error
}

type loaderImplement interface {
	CloseImplement(ctx context.Context) error
	GetImplement(ctx context.Context, confPath string) ([]byte, error)
	WatchImplement(ctx context.Context, confPath string, onContentChange func(string, []byte))
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
	return c.implement.GetImplement(ctx, confPath)
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

func (c *Common) Watch(ctx context.Context, confPath string, onContentChange func(string, []byte)) {
	c.implement.WatchImplement(ctx, confPath, onContentChange)
}
func (c *Common) Name() string { return c.name }
