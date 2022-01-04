package xconf

import (
	"context"
	"fmt"
	"path/filepath"

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

func (x *XConf) WatchFieldPath(fieldPath string, changed OnFieldUpdated) {
	if _, ok := x.fieldPathInfoMap[fieldPath]; !ok {
		panic(fmt.Sprintf("field path:%s not found,valid ones:%v", fieldPath, x.keysList()))
	}
	x.mapOnFieldUpdated[fieldPath] = changed
}

// WatchUpdate confPath不会自动绑定env value,如果需要watch的路径与环境变量相关，先通过ParseEnvValue自行解析替换处理错误
func (x *XConf) WatchUpdate(confPath string, loader kv.Loader) {
	k := &kvLoader{
		confPath: confPath,
		Loader:   loader,
	}
	x.kvs = append(x.kvs, k)
	go func() {
		k.Watch(context.TODO(), confPath, x.onContentChanged)
	}()
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

func (x *XConf) onContentChanged(confPath string, content []byte) {
	x.dynamicUpdate.Lock()
	defer x.dynamicUpdate.Unlock()
	x.cc.LogDebug(fmt.Sprintf("got update:%s", confPath))
	defer func() {
		if reason := recover(); reason == nil {
			x.cc.LogWarning(fmt.Sprintf("onContentChanged with path:%s succ.", confPath))
		} else {
			x.cc.LogWarning(fmt.Sprintf("onContentChanged with path:%s reason:%v content:%s", confPath, reason, string(content)))
		}
	}()
	unmarshal := GetDecodeFunc(filepath.Ext(confPath))
	data := make(map[string]interface{})
	err := unmarshal(content, data)

	panicErrWithWrap(err, "unmarshal_error(%v) ", err)
	panicErr(x.mergeToDest(confPath, data))
	x.notifyChanged()
}
