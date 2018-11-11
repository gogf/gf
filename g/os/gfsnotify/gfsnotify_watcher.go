// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package gfsnotify

import (
    "errors"
    "fmt"
    "gitee.com/johng/gf/g/container/glist"
    "gitee.com/johng/gf/g/os/gtime"
)

// 关闭监听管理对象
func (w *Watcher) Close() {
    w.watcher.Close()
    w.events.Close()
    close(w.closeChan)
}

// 添加对指定文件/目录的监听，并给定回调函数
func (w *Watcher) addWatch(path string, calbackFunc func(event *Event), parentCallback *Callback) (callback *Callback, err error) {
    // 这里统一转换为当前系统的绝对路径，便于统一监控文件名称
    t := fileRealPath(path)
    if t == "" {
        return nil, errors.New(fmt.Sprintf(`"%s" does not exist`, path))
    }
    path = t
    // 添加成功后会注册该callback id到全局的哈希表，并绑定到父级的注册回调中
    defer func() {
        if err == nil {
            if parentCallback == nil {
                // 只有主callback才记录到id map中，因为子callback是自动管理的无需添加到全局id映射map中
                callbackIdMap.Set(callback.Id, callback)
            }
            if parentCallback != nil {
                // 添加到直属父级的subs属性中，建立关联关系，便于后续删除
                parentCallback.subs.PushBack(callback)
            }
        }
    }()
    callback = &Callback {
        Id     : int(gtime.Nanosecond()),
        Func   : calbackFunc,
        Path   : path,
        subs   : glist.New(),
        parent : parentCallback,
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
        callback.elem = result.(*glist.List).PushBack(callback)
    })
    // 添加底层监听
    w.watcher.Add(path)
    return
}

// 添加监控，path参数支持文件或者目录路径，recursive为非必需参数，默认为递归添加监控(当path为目录时)。
// 如果添加目录，这里只会返回目录的callback，按照callback删除时会递归删除。
func (w *Watcher) addWithCallback(parentCallback *Callback, path string, callbackFunc func(event *Event), recursive...bool) (callback *Callback, err error) {
    // 首先添加这个目录
    if callback, err = w.addWatch(path, callbackFunc, parentCallback); err != nil {
        return nil, err
    }
    // 其次递归添加其下的文件/目录
    if fileIsDir(path) && (len(recursive) == 0 || recursive[0]) {
        paths, _ := fileScanDir(path, "*", true)
        for _, v := range paths {
            w.addWatch(v, callbackFunc, callback)
        }
    }
    return
}

// 添加监控，path参数支持文件或者目录路径，recursive为非必需参数，默认为递归添加监控(当path为目录时)。
// 如果添加目录，这里只会返回目录的callback，按照callback删除时会递归删除。
func (w *Watcher) Add(path string, callbackFunc func(event *Event), recursive...bool) (callback *Callback, err error) {
    return w.addWithCallback(nil, path, callbackFunc, recursive...)
}

// 递归移除对指定文件/目录的所有监听回调
func (w *Watcher) Remove(path string) error {
    if fileIsDir(path) {
        paths, _ := fileScanDir(path, "*", true)
        paths     = append(paths, path)
        for _, v := range paths {
            if err := w.removeAll(v); err != nil {
                return err
            }
        }
        return nil
    } else {
        return w.removeAll(path)
    }
}

// 移除对指定文件/目录的所有监听
func (w *Watcher) removeAll(path string) error {
    // 首先移除所有该path的回调注册
    if r := w.callbacks.Get(path); r != nil {
        list := r.(*glist.List)
        for {
            if r := list.PopFront(); r != nil {
                w.removeCallback(r.(*Callback))
            } else {
                break
            }
        }
    }
    // 其次移除该path的监听注册
    w.callbacks.Remove(path)
    // 最后移除底层的监听
    return w.watcher.Remove(path)
}

// 根据指定的回调函数ID，移出指定的inotify回调函数
func (w *Watcher) RemoveCallback(callbackId int) error {
    callback := (*Callback)(nil)
    if r := callbackIdMap.Get(callbackId); r != nil {
        callback = r.(*Callback)
    }
    if callback == nil {
        return errors.New(fmt.Sprintf(`callback for id %d not found`, callbackId))
    }
    w.removeCallback(callback)
    return nil
}

// 移除对指定文件/目录的所有监听
func (w *Watcher) removeCallback(callback *Callback) error {
    if r := w.callbacks.Get(callback.Path); r != nil {
        list := r.(*glist.List)
        list.Remove(callback.elem)
        // 如果存在子级callback，那么也一并递归删除
        if callback.subs.Len() > 0 {
            for {
                if r := callback.subs.PopFront(); r != nil {
                    w.removeCallback(r.(*Callback))
                } else {
                    break
                }
            }
        }
        // 如果该文件/目录的所有回调都被删除，那么移除监听
        if list.Len() == 0 {
            return w.watcher.Remove(callback.Path)
        }
    } else {
        return errors.New(fmt.Sprintf(`callbacks not found for "%s"`, callback.Path))
    }
    return nil
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
                    key := ev.String()
                    if !w.cache.Contains(key) {
                        w.cache.Set(key, struct {}{}, REPEAT_EVENT_FILTER_INTERVAL)
                        w.events.Push(&Event{
                            event   : ev,
                            Path    : ev.Name,
                            Op      : Op(ev.Op),
                            Watcher : w,
                        })
                    }

                case err := <- w.watcher.Errors:
                    panic("error : " + err.Error());
            }
        }
    }()
}

// 检索给定path的回调方法**列表**
func (w *Watcher) getCallbacks(path string) *glist.List {
    for {
        if l := w.callbacks.Get(path); l != nil {
            return l.(*glist.List)
        } else {
            if p := fileDir(path); p != path {
                path = p
            } else {
                break
            }
        }
    }
    return nil
}

// 获得真正监听的文件路径，判断规则：
// 1、在 callbacks 中应当有回调注册函数(否则监听根本没意义)；
// 2、如果该path下不存在回调注册函数，则按照path长度从右往左递减，直到减到目录地址为止(不包含)；
// 3、如果仍旧无法匹配回调函数，那么忽略，否则使用查找到的新path覆盖掉event的path；
func (w *Watcher) getWatchTruePath(path string) string {
    if w.getCallbacks(path) != nil {
        return path
    }
    dirPath := fileDir(path)
    for {
        path = path[0 : len(path) - 1]
        if path == dirPath {
            break
        }
        if w.getCallbacks(path) != nil {
            return path
        }
    }
    return ""
}

// 事件循环
func (w *Watcher) startEventLoop() {
    go func() {
        for {
            if v := w.events.Pop(); v != nil {
                event := v.(*Event)
                if path := w.getWatchTruePath(event.Path); path == "" {
                    continue
                } else {
                    event.Path = path
                }
                switch {
                    // 如果是删除操作，那么需要判断是否文件真正不存在了，如果存在，那么将此事件认为“假删除”
                    case event.IsRemove():
                        if fileExists(event.Path) {
                            // 重新添加监控(底层fsnotify会自动删除掉监控，这里重新添加回去)
                            // 注意这里调用的是底层fsnotify添加监控，只会产生回调事件，并不会使回调函数重复注册
                            w.watcher.Add(event.Path)
                            // 修改事件操作为重命名(相当于重命名为自身名称，最终名称没变)
                            event.Op = RENAME
                        } else {
                            // 如果是真实删除，那么递归删除监控信息
                            w.Remove(event.Path)
                        }

                    // 如果是删除操作，那么需要判断是否文件真正不存在了，如果存在，那么将此事件认为“假命名”
                    // (特别是某些编辑器在编辑文件时会先对文件RENAME再CHMOD)
                    case event.IsRename():
                        if fileExists(event.Path) {
                            // 重新添加监控
                            w.watcher.Add(event.Path)
                        }
                }

                callbacks := w.getCallbacks(event.Path)
                // 如果创建了新的目录，那么将这个目录递归添加到监控中
                if event.IsCreate() && fileIsDir(event.Path) {
                    for _, v := range callbacks.FrontAll() {
                        callback := v.(*Callback)
                        w.addWithCallback(callback, event.Path, callback.Func)
                    }
                }
                // 执行回调处理，异步处理
                if callbacks != nil {
                    go func(callbacks *glist.List) {
                        for _, v := range callbacks.FrontAll() {
                            go v.(*Callback).Func(event)
                        }
                    }(callbacks)
                }
            } else {
                break
            }
        }
    }()
}