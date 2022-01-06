package xfile

import (
	"context"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/sandwich-go/xconf/kv"
	"gopkg.in/fsnotify.v1"
)

// Loader file Loader
type Loader struct {
	watcher *fsnotify.Watcher
	*kv.Common
}

// New new file Loader
func New(opts ...kv.Option) (p kv.Loader, err error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("got error:%v when fsnotify.NewWatcher", err)
	}
	x := &Loader{watcher: watcher}
	x.Common = kv.New("file", x, opts...)
	return x, nil
}

// CloseImplement 实现common.loaderImplement.CloseImplement
func (p *Loader) CloseImplement(ctx context.Context) error { return p.watcher.Close() }

// GetImplement 实现common.loaderImplement.GetImplement
func (p *Loader) GetImplement(ctx context.Context, confPath string) ([]byte, error) {
	return ioutil.ReadFile(confPath)
}

// WatchImplement 实现common.loaderImplement.WatchImplement
func (p *Loader) WatchImplement(ctx context.Context, confPath string, onContentChange kv.ContentChange) {
	go func(pin *Loader, oc kv.ContentChange) {
		watched := false
		for {
			select {
			case <-pin.Done:
				return
			default:
			}
			if !watched {
				if err := pin.watcher.Add(confPath); err != nil {
					pin.CC.OnWatchError(pin.Name(), confPath, err)
				}
			}
			select {
			case event := <-pin.watcher.Events:
				if (event.Op&fsnotify.Write) == fsnotify.Write || (event.Op&fsnotify.Create) == fsnotify.Create {
					confPathChanged := strings.ReplaceAll(event.Name, "\\", "/")
					if b, err := pin.Get(ctx, confPathChanged); err == nil {
						if pin.IsChanged(confPathChanged, b) {
							oc(pin.Name(), confPathChanged, b)
						}
					}
				}
			case err := <-pin.watcher.Errors:
				select {
				case <-pin.Done:
					return
				default:
				}
				pin.CC.OnWatchError(pin.Name(), confPath, err)
			}
		}
	}(p, onContentChange)
}
