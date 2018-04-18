// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// 文件监控.
// 使用时需要注意的是，一旦一个文件被删除，那么对其的监控将会失效。
package gfsnotify

import (
    "errors"
    "github.com/fsnotify/fsnotify"
    "gitee.com/johng/gf/g/os/gfile"
    "gitee.com/johng/gf/g/os/gcache"
    "gitee.com/johng/gf/g/os/grpool"
    "gitee.com/johng/gf/g/util/gconv"
    "gitee.com/johng/gf/g/container/gmap"
    "gitee.com/johng/gf/g/container/glist"
    "gitee.com/johng/gf/g/container/gqueue"
)

// 监听管理对象
type Watcher struct {
    watcher    *fsnotify.Watcher        // 底层fsnotify对象
    events     *gqueue.Queue            // 过滤后的事件通知，不会出现重复事件
    eventCache *gcache.Cache            // 用于进行事件过滤，当同一监听文件在100ms内出现相同事件，则过滤
    closeChan  chan struct{}            // 关闭事件
    callbacks  *gmap.StringInterfaceMap // 监听的回调函数
}

// 监听事件对象
type Event struct {
    Path string // 文件绝对路径
    Op   Op     // 触发监听的文件操作
}

// 按位进行识别的操作集合
type Op uint32

const (
    CREATE Op = 1 << iota
    WRITE
    REMOVE
    RENAME
    CHMOD
)

// 全局监听对象，方便应用端调用
var watcher, _ = New()

// 添加对指定文件/目录的监听，并给定回调函数
func Add(path string, callback func(event *Event)) error {
    if watcher == nil {
        return errors.New("global watcher creating failed")
    }
    return watcher.Add(path, callback)
}

// 移除监听
func Remove(path string) error {
    if watcher == nil {
        return errors.New("global watcher creating failed")
    }
    return watcher.Remove(path)
}



// 创建监听管理对象
func New() (*Watcher, error) {
    if watch, err := fsnotify.NewWatcher(); err == nil {
        w := &Watcher {
            watcher    : watch,
            events     : gqueue.New(),
            eventCache : gcache.New(),
            closeChan  : make(chan struct{}, 1),
            callbacks  : gmap.NewStringInterfaceMap(),
        }
        w.startWatchLoop()
        w.startEventLoop()
        return w, nil
    } else {
        return nil, err
    }
}

// 关闭监听管理对象
func (w *Watcher) Close() {
    w.watcher.Close()
    w.events.Close()
    w.closeChan <- struct{}{}
}

// 添加对指定文件/目录的监听，并给定回调函数
func (w *Watcher) Add(path string, callback func(event *Event)) error {
    if !gfile.Exists(path) {
        return errors.New(path + " does not exist")
    }
    // 注册回调函数
    w.callbacks.LockFunc(func(m map[string]interface{}) {
        var result interface{}
        if v, ok := m[path]; !ok {
            result  = glist.New()
            m[path] = result
        } else {
            result = v
        }
        result.(*glist.List).PushBack(callback)
    })
    // 添加底层监听
    w.watcher.Add(path)
    return nil
}

// 移除监听
func (w *Watcher) Remove(path string) error {
    w.callbacks.Remove(path)
    return w.watcher.Remove(path)
}

// 监听循环
func (w *Watcher) startWatchLoop() {
    go func() {
        for {
            select {
                // 关闭事件
                case <- w.closeChan:
                    return

                // 监听事件
                case ev := <- w.watcher.Events:
                    if !w.eventCache.Lock(ev.Name + ":" + gconv.String(ev.Op), 100) {
                        continue
                    }
                    w.events.PushBack(&Event{
                        Path : ev.Name,
                        Op   : Op(ev.Op),
                    })

                //case err := <- w.watcher.Errors:
                //    log.Println("error : ", err);
                //    return
            }
        }
    }()
}

// 事件循环
func (w *Watcher) startEventLoop() {
    go func() {
        for {
            if v := w.events.PopFront(); v != nil {
                event := v.(*Event)
                // 如果是文件删除时间，判断该文件是否存在，如果存在，那么将此事件认为“假删除”，并重新添加监控
                if event.IsRemove() && gfile.Exists(event.Path){
                    w.watcher.Add(event.Path)
                    continue
                }
                if l := w.callbacks.Get(event.Path); l != nil {
                    grpool.Add(func() {
                        for _, v := range l.(*glist.List).FrontAll() {
                            v.(func(event *Event))(event)
                        }
                    })
                }
            } else {
                break
            }
        }
    }()
}