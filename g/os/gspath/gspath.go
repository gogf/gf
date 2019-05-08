// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gspath implements file index and search for folders.
// 
// 搜索目录管理, 
// 可以添加搜索目录，按照添加的优先级进行文件检索，并在内部进行高效缓存处理(可选)。
// 注意：当开启缓存功能后，在新增/删除文件时，会存在检索延迟。
package gspath

import (
    "errors"
    "fmt"
    "github.com/gogf/gf/g/container/garray"
    "github.com/gogf/gf/g/container/gmap"
    "github.com/gogf/gf/g/os/gfile"
    "github.com/gogf/gf/g/text/gstr"
    "os"
    "sort"
    "strings"
)

// 文件目录搜索管理对象
type SPath struct {
    paths *garray.StringArray  // 搜索路径，按照优先级进行排序
    cache *gmap.StrStrMap      // 搜索结果缓存map(如果未nil表示未启用缓存功能)
}

// 文件搜索缓存项
type SPathCacheItem struct {
    path  string                 // 文件/目录绝对路径
    isDir bool                   // 是否目录
}

var (
    // 单个目录路径对应的SPath对象指针，用于路径检索对象复用
    pathsMap      = gmap.NewStrAnyMap()
    pathsCacheMap = gmap.NewStrAnyMap()
)

// 创建一个搜索对象
func New(path string, cache bool) *SPath {
    sp := &SPath {
        paths : garray.NewStringArray(),
    }
    if cache {
        sp.cache = gmap.NewStrStrMap()
    }
    if len(path) > 0 {
        if _, err := sp.Add(path); err != nil {
            //fmt.Errorf(err.Error())
        }
    }
    return sp
}

// 创建/获取一个单例的搜索对象, root必须为目录的绝对路径
func Get(root string, cache bool) *SPath {
    return pathsMap.GetOrSetFuncLock(root, func() interface{} {
        return New(root, cache)
    }).(*SPath)
}

// 检索root目录(必须为绝对路径)下面的name文件的绝对路径，indexFiles用于指定当检索到的结果为目录时，同时检索是否存在这些indexFiles文件
func Search(root string, name string, indexFiles...string) (filePath string, isDir bool) {
    return Get(root, false).Search(name, indexFiles...)
}

// 检索root目录(必须为绝对路径)下面的name文件的绝对路径，indexFiles用于指定当检索到的结果为目录时，同时检索是否存在这些indexFiles文件
func SearchWithCache(root string, name string, indexFiles...string) (filePath string, isDir bool) {
    return Get(root, true).Search(name, indexFiles...)
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
        if sp.cache != nil {
            sp.cache.Clear()
        }
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
        //fmt.Println("gspath:", realPath, sp.paths.Search(realPath))
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
// 当给定indexFiles时，如果name是一个目录，那么会进一步检索其下对应的indexFiles文件是否存在，存在则返回indexFile绝对路径；
// 否则返回name目录绝对路径。
func (sp *SPath) Search(name string, indexFiles...string) (filePath string, isDir bool) {
    // 不使用缓存
    if sp.cache == nil {
        sp.paths.LockFunc(func(array []string) {
            path := ""
            for _, v := range array {
                path = v + gfile.Separator + name
                if stat, err := os.Stat(path); !os.IsNotExist(err) {
                    filePath = path
                    isDir    = stat.IsDir()
                    break
                }
            }
        })
        if len(indexFiles) > 0 && isDir {
            if name == "/" {
                name = ""
            }
            path := ""
            for _, file := range indexFiles {
                path = filePath + gfile.Separator + file
                if gfile.Exists(path) {
                    filePath = path
                    isDir    = false
                    break
                }
            }
        }
        return
    }
    // 使用缓存功能
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
    if sp.cache == nil {
        return
    }
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
    if sp.cache == nil {
        return nil
    }
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
