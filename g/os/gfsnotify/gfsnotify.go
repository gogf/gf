// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// 文件监控.
package gfsnotify

import (
    "errors"
    "github.com/fsnotify/fsnotify"
    "gitee.com/johng/gf/g/os/gfile"
    "gitee.com/johng/gf/g/os/grpool"
    "gitee.com/johng/gf/g/container/gmap"
    "gitee.com/johng/gf/g/container/glist"
)

// 监听管理对象
type Watcher struct {
    watcher   *fsnotify.Watcher        // 底层fsnotify对象
    closeChan chan struct{}            // 关闭事件
    callbacks *gmap.StringInterfaceMap // 监听的回调函数
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


// 创建监听管理对象
func New() (*Watcher, error) {
    if watch, err := fsnotify.NewWatcher(); err == nil {
        w := &Watcher {
            watcher   : watch,
            closeChan : make(chan struct{}, 1),
            callbacks : gmap.NewStringInterfaceMap(),
        }
        w.startWatchLoop()
        return w, nil
    } else {
        return nil, err
    }
}

// 关闭监听管理对象
func (w *Watcher) Close() {
    w.watcher.Close()
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
                    event := &Event{
                        Path : ev.Name,
                        Op   : Op(ev.Op),
                    }
                    if l := w.callbacks.Get(event.Path); l != nil {
                        grpool.Add(func() {
                            for _, v := range l.(*glist.List).FrontAll() {
                                v.(func(event *Event))(event)
                            }
                        })
                    }

                //case err := <- w.watcher.Errors:
                //    log.Println("error : ", err);
                //    return
            }
        }
    }();
}