// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// 文件监控.
// 使用时需要注意的是，一旦一个文件被删除，那么对其的监控将会失效；如果删除的是目录，那么该目录及其下的文件都将被递归删除监控。
package gfsnotify

import (
    "gitee.com/johng/gf/g/container/gmap"
    "gitee.com/johng/gf/g/container/gqueue"
    "gitee.com/johng/gf/g/encoding/ghash"
    "gitee.com/johng/gf/g/os/gcache"
    "gitee.com/johng/gf/third/github.com/fsnotify/fsnotify"
)

// 监听管理对象
type Watcher struct {
    watcher    *fsnotify.Watcher        // 底层fsnotify对象
    events     *gqueue.Queue            // 过滤后的事件通知，不会出现重复事件
    closeChan  chan struct{}            // 关闭事件
    callbacks  *gmap.StringInterfaceMap // 监听的回调函数
    cache      *gcache.Cache            // 缓存对象，用于事件重复过滤
}

// 监听事件对象
type Event struct {
    event   fsnotify.Event  // 底层事件对象
    Path    string          // 文件绝对路径
    Op      Op              // 触发监听的文件操作
    Watcher *Watcher        // 事件对应的监听对象
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
    DEFAULT_WATCHER_COUNT        = 8 // 默认创建的监控对象数量(使用哈希取模)
)

var (
    // 全局监听对象，方便应用端调用
    watchers = make([]*Watcher, DEFAULT_WATCHER_COUNT)
)


// 包初始化，创建8个watcher对象，用于包默认管理监听
func init() {
    for i := 0; i < DEFAULT_WATCHER_COUNT; i++ {
        if w, err := New(); err == nil {
            watchers[i] = w
        } else {
            panic(err)
        }
    }
}

// 创建监听管理对象，主要注意的是创建监听对象会占用系统的inotify句柄数量，受到 fs.inotify.max_user_instances 的限制
func New() (*Watcher, error) {
    if watch, err := fsnotify.NewWatcher(); err == nil {
        w := &Watcher {
            cache      : gcache.New(),
            watcher    : watch,
            events     : gqueue.New(),
            closeChan  : make(chan struct{}),
            callbacks  : gmap.NewStringInterfaceMap(),
        }
        w.startWatchLoop()
        w.startEventLoop()
        return w, nil
    } else {
        return nil, err
    }
}

// 添加对指定文件/目录的监听，并给定回调函数；如果给定的是一个目录，默认递归监控。
func Add(path string, callback func(event *Event), recursive...bool) error {
    return getWatcherByPath(path).Add(path, callback, recursive...)
}

// 移除监听，默认递归删除。
func Remove(path string) error {
    return getWatcherByPath(path).Remove(path)
}

// 根据path计算对应的watcher对象
func getWatcherByPath(path string) *Watcher {
    return watchers[ghash.BKDRHash([]byte(path)) % DEFAULT_WATCHER_COUNT]
}
