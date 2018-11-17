// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// 搜索目录管理.
// 可以添加搜索目录，按照添加的优先级进行文件检索，并在内部进行高效缓存处理。
package gspath

import (
    "errors"
    "fmt"
    "gitee.com/johng/gf/g/container/garray"
    "gitee.com/johng/gf/g/container/gmap"
    "gitee.com/johng/gf/g/os/gfile"
    "gitee.com/johng/gf/g/os/gfsnotify"
    "gitee.com/johng/gf/g/util/gstr"
    "runtime"
    "sort"
    "strings"
)

// 文件目录搜索管理对象
type SPath struct {
    paths *garray.StringArray       // 搜索路径，按照优先级进行排序
    cache *gmap.StringInterfaceMap  // 搜索结果缓存map
}

// 文件搜索缓存项
type SPathCacheItem struct {
    path  string                    // 文件/目录绝对路径
    isDir bool                      // 是否目录
}

// 创建一个搜索对象
func New () *SPath {
    return &SPath {
        paths : garray.NewStringArray(0, 2),
        cache : gmap.NewStringInterfaceMap(),
    }
}

// 设置搜索路径，只保留当前设置项，其他搜索路径被清空
func (sp *SPath) Set(path string) (realPath string, err error) {
    realPath = gfile.RealPath(path)
    if realPath == "" {
        realPath, _ = sp.Search(path)
        if realPath == "" {
            realPath = gfile.RealPath(gfile.Pwd() + gfile.Separator + path)
        }
    }
    if realPath == "" {
        return realPath, errors.New(fmt.Sprintf(`path "%s" does not exist`, path))
    }
    if realPath == "" {
        return realPath, errors.New("invalid path:" + path)
    }
    // 设置的搜索路径必须为目录
    if gfile.IsDir(realPath) {
        realPath = strings.TrimRight(realPath, gfile.Separator)
        if sp.paths.Search(realPath) != -1 {
            for _, v := range sp.paths.Slice() {
                sp.removeMonitorByPath(v)
            }
        }
        sp.paths.Clear()
        sp.cache.Clear()
        sp.paths.Append(realPath)
        sp.updateCacheByPath(realPath)
        sp.addMonitorByPath(realPath)
        return realPath, nil
    } else {
        return "", errors.New(path + " should be a folder")
    }
}

// 添加搜索路径
func (sp *SPath) Add(path string) (realPath string, err error) {
    realPath = gfile.RealPath(path)
    if realPath == "" {
        realPath, _ = sp.Search(path)
        if realPath == "" {
            realPath = gfile.RealPath(gfile.Pwd() + gfile.Separator + path)
        }
    }
    if realPath == "" {
        return realPath, errors.New(fmt.Sprintf(`path "%s" does not exist`, path))
    }
    if realPath == "" {
        return realPath, errors.New("invalid path:" + path)
    }
    // 添加的搜索路径必须为目录
    if gfile.IsDir(realPath) {
        // 如果已经添加则不再添加
        if sp.paths.Search(realPath) < 0 {
            realPath = strings.TrimRight(realPath, gfile.Separator)
            sp.paths.Append(realPath)
            sp.updateCacheByPath(realPath)
            sp.addMonitorByPath(realPath)
        }
        return realPath, nil
    } else {
        return "", errors.New(path + " should be a folder")
    }
}

// 给定的name只是相对文件路径，找不到该文件时，返回空字符串;
// 当给定indexFiles时，如果name时一个目录，那么会进一步检索其下对应的indexFiles文件是否存在，存在则返回indexFile绝对路径；
// 否则返回name目录绝对路径。
func (sp *SPath) Search(name string, indexFiles...string) (path string, isDir bool) {
    name = sp.formatCacheName(name)
    if v := sp.cache.Get(name); v != nil {
        item := v.(*SPathCacheItem)
        if len(indexFiles) > 0 && item.isDir {
            if name == "/" {
                name = ""
            }
            for _, file := range indexFiles {
                if v := sp.cache.Get(name + "/" + file); v != nil {
                    item := v.(*SPathCacheItem)
                    return item.path, item.isDir
                }
            }
        }
        return item.path, item.isDir
    }
    return "", false
}

// 从搜索路径中移除指定的文件，这样该文件无法给搜索。
// path可以是绝对路径，也可以相对路径。
func (sp *SPath) Remove(path string) {
    if gfile.Exists(path) {
        for _, v := range sp.paths.Slice() {
            name := gstr.Replace(path, v, "")
            name  = sp.formatCacheName(name)
            sp.cache.Remove(name)
        }
    } else {
        name := sp.formatCacheName(path)
        sp.cache.Remove(name)
    }
}

// 返回当前对象缓存的所有路径列表
func (sp *SPath) AllPaths() []string {
    paths := sp.cache.Keys()
    if len(paths) > 0 {
        sort.Strings(paths)
    }
    return paths
}

// 当前的搜索路径数量
func (sp *SPath) Size() int {
    return sp.paths.Len()
}

// 递归添加目录下的文件
func (sp *SPath) updateCacheByPath(path string) {
    sp.addToCache(path, path)
}

// 格式化name返回符合规范的缓存名称，分隔符号统一为'/'，且前缀必须以'/'开头(类似HTTP URI).
func (sp *SPath) formatCacheName(name string) string {
    name = strings.Trim(name, "./")
    if runtime.GOOS != "linux" {
        name = gstr.Replace(name, "\\", "/")
    }
    return "/" + name
}

// 根据path计算出对应的缓存name
func (sp *SPath) nameFromPath(filePath, dirPath string) string {
    name  := gstr.Replace(filePath, dirPath, "")
    name   = sp.formatCacheName(name)
    return name
}

// 添加path到缓存中(递归)
func (sp *SPath) addToCache(filePath, dirPath string) {
    // 首先添加自身
    idDir := gfile.IsDir(filePath)
    sp.cache.SetIfNotExist(sp.nameFromPath(filePath, dirPath), func() interface{} {
        return &SPathCacheItem {
            path  : filePath,
            isDir : idDir,
        }
    })
    // 如果添加的是目录，那么需要递归
    if idDir {
        if files, err := gfile.ScanDir(filePath, "*", true); err == nil {
            for _, path := range files {
                sp.addToCache(path, dirPath)
            }
        }
    }
}

// 添加文件目录监控(递归)，当目录下的文件有更新时，会同时更新缓存。
// 这里需要注意的点是，由于添加监听是递归添加的，那么假如删除一个目录，那么该目录下的文件(包括目录)也会产生一条删除事件，总共会产生N条事件。
func (sp *SPath) addMonitorByPath(path string) {
    gfsnotify.Add(path, func(event *gfsnotify.Event) {
        //glog.Debug(event.String())
        switch {
            case event.IsRemove():
                sp.cache.Remove(sp.nameFromPath(event.Path, path))

            case event.IsRename():
                if !gfile.Exists(event.Path) {
                    sp.cache.Remove(sp.nameFromPath(event.Path, path))
                }

            case event.IsCreate():
                sp.addToCache(event.Path, path)
        }
    }, true)
}

// 删除监听(递归)
func (sp *SPath) removeMonitorByPath(path string) {
    gfsnotify.Remove(path)
}