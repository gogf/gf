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
    "gitee.com/johng/gf/g/encoding/ghash"
    "gitee.com/johng/gf/g/os/gcache"
    "gitee.com/johng/gf/g/os/gcmd"
    "gitee.com/johng/gf/g/os/genv"
    "gitee.com/johng/gf/g/util/gconv"
    "gitee.com/johng/gf/third/github.com/fsnotify/fsnotify"
)

// 监听管理对象
type Watcher struct {
    watcher       *fsnotify.Watcher        // 底层fsnotify对象
    events        *gqueue.Queue            // 过滤后的事件通知，不会出现重复事件
    closeChan     chan struct{}            // 关闭事件
    callbacks     *gmap.StringInterfaceMap // 监听的回调函数
    cache         *gcache.Cache            // 缓存对象，用于事件重复过滤
}

// 注册的监听回调方法
type Callback struct {
    Id     int                 // 唯一ID
    Func   func(event *Event)  // 回调方法
    Path   string              // 监听的文件/目录
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
    DEFAULT_WATCHER_COUNT        = 4 // 默认创建的监控对象数量(使用哈希取模)
)

var (
    // 全局监听对象，方便应用端调用
    watchers     []*Watcher
    // 全局默认的监听watcher数量
    watcherCount int
    // 默认的watchers是否初始化，使用时才创建
    watcherInited  = gtype.NewBool()
    // 回调方法ID与对象指针的映射哈希表，用于根据ID快速查找回调对象
    callbackIdMap  = gmap.NewIntInterfaceMap()
)

// 初始化创建watcher对象，用于包默认管理监听
func initWatcher() {
    if !watcherInited.Set(true) {
        // 默认的创建的inotify数量
        watcherCount = gconv.Int(genv.Get("GF_INOTIFY_COUNT"))
        if watcherCount == 0 {
            watcherCount = gconv.Int(gcmd.Option.Get("gf.inotify-count"))
        }
        if watcherCount == 0 {
            watcherCount = DEFAULT_WATCHER_COUNT
        }
        watchers = make([]*Watcher, watcherCount)
        for i := 0; i < watcherCount; i++ {
            if w, err := New(); err == nil {
                watchers[i] = w
            } else {
                panic(err)
            }
        }
    }
}

// 创建监听管理对象，主要注意的是创建监听对象会占用系统的inotify句柄数量，受到 fs.inotify.max_user_instances 的限制
func New() (*Watcher, error) {
    if watch, err := fsnotify.NewWatcher(); err == nil {
        w := &Watcher {
            cache         : gcache.New(),
            watcher       : watch,
            events        : gqueue.New(),
            closeChan     : make(chan struct{}),
            callbacks     : gmap.NewStringInterfaceMap(),
        }
        w.startWatchLoop()
        w.startEventLoop()
        return w, nil
    } else {
        return nil, err
    }
}

// 添加对指定文件/目录的监听，并给定回调函数；如果给定的是一个目录，默认递归监控。
func Add(path string, callbackFunc func(event *Event), recursive...bool) (callback *Callback, err error) {
    return getWatcherByPath(path).Add(path, callbackFunc, recursive...)
}

// 递归移除对指定文件/目录的所有监听回调
func Remove(path string) error {
    return getWatcherByPath(path).Remove(path)
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
    return getWatcherByPath(callback.Path).RemoveCallback(callbackId)
}

// 根据path计算对应的watcher对象
func getWatcherByPath(path string) *Watcher {
    initWatcher()
    return watchers[ghash.BKDRHash([]byte(path)) % uint32(watcherCount)]
}
