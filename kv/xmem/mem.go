package xmem

import (
	"context"
	"errors"
	"sync"

	"github.com/sandwich-go/xconf/kv"
)

// Loader mem Loader
type Loader struct {
	*kv.Common
	dataMutex       sync.Mutex
	data            map[string][]byte
	onContentChange kv.ContentChange
	confPath        string
}

// New new mem Loader
func New(opts ...kv.Option) (p *Loader, err error) {
	x := &Loader{data: make(map[string][]byte)}
	x.Common = kv.New("mem", x, opts...)
	return x, nil
}

// CloseImplement 实现common.loaderImplement.CloseImplement
func (p *Loader) CloseImplement(ctx context.Context) error { return nil }

// GetImplement 实现common.loaderImplement.GetImplement
func (p *Loader) GetImplement(ctx context.Context, confPath string) ([]byte, error) {
	p.dataMutex.Lock()
	defer p.dataMutex.Unlock()
	v, ok := p.data[confPath]
	if !ok {
		return nil, errors.New("not found with path:" + confPath)
	}
	return v, nil
}

// WatchImplement 实现common.loaderImplement.WatchImplement
func (p *Loader) WatchImplement(ctx context.Context, confPath string, onContentChange kv.ContentChange) {
	p.onContentChange = onContentChange
	p.confPath = confPath
}

// Set 设定数据
func (p *Loader) Set(confPath string, data []byte) {
	p.dataMutex.Lock()
	defer p.dataMutex.Unlock()
	changed := p.Common.IsChanged(confPath, data)
	if !changed {
		return
	}
	p.data[confPath] = data
	if confPath == p.confPath && p.onContentChange != nil {
		p.onContentChange(p.Name(), confPath, data)
	}
}
