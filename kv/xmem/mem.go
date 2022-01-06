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
	dataMutex sync.Mutex
	data      map[string][]byte
	onChanged map[string]kv.ContentChange
}

// New new mem Loader
func New(opts ...kv.Option) (p *Loader, err error) {
	x := &Loader{data: make(map[string][]byte), onChanged: make(map[string]kv.ContentChange)}
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
	p.dataMutex.Lock()
	defer p.dataMutex.Unlock()
	p.onChanged[confPath] = onContentChange
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
	if f, ok := p.onChanged[confPath]; ok {
		f(p.Name(), confPath, data)
	}
}
