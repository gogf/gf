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
    "gitee.com/johng/gf/g/encoding/ghash"
    "gitee.com/johng/gf/third/github.com/fsnotify/fsnotify"
)

// 添加监控，path参数支持文件或者目录路径，recursive为非必需参数，默认为递归添加监控(当path为目录时)。
// 如果添加目录，这里只会返回目录的callback，按照callback删除时会递归删除。
func (w *Watcher) Add(path string, callbackFunc func(event *Event), recursive...bool) (callback *Callback, err error) {
    return w.addWithCallbackFunc(nil, path, callbackFunc, recursive...)
}

// 添加监控，path参数支持文件或者目录路径，recursive为非必需参数，默认为递归添加监控(当path为目录时)。
// 如果添加目录，这里只会返回目录的callback，按照callback删除时会递归删除。
func (w *Watcher) addWithCallbackFunc(parentCallback *Callback, path string, callbackFunc func(event *Event), recursive...bool) (callback *Callback, err error) {
    // 首先添加这个文件/目录
    callback, err = w.doAddWithCallbackFunc(path, callbackFunc, parentCallback)
    if err != nil {
        return nil, err
    }
    // 其次递归添加其下的文件/目录
    if fileIsDir(path) && (len(recursive) == 0 || recursive[0]) {
        // 追加递归监控的回调到recursivePaths中
        w.recursivePaths.LockFunc(func(m map[string]interface{}) {
            list := (*glist.List)(nil)
            if v, ok := m[path]; !ok {
                list    = glist.New()
                m[path] = list
            } else {
                list    = v.(*glist.List)
            }
            list.PushBack(callback)
        })
        // 递归添加监控
        paths, _ := fileScanDir(path, "*", true)
        for _, v := range paths {
            w.doAddWithCallbackFunc(v, callbackFunc, callback)
        }
    }
    return
}

// 添加对指定文件/目录的监听，并给定回调函数
func (w *Watcher) doAddWithCallbackFunc(path string, callbackFunc func(event *Event), parentCallback *Callback) (callback *Callback, err error) {
    // 这里统一转换为当前系统的绝对路径，便于统一监控文件名称
    if t := fileRealPath(path); t == "" {
        return nil, errors.New(fmt.Sprintf(`"%s" does not exist`, path))
    } else {
        path = t
    }
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
        Id     : callbackIdGenerator.Add(1),
        Func   : callbackFunc,
        Path   : path,
        addr   : fmt.Sprintf("%p", callbackFunc)[2:],
        subs   : glist.New(),
        parent : parentCallback,
    }
    // 注册回调函数
    w.callbacks.LockFunc(func(m map[string]interface{}) {
        list := (*glist.List)(nil)
        if v, ok := m[path]; !ok {
            list    = glist.New()
            m[path] = list
        } else {
            list    = v.(*glist.List)
        }
        callback.elem = list.PushBack(callback)
    })
    // 添加底层监听
    w.watcher(path).Add(path)
    return
}

// 根据path查询对应的底层watcher对象
func (w *Watcher) watcher(path string) *fsnotify.Watcher {
    return w.watchers[ghash.BKDRHash([]byte(path)) % uint32(len(w.watchers))]
}

// 关闭监听管理对象
func (w *Watcher) Close() {
    for _, watcher := range w.watchers  {
        watcher.Close()
    }
    w.events.Close()
    close(w.closeChan)
}

// 递归移除对指定文件/目录的所有监听回调
func (w *Watcher) Remove(path string) error {
    if fileIsDir(path) && fileExists(path) {
        paths, _ := fileScanDir(path, "*", true)
        paths     = append(paths, path)
        for _, v := range paths {
            if err := w.removeWatch(v); err != nil {
                return err
            }
        }
        return nil
    } else {
        return w.removeWatch(path)
    }
}

// 移除对指定文件/目录的所有监听
func (w *Watcher) removeWatch(path string) error {
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
    return w.watcher(path).Remove(path)
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

// (递归)移除对指定文件/目录的所有监听
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
        // 如果该文件/目录的所有回调都被删除，那么移除底层的监听
        if list.Len() == 0 {
            return w.watcher(callback.Path).Remove(callback.Path)
        }
    } else {
        return errors.New(fmt.Sprintf(`callbacks not found for "%s"`, callback.Path))
    }
    return nil
}
