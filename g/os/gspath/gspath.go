// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// 搜索目录管理.
// 可以添加搜索目录，按照添加的优先级进行文件检索，并在内部进行高效缓存处理。
package gspath

import (
    "sync"
    "errors"
    "strings"
    "gitee.com/johng/gf/g/os/gfile"
    "gitee.com/johng/gf/g/container/gmap"
)

// 文件目录搜索管理对象
type SPath struct {
    mu    sync.RWMutex
    paths []string              // 搜索路径，按照优先级进行排序
    cache *gmap.StringStringMap // 搜索结果缓存map
}

func New () *SPath {
    return &SPath{
        paths : make([]string, 0),
        cache : gmap.NewStringStringMap(),
    }
}

// 设置搜索路径，只保留当前设置项，其他搜索路径被清空
func (sp *SPath) Set(path string) error {
    r := gfile.RealPath(path)
    if r != "" && gfile.IsDir(r) {
        sp.mu.Lock()
        sp.paths = []string{strings.TrimRight(r, gfile.Separator)}
        sp.mu.Unlock()
        sp.cache.Clear()
        return nil
    }
    return errors.New("invalid path:" + path)
}

// 添加搜索路径
func (sp *SPath) Add(path string) error {
    r := gfile.RealPath(path)
    if r != "" && gfile.IsDir(r) {
        sp.mu.Lock()
        sp.paths = append(sp.paths, r)
        sp.mu.Unlock()
        return nil
    }
    return errors.New("invalid path:" + path)
}

// 按照优先级搜索文件，返回搜索到的文件绝对路径
func (sp *SPath) Search(name string) string {
    path := sp.cache.Get(name)
    if path == "" {
        sp.mu.RLock()
        for _, v := range sp.paths {
            path = v + gfile.Separator + name
            if gfile.Exists(path) {
                break
            }
        }
        sp.mu.RUnlock()
        if path != "" {
            sp.cache.Set(name, path)
        }
    }
    return path
}