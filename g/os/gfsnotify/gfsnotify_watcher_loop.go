// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package gfsnotify

import (
    "fmt"
    "gitee.com/johng/gf/g/container/glist"
    "sync"
    "time"
)

// 监听循环
func (w *Watcher) startWatchLoop() {
    for i := 0; i < len(w.watchers); i++ {
        go func(i int) {
            for {
                select {
                    // 关闭事件
                    case <- w.closeChan: return

                        // 监听事件
                    case ev := <- w.watchers[i].Events:
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

                    case err := <- w.watchers[i].Errors:
                        fmt.Errorf("error: %s\n" + err.Error());
                }
            }
        }(i)
    }
}

// 获得真正监听的文件路径及回调函数列表，假如是临时文件或者新增文件，是无法搜索都监听回调的。
// 判断规则：
// 1、在 callbacks 中应当有回调注册函数(否则监听根本没意义)；
// 2、如果该path下不存在回调注册函数，则按照path长度从右往左递减，直到减到目录地址为止(不包含)；
// 3、如果仍旧无法匹配回调函数，那么忽略，否则使用查找到的新path覆盖掉event的path；
// 解决问题：
// 1、部分IDE修改文件时生成的临时文件，如: /index.html -> /index.html__jdold__；
func (w *Watcher) getWatchPathAndCallbacks(path string) (watchPath string, callbacks *glist.List) {
    if path == "" {
        return "", nil
    }
    dirPath := fileDir(path)
    for {
        if v := w.callbacks.Get(path); v != nil {
            return path, v.(*glist.List)
        }
        path = path[0 : len(path) - 1]
        // 递减到上一级目录为止
        if path == dirPath {
            break
        }
        // 如果不能再继续递减，那么退出
        if len(path) == 0 {
            break
        }
    }
    return "", nil
}

// 事件循环
func (w *Watcher) startEventLoop() {
    go func() {
        for {
            if v := w.events.Pop(); v != nil {
                event := v.(*Event)
                // watchPath是注册回调的路径，可能和event.Path不一样
                watchPath, callbacks := w.getWatchPathAndCallbacks(event.Path)
                if callbacks == nil {
                    continue
                }
                fmt.Println("event:", event.String(), watchPath, fileExists(watchPath))
                switch {
                    // 如果是删除操作，那么需要判断是否文件真正不存在了，如果存在，那么将此事件认为“假删除”
                    case event.IsRemove():
                        if fileExists(watchPath) {
                            // 重新添加监控(底层fsnotify会自动删除掉监控，这里重新添加回去)
                            // 注意这里调用的是底层fsnotify添加监控，只会产生回调事件，并不会使回调函数重复注册
                            w.watcher(watchPath).Add(event.Path)
                            // 修改事件操作为重命名(相当于重命名为自身名称，最终名称没变)
                            event.Op = RENAME
                        } else {
                            // 删除之前需要执行一遍回调，否则Remove之后就无法执行了
                            // 由于是异步回调，这里保证所有回调都开始执行后再执行删除
                            wg := sync.WaitGroup{}
                            for _, v := range callbacks.FrontAll() {
                                wg.Add(1)
                                go func(callback *Callback) {
                                    wg.Done()
                                    callback.Func(event)
                                }(v.(*Callback))
                            }
                            wg.Wait()
                            time.Sleep(time.Second)
                            // 如果是真实删除，那么递归删除监控信息
                            fmt.Println("remove", watchPath)
                            w.Remove(watchPath)
                        }

                    // 如果是重命名操作，那么需要判断是否文件真正不存在了，如果存在，那么将此事件认为“假命名”
                    // (特别是某些编辑器在编辑文件时会先对文件RENAME再CHMOD)
                    case event.IsRename():
                        if fileExists(watchPath) {
                            // 重新添加监控
                            w.watcher(watchPath).Add(watchPath)
                        } else if watchPath != event.Path && fileExists(event.Path) {
                            for _, v := range callbacks.FrontAll() {
                                callback := v.(*Callback)
                                w.addWithCallbackFunc(callback, event.Path, callback.Func)
                            }
                        }

                    // 创建文件/目录
                    case event.IsCreate():
                        for _, v := range callbacks.FrontAll() {
                            callback := v.(*Callback)
                            w.addWithCallbackFunc(callback, event.Path, callback.Func)
                        }

                }
                // 执行回调处理，异步处理
                for _, v := range callbacks.FrontAll() {
                    go v.(*Callback).Func(event)
                }

            } else {
                break
            }
        }
    }()
}