// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// 文件监控.
package gfsnotify

import (
    "container/list"
    "errors"
    "fmt"
    "gitee.com/johng/gf/g/container/glist"
    "gitee.com/johng/gf/g/container/gmap"
    "gitee.com/johng/gf/g/container/gqueue"
    "gitee.com/johng/gf/g/container/gtype"
    "gitee.com/johng/gf/g/os/gcache"
    "gitee.com/johng/gf/g/os/gcmd"
    "gitee.com/johng/gf/g/os/genv"
    "gitee.com/johng/gf/g/util/gconv"
    "gitee.com/johng/gf/third/github.com/fsnotify/fsnotify"
)

// 监听管理对象
type Watcher struct {
    watchers       []*fsnotify.Watcher      // 底层fsnotify对象，支持多个，以避免单个inotify对象监听队列上限问题
    events         *gqueue.Queue            // 过滤后的事件通知，不会出现重复事件
    cache          *gcache.Cache            // 缓存对象，主要用于事件重复过滤
    callbacks      *gmap.StringInterfaceMap // 注册的所有绝对路径机器对应的回调函数列表map
    recursivePaths *gmap.StringInterfaceMap // 支持递归监听的目录绝对路径及其对应的回调函数列表map
    closeChan      chan struct{}            // 关闭事件
}

// 注册的监听回调方法
type Callback struct {
    Id     int                 // 唯一ID
    Func   func(event *Event)  // 回调方法
    Path   string              // 监听的文件/目录
    addr   string              // Func对应的内存地址，用以判断回调的重复
    elem   *list.Element       // 指向监听链表中的元素项位置
    parent *Callback           // 父级callback，有这个属性表示该callback为被自动管理的callback
    subs   *glist.List         // 子级回调对象指针列表
}

// 监听事件对象
type Event struct {
    event   fsnotify.Event   // 底层事件对象
    Path    string           // 文件绝对路径
    Op      Op               // 触发监听的文件操作
    Watcher *Watcher         // 事件对应的监听对象
}

// 按位进行识别的操作集合
type Op uint32

// 必须放到一个const分组里面
const (
    CREATE Op = 1 << iota
    WRITE
    REMOVE
    RENAME
    CHMOD
)

const (
    REPEAT_EVENT_FILTER_INTERVAL = 1 // (毫秒)重复事件过滤间隔
    DEFAULT_WATCHER_COUNT        = 1 // 默认创建的监控对象数量(使用哈希取模)
    gDEFAULT_PKG_WATCHER_COUNT   = 4 // 默认创建的包监控对象数量(使用哈希取模)
)

var (
    // 默认的Watcher对象
    defaultWatcher *Watcher
    // 默认的watchers是否初始化，使用时才创建
    watcherInited  = gtype.NewBool()
    // 回调方法ID与对象指针的映射哈希表，用于根据ID快速查找回调对象
    callbackIdMap  = gmap.NewIntInterfaceMap()
    // 回调函数的ID生成器(原子操作)
    callbackIdGenerator = gtype.NewInt()
)

// 初始化创建watcher对象，用于包默认管理监听
func initWatcher() {
    if !watcherInited.Set(true) {
        pkgWatcherCount := gconv.Int(genv.Get("GF_INOTIFY_COUNT"))
        if pkgWatcherCount == 0 {
            pkgWatcherCount = gconv.Int(gcmd.Option.Get("gf.inotify-count"))
        }
        if pkgWatcherCount == 0 {
            pkgWatcherCount = gDEFAULT_PKG_WATCHER_COUNT
        }
        if w, err := New(pkgWatcherCount); err == nil {
            defaultWatcher = w
        } else {
            panic(err)
        }
    }
}

// 创建监听管理对象，主要注意的是创建监听对象会占用系统的inotify句柄数量，受到 fs.inotify.max_user_instances 的限制
func New(inotifyCount...int) (*Watcher, error) {
    count := DEFAULT_WATCHER_COUNT
    if len(inotifyCount) > 0 {
        count = inotifyCount[0]
    }
    w := &Watcher {
        cache          : gcache.New(),
        watchers       : make([]*fsnotify.Watcher, count),
        events         : gqueue.New(),
        closeChan      : make(chan struct{}),
        callbacks      : gmap.NewStringInterfaceMap(),
        recursivePaths : gmap.NewStringInterfaceMap(),
    }
    for i := 0; i < count; i++ {
        if watcher, err := fsnotify.NewWatcher(); err == nil {
            w.watchers[i] = watcher
        } else {
            // 出错，关闭已创建的底层watcher对象
            for j := 0; j < i; j++ {
                w.watchers[j].Close()
            }
            return nil, err
        }
    }
    w.startWatchLoop()
    w.startEventLoop()
    return w, nil
}

// 添加对指定文件/目录的监听，并给定回调函数；如果给定的是一个目录，默认递归监控。
func Add(path string, callbackFunc func(event *Event), recursive...bool) (callback *Callback, err error) {
    return watcher().Add(path, callbackFunc, recursive...)
}

// 递归移除对指定文件/目录的所有监听回调
func Remove(path string) error {
    return watcher().Remove(path)
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
    return watcher().RemoveCallback(callbackId)
}

// 获得默认的包watcher
func watcher() *Watcher {
    initWatcher()
    return defaultWatcher
}
