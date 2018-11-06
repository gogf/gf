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
                // 只有主callback才记录到id map中，因为子callback是自动管理的
                callbackIdMap.Set(callback.Id, callback)
            }
            if parentCallback != nil {
                // 需要递归查找到顶级的callback
                parent := parentCallback
                for {
                    if p := parent.parent; p != nil {
                        parent = p
                    } else {
                        break
                    }
                }
                parent.subs.PushFront(callback)
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
    w.callbacks.Remove(path)
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
    // 首先删除主callback
    if err := w.removeCallback(callback); err != nil {
        return err
    }
    // 如果存在子级callback，那么也一并删除
    if callback.subs.Len() > 0 {
        for {
            if r := callback.subs.PopBack(); r != nil {
                w.removeCallback(r.(*Callback))
            } else {
                break
            }
        }
        return nil
    }
    return nil
}

// 移除对指定文件/目录的所有监听
func (w *Watcher) removeCallback(callback *Callback) error {
    if r := w.callbacks.Get(callback.Path); r != nil {
        list := r.(*glist.List)
        list.Remove(callback.elem)
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
    for path != "/" {
        if l := w.callbacks.Get(path); l != nil {
            return l.(*glist.List)
        } else {
            path = fileDir(path)
        }
    }
    return nil
}

// 事件循环
func (w *Watcher) startEventLoop() {
    go func() {
        for {
            if v := w.events.Pop(); v != nil {
                event := v.(*Event)
                if event.IsRemove() {
                    if fileExists(event.Path) {
                        // 如果是文件删除事件，判断该文件是否存在，如果存在，那么将此事件认为“假删除”，
                        // 并重新添加监控(底层fsnotify会自动删除掉监控，这里重新添加回去)
                        w.watcher.Add(event.Path)
                        // 修改事件操作为重命名(相当于重命名为自身名称，最终名称没变)
                        event.Op = RENAME
                    } else {
                        // 如果是真实删除，那么递归删除监控信息
                        w.Remove(event.Path)
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