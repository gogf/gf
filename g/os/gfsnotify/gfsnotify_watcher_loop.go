// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package gfsnotify

import (
    "gitee.com/johng/gf/g/container/glist"
)

// 监听循环
func (w *Watcher) startWatchLoop() {
    go func() {
        for {
            select {
                // 关闭事件
                case <- w.closeChan: return

                    // 监听事件
                case ev := <- w.watcher.Events:
                    //fmt.Println("ev:", ev.String())
                    w.cache.SetIfNotExist(ev.String(), func() interface{} {
                        w.events.Push(&Event{
                            event   : ev,
                            Path    : ev.Name,
                            Op      : Op(ev.Op),
                            Watcher : w,
                        })
                        return struct {}{}
                    }, REPEAT_EVENT_FILTER_INTERVAL)

                case <- w.watcher.Errors:
                    //fmt.Fprintf(os.Stderr, "[gfsnotify] error: %s\n", err.Error())
            }
        }
    }()
}

// 获得文件路径的监听回调，包括层级的监听回调。
func (w *Watcher) getCallbacks(path string) (callbacks []*Callback) {
    // 首先检索path对应的回调函数
    if v := w.callbacks.Get(path); v != nil {
        for _, v := range v.(*glist.List).FrontAll() {
            callback := v.(*Callback)
            callbacks = append(callbacks, callback)
        }
    }
    // 其次查找父级目录有无回调注册
    dirPath := fileDir(path)
    if v := w.callbacks.Get(dirPath); v != nil {
        for _, v := range v.(*glist.List).FrontAll() {
            callback := v.(*Callback)
            callbacks = append(callbacks, callback)
        }
    }
    // 最后回溯查找递归回调函数
    for {
        parentDirPath := fileDir(dirPath)
        if parentDirPath == dirPath {
            break
        }
        if v := w.callbacks.Get(parentDirPath); v != nil {
            for _, v := range v.(*glist.List).FrontAll() {
                callback := v.(*Callback)
                if callback.recursive {
                    callbacks = append(callbacks, callback)
                }
            }
        }
        dirPath = parentDirPath
    }
    return
}

// 事件循环
func (w *Watcher) startEventLoop() {
    go func() {
        for {
            if v := w.events.Pop(); v != nil {
                event := v.(*Event)
                // 如果该路径一个回调也没有，那么没有必要执行后续逻辑，删除对该文件的监听
                callbacks := w.getCallbacks(event.Path)
                if len(callbacks) == 0 {
                    w.watcher.Remove(event.Path)
                    continue
                }
                switch {
                    // 如果是删除操作，那么需要判断是否文件真正不存在了，如果存在，那么将此事件认为“假删除”
                    case event.IsRemove():
                        if fileExists(event.Path) {
                            // 底层重新添加监控(不用担心重复添加)
                            w.watcher.Add(event.Path)
                            // 修改事件操作为重命名(相当于重命名为自身名称，最终名称没变)
                            event.Op = RENAME
                        }

                    // 如果是重命名操作，那么需要判断是否文件真正不存在了，如果存在，那么将此事件认为“假命名”
                    // (特别是某些编辑器在编辑文件时会先对文件RENAME再CHMOD)
                    case event.IsRename():
                        if fileExists(event.Path) {
                            // 底层有可能去掉了监控, 这里重新添加监控(不用担心重复添加)
                            w.watcher.Add(event.Path)
                            // 修改事件操作为修改属性
                            event.Op = CHMOD
                        }

                    // 创建文件/目录
                    case event.IsCreate():
                        // =========================================
                        // 注意这里只是添加底层监听，并没有注册任何的回调函数，
                        // 默认的回调函数为父级的递归回调
                        // =========================================
                        if fileIsDir(event.Path) {
                            // 递归添加
                            for _, subPath := range fileAllDirs(event.Path) {
                                if fileIsDir(subPath) {
                                    w.watcher.Add(subPath)
                                }
                            }
                        } else {
                            // 添加文件监听
                            w.watcher.Add(event.Path)
                        }

                }
                // 执行回调处理，异步处理
                for _, callback := range callbacks {
                    go callback.Func(event)
                }

            } else {
                break
            }
        }
    }()
}