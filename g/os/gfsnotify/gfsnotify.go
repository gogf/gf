// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// Package gfsnotify provides a platform-independent interface for file system notifications.
//
// 文件监控.
package gfsnotify

import (
    "errors"
    "fmt"
    "gitee.com/johng/gf/g/container/glist"
    "gitee.com/johng/gf/g/container/gmap"
    "gitee.com/johng/gf/g/container/gqueue"
    "gitee.com/johng/gf/g/container/gtype"
    "gitee.com/johng/gf/g/os/gcache"
    "gitee.com/johng/gf/third/github.com/fsnotify/fsnotify"
)

// 监听管理对象
type Watcher struct {
    watcher        *fsnotify.Watcher        // 底层fsnotify对象
    events         *gqueue.Queue            // 过滤后的事件通知，不会出现重复事件
    cache          *gcache.Cache            // 缓存对象，主要用于事件重复过滤
    callbacks      *gmap.StringInterfaceMap // 注册的所有绝对路径(文件/目录)及其对应的回调函数列表map
    closeChan      chan struct{}            // 关闭事件
}

// 注册的监听回调方法
type Callback struct {
    Id        int                 // 唯一ID
    Func      func(event *Event)  // 回调方法
    Path      string              // 监听的文件/目录
    elem      *glist.Element      // 指向回调函数链表中的元素项位置(便于删除)
    recursive bool                // 当目录时，是否递归监听(使用在子文件/目录回溯查找回调函数时)
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
)

var (
    // 默认的Watcher对象
    defaultWatcher, _   = New()
    // 默认的watchers是否初始化，使用时才创建
    watcherInited       = gtype.NewBool()
    // 回调方法ID与对象指针的映射哈希表，用于根据ID快速查找回调对象
    callbackIdMap       = gmap.NewIntInterfaceMap()
    // 回调函数的ID生成器(原子操作)
    callbackIdGenerator = gtype.NewInt()
)

// 创建监听管理对象，主要注意的是创建监听对象会占用系统的inotify句柄数量，受到 fs.inotify.max_user_instances 的限制
func New() (*Watcher, error) {
    w := &Watcher {
        cache     : gcache.New(),
        events    : gqueue.New(),
        closeChan : make(chan struct{}),
        callbacks : gmap.NewStringInterfaceMap(),
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
