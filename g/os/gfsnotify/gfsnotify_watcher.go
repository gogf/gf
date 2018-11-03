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
)

// 关闭监听管理对象
func (w *Watcher) Close() {
    w.watcher.Close()
    w.events.Close()
    close(w.closeChan)
}

// 添加对指定文件/目录的监听，并给定回调函数
func (w *Watcher) addWatch(path string, callback func(event *Event)) error {
    // 这里统一转换为当前系统的绝对路径，便于统一监控文件名称
    t := fileRealPath(path)
    if t == "" {
        return errors.New(fmt.Sprintf(`"%s" does not exist`, path))
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

// 添加监控，path参数支持文件或者目录路径，recursive为非必需参数，默认为递归添加监控(当path为目录时)
func (w *Watcher) Add(path string, callback func(event *Event), recursive...bool) error {
    if fileIsDir(path) && (len(recursive) == 0 || recursive[0]) {
        paths, _ := fileScanDir(path, "*", true)
        list  := []string{path}
        list   = append(list, paths...)
        for _, v := range list {
            if err := w.addWatch(v, callback); err != nil {
                return err
            }
        }
        return nil
    } else {
        return w.addWatch(path, callback)
    }
}


// 移除监听
func (w *Watcher) removeWatch(path string) error {
    w.callbacks.Remove(path)
    return w.watcher.Remove(path)
}

// 递归移除监听
func (w *Watcher) Remove(path string) error {
    if fileIsDir(path) {
        paths, _ := fileScanDir(path, "*", true)
        list := []string{path}
        list  = append(list, paths...)
        for _, v := range list {
            if err := w.removeWatch(v); err != nil {
                return err
            }
        }
        return nil
    } else {
        return w.removeWatch(path)
    }
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
                    for _, callback := range callbacks.FrontAll() {
                        w.Add(event.Path, callback.(func(event *Event)))
                    }
                }
                if callbacks != nil {
                    go func(callbacks *glist.List) {
                        for _, callback := range callbacks.FrontAll() {
                            callback.(func(event *Event))(event)
                        }
                    }(callbacks)
                }
            } else {
                break
            }
        }
    }()
}