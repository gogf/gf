<<<<<<< HEAD
// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// 文件监控.
// 使用时需要注意的是，一旦一个文件被删除，那么对其的监控将会失效。
=======
// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gfsnotify provides a platform-independent interface for file system notifications.
//
// 文件监控.
>>>>>>> upstream/master
package gfsnotify

import (
    "errors"
<<<<<<< HEAD
    "gitee.com/johng/gf/g/os/glog"
    "github.com/fsnotify/fsnotify"
    "gitee.com/johng/gf/g/os/gfile"
    "gitee.com/johng/gf/g/os/grpool"
    "gitee.com/johng/gf/g/container/gmap"
    "gitee.com/johng/gf/g/container/glist"
    "gitee.com/johng/gf/g/container/gqueue"
=======
    "fmt"
    "github.com/gogf/gf/g/container/glist"
    "github.com/gogf/gf/g/container/gmap"
    "github.com/gogf/gf/g/container/gqueue"
    "github.com/gogf/gf/g/container/gtype"
    "github.com/gogf/gf/g/os/gcache"
    "github.com/gogf/gf/third/github.com/fsnotify/fsnotify"
>>>>>>> upstream/master
)

// 监听管理对象
type Watcher struct {
<<<<<<< HEAD
    watcher    *fsnotify.Watcher        // 底层fsnotify对象
    events     *gqueue.Queue            // 过滤后的事件通知，不会出现重复事件
    closeChan  chan struct{}            // 关闭事件
    callbacks  *gmap.StringInterfaceMap // 监听的回调函数
=======
    watcher        *fsnotify.Watcher        // 底层fsnotify对象
    events         *gqueue.Queue            // 过滤后的事件通知，不会出现重复事件
    cache          *gcache.Cache            // 缓存对象，主要用于事件重复过滤
    callbacks      *gmap.StrAnyMap          // 注册的所有绝对路径(文件/目录)及其对应的回调函数列表map
    closeChan      chan struct{}            // 关闭事件
}

// 注册的监听回调方法
type Callback struct {
    Id        int                 // 唯一ID
    Func      func(event *Event)  // 回调方法
    Path      string              // 监听的文件/目录
    elem      *glist.Element      // 指向回调函数链表中的元素项位置(便于删除)
    recursive bool                // 当目录时，是否递归监听(使用在子文件/目录回溯查找回调函数时)
>>>>>>> upstream/master
}

// 监听事件对象
type Event struct {
<<<<<<< HEAD
    Path string // 文件绝对路径
    Op   Op     // 触发监听的文件操作
=======
    event   fsnotify.Event   // 底层事件对象
    Path    string           // 文件绝对路径
    Op      Op               // 触发监听的文件操作
    Watcher *Watcher         // 事件对应的监听对象
>>>>>>> upstream/master
}

// 按位进行识别的操作集合
type Op uint32

<<<<<<< HEAD
=======
// 必须放到一个const分组里面
>>>>>>> upstream/master
const (
    CREATE Op = 1 << iota
    WRITE
    REMOVE
    RENAME
    CHMOD
)

<<<<<<< HEAD
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
    // 这里统一转换为当前系统的绝对路径，便于统一监控文件名称
    t := gfile.RealPath(path)
    if t == "" {
        return errors.New(path + " does not exist")
    }
    path = t
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
                    w.events.PushBack(&Event{
                        Path : ev.Name,
                        Op   : Op(ev.Op),
                    })

                case err := <- w.watcher.Errors:
                    glog.Error("error : ", err);
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
                // 如果是文件删除事件，判断该文件是否存在，如果存在，那么将此事件认为“假删除”，并重新添加监控
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
=======
const (
    REPEAT_EVENT_FILTER_INTERVAL = 1      // (毫秒)重复事件过滤间隔
    gFSNOTIFY_EVENT_EXIT         = "exit" // 是否退出回调执行
)

var (
    // 默认的Watcher对象
    defaultWatcher, _   = New()
    // 默认的watchers是否初始化，使用时才创建
    watcherInited       = gtype.NewBool()
    // 回调方法ID与对象指针的映射哈希表，用于根据ID快速查找回调对象
    callbackIdMap       = gmap.NewIntAnyMap()
    // 回调函数的ID生成器(原子操作)
    callbackIdGenerator = gtype.NewInt()
)

// 创建监听管理对象，主要注意的是创建监听对象会占用系统的inotify句柄数量，受到 fs.inotify.max_user_instances 的限制
func New() (*Watcher, error) {
    w := &Watcher {
        cache     : gcache.New(),
        events    : gqueue.New(),
        closeChan : make(chan struct{}),
        callbacks : gmap.NewStrAnyMap(),
    }
    if watcher, err := fsnotify.NewWatcher(); err == nil {
        w.watcher = watcher
    } else {
        return nil, err
    }
    w.startWatchLoop()
    w.startEventLoop()
    return w, nil
}

// 添加对指定文件/目录的监听，并给定回调函数；如果给定的是一个目录，默认非递归监控。
func Add(path string, callbackFunc func(event *Event), recursive...bool) (callback *Callback, err error) {
    return defaultWatcher.Add(path, callbackFunc, recursive...)
}

// 递归移除对指定文件/目录的所有监听回调
func Remove(path string) error {
    return defaultWatcher.Remove(path)
}

// 根据指定的回调函数ID，移出指定的inotify回调函数
func RemoveCallback(callbackId int) error {
    callback := (*Callback)(nil)
    if r := callbackIdMap.Get(callbackId); r != nil {
        callback = r.(*Callback)
    }
    if callback == nil {
        return errors.New(fmt.Sprintf(`callback for id %d not found`, callbackId))
    }
    defaultWatcher.RemoveCallback(callbackId)
    return nil
}

// 在回调方法中调用该方法退出回调注册
func Exit() {
    panic(gFSNOTIFY_EVENT_EXIT)
}
>>>>>>> upstream/master
