package xconf

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/sandwich-go/xconf/kv"
	"github.com/sandwich-go/xconf/xutil"
)

type kvLoader struct {
	kv.Loader
	confPath string
}

// OnFieldUpdated 字段发生变化方法签名
type OnFieldUpdated func(fieldPath string, from, to interface{})

// WatchFieldPath 关注特定的字段变化
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
	if ow, ok := loader.(interface {
		CheckOnWatchError(watchError kv.WatchError)
	}); ok {
		ow.CheckOnWatchError(func(name string, confPath string, err error) {
			x.cc.LogWarning(fmt.Sprintf("name:%s confPath:%s watch got error:%s", name, confPath, err))
		})
	}
	// 需要Loader自行维护异步逻辑
	k.Watch(context.TODO(), k.confPath, x.onContentChanged)
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
	for k, v := range x.changes.changed {
		notify, ok := x.mapOnFieldUpdated[k]
		if !ok {
			continue
		}
		notify(v.fieldPath, v.from, v.to)
	}
	x.changes.changed = make(map[string]*fieldValues)
	return nil
}

func (x *XConf) onContentChanged(name string, confPath string, content []byte) {
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

	xutil.PanicErrWithWrap(err, "unmarshal_error(%v) ", err)
	xutil.PanicErr(x.commonUpdateAndNotify(func() error {
		return x.mergeToDest(confPath, data)
	}))
}
