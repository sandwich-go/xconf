package xconf

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"github.com/sandwich-go/xconf/kv"
)

type ProviderType string

const ProviderETCD ProviderType = "etcd"
const ProviderLocalFile ProviderType = "file"

type kvLoader struct {
	kv.Loader
	confPath string
}

type OnFieldUpdated func(from, to interface{})

// WatchFieldPath 关注特定的字段变化
func (x *XConf) WatchFieldPath(fieldPath string, changed OnFieldUpdated) {
	if _, ok := x.fieldPathInfoMap[fieldPath]; !ok {
		panic(fmt.Sprintf("field path:%s not found,valid ones:%v", fieldPath, x.keysList()))
	}
	x.mapOnFieldUpdated[fieldPath] = changed
}

func (x *XConf) autoRecover(tag string, f func()) {
	defer func() {
		if reason := recover(); reason != nil {
			x.cc.LogWarning(fmt.Sprintf("autoRecover %s panic with reason:%v retry after 5 second", tag, reason))
			time.Sleep(time.Second * 5)
			go x.autoRecover(tag, f)
		}
	}()
	f()
}

// WatchUpdate confPath不会自动绑定env value,如果需要watch的路径与环境变量相关，先通过ParseEnvValue自行解析替换处理错误
func (x *XConf) WatchUpdate(confPath string, loader kv.Loader) {
	k := &kvLoader{
		confPath: confPath,
		Loader:   loader,
	}
	x.kvs = append(x.kvs, k)
	if ow, ok := loader.(interface {
		CheckOnWatchError(watchError kv.WatchError)
	}); ok {
		ow.CheckOnWatchError(func(name string, confPath string, err error) {
			x.cc.LogWarning(fmt.Sprintf("name:%s confPath:%s watch got error:%s", name, confPath, err))
		})
	}
	go x.autoRecover(fmt.Sprintf("watch_with_%s_%s"+loader.Name(), confPath), func() {
		k.Watch(context.TODO(), confPath, x.onContentChanged)
	})
}

func (x *XConf) notifyChanged() error {
	latest, err := x.Latest()
	if err != nil {
		return err
	}
	select {
	case <-x.updated:
	default:
	}
	x.updated <- latest
	// 自动更新
	x.atomicSetFunc(latest)
	for k, v := range x.changes.Changed {
		notify, ok := x.mapOnFieldUpdated[k]
		if !ok {
			continue
		}
		notify(v.From, v.To)
	}
	x.changes.Changed = make(map[string]*Values)
	return nil
}

func (x *XConf) onContentChanged(name string, confPath string, content []byte) {
	x.dynamicUpdate.Lock()
	defer x.dynamicUpdate.Unlock()
	x.cc.LogDebug(fmt.Sprintf("got update:%s", confPath))
	defer func() {
		if reason := recover(); reason == nil {
			x.cc.LogWarning(fmt.Sprintf("onContentChanged with name:%s path:%s succ.", name, confPath))
		} else {
			x.cc.LogWarning(fmt.Sprintf("onContentChanged with name:%s path:%s reason:%v content:%s", name, confPath, reason, string(content)))
		}
	}()
	unmarshal := GetDecodeFunc(filepath.Ext(confPath))
	data := make(map[string]interface{})
	err := unmarshal(content, data)

	panicErrWithWrap(err, "unmarshal_error(%v) ", err)
	panicErr(x.mergeToDest(confPath, data))
	x.notifyChanged()
}
