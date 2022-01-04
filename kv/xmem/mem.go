package xmem

import (
	"context"
	"errors"
	"sync"

	"github.com/sandwich-go/xconf/kv"
)

type Loader struct {
	*kv.Common
	dataMutex       sync.Mutex
	data            map[string][]byte
	onContentChange func(string, []byte)
	confPath        string
}

func New(opts ...kv.Option) (p *Loader, err error) {
	x := &Loader{data: make(map[string][]byte)}
	x.Common = kv.New("mem", x, opts...)
	return x, nil
}
func (p *Loader) CloseImplement(ctx context.Context) error { return nil }
func (p *Loader) GetImplement(ctx context.Context, confPath string) ([]byte, error) {
	p.dataMutex.Lock()
	defer p.dataMutex.Unlock()
	v, ok := p.data[confPath]
	if !ok {
		return nil, errors.New("not found with path:" + confPath)
	}
	return v, nil
}

func (p *Loader) WatchImplement(ctx context.Context, confPath string, onContentChange func(string, []byte)) {
	p.onContentChange = onContentChange
	p.confPath = confPath
}

func (p *Loader) Set(confPath string, data []byte) {
	p.dataMutex.Lock()
	defer p.dataMutex.Unlock()
	changed := p.Common.IsChanged(confPath, data)
	if !changed {
		return
	}
	p.data[confPath] = data
	if confPath == p.confPath && p.onContentChange != nil {
		p.onContentChange(confPath, data)
	}
}
