// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// Package gspath implements file index and search for folders.
// 
// 搜索目录管理, 
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
    paths *garray.StringArray    // 搜索路径，按照优先级进行排序
    cache *gmap.StringStringMap  // 搜索结果缓存map
}

// 文件搜索缓存项
type SPathCacheItem struct {
    path  string                 // 文件/目录绝对路径
    isDir bool                   // 是否目录
}

var (
    // 单个目录路径对应的SPath对象指针，用于路径检索对象复用
    pathsMap = gmap.NewStringInterfaceMap()
)

// 创建一个搜索对象
func New(path...string) *SPath {
    sp := &SPath {
        paths : garray.NewStringArray(0, 1),
        cache : gmap.NewStringStringMap(),
    }
    if len(path) > 0 {
        sp.Add(path[0])
    }
    return sp
}

// 创建/获取一个单例的搜索对象, root必须为目录的绝对路径
func Get(root string) *SPath {
    return pathsMap.GetOrSetFuncLock(root, func() interface{} {
        return New(root)
    }).(*SPath)
}

// 检索root目录(必须为绝对路径)下面的name文件的绝对路径，indexFiles用于指定当检索到的结果为目录时，同时检索是否存在这些indexFiles文件
func Search(root string, name string, indexFiles...string) (filePath string, isDir bool) {
    return Get(root).Search(name, indexFiles...)
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
func (sp *SPath) Search(name string, indexFiles...string) (filePath string, isDir bool) {
    name = sp.formatCacheName(name)
    if v := sp.cache.Get(name); v != "" {
        filePath, isDir = sp.parseCacheValue(v)
        if len(indexFiles) > 0 && isDir {
            if name == "/" {
                name = ""
            }
            for _, file := range indexFiles {
                if v := sp.cache.Get(name + "/" + file); v != "" {
                    return sp.parseCacheValue(v)
                }
            }
        }
    }
    return
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

// 返回当前对象搜索目录路径列表
func (sp *SPath) Paths() []string {
    return sp.paths.Slice()
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
    if runtime.GOOS != "linux" {
        name = gstr.Replace(name, "\\", "/")
    }
    return "/" + strings.Trim(name, "./")
}

// 根据path计算出对应的缓存name, dirPath为检索根目录路径
func (sp *SPath) nameFromPath(filePath, rootPath string) string {
    name  := gstr.Replace(filePath, rootPath, "")
    name   = sp.formatCacheName(name)
    return name
}

// 按照一定数据结构生成缓存的数据项字符串
func (sp *SPath) makeCacheValue(filePath string, isDir bool) string {
    if isDir {
        return filePath + "_D_"
    }
    return filePath + "_F_"
}

// 按照一定数据结构解析数据项字符串
func (sp *SPath) parseCacheValue(value string) (filePath string, isDir bool) {
    if value[len(value) - 2 : len(value) - 1][0] == 'F' {
        return value[: len(value) - 3], false
    }
    return value[: len(value) - 3], true
}

// 添加path到缓存中(递归)
func (sp *SPath) addToCache(filePath, rootPath string) {
    // 首先添加自身
    idDir := gfile.IsDir(filePath)
    sp.cache.SetIfNotExist(sp.nameFromPath(filePath, rootPath), sp.makeCacheValue(filePath, idDir))
    // 如果添加的是目录，那么需要递归添加
    if idDir {
        if files, err := gfile.ScanDir(filePath, "*", true); err == nil {
            for _, path := range files {
                sp.cache.SetIfNotExist(sp.nameFromPath(path, rootPath), sp.makeCacheValue(path, gfile.IsDir(path)))
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