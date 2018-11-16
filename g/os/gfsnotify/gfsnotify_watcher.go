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

// 添加监控，path参数支持文件或者目录路径，recursive为非必需参数，默认为非递归监控(当path为目录时)。
// 如果添加目录，这里只会返回目录的callback，按照callback删除时会递归删除。
func (w *Watcher) Add(path string, callbackFunc func(event *Event), recursive...bool) (callback *Callback, err error) {
    // 首先添加这个文件/目录
    callback, err = w.addWithCallbackFunc(path, callbackFunc, recursive...)
    if err != nil {
        return nil, err
    }
    // 如果需要递归，那么递归添加其下的子级目录，
    // 注意!!
    // 1、这里只递归添加**目录**, 而非文件，因为监控了目录即监控了其下一级的文件;
    // 2、这里只是添加底层监控对象对**子级所有目录**的监控，没有任何回调函数的设置，在事件产生时会回溯查找父级的回调函数；
    if fileIsDir(path) && (len(recursive) == 0 || recursive[0]) {
        for _, subPath := range fileAllDirs(path) {
            if fileIsDir(subPath) {
                w.watcher.Add(subPath)
            }
        }
    }
    return
}

// 添加对指定文件/目录的监听，并给定回调函数
func (w *Watcher) addWithCallbackFunc(path string, callbackFunc func(event *Event), recursive...bool) (callback *Callback, err error) {
    // 这里统一转换为当前系统的绝对路径，便于统一监控文件名称
    if t := fileRealPath(path); t == "" {
        return nil, errors.New(fmt.Sprintf(`"%s" does not exist`, path))
    } else {
        path = t
    }
    callback = &Callback {
        Id     : callbackIdGenerator.Add(1),
        Func   : callbackFunc,
        Path   : path,
    }
    if len(recursive) > 0 {
        callback.recursive = recursive[0]
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
    w.watcher.Add(path)
    // 添加成功后会注册该callback id到全局的哈希表
    callbackIdMap.Set(callback.Id, callback)
    return
}

// 关闭监听管理对象
func (w *Watcher) Close() {
    w.events.Close()
    w.watcher.Close()
    close(w.closeChan)
}

// 递归移除对指定文件/目录的所有监听回调
func (w *Watcher) Remove(path string) error {
    // 首先移除path注册的回调注册，以及callbackIdMap中的ID
    if r := w.callbacks.Remove(path); r != nil {
        list := r.(*glist.List)
        for {
            if r := list.PopFront(); r != nil {
                callbackIdMap.Remove(r.(*Callback).Id)
            } else {
                break
            }
        }
    }
    // 其次递归判断所有的子级是否可删除监听
    if subPaths, err := fileScanDir(path, "*", true); err == nil && len(subPaths) > 0 {
        for _, subPath := range subPaths {
            if w.checkPathCanBeRemoved(subPath) {
                w.watcher.Remove(subPath)
            }
        }
    }
    // 最后移除底层的监听
    return w.watcher.Remove(path)
}

// 判断给定的路径是否可以删除监听(只有所有回调函数都没有了才能删除)
func (w *Watcher) checkPathCanBeRemoved(path string) bool {
    // 首先检索path对应的回调函数
    if v := w.callbacks.Get(path); v != nil {
        return false
    }
    // 其次查找父级目录有无回调注册
    dirPath := fileDir(path)
    if v := w.callbacks.Get(dirPath); v != nil {
        return false
    }
    // 最后回溯查找递归回调函数
    for {
        parentDirPath := fileDir(dirPath)
        if parentDirPath == dirPath {
            break
        }
        if v := w.callbacks.Get(parentDirPath); v != nil {
            return false
        }
        dirPath = parentDirPath
    }
    return true
}

// 根据指定的回调函数ID，移出指定的inotify回调函数
func (w *Watcher) RemoveCallback(callbackId int) {
    callback := (*Callback)(nil)
    if r := callbackIdMap.Get(callbackId); r != nil {
        callback = r.(*Callback)
    }
    if callback != nil {
        if r := w.callbacks.Get(callback.Path); r != nil {
            r.(*glist.List).Remove(callback.elem)
        }
        callbackIdMap.Remove(callbackId)
    }
}

