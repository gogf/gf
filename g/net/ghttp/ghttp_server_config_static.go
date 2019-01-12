// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// 静态文件搜索优先级: ServerPaths > ServerRoot > SearchPath

package ghttp

import (
    "fmt"
    "gitee.com/johng/gf/g/container/garray"
    "gitee.com/johng/gf/g/os/gfile"
    "gitee.com/johng/gf/g/os/glog"
    "gitee.com/johng/gf/g/util/gconv"
    "strings"
)

// 静态文件目录映射关系对象
type staticPathItem struct {
    prefix string // 映射的URI前缀
    path   string // 静态文件目录绝对路径
}

// 设置http server参数 - IndexFiles，默认展示文件，如：index.html, index.htm
func (s *Server)SetIndexFiles(index []string) {
    if s.Status() == SERVER_STATUS_RUNNING {
        glog.Error(gCHANGE_CONFIG_WHILE_RUNNING_ERROR)
        return
    }
    s.config.IndexFiles = index
}

// 允许展示访问目录的文件列表
func (s *Server)SetIndexFolder(enabled bool) {
    if s.Status() == SERVER_STATUS_RUNNING {
        glog.Error(gCHANGE_CONFIG_WHILE_RUNNING_ERROR)
        return
    }
    s.config.IndexFolder = enabled
}

// 是否开启/关闭静态文件服务，当关闭时仅提供动态接口服务，路由性能会得到一定提升
func (s *Server) SetFileServerEnabled(enabled bool) {
    if s.Status() == SERVER_STATUS_RUNNING {
        glog.Error(gCHANGE_CONFIG_WHILE_RUNNING_ERROR)
        return
    }
    s.config.FileServerEnabled = enabled
}

// 设置http server参数 - ServerRoot
func (s *Server)SetServerRoot(root string) {
    if s.Status() == SERVER_STATUS_RUNNING {
        glog.Error(gCHANGE_CONFIG_WHILE_RUNNING_ERROR)
        return
    }
    // RealPath的作用除了校验地址正确性以外，还转换分隔符号为当前系统正确的文件分隔符号
    path := gfile.RealPath(root)
    if path == "" {
        path = gfile.RealPath(gfile.MainPkgPath() + gfile.Separator + root)
    }
    if path == "" {
        glog.Fatal(fmt.Sprintf(`[ghttp] SetServerRoot failed: path "%s" does not exist`, root))
    }
    s.config.SearchPaths       = []string{strings.TrimRight(path, gfile.Separator)}
    s.config.FileServerEnabled = true
}

// 添加静态文件搜索**目录**，必须给定目录的绝对路径
func (s *Server) AddSearchPath(path string) {
    if s.Status() == SERVER_STATUS_RUNNING {
        glog.Error(gCHANGE_CONFIG_WHILE_RUNNING_ERROR)
        return
    }
    // RealPath的作用除了校验地址正确性以外，还转换分隔符号为当前系统正确的文件分隔符号
    realPath := gfile.RealPath(path)
    if realPath == "" {
        realPath = gfile.RealPath(gfile.MainPkgPath() + gfile.Separator + path)
    }
    if realPath == "" {
        glog.Fatal(fmt.Sprintf(`[ghttp] AddSearchPath failed: path "%s" does not exist`, path))
    }
    s.config.SearchPaths       = append(s.config.SearchPaths, realPath)
    s.config.FileServerEnabled = true
}

// 添加URI与静态**目录**的映射
func (s *Server) AddStaticPath(prefix string, path string) {
    if s.Status() == SERVER_STATUS_RUNNING {
        glog.Error(gCHANGE_CONFIG_WHILE_RUNNING_ERROR)
        return
    }
    // RealPath的作用除了校验地址正确性以外，还转换分隔符号为当前系统正确的文件分隔符号
    realPath := gfile.RealPath(path)
    if realPath == "" {
        realPath = gfile.RealPath(gfile.MainPkgPath() + gfile.Separator + path)
    }
    if realPath == "" {
        glog.Fatal(fmt.Sprintf(`[ghttp] AddStaticPath failed: path "%s" does not exist`, path))
    }
    addItem := staticPathItem {
        prefix : prefix,
        path   : realPath,
    }
    if len(s.config.StaticPaths) > 0 {
        // 先添加item
        s.config.StaticPaths = append(s.config.StaticPaths, addItem)
        // 按照prefix从长到短进行排序
        array := garray.NewSortedArray(0, func(v1, v2 interface{}) int {
            s1 := gconv.String(v1)
            s2 := gconv.String(v2)
            r  := len(s2) - len(s1)
            if r == 0 {
                r = strings.Compare(s1, s2)
            }
            return r
        }, false)
        for _, v := range s.config.StaticPaths {
            array.Add(v.prefix)
        }
        // 按照重新排序的顺序重新添加item
        paths := make([]staticPathItem, 0)
        for _, v := range array.Slice() {
            for _, item := range s.config.StaticPaths {
                if strings.EqualFold(gconv.String(v), item.prefix) {
                    paths = append(paths, item)
                    break
                }
            }
        }
        s.config.StaticPaths = paths
    } else {
        s.config.StaticPaths = []staticPathItem { addItem }
    }
    s.config.FileServerEnabled = true
}

